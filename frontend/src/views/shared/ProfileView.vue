<script setup>
import { computed, onMounted, ref } from 'vue'
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
  ArrowLeft as BackIcon,
  Camera as CameraIcon,
  CalendarDays as CalendarIcon,
  MapPin as MapIcon,
  Briefcase as WorkIcon,
  BadgeCheck as BadgeIcon
} from 'lucide-vue-next'

const authStore = useAuthStore()
const router = useRouter()
const toast = useToast()
const loading = ref(false)
const savingProfile = ref(false)
const photoLoading = ref(false)
const errors = ref({})
const photoPreview = ref('')

const apiBase = axios.defaults.baseURL
const staticBase = apiBase.replace('/api/', '')

const profile = ref({
  id: null,
  name: '',
  email: '',
  phone_number: '',
  role: '',
  image_path: '',
  nik: '',
  birth_date: '',
  address: '',
  education: '',
  occupation: '',
  income: ''
})

const passwordData = ref({
  current_password: '',
  new_password: '',
  confirm_password: ''
})

const isParent = computed(() => profile.value.role === 'parent')

const photoUrl = computed(() => {
  if (photoPreview.value) return photoPreview.value
  const path = profile.value.image_path
  if (!path) return ''
  return `${staticBase}${String(path).startsWith('/') ? path : `/${path}`}`
})

const initials = computed(() => {
  return String(profile.value.name || 'SP')
    .split(' ')
    .filter(Boolean)
    .slice(0, 2)
    .map(part => part[0])
    .join('')
    .toUpperCase()
})

const assignProfile = (user = {}) => {
  profile.value = {
    id: user.id || null,
    name: user.name || '',
    email: user.email || '',
    phone_number: user.phone_number || '',
    role: user.role || '',
    image_path: user.image_path || '',
    nik: user.nik || '',
    birth_date: user.birth_date || '',
    address: user.address || '',
    education: user.education || '',
    occupation: user.occupation || '',
    income: user.income || ''
  }
}

const syncAuthUser = (user) => {
  authStore.user = user
  localStorage.setItem('user', JSON.stringify(user))
}

const fetchProfile = async () => {
  loading.value = true
  try {
    const res = await axios.get('auth/me')
    const user = res.data.data || {}
    assignProfile(user)
    syncAuthUser(user)
  } catch (err) {
    assignProfile(authStore.user || {})
    toast.error('Gagal', err.response?.data?.message || 'Gagal memuat profil')
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.back()
}

const updateProfile = async () => {
  savingProfile.value = true
  errors.value = {}
  try {
    const payload = {
      name: profile.value.name,
      email: profile.value.email,
      phone_number: profile.value.phone_number,
      nik: profile.value.nik || '',
      birth_date: profile.value.birth_date || null,
      address: profile.value.address || '',
      education: profile.value.education || '',
      occupation: profile.value.occupation || '',
      income: profile.value.income || ''
    }
    const res = await axios.put('auth/profile', payload)
    const user = res.data.data || {}
    assignProfile(user)
    syncAuthUser(user)
    toast.success('Sukses', 'Profil berhasil diperbarui')
  } catch (err) {
    errors.value = err.response?.data?.errors || {}
    toast.error('Gagal', err.response?.data?.message || 'Gagal memperbarui profil')
  } finally {
    savingProfile.value = false
  }
}

const uploadPhoto = async (event) => {
  const file = event.target.files?.[0]
  if (!file) return
  photoPreview.value = URL.createObjectURL(file)
  photoLoading.value = true
  try {
    const payload = new FormData()
    payload.append('photo', file)
    const res = await axios.post('auth/profile/photo', payload, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    const user = res.data.data || {}
    assignProfile(user)
    syncAuthUser(user)
    photoPreview.value = ''
    toast.success('Sukses', 'Foto profil berhasil diperbarui')
  } catch (err) {
    photoPreview.value = ''
    toast.error('Gagal', err.response?.data?.message || 'Gagal mengunggah foto')
  } finally {
    photoLoading.value = false
    event.target.value = ''
  }
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

onMounted(fetchProfile)
</script>

<template>
  <div class="max-w-6xl mx-auto p-4 lg:p-8 space-y-8 animate-fade-in">
    <div class="flex items-center justify-between gap-4">
      <div class="flex items-center gap-4">
        <div class="w-12 h-12 bg-indigo-600 text-white rounded-2xl flex items-center justify-center shadow-xl shadow-indigo-100">
          <UserIcon class="w-6 h-6" />
        </div>
        <div>
          <h2 class="text-2xl font-black text-slate-800 tracking-tight">Profil & Keamanan</h2>
          <p class="text-sm font-bold text-slate-400 uppercase tracking-widest">Kelola identitas akun</p>
        </div>
      </div>

      <button @click="goBack" class="flex items-center gap-2 px-5 py-2.5 bg-white border border-slate-200 text-slate-600 font-black text-[10px] uppercase tracking-widest rounded-xl hover:bg-slate-50 transition-all shadow-sm">
        <BackIcon class="w-4 h-4" />
        <span>Kembali</span>
      </button>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <div class="lg:col-span-1 space-y-6">
        <div class="white-card p-8 text-center">
          <div class="relative w-28 h-28 mx-auto mb-6">
            <div class="w-28 h-28 bg-slate-100 rounded-2xl overflow-hidden flex items-center justify-center text-indigo-600 border border-slate-100">
              <img v-if="photoUrl" :src="photoUrl" class="w-full h-full object-cover" />
              <span v-else class="text-3xl font-black">{{ initials }}</span>
            </div>
            <label class="absolute -bottom-2 -right-2 w-10 h-10 rounded-xl bg-indigo-600 text-white shadow-lg shadow-indigo-100 flex items-center justify-center cursor-pointer hover:bg-indigo-700 transition-all" title="Ubah foto profil">
              <CameraIcon v-if="!photoLoading" class="w-4 h-4" />
              <span v-else class="w-4 h-4 rounded-full border-2 border-white/30 border-t-white animate-spin"></span>
              <input type="file" class="hidden" accept="image/png,image/jpeg,image/jpg,image/webp" @change="uploadPhoto">
            </label>
          </div>
          <h3 class="text-lg font-black text-slate-800 uppercase tracking-tight">{{ profile.name || '-' }}</h3>
          <p class="text-[10px] font-black text-indigo-600 bg-indigo-50 inline-block px-4 py-1 rounded-full uppercase mt-2">
            {{ profile.role || authStore.user?.role }}
          </p>

          <div class="mt-8 space-y-4 text-left border-t border-slate-50 pt-8">
            <div class="flex items-center gap-3 text-slate-500">
              <MailIcon class="w-4 h-4 text-slate-300" />
              <span class="text-xs font-bold truncate">{{ profile.email || '-' }}</span>
            </div>
            <div class="flex items-center gap-3 text-slate-500">
              <PhoneIcon class="w-4 h-4 text-slate-300" />
              <span class="text-xs font-bold">{{ profile.phone_number || '-' }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="lg:col-span-2 space-y-8">
        <div class="white-card p-8">
          <div class="flex items-center gap-3 mb-8">
            <div class="w-8 h-8 bg-indigo-50 text-indigo-600 rounded-xl flex items-center justify-center">
              <BadgeIcon class="w-4 h-4" />
            </div>
            <h3 class="text-sm font-black text-slate-800 uppercase tracking-widest">Data Akun</h3>
          </div>

          <form @submit.prevent="updateProfile" class="space-y-6">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="space-y-2">
                <label class="field-label">Nama Lengkap</label>
                <input v-model="profile.name" type="text" required class="modern-input" placeholder="Nama lengkap">
                <FormError :message="errors.name" />
              </div>
              <div class="space-y-2">
                <label class="field-label">Email</label>
                <input v-model="profile.email" type="email" required class="modern-input" placeholder="nama@email.com">
                <FormError :message="errors.email" />
              </div>
              <div class="space-y-2">
                <label class="field-label">Nomor WhatsApp</label>
                <input v-model="profile.phone_number" type="tel" required class="modern-input" placeholder="628123456789">
                <FormError :message="errors.phone_number" />
              </div>
              <div v-if="isParent" class="space-y-2">
                <label class="field-label">NIK</label>
                <input v-model="profile.nik" type="text" maxlength="16" class="modern-input" placeholder="16 digit">
                <FormError :message="errors.nik" />
              </div>
              <div v-if="isParent" class="space-y-2">
                <label class="field-label">Tanggal Lahir</label>
                <div class="relative">
                  <CalendarIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300" />
                  <input v-model="profile.birth_date" type="date" class="modern-input !pl-12">
                </div>
                <FormError :message="errors.birth_date" />
              </div>
              <div v-if="isParent" class="space-y-2">
                <label class="field-label">Pendidikan</label>
                <input v-model="profile.education" type="text" class="modern-input" placeholder="Pendidikan terakhir">
                <FormError :message="errors.education" />
              </div>
              <div v-if="isParent" class="space-y-2">
                <label class="field-label">Pekerjaan</label>
                <div class="relative">
                  <WorkIcon class="absolute left-4 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-300" />
                  <input v-model="profile.occupation" type="text" class="modern-input !pl-12" placeholder="Pekerjaan">
                </div>
                <FormError :message="errors.occupation" />
              </div>
              <div v-if="isParent" class="space-y-2">
                <label class="field-label">Penghasilan</label>
                <input v-model="profile.income" type="text" class="modern-input" placeholder="Rentang penghasilan">
                <FormError :message="errors.income" />
              </div>
            </div>

            <div v-if="isParent" class="space-y-2">
              <label class="field-label">Alamat</label>
              <div class="relative">
                <MapIcon class="absolute left-4 top-4 w-4 h-4 text-slate-300" />
                <textarea v-model="profile.address" rows="3" class="modern-input !pl-12 !h-auto resize-none" placeholder="Alamat lengkap"></textarea>
              </div>
              <FormError :message="errors.address" />
            </div>

            <div class="pt-2">
              <button type="submit" :disabled="savingProfile || loading" class="w-full btn-primary py-4 gap-3">
                <SaveIcon v-if="!savingProfile" class="w-5 h-5" />
                <span v-else class="animate-spin border-2 border-white/20 border-t-white rounded-full w-5 h-5"></span>
                <span>{{ savingProfile ? 'Menyimpan...' : 'Simpan Profil' }}</span>
              </button>
            </div>
          </form>
        </div>

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
                <label class="field-label">Password Saat Ini</label>
                <input v-model="passwordData.current_password" type="password" required class="modern-input" placeholder="********">
              </div>
              <div class="space-y-2">
                <label class="field-label">Password Baru</label>
                <input v-model="passwordData.new_password" type="password" required class="modern-input" placeholder="********">
              </div>
              <div class="space-y-2">
                <label class="field-label">Konfirmasi Password</label>
                <input v-model="passwordData.confirm_password" type="password" required class="modern-input" placeholder="********">
              </div>
            </div>

            <button type="submit" :disabled="loading" class="w-full btn-primary py-4 gap-3">
              <SaveIcon v-if="!loading" class="w-5 h-5" />
              <span v-else class="animate-spin border-2 border-white/20 border-t-white rounded-full w-5 h-5"></span>
              <span>Perbarui Password</span>
            </button>
          </form>

          <div class="mt-8 p-6 bg-amber-50 rounded-2xl border border-amber-100 flex gap-4">
            <ShieldIcon class="w-6 h-6 text-amber-500 shrink-0" />
            <p class="text-[10px] font-bold text-amber-700 leading-relaxed uppercase">
              Setelah password diperbarui, sesi akan ditutup dan Anda perlu login ulang.
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="postcss">
.white-card {
  @apply bg-white border border-slate-100 rounded-2xl transition-all duration-300 shadow-sm;
}
.field-label {
  @apply text-[10px] font-black text-slate-400 uppercase tracking-widest ml-1;
}
</style>
