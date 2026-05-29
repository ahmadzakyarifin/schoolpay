<script setup>
import { onMounted, onUnmounted, ref } from 'vue'
import { MessageCircle, Send, CheckCircle2, UserCheck, RefreshCw } from 'lucide-vue-next'
import supportService from '../../services/support.service'
import { useToast } from '../../composables/useToast'

const toast = useToast()
const conversations = ref([])
const messages = ref([])
const selected = ref(null)
const loading = ref(false)
const sending = ref(false)
const draft = ref('')
const status = ref('')

const loadConversations = async () => {
  loading.value = true
  try {
    const res = await supportService.getConversations({ status: status.value, limit: 50 })
    conversations.value = res.data.data.data || []
  } catch (err) {
    toast.error('Gagal memuat chat CS', err.response?.data?.message || 'Server tidak merespon')
  } finally {
    loading.value = false
  }
}

const openConversation = async (conversation) => {
  selected.value = conversation
  try {
    const res = await supportService.getMessages(conversation.id)
    messages.value = res.data.data || []
    conversation.unread_count = 0
  } catch (err) {
    toast.error('Gagal memuat pesan', err.response?.data?.message || 'Server tidak merespon')
  }
}

const reply = async () => {
  if (!selected.value || !draft.value.trim()) return
  sending.value = true
  try {
    await supportService.reply(selected.value.id, draft.value.trim())
    draft.value = ''
    await openConversation(selected.value)
    await loadConversations()
    toast.success('Balasan terkirim', 'Pesan dikirim lewat nomor WhatsApp sekolah')
  } catch (err) {
    toast.error('Gagal mengirim balasan', err.response?.data?.message || 'WAHA/server tidak merespon')
  } finally {
    sending.value = false
  }
}

const assign = async () => {
  if (!selected.value) return
  await supportService.assign(selected.value.id)
  await loadConversations()
  toast.success('Percakapan diambil', 'Chat masuk ke antrian Anda')
}

const closeConversation = async () => {
  if (!selected.value) return
  await supportService.close(selected.value.id)
  selected.value.status = 'closed'
  await loadConversations()
  toast.success('Percakapan ditutup', 'Ticket CS selesai')
}

const handleSupportUpdate = async () => {
  await loadConversations()
  if (selected.value) {
    await openConversation(selected.value)
  }
}

onMounted(() => {
  loadConversations()
  window.addEventListener('support-chat-updated', handleSupportUpdate)
})

onUnmounted(() => {
  window.removeEventListener('support-chat-updated', handleSupportUpdate)
})
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-slate-800 tracking-tight">CS WhatsApp Inbox</h1>
        <p class="text-xs font-bold text-slate-400 uppercase tracking-widest mt-1">Balas parent dari dashboard tanpa membuka WhatsApp Web</p>
      </div>
      <button @click="loadConversations" class="btn-secondary flex items-center gap-2 !rounded-2xl">
        <RefreshCw class="w-4 h-4" /> Refresh
      </button>
    </div>

    <div class="grid grid-cols-1 xl:grid-cols-[360px_1fr] gap-6 min-h-[70vh]">
      <section class="bg-white border border-slate-100 rounded-2xl overflow-hidden shadow-sm">
        <div class="p-4 border-b border-slate-100 flex items-center justify-between">
          <div class="flex items-center gap-2 font-black text-slate-700 text-sm uppercase tracking-widest">
            <MessageCircle class="w-4 h-4 text-indigo-600" /> Antrian CS
          </div>
          <select v-model="status" @change="loadConversations" class="text-xs font-bold bg-slate-50 border border-slate-100 rounded-xl px-3 py-2 outline-none">
            <option value="">Semua</option>
            <option value="open">Open</option>
            <option value="pending">Pending</option>
            <option value="closed">Closed</option>
          </select>
        </div>
        <div class="divide-y divide-slate-100 max-h-[70vh] overflow-y-auto">
          <button v-for="item in conversations" :key="item.id" @click="openConversation(item)" :class="['w-full text-left p-4 hover:bg-indigo-50/40 transition-all', selected?.id === item.id ? 'bg-indigo-50' : 'bg-white']">
            <div class="flex items-start justify-between gap-3">
              <div class="min-w-0">
                <p class="font-black text-sm text-slate-800 truncate">{{ item.parent_name || item.phone_number }}</p>
                <p class="text-[10px] font-bold text-slate-400 truncate mt-1">{{ item.last_message || 'Belum ada pesan' }}</p>
              </div>
              <span v-if="item.unread_count" class="bg-rose-500 text-white text-[10px] font-black rounded-full min-w-5 h-5 px-1.5 flex items-center justify-center">{{ item.unread_count }}</span>
            </div>
            <div class="mt-3 flex items-center justify-between">
              <span class="text-[9px] font-black uppercase tracking-widest px-2 py-1 rounded-lg" :class="item.status === 'closed' ? 'bg-slate-100 text-slate-500' : 'bg-emerald-50 text-emerald-600'">{{ item.status }}</span>
              <span class="text-[9px] font-bold text-slate-400">{{ item.phone_number }}</span>
            </div>
          </button>
          <div v-if="!loading && conversations.length === 0" class="p-10 text-center text-xs font-bold text-slate-400">Belum ada chat CS.</div>
        </div>
      </section>

      <section class="bg-white border border-slate-100 rounded-2xl overflow-hidden shadow-sm flex flex-col">
        <template v-if="selected">
          <div class="p-5 border-b border-slate-100 flex items-center justify-between gap-4">
            <div>
              <h2 class="font-black text-slate-800">{{ selected.parent_name || selected.phone_number }}</h2>
              <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">{{ selected.phone_number }}</p>
            </div>
            <div class="flex items-center gap-2">
              <button @click="assign" class="btn-secondary !rounded-xl !px-4 flex items-center gap-2"><UserCheck class="w-4 h-4" /> Ambil</button>
              <button @click="closeConversation" class="btn-secondary !rounded-xl !px-4 flex items-center gap-2"><CheckCircle2 class="w-4 h-4" /> Tutup</button>
            </div>
          </div>

          <div class="flex-1 p-6 space-y-3 overflow-y-auto bg-slate-50/50">
            <div v-for="msg in messages" :key="msg.id" :class="['flex', msg.sender_type === 'admin' ? 'justify-end' : 'justify-start']">
              <div :class="['max-w-[75%] rounded-2xl px-4 py-3 text-sm font-semibold shadow-sm', msg.sender_type === 'admin' ? 'bg-indigo-600 text-white' : 'bg-white border border-slate-100 text-slate-700']">
                <p>{{ msg.message }}</p>
                <p :class="['text-[9px] uppercase tracking-widest mt-2', msg.sender_type === 'admin' ? 'text-indigo-100' : 'text-slate-400']">{{ msg.sender_type }} • {{ msg.delivery_status }}</p>
              </div>
            </div>
          </div>

          <div class="p-5 border-t border-slate-100 flex items-center gap-3">
            <input v-model="draft" @keyup.enter="reply" :disabled="selected.status === 'closed'" class="input-premium flex-1" placeholder="Tulis balasan CS..." />
            <button @click="reply" :disabled="sending || !draft.trim() || selected.status === 'closed'" class="btn-primary !rounded-2xl flex items-center gap-2 disabled:opacity-50">
              <Send class="w-4 h-4" /> Kirim
            </button>
          </div>
        </template>
        <div v-else class="flex-1 flex items-center justify-center text-center p-12 text-slate-400">
          <div>
            <MessageCircle class="w-12 h-12 mx-auto mb-4 opacity-40" />
            <p class="font-black uppercase tracking-widest text-xs">Pilih chat dari antrian</p>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped lang="postcss">
.input-premium { @apply w-full py-4 px-5 bg-slate-50 border-2 border-slate-100 rounded-2xl focus:bg-white focus:border-indigo-500 outline-none font-bold text-xs transition-all shadow-sm; }
</style>
