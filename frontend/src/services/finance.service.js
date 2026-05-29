import axios from 'axios'

const financeService = {
  getBills(params) {
    return axios.get('finance/bills', { params })
  },
  createBill(data) {
    return axios.post('finance/bills', data)
  },
  updateBill(id, data) {
    return axios.put(`finance/bills/${id}`, data)
  },
  deleteBill(id) {
    return axios.delete(`finance/bills/${id}`)
  },
  
  getBillTypes(params) {
    return axios.get('finance/bill-types', { params })
  },
  
  createBillType(data) {
    return axios.post('finance/bill-types', data)
  },
  
  updateBillType(id, data) {
    return axios.put(`finance/bill-types/${id}`, data)
  },

  deleteBillType(id) {
    return axios.delete(`finance/bill-types/${id}`)
  },
  toggleBillTypeStatus(id) {
    return axios.patch(`finance/bill-types/${id}/status`)
  },
  restoreBillType(id) {
    return axios.patch(`finance/bill-types/${id}/restore`)
  },
  bulkDeleteBillTypes(ids) {
    return axios.post('finance/bill-types/bulk-delete', { ids })
  },
  bulkRestoreBillTypes(ids) {
    return axios.patch('finance/bill-types/bulk-restore', { ids })
  },
  getBillTypeDependencyInfo(id) {
    return axios.get(`finance/bill-types/${id}/dependency-info`)
  },

  getBillingRules(params) {
    return axios.get('finance/billing-rules', { params })
  },
  createBillingRule(data) {
    return axios.post('finance/billing-rules', data)
  },
  updateBillingRule(id, data) {
    return axios.put(`finance/billing-rules/${id}`, data)
  },
  deleteBillingRule(id) {
    return axios.delete(`finance/billing-rules/${id}`)
  },
  toggleBillingRuleStatus(id) {
    return axios.patch(`finance/billing-rules/${id}/status`)
  },
  restoreBillingRule(id) {
    return axios.patch(`finance/billing-rules/${id}/restore`)
  },
  bulkDeleteBillingRules(ids) {
    return axios.post('finance/billing-rules/bulk-delete', { ids })
  },
  bulkRestoreBillingRules(ids) {
    return axios.patch('finance/billing-rules/bulk-restore', { ids })
  },
  getBillingRuleDependencyInfo(id) {
    return axios.get(`finance/billing-rules/${id}/dependency-info`)
  },

  generateBills(ruleId, customReason, customMessage, skipNotification) {
    return axios.post('finance/generate-bills', {
      rule_id: ruleId,
      custom_reason: customReason,
      custom_message: customMessage,
      skip_notification: skipNotification
    })
  },

  bulkGenerateBills(ruleIds, customReason, customMessage, skipNotification) {
    return axios.post('finance/generate-bills/bulk', {
      rule_ids: ruleIds,
      custom_reason: customReason,
      custom_message: customMessage,
      skip_notification: skipNotification
    })
  },
  
  processPayment(data) {
    return axios.post('finance/payments', data)
  },

  getClasses(params) {
    return axios.get('academic/class', { params })
  },
  getMajors(params) {
    return axios.get('academic/major', { params })
  },
  getBillsByStudent(studentId) {
    return axios.get(`finance/bills?student_id=${studentId}`)
  },
  getReceipt(paymentId) {
    return axios.get(`finance/payments/${paymentId}/receipt`)
  },
  getPaymentHistory(studentId) {
    return axios.get(`finance/payments?student_id=${studentId}`)
  },
  createPaymentIntent(data) {
    return axios.post('finance/payment-intent', data)
  },
  checkUniqueBillType(name, excludeId) {
    return axios.get('finance/bill-types/check-unique', {
      params: {
        name: name,
        exclude_id: excludeId
      }
    })
  },

  checkUniqueBillingRule(billTypeId, targetType, targetId, classId, excludeId) {
    return axios.get('finance/billing-rules/check-unique', {
      params: {
        bill_type_id: billTypeId,
        target_type: targetType,
        target_id: targetId,
        class_id: classId,
        exclude_id: excludeId
      }
    })
  }
}

export default financeService
