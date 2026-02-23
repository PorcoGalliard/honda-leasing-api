# Honda Leasing API — Dokumentasi

> **Base URL:** `http://localhost:8080/leasing/api`

---

## Daftar Isi

1. [Autentikasi](#1-autentikasi)
2. [Motor](#2-motor)
3. [Order](#3-order)
4. [Staff](#4-staff)
5. [Status & Transisi](#5-status--transisi)
6. [Kode Error](#6-kode-error)

---

## 1. Autentikasi

### `POST /auth/register`
Mendaftarkan akun baru.

**Description:**
- Field `role` opsional. Jika tidak diisi, default ke `CUSTOMER`.
- Role yang tersedia saat register: `CUSTOMER`, `SALES`, `SURVEYOR`, `FINANCE`, `COLLECTION`, `ADMIN_CABANG`.
- Nomor HP tidak boleh duplikat.
- `password` minimal 8 karakter.
- `pin` harus tepat 6 digit angka.

---

### `POST /auth/login`
Login menggunakan nomor HP dan PIN. Mengembalikan `access_token` (berlaku 24 jam) dan `refresh_token` (berlaku 7 hari).

---

## 2. Motor

### `GET /motors`
Menampilkan daftar motor yang tersedia. **Endpoint publik, tidak butuh token.**

**Query Parameters:**

`page` | int | Nomor halaman (default: 1) |
`limit` | int | Jumlah data per halaman (default: 10) |
`merk` | string | Filter berdasarkan merek |
`motor_type` | string | Filter berdasarkan tipe motor |
`status` | string | Filter status: `ready`, `booked`, `sold` |

---

### `POST /motors/credit-simulation`
Simulasi cicilan kredit tanpa perlu login. **Endpoint publik, tidak butuh token.**

**Description:**
- Pilihan tenor: `23`, `29`, atau `35` bulan.
- Bunga: **0%**. Saya tetapkan 0% saja
- Biaya yang sudah ditetapkan (fixed):
  - Biaya Admin: Rp 200.000
  - Asuransi: Rp 250.000
  - Fidusia: Rp 200.000
  - Materai: Rp 10.000

---

## 3. Order

> Semua endpoint order membutuhkan `Authorization: Bearer <access_token>`.

### `POST /orders`
Membuat pengajuan leasing baru.

**Rules:**
- Hanya bisa diakses oleh user yang sudah login (customer).
- Motor harus berstatus `ready`. Jika sudah `booked` atau `sold`, pengajuan ditolak.
- Field `nik` (16 digit) wajib diisi jika customer belum memiliki profil di sistem.
- Tenor valid: `23`, `29`, atau `35` bulan.
- `request_date` format: `YYYY-MM-DD`.
- Setelah order berhasil dibuat:
  - Status motor otomatis berubah menjadi `booked`.
  - Nomor kontrak di-generate otomatis dengan format `KTR-{TAHUN}-{URUTAN}` (contoh: `KTR-2026-0001`).
  - Tasks alur leasing otomatis di-copy dari template sebanyak 11 task.
  - Order berstatus awal `draft`.

---

### `GET /orders/:contract_id/progress`
Melihat detail dan progress task dari sebuah order.

**Rules:**
- Customer hanya bisa melihat order miliknya sendiri. Jika `contract_id` milik orang lain, akan dikembalikan `404`.

---

## 4. Staff

> Semua endpoint staff membutuhkan `Authorization: Bearer <token>` dengan role **selain** `CUSTOMER`.

### `GET /staff/orders`
Melihat daftar semua order. Hasil difilter otomatis berdasarkan role yang login.

**Rules & Akses per Role:**

`SUPER_ADMIN` | Semua order |
`ADMIN_CABANG` | Semua order |
`SALES` | Order yang memiliki task untuk role SALES |
`SURVEYOR` | Order yang memiliki task untuk role SURVEYOR |
`FINANCE` | Order yang memiliki task untuk role FINANCE |
`COLLECTION` | Order yang memiliki task untuk role COLLECTION |

**Query Parameters:**

`page` | int | Nomor halaman (default: 1, min: 1) |
`limit` | int | Jumlah data per halaman (default: 10, max: 100) |
`status` | string | Filter status order (lihat [Status Order](#status-order)) |

---

### `PATCH /staff/orders/:contract_id/status`
Mengubah status sebuah kontrak/order.

**Rules:**
- Transisi status harus mengikuti alur yang valid (lihat [Aturan Transisi](#aturan-transisi-status)).
- Setiap role hanya boleh mengubah ke status tertentu:

`SUPER_ADMIN` | Semua status |
`ADMIN_CABANG` | Semua status |
`FINANCE` | `active`, `paid_off` |
`COLLECTION` | `late` |

---

### `PATCH /staff/orders/:contract_id/tasks/:task_id`
Mengubah status sebuah task dalam order.

**Rules:**
- Setiap task memiliki `role_id` yang menentukan siapa yang bertanggung jawab.
- Staff hanya bisa mengubah task yang di-assign ke role-nya sendiri.
- `SUPER_ADMIN` dapat mengubah task milik role manapun.
- Status yang diizinkan: `completed`, `cancelled`.
- Jika status diubah ke `completed`, `actual_enddate` otomatis diisi dengan waktu saat ini.

**Pembagian task per role (dari template):**

1. Input Pengajuan & Unggah Dokumen | `SALES` |
2. Auto Scoring Awal & Pre-Approval | `ADMIN_CABANG` |
3. Survei Lapangan / Home Visit | `SURVEYOR` |
4. Input Hasil Survei & Rekomendasi | `SURVEYOR` |
5. Review & Approval Final (ACC/Reject) | `ADMIN_CABANG` |
6. Akad & Tanda Tangan Kontrak | `SALES` |
7. Pembayaran DP + Biaya Awal | `FINANCE` |
8. Proses PO & Pembelian Unit ke Dealer | `FINANCE` |
9. Delivery Motor ke Rumah Customer | `SALES` |
10. Mulai Cicilan & Monitoring Pembayaran | `COLLECTION` |
11. System Closed | `SYSTEM` |

