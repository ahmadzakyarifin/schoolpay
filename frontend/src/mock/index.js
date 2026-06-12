import axios from 'axios';
import MockAdapter from 'axios-mock-adapter';

// Mengatur delay untuk mensimulasikan waktu loading jaringan
const mock = new MockAdapter(axios, { delayResponse: 500 });

// === MOCK AUTHENTICATION ===

// Mock Endpoint: Login
mock.onPost('/auth/login').reply(config => {
  const data = JSON.parse(config.data);
  const { email, password } = data;

  if (email === 'admin@demo.com' && password === 'admin123') {
    return [200, {
      status: true,
      message: 'Login berhasil',
      data: {
        access_token: 'mock_token_admin_12345',
        user: {
          id: 1,
          name: 'Admin Demo',
          email: 'admin@demo.com',
          role: 'admin'
        }
      }
    }];
  }

  if (email === 'parent@demo.com' && password === 'parent123') {
    return [200, {
      status: true,
      message: 'Login berhasil',
      data: {
        access_token: 'mock_token_parent_12345',
        user: {
          id: 2,
          name: 'Parent Demo',
          email: 'parent@demo.com',
          role: 'parent'
        }
      }
    }];
  }

  // Jika kredensial salah
  return [401, {
    status: false,
    message: 'Email atau Password Salah'
  }];
});

// Mock Endpoint: Refresh Token (Dipanggil oleh Vue saat refresh halaman)
mock.onPost('/auth/refresh').reply(config => {
  const authHeader = config.headers['Authorization'];
  
  if (authHeader && authHeader.includes('parent')) {
    return [200, {
      status: true,
      data: {
        access_token: 'mock_token_parent_12345',
        user: { id: 2, name: 'Parent Demo', email: 'parent@demo.com', role: 'parent' }
      }
    }];
  }

  if (authHeader && authHeader.includes('admin')) {
     return [200, {
      status: true,
      data: {
        access_token: 'mock_token_admin_12345',
        user: { id: 1, name: 'Admin Demo', email: 'admin@demo.com', role: 'admin' }
      }
    }];
  }

  return [401, { status: false, message: 'Sesi kedaluwarsa' }];
});

// Mock Endpoint: Logout
mock.onPost('/auth/logout').reply(200, { status: true, message: 'Logout berhasil' });

// Pass Through: Biarkan request lain (yang belum di-mock) diteruskan seperti biasa
// Catatan: Jika backend tidak nyala, request ini akan menghasilkan Network Error, 
// tapi teman Anda tetap bisa melihat state "Error" di UI.
mock.onAny().passThrough();

export default mock;
