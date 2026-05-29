const CACHE_NAME = 'schoolpay-shell-v1'
const SHELL_ASSETS = ['/', '/manifest.webmanifest', '/pwa-icon.svg']

self.addEventListener('install', (event) => {
  event.waitUntil(caches.open(CACHE_NAME).then((cache) => cache.addAll(SHELL_ASSETS)))
  self.skipWaiting()
})

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((keys) => Promise.all(keys.filter((key) => key !== CACHE_NAME).map((key) => caches.delete(key))))
  )
  self.clients.claim()
})

self.addEventListener('fetch', (event) => {
  const req = event.request
  const url = new URL(req.url)

  if (req.method !== 'GET') return

  if (url.pathname.startsWith('/api/')) {
    event.respondWith(
      fetch(req).catch(() => new Response(JSON.stringify({ status: 'offline', message: 'Server sedang offline. Mode offline terbatas aktif.' }), {
        status: 503,
        headers: { 'Content-Type': 'application/json' }
      }))
    )
    return
  }

  event.respondWith(
    fetch(req).then((res) => {
      const copy = res.clone()
      caches.open(CACHE_NAME).then((cache) => cache.put(req, copy))
      return res
    }).catch(() => caches.match(req).then((cached) => cached || caches.match('/')))
  )
})
