<script setup>
import { computed, watch } from 'vue'
import { 
  Filter as FilterIcon, 
  X as CloseIcon, 
  ChevronDown as ChevronDownIcon 
} from 'lucide-vue-next'

const props = defineProps({
  modelValue: Boolean,
  filters: Object,
  academicFilters: Object
})

const emit = defineEmits(['update:modelValue', 'apply', 'reset'])

const filteredMajors = computed(() => {
  if (!props.filters.entry_year) return props.academicFilters.majors || []
  const entryYear = parseInt(props.filters.entry_year)
  const ay = props.academicFilters.years?.find(y => y.year === entryYear)
  if (!ay) return props.academicFilters.majors || []
  
  return props.academicFilters.majors?.filter(m => 
    m.year_ids?.includes(ay.id)
  ) || []
})

const filteredClasses = computed(() => {
  const classes = props.academicFilters.classes || []
  let result = classes

  if (props.filters.entry_year) {
    const entryYear = parseInt(props.filters.entry_year)
    const ay = props.academicFilters.years?.find(y => y.year === entryYear)
    if (ay) {
      result = result.filter(c => c.academic_year_ids?.includes(ay.id))
    }
  }

  if (props.filters.major_id) {
    result = result.filter(c => c.major_id === parseInt(props.filters.major_id))
  }

  return result
})

const onYearChange = () => {
  props.filters.major_id = ''
  props.filters.class_id = ''
}

const onMajorChange = () => {
  props.filters.class_id = ''
}

watch(() => props.filters.major_id, (newMajor) => {
  if (newMajor) {
    const majorIdNum = Number(newMajor)
    const selectedClass = props.academicFilters.classes?.find(c => c.id === Number(props.filters.class_id))
    if (selectedClass && selectedClass.major_id !== majorIdNum) {
      props.filters.class_id = ''
    }
  }
})

watch(() => props.filters.class_id, (newClass) => {
  if (newClass) {
    const selectedClass = props.academicFilters.classes?.find(c => c.id === Number(newClass))
    if (selectedClass && selectedClass.major_id) {
      props.filters.major_id = String(selectedClass.major_id)
    }
  }
})



</script>

<template>
  <transition
    enter-active-class="transition duration-300 ease-out"
    enter-from-class="transform -translate-y-4 opacity-0 scale-95"
    enter-to-class="transform translate-y-0 opacity-100 scale-100"
    leave-active-class="transition duration-200 ease-in"
    leave-from-class="transform translate-y-0 opacity-100 scale-100"
    leave-to-class="transform -translate-y-4 opacity-0 scale-95"
  >
    <div v-if="modelValue" class="absolute top-full left-0 mt-2 w-full max-w-5xl bg-white border border-slate-200 rounded-2xl p-6 shadow-2xl shadow-slate-300/50 z-[100]">
      <div class="flex items-center justify-between mb-6 border-b border-slate-50 pb-4">
        <div class="flex items-center gap-3">
          <div class="p-2 bg-indigo-50 text-indigo-600 rounded-lg">
            <FilterIcon class="w-4 h-4" />
          </div>
          <div>
            <h4 class="font-black text-slate-700 text-xs uppercase tracking-widest">Filter Akademik</h4>
            <p class="text-[10px] font-bold text-slate-400">Saring daftar siswa berdasarkan kriteria</p>
          </div>
        </div>
        <button @click="emit('update:modelValue', false)" class="p-2 hover:bg-slate-100 text-slate-400 rounded-lg transition-colors">
          <CloseIcon class="w-5 h-5" />
        </button>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-5 gap-6">
        <!-- Status Filter -->
        <div class="flex flex-col gap-2">
          <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest pl-1">Status</label>
          <div class="relative">
            <select v-model="filters.status" class="w-full py-2.5 px-4 bg-slate-50 border border-slate-100 rounded-xl appearance-none focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-xs text-slate-700 pr-8 shadow-sm">
              <option value="">Semua Status</option>
              <option value="active">Aktif</option>
              <option value="inactive">Keluar</option>
              <option value="graduated">Lulus</option>
            </select>
            <div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
              <ChevronDownIcon class="w-3.5 h-3.5" />
            </div>
          </div>
        </div>

        <!-- Angkatan Filter -->
        <div class="flex flex-col gap-2">
          <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest pl-1">Angkatan</label>
          <div class="relative">
            <select v-model="filters.entry_year" @change="onYearChange" class="w-full py-2.5 px-4 bg-slate-50 border border-slate-100 rounded-xl appearance-none focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-xs text-slate-700 pr-8 shadow-sm">
              <option value="">Semua Angkatan</option>
              <option v-for="y in academicFilters.years" :key="y.id" :value="y.year">{{ y.year }}</option>
            </select>
            <div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
              <ChevronDownIcon class="w-3.5 h-3.5" />
            </div>
          </div>
        </div>

        <!-- Jurusan Filter -->
        <div class="flex flex-col gap-2">
          <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest pl-1">Jurusan</label>
          <div class="relative">
            <select v-model="filters.major_id" @change="onMajorChange" :disabled="!!filters.class_id" class="w-full py-2.5 px-4 bg-slate-50 border border-slate-100 rounded-xl appearance-none focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-xs text-slate-700 pr-8 shadow-sm disabled:opacity-60 disabled:cursor-not-allowed">
              <option value="">Semua Jurusan</option>
              <option v-for="j in filteredMajors" :key="j.id" :value="j.id">{{ j.name }}</option>
            </select>
            <div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
              <ChevronDownIcon class="w-3.5 h-3.5" />
            </div>
          </div>
        </div>

        <!-- Kelas Filter -->
        <div class="flex flex-col gap-2">
          <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest pl-1">Kelas</label>
          <div class="relative">
            <select v-model="filters.class_id" class="w-full py-2.5 px-4 bg-slate-50 border border-slate-100 rounded-xl appearance-none focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-xs text-slate-700 pr-8 shadow-sm">
              <option value="">Semua Kelas</option>
              <option v-for="c in filteredClasses" :key="c.id" :value="c.id">{{ c.name }}</option>
            </select>
            <div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
              <ChevronDownIcon class="w-3.5 h-3.5" />
            </div>
          </div>
        </div>

        <!-- Sort Filter -->
        <div class="flex flex-col gap-2">
          <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest pl-1">Pengurutan</label>
          <div class="relative">
            <select v-model="filters.sort" class="w-full py-2.5 px-4 bg-slate-50 border border-slate-100 rounded-xl appearance-none focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-xs text-slate-700 pr-8 shadow-sm">
              <option value="">Nama (A-Z)</option>
              <option value="name_desc">Nama (Z-A)</option>
              <option value="created_desc">Terbaru</option>
              <option value="created_asc">Terlama</option>
              <option value="entry_year_desc">Angkatan Terbaru</option>
              <option value="entry_year_asc">Angkatan Terlama</option>
            </select>
            <div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
              <ChevronDownIcon class="w-3.5 h-3.5" />
            </div>
          </div>
        </div>
      </div>

      <!-- Apply Button Area -->
      <div class="mt-8 pt-4 border-t border-slate-50 flex justify-end gap-3">
        <button @click="emit('reset')" class="px-6 py-2.5 text-xs font-black text-slate-400 uppercase tracking-widest hover:text-slate-600 transition-colors">
          Reset
        </button>
        <button @click="emit('apply')" class="px-8 py-2.5 bg-indigo-600 hover:bg-indigo-700 text-white text-[10px] font-black uppercase tracking-widest rounded-xl shadow-lg shadow-indigo-100 transition-all">
          Terapkan Filter
        </button>
      </div>
    </div>
  </transition>
</template>
