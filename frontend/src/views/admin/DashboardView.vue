<script setup>
import { ref, reactive, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import * as XLSX from 'xlsx'
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
  Eye as ReadIcon,
  TrendingUp as TrendIcon,
  AlertTriangle as AlertIcon,
  CreditCard as CardIcon,
  PlusCircle as PlusIcon,
  Send as SendIcon,
  LogOut as LogoutIcon,
  RotateCcw as ResetIcon,
  ChevronDown as ChevronDownIcon
} from 'lucide-vue-next'
import { Bar, Pie, Doughnut } from 'vue-chartjs'
import { useAuthStore } from '../../store/auth'
import { Chart as ChartJS, registerables } from 'chart.js'
ChartJS.register(...registerables)
const router = useRouter()
const authStore = useAuthStore()
const isOffline = computed(() => authStore.isOffline || (typeof navigator !== 'undefined' && navigator.onLine === false))
const isMounted = ref(false)
const loading = ref(false)

const stats = ref({
  users: { total: 0, growth: 0 },
  students: { total: 0, growth: 0, total_all: 0 },
  bills: { total: 0, growth: 0 },
  payments: { total: 0, growth: 0 },
  paid_amount: 0,
  unpaid_amount: 0,
  paid_count: 0,
  unpaid_count: 0,
  payments_today: 0,
  failed_reminders: 0,
  critical_bills: { overdue: [], due_soon: [] }
})

const demographics = ref({
  gender: {},
  major: {},
  class: {},
  whatsapp: { pending: 0, sent: 0, delivered: 0, read: 0, failed: 0 },
  email: { pending: 0, sent: 0, delivered: 0, read: 0, failed: 0 }
})

const recentPayments = ref([])
const recentNotifications = ref([])
const paymentTrend = ref([])
const academicYears = ref([])
const availableEntryYears = ref([])
const classes = ref([])
const majors = ref([])
const showDrillDown = ref(false)
const drillDownTitle = ref('')
const drillDownData = ref([])
const drillDownLoading = ref(false)

const filters = reactive({
  period: 'all',
  academic_year_id: '',
  class_id: '',
  major_id: '',
  search: '',
  ref_date: new Date().toISOString().substring(0, 10),
  start_date: '',
  end_date: ''
})

const tempFilters = reactive({
  period: 'all',
  academic_year_id: '',
  class_id: '',
  major_id: '',
  search: '',
  ref_date: new Date().toISOString().substring(0, 10),
  start_date: '',
  end_date: ''
})

const fetchStats = async () => {
  loading.value = true
  try {
    const res = await axios.get('dashboard/stats', { params: { ...filters } })
    if (res.data?.data) {
      const cb = res.data.data.critical_bills || {}
      stats.value = {
        ...stats.value,
        ...res.data.data.stats,
        critical_bills: {
          overdue: cb.overdue || [],
          due_soon: cb.due_soon || []
        }
      }
      demographics.value = res.data.data.demographics || demographics.value
      recentPayments.value = res.data.data.recent_payments || []
      recentNotifications.value = res.data.data.recent_notifications || []
      paymentTrend.value = res.data.data.payment_trend || []
      availableEntryYears.value = res.data.data.entry_years || []
    }
  } catch (err) {} finally {
    loading.value = false
  }
}

const fetchYears = async () => {
  try {
    const res = await axios.get('academic/years', { params: { limit: 100, status: 'active' } })
    academicYears.value = res.data.data.data || []
  } catch (err) {}
}

const fetchClassesMajors = async () => {
  try {
    const [resCls, resMaj] = await Promise.all([
      axios.get('academic/class', { params: { limit: 100, status: 'active' } }),
      axios.get('academic/major', { params: { limit: 100, status: 'active' } })
    ])
    const pCls = resCls.data?.data
    const pMaj = resMaj.data?.data
    classes.value = (pCls?.data || (Array.isArray(pCls) ? pCls : [])).filter(c => c.is_active !== false)
    majors.value = (pMaj?.data || (Array.isArray(pMaj) ? pMaj : [])).filter(m => m.is_active !== false)
  } catch (err) {
    console.error('Gagal fetch classes/majors:', err)
  }
}

const resetFilters = () => {
  const now = new Date()
  const defaults = {
    period: 'all',
    academic_year_id: '',
    class_id: '',
    major_id: '',
    search: '',
    ref_date: now.toISOString().substring(0, 10),
    start_date: '',
    end_date: ''
  }
  Object.assign(tempFilters, defaults)
  Object.assign(filters, defaults)
  fetchStats()
}

const applyFilters = () => {
  Object.assign(filters, tempFilters)
  fetchStats()
}

const fetchCommunicationDetails = async (status, channel = 'whatsapp') => {
  drillDownLoading.value = true
  drillDownTitle.value = `${channelLabel(channel)} - ${statusLabel(status)}`
  showDrillDown.value = true
  try {
    const res = await axios.get('dashboard/communication-details', { 
      params: { 
        status, 
        channel,
        period: filters.period,
        academic_year_id: filters.academic_year_id,
        class_id: filters.class_id,
        major_id: filters.major_id,
        ref_date: filters.ref_date,
        start_date: filters.start_date,
        end_date: filters.end_date
      } 
    })
    drillDownData.value = res.data.data || []
  } catch (err) {} finally {
    drillDownLoading.value = false
  }
}

const sendManualReminder = async (billID) => {
  if (!confirm('Kirim pengingat WhatsApp & Email manual sekarang?')) return
  try {
    await axios.post(`/finance/bills/${billID}/remind`)
    alert('Pengingat berhasil dijadwalkan!')
  } catch (err) {
    alert('Gagal mengirim pengingat.')
  }
}

const printReceipt = (paymentID) => {
  const url = router.resolve({ name: 'receipt-print', params: { id: paymentID } }).href
  window.open(url, '_blank')
}

const exportTrendToExcel = () => {
  if (isOffline.value) return
  if (paymentTrend.value.length === 0) return
  const rows = [
    ['Laporan Tren Pendapatan'],
    ['Periode', filters.period === 'all' ? 'Semua periode' : filters.period],
    ['Tanggal Export', new Date().toLocaleString('id-ID')],
    [],
    ['Tanggal', 'Total Pemasukan']
  ]
  paymentTrend.value.forEach(t => rows.push([t.date, Number(t.total || 0)]))
  const ws = XLSX.utils.aoa_to_sheet(rows)
  ws['!cols'] = [{ wch: 18 }, { wch: 22 }]
  for (let i = 6; i <= paymentTrend.value.length + 5; i++) {
    const cell = ws[`B${i}`]
    if (cell) cell.z = '"Rp" #,##0'
  }
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, 'Tren Pendapatan')
  XLSX.writeFile(wb, `Tren_Pendapatan_${filters.ref_date}.xlsx`)
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
      return { name: name || 'Pembayaran Deposit / Kustom', period: period || '', amount: Number(String(amount || '0').replace(/[^0-9.-]/g, '')) || 0 }
    }).filter(item => item.name)
  }
  return billTypeItems(payment?.bill_type_names).map(name => ({ name, period: '', amount: 0 }))
}

const paymentSummaryTitle = (payment) => {
  const items = paymentDetailItems(payment)
  if (items.length === 0) return 'Pembayaran Deposit / Kustom'
  if (items.length === 1) return items[0].name
  return `${items[0].name} + ${items.length - 1} tagihan`
}

const formatPaymentMethod = (method) => {
  if (!method) return '-'
  return String(method).replace(/_/g, ' ').toUpperCase()
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
  return new Date(dateStr).toLocaleString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

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

const getPercentage = (val, total) => {
  if (!total) return '0%'
  return ((val / total) * 100).toFixed(1) + '%'
}

const communicationStatuses = [
  { key: 'pending', label: 'Menunggu', icon: ClockIcon, note: 'Masuk antrean, belum ada hasil kirim' },
  { key: 'sent', label: 'Terkirim', icon: SendIcon, note: 'Berhasil dikirim dari sistem' },
  { key: 'delivered', label: 'Diterima', icon: CheckIcon, note: 'Diterima perangkat tujuan' },
  { key: 'read', label: 'Dibaca', icon: ReadIcon, note: 'Sudah dibuka/dibaca penerima' },
  { key: 'failed', label: 'Gagal', icon: ErrorIcon, note: 'Gagal dikirim, lihat alasan error' }
]

const channelLabel = (channel) => channel === 'email' ? 'Email' : 'WhatsApp'

const statusLabel = (status) => {
  return communicationStatuses.find(s => s.key === String(status).toLowerCase())?.label || status
}

const normalizedChannelStats = (channel) => {
  const raw = demographics.value?.[channel] || {}
  return communicationStatuses.reduce((acc, status) => {
    acc[status.key] = Number(raw[status.key] || (status.key === 'sent' ? raw.success : 0) || 0)
    return acc
  }, {})
}

const communicationTotal = (channel) => {
  const data = normalizedChannelStats(channel)
  return communicationStatuses.reduce((acc, status) => acc + Number(data[status.key] || 0), 0)
}

const communicationRows = (channel) => {
  const data = normalizedChannelStats(channel)
  const total = communicationTotal(channel)
  return communicationStatuses.map(status => ({
    ...status,
    count: Number(data[status.key] || 0),
    percent: getPercentage(Number(data[status.key] || 0), total)
  }))
}

const communicationTables = computed(() => [
  {
    channel: 'whatsapp',
    title: 'Efikasi WhatsApp',
    icon: WAIcon,
    accent: 'emerald',
    total: communicationTotal('whatsapp'),
    rows: communicationRows('whatsapp')
  },
  {
    channel: 'email',
    title: 'Efikasi Email',
    icon: MailIcon,
    accent: 'sky',
    total: communicationTotal('email'),
    rows: communicationRows('email')
  }
])

const statusBadgeClass = (status) => {
  switch (String(status).toLowerCase()) {
    case 'read':
      return 'bg-emerald-50 text-emerald-600 border-emerald-100'
    case 'delivered':
      return 'bg-sky-50 text-sky-600 border-sky-100'
    case 'failed':
      return 'bg-rose-50 text-rose-600 border-rose-100'
    case 'pending':
      return 'bg-amber-50 text-amber-600 border-amber-100'
    default:
      return 'bg-slate-50 text-slate-600 border-slate-200'
  }
}

const channelBadgeClass = (channel) => {
  return channel === 'email'
    ? 'bg-sky-50 text-sky-600 border-sky-100'
    : 'bg-emerald-50 text-emerald-600 border-emerald-100'
}

// Chart Options
const baseOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { position: 'bottom', labels: { boxWidth: 10, font: { size: 9, weight: 'bold' } } },
    tooltip: {
      callbacks: {
        label: (ctx) => {
          const total = ctx.dataset.data.reduce((a, b) => a + b, 0) || 1
          return `${ctx.label}: ${ctx.raw} (${((ctx.raw / total) * 100).toFixed(1)}%)`
        }
      }
    }
  }
}

// Dedicated bar chart options with Rupiah currency formatting
const trendOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      backgroundColor: 'rgba(15,23,42,0.92)',
      padding: 12,
      cornerRadius: 12,
      titleFont: { size: 9, weight: 'bold', family: 'Inter' },
      bodyFont: { size: 11, weight: 'bold', family: 'Inter' },
      callbacks: {
        title: (ctx) => ctx[0].label,
        label: (ctx) => {
          return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR', minimumFractionDigits: 0 }).format(ctx.raw || 0)
        }
      }
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      grid: { color: 'rgba(0,0,0,0.04)' },
      ticks: {
        font: { size: 9, weight: 'bold' },
        color: '#94a3b8',
        callback: (val) => {
          if (val >= 1_000_000) return 'Rp ' + (val / 1_000_000).toFixed(1) + 'jt'
          if (val >= 1_000) return 'Rp ' + (val / 1_000).toFixed(0) + 'rb'
          return 'Rp ' + val
        }
      }
    },
    x: {
      grid: { display: false },
      ticks: { font: { size: 9, weight: 'bold' }, color: '#94a3b8' }
    }
  }
}

// Data Computed
// Data Computed
const trendData = computed(() => ({
  labels: paymentTrend.value.map(t => {
    if (!t.date) return ''
    const d = new Date(t.date)
    return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}`
  }),
  datasets: [{ 
    label: 'Pemasukan', 
    backgroundColor: '#4f46e5', 
    borderRadius: 12, 
    data: paymentTrend.value.map(t => Number(t.total) || 0),
    datalabels: { display: true, anchor: 'end', align: 'top', font: { size: 9, weight: 'bold' } }
  }]
}))

const genderData = computed(() => {
  const male = (demographics.value.gender?.L || 0) + (demographics.value.gender?.['Laki-laki'] || 0) + (demographics.value.gender?.Male || 0)
  const female = (demographics.value.gender?.P || 0) + (demographics.value.gender?.Perempuan || 0) + (demographics.value.gender?.Female || 0)
  return {
    labels: ['Laki-laki', 'Perempuan'],
    datasets: [{ backgroundColor: ['#4f46e5', '#f43f5e'], data: [male, female] }]
  }
})

const methodData = computed(() => ({
  labels: (demographics.value.payment_methods || []).map(m => m.method),
  datasets: [{ backgroundColor: ['#4f46e5', '#10b981', '#f59e0b', '#3b82f6'], data: (demographics.value.payment_methods || []).map(m => m.count) }]
}))

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

onMounted(async () => {
  isMounted.value = true
  await Promise.all([fetchYears(), fetchClassesMajors(), fetchStats()])
  window.addEventListener('new-payment', fetchStats)
  window.addEventListener('notification-status-changed', fetchStats)
})

onUnmounted(() => {
  window.removeEventListener('new-payment', fetchStats)
  window.removeEventListener('notification-status-changed', fetchStats)
})
</script>

<template>
  <div class="max-w-[1600px] mx-auto p-4 lg:p-8 space-y-10 animate-fade-in relative">

    <!-- 1. Filters — Inline in Navbar -->
    <Teleport v-if="isMounted" to="#header-actions-target">
      <div class="flex items-end gap-2 font-inter pr-2">

        <!-- Search -->
        <div class="flex flex-col gap-1 min-w-[190px]">
          <label class="text-[8px] font-black text-slate-400 uppercase tracking-widest px-1">Cari</label>
          <div class="relative">
            <SearchIcon class="w-3.5 h-3.5 absolute left-2.5 top-1/2 -translate-y-1/2 text-slate-300" />
            <input
              v-model="tempFilters.search"
              @keyup.enter="applyFilters"
              @change="applyFilters"
              type="text"
              placeholder="Siswa, NIS, tagihan..."
              class="h-8 w-full pl-8 pr-3 bg-white border border-slate-200 rounded-lg text-[10px] font-bold text-slate-700 shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-100 focus:border-indigo-400 hover:border-indigo-300 transition-all"
            />
          </div>
        </div>

        <!-- Periode -->
        <div class="flex flex-col gap-1">
          <label class="text-[8px] font-black text-slate-400 uppercase tracking-widest px-1">Periode</label>
          <div class="relative">
            <select v-model="tempFilters.period" @change="applyFilters"
              class="h-8 pl-3 pr-7 bg-white border border-slate-200 rounded-lg text-[10px] font-bold text-slate-700 shadow-sm appearance-none cursor-pointer hover:border-indigo-300 focus:outline-none focus:ring-2 focus:ring-indigo-100 focus:border-indigo-400 transition-all">
              <option value="all">Semua</option>
              <option value="daily">Harian</option>
              <option value="monthly">Bulanan</option>
              <option value="yearly">Tahunan</option>
              <option value="custom">Kustom</option>
            </select>
            <ChevronDownIcon class="w-3 h-3 absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
          </div>
        </div>

        <!-- Ref Date — daily / monthly / yearly -->
        <template v-if="tempFilters.period !== 'all' && tempFilters.period !== 'custom'">
          <div class="flex flex-col gap-1">
            <label class="text-[8px] font-black text-slate-400 uppercase tracking-widest px-1">
              {{ tempFilters.period === 'daily' ? 'Tanggal' : tempFilters.period === 'monthly' ? 'Bulan' : 'Tahun' }}
            </label>
            <input
              v-model="tempFilters.ref_date"
              @change="applyFilters"
              :type="tempFilters.period === 'monthly' ? 'month' : (tempFilters.period === 'daily' ? 'date' : 'number')"
              :placeholder="tempFilters.period === 'yearly' ? '2025' : ''"
              class="h-8 px-3 bg-white border border-slate-200 rounded-lg text-[10px] font-bold text-slate-700 shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-100 focus:border-indigo-400 hover:border-indigo-300 transition-all w-32"
            />
          </div>
        </template>

        <!-- Custom date range -->
        <template v-if="tempFilters.period === 'custom'">
          <div class="flex flex-col gap-1">
            <label class="text-[8px] font-black text-slate-400 uppercase tracking-widest px-1">Dari</label>
            <input v-model="tempFilters.start_date" @change="applyFilters" type="date"
              class="h-8 px-3 bg-white border border-slate-200 rounded-lg text-[10px] font-bold text-slate-700 shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-100 focus:border-indigo-400 hover:border-indigo-300 transition-all w-32" />
          </div>
          <div class="flex flex-col gap-1">
            <label class="text-[8px] font-black text-slate-400 uppercase tracking-widest px-1">Sampai</label>
            <input v-model="tempFilters.end_date" @change="applyFilters" type="date"
              class="h-8 px-3 bg-white border border-slate-200 rounded-lg text-[10px] font-bold text-slate-700 shadow-sm focus:outline-none focus:ring-2 focus:ring-indigo-100 focus:border-indigo-400 hover:border-indigo-300 transition-all w-32" />
          </div>
        </template>

        <!-- Divider -->
        <div class="h-8 w-px bg-slate-200 self-end mb-0.5 mx-0.5"></div>

        <!-- Angkatan -->
        <div class="flex flex-col gap-1">
          <label class="text-[8px] font-black text-slate-400 uppercase tracking-widest px-1">Angkatan</label>
          <div class="relative">
            <select v-model="tempFilters.academic_year_id" @change="applyFilters"
              class="h-8 pl-3 pr-7 bg-white border border-slate-200 rounded-lg text-[10px] font-bold text-slate-700 shadow-sm appearance-none cursor-pointer hover:border-indigo-300 focus:outline-none focus:ring-2 focus:ring-indigo-100 focus:border-indigo-400 transition-all">
              <option value="">Semua</option>
              <option v-for="y in academicYears" :key="y.id" :value="y.id">{{ y.year }}</option>
            </select>
            <ChevronDownIcon class="w-3 h-3 absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
          </div>
        </div>

        <!-- Jurusan -->
        <div class="flex flex-col gap-1">
          <label class="text-[8px] font-black text-slate-400 uppercase tracking-widest px-1">Jurusan</label>
          <div class="relative">
            <select v-model="tempFilters.major_id" @change="applyFilters"
              class="h-8 pl-3 pr-7 bg-white border border-slate-200 rounded-lg text-[10px] font-bold text-slate-700 shadow-sm appearance-none cursor-pointer hover:border-indigo-300 focus:outline-none focus:ring-2 focus:ring-indigo-100 focus:border-indigo-400 transition-all">
              <option value="">Semua</option>
              <option v-for="m in majors" :key="m.id" :value="m.id">{{ m.name }}</option>
            </select>
            <ChevronDownIcon class="w-3 h-3 absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
          </div>
        </div>

        <!-- Kelas -->
        <div class="flex flex-col gap-1">
          <label class="text-[8px] font-black text-slate-400 uppercase tracking-widest px-1">Kelas</label>
          <div class="relative">
            <select v-model="tempFilters.class_id" @change="applyFilters"
              class="h-8 pl-3 pr-7 bg-white border border-slate-200 rounded-lg text-[10px] font-bold text-slate-700 shadow-sm appearance-none cursor-pointer hover:border-indigo-300 focus:outline-none focus:ring-2 focus:ring-indigo-100 focus:border-indigo-400 transition-all">
              <option value="">Semua</option>
              <option v-for="c in classes" :key="c.id" :value="c.id">{{ c.name }}</option>
            </select>
            <ChevronDownIcon class="w-3 h-3 absolute right-2 top-1/2 -translate-y-1/2 text-slate-400 pointer-events-none" />
          </div>
        </div>

        <!-- Divider -->
        <div class="h-8 w-px bg-slate-200 self-end mb-0.5 mx-0.5"></div>

        <!-- Reset Button — sejajar dengan dropdown (items-end sudah handle ini) -->
        <div class="flex flex-col gap-1">
          <button @click="resetFilters"
            class="h-8 w-8 flex items-center justify-center bg-white text-slate-500 hover:text-slate-700 border border-slate-200 rounded-lg shadow-sm hover:bg-slate-50 transition-all group"
            title="Reset semua filter">
            <ResetIcon class="w-3.5 h-3.5 group-hover:rotate-180 transition-transform duration-500" />
          </button>
        </div>

      </div>
    </Teleport>


    <!-- 2. Summary Cards Grid (5 Cards in 1 Row) -->
    <div class="grid grid-cols-2 lg:grid-cols-5 gap-6">
      <div v-for="(v, k) in {
        students: { icon: StudentIcon, label: 'Siswa Terdaftar', color: 'indigo', val: stats.students?.total_all || 0 },
        users: { icon: UsersIcon, label: 'Total Pengguna', color: 'blue', val: stats.users?.total || 0 },
        unpaid_amount: { icon: AlertIcon, label: 'Tagihan Menunggak', color: 'rose', isP: true, val: stats.unpaid_amount },
        paid_amount: { icon: WalletIcon, label: 'Pemasukan Pembayaran', color: 'emerald', isP: true, val: stats.paid_amount },
        paid_count: { icon: BillIcon, label: 'Total Transaksi', color: 'indigo', val: stats.paid_count }
      }" :key="k" class="white-card p-6 group hover:translate-y-[-4px]">
        <div class="flex items-center justify-between mb-4">
          <div :class="[`bg-${v.color}-50 text-${v.color}-600`, 'w-10 h-10 rounded-xl flex items-center justify-center group-hover:scale-110 transition-all']">
            <component :is="v.icon" class="w-5 h-5" />
          </div>
        </div>
        <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">{{ v.label }}</p>
        <h3 class="text-xl font-black text-slate-800 truncate mt-1">
          {{ v.isP ? formatCurrency(v.val) : (v.val || 0) }}
        </h3>
      </div>
    </div>

    <!-- 3. Primary Charts Row -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Activity Log (Compact) -->
      <div class="white-card p-0 flex flex-col h-[550px]">
        <div class="p-6 border-b border-slate-50 flex items-center justify-between">
           <h3 class="text-xs font-black text-slate-800 uppercase tracking-widest">Aktivitas Terbaru</h3>
           <ClockIcon class="w-4 h-4 text-slate-300" />
        </div>
        <div class="flex-1 overflow-y-auto p-6 space-y-6 custom-scrollbar">
            <div v-for="pay in recentPayments" :key="'p'+pay.id" class="flex gap-4 group">
               <div class="w-10 h-10 bg-emerald-50 text-emerald-600 rounded-xl flex items-center justify-center shrink-0">
                  <WalletIcon class="w-4 h-4" />
               </div>
               <div class="flex-1 min-w-0">
                  <div class="flex items-center justify-between">
                    <p class="text-[11px] font-black text-slate-800 uppercase leading-tight truncate">{{ pay.student_name }} Membayar {{ paymentSummaryTitle(pay) }}</p>
                    <button @click="printReceipt(pay.id)" class="opacity-0 group-hover:opacity-100 transition-all text-slate-400 hover:text-indigo-600 shrink-0" title="Cetak kwitansi">
                      <PrintIcon class="w-3.5 h-3.5" />
                    </button>
                  </div>
                  <p class="text-[10px] font-bold text-slate-500 uppercase mt-1">{{ formatCurrency(pay.amount) }} • {{ formatPaymentMethod(pay.method) }}</p>
                  <p class="text-[9px] text-slate-300 mt-0.5">{{ formatDateTime(pay.created_at) }} <span v-if="pay.transaction_ref">• {{ pay.transaction_ref }}</span></p>
               </div>
            </div>
           <div v-for="notif in recentNotifications" :key="'n'+notif.id" class="flex gap-4">
              <div :class="[String(notif.delivery_status).toLowerCase() === 'failed' ? 'bg-rose-50 text-rose-500' : (notif.channel === 'email' ? 'bg-sky-50 text-sky-600' : 'bg-emerald-50 text-emerald-600'), 'w-10 h-10 rounded-xl flex items-center justify-center shrink-0']">
                 <component :is="String(notif.delivery_status).toLowerCase() === 'failed' ? AlertIcon : SendIcon" class="w-4 h-4" />
              </div>
              <div class="min-w-0 flex-1">
                 <div class="flex items-center gap-2">
                   <p class="text-[11px] font-black text-slate-800 uppercase leading-tight truncate">{{ notif.title }}</p>
                   <span class="px-2 py-0.5 rounded-full text-[8px] font-black uppercase tracking-widest border shrink-0" :class="channelBadgeClass(notif.channel || (notif.whatsapp_id ? 'whatsapp' : 'email'))">
                     {{ notif.channel || (notif.whatsapp_id ? 'WA' : 'Email') }}
                   </span>
                 </div>
                 <p class="text-[10px] font-medium text-slate-500 line-clamp-1 mt-1">{{ notif.recipient_name || 'Penerima' }} • {{ notif.message }}</p>
                 <div class="mt-1 flex items-center gap-2">
                   <span class="px-2 py-0.5 rounded-full text-[8px] font-black uppercase tracking-widest border" :class="statusBadgeClass(notif.delivery_status)">{{ statusLabel(notif.delivery_status) }}</span>
                   <span class="text-[9px] text-slate-300">{{ formatDateTime(notif.created_at) }}</span>
                 </div>
              </div>
           </div>
           <div v-if="recentPayments.length === 0 && recentNotifications.length === 0" class="py-16 text-center">
             <ClockIcon class="w-10 h-10 text-slate-200 mx-auto mb-3" />
             <p class="text-[10px] font-black text-slate-300 uppercase tracking-[0.2em]">Belum ada aktivitas sesuai filter</p>
           </div>
        </div>
        <button @click="router.push('/reports')" class="p-4 bg-slate-50 text-indigo-600 text-[10px] font-black uppercase text-center hover:bg-indigo-50 transition-all border-t border-slate-100">
           Lihat Semua Aktivitas
        </button>
      </div>

      <!-- Main Revenue Chart -->
      <div class="lg:col-span-2 white-card p-10 flex flex-col h-[550px]">
        <div class="flex items-center justify-between mb-8">
           <div>
              <h3 class="text-xs font-black text-slate-800 uppercase tracking-widest">Tren Pendapatan</h3>
              <p class="text-[10px] font-bold text-slate-400 uppercase">Visualisasi Aliran Kas Sekolah</p>
           </div>
           <ExcelIcon :class="['w-5 h-5', isOffline ? 'text-amber-500 cursor-not-allowed' : 'text-emerald-600 cursor-pointer']" :title="isOffline ? 'Export Excel membutuhkan server online agar data tren terbaru.' : 'Ekspor Excel'" @click="exportTrendToExcel" />
        </div>
        <div class="flex-1">
           <Bar :data="trendData" :options="trendOptions" />
        </div>
      </div>
    </div>

    <!-- 5. Operational Sections (Activity Style) -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mt-10">
      <!-- Critical Arrears -->
      <div class="white-card p-0 flex flex-col h-[600px] lg:order-2">
        <div class="p-8 border-b border-slate-50 flex items-center justify-between">
           <div class="flex items-center gap-3">
              <div class="w-1.5 h-6 bg-rose-500 rounded-full"></div>
              <h3 class="text-xs font-black text-slate-800 uppercase tracking-widest">Tagihan Lewat Jatuh Tempo ({{ stats.critical_bills?.overdue?.length || 0 }})</h3>
           </div>
           <AlertIcon class="w-4 h-4 text-rose-300" />
        </div>
        <div class="flex-1 overflow-y-auto p-8 space-y-6 custom-scrollbar">
           <div v-for="bill in stats.critical_bills?.overdue" :key="bill.id" class="flex gap-4 group">
              <div class="w-11 h-11 bg-rose-50 text-rose-600 rounded-2xl flex items-center justify-center shrink-0">
                 <AlertIcon class="w-5 h-5" />
              </div>
              <div class="flex-1">
                 <div class="flex items-center justify-between">
                   <p class="text-[11px] font-black text-slate-800 uppercase leading-tight">{{ bill.student_name }}</p>
                   <button @click="sendManualReminder(bill.id)" class="px-3 py-1 bg-indigo-50 text-indigo-600 rounded-lg text-[8px] font-black uppercase hover:bg-indigo-600 hover:text-white transition-all">
                      Kirim Notifikasi
                   </button>
                 </div>
                 <p class="text-[10px] font-bold text-slate-400 uppercase mt-1">{{ bill.parent_name || 'Ortu' }} • {{ bill.parent_phone }}</p>
                 <div class="flex flex-wrap items-center gap-x-2 mt-1.5 text-[9px] font-black uppercase text-slate-400">
                    <span class="text-rose-600">Sisa Tagihan: {{ formatCurrency(bill.amount - bill.total_paid) }}</span>
                    <span class="text-slate-300">|</span>
                    <span>Mulai: {{ formatDate(bill.start_date) }}</span>
                    <span class="text-slate-300">|</span>
                     <span class="text-rose-600">Jatuh Tempo: {{ formatDate(bill.due_date) }} <span v-if="getRemainingDaysText(bill.due_date)">({{ getRemainingDaysText(bill.due_date) }})</span></span>
                 </div>
              </div>
           </div>
        </div>
        <button @click="router.push('/reports?tab=arrears')" class="p-4 bg-slate-50 text-indigo-600 text-[10px] font-black uppercase text-center hover:bg-indigo-50 transition-all border-t border-slate-100">
           Lihat Semua Tagihan Belum Lunas
        </button>
      </div>

      <!-- Upcoming Bills -->
      <div class="white-card p-0 flex flex-col h-[600px] lg:order-1">
        <div class="p-8 border-b border-slate-50 flex items-center justify-between">
           <div class="flex items-center gap-3">
              <div class="w-1.5 h-6 bg-amber-500 rounded-full"></div>
              <h3 class="text-xs font-black text-slate-800 uppercase tracking-widest">Tagihan Belum Jatuh Tempo ({{ stats.critical_bills?.due_soon?.length || 0 }})</h3>
           </div>
           <ClockIcon class="w-4 h-4 text-amber-300" />
        </div>
        <div class="flex-1 overflow-y-auto p-8 space-y-6 custom-scrollbar">
           <div v-for="bill in stats.critical_bills?.due_soon" :key="bill.id" class="flex gap-4 group">
              <div class="w-11 h-11 bg-amber-50 text-amber-600 rounded-2xl flex items-center justify-center shrink-0">
                 <ClockIcon class="w-5 h-5" />
              </div>
              <div class="flex-1">
                 <div class="flex items-center justify-between">
                   <p class="text-[11px] font-black text-slate-800 uppercase leading-tight">{{ bill.student_name }}</p>
                   <button @click="sendManualReminder(bill.id)" class="px-3 py-1 bg-indigo-50 text-indigo-600 rounded-lg text-[8px] font-black uppercase hover:bg-indigo-600 hover:text-white transition-all">
                      Kirim Notifikasi
                   </button>
                 </div>
                 <p class="text-[10px] font-bold text-slate-400 uppercase mt-1">{{ bill.parent_name || 'Ortu' }} • {{ bill.parent_phone }}</p>
                 <div class="flex flex-wrap items-center gap-x-2 mt-1.5 text-[9px] font-black uppercase text-slate-400">
                    <span class="text-amber-600">Sisa Tagihan: {{ formatCurrency(bill.amount - bill.total_paid) }}</span>
                    <span class="text-slate-300">|</span>
                    <span>Mulai: {{ formatDate(bill.start_date) }}</span>
                    <span class="text-slate-300">|</span>
                     <span class="text-amber-600">Jatuh Tempo: {{ formatDate(bill.due_date) }} <span v-if="getRemainingDaysText(bill.due_date)">({{ getRemainingDaysText(bill.due_date) }})</span></span>
                 </div>
              </div>
           </div>
        </div>
        <button @click="router.push('/reports?tab=arrears')" class="p-4 bg-slate-50 text-indigo-600 text-[10px] font-black uppercase text-center hover:bg-indigo-50 transition-all border-t border-slate-100">
           Lihat Semua Tagihan Belum Jatuh Tempo
        </button>
      </div>
    </div>

    <!-- 6. Secondary Analytics Row (2 Charts in 1 Row) -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-8 mt-10">
      <!-- Metode Pembayaran -->
      <div class="white-card p-8 h-[400px] flex flex-col">
        <h3 class="text-[10px] font-black text-slate-800 uppercase tracking-widest mb-6">Metode Pembayaran</h3>
        <div class="flex-1 relative">
          <Doughnut :data="methodData" :options="baseOptions" />
        </div>
      </div>
      <!-- Demografi Gender -->
      <div class="white-card p-8 h-[400px] flex flex-col">
        <h3 class="text-[10px] font-black text-slate-800 uppercase tracking-widest mb-6">Demografi Gender</h3>
        <div class="flex-1 relative">
          <Pie :data="genderData" :options="baseOptions" />
        </div>
      </div>
    </div>

    <!-- 8. Communication Efficacy -->
    <div class="space-y-6 mt-10">
      <div v-for="table in communicationTables" :key="table.channel" class="white-card p-0 overflow-hidden">
        <div class="p-8 border-b border-slate-50 flex items-center justify-between">
          <div class="flex items-center gap-4">
            <div :class="[
              table.channel === 'whatsapp' ? 'bg-emerald-50 text-emerald-600' : 'bg-sky-50 text-sky-600',
              'w-11 h-11 rounded-2xl flex items-center justify-center'
            ]">
              <component :is="table.icon" class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-xs font-black text-slate-800 uppercase tracking-widest">{{ table.title }}</h3>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-1">
                Alur status pengiriman notifikasi - total {{ table.total }} data
              </p>
            </div>
          </div>
          <span :class="[channelBadgeClass(table.channel), 'px-3 py-1 rounded-full border text-[9px] font-black uppercase tracking-widest']">
            {{ channelLabel(table.channel) }}
          </span>
        </div>

        <div class="overflow-x-auto">
          <table class="w-full text-left">
            <thead>
              <tr class="bg-slate-50/60 border-b border-slate-100">
                <th class="px-8 py-4 text-[9px] font-black text-slate-400 uppercase tracking-widest">Alur</th>
                <th class="px-8 py-4 text-[9px] font-black text-slate-400 uppercase tracking-widest">Keterangan</th>
                <th class="px-8 py-4 text-[9px] font-black text-slate-400 uppercase tracking-widest text-center">Jumlah</th>
                <th class="px-8 py-4 text-[9px] font-black text-slate-400 uppercase tracking-widest text-center">Persentase</th>
                <th class="px-8 py-4 text-[9px] font-black text-slate-400 uppercase tracking-widest text-right">Detail</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-50">
              <tr v-for="row in table.rows" :key="table.channel + row.key" class="hover:bg-slate-50/70 transition-all">
                <td class="px-8 py-5">
                  <div class="flex items-center gap-3">
                    <div :class="[statusBadgeClass(row.key), 'w-9 h-9 rounded-xl border flex items-center justify-center']">
                      <component :is="row.icon" class="w-4 h-4" />
                    </div>
                    <div>
                      <p class="text-xs font-black text-slate-700 uppercase tracking-wider">{{ row.label }}</p>
                      <p class="text-[8px] font-bold text-slate-400 uppercase tracking-widest">{{ row.key }}</p>
                    </div>
                  </div>
                </td>
                <td class="px-8 py-5 text-[10px] font-bold text-slate-500 uppercase tracking-wider">
                  {{ row.note }}
                </td>
                <td class="px-8 py-5 text-center">
                  <span class="text-sm font-black text-slate-800">{{ row.count }}</span>
                </td>
                <td class="px-8 py-5 text-center">
                  <div class="flex items-center justify-center gap-3">
                    <div class="w-28 h-2 rounded-full bg-slate-100 overflow-hidden">
                      <div
                        class="h-full rounded-full bg-indigo-500"
                        :style="{ width: row.percent }"
                      ></div>
                    </div>
                    <span class="w-12 text-right text-[10px] font-black text-slate-500">{{ row.percent }}</span>
                  </div>
                </td>
                <td class="px-8 py-5 text-right">
                  <button
                    @click="fetchCommunicationDetails(row.key, table.channel)"
                    class="px-4 py-2 bg-white border border-slate-200 hover:border-indigo-200 hover:bg-indigo-50 hover:text-indigo-600 rounded-xl text-[9px] font-black text-slate-500 uppercase tracking-widest transition-all"
                  >
                    Lihat Pesan
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- 9. Quick Actions Panel -->
    <div class="mt-10 grid grid-cols-1 gap-8">

      <!-- Quick Navigation -->
      <div class="white-card p-8">
        <div class="flex items-center justify-between mb-6">
          <div>
            <h3 class="text-xs font-black text-slate-800 uppercase tracking-widest">Aksi Cepat</h3>
            <p class="text-[10px] font-bold text-slate-400 uppercase mt-0.5">Navigasi Administrasi</p>
          </div>
          <div class="w-10 h-10 bg-slate-50 text-slate-600 rounded-xl flex items-center justify-center">
            <LayoutIcon class="w-5 h-5" />
          </div>
        </div>

        <div class="grid grid-cols-2 gap-3">
          <button @click="router.push('/reports')" class="p-4 bg-slate-50 hover:bg-indigo-50 hover:text-indigo-700 text-slate-600 rounded-2xl transition-all group text-left border border-transparent hover:border-indigo-100">
            <ExcelIcon class="w-5 h-5 mb-2 text-emerald-500 group-hover:scale-110 transition-all" />
            <p class="text-[10px] font-black uppercase tracking-wider">Laporan</p>
            <p class="text-[9px] text-slate-400 font-medium mt-0.5">Ekspor & Analitik</p>
          </button>
          <button @click="router.push('/finance/rules')" class="p-4 bg-slate-50 hover:bg-amber-50 hover:text-amber-700 text-slate-600 rounded-2xl transition-all group text-left border border-transparent hover:border-amber-100">
            <BillIcon class="w-5 h-5 mb-2 text-amber-500 group-hover:scale-110 transition-all" />
            <p class="text-[10px] font-black uppercase tracking-wider">Aturan Tagihan</p>
            <p class="text-[9px] text-slate-400 font-medium mt-0.5">Buat & Kelola</p>
          </button>
          <button @click="router.push('/students')" class="p-4 bg-slate-50 hover:bg-purple-50 hover:text-purple-700 text-slate-600 rounded-2xl transition-all group text-left border border-transparent hover:border-purple-100">
            <StudentIcon class="w-5 h-5 mb-2 text-indigo-500 group-hover:scale-110 transition-all" />
            <p class="text-[10px] font-black uppercase tracking-wider">Kelola Siswa</p>
            <p class="text-[9px] text-slate-400 font-medium mt-0.5">Data Master</p>
          </button>
          <button @click="router.push('/users')" class="p-4 bg-slate-50 hover:bg-emerald-50 hover:text-emerald-700 text-slate-600 rounded-2xl transition-all group text-left border border-transparent hover:border-emerald-100">
            <UsersIcon class="w-5 h-5 mb-2 text-emerald-500 group-hover:scale-110 transition-all" />
            <p class="text-[10px] font-black uppercase tracking-wider">Kelola Pengguna</p>
            <p class="text-[9px] text-slate-400 font-medium mt-0.5">Akun & Wali Murid</p>
          </button>
        </div>
      </div>
    </div>

    <!-- Drill-down Modal -->
    <Teleport to="body">
      <div v-if="showDrillDown" class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-[999] flex items-center justify-center p-4">
         <div class="bg-white w-full max-w-6xl rounded-2xl overflow-hidden shadow-2xl animate-scale-up">
            <div class="p-8 border-b border-slate-50 flex items-center justify-between bg-slate-50/50">
               <div>
                 <h3 class="text-xs font-black text-slate-800 uppercase tracking-widest">{{ drillDownTitle }}</h3>
                 <p class="text-[10px] font-bold text-slate-400 uppercase tracking-wider mt-1">Detail pesan, penerima, status, dan alasan gagal bila ada</p>
               </div>
               <button @click="showDrillDown = false" class="w-8 h-8 rounded-full bg-white flex items-center justify-center text-slate-400 hover:text-rose-500 shadow-sm transition-all">
                  <ErrorIcon class="w-4 h-4" />
               </button>
            </div>
            <div class="max-h-[500px] overflow-y-auto p-4 custom-scrollbar">
               <div v-if="drillDownLoading" class="p-20 text-center text-slate-400 font-black uppercase text-[10px]">Memuat data...</div>
               <div v-else-if="drillDownData.length === 0" class="p-20 text-center">
                 <SendIcon class="w-10 h-10 text-slate-200 mx-auto mb-3" />
                 <p class="text-[10px] font-black text-slate-300 uppercase tracking-[0.2em]">Belum ada pesan dengan status ini</p>
               </div>
               <table v-else class="w-full text-left">
                  <thead>
                    <tr class="text-[9px] font-black text-slate-400 uppercase border-b border-slate-50">
                      <th class="px-4 py-3">Penerima</th>
                      <th class="px-4 py-3">Kontak</th>
                      <th class="px-4 py-3">Judul & Pesan</th>
                      <th class="px-4 py-3">Status</th>
                      <th class="px-4 py-3 text-right">Waktu</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-slate-50">
                    <tr v-for="d in drillDownData" :key="d.id" class="hover:bg-slate-50/50 transition-all">
                      <td class="px-4 py-4 max-w-[180px]">
                         <p class="text-xs font-bold text-slate-600 uppercase tracking-wider truncate">{{ d.student_name || 'Tanpa data siswa' }}</p>
                         <p class="text-[9px] font-bold text-slate-400 uppercase truncate">{{ d.recipient_name || 'Penerima' }}</p>
                      </td>
                      <td class="px-4 py-4 max-w-[170px]">
                         <span :class="[channelBadgeClass(d.channel), 'px-2 py-1 rounded text-[8px] font-black uppercase border']">{{ channelLabel(d.channel) }}</span>
                         <p class="mt-2 text-[10px] font-mono font-bold text-slate-500 truncate">
                           {{ d.channel === 'email' ? (d.recipient_email || '-') : (d.recipient_phone || '-') }}
                         </p>
                      </td>
                      <td class="px-4 py-4 max-w-[420px]">
                         <p class="text-[11px] font-black text-slate-700 uppercase tracking-wider truncate">{{ d.title }}</p>
                         <p class="mt-1 text-[10px] font-medium text-slate-500 line-clamp-2">{{ d.message }}</p>
                         <p v-if="d.delivery_error" class="mt-2 text-[10px] font-bold text-rose-600 bg-rose-50 border border-rose-100 rounded-lg px-2 py-1 line-clamp-2">
                           {{ d.delivery_error }}
                         </p>
                      </td>
                      <td class="px-4 py-4">
                         <span :class="[statusBadgeClass(d.delivery_status), 'px-2 py-1 rounded text-[8px] font-black uppercase border']">{{ statusLabel(d.delivery_status) }}</span>
                      </td>
                      <td class="px-4 py-4 text-right text-[9px] font-bold text-slate-400 uppercase">
                         {{ formatDateTime(d.updated_at || d.created_at) }}
                      </td>
                    </tr>
                  </tbody>
               </table>
            </div>
         </div>
      </div>
    </Teleport>


  </div>
</template>

<style scoped lang="postcss">
.white-card {
  @apply bg-white border border-slate-100 rounded-xl transition-all duration-300 shadow-sm overflow-hidden;
}
.quick-action-btn {
  @apply px-6 py-3 rounded-2xl text-[10px] font-black uppercase tracking-widest flex items-center gap-2 hover:scale-105 active:scale-95 transition-all;
}
.custom-scrollbar::-webkit-scrollbar { width: 3px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(226, 232, 240, 0.4); border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: rgba(203, 213, 225, 0.8); }
.no-scrollbar::-webkit-scrollbar { display: none; }
.no-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }
@keyframes scale-up {
  from { transform: scale(0.95); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}
.animate-scale-up { animation: scale-up 0.3s ease-out; }
</style>
