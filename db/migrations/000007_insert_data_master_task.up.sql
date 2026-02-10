-- =============================================
-- Seed Data
-- =============================================

-- Relasikan role_id dari account.roles
-- Misal: SALES=3, SURVEYOR=4, ADMIN_CABANG=2, FINANCE=6, COLLECTION=5, CUSTOMER=7

-- 1. template_tasks (step alur leasing motor)
INSERT INTO mst.template_tasks (teta_name, teta_role_id, description, sequence_no, is_required) VALUES
('Input Pengajuan & Unggah Dokumen', 
 (SELECT role_id FROM account.roles WHERE role_name = 'SALES'), 
 'Isi form pengajuan, unggah KTP, KK, slip gaji/foto usaha', 
 1, TRUE),

('Auto Scoring Awal & Pre-Approval', 
 (SELECT role_id FROM account.roles WHERE role_name = 'ADMIN_CABANG'), 
 'Sistem cek SLIK OJK, DSR, duplikasi data', 
 2, TRUE),

('Survei Lapangan / Home Visit', 
 (SELECT role_id FROM account.roles WHERE role_name = 'SURVEYOR'), 
 'Kunjungi alamat, verifikasi, foto rumah, wawancara', 
 3, TRUE),

('Input Hasil Survei & Rekomendasi', 
 (SELECT role_id FROM account.roles WHERE role_name = 'SURVEYOR'), 
 'Upload hasil survei & dokumen pendukung ke sistem', 
 4, TRUE),

('Review & Approval Final (ACC/Reject)', 
 (SELECT role_id FROM account.roles WHERE role_name = 'ADMIN_CABANG'), 
 'Analisis hasil survei + scoring, beri keputusan ACC', 
 5, TRUE),

('Akad & Tanda Tangan Kontrak', 
 (SELECT role_id FROM account.roles WHERE role_name = 'SALES'), 
 'Customer tanda tangan perjanjian, polis asuransi di dealer/cabang', 
 6, TRUE),

('Pembayaran DP + Biaya Awal', 
 (SELECT role_id FROM account.roles WHERE role_name = 'FINANCE'), 
 'Verifikasi pembayaran DP, admin, asuransi, fidusia', 
 7, TRUE),

('Proses PO & Pembelian Unit ke Dealer', 
 (SELECT role_id FROM account.roles WHERE role_name = 'FINANCE'), 
 'Leasing bayar ke dealer, terbitkan PO', 
 8, TRUE),

('Delivery Motor ke Rumah Customer', 
 (SELECT role_id FROM account.roles WHERE role_name = 'SALES'), 
 'Unit dikirim, serah terima + dokumen (STNK sementara)', 
 9, TRUE),
('Mulai Cicilan & Monitoring Pembayaran', 
 (SELECT role_id FROM account.roles WHERE role_name = 'COLLECTION'), 
 'Generate schedule angsuran, follow up pembayaran bulanan', 
 10, TRUE),

('System Closed', 
 (SELECT role_id FROM account.roles WHERE role_name = 'SYSTEM'), 
 'System automatically closed', 
 11, FALSE);  

-- 2. template_task_attributes (contoh atribut/dokumen per tahapan)
-- ini contoh saja bisa ditambahkan 
INSERT INTO mst.template_task_attributes (tetat_name, tetat_teta_id, description, is_required, attribute_type) VALUES
-- Tahap 1: Input Pengajuan
('Upload KTP', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Input Pengajuan & Unggah Dokumen'), 'Foto KTP asli', TRUE, 'file'),
('Upload KK', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Input Pengajuan & Unggah Dokumen'), 'Kartu Keluarga', TRUE, 'file'),
('Slip Gaji / Bukti Penghasilan', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Input Pengajuan & Unggah Dokumen'), '3 bulan terakhir atau surat usaha', TRUE, 'file'),

-- Tahap 3: Survei Lapangan
('Foto Rumah Depan & Sekitar', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Survei Lapangan / Home Visit'), 'Minimal 3-5 foto', TRUE, 'file'),
('Foto Selfie + KTP di Lokasi', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Survei Lapangan / Home Visit'), 'Verifikasi identitas', TRUE, 'file'),
('Catatan Wawancara', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Survei Lapangan / Home Visit'), 'Pekerjaan, penghasilan, kondisi rumah', TRUE, 'text'),

-- Tahap 7: Akad
('Tanda Tangan Digital / Fisik', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Akad & Tanda Tangan Kontrak'), 'Perjanjian pembiayaan', TRUE, 'file'),
('Foto Serah Terima Unit', (SELECT teta_id FROM mst.template_tasks WHERE teta_name = 'Delivery Motor ke Rumah Customer'), 'Customer terima motor + kunci', TRUE, 'file');
-- Tahap 10: Delivery