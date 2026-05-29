import { beforeEach, describe, expect, it, vi } from 'vitest'
import axios from 'axios'
import financeService from '../finance.service'

vi.mock('axios', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn()
  }
}))

beforeEach(() => {
  vi.clearAllMocks()
})

describe('financeService API contract', () => {
  it('sends billing rule generation payload expected by backend', () => {
    financeService.generateBills(7, 'manual audit', 'pesan custom', true)

    expect(axios.post).toHaveBeenCalledWith('finance/generate-bills', {
      rule_id: 7,
      custom_reason: 'manual audit',
      custom_message: 'pesan custom',
      skip_notification: true
    })
  })

  it('sends payment payload without mutating caller data', () => {
    const payload = {
      student_id: 3,
      amount: 150000,
      bill_ids: [11, 12],
      is_bypass_rule: true,
      bypass_reason: 'approved by admin'
    }

    financeService.processPayment(payload)

    expect(axios.post).toHaveBeenCalledWith('finance/payments', payload)
  })

  it('creates payment intent through the cross-role finance endpoint', () => {
    const payload = { student_id: 4, amount: 250000, bill_ids: [20] }

    financeService.createPaymentIntent(payload)

    expect(axios.post).toHaveBeenCalledWith('finance/payment-intent', payload)
  })
})
