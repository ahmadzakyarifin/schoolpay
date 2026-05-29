<template>
  <Teleport to="body">
    <transition name="modal">
      <div v-if="isOpen" class="fixed inset-0 z-[200] flex items-center justify-center p-4 sm:p-6 overflow-y-auto">
        <div class="fixed inset-0 bg-slate-900/40 backdrop-blur-sm" @click="close"></div>
        <div class="bg-white w-full max-w-4xl relative z-10 rounded-[2.5rem] shadow-2xl overflow-hidden animate-scale-in flex flex-col max-h-[90vh]">
          <!-- Header -->
          <div class="p-6 sm:p-8 bg-slate-50 border-b border-slate-100 flex items-center justify-between shrink-0">
            <div class="flex items-center gap-4">
              <div class="w-14 h-14 bg-indigo-50 text-indigo-600 rounded-2xl flex items-center justify-center shadow-inner border border-indigo-100/50">
                <DatabaseIcon class="w-7 h-7" />
              </div>
              <div>
                <h3 class="text-xl font-black text-slate-800 tracking-tight">Detail Log Audit #{{ log?.id }}</h3>
                <p class="text-xs font-bold text-slate-400 uppercase tracking-widest mt-1">
                  {{ log?.user_name }} ({{ log?.role }}) • {{ formatDate(log?.created_at) }}
                </p>
              </div>
            </div>
            <button @click="close" class="w-10 h-10 rounded-2xl bg-white border border-slate-100 flex items-center justify-center text-slate-400 hover:text-slate-600 hover:bg-slate-50 transition-all shadow-sm">
              <XIcon class="w-5 h-5" />
            </button>
          </div>

          <!-- Content -->
          <div class="p-6 sm:p-8 overflow-y-auto flex-1">
              <!-- 5W 1H Executive Summary Card -->
              <div class="bg-slate-50/70 border border-slate-100 p-6 sm:p-8 rounded-3xl space-y-6">
                <div class="flex items-center justify-between border-b border-slate-200/60 pb-4">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 bg-indigo-50 text-indigo-600 rounded-2xl flex items-center justify-center border border-indigo-100/50 font-black text-xs">
                      5W1H
                    </div>
                    <div>
                      <h4 class="text-base font-black tracking-tight text-slate-800">Rangkuman Eksekutif Investigasi</h4>
                      <p class="text-[10px] font-bold text-slate-400 uppercase tracking-widest">Metode Analisis 5W 1H (Who, What, Where, When, Why, How)</p>
                    </div>
                  </div>
                  <span class="px-3 py-1 bg-indigo-50 border border-indigo-100/50 rounded-full text-[10px] font-black uppercase tracking-widest text-indigo-600">
                    Verified Log
                  </span>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 gap-6 text-xs text-slate-600">
                  <!-- WHO -->
                  <div class="space-y-1">
                    <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest block">👤 WHO (Pelaku Tindakan)</span>
                    <p class="font-bold text-slate-800 text-sm">{{ log?.user_name }} <span class="text-xs font-normal text-slate-500">({{ log?.role }})</span></p>
                    <p class="text-[11px] text-slate-400">Pengguna yang menginisiasi atau memicu transaksi pada sistem.</p>
                  </div>

                  <!-- WHEN -->
                  <div class="space-y-1">
                    <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest block">🕒 WHEN (Waktu Kejadian)</span>
                    <p class="font-bold text-slate-800 text-sm">{{ formatDate(log?.created_at) }} <span class="text-xs font-normal text-emerald-600">(WIB)</span></p>
                    <p class="text-[11px] text-slate-400">Waktu absolut server yang disinkronisasikan pada zona waktu Asia/Jakarta.</p>
                  </div>

                  <!-- WHERE -->
                  <div class="space-y-1">
                    <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest block">📍 WHERE (Lokasi & Akses)</span>
                    <p class="font-bold text-slate-800 text-sm">{{ log?.ip_address }}</p>
                    <p class="text-[11px] text-slate-400 truncate" :title="log?.user_agent">Modul: <span class="text-indigo-600 font-medium">{{ humanizeEntity(log?.entity_type) }}</span> • Perangkat: {{ log?.user_agent || 'Sistem Internal / API' }}</p>
                  </div>

                  <!-- WHAT -->
                  <div class="space-y-1">
                    <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest block">🎯 WHAT (Aksi & Target Spesifik)</span>
                    <p class="font-bold text-slate-800 text-sm">{{ humanizeAction(log?.action) }}</p>
                    <p class="text-[11px] text-slate-400">Target Spesifik: <span class="text-slate-700 font-semibold">{{ getTargetName(log) }}</span></p>
                  </div>

                  <!-- WHY -->
                  <div class="space-y-1 md:col-span-2 bg-white p-4 rounded-2xl border border-slate-100">
                    <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-1">💡 WHY (Konteks & Tujuan Sistem)</span>
                    <p class="font-semibold text-slate-600 text-xs italic">
                      "{{ getWhyContext(log) }}"
                    </p>
                  </div>

                  <!-- HOW -->
                  <div class="space-y-1 md:col-span-2 bg-white p-4 rounded-2xl border border-slate-100">
                    <span class="text-[10px] font-black text-slate-400 uppercase tracking-widest block mb-1">🔧 HOW (Rincian Mutasi Data / Atribut)</span>
                    <p class="font-semibold text-slate-600 text-xs">
                      {{ getHowSummary(log) }}
                    </p>
                  </div>
                </div>
              </div>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup>
import { 
  Database as DatabaseIcon, 
  X as XIcon 
} from 'lucide-vue-next'

const props = defineProps({
  isOpen: { type: Boolean, default: false },
  log: { type: Object, default: null }
})

const emit = defineEmits(['close'])

const close = () => { emit('close') }

const getTargetName = (log) => {
  if (!log) return '-'
  const newVals = log.new_values || {}
  const oldVals = log.old_values || {}
  
  const name = newVals.name || oldVals.name || 
               newVals.email || oldVals.email || 
               newVals.title || oldVals.title || 
               newVals.bill_type_name || oldVals.bill_type_name || 
               newVals.client_key || oldVals.client_key ||
               newVals.channels || oldVals.channels ||
               log.user_name || 'Entitas Sistem'
               
  return `${name} (ID: #${log.entity_id})`
}

const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  return new Date(dateStr).toLocaleString('id-ID', {
    year: 'numeric', month: 'long', day: 'numeric',
    hour: '2-digit', minute: '2-digit', second: '2-digit'
  })
}

const formatJSON = (val) => {
  if (!val) return 'Tidak ada data'
  try {
    return JSON.stringify(val, null, 2)
  } catch (e) {
    return val
  }
}

const humanizeAction = (action) => {
  if (!action) return '-'
  const dict = {
    CREATE: 'Membuat Baru',
    UPDATE: 'Memperbarui Data',
    DELETE: 'Menghapus Data',
    RESTORE: 'Memulihkan Data',
    BULK_DELETE: 'Menghapus Massal',
    BULK_RESTORE: 'Memulihkan Massal',
    TOGGLE_STATUS: 'Mengubah Status',
    ACTIVATE_ACCOUNT: 'Mengaktifkan Akun',
    SEND_NOTIFICATION: 'Mengirim Notifikasi',
    LOGIN: 'Login Sistem',
    LOGIN_FAILED: 'Login Gagal',
    LOGOUT: 'Logout Sistem',
    REFRESH_TOKEN: 'Perpanjang Sesi',
    FORGOT_PASSWORD: 'Minta Reset Sandi',
    FORGOT_PASSWORD_FAILED: 'Permintaan Reset Gagal',
    RESET_PASSWORD: 'Reset Kata Sandi',
    RESET_PASSWORD_FAILED: 'Reset Kata Sandi Gagal',
    CHANGE_PASSWORD: 'Ubah Kata Sandi',
    CHANGE_PASSWORD_FAILED: 'Ubah Kata Sandi Gagal',
    EXPORT_TREND_EXCEL: 'Export Excel',
    GENERATE_BILL: 'Menerbitkan Tagihan',
    GENERATE_ADJUSTMENT_BILL: 'Menerbitkan Penyesuaian Tagihan',
    REFUND_BILL_REDUCTION: 'Refund Selisih Tagihan',
    BULK_CANCEL_GENERATED_BILL: 'Menarik Tagihan',
    VOID_BILL: 'Membatalkan Tagihan',
    EXPORT_GLOBAL_FINANCE_REPORT: 'Export Laporan Keuangan',
    CREATE_PAYMENT_INTENT: 'Membuat Tagihan Bayar',
    PROCESS_PAYMENT: 'Memproses Pembayaran',
    GATEWAY_PAYMENT_TO_DEPOSIT: 'Dana Gateway ke Saldo'
  }
  return dict[action] || action
}

const humanizeEntity = (entity) => {
  if (!entity) return '-'
  const dict = {
    users: 'Pengguna',
    students: 'Siswa',
    majors: 'Jurusan',
    classes: 'Kelas',
    academic_years: 'Tahun Ajaran',
    bill_types: 'Jenis Tagihan',
    billing_rules: 'Aturan Tagihan',
    student_bills: 'Tagihan Siswa',
    payments: 'Pembayaran',
    finance_reports: 'Laporan Keuangan',
    notifications: 'Notifikasi',
    auth: 'Autentikasi'
  }
  return dict[entity] || entity
}

const getHumanizedChanges = (log) => {
  if (!log) return []
  const oldVals = log.old_values || {}
  const newVals = log.new_values || {}

  const fieldDict = {
    name: 'Nama Lengkap',
    email: 'Alamat Email',
    phone_number: 'Nomor WhatsApp',
    role: 'Peran / Jabatan',
    nik: 'NIK / Nomor Identitas',
    nis: 'NIS',
    nisn: 'NISN',
    status: 'Status Operasional',
    is_active: 'Status Aktif',
    address: 'Alamat',
    parent_id: 'Wali Murid',
    major_id: 'Jurusan',
    class_id: 'Kelas',
    academic_year_id: 'Tahun Ajaran',
    amount: 'Total Dialokasikan (Rp)',
    deposit_applied: 'Saldo Deposit Dipakai (Rp)',
    cash_or_gateway_amount: 'Nominal Tunai/Gateway (Rp)',
    gateway_amount: 'Nominal Gateway (Rp)',
    deposit_refund: 'Saldo Deposit Masuk (Rp)',
    bill_type_name: 'Jenis Tagihan',
    due_date: 'Tanggal Jatuh Tempo',
    title: 'Judul Pesan',
    message: 'Isi Pesan',
    channels: 'Jalur Pengiriman',
    method: 'Metode Pembayaran',
    total_paid: 'Total Dibayar (Rp)',
    is_sandbox: 'Mode Sandbox',
    client_key: 'Client Key API',
    server_key: 'Server Key API',
    frontend_url: 'URL Frontend',
    waha_url: 'URL WAHA Bot'
  }

  const formatVal = (val, key) => {
    if (val === null || val === undefined || val === '') return '(Kosong)'
    if (typeof val === 'boolean') return val ? 'Ya / Aktif' : 'Tidak / Non-Aktif'
    if (Array.isArray(val)) return val.join(', ')
    if (typeof val === 'object') return JSON.stringify(val)
    if (['amount', 'total_paid', 'deposit_applied', 'cash_or_gateway_amount', 'gateway_amount', 'deposit_refund'].includes(key)) return `Rp ${Number(val).toLocaleString('id-ID')}`
    return String(val)
  }

  const changes = []

  if (log.action === 'CREATE' || log.action === 'LOGIN' || log.action === 'LOGIN_FAILED' || log.action === 'REFRESH_TOKEN' || log.action.includes('PASSWORD') || log.action === 'SEND_NOTIFICATION' || ['GENERATE_BILL', 'GENERATE_ADJUSTMENT_BILL', 'REFUND_BILL_REDUCTION', 'BULK_CANCEL_GENERATED_BILL', 'VOID_BILL', 'EXPORT_GLOBAL_FINANCE_REPORT', 'CREATE_PAYMENT_INTENT', 'PROCESS_PAYMENT', 'GATEWAY_PAYMENT_TO_DEPOSIT'].includes(log.action)) {
    for (const [key, val] of Object.entries(newVals)) {
      changes.push({
        field: fieldDict[key] || key.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase()),
        type: 'added',
        oldVal: '-',
        newVal: formatVal(val, key)
      })
    }
  } else if (log.action === 'DELETE' || log.action === 'BULK_DELETE') {
    for (const [key, val] of Object.entries(oldVals)) {
      changes.push({
        field: fieldDict[key] || key.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase()),
        type: 'deleted',
        oldVal: formatVal(val, key),
        newVal: '-'
      })
    }
  } else {
    const allKeys = new Set([...Object.keys(oldVals), ...Object.keys(newVals)])
    for (const key of allKeys) {
      const oldVal = oldVals[key]
      const newVal = newVals[key]
      if (JSON.stringify(oldVal) !== JSON.stringify(newVal)) {
        changes.push({
          field: fieldDict[key] || key.replace(/_/g, ' ').replace(/\b\w/g, l => l.toUpperCase()),
          type: 'modified',
          oldVal: formatVal(oldVal, key),
          newVal: formatVal(newVal, key)
        })
      }
    }
  }

  return changes
}

const getWhyContext = (log) => {
  if (!log) return 'Konteks tidak tersedia.'
  const act = log.action || ''
  const ent = log.entity_type || ''

  if (act === 'REFRESH_TOKEN') {
    return 'Memperpanjang masa aktif sesi login (JWT Refresh) secara otomatis agar pengguna dapat terus mengakses sistem tanpa harus memasukkan ulang email dan kata sandi.'
  }
  if (act === 'LOGIN') {
    return 'Memverifikasi kredensial identitas pengguna untuk memberikan hak akses masuk dan menerbitkan token otorisasi resmi ke dalam sistem SchoolPay.'
  }
  if (act === 'LOGIN_FAILED') {
    return 'Mencatat percobaan masuk yang gagal agar pola kesalahan kredensial, akun nonaktif, atau percobaan akses mencurigakan dapat diaudit.'
  }
  if (act === 'LOGOUT') {
    return 'Mengakhiri sesi akses aktif pengguna dan membersihkan token otorisasi dari peramban (browser) demi keamanan akun.'
  }
  if (act === 'FORGOT_PASSWORD') {
    return 'Mengirimkan tautan (link) darurat ke email pengguna untuk keperluan penyetelan ulang kata sandi yang lupa atau hilang.'
  }
  if (act === 'FORGOT_PASSWORD_FAILED') {
    return 'Mencatat permintaan reset kata sandi yang tidak dapat diproses tanpa membocorkan apakah email tersebut terdaftar kepada pengguna publik.'
  }
  if (act === 'RESET_PASSWORD') {
    return 'Menerapkan kata sandi baru yang diinput oleh pengguna melalui tautan pemulihan resmi untuk memulihkan akses akun.'
  }
  if (act === 'RESET_PASSWORD_FAILED') {
    return 'Mencatat kegagalan reset kata sandi, seperti token tidak valid, kedaluwarsa, atau kegagalan pemrosesan password baru.'
  }
  if (act === 'CHANGE_PASSWORD') {
    return 'Memperbarui kata sandi akun secara mandiri dari dalam pengaturan sistem untuk menjaga kerahasiaan dan keamanan kredensial.'
  }
  if (act === 'CHANGE_PASSWORD_FAILED') {
    return 'Mencatat kegagalan penggantian kata sandi dari profil, misalnya karena password saat ini tidak cocok atau pembaruan database gagal.'
  }
  if (act === 'CREATE') {
    return `Mendaftarkan atau menambahkan entitas rekaman data baru ke dalam database modul ${humanizeEntity(ent)} untuk memperluas cakupan operasional sistem.`
  }
  if (act === 'UPDATE') {
    return `Melakukan modifikasi atau penyesuaian atribut informasi pada rekaman data ${humanizeEntity(ent)} #${log.entity_id} guna memastikan keakuratan dan validitas data terkini.`
  }
  if (act === 'DELETE') {
    return `Menonaktifkan atau menghapus rekaman data ${humanizeEntity(ent)} #${log.entity_id} dari peredaran sistem (Soft Delete) untuk menjaga kebersihan basis data.`
  }
  if (act === 'RESTORE') {
    return `Memulihkan kembali rekaman data ${humanizeEntity(ent)} #${log.entity_id} yang sebelumnya telah dihapus agar dapat digunakan kembali dalam operasional sekolah.`
  }
  if (act === 'BULK_DELETE') {
    return `Menghapus sejumlah rekaman data ${humanizeEntity(ent)} secara massal sekaligus guna efisiensi pengelolaan basis data.`
  }
  if (act === 'BULK_RESTORE') {
    return `Memulihkan sejumlah rekaman data ${humanizeEntity(ent)} secara massal sekaligus untuk mengembalikan status operasional kolektif.`
  }
  if (act === 'TOGGLE_STATUS' || act === 'ACTIVATE_ACCOUNT') {
    return `Mengubah status operasional atau aktivasi pada entitas ${humanizeEntity(ent)} #${log.entity_id} untuk mengontrol hak akses atau validitas transaksi.`
  }
  if (act === 'SEND_NOTIFICATION') {
    return `Mengirimkan pesan pemberitahuan atau tagihan resmi secara otomatis melalui saluran komunikasi (WhatsApp/Email) kepada target penerima.`
  }
  if (act === 'GENERATE_BILL') {
    return `Menerbitkan kewajiban tagihan keuangan baru untuk siswa berdasarkan aturan pembiayaan sekolah yang berlaku.`
  }
  if (act === 'CREATE_PAYMENT_INTENT' || act === 'PROCESS_PAYMENT') {
    return `Memproses transaksi pembayaran atau pembuatan kode bayar tagihan keuangan siswa untuk pencatatan kas sekolah.`
  }
  if (act === 'EXPORT_TREND_EXCEL') {
    return `Mengunduh laporan rekapitulasi data ${humanizeEntity(ent)} ke dalam format spreadsheet Excel untuk kebutuhan analisis atau arsip eksternal.`
  }

  return `Melakukan instruksi eksekusi ${act} pada entitas modul ${humanizeEntity(ent)} sesuai dengan prosedur sistem SchoolPay.`
}

const getHowSummary = (log) => {
  if (!log) return '-'
  const changes = getHumanizedChanges(log)
  if (changes.length === 0) {
    if (log.action === 'REFRESH_TOKEN') {
      const email = log.new_values?.email || 'pengguna'
      return `Sistem memperpanjang masa berlaku sesi untuk akun dengan email '${email}' melalui validasi Refresh Token yang sah tanpa mengubah struktur tabel database utama.`
    }
    if (log.action === 'LOGIN') {
      const email = log.new_values?.email || 'pengguna'
      return `Sistem memvalidasi email '${email}' dan mencatat alamat IP serta User Agent perangkat yang digunakan untuk masuk ke dalam aplikasi.`
    }
    return `Sistem memproses instruksi tanpa melakukan mutasi langsung pada kolom atribut tabel spesifik.`
  }

  const summaries = changes.map(c => {
    if (c.type === 'added') return `Atribut ${c.field} dicatat dengan nilai '${c.newVal}'`
    if (c.type === 'deleted') return `Atribut ${c.field} dengan nilai lama '${c.oldVal}' dihapus/ diarsipkan`
    return `Atribut ${c.field} diubah dari '${c.oldVal}' menjadi '${c.newVal}'`
  })

  return `Terdapat ${changes.length} poin perubahan spesifik: ` + summaries.join('; ') + '.'
}
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
