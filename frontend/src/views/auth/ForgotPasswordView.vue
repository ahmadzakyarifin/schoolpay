<script setup>
import { ref, onUnmounted } from 'vue'
import axios from 'axios'
import TurnstileWidget from '../../components/ui/TurnstileWidget.vue'
import { 
  Key as KeyIcon, 
  Mail as MailIcon, 
  Send as SendIcon, 
  CheckCircle as CheckCircleIcon,
  ArrowLeft as ArrowLeftIcon,
  AlertCircle as AlertIcon
} from 'lucide-vue-next'

const email = ref('')
const loading = ref(false)
const error = ref(null)
const success = ref(false)
const countdown = ref(0)
const captchaToken = ref('')
const captchaRef = ref(null)
const captchaRequired = ref(false)
const turnstileSiteKey = import.meta.env.VITE_TURNSTILE_SITE_KEY || ''
const emailError = ref('')
const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
let timer = null

const validateEmail = () => {
  const cleanEmail = email.value.trim()
  if (!cleanEmail) {
    emailError.value = 'Email wajib diisi.'
    return false
  }
  if (!emailPattern.test(cleanEmail)) {
    emailError.value = 'Format email tidak valid (contoh: user@gmail.com).'
    return false
  }
  email.value = cleanEmail.toLowerCase()
  emailError.value = ''
  return true
}

const handleEmailInput = () => {
  if (emailError.value) {
    validateEmail()
  }
}

const startCountdown = (duration) => {
  if (timer) clearInterval(timer)
  countdown.value = duration
  error.value = `Terlalu banyak request. Coba lagi dalam ${countdown.value} detik.`
  
  timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      clearInterval(timer)
      error.value = null
    } else {
      error.value = `Terlalu banyak request. Coba lagi dalam ${countdown.value} detik.`
    }
  }, 1000)
}

const handleSubmit = async () => {
  error.value = null
  if (!validateEmail()) return

  loading.value = true

  if (captchaRequired.value && !turnstileSiteKey) {
    error.value = 'Verifikasi tambahan diperlukan, tetapi site key Turnstile belum dikonfigurasi.'
    loading.value = false
    return
  }

  if (captchaRequired.value && !captchaToken.value) {
    error.value = 'Selesaikan verifikasi CAPTCHA terlebih dahulu.'
    loading.value = false
    return
  }
  
  try {
    const response = await axios.post('auth/forgot-password', { email: email.value, turnstile_token: captchaToken.value })
    if (response.data?.status === false && response.data?.data?.captcha_required === true) {
      captchaRequired.value = true
      captchaToken.value = ''
      error.value = response.data.message || 'Verifikasi tambahan diperlukan.'
      return
    }

    captchaRequired.value = false
    success.value = true
    if (timer) clearInterval(timer)
  } catch (err) {
    captchaRef.value?.reset()
    const retryAfter = Number(err.response?.data?.data?.retry_after_seconds || err.response?.headers?.['retry-after'] || 0)
    if (err.response?.status === 429 && retryAfter > 0) {
      startCountdown(retryAfter)
    } else {
      error.value = err.response?.data?.message || 'Gagal mengirim link reset'
    }
  } finally {
    loading.value = false
  }
}

const handleCaptchaError = () => {
  captchaToken.value = ''
  error.value = 'CAPTCHA gagal dimuat. Muat ulang halaman atau coba beberapa saat lagi.'
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
        <div v-if="!success">
          <div class="text-center mb-10">
            <div class="w-16 h-16 bg-indigo-600 rounded-2xl flex items-center justify-center shadow-2xl shadow-indigo-200 mx-auto mb-6 transform -rotate-3">
              <KeyIcon class="text-white w-8 h-8" />
            </div>
            <h1 class="text-2xl font-black tracking-tight text-slate-800">Reset Password</h1>
            <p class="text-[10px] font-black text-slate-400 uppercase tracking-[0.25em] mt-1">Enter your email for instructions</p>
          </div>

          <form @submit.prevent="handleSubmit" class="space-y-6">
            <div class="space-y-1.5">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest px-1">Email Address</label>
              <div class="relative group">
                <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none transition-colors group-focus-within:text-indigo-600 text-slate-300">
                  <MailIcon class="w-5 h-5" />
                </div>
                <input 
                  v-model="email" 
                  type="email" 
                  :class="['modern-input !pl-12 !h-[56px] !bg-white !rounded-xl', emailError ? '!border-rose-500 !ring-rose-50' : '']"
                  placeholder="name@schoolpay.id"
                  required
                  @input="handleEmailInput"
                  @blur="validateEmail"
                />
              </div>
              <FormError :message="emailError" />
            </div>

            <TurnstileWidget
              v-if="captchaRequired && turnstileSiteKey"
              ref="captchaRef"
              v-model="captchaToken"
              :site-key="turnstileSiteKey"
              @expired="captchaToken = ''"
              @error="handleCaptchaError"
            />

            <!-- Error Message -->
            <transition name="fade">
              <div v-if="error" class="rounded-lg border border-rose-200 bg-white px-3.5 py-3 text-rose-700 flex items-start gap-2.5">
                <AlertIcon class="mt-0.5 w-4 h-4 shrink-0 text-rose-500" />
                <p class="text-sm font-medium leading-5">{{ error }}</p>
              </div>
            </transition>

            <button 
              type="submit" 
              class="btn-primary w-full !h-[56px] !rounded-xl shadow-2xl shadow-indigo-200/50 flex items-center justify-center gap-3 active:scale-[0.98] transition-all group"
              :disabled="loading || countdown > 0"
            >
              <div v-if="loading" class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              <span class="uppercase text-xs font-black tracking-widest">{{ loading ? 'Sending...' : 'Send Reset Link' }}</span>
              <SendIcon v-if="!loading" class="w-5 h-5 group-hover:translate-x-1 transition-transform" />
            </button>

            <div class="text-center">
              <router-link to="/" class="inline-flex items-center gap-2 text-[10px] font-black text-slate-400 hover:text-indigo-600 uppercase tracking-widest transition-colors">
                <ArrowLeftIcon class="w-3 h-3" />
                <span>Back to Login</span>
              </router-link>
            </div>
          </form>
        </div>

        <div v-else class="text-center py-6 animate-fade-in">
          <div class="w-20 h-20 bg-emerald-600 text-white rounded-3xl flex items-center justify-center mx-auto mb-8 shadow-2xl shadow-emerald-200">
            <CheckCircleIcon class="w-10 h-10" />
          </div>
          <h2 class="text-2xl font-black text-slate-800 mb-4">Check Your Inbox</h2>
          <p class="text-sm font-medium text-slate-500 mb-10 leading-relaxed">
            We've sent a password reset link to your email. Please check your spam folder if you can't find it.
          </p>
          <router-link to="/" class="btn-primary w-full py-4 flex items-center justify-center gap-3">
            <span class="uppercase text-xs font-black tracking-widest">Return to Login</span>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
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

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
