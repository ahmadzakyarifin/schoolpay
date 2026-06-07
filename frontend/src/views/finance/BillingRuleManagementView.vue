<script setup>
import { ref, watch, onMounted, computed, reactive } from 'vue'
import axios from 'axios'
import { 
  Search as SearchIcon, 
  Plus as PlusIcon, 
  ChevronLeft as PrevIcon,
  ChevronRight as NextIcon,
  CheckCircle2 as SuccessIcon,
  AlertCircle as AlertIcon,
  Trash as TrashIcon,
  RotateCcw as ResetIcon,
  History as HistoryIcon,
  RefreshCw as RestoreIcon,
  Play as PlayIcon,
  Undo2 as UndoIcon,
  Filter as FilterIcon
} from 'lucide-vue-next'

// Components & Services
import financeService from '../../services/finance.service'
import { useForm } from '../../composables/useForm'
import BillingRuleTable from '../../components/finance/BillingRuleTable.vue'
import BillingRuleFormModal from '../../components/finance/BillingRuleFormModal.vue'
import BillingRuleBulkActions from '../../components/finance/BillingRuleBulkActions.vue'
import BillingRuleDeleteModal from '../../components/finance/BillingRuleDeleteModal.vue'
import BillingRuleRestoreModal from '../../components/finance/BillingRuleRestoreModal.vue'
import BillingRuleGenerateModal from '../../components/finance/BillingRuleGenerateModal.vue'
import BillingRuleFilter from '../../components/finance/BillingRuleFilter.vue'

// State
const list = ref([])
const total = ref(0)
const page = ref(1)
const limit = ref(10)
const search = ref('')
const loading = ref(false)
const isMounted = ref(false)
const showModal = ref(false)
const isEditing = ref(false)
const selectedIds = ref([])
const showHistory = ref(false)

// Supporting Data
const billTypes = ref([])
const classes = ref([])
const majors = ref([])

// Status & Generate Filter
const showFilter = ref(false)
const tempFilters = reactive({
  status: '',
  generate_status: '',
  sort: ''
})
const activeFilters = reactive({
  status: '',
  generate_status: '',
  sort: ''
})

// Delete Confirmation State
const showDeleteConfirm = ref(false)
const itemToDelete = ref(null)
const isBulkDelete = ref(false)
const deleteLoading = ref(false)
const dependencyLoading = ref(false)
const dependencyInfo = ref(null)

// Status Confirmation State
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
  bill_type_ids: [],
  bill_type_id: '',
  class_ids: [],
  class_id: null,
  target_type: 'all',
  target_ids: [],
  target_id: 0,
  amount: 0,
  period_type: 'bulanan',
  allow_installment: true,
  max_installment: null,
  due_day: 10,
  start_date: '',
  end_date: '',
  is_active: true
}

const { form, errors, submitting, setErrors, clearErrors, clearFieldError, resetForm } = useForm(initialForm)

// Fetch Logic
const fetchData = async () => {
  loading.value = true
  try {
    const response = await financeService.getBillingRules({
      page: page.value,
      limit: limit.value,
      search: search.value || undefined,
      status: showHistory.value ? 'trash' : activeFilters.status || undefined,
      generate_status: showHistory.value ? undefined : activeFilters.generate_status || undefined,
      sort: activeFilters.sort || undefined
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
    console.error('Failed to fetch billing rules:', err)
    list.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const fetchDependencies = async () => {
  try {
    const typesRes = await financeService.getBillTypes({ limit: 100 })
    const typesPayload = typesRes.data?.data
    const typesList = typesPayload?.data || (Array.isArray(typesPayload) ? typesPayload : [])
    billTypes.value = typesList.filter(bt => bt.is_active !== false)
  } catch (err) {
    console.error('Gagal mengambil bill types:', err)
  }

  try {
    const classRes = await financeService.getClasses({ limit: 100 })
    const classesPayload = classRes.data?.data
    const classesList = classesPayload?.data || (Array.isArray(classesPayload) ? classesPayload : [])
    classes.value = classesList.filter(c => c.is_active !== false)
  } catch (err) {
    console.error('Gagal mengambil classes:', err)
  }

  try {
    const majorRes = await financeService.getMajors({ limit: 100 })
    const majorsPayload = majorRes.data?.data
    const majorsList = majorsPayload?.data || (Array.isArray(majorsPayload) ? majorsPayload : [])
    majors.value = majorsList.filter(m => m.is_active !== false)
  } catch (err) {
    console.error('Gagal mengambil majors:', err)
  }
}

// Filter Handlers
const applyFilter = () => {
  activeFilters.status = tempFilters.status
  activeFilters.generate_status = tempFilters.generate_status
  activeFilters.sort = tempFilters.sort
  page.value = 1
  showFilter.value = false
  fetchData()
}

const resetFilters = () => {
  tempFilters.status = ''
  tempFilters.generate_status = ''
  tempFilters.sort = ''
  activeFilters.status = ''
  activeFilters.generate_status = ''
  activeFilters.sort = ''
  search.value = ''
  page.value = 1
  showFilter.value = false
  fetchData()
}

// State Preservation Logic
const STATE_KEY = 'billing_rule_management_state'

const saveState = () => {
  const state = {
    search: search.value,
    statusFilter: activeFilters.status,
    generateStatusFilter: activeFilters.generate_status,
    sortFilter: activeFilters.sort,
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
      activeFilters.status = state.statusFilter || ''
      activeFilters.generate_status = state.generateStatusFilter || ''
      activeFilters.sort = state.sortFilter || ''
      tempFilters.status = state.statusFilter || ''
      tempFilters.generate_status = state.generateStatusFilter || ''
      tempFilters.sort = state.sortFilter || ''
      page.value = state.page || 1
      limit.value = state.limit || 10
    } catch (e) {
      console.error('Gagal memuat state:', e)
    }
  }
}

watch([search, page, limit, () => activeFilters.status, () => activeFilters.generate_status, () => activeFilters.sort], saveState)

watch(showHistory, (newVal) => {
  if (newVal) {
    showFilter.value = false
    Object.assign(activeFilters, { status: '', generate_status: '', sort: '' })
    Object.assign(tempFilters, { status: '', generate_status: '', sort: '' })
  }
  page.value = 1
  selectedIds.value = []
  fetchData()
})

onMounted(() => {
  loadState()
  isMounted.value = true
  fetchData()
  fetchDependencies()
})

const startDateDisplay = ref('')
const endDateDisplay = ref('')

const openAddModal = () => {
  isEditing.value = false
  resetForm()
  form.bill_type_ids = []
  form.target_ids = []
  form.class_ids = []
  form.target_type = 'all'
  form.period_type = 'bulanan'
  form.allow_installment = true
  form.due_day = 10
  startDateDisplay.value = ''
  endDateDisplay.value = ''
  showModal.value = true
}

const openEditModal = (item) => {
  isEditing.value = true
  clearErrors()
  Object.assign(form, item)
  form.bill_type_ids = [item.bill_type_id]
  form.target_ids = item.target_id ? [item.target_id] : []
  form.class_ids = item.class_id ? [item.class_id] : []

  if (form.start_date) {
    const [y, m, d] = new Date(form.start_date).toISOString().split('T')[0].split('-')
    form.start_date = `${y}-${m}-${d}`
    startDateDisplay.value = `${d}/${m}/${y}`
  } else {
    startDateDisplay.value = ''
  }
  if (form.end_date) {
    const [y, m, d] = new Date(form.end_date).toISOString().split('T')[0].split('-')
    form.end_date = `${y}-${m}-${d}`
    endDateDisplay.value = `${d}/${m}/${y}`
  } else {
    endDateDisplay.value = ''
  }
  showModal.value = true
}

const saveBillingRule = async () => {
  clearErrors()
  submitting.value = true
  try {
    const basePayload = { ...form }
    basePayload.amount = Number(basePayload.amount) || 0
    basePayload.due_day = Number(basePayload.due_day) || 10
    basePayload.max_installment = basePayload.max_installment ? Number(basePayload.max_installment) : null
    basePayload.start_date = basePayload.start_date ? new Date(basePayload.start_date).toISOString() : null
    basePayload.end_date = basePayload.end_date ? new Date(basePayload.end_date).toISOString() : null

    if (isEditing.value) {
      const payload = { ...basePayload }
      payload.bill_type_id = Number(form.bill_type_id || form.bill_type_ids[0])
      payload.target_id = Number(form.target_id || form.target_ids[0]) || 0
      payload.class_id = form.class_id ? Number(form.class_id) : null
      await financeService.updateBillingRule(form.id, payload)
      showNotification('Aturan tagihan berhasil diperbarui', 'success')
      showModal.value = false
      fetchData()
    } else {
      const btIds = form.bill_type_ids.length > 0 ? form.bill_type_ids : [form.bill_type_id].filter(Boolean)
      const tIds = form.target_type === 'all' ? [0] : (form.target_ids.length > 0 ? form.target_ids : [form.target_id].filter(Boolean))

      if (btIds.length === 0) {
        errors.value = { bill_type_id: ['Pilih minimal satu kategori tagihan'] }
        submitting.value = false
        return
      }

      if (form.target_type !== 'all' && tIds.length === 0) {
        errors.value = { target_id: [`Pilih minimal satu ${form.target_type === 'major' ? 'jurusan' : 'kelas'}`] }
        submitting.value = false
        return
      }

      let count = 0
      let lastError = null
      for (const btId of btIds) {
        for (const tId of tIds) {
          const payload = { ...basePayload }
          payload.bill_type_id = Number(btId)
          payload.target_id = Number(tId) || 0
          payload.class_id = null
          
          try {
            await financeService.createBillingRule(payload)
            count++
          } catch (err) {
            lastError = err
            console.warn('Gagal membuat aturan:', err)
          }
        }
      }
      if (count > 0) {
        showNotification(`${count} aturan tagihan berhasil ditambahkan`, 'success')
        showModal.value = false
        fetchData()
      } else if (lastError) {
        setErrors(lastError)
        const msg = lastError.response?.data?.message || 'Gagal menambahkan aturan, periksa inputan Anda.'
        showNotification(msg, 'error')
        submitting.value = false
        return
      } else {
        showNotification('Gagal menambahkan aturan (kemungkinan kombinasi target sudah terdaftar)', 'error')
      }
    }
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
    await financeService.toggleBillingRuleStatus(item.id)
    item.is_active = !item.is_active
    showNotification(`Status aturan berhasil diubah`, 'success')
  } catch (err) {
    const errorMsg = err.response?.data?.message || 'Gagal mengubah status'
    showNotification(errorMsg, 'error')
  } finally {
    statusLoading.value = false
  }
}

// Generate Bills Modal State & Triggers
const showGenerateConfirm = ref(false)
const generateActionType = ref('single') // 'single', 'bulk_generate', 'bulk_cancel'
const ruleIdToGenerate = ref(null)
const generateSubmitting = ref(false)
const isPenyesuaian = ref(false)

const confirmGenerate = (ruleId) => {
  const item = list.value.find(i => i.id === ruleId)
  isPenyesuaian.value = item ? item.bill_count > 0 : false
  ruleIdToGenerate.value = ruleId
  generateActionType.value = 'single'
  
  if (!isPenyesuaian.value) {
    executeGenerateAction({ customReason: '', customMessage: '' })
  } else {
    showGenerateConfirm.value = true
  }
}

const confirmBulkGenerate = () => {
  if (selectedIds.value.length === 0) return
  generateActionType.value = 'bulk_generate'
  showGenerateConfirm.value = true
}

const confirmBulkCancel = () => {
  if (selectedIds.value.length === 0) return
  generateActionType.value = 'bulk_cancel'
  showGenerateConfirm.value = true
}

const executeGenerateAction = async (meta = {}) => {
  generateSubmitting.value = true
  try {
    if (generateActionType.value === 'single') {
      await financeService.generateBills(ruleIdToGenerate.value, meta.customReason, meta.customMessage, meta.skipNotification)
      showNotification('Berhasil mengenerate tagihan ke semua siswa target!', 'success')
    } else if (generateActionType.value === 'bulk_generate') {
      await financeService.bulkGenerateBills(selectedIds.value, meta.customReason, meta.customMessage, meta.skipNotification)
      showNotification(`Berhasil mengenerate tagihan untuk ${selectedIds.value.length} aturan!`, 'success')
      selectedIds.value = []
    } else if (generateActionType.value === 'bulk_cancel') {
      await axios.post('finance/generate-bills/bulk-cancel', { rule_ids: selectedIds.value, custom_reason: meta.customReason, custom_message: meta.customMessage, skip_notification: meta.skipNotification })
      showNotification(`Berhasil menarik tagihan untuk ${selectedIds.value.length} aturan!`, 'success')
      selectedIds.value = []
    }
    showGenerateConfirm.value = false
    fetchData()
  } catch (err) {
    const msg = err.response?.data?.message || 'Proses gagal dilakukan'
    showNotification(msg, 'error')
  } finally {
    generateSubmitting.value = false
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
    const res = await financeService.getBillingRuleDependencyInfo(item.id)
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
  if (!isBulkDelete.value && dependencyInfo.value?.has_dependencies) return
  deleteLoading.value = true
  try {
    if (isBulkDelete.value) {
      await financeService.bulkDeleteBillingRules(selectedIds.value)
      showNotification(`${selectedIds.value.length} aturan berhasil dihapus`, 'success')
      selectedIds.value = []
    } else {
      await financeService.deleteBillingRule(itemToDelete.value.id)
      showNotification(`Aturan berhasil dihapus`, 'success')
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
      await financeService.bulkRestoreBillingRules(selectedIds.value)
      showNotification(`${selectedIds.value.length} aturan berhasil dipulihkan`, 'success')
      selectedIds.value = []
    } else {
      await financeService.restoreBillingRule(itemToRestore.value.id)
      showNotification(`Aturan berhasil dipulihkan`, 'success')
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

// Auto-fill nominal from selected Bill Type
watch(() => form.bill_type_id, (val) => {
  if (!val || isEditing.value) return
  const selectedType = billTypes.value.find(bt => Number(bt.id) === Number(val))
  if (selectedType) {
    form.amount = selectedType.default_amount || 0
  }
})

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
      <div class="flex items-center justify-center w-full gap-4 relative mx-auto font-inter">
        <div class="flex items-center justify-center gap-2 flex-1 max-w-2xl mx-auto">
          <div class="relative flex-1 group">
            <SearchIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300 group-focus-within:text-indigo-600" />
            <input v-model="search" type="text" placeholder="Cari aturan tagihan..." class="search-input-premium" />
          </div>

          <!-- Filter Button -->
          <button v-if="!showHistory" @click="showFilter = !showFilter" class="relative p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 flex items-center gap-2 transition-all shadow-sm cursor-pointer">
            <FilterIcon class="w-4 h-4" />
            <span class="text-[10px] font-black uppercase tracking-wider pr-1">Filter</span>
            <span v-if="activeFilters.status || activeFilters.generate_status || activeFilters.sort" class="absolute -top-1 -right-1 w-3 h-3 bg-indigo-600 rounded-full border-2 border-white shadow-sm"></span>
          </button>

          <button @click="resetFilters" class="p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm transition-all group shrink-0 cursor-pointer" title="Reset">
            <ResetIcon class="w-4 h-4 group-hover:rotate-180 transition-transform duration-500" />
          </button>
        </div>

        <!-- Filter Dropdown Component -->
        <BillingRuleFilter
          v-model="showFilter"
          :filters="tempFilters"
          @apply="applyFilter"
          @reset="resetFilters"
        />
      </div>
    </Teleport>

    <!-- Main Content Table Card -->
    <div class="bg-white rounded-xl border border-slate-200 shadow-sm flex flex-col min-h-[710px] overflow-hidden">
      <div class="p-4 border-b border-slate-100 bg-slate-50/30 flex items-center justify-between">
        <div class="flex items-center gap-3">
          <div class="w-1.5 h-5 bg-indigo-500 rounded-full"></div>
          <h3 class="font-black text-slate-700 text-xs uppercase tracking-widest">{{ showHistory ? 'Riwayat Data Terhapus' : 'Otomatisasi Aturan Tagihan' }}</h3>
        </div>

        <div class="flex items-center gap-3">
          <BillingRuleBulkActions 
            :selectedCount="selectedIds.length" 
            :status="showHistory ? 'trash' : activeFilters.status"
            @delete="confirmBulkDelete"
            @restore="confirmBulkRestore"
            @generate="confirmBulkGenerate"
            @cancel-generate="confirmBulkCancel"
          />
          <button @click="showHistory = !showHistory" class="bg-white text-slate-600 border border-slate-200 hover:bg-slate-50 font-bold py-1.5 px-3 rounded-lg text-[10px] flex items-center gap-2 transition-all shadow-sm cursor-pointer animate-fade-in">
            <HistoryIcon v-if="!showHistory" class="w-3.5 h-3.5" />
            <ResetIcon v-else class="w-3.5 h-3.5 rotate-180" />
            <span>{{ showHistory ? 'Kembali ke Data Aktif' : 'Lihat Riwayat Hapus' }}</span>
          </button>
          <button v-if="!showHistory" @click="openAddModal" class="bg-indigo-600 hover:bg-indigo-700 text-white font-black py-1.5 px-3 rounded-lg text-[10px] flex items-center gap-2 shadow-lg shadow-indigo-100 cursor-pointer uppercase tracking-wider">
            <PlusIcon class="w-3.5 h-3.5" />
            <span>Buat Aturan Baru</span>
          </button>
        </div>
      </div>

      <BillingRuleTable 
        :list="list" 
        :loading="loading" 
        :selectedIds="selectedIds"
        :status="showHistory ? 'trash' : activeFilters.status"
        @edit="openEditModal"
        @delete="confirmDelete"
        @restore="confirmRestore"
        @toggle-status="toggleStatus"
        @generate-bills="confirmGenerate"
        @toggle-select-all="toggleSelectAll"
        @toggle-select-item="toggleSelectItem"
      />

      <!-- Pagination -->
      <div class="px-6 py-4 bg-slate-50/50 border-t border-slate-100 flex items-center justify-between mt-auto">
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

          <!-- Page Numbers (Max 3) -->
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
    <BillingRuleFormModal 
      v-model="showModal"
      :isEditing="isEditing"
      :form="form"
      :errors="errors"
      :submitting="submitting"
      :billTypes="billTypes"
      :classes="classes"
      :majors="majors"
      v-model:startDateDisplay="startDateDisplay"
      v-model:endDateDisplay="endDateDisplay"
      @save="saveBillingRule"
      @local-validation-failed="handleLocalValidation"
      @clear-field-error="clearFieldError"
    />

    <!-- Delete Confirmation Modal -->
    <BillingRuleDeleteModal 
      v-model="showDeleteConfirm"
      :isBulkDelete="isBulkDelete"
      :selectedCount="selectedIds.length"
      :dependencyLoading="dependencyLoading"
      :dependencyInfo="dependencyInfo"
      :deleteLoading="deleteLoading"
      :showHistory="showHistory"
      @confirm="handleDelete"
    />

    <!-- Restore Confirmation Modal -->
    <BillingRuleRestoreModal 
      v-model="showRestoreConfirm"
      :isBulkRestore="isBulkRestore"
      :selectedCount="selectedIds.length"
      :restoreLoading="restoreLoading"
      @confirm="handleRestore"
    />

    <!-- Generate Confirmation Modal -->
    <BillingRuleGenerateModal 
      v-model="showGenerateConfirm"
      :generateActionType="generateActionType"
      :selectedCount="selectedIds.length"
      :generateSubmitting="generateSubmitting"
      :isPenyesuaian="isPenyesuaian"
      @confirm="executeGenerateAction"
    />
  </div>
</template>

<style scoped lang="postcss">
.search-input-premium {
  @apply w-full bg-white border border-slate-200 rounded-xl py-2.5 pl-12 pr-4 text-xs font-bold text-slate-700 outline-none transition-all focus:border-indigo-500 focus:ring-4 focus:ring-indigo-50 shadow-sm;
}
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
