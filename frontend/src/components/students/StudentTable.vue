<script setup>
import { 
  User as StudentIcon,
  Search as SearchIcon, 
  Edit as EditIcon, 
  Trash as TrashIcon,
  RotateCcw as RestoreIcon,
  Phone as PhoneIcon,
  MessageCircle as WAIcon,
  Download as DownloadIcon,
  Info as InfoIcon,
  Users as ParentsIcon,
  Check as CheckIcon,
  AlertCircle as AlertIcon,
  Users as StudentsIcon
} from 'lucide-vue-next'

const props = defineProps({
  students: Array,
  loading: Boolean,
  selectedIds: Array,
  showHistory: Boolean,
  staticBase: String
})

const emit = defineEmits([
  'edit', 'delete', 'restore', 'toggle-status', 'view-details',
  'view-parents', 'toggle-select-all', 'toggle-select-user', 'go-to-parent'
])

const getStatusClass = (status) => {
  const classes = {
    'active': 'bg-emerald-50 text-emerald-600 border-emerald-100',
    'graduated': 'bg-blue-50 text-blue-600 border-blue-100',
    'mutasi': 'bg-orange-50 text-orange-600 border-orange-100',
    'dropout': 'bg-red-50 text-red-600 border-red-100'
  }
  return classes[status] || 'bg-slate-50 text-slate-400 border-slate-100'
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const d = new Date(dateString)
  return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear()}`
}
</script>

<template>
  <div class="flex-1 overflow-x-auto custom-scrollbar">
    <table class="w-full text-left border-collapse">
      <thead>
        <tr class="bg-slate-50/50 border-b border-slate-100">
          <th class="w-[50px] px-4 py-4">
            <div @click="emit('toggle-select-all', !(students.length > 0 && selectedIds.length === students.length))" 
              class="w-4 h-4 rounded border-2 flex items-center justify-center cursor-pointer transition-all"
              :class="students.length > 0 && selectedIds.length === students.length ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300'">
              <CheckIcon v-if="students.length > 0 && selectedIds.length === students.length" class="w-2.5 h-2.5 text-white" />
            </div>
          </th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Siswa</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">NIS</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">NISN</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Kelas & Jurusan</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Status</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Tanggal Dibuat</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Orang Tua</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right px-10">Aksi</th>
        </tr>
      </thead>
      <tbody :class="{'opacity-50 pointer-events-none': loading}">
        <tr v-for="s in students" :key="s.id" class="group hover:bg-slate-50/50 transition-colors border-b border-slate-50 last:border-0" 
          :class="{'bg-indigo-50/30': selectedIds.includes(s.id)}">
          <td class="px-4 py-4">
            <div @click="emit('toggle-select-user', s.id)" 
              class="w-4 h-4 rounded border-2 flex items-center justify-center cursor-pointer transition-all"
              :class="selectedIds.includes(s.id) ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300 hover:border-indigo-400'">
              <CheckIcon v-if="selectedIds.includes(s.id)" class="w-2.5 h-2.5 text-white" />
            </div>
          </td>
          <td class="px-4 py-4">
            <div @click="emit('view-details', s.id)" class="flex items-center gap-3 cursor-pointer group/name">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center bg-slate-100 text-slate-500 font-bold border border-slate-200/60 text-xs overflow-hidden shrink-0 shadow-sm group-hover/name:border-indigo-200 group-hover/name:shadow-md transition-all">
                <img v-if="s.image_path" :src="`${staticBase}/${s.image_path}`" class="w-full h-full object-cover" />
                <span v-else>{{ s.name.charAt(0) }}</span>
              </div>
              <div>
                <div class="font-bold text-slate-500 text-xs uppercase tracking-wider truncate flex items-center gap-2 group-hover/name:text-indigo-600 transition-colors">{{ s.name }}</div>
                <div class="flex items-center gap-2 mt-0.5">
                  <div class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">{{ (s.gender === 'L' || s.gender === 'Male' || s.gender === 'Laki-laki') ? 'Laki-laki' : 'Perempuan' }}</div>
                  <div class="w-1 h-1 bg-slate-200 rounded-full"></div>
                  <div class="text-[9px] font-bold text-slate-400 uppercase tracking-widest">Angkatan {{ s.entry_year }}</div>
                </div>
              </div>
            </div>
          </td>
          <td class="px-4 py-4">
            <div class="flex items-center gap-2 text-xs font-medium text-slate-600">{{ s.nis || '-' }}</div>
          </td>
          <td class="px-4 py-4">
            <div class="flex items-center gap-2 text-xs font-medium text-slate-600">{{ s.nisn || '-' }}</div>
          </td>
          <td class="px-4 py-4">
            <div class="flex flex-col text-left min-w-[120px]">
              <span class="text-xs font-medium text-slate-600 uppercase tracking-wider">{{ s.class_name || '-' }}</span>
              <span class="text-[9px] font-bold text-slate-400 uppercase mt-0.5">{{ s.major_name || '-' }}</span>
            </div>
          </td>
           <td class="px-4 py-4">
            <div class="flex items-center">
              <div v-if="!showHistory" class="text-xs font-medium text-slate-600 capitalize">
                {{ s.status === 'active' ? 'Aktif' : (s.status === 'graduated' ? 'Lulus' : 'Non-Aktif') }}
              </div>
              <span v-else class="text-[10px] font-bold text-rose-500 uppercase tracking-widest">Dihapus</span>
            </div>
          </td>
          <td class="px-4 py-4 text-xs font-medium text-slate-600">
            {{ formatDate(s.created_at) }}
          </td>
          <td class="px-4 py-4">
            <div class="flex items-center gap-1.5 relative">
              <button 
                @click="s.parent_id ? emit('go-to-parent', s.parent_id) : emit('view-parents', s)" 
                class="flex items-center justify-center w-6 h-6 rounded-lg transition-all"
                :class="s.parent_id ? 'bg-indigo-50 text-indigo-600 hover:bg-indigo-100' : 'bg-slate-100 text-slate-400 hover:bg-slate-200/50'"
                :title="s.parent_name || 'Belum terhubung ke Wali'"
              >
                <ParentsIcon class="w-3 h-3" />
              </button>
              <span class="text-xs font-bold text-slate-700">{{ s.parent_id ? 1 : 0 }}</span>
            </div>
          </td>
          <td class="px-4 py-4 text-right">
            <div class="flex items-center justify-end gap-1 px-4">
              <template v-if="!showHistory">
                <button @click="emit('edit', s)" class="p-2 hover:bg-amber-50 text-slate-400 hover:text-amber-600 rounded-xl transition-all" title="Ubah">
                  <EditIcon class="w-4 h-4" />
                </button>
                <button @click="emit('delete', s)" class="p-2 hover:bg-rose-50 text-slate-400 hover:text-rose-600 rounded-xl transition-all" title="Hapus">
                  <TrashIcon class="w-4 h-4" />
                </button>
              </template>
              <template v-else>
                <button @click="emit('restore', s)" class="p-2 hover:bg-emerald-50 text-slate-400 hover:text-emerald-600 rounded-xl transition-all" title="Pulihkan">
                  <RestoreIcon class="w-4 h-4" />
                </button>
              </template>
            </div>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- EMPTY STATE (SYNCED WITH USER) -->
    <div v-if="!loading && students.length === 0" class="flex flex-col items-center justify-center py-20 px-6 text-center">
      <div class="w-20 h-20 bg-slate-100 rounded-[2.5rem] flex items-center justify-center text-slate-300 mb-6 border-4 border-white shadow-xl shadow-slate-200/50">
        <StudentsIcon class="w-10 h-10" />
      </div>
      <h3 class="text-lg font-black text-slate-700 tracking-tight mb-2">Tidak Ada Data Siswa</h3>
      <p class="text-slate-400 text-xs font-medium max-w-xs">Belum ada siswa yang terdaftar atau coba sesuaikan filter pencarian Anda.</p>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar { width: 4px; height: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #E2E8F0; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #CBD5E1; }
</style>
