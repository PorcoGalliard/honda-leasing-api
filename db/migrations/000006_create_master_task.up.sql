-- Schema: mst (master data)

DROP TABLE IF EXISTS mst.template_tasks CASCADE;
DROP TABLE IF EXISTS mst.template_task_attributes CASCADE;

-- 1. template_tasks <<mst>>
CREATE TABLE mst.template_tasks (
    teta_id     BIGSERIAL PRIMARY KEY,
    teta_name   VARCHAR(85) NOT NULL,  -- Nama tahapan, misal: "Pengajuan Awal", "Survei Lapangan"
    teta_role_id BIGINT NOT NULL REFERENCES account.roles(role_id) ON DELETE RESTRICT,
    -- optional
    description TEXT,                                     -- Deskripsi detail tahapan
    sequence_no SMALLINT DEFAULT 0,                       -- Urutan tahapan (untuk sorting alur)
    is_required BOOLEAN DEFAULT TRUE,                     -- optional
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    call_function TEXT DEFAULT NULL
);

-- 2. template_task_attributes <<mst>>
-- Atribut / field / dokumen yang diperlukan di setiap tahapan
CREATE TABLE mst.template_task_attributes (
    tetat_id       BIGSERIAL PRIMARY KEY,
    tetat_name     VARCHAR(85) NOT NULL,                 -- Nama atribut, misal: "Upload KTP", "Foto Rumah Depan"
    tetat_teta_id  BIGINT NOT NULL REFERENCES mst.template_tasks(teta_id) ON DELETE CASCADE,
    description    TEXT,
    is_required    BOOLEAN DEFAULT TRUE,
    attribute_type VARCHAR(50),                            -- misal: 'file', 'text', 'date', 'boolean', 'number'
    created_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- INDEX untuk performa query, searching
-- =============================================
CREATE INDEX idx_template_tasks_role ON mst.template_tasks(teta_role_id);
CREATE INDEX idx_template_tasks_seq  ON mst.template_tasks(sequence_no);
CREATE INDEX idx_tetat_teta          ON mst.template_task_attributes(tetat_teta_id);