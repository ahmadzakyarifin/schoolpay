import axios from 'axios';
import MockAdapter from 'axios-mock-adapter';

const mock = new MockAdapter(axios, { delayResponse: 400 });

// ==========================================
// 1. DUMMY DATA BASE (IN-MEMORY DATABASE)
// ==========================================
const db = {
  users: [
    { id: 1, name: 'Admin Utama', email: 'admin@demo.com', role: 'admin', is_active: true, created_at: new Date().toISOString() },
    { id: 2, name: 'Bapak Parent', email: 'parent@demo.com', role: 'parent', is_active: true, created_at: new Date().toISOString() }
  ],
  academicYears: [
    { id: 1, year: '2024/2025', status: 'active', is_active: true },
    { id: 2, year: '2025/2026', status: 'planned', is_active: true }
  ],
  majors: [
    { id: 1, name: 'MIPA', code: 'MIPA', is_active: true },
    { id: 2, name: 'IPS', code: 'IPS', is_active: true }
  ],
  classes: [
    { id: 1, name: 'X MIPA 1', major_id: 1, academic_year_ids: [1], is_active: true },
    { id: 2, name: 'X IPS 1', major_id: 2, academic_year_ids: [1], is_active: true }
  ],
  students: [
    { id: 1, name: 'Budi Santoso', nisn: '0012345678', gender: 'Laki-laki', class_id: 1, class_name: 'X MIPA 1', parent_name: 'Bapak Parent', parent_email: 'parent@demo.com', is_active: true, deposit_balance: 500000, whatsapp: '081234567890' },
    { id: 2, name: 'Andi Saputra', nisn: '0012345679', gender: 'Laki-laki', class_id: 2, class_name: 'X IPS 1', parent_name: 'Bapak Andi', parent_email: 'andi@demo.com', is_active: true, deposit_balance: 0, whatsapp: '081234567891' }
  ],
  billTypes: [
    { id: 1, name: 'SPP Bulanan', type: 'MONTHLY', amount: 350000, is_active: true },
    { id: 2, name: 'Uang Gedung', type: 'ONE_TIME', amount: 2000000, is_active: true }
  ],
  billingRules: [
    { id: 1, name: 'Aturan SPP 2024', bill_type_id: 1, target_type: 'ALL', amount: 350000, is_active: true }
  ],
  bills: [
    { id: 1, student_id: 1, student_name: 'Budi Santoso', name: 'SPP Bulan Ini', amount: 350000, remaining_amount: 350000, status: 'unpaid', due_date: new Date(Date.now() + 5*24*60*60*1000).toISOString(), period: '2026-06' },
    { id: 2, student_id: 1, student_name: 'Budi Santoso', name: 'Uang Gedung', amount: 2000000, remaining_amount: 1500000, status: 'partial', due_date: new Date(Date.now() + 30*24*60*60*1000).toISOString(), period: '2026-06' },
    { id: 3, student_id: 1, student_name: 'Budi Santoso', name: 'Daftar Ulang', amount: 500000, remaining_amount: 0, status: 'paid', due_date: new Date(Date.now() - 10*24*60*60*1000).toISOString(), period: '2026-06' },
    { id: 4, student_id: 2, student_name: 'Andi Saputra', name: 'SPP Bulan Ini', amount: 350000, remaining_amount: 350000, status: 'unpaid', due_date: new Date(Date.now() + 5*24*60*60*1000).toISOString(), period: '2026-06' }
  ],
  payments: [
    { id: 1, student_id: 1, student_name: 'Budi Santoso', method: 'Tunai', status: 'success', amount: 500000, bill_type_names: 'Daftar Ulang', created_at: new Date(Date.now() - 5*24*60*60*1000).toISOString() }
  ],
  notifications: [
    { id: 1, type: 'whatsapp', recipient_name: 'Bapak Parent', recipient_phone: '081234567890', title: 'Tagihan SPP Baru', message: 'Yth. Bapak Parent, terdapat tagihan baru sebesar Rp350.000...', delivery_status: 'delivered', created_at: new Date().toISOString() }
  ],
  auditLogs: [
    { id: 1, user_name: 'Admin Utama', action: 'CREATE', entity: 'STUDENT', description: 'Menambahkan siswa Budi Santoso', created_at: new Date().toISOString() }
  ]
};

// ==========================================
// 2. HELPER FUNGSI CRUD OTOMATIS
// ==========================================
function createMockCRUD(endpoint, collectionName) {
  const collection = db[collectionName];
  
  // GET ALL / GET BY PARAMS
  mock.onGet(new RegExp(`^/${endpoint}(?:\\?.*)?$`)).reply(config => {
    return [200, { status: true, data: { data: collection, total: collection.length, page: 1, totalPages: 1 } }];
  });

  // GET SINGLE
  mock.onGet(new RegExp(`^/${endpoint}/\\d+$`)).reply(config => {
    const id = parseInt(config.url.split('/').pop());
    const item = collection.find(i => i.id === id);
    if (item) return [200, { status: true, data: item }];
    return [404, { status: false, message: 'Data tidak ditemukan' }];
  });

  // GET DEPENDENCY INFO (Bypass hapus)
  mock.onGet(new RegExp(`^/${endpoint}/\\d+/dependency-info$`)).reply(200, { status: true, data: { has_dependencies: false } });

  // POST (CREATE)
  mock.onPost(new RegExp(`^/${endpoint}$`)).reply(config => {
    const data = JSON.parse(config.data || '{}');
    const newItem = { id: Date.now(), ...data, created_at: new Date().toISOString(), is_active: true };
    collection.unshift(newItem);
    return [200, { status: true, message: 'Data berhasil ditambahkan', data: newItem }];
  });

  // PUT (UPDATE)
  mock.onPut(new RegExp(`^/${endpoint}/\\d+$`)).reply(config => {
    const id = parseInt(config.url.split('/').pop());
    const data = JSON.parse(config.data || '{}');
    const idx = collection.findIndex(i => i.id === id);
    if (idx !== -1) {
      collection[idx] = { ...collection[idx], ...data, updated_at: new Date().toISOString() };
      return [200, { status: true, message: 'Data berhasil diperbarui', data: collection[idx] }];
    }
    return [404, { status: false, message: 'Data tidak ditemukan' }];
  });

  // DELETE
  mock.onDelete(new RegExp(`^/${endpoint}/\\d+$`)).reply(config => {
    const id = parseInt(config.url.split('/').pop());
    const idx = collection.findIndex(i => i.id === id);
    if (idx !== -1) collection.splice(idx, 1);
    return [200, { status: true, message: 'Data berhasil dihapus' }];
  });

  // BULK DELETE
  mock.onPost(new RegExp(`^/${endpoint}/bulk-delete$`)).reply(config => {
    const { ids } = JSON.parse(config.data || '{}');
    if (ids && ids.length) {
      ids.forEach(id => {
        const idx = collection.findIndex(i => i.id === id);
        if (idx !== -1) collection.splice(idx, 1);
      });
    }
    return [200, { status: true, message: 'Data terpilih berhasil dihapus' }];
  });
  
  // TOGGLE STATUS
  mock.onPatch(new RegExp(`^/${endpoint}/\\d+/status$`)).reply(config => {
    const id = parseInt(config.url.split('/')[2]);
    const idx = collection.findIndex(i => i.id === id);
    if (idx !== -1) {
      collection[idx].is_active = !collection[idx].is_active;
      return [200, { status: true, message: 'Status diubah', data: collection[idx] }];
    }
    return [404, { status: false, message: 'Tidak ditemukan' }];
  });

  // RESTORE
  mock.onPatch(new RegExp(`^/${endpoint}/\\d+/restore$`)).reply(config => {
    return [200, { status: true, message: 'Berhasil dipulihkan' }];
  });

  // CHECK UNIQUE
  mock.onGet(new RegExp(`^/${endpoint}/check-unique`)).reply(200, { status: true, data: { is_unique: true } });
}

// ==========================================
// 3. TERAPKAN CRUD HELPER KE SEMUA MENU
// ==========================================
createMockCRUD('users', 'users');
createMockCRUD('students', 'students');
createMockCRUD('academic/years', 'academicYears');
createMockCRUD('academic/major', 'majors');
createMockCRUD('academic/class', 'classes');
createMockCRUD('finance/bill-types', 'billTypes');
createMockCRUD('finance/billing-rules', 'billingRules');
createMockCRUD('finance/bills', 'bills');

// ==========================================
// 4. KUSTOMISASI KHUSUS (OVERRIDE)
// ==========================================

// GENERATE BILLS SIMULATION
mock.onPost('/finance/generate-bills').reply(200, { status: true, message: 'Tagihan berhasil digenerate' });
mock.onPost('/finance/generate-bills/bulk').reply(200, { status: true, message: 'Tagihan berhasil digenerate secara masal' });
mock.onPost('/finance/generate-bills/bulk-cancel').reply(200, { status: true, message: 'Tagihan berhasil dibatalkan' });

// A. AUTENTIKASI
mock.onPost('/auth/login').reply(config => {
  const { email, password } = JSON.parse(config.data);
  if (email === 'admin@demo.com' && password === 'admin123') {
    return [200, { status: true, message: 'Login berhasil', data: { access_token: 'token_admin', user: db.users[0] } }];
  }
  if (email === 'parent@demo.com' && parent123) {
    return [200, { status: true, message: 'Login berhasil', data: { access_token: 'token_parent', user: db.users[1] } }];
  }
  return [401, { status: false, message: 'Email atau Password Salah' }];
});
mock.onPost('/auth/refresh').reply(config => {
  const authHeader = config.headers['Authorization'] || '';
  if (authHeader.includes('parent')) return [200, { status: true, data: { access_token: 'token_parent', user: db.users[1] } }];
  return [200, { status: true, data: { access_token: 'token_admin', user: db.users[0] } }];
});
mock.onPost('/auth/logout').reply(200, { status: true });
mock.onGet('/auth/me').reply(200, { status: true, data: db.users[0] });

// B. DASHBOARD ADMIN STATS
mock.onGet(/\/dashboard\/stats.*/).reply(200, {
  status: true,
  data: {
    stats: {
      students: { total: db.students.length },
      users: { total: db.users.length, new_this_period: 0 },
      unpaid_amount: db.bills.filter(b => b.status !== 'paid').reduce((a, b) => a + b.remaining_amount, 0),
      paid_amount: db.payments.reduce((a, p) => a + p.amount, 0),
      paid_count: db.payments.length
    },
    recent_payments: db.payments,
    recent_notifications: db.notifications
  }
});

// C. PORTAL PARENT
mock.onGet(/\/parent\/students.*/).reply(200, { status: true, data: db.students.filter(s => s.id === 1) });
mock.onGet(/\/finance\/my-bills.*/).reply(200, { status: true, data: { data: db.bills.filter(b => b.student_id === 1) } });
mock.onGet(/\/finance\/my-payments.*/).reply(200, { status: true, data: { data: db.payments.filter(p => p.student_id === 1) } });

// Simulasi Pembayaran
mock.onPost('/finance/payment-intent').reply(200, { status: true, data: { snap_token: 'dummy-token' } });
mock.onPost('/finance/payments').reply(config => {
  const payload = JSON.parse(config.data);
  payload.items.forEach(item => {
    const bill = db.bills.find(b => b.id === item.bill_id);
    if (bill) {
      bill.remaining_amount -= item.amount;
      if (bill.remaining_amount <= 0) bill.status = 'paid';
    }
  });
  db.payments.unshift({ id: Date.now(), student_id: 1, student_name: 'Budi Santoso', method: 'Midtrans', status: 'success', amount: payload.total_amount, bill_type_names: 'Pembayaran Online', created_at: new Date().toISOString() });
  return [200, { status: true, message: 'Pembayaran berhasil' }];
});
mock.onPost(/\/finance\/bills\/\d+\/pay-manual/).reply(config => {
  const billId = parseInt(config.url.split('/')[3]);
  const bill = db.bills.find(b => b.id === billId);
  if (bill) {
    bill.remaining_amount = 0;
    bill.status = 'paid';
    db.payments.unshift({ id: Date.now(), student_id: bill.student_id, student_name: bill.student_name, method: 'Tunai', status: 'success', amount: bill.amount, bill_type_names: bill.name, created_at: new Date().toISOString() });
    return [200, { status: true, message: 'Pembayaran manual berhasil' }];
  }
  return [404, { status: false, message: 'Tagihan tidak ditemukan' }];
});

// D. READ-ONLY LOGS & EXTRAS
mock.onGet(/\/notifications(\?.*)?$/).reply(200, { status: true, data: { data: db.notifications, total: db.notifications.length } });
mock.onGet(/\/audit-logs(\?.*)?$/).reply(200, { status: true, data: { data: db.auditLogs, total: db.auditLogs.length } });
mock.onGet(/\/whatsapp\/.*/).reply(200, { status: true, data: { is_connected: true, device_name: 'Mock Device' } });

// ==========================================
// 5. CATCH ALL (MENCEGAH ERROR UI)
// ==========================================
mock.onGet().reply(200, { status: true, data: { data: [], total: 0 } });
mock.onPost().reply(200, { status: true, message: 'Tindakan berhasil disimulasikan' });
mock.onPut().reply(200, { status: true, message: 'Tindakan berhasil disimulasikan' });
mock.onDelete().reply(200, { status: true, message: 'Tindakan berhasil disimulasikan' });

export default mock;
