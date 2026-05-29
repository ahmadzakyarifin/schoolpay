import axios from 'axios'

const academicYearService = {
  getAll(params) {
    return axios.get('academic/years', { params })
  },
  create(data) {
    return axios.post('academic/years', data)
  },
  update(id, data) {
    return axios.put(`academic/years/${id}`, data)
  },
  delete(id) {
    return axios.delete(`academic/years/${id}`)
  },
  restore(id) {
    return axios.patch(`academic/years/${id}/restore`)
  },
  toggleStatus(id) {
    return axios.patch(`academic/years/${id}/status`)
  },
  getDependencyInfo(id) {
    return axios.get(`academic/years/${id}/dependency-info`)
  },
  bulkDelete(ids) {
    return axios.post('academic/years/bulk-delete', { ids })
  },
  bulkRestore(ids) {
    return axios.patch('academic/years/bulk-restore', { ids })
  },
  checkUnique(year, excludeId) {
    return axios.get('academic/years/check-unique', {
      params: {
        year: year,
        exclude_id: excludeId
      }
    })
  }
}

export default academicYearService
