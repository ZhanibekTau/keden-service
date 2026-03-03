-- Migration: 001_init.sql
-- Description: Create all tables for KEDEN service
-- Date: 2026-02-07

BEGIN;

-- ============================================
-- Table: roles
-- ============================================
CREATE TABLE IF NOT EXISTS roles (
    id              BIGSERIAL       PRIMARY KEY,
    name            VARCHAR(50)     NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_roles_name ON roles(name);

INSERT INTO roles (name) VALUES ('admin'), ('client') ON CONFLICT DO NOTHING;

-- ============================================
-- Table: users
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id              BIGSERIAL       PRIMARY KEY,
    email           VARCHAR(255)    NOT NULL,
    password_hash   VARCHAR(255)    NOT NULL,
    first_name      VARCHAR(255)    NOT NULL,
    last_name       VARCHAR(255)    NOT NULL,
    phone           VARCHAR(50)     NOT NULL,
    role_id         BIGINT          NOT NULL REFERENCES roles(id),
    account_type    VARCHAR(20)     NOT NULL,
    is_active       BOOLEAN         NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- ============================================
-- Table: companies
-- ============================================
CREATE TABLE IF NOT EXISTS companies (
    id              BIGSERIAL       PRIMARY KEY,
    user_id         BIGINT          NOT NULL REFERENCES users(id),
    company_name    VARCHAR(255)    NOT NULL,
    legal_name      VARCHAR(255)    NOT NULL,
    bin             VARCHAR(12)     NOT NULL,
    contact_person  VARCHAR(255),
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_companies_bin ON companies(bin);
CREATE INDEX IF NOT EXISTS idx_companies_user_id ON companies(user_id);

-- ============================================
-- Table: subscriptions
-- ============================================
CREATE TABLE IF NOT EXISTS subscriptions (
    id              BIGSERIAL       PRIMARY KEY,
    user_id         BIGINT          NOT NULL REFERENCES users(id),
    status          VARCHAR(50)     NOT NULL DEFAULT 'pending',
    start_date      TIMESTAMPTZ,
    end_date        TIMESTAMPTZ,
    amount          DECIMAL(10,2)   NOT NULL DEFAULT 12990.00,
    admin_comment   TEXT,
    requested_at    TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    approved_at     TIMESTAMPTZ,
    approved_by_id  BIGINT          REFERENCES users(id),
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status);

-- ============================================
-- Table: documents
-- ============================================
CREATE TABLE IF NOT EXISTS documents (
    id                  BIGSERIAL       PRIMARY KEY,
    user_id             BIGINT          NOT NULL REFERENCES users(id),
    original_name       VARCHAR(255)    NOT NULL,
    excel_file_path     VARCHAR(500),
    status              VARCHAR(50)     NOT NULL DEFAULT 'uploaded',
    error_message       TEXT,
    ai_response_json    JSONB,
    file_size           BIGINT          NOT NULL DEFAULT 0,
    queued_at           TIMESTAMPTZ,
    processed_at        TIMESTAMPTZ,
    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_documents_user_id ON documents(user_id);
CREATE INDEX IF NOT EXISTS idx_documents_status ON documents(status);

-- ============================================
-- Table: refresh_tokens
-- ============================================
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id              BIGSERIAL       PRIMARY KEY,
    user_id         BIGINT          NOT NULL REFERENCES users(id),
    token           VARCHAR(500)    NOT NULL,
    expires_at      TIMESTAMPTZ     NOT NULL,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);

COMMIT;
