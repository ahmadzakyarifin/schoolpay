<script setup>
import { ref, onMounted, reactive, watch, computed } from 'vue'
import axios from 'axios'
import { useRoute } from 'vue-router'

import { 
  Users as UsersIcon, 
  GraduationCap as StudentIcon, 
  Receipt as BillIcon, 
  Wallet as WalletIcon,
  RefreshCw as RefreshIcon,
  LayoutGrid as LayoutIcon,
  Clock as ClockIcon,
  Calendar as CalendarIcon,
  Search as SearchIcon,
  ArrowRight as ArrowIcon,
  FileSpreadsheet as ExcelIcon,
  Printer as PrintIcon,
  MessageCircle as WAIcon,
  Mail as MailIcon,
  CheckCircle2 as CheckIcon,
  XCircle as ErrorIcon,
  AlertTriangle as AlertIcon,
  FileText as PDFIcon,
  Filter as FilterIcon,
  ChevronLeft as PrevIcon,
  ChevronRight as NextIcon,
  RotateCcw as ResetIcon,
  X as XIcon
} from 'lucide-vue-next'
import ReportsFilterModal from '../../components/reports/ReportsFilterModal.vue'
import { useAuthStore } from '../../store/auth'

const authStore = useAuthStore()
const route = useRoute()
const isOffline = computed(() => (typeof navigator !== 'undefined' && navigator.onLine === false))

const isMounted = ref(false)
const showFilters = ref(false)
const search = ref('')
const academicYears = ref([])
const loading = ref(false)
const previewData = ref([])
const previewLoading = ref(false)
const classes = ref([])
const majors = ref([])
const billTypes = ref([])

const page = ref(1)
const limit = ref(10)
const totalData = ref(0)
const activeTab = ref(route.query.tab === 'arrears' ? 'arrears' : 'payments')
const showPaymentTypeModal = ref(false)
const selectedPayment = ref(null)

const summary = ref({
  total_paid_amount: 0,
  total_unpaid_amount: 0,
  paid_count: 0,
  unpaid_count: 0,
  total_bills: 0
})

const filters = reactive({
  period: 'all',
  academic_year_id: '',
  class_id: '',
  major_id: '',
  bill_type_id: '',
  ref_date: new Date().toISOString().substring(0, 10),
  start_date: '',
  end_date: ''
})

const tempFilters = reactive({
  period: 'all',
  academic_year_id: '',
  class_id: '',
  major_id: '',
  bill_type_id: '',
  ref_date: new Date().toISOString().substring(0, 10),
  start_date: '',
  end_date: ''
})

watch(() => tempFilters.period, (newP) => {
  const now = new Date()
  if (newP === 'daily') {
    tempFilters.ref_date = now.toISOString().substring(0, 10)
  } else if (newP === 'monthly') {
    tempFilters.ref_date = now.toISOString().substring(0, 7)
  } else if (newP === 'yearly') {
    tempFilters.ref_date = now.getFullYear().toString()
  } else if (newP === 'custom') {
    tempFilters.start_date = now.toISOString().substring(0, 10)
    tempFilters.end_date = now.toISOString().substring(0, 10)
  }
})

const resetFilters = () => {
  search.value = ''
  const defaults = {
    period: 'all',
    academic_year_id: '',
    class_id: '',
    major_id: '',
    bill_type_id: '',
    ref_date: new Date().toISOString().substring(0, 10),
    start_date: '',
    end_date: ''
  }
  Object.assign(tempFilters, defaults)
  Object.assign(filters, defaults)
  page.value = 1
  showFilters.value = false
  fetchPreview()
}

const applyFilters = () => {
  Object.assign(filters, tempFilters)
  page.value = 1
  showFilters.value = false
  fetchPreview()
}

const fetchAcademicYears = async () => {
  try {
    const res = await axios.get('academic/years', { params: { limit: 100 } })
    academicYears.value = res.data.data.data || []
  } catch (err) {}
}

const fetchClassesMajors = async () => {
  try {
    const resCls = await axios.get('academic/class', { params: { limit: 100 } })
    const pCls = resCls.data?.data
    classes.value = (pCls?.data || (Array.isArray(pCls) ? pCls : [])).filter(c => c.is_active !== false)
  } catch (err) { console.error('Gagal fetch classes:', err) }

  try {
    const resMaj = await axios.get('academic/major', { params: { limit: 100 } })
    const pMaj = resMaj.data?.data
    majors.value = (pMaj?.data || (Array.isArray(pMaj) ? pMaj : [])).filter(m => m.is_active !== false)
  } catch (err) { console.error('Gagal fetch majors:', err) }

  try {
    const resBT = await axios.get('finance/bill-types', { params: { limit: 100 } })
    const pBT = resBT.data?.data
    billTypes.value = (pBT?.data || (Array.isArray(pBT) ? pBT : [])).filter(bt => bt.is_active !== false)
  } catch (err) { console.error('Gagal fetch bill types:', err) }
}

const fetchPreview = async () => {
  previewLoading.value = true
  try {
    let endpoint = 'dashboard/stats'
    if (activeTab.value === 'arrears') endpoint = 'finance/arrears'
    
    const res = await axios.get(endpoint, { 
      params: { ...filters, search: search.value, page: page.value, limit: limit.value } 
    })
    
    if (activeTab.value === 'payments') {
      previewData.value = res.data.data.recent_payments || []
      totalData.value = res.data.data.total_payments_count || res.data.data.recent_payments?.length || 0
      summary.value = res.data.data.stats
    } else if (activeTab.value === 'arrears') {
      previewData.value = res.data.data.data || []
      totalData.value = res.data.data.total || 0
    }
  } catch (err) {} finally {
    previewLoading.value = false
  }
}

const downloadReport = async (format, type) => {
  if (isOffline.value) return
  loading.value = true
  try {
    let endpoint = type === 'students' ? 'students/export' : 'dashboard/export'
    const res = await axios.get(endpoint, {
      params: { ...filters, search: search.value, tab: activeTab.value, format },
      responseType: 'blob'
    })
    const url = window.URL.createObjectURL(new Blob([res.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `Laporan_${activeTab.value}_${new Date().toISOString().split('T')[0]}.${format === 'pdf' ? 'pdf' : 'xlsx'}`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (err) {} finally {
    loading.value = false
  }
}

const displayBillTypes = (value) => {
  if (!value) return 'Pembayaran Deposit / Kustom'
  if (typeof value === 'string') return value.split('||').filter(Boolean).join(', ') || 'Pembayaran Deposit / Kustom'
  if (value.String !== undefined) return value.Valid === false ? 'Pembayaran Deposit / Kustom' : value.String
  if (Array.isArray(value)) return value.join(', ')
  return String(value)
}

const billTypeItems = (value) => {
  const text = displayBillTypes(value)
  if (!text || text === 'Pembayaran Deposit / Kustom') return [text]
  return text.split(/\|\||,/).map(item => item.trim()).filter(Boolean)
}

const paymentDetailItems = (payment) => {
  const rawDetails = payment?.bill_type_details
  if (typeof rawDetails === 'string' && rawDetails.trim()) {
    return rawDetails.split('||').map((item) => {
      const [name, period, amount] = item.split('::')
      return {
        name: name || 'Pembayaran Deposit / Kustom',
        period: period || '',
        amount: Number(String(amount || '0').replace(/[^0-9.-]/g, '')) || 0
      }
    }).filter(item => item.name)
  }

  return billTypeItems(payment?.bill_type_names).map(name => ({ name, period: '', amount: 0 }))
}

const paymentItemCount = (payment) => {
  return Number(payment?.bill_item_count || paymentDetailItems(payment).length || 0)
}

const paymentSummaryTitle = (payment) => {
  const items = paymentDetailItems(payment)
  if (items.length === 0) return 'Pembayaran Deposit / Kustom'
  if (items.length === 1) return items[0].name
  return `${items[0].name} + ${items.length - 1} tagihan`
}

const openPaymentTypeModal = (item) => {
  selectedPayment.value = item
  showPaymentTypeModal.value = true
}

const formatCurrency = (val) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency', currency: 'IDR', minimumFractionDigits: 0
  }).format(val || 0)
}

const formatDate = (dateStr) => {
  if (!dateStr || String(dateStr).startsWith('0001-01-01')) return '-'
  const d = new Date(dateStr)
  return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear()}`
}

const formatDateTime = (dateStr) => {
  if (!dateStr || String(dateStr).startsWith('0001-01-01')) return '-'
  const d = new Date(dateStr)
  return d.toLocaleString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatPaymentMethod = (method) => {
  if (!method) return '-'
  return String(method).replace(/_/g, ' ').toUpperCase()
}

const currentPageArrearsTotal = computed(() => {
  return previewData.value.reduce((acc, item) => acc + Math.max(0, Number(item.amount || 0) - Number(item.total_paid || 0)), 0)
})

const summaryCards = computed(() => {
  if (activeTab.value === 'arrears') {
    return [
      { label: 'Sisa Tunggakan', value: formatCurrency(currentPageArrearsTotal.value), helper: 'nominal dari data pada halaman ini' },
      { label: 'Lewat Jatuh Tempo', value: previewData.value.filter(item => getRemainingDaysText(item.due_date).startsWith('Telat')).length, helper: 'tagihan perlu ditindaklanjuti' }
    ]
  }

  return [
    {
      label: 'Total Pembayaran',
      value: formatCurrency(summary.value?.payments?.total || summary.value?.paid_amount || 0),
      helper: `${summary.value?.paid_count || 0} transaksi berhasil`
    },
    {
      label: 'Sisa Tunggakan',
      value: formatCurrency(summary.value?.unpaid_amount || 0),
      helper: `${summary.value?.unpaid_count || 0} tagihan belum lunas`
    }
  ]
})

const getRemainingDaysText = (dateStr) => {
  if (!dateStr || String(dateStr).startsWith('0001-01-01')) return ''
  const match = dateStr.match(/^(\d{4})-(\d{2})-(\d{2})/)
  let d
  if (match) {
    d = new Date(parseInt(match[1]), parseInt(match[2]) - 1, parseInt(match[3]))
  } else {
    d = new Date(dateStr)
  }
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  d.setHours(0, 0, 0, 0)
  
  const diffTime = d.getTime() - today.getTime()
  const diffDays = Math.round(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 0) {
    return 'Hari Ini'
  } else if (diffDays > 0) {
    return `H-${diffDays}`
  } else {
    return `Telat ${Math.abs(diffDays)} Hari`
  }
}

const totalPages = computed(() => Math.ceil(totalData.value / limit.value) || 1)

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
  if (page.value === 1) {
    fetchPreview()
  } else {
    page.value = 1
  }
})

watch(page, () => {
  fetchPreview()
})

watch(limit, () => {
  if (page.value === 1) {
    fetchPreview()
  } else {
    page.value = 1
  }
})

watch(activeTab, () => {
  if (page.value === 1) {
    fetchPreview()
  } else {
    page.value = 1
  }
})

onMounted(() => {
  isMounted.value = true
  fetchAcademicYears()
  fetchClassesMajors()
  fetchPreview()
})
</script>

<template>
  <div class="max-w-[1600px] mx-auto p-4 lg:p-8 space-y-6 animate-fade-in">
    
    <Teleport v-if="isMounted" to="#header-actions-target">
      <div class="flex items-center justify-center w-full gap-4 relative mx-auto font-inter">
        <div class="flex items-center justify-center gap-2 flex-1 max-w-2xl mx-auto">
          <div class="relative flex-1 group">
            <SearchIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300 group-focus-within:text-indigo-600 transition-colors" />
            <input v-model="search" type="text" placeholder="Cari nama siswa atau NIS..." class="w-full py-2.5 pl-11 pr-4 bg-white border border-slate-200 rounded-xl text-xs font-bold text-slate-700 placeholder:text-slate-300 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all shadow-sm outline-none" @keyup.enter="fetchPreview" />
          </div>
          
          <!-- Filter Button -->
          <button @click="showFilters = !showFilters" class="relative p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 flex items-center gap-2 transition-all shadow-sm cursor-pointer">
            <FilterIcon class="w-4 h-4" />
            <span class="text-[10px] font-black uppercase tracking-wider pr-1">Filter</span>
            <span v-if="filters.academic_year_id || filters.class_id || filters.major_id || filters.bill_type_id || filters.period !== 'all'" class="absolute -top-1 -right-1 w-3 h-3 bg-indigo-600 rounded-full border-2 border-white shadow-sm"></span>
          </button>

          <button @click="resetFilters" class="p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm group shrink-0 transition-all cursor-pointer" title="Reset Pencarian & Filter">
            <ResetIcon class="w-4 h-4 group-hover:rotate-180 transition-transform duration-500" />
          </button>
        </div>

        <!-- Filter Modal Component -->
        <ReportsFilterModal 
          v-model="showFilters" 
          :filters="tempFilters" 
          :academicYears="academicYears" 
          :classes="classes" 
          :majors="majors" 
          :billTypes="billTypes" 
          @apply="applyFilters" 
          @reset="resetFilters" 
        />
      </div>
    </Teleport>

    <div class="flex items-center gap-2 bg-slate-100/50 p-1 rounded-2xl w-fit">
      <button v-for="t in [
        { id: 'payments', label: 'Semua Pembayaran' },
        { id: 'arrears', label: 'Data Tunggakan' }
      ]" :key="t.id" @click="activeTab = t.id" :class="[activeTab === t.id ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-500 hover:text-slate-700', 'px-6 py-2.5 rounded-xl text-[10px] font-black uppercase tracking-widest transition-all']">
        {{ t.label }}
      </button>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div v-for="card in summaryCards" :key="card.label" class="bg-white border border-slate-200 rounded-2xl p-5 shadow-sm">
        <p class="text-[9px] font-black text-slate-400 uppercase tracking-[0.22em]">{{ card.label }}</p>
        <p class="mt-2 text-xl font-black text-slate-800">{{ card.value }}</p>
        <p class="mt-1 text-[10px] font-bold text-slate-400 uppercase tracking-wider">{{ card.helper }}</p>
      </div>
    </div>

    <!-- 3. Dynamic Tables -->
    <div class="bg-white rounded-xl border border-slate-200 shadow-sm flex flex-col min-h-[600px] overflow-hidden">
      <div class="p-4 border-b border-slate-100 bg-slate-50/30 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-1.5 h-5 bg-indigo-600 rounded-full"></div>
          <h3 class="text-xs font-black text-slate-800 uppercase tracking-widest">
            {{ activeTab === 'payments' ? 'Detail Transaksi Pembayaran' : (activeTab === 'arrears' ? 'Daftar Siswa Menunggak' : 'Rekapitulasi Data Akademik') }}
          </h3>
        </div>
        <div class="flex items-center gap-3">
          <button @click="downloadReport('xlsx', 'global')" :disabled="isOffline" :title="isOffline ? 'Export membutuhkan server online agar data laporan terbaru.' : 'Ekspor Excel'" :class="['font-bold py-1.5 px-3 rounded-lg border text-[10px] flex items-center gap-2 transition-all shadow-sm', isOffline ? 'bg-amber-50 border-amber-200 text-amber-700 cursor-not-allowed' : 'bg-white text-slate-600 border-slate-200 hover:bg-slate-50']">
            <ExcelIcon class="w-3.5 h-3.5" :class="isOffline ? 'text-amber-600' : 'text-emerald-600'" />
            <span>{{ isOffline ? 'Excel Online Saja' : 'Ekspor Excel' }}</span>
          </button>
          <button @click="downloadReport('pdf', 'global')" :disabled="isOffline" :title="isOffline ? 'Export membutuhkan server online agar data laporan terbaru.' : 'Ekspor PDF'" :class="['font-bold py-1.5 px-3 rounded-lg border text-[10px] flex items-center gap-2 transition-all shadow-sm', isOffline ? 'bg-amber-50 border-amber-200 text-amber-700 cursor-not-allowed' : 'bg-white text-slate-600 border-slate-200 hover:bg-slate-50']">
            <PDFIcon class="w-3.5 h-3.5" :class="isOffline ? 'text-amber-600' : 'text-rose-600'" />
            <span>{{ isOffline ? 'PDF Online Saja' : 'Ekspor PDF' }}</span>
          </button>
        </div>
      </div>

      <div class="flex-1 overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <!-- TAB: PAYMENTS -->
          <template v-if="activeTab === 'payments'">
            <thead>
              <tr class="bg-slate-50/50 border-b border-slate-100 text-[10px] font-black text-slate-400 uppercase tracking-widest">
                <th class="py-3 px-4">Transaksi</th>
                <th class="py-3 px-4">Rincian Tagihan</th>
                <th class="py-3 px-4 text-center">Nominal</th>
                <th class="py-3 px-4 text-center">Metode</th>
                <th class="py-3 px-4 text-right">Tanggal</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-50">
              <tr v-for="item in previewData" :key="item.id" class="border-b border-slate-100 hover:bg-slate-50/30 transition-all text-xs font-semibold text-slate-600">
                <td class="py-3 px-4">
                  <div class="flex items-center gap-4">
                    <div class="w-10 h-10 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center font-black text-[10px] border border-indigo-100 shrink-0">
                      {{ item.student_name?.[0] }}
                    </div>
                    <div class="min-w-0">
                      <p class="text-xs font-black text-slate-700 uppercase tracking-wider truncate">{{ item.student_name }}</p>
                      <p class="mt-1 text-[9px] font-bold text-slate-400 uppercase tracking-widest truncate">
                        {{ item.class_name || 'Tanpa kelas' }} <span v-if="item.transaction_ref">• {{ item.transaction_ref }}</span>
                      </p>
                    </div>
                  </div>
                </td>
                <td class="py-3 px-4 max-w-[360px]">
                  <button @click="openPaymentTypeModal(item)" class="w-full text-left group/detail" :title="displayBillTypes(item.bill_type_names)">
                    <p class="text-[10px] font-black text-slate-700 uppercase tracking-wider truncate group-hover/detail:text-indigo-600">
                      {{ paymentSummaryTitle(item) }}
                    </p>
                    <div class="mt-2 flex flex-wrap gap-1.5">
                      <span v-for="detail in paymentDetailItems(item).slice(0, 2)" :key="detail.name + detail.period" class="max-w-[160px] truncate px-2 py-1 rounded-lg bg-slate-100 text-[8px] font-black text-slate-500 uppercase tracking-wider">
                        {{ detail.name }}
                      </span>
                      <span v-if="paymentItemCount(item) > 2" class="px-2 py-1 rounded-lg bg-indigo-50 text-[8px] font-black text-indigo-600 uppercase tracking-wider">
                        +{{ paymentItemCount(item) - 2 }}
                      </span>
                    </div>
                  </button>
                </td>
                <td class="py-3 px-4 text-center">
                  <span class="text-[10px] font-black text-slate-700">{{ formatCurrency(item.amount) }}</span>
                  <p v-if="Number(item.deposit_applied || 0) > 0" class="text-[8px] font-bold text-slate-400 uppercase mt-1">Saldo {{ formatCurrency(item.deposit_applied) }} • Tunai/Gateway {{ formatCurrency(item.cash_or_gateway_amount) }}</p>
                </td>
                <td class="py-3 px-4 text-center">
                  <span class="px-3 py-1 bg-white border border-slate-100 rounded-full text-[8px] font-black text-slate-500 uppercase">{{ formatPaymentMethod(item.method) }}</span>
                </td>
                <td class="py-3 px-4 text-right">
                  <p class="text-[10px] font-bold text-slate-400 uppercase">{{ formatDateTime(item.created_at) }}</p>
                </td>
              </tr>
              <tr v-if="!previewLoading && previewData.length === 0">
                <td colspan="5" class="py-20 text-center">
                  <ClockIcon class="w-10 h-10 text-slate-200 mx-auto mb-3 animate-pulse" />
                  <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Belum ada data pembayaran sesuai filter</p>
                </td>
              </tr>
            </tbody>
          </template>

          <!-- TAB: ARREARS -->
          <template v-else-if="activeTab === 'arrears'">
            <thead>
              <tr class="bg-slate-50/50 border-b border-slate-100 text-[10px] font-black text-slate-400 uppercase tracking-widest">
                <th class="py-3 px-4">Siswa</th>
                <th class="py-3 px-4 text-center">Kelas</th>
                <th class="py-3 px-4 text-center">Tagihan</th>
                <th class="py-3 px-4 text-center">Mulai</th>
                <th class="py-3 px-4 text-center">Jatuh Tempo</th>
                <th class="py-3 px-4 text-right">Sisa Tagihan</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-50">
              <tr v-for="item in previewData" :key="item.id" class="border-b border-slate-100 hover:bg-slate-50/30 transition-all text-xs font-semibold text-slate-600">
                <td class="py-3 px-4 max-w-[200px] truncate">
                  <p class="text-xs font-black text-slate-800 uppercase tracking-wider truncate">{{ item.student_name }}</p>
                </td>
                <td class="py-3 px-4 text-center">
                  <span class="text-[10px] font-bold text-slate-600">{{ item.class_name }}</span>
                </td>
                <td class="py-3 px-4 text-center">
                  <p class="text-[10px] font-black text-slate-700 uppercase">{{ item.bill_name }}</p>
                  <p class="text-[9px] font-bold text-slate-400 uppercase">{{ item.period || '-' }}</p>
                </td>
                <td class="py-3 px-4 text-center">
                  <span class="text-[10px] font-bold text-slate-500">{{ formatDate(item.start_date) }}</span>
                </td>
                <td class="py-3 px-4 text-center">
                  <span :class="['text-[10px] font-black', getRemainingDaysText(item.due_date).startsWith('Telat') ? 'text-rose-600' : (getRemainingDaysText(item.due_date) === 'Hari Ini' ? 'text-amber-600' : 'text-slate-600')]">
                    {{ formatDate(item.due_date) }}
                    <span v-if="getRemainingDaysText(item.due_date)">({{ getRemainingDaysText(item.due_date) }})</span>
                  </span>
                </td>
                <td class="py-3 px-4 text-right">
                  <span class="text-[10px] font-black text-slate-700">{{ formatCurrency(item.amount - item.total_paid) }}</span>
                </td>
              </tr>
            </tbody>
          </template>

          <!-- TAB: ACADEMIC -->
          <template v-else>
            <thead>
              <tr class="bg-slate-50/50 border-b border-slate-100 text-[10px] font-black text-slate-400 uppercase tracking-widest">
                <th class="py-3 px-4">Nama Kategori</th>
                <th class="py-3 px-4 text-center">Tipe</th>
                <th class="py-3 px-4 text-right">Jumlah Siswa</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-50">
              <tr v-for="(item, idx) in previewData" :key="idx" class="border-b border-slate-100 hover:bg-slate-50/30 transition-all text-xs font-semibold text-slate-600">
                <td class="py-3 px-4 max-w-[200px] truncate">
                  <p class="text-xs font-black text-slate-800 uppercase tracking-wider truncate">{{ item.name }}</p>
                </td>
                <td class="py-3 px-4 text-center">
                  <span class="px-3 py-1 bg-slate-100 rounded-full text-[8px] font-black text-slate-500 uppercase">{{ item.type }}</span>
                </td>
                <td class="py-3 px-4 text-right">
                  <span class="text-[10px] font-black text-indigo-600">{{ item.count }} Siswa</span>
                </td>
              </tr>
            </tbody>
          </template>
        </table>
      </div>

      <!-- Footer with Pagination -->
      <div class="px-6 py-4 bg-slate-50/50 border-t border-slate-100 flex items-center justify-between mt-auto">
        <div class="flex items-center gap-6">
          <div class="flex items-center gap-3">
            <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Tampilkan</span>
            <select v-model="limit" class="bg-white border border-slate-200 rounded-lg text-[10px] font-black text-slate-600 px-2 py-1 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all cursor-pointer shadow-sm">
              <option :value="10">10</option>
              <option :value="25">25</option>
              <option :value="50">50</option>
              <option :value="100">100</option>
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
            class="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 bg-white text-slate-400 hover:text-indigo-600 hover:border-indigo-100 hover:bg-indigo-50/30 disabled:opacity-20 disabled:hover:bg-white disabled:hover:border-slate-200 transition-all cursor-pointer"
          >
            <PrevIcon class="w-3.5 h-3.5" />
          </button>

          <!-- Page Numbers (Max 3) -->
          <div class="flex items-center gap-1">
            <button 
              v-for="p in visiblePages" 
              :key="p"
              @click="page = p"
              class="w-8 h-8 flex items-center justify-center rounded-lg text-[10px] font-black transition-all cursor-pointer"
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
            class="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 bg-white text-slate-400 hover:text-indigo-600 hover:border-indigo-100 hover:bg-indigo-50/30 disabled:opacity-20 disabled:hover:bg-white disabled:hover:border-slate-200 transition-all cursor-pointer"
          >
            <NextIcon class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    </div>

    <Teleport to="body">
      <transition name="page">
        <div v-if="showPaymentTypeModal && selectedPayment" class="fixed inset-0 z-[200] flex items-center justify-center p-6">
          <div class="absolute inset-0 bg-slate-900/40 backdrop-blur-sm" @click="showPaymentTypeModal = false"></div>
          <div class="bg-white w-full max-w-lg relative z-10 rounded-[2rem] shadow-2xl overflow-hidden animate-scale-in">
            <div class="p-7 border-b border-slate-100 flex items-center justify-between">
              <div>
                <h3 class="text-lg font-black text-slate-800 tracking-tight">Detail Jenis Pembayaran</h3>
                <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest mt-1">{{ selectedPayment.student_name }} • {{ formatDateTime(selectedPayment.created_at) }}</p>
              </div>
              <button @click="showPaymentTypeModal = false" class="w-10 h-10 rounded-xl bg-slate-50 text-slate-400 hover:bg-slate-100 hover:text-slate-700 flex items-center justify-center transition-all">
                <XIcon class="w-5 h-5" />
              </button>
            </div>
            <div class="p-7 space-y-4">
              <div class="grid grid-cols-3 gap-3">
                <div class="rounded-2xl bg-slate-50 border border-slate-100 p-4">
                  <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Total</p>
                  <p class="mt-1 text-sm font-black text-slate-800">{{ formatCurrency(selectedPayment.amount) }}</p>
                </div>
                <div class="rounded-2xl bg-slate-50 border border-slate-100 p-4">
                  <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Saldo</p>
                  <p class="mt-1 text-sm font-black text-slate-800">{{ formatCurrency(selectedPayment.deposit_applied) }}</p>
                </div>
                <div class="rounded-2xl bg-slate-50 border border-slate-100 p-4">
                  <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Metode</p>
                  <p class="mt-1 text-sm font-black text-slate-800 uppercase">{{ formatPaymentMethod(selectedPayment.method) }}</p>
                </div>
              </div>

              <div class="rounded-2xl border border-slate-100 overflow-hidden">
                <div class="px-4 py-3 bg-slate-50 text-[10px] font-black text-slate-400 uppercase tracking-widest">Daftar Tagihan Terbayar</div>
                <div class="divide-y divide-slate-100 max-h-72 overflow-y-auto">
                  <div v-for="(detail, idx) in paymentDetailItems(selectedPayment)" :key="idx" class="px-4 py-3 flex items-center justify-between gap-4">
                    <div class="flex items-center gap-3 min-w-0">
                      <div class="w-8 h-8 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center text-[10px] font-black">{{ idx + 1 }}</div>
                      <div class="min-w-0">
                        <p class="text-xs font-black text-slate-700 uppercase tracking-wider truncate">{{ detail.name }}</p>
                        <p v-if="detail.period" class="text-[8px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">Periode {{ detail.period }}</p>
                      </div>
                    </div>
                    <span v-if="detail.amount > 0" class="text-[10px] font-black text-slate-600 shrink-0">{{ formatCurrency(detail.amount) }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>

  </div>
</template>

<style scoped lang="postcss">
.btn-page {
  @apply p-2 text-slate-300 hover:text-indigo-600 hover:bg-slate-50 rounded-xl transition-all disabled:opacity-20;
}
.btn-export {
  @apply flex items-center gap-2 px-4 py-2 bg-emerald-50 border border-emerald-100 text-emerald-600 rounded-xl text-[9px] font-black uppercase hover:bg-emerald-600 hover:text-white transition-all;
}
.white-card {
  @apply bg-white border border-slate-200 rounded-2xl transition-all duration-300 shadow-sm overflow-hidden;
}
.no-scrollbar::-webkit-scrollbar { display: none; }
.no-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }
</style>
