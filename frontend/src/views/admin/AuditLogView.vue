<template>
  <div class="space-y-8 animate-fade-in">
    <!-- Header Teleport Target Content -->
    <Teleport to="#header-actions-target" v-if="isMounted">
      <div class="flex items-center justify-center w-full gap-4 relative mx-auto">
        <div class="flex items-center justify-center gap-2 flex-1 max-w-2xl mx-auto">
          <div class="relative flex-1 group">
            <SearchIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300 group-focus-within:text-indigo-600" />
            <input v-model="filters.user_name" @input="onSearchInput" type="text" placeholder="Cari nama admin, entitas, atau aksi..." class="search-input-premium" />
          </div>
          <button @click="showFilters = !showFilters" class="relative p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 flex items-center gap-2 transition-all shadow-sm">
            <FilterIcon class="w-4 h-4" />
            <span class="text-[10px] font-black uppercase tracking-wider pr-1">Filter</span>
            <span v-if="filters.action || filters.entity_type || filters.role || filters.sort" class="absolute -top-1 -right-1 w-3 h-3 bg-indigo-600 rounded-full border-2 border-white shadow-sm"></span>
          </button>
          <button @click="resetFilters" class="p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm transition-all group" title="Reset">
            <ResetIcon class="w-4 h-4 group-hover:rotate-180 transition-transform duration-500" />
          </button>
        </div>

        <AuditFilter 
          v-model="showFilters" 
          :filters="filters" 
          @apply="applyFilters" 
          @reset="resetFilters" 
        />
      </div>
    </Teleport>

    <!-- Main Table Card -->
    <div class="bg-white rounded border border-slate-200 shadow-sm flex flex-col min-h-[710px] overflow-hidden">
      <div class="px-6 py-6 border-b border-slate-100 bg-slate-50/30 flex items-center justify-between">
        <div class="flex items-center gap-4">
          <div class="w-2 h-6 bg-indigo-500 rounded-full"></div>
          <h3 class="font-black text-slate-700 text-sm uppercase tracking-[0.2em]">Riwayat Aktivitas & Log Audit</h3>
        </div>
        <div class="flex items-center gap-3">
          <span class="text-xs font-bold text-slate-400">Menampilkan {{ logs.length }} dari {{ totalData }} log</span>
        </div>
      </div>

      <!-- Table -->
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-slate-100 text-[10px] font-black text-slate-400 uppercase tracking-widest">
              <th class="py-4 px-4"># ID</th>
              <th class="py-4 px-4">Waktu (WIB)</th>
              <th class="py-4 px-4">Pengguna / Admin</th>
              <th class="py-4 px-4">Role</th>
              <th class="py-4 px-4">Aksi</th>
              <th class="py-4 px-4">Entitas (Target ID)</th>
              <th class="py-4 px-4">IP Address</th>
              <th class="py-4 px-4 text-right">Detail</th>
            </tr>
          </thead>
          <tbody class="text-xs font-medium divide-y divide-slate-100/80">
            <tr v-if="loading" class="animate-pulse">
              <td colspan="8" class="py-12 text-center text-slate-400 font-bold">Memuat data log audit...</td>
            </tr>
            <tr v-else-if="logs.length === 0">
              <td colspan="8" class="py-20 px-6 text-center">
                <div class="flex flex-col items-center justify-center text-center mx-auto">
                  <div class="w-20 h-20 bg-slate-100 rounded-[2.5rem] flex items-center justify-center text-slate-300 mb-6 border-4 border-white shadow-xl shadow-slate-200/50 mx-auto">
                    <DatabaseIcon class="w-10 h-10" />
                  </div>
                  <h3 class="text-lg font-black text-slate-700 tracking-tight mb-2">Riwayat Aktivitas Tidak Ditemukan</h3>
                  <p class="text-slate-400 text-xs font-medium max-w-xs mx-auto">Belum ada log audit yang terekam atau coba sesuaikan filter pencarian Anda.</p>
                </div>
              </td>
            </tr>
            <tr v-else v-for="log in logs" :key="log.id" class="hover:bg-slate-50/80 transition-colors group">
              <td class="py-4 px-4 font-black text-slate-700">#{{ log.id }}</td>
              <td class="py-4 px-4 text-slate-600">{{ formatDate(log.created_at) }}</td>
              <td class="py-4 px-4 font-bold text-slate-500 text-xs uppercase tracking-wider truncate max-w-[180px]">{{ log.user_name }}</td>
              <td class="py-4 px-4">
                <span class="px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-widest bg-slate-100 text-slate-600 border border-slate-200/60">
                  {{ log.role }}
                </span>
              </td>
              <td class="py-4 px-4 font-black text-[10px] uppercase tracking-widest text-slate-700">
                {{ humanizeAction(log.action) }}
              </td>
              <td class="py-4 px-4 font-bold text-slate-700">
                {{ humanizeEntity(log.entity_type) }} <span class="text-slate-400 font-normal">(ID: {{ log.entity_id }})</span>
              </td>
              <td class="py-4 px-4 text-slate-500 font-mono text-[11px]">{{ log.ip_address || '-' }}</td>
              <td class="py-4 px-4 text-right">
                <button 
                  @click="openDetailModal(log)" 
                  class="font-bold text-xs text-slate-600 hover:text-slate-900 transition-colors flex items-center gap-1.5 ml-auto cursor-pointer"
                >
                  <EyeIcon class="w-4 h-4" />
                  <span>Lihat Detail</span>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="px-8 py-6 bg-slate-50/50 border-t border-slate-100 flex items-center justify-between mt-auto">
        <div class="flex items-center gap-6">
          <div class="flex items-center gap-3">
            <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Tampilkan</span>
            <select v-model="limit" @change="fetchLogs(1)" class="bg-white border border-slate-200 rounded-lg text-[10px] font-black text-slate-600 px-2 py-1 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all cursor-pointer shadow-sm">
              <option :value="10">10</option><option :value="25">25</option><option :value="50">50</option><option :value="100">100</option>
            </select>
          </div>
          <div class="h-8 w-px bg-slate-200 hidden sm:block"></div>
          <span class="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em]">
            Halaman <span class="text-indigo-600">{{ page }}</span> dari {{ totalPages }} <span class="mx-2 text-slate-300">|</span> Total <span class="text-indigo-600">{{ totalData }}</span> Data
          </span>
        </div>
        <div class="flex items-center gap-2">
          <button v-if="totalPages > 1" @click="fetchLogs(page - 1)" :disabled="page === 1 || loading" 
            class="w-10 h-10 bg-white border border-slate-100 rounded-xl text-slate-400 hover:text-indigo-600 hover:border-indigo-100 disabled:opacity-30 transition-all shadow-sm flex items-center justify-center cursor-pointer">
            <PrevIcon class="w-4 h-4" />
          </button>
          <div class="flex items-center gap-1.5 mx-1">
            <button v-for="p in visiblePages" :key="p" @click="fetchLogs(p)"
              class="w-10 h-10 rounded-xl text-[10px] font-black transition-all border flex items-center justify-center shadow-sm cursor-pointer"
              :class="page === p ? 'bg-indigo-600 text-white border-indigo-600 shadow-lg shadow-indigo-600/20' : 
                                 'bg-white border-slate-100 text-slate-400 hover:bg-slate-50'">
              {{ p }}
            </button>
          </div>
          <button v-if="totalPages > 1" @click="fetchLogs(page + 1)" :disabled="page >= totalPages || loading" 
            class="w-10 h-10 bg-white border border-slate-100 rounded-xl text-slate-400 hover:text-indigo-600 hover:border-indigo-100 disabled:opacity-30 transition-all shadow-sm flex items-center justify-center cursor-pointer">
            <NextIcon class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>

    <!-- Detail Modal -->
    <AuditLogModal 
      :is-open="showModal" 
      :log="selectedLog" 
      @close="showModal = false" 
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import AuditLogModal from '../../components/audit/AuditLogModal.vue'
import AuditFilter from '../../components/audit/AuditFilter.vue'
import { 
  Database as DatabaseIcon, 
  Search as SearchIcon, 
  Filter as FilterIcon, 
  RefreshCw as RefreshIcon, 
  CheckCircle2 as CheckIcon, 
  ShieldAlert as ShieldAlertIcon, 
  Clock as ClockIcon, 
  Eye as EyeIcon, 
  ChevronLeft as PrevIcon, 
  ChevronRight as NextIcon,
  RotateCcw as ResetIcon
} from 'lucide-vue-next'

const route = useRoute()
const isMounted = ref(false)
const loading = ref(false)
const logs = ref([])
const totalData = ref(0)
const page = ref(1)
const limit = ref(10)

const showModal = ref(false)
const selectedLog = ref(null)
const showFilters = ref(false)

const filters = reactive({
  user_name: '',
  action: '',
  entity_type: '',
  role: '',
  sort: ''
})

let searchTimeout = null

const onSearchInput = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    fetchLogs(1)
  }, 300)
}

const resetFilters = () => {
  filters.user_name = ''
  filters.action = ''
  filters.entity_type = ''
  filters.role = ''
  filters.sort = ''
  showFilters.value = false
  fetchLogs(1)
}

const applyFilters = () => {
  showFilters.value = false
  fetchLogs(1)
}

const totalPages = computed(() => {
  return Math.max(1, Math.ceil(totalData.value / limit.value))
})

const visiblePages = computed(() => {
  const pages = []
  let startPage = Math.max(1, page.value - 1)
  let endPage = Math.min(totalPages.value, startPage + 2)
  if (endPage - startPage < 2) startPage = Math.max(1, endPage - 2)
  for (let i = startPage; i <= endPage; i++) if (i > 0) pages.push(i)
  return pages
})

const fetchLogs = async (p = 1) => {
  page.value = p
  loading.value = true
  try {
    const res = await axios.get('audit-logs', {
      params: { 
        page: page.value, 
        limit: limit.value,
        user_name: filters.user_name,
        action: filters.action,
        entity_type: filters.entity_type,
        role: filters.role,
        sort: filters.sort || undefined
      }
    })
    logs.value = res.data.data.logs || []
    totalData.value = res.data.data.total || 0
  } catch (err) {
    console.error('Gagal memuat log audit', err)
  } finally {
    loading.value = false
  }
}

const openDetailModal = (log) => {
  selectedLog.value = log
  showModal.value = true
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('id-ID', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit', second: '2-digit'
  })
}

const humanizeAction = (action) => {
  if (!action) return '-'
  const dict = {
    CREATE: 'Membuat Baru',
    UPDATE: 'Memperbarui Data',
    DELETE: 'Menghapus Data',
    RESTORE: 'Memulihkan Data',
    BULK_DELETE: 'Menghapus Massal',
    BULK_RESTORE: 'Memulihkan Massal',
    TOGGLE_STATUS: 'Mengubah Status',
    ACTIVATE_ACCOUNT: 'Mengaktifkan Akun',
    SEND_NOTIFICATION: 'Mengirim Notifikasi',
    LOGIN: 'Login Sistem',
    LOGIN_FAILED: 'Login Gagal',
    LOGOUT: 'Logout Sistem',
    REFRESH_TOKEN: 'Perpanjang Sesi',
    FORGOT_PASSWORD: 'Minta Reset Sandi',
    FORGOT_PASSWORD_FAILED: 'Permintaan Reset Gagal',
    RESET_PASSWORD: 'Reset Kata Sandi',
    RESET_PASSWORD_FAILED: 'Reset Kata Sandi Gagal',
    CHANGE_PASSWORD: 'Ubah Kata Sandi',
    CHANGE_PASSWORD_FAILED: 'Ubah Kata Sandi Gagal',
    EXPORT_TREND_EXCEL: 'Export Excel',
    GENERATE_BILL: 'Menerbitkan Tagihan',
    GENERATE_ADJUSTMENT_BILL: 'Menerbitkan Penyesuaian Tagihan',
    REFUND_BILL_REDUCTION: 'Refund Selisih Tagihan',
    BULK_CANCEL_GENERATED_BILL: 'Menarik Tagihan',
    VOID_BILL: 'Membatalkan Tagihan',
    EXPORT_GLOBAL_FINANCE_REPORT: 'Export Laporan Keuangan',
    CREATE_PAYMENT_INTENT: 'Membuat Tagihan Bayar',
    PROCESS_PAYMENT: 'Memproses Pembayaran',
    GATEWAY_PAYMENT_TO_DEPOSIT: 'Dana Gateway ke Saldo'
  }
  return dict[action] || action
}

const humanizeEntity = (entity) => {
  if (!entity) return '-'
  const dict = {
    users: 'Pengguna',
    students: 'Siswa',
    majors: 'Jurusan',
    classes: 'Kelas',
    academic_years: 'Tahun Ajaran',
    bill_types: 'Jenis Tagihan',
    billing_rules: 'Aturan Tagihan',
    student_bills: 'Tagihan Siswa',
    payments: 'Pembayaran',
    finance_reports: 'Laporan Keuangan',
    notifications: 'Notifikasi',
    auth: 'Autentikasi'
  }
  return dict[entity] || entity
}

onMounted(() => {
  isMounted.value = true
  if (route.query.entity_type) {
    filters.entity_type = route.query.entity_type
  }
  fetchLogs(1)
})
</script>

<style scoped lang="postcss">
.search-input-premium {
  @apply w-full bg-white border border-slate-200 rounded-xl py-2.5 pl-12 pr-4 text-xs font-bold text-slate-700 outline-none transition-all focus:border-indigo-500 focus:ring-4 focus:ring-indigo-50 shadow-sm;
}
.btn-primary { @apply bg-indigo-600 hover:bg-indigo-700 text-white transition-all disabled:opacity-50; }
.btn-secondary { @apply bg-white hover:bg-slate-50 text-slate-600 border border-slate-200 transition-all; }
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
