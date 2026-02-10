-- 1. roles <<account>>
INSERT INTO account.roles (role_name, description) VALUES
('SUPER_ADMIN', 'Akses penuh sistem, konfigurasi, user management'),
('ADMIN_CABANG', 'Admin cabang: manage user lokal, approve kontrak, lihat report'),
('SALES', 'Sales/marketing: input pengajuan baru, simulasi cicilan, follow up customer'),
('SURVEYOR', 'Surveyor lapangan: input hasil survei, upload foto rumah/lokasi'),
('COLLECTION', 'Collection/tagih: monitor jadwal bayar, follow up tunggakan'),
('FINANCE', 'Finance/akunting: proses pembayaran DP & cicilan, rekonsiliasi'),
('CUSTOMER', 'Nasabah portal: lihat kontrak sendiri, jadwal cicilan, bukti bayar'),
('SYSTEM', 'Akses untuk trigger pgsql atau external service');

-- 2. permissions
INSERT INTO account.permissions (permission_type, description) VALUES
('view_dashboard', 'Lihat dashboard utama & ringkasan data'),
('view_contract', 'Lihat daftar & detail kontrak leasing'),
('create_contract', 'Buat pengajuan kontrak baru'),
('approve_contract', 'Approve / tolak kontrak setelah survey'),
('view_survey', 'Lihat hasil survei & dokumen pendukung'),
('create_survey', 'Input hasil survei lapangan & upload foto'),
('view_payment', 'Lihat jadwal angsuran & riwayat pembayaran'),
('record_payment', 'Catat pembayaran manual / konfirmasi transfer'),
('manage_user', 'Kelola user, role, dan permission'),
('manage_oauth', 'Setup & edit OAuth provider (Google, Apple, dll)'),
('export_report', 'Download report excel/pdf (kontrak, pembayaran, dll)'),
('send_notif', 'Send notification to users');

-- 3. role_permission (mapping role ke permission)
-- gunakan metode select insert into dengan join with value 'SUPER_ADMIN', karena role_id bisa
--jadi digenerate tidak mulai dari 1
INSERT INTO account.role_permission (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM account.roles r
CROSS JOIN account.permissions p
WHERE r.role_name = 'SUPER_ADMIN';

-- ADMIN_CABANG: sebagian kecuali manage_oauth
INSERT INTO account.role_permission (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM account.roles r
CROSS JOIN account.permissions p
WHERE r.role_name = 'ADMIN_CABANG'
  AND p.permission_type != 'manage_oauth';

-- SALES
INSERT INTO account.role_permission (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM account.roles r, account.permissions p
WHERE r.role_name = 'SALES'
  AND p.permission_type IN ('view_dashboard', 'view_contract', 'create_contract', 'view_payment', 'export_report');

-- SURVEYOR
INSERT INTO account.role_permission (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM account.roles r, account.permissions p
WHERE r.role_name = 'SURVEYOR'
  AND p.permission_type IN ('view_dashboard', 'view_contract', 'view_survey', 'create_survey');

-- COLLECTION
INSERT INTO account.role_permission (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM account.roles r, account.permissions p
WHERE r.role_name = 'COLLECTION'
  AND p.permission_type IN ('view_dashboard', 'view_contract', 'view_payment', 'record_payment');

-- FINANCE
INSERT INTO account.role_permission (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM account.roles r, account.permissions p
WHERE r.role_name = 'FINANCE'
  AND p.permission_type IN ('view_dashboard', 'view_contract', 'view_payment', 'record_payment', 'export_report');

-- CUSTOMER (customer – sangat terbatas)
INSERT INTO account.role_permission (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM account.roles r, account.permissions p
WHERE r.role_name = 'CUSTOMER'
  AND p.permission_type IN ('view_dashboard', 'view_contract', 'view_payment');

  -- CUSTOMER (customer – sangat terbatas)
INSERT INTO account.role_permission (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM account.roles r, account.permissions p
WHERE r.role_name = 'SYSTEM'
  AND p.permission_type IN ('send_notif', 'export_report');


-- 4. users <<account>>
-- Contoh user internal cabang Bandung + 1 nasabah demo
-- PIN contoh: 123456, 654321, dll → hash dulu di app!
INSERT INTO account.users (
    phone_number, email, full_name, password, pin_key, is_active, created_at
) VALUES
('+6281212345678', 'superadmin@leasingbdg.id', 'Andi - Super Admin', NULL, '$argon2id$v=19$m=65536,t=3,p=4$examplehashfor123456', TRUE, CURRENT_TIMESTAMP),
('+6281556789123', 'adminbdg@leasingbdg.id', 'Budi Admin - Admin Cabang', NULL, '$argon2id$v=19$m=65536,t=3,p=4$examplehashfor654321', TRUE, CURRENT_TIMESTAMP),
('+6281788990011', 'sales1@leasingbdg.id', 'Ani Sales - Marketing', NULL, '$argon2id$v=19$m=65536,t=3,p=4$examplehashfor112233', TRUE, CURRENT_TIMESTAMP),
('+6282198765432', 'surveyor1@leasingbdg.id', 'Dedi Surveyor - Lapangan', NULL, '$argon2id$v=19$m=65536,t=3,p=4$examplehashfor445566', TRUE, CURRENT_TIMESTAMP),
('+6285712345678', 'collection@leasingbdg.id', 'Asep Knalpot Collection - Tagih', NULL, '$argon2id$v=19$m=65536,t=3,p=4$examplehashfor778899', TRUE, CURRENT_TIMESTAMP),
('+6289612345678', 'finance@leasingbdg.id', 'Rina Finance - Akunting', NULL, '$argon2id$v=19$m=65536,t=3,p=4$examplehashfor990011', TRUE, CURRENT_TIMESTAMP),
('+6281314151617', 'kangdian@gmail.com', 'Kang Dian - Customer', NULL, '$argon2id$v=19$m=65536,t=3,p=4$examplehashfor121212', TRUE, CURRENT_TIMESTAMP),
('+621', 'system@gmail.com', 'System', NULL, NULL, TRUE, CURRENT_TIMESTAMP);

-- 5. user_roles <<account>>
-- Assign role (asumsi user_id mulai dari 1 sesuai insert di atas)
INSERT INTO account.user_roles (user_id, role_id, assigned_by) VALUES
(1, (SELECT role_id FROM account.roles WHERE role_name = 'SUPER_ADMIN'), 1),
(2, (SELECT role_id FROM account.roles WHERE role_name = 'ADMIN_CABANG'), 1),
(2, (SELECT role_id FROM account.roles WHERE role_name = 'SUPER_ADMIN'), 1),  -- admin cabang juga super
(3, (SELECT role_id FROM account.roles WHERE role_name = 'SALES'), 2),
(4, (SELECT role_id FROM account.roles WHERE role_name = 'SURVEYOR'), 2),
(5, (SELECT role_id FROM account.roles WHERE role_name = 'COLLECTION'), 2),
(6, (SELECT role_id FROM account.roles WHERE role_name = 'FINANCE'), 2),
(7, (SELECT role_id FROM account.roles WHERE role_name = 'CUSTOMER'), 2),
(8, (SELECT role_id FROM account.roles WHERE role_name = 'SYSTEM'), 1);

-- 6. oauth_providers <<account>>
-- Contoh setup untuk login Google & Apple (ganti client_id/secret dengan milikmu)
INSERT INTO account.oauth_providers (
    provider_name, client_id, client_secret, redirect_uri, issuer_url, active
) VALUES
('google', '123456789012-abcde.apps.googleusercontent.com', 'GOCSPX-your-secret-here-very-long', 'https://api.leasingbdg.id/auth/google/callback', 'https://accounts.google.com', TRUE),
('apple', 'com.leasingbdg.signinwithapple.service', 'your-apple-private-key-or-secret', 'https://api.leasingbdg.id/auth/apple/callback', 'https://appleid.apple.com', TRUE);

