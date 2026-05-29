<script setup>
import { computed, ref, watch } from 'vue'
import { 
  X as CloseIcon, 
  CheckCircle2 as CheckCircleIcon,
  Loader2 as LoaderIcon,
  GraduationCap as MajorIcon,
  Users as ClassIcon,
  Calendar as CalendarIcon,
  Clock as ClockIcon,
  Layers as LayersIcon,
  Info as InfoIcon
} from 'lucide-vue-next'
import FormError from '../ui/FormError.vue'
import financeService from '../../services/finance.service'

const props = defineProps({
  modelValue: Boolean,
  isEditing: Boolean,
  form: Object,
  errors: Object,
  submitting: Boolean,
  billTypes: Array,
  classes: Array,
  majors: Array,
  startDateDisplay: String,
  endDateDisplay: String
})

const emit = defineEmits([
  'update:modelValue', 
  'update:startDateDisplay', 
  'update:endDateDisplay', 
  'save', 
  'local-validation-failed', 
  'clear-field-error'
])

const clearFieldError = (field) => {
  emit('clear-field-error', field)
}

const setFieldError = (field, messages) => {
  emit('local-validation-failed', { [field]: messages })
}

const selectedTargetIds = () => {
  if (props.form.target_type === 'all') return [0]
  if (props.isEditing) {
    const selected = props.form.target_id || props.form.target_ids?.[0]
    return selected ? [selected] : []
  }
  return Array.isArray(props.form.target_ids) ? props.form.target_ids : []
}

const targetLabel = () => {
  if (props.form.target_type === 'all') return 'semua siswa'
  return props.form.target_type === 'major' ? 'jurusan' : 'kelas'
}

const checkUniqueRule = async () => {
  if (!props.form.bill_type_id) return true

  const targetIds = selectedTargetIds()
  if (props.form.target_type !== 'all' && targetIds.length === 0) return true

  try {
    const excludeId = props.isEditing ? props.form.id : 0
    for (const targetId of targetIds) {
      const res = await financeService.checkUniqueBillingRule(
        props.form.bill_type_id,
        props.form.target_type,
        targetId,
        null,
        excludeId
      )
      if (res.data?.data?.is_unique === false) {
        setFieldError('target_id', [`Aturan tagihan untuk kategori dan ${targetLabel()} ini sudah terdaftar`])
        return false
      }
    }
    clearFieldError('target_id')
    return true
  } catch (err) {
    console.error('Gagal mengecek duplikasi aturan tagihan:', err)
    return true
  }
}

const handleBillTypeChange = () => {
  clearFieldError('bill_type_id')
  void checkUniqueRule()
}

const setTargetType = (targetType) => {
  props.form.target_type = targetType
  props.form.target_ids = []
  props.form.target_id = 0
  props.form.class_ids = []
  props.form.class_id = null
  clearFieldError('target_id')
  void checkUniqueRule()
}


const validateAmount = () => {
  if (props.form.amount === null || props.form.amount === undefined || props.form.amount <= 0) {
    setFieldError('amount', ['Nominal tagihan wajib diisi dengan angka lebih dari 0'])
  } else {
    clearFieldError('amount')
  }
}

const validateDueDate = () => {
  if (props.form.period_type === 'bulanan' && (!props.form.due_day || props.form.due_day < 1 || props.form.due_day > 31)) {
    setFieldError('due_day', ['Tanggal jatuh tempo harus di antara 1 - 31'])
  } else {
    clearFieldError('due_day')
  }
}

const validateStartDate = () => {
  if (localStartDateDisplay.value && localStartDateDisplay.value.length === 10) {
    const parts = localStartDateDisplay.value.split('/')
    if (parts.length === 3) {
      const day = parts[0].padStart(2, '0')
      const month = parts[1].padStart(2, '0')
      const year = parts[2]
      const numYear = parseInt(year, 10)
      const numMonth = parseInt(month, 10)
      const numDay = parseInt(day, 10)
      const dateObj = new Date(numYear, numMonth - 1, numDay)
      if (dateObj.getFullYear() !== numYear || dateObj.getMonth() + 1 !== numMonth || dateObj.getDate() !== numDay) {
        setFieldError('start_date', ['Tanggal mulai tidak valid sesuai kalender nyata'])
      } else {
        props.form.start_date = `${year}-${month}-${day}`
        clearFieldError('start_date')
      }
    } else {
      setFieldError('start_date', ['Format tanggal mulai tidak valid (HH/BB/TTTT)'])
    }
  } else {
    setFieldError('start_date', ['Tanggal mulai wajib diisi dengan format HH/BB/TTTT'])
  }
}

const validateEndDate = () => {
  if (localEndDateDisplay.value && localEndDateDisplay.value.length === 10) {
    const parts = localEndDateDisplay.value.split('/')
    if (parts.length === 3) {
      const day = parts[0].padStart(2, '0')
      const month = parts[1].padStart(2, '0')
      const year = parts[2]
      const numYear = parseInt(year, 10)
      const numMonth = parseInt(month, 10)
      const numDay = parseInt(day, 10)
      const dateObj = new Date(numYear, numMonth - 1, numDay)
      if (dateObj.getFullYear() !== numYear || dateObj.getMonth() + 1 !== numMonth || dateObj.getDate() !== numDay) {
        setFieldError('end_date', ['Tanggal selesai tidak valid sesuai kalender nyata'])
      } else {
        props.form.end_date = `${year}-${month}-${day}`
        if (props.form.start_date) {
          const start = new Date(props.form.start_date)
          const end = new Date(`${year}-${month}-${day}`)
          if (start > end) {
            setFieldError('end_date', ['Tanggal selesai tidak boleh lebih awal dari tanggal mulai'])
            return
          }
        }
        clearFieldError('end_date')
      }
    } else {
      setFieldError('end_date', ['Format tanggal selesai tidak valid (HH/BB/TTTT)'])
    }
  } else {
    setFieldError('end_date', ['Tanggal selesai wajib diisi dengan format HH/BB/TTTT'])
  }
}

const validateMaxInstallment = () => {
  if (props.form.allow_installment) {
    if (!props.form.max_installment || props.form.max_installment <= 0) {
      setFieldError('max_installment', ['Maksimal cicilan wajib diisi dengan angka lebih dari 0'])
    } else {
      clearFieldError('max_installment')
    }
  } else {
    clearFieldError('max_installment')
  }
}

const formattedAmount = computed({
  get() {
    if (props.form.amount === null || props.form.amount === undefined) return ''
    const parts = props.form.amount.toString().split('.')
    let num = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, '.')
    if (parts.length > 1) {
      num += ',' + parts[1]
    }
    return num
  },
  set(val) {
    clearFieldError('amount')
    if (!val) {
      props.form.amount = 0
      return
    }
    let clean = val.replace(/[^0-9,]/g, '')
    clean = clean.replace(',', '.')
    const parsed = parseFloat(clean)
    props.form.amount = isNaN(parsed) ? 0 : parsed
  }
})

// Multi-select Helpers for Target
const toggleTarget = (id) => {
  if (props.isEditing) {
    props.form.target_id = id
    props.form.target_ids = [id]
    clearFieldError('target_id')
    void checkUniqueRule()
    return
  }
  const idx = props.form.target_ids.indexOf(id)
  if (idx > -1) {
    props.form.target_ids.splice(idx, 1)
  } else {
    props.form.target_ids.push(id)
  }
  clearFieldError('target_id')
  void checkUniqueRule()
}

const selectAllTargets = () => {
  const list = props.form.target_type === 'major' ? props.majors : props.classes
  const allIds = list.map(item => item.id)
  if (props.form.target_ids?.length === allIds.length) {
    props.form.target_ids = []
  } else {
    props.form.target_ids = allIds
  }
  clearFieldError('target_id')
  void checkUniqueRule()
}

const setPeriodType = (periodType) => {
  props.form.period_type = periodType
  if (periodType !== 'bulanan') clearFieldError('due_day')
}

const toggleInstallment = () => {
  props.form.allow_installment = !props.form.allow_installment
  if (!props.form.allow_installment) {
    props.form.max_installment = null
    clearFieldError('max_installment')
  }
}

// Calendar Triggers & Display Logic
const startDateRef = ref(null)
const endDateRef = ref(null)

const triggerStartCalendar = () => {
  if (startDateRef.value) {
    try { startDateRef.value.showPicker() } catch(e) { startDateRef.value.click() }
  }
}

const triggerEndCalendar = () => {
  if (endDateRef.value) {
    try { endDateRef.value.showPicker() } catch(e) { endDateRef.value.click() }
  }
}

const localStartDateDisplay = computed({
  get: () => props.startDateDisplay,
  set: (val) => {
    let cleaned = val.replace(/\D/g, '').substring(0, 8)
    let formatted = cleaned
    if (cleaned.length > 2) formatted = cleaned.substring(0, 2) + '/' + cleaned.substring(2)
    if (cleaned.length > 4) formatted = formatted.substring(0, 5) + '/' + formatted.substring(5)
    emit('update:startDateDisplay', formatted)
  }
})

const localEndDateDisplay = computed({
  get: () => props.endDateDisplay,
  set: (val) => {
    let cleaned = val.replace(/\D/g, '').substring(0, 8)
    let formatted = cleaned
    if (cleaned.length > 2) formatted = cleaned.substring(0, 2) + '/' + cleaned.substring(2)
    if (cleaned.length > 4) formatted = formatted.substring(0, 5) + '/' + formatted.substring(5)
    emit('update:endDateDisplay', formatted)
  }
})

const validateAndSave = async () => {
  const localErrors = {}
  
  if (!props.form.bill_type_id) {
    localErrors.bill_type_id = ['Kategori tagihan wajib dipilih']
  }

  if (props.form.target_type !== 'all' && (!props.form.target_ids || props.form.target_ids.length === 0)) {
    localErrors.target_id = [`Pilih minimal satu ${props.form.target_type === 'major' ? 'jurusan' : 'kelas'}`]
  }

  if (props.form.amount === null || props.form.amount === undefined || props.form.amount <= 0) {
    localErrors.amount = ['Nominal tagihan wajib diisi dengan angka lebih dari 0']
  }

  if (props.form.period_type === 'bulanan' && (!props.form.due_day || props.form.due_day < 1 || props.form.due_day > 31)) {
    localErrors.due_day = ['Tanggal jatuh tempo harus di antara 1 - 31']
  }

  // Parse Start Date
  if (localStartDateDisplay.value && localStartDateDisplay.value.length === 10) {
    const parts = localStartDateDisplay.value.split('/')
    if (parts.length === 3) {
      const day = parts[0].padStart(2, '0')
      const month = parts[1].padStart(2, '0')
      const year = parts[2]
      const numYear = parseInt(year, 10)
      const numMonth = parseInt(month, 10)
      const numDay = parseInt(day, 10)
      const dateObj = new Date(numYear, numMonth - 1, numDay)
      if (dateObj.getFullYear() !== numYear || dateObj.getMonth() + 1 !== numMonth || dateObj.getDate() !== numDay) {
        localErrors.start_date = ['Tanggal mulai tidak valid sesuai kalender nyata']
      } else {
        props.form.start_date = `${year}-${month}-${day}`
      }
    } else {
      localErrors.start_date = ['Format tanggal mulai tidak valid (HH/BB/TTTT)']
    }
  } else {
    localErrors.start_date = ['Tanggal mulai wajib diisi dengan format HH/BB/TTTT']
    props.form.start_date = null
  }

  // Parse End Date
  if (localEndDateDisplay.value && localEndDateDisplay.value.length === 10) {
    const parts = localEndDateDisplay.value.split('/')
    if (parts.length === 3) {
      const day = parts[0].padStart(2, '0')
      const month = parts[1].padStart(2, '0')
      const year = parts[2]
      const numYear = parseInt(year, 10)
      const numMonth = parseInt(month, 10)
      const numDay = parseInt(day, 10)
      const dateObj = new Date(numYear, numMonth - 1, numDay)
      if (dateObj.getFullYear() !== numYear || dateObj.getMonth() + 1 !== numMonth || dateObj.getDate() !== numDay) {
        localErrors.end_date = ['Tanggal selesai tidak valid sesuai kalender nyata']
      } else {
        props.form.end_date = `${year}-${month}-${day}`
      }
    } else {
      localErrors.end_date = ['Format tanggal selesai tidak valid (HH/BB/TTTT)']
    }
  } else {
    localErrors.end_date = ['Tanggal selesai wajib diisi dengan format HH/BB/TTTT']
    props.form.end_date = null
  }

  if (!localErrors.start_date && !localErrors.end_date && props.form.start_date && props.form.end_date) {
    if (new Date(props.form.start_date) > new Date(props.form.end_date)) {
      localErrors.end_date = ['Tanggal selesai tidak boleh lebih awal dari tanggal mulai']
    }
  }

  if (props.form.allow_installment) {
    if (!props.form.max_installment || props.form.max_installment <= 0) {
      localErrors.max_installment = ['Maksimal cicilan wajib diisi dengan angka lebih dari 0']
    }
  } else {
    props.form.max_installment = null
  }

  if (Object.keys(localErrors).length > 0) {
    emit('local-validation-failed', localErrors)
    return
  }

  if (!(await checkUniqueRule())) return

  const hasBlockingError = Object.entries(props.errors || {}).some(([field, value]) => {
    if (field === '_general') return false
    return Array.isArray(value) ? value.some(Boolean) : !!value
  })
  if (hasBlockingError) return

  emit('save')
}

watch(() => props.startDateDisplay, (val) => {
  if (val && val.length === 10) clearFieldError('start_date')
})

watch(() => props.form.start_date, (val) => {
  if (val && val.includes('-')) {
    const [y, m, d] = val.split('T')[0].split('-')
    emit('update:startDateDisplay', `${d}/${m}/${y}`)
    clearFieldError('start_date')
  }
})

watch(() => props.endDateDisplay, (val) => {
  if (val && val.length === 10) clearFieldError('end_date')
})

watch(() => props.form.end_date, (val) => {
  if (val && val.includes('-')) {
    const [y, m, d] = val.split('T')[0].split('-')
    emit('update:endDateDisplay', `${d}/${m}/${y}`)
    clearFieldError('end_date')
  }
})

watch(() => props.form.bill_type_id, (newId) => {
  if (!newId) return
  const sel = props.billTypes?.find(b => Number(b.id) === Number(newId))
  if (sel && sel.default_amount !== undefined) {
    props.form.amount = sel.default_amount
  }
})
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
        <div class="relative bg-white w-full max-w-4xl rounded-[2.5rem] shadow-2xl overflow-hidden animate-scale-in flex flex-col max-h-[90vh] my-8">
          <!-- Header -->
          <div class="px-8 py-6 border-b border-slate-50 flex items-center justify-between bg-white shrink-0">
            <div>
              <h3 class="font-black text-slate-800 text-xl tracking-tight">
                {{ isEditing ? 'Ubah Aturan Tagihan' : 'Tambah Aturan Tagihan' }}
              </h3>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">Lengkapi informasi aturan tagihan siswa</p>
            </div>
            <button @click="emit('update:modelValue', false)" class="p-3 hover:bg-slate-50 text-slate-400 hover:text-slate-600 rounded-2xl transition-all cursor-pointer">
              <CloseIcon class="w-6 h-6" />
            </button>
          </div>

          <!-- Body -->
          <div class="p-8 space-y-8 overflow-y-auto custom-scrollbar bg-slate-50/20">
            <FormError :message="errors._general" />

            <!-- Info Banner -->
            <div class="p-5 bg-blue-50/40 border border-blue-200/60 rounded-2xl flex items-start gap-3 shadow-sm backdrop-blur-sm animate-fade-in">
              <InfoIcon class="w-5 h-5 text-blue-600 shrink-0 mt-0.5" />
              <div class="space-y-1">
                <h4 class="text-xs font-black text-blue-900 uppercase tracking-wider">Petunjuk Pembuatan Aturan Tagihan</h4>
                <p class="text-[11px] font-medium text-blue-800/90 leading-relaxed">
                  Modul ini berfungsi mengatur otomatisasi pembebanan biaya kepada siswa. Pilihan Kategori Tagihan di bawah ini secara otomatis hanya menampilkan jenis tagihan operasional yang berstatus <span class="font-bold underline">Aktif</span>. Pastikan target sasaran (jurusan/kelas) yang dipilih juga dalam kondisi aktif.
                </p>
              </div>
            </div>

            <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
              <!-- Left Column: Category & Target -->
              <div class="space-y-8">
                <!-- Card 1: Informasi Dasar -->
                <div class="bg-white p-6 rounded-[2rem] border border-slate-100 shadow-sm space-y-5">
                  <div class="flex items-center gap-3 mb-2">
                    <div class="w-8 h-8 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center shrink-0">
                      <CheckCircleIcon class="w-4 h-4" />
                    </div>
                    <h4 class="text-[11px] font-black text-slate-700 uppercase tracking-widest">Informasi Dasar</h4>
                  </div>

                  <!-- Kategori Tagihan (Dropdown) -->
                  <div class="space-y-1.5">
                    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1 block">
                      Kategori Tagihan <span class="text-rose-500">*</span>
                    </label>
                    <select 
                      v-model="form.bill_type_id" 
                      @change="handleBillTypeChange"
                      class="w-full py-4 px-6 bg-slate-50 border border-slate-100 rounded-2xl font-bold text-sm text-slate-700 shadow-sm outline-none focus:bg-white focus:ring-4 focus:ring-indigo-50 focus:border-indigo-500 transition-all cursor-pointer"
                      :class="{'!border-rose-500 !ring-rose-50': errors.bill_type_id}"
                    >
                      <option value="" disabled>-- Pilih Kategori Tagihan --</option>
                      <option v-for="bt in billTypes" :key="bt.id" :value="bt.id">{{ bt.name }} (Rp {{ bt.default_amount?.toLocaleString() }})</option>
                    </select>
                    <FormError :message="errors.bill_type_id" />
                  </div>

                  <!-- Nominal Tarif -->
                  <div class="space-y-1.5">
                    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1 block">
                      Nominal Tarif (Rp) <span class="text-rose-500">*</span>
                    </label>
                    <input 
                      type="text" 
                      v-model="formattedAmount" 
                      @blur="validateAmount"
                      placeholder="0" 
                      class="w-full py-4 px-6 bg-slate-50 border border-slate-100 rounded-2xl font-bold text-sm text-slate-700 shadow-sm outline-none focus:bg-white focus:ring-4 focus:ring-indigo-50 focus:border-indigo-500 transition-all"
                      :class="{'!border-rose-500 !ring-rose-50': errors.amount}" 
                    />
                    <FormError :message="errors.amount" />
                  </div>
                </div>

                <!-- Card 2: Cakupan Target Sasaran -->
                <div class="bg-white p-6 rounded-[2rem] border border-slate-100 shadow-sm space-y-5">
                  <div class="flex items-center justify-between mb-2">
                    <div class="flex items-center gap-3">
                      <div class="w-8 h-8 bg-emerald-50 text-emerald-600 rounded-xl flex items-center justify-center shrink-0">
                        <MajorIcon class="w-4 h-4" />
                      </div>
                      <h4 class="text-[11px] font-black text-slate-700 uppercase tracking-widest">Pilih Cakupan Target</h4>
                    </div>
                    <span class="text-[9px] font-bold text-emerald-600 bg-emerald-50 px-3 py-1 rounded-full uppercase tracking-tighter shadow-sm border border-emerald-100">
                      {{ form.target_type === 'all' ? 'Semua Siswa' : (form.target_ids?.length || 0) + ' Terpilih' }}
                    </span>
                  </div>

                  <div class="grid grid-cols-3 gap-2 pt-1">
                    <button v-for="t in ['all', 'major', 'class']" :key="t" @click="setTargetType(t)" 
                      type="button"
                      :class="['py-3 rounded-xl border-2 text-[10px] font-black uppercase tracking-wider transition-all cursor-pointer', form.target_type === t ? 'border-emerald-500 bg-emerald-50 text-emerald-700 shadow-sm' : 'border-slate-100 bg-slate-50 text-slate-400 hover:border-emerald-200']">
                      {{ t === 'all' ? 'Semua Siswa' : t === 'major' ? 'Per Jurusan' : 'Per Kelas' }}
                    </button>
                  </div>

                  <!-- Jika Semua Siswa -->
                  <div v-if="form.target_type === 'all'" class="p-5 bg-slate-50 rounded-2xl border border-slate-100 animate-fade-in text-center">
                    <p class="text-xs font-bold text-slate-500 leading-relaxed">
                      Aturan ini akan berlaku untuk seluruh siswa aktif di sekolah tanpa membedakan jurusan maupun kelas.
                    </p>
                  </div>

                  <!-- Jika Per Jurusan -->
                  <div v-if="form.target_type === 'major'" class="space-y-3 pt-1 animate-fade-in">
                    <div class="flex items-center justify-between px-1">
                      <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Daftar Jurusan Aktif</label>
                      <button v-if="!isEditing" @click="selectAllTargets" type="button" class="text-[10px] font-bold text-emerald-600 hover:text-emerald-700 cursor-pointer">Pilih Semua</button>
                    </div>
                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 max-h-52 overflow-y-auto pr-1 custom-scrollbar">
                      <div v-for="(m, idx) in majors" :key="m.id" 
                        @click="toggleTarget(m.id)"
                        class="p-4 rounded-2xl border-2 transition-all cursor-pointer flex items-center gap-3 group relative overflow-hidden"
                        :class="(isEditing ? form.target_id === m.id : form.target_ids.includes(m.id))
                          ? 'border-emerald-500 bg-emerald-50/50 shadow-sm' 
                          : 'border-slate-100 bg-slate-50 hover:border-emerald-200'">
                        <div class="w-8 h-8 rounded-lg flex items-center justify-center shrink-0 transition-all shadow-sm"
                          :class="(isEditing ? form.target_id === m.id : form.target_ids.includes(m.id)) ? 'bg-emerald-600 text-white' : 'bg-white text-slate-400 group-hover:text-emerald-500'">
                          <span class="text-[10px] font-black">{{ idx + 1 }}</span>
                        </div>
                        <div class="flex flex-col min-w-0">
                          <span class="text-[10px] font-black text-slate-700 uppercase tracking-tight truncate">{{ m.name }}</span>
                          <span class="text-[8px] font-bold text-slate-400 uppercase tracking-widest truncate">{{ m.code }}</span>
                        </div>
                        <div class="ml-auto">
                          <div class="w-4 h-4 rounded-full border-2 flex items-center justify-center transition-all"
                            :class="(isEditing ? form.target_id === m.id : form.target_ids.includes(m.id)) ? 'bg-emerald-600 border-emerald-600' : 'border-slate-300'">
                            <CheckCircleIcon v-if="(isEditing ? form.target_id === m.id : form.target_ids.includes(m.id))" class="w-2.5 h-2.5 text-white" />
                          </div>
                        </div>
                      </div>
                    </div>
                    <FormError :message="errors.target_id" />
                  </div>

                  <!-- Jika Per Kelas -->
                  <div v-if="form.target_type === 'class'" class="space-y-3 pt-1 animate-fade-in">
                    <div class="flex items-center justify-between px-1">
                      <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Daftar Kelas Aktif</label>
                      <button v-if="!isEditing" @click="selectAllTargets" type="button" class="text-[10px] font-bold text-emerald-600 hover:text-emerald-700 cursor-pointer">Pilih Semua</button>
                    </div>
                    <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 max-h-52 overflow-y-auto pr-1 custom-scrollbar">
                      <div v-for="(c, idx) in classes" :key="c.id" 
                        @click="toggleTarget(c.id)"
                        class="p-4 rounded-2xl border-2 transition-all cursor-pointer flex items-center gap-3 group relative overflow-hidden"
                        :class="(isEditing ? form.target_id === c.id : form.target_ids.includes(c.id))
                          ? 'border-emerald-500 bg-emerald-50/50 shadow-sm' 
                          : 'border-slate-100 bg-slate-50 hover:border-emerald-200'">
                        <div class="w-8 h-8 rounded-lg flex items-center justify-center shrink-0 transition-all shadow-sm"
                          :class="(isEditing ? form.target_id === c.id : form.target_ids.includes(c.id)) ? 'bg-emerald-600 text-white' : 'bg-white text-slate-400 group-hover:text-emerald-500'">
                          <span class="text-[10px] font-black">{{ idx + 1 }}</span>
                        </div>
                        <div class="flex flex-col min-w-0">
                          <span class="text-[10px] font-black text-slate-700 uppercase tracking-tight truncate">{{ c.name }}</span>
                        </div>
                        <div class="ml-auto">
                          <div class="w-4 h-4 rounded-full border-2 flex items-center justify-center transition-all"
                            :class="(isEditing ? form.target_id === c.id : form.target_ids.includes(c.id)) ? 'bg-emerald-600 border-emerald-600' : 'border-slate-300'">
                            <CheckCircleIcon v-if="(isEditing ? form.target_id === c.id : form.target_ids.includes(c.id))" class="w-2.5 h-2.5 text-white" />
                          </div>
                        </div>
                      </div>
                    </div>
                    <FormError :message="errors.target_id" />
                  </div>
                </div>
              </div>

              <!-- Right Column: Period, Dates & Installment -->
              <div class="space-y-8">
                <!-- Card 3: Tipe Siklus & Jatuh Tempo -->
                <div class="bg-white p-6 rounded-[2rem] border border-slate-100 shadow-sm space-y-5">
                  <div class="flex items-center gap-3 mb-2">
                    <div class="w-8 h-8 bg-amber-50 text-amber-600 rounded-xl flex items-center justify-center shrink-0">
                      <CalendarIcon class="w-4 h-4" />
                    </div>
                    <h4 class="text-[11px] font-black text-slate-700 uppercase tracking-widest">Tipe Siklus & Jatuh Tempo</h4>
                  </div>

                  <div class="flex gap-3">
                    <button @click="setPeriodType('bulanan')" type="button" :class="['flex-1 py-3.5 rounded-2xl border-2 font-black text-xs uppercase tracking-widest transition-all cursor-pointer', form.period_type === 'bulanan' ? 'border-amber-500 bg-amber-50 text-amber-700 shadow-sm' : 'border-slate-100 bg-slate-50 text-slate-400 hover:border-amber-200']">
                      Bulanan
                    </button>
                    <button @click="setPeriodType('tahunan')" type="button" :class="['flex-1 py-3.5 rounded-2xl border-2 font-black text-xs uppercase tracking-widest transition-all cursor-pointer', form.period_type === 'tahunan' ? 'border-amber-500 bg-amber-50 text-amber-700 shadow-sm' : 'border-slate-100 bg-slate-50 text-slate-400 hover:border-amber-200']">
                      Tahunan
                    </button>
                  </div>

                  <div v-if="form.period_type === 'bulanan'" class="space-y-1.5 pt-2 animate-fade-in">
                    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1 block">
                      Tanggal Jatuh Tempo Rutin (1-31)
                    </label>
                    <input v-model="form.due_day" @input="clearFieldError('due_day')" @blur="validateDueDate" type="number" min="1" max="31" class="w-full py-4 px-6 bg-slate-50 border border-slate-100 rounded-2xl font-bold text-sm text-slate-700 shadow-sm outline-none focus:bg-white focus:ring-4 focus:ring-amber-50 focus:border-amber-500 transition-all" :class="{'!border-rose-500 !ring-rose-50': errors.due_day}" />
                    <FormError :message="errors.due_day" />
                  </div>
                </div>

                <!-- Card 4: Masa Berlaku & Cicilan -->
                <div class="bg-white p-6 rounded-[2rem] border border-slate-100 shadow-sm space-y-5">
                  <div class="flex items-center gap-3 mb-2">
                    <div class="w-8 h-8 bg-purple-50 text-purple-600 rounded-xl flex items-center justify-center shrink-0">
                      <ClockIcon class="w-4 h-4" />
                    </div>
                    <h4 class="text-[11px] font-black text-slate-700 uppercase tracking-widest">Masa Berlaku & Cicilan</h4>
                  </div>

                  <!-- Start Date -->
                  <div class="space-y-1.5">
                    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1 flex items-center justify-between">
                      <span>Berlaku Mulai</span>
                      <span class="text-rose-500">*</span>
                    </label>
                    <div class="relative group">
                      <input 
                        v-model="localStartDateDisplay"
                        @input="clearFieldError('start_date')"
                        @blur="validateStartDate"
                        type="text" 
                        inputmode="numeric"
                        placeholder="HH/BB/TTTT"
                        class="w-full py-4 px-6 bg-slate-50 border border-slate-100 rounded-2xl font-bold text-sm text-slate-700 shadow-sm outline-none focus:bg-white focus:ring-4 focus:ring-purple-50 focus:border-purple-500 transition-all pr-12"
                        :class="errors.start_date ? '!border-red-500 !ring-red-50' : ''"
                      />
                      <div 
                        @click="triggerStartCalendar"
                        class="absolute right-0 top-0 h-full w-12 flex items-center justify-center group-focus-within:text-purple-600 text-slate-400 cursor-pointer overflow-hidden hover:bg-slate-100 transition-colors rounded-r-2xl"
                      >
                        <div class="w-px h-4 bg-slate-200 absolute left-0"></div>
                        <CalendarIcon class="w-4 h-4" />
                        <input 
                          ref="startDateRef"
                          type="date" 
                          v-model="form.start_date" 
                          class="absolute inset-0 opacity-0"
                          style="color-scheme: light;"
                        />
                      </div>
                    </div>
                    <FormError :message="errors.start_date" />
                  </div>

                  <!-- End Date -->
                  <div class="space-y-1.5 pt-1">
                    <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1 flex items-center justify-between">
                      <span>Berlaku Sampai</span>
                      <span class="text-rose-500">*</span>
                    </label>
                    <div class="relative group">
                      <input 
                        v-model="localEndDateDisplay"
                        @input="clearFieldError('end_date')"
                        @blur="validateEndDate"
                        type="text" 
                        inputmode="numeric"
                        placeholder="HH/BB/TTTT"
                        class="w-full py-4 px-6 bg-slate-50 border border-slate-100 rounded-2xl font-bold text-sm text-slate-700 shadow-sm outline-none focus:bg-white focus:ring-4 focus:ring-purple-50 focus:border-purple-500 transition-all pr-12"
                        :class="errors.end_date ? '!border-red-500 !ring-red-50' : ''"
                      />
                      <div 
                        @click="triggerEndCalendar"
                        class="absolute right-0 top-0 h-full w-12 flex items-center justify-center group-focus-within:text-purple-600 text-slate-400 cursor-pointer overflow-hidden hover:bg-slate-100 transition-colors rounded-r-2xl"
                      >
                        <div class="w-px h-4 bg-slate-200 absolute left-0"></div>
                        <CalendarIcon class="w-4 h-4" />
                        <input 
                          ref="endDateRef"
                          type="date" 
                          v-model="form.end_date" 
                          class="absolute inset-0 opacity-0"
                          style="color-scheme: light;"
                        />
                      </div>
                    </div>
                    <FormError :message="errors.end_date" />
                  </div>

                  <!-- Installment Settings -->
                  <div class="space-y-4 pt-4 border-t border-slate-100">
                    <div class="flex items-center justify-between px-1">
                      <label class="text-[11px] font-black text-slate-700 uppercase tracking-widest">Bisa Dicicil?</label>
                      <button @click="toggleInstallment" type="button" :class="['w-12 h-6 rounded-full transition-all relative cursor-pointer shadow-inner', form.allow_installment ? 'bg-purple-600' : 'bg-slate-300']">
                        <div :class="['absolute top-1 w-4 h-4 bg-white rounded-full transition-all shadow', form.allow_installment ? 'left-7' : 'left-1']"></div>
                      </button>
                    </div>

                    <div v-if="form.allow_installment" class="space-y-1.5 animate-fade-in pt-1">
                      <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1 flex items-center justify-between">
                        <span>Maksimal Cicilan</span>
                        <span class="text-rose-500">*</span>
                      </label>
                      <input v-model="form.max_installment" @input="clearFieldError('max_installment')" @blur="validateMaxInstallment" type="number" class="w-full py-4 px-6 bg-slate-50 border border-slate-100 rounded-2xl font-bold text-sm text-slate-700 shadow-sm outline-none focus:bg-white focus:ring-4 focus:ring-purple-50 focus:border-purple-500 transition-all" :class="errors.max_installment ? '!border-red-500 !ring-red-50' : ''" placeholder="Contoh: 3" />
                      <FormError :message="errors.max_installment" />
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Footer -->
          <div class="p-8 bg-slate-50/50 border-t border-slate-50 flex items-center gap-3 shrink-0">
            <button @click="emit('update:modelValue', false)" class="flex-1 py-4 bg-white hover:bg-slate-100 text-slate-600 font-black rounded-2xl transition-all text-[10px] uppercase tracking-widest border border-slate-200 shadow-sm cursor-pointer">
              Batal
            </button>
            <button @click="validateAndSave" :disabled="submitting" class="flex-1 py-4 bg-indigo-600 hover:bg-indigo-700 text-white font-black rounded-2xl transition-all text-[10px] uppercase tracking-widest shadow-xl shadow-indigo-100 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-3 cursor-pointer">
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
