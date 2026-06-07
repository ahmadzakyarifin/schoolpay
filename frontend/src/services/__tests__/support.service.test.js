import { beforeEach, describe, expect, it, vi } from 'vitest'
import axios from 'axios'
import supportService from '../support.service'

vi.mock('axios', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn()
  }
}))

beforeEach(() => {
  vi.clearAllMocks()
})

describe('supportService API contract', () => {
  it('loads conversation queue with filters', () => {
    const params = { status: 'open', limit: 50 }

    supportService.getConversations(params)

    expect(axios.get).toHaveBeenCalledWith('support/conversations', { params })
  })

  it('assigns and closes tickets through PATCH endpoints', () => {
    supportService.assign(9)
    supportService.close(9)

    expect(axios.patch).toHaveBeenNthCalledWith(1, 'support/conversations/9/assign')
    expect(axios.patch).toHaveBeenNthCalledWith(2, 'support/conversations/9/close')
  })

  it('updates ticket status through dropdown endpoint', () => {
    supportService.updateStatus(9, 'pending')

    expect(axios.patch).toHaveBeenCalledWith('support/conversations/9/status', {
      status: 'pending'
    })
  })
})
