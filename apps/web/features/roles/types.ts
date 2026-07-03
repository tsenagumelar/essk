export type Role = {
  id: string;
  tenant_id?: string;
  name: string;
  code: string;
  description?: string;
  is_system: boolean;
  is_active: boolean;
};

export type Permission = {
  id: string;
  code: string;
  name: string;
  description?: string;
  is_active: boolean;
};
