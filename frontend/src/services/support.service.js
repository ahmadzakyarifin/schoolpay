import axios from 'axios'

const supportService = {
  getConversations(params) {
    return axios.get('support/conversations', { params })
  },
  getMessages(id) {
    return axios.get(`support/conversations/${id}/messages`)
  },
  reply(id, message) {
    return axios.post(`support/conversations/${id}/reply`, { message })
  },
  assign(id) {
    return axios.patch(`support/conversations/${id}/assign`)
  },
  close(id) {
    return axios.patch(`support/conversations/${id}/close`)
  }
}

export default supportService
