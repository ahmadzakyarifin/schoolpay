<script setup>
import { 
  X as CloseIcon, 
  Check as CheckIcon, 
  CreditCard as BillIcon, 
  User as StudentIcon, 
  Trash as TrashIcon 
} from 'lucide-vue-next'

const props = defineProps({
  selectedStudent: { type: Object, default: null },
  activeTab: { type: String, default: 'active' },
  selectedBills: { type: Array, default: () => [] },
  recapYears: { type: Array, default: () => [] },
  selectedRecapYear: { type: [Number, String], default: new Date().getFullYear() },
  months: { type: Array, default: () => ['01', '02', '03', '04', '05', '06', '07', '08', '09', '10', '11', '12'] },
  paymentHistory: { type: Array, default: () => [] }
})

const emit = defineEmits([
  'update:activeTab', 
  'update:selectedRecapYear', 
  'close', 
  'toggle-bill', 
  'reset-bill', 
  'print-receipt', 
  'open-payment',
  'pay-manual'
])

const formatCurrency = (val) => {
  if (!val) return 'Rp 0'
  const clean = Number(val).toFixed(0)
  return 'Rp ' + clean.replace(/\B(?=(\d{3})+(?!\d))/g, '.')
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

const displayBillName = (bill) => {
  return bill?.name || (bill?.period ? `${bill.bill_type_name} ${formatPeriod(bill.period)}` : bill?.bill_type_name)
}

const paymentDetailNames = (payment) => {
  const details = payment?.details || []
  if (details.length === 0) return []
  return details.map(detail => detail.bill_type_name || detail.bill_name || 'Tagihan')
}

const getBillByMonth = (month) => {
  if (!props.selectedStudent?.bills) return null
  const period = `${props.selectedRecapYear}-${month}`
  return props.selectedStudent.bills.find(b => b.period === period)
}
</script>

<template>
  <div class="h-full font-inter">
    <div v-if="selectedStudent" class="bg-white rounded-2xl border border-slate-200 shadow-xl overflow-hidden animate-slide-up sticky top-6">
      <div class="p-8 bg-indigo-600 text-white space-y-4">
        <div class="flex items-center justify-between">
          <span class="text-[10px] font-black uppercase tracking-[0.3em] opacity-80">Informasi Keuangan</span>
          <button @click="emit('close')" class="p-1 hover:bg-white/20 rounded-lg transition-colors cursor-pointer"><CloseIcon class="w-5 h-5" /></button>
        </div>
        <div class="flex items-center gap-5">
          <div class="flex flex-col">
            <h4 class="font-black text-xl uppercase tracking-wider">{{ selectedStudent.student_name }}</h4>
          </div>
        </div>
        <div class="grid grid-cols-3 gap-3 pt-2">
          <div class="bg-white/10 rounded-2xl p-4 border border-white/10 backdrop-blur-sm">
            <p class="text-[8px] font-black uppercase tracking-widest opacity-60 mb-1">Total Tunggakan</p>
            <p class="text-sm font-black">{{ formatCurrency(selectedStudent.bills?.reduce((acc, b) => acc + (b.status !== 'paid' && b.status !== 'voided' ? b.amount - b.total_paid : 0), 0) || 0) }}</p>
          </div>
          <div class="bg-white/10 rounded-2xl p-4 border border-white/10 backdrop-blur-sm">
            <p class="text-[8px] font-black uppercase tracking-widest opacity-60 mb-1">Item Aktif</p>
            <p class="text-sm font-black">{{ selectedStudent.bills?.filter(b => b.status !== 'paid' && b.status !== 'voided').length || 0 }} Item</p>
          </div>
          <div class="bg-white/10 rounded-2xl p-4 border border-white/10 backdrop-blur-sm">
            <p class="text-[8px] font-black uppercase tracking-widest opacity-60 mb-1">Saldo</p>
            <p class="text-sm font-black">{{ formatCurrency(selectedStudent.deposit_balance || 0) }}</p>
          </div>
        </div>
      </div>

      <div class="p-4 border-b border-slate-50 flex items-center gap-2">
        <button v-for="t in ['active', 'recap', 'history']" :key="t" @click="emit('update:activeTab', t)"
          :class="['flex-1 py-3 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all cursor-pointer', activeTab === t ? 'bg-slate-100 text-indigo-600 shadow-inner' : 'text-slate-400 hover:text-slate-600']">
          {{ t === 'active' ? 'Aktif' : t === 'recap' ? 'Rekap' : 'Riwayat' }}
        </button>
      </div>

      <div class="p-6 h-[400px] overflow-y-auto">
        <div v-if="activeTab === 'active'" class="space-y-3">
          <div v-for="bill in selectedStudent.bills?.filter(b => b.status !== 'paid' && b.status !== 'voided')" :key="bill.id"
            @click="emit('toggle-bill', bill.id)"
            :class="['group p-4 rounded-2xl border transition-all cursor-pointer flex items-center gap-4', selectedBills.includes(bill.id) ? 'bg-indigo-50 border-indigo-200' : 'bg-white border-slate-100 hover:border-indigo-100']">
            <div :class="['w-5 h-5 rounded-md border flex items-center justify-center transition-all', selectedBills.includes(bill.id) ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300 bg-white']">
              <CheckIcon v-if="selectedBills.includes(bill.id)" class="w-3 h-3 text-white" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-[10px] font-black text-slate-700 uppercase tracking-wider truncate">{{ displayBillName(bill) }}</p>
              <div class="flex items-center gap-2 mt-0.5">
                <p class="text-[8px] font-bold text-slate-400 uppercase tracking-widest">PERIODE {{ formatPeriod(bill.period) }}</p>
                <span v-if="bill.allow_installment" class="px-1.5 py-0.5 bg-indigo-100 text-indigo-600 rounded text-[7px] font-black uppercase">
                  Bisa Dicicil (Maks {{ bill.max_installment }}x)
                </span>
                <span v-if="bill.status === 'overpaid'" class="px-1.5 py-0.5 bg-indigo-100 text-indigo-600 rounded text-[7px] font-black uppercase">Lebih Bayar</span>
              </div>
            </div>
            <div class="text-right shrink-0">
              <div class="flex flex-col items-end">
                <p class="text-xs font-black text-slate-800">{{ formatCurrency(bill.amount - bill.total_paid) }}</p>
                <div class="flex items-center gap-1.5 mt-0.5">
                  <span class="text-[7px] font-bold text-slate-400 uppercase">Terbayar: {{ formatCurrency(bill.total_paid) }}</span>
                  <span v-if="bill.total_paid > 0 && bill.status !== 'paid' && bill.status !== 'voided'" class="px-1 py-0.5 bg-amber-100 text-amber-600 rounded-[4px] text-[6px] font-black uppercase">
                    Cicilan
                  </span>
                  <span :class="['text-[8px] font-black uppercase tracking-widest', 
                    bill.status === 'overdue' ? 'text-rose-500' : 
                    bill.status === 'overpaid' ? 'text-indigo-500' : 
                    bill.status === 'voided' ? 'text-slate-400 line-through' :
                    'text-amber-500']">
                    {{ bill.status }}
                  </span>
                </div>
                <button @click.stop="emit('pay-manual', bill)" 
                  class="mt-1.5 px-2 py-0.5 bg-emerald-50 hover:bg-emerald-100 text-emerald-600 border border-emerald-200 rounded-[4px] text-[7px] font-black uppercase tracking-wider transition-all cursor-pointer"
                  title="Catat pembayaran manual sebagai transaksi dan alokasi">
                  Bayar Manual
                </button>
              </div>
            </div>
          </div>
          <div v-if="selectedStudent.bills?.filter(b => b.status !== 'paid' && b.status !== 'voided').length === 0" class="flex flex-col items-center justify-center py-20 text-center opacity-30">
            <CheckIcon class="w-12 h-12 mb-2" />
            <p class="text-xs font-black uppercase tracking-widest">Semua Tagihan Lunas</p>
          </div>
        </div>

        <!-- Recap Tab (Monthly Matrix) -->
        <div v-if="activeTab === 'recap'" class="p-4 space-y-4">
          <div class="mb-4">
            <h4 class="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Tahun Akademik</h4>
            <select :value="selectedRecapYear" @change="emit('update:selectedRecapYear', Number($event.target.value))" class="w-full bg-slate-100/80 border border-slate-200/60 rounded-2xl px-3 py-2.5 text-xs font-black text-slate-700 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all cursor-pointer">
              <option v-for="y in recapYears" :key="y" :value="y">{{ y }}</option>
            </select>
          </div>

          <div class="grid grid-cols-3 gap-2">
            <div v-for="month in months" :key="month" 
              class="group relative p-3 rounded-2xl border border-slate-50 bg-white shadow-sm flex flex-col items-center text-center transition-all hover:border-indigo-100"
            >
              <!-- Delete/Reset Button -->
              <button v-if="getBillByMonth(month)" 
                @click="emit('reset-bill', getBillByMonth(month).id)"
                class="absolute -top-1 -right-1 w-5 h-5 bg-white border border-slate-100 rounded-full flex items-center justify-center text-slate-300 hover:text-rose-500 hover:border-rose-100 shadow-sm opacity-0 group-hover:opacity-100 transition-all scale-75 group-hover:scale-100 cursor-pointer"
                title="Hapus / Reset Tagihan"
              >
                <TrashIcon class="w-3 h-3" />
              </button>

              <span class="text-[8px] font-black text-slate-300 uppercase mb-2">{{ month }}</span>
              
              <template v-if="getBillByMonth(month)">
                <div :class="[
                  'w-2 h-2 rounded-full mb-1 shadow-sm',
                  getBillByMonth(month).status === 'paid' ? 'bg-emerald-400 shadow-emerald-200' : 
                  getBillByMonth(month).status === 'partial' ? 'bg-amber-400 shadow-amber-200' : 
                  getBillByMonth(month).status === 'voided' ? 'bg-slate-300 shadow-slate-100' : 'bg-rose-400 shadow-rose-200'
                ]"></div>
                <span :class="[
                  'text-[9px] font-black uppercase tracking-tighter',
                  getBillByMonth(month).status === 'paid' ? 'text-emerald-600' : 
                  getBillByMonth(month).status === 'partial' ? 'text-amber-600' : 
                  getBillByMonth(month).status === 'voided' ? 'text-slate-400 line-through' : 'text-rose-600'
                ]">
                  {{ getBillByMonth(month).status === 'paid' ? 'Lunas' : 
                     getBillByMonth(month).status === 'partial' ? 'Cicil' : 
                     getBillByMonth(month).status === 'voided' ? 'Batal' : 'Unpaid' }}
                </span>
                <span class="text-[7px] text-slate-400 mt-0.5">
                  {{ formatCurrency(getBillByMonth(month).total_paid) }}
                </span>
              </template>
              <template v-else>
                <div class="w-2 h-2 rounded-full bg-slate-100 mb-1"></div>
                <span class="text-[9px] font-black text-slate-200 uppercase">-</span>
              </template>
            </div>
          </div>
        </div>

        <div v-if="activeTab === 'history'" class="space-y-3">
          <div v-for="pay in paymentHistory" :key="pay.id"
            class="p-4 rounded-2xl border border-slate-100 bg-white flex items-center gap-4 hover:shadow-md transition-all">
            <div class="w-10 h-10 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center border border-emerald-100 shrink-0">
              <CheckIcon class="w-5 h-5" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-[10px] font-black text-slate-700 uppercase tracking-wider truncate">Pembayaran #{{ pay.id }}</p>
              <p class="text-[8px] font-bold text-slate-400 uppercase tracking-widest">{{ new Date(pay.paid_at).toLocaleString() }} • {{ pay.method }}</p>
              <div v-if="paymentDetailNames(pay).length" class="mt-2 flex flex-wrap gap-1.5">
                <span v-for="name in paymentDetailNames(pay).slice(0, 2)" :key="name" class="max-w-[140px] truncate px-2 py-1 bg-slate-100 rounded-lg text-[7px] font-black text-slate-500 uppercase tracking-wider">
                  {{ name }}
                </span>
                <span v-if="paymentDetailNames(pay).length > 2" class="px-2 py-1 bg-indigo-50 rounded-lg text-[7px] font-black text-indigo-600 uppercase tracking-wider">
                  +{{ paymentDetailNames(pay).length - 2 }}
                </span>
              </div>
            </div>
            <div class="text-right shrink-0">
              <p class="text-xs font-black text-slate-800">{{ formatCurrency(pay.amount) }}</p>
              <button @click="emit('print-receipt', pay.id)" class="mt-1 text-[9px] font-black text-indigo-600 uppercase hover:underline cursor-pointer">
                Cetak Struk
              </button>
            </div>
          </div>
          <div v-if="paymentHistory.length === 0" class="flex flex-col items-center justify-center py-20 text-center opacity-30">
            <BillIcon class="w-12 h-12 mb-2" />
            <p class="text-[10px] font-black uppercase tracking-widest">Belum Ada Riwayat Bayar</p>
          </div>
        </div>
      </div>

      <div v-if="selectedStudent.bills?.some(b => b.status !== 'paid' && b.status !== 'voided')" class="p-6 border-t border-slate-50 bg-slate-50/50">
        <button @click="emit('open-payment')" 
          class="w-full py-4 bg-indigo-600 hover:bg-indigo-700 text-white font-black rounded-2xl shadow-xl shadow-indigo-100 flex items-center justify-center gap-3 group transition-all cursor-pointer">
          <BillIcon class="w-5 h-5 group-hover:scale-110 transition-transform" />
          <span class="text-[11px] uppercase tracking-[0.2em]">Bayar {{ selectedBills.length > 0 ? selectedBills.length + ' Item Terpilih' : 'Semua Tagihan' }}</span>
        </button>
      </div>
    </div>

    <!-- Empty Detail State -->
    <div v-else class="h-full min-h-[500px] border-2 border-dashed border-slate-100 rounded-2xl flex flex-col items-center justify-center text-center p-12 opacity-50">
      <StudentIcon class="w-16 h-16 text-slate-200 mb-4" />
      <h4 class="text-sm font-black text-slate-300 uppercase tracking-[0.2em]">Pilih Siswa</h4>
      <p class="text-[10px] font-bold text-slate-200 uppercase tracking-widest mt-2">Pilih siswa di daftar sebelah kiri untuk melihat detail keuangan & melakukan pembayaran</p>
    </div>
  </div>
</template>

<style scoped lang="postcss">
.animate-slide-up { animation: slideUp 0.4s ease-out; }
@keyframes slideUp { from { opacity: 0; transform: translateY(30px); } to { opacity: 1; transform: translateY(0); } }
</style>
