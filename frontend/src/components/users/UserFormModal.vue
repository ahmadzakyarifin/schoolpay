<script setup>
import { ref, computed, watch } from 'vue'
import { 
  X as CloseIcon, 
  Mail as MailIcon, 
  AlertCircle as AlertIcon, 
  ChevronDown as ChevronDownIcon, 
  Calendar as CalendarAltIcon, 
  MessageCircle as WAIcon, 
  CheckCircle2 as CheckCircleIcon,
  Check as CheckIcon 
} from 'lucide-vue-next'
import FormError from '../ui/FormError.vue'
import userService from '../../services/user.service'

const props = defineProps({
  modelValue: Boolean,
  isEditing: Boolean,
  form: Object,
  errors: Object,
  submitting: Boolean,
  birthDateDisplay: String
})

const emit = defineEmits(['update:modelValue', 'update:birthDateDisplay', 'save', 'local-validation-failed', 'clear-field-error', 'set-field-error'])

const clearFieldError = (field) => {
  emit('clear-field-error', field)
}

const setFieldError = (field, messages) => {
  emit('set-field-error', { field, messages })
}

const validateNameField = () => {
  const value = props.form.name?.trim() || ''
  if (!value) {
    setFieldError('name', ['Nama lengkap wajib diisi'])
    return false
  }
  if (!/^[a-zA-Z\s]*$/.test(value)) {
    setFieldError('name', ['Nama hanya boleh berisi huruf'])
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

const showCountryDropdown = ref(false)

const educationOptions = [
  'SD/Sederajat',
  'SMP/Sederajat',
  'SMA/SMK/Sederajat',
  'D3',
  'D4',
  'S1',
  'S2',
  'S3'
]

const occupationOptions = [
  'PNS',
  'TNI/Polri',
  'Pegawai Swasta',
  'Wiraswasta / Pedagang',
  'Petani / Peternak',
  'Nelayan',
  'Buruh (Harian/Pabrik)',
  'Tenaga Kerja Indonesia (TKI/TKW)',
  'Yang lain:'
]

const incomeOptions = [
  '< 1 juta',
  '1 - 3 juta',
  '3 - 5 juta',
  '> 5 juta'
]

const countries = [
  { code: '62', name: 'Indonesia', flag: '🇮🇩' },
  { code: '60', name: 'Malaysia', flag: '🇲🇾' },
  { code: '65', name: 'Singapore', flag: '🇸🇬' },
  { code: '1', name: 'USA', flag: '🇺🇸' }
]

const selectedCountry = computed(() => {
  return countries.find(c => c.code === props.form.country_code) || countries[0]
})

const birthDateRef = ref(null)

const triggerCalendar = () => {
  if (birthDateRef.value) {
    try {
      birthDateRef.value.showPicker()
    } catch (e) {
      birthDateRef.value.click()
    }
  }
}

const localBirthDateDisplay = computed({
  get: () => props.birthDateDisplay,
  set: (val) => {
    let cleaned = val.replace(/\D/g, '').substring(0, 8)
    
    let formatted = cleaned
    if (cleaned.length > 2) {
      formatted = cleaned.substring(0, 2) + '/' + cleaned.substring(2)
    }
    if (cleaned.length > 4) {
      formatted = formatted.substring(0, 5) + '/' + formatted.substring(5)
    }
    emit('update:birthDateDisplay', formatted)
  }
})

const parseBirthDateDisplay = (value) => {
  const display = String(value || '').trim()
  if (!display) return { valid: true, empty: true, iso: null }

  if (!/^\d{2}\/\d{2}\/\d{4}$/.test(display)) {
    return { valid: false, message: 'Format tanggal lahir harus HH/BB/TTTT' }
  }

  const [dayPart, monthPart, yearPart] = display.split('/')
  const day = Number(dayPart)
  const month = Number(monthPart)
  const year = Number(yearPart)
  const date = new Date(year, month - 1, day)
  const currentYear = new Date().getFullYear()

  if (
    Number.isNaN(date.getTime()) ||
    date.getFullYear() !== year ||
    date.getMonth() !== month - 1 ||
    date.getDate() !== day
  ) {
    return { valid: false, message: 'Tanggal lahir tidak sesuai dengan kalender nyata' }
  }

  if (year < 1900 || year > currentYear) {
    return { valid: false, message: 'Tahun lahir tidak valid' }
  }

  return {
    valid: true,
    iso: `${yearPart}-${monthPart.padStart(2, '0')}-${dayPart.padStart(2, '0')}`
  }
}

const parseIsoBirthDate = (value) => {
  const iso = String(value || '').slice(0, 10)
  const match = iso.match(/^(\d{4})-(\d{2})-(\d{2})$/)
  if (!match) return null

  const [, yearPart, monthPart, dayPart] = match
  const year = Number(yearPart)
  const month = Number(monthPart)
  const day = Number(dayPart)
  const date = new Date(year, month - 1, day)

  if (
    Number.isNaN(date.getTime()) ||
    date.getFullYear() !== year ||
    date.getMonth() !== month - 1 ||
    date.getDate() !== day
  ) {
    return null
  }

  return `${dayPart}/${monthPart}/${yearPart}`
}

const handleBirthDateInput = () => {
  const parsed = parseBirthDateDisplay(localBirthDateDisplay.value)
  props.form.birth_date = parsed.valid ? parsed.iso : null
  if (!localBirthDateDisplay.value || parsed.valid) {
    clearFieldError('birth_date')
  }
}

const onBirthDateBlur = () => {
  const value = localBirthDateDisplay.value || ''
  if (!value) {
    emit('clear-field-error', 'birth_date')
    return
  }
  const result = parseBirthDateDisplay(value)
  if (!result.valid) {
    emit('set-field-error', { field: 'birth_date', messages: [result.message] })
  } else {
    emit('clear-field-error', 'birth_date')
  }
}

const checkFieldUnique = async (field) => {
  let value = props.form[field]
  if (value === undefined || value === null) return

  const cleanValue = value.toString().trim()
  const requiredFields = ['email', 'phone_number']
  if (!cleanValue) {
    if (requiredFields.includes(field)) {
      setFieldError(field, [field === 'email' ? 'Email wajib diisi' : 'Nomor WhatsApp wajib diisi'])
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
  } else if (field === 'email') {
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(cleanValue)) {
      setFieldError(field, ['Format email tidak valid (contoh: user@gmail.com)'])
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
    const res = await userService.checkUnique(field, valueForCheck, excludeId)
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

const validateAndSave = () => {
  const localErrors = {}
  
  if (!props.form.name?.trim()) {
    localErrors.name = ['Nama lengkap wajib diisi']
  } else if (!/^[a-zA-Z\s]*$/.test(props.form.name)) {
    localErrors.name = ['Nama hanya boleh berisi huruf']
  } else if (props.form.name.trim().length < 2) {
    localErrors.name = ['Nama minimal berisi 2 karakter']
  }

  if (!props.form.email?.trim()) {
    localErrors.email = ['Email wajib diisi']
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(props.form.email)) {
    localErrors.email = ['Format email tidak valid (contoh: user@gmail.com)']
  }
  
  let cleanPhone = props.form.phone_number?.toString().trim() || ''
  if (cleanPhone.startsWith('0')) {
    cleanPhone = cleanPhone.replace(/^0+/, '')
  }

  if (!cleanPhone) {
    localErrors.phone_number = ['Nomor WhatsApp wajib diisi']
  } else if (cleanPhone.length < 9) { 
    localErrors.phone_number = ['Nomor WhatsApp minimal 9 digit']
  } else if (cleanPhone.length > 15) {
    localErrors.phone_number = ['Nomor WhatsApp maksimal 15 digit']
  }

 
  if (!localErrors.phone_number) {
    props.form.phone_number = cleanPhone
  }

  if (props.form.role === 'parent' && props.form.nik) {
    const cleanNIK = props.form.nik.toString().trim()
    if (cleanNIK.length !== 16) {
      localErrors.nik = ['NIK harus berjumlah 16 digit']
    } else {
      props.form.nik = cleanNIK 
    }
  }

  const parsedBirthDate = parseBirthDateDisplay(localBirthDateDisplay.value)

  // Validasi Tanggal Lahir (Kalender Nyata)
  if (props.form.role === 'parent' && localBirthDateDisplay.value) {
    if (!parsedBirthDate.valid) {
      localErrors.birth_date = [parsedBirthDate.message]
      props.form.birth_date = null
    } else {
      props.form.birth_date = parsedBirthDate.iso
    }
  }

  if (props.form.role === 'admin') {
    props.form.nik = null
    props.form.birth_date = null
    props.form.address = null
    props.form.education = null
    props.form.occupation = null
    props.form.income = null
    emit('update:birthDateDisplay', '')
  }

  if (Object.keys(localErrors).length > 0) {
    emit('local-validation-failed', localErrors)
    return
  }

  if (!localBirthDateDisplay.value) {
    props.form.birth_date = null
  }

  const hasBlockingError = Object.entries(props.errors || {}).some(([field, value]) => {
    if (field === '_general') return false
    return Array.isArray(value) ? value.some(Boolean) : !!value
  })
  if (hasBlockingError) return

  emit('save')
}

watch(() => props.birthDateDisplay, (val) => {
  const parsed = parseBirthDateDisplay(val)
  if (!val || parsed.valid) {
    clearFieldError('birth_date')
  }
  props.form.birth_date = parsed.valid ? parsed.iso : null
})

watch(() => props.form.birth_date, (val) => {
  const display = parseIsoBirthDate(val)
  if (display) {
    emit('update:birthDateDisplay', display)
    clearFieldError('birth_date')
  }
})


</script>

<template>
  <Teleport to="body">
    <div v-if="modelValue" class="fixed inset-0 z-[500] flex items-center justify-center p-6">
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="emit('update:modelValue', false)"></div>
      <div class="white-card w-full max-w-xl relative z-10 overflow-hidden shadow-[0_20px_50px_rgba(0,0,0,0.2)] animate-scale-in !rounded-[2.5rem]">
        <div class="p-10 border-b border-slate-100 flex items-center justify-between bg-white">
          <div>
            <h2 class="text-2xl font-black text-slate-900 tracking-tight">{{ isEditing ? 'Edit Pengguna' : 'Tambah Pengguna Baru' }}</h2>
            <p class="text-slate-500 font-medium mt-1">Lengkapi informasi akses sistem di bawah ini</p>
          </div>
          <button @click="emit('update:modelValue', false)" class="p-3 hover:bg-slate-50 rounded-2xl text-slate-400 transition-all">
            <CloseIcon class="w-6 h-6" />
          </button>
        </div>

        <div class="p-10 space-y-6 max-h-[70vh] overflow-y-auto scrollbar-thin scrollbar-indigo custom-scrollbar">
          <!-- General Error fallback (minimalist) -->
          <FormError :message="errors._general" />

          <!-- Name Input -->
          <div class="space-y-2">
            <label class="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em] px-1 flex items-center gap-1">
              Nama Lengkap <span class="text-red-500">*</span>
            </label>
            <div class="relative">
              <input v-model="form.name" @input="clearFieldError('name')" @blur="validateNameField" :class="['modern-input !rounded-2xl', errors.name ? '!border-red-500 !ring-red-50' : '']" placeholder="Contoh: Sulaiman" />
            </div>
            <FormError :message="errors.name" />
          </div>

          <!-- Email & Role -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="space-y-2">
              <label class="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em] px-1 flex items-center gap-1">
                Email <span class="text-red-500">*</span>
              </label>
              <div class="relative">
                <MailIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300" />
                <input v-model="form.email" @input="clearFieldError('email')" @blur="checkFieldUnique('email')" :class="['modern-input !pl-12 !rounded-2xl', errors.email ? '!border-red-500 !ring-red-50' : '']" placeholder="name@example.com" />
              </div>
              <FormError :message="errors.email" />
            </div>
            <div class="space-y-2">
              <label class="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em] px-1 flex items-center gap-1">
                Role Akses <span class="text-red-500">*</span>
              </label>
              <div class="relative">
                <select v-model="form.role" @change="clearFieldError('role')" class="modern-input !rounded-2xl appearance-none">
                  <option value="admin">Administrator</option>
                  <option value="parent">Wali Murid (Parent)</option>
                </select>
                <div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
                  <ChevronDownIcon class="w-4 h-4" />
                </div>
              </div>
              <FormError :message="errors.role" />
            </div>
          </div>

          <!-- Parent Specific Fields -->
          <transition name="fade">
            <div v-if="form.role === 'parent'" class="space-y-6 pt-4 border-t border-slate-100 animate-slide-up">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="space-y-2">
                  <label class="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em] px-1 flex items-center gap-1">
                    NIK <span class="text-[8px] font-black bg-slate-100 text-slate-400 px-1.5 py-0.5 rounded uppercase tracking-widest ml-1">Opsional</span>
                  </label>
                  <input type="number" v-model="form.nik" @input="clearFieldError('nik')" @blur="checkFieldUnique('nik')" :class="['modern-input !rounded-2xl', errors.nik ? '!border-red-500 !ring-red-50' : '']" placeholder="Masukkan 16 digit NIK" />
                  <FormError :message="errors.nik" />
                </div>
                <div class="space-y-2">
                  <label class="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em] px-1 flex items-center gap-1">
                    Tanggal Lahir <span class="text-[8px] font-black bg-slate-100 text-slate-400 px-1.5 py-0.5 rounded uppercase tracking-widest ml-1">Opsional</span>
                  </label>
                  <div class="relative group">
                    <input 
                      v-model="localBirthDateDisplay"
                      @input="handleBirthDateInput"
                      @blur="onBirthDateBlur"
                      type="text" 
                      inputmode="numeric"
                      placeholder="HH/BB/TTTT"
                      class="modern-input !rounded-2xl pr-11"
                      :class="errors.birth_date ? '!border-red-500 !ring-red-50' : ''"
                    />
                    <div 
                      @click="triggerCalendar"
                      class="absolute right-0 top-0 h-full w-12 flex items-center justify-center group-focus-within:text-indigo-500 text-slate-400 cursor-pointer overflow-hidden hover:bg-slate-50 transition-colors rounded-r-2xl"
                    >
                      <div class="w-px h-4 bg-slate-200 absolute left-0"></div>
                      <CalendarAltIcon class="w-4 h-4" />
                      <input 
                        ref="birthDateRef"
                        type="date" 
                        v-model="form.birth_date" 
                        class="absolute inset-0 opacity-0"
                        style="color-scheme: light;"
                      />
                    </div>
                  </div>
                  <FormError :message="errors.birth_date" />
                  <p v-if="!errors.birth_date" class="text-[10px] text-slate-400 font-medium px-1 flex items-center gap-1.5">
                    <CalendarAltIcon class="w-3 h-3 shrink-0" />
                    Ketik langsung <span class="font-black text-slate-500">HH/BB/TTTT</span> atau klik ikon kalender
                  </p>
                </div>
              </div>

              <div class="space-y-2">
                <label class="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em] px-1 flex items-center gap-1">
                  Alamat <span class="text-[8px] font-black bg-slate-100 text-slate-400 px-1.5 py-0.5 rounded uppercase tracking-widest ml-1">Opsional</span>
                </label>
                <textarea 
                  v-model="form.address" 
                  @input="clearFieldError('address')"
                  :class="['modern-input !rounded-2xl min-h-[100px] !py-4', errors.address ? '!border-red-500 !ring-red-50' : '']" 
                  placeholder="Contoh: RT 001,RW 001,Dusun Mega, Desa Sukodadi, Kecamatan Paiton, Kabupaten Probolinggo"
                ></textarea>
                <FormError :message="errors.address" />
              </div>

              <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div class="space-y-2">
                  <label class="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em] px-1 flex items-center gap-1">
                    Pendidikan <span class="text-[8px] font-black bg-slate-100 text-slate-400 px-1.5 py-0.5 rounded uppercase tracking-widest ml-1">Opsional</span>
                  </label>
                  <div class="relative">
                    <select v-model="form.education" @change="clearFieldError('education')" class="modern-input !rounded-2xl appearance-none !text-[11px]">
                      <option value="">Pilih Pendidikan</option>
                      <option v-for="opt in educationOptions" :key="opt" :value="opt">{{ opt }}</option>
                    </select>
                    <div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
                      <ChevronDownIcon class="w-4 h-4" />
                    </div>
                  </div>
                  <FormError :message="errors.education" />
                </div>
                <div class="space-y-2">
                  <label class="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em] px-1 flex items-center gap-1">
                    Pekerjaan <span class="text-[8px] font-black bg-slate-100 text-slate-400 px-1.5 py-0.5 rounded uppercase tracking-widest ml-1">Opsional</span>
                  </label>
                  <div class="relative">
                    <select v-model="form.occupation" @change="clearFieldError('occupation')" class="modern-input !rounded-2xl appearance-none !text-[11px]">
                      <option value="">Pilih Pekerjaan</option>
                      <option v-for="opt in occupationOptions" :key="opt" :value="opt">{{ opt }}</option>
                    </select>
                    <div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
                      <ChevronDownIcon class="w-4 h-4" />
                    </div>
                  </div>
                  <FormError :message="errors.occupation" />
                </div>
                <div class="space-y-2">
                  <label class="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em] px-1 flex items-center gap-1">
                    Penghasilan <span class="text-[8px] font-black bg-slate-100 text-slate-400 px-1.5 py-0.5 rounded uppercase tracking-widest ml-1">Opsional</span>
                  </label>
                  <div class="relative">
                    <select v-model="form.income" @change="clearFieldError('income')" class="modern-input !rounded-2xl appearance-none !text-[11px]">
                      <option value="">Pilih Penghasilan</option>
                      <option v-for="opt in incomeOptions" :key="opt" :value="opt">{{ opt }}</option>
                    </select>
                    <div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
                      <ChevronDownIcon class="w-4 h-4" />
                    </div>
                  </div>
                  <FormError :message="errors.income" />
                </div>
              </div>
            </div>
          </transition>

          <!-- Phone with Industry-standard Country Selector -->
          <div class="space-y-2">
            <label class="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em] px-1 flex items-center gap-1">
              Nomor WhatsApp <span class="text-red-500">*</span>
            </label>
            <div class="flex gap-3">
              <div class="relative w-40 shrink-0">
                <button @click="showCountryDropdown = !showCountryDropdown" type="button" class="modern-input !rounded-2xl !bg-slate-50 border-slate-200 flex items-center justify-between px-4">
                  <span class="flex items-center gap-2">
                    <span class="text-lg">{{ selectedCountry.flag }}</span>
                    <span class="text-slate-700">+{{ selectedCountry.code }}</span>
                  </span>
                  <ChevronDownIcon class="w-4 h-4 text-slate-400" />
                </button>
                
                <div v-if="showCountryDropdown" class="absolute top-full left-0 mt-2 w-64 bg-white border border-slate-100 rounded-2xl shadow-2xl z-[110] max-h-60 overflow-y-auto custom-scrollbar animate-scale-in">
                  <div v-for="c in countries" :key="c.code" @click="form.country_code = c.code; showCountryDropdown = false" 
                    class="flex items-center gap-3 px-4 py-3 hover:bg-slate-50 cursor-pointer transition-colors border-b border-slate-50 last:border-0">
                    <span class="text-xl">{{ c.flag }}</span>
                    <div class="flex flex-col">
                      <span class="text-xs font-bold text-slate-700">{{ c.name }}</span>
                      <span class="text-[10px] text-slate-400 font-medium">+{{ c.code }}</span>
                    </div>
                    <CheckIcon v-if="form.country_code === c.code" class="w-4 h-4 text-emerald-500 ml-auto" />
                  </div>
                </div>
              </div>

               <div class="relative flex-1">
                <WAIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300" />
                <input type="number" v-model="form.phone_number" @input="clearFieldError('phone_number')" @blur="checkFieldUnique('phone_number')" :class="['modern-input !pl-12 !rounded-2xl', errors.phone_number ? '!border-red-500 !ring-red-50' : '']" placeholder="8123456..." />
              </div>
            </div>
            <FormError :message="errors.phone_number" />
          </div>

          <div class="bg-indigo-50/50 p-6 rounded-[2rem] border border-indigo-100/50 flex gap-4 items-start">
            <AlertIcon class="w-5 h-5 text-indigo-500 shrink-0 mt-0.5" />
            <div class="text-[10px] leading-relaxed text-indigo-700 font-medium">
              <span class="font-black uppercase tracking-widest block mb-1">Informasi Penting</span>
              Untuk akun **Orang Tua**, pengelolaan data siswa dilakukan melalui menu **Manajemen Siswa** untuk memastikan validitas data wali murid pada setiap siswa.
            </div>
          </div>
        </div>

        <div class="p-10 bg-slate-50/50 border-t border-slate-100 flex gap-4 justify-end">
          <button @click="emit('update:modelValue', false)" class="btn-secondary !rounded-2xl !px-8">Batal</button>
          <button @click="validateAndSave" :disabled="submitting" class="btn-primary !rounded-2xl !px-12 shadow-lg shadow-indigo-100 flex items-center gap-3">
            <CheckCircleIcon v-if="!submitting" class="w-5 h-5" />
            <div v-else class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
            <span>{{ submitting ? 'Menyimpan...' : 'Simpan Data' }}</span>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>
