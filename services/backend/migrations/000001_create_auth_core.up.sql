CREATE TABLE tenants (
    id uuid PRIMARY KEY,
    name varchar(160) NOT NULL,
    slug varchar(120) NOT NULL UNIQUE,
    status varchar(32) NOT NULL DEFAULT 'active',
    is_active boolean NOT NULL DEFAULT true,
    created_by uuid NULL,
    created_date timestamptz NOT NULL DEFAULT now(),
    updated_by uuid NULL,
    updated_date timestamptz NOT NULL DEFAULT now(),
    is_deleted boolean NOT NULL DEFAULT false
);

CREATE TABLE users (
    id uuid PRIMARY KEY,
    tenant_id uuid NULL REFERENCES tenants(id),
    email varchar(255) NOT NULL,
    name varchar(160) NOT NULL,
    password_hash text NOT NULL,
    status varchar(32) NOT NULL DEFAULT 'active',
    last_login_at timestamptz NULL,
    is_active boolean NOT NULL DEFAULT true,
    created_by uuid NULL,
    created_date timestamptz NOT NULL DEFAULT now(),
    updated_by uuid NULL,
    updated_date timestamptz NOT NULL DEFAULT now(),
    is_deleted boolean NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX idx_users_tenant_email_active ON users(tenant_id, lower(email)) WHERE is_deleted = false;
CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_is_active ON users(is_active);
CREATE INDEX idx_users_is_deleted ON users(is_deleted);

CREATE TABLE refresh_tokens (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id),
    token_hash text NOT NULL UNIQUE,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz NULL,
    replaced_by_token_id uuid NULL,
    is_active boolean NOT NULL DEFAULT true,
    created_by uuid NULL,
    created_date timestamptz NOT NULL DEFAULT now(),
    updated_by uuid NULL,
    updated_date timestamptz NOT NULL DEFAULT now(),
    is_deleted boolean NOT NULL DEFAULT false
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_is_deleted ON refresh_tokens(is_deleted);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
