-- Schema: mst 

DROP TABLE IF EXISTS mst.province CASCADE;
DROP TABLE IF EXISTS mst.kabupaten CASCADE;
DROP TABLE IF EXISTS mst.kecamatan CASCADE;
DROP TABLE IF EXISTS mst.kelurahan CASCADE;
DROP TABLE IF EXISTS mst.locations CASCADE;


-- Fokus Provinsi: Jawa Barat saja
-- 1. province <<mst>>
CREATE TABLE mst.province (
    prov_id    BIGSERIAL PRIMARY KEY,
    prov_name  VARCHAR(85) UNIQUE NOT NULL
);

-- 2. kabupaten <<mst>>
CREATE TABLE mst.kabupaten (
    kab_id     BIGSERIAL PRIMARY KEY,
    kab_name   VARCHAR(85) NOT NULL,
    prov_id    BIGINT NOT NULL REFERENCES mst.province(prov_id) ON DELETE RESTRICT
);

-- 3. kecamatan <<mst>>
CREATE TABLE mst.kecamatan (
    kec_id     BIGSERIAL PRIMARY KEY,
    kec_name   VARCHAR(85) NOT NULL,
    kab_id     BIGINT NOT NULL REFERENCES mst.kabupaten(kab_id) ON DELETE RESTRICT
);

-- 4. kelurahan <<mst>>
CREATE TABLE mst.kelurahan (
    kel_id     BIGSERIAL PRIMARY KEY,
    kel_name   VARCHAR(85) NOT NULL,
    kec_id     BIGINT NOT NULL REFERENCES mst.kecamatan(kec_id) ON DELETE RESTRICT
);

-- 5. locations <<mst>>
CREATE TABLE mst.locations (
    location_id    BIGSERIAL PRIMARY KEY,
    street_address TEXT,
    postal_code    VARCHAR(10),
    longitude      numeric(9,6),
    latitude       numeric(9,6),
    kel_id         BIGINT REFERENCES mst.kelurahan(kel_id) ON DELETE SET NULL
);

-- =============================================
-- INDEX untuk performa query 
-- =============================================
CREATE INDEX idx_kabupaten_prov ON mst.kabupaten(prov_id);
CREATE INDEX idx_kecamatan_kab  ON mst.kecamatan(kab_id);
CREATE INDEX idx_kelurahan_kec  ON mst.kelurahan(kec_id);
CREATE INDEX idx_locations_kel  ON mst.locations(kel_id);

