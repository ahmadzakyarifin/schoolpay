<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import { useForm } from '../../composables/useForm'
import TurnstileWidget from '../../components/ui/TurnstileWidget.vue'
import { 
  Lock as LockIcon, 
  ShieldCheck as ShieldCheckIcon, 
  CheckCircle as CheckCircleIcon,
  ArrowRight as ArrowRightIcon,
  AlertCircle as AlertCircleIcon
} from 'lucide-vue-next'

const route = useRoute()
const loading = ref(false)
const success = ref(false)
const captchaToken = ref('')
const captchaRef = ref(null)
const captchaRequired = ref(false)
const turnstileSiteKey = import.meta.env.VITE_TURNSTILE_SITE_KEY || ''

const { form, errors, setErrors, clearErrors } = useForm({
  password: '',
  confirmPassword: ''
})

const setResetPasswordErrors = (err) => {
  setErrors(err)
  if (errors.value.confirm_password && !errors.value.confirmPassword) {
    errors.value = {
      ...errors.value,
      confirmPassword: errors.value.confirm_password
    }
    const normalized = { ...errors.value }
    delete normalized.confirm_password
    errors.value = normalized
  }
}

const handleSubmit = async () => {
  clearErrors()
  const validationErrors = {}
  
  if (!form.password) {
    validationErrors.password = ['Password wajib diisi.']
  } else if (form.password.length < 6) {
    validationErrors.password = ['Password minimal 6 karakter.']
  }
  
  if (!form.confirmPassword) {
    validationErrors.confirmPassword = ['Konfirmasi password wajib diisi.']
  } else if (form.password !== form.confirmPassword) {
    validationErrors.confirmPassword = ['Password konfirmasi tidak cocok.']
  }
  
  if (Object.keys(validationErrors).length > 0) {
    setResetPasswordErrors({ response: { data: { errors: validationErrors } } })
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

  loading.value = true
  
  try {
    const token = route.query.token
    const response = await axios.post('auth/reset-password', {
      token, 
      password: form.password,
      confirm_password: form.confirmPassword,
      turnstile_token: captchaToken.value
    })
    if (response.data?.status === false && response.data?.data?.captcha_required === true) {
      captchaRequired.value = true
      captchaToken.value = ''
      errors.value = { _general: [response.data.message || 'Verifikasi tambahan diperlukan.'] }
      return
    }

    captchaRequired.value = false
    success.value = true
  } catch (err) {
    captchaRef.value?.reset()
    setResetPasswordErrors(err)
  } finally {
    loading.value = false
  }
}

const handleCaptchaError = () => {
  captchaToken.value = ''
  errors.value = { _general: ['CAPTCHA gagal dimuat. Muat ulang halaman atau coba beberapa saat lagi.'] }
}
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
            <div class="w-16 h-16 bg-indigo-600 rounded-2xl flex items-center justify-center shadow-2xl shadow-indigo-200 mx-auto mb-6 transform rotate-3">
              <LockIcon class="text-white w-8 h-8" />
            </div>
            <h1 class="text-2xl font-black tracking-tight text-slate-800">New Password</h1>
            <p class="text-[10px] font-black text-slate-400 uppercase tracking-[0.25em] mt-1">Reset your access credentials</p>
          </div>

          <form @submit.prevent="handleSubmit" class="space-y-5">
            <div class="space-y-1.5">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest px-1">New Password <span class="text-red-500">*</span></label>
              <div class="relative group">
                <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none transition-colors group-focus-within:text-indigo-600 text-slate-300">
                  <LockIcon class="w-5 h-5" />
                </div>
                <input 
                  v-model="form.password" 
                  type="password" 
                  :class="['modern-input !pl-12 !h-[56px] !bg-white !rounded-xl', errors.password ? '!border-rose-500 !ring-rose-50' : '']"
                  placeholder="••••••••"
                />
              </div>
              <FormError :message="errors.password" />
            </div>

            <div class="space-y-1.5">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest px-1">Confirm Password <span class="text-red-500">*</span></label>
              <div class="relative group">
                <div class="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none transition-colors group-focus-within:text-indigo-600 text-slate-300">
                  <ShieldCheckIcon class="w-5 h-5" />
                </div>
                <input 
                  v-model="form.confirmPassword" 
                  type="password" 
                  :class="['modern-input !pl-12 !h-[56px] !bg-white !rounded-xl', errors.confirmPassword ? '!border-rose-500 !ring-rose-50' : '']"
                  placeholder="••••••••"
                />
              </div>
              <FormError :message="errors.confirmPassword" />
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
              <div v-if="errors._general" class="rounded-lg border border-rose-200 bg-white px-3.5 py-3 text-rose-700 flex items-start gap-2.5">
                <AlertCircleIcon class="mt-0.5 w-4 h-4 shrink-0 text-rose-500" />
                <p class="text-sm font-medium leading-5">{{ errors._general[0] }}</p>
              </div>
            </transition>

            <button 
              type="submit" 
              class="btn-primary w-full !h-[56px] !rounded-xl shadow-2xl shadow-indigo-200/50 flex items-center justify-center gap-3 active:scale-[0.98] transition-all group"
              :disabled="loading"
            >
              <div v-if="loading" class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
              <span class="uppercase text-xs font-black tracking-widest">{{ loading ? 'Updating...' : 'Update Password' }}</span>
              <ArrowRightIcon v-if="!loading" class="w-5 h-5 group-hover:translate-x-1 transition-transform" />
            </button>
          </form>
        </div>

        <div v-else class="text-center py-6 animate-fade-in">
          <div class="w-20 h-20 bg-emerald-600 text-white rounded-3xl flex items-center justify-center mx-auto mb-8 shadow-2xl shadow-emerald-200">
            <CheckCircleIcon class="w-10 h-10" />
          </div>
          <h2 class="text-2xl font-black text-slate-800 mb-4">Password Updated!</h2>
          <p class="text-sm font-medium text-slate-500 mb-10 leading-relaxed">Your password has been reset successfully. You can now log in with your new password.</p>
          
          <router-link to="/" class="btn-primary w-full py-4 flex items-center justify-center gap-3">
            <span class="uppercase text-xs font-black tracking-widest">Login Now</span>
            <ArrowRightIcon class="w-4 h-4" />
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
