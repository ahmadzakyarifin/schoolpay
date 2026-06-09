<script setup>
import { ref, watch, onMounted, computed, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
const router = useRouter()
import axios from 'axios'
import { 
  Search as SearchIcon, 
  Plus as PlusIcon, 
  RotateCcw as ResetIcon,
  ChevronLeft as PrevIcon,
  ChevronRight as NextIcon,
  CheckCircle2 as SuccessIcon,
  AlertCircle as AlertIcon,
  FileSpreadsheet as ExcelIcon,
  Filter as FilterIcon,
  Trash as TrashIcon,
  Database as DatabaseIcon,
  History as HistoryIcon
} from 'lucide-vue-next'

// Components
import userService from '../../services/user.service'
import { useForm } from '../../composables/useForm'
import UserTable from '../../components/users/UserTable.vue'
import UserFormModal from '../../components/users/UserFormModal.vue'
import UserFilter from '../../components/users/UserFilter.vue'
import UserBulkActions from '../../components/users/UserBulkActions.vue'
import { useAuthStore } from '../../store/auth'

const authStore = useAuthStore()
const isOffline = computed(() => (typeof navigator !== 'undefined' && navigator.onLine === false))

// State
const users = ref([])
const total = ref(0)
const page = ref(1)
const limit = ref(10)
const search = ref('')
const role = ref('')
const filter = ref('')
const status = ref('')
const sort = ref('')
const loading = ref(false)
const isMounted = ref(false)
const showFilters = ref(false)
const showEditModal = ref(false)
const isEditing = ref(false)
const selectedUserIds = ref([])
const birthDateDisplay = ref('')
const showHistory = ref(false)

// Delete Confirmation State
const showDeleteConfirm = ref(false)
const userToDelete = ref(null)
const isBulkDelete = ref(false)
const deleteLoading = ref(false)
const dependencyLoading = ref(false)
const dependencyInfo = ref(null)

// Status Confirmation State
const showStatusConfirm = ref(false)
const userToToggle = ref(null)
const statusLoading = ref(false)

// Restore Confirmation State
const showRestoreConfirm = ref(false)
const userToRestore = ref(null)
const isBulkRestore = ref(false)
const restoreLoading = ref(false)

const tempFilters = reactive({
  role: '',
  status: '',
  filter: '',
  sort: ''
})

const notification = reactive({
  show: false,
  message: '',
  type: 'success'
})

const bulkDeleteBlockedUsers = computed(() => {
  if (!isBulkDelete.value || selectedUserIds.value.length === 0) return []

  return users.value.filter(user => {
    const hasActiveStudents = Number(user.student_count || 0) > 0
    return selectedUserIds.value.includes(user.id) && user.role === 'parent' && hasActiveStudents
  })
})

// Form Logic
const initialForm = {
  id: null,
  name: '',
  email: '',
  phone_number: '',
  country_code: '62',
  role: 'parent',
  nik: null,
  birth_date: null,
  address: null,
  education: null,
  occupation: null,
  income: null
}

const { form, errors, submitting, setErrors, clearErrors, clearFieldError, resetForm } = useForm(initialForm)

// Fetch Logic
const fetchUsers = async () => {
  loading.value = true
  try {
    const response = await userService.getAll({
      page: page.value,
      limit: limit.value,
      search: search.value || undefined,
      role: role.value || undefined,
      filter: filter.value || undefined,
      status: showHistory.value ? 'trash' : (status.value || undefined),
      sort: sort.value || undefined
    })
    const responseData = response.data?.data
    users.value = responseData?.users || responseData?.data || []
    total.value = responseData?.total ?? users.value.length
    cacheUserData()
  } catch (err) {
    console.error('Failed to fetch users:', err)
    if (!loadCachedUserData()) {
      users.value = []
      total.value = 0
    }
  } finally {
    loading.value = false
  }
}

// Handlers
const applyFilters = () => {
  role.value = tempFilters.role
  status.value = tempFilters.status
  filter.value = tempFilters.filter
  sort.value = tempFilters.sort
  page.value = 1
  showFilters.value = false
  fetchUsers()
}

const resetFilters = () => {
  Object.assign(tempFilters, { role: '', status: '', filter: '', sort: '' })
  search.value = ''
  role.value = ''
  filter.value = ''
  status.value = ''
  sort.value = ''
  page.value = 1
  showFilters.value = false
  fetchUsers()
}

// State Preservation Logic
const STATE_KEY = 'user_management_state'
const DATA_CACHE_KEY = 'user_management_cached_data'

const cacheUserData = () => {
  localStorage.setItem(DATA_CACHE_KEY, JSON.stringify({
    users: users.value,
    total: total.value,
    cached_at: new Date().toISOString()
  }))
}

const loadCachedUserData = () => {
  const cached = localStorage.getItem(DATA_CACHE_KEY)
  if (!cached) return false

  try {
    const parsed = JSON.parse(cached)
    users.value = Array.isArray(parsed.users) ? parsed.users : []
    total.value = Number(parsed.total || users.value.length)
    return true
  } catch (e) {
    localStorage.removeItem(DATA_CACHE_KEY)
    return false
  }
}

const saveState = () => {
  const state = {
    search: search.value,
    tempFilters: { ...tempFilters },
    page: page.value,
    limit: limit.value,
    showHistory: showHistory.value
  }
  localStorage.setItem(STATE_KEY, JSON.stringify(state))
}

const loadState = () => {
  const saved = localStorage.getItem(STATE_KEY)
  if (saved) {
    try {
      const state = JSON.parse(saved)
      search.value = state.search || ''
      Object.assign(tempFilters, state.tempFilters || {})
      // Apply filters to active state
      role.value = tempFilters.role
      status.value = tempFilters.status
      filter.value = tempFilters.filter
      sort.value = tempFilters.sort
      
      page.value = state.page || 1
      limit.value = state.limit || 10
      showHistory.value = !!state.showHistory
    } catch (e) {
      console.error('Gagal memuat state:', e)
    }
  }
}

// Watchers for state preservation
watch([search, page, limit, showHistory, role, status, filter, sort], saveState)

onMounted(() => {
  loadState()
  isMounted.value = true
  fetchUsers()
})

const openAddModal = () => {
  isEditing.value = false
  birthDateDisplay.value = ''
  resetForm()
  showEditModal.value = true
}

const openEditModal = (user) => {
  isEditing.value = true
  clearErrors()
  
  // Extract country code
  let phone = user.phone_number || ''
  let code = '62'
  const countries = [{code:'62'},{code:'60'},{code:'65'},{code:'1'}]
  for (const c of countries) {
    if (phone.startsWith(c.code)) {
      code = c.code
      phone = phone.substring(c.code.length)
      break
    }
  }

  Object.assign(form, { ...user, country_code: code, phone_number: phone })
  if (user.birth_date) {
    const d = new Date(user.birth_date)
    birthDateDisplay.value = `${d.getDate().toString().padStart(2, '0')}/${(d.getMonth()+1).toString().padStart(2, '0')}/${d.getFullYear()}`
  } else {
    birthDateDisplay.value = ''
  }
  showEditModal.value = true
}

const handleGoToStudent = (studentId) => {
  router.push(`/students/${studentId}`)
}

const handleGoToUserDetails = (user) => {
  router.push(`/users/${user.id}`)
}

const saveUser = async () => {
  clearErrors()
  const validationErrors = {}

  // Validasi Nama
  if (!form.name) {
    validationErrors.name = ['Nama wajib diisi.']
  } else if (form.name.length < 2) {
    validationErrors.name = ['Nama minimal 2 karakter.']
  }

  // Validasi Email
  if (!form.email) {
    validationErrors.email = ['Email wajib diisi.']
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    validationErrors.email = ['Format email tidak valid (contoh: user@gmail.com).']
  }

  // Validasi Nomor HP
  if (!form.phone_number) {
    validationErrors.phone_number = ['Nomor HP wajib diisi.']
  } else if (form.phone_number.length < 9) {
    validationErrors.phone_number = ['Nomor HP minimal 9 karakter.']
  }

  // Validasi Role
  if (!form.role) {
    validationErrors.role = ['Role wajib dipilih']
  }

  if (Object.keys(validationErrors).length > 0) {
    setErrors({ response: { data: { errors: validationErrors } } })
    return
  }

	try {
    const payload = { ...form, phone_number: form.country_code + form.phone_number }
    if (isEditing.value) {
      await userService.update(form.id, payload)
    } else {
      await userService.create(payload)
    }
    showNotification('Data berhasil disimpan', 'success')
    showEditModal.value = false
    fetchUsers()
  } catch (err) {
    setErrors(err)
  }
}

const handleLocalValidation = (localErrors) => {
	errors.value = { ...errors.value, ...localErrors }
}

const setFieldError = ({ field, messages }) => {
  errors.value = { ...errors.value, [field]: messages }
}

const toggleStatus = (user) => {
  userToToggle.value = user
  showStatusConfirm.value = true
}

const handleToggleStatus = async () => {
  if (!userToToggle.value) return
  statusLoading.value = true
  try {
    await userService.toggleStatus(userToToggle.value.id)
    userToToggle.value.is_active = !userToToggle.value.is_active
    showNotification(`Status ${userToToggle.value.name} berhasil diubah`, 'success')
    showStatusConfirm.value = false
  } catch (err) {
    const errorMsg = err.response?.data?.message || 'Gagal mengubah status'
    showNotification(errorMsg, 'error')
  } finally {
    statusLoading.value = false
  }
}

// Delete Logic
const confirmDelete = async (user) => {
  userToDelete.value = user
  isBulkDelete.value = false
  showDeleteConfirm.value = true
  dependencyLoading.value = true
  dependencyInfo.value = null
  try {
    const res = await axios.get(`users/${user.id}/dependency-info`)
    dependencyInfo.value = res.data?.data
  } catch (err) {
    dependencyInfo.value = null
  } finally {
    dependencyLoading.value = false
  }
}

const confirmBulkDelete = () => {
  if (selectedUserIds.value.length === 0) return
  isBulkDelete.value = true
  showDeleteConfirm.value = true
}

const handleDelete = async () => {
  deleteLoading.value = true
  try {
    if (isBulkDelete.value) {
      await userService.bulkDelete(selectedUserIds.value)
      showNotification(`${selectedUserIds.value.length} data berhasil dihapus`, 'success')
      selectedUserIds.value = []
    } else {
      await userService.delete(userToDelete.value.id)
      showNotification(`${userToDelete.value.name} berhasil dihapus`, 'success')
    }
    showDeleteConfirm.value = false
    fetchUsers()
  } catch (err) {
    const errorMsg = err.response?.data?.message || 'Gagal menghapus data'
    showNotification(errorMsg, 'error')
  } finally {
    deleteLoading.value = false
  }
}

const confirmRestore = (user) => {
  userToRestore.value = user
  isBulkRestore.value = false
  showRestoreConfirm.value = true
}

const confirmBulkRestore = () => {
  if (selectedUserIds.value.length === 0) return
  isBulkRestore.value = true
  showRestoreConfirm.value = true
}

const handleRestore = async () => {
  restoreLoading.value = true
  try {
    if (isBulkRestore.value) {
      await userService.bulkRestore(selectedUserIds.value)
      showNotification(`${selectedUserIds.value.length} data berhasil dipulihkan`, "success")
      selectedUserIds.value = []
    } else {
      await userService.restore(userToRestore.value.id)
      showNotification(`${userToRestore.value.name} berhasil dipulihkan`, "success")
    }
    showRestoreConfirm.value = false
    fetchUsers()
  } catch (err) {
    const errorMsg = err.response?.data?.message || "Gagal memulihkan data"
    showNotification(errorMsg, "error")
  } finally {
    restoreLoading.value = false
  }
}

const handleBulkResend = async (channel) => {
  if (selectedUserIds.value.length === 0) return
  
  try {
    const response = await userService.bulkResendNotification(selectedUserIds.value, channel)
    const result = response.data.data // Access the nested data field from utils.SuccessResponse

    if (result && result.failed > 0) {
      // Kelompokkan error berdasarkan alasannya
      const groups = result.errors.reduce((acc, err) => {
        const reason = err.split(': ')[1] || 'Lainnya'
        acc[reason] = (acc[reason] || 0) + 1
        return acc
      }, {})

      const errorSummary = Object.entries(groups)
        .map(([reason, count]) => `${count} (${reason})`)
        .join(', ')

      showNotification(
        `✅ ${result.sent} Berhasil Terkirim\n⚠️ Gagal: ${errorSummary}`,
        'warning'
      )
      console.warn('Bulk Notification Errors:', result.errors)
    } else {
      showNotification(`${result.sent || 0} Link aktivasi berhasil dikirim`, 'success')
    }

    selectedUserIds.value = []
    fetchUsers()
  } catch (err) {
    const errorMsg = err.response?.data?.message || 'Gagal mengirim aktivasi massal'
    showNotification(errorMsg, 'error')
  }
}

const handleResendUser = async ({ user, channel }) => {
  if (!user?.id) return

  try {
    await userService.resendNotification(user.id, channel)
    showNotification(`Link aktivasi ${user.name} berhasil dikirim`, 'success')
    fetchUsers()
  } catch (err) {
    const errorMsg = err.response?.data?.message || 'Gagal mengirim link aktivasi'
    showNotification(errorMsg, 'error')
  }
}

const handleExport = async () => {
  if (isOffline.value) {
    showNotification('Export Excel tidak tersedia saat offline. Server harus online agar file berisi data terbaru.', 'error')
    return
  }
  try {
    const response = await axios.get('users/export', {
      params: { search: search.value, role: role.value, status: status.value },
      responseType: 'blob'
    })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    const dateStr = new Date().toISOString().slice(0, 10)
    link.setAttribute('download', `Data_Pengguna_${dateStr}.xlsx`)
    document.body.appendChild(link)
    link.click()
    link.remove()
    setTimeout(() => window.URL.revokeObjectURL(url), 1000)
  } catch (err) {
    showNotification('Gagal ekspor data', 'error')
  }
}


const showNotification = (msg, type = 'success') => {
  notification.message = msg
  notification.type = type
  notification.show = true
  setTimeout(() => notification.show = false, 4000)
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear()}`
}

// Bulk Logic
const toggleSelectAll = () => {
	const selectableIds = users.value.filter(u => u.id !== authStore.user?.id).map(u => u.id)
	if (selectedUserIds.value.length === selectableIds.length && selectableIds.length > 0) {
		selectedUserIds.value = []
	} else {
		selectedUserIds.value = selectableIds
	}
}

const toggleSelectUser = (id) => {
	const index = selectedUserIds.value.indexOf(id)
	if (index > -1) {
		selectedUserIds.value.splice(index, 1)
	} else {
		selectedUserIds.value.push(id)
	}
}

// Lifecycle & Watchers
watch(showHistory, (newVal) => {
	if (newVal) {
		showFilters.value = false
		role.value = ''
		status.value = ''
		filter.value = ''
		sort.value = ''
		Object.assign(tempFilters, { role: '', status: '', filter: '', sort: '' })
	}
	page.value = 1
	selectedUserIds.value = []
	fetchUsers()
})

onMounted(() => {
  isMounted.value = true
  fetchUsers()
})

watch(search, () => {
  if (page.value === 1) {
    fetchUsers()
  } else {
    page.value = 1
  }
})

watch(page, () => {
	fetchUsers()
})

watch(limit, () => {
  if (page.value === 1) {
    fetchUsers()
  } else {
    page.value = 1
  }
})

const totalPages = computed(() => Math.ceil(total.value / limit.value) || 1)

const selectedUsers = computed(() => {
  return users.value.filter(u => selectedUserIds.value.includes(u.id))
})

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

</script>

<template>
  <div class="min-h-screen bg-slate-50/50 p-4 lg:p-6 space-y-6 animate-fade-in font-inter">
    <!-- Notifications -->
    <Teleport v-if="isMounted" to="body">
      <transition name="fade">
        <div v-if="notification.show" class="fixed top-6 right-6 z-[600] px-6 py-4 rounded-2xl shadow-2xl flex items-center gap-3"
          :class="notification.type === 'success' ? 'bg-emerald-600 text-white' : 'bg-rose-600 text-white'">
          <SuccessIcon v-if="notification.type === 'success'" class="w-5 h-5 shrink-0" />
          <AlertIcon v-else class="w-5 h-5 shrink-0" />
          <span class="text-xs font-bold">{{ notification.message }}</span>
        </div>
      </transition>
    </Teleport>

    <Teleport v-if="isMounted" to="#header-actions-target">
      <div class="flex items-center justify-center w-full gap-4 relative mx-auto">
        <div class="flex items-center justify-center gap-2 flex-1 max-w-2xl mx-auto">
          <div class="relative flex-1 group">
            <SearchIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300 group-focus-within:text-indigo-600" />
            <input v-model="search" type="text" placeholder="Cari nama, email, whatsapp..." class="search-input-premium" />
          </div>
          <button v-if="!showHistory" @click="showFilters = !showFilters" class="relative p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 flex items-center gap-2 transition-all shadow-sm cursor-pointer">
            <FilterIcon class="w-4 h-4" />
            <span class="text-[10px] font-black uppercase tracking-wider pr-1">Filter</span>
            <span v-if="role || status || filter || sort" class="absolute -top-1 -right-1 w-3 h-3 bg-indigo-600 rounded-full border-2 border-white shadow-sm"></span>
          </button>
          <button @click="resetFilters" class="p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm transition-all group" title="Reset">
            <ResetIcon class="w-4 h-4 group-hover:rotate-180 transition-transform duration-500" />
          </button>
        </div>

        <UserFilter 
          v-model="showFilters" 
          :filters="tempFilters" 
          @apply="applyFilters" 
          @reset="resetFilters" 
        />
      </div>
    </Teleport>

    <!-- Main Content Table -->
    <div class="bg-white border border-slate-200 rounded-xl shadow-sm flex flex-col min-h-[710px] overflow-hidden">
      <!-- Table Header -->
      <div class="p-4 border-b border-slate-100 bg-slate-50/30 flex items-center justify-between">
        <div class="flex items-center gap-2 font-black text-slate-700 text-xs uppercase tracking-widest">
          <ParentIcon class="w-3.5 h-3.5 text-indigo-600 animate-pulse" />
          <span>{{ showHistory ? 'Riwayat Data Terhapus' : 'Data Operasional Pengguna' }}</span>
        </div>

        <div class="flex items-center gap-3">
          <UserBulkActions 
            :selectedCount="selectedUserIds.length" 
            :selectedUsers="selectedUsers"
            :status="showHistory ? 'trash' : status" 
            @resend="(channel) => handleBulkResend(channel)"
            @delete="confirmBulkDelete"
            @restore="confirmBulkRestore"
          />
          <button @click="handleExport" :disabled="isOffline" :title="isOffline ? 'Export Excel membutuhkan server online agar data terbaru.' : 'Ekspor Excel'" :class="['font-bold py-1.5 px-3 rounded-xl border text-[10px] flex items-center gap-1.5 transition-all shadow-sm', isOffline ? 'bg-amber-50 border-amber-200 text-amber-700 cursor-not-allowed' : 'bg-white text-slate-600 hover:bg-slate-50']">
            <ExcelIcon class="w-3.5 h-3.5" :class="isOffline ? 'text-amber-600' : 'text-emerald-600'" />
            <span>{{ isOffline ? 'Excel Online Saja' : 'Ekspor Excel' }}</span>
          </button>

          <button @click="showHistory = !showHistory" class="bg-white text-slate-600 border border-slate-200 hover:bg-slate-50 font-bold py-1.5 px-3 rounded-xl text-[10px] flex items-center gap-1.5 transition-all shadow-sm">
            <HistoryIcon v-if="!showHistory" class="w-3.5 h-3.5" />
            <ResetIcon v-else class="w-3.5 h-3.5 rotate-180" />
            <span>{{ showHistory ? 'Kembali' : 'Riwayat Hapus' }}</span>
          </button>
          <button v-if="!showHistory" @click="openAddModal" class="bg-indigo-600 hover:bg-indigo-700 text-white font-black py-1.5 px-4 rounded-xl flex items-center gap-1.5 shadow-md shadow-indigo-100 transition-all text-[10px] uppercase tracking-widest shrink-0">
            <PlusIcon class="w-3.5 h-3.5" />
            <span>Tambah Data</span>
          </button>
        </div>
      </div>

      <UserTable 
        :users="users" 
        :loading="loading" 
        :selectedUserIds="selectedUserIds"
        :status="showHistory ? 'trash' : status"
        :formatDate="formatDate"
        @edit="openEditModal"
        @delete="confirmDelete"
        @restore="confirmRestore"
        @toggle-status="toggleStatus"
        @toggle-select-all="toggleSelectAll"
        @toggle-select-user="toggleSelectUser"
        @go-to-student="handleGoToStudent"
        @go-to-details="handleGoToUserDetails"
        @resend-notification="handleResendUser"
      />

      <!-- Pagination -->
      <div class="px-6 py-4 bg-slate-50/50 border-t border-slate-100 flex items-center justify-between">
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
            Halaman <span class="text-indigo-600">{{ page }}</span> dari {{ totalPages }} <span class="mx-2 text-slate-300">|</span> Total <span class="text-indigo-600">{{ total }}</span> Data
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

          <!-- Page Numbers -->
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

    <!-- Modals -->
    <UserFormModal 
			v-model="showEditModal"
			:isEditing="isEditing"
			:form="form"
			:errors="errors"
			:submitting="submitting"
			:birth-date-display="birthDateDisplay"
			@update:birth-date-display="birthDateDisplay = $event"
			@save="saveUser"
			@local-validation-failed="handleLocalValidation"
			@clear-field-error="clearFieldError"
			@set-field-error="setFieldError"
		/>


    <!-- Delete Confirmation Modal -->
    <Teleport v-if="isMounted" to="body">
      <transition name="fade">
        <div v-if="showDeleteConfirm" class="fixed inset-0 z-[700] flex items-center justify-center p-6">
          <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showDeleteConfirm = false"></div>
          <div class="white-card w-full max-w-md relative z-10 overflow-hidden shadow-[0_20px_50px_rgba(0,0,0,0.3)] animate-scale-in !rounded-[2.5rem] p-8 text-center">
            <div class="w-20 h-20 bg-rose-50 text-rose-500 rounded-[2rem] flex items-center justify-center mx-auto mb-6 border border-rose-100 shadow-xl shadow-rose-500/10">
              <TrashIcon class="w-10 h-10" />
            </div>
            
            <h3 class="text-xl font-black text-slate-900 tracking-tight mb-2">
              {{ isBulkDelete ? 'Hapus Terpilih?' : 'Hapus Pengguna?' }}
            </h3>

            <div v-if="dependencyLoading" class="my-4 py-3 px-4 bg-slate-50 rounded-2xl flex items-center justify-center gap-2 text-slate-500 text-[10px] font-bold uppercase tracking-widest">
              <div class="w-3 h-3 border-2 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
              Memeriksa keterhubungan data...
            </div>

            <div v-else-if="(!isBulkDelete && dependencyInfo?.has_dependencies) || (isBulkDelete && bulkDeleteBlockedUsers.length > 0)" class="my-4 p-4 bg-amber-50 border border-amber-200/80 rounded-2xl text-left shadow-sm">
              <div class="flex items-start gap-3">
                <AlertIcon class="w-5 h-5 text-amber-600 shrink-0 mt-0.5" />
                <div>
                  <h4 class="text-xs font-black text-amber-900 uppercase tracking-wider mb-1">Tidak Dapat Dihapus</h4>
                  <p class="text-amber-800 text-[11px] font-medium leading-relaxed mb-2">
                    <template v-if="isBulkDelete">
                      {{ bulkDeleteBlockedUsers.length }} wali murid terpilih masih terhubung dengan siswa aktif.
                    </template>
                    <template v-else>
                      Pengguna ini masih memiliki keterhubungan aktif: <span class="font-bold underline">{{ dependencyInfo.message }}</span>.
                    </template>
                  </p>
                  <p class="text-amber-700/90 text-[10px] font-bold uppercase tracking-wider bg-amber-100/50 py-1 px-2.5 rounded-lg inline-block">
                    Lepaskan atau pindahkan relasi siswa terlebih dahulu sebelum menghapus wali murid ini.
                  </p>
                </div>
              </div>
            </div>
            
            <p class="text-slate-500 text-[10px] font-bold uppercase tracking-widest mb-8 px-4 leading-relaxed">
              {{ isBulkDelete 
                ? `Apakah Anda yakin ingin menghapus ${selectedUserIds.length} data terpilih? Data akan dipindahkan ke riwayat penghapusan (Trash).` 
                : `Apakah Anda yakin ingin menghapus ${userToDelete?.name}? Data akan dipindahkan ke riwayat penghapusan (Trash). Data historis tetap aman di database.` 
              }}
            </p>

            <div class="grid grid-cols-2 gap-4">
              <button 
                @click="showDeleteConfirm = false" 
                class="py-4 bg-slate-100 text-slate-600 font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-slate-200 transition-all"
              >
                Batalkan
              </button>
              <button 
                @click="handleDelete" 
                :disabled="deleteLoading || (!isBulkDelete && dependencyInfo?.has_dependencies) || (isBulkDelete && bulkDeleteBlockedUsers.length > 0)"
                class="py-4 bg-rose-600 text-white font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-rose-700 transition-all shadow-lg shadow-rose-600/20 disabled:opacity-50"
              >
                {{ deleteLoading ? 'Menghapus...' : (((!isBulkDelete && dependencyInfo?.has_dependencies) || (isBulkDelete && bulkDeleteBlockedUsers.length > 0)) ? 'Tidak Bisa Dihapus' : 'Ya, Hapus Data') }}
              </button>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>
    <!-- Status Toggle Confirmation Modal -->
    <Teleport v-if="isMounted" to="body">
      <transition name="fade">
        <div v-if="showStatusConfirm" class="fixed inset-0 z-[700] flex items-center justify-center p-6">
          <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showStatusConfirm = false"></div>
          <div class="white-card w-full max-w-md relative z-10 overflow-hidden shadow-[0_20px_50px_rgba(0,0,0,0.3)] animate-scale-in !rounded-[2.5rem] p-8 text-center">
            <div class="w-20 h-20 rounded-[2rem] flex items-center justify-center mx-auto mb-6 border shadow-xl transition-all duration-500"
              :class="userToToggle?.is_active ? 'bg-rose-50 text-rose-500 border-rose-100 shadow-rose-500/10' : 'bg-emerald-50 text-emerald-600 border-emerald-100 shadow-emerald-500/10'">
              <AlertIcon v-if="userToToggle?.is_active" class="w-10 h-10" />
              <SuccessIcon v-else class="w-10 h-10" />
            </div>
            
            <h3 class="text-xl font-black text-slate-900 tracking-tight mb-2">
              {{ userToToggle?.is_active ? 'Non-aktifkan Akun?' : 'Aktifkan Akun?' }}
            </h3>
            
            <p class="text-slate-500 text-xs font-medium leading-relaxed mb-8 px-4">
              {{ userToToggle?.is_active 
                ? `Apakah Anda yakin ingin menonaktifkan akun ${userToToggle?.name}? Pengguna tidak akan bisa masuk ke aplikasi sampai diaktifkan kembali.` 
                : `Apakah Anda yakin ingin mengaktifkan kembali akun ${userToToggle?.name}? Pengguna akan mendapatkan akses penuh kembali.` 
              }}
            </p>
    
            <div class="grid grid-cols-2 gap-4">
              <button 
                @click="showStatusConfirm = false" 
                class="py-4 bg-slate-100 text-slate-600 font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-slate-200 transition-all"
              >
                Batalkan
              </button>
              <button 
                @click="handleToggleStatus" 
                :disabled="statusLoading"
                class="py-4 text-white font-black rounded-2xl text-[10px] uppercase tracking-widest transition-all shadow-lg disabled:opacity-50"
                :class="userToToggle?.is_active ? 'bg-rose-600 hover:bg-rose-700 shadow-rose-600/20' : 'bg-emerald-600 hover:bg-emerald-700 shadow-emerald-600/20'"
              >
                {{ statusLoading ? 'Memproses...' : (userToToggle?.is_active ? 'Ya, Non-aktifkan' : 'Ya, Aktifkan') }}
              </button>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>

    <!-- Restore Confirmation Modal -->
    <Teleport v-if="isMounted" to="body">
      <transition name="fade">
        <div v-if="showRestoreConfirm" class="fixed inset-0 z-[700] flex items-center justify-center p-6">
          <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showRestoreConfirm = false"></div>
          <div class="white-card w-full max-w-md relative z-10 overflow-hidden shadow-[0_20px_50px_rgba(0,0,0,0.3)] animate-scale-in !rounded-[2.5rem] p-8 text-center">
            <div class="w-20 h-20 rounded-[2rem] flex items-center justify-center mx-auto mb-6 border border-emerald-100 bg-emerald-50 text-emerald-600 shadow-xl shadow-emerald-500/10 transition-all duration-500">
              <SuccessIcon class="w-10 h-10" />
            </div>
            
            <h3 class="text-xl font-black text-slate-900 tracking-tight mb-2">
              {{ isBulkRestore ? 'Pulihkan Data Terpilih?' : 'Pulihkan Akun?' }}
            </h3>
            
            <p class="text-slate-500 text-xs font-medium leading-relaxed mb-8 px-4">
              {{ isBulkRestore 
                ? `Apakah Anda yakin ingin memulihkan ${selectedUserIds.length} data yang terpilih? Data akan dikembalikan ke daftar pengguna aktif.` 
                : `Apakah Anda yakin ingin memulihkan akun ${userToRestore?.name}? Akun ini akan kembali aktif dan bisa digunakan.` 
              }}
            </p>
    
            <div class="grid grid-cols-2 gap-4">
              <button 
                @click="showRestoreConfirm = false" 
                class="py-4 bg-slate-100 text-slate-600 font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-slate-200 transition-all"
              >
                Batalkan
              </button>
              <button 
                @click="handleRestore" 
                :disabled="restoreLoading"
                class="py-4 bg-emerald-600 text-white font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-emerald-700 transition-all shadow-lg shadow-emerald-600/20 disabled:opacity-50"
              >
                {{ restoreLoading ? 'Memulihkan...' : 'Ya, Pulihkan' }}
              </button>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<style scoped lang="postcss">
.search-input-premium {
  @apply w-full bg-white border border-slate-200 rounded-xl py-2.5 pl-12 pr-4 text-xs font-bold text-slate-700 outline-none transition-all focus:border-indigo-500 focus:ring-4 focus:ring-indigo-50 shadow-sm;
}
.btn-primary { @apply bg-indigo-600 hover:bg-indigo-700 text-white transition-all disabled:opacity-50; }
.btn-secondary { @apply bg-white hover:bg-slate-50 text-slate-600 border border-slate-200 transition-all; }
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
