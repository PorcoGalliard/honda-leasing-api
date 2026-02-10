-- =============================================
-- Seed Data (INSERT raw data)
-- =============================================

-- 1. province: hanya Jawa Barat
INSERT INTO mst.province (prov_name) VALUES ('Jawa Barat');

-- 2. kabupaten: semua kab/kota di Jawa Barat 
INSERT INTO mst.kabupaten (kab_name, prov_id) VALUES
('Kabupaten Bandung', 1),
('Kabupaten Bandung Barat', 1),
('Kabupaten Bekasi', 1),
('Kabupaten Bogor', 1),
('Kabupaten Ciamis', 1),
('Kabupaten Cianjur', 1),
('Kabupaten Cirebon', 1),
('Kabupaten Garut', 1),
('Kabupaten Indramayu', 1),
('Kabupaten Karawang', 1),
('Kabupaten Kuningan', 1),
('Kabupaten Majalengka', 1),
('Kabupaten Pangandaran', 1),
('Kabupaten Purwakarta', 1),
('Kabupaten Subang', 1),
('Kabupaten Sukabumi', 1),
('Kabupaten Sumedang', 1),
('Kabupaten Tasikmalaya', 1),
('Kota Bandung', 1),
('Kota Bekasi', 1),
('Kota Bogor', 1),
('Kota Cimahi', 1),
('Kota Cirebon', 1),
('Kota Depok', 1),
('Kota Sukabumi', 1),
('Kota Tasikmalaya', 1),
('Kota Banjar', 1);

-- 3. kecamatan: contoh 5 kecamatan di Kota Bandung (kab_id asumsi 19 = Kota Bandung)
INSERT INTO mst.kecamatan (kec_name, kab_id) VALUES
('Andir', (SELECT kab_id FROM mst.kabupaten WHERE kab_name = 'Kota Bandung')),
('Antapani', (SELECT kab_id FROM mst.kabupaten WHERE kab_name = 'Kota Bandung')),
('Arcamanik', (SELECT kab_id FROM mst.kabupaten WHERE kab_name = 'Kota Bandung')),
('Astanaanyar', (SELECT kab_id FROM mst.kabupaten WHERE kab_name = 'Kota Bandung')),
('Bandung Kulon', (SELECT kab_id FROM mst.kabupaten WHERE kab_name = 'Kota Bandung'));

-- 4. kelurahan: contoh kelurahan di Kecamatan Bandung Kulon 
INSERT INTO mst.kelurahan (kel_name, kec_id) VALUES
('Caringin', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon')),
('Cibuntu', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon')),
('Cijerah', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon')),
('Gempolsari', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon')),
('Cigondewah Kaler', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon')),
('Cigondewah Kidul', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon')),
('Warung Muncang', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon')),
('Cigondewah Rahayu', (SELECT kec_id FROM mst.kecamatan WHERE kec_name = 'Bandung Kulon'));

-- 5. locations: alamat dummy di Bandung
INSERT INTO mst.locations (street_address, postal_code, longitude, latitude, kel_id) VALUES
('Jl. Soekarno Hatta No. 123, Caringin', '40212', '107.612345', '-6.912345', (SELECT kel_id FROM mst.kelurahan WHERE kel_name = 'Caringin')),
('Jl. Terusan Buahbatu No. 45, Cijerah', '40213', '107.598765', '-6.923456', (SELECT kel_id FROM mst.kelurahan WHERE kel_name = 'Cijerah')),
('Komplek Griya Bandung Indah, Warung Muncang', '40211', '107.585678', '-6.935678', (SELECT kel_id FROM mst.kelurahan WHERE kel_name = 'Warung Muncang'));