<script setup>
import { 
  Filter as FilterIcon, 
  X as CloseIcon, 
  ChevronDown as ChevronDownIcon 
} from 'lucide-vue-next'

const props = defineProps({
  modelValue: Boolean,
  filters: Object
})

const emit = defineEmits(['update:modelValue', 'apply', 'reset'])

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
    <div v-if="modelValue" class="absolute top-full left-0 mt-2 w-full max-w-2xl bg-white border border-slate-200 rounded-2xl p-6 shadow-2xl shadow-slate-300/50 z-[100]">
      <div class="flex items-center justify-between mb-6 border-b border-slate-50 pb-4">
        <div class="flex items-center gap-3">
          <div class="p-2 bg-indigo-50 text-indigo-600 rounded-lg">
            <FilterIcon class="w-4 h-4" />
          </div>
          <div>
            <h4 class="font-black text-slate-700 text-xs uppercase tracking-widest">Filter Pengguna</h4>
            <p class="text-[10px] font-bold text-slate-400">Sesuaikan tampilan daftar akun</p>
          </div>
        </div>
        <button @click="emit('update:modelValue', false)" class="p-2 hover:bg-slate-100 text-slate-400 rounded-lg transition-colors">
          <CloseIcon class="w-5 h-5" />
        </button>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
        <!-- Role Filter -->
        <div class="flex flex-col gap-2">
          <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest pl-1">Role Akun</label>
          <div class="relative">
            <select v-model="filters.role" class="w-full py-2.5 px-4 bg-slate-50 border border-slate-100 rounded-xl appearance-none focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-xs text-slate-700 pr-8 shadow-sm">
              <option value="">Semua Role</option>
              <option value="admin">Administrator</option>
              <option value="parent">Wali Murid</option>
            </select>
            <div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
              <ChevronDownIcon class="w-3.5 h-3.5" />
            </div>
          </div>
        </div>

        <!-- Status Filter -->
        <div class="flex flex-col gap-2">
          <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest pl-1">Status Akun</label>
          <div class="relative">
            <select v-model="filters.status" class="w-full py-2.5 px-4 bg-slate-50 border border-slate-100 rounded-xl appearance-none focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-xs text-slate-700 pr-8 shadow-sm">
              <option value="">Semua Status</option>
              <option value="active">Aktif</option>
              <option value="inactive">Non-Aktif</option>
            </select>
            <div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-slate-400">
              <ChevronDownIcon class="w-3.5 h-3.5" />
            </div>
          </div>
        </div>

        <!-- Connection Filter -->
        <div class="flex flex-col gap-2">
          <label class="text-[9px] font-black text-slate-400 uppercase tracking-widest pl-1">Relasi Anak</label>
          <div class="relative">
            <select v-model="filters.filter" class="w-full py-2.5 px-4 bg-slate-50 border border-slate-100 rounded-xl appearance-none focus:bg-white focus:ring-2 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-xs text-slate-700 pr-8 shadow-sm">
              <option value="">Semua Relasi</option>
              <option value="no_child">Tanpa Anak</option>
              <option value="has_child">Punya Anak</option>
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
              <option value="">Terbaru (Default)</option>
              <option value="created_asc">Terlama</option>
              <option value="name_asc">Nama (A-Z)</option>
              <option value="name_desc">Nama (Z-A)</option>
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
