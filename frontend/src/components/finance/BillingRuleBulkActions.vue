<script setup>
import { Trash as TrashIcon, RotateCcw as ResetIcon, Play as PlayIcon, Undo2 as UndoIcon } from 'lucide-vue-next'

const props = defineProps({
  selectedCount: Number,
  status: String
})

const emit = defineEmits(['delete', 'restore', 'generate', 'cancel-generate'])
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
          <button @click="emit('restore')" class="flex items-center gap-2 px-6 hover:bg-emerald-50 text-emerald-600 transition-all group h-full cursor-pointer">
            <ResetIcon class="w-4 h-4 group-hover:rotate-[-45deg] transition-transform duration-500" />
            <span class="text-[10px] font-black uppercase tracking-[0.1em]">Pulihkan Masal</span>
          </button>
        </div>
      </template>
      <!-- Active Bulk Actions -->
      <template v-else>
        <div class="flex items-center bg-white border border-indigo-100 rounded-full shadow-lg shadow-indigo-100/30 overflow-hidden h-[42px]">
          <div class="px-5 flex items-center border-r border-indigo-50 h-full">
            <span class="text-[10px] font-black text-indigo-600 uppercase tracking-widest whitespace-nowrap">
              {{ selectedCount }} Terpilih
            </span>
          </div>
          <button @click="emit('generate')" class="flex items-center gap-2 px-6 hover:bg-indigo-50 text-indigo-600 border-r border-indigo-50 transition-all group h-full cursor-pointer" title="Generate tagihan siswa untuk aturan terpilih">
            <PlayIcon class="w-4 h-4 group-hover:scale-110 transition-transform fill-current" />
            <span class="text-[10px] font-black uppercase tracking-[0.1em]">Generate Tagihan</span>
          </button>
          <button @click="emit('cancel-generate')" class="flex items-center gap-2 px-6 hover:bg-amber-50 text-amber-600 border-r border-indigo-50 transition-all group h-full cursor-pointer" title="Tarik/Batal tagihan siswa yang telah digenerate">
            <UndoIcon class="w-4 h-4 group-hover:scale-110 transition-transform" />
            <span class="text-[10px] font-black uppercase tracking-[0.1em]">Tarik Tagihan</span>
          </button>
          <button @click="emit('delete')" class="flex items-center gap-2 px-6 hover:bg-rose-50 text-rose-500 transition-all group h-full cursor-pointer" title="Hapus master aturan tagihan">
            <TrashIcon class="w-4 h-4 group-hover:scale-110 transition-transform" />
            <span class="text-[10px] font-black uppercase tracking-[0.1em]">Hapus Aturan</span>
          </button>
        </div>
      </template>
    </div>
  </transition>
</template>
