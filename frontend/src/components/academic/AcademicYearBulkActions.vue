<script setup>
import { 
  Trash as TrashIcon,
  Undo2 as RestoreIcon
} from 'lucide-vue-next'

defineProps({
  selectedCount: Number,
  status: String
})

defineEmits(['delete', 'restore'])
</script>

<template>
  <transition name="fade">
    <div v-if="selectedCount > 0" class="flex items-center bg-white border border-indigo-100 rounded-full shadow-lg shadow-indigo-100/30 overflow-hidden h-[42px]">
      <!-- Info Section -->
      <div class="px-5 flex items-center border-r border-indigo-50 h-full">
        <span class="text-[10px] font-black text-indigo-600 uppercase tracking-widest whitespace-nowrap">
          {{ selectedCount }} Terpilih
        </span>
      </div>

      <!-- Action Button -->
      <button 
        v-if="status !== 'trash'" 
        @click="$emit('delete')" 
        class="flex items-center gap-2 px-6 hover:bg-rose-50 text-rose-500 transition-all group h-full"
      >
        <TrashIcon class="w-4 h-4 group-hover:scale-110 transition-transform" />
        <span class="text-[10px] font-black uppercase tracking-[0.1em]">Hapus Masal</span>
      </button>
      
      <button 
        v-else 
        @click="$emit('restore')" 
        class="flex items-center gap-2 px-6 hover:bg-emerald-50 text-emerald-600 transition-all group h-full"
      >
        <RestoreIcon class="w-4 h-4 group-hover:rotate-[-45deg] transition-transform" />
        <span class="text-[10px] font-black uppercase tracking-[0.1em]">Pulihkan Masal</span>
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
