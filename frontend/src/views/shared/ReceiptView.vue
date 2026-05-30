<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { Printer as PrinterIcon, ChevronLeft as BackIcon, CheckCircle2 as SuccessIcon, ShieldCheck as ShieldIcon } from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const receipt = ref(null)
const loading = ref(true)

const fetchReceipt = async () => {
  try {
    const res = await axios.get(`finance/payments/${route.params.id}/receipt`)
    receipt.value = res.data.data
    // Auto print after a short delay
    setTimeout(() => {
      window.print()
    }, 1000)
  } catch (err) {
    console.error('Failed to fetch receipt', err)
  } finally {
    loading.value = false
  }
}

const formatCurrency = (val) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(val)
}

const monthNames = {
  '01': 'Januari',
  '02': 'Februari',
  '03': 'Maret',
  '04': 'April',
  '05': 'Mei',
  '06': 'Juni',
  '07': 'Juli',
  '08': 'Agustus',
  '09': 'September',
  '10': 'Oktober',
  '11': 'November',
  '12': 'Desember'
}

const formatPeriod = (period) => {
  if (!period) return 'Sekali Bayar'
  const match = String(period).match(/^(\d{4})-(\d{2})$/)
  if (!match) return period
  return `${monthNames[match[2]] || match[2]} ${match[1]}`
}

const printWindow = () => {
  window.print()
}

const handleBack = () => {
  if (window.history.length > 1) {
    router.back()
  } else {
    window.close()
    setTimeout(() => {
      router.push({ name: 'login' })
    }, 300)
  }
}

onMounted(() => {
  fetchReceipt()
})
</script>

<template>
  <div class="min-h-screen bg-slate-900/60 backdrop-blur-md flex flex-col items-center py-12 px-4 print:p-0 print:bg-white print:backdrop-blur-none">
    <!-- Back Button (Hidden in Print) -->
    <div class="w-full max-w-md mb-8 flex justify-between items-center print:hidden bg-white/10 backdrop-blur-md p-4 rounded-2xl border border-white/20 shadow-lg">
      <button @click="handleBack" class="flex items-center gap-2 text-xs font-black text-white uppercase tracking-widest hover:text-indigo-200 transition-colors">
        <BackIcon class="w-4 h-4" />
        Kembali
      </button>
      <button @click="printWindow" class="flex items-center gap-2 px-5 py-2.5 bg-indigo-600 text-white rounded-xl text-[10px] font-black uppercase tracking-widest hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-500/30">
        <PrinterIcon class="w-4 h-4" />
        Cetak Struk
      </button>
    </div>

    <!-- Receipt Content -->
    <div v-if="receipt" class="w-full max-w-md bg-white rounded-2xl border border-slate-200 shadow-2xl print:shadow-none print:border-none print:rounded-none p-10 print:p-0 space-y-8 font-mono">
      <!-- Header -->
      <div class="text-center space-y-2 border-b-2 border-dashed border-slate-100 pb-8">
        <h1 class="text-2xl font-black text-slate-800 tracking-tighter">SCHOOLPAY</h1>
        <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Sistem Manajemen Keuangan Sekolah</p>
        <div class="text-[9px] text-slate-400 uppercase mt-4">
          Jl. Pendidikan No. 123, Jakarta Timur<br>
          Telp: (021) 1234567 • Email: support@schoolpay.id
        </div>
      </div>

      <!-- Receipt Info -->
      <div class="space-y-3">
        <div class="flex justify-between text-[11px]">
          <span class="text-slate-400 uppercase">No. Struk</span>
          <span class="font-black text-slate-800">{{ receipt.receipt_number }}</span>
        </div>
        <div class="flex justify-between text-[11px]">
          <span class="text-slate-400 uppercase">Tanggal</span>
          <span class="font-black text-slate-800">{{ new Date(receipt.paid_at).toLocaleString('id-ID') }}</span>
        </div>
        <div class="flex justify-between text-[11px]">
          <span class="text-slate-400 uppercase">Metode</span>
          <span class="font-black text-slate-800">{{ receipt.payment_method }}</span>
        </div>
        <div class="flex justify-between text-[11px] pt-4 border-t border-slate-50">
          <span class="text-slate-400 uppercase">Nama Siswa</span>
          <span class="font-black text-slate-800">{{ receipt.student_name }}</span>
        </div>
        <div class="flex justify-between text-[11px]">
          <span class="text-slate-400 uppercase">NIS</span>
          <span class="font-black text-slate-800">{{ receipt.nis || '-' }}</span>
        </div>
      </div>

      <!-- Items Table -->
      <div class="space-y-4">
        <div class="text-[10px] font-black text-slate-400 uppercase tracking-widest border-b border-slate-100 pb-2">Rincian Pembayaran</div>
        <div class="space-y-3">
          <div v-for="(item, idx) in receipt.items" :key="idx" class="flex justify-between gap-4">
            <div class="flex flex-col flex-1">
              <span class="text-[11px] font-black text-slate-800 uppercase tracking-tight">{{ item.bill_name }}</span>
              <span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Periode: {{ formatPeriod(item.period) }}</span>
            </div>
            <span class="text-[11px] font-black text-slate-800 whitespace-nowrap">{{ formatCurrency(item.amount) }}</span>
          </div>
        </div>
      </div>

      <!-- Footer / Total -->
      <div class="pt-8 border-t-2 border-dashed border-slate-100 space-y-6">
        <div class="flex justify-between items-center">
          <span class="text-lg font-black text-slate-400 uppercase">Total</span>
          <span class="text-2xl font-black text-indigo-600">{{ formatCurrency(receipt.amount) }}</span>
        </div>

        <div class="text-center space-y-4 py-8 border-t border-slate-50">
          <p class="text-[9px] font-bold text-slate-400 uppercase leading-relaxed mb-6">
            Terima kasih atas pembayaran Anda.<br>
            Simpan struk ini sebagai bukti pembayaran yang sah.<br>
            Data telah terverifikasi secara sistem pada unit keuangan.
          </p>
          <div class="pt-6 border-t border-dashed border-slate-100 flex flex-col items-center justify-center gap-2">
             <div class="w-12 h-12 bg-emerald-50 rounded-2xl flex items-center justify-center text-emerald-600 shadow-inner">
               <ShieldIcon class="w-7 h-7" />
             </div>
             <p class="text-[10px] font-black text-emerald-600 tracking-widest uppercase mt-1">VERIFIED BY SCHOOLPAY</p>
             <p class="text-[7px] text-slate-400 tracking-tighter uppercase font-sans">SECURE & VALIDATED TRANSACTION</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-else-if="loading" class="flex flex-col items-center gap-4 pt-20">
      <div class="w-12 h-12 border-4 border-slate-200 border-t-indigo-600 rounded-full animate-spin"></div>
      <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Menyiapkan Struk...</p>
    </div>
  </div>
</template>

<style scoped>
@media print {
  @page {
    margin: 0;
    size: 80mm 200mm;
  }
  body {
    margin: 0;
    padding: 0;
  }
}
</style>
