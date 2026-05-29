<script setup>
import { watch, ref, computed } from 'vue'
import axios from 'axios'
import { 
  X as CloseIcon,
  CheckCircle2 as CheckCircleIcon,
  Loader2 as LoaderIcon,
  Sparkles as MagicIcon
} from 'lucide-vue-next'
import FormError from '../ui/FormError.vue'
import classService from '../../services/class.service'

const props = defineProps({
  modelValue: Boolean,
  isEditing: Boolean,
  form: Object,
  errors: Object,
  submitting: Boolean,
  majors: Array
})

const emit = defineEmits(['update:modelValue', 'save', 'clear-field-error', 'set-field-error'])

const clearFieldError = (field) => {
  emit('clear-field-error', field)
}

const setFieldError = (field, messages) => {
  emit('set-field-error', { field, messages })
}

const checkUniqueName = async () => {
  if (!props.form.name?.trim()) {
    setFieldError('name', ['Nama kelas wajib diisi'])
    return
  }
  if (!props.form.major_id) {
    clearFieldError('name')
    return
  }
  
  try {
    const excludeId = props.isEditing ? props.form.id : 0
    const res = await classService.checkUnique(props.form.name, props.form.major_id, props.form.academic_year_id || 0, excludeId)
    if (res.data?.data?.is_unique === false) {
      setFieldError('name', ['Nama kelas sudah terdaftar di sistem untuk jurusan dan angkatan ini'])
    } else {
      clearFieldError('name')
    }
  } catch (err) {
    console.error('Error checking class uniqueness:', err)
  }
}

const isSuggesting = ref(false)
const nameEditedManually = ref(false)
const activeMajors = computed(() => props.majors.filter(m => m.is_active))

const toRoman = (num) => {
  const map = { M: 1000, CM: 900, D: 500, CD: 400, C: 100, XC: 90, L: 50, XL: 40, X: 10, IX: 9, V: 5, IV: 4, I: 1 }
  let roman = ''
  for (let i in map) {
    while (num >= map[i]) {
      roman += i
      num -= map[i]
    }
  }
  return roman
}

watch(() => props.modelValue, (isOpen) => {
  if (!isOpen) return
  nameEditedManually.value = false
})

// Auto Recommend Name
const recommendName = async (force = false) => {
  if (!props.form.grade || !props.form.major_id) return
  if (!force && nameEditedManually.value) return

  const major = activeMajors.value.find(m => m.id === props.form.major_id)
  if (!major) return

  isSuggesting.value = true
  try {
    // Format: [Grade Roman] [Major Name] [Number]
    const gradeRoman = toRoman(props.form.grade)
    const baseName = `${gradeRoman} ${major.name} `
    const res = await axios.get('academic/class/suggest-name', {
      params: {
        name: baseName,
        major_id: props.form.major_id,
        ay_id: props.form.academic_year_id || undefined,
        exclude_id: props.form.id || undefined
      }
    })
    props.form.name = res.data.data
    nameEditedManually.value = false
  } catch (err) {
    // Fallback logic
    const gradeRoman = toRoman(props.form.grade)
    props.form.name = `${gradeRoman} ${major.name} 1`
    nameEditedManually.value = false
  } finally {
    isSuggesting.value = false
  }
}

const handleSave = () => {
  let hasError = false
  if (!props.form.name?.trim()) {
    setFieldError('name', ['Nama kelas wajib diisi'])
    hasError = true
  }
  if (!props.form.major_id) {
    setFieldError('major_id', ['Jurusan wajib dipilih'])
    hasError = true
  }
  if (!props.form.grade && props.form.grade !== 0) {
    setFieldError('grade', ['Tingkat wajib dipilih'])
    hasError = true
  }

  if (props.errors.name || props.errors.major_id || props.errors.grade) {
    hasError = true
  }

  if (!hasError) {
    emit('save')
  }
}

const grades = [
  { val: 1, lab: 'Kelas I' }, { val: 2, lab: 'Kelas II' }, { val: 3, lab: 'Kelas III' },
  { val: 4, lab: 'Kelas IV' }, { val: 5, lab: 'Kelas V' }, { val: 6, lab: 'Kelas VI' },
  { val: 7, lab: 'Kelas VII' }, { val: 8, lab: 'Kelas VIII' }, { val: 9, lab: 'Kelas IX' },
  { val: 10, lab: 'Kelas X' }, { val: 11, lab: 'Kelas XI' }, { val: 12, lab: 'Kelas XII' }
]
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
                {{ isEditing ? 'Ubah Kelas' : 'Tambah Kelas' }}
              </h3>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest mt-1">Lengkapi informasi rombongan belajar</p>
            </div>
            <button @click="emit('update:modelValue', false)" class="p-3 hover:bg-slate-50 text-slate-400 hover:text-slate-600 rounded-2xl transition-all">
              <CloseIcon class="w-6 h-6" />
            </button>
          </div>

          <div class="p-8 space-y-5">
            <!-- Info Banner -->
            <div class="p-3.5 bg-indigo-50/50 border border-indigo-100 rounded-2xl text-indigo-600 text-[9px] font-bold uppercase tracking-widest flex items-center gap-2 shadow-sm">
              <span>💡 Info: Hanya menampilkan jurusan aktif</span>
            </div>
            <FormError :message="errors._general" />

            <!-- Grade & Major Row -->
            <div class="grid grid-cols-2 gap-4">
              <div class="space-y-1.5">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1">Tingkat <span class="text-rose-500">*</span></label>
                <select 
                  v-model="form.grade"
                  class="w-full py-4 px-6 bg-slate-50 border border-slate-100 rounded-2xl focus:bg-white focus:ring-4 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-sm text-slate-700 shadow-sm appearance-none cursor-pointer"
                  :class="{'!border-rose-500 !ring-rose-50': errors.grade}"
                  @blur="!form.grade && form.grade !== 0 ? setFieldError('grade', ['Tingkat wajib dipilih']) : clearFieldError('grade')"
                  @change="clearFieldError('grade'); recommendName(false)"
                >
                  <option :value="null" disabled>Pilih</option>
                  <option v-for="g in grades" :key="g.val" :value="g.val">{{ g.lab }}</option>
                </select>
                <FormError :message="errors.grade" />
              </div>

              <div class="space-y-1.5">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1">Jurusan <span class="text-rose-500">*</span></label>
                <select 
                  v-model="form.major_id"
                  class="w-full py-4 px-6 bg-slate-50 border border-slate-100 rounded-2xl focus:bg-white focus:ring-4 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-sm text-slate-700 shadow-sm appearance-none cursor-pointer"
                  :class="{'!border-rose-500 !ring-rose-50': errors.major_id}"
                  @blur="!form.major_id ? setFieldError('major_id', ['Jurusan wajib dipilih']) : clearFieldError('major_id')"
                  @change="clearFieldError('major_id'); recommendName(false)"
                >
                  <option :value="null" disabled>Pilih</option>
                  <option v-for="m in activeMajors" :key="m.id" :value="m.id">{{ m.name }}</option>
                </select>
                <FormError :message="errors.major_id" />
              </div>
            </div>

            <div class="space-y-1.5">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest pl-1 flex items-center justify-between">
                <span>Nama Kelas <span class="text-rose-500">*</span></span>
                <button v-if="form.grade && form.major_id" @click="recommendName(true)" 
                  class="text-[9px] font-black text-indigo-500 hover:text-indigo-700 flex items-center gap-1 uppercase transition-colors disabled:opacity-50"
                  :disabled="isSuggesting">
                  <MagicIcon v-if="!isSuggesting" class="w-3 h-3" />
                  <LoaderIcon v-else class="w-3 h-3 animate-spin" />
                  {{ isSuggesting ? '...' : 'Auto' }}
                </button>
              </label>
              <input 
                v-model="form.name"
                type="text" 
                placeholder="Contoh: X IPA 1" 
                class="w-full py-4 px-6 bg-slate-50 border border-slate-100 rounded-2xl focus:bg-white focus:ring-4 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-sm text-slate-700 shadow-sm"
                :class="{'!border-rose-500 !ring-rose-50': errors.name}"
                @input="clearFieldError('name'); nameEditedManually = true"
                @blur="checkUniqueName"
              />
              <FormError :message="errors.name" />
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
