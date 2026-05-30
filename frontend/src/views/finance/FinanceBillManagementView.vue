<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { 
  Plus as PlusIcon, 
  Search as SearchIcon,
  RotateCcw as ResetIcon,
  ChevronLeft as PrevIcon,
  ChevronRight as NextIcon,
  X as CloseIcon,
  Check as CheckIcon,
  CreditCard as BillIcon,
  Calendar as CalendarIcon,
  User as StudentIcon,
  Edit as EditIcon,
  Trash as TrashIcon,
  CheckCircle2 as SuccessIcon,
  AlertCircle as AlertIcon,
  Filter as FilterIcon
} from 'lucide-vue-next'
import financeService from '../../services/finance.service'
import { useForm } from '../../composables/useForm'
import axios from 'axios'
import { useRouter } from 'vue-router'
import StudentBillFilter from '../../components/finance/StudentBillFilter.vue'
import StudentBillTable from '../../components/finance/StudentBillTable.vue'
import StudentBillDetailPanel from '../../components/finance/StudentBillDetailPanel.vue'
import StudentBillPaymentModal from '../../components/finance/StudentBillPaymentModal.vue'
import StudentBillSuccessModal from '../../components/finance/StudentBillSuccessModal.vue'

const router = useRouter()
const bills = ref([])
const billTypes = ref([])
const loading = ref(false)
const showModal = ref(false)
const showPaymentModal = ref(false)
const isMounted = ref(false)
const selectedStudent = ref(null)
const activeTab = ref('active') // active, recap, history
const selectedRecapYear = ref(new Date().getFullYear())

const recapYears = computed(() => {
  if (!selectedStudent.value?.bills) return [new Date().getFullYear()]
  const years = new Set()
  selectedStudent.value.bills.forEach(b => {
    if (b.rule_start_date && b.rule_end_date) {
      const startY = new Date(b.rule_start_date).getFullYear()
      const endY = new Date(b.rule_end_date).getFullYear()
      for (let y = startY; y <= endY; y++) {
        years.add(y)
      }
    } else if (b.period) {
      const y = b.period.split('-')[0]
      years.add(parseInt(y))
    }
  })
  if (years.size === 0) years.add(new Date().getFullYear())
  return Array.from(years).sort().reverse()
})

const months = ['01', '02', '03', '04', '05', '06', '07', '08', '09', '10', '11', '12']

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

const paymentHistory = ref([])

const notification = reactive({ show: false, message: '', type: 'success' })
let notifTimeout
const showNotification = (msg, type = 'success') => {
  clearTimeout(notifTimeout)
  notification.show = true
  notification.message = msg
  notification.type = type
  notifTimeout = setTimeout(() => { notification.show = false }, 4000)
}

const auditAction = reactive({
  show: false,
  type: '',
  bill: null,
  reason: '',
  submitting: false
})

const closeAuditAction = () => {
  auditAction.show = false
  auditAction.type = ''
  auditAction.bill = null
  auditAction.reason = ''
  auditAction.submitting = false
}

// Payment Form
const paymentForm = reactive({
  amount: 0,
  deposit_applied: 0,
  channel: 'cash',
  method: 'Tunai',
  reference: '',
  bill_ids: [],
  is_bypass_rule: false,
  bypass_reason: '',
  note: '',
  proof_attachment: ''
})

const formattedCashierAmount = computed({
  get() {
    if (!paymentForm.amount) return ''
    const parts = paymentForm.amount.toString().split('.')
    return parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, '.')
  },
  set(val) {
    if (!val) {
      paymentForm.amount = 0
      return
    }
    const clean = val.replace(/[^0-9]/g, '')
    const parsed = parseInt(clean, 10)
    paymentForm.amount = isNaN(parsed) ? 0 : parsed
  }
})

// Checkbox Selection
const selectedBills = ref([])
const toggleBill = (billId) => {
  const index = selectedBills.value.indexOf(billId)
  if (index === -1) {
    selectedBills.value.push(billId)
  } else {
    selectedBills.value.splice(index, 1)
  }
}

// Pagination, Search & Filters
const page = ref(1)
const limit = ref(10)
const total = ref(0)
const search = ref('')

const showFilters = ref(false)
const tempFilters = reactive({
  sort: '',
  status: ''
})
const activeFilters = reactive({
  sort: '',
  status: ''
})

const applyFilters = () => {
  activeFilters.sort = tempFilters.sort
  activeFilters.status = tempFilters.status
  page.value = 1
  showFilters.value = false
  fetchBills()
}

const resetFilters = () => {
  tempFilters.sort = ''
  tempFilters.status = ''
  activeFilters.sort = ''
  activeFilters.status = ''
  search.value = ''
  page.value = 1
  showFilters.value = false
  fetchBills()
}

const totalPages = computed(() => Math.ceil(total.value / limit.value) || 1)

const paginatedBills = computed(() => {
  return bills.value
})

watch(search, () => {
  if (page.value === 1) {
    fetchBills()
  } else {
    page.value = 1
  }
})

watch(page, () => {
  fetchBills()
})

watch(limit, () => {
  if (page.value === 1) {
    fetchBills()
  } else {
    page.value = 1
  }
})

const resetSearch = () => {
  search.value = ''
  tempFilters.sort = ''
  tempFilters.status = ''
  activeFilters.sort = ''
  activeFilters.status = ''
  page.value = 1
  fetchBills()
}

const fetchBills = async () => {
  loading.value = true
  try {
    const res = await financeService.getBillSummaries({
      page: page.value,
      limit: limit.value,
      search: search.value || undefined,
      sort: activeFilters.sort || undefined,
      status: activeFilters.status || undefined
    })

    const payload = res.data.data || {}
    bills.value = Array.isArray(payload.data) ? payload.data : []
    total.value = Number(payload.total || bills.value.length)
  } catch (err) {
    console.error('Gagal mengambil data tagihan')
    bills.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const openStudentDetail = async (student) => {
  selectedStudent.value = student
  selectedBills.value = []
  loading.value = true
  try {
    const res = await financeService.getBillsByStudent(student.student_id)
    selectedStudent.value.bills = res.data.data || []
    if (selectedStudent.value.bills.length > 0) {
      selectedStudent.value.deposit_balance = Number(selectedStudent.value.bills[0].deposit_balance || selectedStudent.value.deposit_balance || 0)
    }
    
    // Also fetch payment history
    const historyRes = await financeService.getPaymentHistory(student.student_id)
    paymentHistory.value = historyRes.data.data || []
  } catch (err) {
    showNotification('Gagal mengambil detail tagihan siswa', 'error')
  } finally {
    loading.value = false
  }
}

const allocationError = ref('')

// Payment Allocation Calculator
const allocationResult = computed(() => {
  allocationError.value = ''
  const amount = (Number(paymentForm.amount) || 0) + (Number(paymentForm.deposit_applied) || 0)
  let remaining = amount
  const allocation = []
  
  // Sort target bills: overdue first, then by due_date
  const targets = [...(selectedStudent.value?.bills || [])]
    .filter(b => b.status !== 'paid' && b.status !== 'voided')
    .filter(b => selectedBills.value.length === 0 || selectedBills.value.includes(b.id))
    .sort((a, b) => {
      if (a.status === 'overdue' && b.status !== 'overdue') return -1
      if (a.status !== 'overdue' && b.status === 'overdue') return 1
      return new Date(a.due_date) - new Date(b.due_date)
    })

  for (const bill of targets) {
    if (remaining <= 0) break
    const toPay = bill.amount - bill.total_paid
    if (!paymentForm.is_bypass_rule && !bill.allow_installment && remaining < toPay) {
      allocationError.value = `Tagihan '${bill.bill_type_name}' tidak bisa dicicil. Wajib dibayar lunas ${formatCurrency(toPay)}.`
      break
    }
    const paying = Math.min(remaining, toPay)
    allocation.push({
      id: bill.id,
      name: displayBillName(bill),
      period: bill.period,
      amount: paying,
      is_lunas: paying >= toPay
    })
    remaining -= paying
  }

  return { allocation, change: remaining }
})

const handlePayManual = (bill) => {
  auditAction.show = true
  auditAction.type = 'manual-pay'
  auditAction.bill = bill
  auditAction.reason = ''
}

const handleResetBill = (billOrId) => {
  const bill = typeof billOrId === 'object' ? billOrId : selectedStudent.value?.bills?.find(item => item.id === billOrId)
  auditAction.show = true
  auditAction.type = 'reset-bill'
  auditAction.bill = bill || { id: billOrId }
  auditAction.reason = ''
}

const confirmAuditAction = async () => {
  const reason = auditAction.reason.trim()
  if (!reason) {
    showNotification('Alasan audit wajib diisi', 'error')
    return
  }
  if (!auditAction.bill?.id) {
    showNotification('Data tagihan tidak valid', 'error')
    return
  }

  loading.value = true
  auditAction.submitting = true
  try {
    if (auditAction.type === 'manual-pay') {
      await axios.post(`finance/bills/${auditAction.bill.id}/pay-manual`, { reason, note: reason, payment_method: 'Tunai' })
      showNotification('Pembayaran manual berhasil dicatat', 'success')
    } else {
      await financeService.deleteBill(auditAction.bill.id, { reason })
      showNotification('Tagihan berhasil direset/void', 'success')
    }
    await fetchBills()
    if (selectedStudent.value) {
      const updatedStudent = bills.value.find(s => s.student_id === selectedStudent.value.student_id) || selectedStudent.value
      await openStudentDetail(updatedStudent)
    }
    closeAuditAction()
  } catch (err) {
    const errMsg = err.response?.data?.message || 'Gagal memproses aksi audit'
    showNotification(errMsg, 'error')
  } finally {
    loading.value = false
    auditAction.submitting = false
  }
}

const openPaymentModal = () => {
  if (!selectedStudent.value) return
  paymentForm.channel = 'cash'
  paymentForm.method = 'Tunai'
  paymentForm.reference = ''
  paymentForm.is_bypass_rule = false
  paymentForm.bypass_reason = ''
  paymentForm.note = ''
  paymentForm.proof_attachment = ''
  paymentForm.deposit_applied = 0
  paymentForm.amount = selectedStudent.value.bills
    .filter(b => selectedBills.value.length > 0 ? selectedBills.value.includes(b.id) : (b.status !== 'paid' && b.status !== 'voided'))
    .reduce((acc, b) => acc + (b.amount - b.total_paid), 0)
  showPaymentModal.value = true
}

const handleProcessPayment = async () => {
  const totalPaymentAmount = Number(paymentForm.amount || 0) + Number(paymentForm.deposit_applied || 0)
  if (totalPaymentAmount <= 0) {
    showNotification('Masukkan nominal pembayaran atau saldo yang digunakan', 'error')
    return
  }
  if (paymentForm.channel === 'gateway' && Number(paymentForm.amount || 0) <= 0) {
    showNotification('Pembayaran Midtrans membutuhkan sisa nominal setelah potongan saldo', 'error')
    return
  }
  
  submitting.value = true
  try {
    const payload = {
      student_id: selectedStudent.value.student_id,
      amount: totalPaymentAmount,
      deposit_applied: Number(paymentForm.deposit_applied || 0),
      channel: paymentForm.channel,
      method: paymentForm.method,
      reference: paymentForm.reference,
      bill_ids: selectedBills.value.length > 0 ? selectedBills.value : [],
      is_bypass_rule: paymentForm.is_bypass_rule,
      bypass_reason: paymentForm.bypass_reason,
      note: paymentForm.note,
      proof_attachment: paymentForm.proof_attachment
    }

    if (paymentForm.channel === 'gateway') {
      payload.method = 'Midtrans'
    } else if (Number(paymentForm.amount || 0) <= 0 && Number(paymentForm.deposit_applied || 0) > 0) {
      payload.channel = 'deposit'
      payload.method = 'Saldo Deposit'
    }

    if (paymentForm.channel === 'gateway') {
      const res = await financeService.createPaymentIntent(payload)
      const snapToken = res.data.data?.snap_token
      
      if (snapToken && window.snap) {
        window.snap.pay(snapToken, {
          onSuccess: async () => {
            showPaymentModal.value = false
            lastPaymentId.value = res.data.data?.id
            showSuccessModal.value = true

            await fetchBills()
            if (selectedStudent.value) {
              const updatedStudent = bills.value.find(s => s.student_id === selectedStudent.value.student_id)
              if (updatedStudent) {
                await openStudentDetail(updatedStudent)
              }
            }
          },
          onPending: (result) => {
            showNotification('Menunggu pembayaran...', 'warning')
            showPaymentModal.value = false
            openStudentDetail(selectedStudent.value)
          },
          onError: (result) => {
            showNotification('Pembayaran gagal', 'error')
          },
          onClose: () => {
            showNotification('Anda menutup jendela pembayaran', 'info')
          }
        })
        return
      }
    }

    const res = await financeService.processPayment(payload)
    const paymentId = res.data.data?.id || res.data.data
    
    showPaymentModal.value = false
    lastPaymentId.value = paymentId
    showSuccessModal.value = true
 

    await fetchBills() 
    if (selectedStudent.value) {
      const updatedStudent = bills.value.find(s => s.student_id === selectedStudent.value.student_id)
      if (updatedStudent) {
        await openStudentDetail(updatedStudent)
      }
    }
  } catch (err) {
    showNotification(err.response?.data?.message || 'Gagal memproses pembayaran', 'error')
  } finally {
    submitting.value = false
  }
}


const submitting = ref(false)
const showSuccessModal = ref(false)
const lastPaymentId = ref(null)

const printReceipt = (paymentId) => {
  const url = router.resolve({ name: 'receipt-print', params: { id: paymentId } }).href
  window.open(url, '_blank')
}

onMounted(async () => {
  isMounted.value = true
  fetchBills()
  try {
    const res = await financeService.getBillTypes()
    billTypes.value = (res.data.data.data || res.data.data || []).filter(bt => bt.is_active)
  } catch (err) {
    console.error(err)
  }
})

// Utility
const formatCurrency = (val) => {
  if (!val) return 'Rp 0'
  const clean = Number(val).toFixed(0)
  return 'Rp ' + clean.replace(/\B(?=(\d{3})+(?!\d))/g, '.')
}
</script>

<template>
  <div class="space-y-6 animate-fade-in pb-20 p-2">
    <Teleport v-if="isMounted" to="#header-actions-target">
      <div class="flex items-center justify-center w-full gap-4 relative mx-auto font-inter">
        <div class="flex items-center justify-center gap-2 flex-1 max-w-2xl mx-auto">
          <div class="relative flex-1 group">
            <SearchIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300 group-focus-within:text-indigo-600 transition-colors" />
            <input v-model="search" type="text" placeholder="Cari nama siswa atau NIS..." class="search-input-premium" @keyup.enter="fetchBills" />
          </div>
          
          <!-- Filter Button -->
          <button @click="showFilters = !showFilters" class="relative p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 flex items-center gap-2 transition-all shadow-sm cursor-pointer">
            <FilterIcon class="w-4 h-4" />
            <span class="text-[10px] font-black uppercase tracking-wider pr-1">Filter</span>
            <span v-if="activeFilters.status || activeFilters.sort" class="absolute -top-1 -right-1 w-3 h-3 bg-indigo-600 rounded-full border-2 border-white shadow-sm"></span>
          </button>

          <button @click="resetFilters" class="p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm group shrink-0 transition-all cursor-pointer" title="Reset Pencarian">
            <ResetIcon class="w-4 h-4 group-hover:rotate-180 transition-transform duration-500" />
          </button>
        </div>

        <!-- Filter Dropdown Component -->
        <StudentBillFilter
          v-model="showFilters"
          :filters="tempFilters"
          @apply="applyFilters"
          @reset="resetFilters"
        />
      </div>
    </Teleport>

    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
      <!-- Main List: Students with Arrears -->
      <div class="lg:col-span-7 space-y-6">
        <StudentBillTable 
          :bills="paginatedBills" 
          :loading="loading" 
          :selectedStudent="selectedStudent" 
          :pagination="{ page, limit, total, totalPages }" 
          @select-student="openStudentDetail" 
          @update:page="page = $event" 
          @update:limit="limit = $event" 
        />
      </div>

      <!-- Detail Panel: Student Finance -->
      <div class="lg:col-span-5">
        <StudentBillDetailPanel 
          :selectedStudent="selectedStudent" 
          v-model:activeTab="activeTab" 
          :selectedBills="selectedBills" 
          :recapYears="recapYears" 
          v-model:selectedRecapYear="selectedRecapYear" 
          :months="months" 
          :paymentHistory="paymentHistory" 
          @close="selectedStudent = null" 
          @toggle-bill="toggleBill" 
          @reset-bill="handleResetBill" 
          @print-receipt="printReceipt" 
          @open-payment="openPaymentModal" 
          @pay-manual="handlePayManual"
        />
      </div>
    </div>

    <!-- Payment Modal -->
    <StudentBillPaymentModal 
      :show="showPaymentModal" 
      :selectedStudent="selectedStudent" 
      :selectedBills="selectedBills" 
      :paymentForm="paymentForm" 
      :formattedCashierAmount="formattedCashierAmount" 
      :allocationResult="allocationResult" 
      :allocationError="allocationError" 
      :submitting="submitting" 
      @close="showPaymentModal = false" 
      @update:formattedCashierAmount="formattedCashierAmount = $event" 
      @confirm="handleProcessPayment" 
    />

    <!-- Success Modal -->
    <StudentBillSuccessModal 
      :show="showSuccessModal" 
      :paymentId="lastPaymentId" 
      @close="showSuccessModal = false" 
      @print="printReceipt" 
    />

    <Teleport to="body">
      <transition name="fade">
        <div v-if="auditAction.show" class="fixed inset-0 z-[2500] bg-slate-950/50 backdrop-blur-sm flex items-center justify-center p-4">
          <div class="w-full max-w-md bg-white rounded-2xl border border-slate-200 shadow-2xl overflow-hidden animate-scale-in">
            <div class="px-6 py-5 border-b border-slate-100 flex items-start justify-between">
              <div>
                <p class="text-[10px] font-black uppercase tracking-[0.24em] text-slate-400">Konfirmasi Audit</p>
                <h3 class="mt-1 text-sm font-black text-slate-800">
                  {{ auditAction.type === 'manual-pay' ? 'Catat Pembayaran Manual' : 'Reset / Void Tagihan' }}
                </h3>
              </div>
              <button @click="closeAuditAction" class="p-2 rounded-xl hover:bg-slate-100 text-slate-400 transition-colors">
                <CloseIcon class="w-4 h-4" />
              </button>
            </div>

            <div class="p-6 space-y-5">
              <div class="rounded-xl bg-slate-50 border border-slate-100 p-4">
                <p class="text-[10px] font-black uppercase tracking-widest text-slate-400">Tagihan</p>
                <p class="mt-1 text-xs font-black text-slate-700">{{ displayBillName(auditAction.bill) || auditAction.bill?.bill_type_name || 'Tagihan' }}</p>
                <p class="mt-1 text-[10px] font-bold text-slate-500">
                  Sisa: {{ formatCurrency((auditAction.bill?.amount || 0) - (auditAction.bill?.total_paid || 0)) }}
                </p>
              </div>

              <div v-if="auditAction.type === 'reset-bill'" class="rounded-xl bg-amber-50 border border-amber-100 p-4">
                <p class="text-[10px] font-bold text-amber-700 leading-relaxed">
                  Tagihan akan di-void. Jika sudah ada pembayaran, nominal yang tercatat akan dikembalikan ke saldo deposit siswa.
                </p>
              </div>

              <label class="block space-y-2">
                <span class="text-[10px] font-black uppercase tracking-widest text-slate-400">Alasan Audit Wajib</span>
                <textarea
                  v-model="auditAction.reason"
                  rows="4"
                  class="w-full rounded-xl border border-slate-200 bg-white px-4 py-3 text-xs font-bold text-slate-700 outline-none focus:border-indigo-500 focus:ring-4 focus:ring-indigo-50"
                  placeholder="Contoh: pembayaran tunai di koperasi / koreksi tagihan ganda / beasiswa"
                ></textarea>
              </label>
            </div>

            <div class="px-6 py-5 bg-slate-50 border-t border-slate-100 flex items-center justify-end gap-3">
              <button @click="closeAuditAction" class="px-5 py-2.5 rounded-xl border border-slate-200 bg-white text-[10px] font-black uppercase tracking-widest text-slate-500 hover:bg-slate-50">
                Batal
              </button>
              <button
                @click="confirmAuditAction"
                :disabled="auditAction.submitting"
                class="px-6 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-widest text-white disabled:opacity-60"
                :class="auditAction.type === 'reset-bill' ? 'bg-rose-600 hover:bg-rose-700' : 'bg-indigo-600 hover:bg-indigo-700'"
              >
                {{ auditAction.submitting ? 'Memproses...' : 'Konfirmasi' }}
              </button>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>

    <!-- Notification -->
    <Teleport to="body">
      <transition name="fade">
        <div v-if="notification.show" class="fixed top-6 right-6 z-[3000] px-6 py-4 rounded-2xl shadow-2xl flex items-center gap-3 animate-scale-in"
          :class="notification.type === 'success' ? 'bg-emerald-600 text-white' : 'bg-rose-600 text-white'">
          <SuccessIcon v-if="notification.type === 'success'" class="w-5 h-5 shrink-0" />
          <AlertIcon v-else class="w-5 h-5 shrink-0" />
          <span class="text-xs font-bold">{{ notification.message }}</span>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<style scoped lang="postcss">
.search-input-premium {
  @apply w-full bg-white border border-slate-200 rounded-xl py-2.5 pl-12 pr-4 text-xs font-bold text-slate-700 outline-none transition-all focus:border-indigo-500 focus:ring-4 focus:ring-indigo-50 shadow-sm;
}
.input-premium { @apply w-full py-4 px-5 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:bg-white focus:border-indigo-500 outline-none font-bold text-xs transition-all shadow-sm; }
.label-tiny { @apply text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1; }
.animate-fade-in { animation: fadeIn 0.5s ease-out; }
.animate-slide-up { animation: slideUp 0.4s ease-out; }
.animate-slide-right { animation: slideRight 0.3s ease-out backwards; }
.animate-scale-in { animation: scaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1); }

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideUp { from { opacity: 0; transform: translateY(30px); } to { opacity: 1; transform: translateY(0); } }
@keyframes slideRight { from { opacity: 0; transform: translateX(-20px); } to { opacity: 1; transform: translateX(0); } }
@keyframes scaleIn { from { opacity: 0; transform: scale(0.95); } to { opacity: 1; transform: scale(1); } }

/* Staggered animation for list items */
.animate-slide-right:nth-child(1) { animation-delay: 0.05s; }
.animate-slide-right:nth-child(2) { animation-delay: 0.1s; }
.animate-slide-right:nth-child(3) { animation-delay: 0.15s; }
.animate-slide-right:nth-child(4) { animation-delay: 0.2s; }
.animate-slide-right:nth-child(5) { animation-delay: 0.25s; }
</style>
