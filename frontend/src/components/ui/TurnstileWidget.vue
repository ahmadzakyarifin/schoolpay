<script setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

const props = defineProps({
  siteKey: {
    type: String,
    required: true
  },
  modelValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'error', 'expired'])

const container = ref(null)
let widgetId = null

const loadTurnstileScript = () => {
  if (window.turnstile) return Promise.resolve()

  const existing = document.querySelector('script[data-turnstile-script="true"]')
  if (existing) {
    return new Promise((resolve, reject) => {
      existing.addEventListener('load', resolve, { once: true })
      existing.addEventListener('error', reject, { once: true })
    })
  }

  return new Promise((resolve, reject) => {
    const script = document.createElement('script')
    script.src = 'https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit'
    script.async = true
    script.defer = true
    script.dataset.turnstileScript = 'true'
    script.onload = resolve
    script.onerror = reject
    document.head.appendChild(script)
  })
}

const renderWidget = async () => {
  if (!props.siteKey || widgetId !== null) return

  await loadTurnstileScript()
  await nextTick()

  if (!container.value || !window.turnstile) return

  widgetId = window.turnstile.render(container.value, {
    sitekey: props.siteKey,
    theme: 'light',
    callback: token => emit('update:modelValue', token),
    'expired-callback': () => {
      emit('update:modelValue', '')
      emit('expired')
    },
    'error-callback': () => {
      emit('update:modelValue', '')
      emit('error')
    }
  })
}

const reset = () => {
  emit('update:modelValue', '')
  if (window.turnstile && widgetId !== null) {
    window.turnstile.reset(widgetId)
  }
}

onMounted(() => {
  renderWidget().catch(() => emit('error'))
})

onBeforeUnmount(() => {
  if (window.turnstile && widgetId !== null) {
    window.turnstile.remove(widgetId)
    widgetId = null
  }
})

watch(() => props.siteKey, () => {
  reset()
  renderWidget().catch(() => emit('error'))
})

defineExpose({ reset })
</script>

<template>
  <div class="min-h-[65px] flex items-center justify-center">
    <div ref="container"></div>
  </div>
</template>
