<script setup>
import { ref, watch, computed } from 'vue'
import axios from 'axios'
import { 
  X as CloseIcon,
  CheckCircle2 as CheckCircleIcon,
  Loader2 as LoaderIcon,
  GraduationCap as MajorIcon,
  Users as ClassIcon,
  AlertTriangle as WarningIcon
} from 'lucide-vue-next'
import FormError from '../ui/FormError.vue'
import academicYearService from '../../services/academic_year.service'

const props = defineProps({
  modelValue: Boolean,
  isEditing: Boolean,
  form: Object,
  errors: Object,
  submitting: Boolean
})

const emit = defineEmits(['update:modelValue', 'save', 'clear-field-error', 'set-field-error'])

const clearFieldError = (field) => {
  emit('clear-field-error', field)
}

const setFieldError = (field, messages) => {
  emit('set-field-error', { field, messages })
}

const checkUniqueYear = async () => {
  if (!props.form.year) {
    setFieldError('year', ['Tahun angkatan wajib diisi'])
    return
  } else if (props.form.year < 2000 || props.form.year > 2100) {
    setFieldError('year', ['Tahun angkatan tidak valid'])
    return
  }
  
  try {
    const excludeId = props.isEditing ? props.form.id : 0
    const res = await academicYearService.checkUnique(props.form.year, excludeId)
    if (res.data?.data?.is_unique === false) {
      setFieldError('year', ['Tahun angkatan sudah terdaftar di sistem'])
    } else {
      clearFieldError('year')
    }
  } catch (err) {
    console.error('Error checking academic year uniqueness:', err)
  }
}

const majors = ref([])
const availableClasses = ref([])
const selectedMajorRecords = ref([])
const selectedClassRecords = ref([])
const loadingMajors = ref(false)
const loadingClasses = ref(false)
const selectedMajorIds = ref([])
const selectedClassIds = ref([])
const warningMessage = ref('')

const normalizeID = (value) => Number(value)

const uniqueByID = (items) => {
  const map = new Map()
  items.filter(Boolean).forEach(item => map.set(normalizeID(item.id), item))
  return Array.from(map.values())
}

const fetchMajors = async () => {
  loadingMajors.value = true
  try {
    const res = await axios.get('academic/major', { params: { status: 'active', limit: 1000 } })
    const allMajors = res.data?.data?.data || []
    majors.value = allMajors.filter(m => m.class_count > 0)
  } catch (err) {
    console.error('Gagal mengambil data jurusan:', err)
  } finally {
    loadingMajors.value = false
  }
}

const fetchAllClasses = async () => {
  loadingClasses.value = true
  try {
    const res = await axios.get('academic/class', { params: { status: 'active', limit: 1000 } })
    availableClasses.value = res.data?.data?.data || []
  } catch (err) {
    console.error('Gagal mengambil data kelas:', err)
  } finally {
    loadingClasses.value = false
  }
}

const fetchSelectedData = async (id) => {
  try {
    if (props.form?.major_ids?.length > 0) {
      selectedMajorIds.value = props.form.major_ids.map(normalizeID)
    }
    if (props.form?.class_ids?.length > 0) {
      selectedClassIds.value = props.form.class_ids.map(normalizeID)
    }

    const [majorsRes, classesRes] = await Promise.all([
      axios.get(`academic/years/${id}/majors`),
      axios.get(`academic/years/${id}/classes`)
    ])

    if (majorsRes.data?.data) {
      selectedMajorRecords.value = majorsRes.data.data
      selectedMajorIds.value = majorsRes.data.data.map(m => normalizeID(m.id))
    }
    if (classesRes.data?.data) {
      selectedClassRecords.value = classesRes.data.data
      selectedClassIds.value = classesRes.data.data.map(c => normalizeID(c.id))
    }
  } catch (err) {
    console.error('Gagal mengambil data terpilih:', err)
  }
}

watch(() => props.modelValue, (val) => {
  if (val) {
    warningMessage.value = ''
    selectedMajorRecords.value = []
    selectedClassRecords.value = []
    fetchMajors()
    fetchAllClasses()
    if (props.isEditing && props.form.id) {
      fetchSelectedData(props.form.id)
    } else {
      selectedMajorIds.value = []
      selectedClassIds.value = []
    }
  }
})

const majorOptions = computed(() => uniqueByID([...majors.value, ...selectedMajorRecords.value]))
const allClassOptions = computed(() => uniqueByID([...availableClasses.value, ...selectedClassRecords.value]))

const isMajorInvalid = (major) => {
  const existsAsActiveOption = majors.value.some(m => normalizeID(m.id) === normalizeID(major.id))
  return !major?.is_active || !existsAsActiveOption
}

const isClassInvalid = (cls) => {
  const existsAsActiveOption = availableClasses.value.some(c => normalizeID(c.id) === normalizeID(cls.id))
  return !cls?.is_active || !existsAsActiveOption || !selectedMajorIds.value.includes(normalizeID(cls.major_id))
}

const invalidSelectedMajors = computed(() => majorOptions.value.filter(m => selectedMajorIds.value.includes(normalizeID(m.id)) && isMajorInvalid(m)))
const invalidSelectedClasses = computed(() => allClassOptions.value.filter(c => selectedClassIds.value.includes(normalizeID(c.id)) && isClassInvalid(c)))
const hasInvalidSelections = computed(() => invalidSelectedMajors.value.length > 0 || invalidSelectedClasses.value.length > 0)

// Filter classes based on selected majors, while keeping selected invalid classes visible for review.
const filteredClasses = computed(() => {
  if (selectedMajorIds.value.length === 0) return []
  return allClassOptions.value.filter(c => {
    const id = normalizeID(c.id)
    return selectedMajorIds.value.includes(normalizeID(c.major_id)) || selectedClassIds.value.includes(id)
  })
})

// Group classes by major for better UI
const groupedClasses = computed(() => {
  const groups = {}
  filteredClasses.value.forEach(c => {
    const groupName = c.major_name || 'Perlu Ditinjau'
    if (!groups[groupName]) groups[groupName] = []
    groups[groupName].push(c)
  })
  return groups
})

const clearError = (field) => {
  clearFieldError(field)
}

const handleSave = () => {
  let hasError = false
  warningMessage.value = ''
  if (!props.form.year) {
    setFieldError('year', ['Tahun angkatan wajib diisi'])
    hasError = true
  } else if (props.form.year < 2000 || props.form.year > 2100) {
    setFieldError('year', ['Tahun angkatan tidak valid'])
    hasError = true
  }

  if (selectedMajorIds.value.length === 0) {
    setFieldError('major_ids', ['Pilih minimal satu jurusan'])
    hasError = true
  }

  if (selectedClassIds.value.length === 0) {
    setFieldError('class_ids', ['Pilih minimal satu kelas'])
    hasError = true
  }

  if (props.errors.year || props.errors.major_ids || props.errors.class_ids) {
    hasError = true
  }

  if (hasInvalidSelections.value) {
    warningMessage.value = 'Ada jurusan atau kelas lama yang sudah nonaktif/tidak tersedia. Klik item kuning untuk melepasnya, lalu pilih data aktif sebagai pengganti.'
    hasError = true
  }

  if (!hasError) {
    emit('save', {
      year: parseInt(props.form.year),
      is_active: props.form.is_active || false,
      major_ids: selectedMajorIds.value,
      class_ids: selectedClassIds.value
    })
  }
}

const toggleMajor = (id) => {
  id = normalizeID(id)
  const major = majorOptions.value.find(m => normalizeID(m.id) === id)
  const idx = selectedMajorIds.value.indexOf(id)

  if (major && isMajorInvalid(major)) {
    if (idx > -1) {
      selectedMajorIds.value.splice(idx, 1)
      const majorClasses = allClassOptions.value.filter(c => normalizeID(c.major_id) === id).map(c => normalizeID(c.id))
      selectedClassIds.value = selectedClassIds.value.filter(cid => !majorClasses.includes(normalizeID(cid)))
      warningMessage.value = `Jurusan ${major.name} dilepas karena sudah nonaktif atau tidak tersedia.`
      clearError('major_ids')
      clearError('class_ids')
    } else {
      warningMessage.value = `Jurusan ${major.name} tidak bisa dipilih karena sudah nonaktif atau tidak tersedia.`
    }
    return
  }

  warningMessage.value = ''
  clearError('major_ids')
  if (idx > -1) {
    selectedMajorIds.value.splice(idx, 1)
    const majorClasses = allClassOptions.value.filter(c => normalizeID(c.major_id) === id).map(c => normalizeID(c.id))
    selectedClassIds.value = selectedClassIds.value.filter(cid => !majorClasses.includes(normalizeID(cid)))
  } else {
    selectedMajorIds.value.push(id)
  }
}

const toggleClass = (id) => {
  id = normalizeID(id)
  const cls = allClassOptions.value.find(c => normalizeID(c.id) === id)
  const idx = selectedClassIds.value.indexOf(id)

  if (cls && isClassInvalid(cls)) {
    if (idx > -1) {
      selectedClassIds.value.splice(idx, 1)
      warningMessage.value = `Kelas ${cls.name} dilepas karena sudah nonaktif atau tidak sesuai dengan jurusan aktif yang dipilih.`
      clearError('class_ids')
    } else {
      warningMessage.value = `Kelas ${cls.name} tidak bisa dipilih karena sudah nonaktif atau tidak sesuai dengan jurusan aktif.`
    }
    return
  }

  warningMessage.value = ''
  clearError('class_ids')
  if (idx > -1) {
    selectedClassIds.value.splice(idx, 1)
  } else {
    selectedClassIds.value.push(id)
  }
}
</script>

<template>
  <Teleport to="body">
    <transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0 scale-95"
      enter-to-class="opacity-100 scale-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100 scale-100"
      leave-to-class="opacity-0 scale-95"
    >
      <div v-if="modelValue" class="fixed inset-0 z-[1000] flex items-center justify-center p-4 overflow-y-auto">
        <div class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm" @click="emit('update:modelValue', false)"></div>
        <div class="relative bg-white w-full max-w-2xl rounded-[2.5rem] shadow-2xl overflow-hidden animate-scale-in flex flex-col max-h-[90vh] my-8">
          <div class="px-8 py-6 border-b border-slate-50 flex items-center justify-between bg-white shrink-0">
            <div>
              <h3 class="font-black text-slate-800 text-xl tracking-tight">
                {{ isEditing ? 'Ubah Angkatan' : 'Tambah Angkatan' }}
              </h3>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">Lengkapi informasi angkatan siswa</p>
            </div>
            <button @click="emit('update:modelValue', false)" class="p-3 hover:bg-slate-50 text-slate-400 hover:text-slate-600 rounded-2xl transition-all">
              <CloseIcon class="w-6 h-6" />
            </button>
          </div>

          <div class="p-8 space-y-8 overflow-y-auto custom-scrollbar bg-slate-50/20">
            <!-- Info Banner -->
            <div class="p-3.5 bg-amber-50 border border-amber-200 rounded-2xl text-amber-700 text-[9px] font-bold uppercase tracking-widest flex items-start gap-2 shadow-sm">
              <WarningIcon class="w-4 h-4 shrink-0 mt-0.5" />
              <span>Hanya data aktif yang bisa disimpan. Relasi lama yang nonaktif akan tampil kuning dan perlu dilepas atau diganti.</span>
            </div>
            <div v-if="warningMessage || hasInvalidSelections" class="p-3.5 bg-amber-50 border border-amber-200 rounded-2xl text-amber-700 text-[10px] font-bold leading-relaxed flex items-start gap-2 shadow-sm">
              <WarningIcon class="w-4 h-4 shrink-0 mt-0.5" />
              <span>{{ warningMessage || 'Ada jurusan atau kelas lama yang sudah nonaktif/tidak tersedia. Klik item kuning untuk melepasnya, lalu pilih data aktif sebagai pengganti.' }}</span>
            </div>
            <FormError :message="errors._general" />

            <!-- Tahun -->
            <div class="bg-white p-6 rounded-[2rem] border border-slate-100 shadow-sm space-y-4">
              <div class="flex items-center gap-3 mb-2">
                <div class="w-8 h-8 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center">
                  <CheckCircleIcon class="w-4 h-4" />
                </div>
                <h4 class="text-[11px] font-black text-slate-700 uppercase tracking-widest">Informasi Dasar</h4>
              </div>
              <div class="space-y-1.5">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1">
                  Tahun Angkatan <span class="text-rose-500">*</span>
                </label>
                <input 
                  v-model="form.year"
                  type="number" 
                  placeholder="Contoh: 2024" 
                  class="w-full py-4 px-6 bg-slate-50 border border-slate-100 rounded-2xl focus:bg-white focus:ring-4 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-sm text-slate-700 shadow-sm"
                  :class="{'!border-rose-500 !ring-rose-50': errors.year}"
                  @input="clearFieldError('year')"
                  @blur="checkUniqueYear"
                />
                <FormError :message="errors.year" />
              </div>
            </div>

            <!-- Jurusan -->
            <div class="bg-white p-6 rounded-[2rem] border border-slate-100 shadow-sm space-y-4">
              <div class="flex items-center justify-between mb-2">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 bg-emerald-50 text-emerald-600 rounded-xl flex items-center justify-center shrink-0">
                    <MajorIcon class="w-4 h-4" />
                  </div>
                  <div>
                    <h4 class="text-[11px] font-black text-slate-700 uppercase tracking-widest">Pilih Jurusan Aktif</h4>
                  </div>
                </div>
                <span class="text-[9px] font-bold text-emerald-600 bg-emerald-50 px-3 py-1 rounded-full uppercase tracking-tighter shadow-sm border border-emerald-100">
                  {{ selectedMajorIds.length }} Terpilih<span v-if="invalidSelectedMajors.length"> / {{ invalidSelectedMajors.length }} Perlu Ditinjau</span>
                </span>
              </div>
              
              <div v-if="loadingMajors" class="py-10 text-center bg-slate-50/50 rounded-2xl border border-dashed border-slate-200">
                <LoaderIcon class="w-6 h-6 animate-spin text-indigo-500 mx-auto mb-2" />
                <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Memuat Jurusan...</p>
              </div>
              
              <div v-else class="grid grid-cols-2 gap-3">
                <div v-for="(major, idx) in majorOptions" :key="major.id" 
                  @click="toggleMajor(major.id)"
                  class="p-4 rounded-2xl border-2 transition-all cursor-pointer flex items-center gap-3 group relative overflow-hidden"
                  :title="isMajorInvalid(major) ? 'Jurusan ini sudah nonaktif/tidak tersedia. Klik untuk melepas pilihan.' : ''"
                  :class="isMajorInvalid(major)
                    ? 'border-amber-300 bg-amber-50/80 hover:border-amber-400'
                    : selectedMajorIds.includes(major.id) 
                      ? 'border-emerald-500 bg-emerald-50/50' 
                      : 'border-slate-100 bg-slate-50 hover:border-emerald-200'">
                  <div class="w-8 h-8 rounded-lg flex items-center justify-center shrink-0 transition-all shadow-sm"
                    :class="isMajorInvalid(major)
                      ? 'bg-amber-500 text-white'
                      : selectedMajorIds.includes(major.id) ? 'bg-emerald-600 text-white' : 'bg-white text-slate-400 group-hover:text-emerald-500'">
                    <WarningIcon v-if="isMajorInvalid(major)" class="w-4 h-4" />
                    <span v-else class="text-[10px] font-black">{{ idx + 1 }}</span>
                  </div>
                  <div class="flex flex-col min-w-0">
                    <span class="text-[10px] font-black uppercase tracking-tight truncate" :class="isMajorInvalid(major) ? 'text-amber-900' : 'text-slate-700'">{{ major.name }}</span>
                    <span class="text-[8px] font-bold uppercase tracking-widest truncate" :class="isMajorInvalid(major) ? 'text-amber-700' : 'text-slate-400'">{{ major.code }}<span v-if="isMajorInvalid(major)"> - nonaktif/tidak tersedia</span></span>
                  </div>
                  <div class="ml-auto">
                    <div class="w-4 h-4 rounded-full border-2 flex items-center justify-center transition-all"
                      :class="isMajorInvalid(major)
                        ? 'bg-amber-500 border-amber-500'
                        : selectedMajorIds.includes(major.id) ? 'bg-emerald-600 border-emerald-600' : 'border-slate-300'">
                      <WarningIcon v-if="isMajorInvalid(major)" class="w-2.5 h-2.5 text-white" />
                      <CheckCircleIcon v-else-if="selectedMajorIds.includes(major.id)" class="w-2.5 h-2.5 text-white" />
                    </div>
                  </div>
                </div>
              </div>
              <FormError :message="errors.major_ids" />
            </div>

            <!-- Kelas (Checkbox per Jurusan) -->
            <div v-if="selectedMajorIds.length > 0" class="bg-white p-6 rounded-[2rem] border border-slate-100 shadow-sm space-y-6">
              <div class="flex items-center justify-between mb-2">
                <div class="flex items-center gap-3">
                  <div class="w-8 h-8 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center">
                    <ClassIcon class="w-4 h-4" />
                  </div>
                  <h4 class="text-[11px] font-black text-slate-700 uppercase tracking-widest">Pilih Kelas Aktif</h4>
                </div>
                <span class="text-[9px] font-bold text-indigo-500 bg-indigo-50 px-3 py-1 rounded-full uppercase tracking-tighter shadow-sm border border-indigo-100">
                  {{ selectedClassIds.length }} Terpilih<span v-if="invalidSelectedClasses.length"> / {{ invalidSelectedClasses.length }} Perlu Ditinjau</span>
                </span>
              </div>

              <div v-if="loadingClasses" class="py-10 text-center bg-slate-50/50 rounded-2xl border border-dashed border-slate-200">
                <LoaderIcon class="w-6 h-6 animate-spin text-indigo-500 mx-auto mb-2" />
                <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest">Memuat Kelas...</p>
              </div>

              <div v-else-if="Object.keys(groupedClasses).length === 0" class="py-10 text-center bg-slate-50/50 rounded-2xl border border-dashed border-slate-200">
                <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest leading-relaxed px-10">Belum ada kelas yang terdaftar untuk jurusan terpilih. Silakan buat kelas baru di menu Manajemen Kelas.</p>
              </div>

              <div v-else class="space-y-6">
                <div v-for="(classes, majorName) in groupedClasses" :key="majorName" class="space-y-3">
                  <div class="flex items-center gap-2 px-2">
                    <div class="w-1.5 h-1.5 rounded-full bg-indigo-400"></div>
                    <span class="text-[9px] font-black text-slate-500 uppercase tracking-widest">{{ majorName }}</span>
                  </div>
                  <div class="grid grid-cols-2 gap-3">
                    <div v-for="(cls, idx) in classes" :key="cls.id" 
                      @click="toggleClass(cls.id)"
                      class="p-4 rounded-2xl border-2 transition-all cursor-pointer flex items-center gap-3 group"
                      :title="isClassInvalid(cls) ? 'Kelas ini sudah nonaktif/tidak sesuai. Klik untuk melepas pilihan.' : ''"
                      :class="isClassInvalid(cls)
                        ? 'border-amber-300 bg-amber-50/80 hover:border-amber-400'
                        : selectedClassIds.includes(cls.id) 
                          ? 'border-indigo-600 bg-indigo-50/50' 
                          : 'border-slate-100 hover:border-indigo-200 bg-slate-50/30'">
                      <div class="w-8 h-8 rounded-lg flex items-center justify-center shrink-0 transition-all shadow-sm"
                        :class="isClassInvalid(cls)
                          ? 'bg-amber-500 text-white'
                          : selectedClassIds.includes(cls.id) ? 'bg-indigo-600 text-white' : 'bg-white text-slate-400 group-hover:text-indigo-500'">
                        <WarningIcon v-if="isClassInvalid(cls)" class="w-4 h-4" />
                        <span v-else class="text-[10px] font-black">{{ idx + 1 }}</span>
                      </div>
                      <div class="flex flex-col min-w-0">
                        <span class="text-[10px] font-black uppercase tracking-tight truncate" :class="isClassInvalid(cls) ? 'text-amber-900' : 'text-slate-700'">{{ cls.name }}</span>
                        <span v-if="isClassInvalid(cls)" class="text-[8px] font-bold text-amber-700 uppercase tracking-widest truncate">Nonaktif/tidak sesuai pilihan</span>
                      </div>
                      <div class="ml-auto">
                        <div class="w-4 h-4 rounded border-2 flex items-center justify-center transition-all"
                          :class="isClassInvalid(cls)
                            ? 'bg-amber-500 border-amber-500'
                            : selectedClassIds.includes(cls.id) ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300'">
                          <WarningIcon v-if="isClassInvalid(cls)" class="w-2.5 h-2.5 text-white" />
                          <CheckCircleIcon v-else-if="selectedClassIds.includes(cls.id)" class="w-2.5 h-2.5 text-white" />
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <FormError :message="errors.class_ids" />
            </div>
          </div>

          <div class="p-8 bg-slate-50/50 border-t border-slate-50 flex items-center gap-3 shrink-0">
            <button @click="emit('update:modelValue', false)" class="flex-1 py-4 bg-white hover:bg-slate-100 text-slate-600 font-black rounded-2xl transition-all text-[10px] uppercase tracking-widest border border-slate-200 shadow-sm">
              Batal
            </button>
            <button @click="handleSave" :disabled="submitting" class="flex-1 py-4 bg-indigo-600 hover:bg-indigo-700 text-white font-black rounded-2xl transition-all text-[10px] uppercase tracking-widest shadow-xl shadow-indigo-100 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-3">
              <CheckCircleIcon v-if="!submitting" class="w-4 h-4" />
              <LoaderIcon v-else class="w-4 h-4 animate-spin" />
              <span>{{ submitting ? 'Menyimpan...' : 'Simpan Data' }}</span>
            </button>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<style scoped lang="postcss">
.animate-scale-in { animation: scaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1); }
@keyframes scaleIn { from { opacity: 0; transform: scale(0.9); } to { opacity: 1; transform: scale(1); } }
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #cbd5e1; }
</style>
