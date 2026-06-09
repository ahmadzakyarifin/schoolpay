<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../../store/auth'
import { useForm } from '../../composables/useForm'
import TurnstileWidget from '../../components/ui/TurnstileWidget.vue'
import { 
  Mail as MailIcon, 
  Lock as LockIcon, 
  ArrowRight as ArrowRightIcon, 
  Eye as EyeIcon, 
  EyeOff as EyeOffIcon, 
  GraduationCap as GraduationCapIcon,
  AlertCircle as AlertIcon
} from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const showPassword = ref(false)
const reasonMessage = ref('')
const countdown = ref(0)
const captchaToken = ref('')
const captchaRef = ref(null)
const captchaRequired = ref(false)
const turnstileSiteKey = import.meta.env.VITE_TURNSTILE_SITE_KEY || ''
const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
let timer = null

onMounted(() => {
  const reason = route.query.reason
  
  // Daftar pesan teknis yang tidak perlu ditampilkan ke user
  const technicalMessages = [
    'token tidak ditemukan', 
    'format token salah', 
    'token tidak valid atau kadaluarsa',
    'expired'
  ]

  if (reason === 'forbidden') {
    reasonMessage.value = 'Akses ditolak: Akun Anda telah dinonaktifkan atau hak akses telah berubah.'
  } else if (reason && !technicalMessages.includes(reason.toLowerCase())) {
    reasonMessage.value = reason.charAt(0).toUpperCase() + reason.slice(1)
  }
})

const { form, errors, submitting, setErrors, clearErrors, clearFieldError } = useForm({
  email: '',
  password: ''
})

const setFieldError = (field, messages) => {
  errors.value = { ...errors.value, [field]: messages }
}

const validateEmailField = () => {
  const email = form.email.trim()
  if (!email) {
    setFieldError('email', ['Email wajib diisi.'])
    return false
  }
  if (!emailPattern.test(email)) {
    setFieldError('email', ['Format email tidak valid (contoh: user@gmail.com).'])
    return false
  }
  form.email = email.toLowerCase()
  clearFieldError('email')
  return true
}

const handleEmailInput = () => {
  if (errors.value.email) {
    validateEmailField()
  }
}

const handleLogin = async () => {
  clearErrors()
  reasonMessage.value = '' 
  authStore.error = null
  
  if (route.query.reason) {
    router.replace({ query: {} })
  }
  
  const validationErrors = {}
  const email = form.email.trim()
  
  // Validasi Email
  if (!email) {
    validationErrors.email = ['Email wajib diisi.']
  } else if (!emailPattern.test(email)) {
    validationErrors.email = ['Format email tidak valid (contoh: user@gmail.com).']
  } else {
    form.email = email.toLowerCase()
  }

  // Validasi Password
  if (!form.password) {
    validationErrors.password = ['Password wajib diisi.']
  } else if (form.password.length < 6) {
    validationErrors.password = ['Password minimal 6 karakter.']
  }

  if (Object.keys(validationErrors).length > 0) {
    setErrors({ response: { data: { errors: validationErrors } } })
    return
  }

  if (captchaRequired.value && !turnstileSiteKey) {
    errors.value = { _general: ['Verifikasi tambahan diperlukan, tetapi site key Turnstile belum dikonfigurasi.'] }
    return
  }

  if (captchaRequired.value && !captchaToken.value) {
    errors.value = { _general: ['Selesaikan verifikasi CAPTCHA terlebih dahulu.'] }
    return
  }

  const result = await authStore.login(form.email, form.password, captchaToken.value)
  
  if (result.success) {
    captchaRequired.value = false
    captchaToken.value = ''
    if (authStore.isAdmin) {
      router.push('/dashboard')
    } else {
      router.push('/parent/dashboard')
    }
  } else if (result.captchaRequired) {
    captchaRequired.value = true
    captchaToken.value = ''
    errors.value = { _general: [authStore.error || 'Verifikasi tambahan diperlukan.'] }
  } else {
    setErrors(result.error)
    captchaRequired.value = false
    captchaRef.value?.reset()
    form.email = ''    // Riset email
    form.password = '' // Riset password
    
    const err = result.error
    const retryAfter = Number(err?.response?.data?.data?.retry_after_seconds || err?.response?.headers?.['retry-after'] || 0)
    if (err?.response?.status === 429 && retryAfter > 0) {
      startCountdown(retryAfter)
    }
  }
}

const handleCaptchaError = () => {
  captchaToken.value = ''
  errors.value = { _general: ['CAPTCHA gagal dimuat. Muat ulang halaman atau coba beberapa saat lagi.'] }
}

const startCountdown = (duration) => {
  if (timer) clearInterval(timer)
  countdown.value = duration
  
  const updateMsg = (secs) => {
    const msg = `Terlalu banyak percobaan. Coba lagi dalam ${secs} detik.`
    authStore.error = msg
    errors.value = { _general: [msg] }
  }

  updateMsg(countdown.value)
  
  timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      clearInterval(timer)
      authStore.error = null
      errors.value = {}
    } else {
      updateMsg(countdown.value)
    }
  }, 1000)
}

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="min-h-screen flex items-center justify-center p-6 bg-slate-50 relative overflow-hidden font-inter">
    <!-- Premium Decorative background -->
    <div class="absolute -top-40 -left-40 w-80 h-80 bg-indigo-500/10 rounded-full blur-[100px]"></div>
    <div class="absolute -bottom-40 -right-40 w-80 h-80 bg-blue-500/10 rounded-full blur-[100px]"></div>

    <div class="w-full max-w-md animate-fade-in relative z-10">
      <div class="white-card p-10 md:p-12 shadow-[0_32px_64px_-12px_rgba(0,0,0,0.05)] border-slate-100/50 !rounded-xl">
        <!-- Logo Section -->
        <div class="text-center mb-10">
          <div class="w-16 h-16 bg-indigo-600 rounded-xl flex items-center justify-center shadow-2xl shadow-indigo-200 mx-auto mb-6 transform -rotate-3 transition-transform hover:rotate-0 border-4 border-indigo-400/20">
            <GraduationCapIcon class="text-white w-8 h-8" />
          </div>
          <h1 class="text-3xl font-black tracking-tight text-slate-800">SchoolPay</h1>
          <p class="text-[10px] font-black text-slate-400 uppercase tracking-[0.2em] mt-1">Industrial ecosystem</p>
        </div>

        <!-- Reason Notification -->
        <transition name="fade">
          <div v-if="reasonMessage" class="mb-6 rounded-lg border border-rose-200 bg-white px-3.5 py-3 text-rose-700 flex items-start gap-2.5">
            <AlertIcon class="mt-0.5 w-4 h-4 shrink-0 text-rose-500" />
            <p class="text-sm font-medium leading-5">{{ reasonMessage }}</p>
          </div>
        </transition>

        <form @submit.prevent="handleLogin" novalidate class="space-y-6">
          <!-- Email Input -->
          <div class="space-y-2">
            <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest px-1 flex items-center gap-1">
              Email Address <span class="text-red-500">*</span>
            </label>
            <div class="relative group">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none transition-colors group-focus-within:text-indigo-600 text-slate-300">
                <MailIcon class="w-5 h-5" />
              </div>
              <input 
                v-model="form.email" 
                type="email" 
                :class="['modern-input !pl-12 !h-[56px] !bg-white !rounded-xl', errors.email ? '!border-rose-500 !ring-rose-50' : '']"
                placeholder="name@example.com"
                @input="handleEmailInput"
                @blur="validateEmailField"
              />
            </div>
            <FormError :message="errors.email" />
          </div>

          <!-- Password Input -->
          <div class="space-y-2">
            <div class="flex justify-between items-center px-1">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest flex items-center gap-1">
                Password <span class="text-red-500">*</span>
              </label>
              <router-link to="/forgot-password" class="text-[10px] font-black text-indigo-500 hover:text-indigo-700 uppercase tracking-widest transition-colors">
                Forgot Password?
              </router-link>
            </div>
            <div class="relative group">
              <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none transition-colors group-focus-within:text-indigo-600 text-slate-300">
                <LockIcon class="w-5 h-5" />
              </div>
              <input 
                v-model="form.password" 
                :type="showPassword ? 'text' : 'password'" 
                :class="['modern-input !pl-12 !pr-12 !h-[56px] !bg-white !rounded-xl', errors.password ? '!border-rose-500 !ring-rose-50' : '']"
                placeholder="••••••••"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute right-4 top-1/2 -translate-y-1/2 text-slate-300 hover:text-slate-600 transition-colors"
                tabindex="-1"
              >
                <EyeOffIcon v-if="showPassword" class="w-4 h-4" />
                <EyeIcon v-else class="w-4 h-4" />
              </button>
            </div>
            <FormError :message="errors.password" />
          </div>

          <TurnstileWidget
            v-if="captchaRequired && turnstileSiteKey"
            ref="captchaRef"
            v-model="captchaToken"
            :site-key="turnstileSiteKey"
            @expired="captchaToken = ''"
            @error="handleCaptchaError"
          />

          <!-- General Error Message (Credential & Server Issues) -->
          <transition name="fade">
            <div v-if="errors._general || (authStore.error && !reasonMessage)" class="rounded-lg border border-rose-200 bg-white px-3.5 py-3 text-rose-700 flex items-start gap-2.5">
              <AlertIcon class="mt-0.5 w-4 h-4 shrink-0 text-rose-500" />
              <p class="text-sm font-medium leading-5">{{ errors._general ? errors._general[0] : authStore.error }}</p>
            </div>
          </transition>

          <!-- Submit Button -->
          <div class="pt-2">
            <button 
              type="submit" 
              class="btn-primary w-full !h-[56px] !rounded-xl shadow-2xl shadow-indigo-200/50 flex items-center justify-center gap-3 active:scale-[0.98] transition-all group"
              :disabled="authStore.loading || countdown > 0"
            >
              <div v-if="authStore.loading" class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              <span class="uppercase text-[13px] font-black tracking-[0.2em]">{{ authStore.loading ? 'Signing In...' : 'Sign In' }}</span>
              <ArrowRightIcon v-if="!authStore.loading" class="w-5 h-5 group-hover:translate-x-1 transition-transform" />
            </button>
          </div>
        </form>

        <!-- Footer Links -->
        <div class="mt-10 flex flex-col items-center gap-4">
          <div class="w-8 h-1 bg-slate-100 rounded-full"></div>
          
          <p class="text-[9px] font-bold text-slate-300 uppercase tracking-[0.3em]">
            Secured by SchoolPay Auth
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped lang="postcss">
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
  transform: translateY(-5px);
}

.animate-fade-in {
  animation: fadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

.animate-slide-up {
  animation: slideUp 0.3s ease-out forwards;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes slideUp {
  from { opacity: 0; transform: translateY(5px); }
  to { opacity: 1; transform: translateY(0); }
}

.modern-input {
  @apply w-full py-4 px-6 border-slate-200 focus:ring-4 focus:ring-indigo-50 focus:border-indigo-500 outline-none transition-all font-bold text-slate-700;
}
</style>
