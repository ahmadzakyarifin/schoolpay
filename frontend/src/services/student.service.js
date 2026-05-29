import axios from 'axios'

export default {
  checkUnique(field, value, excludeId = 0) {
    return axios.get('students/check-unique', { params: { field, value, exclude_id: excludeId } })
  },
  getAll(params) {
    return axios.get('students', { params })
  },
  getByID(id) {
    return axios.get(`students/${id}`)
  },
  create(data) {
    return axios.post('students', data)
  },
  update(id, data) {
    return axios.put(`students/${id}`, data)
  },
  delete(id) {
    return axios.delete(`students/${id}`)
  },
  toggleStatus(id) {
    return axios.patch(`students/${id}/status`)
  },
  export(params) {
    return axios.get('students/export', { params, responseType: 'blob' })
  },
  getStudentsByParentID(parentId) {
    return axios.get(`users/${parentId}/students`)
  },
  getParents(studentId) {
    return axios.get(`students/${studentId}/parents`)
  },
  getClassHistory(studentId) {
    return axios.get(`students/${studentId}/class-history`)
  },
  bulkPromote(data) {
    return axios.post('students/bulk-promote', data)
  },
  bulkGraduate(data) {
    return axios.post('students/bulk-graduate', data)
  },
  bulkDelete(ids) {
    return axios.post('students/bulk-delete', { ids })
  },
  restore(id) {
    return axios.patch(`students/${id}/restore`)
  },
  bulkRestore(ids) {
    return axios.patch('students/bulk-restore', { ids })
  }
}
