export interface FormShare {
  id: string;
  resource_name: string;
  token: string;
  has_password: boolean;
  enabled: boolean;
  submission_count: number;
  created_by_user_id: string;
  label: string;
  custom_title: string;
  custom_description: string;
  hidden_fields: string[];
  created_at: string;
  updated_at: string;
}
