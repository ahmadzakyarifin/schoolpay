import axios from 'axios';

class ClassService {
  getAll(params) {
    return axios.get('academic/class', { params });
  }

  getByID(id) {
    return axios.get(`academic/class/${id}`);
  }

  create(data) {
    return axios.post('academic/class', data);
  }

  update(id, data) {
    return axios.put(`academic/class/${id}`, data);
  }

  delete(id) {
    return axios.delete(`academic/class/${id}`);
  }

  restore(id) {
    return axios.patch(`academic/class/${id}/restore`);
  }

  toggleStatus(id) {
    return axios.patch(`academic/class/${id}/status`);
  }

  getDependencyInfo(id) {
    return axios.get(`academic/class/${id}/dependency-info`);
  }

  bulkDelete(ids) {
    return axios.post('academic/class/bulk-delete', { ids });
  }

  bulkRestore(ids) {
    return axios.patch('academic/class/bulk-restore', { ids });
  }

  suggestNextName(name) {
    return axios.get('academic/class/suggest-name', { params: { name } });
  }

  checkUnique(name, majorId, academicYearId, excludeId) {
    return axios.get('academic/class/check-unique', {
      params: {
        name: name,
        major_id: majorId,
        academic_year_id: academicYearId,
        exclude_id: excludeId
      }
    });
  }
}

export default new ClassService();
