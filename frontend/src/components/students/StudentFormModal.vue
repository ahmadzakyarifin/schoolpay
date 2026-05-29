<script setup>
import { ref, computed, watch } from 'vue'
import { 
  X as CloseIcon, 
  User as UserIcon, 
  Calendar as CalendarAltIcon,
  ChevronDown as ChevronDownIcon,
  Mail as MailIcon,
  MessageCircle as WAIcon,
  Check as CheckIcon,
  Search as SearchIcon,
  CheckCircle2 as CheckCircleIcon,
  Plus as PlusIcon
} from 'lucide-vue-next'
import FormError from '../ui/FormError.vue'
import studentService from '../../services/student.service'

const props = defineProps({
  modelValue: Boolean,
  isEditing: Boolean,
  form: Object,
  errors: Object,
  submitting: Boolean,
  academicFilters: Object,
  activeParents: Array,
  photoPreview: String,
  staticBase: String,
  birthDateDisplay: String
})

const emit = defineEmits(['update:modelValue', 'update:birthDateDisplay', 'save', 'fetch-regions', 'local-validation-failed', 'photo-upload', 'clear-field-error', 'set-field-error'])

// Local State
const birthDateRef = ref(null)
const showCountryDropdown = ref(false)
const parentSearchQuery = ref('')

const countries = [
  { name: 'Indonesia', code: '62', flag: '🇮🇩' },
  { name: 'Malaysia', code: '60', flag: '🇲🇾' },
  { name: 'Singapore', code: '65', flag: '🇸🇬' },
  { name: 'Brunei', code: '673', flag: '🇧🇳' }
]

const selectedCountry = computed(() => countries.find(c => c.code === props.form.country_code) || countries[0])

const filteredParents = computed(() => {
  if (!parentSearchQuery.value) return props.activeParents
  const q = parentSearchQuery.value.toLowerCase()
  return props.activeParents.filter(p => 
    p.name.toLowerCase().includes(q) || 
    (p.email && p.email.toLowerCase().includes(q)) ||
    (p.phone_number && p.phone_number.includes(q))
  )
})

const selectedParent = computed(() => props.activeParents.find(p => p.id === props.form.parent_id))

const selectParent = (parent) => {
  props.form.parent_id = parent.id
  clearFieldError('parent_id')
}

const setFieldError = (field, messages) => {
  emit('set-field-error', { field, messages })
}

const validateRequiredField = (field, label) => {
  const value = props.form[field]
  if (value === null || value === undefined || value.toString().trim() === '') {
    setFieldError(field, [`${label} wajib diisi`])
    return false
  }
  clearFieldError(field)
  return true
}

const validateNameField = () => {
  const value = props.form.name?.trim() || ''
  if (!value) {
    setFieldError('name', ['Nama lengkap wajib diisi'])
    return false
  }
  if (value.length < 2) {
    setFieldError('name', ['Nama minimal berisi 2 karakter'])
    return false
  }
  clearFieldError('name')
  return true
}

const getPhoneForUniqueCheck = (cleanPhone) => {
  const countryCode = props.form.country_code || '62'
  return `${countryCode}${cleanPhone}`
}

const triggerCalendar = () => birthDateRef.value?.showPicker()

const handleFileChange = (e) => {
  const file = e.target.files[0]
  if (file) emit('photo-upload', file)
}

// Format HH/BB/TTTT manually while typing
const localBirthDateDisplay = computed({
  get: () => props.birthDateDisplay,
  set: (val) => {
    let cleaned = val.replace(/\D/g, '').substring(0, 8)
    let formatted = cleaned
    if (cleaned.length > 2) formatted = cleaned.substring(0, 2) + '/' + cleaned.substring(2)
    if (cleaned.length > 4) formatted = formatted.substring(0, 5) + '/' + formatted.substring(5)
    emit('update:birthDateDisplay', formatted)
  }
})

const parseBirthDateDisplay = (value) => {
  const raw = (value || '').trim()
  if (!raw) return { error: 'Tanggal lahir wajib diisi' }
  if (!/^\d{2}\/\d{2}\/\d{4}$/.test(raw)) {
    return { error: 'Format tanggal lahir harus HH/BB/TTTT' }
  }

  const [dayText, monthText, yearText] = raw.split('/')
  const day = Number(dayText)
  const month = Number(monthText)
  const year = Number(yearText)

  if (year < 1900) return { error: 'Tahun lahir minimal 1900' }

  const parsed = new Date(Date.UTC(year, month - 1, day))
  const isRealDate = parsed.getUTCFullYear() === year &&
    parsed.getUTCMonth() === month - 1 &&
    parsed.getUTCDate() === day

  if (!isRealDate) return { error: 'Tanggal lahir tidak sesuai kalender nyata' }

  const today = new Date()
  const todayUtc = new Date(Date.UTC(today.getFullYear(), today.getMonth(), today.getDate()))
  if (parsed > todayUtc) return { error: 'Tanggal lahir tidak boleh di masa depan' }

  return { iso: `${yearText}-${monthText}-${dayText}` }
}

const onBirthDateInput = () => {
  const value = localBirthDateDisplay.value || ''
  if (!value) {
    emit('clear-field-error', 'birth_date')
    return
  }
  if (value.length < 10) {
    props.form.birth_date = ''
    return
  }

  const result = parseBirthDateDisplay(value)
  if (result.error) {
    props.form.birth_date = ''
    emit('set-field-error', { field: 'birth_date', messages: [result.error] })
    return
  }

  props.form.birth_date = result.iso
  emit('clear-field-error', 'birth_date')
}

const onBirthDateBlur = () => {
  const value = localBirthDateDisplay.value || ''
  if (!value) {
    setFieldError('birth_date', ['Tanggal lahir wajib diisi'])
    return
  }
  const result = parseBirthDateDisplay(value)
  if (result.error) {
    emit('set-field-error', { field: 'birth_date', messages: [result.error] })
  } else {
    emit('clear-field-error', 'birth_date')
  }
}

const checkFieldUnique = async (field) => {
  let value = props.form[field]
  if (value === undefined || value === null) return

  const cleanValue = value.toString().trim()
  const requiredLabels = {
    nik: 'NIK',
    nisn: 'NISN',
    email: 'Email siswa',
    phone_number: 'Nomor WhatsApp'
  }
  if (!cleanValue) {
    if (requiredLabels[field]) {
      setFieldError(field, [`${requiredLabels[field]} wajib diisi`])
    } else {
      clearFieldError(field)
    }
    return
  }

  let valueForCheck = cleanValue
  if (field === 'nik') {
    if (cleanValue.length !== 16) {
      setFieldError(field, ['NIK harus berjumlah 16 digit'])
      return
    }
  } else if (field === 'nisn') {
    if (cleanValue.length !== 10) {
      setFieldError(field, ['NISN harus berjumlah 10 digit'])
      return
    }
  } else if (field === 'nis') {
    if (cleanValue.length < 3 || cleanValue.length > 20) {
      setFieldError(field, ['NIS minimal 3 dan maksimal 20 karakter'])
      return
    }
  } else if (field === 'email') {
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(cleanValue)) {
      setFieldError(field, ['Format email tidak valid'])
      return
    }
    valueForCheck = cleanValue.toLowerCase()
  } else if (field === 'phone_number') {
    let cleanPhone = cleanValue
    if (cleanPhone.startsWith('0')) cleanPhone = cleanPhone.replace(/^0+/, '')
    if (cleanPhone.length < 9) {
      setFieldError(field, ['Nomor WhatsApp minimal 9 digit'])
      return
    }
    if (cleanPhone.length > 15) {
      setFieldError(field, ['Nomor WhatsApp maksimal 15 digit'])
      return
    }
    props.form.phone_number = cleanPhone
    valueForCheck = getPhoneForUniqueCheck(cleanPhone)
  }

  try {
    const excludeId = props.isEditing ? props.form.id : 0
    const res = await studentService.checkUnique(field, valueForCheck, excludeId)
    if (res.data?.data?.is_unique === false) {
      let fieldLabel = field.toUpperCase()
      if (field === 'phone_number') fieldLabel = 'Nomor WhatsApp'
      setFieldError(field, [`${fieldLabel} sudah terdaftar di sistem`])
    } else {
      clearFieldError(field)
    }
  } catch (err) {
    console.error(`Error checking uniqueness for ${field}:`, err)
  }
}

const onProvinceChange = () => {
  props.form.city = ''
  props.form.district = ''
  props.form.village = ''
  clearFieldError('province')
  emit('fetch-regions', 'province')
}

const onCityChange = () => {
  props.form.district = ''
  props.form.village = ''
  clearFieldError('city')
  emit('fetch-regions', 'city')
}

const onDistrictChange = () => {
  props.form.village = ''
  clearFieldError('district')
  emit('fetch-regions', 'district')
}

const filteredMajors = computed(() => {
  if (!props.form.entry_year) return []
  const entryYear = parseInt(props.form.entry_year)
  const ay = props.academicFilters.years?.find(y => y.year === entryYear)
  if (!ay) return []
  
  // Filter majors that are linked to this AcademicYear ID
  return props.academicFilters.majors?.filter(m => 
    m.year_ids?.includes(ay.id)
  ) || []
})

const filteredClasses = computed(() => {
  if (!props.form.entry_year || !props.form.major_id) return []
  const entryYear = parseInt(props.form.entry_year)
  const ay = props.academicFilters.years?.find(y => y.year === entryYear)
  if (!ay) return []

  return props.academicFilters.classes?.filter(c => 
    c.academic_year_ids?.includes(ay.id) && 
    c.major_id === props.form.major_id
  ) || []
})

const onEntryYearChange = () => {
  props.form.major_id = null
  props.form.class_id = null
  clearFieldError('entry_year')
}

const onMajorChange = () => {
  props.form.class_id = null
  clearFieldError('major_id')
}

const onClassChange = () => {
  clearFieldError('class_id')
}

const validateAndSave = () => {
  const localErrors = {}
  
  // 1. Basic Info
  if (!props.form.name?.trim()) localErrors.name = ['Nama lengkap wajib diisi']
  else if (props.form.name.trim().length < 2) localErrors.name = ['Nama minimal berisi 2 karakter']

  if (!props.form.gender) localErrors.gender = ['Jenis kelamin wajib dipilih']
  if (!props.form.religion) localErrors.religion = ['Agama wajib dipilih']
  if (!props.form.birth_place?.trim()) localErrors.birth_place = ['Tempat lahir wajib diisi']

  // 2. Identity
  const cleanNIK = props.form.nik?.toString().trim() || ''
  if (!cleanNIK) localErrors.nik = ['NIK wajib diisi']
  else if (cleanNIK.length !== 16) localErrors.nik = ['NIK harus berjumlah 16 digit']

  const cleanNISN = props.form.nisn?.toString().trim() || ''
  if (!cleanNISN) localErrors.nisn = ['NISN wajib diisi']
  else if (cleanNISN.length !== 10) localErrors.nisn = ['NISN harus berjumlah 10 digit']

  if (props.form.nis) {
    const cleanNIS = props.form.nis.toString().trim()
    if (cleanNIS.length > 0 && (cleanNIS.length < 3 || cleanNIS.length > 20)) {
      localErrors.nis = ['NIS minimal 3 dan maksimal 20 karakter']
    }
  }

  // 3. Regional Data
  if (!props.form.province) localErrors.province = ['Provinsi wajib dipilih']
  if (!props.form.city) localErrors.city = ['Kota/Kabupaten wajib dipilih']
  if (!props.form.district) localErrors.district = ['Kecamatan wajib dipilih']
  if (!props.form.village) localErrors.village = ['Desa/Kelurahan wajib dipilih']

  if (!props.form.rt?.trim()) localErrors.rt = ['RT wajib']
  else if (props.form.rt.length > 5) localErrors.rt = ['Maks 5 digit']
  
  if (!props.form.rw?.trim()) localErrors.rw = ['RW wajib']
  else if (props.form.rw.length > 5) localErrors.rw = ['Maks 5 digit']

  // 4. Contact Info
  if (!props.form.email?.trim()) localErrors.email = ['Email wajib diisi']
  else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(props.form.email)) localErrors.email = ['Format email tidak valid']

  let cleanPhone = props.form.phone_number?.toString().trim() || ''
  if (cleanPhone.startsWith('0')) cleanPhone = cleanPhone.replace(/^0+/, '')
  if (!cleanPhone) localErrors.phone_number = ['Nomor WhatsApp wajib diisi']
  else if (cleanPhone.length < 9) localErrors.phone_number = ['Nomor WhatsApp minimal 9 digit']
  else if (cleanPhone.length > 15) localErrors.phone_number = ['Nomor WhatsApp maksimal 15 digit']

  // 5. Academic Info
  if (!props.form.entry_year) localErrors.entry_year = ['Angkatan wajib dipilih']
  if (!props.form.class_id) localErrors.class_id = ['Kelas wajib dipilih']
  
  if (!props.form.major_id) {
    localErrors.major_id = ['Jurusan wajib dipilih']
  }

  if (!props.form.parent_id) localErrors.parent_id = ['Wali murid wajib dipilih']

  // 6. Birth Date Logic
  const birthDateResult = parseBirthDateDisplay(localBirthDateDisplay.value)
  if (birthDateResult.error) {
    localErrors.birth_date = [birthDateResult.error]
  } else {
    props.form.birth_date = birthDateResult.iso
  }

  if (Object.keys(localErrors).length > 0) {
    emit('local-validation-failed', localErrors)
    return
  }

  const hasBlockingError = Object.entries(props.errors || {}).some(([field, value]) => {
    if (field === '_general') return false
    return Array.isArray(value) ? value.some(Boolean) : !!value
  })
  if (hasBlockingError) return

  emit('save')
}

// LIVE VALIDATION: Clear or trigger "required" errors in real-time
const clearFieldError = (field) => {
  emit('clear-field-error', field)
}

// Remove aggressive watch that triggers required messages

watch(() => props.birthDateDisplay, (val) => {
  if (!val || val.length < 10) return
  const result = parseBirthDateDisplay(val)
  if (!result.error) clearFieldError('birth_date')
})

watch(() => props.form.birth_date, (val) => {
  if (val && val.includes('-')) {
    const [datePart] = val.split('T')
    const [y, m, d] = datePart.split('-')
    if (y && m && d) {
      emit('update:birthDateDisplay', `${d}/${m}/${y}`)
      clearFieldError('birth_date')
    }
  }
})
</script>

<template>
  <Teleport to="body">
    <transition name="page">
    <div v-if="modelValue" class="fixed inset-0 z-[500] flex items-center justify-center p-6">
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="emit('update:modelValue', false)"></div>
      <div class="white-card w-full max-w-2xl relative z-10 overflow-hidden shadow-[0_20px_50px_rgba(0,0,0,0.2)] animate-scale-in !rounded-[2.5rem]">
        <div class="p-8 border-b border-slate-100 flex items-center justify-between bg-white">
          <div>
            <h2 class="text-xl font-black text-slate-900 tracking-tight">{{ isEditing ? 'Edit Data Siswa' : 'Tambah Siswa Baru' }}</h2>
            <p class="text-slate-500 font-medium text-xs mt-1">Pastikan seluruh data identitas diisi dengan benar</p>
          </div>
          <button @click="emit('update:modelValue', false)" class="p-2.5 hover:bg-slate-50 rounded-2xl text-slate-400 transition-all">
            <CloseIcon class="w-5 h-5" />
          </button>
        </div>

        <div class="p-8 space-y-8 max-h-[75vh] overflow-y-auto custom-scrollbar">
          <!-- General Error fallback -->
          <FormError :message="errors._general" />

          <!-- Photo Upload Section -->
          <div class="flex flex-col items-center justify-center p-6 bg-slate-50/50 rounded-[2.5rem] border-2 border-dashed border-slate-200 hover:border-indigo-400 transition-all group relative overflow-hidden">
            <input type="file" @change="handleFileChange" accept="image/*" class="absolute inset-0 opacity-0 cursor-pointer z-10" />
            
            <div v-if="photoPreview || form.image_path" class="relative group/img">
              <img :src="photoPreview || (form.image_path ? `${staticBase}/${form.image_path}` : '')" 
                   class="w-32 h-32 rounded-[2rem] object-cover border-4 border-white shadow-2xl transition-transform group-hover/img:scale-105" />
              <div class="absolute inset-0 bg-black/40 rounded-[2rem] flex items-center justify-center opacity-0 group-hover/img:opacity-100 transition-opacity">
                <span class="text-[10px] font-black text-white uppercase tracking-widest">Ganti Foto</span>
              </div>
            </div>
            <div v-else class="flex flex-col items-center gap-3">
              <div class="w-16 h-16 bg-white rounded-2xl flex items-center justify-center shadow-sm group-hover:scale-110 transition-transform">
                <PlusIcon class="w-6 h-6 text-slate-300 group-hover:text-indigo-500" />
              </div>
              <div class="text-center">
                <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Unggah Foto Siswa</p>
                <p class="text-[9px] text-slate-300 mt-1 uppercase font-bold">JPG, PNG (Maks 2MB)</p>
              </div>
            </div>
          </div>

          <!-- Basic Info -->
          <div class="space-y-6">
            <div class="flex items-center gap-3">
              <div class="w-1.5 h-4 bg-indigo-500 rounded-full"></div>
              <h4 class="text-[11px] font-black text-slate-700 uppercase tracking-widest">Informasi Dasar</h4>
            </div>

            <div class="space-y-2">
              <label class="label-tiny">Nama Lengkap <span class="text-red-500">*</span></label>
              <div class="relative">
                <UserIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300" />
                <input v-model="form.name" @input="clearFieldError('name')" @blur="validateNameField" class="modern-input !pl-12 !rounded-2xl" :class="{'!border-red-500': errors.name}" placeholder="Contoh: Andi" />
              </div>
              <FormError :message="errors.name" />
            </div>

            <div class="grid grid-cols-2 gap-6">
              <div class="space-y-2">
                <label class="label-tiny">Jenis Kelamin <span class="text-red-500">*</span></label>
                <div class="relative">
                  <select v-model="form.gender" @change="clearFieldError('gender')" class="modern-input !rounded-2xl appearance-none !text-[11px]" :class="{'!border-red-500': errors.gender}">
                    <option value="">Pilih Jenis Kelamin</option>
                    <option value="Laki-laki">Laki-laki</option>
                    <option value="Perempuan">Perempuan</option>
                  </select>
                  <ChevronDownIcon class="absolute right-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                </div>
                <FormError :message="errors.gender" />
              </div>
              <div class="space-y-2">
                <label class="label-tiny">Agama <span class="text-red-500">*</span></label>
                <div class="relative">
                  <select v-model="form.religion" @change="clearFieldError('religion')" class="modern-input !rounded-2xl appearance-none" :class="{'!border-red-500': errors.religion}">
                    <option value="Islam">Islam</option>
                    <option value="Kristen">Kristen</option>
                    <option value="Katolik">Katolik</option>
                    <option value="Hindu">Hindu</option>
                    <option value="Budha">Budha</option>
                    <option value="Lainnya">Lainnya</option>
                  </select>
                  <ChevronDownIcon class="absolute right-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                </div>
                <FormError :message="errors.religion" />
              </div>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="space-y-2">
                <label class="label-tiny">Tempat Lahir <span class="text-red-500">*</span></label>
                <input v-model="form.birth_place" @input="clearFieldError('birth_place')" @blur="validateRequiredField('birth_place', 'Tempat lahir')" class="modern-input !rounded-2xl" :class="{'!border-red-500': errors.birth_place}" placeholder="Contoh: Probolinggo" />
                <FormError :message="errors.birth_place" />
              </div>
              <div class="space-y-2">
                <label class="label-tiny">Tanggal Lahir <span class="text-red-500">*</span></label>
                <div class="relative group">
                  <input v-model="localBirthDateDisplay" @input="onBirthDateInput" @blur="onBirthDateBlur" type="text" inputmode="numeric" placeholder="HH/BB/TTTT" class="modern-input !rounded-2xl pr-11" :class="{'!border-red-500': errors.birth_date}" />
                  <div @click="triggerCalendar" class="absolute right-0 top-0 h-full w-12 flex items-center justify-center text-slate-400 cursor-pointer hover:bg-slate-50 transition-colors rounded-r-2xl">
                    <CalendarAltIcon class="w-4 h-4" />
                    <input ref="birthDateRef" type="date" v-model="form.birth_date" class="absolute inset-0 opacity-0" style="color-scheme: light;" />
                  </div>
                </div>
                <FormError :message="errors.birth_date" />
                <p v-if="!errors.birth_date" class="text-[10px] text-slate-400 font-medium px-1 flex items-center gap-1.5">
                  <CalendarAltIcon class="w-3 h-3 shrink-0" />
                  Ketik <span class="font-black text-slate-500">HH/BB/TTTT</span> atau klik ikon
                </p>
              </div>
            </div>
          </div>

          <!-- Identity Info -->
          <div class="space-y-6 pt-6 border-t border-slate-100">
            <div class="flex items-center gap-3">
              <div class="w-1.5 h-4 bg-emerald-500 rounded-full"></div>
              <h4 class="text-[11px] font-black text-slate-700 uppercase tracking-widest">Identitas Siswa</h4>
            </div>

            <div class="grid grid-cols-2 gap-6">
              <div class="space-y-2">
                <label class="label-tiny">NIK (KTP/KK) <span class="text-red-500">*</span></label>
                <input type="number" v-model="form.nik" @input="clearFieldError('nik')" @blur="checkFieldUnique('nik')" class="modern-input !rounded-2xl" :class="{'!border-red-500': errors.nik}" placeholder="16 digit NIK" />
                <FormError :message="errors.nik" />
              </div>
              <div class="space-y-2">
                <label class="label-tiny">NISN <span class="text-red-500">*</span></label>
                <input type="number" v-model="form.nisn" @input="clearFieldError('nisn')" @blur="checkFieldUnique('nisn')" class="modern-input !rounded-2xl" :class="{'!border-red-500': errors.nisn}" placeholder="10 digit NISN" />
                <FormError :message="errors.nisn" />
              </div>
            </div>

            <div class="space-y-2">
              <label class="label-tiny">Nomor Induk Siswa (NIS) <span class="text-[8px] font-black bg-slate-100 text-slate-400 px-1.5 py-0.5 rounded ml-1 uppercase">Opsional</span></label>
              <input v-model="form.nis" @input="clearFieldError('nis')" @blur="checkFieldUnique('nis')" class="modern-input !rounded-2xl" :class="{'!border-red-500': errors.nis}" placeholder="Contoh: 12345" />
              <FormError :message="errors.nis" />
            </div>
          </div>

          <!-- Address Info -->
          <div class="space-y-6 pt-6 border-t border-slate-100">
            <div class="flex items-center gap-3">
              <div class="w-1.5 h-4 bg-amber-500 rounded-full"></div>
              <h4 class="text-[11px] font-black text-slate-700 uppercase tracking-widest">Wilayah & Alamat</h4>
            </div>

            <div class="grid grid-cols-2 gap-6">
              <div class="space-y-2">
                <label class="label-tiny">Provinsi <span class="text-red-500">*</span></label>
                <div class="relative">
                  <select v-model="form.province" @change="onProvinceChange" class="modern-input !rounded-2xl appearance-none !text-[11px]" :class="{'!border-red-500': errors.province}">
                    <option value="">Pilih Provinsi</option>
                    <option v-for="p in academicFilters.provinces" :key="p.id" :value="p.name">{{ p.name }}</option>
                  </select>
                  <ChevronDownIcon class="absolute right-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                </div>
                <FormError :message="errors.province" />
              </div>
              <div class="space-y-2">
                <label class="label-tiny">Kota / Kabupaten <span class="text-red-500">*</span></label>
                <div class="relative">
                  <select v-model="form.city" @change="onCityChange" :disabled="!form.province" class="modern-input !rounded-2xl appearance-none !text-[11px] disabled:opacity-50 disabled:bg-slate-50 disabled:cursor-not-allowed" :class="{'!border-red-500': errors.city}">
                    <option value="">Pilih Kota/Kab</option>
                    <option v-for="c in academicFilters.regencies" :key="c.id" :value="c.name">{{ c.name }}</option>
                  </select>
                  <ChevronDownIcon class="absolute right-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                </div>
                <FormError :message="errors.city" />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-6">
              <div class="space-y-2">
                <label class="label-tiny">Kecamatan <span class="text-red-500">*</span></label>
                <div class="relative">
                  <select v-model="form.district" @change="onDistrictChange" :disabled="!form.city" class="modern-input !rounded-2xl appearance-none !text-[11px] disabled:opacity-50 disabled:bg-slate-50 disabled:cursor-not-allowed" :class="{'!border-red-500': errors.district}">
                    <option value="">Pilih Kecamatan</option>
                    <option v-for="d in academicFilters.districts" :key="d.id" :value="d.name">{{ d.name }}</option>
                  </select>
                  <ChevronDownIcon class="absolute right-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                </div>
                <FormError :message="errors.district" />
              </div>
              <div class="space-y-2">
                <label class="label-tiny">Desa / Kelurahan <span class="text-red-500">*</span></label>
                <div class="relative">
                  <select v-model="form.village" @change="clearFieldError('village')" :disabled="!form.district" class="modern-input !rounded-2xl appearance-none !text-[11px] disabled:opacity-50 disabled:bg-slate-50 disabled:cursor-not-allowed" :class="{'!border-red-500': errors.village}">
                    <option value="">Pilih Desa/Kel</option>
                    <option v-for="v in academicFilters.villages" :key="v.id" :value="v.name">{{ v.name }}</option>
                  </select>
                  <ChevronDownIcon class="absolute right-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
                </div>
                <FormError :message="errors.village" />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-2">
                <label class="label-tiny text-center">RT <span class="text-red-500">*</span></label>
                <input v-model="form.rt" @input="clearFieldError('rt')" @blur="validateRequiredField('rt', 'RT')" class="modern-input !rounded-2xl text-center" :class="{'!border-red-500': errors.rt}" placeholder="000" />
                <FormError :message="errors.rt" />
              </div>
              <div class="space-y-2">
                <label class="label-tiny text-center">RW <span class="text-red-500">*</span></label>
                <input v-model="form.rw" @input="clearFieldError('rw')" @blur="validateRequiredField('rw', 'RW')" class="modern-input !rounded-2xl text-center" :class="{'!border-red-500': errors.rw}" placeholder="000" />
                <FormError :message="errors.rw" />
              </div>
            </div>

            <div class="space-y-2">
              <label class="label-tiny">Detail Alamat <span class="text-[8px] font-black bg-slate-100 text-slate-400 px-1.5 py-0.5 rounded ml-1 uppercase">Opsional</span></label>
              <textarea v-model="form.address" class="modern-input !rounded-2xl min-h-[80px] !py-4" :class="{'!border-red-500': errors.address}" placeholder="Jalan, Dusun, No. Rumah..."></textarea>
              <FormError :message="errors.address" />
            </div>
          </div>

          <!-- Contact Info -->
          <div class="space-y-6 pt-6 border-t border-slate-100">
            <div class="flex items-center gap-3">
              <div class="w-1.5 h-4 bg-rose-500 rounded-full"></div>
              <h4 class="text-[11px] font-black text-slate-700 uppercase tracking-widest">Informasi Kontak</h4>
            </div>

            <div class="space-y-2">
              <label class="label-tiny">Email Siswa <span class="text-red-500">*</span></label>
              <div class="relative">
                <MailIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300" />
                <input v-model="form.email" @input="clearFieldError('email')" @blur="checkFieldUnique('email')" type="email" class="modern-input !pl-12 !rounded-2xl" :class="{'!border-red-500': errors.email}" placeholder="name@example.com" />
              </div>
              <FormError :message="errors.email" />
            </div>

            <div class="space-y-2">
              <label class="label-tiny">Nomor WhatsApp <span class="text-red-500">*</span></label>
              <div class="flex gap-3">
                <div class="relative w-32 shrink-0">
                  <button @click="showCountryDropdown = !showCountryDropdown" type="button" class="modern-input !rounded-2xl !bg-slate-50 border-slate-200 flex items-center justify-between px-4">
                    <span class="text-lg">{{ selectedCountry.flag }}</span>
                    <span class="text-slate-700 font-bold">+{{ selectedCountry.code }}</span>
                    <ChevronDownIcon class="w-3 h-3 text-slate-400" />
                  </button>
                  <div v-if="showCountryDropdown" class="absolute top-full left-0 mt-2 w-64 bg-white border border-slate-100 rounded-2xl shadow-2xl z-[110] max-h-60 overflow-y-auto custom-scrollbar animate-scale-in">
                    <div v-for="c in countries" :key="c.code" @click="form.country_code = c.code; showCountryDropdown = false" class="flex items-center gap-3 px-4 py-3 hover:bg-slate-50 cursor-pointer border-b border-slate-50 last:border-0 transition-colors">
                      <span class="text-xl">{{ c.flag }}</span>
                      <span class="text-xs font-bold text-slate-700">{{ c.name }} (+{{ c.code }})</span>
                      <CheckIcon v-if="form.country_code === c.code" class="w-4 h-4 text-emerald-500 ml-auto" />
                    </div>
                  </div>
                </div>
                <div class="relative flex-1">
                  <WAIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300" />
                  <input type="number" v-model="form.phone_number" @input="clearFieldError('phone_number')" @blur="checkFieldUnique('phone_number')" class="modern-input !pl-12 !rounded-2xl" :class="{'!border-red-500': errors.phone_number}" placeholder="8123456..." />
                </div>
              </div>
              <FormError :message="errors.phone_number" />
            </div>
          </div>

          <!-- Academic Info -->
          <div class="space-y-6 pt-6 border-t border-slate-100">
            <div class="flex items-center gap-3">
              <div class="w-1.5 h-4 bg-blue-500 rounded-full"></div>
              <h4 class="text-[11px] font-black text-slate-700 uppercase tracking-widest">Informasi Akademik</h4>
            </div>

            <div class="grid grid-cols-4 gap-6">
              <div class="space-y-2">
                <label class="label-tiny">Angkatan <span class="text-red-500">*</span></label>
                <div class="relative">
                  <select v-model="form.entry_year" @change="onEntryYearChange" class="modern-input !rounded-2xl appearance-none !text-[11px]" :class="{'!border-red-500': errors.entry_year}">
                    <option :value="null">Pilih Angkatan</option>
                    <option v-for="y in academicFilters.years" :key="y.id" :value="y.year">{{ y.year }}</option>
                  </select>
                  <ChevronDownIcon class="absolute right-4 top-1/2 -translate-y-1/2 w-3 h-3 text-slate-400 pointer-events-none" />
                </div>
                <FormError :message="errors.entry_year" />
              </div>
              <div class="space-y-2">
                <label class="label-tiny">Jurusan <span class="text-red-500">*</span></label>
                <div class="relative">
                  <select v-model="form.major_id" @change="onMajorChange" :disabled="!form.entry_year" class="modern-input !rounded-2xl appearance-none !text-[11px] disabled:opacity-50 disabled:bg-slate-50" :class="{'!border-red-500': errors.major_id}">
                    <option :value="null">Pilih Jurusan</option>
                    <option v-for="j in filteredMajors" :key="j.id" :value="j.id">{{ j.name }}</option>
                  </select>
                  <ChevronDownIcon class="absolute right-4 top-1/2 -translate-y-1/2 w-3 h-3 text-slate-400 pointer-events-none" />
                </div>
                <FormError :message="errors.major_id" />
              </div>
              <div class="space-y-2">
                <label class="label-tiny">Kelas <span class="text-red-500">*</span></label>
                <div class="relative">
                  <select v-model="form.class_id" @change="onClassChange" :disabled="!form.major_id" class="modern-input !rounded-2xl appearance-none !text-[11px] disabled:opacity-50 disabled:bg-slate-50" :class="{'!border-red-500': errors.class_id}">
                    <option :value="null">Pilih Kelas</option>
                    <option v-for="c in filteredClasses" :key="c.id" :value="c.id">{{ c.name }}</option>
                  </select>
                  <ChevronDownIcon class="absolute right-4 top-1/2 -translate-y-1/2 w-3 h-3 text-slate-400 pointer-events-none" />
                </div>
                <FormError :message="errors.class_id" />
              </div>
              <div class="space-y-2">
                <label class="label-tiny">Status <span class="text-red-500">*</span></label>
                <div class="relative">
                  <select v-model="form.status" @change="clearFieldError('status')" class="modern-input !rounded-2xl appearance-none !text-[11px]" :class="{'!border-red-500': errors.status}">
                    <option value="active">Aktif</option>
                    <option value="inactive">Non-Aktif</option>
                    <option value="graduated">Lulus</option>
                  </select>
                  <ChevronDownIcon class="absolute right-4 top-1/2 -translate-y-1/2 w-3 h-3 text-slate-400 pointer-events-none" />
                </div>
                <FormError :message="errors.status" />
              </div>
            </div>

            <div class="space-y-2 pt-2">
              <label class="label-tiny">Keterangan / Catatan Status <span class="text-[8px] font-black bg-slate-100 text-slate-400 px-1.5 py-0.5 rounded ml-1 uppercase">Opsional</span></label>
              <textarea v-model="form.description" class="modern-input !rounded-2xl min-h-[80px] !py-4" :class="{'!border-red-500': errors.description}" placeholder="Contoh: Pindah sekolah ke luar kota, Cuti sakit 1 semester, Lulus tahun 2026..."></textarea>
              <FormError :message="errors.description" />
            </div>

            <!-- Wali Murid (Orang Tua) Selection Inline -->
            <div class="space-y-4 pt-4">
              <div class="flex items-center justify-between">
                <label class="label-tiny">Wali Murid (Orang Tua) <span class="text-red-500">*</span></label>
                <span v-if="selectedParent" class="text-[9px] font-black text-indigo-500 uppercase bg-indigo-50 px-2 py-0.5 rounded-full animate-pulse">Terpilih: {{ selectedParent.name }}</span>
              </div>
              
              <div class="space-y-3 bg-slate-50/50 p-5 rounded-[2rem] border-2 border-slate-100/50 shadow-inner">
                <!-- Search Box -->
                <div class="relative group">
                  <SearchIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300 group-focus-within:text-indigo-500 transition-colors" />
                  <input v-model="parentSearchQuery" placeholder="Cari Nama / Email / No. HP..." 
                         class="w-full bg-white border-2 border-slate-100 focus:border-indigo-500 rounded-xl py-3 pl-12 pr-4 text-xs font-bold shadow-sm focus:outline-none transition-all" />
                </div>

                <!-- Scrollable List -->
                <div class="max-h-[280px] overflow-y-auto space-y-2 pr-2 custom-scrollbar">
                  <div v-if="filteredParents.length === 0" class="py-10 text-center">
                    <p class="text-[10px] font-bold text-slate-400">Orang tua tidak ditemukan</p>
                  </div>
                  
                  <div v-for="p in filteredParents" :key="p.id" @click="selectParent(p)" 
                       class="flex items-center gap-4 p-3 bg-white border-2 rounded-2xl hover:border-indigo-400 cursor-pointer transition-all group"
                       :class="form.parent_id === p.id ? 'border-indigo-600 bg-indigo-50/20' : 'border-slate-100'">
                    
                    <!-- Checkbox Style -->
                    <div class="w-5 h-5 rounded-lg border-2 flex items-center justify-center transition-all shrink-0"
                         :class="form.parent_id === p.id ? 'bg-indigo-600 border-indigo-600' : 'border-slate-200 group-hover:border-indigo-400'">
                      <CheckIcon v-if="form.parent_id === p.id" class="w-3 h-3 text-white" />
                    </div>

                    <div class="flex-1 min-w-0">
                      <div class="flex items-center justify-between gap-2">
                        <span class="text-[11px] font-black text-slate-700 truncate group-hover:text-indigo-600 transition-colors">{{ p.name }}</span>
                        <span v-if="p.phone_number" class="text-[9px] font-bold text-indigo-400 shrink-0">{{ p.phone_number }}</span>
                      </div>
                      <span class="text-[9px] text-slate-400 font-medium block truncate">{{ p.email || 'Tidak ada email' }}</span>
                    </div>
                  </div>
                </div>
              </div>
              <FormError :message="errors.parent_id" />
            </div>
          </div>
        </div>

        <div class="p-8 bg-slate-50/50 border-t border-slate-100 flex gap-4 justify-end">
          <button @click="emit('update:modelValue', false)" class="btn-secondary !rounded-2xl !px-8">Batal</button>
          <button @click="validateAndSave" :disabled="submitting" class="btn-primary !rounded-2xl !px-12 shadow-lg shadow-indigo-100 flex items-center gap-3">
            <CheckCircleIcon v-if="!submitting" class="w-5 h-5" />
            <div v-else class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
            <span>{{ submitting ? 'Menyimpan...' : 'Simpan Data' }}</span>
          </button>
        </div>
      </div>
    </div>
  </transition>
</Teleport>
</template>
