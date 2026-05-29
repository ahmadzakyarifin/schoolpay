<script setup>
import { computed, ref } from 'vue'
import { 
  Edit as EditIcon, 
  Trash as TrashIcon, 
  Undo2 as RestoreIcon,
  Check as CheckIcon,
  GraduationCap as GraduationCapIcon,
  X as CloseIcon,
  Users as ClassIcon
} from 'lucide-vue-next'

const props = defineProps({
  list: Array,
  loading: Boolean,
  showHistory: Boolean,
  selectedIds: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:selectedIds', 'edit', 'delete', 'restore', 'toggle-status'])

const showClassDetail = ref(false)
const showMajorDetail = ref(false)
const selectedBatch = ref(null)

const openClassDetail = (batch) => {
  selectedBatch.value = batch
  showClassDetail.value = true
}

const openMajorDetail = (batch) => {
  selectedBatch.value = batch
  showMajorDetail.value = true
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const d = new Date(dateString)
  return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear()}`
}

const isAllSelected = computed(() => {
  return props.selectedIds.length === props.list.length && props.list.length > 0
})

const toggleAll = () => {
  if (isAllSelected.value) {
    emit('update:selectedIds', [])
  } else {
    emit('update:selectedIds', props.list.map(i => i.id))
  }
}

const toggleOne = (id) => {
  const current = [...props.selectedIds]
  const index = current.indexOf(id)
  if (index > -1) {
    current.splice(index, 1)
  } else {
    current.push(id)
  }
  emit('update:selectedIds', current)
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
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Tahun Angkatan</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Jurusan Aktif</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Kelas</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Status</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Dibuat Pada</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right pr-10">Aksi</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-slate-50">
        <tr v-if="list.length === 0">
          <td colspan="7" class="px-6 py-20 text-center">
            <div class="flex flex-col items-center justify-center">
              <div class="w-20 h-20 bg-slate-100 rounded-[2.5rem] flex items-center justify-center text-slate-300 mb-6 border-4 border-white shadow-xl shadow-slate-200/50">
                <GraduationCapIcon class="w-10 h-10" />
              </div>
              <h3 class="text-lg font-black text-slate-700 tracking-tight mb-2">Tidak Ada Data Angkatan</h3>
              <p class="text-slate-400 text-xs font-medium max-w-xs mx-auto">
                {{ showHistory ? 'Tidak ada data angkatan di riwayat penghapusan.' : 'Belum ada angkatan yang terdaftar atau coba sesuaikan kata kunci pencarian Anda.' }}
              </p>
            </div>
          </td>
        </tr>
        <template v-else>
          <tr v-for="item in list" :key="item.id" class="group hover:bg-slate-50/50 transition-colors border-b border-slate-50 last:border-0" :class="{'bg-indigo-50/30': selectedIds.includes(item.id)}">
            <td class="px-4 py-4 text-center">
              <div @click="toggleOne(item.id)" 
                class="w-4 h-4 mx-auto rounded border-2 flex items-center justify-center cursor-pointer transition-all"
                :class="selectedIds.includes(item.id) ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300 group-hover:border-indigo-300'">
                <CheckIcon v-if="selectedIds.includes(item.id)" class="w-2.5 h-2.5 text-white" />
              </div>
            </td>
            <td class="px-4 py-4">
              <div class="flex items-center gap-3">
                <div class="w-8 h-8 bg-slate-100 text-slate-500 rounded-lg flex items-center justify-center font-bold text-xs border border-slate-200/60 shadow-sm shrink-0">
                  {{ item.year.toString().slice(-2) }}
                </div>
                <div class="flex flex-col truncate">
                  <span class="font-bold text-slate-500 text-xs uppercase tracking-wider truncate">Angkatan {{ item.year }}</span>
                  <span class="text-[9px] font-bold text-slate-400 uppercase tracking-widest mt-0.5">Tahun Masuk</span>
                </div>
              </div>
            </td>
            <td class="px-4 py-4">
              <div class="flex items-center gap-1.5 relative">
                <button @click="openMajorDetail(item)" 
                  class="flex items-center justify-center w-6 h-6 rounded-lg transition-all"
                  :class="item.major_count > 0 ? 'bg-indigo-50 text-indigo-600 hover:bg-indigo-100' : 'bg-slate-100 text-slate-400 cursor-not-allowed'"
                >
                  <GraduationCapIcon class="w-3 h-3" />
                </button>
                <span class="text-xs font-bold text-slate-700">{{ item.major_count || 0 }}</span>
              </div>
            </td>
            <td class="px-4 py-4">
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
            <td class="px-4 py-4">
              <div v-if="!showHistory" class="flex items-center">
                <button @click="emit('toggle-status', item)" 
                  class="relative w-8 h-4 rounded-full transition-all duration-300 focus:outline-none shadow-inner"
                  :class="item.is_active ? 'bg-indigo-600' : 'bg-slate-300'">
                  <div class="absolute top-0.5 left-0.5 w-3 h-3 bg-white rounded-full shadow transition-transform duration-300"
                    :class="item.is_active ? 'translate-x-4' : 'translate-x-0'"></div>
                </button>
              </div>
              <div v-else class="flex items-center">
                <span class="text-[10px] font-bold text-rose-500 uppercase tracking-widest">Terhapus</span>
              </div>
            </td>
            <td class="px-4 py-4 text-left">
              <span class="text-xs font-medium text-slate-600">{{ formatDate(item.created_at) }}</span>
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
        </template>
      </tbody>
    </table>

    <!-- Major Detail Popup -->
    <transition name="fade">
      <div v-if="showMajorDetail" class="fixed inset-0 z-[1100] flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm" @click="showMajorDetail = false">
        <div class="bg-white w-full max-w-sm rounded-[2.5rem] shadow-2xl overflow-hidden animate-scale-in pb-8" @click.stop>
          <div class="px-8 py-6 border-b border-slate-50 flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-10 h-10 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center">
                <GraduationCapIcon class="w-5 h-5" />
              </div>
              <div>
                <h3 class="font-black text-slate-800 text-lg tracking-tight">Daftar Jurusan</h3>
                <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mt-0.5">Angkatan {{ selectedBatch?.year }}</p>
              </div>
            </div>
            <button @click="showMajorDetail = false" class="p-2 hover:bg-slate-50 text-slate-400 rounded-xl transition-all">
              <CloseIcon class="w-5 h-5" />
            </button>
          </div>
          <div class="p-6 max-h-[60vh] overflow-y-auto custom-scrollbar">
            <div v-if="selectedBatch?.major_count > 0" class="flex flex-col gap-2">
              <div v-for="(name, idx) in selectedBatch.major_names" :key="idx" 
                class="p-4 bg-slate-50 border border-slate-100 rounded-2xl flex items-center gap-4 group hover:border-indigo-200 transition-all">
                <div class="w-8 h-8 bg-white text-indigo-500 rounded-xl flex items-center justify-center border border-slate-100 shadow-sm group-hover:scale-110 transition-transform">
                  <CheckIcon class="w-4 h-4" />
                </div>
                <span class="text-[11px] font-black text-slate-700 uppercase tracking-tight">{{ name }}</span>
              </div>
            </div>
            <div v-else class="py-10 text-center">
              <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Belum ada jurusan aktif</p>
            </div>
          </div>
        </div>
      </div>
    </transition>

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
                <p class="text-[9px] font-black text-slate-400 uppercase tracking-widest mt-0.5">Angkatan {{ selectedBatch?.year }}</p>
              </div>
            </div>
            <button @click="showClassDetail = false" class="p-2 hover:bg-slate-50 text-slate-400 rounded-xl transition-all">
              <CloseIcon class="w-5 h-5" />
            </button>
          </div>
          <div class="p-6 max-h-[60vh] overflow-y-auto custom-scrollbar">
            <div v-if="selectedBatch?.class_count > 0" class="grid grid-cols-2 gap-3">
              <div v-for="(name, idx) in selectedBatch.class_names" :key="idx" 
                class="p-3 bg-slate-50 border border-slate-100 rounded-2xl flex items-center gap-3 group hover:border-indigo-200 transition-all">
                <div class="w-6 h-6 bg-white text-indigo-500 rounded-lg flex items-center justify-center border border-slate-100 shadow-sm">
                  <CheckIcon class="w-3 h-3" />
                </div>
                <span class="text-[10px] font-black text-slate-700 uppercase tracking-tight truncate">{{ name }}</span>
              </div>
            </div>
            <div v-else class="py-10 text-center">
              <p class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Belum ada kelas aktif</p>
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
