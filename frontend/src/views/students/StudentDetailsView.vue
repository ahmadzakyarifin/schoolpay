<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { 
  ArrowLeft as ArrowLeftIcon,
  User as StudentIcon,
  Users as ParentIcon,
  GraduationCap as GraduationIcon,
  Calendar as CalendarIcon,
  MapPin as MapPinIcon,
  CheckCircle2 as SuccessIcon,
  XCircle as ErrorIcon,
  Mail as MailIcon,
  Phone as PhoneIcon,
  BadgeInfo as InfoIcon,
  Fingerprint as IDIcon,
  Hash as HashIcon,
  Stethoscope as GenderIcon,
  Heart as ReligionIcon,
  Smartphone as WAIcon,
  CreditCard as BillingIcon,
  History as HistoryIcon,
  Clock as PendingIcon,
  ChevronLeft as PrevIcon,
  ChevronRight as NextIcon,
  ExternalLink as LinkIcon,
  AlertCircle as AlertIcon,
  Wallet as WalletIcon,
  Target as TargetIcon,
  Library as YearIcon,
  ChevronDown as ChevronDownIcon,
  UserCircle as UserProfileIcon,
  UserCheck as IDBadgeIcon
} from 'lucide-vue-next'
import studentService from '../../services/student.service'
import financeService from '../../services/finance.service'

const route = useRoute()
const router = useRouter()
const student = ref(null)
const parents = ref([])
const classHistory = ref([])
const bills = ref([])

const activeTab = ref('biodata')
const loading = ref(true)
const error = ref(null)

// Pagination & Filter State
const billsPage = ref(1)
const billsLimit = ref(5)
const selectedYear = ref('all')

// Academic History Pagination
const historyPage = ref(1)
const historyLimit = ref(3)

const apiBase = axios.defaults.baseURL
const staticBase = apiBase.replace('/api/', '')

// Helper to format Academic Year with Grade label
const formatAYWithGrade = (ayStr) => {
  if (!ayStr || !student.value) return ayStr
  
  const startYear = parseInt(ayStr.split('/')[0])
  const entryYear = student.value.entry_year
  const diff = startYear - entryYear
  
  // Grade calculation (X, XI, XII)
  let grade = "X"
  if (diff === 1) grade = "XI"
  else if (diff >= 2) grade = "XII"
  
  return `Kelas ${grade} (${ayStr})`
}

// Dropdown options from classHistory: { value: academic_year, label: class_name }
const availableYears = computed(() => {
  const seen = new Set()
  return classHistory.value
    .filter(h => h.academic_year && !seen.has(h.academic_year) && seen.add(h.academic_year))
    .map(h => ({ value: h.academic_year, label: h.class_name, academicYear: h.academic_year }))
    .sort((a, b) => a.academicYear.localeCompare(b.academicYear))
})

// Financial Summary Computed
const financialSummary = computed(() => {
  if (!bills.value.length) return { total: 0, paid: 0, unpaid: 0, percentage: 0, status: 'none' }
  
  const total = bills.value.reduce((acc, b) => acc + (Number(b.amount) || 0), 0)
  const paid = bills.value.reduce((acc, b) => acc + (Number(b.total_paid) || 0), 0)
  const unpaid = total - paid
  const percentage = total > 0 ? Math.round((paid / total) * 100) : 0
  
  // Derive overall status from bill statuses
  const allPaid = bills.value.every(b => b.status === 'paid')
  const anyOverdue = bills.value.some(b => b.status === 'overdue')
  const anyUnpaid = bills.value.some(b => b.status === 'unpaid' || b.status === 'partial')
  
  let status = 'none'
  if (allPaid) status = 'lunas'
  else if (anyOverdue) status = 'overdue'
  else if (anyUnpaid) status = 'berjalan'
  
  return { total, paid, unpaid, percentage, status }
})

// Paginated & Filtered Bills logic
const filteredBills = computed(() => {
  let list = [...bills.value]
  if (selectedYear.value !== 'all') {
    list = list.filter(b => b.academic_year === selectedYear.value)
  }
  
  return list.sort((a, b) => {
    if (a.status === 'unpaid' && b.status !== 'unpaid') return -1
    if (a.status !== 'unpaid' && b.status === 'unpaid') return 1
    return new Date(b.due_date) - new Date(a.due_date)
  })
})

const paginatedBills = computed(() => {
  const start = (billsPage.value - 1) * billsLimit.value
  const end = start + billsLimit.value
  return filteredBills.value.slice(start, end)
})

const billsTotalPages = computed(() => Math.ceil(filteredBills.value.length / billsLimit.value) || 1)

const visibleBillPages = computed(() => {
  const current = billsPage.value
  const last = billsTotalPages.value
  const pages = []
  if (last <= 3) {
    for (let i = 1; i <= last; i++) pages.push(i)
  } else {
    if (current === 1) pages.push(1, 2, 3)
    else if (current === last) pages.push(last - 2, last - 1, last)
    else pages.push(current - 1, current, current + 1)
  }
  return pages
})

// Academic History Pagination Logic
const paginatedHistory = computed(() => {
  const start = (historyPage.value - 1) * historyLimit.value
  const end = start + historyLimit.value
  return [...classHistory.value].sort((a, b) => new Date(b.created_at) - new Date(a.created_at)).slice(start, end)
})

const historyTotalPages = computed(() => Math.ceil(classHistory.value.length / historyLimit.value) || 1)

const visibleHistoryPages = computed(() => {
  const current = historyPage.value
  const last = historyTotalPages.value
  const pages = []
  if (last <= 3) {
    for (let i = 1; i <= last; i++) pages.push(i)
  } else {
    if (current === 1) pages.push(1, 2, 3)
    else if (current === last) pages.push(last - 2, last - 1, last)
    else pages.push(current - 1, current, current + 1)
  }
  return pages
})

// Reset page when filter changes
watch(selectedYear, () => {
  billsPage.value = 1
})

const fetchAllData = async () => {
  const studentId = route.params.id
  if (!studentId) {
    error.value = 'ID Siswa tidak ditemukan'
    loading.value = false
    return
  }

  loading.value = true
  error.value = null
  
  try {
    const [studentRes, parentRes, historyRes, billsRes] = await Promise.all([
      studentService.getByID(studentId),
      studentService.getParents(studentId),
      studentService.getClassHistory(studentId),
      financeService.getBillsByStudent(studentId)
    ])

    student.value = studentRes.data.data
    parents.value = parentRes.data.data || []
    classHistory.value = historyRes.data.data || []
    bills.value = billsRes.data.data || []

    if (classHistory.value.length === 0 && student.value && student.value.class_id) {
      classHistory.value = [{
        id: 'current',
        class_name: student.value.class_name,
        grade: student.value.grade || 'X', 
        is_active: true,
        created_at: student.value.created_at
      }]
    }

  } catch (err) {
    console.error('Error fetching student details:', err)
    error.value = 'Gagal memuat data lengkap siswa.'
  } finally {
    loading.value = false
  }
}

const handlePayManual = async (bill) => {
  const remaining = bill.amount - bill.total_paid
  const reason = prompt(`Catat pembayaran manual untuk tagihan "${bill.bill_type_name}" sebesar ${formatCurrency(remaining)}?\n\nMasukkan keterangan audit (wajib, misal: Tunai di koperasi, Beasiswa, Koreksi admin):`)
  if (reason === null) return // Canceled
  if (!reason.trim()) {
    alert('Alasan pelunasan manual wajib diisi untuk kebutuhan audit!')
    return
  }

  loading.value = true
  try {
    await axios.post(`finance/bills/${bill.id}/pay-manual`, { reason, note: reason, payment_method: 'Tunai' })
    alert('Pembayaran manual berhasil dicatat')
    await fetchAllData()
  } catch (err) {
    const errMsg = err.response?.data?.message || 'Gagal melunasi tagihan secara manual'
    alert(errMsg)
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push({ name: 'student-management' })
  }
}

const goToParent = (id) => {
  router.push({ name: 'user-details', params: { id } })
}

onMounted(() => {
  fetchAllData()
})

watch(() => route.params.id, (newId) => {
  if (newId) fetchAllData()
})

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(amount)
}

const formattedDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear()}`
}

const getStatusColor = (status) => {
  switch (status) {
    case 'paid': return 'bg-emerald-50 text-emerald-600 border-emerald-100'
    case 'unpaid': return 'bg-rose-50 text-rose-600 border-rose-100'
    case 'overdue': return 'bg-red-50 text-red-700 border-red-200'
    case 'partial': return 'bg-amber-50 text-amber-600 border-amber-100'
    case 'pending': return 'bg-amber-50 text-amber-600 border-amber-100'
    default: return 'bg-slate-50 text-slate-500 border-slate-100'
  }
}

const getStatusLabel = (status) => {
  switch (status) {
    case 'paid': return 'Lunas'
    case 'unpaid': return 'Menunggak'
    case 'overdue': return 'Jatuh Tempo'
    case 'partial': return 'Sebagian'
    case 'pending': return 'Proses'
    default: return status || '-'
  }
}
</script>

<template>
  <div class="min-h-screen bg-slate-50/50 p-6 md:p-10 font-inter">
    <!-- Header -->
    <div class="max-w-6xl mx-auto mb-10 flex flex-col md:flex-row md:items-center justify-between gap-6 animate-fade-in">
      <div class="flex items-center gap-6">
        <button @click="goBack" class="w-12 h-12 bg-white border border-slate-200 rounded-2xl flex items-center justify-center text-slate-400 hover:text-indigo-600 hover:border-indigo-100 hover:shadow-xl hover:shadow-indigo-50 transition-all group">
          <ArrowLeftIcon class="w-5 h-5 group-hover:-translate-x-1 transition-transform" />
        </button>
        <div>
          <div class="flex items-center gap-3 mb-1">
            <h1 class="text-3xl font-black text-slate-800 tracking-tight">Detail Siswa</h1>
            <span v-if="student" :class="[
              'px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-widest border',
              student.status === 'active' ? 'bg-emerald-50 text-emerald-600 border-emerald-100' : 'bg-rose-50 text-rose-600 border-rose-100'
            ]">
              {{ student.status === 'active' ? 'Aktif' : 'Non-Aktif' }}
            </span>
          </div>
          <p class="text-slate-400 font-bold text-xs uppercase tracking-[0.2em]">Sistem Informasi Akademik & Keuangan</p>
        </div>
      </div>
    </div>

    <!-- Loading & Error -->
    <div v-if="loading" class="max-w-6xl mx-auto py-20 flex flex-col items-center justify-center">
      <div class="w-12 h-12 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
      <p class="text-slate-400 font-black text-[10px] uppercase tracking-widest">Sinkronisasi Data...</p>
    </div>

    <div v-else-if="error" class="max-w-6xl mx-auto py-20 flex flex-col items-center justify-center text-center">
      <div class="w-20 h-20 bg-rose-50 rounded-3xl flex items-center justify-center mb-6 border border-rose-100 shadow-inner">
        <ErrorIcon class="w-10 h-10 text-rose-500" />
      </div>
      <h3 class="text-xl font-black text-slate-800 mb-2">{{ error }}</h3>
      <button @click="fetchAllData" class="text-indigo-600 font-bold hover:underline">Coba Lagi</button>
    </div>

    <!-- Content -->
    <div v-else-if="student" class="max-w-6xl mx-auto space-y-8 animate-slide-up">
      <!-- Top Financial Summary Cards -->
      <div class="grid grid-cols-1 md:grid-cols-5 gap-6">
        <div class="white-card p-4 rounded-xl flex items-center gap-4 shadow-sm border-slate-100">
           <div class="w-10 h-10 bg-indigo-50 text-indigo-600 rounded-lg flex items-center justify-center"><WalletIcon class="w-5 h-5" /></div>
           <div><p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Total Tagihan</p><p class="text-xs font-black text-slate-700">{{ formatCurrency(financialSummary.total) }}</p></div>
        </div>
        <div class="white-card p-4 rounded-xl flex items-center gap-4 shadow-sm border-slate-100">
           <div class="w-10 h-10 bg-emerald-50 text-emerald-600 rounded-lg flex items-center justify-center"><SuccessIcon class="w-5 h-5" /></div>
           <div><p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Sudah Dibayar</p><p class="text-xs font-black text-slate-700">{{ formatCurrency(financialSummary.paid) }}</p></div>
        </div>
        <div class="white-card p-4 rounded-xl flex items-center gap-4 shadow-sm border-slate-100">
           <div :class="['w-10 h-10 rounded-lg flex items-center justify-center', financialSummary.unpaid > 0 ? 'bg-rose-50 text-rose-600' : 'bg-slate-50 text-slate-400']"><AlertIcon class="w-5 h-5" /></div>
           <div><p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Sisa Tunggakan</p><p :class="['text-xs font-black', financialSummary.unpaid > 0 ? 'text-rose-600' : 'text-slate-700']">{{ formatCurrency(financialSummary.unpaid) }}</p></div>
        </div>
        <div class="white-card p-4 rounded-xl flex items-center gap-4 shadow-sm border-slate-100">
           <div class="w-10 h-10 bg-teal-50 text-teal-600 rounded-lg flex items-center justify-center"><WalletIcon class="w-5 h-5" /></div>
           <div><p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Saldo Deposit</p><p class="text-xs font-black text-teal-600">{{ formatCurrency(student.deposit_balance || 0) }}</p></div>
        </div>
        <div class="white-card p-4 rounded-xl flex items-center gap-4 shadow-sm border-slate-100">
           <div class="flex-1">
             <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1">Status Pelunasan</p>
             <p class="text-xs font-black text-slate-700 mb-1.5">{{ financialSummary.percentage }}% Terbayar</p>
             <span :class="[
               'px-2 py-0.5 rounded text-[8px] font-black uppercase tracking-widest border',
               financialSummary.status === 'lunas' ? 'bg-emerald-50 text-emerald-600 border-emerald-100' :
               financialSummary.status === 'overdue' ? 'bg-red-50 text-red-700 border-red-200' :
               financialSummary.status === 'berjalan' ? 'bg-amber-50 text-amber-600 border-amber-100' :
               'bg-slate-50 text-slate-400 border-slate-100'
             ]">{{ financialSummary.status === 'lunas' ? 'Lunas' : financialSummary.status === 'overdue' ? 'Jatuh Tempo' : financialSummary.status === 'berjalan' ? 'Berjalan' : 'No Data' }}</span>
           </div>
        </div>
      </div>

      <!-- Main Unified Interface -->
      <div class="space-y-6">
        <!-- Tab Navigation -->
        <div class="flex items-center p-1 bg-white border border-slate-100 rounded-xl shadow-sm overflow-x-auto no-scrollbar">
          <button @click="activeTab = 'biodata'" :class="['flex-1 px-6 py-2.5 rounded-lg text-[10px] font-black uppercase tracking-widest transition-all flex items-center justify-center gap-2 whitespace-nowrap', activeTab === 'biodata' ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/20' : 'text-slate-400 hover:text-indigo-600']">
            <StudentIcon class="w-3.5 h-3.5" /> Biodata & Wali
          </button>
          <button @click="activeTab = 'akademik'" :class="['flex-1 px-6 py-2.5 rounded-lg text-[10px] font-black uppercase tracking-widest transition-all flex items-center justify-center gap-2 whitespace-nowrap', activeTab === 'akademik' ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/20' : 'text-slate-400 hover:text-indigo-600']">
            <GraduationIcon class="w-3.5 h-3.5" /> Akademik
          </button>
          <button @click="activeTab = 'keuangan'" :class="['flex-1 px-6 py-2.5 rounded-lg text-[10px] font-black uppercase tracking-widest transition-all flex items-center justify-center gap-2 whitespace-nowrap', activeTab === 'keuangan' ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/20' : 'text-slate-400 hover:text-indigo-600']">
            <WalletIcon class="w-3.5 h-3.5" /> Keuangan
          </button>
        </div>

        <!-- Tab Content Containers -->
        <div class="transition-all duration-500">
          <!-- BIODATA TAB -->
          <div v-if="activeTab === 'biodata'" class="white-card p-6 rounded-xl shadow-sm border-slate-100 min-h-[650px] animate-fade-in">
            <div class="flex flex-col md:flex-row items-center gap-10 mb-12 p-8 bg-slate-50/50 rounded-2xl border border-slate-100">
               <div class="w-44 h-44 rounded-2xl bg-white border-8 border-white shadow-xl overflow-hidden group">
                  <img v-if="student.image_path" :src="`${staticBase}/${student.image_path}`" class="w-full h-full object-cover transition-transform group-hover:scale-105 duration-700" />
                  <div v-else class="w-full h-full flex items-center justify-center bg-slate-50 text-slate-200"><StudentIcon class="w-20 h-20" /></div>
               </div>
               <div class="flex-1 text-center md:text-left">
                  <h2 class="text-4xl font-black text-slate-800 tracking-tight leading-tight uppercase mb-4">{{ student.name }}</h2>
                  <div class="flex flex-wrap justify-center md:justify-start gap-3">
                    <span class="px-4 py-1.5 bg-indigo-600 text-white rounded-xl text-[10px] font-black uppercase tracking-widest">{{ student.class_name || 'Tanpa Kelas' }}</span>
                    <span class="px-4 py-1.5 bg-white text-slate-600 rounded-xl text-[10px] font-black uppercase tracking-widest border border-slate-200">Angkatan {{ student.entry_year }}</span>
                    <span v-if="student.major_name" class="px-4 py-1.5 bg-amber-500 text-white rounded-xl text-[10px] font-black uppercase tracking-widest">{{ student.major_name }}</span>
                  </div>
               </div>
               <div class="hidden lg:grid grid-cols-1 gap-3">
                  <div class="p-4 bg-white rounded-2xl border border-slate-100 shadow-sm flex items-center gap-4 min-w-[200px]">
                    <HashIcon class="w-4 h-4 text-slate-300" />
                    <div><p class="text-[8px] font-black text-slate-400 uppercase tracking-widest">NISN</p><p class="text-xs font-black text-slate-700">{{ student.nisn || '-' }}</p></div>
                  </div>
                  <div class="p-4 bg-white rounded-2xl border border-slate-100 shadow-sm flex items-center gap-4 min-w-[200px]">
                    <IDBadgeIcon class="w-4 h-4 text-slate-300" />
                    <div><p class="text-[8px] font-black text-slate-400 uppercase tracking-widest">NIS</p><p class="text-xs font-black text-slate-700">{{ student.nis || '-' }}</p></div>
                  </div>
               </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-12">
               <div class="p-4 bg-slate-50 border border-slate-100 rounded-xl">
                 <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5">NIK Siswa</p>
                 <p class="text-xs font-black text-slate-700 tracking-wider">{{ student.nik || '-' }}</p>
               </div>
               <div class="p-4 bg-slate-50 border border-slate-100 rounded-xl">
                 <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5">Jenis Kelamin</p>
                 <p class="text-xs font-black text-slate-700 uppercase">{{ student.gender || '-' }}</p>
               </div>
               <div class="p-4 bg-slate-50 border border-slate-100 rounded-xl">
                 <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5">Tempat, Tanggal Lahir</p>
                 <p class="text-xs font-black text-slate-700 uppercase">{{ student.birth_place || '-' }}, {{ formattedDate(student.birth_date) }}</p>
               </div>
               <div class="p-4 bg-slate-50 border border-slate-100 rounded-xl">
                 <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5">Agama</p>
                 <p class="text-xs font-black text-slate-700 uppercase">{{ student.religion || '-' }}</p>
               </div>
               <div class="p-4 bg-slate-50 border border-slate-100 rounded-xl">
                 <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5">Kontak Siswa</p>
                 <p class="text-xs font-black text-slate-700">{{ student.phone_number || '-' }}</p>
               </div>
               <div class="p-4 bg-slate-50 border border-slate-100 rounded-xl">
                 <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5">Email Siswa</p>
                 <p class="text-xs font-black text-slate-700">{{ student.email || '-' }}</p>
               </div>
               <div class="md:col-span-2 lg:col-span-3 p-4 bg-slate-50 border border-slate-100 rounded-xl">
                 <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5">Alamat Lengkap</p>
                 <p class="text-xs font-bold text-slate-700 leading-relaxed capitalize break-words">{{ student.address || '-' }}</p>
               </div>
               <div class="md:col-span-2 lg:col-span-3 p-4 bg-slate-50 border border-slate-100 rounded-xl">
                 <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mb-1.5">Keterangan / Catatan</p>
                 <p class="text-xs font-bold text-slate-700 leading-relaxed break-words">{{ student.description || 'Tidak ada catatan / keterangan' }}</p>
               </div>
            </div>

            <div class="w-full h-px bg-slate-100 my-12"></div>

            <div class="flex items-center gap-4 mb-8"><div class="p-3 bg-emerald-50 text-emerald-600 rounded-xl"><ParentIcon class="w-5 h-5" /></div><h3 class="text-lg font-black text-slate-800 tracking-tight uppercase">Orang Tua / Wali Murid</h3></div>
            <div v-if="parents.length" class="grid grid-cols-1 md:grid-cols-2 gap-6">
               <div v-for="p in parents" :key="p.id" class="p-8 bg-slate-50 border border-slate-100 rounded-xl flex items-center justify-between group hover:border-emerald-200 transition-all">
                  <div class="flex items-center gap-5">
                    <div class="w-14 h-14 bg-white rounded-xl flex items-center justify-center text-slate-300 group-hover:text-emerald-600 transition-colors shadow-sm"><UserProfileIcon class="w-7 h-7" /></div>
                    <div><p class="text-sm font-black text-slate-700 uppercase">{{ p.name }}</p><p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">{{ p.relation || 'Wali' }}</p></div>
                  </div>
                  <button @click="goToParent(p.id)" class="px-6 py-3 bg-white text-[10px] font-black text-slate-500 uppercase tracking-widest rounded-xl border border-slate-200 hover:bg-indigo-600 hover:text-white hover:border-indigo-600 transition-all shadow-sm flex items-center gap-2"><UserProfileIcon class="w-4 h-4" /> Lihat Profil</button>
               </div>
            </div>
          </div>

          <!-- AKADEMIK TAB (PAGINATED & SMALLER UI) -->
          <div v-if="activeTab === 'akademik'" class="white-card rounded-xl shadow-sm border-slate-100 min-h-[650px] flex flex-col animate-fade-in overflow-hidden">
             <div class="p-4 border-b border-slate-100 flex items-center gap-3 bg-slate-50/20">
                <div class="p-2 bg-indigo-50 text-indigo-600 rounded-lg shadow-sm"><HistoryIcon class="w-4 h-4" /></div>
                <div><h3 class="text-xs font-black text-slate-800 uppercase tracking-widest">Riwayat Progres Akademik</h3><p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Rekaman kenaikan dan perpindahan kelas</p></div>
             </div>
             
             <div class="flex-1 p-6 relative pl-16">
                <div v-if="classHistory.length" class="absolute left-8 top-8 bottom-8 w-1 bg-slate-100 rounded-full"></div>
                
                <div class="space-y-8">
                   <div v-for="(history, index) in paginatedHistory" :key="history.id" class="relative group">
                      <!-- Smaller Grade Circle -->
                      <div :class="['absolute -left-12 w-8 h-8 rounded-lg z-10 flex items-center justify-center border-4 border-white shadow-md transition-all duration-700', history.is_active ? 'bg-indigo-600 text-white' : 'bg-slate-100 text-slate-400']">
                        <span class="text-[9px] font-black">{{ classHistory.length - ((historyPage - 1) * historyLimit + index) }}</span>
                      </div>
                      
                      <div class="flex-1 pb-4 border-b border-slate-100 last:border-0 group-hover:translate-x-1.5 transition-transform duration-500">
                         <div class="flex items-center justify-between mb-1.5">
                            <h4 class="text-xs font-black text-slate-700 uppercase tracking-tight">{{ history.class_name }}</h4>
                            <span v-if="history.is_active" class="px-2 py-0.5 bg-emerald-100 text-emerald-700 rounded-md text-[8px] font-black uppercase tracking-widest">Aktif Sekarang</span>
                         </div>
                         <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Ditempatkan pada: {{ formattedDate(history.created_at) }}</p>
                      </div>
                   </div>
                </div>

                <!-- Empty state handled by flex-1 -->
                <div v-if="!classHistory.length" class="h-full flex items-center justify-center italic text-slate-300 text-xs font-bold uppercase tracking-widest">Belum ada riwayat akademik</div>
             </div>

             <!-- History Pagination Footer -->
             <div class="px-6 py-4 border-t border-slate-100 flex items-center justify-between bg-slate-50/30">
                <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Halaman {{ historyPage }} dari {{ historyTotalPages }}</p>
                <div class="flex items-center gap-2">
                   <button @click="historyPage--" :disabled="historyPage <= 1" class="w-8 h-8 border border-slate-200 rounded-lg hover:bg-white hover:shadow-md disabled:opacity-20 disabled:pointer-events-none transition-all flex items-center justify-center cursor-pointer"><PrevIcon class="w-3.5 h-3.5 text-slate-400" /></button>
                   <button v-for="p in visibleHistoryPages" :key="p" @click="historyPage = p" :class="[p === historyPage ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/20' : 'bg-white text-slate-400 border border-slate-200 hover:border-indigo-300']" class="w-8 h-8 rounded-lg text-[10px] font-black transition-all flex items-center justify-center cursor-pointer">{{ p }}</button>
                   <button @click="historyPage++" :disabled="historyPage >= historyTotalPages" class="w-8 h-8 border border-slate-200 rounded-lg hover:bg-white hover:shadow-md disabled:opacity-20 disabled:pointer-events-none transition-all flex items-center justify-center cursor-pointer"><NextIcon class="w-3.5 h-3.5 text-slate-400" /></button>
                </div>
             </div>
          </div>

          <!-- KEUANGAN TAB -->
          <div v-if="activeTab === 'keuangan'" class="white-card rounded-xl shadow-sm border-slate-100 min-h-[650px] flex flex-col overflow-hidden animate-fade-in">
             <div class="p-4 border-b border-slate-100 flex flex-col md:flex-row md:items-center justify-between gap-4 bg-slate-50/20">
                <div class="flex items-center gap-3"><div class="p-2.5 bg-indigo-50 text-indigo-600 rounded-lg shadow-sm"><WalletIcon class="w-4 h-4" /></div><div><h3 class="text-xs font-black text-slate-800 uppercase tracking-widest">Daftar Tagihan Pendidikan</h3><p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Kelola semua kewajiban finansial</p></div></div>
                <!-- Filter Dropdown -->
                <div class="relative min-w-[240px]">
                  <select v-model="selectedYear" class="w-full pl-9 pr-9 py-2 bg-slate-50 border border-slate-200 rounded-lg appearance-none focus:bg-white focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 outline-none transition-all font-black text-[10px] text-slate-700 uppercase tracking-widest cursor-pointer shadow-sm">
                    <option value="all">Semua Kelas</option>
                    <option v-for="opt in availableYears" :key="opt.value" :value="opt.value">{{ opt.label }} ({{ opt.academicYear }})</option>
                  </select>
                  <YearIcon class="absolute left-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-slate-400" /><ChevronDownIcon class="absolute right-3 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-slate-400 pointer-events-none" />
                </div>
             </div>
             <div class="flex-1 overflow-x-auto">
                <table class="w-full text-left border-collapse">
                    <thead>
                      <tr class="bg-slate-50/50 border-b border-slate-100 text-[10px] font-black text-slate-400 uppercase tracking-widest">
                        <th class="py-3 px-4">Tahun Ajaran</th>
                        <th class="py-3 px-4">Nama Tagihan</th>
                        <th class="py-3 px-4 text-right">Nominal</th>
                        <th class="py-3 px-4 text-center">Status</th>
                        <th class="py-3 px-4">Jatuh Tempo</th>
                        <th class="py-3 px-4 text-center w-[120px]">Aksi</th>
                      </tr>
                    </thead>
                    <tbody>
                       <tr v-for="bill in paginatedBills" :key="bill.id" class="border-b border-slate-100 hover:bg-slate-50/30 transition-all text-xs font-semibold text-slate-600">
                          <td class="py-3 px-4"><p class="text-[11px] font-black text-slate-700 uppercase">{{ bill.academic_year || '-' }}</p></td>
                          <td class="py-3 px-4"><p class="text-xs font-bold text-slate-500 uppercase tracking-wider truncate">{{ bill.bill_type_name }}</p><p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">Periode: {{ bill.period }}</p></td>
                          <td class="py-3 px-4 text-right"><span class="text-xs font-black text-slate-800">{{ formatCurrency(bill.amount) }}</span></td>
                          <td class="py-3 px-4 text-center"><span :class="['px-2 py-0.5 rounded text-[8px] font-black uppercase tracking-widest border', getStatusColor(bill.status)]">{{ getStatusLabel(bill.status) }}</span></td>
                          <td class="py-3 px-4"><span class="text-[11px] font-bold text-slate-400 uppercase tracking-tight">{{ formattedDate(bill.due_date) }}</span></td>
                          <td class="py-3 px-4 text-center">
                             <button v-if="bill.status !== 'paid'" @click="handlePayManual(bill)" class="px-2.5 py-1.5 bg-white text-emerald-600 border border-slate-200 hover:bg-slate-50 font-bold rounded-lg text-[9px] uppercase tracking-wider flex items-center justify-center gap-1 transition-all shadow-sm whitespace-nowrap cursor-pointer mx-auto">Bayar</button>
                             <span v-else class="text-[9px] font-black text-slate-300 uppercase">-</span>
                          </td>
                       </tr>
                       <tr v-for="i in (billsLimit - paginatedBills.length)" :key="'empty'+i" class="h-[49px]"><td colspan="6"></td></tr>
                    </tbody>
                 </table>
             </div>
             <div class="px-6 py-4 border-t border-slate-100 flex items-center justify-between bg-slate-50/30">
                <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Menampilkan {{ paginatedBills.length }} dari {{ filteredBills.length }} data</p>
                <div class="flex items-center gap-2">
                   <button @click="billsPage--" :disabled="billsPage <= 1" class="w-8 h-8 border border-slate-200 rounded-lg hover:bg-white hover:shadow-md disabled:opacity-20 disabled:pointer-events-none transition-all flex items-center justify-center cursor-pointer"><PrevIcon class="w-3.5 h-3.5 text-slate-400" /></button>
                   <button v-for="p in visibleBillPages" :key="p" @click="billsPage = p" :class="[p === billsPage ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/20' : 'bg-white text-slate-400 border border-slate-200 hover:border-indigo-300']" class="w-8 h-8 rounded-lg text-[10px] font-black transition-all flex items-center justify-center cursor-pointer">{{ p }}</button>
                   <button @click="billsPage++" :disabled="billsPage >= billsTotalPages" class="w-8 h-8 border border-slate-200 rounded-lg hover:bg-white hover:shadow-md disabled:opacity-20 disabled:pointer-events-none transition-all flex items-center justify-center cursor-pointer"><NextIcon class="w-3.5 h-3.5 text-slate-400" /></button>
                </div>
             </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="postcss">
.animate-fade-in { animation: fadeIn 0.4s cubic-bezier(0.16, 1, 0.3, 1) forwards; }
.animate-slide-up { animation: slideUp 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
@keyframes slideUp { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
.white-card { @apply bg-white border border-slate-100 transition-all; }
.no-scrollbar::-webkit-scrollbar { display: none; }
.no-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }
</style>
