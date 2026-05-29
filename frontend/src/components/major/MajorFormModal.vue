<script setup>
import { computed } from 'vue'
import { 
  X as CloseIcon,
  CheckCircle2 as CheckCircleIcon,
  Loader2 as LoaderIcon,
  AlertTriangle as WarningIcon
} from 'lucide-vue-next'
import FormError from '../ui/FormError.vue'
import majorService from '../../services/major.service'

const props = defineProps({
  modelValue: Boolean,
  isEditing: Boolean,
  form: Object,
  errors: Object,
  submitting: Boolean,
  dependencyInfo: Object
})

const emit = defineEmits(['update:modelValue', 'save', 'clear-field-error', 'set-field-error'])

const clearFieldError = (field) => {
  emit('clear-field-error', field)
}

const setFieldError = (field, messages) => {
  emit('set-field-error', { field, messages })
}

const checkFieldUnique = async (field) => {
  let value = props.form[field]
  if (value === undefined || value === null || String(value).trim() === '') {
    if (field === 'name') {
      setFieldError('name', ['Nama jurusan wajib diisi'])
    } else {
      clearFieldError(field)
    }
    return
  }

  try {
    const excludeId = props.isEditing ? props.form.id : 0
    const res = await majorService.checkUnique(field, value, excludeId)
    if (res.data?.data?.is_unique === false) {
      const fieldLabel = field === 'name' ? 'Nama jurusan' : 'Kode jurusan'
      setFieldError(field, [`${fieldLabel} sudah terdaftar di sistem`])
    } else {
      clearFieldError(field)
    }
  } catch (err) {
    console.error(`Error checking uniqueness for ${field}:`, err)
  }
}

const editLocked = computed(() => props.isEditing && !!props.dependencyInfo?.has_dependencies)
const editLockMessage = computed(() => {
  if (!editLocked.value) return ''
  const message = props.dependencyInfo?.message || 'relasi aktif'
  if (String(message).toLowerCase().includes('tidak dapat diubah')) {
    return `${String(message).replace(/\.$/, '')}.`
  }
  return `Nama jurusan tidak dapat diubah karena masih terhubung dengan ${message}. Kode jurusan tetap bisa disesuaikan sesuai kebutuhan instansi.`
})

const handleSave = () => {
  let hasError = false
  if (!props.form.name?.trim()) {
    setFieldError('name', ['Nama jurusan wajib diisi'])
    hasError = true
  }
  // Check if there are active uniqueness errors
  if (props.errors.name || props.errors.code) {
    hasError = true
  }

  if (!hasError) {
    emit('save')
  }
}
</script>

<template>
  <Teleport to="body">
    <transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="modelValue" class="fixed inset-0 z-[1000] flex items-center justify-center p-4 overflow-y-auto">
        <div class="fixed inset-0 bg-slate-900/60 backdrop-blur-sm" @click="emit('update:modelValue', false)"></div>
        <div class="relative bg-white w-full max-w-md rounded-[2.5rem] shadow-2xl overflow-hidden animate-scale-in my-8">
          <div class="px-8 py-6 border-b border-slate-50 flex items-center justify-between bg-white">
            <div>
              <h3 class="font-black text-slate-800 text-xl tracking-tight">
                {{ isEditing ? 'Ubah Jurusan' : 'Tambah Jurusan' }}
              </h3>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">Lengkapi informasi program keahlian</p>
            </div>
            <button @click="emit('update:modelValue', false)" class="p-3 hover:bg-slate-50 text-slate-400 hover:text-slate-600 rounded-2xl transition-all">
              <CloseIcon class="w-6 h-6" />
            </button>
          </div>

          <div class="p-8 space-y-5">
            <FormError :message="errors._general" />

            <div v-if="editLocked" class="p-4 bg-amber-50 border border-amber-200 rounded-2xl text-left flex gap-3">
              <WarningIcon class="w-5 h-5 text-amber-600 shrink-0 mt-0.5" />
              <div>
                <h4 class="text-[11px] font-black text-amber-900 uppercase tracking-wider mb-1">Data Masih Terhubung</h4>
                <p class="text-[11px] font-bold text-amber-800 leading-relaxed">{{ editLockMessage }}</p>
                <p class="mt-2 text-[10px] font-bold text-amber-700/90 leading-relaxed">Jika memang perlu mengganti nama, buat jurusan baru atau ubah struktur kelas terkait secara terencana.</p>
              </div>
            </div>

            <div class="space-y-1.5">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1">Nama Jurusan <span class="text-rose-500">*</span></label>
              <input 
                v-model="form.name"
                type="text" 
                placeholder="Contoh: Rekayasa Perangkat Lunak" 
                :disabled="editLocked"
                class="w-full py-4 px-6 bg-slate-50 border border-slate-100 rounded-2xl focus:bg-white focus:ring-4 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-sm text-slate-700 shadow-sm disabled:cursor-not-allowed disabled:bg-slate-100 disabled:text-slate-400"
                :class="{'!border-rose-500 !ring-rose-50': errors.name}"
                @input="clearFieldError('name')"
                @blur="checkFieldUnique('name')"
              />
              <FormError :message="errors.name" />
            </div>

            <div class="space-y-1.5">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1">Kode Jurusan <span class="text-slate-300">(opsional)</span></label>
              <input 
                v-model="form.code"
                type="text" 
                placeholder="Bebas, contoh: RPL" 
                :disabled="submitting"
                class="w-full py-4 px-6 bg-slate-50 border border-slate-100 rounded-2xl focus:bg-white focus:ring-4 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-sm text-slate-700 shadow-sm disabled:cursor-not-allowed disabled:bg-slate-100 disabled:text-slate-400"
                :class="{'!border-rose-500 !ring-rose-50': errors.code}"
                @input="clearFieldError('code')"
                @blur="checkFieldUnique('code')"
              />
              <FormError :message="errors.code" />
            </div>
          </div>

          <div class="p-8 bg-slate-50/50 border-t border-slate-50 flex items-center gap-3">
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
</style>
