-- 1. motor_types (kategori utama Honda)
INSERT INTO dealer.motor_types (moty_name) VALUES
('Matic'),    -- skutik harian
('Sport'),    -- CBR, Sonic
('Classic'),  -- Stylo, retro style
('Bebek'),    -- Revo, Supra (jika perlu)
('Maxi');     -- PCX, ADV, Forza

-- 2. motors: contoh lengkap motor Honda saja 
INSERT INTO dealer.motors (
    merk, motor_type, tahun, warna, nomor_rangka, nomor_mesin, cc_mesin, nomor_polisi, 
    status_unit, harga_otr, motor_moty_id
) VALUES
-- Beat 
('Honda', 'Matic', 2026, 'Hitam', 'MH1JF41-001', 'JF41E-001', '109.5 cc', 'D 1234 ABC', 'ready', 18975000, (SELECT moty_id FROM dealer.motor_types WHERE moty_name = 'Matic')),
('Honda', 'Matic', 2026, 'Merah', 'MH1JF41-002', 'JF41E-002', '109.5 cc', 'D 5678 DEF', 'ready', 19775000, (SELECT moty_id FROM dealer.motor_types WHERE moty_name = 'Matic')),  -- Deluxe
('Honda', 'Matic', 2026, 'Putih', 'MH1JF41-003', 'JF41E-003', '109.5 cc', NULL, 'booked', 20575000, (SELECT moty_id FROM dealer.motor_types WHERE moty_name = 'Matic')),  -- Deluxe Smart Key
('Honda', 'Matic', 2026, 'Abu-abu', 'MH1JF41-004', 'JF41E-004', '109.5 cc', 'D 9012 GHI', 'ready', 20075000, (SELECT moty_id FROM dealer.motor_types WHERE moty_name = 'Matic')),  -- Beat Street

-- Genio & Scoopy
('Honda', 'Matic', 2026, 'Biru', 'MH1JF50-001', 'JF50E-001', '110 cc', 'D 3456 JKL', 'ready', 20175000, (SELECT moty_id FROM dealer.motor_types WHERE moty_name = 'Matic')),  -- Genio CBS
('Honda', 'Matic', 2026, 'Hitam Metalik', 'MH1JF50-002', 'JF50E-002', '110 cc', NULL, 'ready', 23000000, (SELECT moty_id FROM dealer.motor_types WHERE moty_name = 'Matic')),  -- Scoopy Fashion

-- Vario 
('Honda', 'Matic', 2026, 'Merah Doff', 'MH1JF60-001', 'JF60E-001', '124.8 cc', 'D 7890 MNO', 'ready', 23775000, (SELECT moty_id FROM dealer.motor_types WHERE moty_name = 'Matic')),  -- Vario 125 CBS
('Honda', 'Matic', 2026, 'Hitam', 'MH1JF60-002', 'JF60E-002', '124.8 cc', 'D 1122 PQR', 'leased', 25600000, (SELECT moty_id FROM dealer.motor_types WHERE moty_name = 'Matic')),  -- Vario 125 CBS-ISS
('Honda', 'Matic', 2026, 'Putih', 'MH1JF70-001', 'JF70E-001', '156.9 cc', NULL, 'ready', 28050000, (SELECT moty_id FROM dealer.motor_types WHERE moty_name = 'Matic')),  -- Vario 160 CBS
('Honda', 'Matic', 2026, 'Abu-abu', 'MH1JF70-002', 'JF70E-002', '156.9 cc', 'D 3344 STU', 'ready', 31075000, (SELECT moty_id FROM dealer.motor_types WHERE moty_name = 'Matic')),  -- Vario 160 ABS

-- Stylo, ADV, PCX (mid-premium)
('Honda', 'Classic', 2026, 'Krem', 'MH1JF80-001', 'JF80E-001', '156.9 cc', NULL, 'ready', 28900000, (SELECT moty_id FROM dealer.motor_types WHERE moty_name = 'Classic')),  -- Stylo 160 CBS
('Honda', 'Matic', 2026, 'Hitam', 'MH1JF90-001', 'JF90E-001', '156.9 cc', 'D 5566 VWX', 'ready', 37459000, (SELECT moty_id FROM dealer.motor_types WHERE moty_name = 'Matic')),  -- ADV 160 CBS
('Honda', 'Maxi', 2026, 'Hitam Metalik', 'MH1JF100-001', 'JF100E-001', '156.9 cc', NULL, 'ready', 33400000, (SELECT moty_id FROM dealer.motor_types WHERE moty_name = 'Maxi'));  -- PCX 160 CBS

-- 3. motor_assets: contoh aset (foto unit)
INSERT INTO dealer.motor_assets (file_name, file_size, file_type, file_url, moas_motor_id) VALUES
('beat_cbs_hitam_2026.jpg', 1250000, 'jpg', 'https://storage.dealer.id/assets/beat_hitam.jpg', 1),
('vario160_abs_putih_side.jpg', 1800000, 'jpg', 'https://storage.dealer.id/assets/vario160_abs.jpg', 10),
('stylo_krem_retro_front.jpg', 2100000, 'jpg', 'https://storage.dealer.id/assets/stylo_krem.jpg', 11),
('pcx160_hitam_3d_view.pdf', 850000, 'pdf', 'https://storage.dealer.id/docs/pcx160_brosur.pdf', 12),
('adv160_hitam_action.jpg', 1500000, 'jpg', 'https://storage.dealer.id/assets/adv160.jpg', 12);

-- 4. customers: contoh 2 customes di Bandung
INSERT INTO dealer.customers (
    nik, nama_lengkap, tanggal_lahir, no_hp, email, pekerjaan, perusahaan, salary, location_id
) VALUES
('3273010101010001', 'Kang Dian', '1990-05-15', '0221234567', 'kangdian@gmail.co.id', 'Programmer', 'PT. Codeid', 150000000, (SELECT location_id FROM mst.locations LIMIT 1)),
('3273010202020002', 'Winona', '1985-08-20', '081234567890', 'winona@gmail.co.id', 'Tiktok influecer', '-', 120000000, (SELECT location_id FROM mst.locations WHERE street_address LIKE '%Soekarno Hatta%' LIMIT 1));