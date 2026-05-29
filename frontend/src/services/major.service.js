import axios from 'axios'

export default {
  getAll(params) {
    return axios.get('academic/major', { params })
  },
  create(data) {
    return axios.post('academic/major', data)
  },
  update(id, data) {
    return axios.put(`academic/major/${id}`, data)
  },
  delete(id) {
    return axios.delete(`academic/major/${id}`)
  },
  restore(id) {
    return axios.patch(`academic/major/${id}/restore`)
  },
  toggleStatus(id) {
    return axios.patch(`academic/major/${id}/status`)
  },
  getDependencyInfo(id) {
    return axios.get(`academic/major/${id}/dependency-info`)
  },
  bulkDelete(ids) {
    return axios.post('academic/major/bulk-delete', { ids })
  },
  bulkRestore(ids) {
    return axios.patch('academic/major/bulk-restore', { ids })
  },
  checkUnique(field, value, excludeId) {
    return axios.get('academic/major/check-unique', {
      params: {
        field: field,
        value: value,
        exclude_id: excludeId
      }
    });
  }
}
