<script setup>
import { 
  Trash as TrashIcon, 
  RotateCcw as ResetIcon
} from 'lucide-vue-next'

defineProps({
  selectedCount: Number,
  status: String
})

const emit = defineEmits(['delete', 'restore'])
</script>

<template>
  <transition name="fade">
    <div v-if="selectedCount > 0" class="flex items-center gap-2">
      <button 
        v-if="status !== 'trash'" 
        @click="emit('delete')" 
        class="bg-white border border-rose-200 hover:bg-rose-50/50 hover:border-rose-300 text-rose-600 font-black py-2 px-4 rounded-xl text-xs flex items-center gap-2 transition-all shadow-sm cursor-pointer"
      >
        <TrashIcon class="w-3.5 h-3.5 text-rose-500" />
        <span>Hapus Terpilih ({{ selectedCount }})</span>
      </button>
      
      <button 
        v-if="status === 'trash'" 
        @click="emit('restore')" 
        class="bg-white border border-emerald-200 hover:bg-emerald-50/50 hover:border-emerald-300 text-emerald-600 font-black py-2 px-4 rounded-xl text-xs flex items-center gap-2 transition-all shadow-sm cursor-pointer"
      >
        <ResetIcon class="w-3.5 h-3.5 text-emerald-500" />
        <span>Pulihkan Terpilih ({{ selectedCount }})</span>
      </button>
    </div>
  </transition>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
