<script setup>
import { ref, onMounted, reactive, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { 
  Plus as PlusIcon, 
  Search as SearchIcon, 
  RotateCcw as ResetIcon,
  ChevronLeft as PrevIcon,
  ChevronRight as NextIcon,
  History as HistoryIcon,
  CheckCircle2 as SuccessIcon,
  AlertCircle as AlertIcon,
  Trash as TrashIcon,
  Undo2 as RestoreIcon,
  Database as DatabaseIcon,
  Filter as FilterIcon
} from 'lucide-vue-next'

// Services & Components
import majorService from '../../services/major.service'
import MajorTable from '../../components/major/MajorTable.vue'
import MajorFormModal from '../../components/major/MajorFormModal.vue'
import MajorBulkActions from '../../components/major/MajorBulkActions.vue'
import MajorFilter from '../../components/major/MajorFilter.vue'

const router = useRouter()

// --- STATE ---
const list = ref([])
const total = ref(0)
const page = ref(1)
const limit = ref(10)
const search = ref('')
const loading = ref(false)
const showHistory = ref(false)
const statusFilter = ref('')
const sortFilter = ref('')
const showFilters = ref(false)
const tempFilters = reactive({ status: '', sort: '' })
const isMounted = ref(false)
const selectedIds = ref([])

const showEditModal = ref(false)
const isEditing = ref(false)
const submitting = ref(false)

const form = reactive({ id: null, code: '', name: '', is_active: true })
const formErrors = ref({})
const editDependencyInfo = ref(null)

// Delete/Restore Confirmation
const showDeleteConfirm = ref(false)
const showRestoreConfirm = ref(false)
const itemToProcess = ref(null)
const isBulkAction = ref(false)
const dependencyLoading = ref(false)
const dependencyInfo = ref(null)
const deleteLoading = ref(false)
const deleteBlocked = computed(() => {
  if (!showDeleteConfirm.value) return false
  if (dependencyLoading.value) return true
  if (!isBulkAction.value) return !!dependencyInfo.value?.has_dependencies
  if (dependencyInfo.value?.has_dependencies) return true
  return selectedIds.value.some(id => {
    const item = list.value.find(row => row.id === id)
    return (item?.class_count || 0) > 0 || (item?.student_count || 0) > 0 || (item?.academic_year_count || 0) > 0
  })
})

const restoreLoading = ref(false)

// Notification
const notification = reactive({ show: false, message: '', type: 'success' })
const showNotification = (msg, type = 'success') => {
  notification.message = msg; notification.type = type; notification.show = true
  setTimeout(() => notification.show = false, 4000)
}

const buildMajorDependencyInfo = (item) => {
  if (!item) return null
  const counts = {
    classes: item.class_count || 0,
    students: item.student_count || 0,
    academic_years: item.academic_year_count || 0
  }
  const messages = []
  if (counts.classes > 0) messages.push(`${counts.classes} kelas aktif`)
  if (counts.students > 0) messages.push(`${counts.students} siswa aktif`)
  if (counts.academic_years > 0) messages.push(`${counts.academic_years} angkatan`)
  return {
    has_dependencies: messages.length > 0,
    message: messages.join(', '),
    counts
  }
}

const dependencyMessageFromError = (message) => {
  const raw = String(message || '').replace(/^gagal:\s*/i, '').trim()
  if (!/tidak dapat dihapus|digunakan oleh|terhubung dengan/i.test(raw)) return ''
  const linked = raw.match(/(?:digunakan oleh|terhubung dengan)\s+(.+)$/i)
  return linked?.[1]?.replace(/\.$/, '') || raw
}

const setDeleteDependencyFromError = (message) => {
  const dependencyMessage = dependencyMessageFromError(message)
  if (!dependencyMessage) return false
  dependencyInfo.value = { has_dependencies: true, message: dependencyMessage }
  dependencyLoading.value = false
  return true
}

// --- ACTIONS ---
const fetchData = async () => {
  loading.value = true
  try {
    const res = await majorService.getAll({
      page: page.value,
      limit: limit.value,
      search: search.value || undefined,
      status: showHistory.value ? 'trash' : statusFilter.value,
      sort: sortFilter.value || undefined
    })
    const responseData = res.data?.data
    list.value = responseData?.data || []
    total.value = responseData?.total || 0
  } catch (err) {
    list.value = []; total.value = 0
  } finally {
    loading.value = false
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
  isEditing.value = false; resetForm(); formErrors.value = {}; editDependencyInfo.value = null; showEditModal.value = true
}

const openEditModal = async (item) => {
  isEditing.value = true
  Object.assign(form, { ...item })
  formErrors.value = {}
  editDependencyInfo.value = buildMajorDependencyInfo(item)
  showEditModal.value = true
  try {
    const res = await majorService.getDependencyInfo(item.id)
    editDependencyInfo.value = res.data?.data || editDependencyInfo.value
  } catch (err) {
    // Fallback count dari tabel tetap dipakai.
  }
}

const resetForm = () => { Object.assign(form, { id: null, code: '', name: '', is_active: true }); formErrors.value = {}; editDependencyInfo.value = null }

const clearFieldError = (field) => {
  if (!formErrors.value?.[field]) return
  formErrors.value = { ...formErrors.value, [field]: '' }
}

const setFieldError = ({ field, messages }) => {
  formErrors.value = { ...formErrors.value, [field]: messages }
}

const saveMajor = async () => {
  submitting.value = true
  formErrors.value = {}
  try {
    const res = isEditing.value
      ? await majorService.update(form.id, form)
      : await majorService.create(form)

    showNotification(isEditing.value ? 'Jurusan berhasil diperbarui' : 'Jurusan berhasil ditambahkan')
    showEditModal.value = false; fetchData()
  } catch (err) {
    const errorData = err.response?.data
    if (errorData?.errors) {
      formErrors.value = errorData.errors
    } else {
      const msg = errorData?.message || ''
      if (msg.toLowerCase().includes('tidak dapat diubah')) {
        editDependencyInfo.value = { has_dependencies: true, message: msg.replace(/^gagal:\s*/i, '') }
      } else if (msg.toLowerCase().includes('nama')) {
        formErrors.value = { name: msg }
      } else if (msg.toLowerCase().includes('kode')) {
        formErrors.value = { code: msg }
      } else {
        showNotification(msg || 'Gagal menyimpan data', 'error')
      }
    }
  } finally {
    submitting.value = false
  }
}

const confirmDelete = async (item) => { 
  isBulkAction.value = false
  itemToProcess.value = item; 
  showDeleteConfirm.value = true 
  dependencyLoading.value = true
  dependencyInfo.value = buildMajorDependencyInfo(item)
  try {
    const res = await majorService.getDependencyInfo(item.id)
    dependencyInfo.value = res.data?.data || dependencyInfo.value
  } catch (err) {
    // Tetap pakai count dari tabel supaya tombol hapus tidak lolos ke backend.
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
        const res = await majorService.getDependencyInfo(id)
        return { item, info: res.data?.data }
      } catch (err) {
        return { item, info: buildMajorDependencyInfo(item) }
      }
    }))
    const blocked = results.filter(row => row.info?.has_dependencies)
    if (blocked.length > 0) {
      dependencyInfo.value = {
        has_dependencies: true,
        message: blocked.map(row => `${row.item?.name || 'Jurusan'} (${row.info.message})`).join(', ')
      }
    }
  } finally {
    dependencyLoading.value = false
  }
}
const confirmRestore = (item) => { 
  isBulkAction.value = false
  itemToProcess.value = item; 
  showRestoreConfirm.value = true 
}
const confirmBulkRestore = () => {
  isBulkAction.value = true
  showRestoreConfirm.value = true
}

const handleDelete = async () => {
  if (deleteBlocked.value) return
  deleteLoading.value = true
  try {
    if (isBulkAction.value) {
      await majorService.bulkDelete(selectedIds.value)
      showNotification(`${selectedIds.value.length} data berhasil dihapus`)
      selectedIds.value = []
    } else {
      await majorService.delete(itemToProcess.value.id)
      showNotification(`Jurusan ${itemToProcess.value.name} berhasil dihapus`)
    }
    showDeleteConfirm.value = false
    fetchData()
  } catch (err) {
    const message = err.response?.data?.message || ''
    if (!setDeleteDependencyFromError(message)) {
      showNotification(message || 'Gagal menghapus data', 'error')
    }
  } finally {
    deleteLoading.value = false
  }
}

const handleRestore = async () => {
  restoreLoading.value = true
  try {
    if (isBulkAction.value) {
      await majorService.bulkRestore(selectedIds.value)
      showNotification(`${selectedIds.value.length} data berhasil dipulihkan`)
      selectedIds.value = []
    } else {
      await majorService.restore(itemToProcess.value.id)
      showNotification(`Jurusan ${itemToProcess.value.name} berhasil dipulihkan`)
    }
    showRestoreConfirm.value = false
    fetchData()
  } catch (err) {
    showNotification(err.response?.data?.message || 'Gagal memulihkan data', 'error')
  } finally {
    restoreLoading.value = false
  }
}

const handleToggleStatus = async (item) => {
  try {
    await majorService.toggleStatus(item.id)
    item.is_active = !item.is_active
    showNotification(`Status ${item.name} diperbarui`)
  } catch (err) { showNotification('Gagal mengubah status', 'error') }
}

// Pagination Logic
const totalPages = computed(() => Math.ceil(total.value / limit.value) || 1)
const visiblePages = computed(() => {
  const pages = []
  let startPage = Math.max(1, page.value - 1)
  let endPage = Math.min(totalPages.value, startPage + 2)
  
  if (endPage - startPage < 2) {
    startPage = Math.max(1, endPage - 2)
  }
  
  for (let i = startPage; i <= endPage; i++) {
    if (i > 0) pages.push(i)
  }
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

watch(sortFilter, () => {
  page.value = 1
  fetchData()
})

// State Preservation Logic
const STATE_KEY = 'major_management_state'

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

onMounted(() => { isMounted.value = true; loadState(); fetchData() })
</script>

<template>
  <div class="space-y-6 animate-fade-in pb-20">
    <!-- Header Teleport -->
    <Teleport v-if="isMounted" to="#header-actions-target">
      <div class="flex items-center justify-center gap-2 w-full max-w-3xl mx-auto">
        <div class="relative flex-1 group">
          <SearchIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300 group-focus-within:text-indigo-600 transition-colors" />
          <input v-model="search" type="text" placeholder="Cari nama atau kode jurusan..." class="search-input-premium" />
        </div>

        <button v-if="!showHistory" @click="showFilters = !showFilters" class="relative p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 flex items-center gap-2 transition-all shadow-sm cursor-pointer">
          <FilterIcon class="w-4 h-4" />
          <span class="text-[10px] font-black uppercase tracking-wider pr-1">Filter</span>
          <span v-if="statusFilter || sortFilter" class="absolute -top-1 -right-1 w-3 h-3 bg-indigo-600 rounded-full border-2 border-white shadow-sm"></span>
        </button>

        <button @click="resetFilters" class="p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm group shrink-0 transition-all" title="Reset">
          <ResetIcon class="w-4 h-4 group-hover:rotate-180 transition-transform duration-500" />
        </button>

        <MajorFilter
          v-model="showFilters"
          :filters="tempFilters"
          @apply="applyFilters"
          @reset="resetFilters"
        />
      </div>
    </Teleport>

    <!-- Main Content Card -->
    <div class="bg-white border border-slate-200 rounded-xl shadow-sm flex flex-col min-h-[710px] transition-all duration-500 overflow-hidden">
      <!-- Table Header -->
      <div class="p-4 border-b border-slate-100 bg-slate-50/30 flex items-center justify-between">
        <div class="flex items-center gap-2 font-black text-slate-700 text-xs uppercase tracking-widest">
          <DatabaseIcon class="w-3.5 h-3.5 text-indigo-600" />
          <span>{{ showHistory ? 'Riwayat Penghapusan' : 'Data Operasional Major' }}</span>
        </div>

        <div class="flex items-center gap-3">
          <!-- Bulk Actions Component -->
          <MajorBulkActions 
            :selectedCount="selectedIds.length" 
            :status="showHistory ? 'trash' : 'active'"
            @delete="confirmBulkDelete"
            @restore="confirmBulkRestore"
          />

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

      <!-- Table Component -->
      <MajorTable 
        :list="list" 
        :loading="loading" 
        :showHistory="showHistory" 
        v-model:selectedIds="selectedIds"
        @edit="openEditModal" 
        @delete="confirmDelete" 
        @restore="confirmRestore"
        @toggle-status="handleToggleStatus"
      />

      <!-- Pagination Footer -->
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

    <!-- Modals -->
    <MajorFormModal 
      v-model="showEditModal" 
      :isEditing="isEditing" 
      :form="form" 
      :errors="formErrors" 
      :submitting="submitting"
      :dependencyInfo="editDependencyInfo"
      @save="saveMajor"
      @clear-field-error="clearFieldError"
      @set-field-error="setFieldError"
    />

    <!-- Confirm Delete -->
    <Teleport v-if="isMounted" to="body">
      <transition name="fade">
        <div v-if="showDeleteConfirm" class="fixed inset-0 z-[1100] flex items-center justify-center p-6">
          <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showDeleteConfirm = false"></div>
          <div class="bg-white w-full max-w-md relative z-10 overflow-hidden shadow-[0_20px_50px_rgba(0,0,0,0.3)] animate-scale-in !rounded-[2.5rem] p-8 text-center border border-slate-100">
            <div class="w-20 h-20 bg-rose-50 text-rose-500 rounded-[2rem] flex items-center justify-center mx-auto mb-6 border border-rose-100 shadow-xl shadow-rose-500/10 transition-all duration-500">
              <TrashIcon class="w-10 h-10" />
            </div>
            
            <h3 class="text-xl font-black text-slate-900 tracking-tight mb-2">
              {{ isBulkAction ? 'Hapus Data Terpilih?' : 'Hapus Jurusan?' }}
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
                    {{ isBulkAction ? 'Beberapa jurusan terpilih masih memiliki kelas atau siswa aktif.' : 'Jurusan ini masih memiliki keterhubungan aktif:' }} <span class="font-bold underline">{{ dependencyInfo.message }}</span>
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
                : `Anda akan memindahkan jurusan ${itemToProcess?.name} ke riwayat data terhapus (Trash). Data historis tetap aman di database.` 
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
                :disabled="deleteLoading || deleteBlocked"
                class="py-4 bg-rose-600 text-white font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-rose-700 transition-all shadow-lg shadow-rose-600/20 disabled:opacity-50"
              >
                {{ deleteLoading ? 'Menghapus...' : (deleteBlocked ? 'Tidak Bisa Dihapus' : 'Ya, Hapus') }}
              </button>
            </div>
          </div>
        </div>
      </transition>
    </Teleport>

    <!-- Confirm Restore -->
    <Teleport v-if="isMounted" to="body">
      <transition name="fade">
        <div v-if="showRestoreConfirm" class="fixed inset-0 z-[1100] flex items-center justify-center p-6">
          <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showRestoreConfirm = false"></div>
          <div class="bg-white w-full max-w-md relative z-10 overflow-hidden shadow-[0_20px_50px_rgba(0,0,0,0.3)] animate-scale-in !rounded-[2.5rem] p-8 text-center border border-slate-100">
            <div class="w-20 h-20 bg-emerald-50 text-emerald-600 rounded-[2rem] flex items-center justify-center mx-auto mb-6 border border-emerald-100 shadow-xl shadow-emerald-500/10 transition-all duration-500">
              <SuccessIcon class="w-10 h-10" />
            </div>
            
            <h3 class="text-xl font-black text-slate-900 tracking-tight mb-2">
              {{ isBulkAction ? 'Pulihkan Data Terpilih?' : 'Pulihkan Data?' }}
            </h3>
            
            <p class="text-slate-500 text-[10px] font-bold uppercase tracking-widest mb-8 px-4 leading-relaxed">
              {{ isBulkAction 
                ? `Apakah Anda yakin ingin memulihkan ${selectedIds.length} data yang terpilih? Data akan dikembalikan ke daftar operasional aktif.` 
                : `Apakah Anda yakin ingin memulihkan kembali jurusan ${itemToProcess?.name} ke daftar operasional aktif?` 
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

    <!-- Global Notification -->
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
