<script setup>
import { Trash as TrashIcon, RotateCcw as ResetIcon } from 'lucide-vue-next'

const props = defineProps({
  selectedCount: Number,
  status: String
})

const emit = defineEmits(['delete', 'restore'])
</script>

<template>
  <transition name="fade">
    <div v-if="selectedCount > 0" class="flex items-center gap-2">
      <!-- Trash Bulk Actions -->
      <template v-if="status === 'trash'">
        <div class="flex items-center bg-white border border-indigo-100 rounded-full shadow-lg shadow-indigo-100/30 overflow-hidden h-[42px]">
          <div class="px-5 flex items-center border-r border-indigo-50 h-full">
            <span class="text-[10px] font-black text-indigo-600 uppercase tracking-widest whitespace-nowrap">
              {{ selectedCount }} Terpilih
            </span>
          </div>
          <button @click="emit('restore')" class="flex items-center gap-2 px-6 hover:bg-emerald-50 text-emerald-600 transition-all group h-full">
            <ResetIcon class="w-4 h-4 group-hover:rotate-[-45deg] transition-transform duration-500" />
            <span class="text-[10px] font-black uppercase tracking-[0.1em]">Pulihkan Masal</span>
          </button>
        </div>
      </template>
      <!-- Active Bulk Delete -->
      <template v-else>
        <div class="flex items-center bg-white border border-indigo-100 rounded-full shadow-lg shadow-indigo-100/30 overflow-hidden h-[42px]">
          <div class="px-5 flex items-center border-r border-indigo-50 h-full">
            <span class="text-[10px] font-black text-indigo-600 uppercase tracking-widest whitespace-nowrap">
              {{ selectedCount }} Terpilih
            </span>
          </div>
          <button @click="emit('delete')" class="flex items-center gap-2 px-6 hover:bg-rose-50 text-rose-500 transition-all group h-full">
            <TrashIcon class="w-4 h-4 group-hover:scale-110 transition-transform" />
            <span class="text-[10px] font-black uppercase tracking-[0.1em]">Hapus Masal</span>
          </button>
        </div>
      </template>
    </div>
  </transition>
</template>
