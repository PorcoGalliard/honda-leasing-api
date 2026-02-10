-- 1. leasing_product 
INSERT INTO leasing.leasing_product (kode_produk, nama_produk, tenor_bulan, dp_persen_min, dp_persen_max, bunga_flat, admin_fee, asuransi) VALUES
('DP-RINGAN-24', 'DP Ringan 24 Bulan', 24, 10.00, 20.00, 1.20, 350000, TRUE),  -- DP min 10%, bunga flat 1.2%/bln
('SUPER-KILAT-12', 'Super Kilat 12 Bulan', 12, 25.00, 40.00, 1.00, 500000, TRUE),
('BEBAS-1TH-36', 'Bebas Angsuran 1 Tahun Pertama (36 bln)', 36, 15.00, 30.00, 1.15, 400000, TRUE),
('STANDAR-48', 'Standar Tenor Panjang 48 Bulan', 48, 20.00, 35.00, 1.30, 300000, TRUE),
('LOW-DP-36', 'Low DP Mulai 500rb (promo Beat)', 36, 5.00, 15.00, 1.25, 450000, TRUE);

-- 2. leasing_contract (contoh 3 kontrak aktif, link ke customer_id=1, motor_id dari sebelumnya)
-- Asumsi customer_id=1 (Kang Dian), motor_id=1 (Beat Hitam), product_id=1
INSERT INTO leasing.leasing_contract (
    contract_number, request_date, tanggal_akad, tanggal_mulai_cicil, tenor_bulan,
    nilai_kendaraan, dp_dibayar, pokok_pinjaman, total_pinjaman, cicilan_per_bulan, status,
    customer_id, motor_id, product_id
) VALUES
('KTR-2026-001', '2026-01-15', '2026-02-01', '2026-03-01', 24,
 18975000, 3000000, 15975000, 18700000, 779166, 'active', 1, 1, (SELECT product_id FROM leasing.leasing_product where kode_produk='DP-RINGAN-24')),  -- Beat Hitam, DP Ringan
('KTR-2026-002', '2026-02-05', '2026-02-20', '2026-03-20', 36,
 28050000, 5000000, 23050000, 28000000, 777777, 'approved', 1, 9, (SELECT product_id FROM leasing.leasing_product where kode_produk='BEBAS-1TH-36')),  -- Vario 160, Bebas 1 Th
('KTR-2026-003', '2026-02-08', NULL, NULL, 12,
 33400000, 8000000, 25400000, 28000000, 2333333, 'draft', 2, 12, (SELECT product_id FROM leasing.leasing_product where kode_produk='SUPER-KILAT-12'));  -- PCX, Super Kilat

-- 3. leasing_tasks (implementasi task untuk kontrak 1, link role dari account.roles)
-- Asumsi role_id SALES=3, SURVEYOR=4, ADMIN_CABANG=2, FINANCE=6
INSERT INTO leasing.leasing_tasks (task_name, startdate, enddate, sequence_no,status, contract_id, role_id) VALUES
('Input Pengajuan & Unggah Dokumen', '2026-01-15', '2026-01-20', 1,'completed', 
(SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001'), 
(SELECT role_id FROM account.roles WHERE role_name='SALES')),  -- SALES
('Auto Scoring Awal & Pre-Approval', '2026-01-25', '2026-01-28', 2,'completed', 
(SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001'), 
(SELECT role_id FROM account.roles WHERE role_name='ADMIN_CABANG')),-- ADMIN_CABANG
('Survei Lapangan / Home Visit', '2026-01-30', '2026-02-01', 3,'completed', 
(SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001'), 
(SELECT role_id FROM account.roles WHERE role_name='SURVEYOR')), -- SURVEYOR
('Input Hasil Survei & Rekomendasi', '2026-02-01', '2026-02-05', 4,'inprogress', 
(SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001'), 
(SELECT role_id FROM account.roles WHERE role_name='SURVEYOR')),-- SURVEYOR
('Review & Approval Final (ACC/Reject)', '2026-02-01', '2026-02-05', 5,'inprogress', 
(SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001'), 
(SELECT role_id FROM account.roles WHERE role_name='ADMIN_CABANG')),-- ADMIN_CABANG
('Akad & Tanda Tangan Kontrak', '2026-02-01', '2026-02-05', 6,'inprogress', 
(SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001'), 
(SELECT role_id FROM account.roles WHERE role_name='SALES')),--SALES
('Pembayaran DP + Biaya Awal', '2026-02-01', '2026-02-05', 7,'inprogress', 
(SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001'), 
(SELECT role_id FROM account.roles WHERE role_name='FINANCE')),
('Proses PO & Pembelian Unit ke Dealer', '2026-02-01', '2026-02-05', 8,'inprogress', 
(SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001'), 
(SELECT role_id FROM account.roles WHERE role_name='FINANCE')),
('Delivery Motor ke Rumah Customer', '2026-02-01', '2026-02-05', 9,'inprogress', 
(SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001'), 
(SELECT role_id FROM account.roles WHERE role_name='SALES')),
('Mulai Cicilan & Monitoring Pembayaran', '2026-02-01', '2026-02-05', 10,'inprogress', 
(SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001'), 
(SELECT role_id FROM account.roles WHERE role_name='COLLECTION')),
('System Closed', '2026-02-01', '2026-02-05', 11,'inprogress', 
(SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001'), 
(SELECT role_id FROM account.roles WHERE role_name='SYSTEM'));


-- 4. leasing_tasks_attributes (contoh atribut untuk task survei kontrak 1)
INSERT INTO leasing.leasing_tasks_attributes (tasa_name, tasa_value, tasa_status, tasa_task_id) VALUES
('Foto Rumah Depan', 'https://storage/.../rumah_depan.jpg', 'completed', 
    (select task_id FROM leasing.leasing_tasks where task_name='Input Pengajuan & Unggah Dokumen')),
('Catatan Wawancara', 'Penghasilan stabil, rumah sendiri', 'completed', (select task_id FROM leasing.leasing_tasks where task_name='Input Pengajuan & Unggah Dokumen')),
('Foto Selfie + KTP', 'https://storage/.../selfie_ktp.jpg', 'completed', (select task_id FROM leasing.leasing_tasks where task_name='Input Pengajuan & Unggah Dokumen'));

-- 5. payment_schedule (untuk kontrak 1, tenor 24 bln)
INSERT INTO finance.payment_schedule (angsuran_ke, jatuh_tempo, pokok, margin, 
total_tagihan, status_pembayaran, contract_id) VALUES
(1, '2026-03-01', 664062, 115104, 779166, 'unpaid', (SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001')),
(2, '2026-04-01', 664062, 115104, 779166, 'unpaid', (SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001')),
(24, '2028-02-01', 664062, 115104, 779166, 'unpaid', (SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001'));

-- 6. payments (contoh pembayaran DP untuk kontrak 1)
INSERT INTO finance.payments (nomor_bukti, jumlah_bayar, tanggal_bayar, metode_pembayaran, provider, contract_id, schedule_id) VALUES
('DP-KTR001-20260201', 3000000, '2026-02-01', 'transfer', 'BCA VA', (SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001'), NULL);

-- 7. leasing_contract_documents (dokumen kontrak 1)
INSERT INTO leasing.leasing_contract_documents (file_name, file_size, file_type, file_url, contract_id) VALUES
('KTP_KangDian.jpg', 450000, 'jpg', 'https://storage/.../ktp.jpg', (SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001')),
('Perjanjian_Akad.pdf', 1200000, 'pdf', 'https://storage/.../akad.pdf', (SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001')),
('Foto_Unit_Beat.jpg', 800000, 'jpg', 'https://storage/.../beat_unit.jpg', (SELECT contract_id from leasing.leasing_contract where contract_number='KTR-2026-001'));