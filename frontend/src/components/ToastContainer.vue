<script setup>
import { useToast } from '../composables/useToast'
import { 
  CheckCircle as SuccessIcon, 
  AlertCircle as ErrorIcon, 
  Info as InfoIcon, 
  AlertTriangle as WarningIcon,
  X as CloseIcon 
} from 'lucide-vue-next'

const { toasts, remove } = useToast()

const getIcon = (type) => {
  switch (type) {
    case 'success': return SuccessIcon
    case 'error': return ErrorIcon
    case 'warning': return WarningIcon
    default: return InfoIcon
  }
}

const getTypeClasses = (type) => {
  switch (type) {
    case 'success': return 'bg-white border-emerald-100 shadow-emerald-100 text-emerald-600'
    case 'error': return 'bg-white border-rose-100 shadow-rose-100 text-rose-600'
    case 'warning': return 'bg-white border-amber-100 shadow-amber-100 text-amber-600'
    default: return 'bg-white border-indigo-100 shadow-indigo-100 text-indigo-600'
  }
}

const getIconClasses = (type) => {
  switch (type) {
    case 'success': return 'bg-emerald-50 text-emerald-600'
    case 'error': return 'bg-rose-50 text-rose-600'
    case 'warning': return 'bg-amber-50 text-amber-600'
    default: return 'bg-indigo-50 text-indigo-600'
  }
}
</script>

<template>
  <Teleport to="body">
    <div class="fixed top-8 right-8 z-[9999] flex flex-col gap-4 pointer-events-none">
      <transition-group name="toast">
        <div 
          v-for="toast in toasts" 
          :key="toast.id"
          class="pointer-events-auto flex items-start gap-4 p-5 rounded-[2rem] border shadow-2xl min-w-[320px] max-w-md animate-slide-in overflow-hidden relative group"
          :class="getTypeClasses(toast.type)"
        >
          <div class="w-12 h-12 shrink-0 rounded-2xl flex items-center justify-center shadow-inner" :class="getIconClasses(toast.type)">
            <component :is="getIcon(toast.type)" class="w-6 h-6" />
          </div>
          
          <div class="flex-1 pr-6">
            <h4 class="text-sm font-black tracking-tight mb-1">{{ toast.title }}</h4>
            <p class="text-[11px] font-bold opacity-70 leading-relaxed">{{ toast.message }}</p>
          </div>

          <button @click="remove(toast.id)" class="absolute top-4 right-4 p-1 rounded-lg hover:bg-slate-50 transition-colors opacity-40 hover:opacity-100">
            <CloseIcon class="w-4 h-4" />
          </button>
        </div>
      </transition-group>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.5s cubic-bezier(0.68, -0.55, 0.265, 1.55);
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%) scale(0.9);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%) scale(0.9);
}

.toast-move {
  transition: transform 0.4s ease;
}
</style>
