<script setup>
import { computed } from 'vue'
import { Trash as TrashIcon, AlertCircle as AlertIcon } from 'lucide-vue-next'

const props = defineProps({
  modelValue: Boolean,
  isBulkDelete: Boolean,
  selectedCount: Number,
  dependencyLoading: Boolean,
  dependencyInfo: Object,
  deleteLoading: Boolean,
  showHistory: Boolean
})

const emit = defineEmits(['update:modelValue', 'confirm'])

const deleteBlocked = computed(() => !props.isBulkDelete && !!props.dependencyInfo?.has_dependencies)

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
          <div class="w-20 h-20 bg-rose-50 text-rose-500 rounded-[2rem] flex items-center justify-center mx-auto mb-6 border border-rose-100 shadow-xl shadow-rose-500/10">
            <TrashIcon class="w-10 h-10" />
          </div>
          
          <h3 class="text-xl font-black text-slate-900 tracking-tight mb-2">
            {{ isBulkDelete ? 'Hapus Aturan Terpilih?' : 'Hapus Aturan Tagihan?' }}
          </h3>

          <div v-if="dependencyLoading" class="my-4 py-3 px-4 bg-slate-50 rounded-2xl flex items-center justify-center gap-2 text-slate-500 text-[10px] font-bold uppercase tracking-widest">
            <div class="w-3 h-3 border-2 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
            Memeriksa keterhubungan data...
          </div>

          <div v-else-if="deleteBlocked" class="my-4 p-4 bg-amber-50 border border-amber-200/80 rounded-2xl text-left shadow-sm">
            <div class="flex items-start gap-3">
              <AlertIcon class="w-5 h-5 text-amber-600 shrink-0 mt-0.5" />
              <div>
                <h4 class="text-xs font-black text-amber-900 uppercase tracking-wider mb-1">Konsep Status & Penghapusan Aturan</h4>
                <p class="text-amber-800 text-[11px] font-medium leading-relaxed mb-2">
                  Aturan ini memiliki <span class="font-bold underline">{{ dependencyInfo.message }}</span> yang telah terbit ke siswa.
                </p>
                <ul class="text-amber-900 text-[10px] space-y-1 list-disc list-inside mb-2">
                  <li><b>Jika dinonaktifkan (OFF):</b> Sistem berhenti menerbitkan tagihan baru di masa depan, namun tagihan lama yang sudah terbit tetap aktif dan bisa dibayar siswa.</li>
                  <li><b>Jika dihapus (Trash):</b> Aturan dipindah ke riwayat terhapus. Tagihan siswa yang sudah terbit tetap aktif dan tidak ikut terhapus.</li>
                </ul>
                <p class="text-amber-700 text-[10px] font-bold bg-amber-100/60 py-1.5 px-3 rounded-lg block">
                  💡 Tips: Jika Anda ingin membatalkan/menarik kembali tagihan siswa yang belum dibayar, gunakan fitur <b>Tarik Tagihan</b> di menu atas.
                </p>
              </div>
            </div>
          </div>
          
          <p class="text-slate-500 text-[10px] font-bold uppercase tracking-widest mb-8 px-4 leading-relaxed">
            {{ isBulkDelete 
              ? `Apakah Anda yakin ingin menghapus ${selectedCount} aturan terpilih? Data akan dipindahkan ke riwayat penghapusan (Trash).` 
              : `Apakah Anda yakin ingin menghapus aturan tagihan ini? Data akan dipindahkan ke riwayat penghapusan (Trash).` 
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
              :disabled="deleteLoading || deleteBlocked"
              class="py-4 bg-rose-600 text-white font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-rose-700 transition-all shadow-lg shadow-rose-600/20 disabled:opacity-50 cursor-pointer"
            >
              {{ deleteLoading ? 'Menghapus...' : (deleteBlocked ? 'Tidak Bisa Dihapus' : (showHistory ? 'Ya, Hapus Permanen' : 'Ya, Hapus Aturan')) }}
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
