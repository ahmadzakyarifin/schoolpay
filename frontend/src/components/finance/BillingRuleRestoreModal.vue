<script setup>
import { RefreshCw as RestoreIcon } from 'lucide-vue-next'

const props = defineProps({
  modelValue: Boolean,
  isBulkRestore: Boolean,
  selectedCount: Number,
  restoreLoading: Boolean
})

const emit = defineEmits(['update:modelValue', 'confirm'])

const close = () => {
  emit('update:modelValue', false)
}
</script>

<template>
  <Teleport to="body">
    <transition name="fade">
      <div v-if="modelValue" class="fixed inset-0 z-[700] flex items-center justify-center p-6">
        <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="close"></div>
        <div class="white-card w-full max-w-md relative z-10 overflow-hidden shadow-[0_20px_50px_rgba(0,0,0,0.3)] animate-scale-in !rounded-[2.5rem] p-8 text-center bg-white font-inter">
          <div class="w-20 h-20 bg-indigo-50 text-indigo-600 rounded-[2rem] flex items-center justify-center mx-auto mb-6 border border-indigo-100 shadow-xl shadow-indigo-600/10">
            <RestoreIcon class="w-10 h-10" />
          </div>
          
          <h3 class="text-xl font-black text-slate-900 tracking-tight mb-2">
            {{ isBulkRestore ? 'Pulihkan Aturan Terpilih?' : 'Pulihkan Aturan Tagihan?' }}
          </h3>
          
          <p class="text-slate-500 text-[10px] font-bold uppercase tracking-widest mb-8 px-4 leading-relaxed">
            {{ isBulkRestore 
              ? `Apakah Anda yakin ingin memulihkan ${selectedCount} aturan terpilih kembali ke daftar aktif?` 
              : `Apakah Anda yakin ingin memulihkan aturan tagihan ini kembali ke daftar aktif?` 
            }}
          </p>

          <div class="grid grid-cols-2 gap-4">
            <button 
              @click="close" 
              class="py-4 bg-slate-100 text-slate-600 font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-slate-200 transition-all cursor-pointer"
            >
              Batalkan
            </button>
            <button 
              @click="emit('confirm')" 
              :disabled="restoreLoading"
              class="py-4 bg-indigo-600 text-white font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-600/20 disabled:opacity-50 cursor-pointer"
            >
              {{ restoreLoading ? 'Memulihkan...' : 'Ya, Pulihkan Aturan' }}
            </button>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<style scoped lang="postcss">
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
