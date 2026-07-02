CREATE TABLE permissions (
    id uuid PRIMARY KEY,
    code varchar(160) NOT NULL UNIQUE,
    name varchar(160) NOT NULL,
    description text NULL,
    is_active boolean NOT NULL DEFAULT true,
    created_by uuid NULL,
    created_date timestamptz NOT NULL DEFAULT now(),
    updated_by uuid NULL,
    updated_date timestamptz NOT NULL DEFAULT now(),
    is_deleted boolean NOT NULL DEFAULT false
);

CREATE TABLE roles (
    id uuid PRIMARY KEY,
    tenant_id uuid NULL REFERENCES tenants(id),
    name varchar(120) NOT NULL,
    code varchar(120) NOT NULL,
    description text NULL,
    is_system boolean NOT NULL DEFAULT false,
    is_active boolean NOT NULL DEFAULT true,
    created_by uuid NULL,
    created_date timestamptz NOT NULL DEFAULT now(),
    updated_by uuid NULL,
    updated_date timestamptz NOT NULL DEFAULT now(),
    is_deleted boolean NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX idx_roles_tenant_code_active ON roles(tenant_id, lower(code)) WHERE is_deleted = false;
CREATE INDEX idx_roles_tenant_id ON roles(tenant_id);
CREATE INDEX idx_roles_is_deleted ON roles(is_deleted);

CREATE TABLE role_permissions (
    role_id uuid NOT NULL REFERENCES roles(id),
    permission_id uuid NOT NULL REFERENCES permissions(id),
    is_active boolean NOT NULL DEFAULT true,
    created_by uuid NULL,
    created_date timestamptz NOT NULL DEFAULT now(),
    updated_by uuid NULL,
    updated_date timestamptz NOT NULL DEFAULT now(),
    is_deleted boolean NOT NULL DEFAULT false,
    PRIMARY KEY (role_id, permission_id)
);

CREATE INDEX idx_role_permissions_is_deleted ON role_permissions(is_deleted);

CREATE TABLE user_roles (
    user_id uuid NOT NULL REFERENCES users(id),
    role_id uuid NOT NULL REFERENCES roles(id),
    is_active boolean NOT NULL DEFAULT true,
    created_by uuid NULL,
    created_date timestamptz NOT NULL DEFAULT now(),
    updated_by uuid NULL,
    updated_date timestamptz NOT NULL DEFAULT now(),
    is_deleted boolean NOT NULL DEFAULT false,
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX idx_user_roles_is_deleted ON user_roles(is_deleted);
