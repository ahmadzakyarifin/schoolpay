<script setup>
import { ref, watch, onMounted, computed, reactive } from 'vue'
import { 
  Search as SearchIcon, 
  Plus as PlusIcon, 
  ChevronLeft as PrevIcon,
  ChevronRight as NextIcon,
  CheckCircle2 as SuccessIcon,
  AlertCircle as AlertIcon,
  Filter as FilterIcon,
  Trash as TrashIcon,
  RotateCcw as ResetIcon,
  History as HistoryIcon,
  RefreshCw as RestoreIcon
} from 'lucide-vue-next'

// Components & Services
import financeService from '../../services/finance.service'
import { useForm } from '../../composables/useForm'
import BillTypeTable from '../../components/finance/BillTypeTable.vue'
import BillTypeFormModal from '../../components/finance/BillTypeFormModal.vue'
import BillTypeFilter from '../../components/finance/BillTypeFilter.vue'
import BillTypeBulkActions from '../../components/finance/BillTypeBulkActions.vue'

// State
const list = ref([])
const total = ref(0)
const page = ref(1)
const limit = ref(10)
const search = ref('')
const loading = ref(false)
const isMounted = ref(false)
const showFilters = ref(false)
const showModal = ref(false)
const isEditing = ref(false)
const selectedIds = ref([])
const showHistory = ref(false)

// Temp Filters
const tempFilters = reactive({
  type: '',
  status: '',
  sort: ''
})

// Delete Confirmation State
const showDeleteConfirm = ref(false)
const itemToDelete = ref(null)
const isBulkDelete = ref(false)
const deleteLoading = ref(false)
const dependencyLoading = ref(false)
const dependencyInfo = ref(null)
const deleteBlocked = computed(() => {
  if (!showDeleteConfirm.value) return false
  return !isBulkDelete.value && !!dependencyInfo.value?.has_dependencies
})

// Status Confirmation State
const showStatusConfirm = ref(false)
const itemToToggle = ref(null)
const statusLoading = ref(false)

const notification = reactive({
  show: false,
  message: '',
  type: 'success'
})

const showNotification = (msg, type = 'success') => {
  notification.message = msg
  notification.type = type
  notification.show = true
  setTimeout(() => notification.show = false, 4000)
}

// Form Logic
const initialForm = {
  id: null,
  name: '',
  description: '',
  type: 'recurring',
  default_amount: 0
}

const { form, errors, submitting, setErrors, clearErrors, clearFieldError, resetForm } = useForm(initialForm)

// Fetch Logic
const fetchData = async () => {
  loading.value = true
  try {
    const response = await financeService.getBillTypes({
      page: page.value,
      limit: limit.value,
      search: search.value || undefined,
      type: tempFilters.type || undefined,
      status: showHistory.value ? 'trash' : tempFilters.status || undefined,
      sort: tempFilters.sort || undefined
    })
    const payload = response.data?.data
    if (payload && typeof payload === 'object' && !Array.isArray(payload)) {
      list.value = payload.data || []
      total.value = payload.total ?? list.value.length
    } else {
      list.value = Array.isArray(payload) ? payload : []
      total.value = list.value.length
    }
  } catch (err) {
    console.error('Failed to fetch bill types:', err)
    list.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// Filter Handlers
const applyFilters = () => {
  page.value = 1
  showFilters.value = false
  fetchData()
}

const resetFilters = () => {
  Object.assign(tempFilters, { type: '', status: '', sort: '' })
  search.value = ''
  page.value = 1
  showFilters.value = false
  fetchData()
}

// State Preservation Logic
const STATE_KEY = 'bill_type_management_state'

const saveState = () => {
  const state = {
    search: search.value,
    tempFilters: { ...tempFilters },
    page: page.value,
    limit: limit.value
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
      page.value = state.page || 1
      limit.value = state.limit || 10
    } catch (e) {
      console.error('Gagal memuat state:', e)
    }
  }
}

watch([search, page, limit, () => tempFilters.type, () => tempFilters.status, () => tempFilters.sort], saveState)

watch(showHistory, (newVal) => {
  if (newVal) {
    showFilters.value = false
    Object.assign(tempFilters, { type: '', status: '', sort: '' })
  }
  page.value = 1
  selectedIds.value = []
  fetchData()
})

onMounted(() => {
  loadState()
  isMounted.value = true
  fetchData()
})

const openAddModal = () => {
  isEditing.value = false
  resetForm()
  showModal.value = true
}

const openEditModal = (item) => {
  isEditing.value = true
  clearErrors()
  Object.assign(form, item)
  showModal.value = true
}

const saveBillType = async () => {
  clearErrors()
  submitting.value = true
  try {
    const payload = { ...form }
    payload.default_amount = Number(payload.default_amount) || 0

    if (isEditing.value) {
      await financeService.updateBillType(form.id, payload)
      showNotification('Kategori tagihan berhasil diperbarui', 'success')
    } else {
      await financeService.createBillType(payload)
      showNotification('Kategori tagihan berhasil ditambahkan', 'success')
    }
    showModal.value = false
    fetchData()
  } catch (err) {
    setErrors(err)
    const msg = err.response?.data?.message || 'Gagal menyimpan data, periksa inputan Anda.'
    showNotification(msg, 'error')
  } finally {
    submitting.value = false
  }
}

const handleLocalValidation = (localErrors) => {
  errors.value = { ...errors.value, ...localErrors }
}

// Status Toggle Logic
const toggleStatus = async (item) => {
  statusLoading.value = true
  try {
    await financeService.toggleBillTypeStatus(item.id)
    item.is_active = !item.is_active
    showNotification(`Status ${item.name} berhasil diubah`, 'success')
  } catch (err) {
    const errorMsg = err.response?.data?.message || 'Gagal mengubah status'
    showNotification(errorMsg, 'error')
  } finally {
    statusLoading.value = false
  }
}


// Delete Logic
const confirmDelete = async (item) => {
  itemToDelete.value = item
  isBulkDelete.value = false
  showDeleteConfirm.value = true
  dependencyLoading.value = true
  dependencyInfo.value = null
  try {
    const res = await financeService.getBillTypeDependencyInfo(item.id)
    dependencyInfo.value = res.data?.data
  } catch (err) {
    dependencyInfo.value = null
  } finally {
    dependencyLoading.value = false
  }
}

const confirmBulkDelete = () => {
  if (selectedIds.value.length === 0) return
  isBulkDelete.value = true
  showDeleteConfirm.value = true
}

const handleDelete = async () => {
  if (deleteBlocked.value) return
  deleteLoading.value = true
  try {
    if (isBulkDelete.value) {
      await financeService.bulkDeleteBillTypes(selectedIds.value)
      showNotification(`${selectedIds.value.length} kategori berhasil dihapus`, 'success')
      selectedIds.value = []
    } else {
      await financeService.deleteBillType(itemToDelete.value.id)
      showNotification(`${itemToDelete.value.name} berhasil dihapus`, 'success')
    }
    showDeleteConfirm.value = false
    fetchData()
  } catch (err) {
    const errorMsg = err.response?.data?.message || 'Gagal menghapus data'
    showNotification(errorMsg, 'error')
  } finally {
    deleteLoading.value = false
  }
}

// Restore Logic
const showRestoreConfirm = ref(false)
const itemToRestore = ref(null)
const isBulkRestore = ref(false)
const restoreLoading = ref(false)

const confirmRestore = (item) => {
  itemToRestore.value = item
  isBulkRestore.value = false
  showRestoreConfirm.value = true
}

const confirmBulkRestore = () => {
  if (selectedIds.value.length === 0) return
  isBulkRestore.value = true
  showRestoreConfirm.value = true
}

const handleRestore = async () => {
  restoreLoading.value = true
  try {
    if (isBulkRestore.value) {
      await financeService.bulkRestoreBillTypes(selectedIds.value)
      showNotification(`${selectedIds.value.length} kategori berhasil dipulihkan`, 'success')
      selectedIds.value = []
    } else {
      await financeService.restoreBillType(itemToRestore.value.id)
      showNotification(`${itemToRestore.value.name} berhasil dipulihkan`, 'success')
    }
    showRestoreConfirm.value = false
    fetchData()
  } catch (err) {
    const errorMsg = err.response?.data?.message || 'Gagal memulihkan data'
    showNotification(errorMsg, 'error')
  } finally {
    restoreLoading.value = false
  }
}

// Bulk Selection Handlers
const toggleSelectAll = () => {
  const selectableIds = list.value.map(item => item.id)
  if (selectedIds.value.length === selectableIds.length && selectableIds.length > 0) {
    selectedIds.value = []
  } else {
    selectedIds.value = selectableIds
  }
}

const toggleSelectItem = (id) => {
  const index = selectedIds.value.indexOf(id)
  if (index > -1) {
    selectedIds.value.splice(index, 1)
  } else {
    selectedIds.value.push(id)
  }
}

// Lifecycle & Watchers

watch(search, () => {
  if (page.value === 1) {
    fetchData()
  } else {
    page.value = 1
  }
})

watch(page, () => {
  fetchData()
})

watch(limit, () => {
  if (page.value === 1) {
    fetchData()
  } else {
    page.value = 1
  }
})

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

    <!-- Header Actions Teleport -->
    <Teleport v-if="isMounted" to="#header-actions-target">
      <div class="flex items-center justify-center w-full gap-4 relative mx-auto">
        <div class="flex items-center justify-center gap-2 flex-1 max-w-2xl mx-auto">
          <div class="relative flex-1 group">
            <SearchIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300 group-focus-within:text-indigo-600" />
            <input v-model="search" type="text" placeholder="Cari nama kategori atau deskripsi..." class="search-input-premium" />
          </div>
          <button v-if="!showHistory" @click="showFilters = !showFilters" class="relative p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 flex items-center gap-2 transition-all shadow-sm cursor-pointer">
            <FilterIcon class="w-4 h-4" />
            <span class="text-[10px] font-black uppercase tracking-wider pr-1">Filter</span>
            <span v-if="tempFilters.type || tempFilters.status || tempFilters.sort" class="absolute -top-1 -right-1 w-3 h-3 bg-indigo-600 rounded-full border-2 border-white shadow-sm"></span>
          </button>
          <button @click="resetFilters" class="p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm transition-all group" title="Reset">
            <ResetIcon class="w-4 h-4 group-hover:rotate-180 transition-transform duration-500" />
          </button>
        </div>

        <BillTypeFilter 
          v-model="showFilters" 
          :filters="tempFilters" 
          @apply="applyFilters" 
          @reset="resetFilters" 
        />
      </div>
    </Teleport>

    <!-- Main Content Table Card -->
    <div class="bg-white rounded border border-slate-200 shadow-sm flex flex-col min-h-[710px] overflow-hidden">
      <div class="px-6 py-6 border-b border-slate-100 bg-slate-50/30 flex items-center justify-between">
        <div class="flex items-center gap-4">
          <div class="w-2 h-6 bg-indigo-500 rounded-full"></div>
          <h3 class="font-black text-slate-700 text-sm uppercase tracking-[0.2em]">{{ showHistory ? 'Riwayat Data Terhapus' : 'Master Data Jenis Tagihan' }}</h3>
        </div>

        <div class="flex items-center gap-3">
          <BillTypeBulkActions 
            :selectedCount="selectedIds.length" 
            :status="showHistory ? 'trash' : tempFilters.status" 
            @delete="confirmBulkDelete"
            @restore="confirmBulkRestore"
          />
          <button @click="showHistory = !showHistory" class="bg-white text-slate-600 border border-slate-200 hover:bg-slate-50 font-bold py-2 px-4 rounded-xl text-[10px] flex items-center gap-2 transition-all shadow-sm cursor-pointer">
            <HistoryIcon v-if="!showHistory" class="w-3.5 h-3.5" />
            <ResetIcon v-else class="w-3.5 h-3.5 rotate-180" />
            <span>{{ showHistory ? 'Kembali ke Data Aktif' : 'Lihat Riwayat Hapus' }}</span>
          </button>
          <button v-if="!showHistory" @click="openAddModal" class="bg-indigo-600 hover:bg-indigo-700 text-white font-black py-2 px-5 rounded-xl text-xs flex items-center gap-2 shadow-lg shadow-indigo-100">
            <PlusIcon class="w-4 h-4" />
            <span>Tambah Data</span>
          </button>
        </div>
      </div>

      <BillTypeTable 
        :list="list" 
        :loading="loading" 
        :selectedIds="selectedIds"
        :status="showHistory ? 'trash' : tempFilters.status"
        @edit="openEditModal"
        @delete="confirmDelete"
        @restore="confirmRestore"
        @toggle-status="toggleStatus"
        @toggle-select-all="toggleSelectAll"
        @toggle-select-item="toggleSelectItem"
      />

      <!-- Pagination -->
      <div class="px-8 py-6 bg-slate-50/50 border-t border-slate-100 flex items-center justify-between mt-auto">
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
            class="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-400 hover:text-indigo-600 hover:border-indigo-100 hover:bg-indigo-50/30 disabled:opacity-20 disabled:hover:bg-white disabled:hover:border-slate-200 transition-all cursor-pointer"
          >
            <PrevIcon class="w-4 h-4" />
          </button>

          <!-- Page Numbers (Max 3) -->
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

    <!-- Modals -->
    <BillTypeFormModal 
      v-model="showModal"
      :isEditing="isEditing"
      :form="form"
      :errors="errors"
      :submitting="submitting"
      @save="saveBillType"
      @local-validation-failed="handleLocalValidation"
      @clear-field-error="clearFieldError"
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
              {{ isBulkDelete ? 'Hapus Kategori Terpilih?' : 'Hapus Kategori Tagihan?' }}
            </h3>

            <div v-if="dependencyLoading" class="my-4 py-3 px-4 bg-slate-50 rounded-2xl flex items-center justify-center gap-2 text-slate-500 text-[10px] font-bold uppercase tracking-widest">
              <div class="w-3 h-3 border-2 border-indigo-600 border-t-transparent rounded-full animate-spin"></div>
              Memeriksa keterhubungan data...
            </div>

            <div v-else-if="deleteBlocked" class="my-4 p-4 bg-amber-50 border border-amber-200/80 rounded-2xl text-left shadow-sm">
              <div class="flex items-start gap-3">
                <AlertIcon class="w-5 h-5 text-amber-600 shrink-0 mt-0.5" />
                <div>
                  <h4 class="text-xs font-black text-amber-900 uppercase tracking-wider mb-1">Perhatian: Kategori Masih Digunakan</h4>
                  <p class="text-amber-800 text-[11px] font-medium leading-relaxed mb-2">
                    Kategori tagihan ini masih terhubung dengan: <span class="font-bold underline">{{ dependencyInfo.message }}</span>.
                  </p>
                  <p class="text-amber-700/90 text-[10px] font-bold uppercase tracking-wider bg-amber-100/50 py-1 px-2.5 rounded-lg inline-block">
                    💡 Tips: Jika ingin menonaktifkan pembuatan tagihan baru, gunakan tombol toggle status di tabel menjadi Non-Aktif. Menghapus akan memindahkannya ke Riwayat Penghapusan (Trash).
                  </p>
                </div>
              </div>
            </div>
            
            <p class="text-slate-500 text-[10px] font-bold uppercase tracking-widest mb-8 px-4 leading-relaxed">
              {{ isBulkDelete 
                ? `Apakah Anda yakin ingin menghapus ${selectedIds.length} kategori terpilih? Data akan dipindahkan ke riwayat penghapusan (Trash).` 
                : `Apakah Anda yakin ingin menghapus kategori ${itemToDelete?.name}? Data akan dipindahkan ke riwayat penghapusan (Trash). Data historis tetap aman di database.` 
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
                {{ deleteLoading ? 'Menghapus...' : (deleteBlocked ? 'Tidak Bisa Dihapus' : (showHistory ? 'Ya, Hapus Permanen' : 'Ya, Hapus Kategori')) }}
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
            <div class="w-20 h-20 bg-indigo-50 text-indigo-600 rounded-[2rem] flex items-center justify-center mx-auto mb-6 border border-indigo-100 shadow-xl shadow-indigo-600/10">
              <RestoreIcon class="w-10 h-10" />
            </div>
            
            <h3 class="text-xl font-black text-slate-900 tracking-tight mb-2">
              {{ isBulkRestore ? 'Pulihkan Kategori Terpilih?' : 'Pulihkan Kategori Tagihan?' }}
            </h3>
            
            <p class="text-slate-500 text-[10px] font-bold uppercase tracking-widest mb-8 px-4 leading-relaxed">
              {{ isBulkRestore 
                ? `Apakah Anda yakin ingin memulihkan ${selectedIds.length} kategori terpilih kembali ke daftar aktif?` 
                : `Apakah Anda yakin ingin memulihkan kategori ${itemToRestore?.name} kembali ke daftar aktif?` 
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
                class="py-4 bg-indigo-600 text-white font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-indigo-700 transition-all shadow-lg shadow-indigo-600/20 disabled:opacity-50"
              >
                {{ restoreLoading ? 'Memulihkan...' : 'Ya, Pulihkan Kategori' }}
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
