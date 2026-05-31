<script setup>
import { 
  MessageSquare as WhatsAppIcon,
  Mail as MailIcon,
  Send as SendIcon,
  Trash as TrashIcon, 
  RotateCcw as ResetIcon,
  ChevronDown as ChevronIcon,
  X as CloseIcon
} from 'lucide-vue-next'
import { ref, computed } from 'vue'

const props = defineProps({
  selectedCount: Number,
  selectedUsers: {
    type: Array,
    default: () => []
  },
  status: String
})

const emit = defineEmits(['resend', 'delete', 'restore'])
const showMenu = ref(false)

// Validasi: Berapa banyak dari yang dipilih yang BENAR-BENAR bisa dikirimi notifikasi?
const eligibleUsers = computed(() => {
  return props.selectedUsers.filter(u => {
    const isActive = u.is_active
    const isParentWithChild = u.role !== 'parent' || (u.student_count && u.student_count > 0)
    const hasNoPassword = !u.has_password
    return isActive && isParentWithChild && hasNoPassword
  })
})

const hasIneligible = computed(() => eligibleUsers.value.length < props.selectedCount)

const handleResend = (channel) => {
  emit('resend', channel)
  showMenu.value = false
}
</script>

<template>
  <transition name="fade">
    <div v-if="selectedCount > 0" class="flex items-center gap-2">
      <!-- Bulk Notification Button with Dropdown -->
      <div v-if="status !== 'trash'" class="relative">
        <button 
          @click="showMenu = !showMenu" 
          class="bg-white border border-indigo-200 hover:bg-indigo-50/50 hover:border-indigo-300 text-indigo-600 font-black py-2 px-4 rounded-xl flex items-center gap-2 transition-all text-xs shadow-sm cursor-pointer"
          :class="{ 'ring-2 ring-indigo-100': showMenu }"
        >
          <SendIcon class="w-3.5 h-3.5 text-indigo-500" />
          <span>Kirim Aktivasi ({{ eligibleUsers.length }})</span>
          <ChevronIcon class="w-3.5 h-3.5 text-indigo-400 transition-transform" :class="{ 'rotate-180': showMenu }" />
        </button>

        <!-- Warning badge if some users are ineligible -->
        <div v-if="hasIneligible" class="absolute -top-2 -right-2 w-5 h-5 bg-amber-500 text-white rounded-full flex items-center justify-center border-2 border-white shadow-sm cursor-help" :title="`${selectedCount - eligibleUsers.length} user dilewati karena Non-Aktif, sudah punya password, atau Wali tanpa anak` ">
          <span class="text-[9px] font-black">!</span>
        </div>

        <!-- Options Dropdown -->
        <transition name="fade-scale">
          <div v-if="showMenu" 
            class="absolute top-full left-0 mt-2 w-56 bg-white rounded-2xl shadow-[0_15px_50px_rgba(0,0,0,0.2)] border border-slate-100 z-[200] overflow-hidden p-2 origin-top-left"
          >
            <div class="px-3 py-2 border-b border-slate-50 mb-1 flex items-center justify-between">
              <h4 class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Kirim Aktivasi</h4>
              <button @click="showMenu = false" class="p-1 hover:bg-slate-100 text-slate-400 rounded-md">
                <CloseIcon class="w-3 h-3" />
              </button>
            </div>
            
            <div class="space-y-0.5">
              <button 
                @click="handleResend('whatsapp')"
                class="w-full flex items-center gap-3 px-3 py-3 rounded-xl hover:bg-emerald-50 text-slate-600 hover:text-emerald-700 transition-all group text-left"
              >
                <div class="w-8 h-8 bg-emerald-100/50 rounded-lg flex items-center justify-center group-hover:bg-emerald-100">
                  <WhatsAppIcon class="w-4 h-4 text-emerald-600" />
                </div>
                <div class="flex flex-col">
                  <span class="text-[11px] font-bold">WhatsApp</span>
                  <span class="text-[8px] font-medium text-slate-400">Kirim ke nomor WA</span>
                </div>
              </button>

              <button 
                @click="handleResend('email')"
                class="w-full flex items-center gap-3 px-3 py-3 rounded-xl hover:bg-indigo-50 text-slate-600 hover:text-indigo-700 transition-all group text-left"
              >
                <div class="w-8 h-8 bg-indigo-100/50 rounded-lg flex items-center justify-center group-hover:bg-indigo-100">
                  <MailIcon class="w-4 h-4 text-indigo-600" />
                </div>
                <div class="flex flex-col">
                  <span class="text-[11px] font-bold">Email</span>
                  <span class="text-[8px] font-medium text-slate-400">Kirim ke alamat email</span>
                </div>
              </button>

              <div class="h-px bg-slate-50 my-1"></div>

              <button 
                @click="handleResend('all')"
                class="w-full flex items-center gap-3 px-3 py-3 rounded-xl hover:bg-slate-50 text-slate-700 hover:text-slate-900 transition-all group text-left"
              >
                <div class="w-8 h-8 bg-slate-100 rounded-lg flex items-center justify-center group-hover:bg-slate-200">
                  <SendIcon class="w-4 h-4 text-slate-600" />
                </div>
                <div class="flex flex-col">
                  <span class="text-[11px] font-black italic">Semua Media</span>
                  <span class="text-[8px] font-medium text-slate-400 italic">WA & Email sekaligus</span>
                </div>
              </button>
            </div>
          </div>
        </transition>
      </div>

      <button 
        v-if="status !== 'trash'"
        @click="emit('delete')" 
        class="bg-white border border-rose-200 hover:bg-rose-50/50 hover:border-rose-300 text-rose-600 font-black py-2 px-4 rounded-xl text-xs flex items-center gap-2 transition-all shadow-sm cursor-pointer"
      >
        <TrashIcon class="w-3.5 h-3.5 text-rose-500" />
        <span>Hapus Terpilih ({{ selectedCount }})</span>
      </button>
      
      <button 
        v-if="status === 'trash'"
        @click="emit('restore')" 
        class="bg-white border border-emerald-200 hover:bg-emerald-50/50 hover:border-emerald-300 text-emerald-600 font-black py-2 px-4 rounded-xl text-xs flex items-center gap-2 transition-all shadow-sm cursor-pointer"
      >
        <ResetIcon class="w-3.5 h-3.5 text-emerald-500" />
        <span>Pulihkan Terpilih ({{ selectedCount }})</span>
      </button>
    </div>
  </transition>
</template>

<style scoped>
.fade-scale-enter-active, .fade-scale-leave-active {
  transition: all 0.2s cubic-bezier(0.34, 1.56, 0.64, 1);
}
.fade-scale-enter-from, .fade-scale-leave-to {
  opacity: 0;
  transform: scale(0.9) translateY(-10px);
}
</style>
