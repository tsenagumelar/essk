CREATE TABLE products (
    id uuid PRIMARY KEY,
    tenant_id uuid NOT NULL REFERENCES tenants(id),
    code varchar(80) NOT NULL,
    name varchar(160) NOT NULL,
    category varchar(120) NOT NULL,
    price_cents integer NOT NULL DEFAULT 0,
    status varchar(40) NOT NULL DEFAULT 'active',
    is_active boolean NOT NULL DEFAULT true,
    created_by uuid NULL REFERENCES users(id),
    created_date timestamptz NOT NULL DEFAULT now(),
    updated_by uuid NULL REFERENCES users(id),
    updated_date timestamptz NOT NULL DEFAULT now(),
    is_deleted boolean NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX idx_products_tenant_code_active ON products(tenant_id, lower(code)) WHERE is_deleted = false;
CREATE INDEX idx_products_tenant_id ON products(tenant_id);
CREATE INDEX idx_products_is_deleted ON products(is_deleted);
