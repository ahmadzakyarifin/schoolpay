<template>
  <div class="space-y-8 animate-fade-in">
    <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
      <div v-for="card in statCards" :key="card.label" class="bg-white rounded-2xl p-6 border border-slate-200 shadow-sm flex items-center gap-4">
        <div :class="[card.color, 'w-12 h-12 rounded-xl flex items-center justify-center']">
          <component :is="card.icon" class="w-6 h-6" />
        </div>
        <div>
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest leading-none mb-1.5">{{ card.label }}</p>
          <h3 class="text-2xl font-black text-slate-800 leading-none">{{ card.value }}</h3>
        </div>
      </div>
    </div>

    <Teleport to="#header-actions-target" v-if="isMounted">
      <div class="flex items-center justify-center w-full gap-4 relative mx-auto">
        <div class="flex items-center justify-center gap-2 flex-1 max-w-2xl mx-auto">
          <div class="relative flex-1 group">
            <SearchIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300 group-focus-within:text-indigo-600" />
            <input v-model="searchQuery" @input="onSearchInput" type="text" placeholder="Cari penerima, tujuan, judul, atau pesan..." class="search-input-premium" />
          </div>
          <button @click="resetFilters" class="p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm transition-all group" title="Reset">
            <ResetIcon class="w-4 h-4 group-hover:rotate-180 transition-transform duration-500" />
          </button>
        </div>
      </div>
    </Teleport>

    <div class="flex items-center gap-2 overflow-x-auto pb-1">
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

    <section v-for="table in channelTables" :key="table.channel" class="bg-white rounded-2xl border border-slate-200 shadow-sm flex flex-col min-h-[520px] overflow-hidden">
      <div class="px-8 py-6 border-b border-slate-100 bg-slate-50/30 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div class="flex items-center gap-4">
          <div :class="[table.channel === 'whatsapp' ? 'bg-emerald-500' : 'bg-sky-500', 'w-2 h-6 rounded-full']"></div>
          <div>
            <h3 class="font-black text-slate-700 text-sm uppercase tracking-[0.2em]">{{ table.title }}</h3>
            <p class="mt-1 text-[10px] font-bold uppercase tracking-wider text-slate-400">
              Total <span :class="table.channel === 'whatsapp' ? 'text-emerald-600' : 'text-sky-600'">{{ totalByChannel[table.channel] }}</span> data sesuai filter
            </p>
          </div>
        </div>
        <div :class="[table.channel === 'whatsapp' ? 'bg-emerald-50 text-emerald-600 border-emerald-100' : 'bg-sky-50 text-sky-600 border-sky-100', 'px-3 py-1.5 rounded-full border text-[10px] font-black uppercase tracking-widest flex items-center gap-2']">
          <component :is="table.icon" class="w-3.5 h-3.5" />
          {{ table.badge }}
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="border-b border-slate-100 text-[10px] font-black text-slate-400 uppercase tracking-widest">
              <th class="py-4 px-6"># ID</th>
              <th class="py-4 px-6">Waktu</th>
              <th class="py-4 px-6">Penerima</th>
              <th class="py-4 px-6">Tujuan</th>
              <th class="py-4 px-6">Judul / Subjek</th>
              <th class="py-4 px-6">Pesan</th>
              <th class="py-4 px-6">Status</th>
              <th class="py-4 px-6 text-right">Aksi</th>
            </tr>
          </thead>
          <tbody class="text-xs font-semibold divide-y divide-slate-100/80">
            <tr v-if="loadingByChannel[table.channel]" class="animate-pulse">
              <td colspan="8" class="py-12 text-center text-slate-400 font-bold">Memuat log {{ table.badge }}...</td>
            </tr>
            <tr v-else-if="logsByChannel[table.channel].length === 0">
              <td colspan="8" class="py-20 px-6 text-center">
                <div class="flex flex-col items-center justify-center text-center mx-auto max-w-sm">
                  <div class="w-20 h-20 bg-slate-50 rounded-2xl flex items-center justify-center text-slate-300 mb-6 border border-slate-100 mx-auto">
                    <component :is="table.icon" class="w-10 h-10" />
                  </div>
                  <h3 class="text-lg font-black text-slate-700 tracking-tight mb-2">Log {{ table.badge }} Kosong</h3>
                  <p class="text-slate-400 text-xs font-medium">Belum ada data pengiriman yang sesuai dengan filter saat ini.</p>
                </div>
              </td>
            </tr>
            <tr v-else v-for="log in logsByChannel[table.channel]" :key="`${table.channel}-${log.id}`" class="hover:bg-slate-50/50 transition-colors group">
              <td class="py-4 px-6 font-black text-slate-700">#{{ log.id }}</td>
              <td class="py-4 px-6 text-slate-500 whitespace-nowrap">{{ formatDate(log.created_at) }}</td>
              <td class="py-4 px-6 font-bold text-slate-700 whitespace-nowrap">{{ log.recipient_name || 'Orang Tua' }}</td>
              <td class="py-4 px-6 text-slate-500 font-mono text-[11px]">{{ table.channel === 'whatsapp' ? (log.recipient_phone || 'WhatsApp') : (log.recipient_email || 'Email') }}</td>
              <td class="py-4 px-6 font-bold text-slate-600 truncate max-w-[170px]">{{ log.title }}</td>
              <td class="py-4 px-6 text-slate-400 truncate max-w-[240px]">{{ log.message }}</td>
              <td class="py-4 px-6">
                <span class="px-2.5 py-1 rounded-full text-[9px] font-black uppercase tracking-widest border" :class="getStatusBadgeClass(log.delivery_status)">
                  {{ statusLabel(log.delivery_status) }}
                </span>
                <p v-if="log.delivery_error" class="mt-1 max-w-[180px] truncate text-[10px] font-bold text-rose-500" :title="log.delivery_error">{{ log.delivery_error }}</p>
              </td>
              <td class="py-4 px-6 text-right whitespace-nowrap">
                <div class="flex items-center justify-end gap-2">
                  <button @click="openDetailModal(log, table.channel)" class="p-2 bg-slate-50 hover:bg-slate-100 text-slate-600 rounded-xl transition-all shadow-sm flex items-center justify-center cursor-pointer" title="Lihat Pesan">
                    <EyeIcon class="w-3.5 h-3.5" />
                  </button>
                  <button @click="resendNotification(log)" :disabled="resendingId === log.id" class="px-3 py-2 bg-indigo-50 hover:bg-indigo-100 text-indigo-600 rounded-xl transition-all font-black text-[10px] uppercase tracking-wider flex items-center gap-1.5 cursor-pointer disabled:opacity-50">
                    <RefreshIcon class="w-3 h-3" :class="{ 'animate-spin': resendingId === log.id }" />
                    <span>Kirim Ulang</span>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

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
          <span class="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em]">
            Halaman <span class="text-indigo-600">{{ pageByChannel[table.channel] }}</span> dari {{ totalPages(table.channel) }}
          </span>
        </div>
        <div class="flex items-center gap-2">
          <button v-if="totalPages(table.channel) > 1" @click="fetchLogs(pageByChannel[table.channel] - 1, table.channel)" :disabled="pageByChannel[table.channel] === 1 || loadingByChannel[table.channel]" class="w-10 h-10 bg-white border border-slate-100 rounded-xl text-slate-400 hover:text-indigo-600 disabled:opacity-30 transition-all shadow-sm flex items-center justify-center cursor-pointer">
            <PrevIcon class="w-4 h-4" />
          </button>
          <button v-if="totalPages(table.channel) > 1" @click="fetchLogs(pageByChannel[table.channel] + 1, table.channel)" :disabled="pageByChannel[table.channel] >= totalPages(table.channel) || loadingByChannel[table.channel]" class="w-10 h-10 bg-white border border-slate-100 rounded-xl text-slate-400 hover:text-indigo-600 disabled:opacity-30 transition-all shadow-sm flex items-center justify-center cursor-pointer">
            <NextIcon class="w-4 h-4" />
          </button>
        </div>
      </div>
    </section>

    <Teleport to="body">
      <transition name="page">
        <div v-if="showDetailModal && selectedLog" class="fixed inset-0 z-[200] flex items-center justify-center p-6">
          <div class="absolute inset-0 bg-slate-900/40 backdrop-blur-sm" @click="showDetailModal = false"></div>
          <div class="bg-white w-full max-w-xl relative z-10 rounded-2xl shadow-2xl overflow-hidden animate-scale-in">
            <div class="p-8 space-y-6">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center">
                    <MessageCircleIcon class="w-5 h-5" />
                  </div>
                  <div>
                    <h3 class="text-base font-black text-slate-800 tracking-tight">Detail Pesan Notifikasi</h3>
                    <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mt-0.5">{{ selectedChannelLabel }}</p>
                  </div>
                </div>
                <span class="px-2.5 py-1 rounded-full text-[9px] font-black uppercase tracking-widest border" :class="getStatusBadgeClass(selectedLog.delivery_status)">
                  {{ statusLabel(selectedLog.delivery_status) }}
                </span>
              </div>

              <div class="bg-slate-50 border border-slate-100 p-6 rounded-2xl space-y-4">
                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="field-label">Penerima</label>
                    <p class="field-value">{{ selectedLog.recipient_name || 'Orang Tua' }}</p>
                  </div>
                  <div>
                    <label class="field-label">Kontak Tujuan</label>
                    <p class="field-value font-mono">{{ selectedLog.channel === 'whatsapp' ? (selectedLog.recipient_phone || 'WhatsApp') : (selectedLog.recipient_email || 'Email') }}</p>
                  </div>
                </div>
                <div>
                  <label class="field-label">Judul / Subjek</label>
                  <p class="field-value text-slate-800">{{ selectedLog.title }}</p>
                </div>
                <div v-if="selectedLog.delivery_error">
                  <label class="field-label text-rose-400">Alasan Gagal</label>
                  <div class="text-xs font-bold text-rose-600 bg-rose-50 border border-rose-100 rounded-xl p-3">{{ selectedLog.delivery_error }}</div>
                </div>
                <div>
                  <label class="field-label">Isi Pesan</label>
                  <div class="text-xs font-semibold text-slate-700 whitespace-pre-wrap leading-relaxed border-t border-slate-200/50 pt-2">{{ selectedLog.message }}</div>
                </div>
              </div>

              <div class="flex items-center gap-3">
                <button @click="showDetailModal = false" class="flex-1 py-3 border border-slate-200 text-slate-600 rounded-xl font-black text-[10px] uppercase tracking-widest hover:bg-slate-50 transition-all text-center">Tutup</button>
                <button @click="resendNotification(selectedLog); showDetailModal = false" class="flex-1 py-3 bg-indigo-600 text-white rounded-xl font-black text-[10px] uppercase tracking-widest hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-100 flex items-center justify-center gap-1.5">
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
import { ref, reactive, onMounted, onUnmounted, computed } from 'vue'
import axios from 'axios'
import { useToast } from '../../composables/useToast'
import {
  MessageCircle as MessageCircleIcon,
  Mail as MailIcon,
  Search as SearchIcon,
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
const logsByChannel = reactive({ whatsapp: [], email: [] })
const loadingByChannel = reactive({ whatsapp: false, email: false })
const totalByChannel = reactive({ whatsapp: 0, email: 0 })
const pageByChannel = reactive({ whatsapp: 1, email: 1 })
const stats = ref({})
const limit = ref(10)
const searchQuery = ref('')
const selectedStatus = ref('')
const resendingId = ref(null)
const toast = useToast()
const showDetailModal = ref(false)
const selectedLog = ref(null)
const selectedChannelLabel = ref('')

const channelTables = [
  { channel: 'whatsapp', title: 'Efikasi WhatsApp', badge: 'WhatsApp', icon: MessageCircleIcon },
  { channel: 'email', title: 'Efikasi Email', badge: 'Email', icon: MailIcon }
]

const statusTabs = [
  { label: 'Semua', value: '' },
  { label: 'Menunggu', value: 'pending' },
  { label: 'Terkirim', value: 'sent' },
  { label: 'Diterima', value: 'delivered' },
  { label: 'Dibaca', value: 'read' },
  { label: 'Gagal', value: 'failed' }
]

const statCards = computed(() => [
  { label: 'Terkirim', value: stats.value.SENT || 0, icon: SendIcon, color: 'bg-slate-50 text-slate-600' },
  { label: 'Diterima', value: stats.value.DELIVERED || 0, icon: CheckIcon, color: 'bg-sky-50 text-sky-600' },
  { label: 'Dibaca', value: stats.value.READ || 0, icon: CheckIcon, color: 'bg-emerald-50 text-emerald-600' },
  { label: 'Gagal', value: stats.value.FAILED || 0, icon: AlertCircleIcon, color: 'bg-rose-50 text-rose-600' }
])

let searchTimeout = null

const onSearchInput = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => fetchLogs(1), 300)
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
  fetchLogs()
}

const totalPages = (channel) => Math.max(1, Math.ceil(totalByChannel[channel] / limit.value))

const fetchStats = async () => {
  try {
    const res = await axios.get('whatsapp/stats')
    stats.value = res.data.data || {}
  } catch (err) {
    console.error('Failed to fetch notification stats', err)
  }
}

const fetchOneChannel = async (channel, p = pageByChannel[channel]) => {
  pageByChannel[channel] = Math.max(1, p)
  loadingByChannel[channel] = true
  try {
    const res = await axios.get('whatsapp/logs', {
      params: {
        page: pageByChannel[channel],
        limit: limit.value,
        status: selectedStatus.value,
        search: searchQuery.value.trim(),
        channel
      }
    })
    logsByChannel[channel] = res.data.data.data || []
    totalByChannel[channel] = res.data.data.total || 0
  } catch (err) {
    console.error(`Failed to load ${channel} logs`, err)
    toast.error('Gagal', `Gagal memuat log ${channel === 'whatsapp' ? 'WhatsApp' : 'Email'}`)
  } finally {
    loadingByChannel[channel] = false
  }
}

const fetchLogs = async (p = null, channel = '') => {
  if (channel) {
    await fetchOneChannel(channel, p ?? pageByChannel[channel])
    return
  }
  await Promise.all(channelTables.map((item) => fetchOneChannel(item.channel, p ?? pageByChannel[item.channel])))
}

const resendNotification = async (log) => {
  resendingId.value = log.id
  try {
    await axios.post(`whatsapp/notifications/${log.id}/resend`)
    toast.success('Sukses', 'Notifikasi berhasil dikirim ulang')
    fetchLogs()
    fetchStats()
  } catch (err) {
    console.error('Failed to resend notification', err)
    toast.error('Gagal', 'Gagal mengirim ulang notifikasi')
  } finally {
    resendingId.value = null
  }
}

const openDetailModal = (log, channel) => {
  selectedLog.value = { ...log, channel }
  selectedChannelLabel.value = channel === 'whatsapp' ? 'WhatsApp' : 'Email'
  showDetailModal.value = true
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('id-ID', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}

const statusLabel = (status) => {
  const norm = String(status || 'sent').toLowerCase()
  return statusTabs.find((item) => item.value === norm)?.label || norm.toUpperCase()
}

const getStatusBadgeClass = (status) => {
  const norm = String(status || 'sent').toLowerCase()
  switch (norm) {
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
.field-label {
  @apply text-[9px] font-black text-slate-400 uppercase tracking-widest block mb-1;
}
.field-value {
  @apply text-xs font-bold text-slate-700;
}
</style>
