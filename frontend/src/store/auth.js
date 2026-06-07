import { defineStore } from 'pinia'
import axios from 'axios'

export const useAuthStore = defineStore('auth', {
    state: () => ({
        user: JSON.parse(localStorage.getItem('user')) || null,
        token: null, 
        loading: false,
        error: null,
        isInitialized: false,
        isOffline: false,
        initPromise: null,
        refreshPromise: null 
    }),

    getters: {
        isAuthenticated: (state) => Boolean((state.token && state.user) || (state.isOffline && state.user)),
        isAdmin: (state) => state.user?.role === 'admin',
        userRole: (state) => state.user?.role
    },

    actions: {
        async login(email, password) {
            this.loading = true
            this.error = null
            try {
                const response = await axios.post('/auth/login', { email, password })
                const { access_token, user } = response.data.data
                this.token = access_token
                this.user = user
                localStorage.setItem('user', JSON.stringify(user))
                return { success: true }
            } catch (err) {
                if (err.response?.status === 401) { 
                    this.error = "Email atau Password Salah" 
                } else if (err.response?.status === 429) {
                    const retryAfter = Number(err.response?.data?.data?.retry_after_seconds || err.response?.headers?.['retry-after'] || 0)
                    this.error = retryAfter > 0
                        ? `Terlalu banyak percobaan. Coba lagi dalam ${retryAfter} detik.`
                        : (err.response?.data?.message || "Terlalu banyak percobaan. Coba lagi sebentar lagi.")
                } else { 
                    this.error = err.response?.data?.message || "Login gagal" 
                }
                return { success: false, error: err }
            } finally {
                this.loading = false
            }
        },

        async refreshToken() {
            if (this.refreshPromise) return this.refreshPromise
            this.refreshPromise = (async () => {
                try {
                    const response = await axios.post("/auth/refresh")
                    const { access_token } = response.data.data
                    this.token = access_token
                    return true
                } catch (err) {
                    const isNetworkError = err.code === "ERR_NETWORK" || !err.response;
                    if (isNetworkError) {
                        this.isOffline = true
                        return "network_error"
                    } else {
                        this.isOffline = false
                        this.clearAuth()
                        return false
                    }
                } finally {
                    this.refreshPromise = null
                }
            })()
            return this.refreshPromise
        },

        async logout() {
            try {
                await axios.post('/auth/logout')
            } catch (err) {
                console.error('Logout API failed', err)
            } finally {
                this.clearAuth()
            }
        },

        clearAuth() {
            this.token = null
            this.user = null
            this.isInitialized = false
            this.initPromise = null
            this.refreshPromise = null
            localStorage.removeItem('user')
        },

        async initializeAuth() {
            if (this.isInitialized && this.token) return true
            if (this.initPromise) return this.initPromise
            this.initPromise = (async () => {
                try {
                    const result = await this.refreshToken()
                    this.isInitialized = true
                    return result === true || result === "network_error"
                } catch (e) {
                    this.isInitialized = true
                    return false
                } finally {
                    this.initPromise = null
                }
            })()
            return this.initPromise
        }
    }
})
