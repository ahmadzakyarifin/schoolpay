<script setup>
import { computed } from 'vue'
import { 
  X as CloseIcon, 
  AlertCircle as AlertIcon, 
  CheckCircle2 as CheckCircleIcon,
  HelpCircle as HelpIcon,
  Calendar as RecurringIcon,
  Zap as OneTimeIcon
} from 'lucide-vue-next'
import FormError from '../ui/FormError.vue'
import financeService from '../../services/finance.service'

const props = defineProps({
  modelValue: Boolean,
  isEditing: Boolean,
  form: Object,
  errors: Object,
  submitting: Boolean
})

const emit = defineEmits(['update:modelValue', 'save', 'local-validation-failed', 'clear-field-error'])

const clearFieldError = (field) => {
  emit('clear-field-error', field)
}

const setFieldError = (field, messages) => {
  emit('local-validation-failed', { [field]: messages })
}

const checkUniqueName = async () => {
  if (!props.form.name?.trim()) {
    setFieldError('name', ['Nama kategori tagihan wajib diisi'])
    return
  }
  
  try {
    const excludeId = props.isEditing ? props.form.id : 0
    const res = await financeService.checkUniqueBillType(props.form.name, excludeId)
    if (res.data?.data?.is_unique === false) {
      setFieldError('name', ['Nama kategori tagihan sudah terdaftar di sistem'])
    } else {
      clearFieldError('name')
    }
  } catch (err) {
    console.error('Error checking bill type uniqueness:', err)
  }
}

const validateDefaultAmount = () => {
  if (props.form.default_amount === null || props.form.default_amount === undefined || props.form.default_amount <= 0) {
    setFieldError('default_amount', ['Nominal dasar wajib diisi dengan angka lebih dari 0'])
  } else {
    clearFieldError('default_amount')
  }
}

const formattedAmount = computed({
  get() {
    if (props.form.default_amount === null || props.form.default_amount === undefined) return ''
    const parts = props.form.default_amount.toString().split('.')
    let num = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, '.')
    if (parts.length > 1) {
      num += ',' + parts[1]
    }
    return num
  },
  set(val) {
    clearFieldError('default_amount')
    if (!val) {
      props.form.default_amount = 0
      return
    }
    let clean = val.replace(/[^0-9,]/g, '')
    clean = clean.replace(',', '.')
    const parsed = parseFloat(clean)
    props.form.default_amount = isNaN(parsed) ? 0 : parsed
  }
})

const validateAndSave = () => {
  const localErrors = {}
  
  if (!props.form.name?.trim()) {
    localErrors.name = ['Nama jenis tagihan wajib diisi']
  } else if (props.form.name.trim().length < 2) {
    localErrors.name = ['Nama jenis tagihan minimal 2 karakter']
  }

  if (!props.form.type) {
    localErrors.type = ['Tipe tagihan wajib dipilih']
  }

  if (props.form.default_amount === null || props.form.default_amount === undefined || props.form.default_amount <= 0) {
    localErrors.default_amount = ['Nominal dasar wajib diisi dengan angka lebih dari 0']
  }

  if (Object.keys(localErrors).length > 0) {
    emit('local-validation-failed', localErrors)
    return
  }
  
  if (props.errors.name || props.errors.default_amount || props.errors.type) {
    return
  }

  emit('save')
}
</script>

<template>
  <Teleport to="body">
    <div v-if="modelValue" class="fixed inset-0 z-[500] flex items-center justify-center p-6">
      <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="emit('update:modelValue', false)"></div>
      <div class="white-card w-full max-w-xl relative z-10 overflow-hidden shadow-[0_20px_50px_rgba(0,0,0,0.2)] animate-scale-in !rounded-[2.5rem]">
        <div class="p-10 border-b border-slate-100 flex items-center justify-between bg-white">
          <div>
            <h2 class="text-2xl font-black text-slate-900 tracking-tight">{{ isEditing ? 'Ubah Jenis Tagihan' : 'Tambah Jenis Tagihan' }}</h2>
            <p class="text-slate-500 font-medium mt-1">Lengkapi informasi master data tagihan di bawah ini</p>
          </div>
          <button @click="emit('update:modelValue', false)" class="p-3 hover:bg-slate-50 rounded-2xl text-slate-400 transition-all">
            <CloseIcon class="w-6 h-6" />
          </button>
        </div>

        <div class="p-10 space-y-6 max-h-[70vh] overflow-y-auto scrollbar-thin scrollbar-indigo custom-scrollbar">
          <!-- General Error fallback -->
          <FormError :message="errors._general" />

          <!-- Name Input -->
          <div class="space-y-2">
            <label class="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em] px-1 flex items-center gap-1">
              Nama Kategori Tagihan <span class="text-red-500">*</span>
            </label>
            <input v-model="form.name" @input="clearFieldError('name')" @blur="checkUniqueName" :class="['modern-input !rounded-2xl', errors.name ? '!border-red-500 !ring-red-50' : '']" placeholder="Contoh: SPP Bulanan, Uang Gedung, Seragam" />
            <FormError :message="errors.name" />
          </div>

          <!-- Description Input with Help Text -->
          <div class="space-y-2">
            <div class="flex items-center justify-between px-1">
              <label class="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em]">
                Deskripsi <span class="text-[8px] font-black bg-slate-100 text-slate-400 px-1.5 py-0.5 rounded uppercase tracking-widest ml-1">Opsional</span>
              </label>
            </div>
            <textarea 
              v-model="form.description" 
              @input="clearFieldError('description')"
              :class="['modern-input !rounded-2xl min-h-[90px] !py-3', errors.description ? '!border-red-500 !ring-red-50' : '']" 
              placeholder="Tuliskan keterangan singkat kegunaan tagihan ini..."
            ></textarea>
            <FormError :message="errors.description" />
            <div class="flex items-start gap-2 px-1 text-[10px] text-slate-400 leading-relaxed font-medium">
              <HelpIcon class="w-3.5 h-3.5 text-indigo-500 shrink-0 mt-0.5" />
              <span>**Fungsi Deskripsi:** Digunakan sebagai catatan internal bagi admin atau bendahara untuk membedakan peruntukan jenis tagihan (misal: "Khusus untuk siswa baru tahun ajaran 2026/2027").</span>
            </div>
          </div>

          <!-- Type Selection with Explanations -->
          <div class="space-y-3">
            <label class="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em] px-1 flex items-center gap-1">
              Tipe Sifat Tagihan <span class="text-red-500">*</span>
            </label>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <!-- Recurring Option -->
              <div 
                @click="form.type = 'recurring'; clearFieldError('type')"
                :class="[
                  'p-5 rounded-2xl border-2 transition-all cursor-pointer flex flex-col gap-2',
                  form.type === 'recurring' ? 'border-indigo-600 bg-indigo-50/40 shadow-lg shadow-indigo-100/50' : 'border-slate-100 bg-slate-50 hover:border-slate-200 hover:bg-slate-50/80'
                ]"
              >
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-2.5">
                    <div :class="['p-2 rounded-xl', form.type === 'recurring' ? 'bg-indigo-600 text-white' : 'bg-slate-200 text-slate-500']">
                      <RecurringIcon class="w-4 h-4" />
                    </div>
                    <span :class="['font-black text-xs uppercase tracking-wider', form.type === 'recurring' ? 'text-indigo-900' : 'text-slate-700']">Rutin Bayar</span>
                  </div>
                  <div :class="['w-4 h-4 rounded-full border-2 flex items-center justify-center', form.type === 'recurring' ? 'border-indigo-600 bg-indigo-600' : 'border-slate-300']">
                    <div v-if="form.type === 'recurring'" class="w-1.5 h-1.5 bg-white rounded-full"></div>
                  </div>
                </div>
                <p class="text-[10px] text-slate-500 font-medium leading-relaxed mt-1">
                  Tagihan yang dibebankan secara berulang secara periodik (misal: **SPP Bulanan**, Iuran Ekstrakurikuler).
                </p>
              </div>

              <!-- One Time Option -->
              <div 
                @click="form.type = 'one_time'; clearFieldError('type')"
                :class="[
                  'p-5 rounded-2xl border-2 transition-all cursor-pointer flex flex-col gap-2',
                  form.type === 'one_time' ? 'border-amber-600 bg-amber-50/40 shadow-lg shadow-amber-100/50' : 'border-slate-100 bg-slate-50 hover:border-slate-200 hover:bg-slate-50/80'
                ]"
              >
                <div class="flex items-center justify-between">
                  <div class="flex items-center gap-2.5">
                    <div :class="['p-2 rounded-xl', form.type === 'one_time' ? 'bg-amber-600 text-white' : 'bg-slate-200 text-slate-500']">
                      <OneTimeIcon class="w-4 h-4" />
                    </div>
                    <span :class="['font-black text-xs uppercase tracking-wider', form.type === 'one_time' ? 'text-amber-900' : 'text-slate-700']">Sekali Bayar</span>
                  </div>
                  <div :class="['w-4 h-4 rounded-full border-2 flex items-center justify-center', form.type === 'one_time' ? 'border-amber-600 bg-amber-600' : 'border-slate-300']">
                    <div v-if="form.type === 'one_time'" class="w-1.5 h-1.5 bg-white rounded-full"></div>
                  </div>
                </div>
                <p class="text-[10px] text-slate-500 font-medium leading-relaxed mt-1">
                  Tagihan insidental yang hanya dibebankan satu kali selama masa sekolah (misal: **Uang Gedung**, **Seragam**, Buku).
                </p>
              </div>
            </div>
            <FormError :message="errors.type" />
          </div>

          <!-- Default Amount Input -->
          <div class="space-y-2">
            <label class="text-[11px] font-black text-slate-400 uppercase tracking-[0.2em] px-1 flex items-center gap-1">
              Nominal Dasar (Rp) <span class="text-red-500">*</span>
            </label>
            <input 
              type="text" 
              v-model="formattedAmount" 
              @blur="validateDefaultAmount"
              :class="['modern-input !rounded-2xl font-black text-slate-800', errors.default_amount ? '!border-red-500 !ring-red-50' : '']" 
              placeholder="0" 
            />
            <FormError :message="errors.default_amount" />
            <p class="text-[10px] text-slate-400 font-medium px-1">
              *Nominal ini akan menjadi acuan awal saat Anda membuat Aturan Tagihan (Billing Rules) untuk siswa/kelas.
            </p>
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
