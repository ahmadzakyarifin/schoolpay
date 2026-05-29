<script setup>
import { Filter as FilterIcon, X as CloseIcon, ChevronDown as ChevronDownIcon, Calendar as CalendarIcon, Clock as ClockIcon, GraduationCap as StudentIcon } from 'lucide-vue-next'

const props = defineProps({
  modelValue: Boolean,
  filters: Object,
  academicYears: Array,
  classes: Array,
  majors: Array,
  billTypes: Array,
  hideBillType: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'apply', 'reset'])
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
      <div v-if="modelValue" class="fixed inset-0 z-[999] flex items-center justify-center p-4 sm:p-6 font-inter">
        <!-- Backdrop (blocks the whole screen) -->
        <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="emit('update:modelValue', false)"></div>
        
        <!-- Modal Card Dialog -->
        <div class="bg-white w-full max-w-3xl relative z-10 rounded-[2rem] p-8 shadow-2xl border border-slate-100 max-h-[90vh] overflow-y-auto transform transition-all animate-scale-in">
          
          <!-- Header -->
          <div class="flex items-center justify-between mb-8 border-b border-slate-50 pb-4">
            <div class="flex items-center gap-3">
              <div class="p-2.5 bg-indigo-50 text-indigo-600 rounded-xl">
                <FilterIcon class="w-5 h-5" />
              </div>
              <div>
                <h4 class="font-black text-slate-700 text-sm uppercase tracking-widest">Filter Laporan & Analitik</h4>
                <p class="text-[10px] font-bold text-slate-400">Sesuaikan rentang waktu dan demografi data laporan</p>
              </div>
            </div>
            <button @click="emit('update:modelValue', false)" class="p-2 hover:bg-slate-100 text-slate-400 rounded-xl transition-colors cursor-pointer">
              <CloseIcon class="w-6 h-6" />
            </button>
          </div>

          <!-- Body Fields -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <!-- Periode Filter -->
            <div class="flex flex-col gap-2">
              <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest pl-1 flex items-center gap-1.5">
                <ClockIcon class="w-3 h-3 text-indigo-500" /> Rentang Periode
              </label>
              <div class="relative">
                <select v-model="filters.period" class="w-full py-3 px-4 bg-slate-50 border border-slate-100 rounded-xl appearance-none focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-50 text-xs font-bold text-slate-700 pr-8 shadow-sm cursor-pointer outline-none transition-all">
                  <option value="all">Semua Periode</option>
                  <option value="daily">Harian</option>
                  <option value="monthly">Bulanan</option>
                  <option value="yearly">Tahunan</option>
                  <option value="custom">Kustom (Rentang Tanggal)</option>
                </select>
                <div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
                  <ChevronDownIcon class="w-3.5 h-3.5" />
                </div>
              </div>
            </div>

            <!-- Reference Date / Custom Date Filter -->
            <div class="flex flex-col gap-2 md:col-span-2">
              <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest pl-1 flex items-center gap-1.5">
                <CalendarIcon class="w-3 h-3 text-indigo-500" /> Tanggal / Bulan / Tahun Acuan
              </label>
              <div v-if="filters.period === 'all'" class="relative">
                <input 
                  type="text" 
                  disabled 
                  value="Semua Periode (Tanpa Acuan Tanggal)" 
                  class="w-full py-3 px-4 bg-slate-100 border border-slate-200 rounded-xl text-xs font-bold text-slate-400 shadow-sm outline-none cursor-not-allowed" 
                />
              </div>
              <div v-else-if="filters.period !== 'custom'" class="relative">
                <input 
                  v-model="filters.ref_date" 
                  :type="filters.period === 'monthly' ? 'month' : (filters.period === 'daily' ? 'date' : 'number')" 
                  class="w-full py-3 px-4 bg-slate-50 border border-slate-100 rounded-xl focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-50 text-xs font-bold text-slate-700 shadow-sm outline-none transition-all" 
                  :placeholder="filters.period === 'yearly' ? 'Contoh: 2026' : ''"
                />
              </div>
              <div v-else class="grid grid-cols-2 gap-4">
                <div class="relative">
                  <span class="absolute left-3 top-[-8px] bg-white px-1 text-[8px] font-black text-slate-400 uppercase tracking-widest">Mulai</span>
                  <input v-model="filters.start_date" type="date" class="w-full py-3 px-4 bg-slate-50 border border-slate-100 rounded-xl focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-50 text-xs font-bold text-slate-700 shadow-sm outline-none transition-all" />
                </div>
                <div class="relative">
                  <span class="absolute left-3 top-[-8px] bg-white px-1 text-[8px] font-black text-slate-400 uppercase tracking-widest">Selesai</span>
                  <input v-model="filters.end_date" type="date" class="w-full py-3 px-4 bg-slate-50 border border-slate-100 rounded-xl focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-50 text-xs font-bold text-slate-700 shadow-sm outline-none transition-all" />
                </div>
              </div>
            </div>

            <!-- Angkatan Filter -->
            <div class="flex flex-col gap-2">
              <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest pl-1 flex items-center gap-1.5">
                <StudentIcon class="w-3 h-3 text-indigo-500" /> Angkatan (Tahun Masuk)
              </label>
              <div class="relative">
                <select v-model="filters.academic_year_id" class="w-full py-3 px-4 bg-slate-50 border border-slate-100 rounded-xl appearance-none focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-50 text-xs font-bold text-slate-700 pr-8 shadow-sm cursor-pointer outline-none transition-all">
                  <option value="">Semua Angkatan</option>
                  <option v-for="y in academicYears" :key="y.id" :value="y.id">{{ y.year }}</option>
                </select>
                <div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
                  <ChevronDownIcon class="w-3.5 h-3.5" />
                </div>
              </div>
            </div>

            <!-- Major Filter -->
            <div class="flex flex-col gap-2">
              <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest pl-1">Jurusan / Major</label>
              <div class="relative">
                <select v-model="filters.major_id" class="w-full py-3 px-4 bg-slate-50 border border-slate-100 rounded-xl appearance-none focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-50 text-xs font-bold text-slate-700 pr-8 shadow-sm cursor-pointer outline-none transition-all">
                  <option value="">Semua Jurusan</option>
                  <option v-for="m in majors" :key="m.id" :value="m.id">{{ m.name }}</option>
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
                <select v-model="filters.class_id" class="w-full py-3 px-4 bg-slate-50 border border-slate-100 rounded-xl appearance-none focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-50 text-xs font-bold text-slate-700 pr-8 shadow-sm cursor-pointer outline-none transition-all">
                  <option value="">Semua Kelas</option>
                  <option v-for="c in classes" :key="c.id" :value="c.id">{{ c.name }}</option>
                </select>
                <div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
                  <ChevronDownIcon class="w-3.5 h-3.5" />
                </div>
              </div>
            </div>

            <!-- Jenis Tagihan Filter -->
            <div v-if="!hideBillType" class="flex flex-col gap-2 md:col-span-3">
              <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest pl-1">Jenis Tagihan</label>
              <div class="relative">
                <select v-model="filters.bill_type_id" class="w-full py-3 px-4 bg-slate-50 border border-slate-100 rounded-xl appearance-none focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-50 text-xs font-bold text-slate-700 pr-8 shadow-sm cursor-pointer outline-none transition-all">
                  <option value="">Semua Jenis Tagihan</option>
                  <option v-for="b in billTypes" :key="b.id" :value="b.id">{{ b.name }}</option>
                </select>
                <div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
                  <ChevronDownIcon class="w-3.5 h-3.5" />
                </div>
              </div>
            </div>
          </div>

          <!-- Apply Button Area -->
          <div class="mt-8 pt-6 border-t border-slate-50 flex justify-end gap-3">
            <button @click="emit('reset')" class="px-6 py-3 text-xs font-black text-slate-400 uppercase tracking-widest hover:text-slate-600 transition-colors cursor-pointer">
              Reset Filter
            </button>
            <button @click="emit('apply')" class="px-8 py-3 bg-indigo-600 hover:bg-indigo-700 text-white text-[11px] font-black uppercase tracking-widest rounded-xl shadow-lg shadow-indigo-100 transition-all cursor-pointer">
              Terapkan Filter
            </button>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<style scoped>
@keyframes scale-in {
  from { transform: scale(0.95); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}
.animate-scale-in {
  animation: scale-in 0.25s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}
</style>
