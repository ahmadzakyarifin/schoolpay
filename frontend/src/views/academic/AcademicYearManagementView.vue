<script setup>
import { ref, onMounted, reactive, computed, watch } from 'vue'
import axios from 'axios'
import AcademicYearTable from '../../components/academic/AcademicYearTable.vue'
import AcademicYearFormModal from '../../components/academic/AcademicYearFormModal.vue'
import AcademicYearBulkActions from '../../components/academic/AcademicYearBulkActions.vue'
import AcademicYearFilter from '../../components/academic/AcademicYearFilter.vue'
import { 
  Plus as PlusIcon, 
  Calendar as YearIcon,
  Search as SearchIcon,
  RotateCcw as ResetIcon,
  ChevronLeft as PrevIcon,
  ChevronRight as NextIcon,
  History as HistoryIcon,
  CheckCircle2 as SuccessIcon,
  Undo2 as RestoreIcon,
  AlertCircle as AlertIcon,
  Trash as TrashIcon,
  Filter as FilterIcon
} from 'lucide-vue-next'

const list = ref([])
const total = ref(0)
const page = ref(1)
const limit = ref(10)
const loading = ref(false)
const search = ref('')
const statusFilter = ref('')
const sortFilter = ref('')
const showFilters = ref(false)
const tempFilters = reactive({ status: '', sort: '' })
const showModal = ref(false)
const isEditing = ref(false)
const submitting = ref(false)
const showDeleteConfirm = ref(false)
const showRestoreConfirm = ref(false)
const itemToProcess = ref(null)
const selectedIds = ref([])
const isBulkAction = ref(false)
const showHistory = ref(false)
const isMounted = ref(false)
const dependencyLoading = ref(false)
const dependencyInfo = ref(null)
const deleteBlocked = computed(() => {
  if (!showDeleteConfirm.value) return false
  if (dependencyLoading.value) return true
  if (!isBulkAction.value) return !!dependencyInfo.value?.has_dependencies
  if (dependencyInfo.value?.has_dependencies) return true
  return selectedIds.value.some(id => {
    const item = list.value.find(row => row.id === id)
    return (item?.major_count || 0) > 0 || (item?.class_count || 0) > 0 || (item?.student_count || 0) > 0
  })
})

const form = reactive({
  id: null,
  year: new Date().getFullYear(),
  is_active: false,
  major_ids: [],
  class_ids: []
})

const errors = ref({})

const clearFieldError = (field) => {
  if (!errors.value?.[field]) return
  errors.value = { ...errors.value, [field]: '' }
}

const setFieldError = ({ field, messages }) => {
  errors.value = { ...errors.value, [field]: messages }
}

const notification = reactive({
  show: false,
  message: '',
  type: 'success'
})

const showNotification = (message, type = 'success') => {
  notification.message = message
  notification.type = type
  notification.show = true
  setTimeout(() => {
    notification.show = false
  }, 4000)
}

const isOfflineQueuedResponse = (res) => res?.status === 202 || res?.data?.status === 'queued'

const buildYearDependencyInfo = (item) => {
  if (!item) return null
  const messages = []
  if ((item.major_count || 0) > 0) messages.push(`${item.major_count} jurusan`)
  if ((item.class_count || 0) > 0) messages.push(`${item.class_count} kelas`)
  if ((item.student_count || 0) > 0) messages.push(`${item.student_count} siswa`)
  return { has_dependencies: messages.length > 0, message: messages.join(', ') }
}

const dependencyMessageFromError = (message) => {
  const raw = String(message || '').replace(/^gagal:\s*/i, '').trim()
  if (!/tidak dapat dihapus|digunakan oleh|terhubung dengan/i.test(raw)) return ''
  return raw.replace(/^.*?(?:digunakan oleh|terhubung dengan)\s+/i, '').trim()
}

const setDeleteDependencyFromError = (message) => {
  const dependencyMessage = dependencyMessageFromError(message)
  if (!dependencyMessage) return false
  dependencyInfo.value = { has_dependencies: true, message: dependencyMessage }
  dependencyLoading.value = false
  return true
}

const fetchData = async () => {
  try {
    const res = await axios.get('academic/years', {
      params: {
        page: page.value,
        limit: limit.value,
        search: search.value || undefined,
        status: showHistory.value ? 'trash' : statusFilter.value || undefined,
        sort: sortFilter.value || undefined
      }
    })
    
    const responseData = res.data?.data
    list.value = responseData?.data || []
    total.value = responseData?.total || 0
  } catch (err) {
    list.value = []
    total.value = 0
  } finally {
    selectedIds.value = []
  }
}

const applyFilters = () => {
  statusFilter.value = tempFilters.status
  sortFilter.value = tempFilters.sort
  page.value = 1
  showFilters.value = false
  fetchData()
}

const resetFilters = () => {
  search.value = ''
  statusFilter.value = ''
  sortFilter.value = ''
  Object.assign(tempFilters, { status: '', sort: '' })
  page.value = 1
  showFilters.value = false
  fetchData()
}

const openAddModal = () => {
  isEditing.value = false
  form.id = null
  form.year = new Date().getFullYear()
  form.is_active = false
  form.major_ids = []
  form.class_ids = []
  errors.value = {}
  showModal.value = true
}

const openEditModal = (item) => {
  isEditing.value = true
  form.id = item.id
  form.year = item.year
  form.is_active = item.is_active
  form.major_ids = item.major_ids || []
  form.class_ids = item.class_ids || []
  errors.value = {}
  showModal.value = true
}

const handleSubmit = async (payload) => {
  submitting.value = true
  errors.value = {}
  
  try {
    const res = isEditing.value
      ? await axios.put(`academic/years/${form.id}`, payload)
      : await axios.post('academic/years', payload)

    if (isOfflineQueuedResponse(res)) {
      showNotification('Perubahan angkatan disimpan sementara dan akan disinkronkan saat server online')
    } else {
      showNotification(isEditing.value ? 'Angkatan berhasil diperbarui' : 'Angkatan berhasil ditambahkan')
    }
    showModal.value = false
    fetchData()
  } catch (err) {
    const errorData = err.response?.data
    if (errorData?.errors) {
      errors.value = errorData.errors
    } else {
      const msg = errorData?.message || ''
      const msgLower = msg.toLowerCase()
      if (msgLower.includes('jurusan')) {
        errors.value = { major_ids: msg }
      } else if (msgLower.includes('kelas')) {
        errors.value = { class_ids: msg }
      } else if (msgLower.includes('tahun') || msgLower.includes('angkatan')) {
        errors.value = { year: msg }
      } else {
        errors.value = { _general: msg }
        showNotification(msg || 'Gagal menyimpan data', 'error')
      }
    }
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async (item) => {
  isBulkAction.value = false
  itemToProcess.value = item
  showDeleteConfirm.value = true
  dependencyLoading.value = true
  dependencyInfo.value = buildYearDependencyInfo(item)
  try {
    const res = await axios.get(`academic/years/${item.id}/dependency-info`)
    dependencyInfo.value = res.data?.data
  } catch (err) {
    // Fallback count dari tabel tetap dipakai.
  } finally {
    dependencyLoading.value = false
  }
}

const confirmBulkDelete = async () => {
  isBulkAction.value = true
  showDeleteConfirm.value = true
  dependencyLoading.value = true
  dependencyInfo.value = null
  try {
    const results = await Promise.all(selectedIds.value.map(async id => {
      const item = list.value.find(row => row.id === id)
      try {
        const res = await axios.get(`academic/years/${id}/dependency-info`)
        return { item, info: res.data?.data }
      } catch (err) {
        return { item, info: buildYearDependencyInfo(item) }
      }
    }))
    const blocked = results.filter(row => row.info?.has_dependencies)
    if (blocked.length > 0) {
      dependencyInfo.value = {
        has_dependencies: true,
        message: blocked.map(row => `${row.item?.year || 'Angkatan'} (${row.info.message})`).join(', ')
      }
    }
  } finally {
    dependencyLoading.value = false
  }
}

const confirmRestore = (item) => {
  isBulkAction.value = false
  itemToProcess.value = item
  showRestoreConfirm.value = true
}

const confirmBulkRestore = () => {
  isBulkAction.value = true
  showRestoreConfirm.value = true
}

const handleDelete = async () => {
  if (deleteBlocked.value) return
  try {
    if (isBulkAction.value) {
      const res = await axios.post('academic/years/bulk-delete', { ids: selectedIds.value })
      showNotification(isOfflineQueuedResponse(res) ? `${selectedIds.value.length} penghapusan angkatan disimpan sementara untuk sinkron` : `${selectedIds.value.length} data berhasil dihapus`)
      selectedIds.value = []
    } else {
      const res = await axios.delete(`academic/years/${itemToProcess.value.id}`)
      showNotification(isOfflineQueuedResponse(res) ? `Penghapusan angkatan ${itemToProcess.value.year} disimpan sementara untuk sinkron` : 'Angkatan berhasil dipindahkan ke riwayat')
    }
    showDeleteConfirm.value = false
    fetchData()
  } catch (err) {
    const message = err.response?.data?.message || ''
    if (!setDeleteDependencyFromError(message)) {
      showNotification(message || 'Gagal menghapus data', 'error')
    }
  }
}

const handleRestore = async () => {
  try {
    if (isBulkAction.value) {
      const res = await axios.patch('academic/years/bulk-restore', { ids: selectedIds.value })
      showNotification(isOfflineQueuedResponse(res) ? `${selectedIds.value.length} pemulihan angkatan disimpan sementara untuk sinkron` : `${selectedIds.value.length} data berhasil dipulihkan`)
      selectedIds.value = []
    } else {
      const res = await axios.patch(`academic/years/${itemToProcess.value.id}/restore`)
      showNotification(isOfflineQueuedResponse(res) ? `Pemulihan angkatan ${itemToProcess.value.year} disimpan sementara untuk sinkron` : 'Angkatan berhasil dipulihkan')
    }
    showRestoreConfirm.value = false
    fetchData()
  } catch (err) {
    showNotification(err.response?.data?.message || 'Gagal memulihkan data', 'error')
  }
}

const toggleActive = async (item) => {
  try {
    const payload = { ...item, is_active: !item.is_active }
    const res = await axios.put(`academic/years/${item.id}`, payload)
    if (!isOfflineQueuedResponse(res)) item.is_active = !item.is_active
    showNotification(isOfflineQueuedResponse(res) ? `Perubahan status angkatan ${item.year} disimpan sementara untuk sinkron` : `Status angkatan ${item.year} diperbarui`)
  } catch (err) {
    showNotification('Gagal mengubah status', 'error')
  }
}

const totalPages = computed(() => Math.ceil(total.value / limit.value) || 1)
const visiblePages = computed(() => {
  const pages = []
  let startPage = Math.max(1, page.value - 1)
  let endPage = Math.min(totalPages.value, startPage + 2)
  if (endPage - startPage < 2) startPage = Math.max(1, endPage - 2)
  for (let i = startPage; i <= endPage; i++) if (i > 0) pages.push(i)
  return pages
})

watch(page, () => {
  fetchData()
})

watch(limit, () => {
  if (page.value === 1) fetchData()
  else page.value = 1
})

let searchTimeout; 
watch(search, () => { 
  clearTimeout(searchTimeout); 
  searchTimeout = setTimeout(() => { 
    if (page.value === 1) fetchData()
    else page.value = 1
  }, 500) 
})

watch(showHistory, (newVal) => {
  if (newVal) {
    showFilters.value = false
    statusFilter.value = ''
    sortFilter.value = ''
    Object.assign(tempFilters, { status: '', sort: '' })
  }
  page.value = 1
  selectedIds.value = []
  fetchData()
})

watch(statusFilter, () => {
  page.value = 1
  fetchData()
})

// State Preservation Logic
const STATE_KEY = 'academic_year_management_state'

const saveState = () => {
  const state = {
    search: search.value,
    statusFilter: statusFilter.value,
    sortFilter: sortFilter.value,
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
      statusFilter.value = state.statusFilter || ''
      sortFilter.value = state.sortFilter || ''
      Object.assign(tempFilters, { status: statusFilter.value, sort: sortFilter.value })
      page.value = state.page || 1
      limit.value = state.limit || 10
      showHistory.value = !!state.showHistory
    } catch (e) {
      console.error('Gagal memuat state:', e)
    }
  }
}

watch([search, page, limit, showHistory, statusFilter, sortFilter], saveState)

onMounted(() => {
  isMounted.value = true
  loadState()
  fetchData()
})
</script>

<template>
  <div class="space-y-6 animate-fade-in pb-20">
    <Teleport v-if="isMounted" to="#header-actions-target">
      <div class="flex items-center justify-center gap-2 w-full max-w-2xl mx-auto">
        <div class="relative flex-1 group">
          <SearchIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300 group-focus-within:text-indigo-600 transition-colors" />
          <input 
            v-model="search"
            type="text" 
            placeholder="Cari angkatan..." 
            class="search-input-premium"
          />
        </div>

        <button v-if="!showHistory" @click="showFilters = !showFilters" class="relative p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 flex items-center gap-2 transition-all shadow-sm cursor-pointer">
          <FilterIcon class="w-4 h-4" />
          <span class="text-[10px] font-black uppercase tracking-wider pr-1">Filter</span>
          <span v-if="statusFilter || sortFilter" class="absolute -top-1 -right-1 w-3 h-3 bg-indigo-600 rounded-full border-2 border-white shadow-sm"></span>
        </button>

        <button @click="resetFilters" class="p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm group shrink-0 transition-all" title="Reset">
          <ResetIcon class="w-4 h-4 group-hover:rotate-180 transition-transform duration-500" />
        </button>

        <AcademicYearFilter
          v-model="showFilters"
          :filters="tempFilters"
          @apply="applyFilters"
          @reset="resetFilters"
        />
      </div>
    </Teleport>

    <div class="bg-white border border-slate-200 rounded-xl shadow-sm flex flex-col min-h-[710px] transition-all duration-500 overflow-hidden">
      <div class="p-4 border-b border-slate-100 bg-slate-50/30 flex items-center justify-between">
        <div class="flex items-center gap-2 font-black text-slate-700 text-xs uppercase tracking-widest">
          <GraduationCapIcon class="w-3.5 h-3.5 text-indigo-600 animate-pulse" />
          <span>{{ showHistory ? 'Riwayat Penghapusan' : 'Data Operasional Angkatan' }}</span>
        </div>

        <div class="flex items-center gap-3">
          <AcademicYearBulkActions 
            :selectedCount="selectedIds.length" 
            :status="showHistory ? 'trash' : 'active'"
            @delete="confirmBulkDelete"
            @restore="confirmBulkRestore"
          />

          <div class="h-6 w-px bg-slate-200 mx-1" v-if="selectedIds.length > 0"></div>

          <div class="flex items-center gap-2">
            <button @click="showHistory = !showHistory" 
              class="flex items-center gap-1.5 px-3 py-1.5 bg-white hover:bg-slate-50 text-slate-600 border border-slate-200 rounded-xl transition-all shadow-sm group font-bold">
              <HistoryIcon v-if="!showHistory" class="w-3.5 h-3.5 text-slate-600" />
              <RestoreIcon v-else class="w-3.5 h-3.5 text-slate-600" />
              <span class="text-[9px] font-black uppercase tracking-widest">{{ showHistory ? 'Kembali ke Data Aktif' : 'Lihat Riwayat Hapus' }}</span>
            </button>

            <button v-if="!showHistory" @click="openAddModal" class="bg-indigo-600 hover:bg-indigo-700 text-white font-black py-1.5 px-4 rounded-xl flex items-center gap-1.5 shadow-md shadow-indigo-100 transition-all text-[10px] uppercase tracking-widest shrink-0">
              <PlusIcon class="w-3.5 h-3.5" />
              <span>Tambah Data</span>
            </button>
          </div>
        </div>
      </div>

      <AcademicYearTable 
        :list="list" 
        :loading="loading" 
        :showHistory="showHistory" 
        v-model:selectedIds="selectedIds"
        @edit="openEditModal" 
        @delete="confirmDelete" 
        @restore="confirmRestore"
        @toggle-status="toggleActive"
      />

      <div class="px-6 py-4 bg-slate-50/50 border-t border-slate-100 flex items-center justify-between">
        <div class="flex items-center gap-6">
          <div class="flex items-center gap-3">
            <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Tampilkan</span>
            <select v-model="limit" class="bg-white border border-slate-200 rounded-lg text-[10px] font-black text-slate-600 px-2 py-1 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all cursor-pointer shadow-sm">
              <option :value="10">10</option><option :value="25">25</option><option :value="50">50</option><option :value="100">100</option>
            </select>
          </div>
          <div class="h-8 w-px bg-slate-200 hidden sm:block"></div>
          <span class="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em]">
            Halaman <span class="text-indigo-600">{{ page }}</span> dari {{ totalPages }} <span class="mx-2 text-slate-300">|</span> Total <span class="text-indigo-600">{{ total }}</span> Data
          </span>
        </div>

        <div class="flex items-center gap-2">
          <button v-if="totalPages > 1" @click="page--" :disabled="page === 1 || loading" 
            class="w-8 h-8 bg-white border border-slate-200 rounded-lg text-slate-400 hover:text-indigo-600 hover:border-indigo-100 disabled:opacity-20 transition-all shadow-sm flex items-center justify-center cursor-pointer">
            <PrevIcon class="w-3.5 h-3.5" />
          </button>
          
          <div class="flex items-center gap-1">
            <button v-for="p in visiblePages" :key="p" @click="page = p"
              class="w-8 h-8 rounded-lg text-[10px] font-black transition-all flex items-center justify-center cursor-pointer"
              :class="page === p ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/20' : 
                                 'bg-white border border-slate-200 text-slate-500 hover:bg-slate-50 hover:border-slate-300'">
              {{ p }}
            </button>
          </div>

          <button v-if="totalPages > 1" @click="page++" :disabled="page >= totalPages || loading" 
            class="w-8 h-8 bg-white border border-slate-200 rounded-lg text-slate-400 hover:text-indigo-600 hover:border-indigo-100 disabled:opacity-20 transition-all shadow-sm flex items-center justify-center cursor-pointer">
            <NextIcon class="w-3.5 h-3.5" />
          </button>
        </div>
      </div>
    </div>

    <AcademicYearFormModal 
      v-model="showModal" 
      :isEditing="isEditing" 
      :form="form" 
      :errors="errors" 
      :submitting="submitting"
      @save="handleSubmit"
      @clear-field-error="clearFieldError"
      @set-field-error="setFieldError"
    />

    <Teleport v-if="isMounted" to="body">
      <transition name="fade">
        <div v-if="showDeleteConfirm" class="fixed inset-0 z-[1100] flex items-center justify-center p-6">
          <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showDeleteConfirm = false"></div>
          <div class="bg-white w-full max-w-md relative z-10 overflow-hidden shadow-[0_20px_50px_rgba(0,0,0,0.3)] animate-scale-in !rounded-[2.5rem] p-8 text-center border border-slate-100">
            <div class="w-20 h-20 bg-rose-50 text-rose-500 rounded-[2rem] flex items-center justify-center mx-auto mb-6 border border-rose-100 shadow-xl shadow-rose-500/10 transition-all duration-500">
              <TrashIcon class="w-10 h-10" />
            </div>
            
            <h3 class="text-xl font-black text-slate-900 tracking-tight mb-1">
              Hapus Angkatan?
            </h3>
            <div v-if="dependencyLoading" class="my-4 py-3 px-4 bg-slate-50 rounded-2xl flex items-center justify-center gap-2 text-slate-500 text-[10px] font-bold uppercase tracking-widest">
              <div class="w-3 h-3 border-2 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
              Memeriksa keterhubungan data...
            </div>

            <div v-else-if="deleteBlocked" class="my-4 p-4 bg-amber-50 border border-amber-200/80 rounded-2xl text-left shadow-sm">
              <div class="flex items-start gap-3">
                <AlertIcon class="w-5 h-5 text-amber-600 shrink-0 mt-0.5" />
                <div>
                  <h4 class="text-xs font-black text-amber-900 uppercase tracking-wider mb-1">Perhatian: Data Masih Terhubung</h4>
                  <p class="text-amber-800 text-[11px] font-medium leading-relaxed mb-2">
                    {{ isBulkAction ? 'Beberapa angkatan terpilih masih memiliki jurusan, kelas, atau siswa aktif.' : 'Angkatan ini masih memiliki keterhubungan aktif:' }} <span class="font-bold underline">{{ dependencyInfo.message }}</span>
                  </p>
                  <p class="text-amber-700/90 text-[10px] font-bold uppercase tracking-wider bg-amber-100/50 py-1 px-2.5 rounded-lg inline-block">
                    Tips: Jika hanya ingin menyembunyikan dari form, gunakan tombol Edit & ubah status menjadi Non-Aktif. Menghapus akan memindahkannya ke Riwayat Penghapusan (Trash).
                  </p>
                </div>
              </div>
            </div>
            
            <p class="text-slate-500 text-[10px] font-bold uppercase tracking-widest mb-8 px-4 leading-relaxed">
              {{ isBulkAction 
                ? `Apakah Anda yakin ingin menghapus ${selectedIds.length} data yang terpilih? Data akan dipindahkan ke riwayat penghapusan (Trash).` 
                : `Anda akan memindahkan angkatan ${itemToProcess?.year} ke riwayat data terhapus (Trash). Data historis tetap aman di database.` 
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
                :disabled="deleteBlocked"
                class="py-4 bg-rose-600 text-white font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-rose-700 transition-all shadow-lg shadow-rose-600/20 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ deleteBlocked ? 'Tidak Bisa Dihapus' : 'Ya, Hapus' }}
              </button>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>

    <Teleport v-if="isMounted" to="body">
      <transition name="fade">
        <div v-if="showRestoreConfirm" class="fixed inset-0 z-[1100] flex items-center justify-center p-6">
          <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showRestoreConfirm = false"></div>
          <div class="bg-white w-full max-w-md relative z-10 overflow-hidden shadow-[0_20px_50px_rgba(0,0,0,0.3)] animate-scale-in !rounded-[2.5rem] p-8 text-center border border-slate-100">
            <div class="w-20 h-20 bg-emerald-50 text-emerald-600 rounded-[2rem] flex items-center justify-center mx-auto mb-6 border border-emerald-100 shadow-xl shadow-emerald-500/10 transition-all duration-500">
              <SuccessIcon class="w-10 h-10" />
            </div>
            
            <h3 class="text-xl font-black text-slate-900 tracking-tight mb-2">
              {{ isBulkAction ? 'Pulihkan Data Terpilih?' : 'Pulihkan Angkatan?' }}
            </h3>
            
            <p class="text-slate-500 text-[10px] font-bold uppercase tracking-widest mb-8 px-4 leading-relaxed">
              {{ isBulkAction 
                ? `Apakah Anda yakin ingin memulihkan ${selectedIds.length} data yang terpilih? Data akan dikembalikan ke daftar operasional aktif.` 
                : `Apakah Anda yakin ingin memulihkan kembali angkatan ${itemToProcess?.year} ke daftar operasional aktif?` 
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
                class="py-4 bg-emerald-600 text-white font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-emerald-700 transition-all shadow-lg shadow-emerald-600/20"
              >
                Ya, Pulihkan
              </button>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>

    <!-- Notification -->
    <Teleport to="body">
      <transition name="fade">
        <div v-if="notification.show" class="fixed top-6 right-6 z-[2000] px-6 py-4 rounded-2xl shadow-2xl flex items-center gap-3 animate-scale-in"
          :class="notification.type === 'success' ? 'bg-emerald-600 text-white' : 'bg-rose-600 text-white'">
          <SuccessIcon v-if="notification.type === 'success'" class="w-5 h-5 shrink-0" />
          <AlertIcon v-else class="w-5 h-5 shrink-0" />
          <span class="text-[10px] font-black uppercase tracking-widest">{{ notification.message }}</span>
        </div>
      </transition>
    </Teleport>
  </div>
</template>

<style scoped lang="postcss">
.search-input-premium {
  @apply w-full bg-white border border-slate-200 rounded-xl py-2.5 pl-12 pr-4 text-xs font-bold text-slate-700 outline-none transition-all focus:border-indigo-500 focus:ring-4 focus:ring-indigo-50 shadow-sm;
}
.animate-fade-in { animation: fadeIn 0.5s ease-out; }
.animate-scale-in { animation: scaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1); }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
@keyframes scaleIn { from { opacity: 0; transform: scale(0.9); } to { opacity: 1; transform: scale(1); } }
</style>
