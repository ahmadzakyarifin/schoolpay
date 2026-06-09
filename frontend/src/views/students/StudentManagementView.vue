<script setup>
import { ref, watch, onMounted, computed, reactive } from 'vue'
import { useRouter } from 'vue-router'
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
  History as HistoryIcon,
  ArrowUp as ArrowUpIcon,
  GraduationCap as GraduationCapIcon,
  Database as DatabaseIcon
} from 'lucide-vue-next'

// Services & Composables
import studentService from '../../services/student.service'
import { useForm } from '../../composables/useForm'
import { useAuthStore } from '../../store/auth'

// Components
import StudentTable from '../../components/students/StudentTable.vue'
import StudentFilter from '../../components/students/StudentFilter.vue'
import StudentFormModal from '../../components/students/StudentFormModal.vue'
import StudentBulkActions from '../../components/students/StudentBulkActions.vue'
import ParentListModal from '../../components/students/ParentListModal.vue'

// --- STATE ---
const list = ref([])
const total = ref(0)
const page = ref(1)
const limit = ref(10)
const search = ref('')
const loading = ref(false)
const isMounted = ref(false)
const showFilters = ref(false)
const showEditModal = ref(false)
const isEditing = ref(false)
const selectedStudentIds = ref([])
const showHistory = ref(false)
const birthDateDisplay = ref('')
const deleteLoading = ref(false)
const restoreLoading = ref(false)
const dependencyLoading = ref(false)
const dependencyInfo = ref(null)

const academicFilters = ref({ 
  years: [], majors: [], classes: [],
  provinces: [], regencies: [], districts: [], villages: []
})
const filterParams = reactive({ status: '', entry_year: '', class_id: '', major_id: '', sort: '' })
const tempFilters = reactive({ ...filterParams })

// Modals
const showParentModal = ref(false)
const showBulkGraduateModal = ref(false)
const showPromoteModal = ref(false)

// Data
const activeParents = ref([])
const selectedStudent = ref(null)
const photoPreview = ref(null)

const authStore = useAuthStore()
const isOffline = computed(() => (typeof navigator !== 'undefined' && navigator.onLine === false))
const offlineLockedMessage = 'Aksi ini membutuhkan koneksi server aktif agar data akademik dan export tetap akurat.'

const apiBase = axios.defaults.baseURL
const staticBase = apiBase.replace('/api/', '')

// Form Setup
const initialForm = {
  id: null,
  name: '',
  nisn: '',
  nis: '',
  nik: '',
  gender: 'Laki-laki',
  birth_place: '',
  birth_date: '',
  religion: 'Islam',
  rt: '',
  rw: '',
  village: '',
  district: '',
  city: '',
  province: '',
  address: '',
  email: '',
  phone_number: '',
  country_code: '62',
  entry_year: null,
  status: 'active',
  parent_id: null,
  class_id: null,
  major_id: null,
  image_path: null,
  description: ''
}

const { form, errors, submitting, resetForm, setErrors, clearErrors, clearFieldError } = useForm(initialForm)

// --- ACTIONS ---
const fetchStudents = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      limit: limit.value,
      search: search.value || undefined,
      status: showHistory.value ? 'trash' : (filterParams.status || undefined),
      entry_year: filterParams.entry_year || undefined,
      class_id: filterParams.class_id || undefined,
      major_id: filterParams.major_id || undefined,
      sort: filterParams.sort || undefined
    }
    const res = await studentService.getAll(params)
    const responseData = res.data?.data
    list.value = responseData?.data || []
    total.value = responseData?.total || 0
    totalPages.value = Math.ceil(total.value / limit.value)
  } catch (err) { 
    console.error('Fetch error:', err)
    list.value = []; total.value = 0; totalPages.value = 0 
  } finally { loading.value = false }
}

const fetchAcademicOptions = async () => {
  try {
    const res = await axios.get('students/filters')
    Object.assign(academicFilters.value, res.data.data)
  } catch (err) { console.error('Gagal mengambil opsi akademik') }
}

const fetchActiveParents = async () => {
  if (activeParents.value.length > 0) return
  try {
    const res = await axios.get('users', { params: { role: 'parent', status: 'active', limit: 100 } })
    activeParents.value = res.data.data.users || []
  } catch (err) { console.error('Gagal mengambil data orang tua') }
}

const fetchProvinces = async () => {
  if (academicFilters.value.provinces?.length > 0) return 
  try {
    const res = await fetch('https://ibnux.github.io/data-indonesia/provinsi.json')
    const data = await res.json() || []
    academicFilters.value.provinces = data.map(x => ({ id: x.id, name: x.nama }))
  } catch (err) { console.error('Gagal mengambil data provinsi') }
}

const handleFetchRegions = async (type) => {
  try {
    if (type === 'province') {
      const p = academicFilters.value.provinces.find(x => x.name === form.province)
      if (p) {
        const res = await fetch(`https://ibnux.github.io/data-indonesia/kabupaten/${p.id}.json`)
        const data = await res.json() || []
        academicFilters.value.regencies = data.map(x => ({ id: x.id, name: x.nama }))
      }
    } else if (type === 'city') {
      const c = academicFilters.value.regencies.find(x => x.name === form.city)
      if (c) {
        const res = await fetch(`https://ibnux.github.io/data-indonesia/kecamatan/${c.id}.json`)
        const data = await res.json() || []
        academicFilters.value.districts = data.map(x => ({ id: x.id, name: x.nama }))
      }
    } else if (type === 'district') {
      const d = academicFilters.value.districts.find(x => x.name === form.district)
      if (d) {
        const res = await fetch(`https://ibnux.github.io/data-indonesia/kelurahan/${d.id}.json`)
        const data = await res.json() || []
        academicFilters.value.villages = data.map(x => ({ id: x.id, name: x.nama }))
      }
    }
  } catch (err) { console.error('Gagal mengambil data wilayah') }
}

const applyFilters = () => { Object.assign(filterParams, tempFilters); page.value = 1; showFilters.value = false; fetchStudents() }
const resetFilters = () => {
  Object.assign(tempFilters, { status: '', entry_year: '', class_id: '', major_id: '', sort: '' })
  Object.assign(filterParams, tempFilters)
  search.value = ''
  showFilters.value = false
  page.value = 1
  fetchStudents()
}

const openAddModal = () => {
  isEditing.value = false; resetForm(); photoPreview.value = null; birthDateDisplay.value = ''
  showEditModal.value = true
  // Background fetch
  fetchProvinces(); fetchActiveParents()
}

const openEditModal = async (student) => {
  isEditing.value = true; clearErrors(); photoPreview.value = null
  const countryCodes = ['673', '62', '60', '65', '1']
  let code = '62'
  let phone = student.phone_number || ''
  phone = String(phone).replace(/^\+/, '')
  for (const c of countryCodes) {
    if (phone.startsWith(c)) {
      code = c
      phone = phone.substring(c.length)
      break
    }
  }
  Object.assign(form, { ...student, country_code: code, phone_number: phone })
  
  if (student.birth_date) {
    const [datePart] = String(student.birth_date).split('T')
    const [year, month, day] = datePart.split('-')
    birthDateDisplay.value = year && month && day ? `${day}/${month}/${year}` : ''
  }
  
  showEditModal.value = true
  // Background fetch
  fetchActiveParents()
  
  await fetchProvinces()
  if (form.province) await handleFetchRegions('province')
  if (form.city) await handleFetchRegions('city')
  if (form.district) await handleFetchRegions('district')
}

const handleLocalValidation = (localErrors) => {
  errors.value = { ...errors.value, ...localErrors }
}

const setFieldError = ({ field, messages }) => {
  errors.value = { ...errors.value, [field]: messages }
}

const saveStudent = async () => {
  submitting.value = true
  try {
    const formData = new FormData()
    const normalizedPhone = `${form.country_code || '62'}${String(form.phone_number || '').replace(/^0+/, '')}`
    
    // Append all form fields to FormData
    Object.keys(form).forEach(key => {
      if (key === 'image_path') {
        // Hanya kirim jika ini adalah File (foto baru yang diunggah)
        if (form[key] instanceof File) {
          formData.append(key, form[key])
        }
      } else if (key === 'phone_number') {
        formData.append(key, normalizedPhone)
      } else if (key !== 'country_code' && form[key] !== null && form[key] !== undefined) {
        formData.append(key, form[key])
      }
    })

    if (isEditing.value) {
      await studentService.update(form.id, formData)
    } else {
      await studentService.create(formData)
    }

    showNotification(isEditing.value ? 'Siswa berhasil diperbarui' : 'Siswa berhasil ditambahkan', 'success')
    showEditModal.value = false; fetchStudents()
  } catch (err) {
    setErrors(err)
    const msg = err.response?.data?.message || 'Gagal menyimpan data, periksa isian form Anda'
    showNotification(msg, 'error')
  } finally {
    submitting.value = false
  }
}



const handlePhotoUpload = (file) => {
  form.image_path = file
  const reader = new FileReader()
  reader.onload = (e) => photoPreview.value = e.target.result
  reader.readAsDataURL(file)
}

// Notifications
const notification = reactive({ show: false, message: '', type: 'success' })
const showNotification = (msg, type = 'success') => {
  notification.message = msg; notification.type = type; notification.show = true
  setTimeout(() => notification.show = false, 4000)
}

// Bulk Actions
const handlePromote = async (data) => {
  try {
    await studentService.bulkPromote(data)
    showNotification('Perpindahan kelas berhasil diproses', 'success'); showPromoteModal.value = false; fetchStudents()
  } catch (err) { showNotification(err.response?.data?.message || 'Gagal proses perpindahan kelas', 'error') }
}

const handleGraduate = async (data) => {
  try {
    await studentService.bulkGraduate(data)
    showNotification('Siswa berhasil diluluskan', 'success'); showBulkGraduateModal.value = false; fetchStudents()
  } catch (err) { showNotification(err.response?.data?.message || 'Gagal proses kelulusan', 'error') }
}

const handleDelete = async () => {
  if (deleteBlocked.value) return
  deleteLoading.value = true
  try {
    if (isBulkDelete.value) {
      await studentService.bulkDelete(selectedStudentIds.value)
      showNotification(`${selectedStudentIds.value.length} data siswa berhasil dihapus`, 'success')
      selectedStudentIds.value = []
    } else {
      await studentService.delete(studentToDelete.value.id)
      showNotification(`Siswa ${studentToDelete.value.name} berhasil dihapus`, 'success')
    }
    showDeleteConfirm.value = false
    fetchStudents()
  } catch (err) {
    const errorMsg = err.response?.data?.message || 'Gagal menghapus data'
    showNotification(errorMsg, 'error')
  } finally {
    deleteLoading.value = false
  }
}

const handleRestore = async () => {
  restoreLoading.value = true
  try {
    if (isBulkRestore.value) {
      await studentService.bulkRestore(selectedStudentIds.value)
      showNotification(`${selectedStudentIds.value.length} data siswa berhasil dipulihkan`, 'success')
      selectedStudentIds.value = []
    } else {
      await studentService.restore(studentToRestore.value.id)
      showNotification(`Siswa ${studentToRestore.value.name} berhasil dipulihkan`, 'success')
    }
    showRestoreConfirm.value = false
    fetchStudents()
  } catch (err) {
    const errorMsg = err.response?.data?.message || 'Gagal memulihkan data'
    showNotification(errorMsg, 'error')
  } finally {
    restoreLoading.value = false
  }
}

const handleToggleStatus = async (student) => {
  try {
    await studentService.toggleStatus(student.id)
    showNotification(`Status ${student.name} berhasil diubah`, 'success')
    fetchStudents()
  } catch (err) {
    showNotification('Gagal mengubah status', 'error')
  }
}

const goToParent = (parentId) => {
  router.push({ name: 'user-details', params: { id: parentId } })
}

// State Preservation Logic
const STATE_KEY = 'student_management_state'

const saveState = () => {
  const state = {
    search: search.value,
    tempFilters: tempFilters,
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
      Object.assign(filterParams, state.tempFilters || {})
      page.value = state.page || 1
      limit.value = state.limit || 10
      showHistory.value = !!state.showHistory
    } catch (e) {
      console.error('Gagal memuat state:', e)
    }
  }
}

// Watchers & Lifecycle
watch(showHistory, (newVal) => {
  if (newVal) {
    showFilters.value = false
    Object.assign(tempFilters, { status: '', entry_year: '', class_id: '', major_id: '', sort: '' })
    Object.assign(filterParams, { status: '', entry_year: '', class_id: '', major_id: '', sort: '' })
  }
  page.value = 1
  selectedStudentIds.value = []
  fetchStudents()
})

watch([search, page, limit, showHistory], saveState)
watch(tempFilters, saveState, { deep: true })

watch(search, () => {
  if (page.value === 1) fetchStudents()
  else page.value = 1
})

watch(page, () => {
  fetchStudents()
})

watch(limit, () => {
  if (page.value === 1) fetchStudents()
  else page.value = 1
})

onMounted(async () => {
  loadState()
  isMounted.value = true
  await fetchAcademicOptions()
  fetchActiveParents()
  fetchStudents()
})

// Pagination Logic

const totalPages = computed(() => Math.ceil(total.value / limit.value) || 1)
const visiblePages = computed(() => {
  const pages = []
  let startPage = Math.max(1, page.value - 1)
  let endPage = Math.min(totalPages.value, startPage + 2)
  if (endPage - startPage < 2) startPage = Math.max(1, endPage - 2)
  for (let i = startPage; i <= endPage; i++) if (i > 0) pages.push(i)
  return pages
})

const toggleSelectAll = (checked) => { selectedStudentIds.value = checked ? list.value.map(s => s.id) : [] }
const toggleSelectUser = (id) => {
  const idx = selectedStudentIds.value.indexOf(id)
  if (idx > -1) selectedStudentIds.value.splice(idx, 1)
  else selectedStudentIds.value.push(id)
}

const studentToDelete = ref(null); const isBulkDelete = ref(false); const showDeleteConfirm = ref(false)
const deleteBlocked = computed(() => {
  if (!showDeleteConfirm.value) return false
  return !isBulkDelete.value && !!dependencyInfo.value?.has_dependencies
})
const confirmDelete = async (s) => { 
  studentToDelete.value = s; 
  isBulkDelete.value = false; 
  showDeleteConfirm.value = true 
  dependencyLoading.value = true
  dependencyInfo.value = null
  try {
    const res = await axios.get(`students/${s.id}/dependency-info`)
    dependencyInfo.value = res.data?.data
  } catch (err) {
    dependencyInfo.value = null
  } finally {
    dependencyLoading.value = false
  }
}
const confirmBulkDelete = () => { isBulkDelete.value = true; showDeleteConfirm.value = true }
const studentToRestore = ref(null); const isBulkRestore = ref(false); const showRestoreConfirm = ref(false)
const confirmRestore = (s) => { studentToRestore.value = s; isBulkRestore.value = false; showRestoreConfirm.value = true }
const confirmBulkRestore = () => { isBulkRestore.value = true; showRestoreConfirm.value = true }
const parents = ref([])

const handleExport = async () => {
  if (isOffline.value) {
    showNotification('Export Excel tidak tersedia saat offline. Server harus online agar file berisi data terbaru.', 'error')
    return
  }
  try {
    const response = await axios.get('students/export', {
      params: { 
        search: search.value, 
        status: showHistory.value ? 'trash' : filterParams.status,
        entry_year: filterParams.entry_year,
        class_id: filterParams.class_id,
        major_id: filterParams.major_id
      },
      responseType: 'blob'
    })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    const dateStr = new Date().toISOString().slice(0, 10)
    link.setAttribute('download', `Data_Siswa_${dateStr}.xlsx`)
    document.body.appendChild(link)
    link.click()
    link.remove()
    setTimeout(() => window.URL.revokeObjectURL(url), 1000)
  } catch (err) {
    showNotification('Gagal ekspor data siswa', 'error')
  }
}

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
            <SearchIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300 group-focus-within:text-indigo-600 transition-colors" />
            <input v-model="search" type="text" placeholder="Cari NISN, NIS, atau Nama siswa..." class="search-input-premium" />
          </div>
          <button v-if="!showHistory" @click="showFilters = !showFilters" class="relative p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 flex items-center gap-2 transition-all shadow-sm cursor-pointer">
            <FilterIcon class="w-4 h-4" />
            <span class="text-[10px] font-black uppercase tracking-wider pr-1">Filter</span>
            <span v-if="filterParams.status || filterParams.entry_year || filterParams.class_id || filterParams.major_id || filterParams.sort" class="absolute -top-1 -right-1 w-3 h-3 bg-indigo-600 rounded-full border-2 border-white shadow-sm"></span>
          </button>
          <button @click="resetFilters" class="p-2.5 bg-white text-slate-600 hover:bg-slate-50 rounded-xl border border-slate-200 shadow-sm transition-all group" title="Reset">
            <ResetIcon class="w-4 h-4 group-hover:rotate-180 transition-transform duration-500" />
          </button>
        </div>
        <StudentFilter v-model="showFilters" :filters="tempFilters" :academicFilters="academicFilters" @apply="applyFilters" @reset="resetFilters" />
      </div>
    </Teleport>

    <!-- Main Content -->
    <div class="bg-white border border-slate-200 rounded-xl shadow-sm flex flex-col min-h-[710px] overflow-hidden">
      <!-- Table Header -->
      <div class="p-4 border-b border-slate-100 bg-slate-50/30 flex items-center justify-between">
        <div class="flex items-center gap-2 font-black text-slate-700 text-xs uppercase tracking-widest">
          <StudentIcon class="w-3.5 h-3.5 text-indigo-600 animate-pulse" />
          <span>{{ showHistory ? 'Riwayat Data Terhapus' : 'Data Operasional Siswa' }}</span>
        </div>
        <div class="flex items-center gap-3">
          <StudentBulkActions :selectedCount="selectedStudentIds.length" :status="showHistory ? 'trash' : 'active'" :academicFilters="academicFilters" :showPromote="showPromoteModal" :showGraduate="showBulkGraduateModal" :loading="submitting" :isOffline="isOffline" @close="showPromoteModal = false; showBulkGraduateModal = false" @promote="handlePromote" @graduate="handleGraduate" @delete="confirmBulkDelete" @restore="confirmBulkRestore" />
          <template v-if="!showHistory">
            <button @click="!isOffline && (showPromoteModal = true)" :disabled="isOffline" :title="isOffline ? offlineLockedMessage : 'Pindah / naik kelas'" :class="['font-bold py-1.5 px-3 rounded-xl border text-[10px] flex items-center gap-1.5 transition-all shadow-sm', isOffline ? 'bg-amber-50 border-amber-200 text-amber-700 cursor-not-allowed' : 'bg-white text-slate-600 hover:bg-slate-50']">
              <ArrowUpIcon class="w-3.5 h-3.5" /> <span>Pindah/Naik</span>
            </button>
            <button @click="!isOffline && (showBulkGraduateModal = true)" :disabled="isOffline" :title="isOffline ? offlineLockedMessage : 'Kelulusan masal'" :class="['font-bold py-1.5 px-3 rounded-xl border text-[10px] flex items-center gap-1.5 transition-all', isOffline ? 'bg-amber-50 border-amber-200 text-amber-700 cursor-not-allowed' : 'bg-white text-slate-600 hover:bg-slate-50']">
              <GraduationCapIcon class="w-3.5 h-3.5" /> <span>Lulus Masal</span>
            </button>
          </template>
          <button @click="handleExport" :disabled="isOffline" :title="isOffline ? 'Export Excel membutuhkan server online agar data terbaru.' : 'Ekspor Excel'" :class="['font-bold py-1.5 px-3 rounded-xl border text-[10px] flex items-center gap-1.5 transition-all shadow-sm', isOffline ? 'bg-amber-50 border-amber-200 text-amber-700 cursor-not-allowed' : 'bg-white text-slate-600 hover:bg-slate-50']">
            <ExcelIcon class="w-3.5 h-3.5" :class="isOffline ? 'text-amber-600' : 'text-emerald-600'" /> <span>{{ isOffline ? 'Excel' : 'Ekspor' }}</span>
          </button>
          <button @click="showHistory = !showHistory" class="bg-white text-slate-600 border border-slate-200 hover:bg-slate-50 font-bold py-1.5 px-3 rounded-xl text-[10px] flex items-center gap-1.5 transition-all shadow-sm">
            <HistoryIcon v-if="!showHistory" class="w-3.5 h-3.5" /> <ResetIcon v-else class="w-3.5 h-3.5 rotate-180" />
            <span>{{ showHistory ? 'Kembali' : 'Riwayat Hapus' }}</span>
          </button>
          <button v-if="!showHistory" @click="openAddModal" class="bg-indigo-600 hover:bg-indigo-700 text-white font-black py-1.5 px-4 rounded-xl flex items-center gap-1.5 shadow-md shadow-indigo-100 transition-all text-[10px] uppercase tracking-widest shrink-0">
            <PlusIcon class="w-3.5 h-3.5" /> <span>Tambah Data</span>
          </button>
        </div>
      </div>

      <StudentTable :students="list" :loading="loading" :selectedIds="selectedStudentIds" :showHistory="showHistory" :pagination="{ page, limit, total, totalPages }" :staticBase="staticBase" @view-details="(id) => router.push({ name: 'student-details', params: { id } })" @edit="openEditModal" @delete="confirmDelete" @restore="confirmRestore" @toggle-status="handleToggleStatus" @go-to-parent="goToParent" @view-parents="selectedStudent = $event; showParentModal = true; parents = []; axios.get(`students/${$event.id}/parents`).then(res => parents = res.data.data)" @toggle-select-user="toggleSelectUser" @toggle-select-all="toggleSelectAll" @page-change="(ev) => { if(ev.page) page = ev.page; if(ev.limit) limit = ev.limit; }" />

      <!-- Pagination -->
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
    <StudentFormModal 
      v-model="showEditModal"
      v-model:birthDateDisplay="birthDateDisplay"
      :isEditing="isEditing" :form="form" :errors="errors" :submitting="submitting"
      :academicFilters="academicFilters" :activeParents="activeParents"
      :photoPreview="photoPreview" :staticBase="staticBase"
      @save="saveStudent" @fetch-regions="handleFetchRegions"
      @photo-upload="handlePhotoUpload"
      @local-validation-failed="handleLocalValidation"
      @clear-field-error="clearFieldError"
      @set-field-error="setFieldError"
    />

    <ParentListModal :show="showParentModal" :student="selectedStudent" :parents="parents" @close="showParentModal = false" />


    <!-- Confirm Modals -->
    <Teleport v-if="isMounted" to="body">
      <transition name="fade">
        <div v-if="showDeleteConfirm" class="fixed inset-0 z-[1100] flex items-center justify-center p-6">
          <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showDeleteConfirm = false"></div>
          <div class="bg-white w-full max-w-md relative z-10 overflow-hidden shadow-[0_20px_50px_rgba(0,0,0,0.3)] animate-scale-in !rounded-[2.5rem] p-8 text-center border border-slate-100">
            <div class="w-20 h-20 bg-rose-50 text-rose-500 rounded-[2rem] flex items-center justify-center mx-auto mb-6 border border-rose-100 shadow-xl shadow-rose-500/10 transition-all duration-500">
              <TrashIcon class="w-10 h-10" />
            </div>
            
            <h3 class="text-xl font-black text-slate-900 tracking-tight mb-2">
              {{ isBulkDelete ? 'Hapus Data Terpilih?' : 'Hapus Siswa?' }}
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
                    Siswa ini masih memiliki keterhubungan aktif: <span class="font-bold underline">{{ dependencyInfo.message }}</span>.
                  </p>
                  <p class="text-amber-700/90 text-[10px] font-bold uppercase tracking-wider bg-amber-100/50 py-1 px-2.5 rounded-lg inline-block">
                    💡 Tips: Jika hanya ingin menyembunyikan dari form, gunakan tombol Edit & ubah status menjadi Non-Aktif. Menghapus akan memindahkannya ke Riwayat Penghapusan (Trash).
                  </p>
                </div>
              </div>
            </div>
            
            <p class="text-slate-500 text-[10px] font-bold uppercase tracking-widest mb-8 px-4 leading-relaxed">
              {{ isBulkDelete 
                ? `Apakah Anda yakin ingin menghapus ${selectedStudentIds.length} data yang terpilih? Data akan dipindahkan ke riwayat penghapusan (Trash).` 
                : `Anda akan memindahkan siswa ${studentToDelete?.name} ke riwayat data terhapus (Trash). Data historis tetap aman di database.` 
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

    <Teleport v-if="isMounted" to="body">
      <transition name="fade">
        <div v-if="showRestoreConfirm" class="fixed inset-0 z-[1100] flex items-center justify-center p-6">
          <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="showRestoreConfirm = false"></div>
          <div class="bg-white w-full max-w-md relative z-10 overflow-hidden shadow-[0_20px_50px_rgba(0,0,0,0.3)] animate-scale-in !rounded-[2.5rem] p-8 text-center border border-slate-100">
            <div class="w-20 h-20 bg-emerald-50 text-emerald-600 rounded-[2rem] flex items-center justify-center mx-auto mb-6 border border-emerald-100 shadow-xl shadow-emerald-500/10 transition-all duration-500">
              <SuccessIcon class="w-10 h-10" />
            </div>
            
            <h3 class="text-xl font-black text-slate-900 tracking-tight mb-2">
              {{ isBulkRestore ? 'Pulihkan Data Terpilih?' : 'Pulihkan Siswa?' }}
            </h3>
            
            <p class="text-slate-500 text-[10px] font-bold uppercase tracking-widest mb-8 px-4 leading-relaxed">
              {{ isBulkRestore 
                ? `Apakah Anda yakin ingin memulihkan ${selectedStudentIds.length} data yang terpilih? Data akan dikembalikan ke daftar operasional aktif.` 
                : `Apakah Anda yakin ingin memulihkan kembali siswa ${studentToRestore?.name} ke daftar operasional aktif?` 
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
.animate-fade-in { animation: fadeIn 0.5s ease-out; }
.animate-scale-in { animation: scaleIn 0.3s ease-out; }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
@keyframes scaleIn { from { opacity: 0; transform: scale(0.95); } to { opacity: 1; transform: scale(1); } }
</style>
