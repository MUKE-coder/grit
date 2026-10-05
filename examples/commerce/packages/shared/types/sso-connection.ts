export interface SSOConnection {
  id: string;
  slug: string;
  name: string;
  protocol: string;
  domains: string;
  issuer_url: string;
  discovery_url: string;
  client_id: string;
  has_secret: boolean;
  scopes: string;
  enabled: boolean;
  jit_provisioning: boolean;
  default_role_id: string;
  groups_claim: string;
  group_mappings: string;
  metadata_url: string;
  metadata_xml: string;
  email_attribute: string;
  first_name_attribute: string;
  last_name_attribute: string;
  groups_attribute: string;
  allow_idp_initiated: boolean;
  last_used_at: string | null;
  created_at: string;
  updated_at: string;
}
