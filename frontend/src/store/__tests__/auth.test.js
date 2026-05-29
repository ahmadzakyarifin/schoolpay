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
  localStorage.clear()
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
    expect(JSON.parse(localStorage.getItem('user'))).toEqual(user)
  })

  it('shows Indonesian wrong credential message on 401 login', async () => {
    axios.post.mockRejectedValueOnce({ response: { status: 401, data: { message: 'invalid' } } })

    const auth = useAuthStore()
    const result = await auth.login('parent@school.test', 'wrong')

    expect(result.success).toBe(false)
    expect(auth.error).toBe('Email atau Password Salah')
  })

  it('keeps cached user and enters offline mode when refresh cannot reach server', async () => {
    const cachedUser = { id: 2, name: 'Parent', role: 'parent' }
    localStorage.setItem('user', JSON.stringify(cachedUser))
    axios.post.mockRejectedValueOnce({ code: 'ERR_NETWORK' })

    const auth = useAuthStore()
    const result = await auth.refreshToken()

    expect(result).toBe('network_error')
    expect(auth.isOffline).toBe(true)
    expect(auth.user).toEqual(cachedUser)
    expect(auth.isAuthenticated).toBe(true)
  })

  it('clears auth on failed refresh with server response', async () => {
    localStorage.setItem('user', JSON.stringify({ id: 2, name: 'Parent', role: 'parent' }))
    axios.post.mockRejectedValueOnce({ response: { status: 401 } })

    const auth = useAuthStore()
    const result = await auth.refreshToken()

    expect(result).toBe(false)
    expect(auth.user).toBeNull()
    expect(auth.token).toBeNull()
    expect(localStorage.getItem('user')).toBeNull()
  })
})
