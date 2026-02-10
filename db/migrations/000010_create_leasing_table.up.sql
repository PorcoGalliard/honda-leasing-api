-- Schema: leasing (produk, kontrak, task implementasi, dokumen)
-- Schema: payment (jadwal & pembayaran)

DROP TABLE IF EXISTS leasing.leasing_product CASCADE;
DROP TABLE IF EXISTS leasing.leasing_contract CASCADE;
DROP TABLE IF EXISTS leasing.leasing_tasks CASCADE;
DROP TABLE IF EXISTS leasing.leasing_tasks_attributes CASCADE;
DROP TABLE IF EXISTS leasing.leasing_contract_documents CASCADE;
DROP TABLE IF EXISTS finance.payments CASCADE;
DROP TABLE IF EXISTS finance.payment_schedule CASCADE;


-- 1. leasing_product <<leasing>>
CREATE TABLE leasing.leasing_product (
    product_id     BIGSERIAL PRIMARY KEY,
    kode_produk    VARCHAR(20) UNIQUE NOT NULL,
    nama_produk    VARCHAR(100) NOT NULL,               -- "Bebas 1 Tahun", "DP Ringan", dll
    tenor_bulan    SMALLINT NOT NULL CHECK (tenor_bulan BETWEEN 12 AND 60),
    dp_persen_min  NUMERIC(5,2) NOT NULL,
    dp_persen_max  NUMERIC(5,2) NOT NULL,
    bunga_flat     NUMERIC(5,2) NOT NULL,                -- % flat per tahun
    admin_fee      NUMERIC(12,2) DEFAULT 0,
    asuransi       BOOLEAN DEFAULT TRUE,
    created_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 2. leasing_contract <<leasing>>
CREATE TABLE leasing.leasing_contract (
    contract_id         BIGSERIAL PRIMARY KEY,
    contract_number     VARCHAR(30) UNIQUE,
    request_date        DATE NOT NULL DEFAULT CURRENT_DATE,
    tanggal_akad        DATE,
    tanggal_mulai_cicil DATE,
    tenor_bulan         SMALLINT NOT NULL,
    nilai_kendaraan     NUMERIC(15,2) NOT NULL,
    dp_dibayar          NUMERIC(15,2) NOT NULL,
    pokok_pinjaman      NUMERIC(15,2) NOT NULL,
    total_pinjaman      NUMERIC(15,2) NOT NULL,         -- pokok + margin/bunga
    cicilan_per_bulan   NUMERIC(15,2) NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'draft'
                        CHECK (status IN ('draft', 'approved', 'active', 'late', 'paid_off', 'repo', 'canceled')),
    created_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    customer_id         BIGINT NOT NULL REFERENCES dealer.customers(customer_id) ON DELETE RESTRICT,  -- nasabah
    motor_id            BIGINT NOT NULL REFERENCES dealer.motors(motor_id) ON DELETE RESTRICT,
    product_id          BIGINT NOT NULL REFERENCES leasing.leasing_product(product_id) ON DELETE RESTRICT
);

-- 3. leasing_tasks <<leasing>> (implementasi task dari template)
CREATE TABLE leasing.leasing_tasks (
    task_id            BIGSERIAL PRIMARY KEY,
    task_name          VARCHAR(85) NOT NULL,
    startdate          DATE,
    enddate            DATE,
    actual_startdate   DATE,
    actual_enddate     DATE,
    status             VARCHAR(15) DEFAULT 'inprogress'
                       CHECK (status IN ('inprogress', 'completed', 'cancelled')),
    contract_id        BIGINT NOT NULL REFERENCES leasing.leasing_contract(contract_id) ON DELETE CASCADE,
    role_id            BIGINT NOT NULL REFERENCES account.roles(role_id) ON DELETE RESTRICT,
    created_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    sequence_no SMALLINT DEFAULT 0  ,
    call_function TEXT DEFAULT NULL
);

-- 4. leasing_tasks_attributes <<leasing>>
CREATE TABLE leasing.leasing_tasks_attributes (
    tasa_id       BIGSERIAL PRIMARY KEY,
    tasa_name     VARCHAR(55) NOT NULL,
    tasa_value    VARCHAR(255),                            -- nilai atribut (bisa JSON/text)
    tasa_status   VARCHAR(15) DEFAULT 'inprogress'
                  CHECK (tasa_status IN ('inprogress', 'pending', 'completed', 'cancelled')),
    tasa_task_id  BIGINT NOT NULL REFERENCES leasing.leasing_tasks(task_id) ON DELETE CASCADE,  -- tasa_task_id = task_id
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 5. payment_schedule <<payment>>
CREATE TABLE finance.payment_schedule (
    schedule_id        BIGSERIAL PRIMARY KEY,
    angsuran_ke        SMALLINT NOT NULL,
    jatuh_tempo        DATE NOT NULL,
    pokok              NUMERIC(15,2) NOT NULL,
    margin             NUMERIC(15,2) NOT NULL,
    total_tagihan      NUMERIC(15,2) NOT NULL,
    status_pembayaran  VARCHAR(20) DEFAULT 'unpaid'
                       CHECK (status_pembayaran IN ('unpaid', 'partial', 'paid', 'overdue')),
    tanggal_bayar      DATE,
    created_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    contract_id        BIGINT NOT NULL REFERENCES leasing.leasing_contract(contract_id) ON DELETE CASCADE
);

-- 6. payments <<payment>>
CREATE TABLE finance.payments (
    payment_id         BIGSERIAL PRIMARY KEY,
    nomor_bukti        VARCHAR(40) UNIQUE NOT NULL,
    jumlah_bayar       NUMERIC(15,2) NOT NULL,
    tanggal_bayar      DATE NOT NULL,
    metode_pembayaran  VARCHAR(30) NOT NULL,             -- transfer, tunai, virtual account, minimarket
    provider           VARCHAR(50),                       -- BCA VA, Alfamart, Indomaret, dll
    created_at         TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    contract_id        BIGINT NOT NULL REFERENCES leasing.leasing_contract(contract_id) ON DELETE CASCADE,
    schedule_id        BIGINT REFERENCES finance.payment_schedule(schedule_id) ON DELETE SET NULL
);

-- 7. leasing_contract_documents <<leasing>>
CREATE TABLE leasing.leasing_contract_documents (
    ldoc_id     BIGSERIAL PRIMARY KEY,
    file_name   VARCHAR(125) NOT NULL,
    file_size   DOUBLE PRECISION,
    file_type   VARCHAR(15) CHECK (file_type IN ('png','jpg','jpeg','pdf','doc','docx')),
    file_url    VARCHAR(255) NOT NULL,
    contract_id BIGINT NOT NULL REFERENCES leasing.leasing_contract(contract_id) ON DELETE CASCADE,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Index performa query & searching
CREATE INDEX idx_kontrak_customer ON leasing.leasing_contract(customer_id);
CREATE INDEX idx_kontrak_motor    ON leasing.leasing_contract(motor_id);
CREATE INDEX idx_kontrak_product  ON leasing.leasing_contract(product_id);
CREATE INDEX idx_kontrak_status   ON leasing.leasing_contract(status);
CREATE INDEX idx_jadwal_kontrak   ON finance.payment_schedule(contract_id);
CREATE INDEX idx_bayar_kontrak    ON finance.payments(contract_id);
CREATE INDEX idx_task_kontrak     ON leasing.leasing_tasks(contract_id);
CREATE INDEX idx_task_role        ON leasing.leasing_tasks(role_id);