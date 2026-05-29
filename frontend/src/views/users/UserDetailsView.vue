<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { 
  ArrowLeft as ArrowLeftIcon,
  User as UserIcon,
  Mail as MailIcon,
  Phone as PhoneIcon,
  ShieldCheck as ShieldCheckIcon,
  ShieldCheck as AdminIcon,
  UserCheck as ParentIcon,
  Calendar as CalendarIcon,
  MapPin as MapPinIcon,
  GraduationCap as EducationIcon,
  Briefcase as OccupationIcon,
  Wallet as IncomeIcon,
  Users as StudentsIcon,
  CheckCircle2 as SuccessIcon,
  XCircle as ErrorIcon
} from 'lucide-vue-next'
import axios from 'axios'
import userService from '../../services/user.service'
import studentService from '../../services/student.service'

const route = useRoute()
const router = useRouter()
const user = ref(null)
const students = ref([])
const loading = ref(true)
const error = ref(null)

const fetchUserDetails = async () => {
  loading.value = true
  error.value = null
  try {
    const response = await userService.getByID(route.params.id)
    user.value = response.data.data
    
    if (user.value.role === 'parent') {
      const studentRes = await studentService.getStudentsByParentID(user.value.id)
      students.value = studentRes.data.data || []
    }
  } catch (err) {
    console.error('Error fetching user details:', err)
    error.value = 'Gagal memuat data pengguna. Pastikan koneksi server aktif.'
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push({ name: 'user-management' })
  }
}

onMounted(() => {
  fetchUserDetails()
})

const formattedDate = (dateStr) => {
  if (!dateStr) return '-'
  const d = new Date(dateStr)
  return `${String(d.getDate()).padStart(2, '0')}/${String(d.getMonth() + 1).padStart(2, '0')}/${d.getFullYear()}`
}
</script>

<template>
  <div class="min-h-screen bg-slate-50/50 p-6 md:p-10 font-inter">
    <!-- Header -->
    <div class="max-w-5xl mx-auto mb-10 flex flex-col md:flex-row md:items-center justify-between gap-6 animate-fade-in">
      <div class="flex items-center gap-6">
        <button @click="goBack" class="w-12 h-12 bg-white border border-slate-200 rounded-2xl flex items-center justify-center text-slate-400 hover:text-indigo-600 hover:border-indigo-100 hover:shadow-xl hover:shadow-indigo-50 transition-all group">
          <ArrowLeftIcon class="w-5 h-5 group-hover:-translate-x-1 transition-transform" />
        </button>
        <div>
          <div class="flex items-center gap-3 mb-1">
            <h1 class="text-3xl font-black text-slate-800 tracking-tight">Detail Pengguna</h1>
            <span v-if="user" :class="[
              'px-3 py-1 rounded-full text-[10px] font-black uppercase tracking-widest border',
              user.role === 'admin' ? 'bg-indigo-50 text-indigo-600 border-indigo-100' : 'bg-emerald-50 text-emerald-600 border-emerald-100'
            ]">
              {{ user.role === 'admin' ? 'Administrator' : 'Wali Murid' }}
            </span>
          </div>
          <p class="text-slate-400 font-bold text-xs uppercase tracking-[0.2em]">Informasi Lengkap Akun & Personalia</p>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="max-w-5xl mx-auto py-20 flex flex-col items-center justify-center">
      <div class="w-12 h-12 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin mb-4"></div>
      <p class="text-slate-400 font-black text-[10px] uppercase tracking-widest">Memuat Data Pengguna...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="max-w-5xl mx-auto py-20 flex flex-col items-center justify-center text-center">
      <div class="w-20 h-20 bg-rose-50 rounded-3xl flex items-center justify-center mb-6 border border-rose-100 shadow-inner">
        <ErrorIcon class="w-10 h-10 text-rose-500" />
      </div>
      <h3 class="text-xl font-black text-slate-800 mb-2">{{ error }}</h3>
      <button @click="fetchUserDetails" class="text-indigo-600 font-bold hover:underline">Coba Lagi</button>
    </div>

    <!-- Content -->
    <div v-else-if="user" class="max-w-5xl mx-auto space-y-8 animate-slide-up">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Left Column: Main Profile -->
        <div class="lg:col-span-1 space-y-8">
          <div class="white-card p-10 !rounded-[2.5rem] shadow-[0_20px_50px_rgba(0,0,0,0.03)] border-slate-100/50 flex flex-col items-center text-center">
          <div class="w-32 h-32 rounded-[3rem] flex items-center justify-center mb-8 relative group cursor-pointer"
            :class="user.role === 'admin' ? 'bg-indigo-600 text-white shadow-2xl shadow-indigo-200' : 'bg-emerald-600 text-white shadow-2xl shadow-emerald-200'">
            <AdminIcon v-if="user.role === 'admin'" class="w-14 h-14" />
            <ParentIcon v-else class="w-14 h-14" />
            <div class="absolute -bottom-2 -right-2 w-10 h-10 bg-white rounded-2xl flex items-center justify-center shadow-lg border border-slate-50">
              <SuccessIcon v-if="user.is_active" class="w-6 h-6 text-emerald-500" />
              <ErrorIcon v-else class="w-6 h-6 text-slate-300" />
            </div>
          </div>
          
          <h2 class="text-2xl font-black text-slate-800 tracking-tight leading-tight">{{ user.name }}</h2>
          <p class="text-slate-400 font-bold text-[10px] uppercase tracking-widest mt-2">User ID: #{{ user.id }}</p>

          <div class="w-full h-px bg-slate-50 my-8"></div>

          <div class="w-full space-y-6">
            <div class="flex items-center gap-4 group">
              <div class="p-3 bg-slate-50 text-slate-400 rounded-2xl group-hover:bg-indigo-50 group-hover:text-indigo-600 transition-colors">
                <MailIcon class="w-4 h-4" />
              </div>
              <div class="text-left overflow-hidden">
                <p class="text-[9px] font-black text-slate-300 uppercase tracking-widest mb-0.5">Email</p>
                <p class="text-xs font-bold text-slate-700 truncate">{{ user.email || '-' }}</p>
              </div>
            </div>
            <div class="flex items-center gap-4 group">
              <div class="p-3 bg-slate-50 text-slate-400 rounded-2xl group-hover:bg-emerald-50 group-hover:text-emerald-600 transition-colors">
                <PhoneIcon class="w-4 h-4" />
              </div>
              <div class="text-left overflow-hidden">
                <p class="text-[9px] font-black text-slate-300 uppercase tracking-widest mb-0.5">WhatsApp</p>
                <p class="text-xs font-bold text-slate-700 truncate">+{{ user.phone_number || '-' }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Role Badge Info -->
        <div class="p-8 rounded-[2.5rem] border border-dashed border-slate-200 text-center">
          <p class="text-[10px] font-bold text-slate-400 leading-relaxed uppercase tracking-widest">
            Akses Akun Dibuat Pada:<br>
            <span class="text-slate-700 font-black">{{ formattedDate(user.created_at) }}</span>
          </p>
        </div>
      </div>

      <!-- Right Column: Details & Stats -->
      <div class="lg:col-span-2 space-y-8">
        <!-- Parent Specific Details Card -->
        <div v-if="user.role === 'parent'" class="white-card p-10 !rounded-[3rem] shadow-[0_20px_50px_rgba(0,0,0,0.03)] border-slate-100/50">
          <div class="flex items-center gap-4 mb-10">
            <div class="p-3 bg-indigo-50 text-indigo-600 rounded-2xl">
              <UserIcon class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-lg font-black text-slate-800 tracking-tight">Data Personal Wali</h3>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Biodata lengkap sesuai kartu identitas</p>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-y-10 gap-x-12">
            <div class="space-y-1.5">
              <label class="text-[10px] font-black text-slate-300 uppercase tracking-widest flex items-center gap-2">
                <ShieldCheckIcon class="w-3 h-3" /> NIK Wali
              </label>
              <p class="text-sm font-black text-slate-700">{{ user.nik || '-' }}</p>
            </div>
            <div class="space-y-1.5">
              <label class="text-[10px] font-black text-slate-300 uppercase tracking-widest flex items-center gap-2">
                <CalendarIcon class="w-3 h-3" /> Tanggal Lahir
              </label>
              <p class="text-sm font-black text-slate-700">{{ formattedDate(user.birth_date) }}</p>
            </div>
            <div class="space-y-1.5 md:col-span-2">
              <label class="text-[10px] font-black text-slate-300 uppercase tracking-widest flex items-center gap-2">
                <MapPinIcon class="w-3 h-3" /> Alamat Lengkap
              </label>
              <p class="text-sm font-bold text-slate-700 leading-relaxed">{{ user.address || '-' }}</p>
            </div>
            <div class="space-y-1.5">
              <label class="text-[10px] font-black text-slate-300 uppercase tracking-widest flex items-center gap-2">
                <EducationIcon class="w-3 h-3" /> Pendidikan
              </label>
              <p class="text-sm font-black text-slate-700">{{ user.education || '-' }}</p>
            </div>
            <div class="space-y-1.5">
              <label class="text-[10px] font-black text-slate-300 uppercase tracking-widest flex items-center gap-2">
                <OccupationIcon class="w-3 h-3" /> Pekerjaan
              </label>
              <p class="text-sm font-black text-slate-700">{{ user.occupation || '-' }}</p>
            </div>
            <div class="space-y-1.5">
              <label class="text-[10px] font-black text-slate-300 uppercase tracking-widest flex items-center gap-2">
                <IncomeIcon class="w-3 h-3" /> Penghasilan
              </label>
              <p class="text-sm font-black text-slate-700">{{ user.income || '-' }}</p>
            </div>
          </div>
        </div>

        <!-- Linked Students Section -->
        <div v-if="user.role === 'parent'" class="white-card p-10 !rounded-[3rem] shadow-[0_20px_50px_rgba(0,0,0,0.03)] border-slate-100/50">
          <div class="flex items-center gap-4 mb-8">
            <div class="p-3 bg-emerald-50 text-emerald-600 rounded-2xl">
              <StudentsIcon class="w-5 h-5" />
            </div>
            <div>
              <h3 class="text-lg font-black text-slate-800 tracking-tight">Daftar Siswa Terhubung</h3>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Anak/Wali yang terdaftar di sistem</p>
            </div>
          </div>

          <div v-if="students.length > 0" class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div v-for="student in students" :key="student.id" 
              @click="router.push({ name: 'student-details', params: { id: student.id } })"
              class="p-6 bg-slate-50/50 border border-slate-100 rounded-[2rem] flex items-center gap-5 hover:bg-white hover:shadow-xl hover:shadow-indigo-50 transition-all group cursor-pointer"
            >
              <div class="w-12 h-12 bg-white rounded-2xl flex items-center justify-center text-slate-300 group-hover:text-indigo-600 transition-colors border border-slate-100 shadow-sm">
                <UserIcon class="w-6 h-6" />
              </div>
              <div class="overflow-hidden">
                <h4 class="text-xs font-black text-slate-700 uppercase tracking-tight truncate">{{ student.name }}</h4>
                <p class="text-[9px] font-bold text-slate-400 mt-1 uppercase tracking-widest">{{ student.class_name || 'Tanpa Kelas' }} • NISN: {{ student.nisn }}</p>
                <p class="text-[8px] font-black text-indigo-400 uppercase tracking-[0.2em] mt-1 opacity-0 group-hover:opacity-100 transition-all">Klik Untuk Detail</p>
              </div>
            </div>
          </div>
          <div v-else class="py-10 text-center bg-slate-50/50 rounded-[2.5rem] border border-dashed border-slate-200">
            <p class="text-[10px] font-black text-slate-300 uppercase tracking-widest">Belum Ada Siswa Terhubung</p>
          </div>
        </div>

        <!-- Admin Info Placeholder -->
        <div v-if="user.role === 'admin'" class="white-card p-20 !rounded-[3rem] text-center border-dashed border-slate-200">
          <div class="w-20 h-20 bg-indigo-50 text-indigo-600 rounded-full flex items-center justify-center mx-auto mb-6">
            <AdminIcon class="w-10 h-10" />
          </div>
          <h3 class="text-xl font-black text-slate-800 mb-2">Akses Administrator</h3>
          <p class="text-sm text-slate-500 max-w-sm mx-auto">Akun ini memiliki hak akses penuh ke seluruh modul sistem SchoolPay. Tidak ada data personal tambahan yang tersimpan untuk role ini.</p>
        </div>
      </div>
    </div>
  </div>
  </div>
</template>

<style scoped lang="postcss">
.animate-fade-in {
  animation: fadeIn 0.5s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.animate-slide-up {
  animation: slideUp 0.6s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(30px); }
  to { opacity: 1; transform: translateY(0); }
}

.white-card {
  @apply bg-white border border-slate-100 transition-all;
}
</style>
