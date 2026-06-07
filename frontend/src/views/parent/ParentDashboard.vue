<script setup>
import { ref, onMounted, computed, watch, reactive } from 'vue'
import axios from 'axios'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../../store/auth'
import { 
  User as UserIcon, 
  CreditCard as BillIcon, 
  Calendar as CalendarIcon,
  ChevronDown as ChevronDownIcon,
  CheckCircle2 as PaidIcon,
  AlertCircle as AlertIcon,
  Search as SearchIcon,
  Filter as FilterIcon,
  ArrowUpDown as SortIcon,
  RotateCcw as ResetIcon,
  ChevronLeft as PrevIcon,
  ChevronRight as NextIcon,
  Receipt as ReceiptIcon,
  MessageCircle as ChatIcon,
  Send as SendIcon,
  X as CloseIcon
} from 'lucide-vue-next'

const authStore = useAuthStore()
const route = useRoute()
const router = useRouter()
const students = ref([])
const selectedStudent = ref(null)
const allBills = ref([])
const paymentHistory = ref([])
const classHistory = ref([])
const loading = ref(false)
const payingBillId = ref(null)
const showStudentDropdown = ref(false)

// Filter & Pagination State
const search = ref('')
const statusFilter = ref('')
const academicYearFilter = ref('')
const classFilter = ref('')
const sortFilter = ref('')
const page = ref(1)
const limit = ref(10)
const isMounted = ref(false)
const notification = reactive({ show: false, message: '', type: 'success' })
let notificationTimer

const showNotification = (message, type = 'success') => {
  clearTimeout(notificationTimer)
  notification.message = message
  notification.type = type
  notification.show = true
  notificationTimer = setTimeout(() => {
    notification.show = false
  }, 4500)
}

const fetchParentData = async () => {
  loading.value = true
  try {
    const res = await axios.get('parent/students/me') 
    students.value = res.data.data
    
    if (students.value.length > 0) {
      selectedStudent.value = students.value[0]
      await Promise.all([
        fetchBills(),
        fetchClassHistory(selectedStudent.value.id),
        fetchPaymentHistory()
      ])
    }
  } catch (err) {
    console.error('Gagal mengambil data siswa')
  } finally {
    loading.value = false
  }
}

const fetchBills = async () => {
  loading.value = true
  try {
    const billsRes = await axios.get('finance/my-bills')
    allBills.value = billsRes.data.data
    syncSelectedStudentDeposit()
  } catch (err) {
    console.error('Gagal mengambil data tagihan')
  } finally {
    loading.value = false
  }
}

const fetchClassHistory = async (studentId = selectedStudent.value?.id) => {
  if (!studentId) {
    classHistory.value = []
    return
  }

  try {
    const res = await axios.get(`parent/students/${studentId}/class-history`)
    classHistory.value = res.data.data || []
  } catch (err) {
    classHistory.value = []
    showNotification(err.response?.data?.message || 'Gagal mengambil riwayat kelas anak', 'error')
  }
}

const syncSelectedStudentDeposit = () => {
  if (!selectedStudent.value) return
  const bill = allBills.value.find(item => item.student_id === selectedStudent.value.id)
  if (bill) {
    selectedStudent.value.deposit_balance = Number(bill.deposit_balance || selectedStudent.value.deposit_balance || 0)
  }
}

const fetchPaymentHistory = async () => {
  if (!selectedStudent.value?.id) {
    paymentHistory.value = []
    return
  }
  try {
    const res = await axios.get(`finance/my-payments?student_id=${selectedStudent.value.id}`)
    paymentHistory.value = res.data.data || []
  } catch (err) {
    paymentHistory.value = []
    showNotification(err.response?.data?.message || 'Gagal mengambil riwayat pembayaran', 'error')
  }
}

const contactWhatsApp = () => {
  const schoolWhatsApp = import.meta.env.VITE_SCHOOL_WHATSAPP || '6283120309758'
  const studentName = selectedStudent.value ? selectedStudent.value.name : ''
  const parentName = authStore.user?.name || ''
  
  let msg = 'Halo Admin, saya orang tua'
  if (parentName) msg += ` atas nama ${parentName}`
  if (studentName) msg += `, wali dari siswa ${studentName}`
  msg += '. Saya butuh bantuan/CS terkait layanan SchoolPay.'
  
  const text = encodeURIComponent(msg)
  const url = `https://wa.me/${schoolWhatsApp}?text=${text}`
  window.open(url, '_blank')
}

const selectStudent = async (student) => {
  selectedStudent.value = student
  showStudentDropdown.value = false
  academicYearFilter.value = ''
  classFilter.value = ''
  page.value = 1
  await Promise.all([
    fetchClassHistory(student.id),
    fetchPaymentHistory()
  ])
}

const resetFilters = () => {
  search.value = ''
  statusFilter.value = ''
  academicYearFilter.value = ''
  classFilter.value = ''
  sortFilter.value = ''
  page.value = 1
}

const paymentAcademicYears = (payment) => {
  const years = new Set()
  ;(payment?.details || []).forEach(detail => {
    if (detail.academic_year) years.add(String(detail.academic_year))
  })
  return [...years]
}

const availableAcademicYears = computed(() => {
  const years = new Set()
  allBills.value
    .filter(bill => !selectedStudent.value || bill.student_id === selectedStudent.value.id)
    .forEach(bill => {
      if (bill.academic_year) years.add(String(bill.academic_year))
    })
  paymentHistory.value.forEach(payment => {
    paymentAcademicYears(payment).forEach(year => years.add(year))
  })
  return [...years].filter(Boolean).sort().reverse()
})

const classHistoryLabel = (item) => {
  const name = item?.class_name || 'Kelas'
  const year = item?.academic_year ? ` - ${item.academic_year}` : ''
  return `${name}${year}`
}

const availableClassFilters = computed(() => {
  return (classHistory.value || [])
    .filter(item => item?.id)
    .map(item => ({
      id: String(item.id),
      label: classHistoryLabel(item),
      academic_year: String(item.academic_year || '')
    }))
})

const selectedClassAcademicYear = computed(() => {
  const selected = availableClassFilters.value.find(item => item.id === classFilter.value)
  return selected?.academic_year || ''
})

const currentStudentClassText = computed(() => {
  const className = selectedStudent.value?.class_name
  const majorName = selectedStudent.value?.major_name
  if (className && majorName) return `${className} - ${majorName}`
  return className || majorName || 'Kelas belum diatur'
})

const classLabelForAcademicYear = (year) => {
  if (!year) return ''
  const item = classHistory.value.find(row => String(row.academic_year || '') === String(year))
  return item ? classHistoryLabel(item) : ''
}

const isHistoryView = computed(() => route.path.includes('/history'))

// Computed Bills
const filteredBills = computed(() => {
  let filtered = allBills.value

  // Filter by selected student
  if (selectedStudent.value) {
    filtered = filtered.filter(b => b.student_id === selectedStudent.value.id)
  }

  // Search filter
  if (search.value) {
    const s = search.value.toLowerCase()
    filtered = filtered.filter(b => 
      (b.bill_type_name && b.bill_type_name.toLowerCase().includes(s)) ||
      (b.name && b.name.toLowerCase().includes(s)) ||
      (b.period && b.period.toLowerCase().includes(s)) ||
      (b.amount && b.amount.toString().includes(s))
    )
  }

  if (statusFilter.value === 'overdue') {
    filtered = filtered.filter(b => b.status !== 'paid' && isOverdue(b))
  } else if (statusFilter.value === 'outstanding') {
    filtered = filtered.filter(b => b.status !== 'paid')
  } else if (statusFilter.value) {
    filtered = filtered.filter(b => b.status === statusFilter.value)
  }

  if (academicYearFilter.value) {
    filtered = filtered.filter(b => String(b.academic_year || '') === academicYearFilter.value)
  }

  if (selectedClassAcademicYear.value) {
    filtered = filtered.filter(b => String(b.academic_year || '') === selectedClassAcademicYear.value)
  }

  return sortBills(filtered)
})

const paymentDetailNames = (payment) => {
  const details = payment?.details || []
  if (details.length === 0) return []
  return details.map(detail => {
    const name = detail.bill_type_name || detail.bill_name || 'Tagihan'
    const period = detail.period ? ` ${formatPeriod({ period: detail.period })}` : ''
    return `${name}${period}`
  })
}

const dateValue = (value) => {
  if (!value) return 0
  const parsed = new Date(value).getTime()
  return Number.isFinite(parsed) ? parsed : 0
}

const textCompare = (a, b) => String(a || '').localeCompare(String(b || ''), 'id-ID', { sensitivity: 'base' })

const sortBills = (rows) => {
  const sorted = [...rows]
  const option = sortFilter.value || 'due_asc'

  sorted.sort((a, b) => {
    switch (option) {
      case 'due_desc':
        return dateValue(b.due_date) - dateValue(a.due_date)
      case 'amount_desc':
        return remainingAmount(b) - remainingAmount(a)
      case 'amount_asc':
        return remainingAmount(a) - remainingAmount(b)
      case 'name_asc':
        return textCompare(a.name || a.bill_type_name, b.name || b.bill_type_name)
      case 'status_asc':
        return textCompare(a.status, b.status)
      case 'due_asc':
      default:
        return dateValue(a.due_date) - dateValue(b.due_date)
    }
  })

  return sorted
}

const sortPayments = (rows) => {
  const sorted = [...rows]
  const option = sortFilter.value || 'date_desc'

  sorted.sort((a, b) => {
    switch (option) {
      case 'date_asc':
        return dateValue(a.paid_at || a.created_at) - dateValue(b.paid_at || b.created_at)
      case 'amount_desc':
        return Number(b.amount || 0) - Number(a.amount || 0)
      case 'amount_asc':
        return Number(a.amount || 0) - Number(b.amount || 0)
      case 'method_asc':
        return textCompare(a.method || a.channel, b.method || b.channel)
      case 'date_desc':
      default:
        return dateValue(b.paid_at || b.created_at) - dateValue(a.paid_at || a.created_at)
    }
  })

  return sorted
}

const sortOptions = computed(() => {
  if (isHistoryView.value) {
    return [
      { value: 'date_desc', label: 'Terbaru' },
      { value: 'date_asc', label: 'Terlama' },
      { value: 'amount_desc', label: 'Nominal Besar' },
      { value: 'amount_asc', label: 'Nominal Kecil' },
      { value: 'method_asc', label: 'Metode A-Z' }
    ]
  }

  return [
    { value: 'due_asc', label: 'Jatuh Tempo Dekat' },
    { value: 'due_desc', label: 'Jatuh Tempo Jauh' },
    { value: 'amount_desc', label: 'Sisa Besar' },
    { value: 'amount_asc', label: 'Sisa Kecil' },
    { value: 'name_asc', label: 'Tagihan A-Z' },
    { value: 'status_asc', label: 'Status A-Z' }
  ]
})

const filteredPayments = computed(() => {
  let filtered = paymentHistory.value

  if (search.value) {
    const s = search.value.toLowerCase()
    filtered = filtered.filter(payment => {
      const detailText = paymentDetailNames(payment).join(' ').toLowerCase()
      return (
        String(payment.id || '').includes(s) ||
        String(payment.transaction_ref || '').toLowerCase().includes(s) ||
        String(payment.method || '').toLowerCase().includes(s) ||
        String(payment.amount || '').includes(s) ||
        detailText.includes(s)
      )
    })
  }

  if (academicYearFilter.value) {
    filtered = filtered.filter(payment => paymentAcademicYears(payment).includes(academicYearFilter.value))
  }

  if (selectedClassAcademicYear.value) {
    filtered = filtered.filter(payment => paymentAcademicYears(payment).includes(selectedClassAcademicYear.value))
  }

  return sortPayments(filtered)
})

const currentRows = computed(() => isHistoryView.value ? filteredPayments.value : filteredBills.value)

watch(isHistoryView, () => {
  statusFilter.value = ''
  sortFilter.value = ''
  page.value = 1
  if (isHistoryView.value) {
    fetchPaymentHistory()
  }
})

const paginatedBills = computed(() => {
  const start = (page.value - 1) * limit.value
  const end = start + limit.value
  return currentRows.value.slice(start, end)
})

const totalPages = computed(() => Math.ceil(currentRows.value.length / limit.value) || 1)
const totalData = computed(() => currentRows.value.length)

const visiblePages = computed(() => {
  const pages = []
  let startPage = Math.max(1, page.value - 1)
  let endPage = Math.min(totalPages.value, startPage + 2)
  
  if (endPage - startPage < 2) {
    startPage = Math.max(1, endPage - 2)
  }
  
  for (let i = startPage; i <= endPage; i++) {
    pages.push(i)
  }
  return pages
})

watch(search, () => {
  page.value = 1
})

watch(statusFilter, () => {
  page.value = 1
})

watch(academicYearFilter, () => {
  page.value = 1
})

watch(classFilter, () => {
  if (selectedClassAcademicYear.value) {
    academicYearFilter.value = selectedClassAcademicYear.value
  }
  page.value = 1
})

watch(sortFilter, () => {
  page.value = 1
})

watch(limit, () => {
  page.value = 1
})

const isOverdue = (bill) => {
  if (!bill?.due_date) return false
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const due = new Date(bill.due_date)
  due.setHours(0, 0, 0, 0)
  return due < today
}

const remainingAmount = (bill) => Math.max(0, Number(bill?.amount || 0) - Number(bill?.total_paid || 0))
const canUseDeposit = (bill) => Number(selectedStudent.value?.deposit_balance || bill?.deposit_balance || 0) > 0 && remainingAmount(bill) > 0

const selectedSummary = computed(() => {
  if (isHistoryView.value) {
    const payments = filteredPayments.value
    return {
      total: payments.length,
      unpaid: payments.filter(p => p.channel === 'gateway').length,
      overdue: payments.filter(p => p.deposit_applied > 0).length,
      remaining: payments.reduce((acc, payment) => acc + Number(payment.amount || 0), 0)
    }
  }
  const bills = filteredBills.value
  return {
    total: bills.length,
    unpaid: bills.filter(b => b.status !== 'paid').length,
    overdue: bills.filter(b => b.status !== 'paid' && isOverdue(b)).length,
    remaining: bills.reduce((acc, bill) => acc + remainingAmount(bill), 0)
  }
})

const payWithMidtrans = async (bill, useDeposit = false) => {
  if (isOverdue(bill)) {
    showNotification('Tagihan sudah jatuh tempo. Pembayaran online ditutup, silakan bayar langsung ke admin/kasir sekolah.', 'error')
    return
  }

  const remaining = Number(bill.amount || 0) - Number(bill.total_paid || 0)
  if (remaining <= 0) return
  const availableDeposit = useDeposit ? Number(selectedStudent.value?.deposit_balance || bill.deposit_balance || 0) : 0
  const depositApplied = Math.min(availableDeposit, remaining)

  payingBillId.value = bill.id
  try {
    if (useDeposit && depositApplied >= remaining) {
      await axios.post('finance/payments', {
        student_id: bill.student_id,
        amount: remaining,
        deposit_applied: remaining,
        channel: 'deposit',
        method: 'Saldo Deposit',
        bill_ids: [bill.id],
        note: 'Pembayaran tagihan dari saldo deposit parent web'
      })
      showNotification('Tagihan berhasil dibayar memakai saldo deposit.')
      await fetchBills()
      await fetchPaymentHistory()
      return
    }

    const res = await axios.post('finance/payment-intent', {
      student_id: bill.student_id,
      amount: remaining,
      deposit_applied: depositApplied,
      bill_ids: [bill.id],
      is_bypass_rule: false
    })
    const payment = res.data?.data
    if (payment?.snap_token && window.snap) {
      window.snap.pay(payment.snap_token, {
        onSuccess: async () => {
          showNotification('Pembayaran diterima Midtrans. Status akan diverifikasi otomatis oleh sistem.')
          await fetchBills()
          await fetchPaymentHistory()
        },
        onPending: async () => {
          showNotification('Pembayaran masih menunggu penyelesaian.', 'warning')
          await fetchBills()
          await fetchPaymentHistory()
        },
        onError: () => showNotification('Pembayaran gagal diproses oleh Midtrans.', 'error'),
        onClose: () => showNotification('Anda menutup halaman pembayaran.', 'warning')
      })
    } else if (payment?.payment_url) {
      window.location.href = payment.payment_url
    } else {
      showNotification('Link pembayaran belum tersedia. Silakan coba lagi.', 'error')
    }
  } catch (err) {
    showNotification(err.response?.data?.message || 'Gagal membuat pembayaran Midtrans', 'error')
  } finally {
    payingBillId.value = null
  }
}

const printReceipt = (paymentId) => {
  const url = router.resolve({ name: 'receipt-print', params: { id: paymentId } }).href
  window.open(url, '_blank')
}

const openHistory = () => {
  router.push('/parent/history')
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear()}`
}

const formatCurrency = (value) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value || 0)
}

const formatPeriod = (bill) => {
  if (bill?.period_month && bill?.period_year) {
    const date = new Date(Number(bill.period_year), Number(bill.period_month) - 1, 1)
    return date.toLocaleDateString('id-ID', { month: 'long', year: 'numeric' })
  }
  if (bill?.period && /^\d{4}-\d{2}$/.test(bill.period)) {
    const [year, month] = bill.period.split('-')
    const date = new Date(Number(year), Number(month) - 1, 1)
    return date.toLocaleDateString('id-ID', { month: 'long', year: 'numeric' })
  }
  return bill?.period || bill?.academic_year || '-'
}

onMounted(() => {
  isMounted.value = true
  fetchParentData()
})

</script>

<template>
  <div class="min-h-screen bg-slate-50/50 p-4 lg:p-6 space-y-6 animate-fade-in font-inter">

    <!-- Teleport to Header for Global Search and Filters -->
    <Teleport v-if="isMounted" to="#header-actions-target">
      <div class="flex items-center justify-center w-full gap-4 relative mx-auto">
        <div class="flex flex-wrap items-center justify-center gap-2 flex-1 max-w-5xl mx-auto">
          
          <!-- Student Selector Dropdown -->
          <div class="relative w-64 group/student">
            <button @click="showStudentDropdown = !showStudentDropdown" class="w-full relative p-2.5 pl-10 pr-8 bg-white text-slate-700 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm transition-all text-left flex items-center justify-between">
              <UserIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-indigo-500" />
              <span class="text-xs font-bold truncate">
                {{ selectedStudent ? selectedStudent.name : 'Pilih Siswa' }}
              </span>
              <ChevronDownIcon class="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400" :class="{ 'rotate-180': showStudentDropdown }" />
            </button>
            
            <transition name="dropdown">
              <div v-if="showStudentDropdown" class="absolute top-full left-0 mt-2 w-full bg-white rounded-2xl shadow-xl border border-slate-100 z-50 py-2 overflow-hidden">
                <button 
                  v-for="s in students" 
                  :key="s.id"
                  @click="selectStudent(s)"
                  class="w-full text-left px-4 py-3 hover:bg-indigo-50 flex flex-col transition-all"
                  :class="{'bg-indigo-50/50': selectedStudent?.id === s.id}"
                >
                  <span class="text-xs font-bold text-slate-800" :class="{'text-indigo-600': selectedStudent?.id === s.id}">{{ s.name }}</span>
                  <span class="text-[10px] font-medium text-slate-400">
                    {{ s.class_name || 'Kelas belum diatur' }} • NISN: {{ s.nisn || '-' }}
                  </span>
                </button>
              </div>
            </transition>
          </div>

          <div class="relative flex-1 min-w-[180px] group">
            <SearchIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300 group-focus-within:text-indigo-600" />
            <input v-model="search" type="text" :placeholder="isHistoryView ? 'Cari transaksi...' : 'Cari tagihan...'" class="search-input-premium" />
          </div>
          
          <div v-if="availableAcademicYears.length" class="relative">
            <select v-model="academicYearFilter" class="appearance-none p-2.5 pr-8 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm transition-all text-[10px] font-black uppercase tracking-wider cursor-pointer focus:outline-none focus:border-indigo-500">
              <option value="">Semua Tahun</option>
              <option v-for="year in availableAcademicYears" :key="year" :value="year">{{ year }}</option>
            </select>
            <ChevronDownIcon class="absolute right-2 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-slate-400 pointer-events-none" />
          </div>

          <div v-if="availableClassFilters.length" class="relative">
            <select v-model="classFilter" class="appearance-none p-2.5 pr-8 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm transition-all text-[10px] font-black uppercase tracking-wider cursor-pointer focus:outline-none focus:border-indigo-500 max-w-[170px]">
              <option value="">Semua Kelas</option>
              <option v-for="item in availableClassFilters" :key="item.id" :value="item.id">{{ item.label }}</option>
            </select>
            <ChevronDownIcon class="absolute right-2 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-slate-400 pointer-events-none" />
          </div>

          <div v-if="!isHistoryView" class="relative">
            <FilterIcon class="absolute left-2.5 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-slate-400 pointer-events-none" />
            <select v-model="statusFilter" class="appearance-none p-2.5 pl-8 pr-8 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm transition-all text-[10px] font-black uppercase tracking-wider cursor-pointer focus:outline-none focus:border-indigo-500">
              <option value="">Semua Status</option>
              <option value="outstanding">Ada Tunggakan</option>
              <option value="overdue">Lewat Tempo</option>
              <option value="unpaid">Belum Lunas</option>
              <option value="partial">Sebagian</option>
              <option value="paid">Lunas</option>
            </select>
            <ChevronDownIcon class="absolute right-2 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-slate-400 pointer-events-none" />
          </div>

          <div class="relative">
            <SortIcon class="absolute left-2.5 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-slate-400 pointer-events-none" />
            <select v-model="sortFilter" class="appearance-none p-2.5 pl-8 pr-8 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm transition-all text-[10px] font-black uppercase tracking-wider cursor-pointer focus:outline-none focus:border-indigo-500">
              <option value="">Urutan Default</option>
              <option v-for="item in sortOptions" :key="item.value" :value="item.value">{{ item.label }}</option>
            </select>
            <ChevronDownIcon class="absolute right-2 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-slate-400 pointer-events-none" />
          </div>

          <button @click="resetFilters" class="p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm transition-all group" title="Reset Filter">
            <ResetIcon class="w-4 h-4 group-hover:rotate-180 transition-transform duration-500" />
          </button>
        </div>
      </div>
    </Teleport>

    <!-- Welcome Section (Optional, smaller since Header has context) -->
    <div class="flex items-center justify-between mb-2">
      <div>
        <h1 class="text-2xl font-black text-slate-800 tracking-tight">Halo, {{ authStore.user?.name }}</h1>
        <p class="text-[11px] text-slate-500 font-bold mt-1">{{ isHistoryView ? 'Lihat riwayat tagihan yang sudah selesai dibayar.' : 'Pantau dan kelola tagihan pembayaran putra-putri Anda.' }}</p>
      </div>
      <div v-if="selectedStudent" class="bg-indigo-50 px-4 py-2.5 rounded-xl border border-indigo-100 flex items-center gap-3">
         <div class="w-8 h-8 rounded-lg bg-indigo-600 text-white flex items-center justify-center shrink-0">
           <UserIcon class="w-4 h-4" />
         </div>
         <div>
           <p class="text-[9px] font-black uppercase tracking-widest text-indigo-400">Siswa Terpilih</p>
           <p class="text-xs font-black text-indigo-700 leading-tight">{{ selectedStudent.name }}</p>
           <p class="text-[9px] font-bold text-indigo-500 uppercase tracking-widest mt-0.5">{{ currentStudentClassText }}</p>
           <p class="text-[9px] font-bold text-emerald-600 uppercase tracking-widest mt-0.5">Saldo: {{ formatCurrency(selectedStudent.deposit_balance || 0) }}</p>
         </div>
      </div>
    </div>

    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <div v-for="card in [
        { label: isHistoryView ? 'Total Transaksi' : 'Total Tagihan', value: selectedSummary.total },
        { label: isHistoryView ? 'Via Midtrans' : 'Ada Tunggakan', value: selectedSummary.unpaid },
        { label: isHistoryView ? 'Pakai Saldo' : 'Lewat Tempo', value: selectedSummary.overdue },
        { label: isHistoryView ? 'Total Dibayar' : 'Sisa Bayar', value: formatCurrency(selectedSummary.remaining), wide: true }
      ]" :key="card.label" class="bg-white border border-slate-200 rounded-2xl p-5 shadow-sm">
        <p class="text-[9px] font-black text-slate-400 uppercase tracking-[0.22em]">{{ card.label }}</p>
        <p class="mt-2 text-xl font-black text-slate-800">{{ card.value }}</p>
      </div>
    </div>

    <!-- Main Content Table -->
    <div class="bg-white rounded-2xl border border-slate-200 shadow-sm flex flex-col min-h-[710px] overflow-hidden">
      <div class="px-6 py-6 border-b border-slate-100 bg-slate-50/30 flex items-center justify-between">
        <div class="flex items-center gap-4">
          <div class="w-2 h-6 bg-indigo-500 rounded-full"></div>
          <h3 class="font-black text-slate-700 text-sm uppercase tracking-[0.2em]">{{ isHistoryView ? 'Riwayat Pembayaran Siswa' : 'Daftar Tagihan Siswa' }}</h3>
        </div>
        <div class="flex items-center gap-3">
           <router-link :to="isHistoryView ? '/parent/bills' : '/parent/history'" class="bg-indigo-50 text-indigo-600 hover:bg-indigo-100 font-bold py-2 px-4 rounded-xl text-[10px] flex items-center gap-2 transition-all shadow-sm">
             <ReceiptIcon class="w-3.5 h-3.5" />
             <span>{{ isHistoryView ? 'Lihat Tagihan' : 'Riwayat Pembayaran' }}</span>
           </router-link>
        </div>
      </div>

      <div class="flex-1 overflow-x-auto custom-scrollbar relative">
        <!-- Loading Overlay -->
        <div v-if="loading" class="absolute inset-0 bg-white/60 backdrop-blur-sm z-10 flex items-center justify-center">
          <div class="flex flex-col items-center gap-3">
            <div class="w-10 h-10 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin"></div>
            <p class="text-xs font-black text-indigo-600 uppercase tracking-widest">Memuat Data...</p>
          </div>
        </div>

        <table class="premium-table w-full">
          <thead>
            <tr v-if="!isHistoryView">
              <th class="w-16">No</th>
              <th>Detail Tagihan</th>
              <th>Periode</th>
              <th>Jatuh Tempo</th>
              <th>Nominal</th>
              <th>Status</th>
              <th class="text-center w-32">Aksi</th>
            </tr>
            <tr v-else>
              <th class="w-16">No</th>
              <th>Detail Pembayaran</th>
              <th>Tanggal</th>
              <th>Metode</th>
              <th>Nominal</th>
              <th>Status</th>
              <th class="text-center w-32">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!loading && paginatedBills.length === 0">
              <td colspan="7" class="py-20 text-center">
                <div class="w-24 h-24 bg-slate-50 rounded-full flex items-center justify-center mx-auto mb-4 border border-slate-100">
                  <ReceiptIcon class="w-10 h-10 text-slate-300" />
                </div>
                <h3 class="text-sm font-black text-slate-700 uppercase tracking-wider mb-1">{{ isHistoryView ? 'Belum Ada Riwayat' : 'Tidak Ada Tagihan' }}</h3>
                <p class="text-xs text-slate-400 font-medium">{{ isHistoryView ? 'Belum ada pembayaran lunas yang sesuai dengan filter.' : 'Belum ada data tagihan yang sesuai dengan pencarian Anda.' }}</p>
              </td>
            </tr>
            <tr v-for="(b, i) in paginatedBills" :key="`${isHistoryView ? 'payment' : 'bill'}-${b.id}`" class="hover:bg-slate-50/50 transition-colors group">
              <td class="text-center text-xs font-bold text-slate-400">
                {{ (page - 1) * limit + i + 1 }}
              </td>
              <template v-if="!isHistoryView">
              <td>
                <div>
                  <p class="text-xs font-black text-slate-800 group-hover:text-indigo-600 transition-colors leading-tight">
                    {{ b.name || b.bill_type_name }}
                  </p>
                  <p class="text-[10px] font-bold text-slate-400 mt-1 truncate max-w-[200px]">
                    {{ b.description || 'Tidak ada deskripsi' }}
                  </p>
                  <div class="mt-2 flex flex-wrap gap-1.5">
                    <span v-if="b.academic_year" class="px-2 py-1 bg-indigo-50 text-indigo-600 rounded-lg text-[8px] font-black uppercase tracking-wider">
                      {{ b.academic_year }}
                    </span>
                    <span v-if="classLabelForAcademicYear(b.academic_year)" class="px-2 py-1 bg-slate-100 text-slate-500 rounded-lg text-[8px] font-black uppercase tracking-wider max-w-[180px] truncate">
                      {{ classLabelForAcademicYear(b.academic_year) }}
                    </span>
                  </div>
                </div>
              </td>
              <td>
                <span class="px-2.5 py-1 bg-slate-100 text-slate-600 rounded-lg text-[10px] font-black uppercase tracking-wider">
                  {{ formatPeriod(b) }}
                </span>
              </td>
              <td>
                <div class="flex items-center gap-1.5">
                  <CalendarIcon class="w-3.5 h-3.5" :class="isOverdue(b) && b.status !== 'paid' ? 'text-amber-500' : 'text-slate-400'" />
                  <span class="text-xs font-bold" :class="isOverdue(b) && b.status !== 'paid' ? 'text-amber-600' : 'text-slate-600'">
                    {{ formatDate(b.due_date) }}
                  </span>
                </div>
              </td>
              <td>
                <div class="space-y-1">
                  <p class="text-xs font-black text-slate-800">{{ formatCurrency(b.amount) }}</p>
                  <p v-if="b.total_paid > 0 && b.status !== 'paid'" class="text-[9px] font-black text-emerald-600 uppercase tracking-widest">
                    Dibayar: {{ formatCurrency(b.total_paid) }}
                  </p>
                  <p v-if="b.status !== 'paid'" class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">
                    Sisa: {{ formatCurrency(remainingAmount(b)) }}
                  </p>
                </div>
              </td>
              <td>
                <div class="flex flex-col gap-1 items-start">
                  <span v-if="b.status === 'paid'" class="px-2.5 py-1 bg-emerald-50 text-emerald-600 border border-emerald-100/50 rounded-lg text-[9px] font-black uppercase tracking-widest flex items-center gap-1.5 shadow-sm shadow-emerald-500/10">
                    <div class="w-1.5 h-1.5 rounded-full bg-emerald-500"></div>Lunas
                  </span>
                  <span v-else-if="b.status === 'partial'" class="px-2.5 py-1 bg-amber-50 text-amber-600 border border-amber-100/50 rounded-lg text-[9px] font-black uppercase tracking-widest flex items-center gap-1.5 shadow-sm shadow-amber-500/10">
                    <div class="w-1.5 h-1.5 rounded-full bg-amber-500"></div>Sebagian
                  </span>
                  <span v-else class="px-2.5 py-1 bg-rose-50 text-rose-600 border border-rose-100/50 rounded-lg text-[9px] font-black uppercase tracking-widest flex items-center gap-1.5 shadow-sm shadow-rose-500/10">
                    <div class="w-1.5 h-1.5 rounded-full bg-rose-500"></div>Belum Lunas
                  </span>
                  <span v-if="isOverdue(b) && b.status !== 'paid'" class="text-[8px] font-black text-amber-600 uppercase tracking-widest flex items-center gap-1 mt-0.5 ml-1">
                    <AlertIcon class="w-2.5 h-2.5" /> Jatuh Tempo
                  </span>
                </div>
              </td>
              <td class="text-center">
                <div v-if="b.status !== 'paid'" class="space-y-2">
                 <button
                  @click="payWithMidtrans(b)"
                  :disabled="payingBillId === b.id"
                  :title="isOverdue(b) ? 'Tagihan jatuh tempo harus dibayar ke kasir sekolah' : 'Bayar via Midtrans'"
                  :class="[
                    'w-full font-black py-2.5 px-3 rounded-xl text-[10px] uppercase tracking-widest transition-all disabled:opacity-60 flex items-center justify-center gap-1.5 shadow-sm',
                    isOverdue(b)
                      ? 'bg-slate-50 hover:bg-slate-100 text-slate-500 border border-slate-200 cursor-not-allowed'
                      : 'bg-indigo-600 hover:bg-indigo-700 text-white hover:shadow-md hover:shadow-indigo-500/20'
                  ]"
                >
                  <BillIcon class="w-3.5 h-3.5" />
                  {{ payingBillId === b.id ? 'Loading...' : (isOverdue(b) ? 'Ke Kasir' : 'Bayar') }}
                </button>
                <button
                  v-if="canUseDeposit(b) && !isOverdue(b)"
                  @click="payWithMidtrans(b, true)"
                  :disabled="payingBillId === b.id"
                  class="w-full bg-emerald-50 hover:bg-emerald-100 text-emerald-700 border border-emerald-100 font-black py-2.5 px-3 rounded-xl text-[10px] uppercase tracking-widest transition-all flex items-center justify-center gap-1.5 shadow-sm disabled:opacity-60"
                  :title="selectedStudent.deposit_balance >= remainingAmount(b) ? 'Bayar penuh memakai saldo deposit' : 'Pakai saldo deposit sebagai potongan Midtrans'"
                >
                  <BillIcon class="w-3.5 h-3.5" />
                  {{ selectedStudent.deposit_balance >= remainingAmount(b) ? 'Pakai Saldo' : 'Saldo + Bayar' }}
                </button>
                </div>
                <button v-else @click="openHistory" class="w-full bg-white border border-slate-200 text-slate-500 hover:bg-slate-50 hover:text-indigo-600 font-black py-2.5 px-3 rounded-xl text-[10px] uppercase tracking-widest transition-all flex items-center justify-center gap-1.5 shadow-sm">
                  <ReceiptIcon class="w-3.5 h-3.5" />
                  Kwitansi
                </button>
              </td>
              </template>
              <template v-else>
                <td>
                  <div>
                    <p class="text-xs font-black text-slate-800 group-hover:text-indigo-600 transition-colors leading-tight">
                      Pembayaran #{{ b.id }}
                    </p>
                    <div class="mt-1 flex flex-wrap gap-1.5">
                      <span v-for="name in paymentDetailNames(b).slice(0, 2)" :key="name" class="max-w-[180px] truncate px-2 py-1 bg-slate-100 rounded-lg text-[8px] font-black text-slate-500 uppercase tracking-wider">
                        {{ name }}
                      </span>
                      <span v-if="paymentDetailNames(b).length > 2" class="px-2 py-1 bg-indigo-50 rounded-lg text-[8px] font-black text-indigo-600 uppercase tracking-wider">
                        +{{ paymentDetailNames(b).length - 2 }}
                      </span>
                      <span v-if="paymentDetailNames(b).length === 0" class="px-2 py-1 bg-slate-100 rounded-lg text-[8px] font-black text-slate-500 uppercase tracking-wider">
                        Alokasi otomatis
                      </span>
                    </div>
                    <div v-if="paymentAcademicYears(b).length" class="mt-1 flex flex-wrap gap-1.5">
                      <span v-for="year in paymentAcademicYears(b).slice(0, 2)" :key="year" class="px-2 py-1 bg-indigo-50 text-indigo-600 rounded-lg text-[8px] font-black uppercase tracking-wider">
                        {{ classLabelForAcademicYear(year) || year }}
                      </span>
                    </div>
                  </div>
                </td>
                <td>
                  <div class="flex items-center gap-1.5">
                    <CalendarIcon class="w-3.5 h-3.5 text-slate-400" />
                    <span class="text-xs font-bold text-slate-600">{{ formatDate(b.paid_at || b.created_at) }}</span>
                  </div>
                </td>
                <td>
                  <span class="px-2.5 py-1 bg-slate-100 text-slate-600 rounded-lg text-[10px] font-black uppercase tracking-wider">
                    {{ b.method || b.channel || '-' }}
                  </span>
                </td>
                <td>
                  <div class="space-y-1">
                    <p class="text-xs font-black text-slate-800">{{ formatCurrency(b.amount) }}</p>
                    <p v-if="b.deposit_applied > 0" class="text-[9px] font-black text-emerald-600 uppercase tracking-widest">
                      Saldo: {{ formatCurrency(b.deposit_applied) }}
                    </p>
                  </div>
                </td>
                <td>
                  <span class="px-2.5 py-1 bg-emerald-50 text-emerald-600 border border-emerald-100/50 rounded-lg text-[9px] font-black uppercase tracking-widest flex items-center gap-1.5 shadow-sm shadow-emerald-500/10 w-fit">
                    <div class="w-1.5 h-1.5 rounded-full bg-emerald-500"></div>Lunas
                  </span>
                </td>
                <td class="text-center">
                  <button
                    @click="printReceipt(b.id)"
                    class="w-full bg-white border border-slate-200 text-slate-500 hover:bg-slate-50 hover:text-indigo-600 font-black py-2.5 px-3 rounded-xl text-[10px] uppercase tracking-widest transition-all flex items-center justify-center gap-1.5 shadow-sm"
                  >
                    <ReceiptIcon class="w-3.5 h-3.5" />
                    Kwitansi
                  </button>
                </td>
              </template>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="px-8 py-6 bg-slate-50/50 border-t border-slate-100 flex items-center justify-between">
        <div class="flex items-center gap-6">
          <div class="flex items-center gap-3">
            <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Tampilkan</span>
            <select v-model="limit" class="bg-white border border-slate-200 rounded-lg text-[10px] font-black text-slate-600 px-2 py-1 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all cursor-pointer shadow-sm">
              <option :value="10">10</option>
              <option :value="25">25</option>
              <option :value="50">50</option>
            </select>
          </div>
          <div class="h-8 w-px bg-slate-200 hidden sm:block"></div>
          <span class="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em]">
            Halaman <span class="text-indigo-600">{{ page }}</span> dari {{ totalPages }} <span class="mx-2 text-slate-300">|</span> Total <span class="text-indigo-600">{{ totalData }}</span> Data
          </span>
        </div>
        
        <!-- Pagination Control -->
        <div class="flex items-center gap-2">
          <!-- Previous Button -->
          <button 
            v-if="totalPages > 1"
            @click="page > 1 && (page--)" 
            :disabled="page <= 1" 
            class="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-400 hover:text-indigo-600 hover:border-indigo-100 hover:bg-indigo-50/30 disabled:opacity-20 disabled:hover:bg-white disabled:hover:border-slate-200 transition-all cursor-pointer"
          >
            <PrevIcon class="w-4 h-4" />
          </button>

          <!-- Page Numbers -->
          <div class="flex items-center gap-1">
            <button 
              v-for="p in visiblePages" 
              :key="p"
              @click="page = p"
              class="w-10 h-10 flex items-center justify-center rounded-xl text-[10px] font-black transition-all cursor-pointer"
              :class="p === page 
                ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/20' 
                : 'bg-white border border-slate-200 text-slate-500 hover:bg-slate-50 hover:border-slate-300'"
            >
              {{ p }}
            </button>
          </div>

          <!-- Next Button -->
          <button 
            v-if="totalPages > 1"
            @click="page < totalPages && (page++)" 
            :disabled="page >= totalPages" 
            class="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-400 hover:text-indigo-600 hover:border-indigo-100 hover:bg-indigo-50/30 disabled:opacity-20 disabled:hover:bg-white disabled:hover:border-slate-200 transition-all cursor-pointer"
          >
            <NextIcon class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>

    <Teleport to="body">
      <transition name="fade">
        <div
          v-if="notification.show"
          class="fixed top-6 right-6 z-[3000] px-5 py-4 rounded-2xl shadow-2xl flex items-center gap-3 max-w-md"
          :class="notification.type === 'error' ? 'bg-rose-600 text-white' : notification.type === 'warning' ? 'bg-amber-500 text-white' : 'bg-emerald-600 text-white'"
        >
          <AlertIcon v-if="notification.type === 'error' || notification.type === 'warning'" class="w-5 h-5 shrink-0" />
          <PaidIcon v-else class="w-5 h-5 shrink-0" />
          <span class="text-xs font-bold leading-relaxed">{{ notification.message }}</span>
        </div>
      </transition>
    </Teleport>

    <Teleport to="body">
      <button
        @click="contactWhatsApp"
        class="fixed bottom-6 right-6 z-[2200] w-14 h-14 rounded-2xl bg-emerald-500 text-white shadow-2xl shadow-emerald-300/60 flex items-center justify-center hover:bg-emerald-600 transition-all hover:scale-105 active:scale-95 group hover:-translate-y-1"
        title="Hubungi CS Admin via WhatsApp"
      >
        <ChatIcon class="w-6 h-6 group-hover:rotate-12 transition-transform" />
      </button>
    </Teleport>
  </div>
</template>

<style scoped>
.dropdown-enter-active, .dropdown-leave-active {
  transition: opacity 0.2s, transform 0.2s;
}
.dropdown-enter-from, .dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
