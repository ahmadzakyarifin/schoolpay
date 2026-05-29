<script setup>
import { computed } from 'vue'
import { User as StudentIcon, ChevronLeft as PrevIcon } from 'lucide-vue-next'

const props = defineProps({
  bills: { type: Array, required: true },
  loading: { type: Boolean, default: false },
  selectedStudent: { type: Object, default: null },
  pagination: { type: Object, required: true }
})

defineEmits(['select-student', 'update:page', 'update:limit'])

const formatCurrency = (val) => {
  if (!val) return 'Rp 0'
  const clean = Number(val).toFixed(0)
  return 'Rp ' + clean.replace(/\B(?=(\d{3})+(?!\d))/g, '.')
}

const visiblePages = computed(() => {
  const pages = []
  let startPage = Math.max(1, props.pagination.page - 1)
  let endPage = Math.min(props.pagination.totalPages, startPage + 2)
  
  if (endPage - startPage < 2) {
    startPage = Math.max(1, endPage - 2)
  }
  
  for (let i = startPage; i <= endPage; i++) {
    pages.push(i)
  }
  return pages
})
</script>

<template>
  <div class="bg-white rounded-2xl border border-slate-200 shadow-sm overflow-hidden flex flex-col min-h-[700px]">
    <div class="px-8 py-8 border-b border-slate-100 bg-slate-50/30 flex items-center justify-between">
      <div class="flex items-center gap-4">
        <div class="w-2 h-6 bg-indigo-500 rounded-full"></div>
        <h3 class="font-black text-slate-700 text-sm uppercase tracking-[0.2em]">Daftar Tagihan Siswa</h3>
      </div>
    </div>

    <div class="flex-1 overflow-x-auto">
      <table class="w-full text-left border-collapse table-fixed">
        <thead>
          <tr class="bg-slate-50/50 border-b border-slate-100">
            <th class="w-[45%] px-6 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Siswa</th>
            <th class="w-[30%] px-6 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Status</th>
            <th class="w-[25%] px-6 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right pr-8">Detail</th>
          </tr>
        </thead>
        <tbody :class="{'opacity-50 pointer-events-none': loading}">
          <tr v-for="b in bills" :key="b.id" @click="$emit('select-student', b)"
            class="group hover:bg-indigo-50/30 cursor-pointer transition-all border-b border-slate-50 last:border-0"
            :class="selectedStudent?.student_id === b.student_id ? 'bg-indigo-50/50 ring-1 ring-inset ring-indigo-100' : ''">
            <td class="px-6 py-5 text-left">
              <div class="flex items-center gap-4">
                <div class="w-12 h-12 rounded-2xl bg-white border border-slate-100 flex items-center justify-center text-slate-400 shadow-sm group-hover:border-indigo-200 transition-colors">
                  <StudentIcon class="w-6 h-6" />
                </div>
                <div class="flex flex-col truncate">
                  <span class="font-bold text-slate-500 text-xs uppercase tracking-wider truncate group-hover:text-indigo-600">{{ b.student_name }}</span>
                  <span class="text-[9px] font-black text-slate-400 uppercase tracking-widest mt-0.5">NIS #{{ b.student_id }}</span>
                </div>
              </div>
            </td>
            <td class="px-6 py-5 text-left">
              <div class="flex flex-col gap-1">
                <span :class="[
                  'inline-flex w-fit px-2.5 py-1 rounded text-[8px] font-black uppercase tracking-widest border',
                  b.status === 'paid' ? 'bg-emerald-50 text-emerald-600 border-emerald-100' : 'bg-rose-50 text-rose-600 border-rose-100'
                ]">
                  {{ b.status === 'paid' ? 'Lunas' : 'Belum Lunas' }}
                </span>
                <div class="flex items-center gap-2">
                  <span class="text-[10px] font-black text-slate-800">{{ formatCurrency(b.amount - (b.total_paid || 0)) }}</span>
                  <span class="text-[8px] font-bold text-slate-400">({{ b.bill_count }} Item)</span>
                </div>
              </div>
            </td>
            <td class="px-6 py-5 text-right pr-8">
              <div class="p-2 bg-white rounded-xl border border-slate-100 group-hover:border-indigo-200 group-hover:text-indigo-600 transition-all inline-block shadow-sm">
                <PrevIcon class="w-4 h-4 rotate-180" />
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div class="px-8 py-6 bg-slate-50/50 border-t border-slate-100 flex items-center justify-between">
      <div class="flex items-center gap-6">
        <div class="flex items-center gap-3">
          <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Tampilkan</span>
          <select :value="pagination.limit" @change="$emit('update:limit', Number($event.target.value))" class="bg-white border border-slate-200 rounded-lg text-[10px] font-black text-slate-600 px-2 py-1 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 transition-all cursor-pointer shadow-sm">
            <option :value="10">10</option>
            <option :value="25">25</option>
            <option :value="50">50</option>
            <option :value="100">100</option>
          </select>
        </div>
        <div class="h-8 w-px bg-slate-200 hidden sm:block"></div>
        <span class="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em]">
          Halaman <span class="text-indigo-600">{{ pagination.page }}</span> dari {{ pagination.totalPages }} <span class="mx-2 text-slate-300">|</span> Total <span class="text-indigo-600">{{ pagination.total }}</span> Siswa
        </span>
      </div>
      <!-- Pagination Control -->
      <div class="flex items-center gap-2">
        <button 
          v-if="pagination.totalPages > 1"
          @click="pagination.page > 1 && $emit('update:page', pagination.page - 1)" 
          :disabled="pagination.page <= 1" 
          class="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-400 hover:text-indigo-600 hover:border-indigo-100 hover:bg-indigo-50/30 disabled:opacity-20 disabled:hover:bg-white disabled:hover:border-slate-200 transition-all cursor-pointer"
        >
          <PrevIcon class="w-4 h-4" />
        </button>

        <!-- Page Numbers -->
        <div class="flex items-center gap-1">
          <button 
            v-for="p in visiblePages" 
            :key="p"
            @click="$emit('update:page', p)"
            class="w-10 h-10 flex items-center justify-center rounded-xl text-[10px] font-black transition-all cursor-pointer"
            :class="p === pagination.page 
              ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/20' 
              : 'bg-white border border-slate-200 text-slate-500 hover:bg-slate-50 hover:border-slate-300'"
          >
            {{ p }}
          </button>
        </div>

        <button 
          v-if="pagination.totalPages > 1"
          @click="pagination.page < pagination.totalPages && $emit('update:page', pagination.page + 1)" 
          :disabled="pagination.page >= pagination.totalPages" 
          class="w-10 h-10 flex items-center justify-center rounded-xl border border-slate-200 bg-white text-slate-400 hover:text-indigo-600 hover:border-indigo-100 hover:bg-indigo-50/30 disabled:opacity-20 disabled:hover:bg-white disabled:hover:border-slate-200 transition-all cursor-pointer"
        >
          <PrevIcon class="w-4 h-4 rotate-180" />
        </button>
      </div>
    </div>
  </div>
</template>
