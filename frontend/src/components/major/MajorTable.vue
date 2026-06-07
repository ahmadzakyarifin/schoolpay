<script setup>
import { computed, ref } from 'vue'
import { 
  Edit as EditIcon, 
  Trash as TrashIcon,
  Undo2 as RestoreIcon,
  GraduationCap as MajorIcon,
  Check as CheckIcon,
  Users as ClassIcon,
  X as CloseIcon
} from 'lucide-vue-next'

const props = defineProps({
  list: Array,
  loading: Boolean,
  showHistory: Boolean,
  selectedIds: Array
})

const emit = defineEmits(['edit', 'delete', 'restore', 'toggle-status', 'update:selectedIds'])

const showClassDetail = ref(false)
const selectedMajor = ref(null)

const openClassDetail = (item) => {
  selectedMajor.value = item
  showClassDetail.value = true
}

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

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const d = new Date(dateString)
  return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear()}`
}
</script>

<template>
  <div class="flex-1 overflow-x-auto">
    <table class="w-full text-left border-collapse">
      <thead>
        <tr class="bg-slate-50/50 border-b border-slate-100 text-[10px] font-black text-slate-400 uppercase tracking-widest">
          <th class="w-12 py-3 px-4 text-center">
            <div @click="toggleAll" class="w-4 h-4 mx-auto rounded border-2 flex items-center justify-center cursor-pointer transition-all"
              :class="isAllSelected ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300'">
              <CheckIcon v-if="isAllSelected" class="w-2.5 h-2.5 text-white" />
            </div>
          </th>
          <th class="py-3 px-4">Kode</th>
          <th class="py-3 px-4">Nama Jurusan</th>
          <th class="py-3 px-4">Kelas</th>
          <th v-if="!showHistory" class="py-3 px-4">Status</th>
          <th class="py-3 px-4">{{ showHistory ? 'Dihapus' : 'Dibuat' }}</th>
          <th class="py-3 px-4 text-center w-[200px]">Aksi</th>
        </tr>
      </thead>
      <tbody :class="{'opacity-50 pointer-events-none': loading}">
        <tr v-for="item in list" :key="item.id" 
          class="border-b border-slate-100 hover:bg-slate-50/30 transition-all text-xs font-semibold text-slate-600"
          :class="{'bg-indigo-50/30': selectedIds.includes(item.id)}">
          <td class="py-3 px-4 text-center">
            <div @click="toggleOne(item.id)" 
              class="w-4 h-4 mx-auto rounded border-2 flex items-center justify-center cursor-pointer transition-all"
              :class="selectedIds.includes(item.id) ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300 hover:border-indigo-300'">
              <CheckIcon v-if="selectedIds.includes(item.id)" class="w-2.5 h-2.5 text-white" />
            </div>
          </td>
          <td class="py-3 px-4">
            <span class="bg-slate-100 text-slate-700 px-2 py-0.5 rounded text-[10px] font-mono font-bold">{{ item.code || '-' }}</span>
          </td>
          <td class="py-3 px-4">
            <div class="font-black text-slate-800 text-xs uppercase tracking-wider truncate max-w-[200px]">{{ item.name }}</div>
          </td>
          <td class="py-3 px-4">
            <div class="flex items-center gap-1.5 relative">
              <button @click="openClassDetail(item)" 
                class="flex items-center justify-center w-6 h-6 rounded-lg transition-all"
                :class="item.class_count > 0 ? 'bg-indigo-50 text-indigo-600 hover:bg-indigo-100' : 'bg-slate-100 text-slate-400 cursor-not-allowed'"
              >
                <ClassIcon class="w-3 h-3" />
              </button>
              <span class="text-xs font-bold text-slate-700">{{ item.class_count || 0 }}</span>
            </div>
          </td>
          <td v-if="!showHistory" class="py-3 px-4">
            <div class="flex items-center">
              <button @click="emit('toggle-status', item)" 
                class="relative w-8 h-4 rounded-full transition-all duration-300 focus:outline-none shadow-inner"
                :class="item.is_active ? 'bg-indigo-600' : 'bg-slate-300'">
                <div class="absolute top-0.5 left-0.5 w-3 h-3 bg-white rounded-full shadow transition-transform duration-300"
                  :class="item.is_active ? 'translate-x-4' : 'translate-x-0'"></div>
              </button>
            </div>
          </td>
          <td class="py-3 px-4 text-left">
            <span class="text-slate-500 text-[11px]">{{ showHistory ? formatDate(item.deleted_at) : formatDate(item.created_at) }}</span>
          </td>
          <td class="py-3 px-4 text-center">
            <div class="flex items-center justify-center gap-1.5 flex-nowrap">
              <template v-if="!showHistory">
                <button @click="emit('edit', item)" title="Ubah" class="p-2 bg-white text-slate-600 border border-slate-200 hover:bg-slate-50 rounded-lg flex items-center justify-center transition-all shadow-sm">
                  <EditIcon class="w-3.5 h-3.5 text-slate-500" />
                </button>
                <button @click="emit('delete', item)" title="Hapus" class="p-2 bg-white text-rose-600 border border-slate-200 hover:bg-rose-50 rounded-lg flex items-center justify-center transition-all shadow-sm">
                  <TrashIcon class="w-3.5 h-3.5 text-rose-500" />
                </button>
              </template>
              <template v-else>
                <button @click="emit('restore', item)" title="Pulihkan" class="p-2 bg-white text-emerald-600 border border-slate-200 hover:bg-emerald-50 rounded-lg flex items-center justify-center transition-all shadow-sm">
                  <RestoreIcon class="w-3.5 h-3.5 text-emerald-600" />
                </button>
              </template>
            </div>
          </td>
        </tr>

        <!-- Empty State -->
        <tr v-if="list.length === 0 && !loading">
          <td colspan="7" class="py-24">
            <div class="flex flex-col items-center justify-center text-center animate-scale-in px-6">
              <div class="w-20 h-20 bg-slate-100 rounded-[2.5rem] flex items-center justify-center text-slate-300 mb-6 border-4 border-white shadow-xl shadow-slate-200/50">
                <MajorIcon class="w-10 h-10" />
              </div>
              <h3 class="text-lg font-black text-slate-700 tracking-tight mb-2">{{ showHistory ? 'Tidak Ada Riwayat' : 'Data Belum Tersedia' }}</h3>
              <p class="text-slate-400 text-xs font-medium max-w-xs">{{ showHistory ? 'Belum ada data jurusan yang dihapus dari daftar operasional.' : 'Silakan tambahkan data jurusan baru untuk mulai mengelola data akademik.' }}</p>
            </div>
          </td>
        </tr>
      </tbody>
    </table>

    <!-- Class Detail Popup -->
    <transition name="fade">
      <div v-if="showClassDetail" class="fixed inset-0 z-[1100] flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm" @click="showClassDetail = false">
        <div class="bg-white w-full max-w-sm rounded-[2.5rem] shadow-2xl overflow-hidden animate-scale-in pb-8" @click.stop>
          <div class="px-8 py-6 border-b border-slate-50 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center">
                <ClassIcon class="w-5 h-5" />
              </div>
              <div>
                <h3 class="font-black text-slate-800 text-lg tracking-tight">Daftar Kelas</h3>
                <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mt-0.5">Jurusan {{ selectedMajor?.name }}</p>
              </div>
            </div>
            <button @click="showClassDetail = false" class="p-2 hover:bg-slate-50 text-slate-400 rounded-xl transition-all">
              <CloseIcon class="w-5 h-5" />
            </button>
          </div>
          <div class="p-6 max-h-[60vh] overflow-y-auto custom-scrollbar">
            <div v-if="selectedMajor?.class_count > 0" class="grid grid-cols-2 gap-3">
              <div v-for="(name, idx) in selectedMajor.class_names" :key="idx" 
                class="p-3 bg-slate-50 border border-slate-100 rounded-2xl flex items-center gap-3 group hover:border-indigo-200 transition-all">
                <div class="w-6 h-6 bg-white text-indigo-500 rounded-lg flex items-center justify-center border border-slate-100 shadow-sm">
                  <CheckIcon class="w-3 h-3" />
                </div>
                <span class="text-[10px] font-black text-slate-700 uppercase tracking-tight truncate">{{ name }}</span>
              </div>
            </div>
            <div v-else class="py-10 text-center">
              <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Belum ada kelas di jurusan ini</p>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<style scoped lang="postcss">
.animate-scale-in { animation: scaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1); }
@keyframes scaleIn { from { opacity: 0; transform: scale(0.9); } to { opacity: 1; transform: scale(1); } }
.custom-scrollbar::-webkit-scrollbar { width: 4px; }
.custom-scrollbar::-webkit-scrollbar-track { background: transparent; }
.custom-scrollbar::-webkit-scrollbar-thumb { background: #e2e8f0; border-radius: 10px; }
.custom-scrollbar::-webkit-scrollbar-thumb:hover { background: #cbd5e1; }
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
