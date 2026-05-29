const XLSX = require('xlsx-js-style');
const fs = require('fs');

const validData = [
  {
    'Timestamp': '21/04/2026 13:24:44',
    'Email Address': 'siswa1@example.com',
    'Nama Lengkap': 'Siswa Contoh Valid 1',
    'Jenis Kelamin': 'Laki - laki',
    'Tempat Lahir': 'PROBOLINGGO',
    'Tanggal Lahir': '01/04/2006',
    'NISN': '0012345671',
    'NIK': '3513124567234567',
    'Agama': 'Islam',
    'RT/RW': '01/01',
    'DESA/KELURAHAN': 'Sukodadi',
    'Kecamatan': 'Paiton',
    'Kabupaten/Kota': 'Probolinggo',
    'Provinsi': 'Jawa Timur',
    'Kode Pos': '67291',
    'No Hp': '0823455789',
    'Email': 'siswa1@example.com',
    'Kelas (Nama Kelas)': 'X-A',
    'Jurusan (Nama Singkat)': 'RPL',
    'NIS': '10001',
    'Angkatan': '2024',
    'Nama Lengkap Ayah': 'Wali Ahmad Zaky',
    'NIK Ayah': '3513125670234567',
    'Tanggal Lahir Ayah': '18/04/1979',
    'Alamat Ayah': 'Probolinggo',
    'Pendidikan Terakhir Ayah': 'S1',
    'Pekerjaan Ayah': 'Programmer',
    'Penghasilan Ayah': '> 5 juta',
    'Nomor WhatsApp (WA) Ayah': '083150507691',
    'Email Aktif Ayah': 'ayah1@example.com',
  },
  {
    'Timestamp': '21/04/2026 13:25:00',
    'Email Address': 'siswa2@example.com',
    'Nama Lengkap': 'Siswa Contoh Valid 2',
    'Jenis Kelamin': 'Perempuan',
    'Tempat Lahir': 'MALANG',
    'Tanggal Lahir': '15/05/2007',
    'NISN': '0012345672',
    'NIK': '3513124567234568',
    'Agama': 'Kristen',
    'RT/RW': '02/03',
    'DESA/KELURAHAN': 'Tlogomas',
    'Kecamatan': 'Lowokwaru',
    'Kabupaten/Kota': 'Malang',
    'Provinsi': 'Jawa Timur',
    'Kode Pos': '65144',
    'No Hp': '0823455790',
    'Email': 'siswa2@example.com',
    'Kelas (Nama Kelas)': 'X-B',
    'Jurusan (Nama Singkat)': 'TKJ',
    'NIS': '10002',
    'Angkatan': '2024',
    'Nama Lengkap Ayah': 'Wali Budi Santoso',
    'NIK Ayah': '3513125670234568',
    'Tanggal Lahir Ayah': '20/05/1980',
    'Alamat Ayah': 'Malang',
    'Pendidikan Terakhir Ayah': 'SMA',
    'Pekerjaan Ayah': 'Swasta',
    'Penghasilan Ayah': '3 - 5 juta',
    'Nomor WhatsApp (WA) Ayah': '083150507692',
    'Email Aktif Ayah': 'ayah2@example.com',
  }
];

const invalidData = [
  {
    'Timestamp': '21/04/2026 13:24:44',
    'Email Address': 'gagal1@example.com',
    'Nama Lengkap': 'Siswa Gagal 1',
    'Jenis Kelamin': 'Laki - laki',
    'Tempat Lahir': 'PROBOLINGGO',
    'Tanggal Lahir': '01/04/2006',
    'NISN': '', // GAGAL: NISN kosong
    'NIK': '3513124567234569',
    'Agama': 'Islam',
    'RT/RW': '01/01',
    'DESA/KELURAHAN': 'Sukodadi',
    'Kecamatan': 'Paiton',
    'Kabupaten/Kota': 'Probolinggo',
    'Provinsi': 'Jawa Timur',
    'Kode Pos': '67291',
    'No Hp': '0823455789',
    'Email': 'gagal1@example.com',
    'Kelas (Nama Kelas)': 'X-A',
    'Jurusan (Nama Singkat)': 'RPL',
    'NIS': '10003',
    'Angkatan': '2024',
    'Nama Lengkap Ayah': 'Wali Ahmad Zaky',
    'NIK Ayah': '3513125670234567',
    'Tanggal Lahir Ayah': '18/04/1979',
    'Alamat Ayah': 'Probolinggo',
    'Pendidikan Terakhir Ayah': 'S1',
    'Pekerjaan Ayah': 'Programmer',
    'Penghasilan Ayah': '> 5 juta',
    'Nomor WhatsApp (WA) Ayah': '083150507691',
    'Email Aktif Ayah': 'ayah1@example.com',
  },
  {
    'Timestamp': '21/04/2026 13:25:00',
    'Email Address': 'gagal2@example.com',
    'Nama Lengkap': 'Siswa Gagal 2',
    'Jenis Kelamin': 'Perempuan',
    'Tempat Lahir': 'MALANG',
    'Tanggal Lahir': '15/05/2007',
    'NISN': '0012345674',
    'NIK': '3513124567234560',
    'Agama': 'Kristen',
    'RT/RW': '02/03',
    'DESA/KELURAHAN': 'Tlogomas',
    'Kecamatan': 'Lowokwaru',
    'Kabupaten/Kota': 'Malang',
    'Provinsi': 'Jawa Timur',
    'Kode Pos': '65144',
    'No Hp': '0823455790',
    'Email': 'gagal2@example.com',
    'Kelas (Nama Kelas)': 'Kelas Hantu', // GAGAL: Kelas Hantu tidak ada
    'Jurusan (Nama Singkat)': 'TKJ',
    'NIS': '10004',
    'Angkatan': '2024',
    'Nama Lengkap Ayah': '', // GAGAL: Nama wali kosong
    'NIK Ayah': '3513125670234568',
    'Tanggal Lahir Ayah': '20/05/1980',
    'Alamat Ayah': 'Malang',
    'Pendidikan Terakhir Ayah': 'SMA',
    'Pekerjaan Ayah': 'Swasta',
    'Penghasilan Ayah': '3 - 5 juta',
    'Nomor WhatsApp (WA) Ayah': '083150507692',
    'Email Aktif Ayah': 'ayah2@example.com',
  }
];


function createExcel(data, sheetName, fileName) {
  const ws = XLSX.utils.json_to_sheet(data);
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, ws, sheetName);
  XLSX.writeFile(wb, fileName);
  console.log('Created:', fileName);
}

const dir = '../test_import_data';
if (!fs.existsSync(dir)){
    fs.mkdirSync(dir);
}

createExcel(validData, 'Data Import', '../test_import_data/Import_Siswa_Wali_Gabungan_Valid.xlsx');
createExcel(invalidData, 'Data Import', '../test_import_data/Import_Siswa_Wali_Gabungan_Gagal.xlsx');
