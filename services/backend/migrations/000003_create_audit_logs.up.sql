CREATE TABLE audit_logs (
    id uuid PRIMARY KEY,
    tenant_id uuid NULL,
    actor_user_id uuid NULL,
    action varchar(120) NOT NULL,
    resource_type varchar(120) NOT NULL,
    resource_id varchar(120) NULL,
    ip_address inet NULL,
    user_agent text NULL,
    metadata jsonb NOT NULL DEFAULT '{}'::jsonb,
    is_active boolean NOT NULL DEFAULT true,
    created_by uuid NULL,
    created_date timestamptz NOT NULL DEFAULT now(),
    updated_by uuid NULL,
    updated_date timestamptz NOT NULL DEFAULT now(),
    is_deleted boolean NOT NULL DEFAULT false
);

CREATE INDEX idx_audit_logs_tenant_id ON audit_logs(tenant_id);
CREATE INDEX idx_audit_logs_actor_user_id ON audit_logs(actor_user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_resource ON audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_created_date ON audit_logs(created_date);
CREATE INDEX idx_audit_logs_is_deleted ON audit_logs(is_deleted);
