-- Schema: dealer

-- drop table 
DROP TABLE IF EXISTS dealer.customers CASCADE;
DROP TABLE IF EXISTS dealer.motor_types CASCADE;
DROP TABLE IF EXISTS dealer.motors CASCADE;
DROP TABLE IF EXISTS dealer.motor_assets CASCADE;

-- 1. customers kita masukan di schema dealer 
CREATE TABLE dealer.customers (
    customer_id      BIGSERIAL PRIMARY KEY,
    nik            VARCHAR(16) UNIQUE NOT NULL,
    nama_lengkap   VARCHAR(100) NOT NULL,
    tanggal_lahir  DATE,
    no_hp          VARCHAR(15) UNIQUE NOT NULL,
    email          VARCHAR(100) UNIQUE,
    pekerjaan      VARCHAR(80),
    perusahaan     VARCHAR(120),
    salary         NUMERIC(15,2),
    location_id    BIGINT REFERENCES mst.locations(location_id) ON DELETE SET NULL,
    created_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 2. motor_types <<dealer>> (kategori besar: Classic, Sport, Matic, dll)
CREATE TABLE dealer.motor_types (
    moty_id    BIGSERIAL PRIMARY KEY,
    moty_name  VARCHAR(55) UNIQUE NOT NULL  -- Classic, Sport, Matic, Bebek, Maxi, dll
);

-- 3. motors <<dealer>> (stok unit motor, fokus Honda saja)
-- tuk nambahin value di motor_type -> ALTER TYPE motor_type_enum ADD VALUE 'Bebek';
CREATE TABLE dealer.motors (
    motor_id       BIGSERIAL PRIMARY KEY,
    merk           VARCHAR(50) DEFAULT 'Honda' NOT NULL,
    motor_type     VARCHAR(15) CHECK (motor_type IN ('Classic', 'Sport', 'Matic','Maxi','Bebek')),  
    tahun          SMALLINT NOT NULL,
    warna          VARCHAR(30),
    nomor_rangka   VARCHAR(30) UNIQUE NOT NULL,
    nomor_mesin    VARCHAR(30) UNIQUE NOT NULL,
    cc_mesin       VARCHAR(30),  -- misal '109.5 cc', '156.9 cc'
    nomor_polisi   VARCHAR(12) UNIQUE,
    status_unit    VARCHAR(20) DEFAULT 'ready' 
                   CHECK (status_unit IN ('ready', 'booked', 'leased', 'returned', 'repo')),
    harga_otr      NUMERIC(15,2) NOT NULL,  -- OTR Bandung 2026
    motor_moty_id  BIGINT NOT NULL REFERENCES dealer.motor_types(moty_id) ON DELETE RESTRICT,
    created_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 4. motor_assets <<dealer>> (foto, dokumen unit motor)
CREATE TABLE dealer.motor_assets (
    moas_id        BIGSERIAL PRIMARY KEY,
    file_name      VARCHAR(125) NOT NULL,
    file_size      DOUBLE PRECISION,  -- dalam bytes
    file_type      VARCHAR(15) CHECK (file_type IN ('png','jpg','jpeg','pdf','doc','docx')), 
    file_url       VARCHAR(255) NOT NULL,  -- path atau URL ke S3 atau local
    moas_motor_id  BIGINT NOT NULL REFERENCES dealer.motors(motor_id) ON DELETE CASCADE,
    created_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Index untuk performa query & searching
CREATE INDEX idx_motors_merk_type    ON dealer.motors(merk, motor_type);
CREATE INDEX idx_motors_status       ON dealer.motors(status_unit);
CREATE INDEX idx_motors_moty_id      ON dealer.motors(motor_moty_id);
CREATE INDEX idx_motor_assets_motor  ON dealer.motor_assets(moas_motor_id);
CREATE INDEX idx_customers_location  ON dealer.customers(location_id);