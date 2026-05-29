<script setup>
import { reactive, ref, watch } from 'vue'
import axios from 'axios'
import { 
  X as CloseIcon,
  ArrowUp as ArrowUpIcon,
  GraduationCap as GraduationCapIcon,
  Info as InfoIcon,
  CheckCircle2 as SuccessIcon,
  AlertCircle as AlertIcon,
  ChevronDown as ChevronDownIcon,
  Trash as TrashIcon,
  RotateCcw as ResetIcon
} from 'lucide-vue-next'

const props = defineProps({
  showPromote: Boolean,
  showGraduate: Boolean,
  loading: Boolean,
  academicFilters: Object,
  selectedCount: Number,
  status: String,
  isOffline: Boolean
})

const emit = defineEmits(['close', 'promote', 'graduate', 'delete', 'restore'])

const promoteForm = reactive({
  source_class_id: '',
  target_class_id: ''
})

const graduateForm = reactive({
  class_id: ''
})

const eligibleStudents = ref([])
const loadingStudents = ref(false)
const selectedStudentIdsForGraduation = ref([])

const eligiblePromoteStudents = ref([])
const loadingPromoteStudents = ref(false)
const selectedStudentIdsForPromotion = ref([])

watch(() => props.isOffline, (offline) => {
  if (!offline) return
  eligibleStudents.value = []
  selectedStudentIdsForGraduation.value = []
  eligiblePromoteStudents.value = []
  selectedStudentIdsForPromotion.value = []
})

watch(() => graduateForm.class_id, async (newClassId) => {
  if (props.isOffline) {
    eligibleStudents.value = []
    selectedStudentIdsForGraduation.value = []
    return
  }
  if (!newClassId) {
    eligibleStudents.value = []
    selectedStudentIdsForGraduation.value = []
    return
  }
  loadingStudents.value = true
  try {
    const params = {
      class_id: newClassId,
      status: 'active',
      limit: 500
    }
    const res = await axios.get('students', { params })
    eligibleStudents.value = res.data?.data?.data || []
    selectedStudentIdsForGraduation.value = eligibleStudents.value.map(s => s.id)
  } catch (err) {
    console.error('Gagal mengambil data siswa kelas:', err)
    eligibleStudents.value = []
    selectedStudentIdsForGraduation.value = []
  } finally {
    loadingStudents.value = false
  }
})

watch(() => promoteForm.source_class_id, async (newSourceId) => {
  if (props.isOffline) {
    eligiblePromoteStudents.value = []
    selectedStudentIdsForPromotion.value = []
    return
  }
  if (!newSourceId) {
    eligiblePromoteStudents.value = []
    selectedStudentIdsForPromotion.value = []
    return
  }
  loadingPromoteStudents.value = true
  try {
    const params = {
      class_id: newSourceId,
      status: 'active',
      limit: 500
    }
    const res = await axios.get('students', { params })
    eligiblePromoteStudents.value = res.data?.data?.data || []
    selectedStudentIdsForPromotion.value = eligiblePromoteStudents.value.map(s => s.id)
  } catch (err) {
    console.error('Gagal mengambil data siswa kelas asal:', err)
    eligiblePromoteStudents.value = []
    selectedStudentIdsForPromotion.value = []
  } finally {
    loadingPromoteStudents.value = false
  }
})

watch(() => props.showGraduate, (val) => {
  if (!val) {
    graduateForm.class_id = ''
    eligibleStudents.value = []
    selectedStudentIdsForGraduation.value = []
  }
})

watch(() => props.showPromote, (val) => {
  if (!val) {
    promoteForm.source_class_id = ''
    promoteForm.target_class_id = ''
    eligiblePromoteStudents.value = []
    selectedStudentIdsForPromotion.value = []
  }
})

const toggleSelectAllGraduation = (checked) => {
  selectedStudentIdsForGraduation.value = checked ? eligibleStudents.value.map(s => s.id) : []
}

const toggleSelectOneGraduation = (id) => {
  const idx = selectedStudentIdsForGraduation.value.indexOf(id)
  if (idx > -1) selectedStudentIdsForGraduation.value.splice(idx, 1)
  else selectedStudentIdsForGraduation.value.push(id)
}

const toggleSelectAllPromotion = (checked) => {
  selectedStudentIdsForPromotion.value = checked ? eligiblePromoteStudents.value.map(s => s.id) : []
}

const toggleSelectOnePromotion = (id) => {
  const idx = selectedStudentIdsForPromotion.value.indexOf(id)
  if (idx > -1) selectedStudentIdsForPromotion.value.splice(idx, 1)
  else selectedStudentIdsForPromotion.value.push(id)
}

const handlePromote = () => {
  if (props.isOffline) return
  emit('promote', { 
    source_class_id: promoteForm.source_class_id,
    target_class_id: promoteForm.target_class_id,
    student_ids: selectedStudentIdsForPromotion.value
  })
}

const handleGraduate = () => {
  if (props.isOffline) return
  emit('graduate', { 
    class_id: graduateForm.class_id,
    student_ids: selectedStudentIdsForGraduation.value
  })
}
</script>

<template>
  <div class="flex items-center gap-2">
    <!-- Bulk Delete/Restore Buttons (Visible when items selected) -->
    <transition name="fade">
      <div v-if="selectedCount > 0" class="flex items-center gap-2">
        <button 
          v-if="status !== 'trash'" 
          @click="emit('delete')" 
          class="bg-white border border-rose-200 hover:bg-rose-50/50 hover:border-rose-300 text-rose-600 font-black py-2 px-4 rounded-xl text-xs flex items-center gap-2 transition-all shadow-sm cursor-pointer"
        >
          <TrashIcon class="w-3.5 h-3.5 text-rose-500" />
          <span>Hapus Terpilih ({{ selectedCount }})</span>
        </button>
        
        <button 
          v-if="status === 'trash'" 
          @click="emit('restore')" 
          class="bg-white border border-emerald-200 hover:bg-emerald-50/50 hover:border-emerald-300 text-emerald-600 font-black py-2 px-4 rounded-xl text-xs flex items-center gap-2 transition-all shadow-sm cursor-pointer"
        >
          <ResetIcon class="w-3.5 h-3.5 text-emerald-500" />
          <span>Pulihkan Terpilih ({{ selectedCount }})</span>
        </button>
      </div>
    </transition>

    <!-- Bulk Promote Modal -->
    <transition enter-active-class="transition duration-300 ease-out" enter-from-class="opacity-0" enter-to-class="opacity-100" leave-active-class="transition duration-200 ease-in" leave-from-class="opacity-100" leave-to-class="opacity-0">
      <div v-if="showPromote" class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
        <div class="bg-white w-full max-w-2xl rounded-[40px] shadow-2xl overflow-hidden animate-scale-in flex flex-col max-h-[90vh]">
          <div class="px-8 py-8 border-b border-slate-50 flex items-center justify-between bg-indigo-50/30 shrink-0">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-2xl bg-indigo-600 flex items-center justify-center shadow-lg shadow-indigo-100">
                <ArrowUpIcon class="w-6 h-6 text-white" />
              </div>
              <div>
                <h2 class="text-lg font-black text-slate-800 tracking-tight">Pindah / Naik Kelas (Student Wizard)</h2>
                <p class="text-[10px] text-slate-400 font-bold uppercase tracking-widest mt-0.5">Proses masal perpindahan kelas, jenjang, atau rombel</p>
              </div>
            </div>
            <button @click="$emit('close')" class="w-10 h-10 rounded-xl hover:bg-white text-slate-400 hover:text-slate-600 transition-all shadow-sm border border-transparent">
              <CloseIcon class="w-5 h-5" />
            </button>
          </div>

          <div class="p-8 space-y-6 overflow-y-auto custom-scrollbar bg-slate-50/30 flex-1">
            <div v-if="isOffline" class="p-5 bg-amber-50 border border-amber-200 rounded-3xl flex items-start gap-4 shadow-sm shadow-amber-100/60">
              <div class="w-10 h-10 rounded-xl bg-white flex items-center justify-center shrink-0 shadow-sm text-amber-600">
                <AlertIcon class="w-5 h-5" />
              </div>
              <p class="text-[10px] font-black text-amber-800 leading-relaxed uppercase tracking-tight">
                Server sedang offline. Pindah atau naik kelas dikunci karena sistem harus mengecek data akademik dan tagihan aktif terbaru.
              </p>
            </div>

            <div class="p-5 bg-indigo-50/50 border border-indigo-100 rounded-3xl flex items-start gap-4 shadow-sm">
              <div class="w-10 h-10 rounded-xl bg-white flex items-center justify-center shrink-0 shadow-sm shadow-indigo-100">
                <InfoIcon class="w-5 h-5 text-indigo-500" />
              </div>
              <p class="text-[10px] font-bold text-indigo-700 leading-relaxed uppercase tracking-tight">
                Pilih kelas asal dan kelas tujuan. Sistem akan memuat seluruh siswa aktif di kelas asal. Anda dapat menghapus centang pada siswa yang tidak ikut dipindahkan (misal tinggal kelas atau beda jurusan).
              </p>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-3">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1">Dari Kelas Asal <span class="text-rose-500">*</span></label>
                <div class="relative group">
                  <select v-model="promoteForm.source_class_id" class="w-full py-4 px-6 bg-white border border-slate-200 rounded-2xl appearance-none focus:ring-4 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-xs text-slate-700 shadow-sm pr-12 cursor-pointer">
                    <option value="">Pilih kelas asal</option>
                    <option v-for="c in academicFilters.classes" :key="c.id" :value="c.id">{{ c.name }}</option>
                  </select>
                  <div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400 group-focus-within:text-indigo-500 transition-colors">
                    <ChevronDownIcon class="w-4 h-4" />
                  </div>
                </div>
              </div>

              <div class="space-y-3">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1">Ke Kelas Tujuan <span class="text-rose-500">*</span></label>
                <div class="relative group">
                  <select v-model="promoteForm.target_class_id" class="w-full py-4 px-6 bg-white border border-slate-200 rounded-2xl appearance-none focus:ring-4 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-xs text-slate-700 shadow-sm pr-12 cursor-pointer">
                    <option value="">Pilih kelas tujuan</option>
                    <option v-for="c in academicFilters.classes" :key="c.id" :value="c.id">{{ c.name }}</option>
                  </select>
                  <div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400 group-focus-within:text-indigo-500 transition-colors">
                    <ChevronDownIcon class="w-4 h-4" />
                  </div>
                </div>
              </div>
            </div>

            <!-- Student Promotion Selection List -->
            <div v-if="promoteForm.source_class_id" class="bg-white rounded-[2rem] border border-slate-200 shadow-sm overflow-hidden space-y-4 p-6">
              <div class="flex items-center justify-between pb-4 border-b border-slate-100">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center font-black text-xs">
                    {{ selectedStudentIdsForPromotion.length }}
                  </div>
                  <div>
                    <h4 class="text-xs font-black text-slate-800 uppercase tracking-wider">Daftar Siswa Kandidat Pindah/Naik Kelas</h4>
                    <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Centang siswa yang berhak dipindahkan ke kelas tujuan</p>
                  </div>
                </div>
                <label class="flex items-center gap-2 text-[10px] font-black text-slate-600 uppercase tracking-widest cursor-pointer hover:text-indigo-600 transition-colors bg-slate-50 px-4 py-2 rounded-xl border border-slate-200">
                  <input 
                    type="checkbox" 
                    :checked="selectedStudentIdsForPromotion.length === eligiblePromoteStudents.length && eligiblePromoteStudents.length > 0"
                    @change="e => toggleSelectAllPromotion(e.target.checked)"
                    class="rounded border-slate-300 text-indigo-600 focus:ring-indigo-500 w-4 h-4"
                  />
                  <span>Pilih Semua ({{ eligiblePromoteStudents.length }})</span>
                </label>
              </div>

              <div v-if="loadingPromoteStudents" class="py-12 text-center space-y-3">
                <div class="w-8 h-8 border-4 border-indigo-600/20 border-t-indigo-600 rounded-full animate-spin mx-auto"></div>
                <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Memuat daftar siswa aktif di kelas asal...</p>
              </div>

              <div v-else-if="eligiblePromoteStudents.length === 0" class="py-12 text-center space-y-2 bg-slate-50/50 rounded-2xl border border-dashed border-slate-200">
                <p class="text-xs font-black text-slate-600 uppercase tracking-wider">Tidak Ada Siswa Aktif</p>
                <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Tidak ditemukan siswa dengan status aktif di kelas asal terpilih.</p>
              </div>

              <div v-else class="space-y-2.5 max-h-[280px] overflow-y-auto custom-scrollbar pr-2">
                <div v-for="(stu, idx) in eligiblePromoteStudents" :key="stu.id"
                  @click="toggleSelectOnePromotion(stu.id)"
                  class="flex items-center justify-between p-4 rounded-2xl border-2 transition-all cursor-pointer group"
                  :class="selectedStudentIdsForPromotion.includes(stu.id) ? 'border-indigo-600 bg-indigo-50/30' : 'border-slate-100 bg-slate-50 hover:border-indigo-200'">
                  <div class="flex items-center gap-4 min-w-0">
                    <div class="w-7 h-7 rounded-lg flex items-center justify-center shrink-0 font-black text-[10px] transition-colors"
                      :class="selectedStudentIdsForPromotion.includes(stu.id) ? 'bg-indigo-600 text-white' : 'bg-white text-slate-400'">
                      {{ idx + 1 }}
                    </div>
                    <div class="flex flex-col min-w-0">
                      <span class="text-xs font-black text-slate-800 uppercase tracking-tight truncate">{{ stu.name }}</span>
                      <div class="flex items-center gap-2 mt-0.5">
                        <span class="text-[9px] font-bold text-slate-500 bg-white px-2 py-0.5 rounded border border-slate-200 shadow-sm uppercase tracking-wider">NISN: {{ stu.nisn || '-' }}</span>
                        <span class="text-[9px] font-bold text-indigo-600 bg-indigo-50 px-2 py-0.5 rounded border border-indigo-100 uppercase tracking-wider">Angkatan Masuk: {{ stu.entry_year }}</span>
                      </div>
                    </div>
                  </div>
                  <div class="flex items-center gap-3 shrink-0">
                    <div class="w-5 h-5 rounded-md border-2 flex items-center justify-center transition-colors"
                      :class="selectedStudentIdsForPromotion.includes(stu.id) ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300 bg-white'">
                      <SuccessIcon v-if="selectedStudentIdsForPromotion.includes(stu.id)" class="w-3.5 h-3.5 text-white" />
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="px-8 py-6 bg-slate-50/80 border-t border-slate-100 flex items-center justify-between shrink-0">
            <button @click="$emit('close')" class="text-[10px] font-black text-slate-400 hover:text-slate-600 uppercase tracking-widest transition-all px-6 py-3 hover:bg-slate-100 rounded-xl border border-transparent hover:border-slate-200">Batal</button>
            <button @click="handlePromote" :disabled="isOffline || loading || !promoteForm.source_class_id || !promoteForm.target_class_id || selectedStudentIdsForPromotion.length === 0" :title="isOffline ? 'Server harus online untuk memproses kenaikan/pindah kelas.' : ''" class="bg-indigo-600 hover:bg-indigo-700 text-white font-black px-8 py-4 rounded-2xl shadow-xl shadow-indigo-100 transition-all text-[11px] uppercase tracking-widest flex items-center gap-3 disabled:opacity-50 disabled:cursor-not-allowed">
              <SuccessIcon v-if="!loading" class="w-5 h-5" />
              <span v-else class="w-5 h-5 border-2 border-white/20 border-t-white rounded-full animate-spin"></span>
              <span>{{ isOffline ? 'Online Diperlukan' : `Proses Perpindahan (${selectedStudentIdsForPromotion.length} Siswa)` }}</span>
            </button>
          </div>
        </div>
      </div>
    </transition>

    <!-- Bulk Graduate Modal -->
    <transition enter-active-class="transition duration-300 ease-out" enter-from-class="opacity-0" enter-to-class="opacity-100" leave-active-class="transition duration-200 ease-in" leave-from-class="opacity-100" leave-to-class="opacity-0">
      <div v-if="showGraduate" class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
        <div class="bg-white w-full max-w-2xl rounded-[40px] shadow-2xl overflow-hidden animate-scale-in flex flex-col max-h-[90vh]">
          <div class="px-8 py-8 border-b border-slate-50 flex items-center justify-between bg-blue-50/30 shrink-0">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-2xl bg-blue-600 flex items-center justify-center shadow-lg shadow-blue-100">
                <GraduationCapIcon class="w-6 h-6 text-white" />
              </div>
              <div>
                <h2 class="text-lg font-black text-slate-800 tracking-tight">Kelulusan Masal (Student Wizard)</h2>
                <p class="text-[10px] text-slate-400 font-bold uppercase tracking-widest mt-0.5">Penetapan status alumni berbasis rombel nyata</p>
              </div>
            </div>
            <button @click="$emit('close')" class="w-10 h-10 rounded-xl hover:bg-white text-slate-400 hover:text-slate-600 transition-all shadow-sm border border-transparent">
              <CloseIcon class="w-5 h-5" />
            </button>
          </div>

          <div class="p-8 space-y-6 overflow-y-auto custom-scrollbar bg-slate-50/30 flex-1">
            <div v-if="isOffline" class="p-5 bg-amber-50 border border-amber-200 rounded-3xl flex items-start gap-4 shadow-sm shadow-amber-100/60">
              <div class="w-10 h-10 rounded-xl bg-white flex items-center justify-center shrink-0 shadow-sm text-amber-600">
                <AlertIcon class="w-5 h-5" />
              </div>
              <p class="text-[10px] font-black text-amber-800 leading-relaxed uppercase tracking-tight">
                Server sedang offline. Kelulusan masal dikunci karena sistem harus mengecek data siswa dan tagihan aktif terbaru.
              </p>
            </div>

            <div class="p-5 bg-blue-50/50 border border-blue-100 rounded-3xl flex items-start gap-4 shadow-sm">
              <div class="w-10 h-10 rounded-xl bg-white flex items-center justify-center shrink-0 shadow-sm shadow-blue-100">
                <InfoIcon class="w-5 h-5 text-blue-500" />
              </div>
              <p class="text-[10px] font-bold text-blue-700 leading-relaxed uppercase tracking-tight">
                Pilih rombongan belajar (kelas aktif). Sistem akan mendeteksi seluruh siswa di kelas tersebut (termasuk siswa angkatan lama yang tinggal kelas). Anda dapat menghapus centang pada siswa yang tidak lulus tahun ini.
              </p>
            </div>

            <div class="space-y-3">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1">Pilih Kelas Aktif <span class="text-rose-500">*</span></label>
              <div class="relative group">
                <select v-model="graduateForm.class_id" class="w-full py-4 px-6 bg-white border border-slate-200 rounded-2xl appearance-none focus:ring-4 focus:ring-blue-50 focus:border-blue-500 outline-none transition-all font-bold text-xs text-slate-700 shadow-sm pr-12 cursor-pointer">
                  <option value="">Pilih kelas yang akan diluluskan</option>
                  <option v-for="c in academicFilters.classes" :key="c.id" :value="c.id">{{ c.name }}</option>
                </select>
                <div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400 group-focus-within:text-blue-500 transition-colors">
                  <ChevronDownIcon class="w-4 h-4" />
                </div>
              </div>
            </div>

            <!-- Student Selection List -->
            <div v-if="graduateForm.class_id" class="bg-white rounded-[2rem] border border-slate-200 shadow-sm overflow-hidden space-y-4 p-6">
              <div class="flex items-center justify-between pb-4 border-b border-slate-100">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 bg-blue-50 text-blue-600 rounded-xl flex items-center justify-center font-black text-xs">
                    {{ selectedStudentIdsForGraduation.length }}
                  </div>
                  <div>
                    <h4 class="text-xs font-black text-slate-800 uppercase tracking-wider">Daftar Siswa Kandidat Lulus</h4>
                    <p class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Centang siswa yang berhak lulus</p>
                  </div>
                </div>
                <label class="flex items-center gap-2 text-[10px] font-black text-slate-600 uppercase tracking-widest cursor-pointer hover:text-blue-600 transition-colors bg-slate-50 px-4 py-2 rounded-xl border border-slate-200">
                  <input 
                    type="checkbox" 
                    :checked="selectedStudentIdsForGraduation.length === eligibleStudents.length && eligibleStudents.length > 0"
                    @change="e => toggleSelectAllGraduation(e.target.checked)"
                    class="rounded border-slate-300 text-blue-600 focus:ring-blue-500 w-4 h-4"
                  />
                  <span>Pilih Semua ({{ eligibleStudents.length }})</span>
                </label>
              </div>

              <div v-if="loadingStudents" class="py-12 text-center space-y-3">
                <div class="w-8 h-8 border-4 border-blue-600/20 border-t-blue-600 rounded-full animate-spin mx-auto"></div>
                <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Memuat daftar siswa aktif di kelas...</p>
              </div>

              <div v-else-if="eligibleStudents.length === 0" class="py-12 text-center space-y-2 bg-slate-50/50 rounded-2xl border border-dashed border-slate-200">
                <p class="text-xs font-black text-slate-600 uppercase tracking-wider">Tidak Ada Siswa Aktif</p>
                <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Tidak ditemukan siswa dengan status aktif di rombel terpilih.</p>
              </div>

              <div v-else class="space-y-2.5 max-h-[280px] overflow-y-auto custom-scrollbar pr-2">
                <div v-for="(stu, idx) in eligibleStudents" :key="stu.id"
                  @click="toggleSelectOneGraduation(stu.id)"
                  class="flex items-center justify-between p-4 rounded-2xl border-2 transition-all cursor-pointer group"
                  :class="selectedStudentIdsForGraduation.includes(stu.id) ? 'border-blue-600 bg-blue-50/30' : 'border-slate-100 bg-slate-50 hover:border-blue-200'">
                  <div class="flex items-center gap-4 min-w-0">
                    <div class="w-7 h-7 rounded-lg flex items-center justify-center shrink-0 font-black text-[10px] transition-colors"
                      :class="selectedStudentIdsForGraduation.includes(stu.id) ? 'bg-blue-600 text-white' : 'bg-white text-slate-400'">
                      {{ idx + 1 }}
                    </div>
                    <div class="flex flex-col min-w-0">
                      <span class="text-xs font-black text-slate-800 uppercase tracking-tight truncate">{{ stu.name }}</span>
                      <div class="flex items-center gap-2 mt-0.5">
                        <span class="text-[9px] font-bold text-slate-500 bg-white px-2 py-0.5 rounded border border-slate-200 shadow-sm uppercase tracking-wider">NISN: {{ stu.nisn || '-' }}</span>
                        <span class="text-[9px] font-bold text-blue-600 bg-blue-50 px-2 py-0.5 rounded border border-blue-100 uppercase tracking-wider">Angkatan Masuk: {{ stu.entry_year }}</span>
                      </div>
                    </div>
                  </div>
                  <div class="flex items-center gap-3 shrink-0">
                    <div class="w-5 h-5 rounded-md border-2 flex items-center justify-center transition-colors"
                      :class="selectedStudentIdsForGraduation.includes(stu.id) ? 'bg-blue-600 border-blue-600' : 'border-slate-300 bg-white'">
                      <SuccessIcon v-if="selectedStudentIdsForGraduation.includes(stu.id)" class="w-3.5 h-3.5 text-white" />
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="px-8 py-6 bg-slate-50/80 border-t border-slate-100 flex items-center justify-between shrink-0">
            <button @click="$emit('close')" class="text-[10px] font-black text-slate-400 hover:text-slate-600 uppercase tracking-widest transition-all px-6 py-3 hover:bg-slate-100 rounded-xl border border-transparent hover:border-slate-200">Batal</button>
            <button @click="handleGraduate" :disabled="isOffline || loading || !graduateForm.class_id || selectedStudentIdsForGraduation.length === 0" :title="isOffline ? 'Server harus online untuk memproses kelulusan masal.' : ''" class="bg-blue-600 hover:bg-blue-700 text-white font-black px-8 py-4 rounded-2xl shadow-xl shadow-blue-100 transition-all text-[11px] uppercase tracking-widest flex items-center gap-3 disabled:opacity-50 disabled:cursor-not-allowed">
              <SuccessIcon v-if="!loading" class="w-5 h-5" />
              <span v-else class="w-5 h-5 border-2 border-white/20 border-t-white rounded-full animate-spin"></span>
              <span>{{ isOffline ? 'Online Diperlukan' : `Proses Kelulusan (${selectedStudentIdsForGraduation.length} Siswa)` }}</span>
            </button>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<style scoped>
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
