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

const statusLabel = (status) => {
  const labels = {
    paid: 'Lunas',
    partial: 'Sebagian',
    overdue: 'Menunggak',
    unpaid: 'Belum Lunas'
  }
  return labels[status] || 'Belum Lunas'
}

const statusClass = (status) => {
  const classes = {
    paid: 'bg-emerald-50 text-emerald-600 border-emerald-100',
    partial: 'bg-amber-50 text-amber-600 border-amber-100',
    overdue: 'bg-rose-50 text-rose-600 border-rose-100',
    unpaid: 'bg-slate-50 text-slate-600 border-slate-200'
  }
  return classes[status] || classes.unpaid
}
</script>

<template>
  <div class="bg-white rounded-xl border border-slate-200 shadow-sm overflow-hidden flex flex-col min-h-[700px]">
    <div class="p-4 border-b border-slate-100 bg-slate-50/30 flex items-center justify-between">
      <div class="flex items-center gap-2 font-black text-slate-700 text-xs uppercase tracking-widest">
        <div class="w-2 h-6 bg-indigo-500 rounded-full"></div>
        <span>Daftar Tagihan Siswa</span>
      </div>
    </div>

    <div class="flex-1 overflow-x-auto">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr class="bg-slate-50/50 border-b border-slate-100 text-[10px] font-black text-slate-400 uppercase tracking-widest">
            <th class="w-[45%] py-3 px-4">Siswa</th>
            <th class="w-[30%] py-3 px-4">Status</th>
            <th class="w-[25%] py-3 px-4 text-center">Detail</th>
          </tr>
        </thead>
        <tbody :class="{'opacity-50 pointer-events-none': loading}">
          <tr v-for="b in bills" :key="b.id" @click="$emit('select-student', b)"
            class="border-b border-slate-100 hover:bg-slate-50/30 cursor-pointer transition-all text-xs font-semibold text-slate-600"
            :class="selectedStudent?.student_id === b.student_id ? 'bg-indigo-50/50 ring-1 ring-inset ring-indigo-100' : ''">
            <td class="py-3 px-4">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 rounded-lg bg-white border border-slate-100 flex items-center justify-center text-slate-400 shadow-sm group-hover:border-indigo-200 transition-colors shrink-0">
                  <StudentIcon class="w-4 h-4" />
                </div>
                <div class="flex flex-col truncate">
                  <span class="font-black text-slate-800 text-xs uppercase tracking-wider truncate group-hover:text-indigo-600">{{ b.student_name }}</span>
                  <span class="text-[9px] font-black text-slate-400 uppercase tracking-widest mt-0.5">NIS #{{ b.student_id }}</span>
                </div>
              </div>
            </td>
            <td class="py-3 px-4">
              <div class="flex flex-col gap-1">
                <span :class="[
                  'inline-flex w-fit px-2 py-0.5 rounded text-[8px] font-black uppercase tracking-wider border',
                  statusClass(b.status)
                ]">
                  {{ statusLabel(b.status) }}
                </span>
                <div class="flex items-center gap-2">
                  <span class="text-[10px] font-black text-slate-800">{{ formatCurrency(b.outstanding ?? (b.amount - (b.total_paid || 0))) }}</span>
                  <span class="text-[8px] font-bold text-slate-400">({{ b.bill_count }} Item)</span>
                </div>
                <div v-if="b.overdue_count || b.partial_count" class="flex items-center gap-1.5">
                  <span v-if="b.overdue_count" class="text-[8px] font-black text-rose-500 uppercase">{{ b.overdue_count }} lewat tempo</span>
                  <span v-if="b.partial_count" class="text-[8px] font-black text-amber-500 uppercase">{{ b.partial_count }} cicilan</span>
                </div>
              </div>
            </td>
            <td class="py-3 px-4 text-center">
              <div class="p-1 bg-white rounded-lg border border-slate-100 group-hover:border-indigo-200 group-hover:text-indigo-600 transition-all inline-block shadow-sm">
                <PrevIcon class="w-3.5 h-3.5 rotate-180 text-slate-500" />
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div class="px-6 py-4 bg-slate-50/50 border-t border-slate-100 flex items-center justify-between">
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
          class="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 bg-white text-slate-400 hover:text-indigo-600 hover:border-indigo-100 hover:bg-indigo-50/30 disabled:opacity-20 disabled:hover:bg-white disabled:hover:border-slate-200 transition-all cursor-pointer"
        >
          <PrevIcon class="w-3.5 h-3.5" />
        </button>

        <!-- Page Numbers -->
        <div class="flex items-center gap-1">
          <button 
            v-for="p in visiblePages" 
            :key="p"
            @click="$emit('update:page', p)"
            class="w-8 h-8 flex items-center justify-center rounded-lg text-[10px] font-black transition-all cursor-pointer"
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
          class="w-8 h-8 flex items-center justify-center rounded-lg border border-slate-200 bg-white text-slate-400 hover:text-indigo-600 hover:border-indigo-100 hover:bg-indigo-50/30 disabled:opacity-20 disabled:hover:bg-white disabled:hover:border-slate-200 transition-all cursor-pointer"
        >
          <PrevIcon class="w-3.5 h-3.5 rotate-180" />
        </button>
      </div>
    </div>
  </div>
</template>
