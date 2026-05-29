import { reactive } from 'vue'

const state = reactive({
  toasts: []
})

export const useToast = () => {
  const show = (title, message, type = 'info', duration = 5000) => {
    const id = Date.now()
    state.toasts.push({ id, title, message, type })
    
    setTimeout(() => {
      remove(id)
    }, duration)
  }

  const remove = (id) => {
    const index = state.toasts.findIndex(t => t.id === id)
    if (index > -1) {
      state.toasts.splice(index, 1)
    }
  }

  return {
    toasts: state.toasts,
    show,
    remove,
    success: (title, message) => show(title, message, 'success'),
    error: (title, message) => show(title, message, 'error'),
    info: (title, message) => show(title, message, 'info'),
    warning: (title, message) => show(title, message, 'warning')
  }
}
