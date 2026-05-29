<script setup>
import { ref, computed } from 'vue'
import { Play as PlayIcon, Undo2 as UndoIcon, MessageSquare as MsgIcon, Check as CheckIcon, Edit3 as EditIcon, Smartphone as PhoneIcon } from 'lucide-vue-next'

const props = defineProps({
  modelValue: Boolean,
  generateActionType: String,
  selectedCount: Number,
  generateSubmitting: Boolean,
  isPenyesuaian: Boolean
})

const emit = defineEmits(['update:modelValue', 'confirm'])

const customReason = ref('')
const customMessageOverride = ref('')
const isManualOverride = ref(false)
const selectedTemplate = ref('standard')
const skipNotification = ref(false)

const defaultPreview = computed(() => {
  if (!props.isPenyesuaian) {
    let base = `📢 *PEMBERITAHUAN RESMI SCHOOLPAY*\n*Tagihan Baru Diterbitkan*\n\n` +
      `Yth. Orang Tua / Wali Siswa,\n\n` +
      `Melalui pesan ini, kami menginformasikan bahwa terdapat tagihan baru yang telah diterbitkan pada sistem untuk putra/putri Anda.\n\n`
    
    if (customReason.value) {
      base += `*Keterangan Tambahan:*\n_${customReason.value}_\n\n`
    }
    
    base += `Tagihan tersebut kini sudah tersedia di akun SchoolPay Anda. Pembayaran dapat dilakukan melalui portal orang tua. Terima kasih atas perhatian dan kerja sama Bapak/Ibu.`
    return base
  } else {
    let base = `📢 *PEMBERITAHUAN RESMI SCHOOLPAY*\n*Penyesuaian Tarif Tagihan Sekolah*\n\n` +
      `Yth. Orang Tua / Wali Siswa,\n\n` +
      `Melalui pesan ini, kami menginformasikan adanya penyesuaian kebijakan tarif tagihan sekolah dengan rincian penyesuaian yang telah diterbitkan pada sistem.\n\n`
    
    if (customReason.value) {
      base += `*Mengapa ada penyesuaian ini?*\n_${customReason.value}_\n\n`
    }
    
    base += `Sistem kami telah menerbitkan tagihan penyesuaian pada akun SchoolPay Anda. Pembayaran dapat dilakukan melalui portal orang tua. Terima kasih atas perhatian dan kerja sama Bapak/Ibu.`
    return base
  }
})

const activePreview = computed(() => {
  return isManualOverride.value ? (customMessageOverride.value || 'Ketik pesan manual Anda di sini...') : defaultPreview.value
})

const close = () => {
  emit('update:modelValue', false)
}

const confirmAction = () => {
  emit('confirm', {
    customReason: customReason.value,
    customMessage: isManualOverride.value ? customMessageOverride.value : '',
    skipNotification: skipNotification.value
  })
}
</script>

<template>
  <Teleport to="body">
    <transition name="fade">
      <div v-if="modelValue" class="fixed inset-0 z-[700] flex items-center justify-center p-6">
        <div class="absolute inset-0 bg-slate-900/60 backdrop-blur-sm" @click="close"></div>
        <div class="white-card w-full relative z-10 overflow-hidden shadow-[0_20px_50px_rgba(0,0,0,0.3)] animate-scale-in !rounded-[2.5rem] p-8 bg-white font-inter"
             :class="generateActionType !== 'bulk_cancel' ? 'max-w-4xl' : 'max-w-md text-center'">
          
          <!-- Bulk Cancel View -->
          <div v-if="generateActionType === 'bulk_cancel'">
            <div class="w-20 h-20 bg-amber-50 text-amber-600 rounded-[2rem] flex items-center justify-center mx-auto mb-6 border border-amber-100 shadow-xl shadow-amber-600/10">
              <UndoIcon class="w-10 h-10" />
            </div>
            <h3 class="text-xl font-black text-slate-900 tracking-tight mb-2">Tarik Tagihan Siswa?</h3>
            <p class="text-slate-500 text-[10px] font-bold uppercase tracking-widest mb-8 px-4 leading-relaxed">
              Sistem akan menarik tagihan dari {{ selectedCount }} aturan terpilih. Jika sudah ada pembayaran, dana akan dipindahkan ke saldo deposit siswa dan riwayat transaksi tetap aman.
            </p>
            <div class="text-left space-y-2 mb-5">
              <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Alasan Penarikan (Opsional)</label>
              <textarea v-model="customReason" rows="3" class="w-full p-4 bg-amber-50/60 border border-amber-100 rounded-2xl text-xs font-bold text-slate-700 outline-none focus:bg-white focus:border-amber-400" placeholder="Contoh: Koreksi aturan tagihan karena nominal/periode salah."></textarea>
              <label class="flex items-center gap-2 text-[10px] font-black text-rose-600 uppercase tracking-wider cursor-pointer">
                <input v-model="skipNotification" type="checkbox" class="rounded border-rose-200 text-rose-500" />
                Jangan kirim notifikasi ke orang tua
              </label>
            </div>
            <div class="grid grid-cols-2 gap-4">
              <button @click="close" class="py-4 bg-slate-100 text-slate-600 font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-slate-200 transition-all cursor-pointer">Batalkan</button>
              <button @click="confirmAction" :disabled="generateSubmitting" class="py-4 bg-amber-600 hover:bg-amber-700 text-white font-black rounded-2xl text-[10px] uppercase tracking-widest transition-all shadow-lg shadow-amber-600/20 disabled:opacity-50 cursor-pointer">
                {{ generateSubmitting ? 'Memproses...' : 'Ya, Tarik Tagihan' }}
              </button>
            </div>
          </div>

          <!-- Generate View with Broadcast Editor -->
          <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-8 items-start text-left">
            <!-- Left Column: Editor Controls -->
            <div class="flex flex-col gap-6">
              <div class="flex items-center gap-4 border-b border-slate-100 pb-4">
                <div class="w-12 h-12 bg-indigo-50 text-indigo-600 rounded-2xl flex items-center justify-center border border-indigo-100 shadow-md shadow-indigo-600/10">
                  <PlayIcon class="w-6 h-6 fill-current" />
                </div>
                <div>
                  <h3 class="text-lg font-black text-slate-900 tracking-tight">
                    {{ generateActionType === 'single' ? (isPenyesuaian ? 'Generate & Siaran Penyesuaian' : 'Generate Tagihan Baru & Siaran') : 'Generate Masal & Siaran' }}
                  </h3>
                  <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Atur pesan notifikasi WhatsApp & Email</p>
                </div>
              </div>

              <!-- Template Selector -->
              <div class="flex flex-col gap-2">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest flex items-center gap-1.5">
                  <MsgIcon class="w-3.5 h-3.5 text-indigo-500" /> Mode Pesan Siaran
                </label>
                <div class="grid grid-cols-2 gap-3">
                  <button @click="isManualOverride = false" class="p-3.5 rounded-2xl border text-xs font-bold flex items-center gap-2 transition-all cursor-pointer"
                          :class="!isManualOverride ? 'border-indigo-600 bg-indigo-50/50 text-indigo-600 shadow-sm' : 'border-slate-200 bg-white text-slate-500 hover:bg-slate-50'">
                    <CheckIcon v-if="!isManualOverride" class="w-4 h-4 text-indigo-600" />
                    <span>Template Otomatis</span>
                  </button>
                  <button @click="isManualOverride = true" class="p-3.5 rounded-2xl border text-xs font-bold flex items-center gap-2 transition-all cursor-pointer"
                          :class="isManualOverride ? 'border-indigo-600 bg-indigo-50/50 text-indigo-600 shadow-sm' : 'border-slate-200 bg-white text-slate-500 hover:bg-slate-50'">
                    <CheckIcon v-if="isManualOverride" class="w-4 h-4 text-indigo-600" />
                    <span>Teks Manual (Bebas)</span>
                  </button>
                </div>
              </div>

              <!-- Custom Reason (If Template Mode) -->
              <div v-if="!isManualOverride" class="flex flex-col gap-2 animate-fade-in">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest flex items-center justify-between">
                  <span>Alasan Penyesuaian (Opsional)</span>
                  <span class="text-[9px] text-slate-400 font-normal">Akan disisipkan ke pesan</span>
                </label>
                <textarea v-model="customReason" rows="3" placeholder="Contoh: Terdapat penambahan fasilitas paket kelulusan berupa plakat eksklusif dan dokumentasi video 4K..."
                          class="w-full p-4 bg-slate-50 border border-slate-200 rounded-2xl text-xs font-bold text-slate-700 outline-none focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-50 shadow-sm transition-all"></textarea>
              </div>

              <!-- Manual Message Override (If Manual Mode) -->
              <div v-else class="flex flex-col gap-2 animate-fade-in">
                <label class="text-[10px] font-black text-slate-400 uppercase tracking-widest flex items-center justify-between">
                  <span>Isi Pesan Siaran Lengkap</span>
                  <span class="text-[9px] text-slate-400 font-normal">Tulis pesan kustom Anda</span>
                </label>
                <textarea v-model="customMessageOverride" rows="5" placeholder="Ketik pengumuman lengkap Anda di sini..."
                          class="w-full p-4 bg-slate-50 border border-slate-200 rounded-2xl text-xs font-bold text-slate-700 outline-none focus:bg-white focus:border-indigo-500 focus:ring-4 focus:ring-indigo-50 shadow-sm transition-all"
                          :disabled="skipNotification"></textarea>
              </div>

              <!-- Skip Notification Toggle -->
              <div class="flex items-center gap-3 p-4 bg-rose-50/50 border border-rose-100 rounded-2xl">
                <button 
                  @click="skipNotification = !skipNotification"
                  class="relative w-10 h-5 rounded-full transition-all duration-300 focus:outline-none shadow-inner shrink-0 cursor-pointer"
                  :class="skipNotification ? 'bg-rose-500' : 'bg-slate-300'"
                >
                  <div class="absolute top-0.5 left-0.5 w-4 h-4 bg-white rounded-full shadow transition-transform duration-300"
                    :class="skipNotification ? 'translate-x-5' : 'translate-x-0'"></div>
                </button>
                <div class="flex flex-col cursor-pointer" @click="skipNotification = !skipNotification">
                  <span class="text-xs font-black text-rose-700 uppercase tracking-widest">Generate Tanpa Notifikasi</span>
                  <span class="text-[9px] font-bold text-rose-500">Jangan kirim pesan WA/Email ke orang tua (Bisu)</span>
                </div>
              </div>

              <!-- Action Buttons -->
              <div class="grid grid-cols-2 gap-4 mt-auto pt-4 border-t border-slate-100">
                <button @click="close" class="py-4 bg-slate-100 text-slate-600 font-black rounded-2xl text-[10px] uppercase tracking-widest hover:bg-slate-200 transition-all cursor-pointer text-center">Batalkan</button>
                <button @click="confirmAction" :disabled="generateSubmitting" class="py-4 bg-indigo-600 hover:bg-indigo-700 text-white font-black rounded-2xl text-[10px] uppercase tracking-widest transition-all shadow-lg shadow-indigo-600/20 disabled:opacity-50 cursor-pointer text-center">
                  {{ generateSubmitting ? 'Memproses...' : 'Kirim Siaran Sekarang' }}
                </button>
              </div>
            </div>

            <!-- Right Column: Live Phone Preview -->
            <div class="bg-slate-50 rounded-[2rem] p-6 border border-slate-100 flex flex-col items-center relative overflow-hidden h-full shadow-inner" :class="{'opacity-50 grayscale': skipNotification}">
              <div v-if="skipNotification" class="absolute inset-0 z-10 flex items-center justify-center bg-slate-100/80 backdrop-blur-sm">
                <span class="text-sm font-black text-slate-500 uppercase tracking-widest">Siaran Dimatikan</span>
              </div>
              <div class="flex items-center justify-between border-b border-slate-200/60 pb-3 w-full">
                <div class="flex items-center gap-2">
                  <PhoneIcon class="w-4 h-4 text-indigo-500" />
                  <span class="text-[10px] font-black text-slate-600 uppercase tracking-widest">Live Phone Preview</span>
                </div>
                <span class="px-2.5 py-1 bg-emerald-100 text-emerald-700 rounded-full text-[9px] font-black uppercase tracking-widest">WhatsApp & Email</span>
              </div>
              <div class="bg-white border border-slate-200/80 rounded-2xl p-5 shadow-sm flex-1 overflow-y-auto font-inter text-xs leading-relaxed text-slate-700 whitespace-pre-wrap">
                {{ activePreview }}
              </div>
              <div class="text-[9px] font-bold text-slate-400 text-center bg-slate-100/80 py-2 rounded-xl">
                💡 Variabel nama siswa & nominal akan diisi otomatis oleh sistem
              </div>
            </div>
          </div>

        </div>
      </div>
    </transition>
  </Teleport>
</template>

<style scoped lang="postcss">
.fade-enter-active, .fade-leave-active { transition: opacity 0.3s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
