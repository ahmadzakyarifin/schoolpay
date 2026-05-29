<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../../store/auth'
import { useToast } from '../../composables/useToast'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { 
  User as UserIcon, 
  Lock as LockIcon, 
  Mail as MailIcon, 
  Phone as PhoneIcon,
  ShieldCheck as ShieldIcon,
  Save as SaveIcon,
  ArrowLeft as BackIcon
} from 'lucide-vue-next'

const authStore = useAuthStore()
const router = useRouter()
const toast = useToast()
const loading = ref(false)

const profile = ref({
  name: authStore.user?.name || '',
  email: authStore.user?.email || '',
  phone: authStore.user?.phone || ''
})

const passwordData = ref({
  current_password: '',
  new_password: '',
  confirm_password: ''
})

const goBack = () => {
  router.back()
}

const updatePassword = async () => {
  if (passwordData.value.new_password !== passwordData.value.confirm_password) {
    return toast.error('Error', 'Konfirmasi password tidak cocok')
  }

  loading.value = true
  try {
    await axios.post('auth/change-password', {
      current_password: passwordData.value.current_password,
      new_password: passwordData.value.new_password
    })
    toast.success('Sukses', 'Password berhasil diperbarui. Silakan login ulang.')
    passwordData.value = { current_password: '', new_password: '', confirm_password: '' }
    authStore.clearAuth()
    router.push({ name: 'login' })
  } catch (err) {
    toast.error('Gagal', err.response?.data?.message || 'Gagal memperbarui password')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-4xl mx-auto p-4 lg:p-8 space-y-8 animate-fade-in">
    <!-- Header with Back Button -->
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <div class="w-12 h-12 bg-indigo-600 text-white rounded-2xl flex items-center justify-center shadow-xl shadow-indigo-100">
          <UserIcon class="w-6 h-6" />
        </div>
        <div>
          <h2 class="text-2xl font-black text-slate-800 tracking-tight">Profil & Keamanan</h2>
          <p class="text-sm font-bold text-slate-400 uppercase tracking-widest">Kelola akun Anda</p>
        </div>
      </div>
      
      <button @click="goBack" class="flex items-center gap-2 px-5 py-2.5 bg-white border border-slate-200 text-slate-600 font-black text-[10px] uppercase tracking-widest rounded-xl hover:bg-slate-50 transition-all shadow-sm">
        <BackIcon class="w-4 h-4" />
        <span>Kembali</span>
      </button>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Info Section -->
      <div class="lg:col-span-1 space-y-6">
        <div class="white-card p-8 text-center">
          <div class="w-24 h-24 bg-slate-100 rounded-[2rem] mx-auto mb-6 flex items-center justify-center text-indigo-600">
             <UserIcon class="w-12 h-12" />
          </div>
          <h3 class="text-lg font-black text-slate-800 uppercase tracking-tight">{{ profile.name }}</h3>
          <p class="text-[10px] font-black text-indigo-600 bg-indigo-50 inline-block px-4 py-1 rounded-full uppercase mt-2">
            {{ authStore.user?.role }}
          </p>
          
          <div class="mt-8 space-y-4 text-left border-t border-slate-50 pt-8">
            <div class="flex items-center gap-3 text-slate-500">
              <MailIcon class="w-4 h-4 text-slate-300" />
              <span class="text-xs font-bold truncate">{{ profile.email }}</span>
            </div>
            <div class="flex items-center gap-3 text-slate-500">
              <PhoneIcon class="w-4 h-4 text-slate-300" />
              <span class="text-xs font-bold">{{ profile.phone || '-' }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Change Password Section -->
      <div class="lg:col-span-2">
        <div class="white-card p-8">
          <div class="flex items-center gap-3 mb-8">
            <div class="w-8 h-8 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center">
              <LockIcon class="w-4 h-4" />
            </div>
            <h3 class="text-sm font-black text-slate-800 uppercase tracking-widest">Ganti Password</h3>
          </div>

          <form @submit.prevent="updatePassword" class="space-y-6">
            <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div class="space-y-2">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Password Saat Ini</label>
                <input v-model="passwordData.current_password" type="password" required class="modern-input" placeholder="••••••••">
              </div>
              <div class="space-y-2">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Password Baru</label>
                <input v-model="passwordData.new_password" type="password" required class="modern-input" placeholder="••••••••">
              </div>
              <div class="space-y-2">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1">Konfirmasi Password</label>
                <input v-model="passwordData.confirm_password" type="password" required class="modern-input" placeholder="••••••••">
              </div>
            </div>

            <div class="pt-4">
              <button type="submit" :disabled="loading" class="w-full btn-primary py-4 gap-3">
                <SaveIcon v-if="!loading" class="w-5 h-5" />
                <span v-else class="animate-spin border-2 border-white/20 border-t-white rounded-full w-5 h-5"></span>
                <span>Perbarui Password</span>
              </button>
            </div>
          </form>
          
          <div class="mt-8 p-6 bg-amber-50 rounded-2xl border border-amber-100 flex gap-4">
            <ShieldIcon class="w-6 h-6 text-amber-500 shrink-0" />
            <p class="text-[10px] font-bold text-amber-700 leading-relaxed uppercase">
              Masukkan password saat ini untuk memastikan perubahan dilakukan oleh pemilik akun.
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="postcss">
.white-card {
  @apply bg-white border border-slate-100 rounded-[2.5rem] transition-all duration-300 shadow-sm;
}
</style>
