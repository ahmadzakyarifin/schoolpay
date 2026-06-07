import { ref, reactive } from 'vue'

export function useForm(initialData = {}) {
  const form = reactive({ ...initialData })
  const errors = ref({})
  const submitting = ref(false)

  const setErrors = (err) => {
    // Jika input adalah object error mentah (dari frontend validation)
    if (err && !err.response && typeof err === 'object' && !Array.isArray(err)) {
        errors.value = err
        return
    }

    // Jika input dari Axios response
    if (err.response?.status === 429) {
      const retryAfter = Number(err.response?.data?.data?.retry_after_seconds || err.response?.headers?.['retry-after'] || 0)
      errors.value = {
        _general: [
          retryAfter > 0
            ? `Terlalu banyak request. Coba lagi dalam ${retryAfter} detik.`
            : (err.response?.data?.message || 'Terlalu banyak request. Coba lagi sebentar lagi.')
        ]
      }
    } else if (err.response?.data?.errors) {
      const rawErrors = err.response.data.errors
      const normalized = {}
      
      if (typeof rawErrors === 'string') {
        normalized._general = [rawErrors]
      } else if (typeof rawErrors === 'object') {
        Object.keys(rawErrors).forEach(key => {
          const val = rawErrors[key]
          normalized[key] = Array.isArray(val) ? val : [val]
        })
      }
      
      errors.value = normalized
    } else if (err.response?.data?.message) {
      errors.value = { _general: [err.response.data.message] }
    } else {
      console.error('Unhandled form error:', err)
    }
  }

  const clearErrors = () => {
    errors.value = {}
  }

  const clearFieldError = (field) => {
    if (errors.value[field]) {
      const newErrors = { ...errors.value }
      delete newErrors[field]
      errors.value = newErrors
    }
  }

  const resetForm = () => {
    Object.assign(form, initialData)
    clearErrors()
  }

  return {
    form,
    errors,
    submitting,
    setErrors,
    clearErrors,
    clearFieldError,
    resetForm
  }
}
