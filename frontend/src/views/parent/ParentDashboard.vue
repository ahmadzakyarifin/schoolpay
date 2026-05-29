<script setup>
import { ref, onMounted, computed, watch, reactive } from 'vue'
import axios from 'axios'
import { useAuthStore } from '../../store/auth'
import { 
  User as UserIcon, 
  CreditCard as BillIcon, 
  Calendar as CalendarIcon,
  ChevronRight as ChevronRightIcon,
  ChevronDown as ChevronDownIcon,
  CheckCircle2 as PaidIcon,
  AlertCircle as AlertIcon,
  Search as SearchIcon,
  Filter as FilterIcon,
  RotateCcw as ResetIcon,
  ChevronLeft as PrevIcon,
  ChevronRight as NextIcon,
  Receipt as ReceiptIcon
} from 'lucide-vue-next'

const authStore = useAuthStore()
const students = ref([])
const selectedStudent = ref(null)
const allBills = ref([])
const loading = ref(false)
const payingBillId = ref(null)
const showStudentDropdown = ref(false)

// Filter & Pagination State
const search = ref('')
const statusFilter = ref('')
const page = ref(1)
const limit = ref(10)
const isMounted = ref(false)
const showFilters = ref(false)

const fetchParentData = async () => {
  loading.value = true
  try {
    const res = await axios.get('parent/students/me') 
    students.value = res.data.data
    
    if (students.value.length > 0) {
      selectedStudent.value = students.value[0]
      await fetchBills()
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
  } catch (err) {
    console.error('Gagal mengambil data tagihan')
  } finally {
    loading.value = false
  }
}

const selectStudent = (student) => {
  selectedStudent.value = student
  showStudentDropdown.value = false
  page.value = 1
}

const resetFilters = () => {
  search.value = ''
  statusFilter.value = ''
  page.value = 1
}

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

  // Status filter
  if (statusFilter.value) {
    filtered = filtered.filter(b => b.status === statusFilter.value)
  }

  return filtered
})

const paginatedBills = computed(() => {
  const start = (page.value - 1) * limit.value
  const end = start + limit.value
  return filteredBills.value.slice(start, end)
})

const totalPages = computed(() => Math.ceil(filteredBills.value.length / limit.value) || 1)
const totalData = computed(() => filteredBills.value.length)

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

const isOverdue = (bill) => {
  if (!bill?.due_date) return false
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const due = new Date(bill.due_date)
  due.setHours(0, 0, 0, 0)
  return due < today
}

const payWithMidtrans = async (bill) => {
  if (isOverdue(bill)) {
    alert('Tagihan sudah jatuh tempo. Pembayaran online ditutup, silakan bayar langsung ke admin/kasir sekolah secara cash atau transfer manual.')
    return
  }

  const remaining = Number(bill.amount || 0) - Number(bill.total_paid || 0)
  if (remaining <= 0) return

  payingBillId.value = bill.id
  try {
    const res = await axios.post('finance/payment-intent', {
      student_id: bill.student_id,
      amount: remaining,
      bill_ids: [bill.id],
      is_bypass_rule: false
    })
    const payment = res.data?.data
    if (payment?.snap_token && window.snap) {
      window.snap.pay(payment.snap_token, {
        onSuccess: async () => {
          alert('Pembayaran diterima Midtrans. Status akan diverifikasi otomatis oleh sistem.')
          await fetchBills()
        },
        onPending: async () => {
          alert('Pembayaran masih menunggu penyelesaian.')
          await fetchBills()
        },
        onError: () => alert('Pembayaran gagal diproses oleh Midtrans.'),
        onClose: () => alert('Anda menutup halaman pembayaran.')
      })
    } else if (payment?.payment_url) {
      window.location.href = payment.payment_url
    } else {
      alert('Link pembayaran belum tersedia. Silakan coba lagi.')
    }
  } catch (err) {
    alert(err.response?.data?.message || 'Gagal membuat pembayaran Midtrans')
  } finally {
    payingBillId.value = null
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear()}`
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
        <div class="flex items-center justify-center gap-2 flex-1 max-w-2xl mx-auto">
          
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
                  <span class="text-[10px] font-medium text-slate-400">NISN: {{ s.nisn || '-' }}</span>
                </button>
              </div>
            </transition>
          </div>

          <div class="relative flex-1 group">
            <SearchIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300 group-focus-within:text-indigo-600" />
            <input v-model="search" type="text" placeholder="Cari tagihan..." class="search-input-premium" />
          </div>
          
          <div class="relative">
            <select v-model="statusFilter" class="appearance-none p-2.5 pr-8 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm transition-all text-[10px] font-black uppercase tracking-wider cursor-pointer focus:outline-none focus:border-indigo-500">
              <option value="">Semua Status</option>
              <option value="unpaid">Belum Lunas</option>
              <option value="partial">Sebagian</option>
              <option value="paid">Lunas</option>
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
        <p class="text-[11px] text-slate-500 font-bold mt-1">Pantau dan kelola tagihan pembayaran putra-putri Anda.</p>
      </div>
      <div v-if="selectedStudent" class="bg-indigo-50 px-4 py-2.5 rounded-xl border border-indigo-100 flex items-center gap-3">
         <div class="w-8 h-8 rounded-lg bg-indigo-600 text-white flex items-center justify-center shrink-0">
           <UserIcon class="w-4 h-4" />
         </div>
         <div>
           <p class="text-[9px] font-black uppercase tracking-widest text-indigo-400">Siswa Terpilih</p>
           <p class="text-xs font-black text-indigo-700 leading-tight">{{ selectedStudent.name }}</p>
           <p class="text-[9px] font-bold text-emerald-600 uppercase tracking-widest mt-0.5">Saldo: Rp {{ (selectedStudent.deposit_balance || 0).toLocaleString('id-ID') }}</p>
         </div>
      </div>
    </div>

    <!-- Main Content Table -->
    <div class="bg-white rounded-[2rem] border border-slate-200 shadow-sm flex flex-col min-h-[710px] overflow-hidden">
      <div class="px-6 py-6 border-b border-slate-100 bg-slate-50/30 flex items-center justify-between">
        <div class="flex items-center gap-4">
          <div class="w-2 h-6 bg-indigo-500 rounded-full"></div>
          <h3 class="font-black text-slate-700 text-sm uppercase tracking-[0.2em]">Daftar Tagihan Siswa</h3>
        </div>
        <div class="flex items-center gap-3">
           <button class="bg-indigo-50 text-indigo-600 hover:bg-indigo-100 font-bold py-2 px-4 rounded-xl text-[10px] flex items-center gap-2 transition-all shadow-sm">
             <ReceiptIcon class="w-3.5 h-3.5" />
             <span>Riwayat Pembayaran</span>
           </button>
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
            <tr>
              <th class="w-16">No</th>
              <th>Detail Tagihan</th>
              <th>Periode</th>
              <th>Jatuh Tempo</th>
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
                <h3 class="text-sm font-black text-slate-700 uppercase tracking-wider mb-1">Tidak Ada Tagihan</h3>
                <p class="text-xs text-slate-400 font-medium">Belum ada data tagihan yang sesuai dengan pencarian Anda.</p>
              </td>
            </tr>
            <tr v-for="(b, i) in paginatedBills" :key="b.id" class="hover:bg-slate-50/50 transition-colors group">
              <td class="text-center text-xs font-bold text-slate-400">
                {{ (page - 1) * limit + i + 1 }}
              </td>
              <td>
                <div>
                  <p class="text-xs font-black text-slate-800 group-hover:text-indigo-600 transition-colors leading-tight">
                    {{ b.name || b.bill_type_name }}
                  </p>
                  <p class="text-[10px] font-bold text-slate-400 mt-1 truncate max-w-[200px]">
                    {{ b.description || 'Tidak ada deskripsi' }}
                  </p>
                </div>
              </td>
              <td>
                <span class="px-2.5 py-1 bg-slate-100 text-slate-600 rounded-lg text-[10px] font-black uppercase tracking-wider">
                  {{ b.period || b.academic_year || '-' }}
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
                  <p class="text-xs font-black text-slate-800">Rp {{ b.amount.toLocaleString('id-ID') }}</p>
                  <p v-if="b.total_paid > 0 && b.status !== 'paid'" class="text-[9px] font-black text-emerald-600 uppercase tracking-widest">
                    Dibayar: Rp {{ b.total_paid.toLocaleString('id-ID') }}
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
                 <button
                  v-if="b.status !== 'paid'"
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
                <button v-else class="w-full bg-white border border-slate-200 text-slate-500 hover:bg-slate-50 hover:text-indigo-600 font-black py-2.5 px-3 rounded-xl text-[10px] uppercase tracking-widest transition-all flex items-center justify-center gap-1.5 shadow-sm">
                  <ReceiptIcon class="w-3.5 h-3.5" />
                  Kwitansi
                </button>
              </td>
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
