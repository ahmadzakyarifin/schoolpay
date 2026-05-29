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
  <transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="opacity-0 translate-y-4"
    enter-to-class="opacity-100 translate-y-0"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="opacity-100 translate-y-0"
    leave-to-class="opacity-0 translate-y-4"
  >
    <div v-if="selectedCount > 0" class="flex items-center gap-2">
      <button 
        v-if="status !== 'trash'" 
        @click="$emit('delete')" 
        class="bg-white border border-rose-200 hover:bg-rose-50/50 hover:border-rose-300 text-rose-600 font-black py-2 px-4 rounded-xl text-xs flex items-center gap-2 transition-all shadow-sm cursor-pointer"
      >
        <TrashIcon class="w-3.5 h-3.5 text-rose-500" />
        <span>Hapus Terpilih ({{ selectedCount }})</span>
      </button>
      
      <button 
        v-else 
        @click="$emit('restore')" 
        class="bg-white border border-emerald-200 hover:bg-emerald-50/50 hover:border-emerald-300 text-emerald-600 font-black py-2 px-4 rounded-xl text-xs flex items-center gap-2 transition-all shadow-sm cursor-pointer"
      >
        <RestoreIcon class="w-3.5 h-3.5 text-emerald-500" />
        <span>Pulihkan Terpilih ({{ selectedCount }})</span>
      </button>
    </div>
  </transition>
</template>
