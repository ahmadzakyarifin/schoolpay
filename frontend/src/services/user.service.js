import axios from 'axios'

const userService = {
	checkUnique(field, value, excludeId = 0) {
		return axios.get('users/check-unique', { params: { field, value, exclude_id: excludeId } })
	},

	getAll(params) {
		return axios.get('users', { params })
	},

	getByID(id) {
		return axios.get(`users/${id}`)
	},
	
	getParents(params) {
		return axios.get('users/parents', { params })
	},
	
	create(data) {
		return axios.post('users', data)
	},
	
	update(id, data) {
		return axios.put(`users/${id}`, data)
	},
	
	delete(id) {
		return axios.delete(`users/${id}`)
	},
	
	toggleStatus(id) {
		return axios.patch(`users/${id}/status`)
	},
	
	resetPassword(id, data) {
		return axios.post(`users/${id}/reset-password`, data)
	},

	resendNotification(id, channel) {
		return axios.post(`users/${id}/resend-notification`, { channel })
	},

	bulkDelete(ids) {
		return axios.post('users/bulk-delete', { ids })
	},

	bulkResendNotification(ids, channel) {
		return axios.post('users/bulk-resend-notification', { ids, channel })
	},

	restore(id) {
		return axios.patch(`users/${id}/restore`)
	},

	bulkRestore(ids) {
		return axios.patch('users/bulk-restore', { ids })
	}
}

export default userService
