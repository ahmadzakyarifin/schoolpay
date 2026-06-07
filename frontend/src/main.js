import { createApp } from 'vue'
import { createPinia } from 'pinia'
import axios from 'axios'
import './assets/style.css'
import App from './App.vue'
import router from './router'
import FormError from './components/ui/FormError.vue'
axios.defaults.baseURL = import.meta.env.VITE_API_URL || "http://localhost:8080/api/"
axios.defaults.withCredentials = true
import { useAuthStore } from './store/auth'
import { queueOfflineRequest, syncOfflineOutbox } from './utils/offlineSync'

const matchesApiPath = (pathOnly, paths) => {
    return paths.some(path => pathOnly === path || pathOnly.startsWith(`${path}/`))
}

const getApiPath = (url = '') => {
    return String(url || '').replace(axios.defaults.baseURL, '').replace(/^\/+/, '').split('?')[0]
}

const getDeviceID = () => {
    const key = 'schoolpay_device_id'
    try {
        let deviceID = localStorage.getItem(key)
        if (!deviceID) {
            deviceID = crypto.randomUUID ? crypto.randomUUID() : `${Date.now()}-${Math.random().toString(36).slice(2)}`
            localStorage.setItem(key, deviceID)
        }
        return deviceID
    } catch (e) {
        return ''
    }
}

const getOfflineWritePolicy = (config) => {
    const method = (config.method || 'get').toLowerCase()
    if (method === 'get') return { type: 'read' }

    const pathOnly = getApiPath(config.url)

    const transactionPaths = [
        'finance/payments',
        'finance/payment-intent',
        'finance/bills',
        'finance/generate-bills'
    ]

    if (matchesApiPath(pathOnly, transactionPaths)) {
        return {
            type: 'blocked',
            message: 'Aksi transaksi dikunci saat offline untuk mencegah data pembayaran atau tagihan ganda.'
        }
    }

    const studentOnlineOnlyPaths = [
        'students/bulk-promote',
        'students/bulk-graduate'
    ]

    if (matchesApiPath(pathOnly, studentOnlineOnlyPaths)) {
        return {
            type: 'blocked',
            message: 'Kenaikan kelas dan kelulusan memerlukan koneksi online karena sistem harus mengecek tagihan aktif.'
        }
    }

    const onlineOnlyWritePaths = [
        'whatsapp',
        'support/conversations',
        'parent/support'
    ]

    if (matchesApiPath(pathOnly, onlineOnlyWritePaths)) {
        return {
            type: 'blocked',
            message: 'Aksi ini memerlukan koneksi online.'
        }
    }

    const offlineQueuePaths = [
        'users',
        'students',
        'academic',
        'finance',
        'dashboard',
        'parent'
    ]

    if (matchesApiPath(pathOnly, offlineQueuePaths)) {
        return {
            type: 'queued',
            message: 'Perubahan data disimpan sementara dan akan disinkronkan saat server online.'
        }
    }

    return { type: 'passthrough' }
}

const getOfflineCacheKey = (config) => {
    if (config.skipOfflineCache) return ''
    const method = (config.method || 'get').toLowerCase()
    if (method !== 'get' || config.responseType === 'blob') return ''

    const pathOnly = getApiPath(config.url)
    const cacheablePaths = ['users', 'students', 'academic', 'dashboard', 'finance', 'parent', 'audit-logs']
    if (!matchesApiPath(pathOnly, cacheablePaths)) return ''

    let paramsObj = config.params ? { ...config.params } : {}
    if (pathOnly.startsWith('dashboard')) {
        delete paramsObj.ref_date
    }
    const params = Object.keys(paramsObj).length ? JSON.stringify(paramsObj) : ''
    return `offline_api_cache:${pathOnly}:${params}`
}

const getCachedOfflineResponse = (config) => {
    const cacheKey = getOfflineCacheKey(config)
    if (!cacheKey) return null

    try {
        const cached = JSON.parse(localStorage.getItem(cacheKey) || 'null')
        if (!cached) return null
        return {
            data: cached.data,
            status: 200,
            statusText: 'OK (offline cache)',
            headers: {},
            config,
            request: null
        }
    } catch (e) {
        localStorage.removeItem(cacheKey)
        return null
    }
}

axios.interceptors.request.use(
    async config => {
        const authStore = useAuthStore()

        if ((authStore.isOffline || navigator.onLine === false) && (config.method || 'get').toLowerCase() === 'get' && !config.skipOfflineCache) {
            const cachedResponse = getCachedOfflineResponse(config)
            if (cachedResponse) {
                config.adapter = () => Promise.resolve(cachedResponse)
                config.offlineCached = true
            }
        }

        const offlinePolicy = getOfflineWritePolicy(config)
        if ((authStore.isOffline || navigator.onLine === false) && offlinePolicy.type === 'blocked') {
            authStore.isOffline = true
            const offlineError = new Error(offlinePolicy.message)
            offlineError.isOfflineBlocked = true
            window.dispatchEvent(new CustomEvent('network-error', { detail: offlineError.message }))
            return Promise.reject(offlineError)
        }

        if ((authStore.isOffline || navigator.onLine === false) && offlinePolicy.type === 'queued' && !config.skipOfflineQueue) {
            authStore.isOffline = true
            try {
                const item = await queueOfflineRequest(config)
                config.offlineQueued = true
                config.adapter = () => Promise.resolve({
                    data: {
                        status: 'queued',
                        message: offlinePolicy.message,
                        data: { offline_queue_id: item.id }
                    },
                    status: 202,
                    statusText: 'Queued Offline',
                    headers: {},
                    config,
                    request: null
                })
            } catch (err) {
                const offlineError = new Error(err.message || 'Gagal menyimpan perubahan offline')
                offlineError.isOfflineBlocked = true
                window.dispatchEvent(new CustomEvent('network-error', { detail: offlineError.message }))
                return Promise.reject(offlineError)
            }
        }

        config.headers = config.headers || {}
        const deviceID = getDeviceID()
        if (deviceID) {
            config.headers['X-Device-ID'] = deviceID
        }
        config.headers['X-App-Platform'] = 'Web'
        config.headers['X-App-Version'] = '1.0.0'

        const token = authStore.token
        if (token) {
            config.headers['Authorization'] = `Bearer ${token}`
        }
        return config
    },
    error => Promise.reject(error)
)

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(router)
app.component('FormError', FormError)


const authStore = useAuthStore(pinia)
axios.interceptors.response.use(
    response => {
        const authStore = useAuthStore()

        const cacheKey = getOfflineCacheKey(response.config || {})
        if (cacheKey && !response.config?.offlineCached) {
            try {
                localStorage.setItem(cacheKey, JSON.stringify({
                    data: response.data,
                    cached_at: new Date().toISOString()
                }))
            } catch (e) { }
        }

        if (!response.config?.offlineQueued && !response.config?.offlineCached) {
            authStore.isOffline = false
            if (!response.config?.skipOfflineQueue) {
                syncOfflineOutbox().catch(() => { })
            }
        }
        return response
    },
    async error => {
        const originalRequest = error.config

        if (error.code === "ERR_NETWORK" || !error.response) {
            const url = error.config?.url || "";
            const isLocalAPI = !url || !url.startsWith("http") || url.startsWith(axios.defaults.baseURL);
            if (isLocalAPI) {
                authStore.isOffline = true;
                const customError = new Error("Server sedang offline. Mode offline terbatas aktif.");
                customError.isNetworkError = true;
                const event = new CustomEvent("network-error", { detail: customError.message });
                window.dispatchEvent(event);

                const cachedResponse = getCachedOfflineResponse(originalRequest || {})
                if (cachedResponse) return cachedResponse

                const offlinePolicy = getOfflineWritePolicy(originalRequest || {})
                if (offlinePolicy.type === 'queued' && originalRequest && !originalRequest.skipOfflineQueue) {
                    try {
                        const item = await queueOfflineRequest(originalRequest)
                        originalRequest.offlineQueued = true
                        return {
                            data: {
                                status: 'queued',
                                message: offlinePolicy.message,
                                data: { offline_queue_id: item.id }
                            },
                            status: 202,
                            statusText: 'Queued Offline',
                            headers: {},
                            config: originalRequest,
                            request: null
                        }
                    } catch (queueError) {
                        return Promise.reject(queueError)
                    }
                }
            }
            return Promise.reject(error);
        }


        if (error.response?.status === 503 && error.response?.data?.status === 'offline') {
            authStore.isOffline = true
            const customError = new Error('Server sedang offline. Mode offline terbatas aktif.')
            customError.isNetworkError = true
            window.dispatchEvent(new CustomEvent('network-error', { detail: customError.message }))

            const cachedResponse = getCachedOfflineResponse(originalRequest || {})
            if (cachedResponse) return cachedResponse

            return Promise.reject(error)
        }

        if (error.response?.status === 429) {
            const retryAfter = Number(error.response?.data?.data?.retry_after_seconds || error.response?.headers?.['retry-after'] || 0)
            const message = retryAfter > 0
                ? `Terlalu banyak request. Coba lagi dalam ${retryAfter} detik.`
                : (error.response?.data?.message || 'Terlalu banyak request. Coba lagi sebentar lagi.')

            window.dispatchEvent(new CustomEvent('rate-limit-error', {
                detail: {
                    message,
                    retryAfterSeconds: retryAfter
                }
            }))
        }

        if (error.response && error.response.status === 401 && !originalRequest._retry &&
            !originalRequest.url.includes('/auth/refresh') && !originalRequest.url.includes('/auth/login')) {
            originalRequest._retry = true
            const result = await authStore.refreshToken()

            if (result === true) {
                originalRequest.headers['Authorization'] = `Bearer ${authStore.token}`
                return axios(originalRequest)
            } else if (result === "network_error") {
                return Promise.reject(error)
            } else {
                const message = error.response?.data?.message || 'expired'
                authStore.logout()
                router.push({ name: 'login', query: { reason: message } })
            }
        }

        if (error.response && error.response.status === 403) {
            // Role changed or access revoked
            const message = error.response.data?.message || 'forbidden'
            authStore.logout()
            router.push({ name: 'login', query: { reason: message } })
        }
        return Promise.reject(error)
    }
)

const checkServerOnline = async () => {
    if (typeof navigator !== 'undefined' && navigator.onLine === false) return false
    try {
        await axios.get('health', {
            skipOfflineCache: true,
            skipOfflineQueue: true,
            timeout: 3000
        })
        const store = useAuthStore()
        store.isOffline = false
        syncOfflineOutbox().catch(() => { })
        return true
    } catch (err) {
        return false
    }
}

let offlineProbeTimer = null
const startOfflineProbe = () => {
    if (offlineProbeTimer) return
    offlineProbeTimer = window.setInterval(async () => {
        const store = useAuthStore()
        if (!store.isOffline) {
            window.clearInterval(offlineProbeTimer)
            offlineProbeTimer = null
            return
        }
        await checkServerOnline()
    }, 5000)
}

// Initialize auth and then mount
authStore.initializeAuth().finally(() => {
    app.mount('#app')
    if (authStore.isOffline) {
        startOfflineProbe()
        checkServerOnline()
    }
})


window.addEventListener('online', () => {
    checkServerOnline().then((ok) => {
        if (!ok) startOfflineProbe()
    })
})

window.addEventListener('offline', () => {
    const authStore = useAuthStore()
    authStore.isOffline = true
    startOfflineProbe()
})

window.addEventListener('focus', () => {
    const authStore = useAuthStore()
    if (authStore.isOffline) {
        checkServerOnline().then((ok) => {
            if (!ok) startOfflineProbe()
        })
    }
})

window.addEventListener('network-error', () => {
    startOfflineProbe()
})

if ('serviceWorker' in navigator) {
    window.addEventListener('load', () => {
        navigator.serviceWorker.register('/sw.js').catch((err) => {
            console.warn('Service worker registration failed', err)
        })
    })
}
