import { createApp } from 'vue'
import { createPinia } from 'pinia'
import axios from 'axios'
import './assets/style.css'
import App from './App.vue'
import router from './router'
import FormError from './components/ui/FormError.vue'
import { useAuthStore } from './store/auth'
// import './mock/index.js' 
axios.defaults.baseURL = import.meta.env.VITE_API_URL || "http://localhost:8080/api/"
axios.defaults.withCredentials = true

const legacyOfflineDatabases = [
    'schoolpay_offline',
    'schoolpay_offline_cache',
    'schoolpay_auth_cache'
]

const cleanupLegacyOfflineArtifacts = () => {
    if ('serviceWorker' in navigator) {
        navigator.serviceWorker.getRegistrations()
            .then(registrations => registrations.forEach(registration => registration.unregister()))
            .catch(() => { })
    }

    if (typeof caches !== 'undefined') {
        caches.keys()
            .then(keys => Promise.all(keys.filter(key => key.startsWith('schoolpay')).map(key => caches.delete(key))))
            .catch(() => { })
    }

    if (typeof indexedDB !== 'undefined') {
        legacyOfflineDatabases.forEach(name => {
            try {
                indexedDB.deleteDatabase(name)
            } catch (e) { }
        })
    }

    try {
        Object.keys(localStorage)
            .filter(key => key === 'user' || key.startsWith('offline_api_cache:'))
            .forEach(key => localStorage.removeItem(key))
    } catch (e) { }
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

axios.interceptors.request.use(
    config => {
        const authStore = useAuthStore()
        config.headers = config.headers || {}

        const deviceID = getDeviceID()
        if (deviceID) {
            config.headers['X-Device-ID'] = deviceID
        }
        config.headers['X-App-Platform'] = 'Web'
        config.headers['X-App-Version'] = '1.0.0'

        if (authStore.token) {
            config.headers['Authorization'] = `Bearer ${authStore.token}`
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
    response => response,
    async error => {
        const originalRequest = error.config || {}

        if (error.code === "ERR_NETWORK" || !error.response) {
            window.dispatchEvent(new CustomEvent('network-error', {
                detail: 'Server tidak dapat dijangkau. Pastikan koneksi dan layanan backend aktif.'
            }))
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

        if (error.response?.status === 401 && !originalRequest._retry &&
            !originalRequest.url?.includes('/auth/refresh') && !originalRequest.url?.includes('/auth/login')) {
            originalRequest._retry = true
            const result = await authStore.refreshToken()

            if (result === true) {
                originalRequest.headers = originalRequest.headers || {}
                originalRequest.headers['Authorization'] = `Bearer ${authStore.token}`
                return axios(originalRequest)
            }

            const message = error.response?.data?.message || 'expired'
            await authStore.logout()
            router.push({ name: 'login', query: { reason: message } })
        }

        if (error.response?.status === 403) {
            const message = error.response.data?.message || 'forbidden'
            await authStore.logout()
            router.push({ name: 'login', query: { reason: message } })
        }

        return Promise.reject(error)
    }
)

cleanupLegacyOfflineArtifacts()

authStore.initializeAuth().finally(() => {
    app.mount('#app')
})
