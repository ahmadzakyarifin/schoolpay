import { beforeEach, describe, expect, it, vi } from 'vitest'
import { createPinia, setActivePinia } from 'pinia'
import axios from 'axios'
import { useAuthStore } from '../auth'

vi.mock('axios', () => ({
  default: {
    post: vi.fn()
  }
}))

beforeEach(() => {
  setActivePinia(createPinia())
  vi.clearAllMocks()
})

describe('auth store', () => {
  it('stores user and token after successful login', async () => {
    const user = { id: 1, name: 'Admin', role: 'admin' }
    axios.post.mockResolvedValueOnce({ data: { data: { access_token: 'token-123', user } } })

    const auth = useAuthStore()
    const result = await auth.login('admin@school.test', 'password123')

    expect(result.success).toBe(true)
    expect(auth.token).toBe('token-123')
    expect(auth.user).toEqual(user)
    expect(auth.isAuthenticated).toBe(true)
  })

  it('shows Indonesian wrong credential message on 401 login', async () => {
    axios.post.mockRejectedValueOnce({ response: { status: 401, data: { message: 'invalid' } } })

    const auth = useAuthStore()
    const result = await auth.login('parent@school.test', 'wrong')

    expect(result.success).toBe(false)
    expect(auth.error).toBe('Email atau Password Salah')
  })

  it('updates user after successful refresh', async () => {
    const freshUser = { id: 2, name: 'Fresh Parent', email: 'parent@school.test', role: 'parent' }
    axios.post.mockResolvedValueOnce({ data: { data: { access_token: 'fresh-token', user: freshUser } } })

    const auth = useAuthStore()
    const result = await auth.refreshToken()

    expect(result).toBe(true)
    expect(auth.token).toBe('fresh-token')
    expect(auth.user).toEqual(freshUser)
    expect(auth.isAuthenticated).toBe(true)
  })

  it('clears auth when refresh cannot reach server', async () => {
    axios.post.mockRejectedValueOnce({ code: 'ERR_NETWORK' })

    const auth = useAuthStore()
    auth.user = { id: 2, name: 'Parent', role: 'parent' }
    auth.token = 'old-token'
    const result = await auth.refreshToken()

    expect(result).toBe(false)
    expect(auth.user).toBeNull()
    expect(auth.token).toBeNull()
    expect(auth.isAuthenticated).toBe(false)
  })

  it('clears auth on failed refresh with server response', async () => {
    axios.post.mockRejectedValueOnce({ response: { status: 401 } })

    const auth = useAuthStore()
    auth.user = { id: 2, name: 'Parent', role: 'parent' }
    auth.token = 'old-token'
    const result = await auth.refreshToken()

    expect(result).toBe(false)
    expect(auth.user).toBeNull()
    expect(auth.token).toBeNull()
  })
})
