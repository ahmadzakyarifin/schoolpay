import axios from 'axios'

const supportService = {
  getConversations(params) {
    return axios.get('support/conversations', { params })
  },
  assign(id) {
    return axios.patch(`support/conversations/${id}/assign`)
  },
  close(id) {
    return axios.patch(`support/conversations/${id}/close`)
  },
  updateStatus(id, status) {
    return axios.patch(`support/conversations/${id}/status`, { status })
  }
}

export default supportService
