<template>
  <div class="space-y-8 animate-fade-in">
    <!-- Efficacy Stats Header Cards -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <!-- Total Sent -->
      <div class="bg-white rounded-3xl p-6 border border-slate-200 shadow-sm flex items-center gap-4">
        <div class="w-12 h-12 bg-slate-50 text-slate-600 rounded-2xl flex items-center justify-center">
          <SendIcon class="w-6 h-6" />
        </div>
        <div>
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest leading-none mb-1.5">Terkirim</p>
          <h3 class="text-2xl font-black text-slate-800 leading-none">{{ stats.SENT || 0 }}</h3>
        </div>
      </div>

      <!-- Total Delivered -->
      <div class="bg-white rounded-3xl p-6 border border-slate-200 shadow-sm flex items-center gap-4">
        <div class="w-12 h-12 bg-sky-50 text-sky-600 rounded-2xl flex items-center justify-center">
          <CheckIcon class="w-6 h-6" />
        </div>
        <div>
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest leading-none mb-1.5">Diterima</p>
          <h3 class="text-2xl font-black text-slate-800 leading-none">{{ stats.DELIVERED || 0 }}</h3>
        </div>
      </div>

      <!-- Total Read -->
      <div class="bg-white rounded-3xl p-6 border border-slate-200 shadow-sm flex items-center gap-4">
        <div class="w-12 h-12 bg-emerald-50 text-emerald-600 rounded-2xl flex items-center justify-center">
          <CheckIcon class="w-6 h-6" />
        </div>
        <div>
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest leading-none mb-1.5">Dibaca</p>
          <h3 class="text-2xl font-black text-slate-800 leading-none">{{ stats.READ || 0 }}</h3>
        </div>
      </div>

      <!-- Total Failed -->
      <div class="bg-white rounded-3xl p-6 border border-slate-200 shadow-sm flex items-center gap-4">
        <div class="w-12 h-12 bg-rose-50 text-rose-600 rounded-2xl flex items-center justify-center">
          <AlertCircleIcon class="w-6 h-6" />
        </div>
        <div>
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest leading-none mb-1.5">Gagal</p>
          <h3 class="text-2xl font-black text-slate-800 leading-none">{{ stats.FAILED || 0 }}</h3>
        </div>
      </div>
    </div>

    <!-- Header Teleport Target Content -->
    <Teleport to="#header-actions-target" v-if="isMounted">
      <div class="flex items-center justify-center w-full gap-4 relative mx-auto">
        <div class="flex items-center justify-center gap-2 flex-1 max-w-2xl mx-auto">
          <div class="relative flex-1 group">
            <SearchIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300 group-focus-within:text-indigo-600" />
            <input v-model="searchQuery" @input="onSearchInput" type="text" placeholder="Cari nama penerima atau pesan..." class="search-input-premium" />
          </div>
          <button @click="resetFilters" class="p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm transition-all group" title="Reset">
            <ResetIcon class="w-4 h-4 group-hover:rotate-180 transition-transform duration-500" />
          </button>
        </div>
      </div>
    </Teleport>

    <!-- Main Table Card -->
    <div class="bg-white rounded-3xl border border-slate-200 shadow-sm flex flex-col min-h-[600px] overflow-hidden">
      <!-- Table Filter Bar -->
      <div class="px-8 py-6 border-b border-slate-100 bg-slate-50/30 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div class="flex items-center gap-4">
          <div class="w-2 h-6 bg-indigo-500 rounded-full"></div>
          <h3 class="font-black text-slate-700 text-sm uppercase tracking-[0.2em]">Log Notifikasi & Efikasi</h3>
        </div>
        <!-- Status Filters Tab -->
        <div class="flex items-center gap-2 self-start sm:self-center overflow-x-auto pb-1 sm:pb-0">
          <button 
            v-for="statusTab in statusTabs" 
            :key="statusTab.value" 
            @click="selectStatus(statusTab.value)" 
            class="px-4 py-2 rounded-xl text-[10px] font-black uppercase tracking-wider border transition-all cursor-pointer whitespace-nowrap"
            :class="selectedStatus === statusTab.value ? 'bg-indigo-600 border-indigo-600 text-white shadow-lg shadow-indigo-600/10' : 'bg-white border-slate-200 text-slate-600 hover:bg-slate-50'"
          >
            {{ statusTab.label }}
          </button>
        </div>
      </div>

      <!-- Table -->
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-slate-100 text-[10px] font-black text-slate-400 uppercase tracking-widest">
              <th class="py-4 px-6"># ID</th>
              <th class="py-4 px-6">Waktu</th>
              <th class="py-4 px-6">Penerima</th>
              <th class="py-4 px-6">Channel</th>
              <th class="py-4 px-6">Tujuan</th>
              <th class="py-4 px-6">Judul / Subjek</th>
              <th class="py-4 px-6">Pesan</th>
              <th class="py-4 px-6">Status</th>
              <th class="py-4 px-6 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="text-xs font-semibold divide-y divide-slate-100/80">
            <tr v-if="loading" class="animate-pulse">
              <td colspan="9" class="py-12 text-center text-slate-400 font-bold">Memuat log notifikasi...</td>
            </tr>
            <tr v-else-if="logs.length === 0">
              <td colspan="9" class="py-20 px-6 text-center">
                <div class="flex flex-col items-center justify-center text-center mx-auto max-w-sm">
                  <div class="w-20 h-20 bg-slate-50 rounded-[2.5rem] flex items-center justify-center text-slate-300 mb-6 border border-slate-100 mx-auto">
                    <MessageCircleIcon class="w-10 h-10" />
                  </div>
                  <h3 class="text-lg font-black text-slate-700 tracking-tight mb-2">Log Notifikasi Kosong</h3>
                  <p class="text-slate-400 text-xs font-medium">Belum ada data pengiriman notifikasi yang sesuai dengan filter saat ini.</p>
                </div>
              </td>
            </tr>
            <tr v-else v-for="log in filteredLogs" :key="log.id" class="hover:bg-slate-50/50 transition-colors group">
              <td class="py-4 px-6 font-black text-slate-700">#{{ log.id }}</td>
              <td class="py-4 px-6 text-slate-500 whitespace-nowrap">{{ formatDate(log.created_at) }}</td>
              <td class="py-4 px-6 font-bold text-slate-700 whitespace-nowrap">{{ log.recipient_name || 'Orang Tua' }}</td>
              <td class="py-4 px-6">
                <span class="px-2.5 py-1 rounded-full text-[9px] font-black uppercase tracking-widest border" :class="(log.channel || (log.whatsapp_id ? 'whatsapp' : 'email')) === 'whatsapp' ? 'bg-emerald-50 text-emerald-600 border-emerald-100' : 'bg-sky-50 text-sky-600 border-sky-100'">
                  {{ log.channel || (log.whatsapp_id ? 'whatsapp' : 'email') }}
                </span>
              </td>
              <td class="py-4 px-6 text-slate-500 font-mono text-[11px]">{{ log.whatsapp_id ? (log.recipient_phone || 'WhatsApp') : (log.recipient_email || 'Email') }}</td>
              <td class="py-4 px-6 font-bold text-slate-600 truncate max-w-[150px]">{{ log.title }}</td>
              <td class="py-4 px-6 text-slate-400 truncate max-w-[200px]">{{ log.message }}</td>
              <td class="py-4 px-6">
                <span class="px-2.5 py-1 rounded-full text-[9px] font-black uppercase tracking-widest border" :class="getStatusBadgeClass(log.delivery_status)">
                  {{ log.delivery_status || 'SENT' }}
                </span>
                <p v-if="log.delivery_error" class="mt-1 max-w-[180px] truncate text-[10px] font-bold text-rose-500" :title="log.delivery_error">{{ log.delivery_error }}</p>
              </td>
              <td class="py-4 px-6 text-right whitespace-nowrap">
                <div class="flex items-center justify-end gap-2">
                  <button 
                    @click="openDetailModal(log)" 
                    class="p-2 bg-slate-50 hover:bg-slate-100 text-slate-600 rounded-xl transition-all shadow-sm flex items-center justify-center cursor-pointer"
                    title="Lihat Pesan"
                  >
                    <EyeIcon class="w-3.5 h-3.5" />
                  </button>
                  <button 
                    @click="resendNotification(log)" 
                    :disabled="resendingId === log.id"
                    class="px-3 py-2 bg-indigo-50 hover:bg-indigo-100 text-indigo-600 rounded-xl transition-all font-black text-[10px] uppercase tracking-wider flex items-center gap-1.5 cursor-pointer disabled:opacity-50"
                  >
                    <RefreshIcon class="w-3 h-3" :class="{ 'animate-spin': resendingId === log.id }" />
                    <span>Kirim Ulang</span>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="px-8 py-6 bg-slate-50/30 border-t border-slate-100 flex items-center justify-between mt-auto">
        <div class="flex items-center gap-6">
          <div class="flex items-center gap-3">
            <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Tampilkan</span>
            <select v-model="limit" @change="fetchLogs(1)" class="bg-white border border-slate-200 rounded-lg text-[10px] font-black text-slate-600 px-2 py-1 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all cursor-pointer shadow-sm">
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

    <!-- Detail Message Modal -->
    <Teleport to="body">
      <transition name="page">
        <div v-if="showDetailModal && selectedLog" class="fixed inset-0 z-[200] flex items-center justify-center p-6">
          <div class="absolute inset-0 bg-slate-900/40 backdrop-blur-sm" @click="showDetailModal = false"></div>
          <div class="bg-white w-full max-w-xl relative z-10 rounded-[2.5rem] shadow-2xl overflow-hidden animate-scale-in">
            <div class="p-8 space-y-6">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center">
                    <MessageCircleIcon class="w-5 h-5" />
                  </div>
                  <div>
                    <h3 class="text-base font-black text-slate-800 tracking-tight">Detail Pesan Notifikasi</h3>
                    <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mt-0.5">Informasi Pengiriman</p>
                  </div>
                </div>
                <span class="px-2.5 py-1 rounded-full text-[9px] font-black uppercase tracking-widest border" :class="getStatusBadgeClass(selectedLog.delivery_status)">
                  {{ selectedLog.delivery_status || 'SENT' }}
                </span>
              </div>

              <!-- Content Card -->
              <div class="bg-slate-50 border border-slate-100 p-6 rounded-2xl space-y-4">
                <div>
                  <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest block mb-1">Penerima</label>
                  <p class="text-xs font-bold text-slate-700">{{ selectedLog.recipient_name || 'Orang Tua' }}</p>
                </div>

                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest block mb-1">Kontak Tujuan</label>
                    <p class="text-xs font-bold text-slate-600 font-mono">{{ selectedLog.whatsapp_id ? (selectedLog.recipient_phone || 'WhatsApp') : (selectedLog.recipient_email || 'Email') }}</p>
                    <p class="mt-1 text-[9px] font-black uppercase tracking-widest text-slate-400">{{ selectedLog.channel || (selectedLog.whatsapp_id ? 'whatsapp' : 'email') }}</p>
                  </div>
                  <div>
                    <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest block mb-1">Waktu Kirim</label>
                    <p class="text-xs font-bold text-slate-600">{{ formatDate(selectedLog.created_at) }}</p>
                  </div>
                </div>

                <div>
                  <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest block mb-1">Judul / Subjek</label>
                  <p class="text-xs font-bold text-slate-800">{{ selectedLog.title }}</p>
                </div>

                <div v-if="selectedLog.delivery_error">
                  <label class="text-[9px] font-black text-rose-400 uppercase tracking-widest block mb-1">Alasan Gagal</label>
                  <div class="text-xs font-bold text-rose-600 bg-rose-50 border border-rose-100 rounded-xl p-3">{{ selectedLog.delivery_error }}</div>
                </div>

                <div>
                  <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest block mb-1">Isi Pesan</label>
                  <div class="text-xs font-semibold text-slate-700 whitespace-pre-wrap leading-relaxed border-t border-slate-200/50 pt-2">{{ selectedLog.message }}</div>
                </div>
              </div>

              <div class="flex items-center gap-3">
                <button 
                  @click="showDetailModal = false" 
                  class="flex-1 py-3 border border-slate-200 text-slate-600 rounded-2xl font-black text-[10px] uppercase tracking-widest hover:bg-slate-50 transition-all text-center"
                >
                  Tutup
                </button>
                <button 
                  @click="resendNotification(selectedLog); showDetailModal = false" 
                  class="flex-1 py-3 bg-indigo-600 text-white rounded-2xl font-black text-[10px] uppercase tracking-widest hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-100 flex items-center justify-center gap-1.5"
                >
                  <RefreshIcon class="w-3.5 h-3.5" />
                  <span>Kirim Ulang</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted, computed, watch } from 'vue'
import axios from 'axios'
import { useToast } from '../../composables/useToast'
import { 
  MessageCircle as MessageCircleIcon,
  Search as SearchIcon,
  Filter as FilterIcon,
  RefreshCw as RefreshIcon,
  CheckCircle2 as CheckIcon,
  AlertCircle as AlertCircleIcon,
  Eye as EyeIcon,
  ChevronLeft as PrevIcon,
  ChevronRight as NextIcon,
  RotateCcw as ResetIcon,
  Send as SendIcon
} from 'lucide-vue-next'

const isMounted = ref(false)
const loading = ref(false)
const logs = ref([])
const stats = ref({})
const totalData = ref(0)
const page = ref(1)
const limit = ref(10)
const searchQuery = ref('')
const selectedStatus = ref('')
const resendingId = ref(null)
const toast = useToast()

// Detail Modal State
const showDetailModal = ref(false)
const selectedLog = ref(null)

const statusTabs = [
  { label: 'Semua', value: '' },
  { label: 'Terkirim', value: 'sent' },
  { label: 'Diterima', value: 'delivered' },
  { label: 'Dibaca', value: 'read' },
  { label: 'Gagal', value: 'failed' }
]

let searchTimeout = null

const onSearchInput = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    fetchLogs(1)
  }, 300)
}

const selectStatus = (status) => {
  selectedStatus.value = status
  fetchLogs(1)
}

const resetFilters = () => {
  searchQuery.value = ''
  selectedStatus.value = ''
  fetchLogs(1)
}

const handleNotificationStatusChanged = () => {
  fetchStats()
  fetchLogs(page.value)
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

const fetchStats = async () => {
  try {
    const res = await axios.get('whatsapp/stats')
    stats.value = res.data.data || {}
  } catch (err) {
    console.error('Failed to fetch notification stats', err)
  }
}

const fetchLogs = async (p = 1) => {
  page.value = p
  loading.value = true
  try {
    const res = await axios.get('whatsapp/logs', {
      params: { 
        page: page.value, 
        limit: limit.value,
        status: selectedStatus.value,
        search: searchQuery.value.trim()
      }
    })
    logs.value = res.data.data.data || []
    totalData.value = res.data.data.total || 0
  } catch (err) {
    console.error('Failed to load notification logs', err)
    toast.error('Gagal', 'Gagal memuat log notifikasi')
  } finally {
    loading.value = false
  }
}

// Local search filter to filter down UI results dynamically as helper
const filteredLogs = computed(() => {
  if (!searchQuery.value) return logs.value
  const query = searchQuery.value.toLowerCase()
  return logs.value.filter(log => {
    return (
      (log.recipient_name && log.recipient_name.toLowerCase().includes(query)) ||
      (log.message && log.message.toLowerCase().includes(query)) ||
      (log.title && log.title.toLowerCase().includes(query))
    )
  })
})

const resendNotification = async (log) => {
  resendingId.value = log.id
  try {
    await axios.post(`whatsapp/notifications/${log.id}/resend`)
    toast.success('Sukses', 'Notifikasi berhasil dikirim ulang')
    fetchLogs(page.value)
    fetchStats()
  } catch (err) {
    console.error('Failed to resend notification', err)
    toast.error('Gagal', 'Gagal mengirim ulang notifikasi')
  } finally {
    resendingId.value = null
  }
}

const openDetailModal = (log) => {
  selectedLog.value = log
  showDetailModal.value = true
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('id-ID', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}

const getStatusBadgeClass = (status) => {
  const norm = (status || 'sent').toLowerCase()
  switch (norm) {
    case 'read':
      return 'bg-emerald-50 text-emerald-600 border-emerald-100'
    case 'delivered':
      return 'bg-sky-50 text-sky-600 border-sky-100'
    case 'sent':
      return 'bg-slate-50 text-slate-600 border-slate-200'
    case 'failed':
      return 'bg-rose-50 text-rose-600 border-rose-100'
    default:
      return 'bg-slate-50 text-slate-600 border-slate-200'
  }
}

onMounted(() => {
  isMounted.value = true
  fetchStats()
  fetchLogs(1)
  window.addEventListener('notification-status-changed', handleNotificationStatusChanged)
})

onUnmounted(() => {
  window.removeEventListener('notification-status-changed', handleNotificationStatusChanged)
})
</script>

<style scoped lang="postcss">
.search-input-premium {
  @apply w-full bg-white border border-slate-200 rounded-xl py-2.5 pl-12 pr-4 text-xs font-bold text-slate-700 outline-none transition-all focus:border-indigo-500 focus:ring-4 focus:ring-indigo-50 shadow-sm;
}
</style>
