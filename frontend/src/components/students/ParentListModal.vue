<script setup>
import { 
  X as CloseIcon,
  Users as ParentsIcon,
  Mail as MailIcon,
  Phone as PhoneIcon,
  MessageCircle as WAIcon,
  User as UserIcon,
  ArrowRight as ArrowRightIcon
} from 'lucide-vue-next'

const props = defineProps({
  show: Boolean,
  student: Object,
  parents: Array,
  loading: Boolean
})

const emit = defineEmits(['close'])
</script>

<template>
  <transition enter-active-class="transition duration-300 ease-out" enter-from-class="opacity-0" enter-to-class="opacity-100" leave-active-class="transition duration-200 ease-in" leave-from-class="opacity-100" leave-to-class="opacity-0">
    <div v-if="show" class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-slate-900/60 backdrop-blur-sm">
      <div class="bg-white w-full max-w-xl rounded-[40px] shadow-2xl flex flex-col overflow-hidden animate-scale-in">
        <!-- Header -->
        <div class="px-8 py-8 border-b border-slate-50 flex items-center justify-between bg-slate-50/50">
          <div class="flex items-center gap-4">
            <div class="w-12 h-12 rounded-2xl bg-indigo-600 flex items-center justify-center shadow-lg shadow-indigo-100">
              <ParentsIcon class="w-6 h-6 text-white" />
            </div>
            <div>
              <h2 class="text-lg font-black text-slate-800 tracking-tight">Wali Murid</h2>
              <p class="text-[10px] text-slate-400 font-bold uppercase tracking-widest mt-0.5">Siswa: {{ student?.name }}</p>
            </div>
          </div>
          <button @click="$emit('close')" class="w-10 h-10 rounded-xl hover:bg-white hover:shadow-xl hover:shadow-slate-200/50 flex items-center justify-center text-slate-400 hover:text-rose-500 transition-all">
            <CloseIcon class="w-5 h-5" />
          </button>
        </div>

        <div class="p-8">
          <div v-if="loading" class="flex flex-col items-center justify-center py-12 gap-4">
            <div class="w-10 h-10 border-4 border-indigo-100 border-t-indigo-600 rounded-full animate-spin"></div>
            <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest">Mengambil data...</span>
          </div>

          <div v-else-if="parents.length === 0" class="flex flex-col items-center justify-center py-12 gap-4 bg-slate-50/50 rounded-[32px] border border-dashed border-slate-200">
            <div class="w-16 h-16 rounded-full bg-white flex items-center justify-center text-slate-200 border border-slate-100">
              <ParentsIcon class="w-8 h-8" />
            </div>
            <div class="text-center">
              <p class="text-sm font-bold text-slate-500">Belum Ada Data Orang Tua</p>
              <p class="text-[10px] text-slate-400 font-medium mt-1">Gunakan fitur edit siswa untuk menghubungkan wali</p>
            </div>
          </div>

          <div v-else class="space-y-4">
            <div v-for="p in parents" :key="p.id" class="p-6 bg-white border border-slate-100 rounded-[32px] shadow-sm hover:shadow-xl hover:shadow-indigo-100/20 hover:border-indigo-100 transition-all duration-500 group">
              <div class="flex items-center gap-5">
                <div class="w-16 h-16 rounded-2xl bg-slate-50 border border-slate-100 flex items-center justify-center overflow-hidden group-hover:bg-indigo-50 transition-colors">
                  <img v-if="p.image_path" :src="p.image_path" class="w-full h-full object-cover" />
                  <UserIcon v-else class="w-7 h-7 text-slate-300 group-hover:text-indigo-300 transition-colors" />
                </div>
                <div class="flex-1">
                  <div class="flex items-center justify-between">
                    <span class="text-[10px] font-black text-indigo-500 uppercase tracking-widest">Wali Utama</span>
                    <span class="px-2 py-0.5 bg-emerald-50 text-emerald-600 text-[8px] font-black uppercase rounded border border-emerald-100">Aktif</span>
                  </div>
                  <h4 class="text-md font-black text-slate-800 capitalize mt-1">{{ p.name }}</h4>
                  <p class="text-[10px] text-slate-400 font-bold uppercase tracking-tighter mt-0.5">{{ p.occupation || 'Wiraswasta' }}</p>
                </div>
              </div>

              <div class="grid grid-cols-2 gap-3 mt-6 pt-6 border-t border-slate-50">
                <a :href="`tel:${p.phone_number}`" class="flex items-center gap-3 p-3 rounded-2xl bg-slate-50 hover:bg-indigo-50 text-slate-400 hover:text-indigo-600 transition-all border border-transparent hover:border-indigo-100">
                  <PhoneIcon class="w-4 h-4" />
                  <span class="text-[10px] font-black uppercase tracking-widest">Telepon</span>
                </a>
                <a :href="`https://wa.me/${p.phone_number}`" target="_blank" class="flex items-center gap-3 p-3 rounded-2xl bg-slate-50 hover:bg-emerald-50 text-slate-400 hover:text-emerald-600 transition-all border border-transparent hover:border-emerald-100">
                  <WAIcon class="w-4 h-4" />
                  <span class="text-[10px] font-black uppercase tracking-widest">WhatsApp</span>
                </a>
              </div>
            </div>
          </div>
        </div>

        <div class="px-8 py-6 bg-slate-50/50 border-t border-slate-50 flex justify-center">
          <button @click="$emit('close')" class="text-[10px] font-black text-slate-400 hover:text-slate-600 uppercase tracking-[0.2em] transition-all">Tutup Jendela</button>
        </div>
      </div>
    </div>
  </transition>
</template>
