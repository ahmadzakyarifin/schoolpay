<script setup>
import { computed, watch } from 'vue'
import { 
  X as CloseIcon, 
  Check as CheckIcon, 
  CreditCard as BillIcon 
} from 'lucide-vue-next'

const props = defineProps({
  show: { type: Boolean, default: false },
  selectedStudent: { type: Object, required: true },
  selectedBills: { type: Array, default: () => [] },
  paymentForm: { type: Object, required: true },
  formattedCashierAmount: { type: String, default: '' },
  allocationResult: { type: Object, required: true },
  allocationError: { type: String, default: '' },
  submitting: { type: Boolean, default: false }
})

const emit = defineEmits(['close', 'update:formattedCashierAmount', 'confirm'])

const selectedDepositBalance = computed(() => Number(props.selectedStudent?.deposit_balance || 0))
const selectedPayableBills = computed(() => {
  const bills = props.selectedStudent?.bills || []
  return bills
    .filter(b => b.status !== 'paid' && b.status !== 'voided')
    .filter(b => props.selectedBills.length === 0 || props.selectedBills.includes(b.id))
})

const outstandingAmount = computed(() => {
  return selectedPayableBills.value.reduce((acc, b) => acc + (Number(b.amount || 0) - Number(b.total_paid || 0)), 0)
})

const hasOverdueSelectedBill = computed(() => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return selectedPayableBills.value.some(b => {
    if (!b.due_date) return false
    const due = new Date(b.due_date)
    due.setHours(0, 0, 0, 0)
    return due < today
  })
})

const totalPaymentInput = computed(() => (Number(props.paymentForm.amount) || 0) + (Number(props.paymentForm.deposit_applied) || 0))
const maxDepositUsable = computed(() => Math.min(selectedDepositBalance.value, outstandingAmount.value))
const cashOrGatewayAmount = computed(() => Math.max(0, outstandingAmount.value - (Number(props.paymentForm.deposit_applied) || 0)))

const syncDepositAsDiscount = () => {
  let deposit = Number(props.paymentForm.deposit_applied || 0)
  if (deposit < 0) deposit = 0
  if (deposit > maxDepositUsable.value) deposit = maxDepositUsable.value
  if (deposit !== Number(props.paymentForm.deposit_applied || 0)) {
    props.paymentForm.deposit_applied = deposit
  }
  props.paymentForm.amount = Math.max(0, outstandingAmount.value - deposit)
  if (props.paymentForm.channel === 'gateway' && props.paymentForm.amount <= 0) {
    props.paymentForm.channel = 'cash'
    props.paymentForm.method = 'Tunai'
  }
}

watch(() => props.paymentForm.deposit_applied, syncDepositAsDiscount)
watch(outstandingAmount, syncDepositAsDiscount)

const selectGateway = () => {
  props.paymentForm.channel = 'gateway'
  props.paymentForm.method = 'Midtrans'
}

const applyMaxDeposit = () => {
  if (maxDepositUsable.value <= 0) return
  props.paymentForm.deposit_applied = maxDepositUsable.value
  props.paymentForm.amount = Math.max(0, outstandingAmount.value - maxDepositUsable.value)
  if (props.paymentForm.channel === 'gateway' && props.paymentForm.amount <= 0) {
    props.paymentForm.channel = 'cash'
    props.paymentForm.method = 'Tunai'
  }
}

const clearDeposit = () => {
  props.paymentForm.deposit_applied = 0
  props.paymentForm.amount = outstandingAmount.value
}

const formatCurrency = (val) => {
  if (!val) return 'Rp 0'
  const clean = Number(val).toFixed(0)
  return 'Rp ' + clean.replace(/\B(?=(\d{3})+(?!\d))/g, '.')
}
</script>

<template>
  <Teleport to="body">
    <div v-if="show" class="fixed inset-0 z-[2000] flex items-center justify-center p-4 font-inter">
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-md" @click="emit('close')"></div>
      <div class="relative bg-white w-full max-w-lg rounded-2xl shadow-2xl overflow-hidden animate-scale-in flex flex-col max-h-[90vh]">
        <div class="p-8 border-b border-slate-50 flex items-center justify-between bg-slate-50/30 shrink-0">
          <div>
            <h3 class="font-black text-slate-800 text-xl tracking-tight">Kalkulator Pembayaran</h3>
            <p class="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] mt-1">{{ selectedStudent.student_name }}</p>
          </div>
          <button @click="emit('close')" class="p-2 hover:bg-white rounded-xl text-slate-400 cursor-pointer"><CloseIcon class="w-6 h-6" /></button>
        </div>

        <div class="flex-1 overflow-y-auto p-8 space-y-6">
          <div class="grid grid-cols-2 gap-4">
            <div class="bg-indigo-50/50 p-6 rounded-2xl border border-indigo-100">
              <p class="text-[8px] font-black uppercase tracking-widest text-indigo-400 mb-1">Total Tagihan</p>
              <p class="text-xl font-black text-indigo-600">{{ formatCurrency(selectedStudent.bills.filter(b => b.status !== 'paid' && b.status !== 'voided').filter(b => selectedBills.length > 0 ? selectedBills.includes(b.id) : true).reduce((acc, b) => acc + (b.amount - b.total_paid), 0)) }}</p>
            </div>
            <div class="bg-emerald-50/50 p-6 rounded-2xl border border-emerald-100">
              <p class="text-[8px] font-black uppercase tracking-widest text-emerald-400 mb-1">Total Item</p>
              <p class="text-xl font-black text-emerald-600">{{ selectedBills.length > 0 ? selectedBills.length : selectedStudent.bills.filter(b => b.status !== 'paid' && b.status !== 'voided').length }} <span class="text-[10px] uppercase">Item</span></p>
            </div>
          </div>

          <div class="space-y-3">
            <label class="label-tiny">Uang Diterima / Dibayar Setelah Saldo (Rp)</label>
            <div class="relative">
              <span class="absolute left-6 top-1/2 -translate-y-1/2 font-black text-indigo-600 text-lg">Rp</span>
              <input :value="formattedCashierAmount" @input="emit('update:formattedCashierAmount', $event.target.value)" type="text" class="w-full py-6 pl-16 pr-6 bg-slate-50 border-2 border-slate-100 rounded-3xl focus:bg-white focus:border-indigo-500 outline-none transition-all font-black text-2xl text-slate-700 shadow-inner" placeholder="0" />
            </div>
            <div v-if="allocationResult.change > 0 && allocationResult.allocation.length > 0 && !allocationError" class="flex items-center justify-between px-2 animate-fade-in">
              <span class="text-[10px] font-bold text-emerald-600 uppercase tracking-widest">Lebih bayar masuk saldo:</span>
              <span class="text-sm font-black text-emerald-600">{{ formatCurrency(allocationResult.change) }}</span>
            </div>
          </div>

          <div class="space-y-3 p-4 bg-indigo-50/40 border border-indigo-100 rounded-2xl">
            <div class="flex items-center justify-between gap-3">
              <div>
                <p class="label-tiny !pl-0">Potongan Saldo Deposit</p>
                <p class="text-[9px] font-bold text-slate-500 uppercase tracking-wider mt-1">Saldo tersedia: {{ formatCurrency(selectedDepositBalance) }}</p>
              </div>
              <div class="flex gap-2 shrink-0">
                <button type="button" @click="applyMaxDeposit" :disabled="maxDepositUsable <= 0" class="px-3 py-2 rounded-xl bg-white border border-indigo-100 text-[9px] font-black uppercase tracking-wider text-indigo-600 disabled:opacity-50 cursor-pointer">Pakai Saldo</button>
                <button type="button" @click="clearDeposit" :disabled="!paymentForm.deposit_applied" class="px-3 py-2 rounded-xl bg-white border border-slate-100 text-[9px] font-black uppercase tracking-wider text-slate-500 disabled:opacity-50 cursor-pointer">Reset</button>
              </div>
            </div>
            <input v-model.number="paymentForm.deposit_applied" type="number" min="0" :max="maxDepositUsable" class="w-full py-3 px-4 bg-white border border-indigo-100 rounded-2xl font-black text-sm text-slate-700 outline-none focus:border-indigo-500" />
            <div class="grid grid-cols-2 gap-3 text-[10px] font-black uppercase tracking-wider">
              <div class="text-slate-500">Saldo Dipakai: <span class="text-indigo-600">{{ formatCurrency(paymentForm.deposit_applied || 0) }}</span></div>
              <div class="text-slate-500 text-right">Sisa Bayar: <span class="text-slate-800">{{ formatCurrency(cashOrGatewayAmount) }}</span></div>
            </div>
            <div v-if="paymentForm.deposit_applied > selectedDepositBalance" class="text-[10px] font-black text-rose-600 uppercase tracking-wider">Saldo tidak mencukupi.</div>
            <div v-else-if="paymentForm.deposit_applied > outstandingAmount" class="text-[10px] font-black text-rose-600 uppercase tracking-wider">Saldo yang digunakan tidak boleh melebihi total tagihan.</div>
          </div>

          <div class="space-y-4">
            <label class="label-tiny">Pilih Metode Pembayaran Akhir</label>
            <div class="grid grid-cols-2 gap-4">
              <button 
                @click="paymentForm.channel = 'cash'; paymentForm.method = 'Tunai'"
                :class="['flex flex-col items-center justify-center p-6 rounded-2xl border-2 transition-all cursor-pointer', 
                  paymentForm.channel === 'cash' ? 'bg-indigo-50 border-indigo-500 ring-4 ring-indigo-50' : 'bg-white border-slate-100 hover:border-indigo-100']"
              >
                <div :class="['w-12 h-12 rounded-2xl flex items-center justify-center mb-3 transition-colors', 
                  paymentForm.channel === 'cash' ? 'bg-indigo-600 text-white' : 'bg-slate-100 text-slate-400']">
                  <CheckIcon class="w-6 h-6" />
                </div>
                <span :class="['text-[10px] font-black uppercase tracking-widest', paymentForm.channel === 'cash' ? 'text-indigo-600' : 'text-slate-400']">Tunai / Kasir</span>
              </button>


              <button 
                @click="selectGateway"
                :disabled="paymentForm.amount <= 0"
                title="Bayar melalui Midtrans setelah potongan saldo"
                :class="['flex flex-col items-center justify-center p-6 rounded-2xl border-2 transition-all cursor-pointer disabled:cursor-not-allowed disabled:opacity-50', 
                  paymentForm.channel === 'gateway' ? 'bg-indigo-50 border-indigo-500 ring-4 ring-indigo-50' : 'bg-white border-slate-100 hover:border-indigo-100']"
              >
                <div :class="['w-12 h-12 rounded-2xl flex items-center justify-center mb-3 transition-colors', 
                  paymentForm.channel === 'gateway' ? 'bg-indigo-600 text-white' : 'bg-slate-100 text-slate-400']">
                  <BillIcon class="w-6 h-6" />
                </div>
                <span :class="['text-[10px] font-black uppercase tracking-widest', paymentForm.channel === 'gateway' ? 'text-indigo-600' : 'text-slate-400']">Online (Midtrans)</span>
              </button>
            </div>
            <div v-if="hasOverdueSelectedBill" class="p-3 bg-amber-50 border border-amber-200 rounded-xl text-[10px] font-black text-amber-700 uppercase tracking-wider">
              Ada tagihan yang melewati jatuh tempo. Untuk orang tua, pembayaran online ditutup; admin tetap bisa mencatat pembayaran sesuai kebijakan sekolah.
            </div>
          </div>

          <div class="flex items-center justify-between p-4 bg-amber-50/50 rounded-2xl border border-amber-100">
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-xl bg-amber-100/50 flex items-center justify-center text-amber-600 font-bold shrink-0">
                ⚠️
              </div>
              <div>
                <p class="text-[10px] font-black text-amber-800 uppercase tracking-wider">Mode Kebijakan Khusus (Bypass Aturan)</p>
                <p class="text-[8px] font-bold text-amber-600 uppercase tracking-widest">Izinkan cicil tagihan yang diatur 'Tidak Bisa Dicicil'</p>
              </div>
            </div>
            <button @click="paymentForm.is_bypass_rule = !paymentForm.is_bypass_rule" type="button" :class="['w-10 h-5 rounded-full transition-all relative cursor-pointer shadow-inner shrink-0', paymentForm.is_bypass_rule ? 'bg-amber-500' : 'bg-slate-300']">
              <div :class="['absolute top-1 w-3 h-3 bg-white rounded-full transition-all shadow', paymentForm.is_bypass_rule ? 'left-6' : 'left-1']"></div>
            </button>
          </div>

          <div v-if="paymentForm.is_bypass_rule" class="space-y-2 animate-fade-in">
            <label class="label-tiny">Alasan Bypass (Wajib)</label>
            <textarea v-model="paymentForm.bypass_reason" rows="3" class="input-premium resize-none" placeholder="Contoh: orang tua meminta pembayaran sebagian karena kondisi khusus, disetujui admin."></textarea>
          </div>

          <div class="space-y-3 p-6 bg-slate-50 rounded-2xl border border-slate-100 shadow-inner">
            <div class="flex items-center justify-between mb-2">
              <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Alokasi Otomatis</span>
              <span class="text-[9px] font-bold text-indigo-600 uppercase">Prioritas: Overdue & Tgl Terlama</span>
            </div>
            <div class="space-y-2">
              <div v-for="item in allocationResult.allocation" :key="item.id" class="flex items-center justify-between bg-white p-3 rounded-xl border border-slate-100/50 shadow-sm animate-slide-right">
                <div class="flex items-center gap-3">
                  <div :class="['w-1.5 h-1.5 rounded-full', item.is_lunas ? 'bg-emerald-500' : 'bg-amber-500']"></div>
                  <div>
                    <p class="text-[9px] font-black text-slate-700 uppercase tracking-wider">{{ item.name }}</p>
                    <p class="text-[7px] font-bold text-slate-400 uppercase">
                      {{ item.is_lunas ? 'PELUNASAN' : 'DICICIL' }} • {{ item.period || 'SEKALI BAYAR' }}
                    </p>
                  </div>
                </div>
                <span class="text-[10px] font-black text-indigo-600">{{ formatCurrency(item.amount) }}</span>
              </div>
              <div v-if="allocationResult.allocation.length === 0" class="text-center py-4 opacity-30 italic text-[10px] font-bold">Masukkan nominal untuk melihat alokasi</div>
            </div>
            
            <div v-if="allocationError" class="p-3 bg-rose-50 border border-rose-100 rounded-xl text-[10px] font-black text-rose-600 animate-shake flex items-center gap-2">
              <span>🔴</span> {{ allocationError }}
            </div>
          </div>

          <div class="space-y-3">
            <div class="space-y-1">
              <label class="label-tiny">Nomor Referensi (Opsional)</label>
              <input v-model="paymentForm.reference" class="input-premium" placeholder="Nomor struk, nama pengirim, dsb." />
            </div>
            <div class="space-y-1">
              <label class="label-tiny">Catatan Audit (Opsional)</label>
              <textarea v-model="paymentForm.note" rows="2" class="input-premium resize-none" placeholder="Keterangan pembayaran untuk audit internal."></textarea>
            </div>
          </div>
        </div>

        <div class="p-8 border-t border-slate-50 bg-slate-50/50 shrink-0">
          <button @click="emit('confirm')" :disabled="submitting || totalPaymentInput <= 0 || !!allocationError || (paymentForm.is_bypass_rule && !paymentForm.bypass_reason?.trim()) || paymentForm.deposit_applied > selectedDepositBalance || paymentForm.deposit_applied > outstandingAmount || (paymentForm.channel === 'gateway' && paymentForm.amount <= 0)" 
            class="w-full py-5 bg-emerald-600 hover:bg-emerald-700 text-white font-black rounded-3xl shadow-xl shadow-emerald-100 flex items-center justify-center gap-3 transition-all active:scale-95 disabled:opacity-50 cursor-pointer">
            <div v-if="submitting" class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
            <CheckIcon v-else class="w-6 h-6" />
            <span class="text-xs uppercase tracking-[0.2em]">Konfirmasi Pembayaran</span>
          </button>
          <p class="text-[8px] text-center text-slate-400 font-bold uppercase mt-4 tracking-widest">Pastikan uang yang diterima sudah sesuai sebelum konfirmasi</p>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped lang="postcss">
.input-premium { @apply w-full py-4 px-5 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:bg-white focus:border-indigo-500 outline-none font-bold text-xs transition-all shadow-sm; }
.label-tiny { @apply text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1; }
.animate-fade-in { animation: fadeIn 0.5s ease-out; }
.animate-slide-right { animation: slideRight 0.3s ease-out backwards; }
.animate-scale-in { animation: scaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1); }

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideRight { from { opacity: 0; transform: translateX(-20px); } to { opacity: 1; transform: translateX(0); } }
@keyframes scaleIn { from { opacity: 0; transform: scale(0.95); } to { opacity: 1; transform: scale(1); } }

.animate-slide-right:nth-child(1) { animation-delay: 0.05s; }
.animate-slide-right:nth-child(2) { animation-delay: 0.1s; }
.animate-slide-right:nth-child(3) { animation-delay: 0.15s; }
.animate-slide-right:nth-child(4) { animation-delay: 0.2s; }
.animate-slide-right:nth-child(5) { animation-delay: 0.25s; }
</style>
