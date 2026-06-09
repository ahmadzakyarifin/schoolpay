<template>
  <div class="flex min-h-screen bg-[#f8fafc]">
    <!-- Sidebar -->
    <aside class="admin-sidebar flex flex-col justify-between">
      <div class="flex items-center gap-3 mb-10 px-2 shrink-0">
        <div class="w-12 h-12 bg-indigo-600 rounded-2xl flex items-center justify-center shadow-xl shadow-indigo-200">
          <GraduationCapIcon class="text-white w-6 h-6" />
        </div>
        <div>
          <h1 class="text-xl font-black text-slate-800 tracking-tight">SchoolPay</h1>
          <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">{{ panelLabel }}</p>
        </div>
      </div>
      
      <nav class="flex-1 space-y-2 overflow-y-auto px-2 custom-scrollbar">
        <!-- Admin Menu -->
        <template v-if="authStore.user?.role === 'admin'">
          <template v-for="group in computedMenuGroups" :key="group.title">
            <!-- Jika isDirectGroup (termasuk hasil Auto-Flatten) -->
            <template v-if="group.isDirectGroup">
              <router-link v-for="child in group.children" :key="child.name" :to="child.to" class="nav-link" :exact-active-class="child.exact ? 'active' : ''" :active-class="!child.exact ? 'active' : ''">
                <component :is="child.icon" class="w-5 h-5" />
                <span>{{ child.label }}</span>
              </router-link>
            </template>

            <!-- Jika Dropdown Group (children > 1) -->
            <div v-else class="space-y-1">
              <button @click="openDropdowns[group.title] = !openDropdowns[group.title]" class="w-full flex items-center justify-between px-4 py-3 rounded-2xl text-slate-600 hover:bg-indigo-50/50 hover:text-indigo-600 transition-all font-bold text-xs group">
                <div class="flex items-center gap-2.5">
                  <component :is="group.icon" class="w-5 h-5 text-slate-400 group-hover:text-indigo-600 transition-colors" />
                  <span>{{ group.title }}</span>
                </div>
                <ChevronRightIcon class="w-4 h-4 text-slate-400 transition-transform duration-200" :class="{ 'rotate-90 text-indigo-600': openDropdowns[group.title] }" />
              </button>

              <!-- Sub-menu list -->
              <transition name="dropdown">
                <div v-if="openDropdowns[group.title]" class="pl-7 pr-1 py-1 space-y-1">
                  <router-link v-for="child in group.children" :key="child.name" :to="child.to" class="nav-link" active-class="active">
                    <component :is="child.icon" class="w-4 h-4 opacity-75" />
                    <span>{{ child.label }}</span>
                  </router-link>
                </div>
              </transition>
            </div>
          </template>
        </template>

        <!-- Parent Menu -->
        <template v-else-if="authStore.user?.role === 'parent'">
          <router-link to="/parent/dashboard" class="nav-link" active-class="active">
            <LayoutGridIcon class="w-5 h-5" />
            <span>Dashboard</span>
          </router-link>
          <router-link to="/parent/bills" class="nav-link" active-class="active">
            <CreditCardIcon class="w-5 h-5" />
            <span>Tagihan Ananda</span>
          </router-link>
          <router-link to="/parent/history" class="nav-link" active-class="active">
            <ReceiptIcon class="w-5 h-5" />
            <span>Riwayat Bayar</span>
          </router-link>
          <button @click="openParentSupport" class="nav-link w-full text-left">
            <MessageCircleIcon class="w-5 h-5" />
            <span>CS Admin</span>
          </button>
        </template>
      </nav>

      <div class="mt-auto pt-6 pb-4 px-2 shrink-0">
        <button @click="handleLogout" class="w-full flex items-center justify-center gap-3 px-6 py-4 rounded-2xl text-red-500 bg-red-50 hover:bg-red-100 transition-all duration-300 font-bold text-xs">
          <LogOutIcon class="w-4 h-4" />
          <span>Logout</span>
        </button>
      </div>
    </aside>

    <!-- Main Content -->
    <main class="flex-1" style="margin-left: 17rem">
      <header class="flex items-center justify-between bg-white/95 backdrop-blur-md sticky top-0 py-4 px-6 lg:px-10 z-[50] border-b border-slate-100 shadow-sm gap-6">
        <!-- Left: Page Title -->
        <div class="shrink-0 min-w-[150px]">
          <h2 class="text-lg font-black text-slate-800 tracking-tight capitalize leading-tight">{{ $route.name?.replace('-', ' ') || 'Management' }}</h2>
          <p class="text-slate-400 text-[10px] font-black uppercase tracking-widest mt-0.5">{{ headerSubtitle }}</p>
        </div>
        
        <!-- Middle: Search & Filters (Teleport Target) -->
        <div class="flex-1 max-w-3xl flex justify-center items-center">
          <div id="header-actions-target" class="w-full"></div>
        </div>

        <!-- Right: WAHA & Profile -->
        <div class="flex items-center gap-6 shrink-0 min-w-[150px] justify-end">
          <!-- WhatsApp Connection Status Dropdown (New) -->
          <div v-if="authStore.user?.role === 'admin'" class="relative group">
            <button class="flex items-center gap-3 p-2 rounded-xl bg-slate-50 border border-slate-100 hover:bg-white hover:shadow-sm transition-all">
              <div :class="['WORKING', 'CONNECTED'].includes(waStatus) ? 'bg-emerald-50 text-emerald-600' : waStatus === 'SCANNING' || waStatus === 'STARTING' ? 'bg-amber-50 text-amber-600' : 'bg-rose-50 text-rose-600'" class="w-8 h-8 rounded-lg flex items-center justify-center transition-all">
                <MessageCircleIcon class="w-4 h-4" />
              </div>
              <div class="text-left pr-2">
                <p class="text-[8px] font-black text-slate-400 uppercase tracking-widest leading-none mb-1">WA Status</p>
                <div class="flex items-center gap-2">
                  <p class="text-[10px] font-black uppercase tracking-tight" :class="['WORKING', 'CONNECTED'].includes(waStatus) ? 'text-emerald-600' : 'text-slate-600'">
                    {{ ['WORKING', 'CONNECTED'].includes(waStatus) ? 'Connected' : 'Disconnected' }}
                  </p>
                  <button @click.stop="fetchWAStatus" class="p-1 hover:bg-slate-200 rounded-md transition-all group/refresh" title="Refresh Status">
                    <RotateCcwIcon class="w-2.5 h-2.5 text-slate-400 group-hover/refresh:text-indigo-600" />
                  </button>
                </div>
              </div>
            </button>

            <!-- Dropdown Content -->
            <div class="absolute right-0 top-full mt-2 w-64 bg-white rounded-3xl shadow-2xl border border-slate-100 p-4 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-[100] translate-y-2 group-hover:translate-y-0">
              <div class="flex flex-col gap-3">
                <div class="flex items-center justify-between p-3 bg-slate-50 rounded-2xl border border-slate-100">
                  <div class="flex items-center gap-3">
                    <div :class="['WORKING', 'CONNECTED'].includes(waStatus) ? 'bg-emerald-50 text-emerald-600' : 'bg-rose-50 text-rose-600'" class="w-10 h-10 rounded-xl flex items-center justify-center">
                      <ActivityIcon class="w-5 h-5" />
                    </div>
                    <div>
                      <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Bot Service</p>
                      <p class="text-xs font-bold text-slate-700">{{ waStatus }}</p>
                    </div>
                  </div>
                </div>

                <router-link to="/support/chat" class="flex items-center justify-between p-3 hover:bg-indigo-50 rounded-2xl transition-all border border-transparent hover:border-indigo-100 group/link">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 bg-white rounded-xl flex items-center justify-center border border-slate-100 group-hover/link:border-indigo-200">
                      <MessageCircleIcon class="w-5 h-5 text-slate-400 group-hover/link:text-indigo-600" />
                    </div>
                    <div>
                      <p class="text-xs font-black text-slate-700">CS WhatsApp Inbox</p>
                      <p class="text-[10px] font-medium text-slate-400">Balas dari dashboard</p>
                    </div>
                  </div>
                  <ChevronRightIcon class="w-4 h-4 text-slate-300 group-hover/link:text-indigo-400" />
                </router-link>

                <div class="grid grid-cols-1 gap-2">
                  <button v-if="!['WORKING', 'CONNECTED'].includes(waStatus)" @click="openQRModal" class="w-full py-3 bg-indigo-600 text-white rounded-2xl font-black text-[10px] uppercase tracking-widest hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-100">
                    Hubungkan WhatsApp
                  </button>
                  <button v-else @click="changeWANumber" :disabled="waActionLoading" class="w-full py-3 bg-indigo-600 text-white rounded-2xl font-black text-[10px] uppercase tracking-widest hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-100 disabled:opacity-50">
                    Ganti Nomor WhatsApp
                  </button>
                  <button @click="restartWA" :disabled="waActionLoading" class="w-full py-3 bg-white text-slate-600 border border-slate-200 rounded-2xl font-black text-[10px] uppercase tracking-widest hover:bg-slate-50 transition-all disabled:opacity-50">
                    Restart Session
                  </button>
                  <button v-if="['WORKING', 'CONNECTED', 'STARTING', 'SCANNING'].includes(waStatus)" @click="disconnectWA" :disabled="waActionLoading" class="w-full py-3 bg-white text-rose-600 border border-rose-100 rounded-2xl font-black text-[10px] uppercase tracking-widest hover:bg-rose-50 transition-all disabled:opacity-50">
                    Logout WhatsApp
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Profile Link -->
          <router-link to="/profile" class="flex items-center gap-4 hover:opacity-80 transition-all group bg-slate-50/80 hover:bg-slate-100/80 px-4 py-2.5 rounded-2xl border border-slate-100">
            <div class="text-right hidden sm:block">
              <p class="text-sm font-black text-slate-800 leading-none group-hover:text-indigo-600 transition-colors">{{ authStore.user?.name }}</p>
              <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mt-1">{{ authStore.user?.role }}</p>
            </div>
          </router-link>
        </div>
      </header>

      <div class="w-full mx-auto px-6 py-8 lg:px-10">


        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" :key="$route.fullPath" />
          </transition>
        </router-view>
      </div>
    </main>

    <!-- Global Toast Notifications -->
    <ToastContainer />

    <!-- WhatsApp QR Modal -->
    <Teleport to="body">
      <transition name="page">
        <div v-if="showQRModal" class="fixed inset-0 z-[200] flex items-center justify-center p-6">
          <div class="absolute inset-0 bg-slate-900/40 backdrop-blur-sm" @click="showQRModal = false"></div>
          <div class="bg-white w-full max-w-sm relative z-10 rounded-[3rem] shadow-2xl overflow-hidden animate-scale-in">
            <div class="p-8 text-center space-y-6">
              <div class="w-20 h-20 bg-indigo-50 text-indigo-600 rounded-[2rem] flex items-center justify-center mx-auto shadow-inner">
                <MessageCircleIcon class="w-10 h-10" />
              </div>
              <div>
                <h3 class="text-xl font-black text-slate-800 tracking-tight">Scan WhatsApp</h3>
                <p class="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] mt-2">Hubungkan Perangkat</p>
              </div>

              <!-- QR Code Frame -->
              <div class="relative bg-slate-50 p-6 rounded-[2.5rem] border-2 border-dashed border-slate-100 flex items-center justify-center min-h-[250px]">
                <img v-if="qrCodeUrl" :src="qrCodeUrl" class="w-full h-full object-contain rounded-2xl shadow-lg" alt="WA QR" />
                <div v-else class="flex flex-col items-center gap-4">
                  <div class="w-12 h-12 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin"></div>
                  <div class="space-y-1">
                    <p class="text-[10px] font-black text-slate-700 uppercase tracking-widest">Menyiapkan Sesi...</p>
                    <p class="text-[9px] font-medium text-slate-400">Server sedang merespon, mohon tunggu.</p>
                  </div>
                  <button @click="openQRModal" class="mt-4 text-[9px] font-black text-indigo-600 uppercase tracking-widest hover:underline">Coba Lagi</button>
                </div>
              </div>

              <p class="text-[10px] font-bold text-slate-500 leading-relaxed bg-slate-50 p-4 rounded-2xl">
                Buka WhatsApp di HP Anda, masuk ke Perangkat Tertaut, dan arahkan kamera ke kode di atas.
              </p>

              <button @click="showQRModal = false" class="w-full py-4 bg-slate-900 text-white rounded-2xl font-black text-xs uppercase tracking-widest hover:bg-slate-800 transition-all shadow-xl">
                Selesai
              </button>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<script setup>
import { useAuthStore } from '../store/auth'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from '../composables/useToast'
import ToastContainer from '../components/ToastContainer.vue'
import { 
  Users as UsersIcon, 
  GraduationCap as GraduationCapIcon, 
  CreditCard as CreditCardIcon, 
  Receipt as ReceiptIcon,
  LogOut as LogOutIcon,
  User as UserIcon,
  LayoutGrid as LayoutGridIcon,
  Calendar as CalendarIcon,
  MapPin as MapPinIcon,
  Search as SearchIcon,
  MessageCircle as MessageCircleIcon,
  Activity as ActivityIcon,
  ChevronRight as ChevronRightIcon,
  RotateCcw as RotateCcwIcon,
  X as XIcon,
  FileText as FileTextIcon,
  Database as DatabaseIcon
} from 'lucide-vue-next'
import axios from 'axios'
import { nextTick, onMounted, onUnmounted, ref, computed, watch } from 'vue'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

const waStatus = ref('OFFLINE')
const showQRModal = ref(false)
const qrCodeUrl = ref(null)
const toast = useToast()
const waActionLoading = ref(false)
const panelLabel = computed(() => authStore.user?.role === 'parent' ? 'Parent Portal' : 'Admin Panel')
const headerSubtitle = computed(() => authStore.user?.role === 'parent' ? 'Portal Orang Tua' : 'Industrial Ecosystem')

let waInterval = null
let ws = null
let lastSupportToastAt = 0

const openDropdowns = ref({
  'Data Master': false,
  'Keuangan & Tagihan': false,
  'Manajemen Sistem': false
})

const adminMenuGroups = [
  {
    title: 'Utama',
    isDirectGroup: true,
    children: [
      { name: 'dashboard', label: 'Dashboard', icon: LayoutGridIcon, to: { name: 'dashboard' }, exact: true },
    ]
  },
  {
    title: 'Data Master',
    icon: UsersIcon,
    children: [
      { name: 'users', label: 'Pengguna', icon: UsersIcon, to: '/users' },
      { name: 'students', label: 'Siswa', icon: UserIcon, to: '/students' },
      { name: 'majors', label: 'Jurusan', icon: GraduationCapIcon, to: '/academic/major' },
      { name: 'classes', label: 'Kelas', icon: LayoutGridIcon, to: '/academic/class' },
      { name: 'years', label: 'Angkatan', icon: CalendarIcon, to: '/academic/years' },
    ]
  },
  {
    title: 'Keuangan & Tagihan',
    icon: CreditCardIcon,
    children: [
      { name: 'bill-types', label: 'Jenis Tagihan', icon: CreditCardIcon, to: { name: 'bill-types' } },
      { name: 'billing-rules', label: 'Aturan Tagihan', icon: CalendarIcon, to: { name: 'billing-rules' } },
      { name: 'all-bills', label: 'Data Tagihan', icon: ReceiptIcon, to: { name: 'all-bills' } },
    ]
  },
  {
    title: 'Manajemen Sistem',
    icon: DatabaseIcon,
    children: [
      { name: 'reports', label: 'Laporan', icon: FileTextIcon, to: '/reports' },
      { name: 'notification-logs', label: 'Efikasi Notifikasi', icon: MessageCircleIcon, to: '/notifications' },
      { name: 'support-chat', label: 'CS WhatsApp', icon: MessageCircleIcon, to: '/support/chat' },
      { name: 'audit-logs', label: 'Log Audit', icon: DatabaseIcon, to: '/audit-logs' },
    ]
  }
]

const canAccess = (item) => {
  // Di masa depan, di sinilah logika RBAC / permission dipasang.
  // Contoh: if (authStore.user?.role === 'bendahara' && !item.allowBendahara) return false;
  // Saat ini admin memiliki akses penuh ke semua menu:
  return true
}

const computedMenuGroups = computed(() => {
  return adminMenuGroups.map(group => {
    if (group.isDirectGroup) {
      return {
        ...group,
        children: group.children.filter(child => canAccess(child))
      }
    }

    const filteredChildren = group.children.filter(child => canAccess(child))

    // LOGIKA AUTO-FLATTEN (Jika sub-menu hanya tersisa 1, jadikan direct group tanpa dropdown)
    if (filteredChildren.length === 1) {
      return {
        title: group.title,
        isDirectGroup: true,
        children: filteredChildren
      }
    }

    return {
      ...group,
      children: filteredChildren
    }
  }).filter(group => group.children.length > 0)
})

// Auto open group if active route is inside
const checkActiveGroup = () => {
  adminMenuGroups.forEach(group => {
    if (!group.isDirectGroup && group.children) {
      const hasActive = group.children.some(child => {
        if (typeof child.to === 'string') return route.path.startsWith(child.to)
        if (child.to.name) return route.name === child.to.name
        return false
      })
      if (hasActive) {
        openDropdowns.value[group.title] = true
      }
    }
  })
}

const fetchWAStatus = async () => {
  if (authStore.user?.role !== 'admin') return
  try {
    const res = await axios.get('whatsapp/status')
    waStatus.value = res.data.data.status
    // Auto close modal if connected
    if (['WORKING', 'CONNECTED'].includes(waStatus.value) && showQRModal.value) {
      showQRModal.value = false
    }
  } catch (err) {
    waStatus.value = 'OFFLINE'
  }
}

const openQRModal = async () => {
  showQRModal.value = true
  qrCodeUrl.value = null
  try {
    // Fetch QR as Blob to include Authorization headers
    const res = await axios.get('whatsapp/qr', { responseType: 'blob' })
    if (qrCodeUrl.value) URL.revokeObjectURL(qrCodeUrl.value)
    qrCodeUrl.value = URL.createObjectURL(res.data)
  } catch (err) {
    console.error('Failed to load QR', err)
    toast.error('QR belum siap', err.response?.data?.message || 'Coba lagi beberapa detik lagi.')
  }
}

const disconnectWA = async () => {
  if (waActionLoading.value) return
  waActionLoading.value = true
  try {
    await axios.post('whatsapp/logout')
    waStatus.value = 'STOPPED'
    toast.success('WhatsApp logout', 'Nomor lama sudah dilepas dari SchoolPay.')
    await fetchWAStatus()
  } catch (err) {
    toast.error('Gagal logout WhatsApp', err.response?.data?.message || 'WAHA/server tidak merespon')
  } finally {
    waActionLoading.value = false
  }
}

const changeWANumber = async () => {
  if (waActionLoading.value) return
  waActionLoading.value = true
  try {
    await axios.post('whatsapp/logout')
    waStatus.value = 'STARTING'
    toast.success('Siap ganti nomor', 'Scan QR dengan nomor WhatsApp sekolah yang baru.')
    setTimeout(() => openQRModal(), 1200)
  } catch (err) {
    toast.error('Gagal menyiapkan ganti nomor', err.response?.data?.message || 'WAHA/server tidak merespon')
  } finally {
    waActionLoading.value = false
  }
}

const restartWA = async () => {
  if (waActionLoading.value) return
  waActionLoading.value = true
  try {
    await axios.post('whatsapp/restart')
    waStatus.value = 'STARTING'
    toast.success('Session direstart', 'SchoolPay sedang menyambungkan ulang WhatsApp.')
    setTimeout(fetchWAStatus, 1500)
  } catch (err) {
    toast.error('Gagal restart WhatsApp', err.response?.data?.message || 'WAHA/server tidak merespon')
  } finally {
    waActionLoading.value = false
  }
}

const handleRateLimitError = (event) => {
  toast.warning('Tunggu sebentar', event.detail?.message || 'Terlalu banyak request. Coba lagi sebentar lagi.')
}

const handleNetworkError = (event) => {
  toast.error('Server tidak terhubung', event.detail || 'Server tidak dapat dijangkau.')
}

const initWebSocket = () => {
  if (!['admin', 'parent'].includes(authStore.user?.role)) return

  const token = authStore.token
  if (!token) return

  if (ws) {
    try {
      ws.close()
    } catch (e) {}
  }

  let wsBase = axios.defaults.baseURL || ''
  if (!wsBase.startsWith('http')) {
    // Jika relative path, bangun menggunakan window.location
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    wsBase = `${protocol}//${window.location.host}${wsBase}`
  } else {
    wsBase = wsBase.replace(/^http/, 'ws')
  }

  const wsUrl = wsBase.endsWith('/') 
    ? `${wsBase}ws?token=${token}` 
    : `${wsBase}/ws?token=${token}`

  ws = new WebSocket(wsUrl)

  ws.onmessage = (event) => {
    try {
      const msg = JSON.parse(event.data)
      if (msg.topic === 'NEW_PAYMENT' && authStore.user?.role === 'admin') {
        const data = msg.data
        toast.success(
          'Pembayaran Baru!',
          `Siswa: ${data.student_name} - Rp ${data.amount.toLocaleString('id-ID')}`
        )
        // Dispatch global event for views to refresh
        window.dispatchEvent(new CustomEvent('new-payment', { detail: data }))
        
        // Play notification sound
        const audio = new Audio('https://assets.mixkit.co/active_storage/sfx/2358/2358-preview.mp3')
        audio.play().catch(() => {}) // Ignore if browser blocks auto-play
      } else if (msg.topic === 'WA_STATUS_CHANGED') {
        waStatus.value = msg.data.status
        // Auto close modal if connected
        if (['WORKING', 'CONNECTED'].includes(waStatus.value) && showQRModal.value) {
          showQRModal.value = false
        }
      } else if (msg.topic === 'NOTIFICATION_STATUS_CHANGED') {
        window.dispatchEvent(new CustomEvent('notification-status-changed', { detail: msg.data }))
      } else if (msg.topic === 'SUPPORT_CHAT_UPDATED') {
        if (authStore.user?.role === 'admin') {
          const now = Date.now()
          const isSupportPage = route.path === '/support/chat'
          if (!isSupportPage && now - lastSupportToastAt > 8000) {
            toast.info('Chat CS baru', msg.data.phone ? `Pesan dari ${msg.data.phone}` : 'Ada pembaruan chat CS')
            lastSupportToastAt = now
          }
        }
        window.dispatchEvent(new CustomEvent('support-chat-updated', { detail: msg.data }))
      }
    } catch (err) {
      console.error('WS Message Error:', err)
    }
  }

  ws.onclose = () => {
    console.warn('WS disconnected, retrying in 5s...')
    setTimeout(initWebSocket, 5000)
  }

  ws.onerror = (err) => {
    console.error('WS Error:', err)
  }
}

const openParentSupport = async () => {
  if (!route.path.startsWith('/parent')) {
    await router.push('/parent/dashboard')
  }
  await nextTick()
  window.dispatchEvent(new CustomEvent('open-parent-support'))
}

const handleLogout = () => {
  authStore.logout()
  router.push('/')
}

onMounted(() => {
  checkActiveGroup()
  window.addEventListener('rate-limit-error', handleRateLimitError)
  window.addEventListener('network-error', handleNetworkError)
  fetchWAStatus()
  initWebSocket()
  if (authStore.user?.role === 'admin') {
    waInterval = setInterval(fetchWAStatus, 30000)
  }
})

watch(() => route.path, () => {
  checkActiveGroup()
})

onUnmounted(() => {
  window.removeEventListener('rate-limit-error', handleRateLimitError)
  window.removeEventListener('network-error', handleNetworkError)
  if (waInterval) clearInterval(waInterval)
  if (qrCodeUrl.value) URL.revokeObjectURL(qrCodeUrl.value)
  if (ws) ws.close()
})
</script>

<style scoped>
.page-enter-active,
.page-leave-active {
  transition: opacity 0.2s ease;
}

.page-enter-from,
.page-leave-to {
  opacity: 0;
}

.slide-down-enter-active, .slide-down-leave-active { transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1); }
.slide-down-enter-from, .slide-down-leave-to { transform: translate(-50%, -100px); opacity: 0; }

.dropdown-enter-active, .dropdown-leave-active { 
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); 
  max-height: 260px; 
  opacity: 1;
  transform: translateY(0);
  overflow: hidden; 
}
.dropdown-enter-from, .dropdown-leave-to { 
  opacity: 0; 
  max-height: 0; 
  transform: translateY(-5px); 
  overflow: hidden;
}
</style>
