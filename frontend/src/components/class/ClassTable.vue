<script setup>
import { computed } from 'vue'
import { 
  Edit as EditIcon, 
  Trash as TrashIcon,
  Undo2 as RestoreIcon,
  Users as ClassIcon,
  Check as CheckIcon
} from 'lucide-vue-next'

const props = defineProps({
  list: Array,
  loading: Boolean,
  showHistory: Boolean,
  selectedIds: Array
})

const emit = defineEmits(['edit', 'delete', 'restore', 'toggle-status', 'update:selectedIds'])

const isAllSelected = computed(() => {
  return props.list.length > 0 && props.selectedIds.length === props.list.length
})

const toggleAll = () => {
  if (isAllSelected.value) {
    emit('update:selectedIds', [])
  } else {
    emit('update:selectedIds', props.list.map(item => item.id))
  }
}

const toggleOne = (id) => {
  const newSelection = [...props.selectedIds]
  const index = newSelection.indexOf(id)
  if (index > -1) {
    newSelection.splice(index, 1)
  } else {
    newSelection.push(id)
  }
  emit('update:selectedIds', newSelection)
}

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

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const d = new Date(dateString)
  return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear()}`
}
</script>

<template>
  <div class="flex-1 overflow-x-auto">
    <table class="w-full text-left border-collapse table-fixed">
      <thead>
        <tr class="bg-slate-50/50 border-b border-slate-100">
          <th class="w-[50px] px-4 py-4 text-center">
            <div @click="toggleAll" class="w-4 h-4 mx-auto rounded border-2 flex items-center justify-center cursor-pointer transition-all"
              :class="isAllSelected ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300'">
              <CheckIcon v-if="isAllSelected" class="w-2.5 h-2.5 text-white" />
            </div>
          </th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Nama Kelas</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Jurusan</th>
          <th v-if="!showHistory" class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-center">Status</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">{{ showHistory ? 'Dihapus' : 'Dibuat' }}</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right pr-10">Aksi</th>
        </tr>
      </thead>
      <tbody :class="{'opacity-50 pointer-events-none': loading}">
        <tr v-for="item in list" :key="item.id" 
          class="group hover:bg-slate-50/50 transition-colors border-b border-slate-50 last:border-0"
          :class="{'bg-indigo-50/30': selectedIds.includes(item.id)}">
          <td class="px-4 py-4 text-center">
            <div @click="toggleOne(item.id)" 
              class="w-4 h-4 mx-auto rounded border-2 flex items-center justify-center cursor-pointer transition-all"
              :class="selectedIds.includes(item.id) ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300 group-hover:border-indigo-300'">
              <CheckIcon v-if="selectedIds.includes(item.id)" class="w-2.5 h-2.5 text-white" />
            </div>
          </td>
          <td class="px-4 py-4">
            <div class="flex flex-col truncate">
              <span class="font-bold text-slate-500 text-xs uppercase tracking-wider truncate">{{ item.name }}</span>
              <span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">Tingkat {{ toRoman(item.grade) }}</span>
            </div>
          </td>
          <td class="px-4 py-4">
            <span class="text-xs font-medium text-slate-600 uppercase">{{ item.major_name || '-' }}</span>
          </td>
          <td v-if="!showHistory" class="px-4 py-4 text-center">
            <div class="flex justify-center">
              <button @click="emit('toggle-status', item)" 
                class="relative w-8 h-4 rounded-full transition-all duration-300 focus:outline-none shadow-inner"
                :class="item.is_active ? 'bg-indigo-600' : 'bg-slate-300'">
                <div class="absolute top-0.5 left-0.5 w-3 h-3 bg-white rounded-full shadow transition-transform duration-300"
                  :class="item.is_active ? 'translate-x-4' : 'translate-x-0'"></div>
              </button>
            </div>
          </td>
          <td class="px-4 py-4 text-left">
            <span class="text-xs font-medium text-slate-600">{{ showHistory ? formatDate(item.deleted_at) : formatDate(item.created_at) }}</span>
          </td>
          <td class="px-4 py-4 text-right pr-10">
            <div class="flex items-center justify-end gap-1 px-4">
              <template v-if="!showHistory">
                <button @click="emit('edit', item)" class="p-2 hover:bg-amber-50 text-slate-400 hover:text-amber-600 rounded-xl transition-all" title="Ubah">
                  <EditIcon class="w-4 h-4" />
                </button>
                <button @click="emit('delete', item)" class="p-2 hover:bg-rose-50 text-slate-400 hover:text-rose-600 rounded-xl transition-all" title="Hapus">
                  <TrashIcon class="w-4 h-4" />
                </button>
              </template>
              <template v-else>
                <button @click="emit('restore', item)" class="p-2 hover:bg-emerald-50 text-slate-400 hover:text-emerald-600 rounded-xl transition-all" title="Pulihkan">
                  <RestoreIcon class="w-4 h-4" />
                </button>
              </template>
            </div>
          </td>
        </tr>

        <!-- Empty State -->
        <tr v-if="list.length === 0 && !loading">
          <td colspan="6" class="py-32">
            <div class="flex flex-col items-center justify-center text-center animate-scale-in px-6">
              <div class="w-20 h-20 bg-slate-100 rounded-[2.5rem] flex items-center justify-center text-slate-300 mb-6 border-4 border-white shadow-xl shadow-slate-200/50">
                <ClassIcon class="w-10 h-10" />
              </div>
              <h3 class="text-lg font-black text-slate-700 tracking-tight mb-2">{{ showHistory ? 'Tidak Ada Riwayat' : 'Data Belum Tersedia' }}</h3>
              <p class="text-slate-400 text-xs font-medium max-w-xs">{{ showHistory ? 'Belum ada data yang dihapus dari daftar operasional.' : 'Silakan tambahkan data kelas baru untuk mulai mengelola data akademik.' }}</p>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
