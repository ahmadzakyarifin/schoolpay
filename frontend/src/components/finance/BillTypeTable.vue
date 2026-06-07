<script setup>
import { 
  CreditCard as BillIcon,
  Trash as TrashIcon,
  Edit as EditIcon,
  Check as CheckIcon,
  Link as LinkIcon,
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
  'toggle-select-all', 
  'toggle-select-item'
])

const isAllSelected = () => {
  return props.list.length > 0 && props.selectedIds.length === props.list.length
}

const formatCurrency = (val) => {
  if (val === null || val === undefined) return 'Rp.0,00'
  return new Intl.NumberFormat('id-ID', { style: 'currency', currency: 'IDR' }).format(val).replace('Rp', 'Rp.').replace(/\s+/g, '')
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return d.toLocaleDateString('id-ID', { day: '2-digit', month: 'short', year: 'numeric' })
}
</script>

<template>
  <div class="flex-1 overflow-x-auto custom-scrollbar">
    <table class="w-full text-left border-collapse">
      <thead>
        <tr class="bg-slate-50/50 border-b border-slate-100 text-[10px] font-black text-slate-400 uppercase tracking-widest">
          <th class="w-12 py-3 px-4 text-center">
            <div @click="emit('toggle-select-all')" class="w-4 h-4 mx-auto rounded border-2 flex items-center justify-center cursor-pointer transition-all"
              :class="isAllSelected() ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300'">
              <CheckIcon v-if="isAllSelected()" class="w-2.5 h-2.5 text-white" />
            </div>
          </th>
          <th class="py-3 px-4">Kategori & Ikon</th>
          <th class="py-3 px-4">Tipe</th>
          <th class="py-3 px-4">Biaya Dasar</th>
          <th class="py-3 px-4">Keterhubungan</th>
          <th class="py-3 px-4">Tanggal Dibuat</th>
          <th class="py-3 px-4">Status</th>
          <th class="py-3 px-4 text-center w-[200px]">Aksi</th>
        </tr>
      </thead>
      <tbody :class="{'opacity-50 pointer-events-none': loading}">
        <tr v-for="item in list" :key="item.id" class="border-b border-slate-100 hover:bg-slate-50/30 transition-all text-xs font-semibold text-slate-600"
          :class="{'bg-indigo-50/30': selectedIds.includes(item.id)}">
          <td class="py-3 px-4 text-center">
            <div @click="emit('toggle-select-item', item.id)" 
              class="w-4 h-4 mx-auto rounded border-2 flex items-center justify-center cursor-pointer transition-all"
              :class="selectedIds.includes(item.id) ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300 hover:border-indigo-300'">
              <CheckIcon v-if="selectedIds.includes(item.id)" class="w-2.5 h-2.5 text-white" />
            </div>
          </td>
          <td class="py-3 px-4">
            <div class="flex flex-col truncate">
              <span class="font-black text-slate-800 text-xs uppercase tracking-wider truncate">{{ item.name }}</span>
              <span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">ID #{{ item.id }}</span>
              <span v-if="item.description" class="text-[10px] font-medium text-slate-500 truncate mt-0.5" :title="item.description">{{ item.description }}</span>
            </div>
          </td>
          <td class="py-3 px-4">
            <span class="bg-slate-100 text-slate-700 px-2 py-0.5 rounded text-[10px] font-mono font-bold uppercase tracking-wider">
              {{ item.type === 'recurring' ? 'Rutin' : 'Sekali' }}
            </span>
          </td>
          <td class="py-3 px-4 font-black text-slate-800 text-xs">
            {{ formatCurrency(item.default_amount) }}
          </td>
          <td class="py-3 px-4">
            <span class="text-xs font-bold capitalize"
              :class="item.rule_count > 0 ? 'text-indigo-600' : 'text-slate-500'">
              {{ item.rule_count > 0 ? `${item.rule_count} Aturan Aktif` : 'Belum Dipakai' }}
            </span>
          </td>
          <td class="py-3 px-4 text-slate-500 text-[11px]">
            <span>{{ formatDate(item.created_at) }}</span>
          </td>
          <td class="py-3 px-4">
            <div class="flex items-center">
              <button 
                @click.stop="emit('toggle-status', item)"
                class="relative w-8 h-4 rounded-full transition-all duration-300 focus:outline-none shadow-inner cursor-pointer"
                :class="item.is_active ? 'bg-indigo-600' : 'bg-slate-300'"
              >
                <div class="absolute top-0.5 left-0.5 w-3 h-3 bg-white rounded-full shadow transition-transform duration-300"
                  :class="item.is_active ? 'translate-x-4' : 'translate-x-0'"></div>
              </button>
            </div>
          </td>
          <td class="py-3 px-4 text-center">
            <div class="flex items-center justify-center gap-1.5 flex-nowrap">
              <template v-if="status === 'trash'">
                <button @click="emit('restore', item)" title="Pulihkan" class="p-2 bg-white text-emerald-600 border border-slate-200 hover:bg-emerald-50 rounded-lg flex items-center justify-center transition-all shadow-sm">
                  <ResetIcon class="w-3.5 h-3.5 text-emerald-600" />
                </button>
              </template>
              <template v-else>
                <button @click="emit('edit', item)" title="Ubah" class="p-2 bg-white text-slate-600 border border-slate-200 hover:bg-slate-50 rounded-lg flex items-center justify-center transition-all shadow-sm">
                  <EditIcon class="w-3.5 h-3.5 text-slate-500" />
                </button>
                <button @click="emit('delete', item)" title="Hapus" class="p-2 bg-white text-rose-600 border border-slate-200 hover:bg-rose-50 rounded-lg flex items-center justify-center transition-all shadow-sm">
                  <TrashIcon class="w-3.5 h-3.5 text-rose-500" />
                </button>
              </template>
            </div>
          </td>
        </tr>

        <!-- Empty State -->
        <tr v-if="list.length === 0 && !loading">
          <td colspan="8" class="py-24">
            <div class="flex flex-col items-center justify-center text-center animate-scale-in">
              <div class="w-20 h-20 bg-slate-50 rounded-[2.5rem] flex items-center justify-center mb-6 border border-slate-100 shadow-inner">
                <BillIcon class="w-8 h-8 text-slate-300" />
              </div>
              <h3 class="font-black text-slate-700 text-sm uppercase tracking-[0.2em] mb-2">Data Kategori Tidak Ditemukan</h3>
              <p class="text-slate-400 text-xs font-medium max-w-xs mx-auto">Belum ada master data tagihan yang terekam atau coba sesuaikan filter pencarian Anda.</p>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
