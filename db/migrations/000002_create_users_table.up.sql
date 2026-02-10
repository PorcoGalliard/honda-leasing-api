-- Schema: account (autentikasi & authorization)

DROP TABLE IF EXISTS account.users CASCADE;
DROP TABLE IF EXISTS account.oauth_providers CASCADE;
DROP TABLE IF EXISTS account.user_oauth_provider CASCADE;
DROP TABLE IF EXISTS account.roles CASCADE;
DROP TABLE IF EXISTS account.user_roles CASCADE;
DROP TABLE IF EXISTS account.permissions CASCADE;
DROP TABLE IF EXISTS account.role_permission CASCADE;

-- 1. Tabel utama: users <<account>>
CREATE TABLE account.users (
    user_id         BIGSERIAL PRIMARY KEY,
    username        VARCHAR(50) UNIQUE,
    phone_number    VARCHAR(15) UNIQUE NOT NULL, -- utama untuk login OTP/PIN
    email           VARCHAR(100) UNIQUE,
    full_name       VARCHAR(100) NOT NULL,
    password        VARCHAR(255),-- hashed password
    pin_key         VARCHAR(255),-- hashed 6-digit PIN
    is_active       BOOLEAN DEFAULT TRUE,
    last_login      TIMESTAMPTZ,
    failed_attempts SMALLINT DEFAULT 0,
    locked_until    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 2. oauth_providers <<account>>
CREATE TABLE account.oauth_providers (
    provider_id     BIGSERIAL PRIMARY KEY,
    provider_name   VARCHAR(50) NOT NULL UNIQUE,-- google|apple|facebook|keycloak
    client_id       VARCHAR(255) NOT NULL,
    client_secret   VARCHAR(255) NOT NULL,-- ENCRYPT di app level!
    redirect_uri    VARCHAR(255) NOT NULL,
    issuer_url      VARCHAR(255),
    active          BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 3. user_oauth_provider <<account>>
CREATE TABLE account.user_oauth_provider (
    user_oauth_id   BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES account.users(user_id) ON DELETE CASCADE,
    provider_id     BIGINT NOT NULL REFERENCES account.oauth_providers(provider_id) ON DELETE RESTRICT,
    access_token    TEXT,
    refresh_token   TEXT,
    expires_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, provider_id)
);

-- 4. roles <<account>>
CREATE TABLE account.roles (
    role_id         BIGSERIAL PRIMARY KEY,
    role_name       VARCHAR(50) NOT NULL UNIQUE,-- ADMIN|SALES|SURVEYOR|COLLECTION|FINANCE|CUSTOMER|SYSTEM
    description     TEXT,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 5. user_roles <<account>>
CREATE TABLE account.user_roles (
    user_role_id    BIGSERIAL PRIMARY KEY,
    user_id         BIGINT NOT NULL REFERENCES account.users(user_id) ON DELETE CASCADE,
    role_id         BIGINT NOT NULL REFERENCES account.roles(role_id) ON DELETE RESTRICT,
    assigned_at     TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    assigned_by     BIGINT REFERENCES account.users(user_id),-- admin yang assign untuk roles
    UNIQUE (user_id, role_id)
);

-- 6. permissions
CREATE TABLE account.permissions (
    permission_id   BIGSERIAL PRIMARY KEY,
    permission_type VARCHAR(100) NOT NULL UNIQUE,-- view_contract|approve_survey|dll
    description     TEXT,
    created_at      TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- 7. role_permission
CREATE TABLE account.role_permission (
    role_permission_id BIGSERIAL PRIMARY KEY,
    role_id            BIGINT NOT NULL REFERENCES account.roles(role_id) ON DELETE CASCADE,
    permission_id      BIGINT NOT NULL REFERENCES account.permissions(permission_id) ON DELETE CASCADE,
    UNIQUE (role_id, permission_id)
);

-- =============================================
-- CREAETE INDEX untuk performa (login, query role, oauth)
-- =============================================
CREATE INDEX idx_users_phone_login    ON account.users(phone_number) WHERE is_active = TRUE;
CREATE INDEX idx_users_email          ON account.users(email);
CREATE INDEX idx_user_oauth_user      ON account.user_oauth_provider(user_id);
CREATE INDEX idx_user_oauth_provider  ON account.user_oauth_provider(provider_id);
CREATE INDEX idx_user_roles_user      ON account.user_roles(user_id);
CREATE INDEX idx_role_permission_role ON account.role_permission(role_id);


