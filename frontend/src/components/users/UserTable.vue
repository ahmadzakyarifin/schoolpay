<script setup>
import { ref } from 'vue'
import { 
  ShieldCheck as AdminIcon, 
  UserCheck as ParentIcon, 
  Users as StudentsIcon,
  RotateCcw as ResetIcon,
  Trash as TrashIcon,
  Edit as EditIcon,
  Send as SendIcon,
  Mail as MailIcon,
  MessageSquare as WhatsAppIcon,
  ChevronDown as ChevronDownIcon,
  Check as CheckIcon,
  X as CloseIcon
} from 'lucide-vue-next'
import { useAuthStore } from '../../store/auth'

const authStore = useAuthStore()

const props = defineProps({
  users: Array,
  loading: Boolean,
  selectedUserIds: Array,
  status: String,
  formatDate: Function
})

const activeStudentListId = ref(null)
const activeResendMenuId = ref(null)
const emit = defineEmits([
	'edit', 
	'delete', 
	'restore', 
	'toggle-status', 
	'toggle-select-all', 
	'toggle-select-user', 
	'go-to-student', 
	'go-to-details',
  'resend-notification'
])

const parseStudents = (studentNamesStr) => {
  if (!studentNamesStr) return []
  return studentNamesStr.split('||').map(s => {
    const [id, name] = s.split('::')
    return { id, name }
  })
}

const toggleStudentList = (userId) => {
  if (activeStudentListId.value === userId) {
    activeStudentListId.value = null
  } else {
    activeStudentListId.value = userId
  }
}

const toggleResendMenu = (userId) => {
  if (activeResendMenuId.value === userId) {
    activeResendMenuId.value = null
  } else {
    activeResendMenuId.value = userId
    activeStudentListId.value = null // Close student list if open
  }
}

const canResendActivation = (user) => {
  if (props.status === 'trash' || !user?.is_active || user?.has_password) return false
  return user.role !== 'parent' || Number(user.student_count || 0) > 0
}

const resendActivation = (user, channel) => {
  emit('resend-notification', { user, channel })
  activeResendMenuId.value = null
}

const isAllSelected = () => {
  const selectable = props.users.filter(u => u.id !== authStore.user?.id)
  return selectable.length > 0 && props.selectedUserIds.length === selectable.length
}

</script>

<template>
  <div class="flex-1 overflow-x-auto custom-scrollbar">
    <table class="w-full text-left border-collapse">
      <thead>
        <tr class="bg-slate-50/50 border-b border-slate-100">
          <th class="w-[50px] px-4 py-4">
            <div @click="emit('toggle-select-all')" class="w-4 h-4 rounded border-2 flex items-center justify-center cursor-pointer transition-all"
              :class="isAllSelected() ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300'">
              <CheckIcon v-if="isAllSelected()" class="w-2.5 h-2.5 text-white" />
            </div>
          </th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Pengguna</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Email</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">WhatsApp</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Role</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Tanggal Dibuat</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Status</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-left">Siswa</th>
          <th class="px-4 py-4 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right px-10">Aksi</th>
        </tr>
      </thead>
      <tbody :class="{'opacity-50 pointer-events-none': loading}">
        <tr v-for="user in users" :key="user.id" class="group hover:bg-slate-50/50 transition-colors border-b border-slate-50 last:border-0"
          :class="{'bg-indigo-50/30': selectedUserIds.includes(user.id)}">
          <td class="px-4 py-4">
            <div v-if="user.id !== authStore.user?.id" @click="emit('toggle-select-user', user.id)" 
              class="w-4 h-4 rounded border-2 flex items-center justify-center cursor-pointer transition-all"
              :class="selectedUserIds.includes(user.id) ? 'bg-indigo-600 border-indigo-600' : 'border-slate-300 group-hover:border-indigo-300'">
              <CheckIcon v-if="selectedUserIds.includes(user.id)" class="w-2.5 h-2.5 text-white" />
            </div>
          </td>
          <td class="px-4 py-4">
              <div class="flex items-center gap-3 cursor-pointer" @click="emit('go-to-details', user)">
                <div class="w-8 h-8 rounded-lg flex items-center justify-center transition-all shadow-sm shrink-0"
                  :class="user.role === 'admin' ? 'bg-indigo-50/50 text-indigo-500 border border-indigo-100/50 group-hover:bg-indigo-50 group-hover:text-indigo-600' : 'bg-slate-100/70 text-slate-500 border border-slate-200/60 group-hover:bg-slate-200/50 group-hover:text-indigo-600'">
                  <AdminIcon v-if="user.role === 'admin'" class="w-4 h-4" />
                  <ParentIcon v-else class="w-4 h-4" />
                </div>
                <div>
                  <div class="font-bold text-slate-500 text-xs uppercase tracking-wider truncate flex items-center gap-2 group-hover:text-indigo-600 transition-colors">
                    {{ user.name }}
                    <span v-if="user.role === 'parent' && !user.student_count" class="text-[9px] font-normal text-slate-400 italic lowercase tracking-normal" title="Belum terhubung ke Siswa">(tanpa anak)</span>
                  </div>
                  <div class="text-[9px] font-bold text-slate-400 mt-0.5">ID: #{{ user.id }}</div>
                </div>
              </div>
          </td>
          <td class="px-4 py-4">
            <div class="flex items-center gap-2 text-xs font-medium text-slate-600">
              <span class="truncate max-w-[120px]">{{ user.email || '-' }}</span>
            </div>
          </td>
          <td class="px-4 py-4">
            <div class="flex items-center gap-2 text-xs font-medium text-slate-600">
              <span class="truncate">+{{ user.phone_number || '-' }}</span>
            </div>
          </td>
          <td class="px-4 py-4 text-xs font-medium text-slate-600 capitalize">
            {{ user.role === 'admin' ? 'Admin' : 'Wali' }}
          </td>
          <td class="px-4 py-4 text-xs font-medium text-slate-600">
            {{ formatDate ? formatDate(user.created_at) : new Date(user.created_at).toLocaleDateString() }}
          </td>
          <td class="px-4 py-4">
            <div class="flex items-center">
              <button 
                v-if="status !== 'trash'"
                @click.stop="user.id !== authStore.user?.id && emit('toggle-status', user)"
                class="relative w-8 h-4 rounded-full transition-all duration-300 focus:outline-none shadow-inner"
                :class="[
                  user.is_active ? 'bg-indigo-600' : 'bg-slate-300',
                  user.id === authStore.user?.id ? 'opacity-50 cursor-not-allowed grayscale' : 'cursor-pointer'
                ]"
                :title="user.id === authStore.user?.id ? 'Anda tidak bisa menonaktifkan akun sendiri' : ''"
              >
                <div class="absolute top-0.5 left-0.5 w-3 h-3 bg-white rounded-full shadow transition-transform duration-300"
                  :class="user.is_active ? 'translate-x-4' : 'translate-x-0'"></div>
              </button>
              <span v-else class="text-[10px] font-bold text-rose-500 uppercase tracking-widest">Dihapus</span>
            </div>
          </td>
          <td class="px-4 py-4">
            <div v-if="user.role === 'parent'" class="flex items-center gap-1.5 relative">
              <button 
                @click.stop="toggleStudentList(user.id)"
                class="flex items-center justify-center w-6 h-6 rounded-lg transition-all"
                :class="user.student_count > 0 ? 'bg-indigo-50 text-indigo-600 hover:bg-indigo-100' : 'bg-slate-100 text-slate-400 cursor-not-allowed'"
              >
                <StudentsIcon class="w-3 h-3" />
              </button>
              <span class="text-xs font-bold text-slate-700">{{ user.student_count || 0 }}</span>

              <!-- Student List Popover -->
              <transition name="fade-scale">
                <div v-if="activeStudentListId === user.id" 
                  class="absolute bottom-full left-0 mb-2 w-52 bg-white rounded-2xl shadow-[0_10px_30px_rgba(0,0,0,0.15)] border border-slate-100 z-[100] overflow-hidden p-2 origin-bottom-left"
                >
                  <div class="px-3 py-2 border-b border-slate-50 mb-1 flex items-center justify-between">
                    <h4 class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Siswa Terhubung</h4>
                    <button @click.stop="activeStudentListId = null" class="p-1 hover:bg-slate-100 text-slate-400 hover:text-slate-600 rounded-md transition-all">
                      <CloseIcon class="w-3 h-3" />
                    </button>
                  </div>
                  <div class="max-h-40 overflow-y-auto custom-scrollbar">
                    <div v-if="user.student_count > 0" class="space-y-0.5">
                      <button 
                        v-for="student in parseStudents(user.student_names)" 
                        :key="student.id"
                        @click.stop="emit('go-to-student', student.id)"
                        class="w-full text-left px-3 py-2 rounded-xl hover:bg-indigo-50 group transition-all"
                      >
                        <div class="text-[11px] font-bold text-slate-700 group-hover:text-indigo-600 truncate">{{ student.name }}</div>
                        <div class="text-[8px] font-medium text-slate-400 group-hover:text-indigo-400">Klik untuk detail</div>
                      </button>
                    </div>
                    <div v-else class="px-3 py-4 text-center">
                      <p class="text-[10px] font-medium text-slate-400 italic text-center">Tidak ada siswa</p>
                    </div>
                  </div>
                </div>
              </transition>
            </div>
            <span v-else class="text-[10px] text-slate-300 font-medium italic">N/A</span>
          </td>
          <td class="px-4 py-4 text-right">
            <div class="flex items-center justify-end gap-1 px-4">
              <template v-if="status !== 'trash'">
                <div v-if="canResendActivation(user)" class="relative">
                  <button
                    @click.stop="toggleResendMenu(user.id)"
                    class="p-2 hover:bg-indigo-50 text-slate-400 hover:text-indigo-600 rounded-xl transition-all flex items-center gap-0.5"
                    title="Kirim ulang aktivasi"
                  >
                    <SendIcon class="w-4 h-4" />
                    <ChevronDownIcon class="w-2.5 h-2.5" />
                  </button>

                  <transition name="fade-scale">
                    <div
                      v-if="activeResendMenuId === user.id"
                      class="absolute right-0 top-full mt-2 w-52 bg-white rounded-2xl shadow-[0_15px_45px_rgba(15,23,42,0.18)] border border-slate-100 z-[120] overflow-hidden p-2 origin-top-right"
                    >
                      <div class="px-3 py-2 border-b border-slate-50 mb-1 flex items-center justify-between">
                        <h4 class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Aktivasi Akun</h4>
                        <button @click.stop="activeResendMenuId = null" class="p-1 hover:bg-slate-100 text-slate-400 hover:text-slate-600 rounded-md transition-all">
                          <CloseIcon class="w-3 h-3" />
                        </button>
                      </div>
                      <button
                        @click.stop="resendActivation(user, 'whatsapp')"
                        class="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl hover:bg-emerald-50 text-slate-600 hover:text-emerald-700 transition-all text-left"
                      >
                        <WhatsAppIcon class="w-4 h-4 text-emerald-600" />
                        <span class="text-[11px] font-bold">WhatsApp</span>
                      </button>
                      <button
                        @click.stop="resendActivation(user, 'email')"
                        class="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl hover:bg-indigo-50 text-slate-600 hover:text-indigo-700 transition-all text-left"
                      >
                        <MailIcon class="w-4 h-4 text-indigo-600" />
                        <span class="text-[11px] font-bold">Email</span>
                      </button>
                      <button
                        @click.stop="resendActivation(user, 'all')"
                        class="w-full flex items-center gap-3 px-3 py-2.5 rounded-xl hover:bg-slate-50 text-slate-700 transition-all text-left"
                      >
                        <SendIcon class="w-4 h-4 text-slate-600" />
                        <span class="text-[11px] font-black">WA & Email</span>
                      </button>
                    </div>
                  </transition>
                </div>
                <button @click="emit('edit', user)" class="p-2 hover:bg-amber-50 text-slate-400 hover:text-amber-600 rounded-xl transition-all" title="Ubah">
                  <EditIcon class="w-4 h-4" />
                </button>
                <button 
                  v-if="user.id !== authStore.user?.id" 
                  @click="emit('delete', user)" 
                  class="p-2 hover:bg-rose-50 text-slate-400 hover:text-rose-600 rounded-xl transition-all" 
                  title="Hapus"
                >
                  <TrashIcon class="w-4 h-4" />
                </button>
              </template>
              <template v-else>
                <button @click="emit('restore', user)" class="p-2 hover:bg-emerald-50 text-slate-400 hover:text-emerald-600 rounded-xl transition-all" title="Pulihkan">
                  <ResetIcon class="w-4 h-4" />
                </button>
              </template>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
    
    <div v-if="!loading && users.length === 0" class="flex flex-col items-center justify-center py-20 px-6 text-center">
      <div class="w-20 h-20 bg-slate-100 rounded-[2.5rem] flex items-center justify-center text-slate-300 mb-6 border-4 border-white shadow-xl shadow-slate-200/50">
        <StudentsIcon class="w-10 h-10" />
      </div>
      <h3 class="text-lg font-black text-slate-700 tracking-tight mb-2">Tidak Ada Data Pengguna</h3>
      <p class="text-slate-400 text-xs font-medium max-w-xs">Belum ada akun yang terdaftar atau coba sesuaikan filter pencarian Anda.</p>
    </div>
  </div>
</template>

<style scoped>
.fade-scale-enter-active,
.fade-scale-leave-active {
  transition: all 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.fade-scale-enter-from,
.fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.9) translateY(10px);
}

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #e2e8f0;
  border-radius: 10px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #cbd5e1;
}
</style>
