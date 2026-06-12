import axios from 'axios';
import MockAdapter from 'axios-mock-adapter';

const mock = new MockAdapter(axios, { delayResponse: 400 });

if (typeof window !== 'undefined') {
  window.__SCHOOLPAY_MOCK__ = true;
}

// ============================================================
// 1. IN-MEMORY DATABASE (DATA SANGAT KOMPLEKS)
// ============================================================
const db = {
  users: [
    { id: 1, name: 'Admin Utama', email: 'admin@demo.com', role: 'admin', is_active: true, phone_number: '6281200001111', student_count: 0, created_at: '2025-01-15T08:00:00Z', deleted_at: null },
    { id: 2, name: 'Bapak Budi', email: 'parent@demo.com', role: 'parent', is_active: true, phone_number: '6281234567890', student_count: 1, created_at: '2025-02-10T09:00:00Z', deleted_at: null },
    { id: 3, name: 'Ibu Siti Aminah', email: 'siti@demo.com', role: 'parent', is_active: true, phone_number: '6281234567891', student_count: 1, created_at: '2025-02-12T10:00:00Z', deleted_at: null },
    { id: 4, name: 'Bapak Joko', email: 'joko@demo.com', role: 'parent', is_active: true, phone_number: '6281234567892', student_count: 2, created_at: '2025-03-01T11:00:00Z', deleted_at: null },
    { id: 5, name: 'Ibu Ratna', email: 'ratna@demo.com', role: 'parent', is_active: false, phone_number: '6281234567893', student_count: 1, created_at: '2025-03-05T14:00:00Z', deleted_at: null },
    { id: 6, name: 'Bapak Wahyu (Terhapus)', email: 'wahyu@demo.com', role: 'parent', is_active: false, phone_number: '6281234567800', student_count: 0, created_at: '2025-01-20T08:00:00Z', deleted_at: '2025-06-01T10:00:00Z' }
  ],
  academicYears: [
    { id: 1, year: '2024/2025', status: 'active', is_active: true, deleted_at: null },
    { id: 2, year: '2025/2026', status: 'planned', is_active: true, deleted_at: null }
  ],
  majors: [
    { id: 1, name: 'MIPA', code: 'MIPA', is_active: true, deleted_at: null },
    { id: 2, name: 'IPS', code: 'IPS', is_active: true, deleted_at: null },
    { id: 3, name: 'Bahasa', code: 'BHS', is_active: false, deleted_at: null }
  ],
  classes: [
    { id: 1, name: 'X MIPA 1', major_id: 1, major_name: 'MIPA', academic_year_id: 1, is_active: true, deleted_at: null },
    { id: 2, name: 'X IPS 1', major_id: 2, major_name: 'IPS', academic_year_id: 1, is_active: true, deleted_at: null },
    { id: 3, name: 'XI MIPA 1', major_id: 1, major_name: 'MIPA', academic_year_id: 1, is_active: true, deleted_at: null },
    { id: 4, name: 'XI IPS 1', major_id: 2, major_name: 'IPS', academic_year_id: 1, is_active: true, deleted_at: null }
  ],
  students: [
    { id: 1, name: 'Budi Santoso', nisn: '0012345678', nis: '10001', nik: '3201010101010001', gender: 'Laki-laki', class_id: 1, class_name: 'X MIPA 1', major_id: 1, major_name: 'MIPA', parent_id: 2, parent_name: 'Bapak Budi', parent_email: 'parent@demo.com', is_active: true, status: 'active', deposit_balance: 500000, phone_number: '081234567890', entry_year: 2024, birth_place: 'Jakarta', birth_date: '2008-05-15', religion: 'Islam', address: 'Jl. Merdeka No. 10, Jakarta Selatan', created_at: '2025-07-01T08:00:00Z', deleted_at: null },
    { id: 2, name: 'Andi Saputra', nisn: '0012345679', nis: '10002', nik: '3201010101010002', gender: 'Laki-laki', class_id: 1, class_name: 'X MIPA 1', major_id: 1, major_name: 'MIPA', parent_id: 3, parent_name: 'Ibu Siti Aminah', parent_email: 'siti@demo.com', is_active: true, status: 'active', deposit_balance: 0, phone_number: '081234567891', entry_year: 2024, birth_place: 'Bandung', birth_date: '2008-08-20', religion: 'Islam', address: 'Jl. Asia Afrika No. 5, Bandung', created_at: '2025-07-01T08:30:00Z', deleted_at: null },
    { id: 3, name: 'Citra Dewi', nisn: '0012345680', nis: '10003', nik: '3201010101010003', gender: 'Perempuan', class_id: 2, class_name: 'X IPS 1', major_id: 2, major_name: 'IPS', parent_id: 4, parent_name: 'Bapak Joko', parent_email: 'joko@demo.com', is_active: true, status: 'active', deposit_balance: 250000, phone_number: '081234567892', entry_year: 2024, birth_place: 'Surabaya', birth_date: '2008-03-12', religion: 'Kristen', address: 'Jl. Diponegoro No. 22, Surabaya', created_at: '2025-07-02T09:00:00Z', deleted_at: null },
    { id: 4, name: 'Dian Permata', nisn: '0012345681', nis: '10004', nik: '3201010101010004', gender: 'Perempuan', class_id: 2, class_name: 'X IPS 1', major_id: 2, major_name: 'IPS', parent_id: 4, parent_name: 'Bapak Joko', parent_email: 'joko@demo.com', is_active: true, status: 'active', deposit_balance: 100000, phone_number: '081234567893', entry_year: 2024, birth_place: 'Semarang', birth_date: '2009-01-25', religion: 'Islam', address: 'Jl. Pemuda No. 8, Semarang', created_at: '2025-07-02T09:30:00Z', deleted_at: null },
    { id: 5, name: 'Eka Pratama', nisn: '0012345682', nis: '10005', nik: '3201010101010005', gender: 'Laki-laki', class_id: 3, class_name: 'XI MIPA 1', major_id: 1, major_name: 'MIPA', parent_id: 5, parent_name: 'Ibu Ratna', parent_email: 'ratna@demo.com', is_active: true, status: 'active', deposit_balance: 0, phone_number: '081234567894', entry_year: 2023, birth_place: 'Yogyakarta', birth_date: '2007-11-30', religion: 'Hindu', address: 'Jl. Malioboro No. 1, Yogyakarta', created_at: '2025-07-03T10:00:00Z', deleted_at: null },
    { id: 6, name: 'Farhan Rizky (Terhapus)', nisn: '0012345699', nis: '10099', nik: '3201010101010099', gender: 'Laki-laki', class_id: 1, class_name: 'X MIPA 1', major_id: 1, major_name: 'MIPA', parent_id: 6, parent_name: 'Bapak Wahyu', parent_email: 'wahyu@demo.com', is_active: false, status: 'inactive', deposit_balance: 0, phone_number: '081234567800', entry_year: 2024, birth_place: 'Medan', birth_date: '2008-07-07', religion: 'Islam', address: 'Jl. Gatot Subroto No. 3, Medan', created_at: '2025-07-01T08:00:00Z', deleted_at: '2025-10-15T10:00:00Z' }
  ],
  billTypes: [
    { id: 1, name: 'SPP Bulanan', type: 'MONTHLY', default_amount: 350000, is_active: true, deleted_at: null },
    { id: 2, name: 'Uang Gedung', type: 'ONE_TIME', default_amount: 2000000, is_active: true, deleted_at: null },
    { id: 3, name: 'Daftar Ulang', type: 'YEARLY', default_amount: 500000, is_active: true, deleted_at: null },
    { id: 4, name: 'Seragam (Terhapus)', type: 'ONE_TIME', default_amount: 750000, is_active: false, deleted_at: '2025-09-01T10:00:00Z' }
  ],
  billingRules: [
    { id: 1, name: 'SPP Kelas X MIPA', bill_type_id: 1, bill_type_name: 'SPP Bulanan', target_type: 'class', target_id: 1, target_name: 'X MIPA 1', class_id: null, amount: 350000, period_type: 'bulanan', due_day: 10, is_active: true, bill_count: 2, start_date: '2025-07-01T00:00:00Z', end_date: '2026-06-30T23:59:59Z', deleted_at: null },
    { id: 2, name: 'SPP Kelas X IPS', bill_type_id: 1, bill_type_name: 'SPP Bulanan', target_type: 'class', target_id: 2, target_name: 'X IPS 1', class_id: null, amount: 350000, period_type: 'bulanan', due_day: 10, is_active: true, bill_count: 2, start_date: '2025-07-01T00:00:00Z', end_date: '2026-06-30T23:59:59Z', deleted_at: null },
    { id: 3, name: 'Uang Gedung Semua Siswa', bill_type_id: 2, bill_type_name: 'Uang Gedung', target_type: 'all', target_id: 0, target_name: 'Semua Siswa', class_id: null, amount: 2000000, period_type: 'sekali', due_day: 15, is_active: true, bill_count: 5, start_date: '2025-07-01T00:00:00Z', end_date: null, deleted_at: null }
  ],
  bills: [
    { id: 1, student_id: 1, student_name: 'Budi Santoso', bill_type_name: 'SPP Bulanan', name: 'SPP Juni 2026', amount: 350000, total_paid: 0, remaining_amount: 350000, status: 'unpaid', due_date: new Date(Date.now() + 5*86400000).toISOString(), period: '2026-06', academic_year: '2024/2025', deleted_at: null },
    { id: 2, student_id: 1, student_name: 'Budi Santoso', bill_type_name: 'Uang Gedung', name: 'Uang Gedung', amount: 2000000, total_paid: 500000, remaining_amount: 1500000, status: 'partial', due_date: new Date(Date.now() + 30*86400000).toISOString(), period: '2025-07', academic_year: '2024/2025', deleted_at: null },
    { id: 3, student_id: 1, student_name: 'Budi Santoso', bill_type_name: 'Daftar Ulang', name: 'Daftar Ulang', amount: 500000, total_paid: 500000, remaining_amount: 0, status: 'paid', due_date: new Date(Date.now() - 60*86400000).toISOString(), period: '2025-07', academic_year: '2024/2025', deleted_at: null },
    { id: 4, student_id: 2, student_name: 'Andi Saputra', bill_type_name: 'SPP Bulanan', name: 'SPP Juni 2026', amount: 350000, total_paid: 0, remaining_amount: 350000, status: 'unpaid', due_date: new Date(Date.now() + 5*86400000).toISOString(), period: '2026-06', academic_year: '2024/2025', deleted_at: null },
    { id: 5, student_id: 2, student_name: 'Andi Saputra', bill_type_name: 'Uang Gedung', name: 'Uang Gedung', amount: 2000000, total_paid: 2000000, remaining_amount: 0, status: 'paid', due_date: new Date(Date.now() - 30*86400000).toISOString(), period: '2025-07', academic_year: '2024/2025', deleted_at: null },
    { id: 6, student_id: 3, student_name: 'Citra Dewi', bill_type_name: 'SPP Bulanan', name: 'SPP Juni 2026', amount: 350000, total_paid: 0, remaining_amount: 350000, status: 'unpaid', due_date: new Date(Date.now() + 5*86400000).toISOString(), period: '2026-06', academic_year: '2024/2025', deleted_at: null },
    { id: 7, student_id: 3, student_name: 'Citra Dewi', bill_type_name: 'SPP Bulanan', name: 'SPP Mei 2026', amount: 350000, total_paid: 0, remaining_amount: 350000, status: 'overdue', due_date: new Date(Date.now() - 15*86400000).toISOString(), period: '2026-05', academic_year: '2024/2025', deleted_at: null },
    { id: 8, student_id: 4, student_name: 'Dian Permata', bill_type_name: 'SPP Bulanan', name: 'SPP Juni 2026', amount: 350000, total_paid: 350000, remaining_amount: 0, status: 'paid', due_date: new Date(Date.now() + 5*86400000).toISOString(), period: '2026-06', academic_year: '2024/2025', deleted_at: null },
    { id: 9, student_id: 5, student_name: 'Eka Pratama', bill_type_name: 'SPP Bulanan', name: 'SPP Juni 2026', amount: 350000, total_paid: 0, remaining_amount: 350000, status: 'unpaid', due_date: new Date(Date.now() + 5*86400000).toISOString(), period: '2026-06', academic_year: '2024/2025', deleted_at: null }
  ],
  payments: [
    { id: 1, student_id: 1, student_name: 'Budi Santoso', method: 'Tunai', status: 'success', amount: 500000, bill_type_names: 'Daftar Ulang', created_at: new Date(Date.now() - 60*86400000).toISOString() },
    { id: 2, student_id: 1, student_name: 'Budi Santoso', method: 'Midtrans', status: 'success', amount: 500000, bill_type_names: 'Uang Gedung (Cicilan 1)', created_at: new Date(Date.now() - 30*86400000).toISOString() },
    { id: 3, student_id: 2, student_name: 'Andi Saputra', method: 'Tunai', status: 'success', amount: 2000000, bill_type_names: 'Uang Gedung', created_at: new Date(Date.now() - 45*86400000).toISOString() },
    { id: 4, student_id: 4, student_name: 'Dian Permata', method: 'Transfer Bank', status: 'success', amount: 350000, bill_type_names: 'SPP Juni', created_at: new Date(Date.now() - 2*86400000).toISOString() }
  ],
  notifications: [
    { id: 1, type: 'whatsapp', recipient_name: 'Bapak Budi', recipient_phone: '081234567890', title: 'Tagihan SPP Baru', message: 'Yth. Bapak Budi, tagihan SPP bulan Juni sebesar Rp350.000 telah diterbitkan. Jatuh tempo: 10/06/2026', delivery_status: 'delivered', created_at: new Date(Date.now() - 3*86400000).toISOString() },
    { id: 2, type: 'whatsapp', recipient_name: 'Ibu Siti Aminah', recipient_phone: '081234567891', title: 'Tagihan SPP Baru', message: 'Yth. Ibu Siti, tagihan SPP bulan Juni sebesar Rp350.000 telah diterbitkan.', delivery_status: 'delivered', created_at: new Date(Date.now() - 3*86400000).toISOString() },
    { id: 3, type: 'whatsapp', recipient_name: 'Bapak Joko', recipient_phone: '081234567892', title: 'Peringatan Tunggakan', message: 'Yth. Bapak Joko, Citra Dewi memiliki tunggakan SPP Mei 2026 sebesar Rp350.000.', delivery_status: 'failed', created_at: new Date(Date.now() - 1*86400000).toISOString() }
  ],
  emailNotifications: [
    { id: 101, type: 'email', recipient_name: 'Bapak Budi', recipient_email: 'parent@demo.com', title: 'Bukti Pembayaran SchoolPay', message: 'Pembayaran Daftar Ulang telah berhasil dikonfirmasi.', delivery_status: 'delivered', created_at: new Date(Date.now() - 2*86400000).toISOString() },
    { id: 102, type: 'email', recipient_name: 'Ibu Siti Aminah', recipient_email: 'siti@demo.com', title: 'Tagihan Baru SchoolPay', message: 'Tagihan Uang Gedung telah diterbitkan untuk Andi Saputra.', delivery_status: 'sent', created_at: new Date(Date.now() - 4*86400000).toISOString() },
    { id: 103, type: 'email', recipient_name: 'Bapak Joko', recipient_email: 'joko@demo.com', title: 'Pengingat Tunggakan', message: 'Mohon segera menyelesaikan tunggakan SPP Mei 2026.', delivery_status: 'failed', delivery_error: 'Mailbox sementara tidak tersedia', created_at: new Date(Date.now() - 1*86400000).toISOString() }
  ],
  supportConversations: [
    { id: 1, parent_name: 'Bapak Budi', student_names: 'Budi Santoso', phone_number: '6281234567890', status: 'open', last_message: 'Saya ingin bertanya soal cicilan uang gedung.', created_at: new Date(Date.now() - 45*60000).toISOString() },
    { id: 2, parent_name: 'Bapak Joko', student_names: 'Citra Dewi, Dian Permata', phone_number: '6281234567892', status: 'pending', last_message: 'Bukti transfer sudah saya kirim lewat WhatsApp.', created_at: new Date(Date.now() - 4*3600000).toISOString() },
    { id: 3, parent_name: 'Ibu Siti Aminah', student_names: 'Andi Saputra', phone_number: '6281234567891', status: 'closed', last_message: 'Terima kasih, sudah jelas.', created_at: new Date(Date.now() - 26*3600000).toISOString() }
  ],
  auditLogs: [
    { id: 1, user_name: 'Admin Utama', action: 'CREATE', entity: 'STUDENT', description: 'Menambahkan siswa Budi Santoso', created_at: new Date(Date.now() - 90*86400000).toISOString() },
    { id: 2, user_name: 'Admin Utama', action: 'GENERATE', entity: 'BILLING_RULE', description: 'Generate tagihan SPP untuk 5 siswa', created_at: new Date(Date.now() - 5*86400000).toISOString() },
    { id: 3, user_name: 'Admin Utama', action: 'PAYMENT', entity: 'BILL', description: 'Pembayaran manual Daftar Ulang Budi Santoso (Rp500.000)', created_at: new Date(Date.now() - 60*86400000).toISOString() },
    { id: 4, user_name: 'Admin Utama', action: 'DELETE', entity: 'STUDENT', description: 'Menghapus siswa Farhan Rizky ke Trash', created_at: new Date(Date.now() - 30*86400000).toISOString() }
  ]
};

// ============================================================
// 2. HELPER FUNCTIONS
// ============================================================
let _nextId = 1000;
const nextId = () => ++_nextId;

const getActive = (collection) => collection.filter(i => !i.deleted_at);
const getTrashed = (collection) => collection.filter(i => !!i.deleted_at);

const parseData = (config) => {
  try {
    if (config.data instanceof FormData) {
      const obj = {};
      config.data.forEach((v, k) => { obj[k] = v; });
      return obj;
    }
    return JSON.parse(config.data || '{}');
  } catch { return {}; }
};

const extractId = (url, segment) => {
  const parts = url.replace(/\?.*/, '').split('/');
  if (segment !== undefined) return parseInt(parts[segment]);
  // Find last numeric segment
  for (let i = parts.length - 1; i >= 0; i--) {
    if (/^\d+$/.test(parts[i])) return parseInt(parts[i]);
  }
  return null;
};

const addAuditLog = (action, entity, description) => {
  db.auditLogs.unshift({ id: nextId(), user_name: 'Admin Utama', action, entity, description, created_at: new Date().toISOString() });
};

const addNotification = (recipientName, recipientPhone, title, message) => {
  db.notifications.unshift({ id: nextId(), type: 'whatsapp', recipient_name: recipientName, recipient_phone: recipientPhone, title, message, delivery_status: 'delivered', created_at: new Date().toISOString() });
};

const paramsFrom = (config) => {
  const query = config.url?.includes('?') ? config.url.split('?')[1] : '';
  return { ...Object.fromEntries(new URLSearchParams(query)), ...(config.params || {}) };
};

const paginate = (items, params = {}) => {
  const page = parseInt(params.page) || 1;
  const limit = parseInt(params.limit) || 10;
  const total = items.length;
  const totalPages = Math.ceil(total / limit) || 1;
  const start = (page - 1) * limit;
  return { data: items.slice(start, start + limit), total, page, totalPages };
};

const paymentTrend = () => {
  return Array.from({ length: 6 }, (_, index) => {
    const date = new Date();
    date.setMonth(date.getMonth() - (5 - index));
    return {
      date: date.toLocaleDateString('id-ID', { month: 'short', year: 'numeric' }),
      total: Math.max(0, db.payments.reduce((sum, p) => sum + p.amount, 0) / 6 + index * 125000)
    };
  });
};

// ============================================================
// 3. GENERIC CRUD FACTORY (WITH SOFT DELETE)
// ============================================================
function createCRUD(endpoint, collectionName) {
  const col = () => db[collectionName];

  // GET ALL (supports status=trash)
  mock.onGet(new RegExp(`\\/?${endpoint}(\\?.*)?$`)).reply(config => {
    const params = config.params || {};
    const isTrash = params.status === 'trash';
    let data = isTrash ? getTrashed(col()) : getActive(col());
    // Search
    if (params.search) {
      const q = params.search.toLowerCase();
      data = data.filter(i => JSON.stringify(i).toLowerCase().includes(q));
    }
    // For users, support nested users key
    if (collectionName === 'users') {
      return [200, { status: true, data: { users: data, data: data, total: data.length, page: 1, totalPages: 1 } }];
    }
    return [200, { status: true, data: { data, total: data.length, page: 1, totalPages: 1 } }];
  });

  // GET SINGLE
  mock.onGet(new RegExp(`\\/?${endpoint}/\\d+$`)).reply(config => {
    const id = extractId(config.url);
    const item = col().find(i => i.id === id);
    if (item) return [200, { status: true, data: item }];
    return [404, { status: false, message: 'Tidak ditemukan' }];
  });

  // GET DEPENDENCY INFO
  mock.onGet(new RegExp(`\\/?${endpoint}/\\d+/dependency-info`)).reply(200, { status: true, data: { has_dependencies: false, message: '' } });

  // CHECK UNIQUE
  mock.onGet(new RegExp(`\\/?${endpoint}/check-unique`)).reply(200, { status: true, data: { is_unique: true } });

  // CREATE
  mock.onPost(new RegExp(`\\/?${endpoint}$`)).reply(config => {
    const data = parseData(config);
    const newItem = { id: nextId(), ...data, is_active: true, created_at: new Date().toISOString(), deleted_at: null };
    col().unshift(newItem);
    addAuditLog('CREATE', collectionName.toUpperCase(), `Menambahkan data baru: ${data.name || data.year || 'item'}`);
    return [200, { status: true, message: 'Data berhasil ditambahkan', data: newItem }];
  });

  // UPDATE
  mock.onPut(new RegExp(`\\/?${endpoint}/\\d+$`)).reply(config => {
    const id = extractId(config.url);
    const data = parseData(config);
    const idx = col().findIndex(i => i.id === id);
    if (idx !== -1) {
      col()[idx] = { ...col()[idx], ...data, updated_at: new Date().toISOString() };
      addAuditLog('UPDATE', collectionName.toUpperCase(), `Memperbarui data: ${col()[idx].name || col()[idx].year || id}`);
      return [200, { status: true, message: 'Data berhasil diperbarui', data: col()[idx] }];
    }
    return [404, { status: false, message: 'Tidak ditemukan' }];
  });

  // SOFT DELETE
  mock.onDelete(new RegExp(`\\/?${endpoint}/\\d+$`)).reply(config => {
    const id = extractId(config.url);
    const item = col().find(i => i.id === id);
    if (item) {
      item.deleted_at = new Date().toISOString();
      addAuditLog('DELETE', collectionName.toUpperCase(), `Menghapus: ${item.name || item.year || id}`);
    }
    return [200, { status: true, message: 'Data berhasil dihapus' }];
  });

  // BULK DELETE
  mock.onPost(new RegExp(`\\/?${endpoint}/bulk-delete`)).reply(config => {
    const { ids } = parseData(config);
    if (ids) ids.forEach(id => { const item = col().find(i => i.id === id); if (item) item.deleted_at = new Date().toISOString(); });
    addAuditLog('BULK_DELETE', collectionName.toUpperCase(), `Menghapus ${ids?.length || 0} data`);
    return [200, { status: true, message: 'Data terpilih berhasil dihapus' }];
  });

  // TOGGLE STATUS
  mock.onPatch(new RegExp(`\\/?${endpoint}/\\d+/status`)).reply(config => {
    const id = extractId(config.url);
    const item = col().find(i => i.id === id);
    if (item) { item.is_active = !item.is_active; return [200, { status: true, message: 'Status diubah', data: item }]; }
    return [404, { status: false, message: 'Tidak ditemukan' }];
  });

  // RESTORE
  mock.onPatch(new RegExp(`\\/?${endpoint}/\\d+/restore`)).reply(config => {
    const id = extractId(config.url);
    const item = col().find(i => i.id === id);
    if (item) { item.deleted_at = null; item.is_active = true; addAuditLog('RESTORE', collectionName.toUpperCase(), `Memulihkan: ${item.name || id}`); }
    return [200, { status: true, message: 'Berhasil dipulihkan' }];
  });

  // BULK RESTORE
  mock.onPatch(new RegExp(`\\/?${endpoint}/bulk-restore`)).reply(config => {
    const { ids } = parseData(config);
    if (ids) ids.forEach(id => { const item = col().find(i => i.id === id); if (item) { item.deleted_at = null; item.is_active = true; } });
    addAuditLog('BULK_RESTORE', collectionName.toUpperCase(), `Memulihkan ${ids?.length || 0} data`);
    return [200, { status: true, message: 'Data terpilih berhasil dipulihkan' }];
  });

  // EXPORT (return blob)
  mock.onGet(new RegExp(`\\/?${endpoint}/export`)).reply(200, new Blob(['Demo Export'], { type: 'application/octet-stream' }));
}

// ============================================================
// 4. REGISTER GENERIC CRUD FOR ALL ENTITIES
// ============================================================
createCRUD('users', 'users');
createCRUD('students', 'students');
createCRUD('academic/years', 'academicYears');
createCRUD('academic/major', 'majors');
createCRUD('academic/class', 'classes');
createCRUD('finance/bill-types', 'billTypes');
createCRUD('finance/billing-rules', 'billingRules');
createCRUD('finance/bills', 'bills');

// ============================================================
// 5. CUSTOM / SPECIALIZED ENDPOINTS
// ============================================================

// --- AUTH ---
mock.onPost('/auth/login').reply(config => {
  const { email, password } = parseData(config);
  if (email === 'admin@demo.com' && password === 'admin123') return [200, { status: true, message: 'Login berhasil', data: { access_token: 'token_admin', user: { ...db.users[0], role: 'admin' } } }];
  if (email === 'parent@demo.com' && password === 'parent123') return [200, { status: true, message: 'Login berhasil', data: { access_token: 'token_parent', user: { ...db.users[1], role: 'parent' } } }];
  return [401, { status: false, message: 'Email atau Password Salah' }];
});
mock.onPost('/auth/refresh').reply(config => {
  const auth = config.headers?.['Authorization'] || '';
  if (auth.includes('parent')) return [200, { status: true, data: { access_token: 'token_parent', user: { ...db.users[1], role: 'parent' } } }];
  return [200, { status: true, data: { access_token: 'token_admin', user: { ...db.users[0], role: 'admin' } } }];
});
mock.onPost('/auth/logout').reply(200, { status: true });
mock.onGet(/\/?auth\/me/).reply(200, { status: true, data: db.users[0] });
mock.onPut(/\/?auth\/profile/).reply(200, { status: true, message: 'Profil berhasil diperbarui' });
mock.onPost(/\/?auth\/profile\/photo/).reply(200, { status: true, message: 'Foto profil berhasil diperbarui' });
mock.onPost(/\/?auth\/change-password/).reply(200, { status: true, message: 'Password berhasil diubah' });
mock.onPost(/\/?auth\/forgot-password/).reply(200, { status: true, message: 'Email reset password terkirim' });
mock.onPost(/\/?auth\/reset-password/).reply(200, { status: true, message: 'Password berhasil direset' });
mock.onPost(/\/?users\/activate/).reply(200, { status: true, message: 'Akun berhasil diaktifkan' });

// --- DASHBOARD ---
mock.onGet(/\/?dashboard\/stats/).reply(() => {
  const activeStudents = getActive(db.students);
  const activeBills = getActive(db.bills);
  const paidBills = activeBills.filter(b => b.status === 'paid');
  const unpaidBills = activeBills.filter(b => b.status !== 'paid');
  const overdueBills = activeBills.filter(b => b.status === 'overdue');
  const dueSoonBills = unpaidBills.filter(b => {
    const diff = new Date(b.due_date).getTime() - Date.now();
    return diff >= 0 && diff <= 7 * 86400000;
  });
  return [200, { status: true, data: {
    stats: {
      students: { total: activeStudents.length, total_all: db.students.length, growth: 8 },
      users: { total: getActive(db.users).length, new_this_period: 1, growth: 4 },
      bills: { total: activeBills.length, growth: 12 },
      payments: { total: db.payments.length, growth: 15 },
      unpaid_amount: unpaidBills.reduce((a, b) => a + b.remaining_amount, 0),
      paid_amount: db.payments.reduce((a, p) => a + p.amount, 0),
      paid_count: paidBills.length,
      unpaid_count: unpaidBills.length,
      total_bills: activeBills.length,
      payments_today: db.payments.filter(p => new Date(p.created_at).toDateString() === new Date().toDateString()).length,
      failed_reminders: db.notifications.filter(n => n.delivery_status === 'failed').length
    },
    critical_bills: {
      overdue: overdueBills.slice(0, 5),
      due_soon: dueSoonBills.slice(0, 5)
    },
    demographics: {
      gender: {
        'Laki-laki': activeStudents.filter(s => s.gender === 'Laki-laki').length,
        'Perempuan': activeStudents.filter(s => s.gender === 'Perempuan').length
      },
      major: Object.fromEntries(db.majors.map(m => [m.name, activeStudents.filter(s => s.major_id === m.id).length])),
      class: Object.fromEntries(db.classes.map(c => [c.name, activeStudents.filter(s => s.class_id === c.id).length])),
      whatsapp: {
        pending: db.notifications.filter(n => n.delivery_status === 'pending').length,
        sent: db.notifications.filter(n => n.delivery_status === 'sent').length,
        delivered: db.notifications.filter(n => n.delivery_status === 'delivered').length,
        read: db.notifications.filter(n => n.delivery_status === 'read').length,
        failed: db.notifications.filter(n => n.delivery_status === 'failed').length
      },
      email: {
        pending: db.emailNotifications.filter(n => n.delivery_status === 'pending').length,
        sent: db.emailNotifications.filter(n => n.delivery_status === 'sent').length,
        delivered: db.emailNotifications.filter(n => n.delivery_status === 'delivered').length,
        read: db.emailNotifications.filter(n => n.delivery_status === 'read').length,
        failed: db.emailNotifications.filter(n => n.delivery_status === 'failed').length
      }
    },
    payment_trend: paymentTrend(),
    entry_years: [...new Set(activeStudents.map(s => s.entry_year))].sort(),
    total_payments_count: db.payments.length,
    recent_payments: db.payments.slice(0, 5),
    recent_notifications: db.notifications.slice(0, 5)
  }}];
});

mock.onGet(/\/?dashboard\/communication-details/).reply(config => {
  const params = paramsFrom(config);
  const logs = params.channel === 'email' ? db.emailNotifications : db.notifications;
  return [200, { status: true, data: logs.filter(log => !params.status || log.delivery_status === params.status) }];
});

mock.onGet(/\/?dashboard\/export/).reply(200, new Blob(['Laporan demo SchoolPay'], { type: 'application/octet-stream' }));

// --- STUDENT RELATIONS ---
mock.onGet(/\/?students\/filters/).reply(() => {
  return [200, { status: true, data: {
    years: db.academicYears.map(y => ({ value: y.year, label: y.year })),
    majors: getActive(db.majors).map(m => ({ value: m.id, label: m.name })),
    classes: getActive(db.classes).map(c => ({ value: c.id, label: c.name }))
  }}];
});

mock.onGet(/\/?students\/\d+\/parents/).reply(config => {
  const studentId = extractId(config.url);
  const student = db.students.find(s => s.id === studentId);
  if (student) {
    const parent = db.users.find(u => u.id === student.parent_id);
    return [200, { status: true, data: parent ? [{ ...parent, relation: 'Wali Murid' }] : [] }];
  }
  return [200, { status: true, data: [] }];
});

mock.onGet(/\/?students\/\d+\/class-history/).reply(config => {
  const studentId = extractId(config.url);
  const student = db.students.find(s => s.id === studentId);
  if (student) {
    return [200, { status: true, data: [
      { id: 1, class_name: student.class_name, grade: 'X', academic_year: '2024/2025', is_active: true, created_at: student.created_at }
    ]}];
  }
  return [200, { status: true, data: [] }];
});

mock.onGet(/\/?users\/\d+\/students/).reply(config => {
  const parentId = extractId(config.url);
  const children = getActive(db.students).filter(s => s.parent_id === parentId);
  return [200, { status: true, data: children }];
});

mock.onGet(/\/?users\/parents/).reply(config => {
  const params = paramsFrom(config);
  let parents = getActive(db.users).filter(user => user.role === 'parent');
  if (params.search) {
    const q = String(params.search).toLowerCase();
    parents = parents.filter(user => `${user.name} ${user.email} ${user.phone_number}`.toLowerCase().includes(q));
  }
  return [200, { status: true, data: { ...paginate(parents, params), users: paginate(parents, params).data } }];
});

// Student Bulk Actions
mock.onPost(/\/?students\/bulk-promote/).reply(200, { status: true, message: 'Perpindahan kelas berhasil' });
mock.onPost(/\/?students\/bulk-graduate/).reply(200, { status: true, message: 'Kelulusan masal berhasil' });

// --- ACADEMIC YEAR RELATIONS ---
mock.onGet(/\/?academic\/years\/\d+\/majors/).reply(200, { status: true, data: getActive(db.majors) });
mock.onGet(/\/?academic\/years\/\d+\/classes/).reply(200, { status: true, data: getActive(db.classes) });
mock.onGet(/\/?academic\/class\/suggest-name/).reply(config => {
  return [200, { status: true, data: { suggested_name: (config.params?.name || 'Kelas') + ' 2' } }];
});

// --- FINANCE: GENERATE BILLS (COMPLEX LOGIC) ---
mock.onPost(/\/?finance\/generate-bills$/).reply(config => {
  const { rule_id } = parseData(config);
  const rule = db.billingRules.find(r => r.id === rule_id);
  if (!rule) return [404, { status: false, message: 'Aturan tidak ditemukan' }];

  // Find target students
  let targets = getActive(db.students);
  if (rule.target_type === 'class') targets = targets.filter(s => s.class_id === rule.target_id);
  else if (rule.target_type === 'major') targets = targets.filter(s => s.major_id === rule.target_id);

  let generated = 0;
  targets.forEach(student => {
    const newBill = {
      id: nextId(), student_id: student.id, student_name: student.name,
      bill_type_name: rule.bill_type_name || 'Tagihan', name: `${rule.bill_type_name || 'Tagihan'} - Generated`,
      amount: rule.amount, total_paid: 0, remaining_amount: rule.amount,
      status: 'unpaid', due_date: new Date(Date.now() + 30*86400000).toISOString(),
      period: new Date().toISOString().slice(0, 7), academic_year: '2024/2025', deleted_at: null
    };
    db.bills.unshift(newBill);
    generated++;
    // Send notification to parent
    const parent = db.users.find(u => u.id === student.parent_id);
    if (parent) addNotification(parent.name, parent.phone_number, 'Tagihan Baru Diterbitkan', `Yth. ${parent.name}, tagihan ${rule.bill_type_name} sebesar Rp${rule.amount.toLocaleString('id-ID')} telah diterbitkan untuk ${student.name}.`);
  });

  rule.bill_count = (rule.bill_count || 0) + generated;
  addAuditLog('GENERATE', 'BILLING_RULE', `Generate tagihan "${rule.name}" untuk ${generated} siswa`);
  return [200, { status: true, message: `Berhasil generate ${generated} tagihan!` }];
});

mock.onPost(/\/?finance\/generate-bills\/bulk/).reply(config => {
  const { rule_ids } = parseData(config);
  let totalGenerated = 0;
  (rule_ids || []).forEach(ruleId => {
    const rule = db.billingRules.find(r => r.id === ruleId);
    if (!rule) return;
    let targets = getActive(db.students);
    if (rule.target_type === 'class') targets = targets.filter(s => s.class_id === rule.target_id);
    else if (rule.target_type === 'major') targets = targets.filter(s => s.major_id === rule.target_id);
    targets.forEach(student => {
      db.bills.unshift({ id: nextId(), student_id: student.id, student_name: student.name, bill_type_name: rule.bill_type_name, name: `${rule.bill_type_name} - Bulk`, amount: rule.amount, total_paid: 0, remaining_amount: rule.amount, status: 'unpaid', due_date: new Date(Date.now() + 30*86400000).toISOString(), period: new Date().toISOString().slice(0, 7), academic_year: '2024/2025', deleted_at: null });
      totalGenerated++;
    });
    rule.bill_count = (rule.bill_count || 0) + targets.length;
  });
  addAuditLog('BULK_GENERATE', 'BILLING_RULE', `Bulk generate untuk ${rule_ids?.length} aturan, total ${totalGenerated} tagihan`);
  return [200, { status: true, message: `Berhasil generate ${totalGenerated} tagihan!` }];
});

mock.onPost(/\/?finance\/generate-bills\/bulk-cancel/).reply(config => {
  addAuditLog('BULK_CANCEL', 'BILLING_RULE', 'Membatalkan generate tagihan secara masal');
  return [200, { status: true, message: 'Tagihan berhasil dibatalkan' }];
});

// --- FINANCE: PAYMENTS ---
mock.onPost(/\/?finance\/bills\/\d+\/pay-manual/).reply(config => {
  const billId = extractId(config.url);
  const bill = db.bills.find(b => b.id === billId);
  if (bill) {
    const paid = bill.remaining_amount;
    bill.total_paid = bill.amount;
    bill.remaining_amount = 0;
    bill.status = 'paid';
    const paymentId = nextId();
    db.payments.unshift({ id: paymentId, student_id: bill.student_id, student_name: bill.student_name, method: 'Tunai (Manual)', status: 'success', amount: paid, bill_type_names: bill.bill_type_name || bill.name, created_at: new Date().toISOString() });
    addAuditLog('PAYMENT', 'BILL', `Pembayaran manual ${bill.name} untuk ${bill.student_name} (Rp${paid.toLocaleString('id-ID')})`);
    return [200, { status: true, message: 'Pembayaran manual berhasil dicatat', data: { id: paymentId } }];
  }
  return [404, { status: false, message: 'Tagihan tidak ditemukan' }];
});

mock.onPost(/\/?finance\/bills\/\d+\/remind/).reply(config => {
  const billId = extractId(config.url);
  const bill = db.bills.find(b => b.id === billId);
  if (bill) {
    const student = db.students.find(s => s.id === bill.student_id);
    const parent = db.users.find(u => u.id === student?.parent_id);
    if (parent) {
      addNotification(parent.name, parent.phone_number, 'Pengingat Tagihan', `Yth. ${parent.name}, tagihan ${bill.name} untuk ${bill.student_name} masih belum lunas.`);
    }
    addAuditLog('REMINDER', 'BILL', `Mengirim pengingat tagihan ${bill.name} untuk ${bill.student_name}`);
  }
  return [200, { status: true, message: 'Pengingat demo berhasil dikirim' }];
});

mock.onPost(/\/?finance\/payment-intent/).reply(config => {
  const payload = parseData(config);
  const snapToken = `demo-snap-token-${nextId()}`;
  
  if (typeof window !== 'undefined') {
    window._latestPaymentIntent = {
      student_id: payload.student_id,
      amount: payload.amount,
      deposit_applied: payload.deposit_applied || 0,
      bill_ids: payload.bill_ids || []
    };
  }
  
  return [200, { status: true, data: { snap_token: snapToken } }];
});

mock.onPost(/\/?finance\/payments$/).reply(config => {
  const payload = parseData(config);
  const student = db.students.find(s => s.id === payload.student_id);
  const studentName = student ? student.name : 'Siswa';
  
  const paymentId = nextId();
  const items = [];
  
  const targetBillIds = payload.bill_ids || [];
  if (targetBillIds.length > 0) {
    targetBillIds.forEach(billId => {
      const bill = db.bills.find(b => b.id === billId);
      if (bill) {
        const remainingBefore = bill.remaining_amount;
        bill.total_paid = bill.amount;
        bill.remaining_amount = 0;
        bill.status = 'paid';
        items.push({ bill_name: bill.name || bill.bill_type_name, period: bill.period, amount: remainingBefore });
      }
    });
  } else {
    let remaining = payload.amount;
    const unpaidBills = db.bills.filter(b => b.student_id === payload.student_id && b.status !== 'paid');
    unpaidBills.forEach(bill => {
      if (remaining <= 0) return;
      const toPay = bill.remaining_amount;
      bill.total_paid = bill.amount;
      bill.remaining_amount = 0;
      bill.status = 'paid';
      items.push({ bill_name: bill.name || bill.bill_type_name, period: bill.period, amount: toPay });
      remaining -= toPay;
    });
  }

  const newPayment = {
    id: paymentId,
    student_id: payload.student_id || 1,
    student_name: studentName,
    method: payload.method || 'Midtrans',
    status: 'success',
    amount: payload.amount || 0,
    bill_type_names: items.map(item => item.bill_name).join(', ') || 'Pembayaran Tagihan',
    created_at: new Date().toISOString()
  };
  
  db.payments.unshift(newPayment);
  addAuditLog('PAYMENT', 'BILL', `Pembayaran ${payload.method || 'Online'} untuk ${studentName} sebesar Rp${(payload.amount || 0).toLocaleString('id-ID')}`);
  
  if (payload.deposit_applied > 0 && student) {
    student.deposit_balance = Math.max(0, student.deposit_balance - payload.deposit_applied);
  }
  
  return [200, { status: true, message: 'Pembayaran berhasil dikonfirmasi', data: { id: paymentId } }];
});

mock.onGet(/\/?finance\/payments\/([^/]+)\/receipt/).reply(config => {
  const parts = config.url.split('/');
  const paymentId = parseInt(parts[parts.length - 2]);
  const payment = db.payments.find(p => p.id === paymentId);
  
  if (!payment) {
    return [200, {
      status: true,
      data: {
        receipt_number: `SP-${Date.now()}`,
        paid_at: new Date().toISOString(),
        payment_method: 'Tunai',
        student_name: 'Siswa Demo',
        nis: '12345',
        amount: 350000,
        items: [
          { bill_name: 'SPP Bulanan', period: '2026-06', amount: 350000 }
        ]
      }
    }];
  }

  const student = db.students.find(s => s.id === payment.student_id) || {};
  
  return [200, {
    status: true,
    data: {
      receipt_number: `SP-${payment.id}`,
      paid_at: payment.created_at,
      payment_method: payment.method,
      student_name: payment.student_name,
      nis: student.nis || String(payment.student_id),
      amount: payment.amount,
      items: [
        { bill_name: payment.bill_type_names || 'Pembayaran Tagihan', period: '', amount: payment.amount }
      ]
    }
  }];
});

mock.onGet(/\/?finance\/payments/).reply(config => {
  const studentId = config.params?.student_id;
  let data = db.payments;
  if (studentId) data = data.filter(p => p.student_id === parseInt(studentId));
  return [200, { status: true, data: { data, total: data.length } }];
});

mock.onGet(/\/?finance\/bill-summaries/).reply(config => {
  const params = config.params || {};
  const activeBills = getActive(db.bills);
  
  const studentBillsMap = {};
  activeBills.forEach(bill => {
    if (!studentBillsMap[bill.student_id]) {
      studentBillsMap[bill.student_id] = [];
    }
    studentBillsMap[bill.student_id].push(bill);
  });

  const summaries = [];
  const activeStudents = getActive(db.students);
  
  activeStudents.forEach(student => {
    const studentBills = studentBillsMap[student.id] || [];
    if (studentBills.length === 0) return;
    
    let totalAmount = 0;
    let totalPaid = 0;
    let billCount = studentBills.length;
    let paidCount = 0;
    let overdueCount = 0;
    let partialCount = 0;
    let unpaidCount = 0;

    studentBills.forEach(b => {
      totalAmount += b.amount;
      totalPaid += b.total_paid;
      if (b.status === 'paid') paidCount++;
      else if (b.status === 'overdue') overdueCount++;
      else if (b.status === 'partial') partialCount++;
      else unpaidCount++;
    });

    const outstanding = totalAmount - totalPaid;
    let status = 'unpaid';
    if (outstanding === 0 && billCount > 0) status = 'paid';
    else if (overdueCount > 0) status = 'overdue';
    else if (partialCount > 0 || totalPaid > 0) status = 'partial';

    summaries.push({
      id: student.id,
      student_id: student.id,
      student_name: student.name,
      status,
      outstanding,
      bill_count: billCount,
      overdue_count: overdueCount,
      partial_count: partialCount,
      unpaid_count: unpaidCount
    });
  });

  let filtered = summaries;
  if (params.search) {
    const q = params.search.toLowerCase();
    filtered = filtered.filter(s => s.student_name.toLowerCase().includes(q) || String(s.student_id).includes(q));
  }
  if (params.status && params.status !== 'outstanding' && params.status !== 'overdue') {
    filtered = filtered.filter(s => s.status === params.status);
  } else if (params.status === 'outstanding') {
    filtered = filtered.filter(s => s.status !== 'paid');
  } else if (params.status === 'overdue') {
    filtered = filtered.filter(s => s.overdue_count > 0);
  }
  
  const total = filtered.length;
  const page = parseInt(params.page) || 1;
  const limit = parseInt(params.limit) || 10;
  const totalPages = Math.ceil(total / limit) || 1;
  const start = (page - 1) * limit;
  const paginated = filtered.slice(start, start + limit);

  return [200, { status: true, data: { data: paginated, total } }];
});

mock.onGet(/\/?finance\/arrears/).reply(config => {
  const params = paramsFrom(config);
  let arrears = getActive(db.bills)
    .filter(b => b.status !== 'paid')
    .map(b => {
      const student = db.students.find(s => s.id === b.student_id) || {};
      const parent = db.users.find(u => u.id === student.parent_id) || {};
      return {
        ...b,
        parent_name: parent.name || '-',
        parent_phone: parent.phone_number || '-',
        class_name: student.class_name || '-',
        major_name: student.major_name || '-',
        outstanding: b.remaining_amount
      };
    });

  if (params.search) {
    const q = String(params.search).toLowerCase();
    arrears = arrears.filter(item => JSON.stringify(item).toLowerCase().includes(q));
  }

  return [200, { status: true, data: paginate(arrears, params) }];
});

// --- PARENT PORTAL ---
mock.onGet(/\/?parent\/students\/me/).reply(() => [200, { status: true, data: getActive(db.students).filter(s => s.parent_id === 2) }]);
mock.onGet(/\/?parent\/students/).reply(() => [200, { status: true, data: getActive(db.students).filter(s => s.parent_id === 2) }]);
mock.onGet(/\/?parent\/students\/\d+\/class-history/).reply(config => {
  const studentId = extractId(config.url);
  const student = db.students.find(s => s.id === studentId);
  if (student) {
    return [200, { status: true, data: [
      { id: 1, class_name: student.class_name, grade: 'X', academic_year: '2024/2025', is_active: true, created_at: student.created_at }
    ]}];
  }
  return [200, { status: true, data: [] }];
});

mock.onGet(/\/?students\/\d+\/class-history/).reply(config => {
  const studentId = extractId(config.url);
  const student = db.students.find(s => s.id === studentId);
  if (student) {
    return [200, { status: true, data: [
      { id: 1, class_name: student.class_name, grade: 'X', academic_year: '2024/2025', is_active: true, created_at: student.created_at }
    ]}];
  }
  return [200, { status: true, data: [] }];
});

mock.onGet(/\/?finance\/my-bills/).reply(() => {
  const parentStudentIds = getActive(db.students).filter(s => s.parent_id === 2).map(s => s.id);
  const parentBills = getActive(db.bills).filter(b => parentStudentIds.includes(b.student_id));
  return [200, { status: true, data: parentBills }];
});

mock.onGet(/\/?finance\/my-payments/).reply(() => {
  const parentStudentIds = getActive(db.students).filter(s => s.parent_id === 2).map(s => s.id);
  return [200, { status: true, data: db.payments.filter(p => parentStudentIds.includes(p.student_id)) }];
});

// --- USER SPECIAL ACTIONS ---
mock.onPost(/\/?users\/\d+\/reset-password/).reply(200, { status: true, message: 'Password berhasil direset' });
mock.onPost(/\/?users\/\d+\/resend-notification/).reply(200, { status: true, message: 'Notifikasi berhasil dikirim' });
mock.onPost(/\/?users\/bulk-resend-notification/).reply(config => {
  const { ids } = parseData(config);
  return [200, { status: true, data: { sent: ids?.length || 0, failed: 0, errors: [] } }];
});

// --- READ-ONLY LOGS ---
mock.onGet(/\/?notifications(\?.*)?$/).reply(() => [200, { status: true, data: { data: db.notifications, total: db.notifications.length } }]);
mock.onGet(/\/?audit-logs(\?.*)?$/).reply(() => [200, { status: true, data: { data: db.auditLogs, total: db.auditLogs.length } }]);
mock.onGet(/\/?whatsapp\/stats/).reply(() => {
  const countStatus = (items, status) => items.filter(item => item.delivery_status === status).length;
  const sent = countStatus(db.notifications, 'sent') + countStatus(db.notifications, 'delivered') + countStatus(db.notifications, 'read') + countStatus(db.emailNotifications, 'sent') + countStatus(db.emailNotifications, 'delivered') + countStatus(db.emailNotifications, 'read');
  const delivered = countStatus(db.notifications, 'delivered') + countStatus(db.notifications, 'read');
  const read = countStatus(db.notifications, 'read');
  const failed = countStatus(db.notifications, 'failed') + countStatus(db.emailNotifications, 'failed');
  return [200, { status: true, data: {
    SENT: sent,
    DELIVERED: delivered,
    READ: read,
    FAILED: failed,
    whatsapp: {
      total: db.notifications.length,
      delivered: countStatus(db.notifications, 'delivered'),
      sent: countStatus(db.notifications, 'sent'),
      pending: countStatus(db.notifications, 'pending'),
      failed: countStatus(db.notifications, 'failed')
    },
    email: {
      total: db.emailNotifications.length,
      delivered: countStatus(db.emailNotifications, 'delivered'),
      sent: countStatus(db.emailNotifications, 'sent'),
      pending: countStatus(db.emailNotifications, 'pending'),
      failed: countStatus(db.emailNotifications, 'failed')
    }
  }}];
});
mock.onGet(/\/?whatsapp\/logs/).reply(config => {
  const params = paramsFrom(config);
  const channel = params.channel || 'whatsapp';
  let logs = channel === 'email' ? db.emailNotifications : db.notifications;
  if (params.status) logs = logs.filter(log => log.delivery_status === params.status);
  if (params.search) {
    const q = String(params.search).toLowerCase();
    logs = logs.filter(log => JSON.stringify(log).toLowerCase().includes(q));
  }
  return [200, { status: true, data: paginate(logs, params) }];
});
mock.onPost(/\/?whatsapp\/notifications\/\d+\/resend/).reply(config => {
  const id = extractId(config.url);
  const log = [...db.notifications, ...db.emailNotifications].find(item => item.id === id);
  if (log) {
    log.delivery_status = 'delivered';
    log.delivery_error = '';
    log.updated_at = new Date().toISOString();
  }
  return [200, { status: true, message: 'Notifikasi demo berhasil dikirim ulang' }];
});
mock.onGet(/\/?whatsapp\/status/).reply(200, { status: true, data: { status: 'CONNECTED', is_connected: true, device_name: 'Demo WhatsApp Gateway' } });
mock.onGet(/\/?whatsapp\/qr/).reply(() => [200, new Blob([`<svg xmlns="http://www.w3.org/2000/svg" width="240" height="240" viewBox="0 0 240 240"><rect width="240" height="240" fill="#fff"/><rect x="24" y="24" width="64" height="64" fill="#0f172a"/><rect x="152" y="24" width="64" height="64" fill="#0f172a"/><rect x="24" y="152" width="64" height="64" fill="#0f172a"/><path d="M104 104h16v16h-16zm32 0h16v16h-16zm32 0h16v16h-16zM104 136h48v16h-48zm64 0h16v48h-16zm-64 32h16v16h-16zm32 0h16v32h-16zm48 32h32v16h-32z" fill="#0f172a"/><text x="120" y="226" text-anchor="middle" font-family="Arial" font-size="11" font-weight="700" fill="#64748b">Demo QR</text></svg>`], { type: 'image/svg+xml' })]);
mock.onGet(/\/?support\/conversations/).reply(config => {
  const params = paramsFrom(config);
  let conversations = db.supportConversations;
  if (params.status) conversations = conversations.filter(item => item.status === params.status);
  return [200, { status: true, data: paginate(conversations, params) }];
});
mock.onPatch(/\/?support\/conversations\/\d+\/assign/).reply(config => {
  const item = db.supportConversations.find(c => c.id === extractId(config.url));
  if (item) item.status = 'pending';
  return [200, { status: true, message: 'Tiket demo berhasil diambil' }];
});
mock.onPatch(/\/?support\/conversations\/\d+\/close/).reply(config => {
  const item = db.supportConversations.find(c => c.id === extractId(config.url));
  if (item) item.status = 'closed';
  return [200, { status: true, message: 'Tiket demo berhasil ditutup' }];
});
mock.onPatch(/\/?support\/conversations\/\d+\/status/).reply(config => {
  const item = db.supportConversations.find(c => c.id === extractId(config.url));
  if (item) item.status = parseData(config).status || item.status;
  return [200, { status: true, message: 'Status tiket demo diperbarui' }];
});
mock.onPost(/\/?whatsapp/).reply(200, { status: true, message: 'Aksi WhatsApp berhasil' });

// ============================================================
// 6. CATCH ALL (NEVER NETWORK ERROR)
// ============================================================
mock.onGet().reply(200, { status: true, data: { data: [], total: 0 } });
mock.onPost().reply(200, { status: true, message: 'Aksi simulasi berhasil' });
mock.onPut().reply(200, { status: true, message: 'Aksi simulasi berhasil' });
mock.onPatch().reply(200, { status: true, message: 'Aksi simulasi berhasil' });
mock.onDelete().reply(200, { status: true, message: 'Aksi simulasi berhasil' });

if (typeof window !== 'undefined') {
  const injectMidtransSnapMock = () => {
    window.snap = {
      pay(token, callbacks) {
        console.log('Midtrans Snap mock triggered with token:', token);
        
        const existing = document.getElementById('midtrans-snap-mock-modal');
        if (existing) existing.remove();
        
        const intent = window._latestPaymentIntent || { student_id: 1, amount: 350000, deposit_applied: 0, bill_ids: [] };
        
        const modal = document.createElement('div');
        modal.id = 'midtrans-snap-mock-modal';
        modal.style.position = 'fixed';
        modal.style.inset = '0';
        modal.style.zIndex = '99999';
        modal.style.backgroundColor = 'rgba(15, 23, 42, 0.6)';
        modal.style.backdropFilter = 'blur(4px)';
        modal.style.display = 'flex';
        modal.style.alignItems = 'center';
        modal.style.justifyContent = 'center';
        modal.style.fontFamily = 'Inter, sans-serif';
        modal.style.padding = '16px';
        
        modal.innerHTML = `
          <div style="background-color: white; width: 100%; max-w: 420px; border-radius: 16px; overflow: hidden; box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.25); border: 1px solid #e2e8f0; animation: midtransScaleIn 0.3s cubic-bezier(0.34, 1.56, 0.64, 1);">
            <!-- Header -->
            <div style="background-color: #002c5f; padding: 18px 24px; display: flex; align-items: center; justify-content: space-between; color: white;">
              <div style="display: flex; align-items: center; gap: 8px;">
                <span style="font-weight: 800; font-size: 16px; tracking: -0.02em;">midtrans</span>
                <span style="background-color: #00db72; color: #002c5f; font-size: 8px; font-weight: 900; padding: 2px 6px; border-radius: 99px; text-transform: uppercase; letter-spacing: 0.05em;">Sandbox</span>
              </div>
              <div style="font-size: 11px; font-weight: 700; opacity: 0.8;">Simulator Pembayaran</div>
            </div>
            
            <!-- Body -->
            <div style="padding: 24px; text-align: left;">
              <div style="text-align: center; margin-bottom: 24px;">
                <div style="font-size: 11px; color: #64748b; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; margin-bottom: 4px;">Total Tagihan</div>
                <div style="font-size: 24px; font-weight: 800; color: #0f172a;">Rp ${(intent.amount || 0).toLocaleString('id-ID')}</div>
                <div style="font-size: 10px; color: #94a3b8; font-weight: 600; margin-top: 4px;">Token: ${token}</div>
              </div>
              
              <div style="background-color: #f8fafc; border-radius: 12px; padding: 16px; border: 1px solid #f1f5f9; margin-bottom: 24px;">
                <div style="font-size: 11px; font-weight: 800; color: #475569; text-transform: uppercase; margin-bottom: 12px;">Pilih Metode Pembayaran Demo:</div>
                <label style="display: flex; align-items: center; gap: 12px; padding: 10px; background: white; border: 1px solid #cbd5e1; border-radius: 8px; margin-bottom: 8px; cursor: pointer;">
                  <input type="radio" name="pay_method" value="Virtual Account" checked style="accent-color: #002c5f;">
                  <span style="font-size: 12px; font-weight: 700; color: #334155;">Bank Transfer / Virtual Account (Demo)</span>
                </label>
                <label style="display: flex; align-items: center; gap: 12px; padding: 10px; background: white; border: 1px solid #e2e8f0; border-radius: 8px; margin-bottom: 8px; cursor: pointer;">
                  <input type="radio" name="pay_method" value="GoPay" style="accent-color: #002c5f;">
                  <span style="font-size: 12px; font-weight: 700; color: #334155;">GoPay / QRIS (Demo)</span>
                </label>
                <label style="display: flex; align-items: center; gap: 12px; padding: 10px; background: white; border: 1px solid #e2e8f0; border-radius: 8px; cursor: pointer;">
                  <input type="radio" name="pay_method" value="Kartu Kredit" style="accent-color: #002c5f;">
                  <span style="font-size: 12px; font-weight: 700; color: #334155;">Kartu Kredit (Demo)</span>
                </label>
              </div>
              
              <div style="display: flex; flex-direction: column; gap: 10px;">
                <button id="midtrans-pay-btn" style="width: 100%; background-color: #00db72; color: #002c5f; border: none; padding: 12px; border-radius: 10px; font-size: 12px; font-weight: 800; cursor: pointer; transition: all 0.2s; box-shadow: 0 4px 12px rgba(0, 219, 114, 0.2);">
                  Simulasikan Pembayaran Sukses
                </button>
                <button id="midtrans-cancel-btn" style="width: 100%; background-color: #f1f5f9; color: #475569; border: none; padding: 12px; border-radius: 10px; font-size: 12px; font-weight: 800; cursor: pointer; transition: all 0.2s;">
                  Batalkan Pembayaran
                </button>
              </div>
            </div>
            
            <!-- Footer -->
            <div style="background-color: #f8fafc; padding: 14px; text-align: center; border-t: 1px solid #e2e8f0; font-size: 9px; color: #94a3b8; font-weight: 600;">
              Securely processed by Midtrans Simulator
            </div>
          </div>
          
          <style>
            @keyframes midtransScaleIn {
              from { opacity: 0; transform: scale(0.95); }
              to { opacity: 1; transform: scale(1); }
            }
          </style>
        `;
        
        document.body.appendChild(modal);
        
        document.getElementById('midtrans-cancel-btn').onclick = () => {
          modal.remove();
          if (callbacks.onClose) callbacks.onClose();
        };
        
        document.getElementById('midtrans-pay-btn').onclick = async () => {
          const btn = document.getElementById('midtrans-pay-btn');
          btn.disabled = true;
          btn.innerText = 'Memproses Pembayaran...';
          btn.style.opacity = '0.7';
          
          const selectedMethod = document.querySelector('input[name="pay_method"]:checked').value;
          
          try {
            const res = await axios.post('finance/payments', {
              student_id: intent.student_id,
              amount: intent.amount,
              deposit_applied: intent.deposit_applied,
              bill_ids: intent.bill_ids,
              method: `Midtrans (${selectedMethod})`,
              channel: 'gateway',
              reference: `MID-${Date.now()}`
            });
            
            modal.remove();
            
            if (callbacks.onSuccess) {
              callbacks.onSuccess(res.data);
            }
          } catch (err) {
            console.error(err);
            btn.disabled = false;
            btn.innerText = 'Simulasikan Pembayaran Sukses';
            btn.style.opacity = '1';
            alert('Gagal memproses pembayaran simulasi');
          }
        };
      }
    };
  };
  
  if (document.readyState === 'complete' || document.readyState === 'interactive') {
    injectMidtransSnapMock();
  } else {
    window.addEventListener('DOMContentLoaded', injectMidtransSnapMock);
  }
}

export default mock;
