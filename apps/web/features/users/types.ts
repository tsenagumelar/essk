export type AdminUser = {
  id: string;
  tenant_id?: string;
  email: string;
  name: string;
  status: string;
  is_active: boolean;
  role_ids: string[];
};
