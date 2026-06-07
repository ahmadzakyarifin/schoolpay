<template>
  <div class="space-y-8 animate-fade-in">

    <!-- ─── STAT CARDS ─── -->
    <div class="grid grid-cols-2 lg:grid-cols-4 gap-4">
      <div
        v-for="card in statCards"
        :key="card.label"
        class="relative bg-white rounded-2xl p-5 border border-slate-100 shadow-sm overflow-hidden group hover:shadow-md transition-all duration-300"
      >
        <div class="absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity duration-300" :class="card.gradientBg"></div>
        <div class="relative flex items-start justify-between mb-4">
          <div :class="[card.iconBg, 'w-11 h-11 rounded-xl flex items-center justify-center shadow-sm']">
            <component :is="card.icon" class="w-5 h-5" :class="card.iconColor" />
          </div>
          <span class="text-[9px] font-black uppercase tracking-widest px-2 py-0.5 rounded-full border" :class="card.badgeClass">
            {{ card.badge }}
          </span>
        </div>
        <div class="relative">
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest leading-none mb-2">{{ card.label }}</p>
          <h3 class="text-3xl font-black text-slate-800 leading-none">{{ card.value.toLocaleString('id-ID') }}</h3>
        </div>
      </div>
    </div>

    <!-- ─── HEADER SEARCH (teleport) ─── -->
    <Teleport to="#header-actions-target" v-if="isMounted">
      <div class="flex items-center justify-center w-full gap-3 mx-auto">
        <div class="relative flex-1 max-w-lg group">
          <SearchIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300 group-focus-within:text-indigo-500 transition-colors" />
          <input
            v-model="searchQuery"
            @input="onSearchInput"
            type="text"
            placeholder="Cari penerima, tujuan, judul, atau pesan..."
            class="search-input-premium"
          />
        </div>
        <button @click="resetFilters" class="p-2.5 bg-white text-slate-500 hover:text-indigo-600 hover:bg-indigo-50 rounded-xl border border-slate-200 shadow-sm transition-all group" title="Reset Filter">
          <ResetIcon class="w-4 h-4 group-hover:rotate-180 transition-transform duration-500" />
        </button>
      </div>
    </Teleport>

    <!-- ─── TABLES PER CHANNEL ─── -->
    <div class="space-y-8">
      <div
        v-for="table in channelTables"
        :key="table.channel"
        class="bg-white border border-slate-100 rounded-2xl shadow-sm overflow-hidden"
      >
        <!-- Table Header -->
        <div class="px-6 py-4 border-b border-slate-100 flex flex-wrap items-center justify-between gap-4">
          <div class="flex items-center gap-3">
            <div
              class="w-9 h-9 rounded-xl flex items-center justify-center shadow-sm"
              :class="table.channel === 'whatsapp' ? 'bg-emerald-500 text-white' : 'bg-sky-500 text-white'"
            >
              <component :is="table.icon" class="w-4.5 h-4.5" />
            </div>
            <div>
              <h2 class="text-sm font-black text-slate-800 tracking-tight leading-none">{{ table.title }}</h2>
              <p class="text-[10px] font-bold text-slate-400 mt-0.5">
                {{ totalByChannel[table.channel].toLocaleString('id-ID') }} entri ditemukan
              </p>
            </div>
          </div>

          <!-- Status Tabs + Limit Selector -->
          <div class="flex items-center gap-2 flex-wrap">
            <div class="flex items-center rounded-xl border border-slate-200 bg-slate-50 p-0.5 gap-0.5">
              <button
                v-for="status in statusOptions(table.channel)"
                :key="status.value"
                @click="selectStatus(table.channel, status.value)"
                :class="[
                  'px-3 py-1.5 rounded-lg text-[10px] font-bold uppercase tracking-wider transition-all',
                  selectedStatusByChannel[table.channel] === status.value
                    ? 'bg-white text-slate-800 shadow-sm'
                    : 'text-slate-400 hover:text-slate-600'
                ]"
              >
                {{ status.label }}
              </button>
            </div>

            <select
              v-model.number="limitByChannel[table.channel]"
              @change="onLimitChange(table.channel)"
              class="limit-select"
            >
              <option v-for="n in [10, 25, 50, 100]" :key="n" :value="n">{{ n }} / hal</option>
            </select>
          </div>
        </div>

        <!-- Loading State -->
        <div v-if="loadingByChannel[table.channel]" class="py-16 flex flex-col items-center justify-center gap-3">
          <div class="w-8 h-8 border-4 border-indigo-100 border-t-indigo-500 rounded-full animate-spin"></div>
          <p class="text-xs font-bold text-slate-400">Memuat log {{ table.badge }}…</p>
        </div>

        <!-- Empty State -->
        <div v-else-if="!logsByChannel[table.channel]?.length" class="py-16 flex flex-col items-center justify-center gap-3">
          <div class="w-14 h-14 rounded-2xl flex items-center justify-center" :class="table.channel === 'whatsapp' ? 'bg-emerald-50' : 'bg-sky-50'">
            <InboxIcon class="w-7 h-7" :class="table.channel === 'whatsapp' ? 'text-emerald-300' : 'text-sky-300'" />
          </div>
          <p class="text-sm font-bold text-slate-400">Belum ada log {{ table.badge }}</p>
          <p class="text-xs text-slate-300">Notifikasi yang terkirim akan muncul di sini</p>
        </div>

        <!-- Table -->
        <div v-else class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr class="border-b border-slate-100 bg-slate-50/60">
                <th class="px-6 py-3.5 text-[9px] font-black text-slate-400 uppercase tracking-widest whitespace-nowrap">Waktu</th>
                <th class="px-6 py-3.5 text-[9px] font-black text-slate-400 uppercase tracking-widest whitespace-nowrap">Penerima</th>
                <th class="px-6 py-3.5 text-[9px] font-black text-slate-400 uppercase tracking-widest whitespace-nowrap">Kontak</th>
                <th class="px-6 py-3.5 text-[9px] font-black text-slate-400 uppercase tracking-widest">Pesan</th>
                <th class="px-6 py-3.5 text-[9px] font-black text-slate-400 uppercase tracking-widest whitespace-nowrap">Status</th>
                <th class="px-6 py-3.5 text-[9px] font-black text-slate-400 uppercase tracking-widest text-right whitespace-nowrap">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-50">
              <tr
                v-for="log in logsByChannel[table.channel]"
                :key="log.id"
                class="hover:bg-slate-50/70 transition-colors group"
              >
                <!-- Waktu -->
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-xs font-bold text-slate-700">{{ formatDateShort(log.created_at) }}</div>
                  <div class="text-[10px] text-slate-400 mt-0.5">{{ formatTime(log.created_at) }}</div>
                </td>

                <!-- Penerima -->
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center gap-2.5">
                    <div class="w-7 h-7 rounded-full flex items-center justify-center text-white text-[10px] font-black shrink-0" :class="table.channel === 'whatsapp' ? 'bg-emerald-400' : 'bg-sky-400'">
                      {{ (log.recipient_name || 'OT').charAt(0).toUpperCase() }}
                    </div>
                    <div class="text-xs font-bold text-slate-700 max-w-[120px] truncate">{{ log.recipient_name || 'Orang Tua' }}</div>
                  </div>
                </td>

                <!-- Kontak -->
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-xs font-mono font-semibold text-slate-500">
                    {{ table.channel === 'whatsapp' ? log.recipient_phone : log.recipient_email }}
                  </div>
                </td>

                <!-- Pesan -->
                <td class="px-6 py-4 max-w-xs">
                  <div class="text-xs font-bold text-slate-800 leading-snug line-clamp-1">
                    {{ log.title || humanReadableMessage(log.message) }}
                  </div>
                  <div class="text-[10px] text-slate-400 mt-0.5 line-clamp-1 leading-snug">
                    {{ humanReadableMessage(log.message) }}
                  </div>
                  <!-- Error alert untuk yang gagal -->
                  <div v-if="log.delivery_status === 'failed' && log.delivery_error" class="mt-1.5 text-[10px] font-bold text-rose-500 flex items-center gap-1">
                    <AlertCircleIcon class="w-3 h-3 shrink-0" />
                    <span class="line-clamp-1">{{ log.delivery_error }}</span>
                  </div>
                </td>

                <!-- Status Badge -->
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center gap-1.5">
                    <div class="w-1.5 h-1.5 rounded-full shrink-0" :class="getStatusDotClass(log.delivery_status)"></div>
                    <span class="px-2.5 py-1 rounded-full text-[9px] font-black uppercase tracking-wider border" :class="getStatusBadgeClass(log.delivery_status)">
                      {{ statusLabel(log.delivery_status, table.channel) }}
                    </span>
                  </div>
                </td>

                <!-- Aksi -->
                <td class="px-6 py-4 whitespace-nowrap text-right">
                  <div class="flex items-center justify-end gap-1">
                    <button
                      @click="openDetailModal(log, table.channel)"
                      class="p-1.5 text-slate-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors"
                      title="Lihat Detail"
                    >
                      <EyeIcon class="w-4 h-4" />
                    </button>
                    <button
                      v-if="canResend(log)"
                      @click="resendNotification(log)"
                      :disabled="resendingId === log.id"
                      class="flex items-center gap-1 px-2.5 py-1.5 bg-rose-50 text-rose-600 hover:bg-rose-500 hover:text-white border border-rose-100 hover:border-rose-500 rounded-lg text-[10px] font-black transition-all disabled:opacity-60 disabled:cursor-not-allowed"
                      title="Kirim Ulang"
                    >
                      <RefreshIcon class="w-3 h-3" :class="{ 'animate-spin': resendingId === log.id }" />
                      <span>{{ resendingId === log.id ? 'Mengirim…' : 'Kirim Ulang' }}</span>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Pagination Footer -->
        <div
          v-if="totalByChannel[table.channel] > 0"
          class="px-6 py-4 border-t border-slate-100 flex items-center justify-between bg-slate-50/40"
        >
          <p class="text-[10px] font-bold text-slate-400">
            Menampilkan
            <span class="text-slate-600">{{ rangeStart(table.channel) }}–{{ rangeEnd(table.channel) }}</span>
            dari <span class="text-slate-600">{{ totalByChannel[table.channel].toLocaleString('id-ID') }}</span> entri
          </p>

          <div class="flex items-center gap-1">
            <button
              @click="gotoPage(table.channel, 1)"
              :disabled="pageByChannel[table.channel] === 1"
              class="pagination-btn"
              title="Halaman Pertama"
            >
              <ChevronsLeftIcon class="w-3.5 h-3.5" />
            </button>
            <button
              @click="gotoPage(table.channel, pageByChannel[table.channel] - 1)"
              :disabled="pageByChannel[table.channel] === 1"
              class="pagination-btn"
            >
              <PrevIcon class="w-3.5 h-3.5" />
            </button>

            <div class="flex items-center gap-1 mx-1">
              <button
                v-for="p in visiblePages(table.channel)"
                :key="p"
                @click="p !== '…' && gotoPage(table.channel, p)"
                :class="[
                  'min-w-[32px] h-8 text-[11px] font-black rounded-lg transition-all',
                  p === pageByChannel[table.channel]
                    ? 'bg-indigo-600 text-white shadow-md shadow-indigo-100'
                    : p === '…'
                    ? 'text-slate-300 cursor-default'
                    : 'text-slate-500 hover:bg-slate-100'
                ]"
              >{{ p }}</button>
            </div>

            <button
              @click="gotoPage(table.channel, pageByChannel[table.channel] + 1)"
              :disabled="pageByChannel[table.channel] >= totalPages(table.channel)"
              class="pagination-btn"
            >
              <NextIcon class="w-3.5 h-3.5" />
            </button>
            <button
              @click="gotoPage(table.channel, totalPages(table.channel))"
              :disabled="pageByChannel[table.channel] >= totalPages(table.channel)"
              class="pagination-btn"
              title="Halaman Terakhir"
            >
              <ChevronsRightIcon class="w-3.5 h-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- ─── DETAIL MODAL ─── -->
    <Teleport to="body">
      <transition name="page">
        <div v-if="showDetailModal && selectedLog" class="fixed inset-0 z-[200] flex items-center justify-center p-4 sm:p-6">
          <div class="absolute inset-0 bg-slate-900/50 backdrop-blur-sm" @click="showDetailModal = false"></div>
          <div class="bg-white w-full max-w-lg relative z-10 rounded-2xl shadow-2xl overflow-hidden animate-scale-in">
            <!-- Modal Header -->
            <div class="px-6 py-5 border-b border-slate-100 flex items-center justify-between">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-xl flex items-center justify-center" :class="selectedLog.channel === 'whatsapp' ? 'bg-emerald-500 text-white' : 'bg-sky-500 text-white'">
                  <component :is="selectedLog.channel === 'whatsapp' ? MessageCircleIcon : MailIcon" class="w-5 h-5" />
                </div>
                <div>
                  <h3 class="text-sm font-black text-slate-800 tracking-tight">Detail Pesan Notifikasi</h3>
                  <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mt-0.5">
                    Via {{ selectedLog.channel === 'whatsapp' ? 'WhatsApp' : 'Email' }}
                  </p>
                </div>
              </div>
              <button @click="showDetailModal = false" class="p-2 text-slate-400 hover:text-slate-600 hover:bg-slate-100 rounded-xl transition-colors">
                <XIcon class="w-4 h-4" />
              </button>
            </div>

            <!-- Modal Body -->
            <div class="p-6 space-y-4">
              <!-- Status Row -->
              <div class="flex items-center justify-between p-4 bg-slate-50 rounded-xl border border-slate-100">
                <div class="flex items-center gap-2">
                  <div class="w-2 h-2 rounded-full" :class="getStatusDotClass(selectedLog.delivery_status)"></div>
                  <span class="text-xs font-bold text-slate-600">Status Pengiriman</span>
                </div>
                <span class="px-3 py-1 rounded-full text-[9px] font-black uppercase tracking-wider border" :class="getStatusBadgeClass(selectedLog.delivery_status)">
                  {{ statusLabel(selectedLog.delivery_status, selectedLog.channel) }}
                </span>
              </div>

              <!-- Info Grid -->
              <div class="grid grid-cols-2 gap-3">
                <div class="p-3.5 bg-slate-50 rounded-xl border border-slate-100">
                  <label class="field-label">Penerima</label>
                  <p class="field-value">{{ selectedLog.recipient_name || 'Orang Tua' }}</p>
                </div>
                <div class="p-3.5 bg-slate-50 rounded-xl border border-slate-100">
                  <label class="field-label">Kontak Tujuan</label>
                  <p class="field-value font-mono text-[11px]">
                    {{ selectedLog.channel === 'whatsapp' ? (selectedLog.recipient_phone || '-') : (selectedLog.recipient_email || '-') }}
                  </p>
                </div>
                <div class="p-3.5 bg-slate-50 rounded-xl border border-slate-100">
                  <label class="field-label">Waktu Kirim</label>
                  <p class="field-value">{{ formatDate(selectedLog.created_at) }}</p>
                </div>
                <div class="p-3.5 bg-slate-50 rounded-xl border border-slate-100">
                  <label class="field-label">Jenis Pesan</label>
                  <p class="field-value capitalize">{{ selectedLog.type || 'notifikasi' }}</p>
                </div>
              </div>

              <!-- Judul -->
              <div class="p-3.5 bg-slate-50 rounded-xl border border-slate-100" v-if="selectedLog.title">
                <label class="field-label">Judul / Subjek</label>
                <p class="field-value text-slate-800 font-bold">{{ selectedLog.title }}</p>
              </div>

              <!-- Error Alert -->
              <div v-if="selectedLog.delivery_status === 'failed' && selectedLog.delivery_error" class="p-3.5 bg-rose-50 rounded-xl border border-rose-100">
                <label class="field-label text-rose-400">Alasan Gagal</label>
                <p class="text-xs font-bold text-rose-600 mt-1 leading-relaxed">{{ selectedLog.delivery_error }}</p>
              </div>

              <!-- Isi Pesan -->
              <div class="p-3.5 bg-slate-50 rounded-xl border border-slate-100">
                <label class="field-label">Isi Pesan</label>
                <div class="text-xs font-semibold text-slate-700 whitespace-pre-wrap leading-relaxed mt-1 max-h-40 overflow-y-auto">{{ selectedLog.message }}</div>
              </div>
            </div>

            <!-- Modal Footer -->
            <div class="px-6 py-4 border-t border-slate-100 bg-slate-50/50 flex items-center gap-3">
              <button
                @click="showDetailModal = false"
                class="flex-1 py-2.5 border border-slate-200 text-slate-600 rounded-xl font-black text-[10px] uppercase tracking-widest hover:bg-slate-50 transition-all"
              >Tutup</button>
              <button
                v-if="canResend(selectedLog)"
                @click="resendNotification(selectedLog); showDetailModal = false"
                :disabled="resendingId === selectedLog?.id"
                class="flex-1 py-2.5 bg-rose-500 text-white rounded-xl font-black text-[10px] uppercase tracking-widest hover:bg-rose-600 transition-all shadow-lg shadow-rose-100 flex items-center justify-center gap-1.5 disabled:opacity-60"
              >
                <RefreshIcon class="w-3.5 h-3.5" :class="{ 'animate-spin': resendingId === selectedLog?.id }" />
                <span>Kirim Ulang</span>
              </button>
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
  ChevronsLeft as ChevronsLeftIcon,
  ChevronsRight as ChevronsRightIcon,
  RotateCcw as ResetIcon,
  Send as SendIcon,
  Inbox as InboxIcon,
  X as XIcon
} from 'lucide-vue-next'

const isMounted = ref(false)
const logsByChannel = reactive({ whatsapp: [], email: [] })
const loadingByChannel = reactive({ whatsapp: false, email: false })
const totalByChannel = reactive({ whatsapp: 0, email: 0 })
const pageByChannel = reactive({ whatsapp: 1, email: 1 })
const limitByChannel = reactive({ whatsapp: 10, email: 10 })
const stats = ref({})
const searchQuery = ref('')
const selectedStatusByChannel = reactive({ whatsapp: '', email: '' })
const resendingId = ref(null)
const toast = useToast()
const showDetailModal = ref(false)
const selectedLog = ref(null)
const selectedChannelLabel = ref('')

const channelTables = [
  { channel: 'whatsapp', title: 'Efikasi WhatsApp', badge: 'WhatsApp', icon: MessageCircleIcon },
  { channel: 'email', title: 'Efikasi Email', badge: 'Email', icon: MailIcon }
]

const whatsappStatusTabs = [
  { label: 'Semua', value: '' },
  { label: 'Menunggu', value: 'pending' },
  { label: 'Terkirim', value: 'sent' },
  { label: 'Diterima', value: 'delivered' },
  { label: 'Dibaca', value: 'read' },
  { label: 'Gagal', value: 'failed' }
]

const emailStatusTabs = [
  { label: 'Semua', value: '' },
  { label: 'Menunggu', value: 'pending' },
  { label: 'Terkirim', value: 'sent' },
  { label: 'Gagal', value: 'failed' }
]

const statusOptions = (channel) => channel === 'email' ? emailStatusTabs : whatsappStatusTabs

const statCards = computed(() => [
  {
    label: 'Total Terkirim',
    value: stats.value.SENT || 0,
    icon: SendIcon,
    iconBg: 'bg-indigo-50',
    iconColor: 'text-indigo-600',
    gradientBg: 'bg-gradient-to-br from-indigo-50/40 to-transparent',
    badge: 'WA + Email',
    badgeClass: 'bg-indigo-50 text-indigo-600 border-indigo-100'
  },
  {
    label: 'WA Diterima',
    value: stats.value.DELIVERED || 0,
    icon: CheckIcon,
    iconBg: 'bg-sky-50',
    iconColor: 'text-sky-600',
    gradientBg: 'bg-gradient-to-br from-sky-50/40 to-transparent',
    badge: 'WhatsApp',
    badgeClass: 'bg-sky-50 text-sky-600 border-sky-100'
  },
  {
    label: 'WA Dibaca',
    value: stats.value.READ || 0,
    icon: CheckIcon,
    iconBg: 'bg-emerald-50',
    iconColor: 'text-emerald-600',
    gradientBg: 'bg-gradient-to-br from-emerald-50/40 to-transparent',
    badge: 'WhatsApp',
    badgeClass: 'bg-emerald-50 text-emerald-600 border-emerald-100'
  },
  {
    label: 'Gagal Kirim',
    value: stats.value.FAILED || 0,
    icon: AlertCircleIcon,
    iconBg: 'bg-rose-50',
    iconColor: 'text-rose-600',
    gradientBg: 'bg-gradient-to-br from-rose-50/40 to-transparent',
    badge: 'Perlu Aksi',
    badgeClass: 'bg-rose-50 text-rose-600 border-rose-100'
  }
])

let searchTimeout = null

const onSearchInput = () => {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => fetchLogs(1), 350)
}

const onLimitChange = (channel) => {
  pageByChannel[channel] = 1
  fetchOneChannel(channel, 1)
}

const selectStatus = (channel, status) => {
  selectedStatusByChannel[channel] = status
  fetchLogs(1, channel)
}

const resetFilters = () => {
  searchQuery.value = ''
  selectedStatusByChannel.whatsapp = ''
  selectedStatusByChannel.email = ''
  fetchLogs(1)
}

const handleNotificationStatusChanged = () => {
  fetchStats()
  fetchLogs()
}

const totalPages = (channel) => Math.max(1, Math.ceil(totalByChannel[channel] / limitByChannel[channel]))

const rangeStart = (channel) => {
  if (totalByChannel[channel] === 0) return 0
  return (pageByChannel[channel] - 1) * limitByChannel[channel] + 1
}

const rangeEnd = (channel) => {
  return Math.min(pageByChannel[channel] * limitByChannel[channel], totalByChannel[channel])
}

const visiblePages = (channel) => {
  const total = totalPages(channel)
  const current = pageByChannel[channel]
  if (total <= 7) return Array.from({ length: total }, (_, i) => i + 1)
  const pages = []
  if (current <= 4) {
    pages.push(1, 2, 3, 4, 5, '…', total)
  } else if (current >= total - 3) {
    pages.push(1, '…', total - 4, total - 3, total - 2, total - 1, total)
  } else {
    pages.push(1, '…', current - 1, current, current + 1, '…', total)
  }
  return pages
}

const gotoPage = (channel, page) => {
  const p = Math.max(1, Math.min(page, totalPages(channel)))
  fetchOneChannel(channel, p)
}

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
        limit: limitByChannel[channel],
        status: selectedStatusByChannel[channel],
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

const canResend = (log) => String(log?.delivery_status || '').toLowerCase() === 'failed'

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

// Membuat pesan yang bisa dibaca manusia - membersihkan JSON / template literal
const humanReadableMessage = (msg) => {
  if (!msg) return '-'
  try {
    // Jika pesan berbentuk JSON string, coba parse
    const parsed = JSON.parse(msg)
    if (typeof parsed === 'string') return parsed
    if (parsed.text) return parsed.text
    if (parsed.body) return parsed.body
    return msg
  } catch {
    // Bukan JSON - bersihkan tag HTML dan backtick
    return msg.replace(/<[^>]*>/g, '').replace(/\*([^*]+)\*/g, '$1').trim()
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('id-ID', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit'
  })
}

const formatDateShort = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleDateString('id-ID', {
    day: '2-digit', month: 'short', year: 'numeric'
  })
}

const formatTime = (dateStr) => {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

const statusLabel = (status, channel = 'whatsapp') => {
  const norm = String(status || 'sent').toLowerCase()
  return statusOptions(channel).find((item) => item.value === norm)?.label || norm.toUpperCase()
}

const getStatusDotClass = (status) => {
  const norm = String(status || '').toLowerCase()
  switch (norm) {
    case 'read': return 'bg-emerald-500'
    case 'delivered': return 'bg-sky-500'
    case 'failed': return 'bg-rose-500'
    case 'pending': return 'bg-amber-400'
    default: return 'bg-slate-400'
  }
}

const getStatusBadgeClass = (status) => {
  const norm = String(status || '').toLowerCase()
  switch (norm) {
    case 'read': return 'bg-emerald-50 text-emerald-700 border-emerald-200'
    case 'delivered': return 'bg-sky-50 text-sky-700 border-sky-200'
    case 'failed': return 'bg-rose-50 text-rose-700 border-rose-200'
    case 'pending': return 'bg-amber-50 text-amber-700 border-amber-200'
    default: return 'bg-slate-50 text-slate-600 border-slate-200'
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
  @apply w-full bg-white border border-slate-200 rounded-xl py-2.5 pl-11 pr-4 text-xs font-bold text-slate-700 outline-none transition-all focus:border-indigo-500 focus:ring-4 focus:ring-indigo-50 shadow-sm;
}
.field-label {
  @apply text-[9px] font-black text-slate-400 uppercase tracking-widest block mb-1;
}
.field-value {
  @apply text-xs font-bold text-slate-700;
}
.limit-select {
  @apply h-9 pl-3 pr-8 text-[11px] font-bold text-slate-600 bg-white border border-slate-200 rounded-xl outline-none cursor-pointer transition-all focus:border-indigo-400 focus:ring-4 focus:ring-indigo-50 shadow-sm appearance-none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%2394a3b8' stroke-width='2.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 10px center;
}
.pagination-btn {
  @apply p-2 bg-white text-slate-500 hover:text-indigo-600 hover:bg-indigo-50 disabled:opacity-40 disabled:cursor-not-allowed rounded-lg border border-slate-200 shadow-sm transition-all;
}
</style>
