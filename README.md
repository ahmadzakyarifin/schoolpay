# 🏛️ Catatan Arsitektur Backend

Dokumen ini menjelaskan dokumentasi arsitektur backend, alur inisialisasi aplikasi, serta pembagian tanggung jawab (*Separation of Concerns*) pada *entry point* aplikasi.

---

## 🚀 Peran `main.go`

`main.go` bertindak sebagai **pintu masuk utama (entry point)** aplikasi. Mengikuti prinsip *Single Responsibility*, tugas utamanya hanyalah menyiapkan *dependency* global yang krusial dan menjalankan server.

### 🔄 Alur Inisialisasi:
1. **Load Config** — Membaca seluruh konfigurasi environment.
2. **Init Logger** — Menyiapkan sistem logging aplikasi.
3. **Connect Database** — Membuka koneksi ke *database engine*.
4. **Build Application** — Merakit instance aplikasi melalui `app.NewApp(db, cfg)`.
5. **Run Server** — Menjalankan HTTP server/gRPC server.

---

## 🏗️ Keputusan Desain & Arsitektur

### 1. Kenapa koneksi database dilakukan di `main.go`?
Database merupakan *core dependency* (ketergantungan utama). Jika koneksi database gagal atau mengalami *timeout*, aplikasi harus langsung dihentikan (*fail fast*). Dengan menaruhnya di `main.go`, kita mencegah aplikasi berjalan dalam kondisi pincang, dan fungsi penentu seperti `NewApp()` tidak perlu dieksekusi secara sia-sia.

### 2. Kenapa tidak melakukan koneksi database di dalam `NewApp()`?
* **Fokus Tanggung Jawab:** `NewApp()` dirancang khusus untuk fokus merakit komponen internal aplikasi (seperti routing HTTP framework, middleware, inisialisasi repository, service, handler, websocket, worker, hingga scheduler).
* **Kemudahan Pengujian (Testability):** Dengan memisahkan inisialisasi database di luar `NewApp()`, kita menerapkan pola *Dependency Injection*. Saat melakukan *Integration Testing*, kita bisa dengan sangat mudah menyuntikkan (inject) database tiruan (*test database* atau *mock*) ke dalam `NewApp()` tanpa mengubah kode core aplikasi.

---

## 🧭 Peran `app.NewApp()`

`app.NewApp()` berperan sebagai **komposer aplikasi**, bukan tempat menaruh seluruh detail fitur. File ini dibuat agar mudah dibaca dari atas ke bawah seperti alur startup server.

### Urutan kerja yang diharapkan:
1. **Setup HTTP Engine** — Menyiapkan Gin, mode production, CORS, static files, dan middleware global.
2. **Buat Komponen Global** — Menyiapkan messenger dan websocket hub yang dipakai banyak fitur.
3. **Buat API Group** — Menyiapkan `/api`, rate limit global, dan idempotency middleware.
4. **Rakit Dependency Bersama** — Memanggil `buildAppServices()` untuk service/repository yang dipakai lintas router.
5. **Daftarkan WebSocket Global** — Membuka `/api/ws` untuk event realtime lintas role.
6. **Jalankan Background Jobs** — Menjalankan scheduler, database worker, dan cleanup idempotency.
7. **Daftarkan Router Fitur** — Auth, admin, parent, finance, webhook, dan swagger.

### Apa yang boleh ada di `app.go`?
* **Middleware global:** CORS, rate limit global, idempotency middleware.
* **Komponen global:** Messenger, websocket hub, API group.
* **Route infrastruktur:** WebSocket global dan Swagger.
* **Pemanggilan router:** `auth`, `admin`, `parent`, `finance`, `webhook`.

### Apa yang sebaiknya tidak ditaruh di `app.go`?
* Detail endpoint fitur seperti `POST /finance/payments`.
* Logic bisnis pembayaran, tagihan, user, atau notifikasi.
* Query database langsung.
* Handler inline yang panjang.

Detail fitur ditempatkan di router atau module masing-masing agar pembacaan lebih jelas:

| Area | Lokasi | Alasan |
| :--- | :--- | :--- |
| Route pembayaran parent/admin | `internal/router/finance` | Endpoint finance dikumpulkan dekat handler finance. |
| Business logic pembayaran | `internal/module/finance/usecase` | Aturan pembayaran, transaksi DB, dan validasi bisnis berada di service. |
| Query pembayaran | `internal/module/finance/repository` | Query DB dipisahkan dari handler dan service. |
| Dependency bersama | `internal/app/services.go` | Service yang dipakai lintas router dibuat sekali agar tidak duplikatif. |

---

## 🧱 Kenapa Ada `services.go`?

`services.go` adalah tempat merakit **dependency bersama** (*shared dependencies*). Dependency seperti `PaymentService`, `StudentBillService`, `FinanceNotificationService`, `AuditLogService`, dan `UserRepo` dipakai oleh lebih dari satu router, misalnya admin, finance, webhook, worker, dan middleware auth.

Jika semua dependency ini dibuat ulang di masing-masing router, aplikasi tetap bisa berjalan, tetapi alurnya menjadi sulit dilacak karena service yang sama dibuat di banyak tempat. Dengan `services.go`, aplikasi punya satu titik komposisi yang jelas.

> [!NOTE]
> Router tetap boleh membuat repository/service sendiri jika dependency tersebut hanya dipakai oleh router itu saja. Dependency yang dipakai lintas fitur lebih baik dirakit di `services.go`.

---

## 🧩 Perbedaan `config` dan `infrastructure`

Untuk menjaga kode tetap modular, tanggung jawab pembacaan data dan koneksi teknis dipisahkan ke dalam dua entitas yang berbeda:

| Komponen | Tanggung Jawab Utama | Contoh Batasan Kerja |
| :--- | :--- | :--- |
| **`config`** | Membaca dan memvalidasi environment variable dari sistem (`.env` atau OS). | Memuat kredensial DB, JWT Secret, konfigurasi WAHA, SMTP, Midtrans, dan Frontend URL. |
| **`infrastructure`** | Menyediakan driver dan membuat koneksi teknis ke sistem eksternal atau pihak ketiga. | Menginisialisasi koneksi MySQL/MariaDB menggunakan database driver (misal: Bun, `sqlx`, atau driver native). |

---

## 🛡️ Lapisan Keamanan: CORS Middleware

Aplikasi kami menggunakan *middleware* kustom untuk mengelola **Cross-Origin Resource Sharing (CORS)**. Ini adalah lapisan keamanan yang membatasi situs mana yang diizinkan untuk berinteraksi dengan API kami.

### Apa saja yang diizinkan?
1. **Origin:** Hanya domain yang terdaftar dalam `cfg.FrontendURL` yang diizinkan mengirim *request*.
2. **Methods:** `POST`, `GET`, `OPTIONS`, `PUT`, `DELETE`, `PATCH`.
3. **Headers:** `Content-Type`, `Authorization`, `X-Requested-With`, `Accept`, `Origin`, `X-Idempotency-Key`.
4. **Credentials:** `Access-Control-Allow-Credentials: true` (mengizinkan pengiriman *cookie* atau *auth token*).

### Mengapa perlu penanganan `OPTIONS`?
```go
if c.Request.Method == http.MethodOptions {
    c.AbortWithStatus(http.StatusNoContent) 
    return
}
```

* **Fungsi:** Sebelum browser mengirim data sensitif (seperti POST atau DELETE), browser akan mengirim request "tes" dengan method OPTIONS (Preflight Request).

* **Peran Kode:** Kode di atas memberikan jawaban cepat "Ya, Anda diizinkan" dengan status 204 No Content tanpa harus menjalankan business logic (menghemat sumber daya server).

* **Pentingnya return:** return memastikan proses dihentikan seketika setelah memberikan izin, sehingga server tidak lanjut memproses request OPTIONS sebagai request data biasa yang bisa menyebabkan error.

## ⚠️ Analisis Kombinasi Rate Limiter Client Identity

Bagian ini mendokumentasikan **analisis kelemahan** dari setiap kombinasi *client identity* yang digunakan pada Rate Limiter. Tujuannya adalah menjelaskan mengapa kombinasi tertentu **tidak ideal** dan bagaimana *attacker* dapat mengeksploitasinya.

> [!IMPORTANT]
> Pemilihan kombinasi client identity pada rate limiter sangat berpengaruh terhadap **keamanan akun pengguna** dan **ketahanan server**. Kombinasi yang salah bisa membuat pengguna sah terkunci (*false positive*) atau justru memberi celah bagi attacker untuk mem-*bypass* limit.

---

### 1. 🌐 IP Saja

**Ruang lingkup limit:** Berdasarkan alamat IP pengirim request saja.

**Cara kerja:** Semua request dari IP yang sama dianggap berasal dari satu entitas. Jika limit tercapai, seluruh IP tersebut diblokir.

| Aspek | Penjelasan |
| :--- | :--- |
| **Dampak ke pemilik akun** | Jika ada orang lain di **jaringan yang sama** (misal Wi-Fi kampus) yang melakukan spam login, maka **seluruh pengguna di IP tersebut ikut kena limit**, padahal mereka tidak melakukan kesalahan apa pun. Pemilik akun yang sah harus **pindah ke IP lain** (misal pakai data seluler) agar bisa login kembali. |
| **Dampak ke pengguna lain** | Pengguna lain di jaringan yang sama juga **terkena dampak** meskipun tidak ada hubungannya. Ini adalah *false positive* yang tinggi. |

**🔓 Cara Attacker Mem-bypass:**

```
Skenario: Attacker ingin brute-force login akun korban

1. Attacker mengirim request login sebanyak N kali dari IP-A
2. IP-A kena rate limit → Attacker DIBLOKIR
3. Attacker pindah ke IP-B (ganti Wi-Fi, pakai VPN, atau pakai data seluler)
4. Counter reset → Attacker bisa spam lagi dari awal
5. Ulangi: ganti IP → spam → kena limit → ganti IP lagi
```

> [!WARNING]
> **Kesimpulan:** IP yang "terkontaminasi" merugikan pengguna sah yang berada di jaringan yang sama. Sementara attacker cukup **ganti IP** untuk melanjutkan serangan.

---

### 2. 📧 Email Saja

**Ruang lingkup limit:** Berdasarkan alamat email target login saja.

**Cara kerja:** Semua request yang mencoba login ke email tertentu dihitung, **tanpa memandang IP atau device** yang digunakan.

| Aspek | Penjelasan |
| :--- | :--- |
| **Dampak ke pemilik akun** | Pemilik akun akan **terkunci dari akunnya sendiri** jika ada attacker yang mengirim banyak request login ke emailnya. Meskipun pemilik akun berpindah IP atau ganti device, dia tetap tidak bisa login karena limit-nya ada di level email. |
| **Keunggulan** | Attacker **tidak bisa mengakali** limit hanya dengan ganti IP atau ganti device untuk menyerang email yang sama. |

**🔓 Cara Attacker Mem-bypass:**

```
Skenario: Attacker ingin mengincar banyak akun

1. Attacker spam login ke email-korban-1@mail.com → kena limit
2. Attacker TIDAK bisa lanjut ke email-korban-1 (limit per email)
3. Tapi attacker bisa PINDAH TARGET ke email-korban-2@mail.com
4. Spam lagi → kena limit → pindah ke email-korban-3@mail.com
5. Ulangi: ganti target email → spam → kena limit → ganti target lagi
```

> [!WARNING]
> **Kesimpulan:** Sangat efektif melindungi **satu akun dari brute-force**, tapi attacker bisa berpindah target ke email lain. Dan yang paling bahaya: pemilik akun yang sah **ikut terkunci** karena limit ada di level email — ini bisa menjadi serangan *Denial of Service* terhadap akun tertentu.

---

### 3. 📱 Device Saja

**Ruang lingkup limit:** Berdasarkan *fingerprint* device/browser saja.

**Cara kerja:** Semua request dari device (browser) yang sama dihitung. Tidak berpengaruh pada IP atau email manapun.

| Aspek | Penjelasan |
| :--- | :--- |
| **Dampak ke pemilik akun** | **Tidak ada dampak langsung** ke pemilik akun. Limit hanya berlaku di device attacker, bukan di akun korban. Pemilik akun tetap bisa login dari device miliknya sendiri tanpa masalah. |
| **Kelemahan** | Tidak memberikan perlindungan nyata terhadap akun target. Attacker bebas mencoba berbagai email dari device berbeda. |

**🔓 Cara Attacker Mem-bypass:**

```
Skenario: Attacker ingin brute-force login akun korban

1. Attacker spam login dari Chrome → kena limit di Chrome
2. Attacker buka Firefox → counter reset dari 0
3. Spam lagi dari Firefox → kena limit → buka Safari / Edge
4. Atau buka Incognito mode / ganti User-Agent
5. Ulangi: ganti browser/device → spam → kena limit → ganti lagi
```

> [!WARNING]
> **Kesimpulan:** Limit di device **tidak melindungi akun korban** sama sekali. Attacker cukup **ganti browser** atau ubah *fingerprint* device untuk melanjutkan serangan. Tidak ada pengaruh ke level IP maupun email.

---

### 4. 🌐📧 IP + Email

**Ruang lingkup limit:** Berdasarkan kombinasi IP **dan** Email secara bersamaan.

**Cara kerja:** Limit hanya berlaku ketika **IP tertentu** mencoba login ke **email tertentu**. Jika salah satu berubah (IP berbeda atau email berbeda), counter dianggap baru.

| Aspek | Penjelasan |
| :--- | :--- |
| **Dampak ke pemilik akun** | Pemilik akun hanya kena limit **jika berada di IP yang sama dengan attacker**. Contoh: Jika attacker dan korban sama-sama di Wi-Fi kampus, maka korban tidak bisa login dari Wi-Fi tersebut. Tapi korban **bisa login dari IP lain** (misal data seluler). |
| **Cakupan limit** | Lebih sempit dibanding IP saja atau Email saja. Hanya kombinasi IP+Email spesifik yang diblokir. |

**🔓 Cara Attacker Mem-bypass:**

```
Skenario: Attacker ingin brute-force akun korban@mail.com

1. Attacker spam login ke korban@mail.com dari IP-A
2. Kombinasi (IP-A + korban@mail.com) kena limit
3. Attacker pindah ke IP-B (ganti Wi-Fi/VPN)
4. Kombinasi (IP-B + korban@mail.com) → counter mulai dari 0!
5. Spam lagi → kena limit → pindah ke IP-C
6. Ulangi: ganti IP → spam email yang sama → kena limit → ganti IP lagi
```

**Skenario bahaya — Attacker & Korban di jaringan yang sama:**

```
Kondisi: Attacker dan korban sama-sama di Wi-Fi kampus (IP-KAMPUS)

1. Attacker spam login ke korban@mail.com dari IP-KAMPUS
2. Kombinasi (IP-KAMPUS + korban@mail.com) kena limit
3. Korban yang juga di IP-KAMPUS → TIDAK bisa login ke akunnya sendiri!
4. Solusi korban: pindah ke data seluler (IP berbeda)
5. Jika attacker juga pindah IP → serangan bisa dilanjutkan
```

> [!WARNING]
> **Kesimpulan:** Attacker cukup **ganti IP** untuk meneruskan serangan ke email yang sama. Korban terdampak hanya ketika berada di **jaringan yang sama** dengan attacker.

---

### 5. 📱📧 Device + Email

**Ruang lingkup limit:** Berdasarkan kombinasi Device **dan** Email secara bersamaan.

**Cara kerja:** Limit berlaku ketika **device tertentu** mencoba login ke **email tertentu**. Mirip konsepnya dengan IP+Email, tapi scope-nya di device bukan IP.

| Aspek | Penjelasan |
| :--- | :--- |
| **Dampak ke pemilik akun** | Pemilik akun **tidak terdampak** selama menggunakan device-nya sendiri. Limit hanya berlaku di device attacker untuk email korban. |
| **Cakupan limit** | Terbatas pada kombinasi device+email yang spesifik. |

**🔓 Cara Attacker Mem-bypass:**

```
Skenario: Attacker ingin brute-force akun korban@mail.com

1. Attacker spam login ke korban@mail.com dari Chrome
2. Kombinasi (Chrome + korban@mail.com) kena limit
3. Attacker buka Firefox → kombinasi baru (Firefox + korban@mail.com)
4. Counter mulai dari 0 → spam lagi
5. Kena limit → buka Safari / Edge / Incognito
6. Ulangi: ganti device/browser → spam email sama → kena limit → ganti lagi

Untuk menyerang banyak email:
1. Cukup pakai 1 device saja
2. Spam email-1 → kena limit → pindah ke email-2
3. Device yang sama bisa dipakai karena kombinasinya baru
```

> [!WARNING]
> **Kesimpulan:** Cakupannya mirip dengan IP+Email. Attacker cukup **ganti device/browser** untuk melanjutkan serangan ke email yang sama. Atau tetap di device yang sama tapi **ganti target email**.

---

### 6. 🌐📱 IP + Device

**Ruang lingkup limit:** Berdasarkan kombinasi IP **dan** Device saja (**tanpa memperhitungkan email**).

**Cara kerja:** Limit berlaku ketika IP tertentu dan device tertentu mengirim request. **Tidak ada perlindungan di level akun/email.**

| Aspek | Penjelasan |
| :--- | :--- |
| **Dampak ke pemilik akun** | **Tidak ada perlindungan** untuk akun. Karena email tidak diperhitungkan, attacker bisa mencoba login ke email manapun tanpa batas selama kombinasi IP+Device-nya berbeda. |
| **Kelemahan fatal** | Kombinasi IP+Device sangat mudah diubah dan **tidak melindungi akun sama sekali**. |

**🔓 Cara Attacker Mem-bypass:**

```
Skenario: Attacker ingin brute-force akun korban@mail.com

1. Attacker spam dari (IP-A + Chrome) → kena limit
2. Ganti browser: (IP-A + Firefox) → counter baru! Spam lagi
3. Ganti IP: (IP-B + Chrome) → counter baru! Spam lagi
4. Ganti keduanya: (IP-B + Firefox) → counter baru lagi!

Jumlah kombinasi yang bisa dibuat:
   N (jumlah IP) × M (jumlah Device) = total kombinasi unik

Contoh: 5 IP × 4 Browser = 20 kombinasi unik
→ Jika limit = 5 request, attacker bisa kirim 100 request!
```

> [!CAUTION]
> **Kesimpulan — PALING BERBAHAYA:** Kombinasi IP+Device **sama sekali tidak melindungi akun**. Attacker bisa melakukan spam dengan sangat banyak kombinasi sampai **server overload/jebol**. Ini bukan hanya ancaman di level akun, tapi ancaman **di level server** karena banyaknya kombinasi unik yang bisa dibuat. Attacker bisa mengirim request dalam jumlah sangat besar tanpa pernah kena limit yang berarti.

---

### 📊 Tabel Ringkasan Perbandingan

| Kombinasi | Pemilik Akun Terkunci? | Cara Attacker Bypass | Tingkat Bahaya |
| :--- | :--- | :--- | :--- |
| **IP saja** | ✅ Ya (jika 1 jaringan) | Ganti IP (VPN/data seluler) | ⚠️ Sedang |
| **Email saja** | ✅ Ya (selalu) | Ganti target email | ⚠️ Sedang |
| **Device saja** | ❌ Tidak | Ganti browser/device | ⚠️ Sedang |
| **IP + Email** | ✅ Ya (jika 1 jaringan) | Ganti IP | ⚠️ Sedang |
| **Device + Email** | ❌ Tidak | Ganti browser/device | ⚠️ Sedang |
| **IP + Device** | ❌ Tidak | Ganti IP dan/atau browser | 🔴 **Sangat Tinggi** |

> [!CAUTION]
> Kombinasi **IP + Device** adalah yang **paling berbahaya** karena:
> 1. Tidak ada perlindungan akun (email tidak diperhitungkan)
> 2. Attacker bisa membuat **sangat banyak kombinasi unik**
> 3. Spam bisa terjadi sampai **server jebol/overload** — bukan hanya akun yang terancam, tapi **seluruh sistem**

---

### 🛡️ Kombinasi Langsung (IP + Device + Email) vs. Rate Limiter Bertingkat (Multi-tier)

Seringkali muncul pertanyaan: **"Mengapa tidak menggabungkan saja IP + Device + Email secara langsung menjadi satu key?"** atau **"Apa bedanya jika menggunakan Rate Limiter Bertingkat?"**

Berikut adalah penjelasan perbandingannya secara mendalam:

#### 1. Kombinasi Langsung (IP + Device + Email)
Dalam skenario ini, kita menggabungkan ketiga parameter tersebut menjadi **satu kunci tunggal** pada server Redis/Memory Rate Limiter.
* **Format Key:** `limiter:auth:ip_device_email:<ip_address>:<device_fingerprint>:<email>`
* **Cara Kerja:** Rate limiter hanya akan memblokir request jika **ketiga parameter tersebut sama persis**. Jika ada satu saja yang berubah, counter dianggap baru dari 0.

| Aspek | Penjelasan |
| :--- | :--- |
| **Kelemahan Utama** | **Sangat mudah di-bypass oleh attacker.** Attacker hanya perlu mengubah salah satu variabel saja (misalnya mengganti IP menggunakan VPN atau mengganti User-Agent browser) untuk mereset counter rate limit. |
| **Efek Brute-force** | Attacker dapat melakukan brute-force bertubi-tubi ke email korban dengan aman selama dia memutar IP atau berpindah device karena kombinasi kunci akan selalu terlihat baru bagi server. |

#### 2. Rate Limiter Bertingkat (Layered/Multi-tier Rate Limiter)
Alih-alih menyatukan semua parameter dalam satu key, kita menerapkan **beberapa aturan rate limiter secara terpisah (berlapis)** pada endpoint login.
* **Layer 1 (IP Limiter):** `limiter:auth:ip:<ip_address>` -> Limit 60 request/menit. *(Melindungi server dari spam DDoS global)*.
* **Layer 2 (Device/Browser Limiter):** `limiter:auth:device:<device_fingerprint>` -> Limit 20 request/menit. *(Membatasi spam dari satu browser)*.
* **Layer 3 (Email Target Limiter):** `limiter:auth:email:<email>` -> Limit 5 request/menit. *(Melindungi akun individu dari brute-force)*.

| Aspek | Penjelasan |
| :--- | :--- |
| **Mengapa Lebih Aman?** | **Attacker terkurung dari segala arah.** Jika attacker mengganti IP (untuk mem-bypass Layer 1), mereka tetap diblokir oleh Layer 3 (Email Target) karena limit email target sudah habis. Jika attacker mengganti email target (untuk menyerang akun lain), mereka akan segera diblokir oleh Layer 1 (IP) atau Layer 2 (Device) karena limit di level IP/device mereka telah terlampaui. |
| **Dampak ke Pengguna Sah** | Minimal. Jika terjadi brute-force massal pada satu email, hanya email tersebut yang dikunci (Layer 3) tanpa mengganggu IP kampus (Layer 1 masih memiliki kuota tinggi untuk IP bersama). |

> [!TIP]
> **Kesimpulan:** 
> * **Kombinasi Langsung** membuat ruang lingkup blokir menjadi terlalu spesifik, memudahkan penyerang untuk mem-bypass dengan mengubah salah satu elemen saja.
> * **Rate Limiter Bertingkat** memberikan perlindungan berlapis di setiap level (infrastruktur server, client device, dan akun pengguna), sehingga memberikan keamanan optimal tanpa mengorbankan kenyamanan pengguna sah.

---


## 📌 Ringkasan Eksekutif

> [!NOTE]
> * **`main.go`** menyiapkan pondasi awal.
> * **`config`** bertugas membaca pengaturan.
> * **`infrastructure`** membuat jembatan koneksi teknis.
> * **`app.NewApp()`** merakit seluruh logika sirkuit aplikasi.

* **Fail-Fast System:** Jika database gagal terkoneksi di awal, `NewApp()` tidak akan pernah berjalan.
* **Tightly Coupled Avoided:** Membuat database di luar `NewApp()` memberikan fleksibilitas penuh saat *testing* karena *resource* bisa di-swap dengan mudah.
* **Centralized State:** Pemisahan `config` memastikan komponen lain dapat menggunakan pengaturan aplikasi secara terpusat tanpa perlu memanggil `os.Getenv()` secara berulang-ulang di berbagai file (*anti-pattern*).
