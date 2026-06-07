<script setup>
import { onMounted, onUnmounted, ref, computed, watch } from 'vue'
import { 
  MessageCircle, 
  CheckCircle2, 
  UserCheck, 
  RefreshCw, 
  ExternalLink,
  Bell,
  Link as LinkIcon,
  ChevronLeft as PrevIcon,
  ChevronRight as NextIcon
} from 'lucide-vue-next'
import supportService from '../../services/support.service'
import { useToast } from '../../composables/useToast'
import { useRoute } from 'vue-router'

const toast = useToast()
const route = useRoute()
const conversations = ref([])
const total = ref(0)
const page = ref(1)
const limit = ref(10)
const loading = ref(false)
const status = ref('open')
let lastNotificationAt = 0

const statusOptions = [
  { value: 'open', label: 'Menunggu Admin' },
  { value: 'pending', label: 'Sedang Ditangani' },
  { value: 'closed', label: 'Selesai' }
]

// Memutar bunyi bel notifikasi
const playNotificationSound = () => {
  try {
    const audioContext = new (window.AudioContext || window.webkitAudioContext)()
    const osc = audioContext.createOscillator()
    const gain = audioContext.createGain()
    osc.connect(gain)
    gain.connect(audioContext.destination)
    
    // Bunyi ting-ting yang ramah
    osc.frequency.setValueAtTime(880, audioContext.currentTime) // A5 note
    gain.gain.setValueAtTime(0.3, audioContext.currentTime)
    osc.start()
    
    // Bunyi kedua 0.15 detik setelahnya
    setTimeout(() => {
      osc.frequency.setValueAtTime(1200, audioContext.currentTime)
    }, 150)

    osc.stop(audioContext.currentTime + 0.35)
  } catch (e) {
    console.error('Gagal memainkan bunyi bel notifikasi', e)
  }
}

const loadConversations = async () => {
  loading.value = true
  try {
    const res = await supportService.getConversations({ 
      status: status.value, 
      page: page.value, 
      limit: limit.value 
    })
    conversations.value = res.data.data.data || []
    total.value = res.data.data.total || 0
  } catch (err) {
    toast.error('Gagal memuat antrean CS', err.response?.data?.message || 'Server tidak merespon')
  } finally {
    loading.value = false
  }
}

const assignTicket = async (id) => {
  try {
    await supportService.assign(id)
    toast.success('Tiket diambil', 'Tiket telah masuk ke daftar penanganan Anda')
    await loadConversations()
  } catch (err) {
    toast.error('Gagal mengambil tiket', err.response?.data?.message || 'Server tidak merespon')
  }
}

const closeTicket = async (id) => {
  try {
    await supportService.close(id)
    toast.success('Tiket selesai', 'Status bantuan ditutup, bot otomatis kembali normal')
    await loadConversations()
  } catch (err) {
    toast.error('Gagal menutup tiket', err.response?.data?.message || 'Server tidak merespon')
  }
}

const updateTicketStatus = async (item, nextStatus) => {
  if (!item || item.status === nextStatus) return
  const previousStatus = item.status
  item.status = nextStatus
  try {
    await supportService.updateStatus(item.id, nextStatus)
    toast.success('Status diperbarui', `Tiket sekarang ${statusOptions.find(s => s.value === nextStatus)?.label || nextStatus}`)
    await loadConversations()
  } catch (err) {
    item.status = previousStatus
    toast.error('Gagal memperbarui status', err.response?.data?.message || 'Server tidak merespon')
  }
}

// Buka link chat WhatsApp Web resmi langsung ke target nomor
const buildWhatsAppWebURL = (phoneNumber) => {
  let cleanNumber = phoneNumber.replace(/[^0-9]/g, '')
  if (cleanNumber.startsWith('0')) {
    cleanNumber = '62' + cleanNumber.substring(1)
  }
  const text = encodeURIComponent('Halo Bapak/Ibu, ada yang bisa kami bantu terkait SchoolPay?')
  return `https://web.whatsapp.com/send?phone=${cleanNumber}&text=${text}`
}

const openWhatsAppWeb = (item) => {
  const url = item.whatsapp_web_url || buildWhatsAppWebURL(item.phone_number)
  window.open(url, '_blank')
}

const copyAdminLink = async (item) => {
  const adminURL = `${window.location.origin}/support/chat?status=${item.status || 'open'}`
  try {
    await navigator.clipboard.writeText(adminURL)
    toast.success('Link CS disalin', 'Link antrean admin siap dibagikan.')
  } catch (err) {
    toast.info('Link UI CS', adminURL)
  }
}

const handleSupportUpdate = async (event) => {
  const now = Date.now()
  if (now - lastNotificationAt > 8000) {
    playNotificationSound()
    toast.info('Antrean CS Baru', `Ada panggilan bantuan masuk dari ${event.detail?.parent_name || 'Orang Tua'}`)
    lastNotificationAt = now
  }
  await loadConversations()
}

// Computed Properties for Pagination
const totalPages = computed(() => Math.ceil(total.value / limit.value) || 1)
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

// Watchers for resetting page on filters change
watch(status, () => {
  page.value = 1
  loadConversations()
})

onMounted(() => {
  const queryStatus = String(route.query.status || '')
  if (statusOptions.some(option => option.value === queryStatus)) {
    status.value = queryStatus
  }
  loadConversations()
  window.addEventListener('support-chat-updated', handleSupportUpdate)
})

onUnmounted(() => {
  window.removeEventListener('support-chat-updated', handleSupportUpdate)
})
</script>

<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between gap-4">
      <div>
        <h1 class="text-lg font-black text-slate-800 tracking-tight flex items-center gap-2">
          <Bell class="w-5 h-5 text-amber-500" />
          Antrean Panggilan CS WhatsApp
        </h1>
        <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">
          Hubungi wali murid via WhatsApp Web resmi secara langsung saat lonceng berbunyi
        </p>
      </div>
      <button @click="loadConversations" class="px-3 py-1.5 bg-white border border-slate-200 text-slate-600 hover:bg-slate-50 font-bold rounded-xl text-[10px] flex items-center gap-1.5 transition-all shadow-sm">
        <RefreshCw class="w-3.5 h-3.5" /> Perbarui
      </button>
    </div>

    <div class="bg-white border border-slate-200 rounded-xl overflow-hidden shadow-sm flex flex-col min-h-[550px]">
      <div class="p-4 border-b border-slate-100 flex items-center justify-between bg-slate-50/30">
        <div class="flex items-center gap-2 font-black text-slate-700 text-xs uppercase tracking-widest">
          <MessageCircle class="w-3.5 h-3.5 text-indigo-600 animate-pulse" /> Status Antrean Tiket
        </div>
        <select v-model="status" class="text-[10px] font-bold bg-white border border-slate-200 rounded-lg px-2 py-1.5 outline-none cursor-pointer focus:ring-2 focus:ring-indigo-500/20 transition-all">
          <option v-for="option in statusOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
        </select>
      </div>

      <div class="flex-1 overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-slate-50/50 border-b border-slate-100 text-[10px] font-black text-slate-400 uppercase tracking-widest">
              <th class="py-3 px-4 w-12 text-center">No</th>
              <th class="py-3 px-4">Orang Tua / Wali</th>
              <th class="py-3 px-4">Siswa</th>
              <th class="py-3 px-4">Nomor WhatsApp</th>
              <th class="py-3 px-4">Waktu Masuk</th>
              <th class="py-3 px-4">Status</th>
              <th class="py-3 px-4 text-center w-[300px]">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, idx) in conversations" :key="item.id" class="border-b border-slate-100 hover:bg-slate-50/30 transition-all text-xs font-semibold text-slate-600">
              <td class="py-3 px-4 text-center text-slate-400 font-bold">{{ (page - 1) * limit + idx + 1 }}</td>
              <td class="py-3 px-4">
                <div class="font-black text-slate-800 text-xs">{{ item.parent_name || 'Orang Tua / Wali Siswa' }}</div>
              </td>
              <td class="py-3 px-4">
                <div class="text-slate-600 text-[11px] font-bold">{{ item.student_names || '-' }}</div>
              </td>
              <td class="py-3 px-4">
                <span class="bg-slate-100 text-slate-700 px-2 py-0.5 rounded text-[10px] font-mono font-bold">{{ item.phone_number }}</span>
              </td>
              <td class="py-3 px-4 text-slate-500 text-[11px]">{{ new Date(item.created_at).toLocaleString('id-ID') }}</td>
              <td class="py-3 px-4">
                <select
                  :value="item.status"
                  @change="updateTicketStatus(item, $event.target.value)"
                  class="text-[10px] font-black bg-white border border-slate-200 rounded-lg px-2 py-1.5 outline-none cursor-pointer focus:ring-2 focus:ring-indigo-500/20 transition-all"
                >
                  <option v-for="option in statusOptions" :key="option.value" :value="option.value">{{ option.label }}</option>
                </select>
              </td>
              <td class="py-3 px-4 text-center">
                <div class="flex items-center justify-center gap-1.5 flex-nowrap">
                  <button @click="openWhatsAppWeb(item)" class="px-2.5 py-1.5 bg-emerald-600 text-white hover:bg-emerald-700 font-bold rounded-lg text-[9px] uppercase tracking-wider flex items-center gap-1 transition-all shadow-sm whitespace-nowrap">
                    <ExternalLink class="w-3 h-3" /> WA Web
                  </button>
                  <button @click="copyAdminLink(item)" class="px-2.5 py-1.5 bg-white text-slate-600 border border-slate-200 hover:bg-slate-50 font-bold rounded-lg text-[9px] uppercase tracking-wider flex items-center gap-1 transition-all shadow-sm whitespace-nowrap">
                    <LinkIcon class="w-3 h-3" /> Link CS
                  </button>
                  <button v-if="item.status === 'open'" @click="assignTicket(item.id)" class="px-2.5 py-1.5 bg-white text-slate-600 border border-slate-200 hover:bg-slate-50 font-bold rounded-lg text-[9px] uppercase tracking-wider flex items-center gap-1 transition-all shadow-sm whitespace-nowrap">
                    <UserCheck class="w-3 h-3" /> Tangani
                  </button>
                  <button v-if="item.status !== 'closed'" @click="closeTicket(item.id)" class="px-2.5 py-1.5 bg-white text-emerald-600 border border-slate-200 hover:bg-emerald-50 font-bold rounded-lg text-[9px] uppercase tracking-wider flex items-center gap-1 transition-all shadow-sm whitespace-nowrap">
                    <CheckCircle2 class="w-3.5 h-3.5 text-emerald-600" /> Selesai
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="!loading && conversations.length === 0">
              <td colspan="7" class="p-12 text-center">
                <MessageCircle class="w-10 h-10 mx-auto mb-2 text-slate-200 animate-pulse" />
                <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">
                  Tidak ada panggilan bantuan aktif saat ini.
                </p>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div class="px-6 py-4 bg-slate-50/50 border-t border-slate-100 flex items-center justify-between">
        <div class="flex items-center gap-6">
          <div class="flex items-center gap-3">
            <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Tampilkan</span>
            <select v-model="limit" @change="page = 1; loadConversations()" class="bg-white border border-slate-200 rounded-lg text-[10px] font-black text-slate-600 px-2 py-1 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all cursor-pointer shadow-sm">
              <option :value="10">10</option>
              <option :value="25">25</option>
              <option :value="50">50</option>
              <option :value="100">100</option>
            </select>
          </div>
          <div class="h-8 w-px bg-slate-200 hidden sm:block"></div>
          <span class="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em]">
            Halaman <span class="text-indigo-600">{{ page }}</span> dari {{ totalPages }} <span class="mx-2 text-slate-300">|</span> Total <span class="text-indigo-600">{{ total }}</span> Data
          </span>
        </div>
        <!-- Pagination Control -->
        <div class="flex items-center gap-2">
          <!-- Previous Button -->
          <button 
            v-if="totalPages > 1"
            @click="page > 1 && (page--) && loadConversations()" 
            :disabled="page <= 1" 
            class="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 bg-white text-slate-400 hover:text-indigo-600 hover:border-indigo-100 hover:bg-indigo-50/30 disabled:opacity-20 disabled:hover:bg-white disabled:hover:border-slate-200 transition-all cursor-pointer"
          >
            <PrevIcon class="w-3.5 h-3.5" />
          </button>

          <!-- Page Numbers -->
          <div class="flex items-center gap-1">
            <button 
              v-for="p in visiblePages" 
              :key="p"
              @click="page = p; loadConversations()"
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
            @click="page < totalPages && (page++) && loadConversations()" 
            :disabled="page >= totalPages" 
            class="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 bg-white text-slate-400 hover:text-indigo-600 hover:border-indigo-100 hover:bg-indigo-50/30 disabled:opacity-20 disabled:hover:bg-white disabled:hover:border-slate-200 transition-all cursor-pointer"
          >
            <NextIcon class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
