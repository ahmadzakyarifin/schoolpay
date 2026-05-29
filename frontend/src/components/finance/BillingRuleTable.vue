<script setup>
import { 
  Calendar as CalendarIcon,
  Trash as TrashIcon,
  Edit as EditIcon,
  Check as CheckIcon,
  Play as PlayIcon,
  CreditCard as BillIcon,
  RotateCcw as ResetIcon
} from 'lucide-vue-next'

const props = defineProps({
  list: Array,
  loading: Boolean,
  selectedIds: Array,
  status: String
})

const emit = defineEmits([
  'edit', 
  'delete', 
  'restore',
  'toggle-status', 
  'generate-bills',
  'toggle-select-all', 
  'toggle-select-item'
])

const isAllSelected = () => {
  return props.list.length > 0 && props.selectedIds.length === props.list.length
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return new Intl.DateTimeFormat('id-ID', { day: 'numeric', month: 'short', year: 'numeric' }).format(date)
}

const formatCurrency = (val) => {
  if (!val) return 'Rp 0'
  const clean = Number(val).toFixed(0)
  return 'Rp ' + clean.replace(/\B(?=(\d{3})+(?!\d))/g, '.')
}
</script>

<template>
  <div class="overflow-x-auto">
    <table class="w-full border-collapse font-inter">
      <thead>
        <tr class="border-b border-slate-100 bg-slate-50/50">
          <th class="w-[50px] px-4 py-4">
            <div @click="emit('toggle-select-all')" class="w-4 h-4 rounded border-2 flex items-center justify-center cursor-pointer transition-all"
              :class="isAllSelected() ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300'">
              <CheckIcon v-if="isAllSelected()" class="w-2.5 h-2.5 text-white" />
            </div>
          </th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Nama Tagihan</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Target</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Periode</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Nominal</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Masa Aktif</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Status</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Status Distribusi</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right px-10">Aksi</th>
        </tr>
      </thead>
      <tbody :class="{'opacity-50 pointer-events-none': loading}">
        <tr v-for="item in list" :key="item.id" class="group hover:bg-slate-50/50 transition-colors border-b border-slate-50 last:border-0"
          :class="{'bg-indigo-50/30': selectedIds.includes(item.id)}">
          <td class="px-4 py-4">
            <div @click="emit('toggle-select-item', item.id)" 
              class="w-4 h-4 rounded border-2 flex items-center justify-center cursor-pointer transition-all"
              :class="selectedIds.includes(item.id) ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300 group-hover:border-indigo-300'">
              <CheckIcon v-if="selectedIds.includes(item.id)" class="w-2.5 h-2.5 text-white" />
            </div>
          </td>
          <td class="px-4 py-4 text-left">
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-lg bg-slate-100 flex items-center justify-center text-slate-500 border border-slate-200/60 shrink-0 shadow-sm">
                <CalendarIcon class="w-4 h-4" />
              </div>
              <div class="flex flex-col truncate">
                <span class="font-bold text-slate-500 text-xs uppercase tracking-wider truncate">{{ item.bill_type_name }}</span>
                <span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">Rule #{{ item.id }}</span>
              </div>
            </div>
          </td>
          <td class="px-4 py-4 text-left">
            <div class="flex flex-col gap-0.5">
              <span class="text-xs font-medium text-slate-600 uppercase tracking-wider">
                {{ item.target_type }}
              </span>
              <span v-if="item.class_name" class="text-[9px] font-bold text-slate-400 uppercase truncate block">
                {{ item.class_name }}
              </span>
            </div>
          </td>
          <td class="px-4 py-4 text-left">
            <span class="text-xs font-medium text-slate-600 uppercase tracking-wider block">
              {{ item.period_type }}
            </span>
            <p v-if="item.period_type === 'bulanan'" class="text-[9px] text-slate-400 font-bold uppercase mt-0.5">Tiap Tgl {{ item.due_day }}</p>
          </td>
          <td class="px-4 py-4 text-left font-bold text-slate-700 text-xs">
            {{ formatCurrency(item.amount) }}
          </td>
          <td class="px-4 py-4 text-left">
            <div v-if="item.start_date" class="text-[9px] font-bold text-slate-500 flex flex-col gap-0.5">
              <span class="uppercase tracking-widest">Mulai: {{ formatDate(item.start_date) }}</span>
              <span class="uppercase tracking-widest">Selesai: {{ formatDate(item.end_date) }}</span>
            </div>
            <span v-else class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Selamanya</span>
          </td>
          <td class="px-4 py-4 text-left">
            <button 
              @click.stop="emit('toggle-status', item)"
              class="relative w-8 h-4 rounded-full transition-all duration-300 focus:outline-none shadow-inner cursor-pointer"
              :class="item.is_active ? 'bg-indigo-600' : 'bg-slate-300'"
              :title="item.is_active ? 'Klik untuk Non-aktifkan' : 'Klik untuk Aktifkan'"
            >
              <div class="absolute top-0.5 left-0.5 w-3 h-3 bg-white rounded-full shadow transition-transform duration-300"
                :class="item.is_active ? 'translate-x-4' : 'translate-x-0'"></div>
            </button>
          </td>
          <td class="px-4 py-4 text-left">
            <span class="text-xs font-medium capitalize"
              :class="item.bill_count > 0 ? 'text-emerald-600 font-bold' : 'text-slate-500'">
              {{ item.bill_count > 0 ? `Terdistribusi (${item.bill_count})` : 'Belum Distribusi' }}
            </span>
          </td>
          <td class="px-4 py-4 text-right px-10">
            <div class="flex items-center justify-end gap-1 px-4">
              <template v-if="status === 'trash'">
                <button @click="emit('restore', item)" class="p-2 hover:bg-emerald-50 text-slate-400 hover:text-emerald-600 rounded-xl transition-all" title="Pulihkan">
                  <ResetIcon class="w-4 h-4" />
                </button>
              </template>
              <template v-else>
                <button @click="emit('generate-bills', item.id)" class="generate-btn" title="Generate Tagihan">
                  <PlayIcon class="w-3 h-3 fill-current" />
                  <span>Generate</span>
                </button>
                <button @click="emit('edit', item)" class="p-2 hover:bg-amber-50 text-slate-400 hover:text-amber-600 rounded-xl transition-all" title="Ubah">
                  <EditIcon class="w-4 h-4" />
                </button>
                <button @click="emit('delete', item)" class="p-2 hover:bg-rose-50 text-slate-400 hover:text-rose-600 rounded-xl transition-all" title="Hapus">
                  <TrashIcon class="w-4 h-4" />
                </button>
              </template>
            </div>
          </td>
        </tr>

        <!-- Empty State -->
        <tr v-if="list.length === 0 && !loading">
          <td colspan="9" class="py-32">
            <div class="flex flex-col items-center justify-center text-center animate-scale-in">
              <div class="w-20 h-20 bg-slate-50 rounded-[2.5rem] flex items-center justify-center mb-6 border border-slate-100 shadow-inner">
                <BillIcon class="w-8 h-8 text-slate-300" />
              </div>
              <h3 class="font-black text-slate-700 text-sm uppercase tracking-[0.2em] mb-2">Data Aturan Tidak Ditemukan</h3>
              <p class="text-slate-400 text-xs font-medium max-w-xs mx-auto">Belum ada master aturan tagihan yang terekam atau coba sesuaikan filter pencarian Anda.</p>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped lang="postcss">
.generate-btn { @apply bg-indigo-600 hover:bg-indigo-700 text-white px-3 py-2 rounded-xl flex items-center gap-1.5 font-black text-[9px] uppercase tracking-widest shadow-lg shadow-indigo-100 transition-all active:scale-95 cursor-pointer; }
</style>
