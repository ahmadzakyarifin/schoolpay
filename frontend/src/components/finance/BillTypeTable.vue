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
    <table class="w-full text-left border-collapse table-fixed">
      <thead>
        <tr class="bg-slate-50/50 border-b border-slate-100">
          <th class="w-[50px] px-4 py-4">
            <div @click="emit('toggle-select-all')" class="w-4 h-4 rounded border-2 flex items-center justify-center cursor-pointer transition-all"
              :class="isAllSelected() ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300'">
              <CheckIcon v-if="isAllSelected()" class="w-2.5 h-2.5 text-white" />
            </div>
          </th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Kategori & Ikon</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Tipe</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Biaya Dasar</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Keterhubungan</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Tanggal Dibuat</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Status</th>
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
            <div class="flex flex-col truncate">
              <span class="font-bold text-slate-500 text-xs uppercase tracking-wider truncate">{{ item.name }}</span>
              <span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">ID #{{ item.id }}</span>
              <span v-if="item.description" class="text-[10px] font-medium text-slate-500 truncate mt-0.5" :title="item.description">{{ item.description }}</span>
            </div>
          </td>
          <td class="px-4 py-4 text-left">
            <span class="text-xs font-medium text-slate-600 uppercase tracking-wider block">
              {{ item.type === 'recurring' ? 'Rutin Bayar' : 'Sekali Bayar' }}
            </span>
          </td>
          <td class="px-4 py-4 text-left font-bold text-slate-700 text-xs">
            {{ formatCurrency(item.default_amount) }}
          </td>
          <td class="px-4 py-4 text-left">
            <span class="text-xs font-medium capitalize"
              :class="item.rule_count > 0 ? 'text-indigo-600 font-bold' : 'text-slate-500'">
              {{ item.rule_count > 0 ? `${item.rule_count} Aturan Aktif` : 'Belum Dipakai' }}
            </span>
          </td>
          <td class="px-4 py-4 text-left">
            <span class="text-xs font-medium text-slate-600">{{ formatDate(item.created_at) }}</span>
          </td>
          <td class="px-4 py-4 text-left">
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
          <td class="px-4 py-4 text-right px-10">
            <div class="flex items-center justify-end gap-1 px-4">
              <template v-if="status === 'trash'">
                <button @click="emit('restore', item)" class="p-2 hover:bg-emerald-50 text-slate-400 hover:text-emerald-600 rounded-xl transition-all" title="Pulihkan">
                  <ResetIcon class="w-4 h-4" />
                </button>
              </template>
              <template v-else>
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
          <td colspan="8" class="py-32">
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
