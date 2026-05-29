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

const getBillByMonth = (month) => {
  if (!selectedStudent.value?.bills) return null
  const period = `${selectedRecapYear.value}-${month}`
  return selectedStudent.value.bills.find(b => b.period === period)
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
  const start = (page.value - 1) * limit.value
  const end = start + limit.value
  return bills.value.slice(start, end)
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
    const res = await financeService.getBills({
      page: page.value,
      limit: 50, // Increase limit for local grouping
      search: search.value || undefined,
      sort: activeFilters.sort || undefined
    })
    
    const rawBills = res.data.data || []
    
    // Grouping by student_id to show 1 row per student
    const grouped = rawBills.reduce((acc, bill) => {
      const sId = bill.student_id
      if (!acc[sId]) {
        acc[sId] = {
          id: bill.id, // For :key
          student_id: sId,
          student_name: bill.student_name,
          status: 'paid',
          amount: 0,
          total_paid: 0,
          bill_count: 0,
          deposit_balance: Number(bill.deposit_balance || 0)
        }
      }
      acc[sId].deposit_balance = Number(bill.deposit_balance || acc[sId].deposit_balance || 0)
      
      acc[sId].amount += bill.amount
      acc[sId].total_paid += (bill.total_paid || 0)
      acc[sId].bill_count++
      
      // If any bill is not paid, status is "unpaid"
      if (bill.status !== 'paid') {
        acc[sId].status = 'unpaid'
      }
      
      return acc
    }, {})

    let resultList = Object.values(grouped)

    if (activeFilters.status === 'paid') {
      resultList = resultList.filter(item => item.status === 'paid')
    } else if (activeFilters.status === 'unpaid') {
      resultList = resultList.filter(item => item.status !== 'paid')
    }

    bills.value = resultList
    total.value = bills.value.length
  } catch (err) {
    console.error('Gagal mengambil data tagihan')
    bills.value = []
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
      name: bill.bill_type_name,
      period: bill.period,
      amount: paying,
      is_lunas: paying >= toPay
    })
    remaining -= paying
  }

  return { allocation, change: remaining }
})

const handlePayManual = async (bill) => {
  const remaining = bill.amount - bill.total_paid
  const reason = prompt(`Catat pembayaran manual untuk tagihan "${bill.bill_type_name}" sebesar ${formatCurrency(remaining)}?\n\nMasukkan keterangan audit (wajib, misal: Tunai di koperasi, Beasiswa, Koreksi admin):`)
  if (reason === null) return // Canceled
  if (!reason.trim()) {
    showNotification('Alasan pelunasan manual wajib diisi untuk kebutuhan audit!', 'error')
    return
  }

  loading.value = true
  try {
    await axios.post(`finance/bills/${bill.id}/pay-manual`, { reason, note: reason, payment_method: 'Tunai' })
    if (selectedStudent.value) {
      await openStudentDetail(selectedStudent.value)
    }
    await fetchBills()
  } catch (err) {
    const errMsg = err.response?.data?.message || 'Gagal melunasi tagihan secara manual'
    showNotification(errMsg, 'error')
  } finally {
    loading.value = false
  }
}

const handleResetBill = async (id) => {
  if (!confirm('Apakah Anda yakin ingin menghapus/reset tagihan ini? Status pembayaran di bulan ini akan hilang.')) return
  
  try {
    await financeService.deleteBill(id)
    showNotification('Tagihan berhasil direset', 'success')
    if (selectedStudent.value) {
      openStudentDetail(selectedStudent.value)
    }
  } catch (err) {
    showNotification('Gagal reset tagihan', 'error')
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
