<template>
  <transition name="fade">
    <div v-if="errorList.length > 0" class="space-y-0.5 mt-1 animate-fade-in pl-1">
      <p v-for="(msg, index) in errorList" :key="index" class="text-[10px] font-bold text-rose-500 leading-tight">
        <span v-if="errorList.length > 1" class="mr-1">•</span>{{ msg }}
      </p>
    </div>
  </transition>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  message: [String, Array, Object]
})

const errorList = computed(() => {
  if (!props.message) return []
  if (Array.isArray(props.message)) return props.message.filter(m => !!m)
  if (typeof props.message === 'string') return [props.message]
  return [String(props.message)]
})
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s, transform 0.2s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
  transform: translateX(-5px);
}
</style>
