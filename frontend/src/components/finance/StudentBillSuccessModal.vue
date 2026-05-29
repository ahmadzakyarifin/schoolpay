<script setup>
import { 
  CheckCircle2 as SuccessIcon, 
  RotateCcw as ResetIcon 
} from 'lucide-vue-next'

defineProps({
  show: { type: Boolean, default: false },
  paymentId: { type: [Number, String], default: null }
})

const emit = defineEmits(['close', 'print'])
</script>

<template>
  <Teleport to="body">
    <div v-if="show" class="fixed inset-0 z-[3000] flex items-center justify-center p-4 font-inter">
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-md" @click="emit('close')"></div>
      <div class="relative bg-white w-full max-w-sm rounded-2xl shadow-2xl overflow-hidden animate-scale-in p-8 text-center">
        <div class="w-20 h-20 bg-emerald-100 text-emerald-600 rounded-3xl flex items-center justify-center mx-auto mb-6">
          <SuccessIcon class="w-10 h-10" />
        </div>
        <h3 class="text-xl font-black text-slate-800 mb-2">Pembayaran Berhasil!</h3>
        <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mb-8">Transaksi telah berhasil diproses dan dicatat dalam sistem</p>
        
        <div class="space-y-3">
          <button @click="emit('print', paymentId)" class="w-full py-4 bg-emerald-600 hover:bg-emerald-700 text-white font-black rounded-2xl flex items-center justify-center gap-3 transition-all cursor-pointer">
            <ResetIcon class="w-5 h-5 rotate-180" />
            <span class="text-[10px] uppercase tracking-widest">Cetak Struk Sekarang</span>
          </button>
          <button @click="emit('close')" class="w-full py-4 bg-slate-100 hover:bg-slate-200 text-slate-600 font-black rounded-2xl transition-all cursor-pointer">
            <span class="text-[10px] uppercase tracking-widest">Tutup</span>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped lang="postcss">
.animate-scale-in { animation: scaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1); }
@keyframes scaleIn { from { opacity: 0; transform: scale(0.95); } to { opacity: 1; transform: scale(1); } }
</style>
