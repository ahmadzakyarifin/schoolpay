import axios from 'axios'

const DB_NAME = 'schoolpay_offline'
const DB_VERSION = 1
const OUTBOX_STORE = 'outbox'

let syncPromise = null

const openDB = () => new Promise((resolve, reject) => {
  const request = indexedDB.open(DB_NAME, DB_VERSION)

  request.onupgradeneeded = () => {
    const db = request.result
    if (!db.objectStoreNames.contains(OUTBOX_STORE)) {
      const store = db.createObjectStore(OUTBOX_STORE, { keyPath: 'id' })
      store.createIndex('created_at', 'created_at')
    }
  }

  request.onsuccess = () => resolve(request.result)
  request.onerror = () => reject(request.error)
})

const withStore = async (mode, callback) => {
  const db = await openDB()
  return new Promise((resolve, reject) => {
    const tx = db.transaction(OUTBOX_STORE, mode)
    const store = tx.objectStore(OUTBOX_STORE)
    let result

    tx.oncomplete = () => {
      db.close()
      resolve(result)
    }
    tx.onerror = () => {
      db.close()
      reject(tx.error)
    }

    result = callback(store)
  })
}

const requestToPromise = (request) => new Promise((resolve, reject) => {
  request.onsuccess = () => resolve(request.result)
  request.onerror = () => reject(request.error)
})

const hasFileValue = (value) => {
  return typeof File !== 'undefined' && value instanceof File
}

const serializeData = (data) => {
  if (typeof FormData !== 'undefined' && data instanceof FormData) {
    const entries = []
    for (const [key, value] of data.entries()) {
      if (hasFileValue(value)) {
        throw new Error('Upload file belum didukung saat offline. Simpan perubahan ini ketika online kembali.')
      }
      entries.push([key, value])
    }
    return { kind: 'form', entries }
  }

  return { kind: 'json', value: data ?? null }
}

const deserializeData = (payload) => {
  if (payload?.kind === 'form') {
    const formData = new FormData()
    for (const [key, value] of payload.entries || []) {
      formData.append(key, value)
    }
    return formData
  }

  return payload?.value ?? null
}

const makeID = () => {
  if (typeof crypto !== 'undefined' && crypto.randomUUID) return crypto.randomUUID()
  return `${Date.now()}-${Math.random().toString(16).slice(2)}`
}

const sanitizeHeaders = (headers = {}) => {
  const clean = {}
  Object.entries(headers || {}).forEach(([key, value]) => {
    const lower = key.toLowerCase()
    if (lower === 'authorization' || lower === 'content-type') return
    clean[key] = value
  })
  return clean
}

export const isMasterWritePath = (pathOnly) => {
  return ['users', 'students', 'academic'].some(path => pathOnly === path || pathOnly.startsWith(`${path}/`))
}

export const queueOfflineRequest = async (config) => {
  const item = {
    id: makeID(),
    idempotency_key: makeID(),
    method: (config.method || 'get').toLowerCase(),
    url: config.url,
    baseURL: config.baseURL,
    params: config.params || null,
    data: serializeData(config.data),
    headers: sanitizeHeaders(config.headers),
    created_at: new Date().toISOString(),
    status: 'pending'
  }

  await withStore('readwrite', (store) => store.put(item))
  const count = await getOfflineOutboxCount()
  window.dispatchEvent(new CustomEvent('offline-sync-queued', { detail: { item, count } }))
  return item
}

export const getOfflineOutboxItems = async () => {
  return withStore('readonly', async (store) => {
    const items = await requestToPromise(store.getAll())
    return items.sort((a, b) => new Date(a.created_at) - new Date(b.created_at))
  })
}

export const getOfflineOutboxCount = async () => {
  return withStore('readonly', (store) => requestToPromise(store.count()))
}

const deleteOutboxItem = async (id) => {
  return withStore('readwrite', (store) => store.delete(id))
}

export const syncOfflineOutbox = async () => {
  if (syncPromise) return syncPromise
  if (typeof navigator !== 'undefined' && navigator.onLine === false) return { synced: 0, pending: await getOfflineOutboxCount() }

  syncPromise = (async () => {
    const items = await getOfflineOutboxItems()
    let synced = 0

    for (const item of items) {
      try {
        await axios.request({
          method: item.method,
          url: item.url,
          baseURL: item.baseURL,
          params: item.params || undefined,
          data: deserializeData(item.data),
          headers: {
            'X-Idempotency-Key': item.idempotency_key
          },
          skipOfflineQueue: true,
          skipOfflineCache: true
        })
        await deleteOutboxItem(item.id)
        synced += 1
      } catch (err) {
        if (!err.response || err.code === 'ERR_NETWORK' || err.response?.status === 503) {
          break
        }

        item.status = 'failed'
        item.error = err.response?.data?.message || 'Sinkron gagal'
        item.failed_at = new Date().toISOString()
        await withStore('readwrite', (store) => store.put(item))
      }
    }

    const pending = await getOfflineOutboxCount()
    window.dispatchEvent(new CustomEvent('offline-sync-complete', { detail: { synced, pending } }))
    return { synced, pending }
  })().finally(() => {
    syncPromise = null
  })

  return syncPromise
}
