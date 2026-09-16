package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review M41: the admin's screens redeclared the models' types
// instead of importing them from @repo/shared.
//
// grit sync writes packages/shared/types/<model>.ts for every Go model, and a
// page holding its own copy of Ticket or Notification drifts from it silently:
// a field added to the model reaches the shared type and never the page. The
// catch was that a new project had none of these files until sync first ran,
// so there was nothing to import.
//
// They are written with the project now, byte for byte what grit sync produces
// from the same models (internal/generate checks that), so the first sync
// changes nothing. A page imports the row from @repo/shared/types and narrows
// only what a Go struct cannot say: a string that is really a union, and a
// relation sync types as unknown.

// sharedModelType is one model's shared type file.
type sharedModelType struct {
	// Kebab is the file name under packages/shared/types, without .ts.
	Kebab string
	// Name is the exported interface.
	Name string
	// Source is the file, as grit sync writes it.
	Source string
	// Model is the Go model file it is generated from.
	Model func() string
}

// SharedModelTypes are the shared types a new project gets for the scaffold's
// own models. Exported for internal/generate's test that they match sync.
var SharedModelTypes = []sharedModelType{
	{Kebab: "api-key", Name: "APIKey", Model: apiAPIKeyModelGo, Source: `export interface APIKey {
  id: string;
  user_id: string;
  name: string;
  prefix: string;
  kind: string;
  token: string;
  scopes: string[];
  endpoints: string[];
  origins: string[];
  rate_limit: number;
  last_used_at: string | null;
  expires_at: string | null;
  revoked_at: string | null;
  created_at: string;
}
`},
	{Kebab: "form-share", Name: "FormShare", Model: formShareModelGo, Source: `export interface FormShare {
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
`},
	{Kebab: "form-submission", Name: "FormSubmission", Model: formSubmissionModelGo, Source: `export interface FormSubmission {
  id: string;
  share_id: string;
  resource_name: string;
  record_id: string;
  ip: string;
  user_agent: string;
  created_at: string;
}
`},
	{Kebab: "notification", Name: "Notification", Model: notificationModelGo, Source: `export interface Notification {
  id: string;
  user_id: string;
  source: string;
  severity: string;
  title: string;
  body: string;
  link: string;
  count: number;
  read_at: string | null;
  created_at: string;
  updated_at: string;
}
`},
	{Kebab: "sso-connection", Name: "SSOConnection", Model: apiSSOModelGo, Source: `export interface SSOConnection {
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
`},
	{Kebab: "ticket", Name: "Ticket", Model: ticketModelGo, Source: `export interface Ticket {
  id: string;
  user_id: string;
  subject: string;
  description: string;
  status: string;
  priority: string;
  labels: string;
  assignee_id: string;
  last_reply_at: string | null;
  closed_at: string | null;
  created_at: string;
  updated_at: string;
  user: unknown | null;
  assignee: unknown | null;
  replies: unknown[];
}
`},
	{Kebab: "ticket-reply", Name: "TicketReply", Model: ticketModelGo, Source: `export interface TicketReply {
  id: string;
  ticket_id: string;
  user_id: string;
  body: string;
  is_admin_reply: boolean;
  created_at: string;
  updated_at: string;
  user: unknown | null;
}
`},
}

// sharedModelTypeExports is the lines types/index.ts exports them with.
func sharedModelTypeExports() string {
	var b strings.Builder
	for _, t := range SharedModelTypes {
		fmt.Fprintf(&b, "export type { %s } from \"./%s\";\n", t.Name, t.Kebab)
	}
	return b.String()
}

// sharedModelTypeFiles maps each shared type file under sharedRoot to its source.
func sharedModelTypeFiles(sharedRoot string) map[string]string {
	files := map[string]string{}
	for _, t := range SharedModelTypes {
		files[filepath.Join(sharedRoot, "types", t.Kebab+".ts")] = t.Source
	}
	return files
}

// sharedProfileSchemas is packages/shared/schemas/profile.ts: the profile
// page's three forms, which were declared inside the page.
func sharedProfileSchemas() string {
	return `import { z } from "zod";

// The admin profile page's forms. Here rather than in the page so an app, a
// mobile client or a test validates a profile change the way the panel does.

export const PersonalInfoSchema = z.object({
  first_name: z.string().min(2, "First name must be at least 2 characters"),
  last_name: z.string().min(2, "Last name must be at least 2 characters"),
  email: z.string().email("Please enter a valid email"),
  // Only needed when the email changes.
  current_password: z.string().optional(),
});
export type PersonalInfoInput = z.infer<typeof PersonalInfoSchema>;

export const ProfessionalInfoSchema = z.object({
  job_title: z.string().optional().default(""),
  bio: z.string().optional().default(""),
});
export type ProfessionalInfoInput = z.infer<typeof ProfessionalInfoSchema>;

export const ChangePasswordSchema = z
  .object({
    current_password: z.string().min(1, "Enter your current password"),
    password: z.string().min(8, "Password must be at least 8 characters"),
    confirm_password: z.string().min(1, "Please confirm your password"),
  })
  .refine((d) => d.password === d.confirm_password, {
    message: "Passwords do not match",
    path: ["confirm_password"],
  });
export type ChangePasswordInput = z.infer<typeof ChangePasswordSchema>;
`
}

// sharedProfileSchemaExport is the line schemas/index.ts exports them with.
const sharedProfileSchemaExport = `export {
  PersonalInfoSchema,
  ProfessionalInfoSchema,
  ChangePasswordSchema,
  type PersonalInfoInput,
  type ProfessionalInfoInput,
  type ChangePasswordInput,
} from "./profile";
`

// ─── upgrade ──────────────────────────────────────────────────────────────────

// sharedRoots are the shared packages a project has: the workspace package of
// a monorepo, and the mirror a single's SPA keeps under src/shared.
func sharedRoots(root string) []string {
	var roots []string
	for _, dir := range []string{
		filepath.Join(root, "packages", "shared"),
		filepath.Join(root, "frontend", "src", "shared"),
	} {
		if fileExists(filepath.Join(dir, "types", "index.ts")) {
			roots = append(roots, dir)
		}
	}
	return roots
}

// addMissingExports inserts each export line whose name the index does not
// export yet, before the marker grit generate appends at. Names already
// exported, from wherever, are left alone.
func addMissingExports(src, marker string, exports map[string]string) (string, []string) {
	if !strings.Contains(src, marker) {
		return src, nil
	}
	names := make([]string, 0, len(exports))
	for name := range exports {
		names = append(names, name)
	}
	sort.Strings(names)
	var added []string
	var lines strings.Builder
	for _, name := range names {
		if exportsName(src, name) {
			continue
		}
		lines.WriteString(exports[name])
		added = append(added, name)
	}
	if len(added) == 0 {
		return src, nil
	}
	return strings.Replace(src, marker, lines.String()+marker, 1), added
}

// exportsName reports whether an index already exports name, in either of the
// forms the scaffold and the generator write.
func exportsName(src, name string) bool {
	for _, form := range []string{"{ " + name + " }", "  " + name + ",", "type " + name + ",", "type " + name + " }"} {
		if strings.Contains(src, form) {
			return true
		}
	}
	return false
}

// repairSharedTypes gives an existing project the shared model types and the
// profile schemas the admin now imports.
//
// A file that exists is left exactly as it is: it is either grit sync's, which
// is the same thing, or the project's own. Only the index gains exports.
func repairSharedTypes(root string) error {
	for _, shared := range sharedRoots(root) {
		var created []string
		files := sharedModelTypeFiles(shared)
		files[filepath.Join(shared, "schemas", "profile.ts")] = sharedProfileSchemas()
		paths := make([]string, 0, len(files))
		for path := range files {
			paths = append(paths, path)
		}
		sort.Strings(paths)
		for _, path := range paths {
			made, err := createIfMissing(path, files[path])
			if err != nil {
				return err
			}
			if made {
				created = append(created, filepath.Base(path))
			}
		}

		typeExports := map[string]string{}
		for _, t := range SharedModelTypes {
			typeExports[t.Name] = fmt.Sprintf("export type { %s } from \"./%s\";\n", t.Name, t.Kebab)
		}
		exported, err := addExportsToIndex(filepath.Join(shared, "types", "index.ts"), "// grit:types", typeExports)
		if err != nil {
			return err
		}
		schemaExports := map[string]string{"PersonalInfoSchema": sharedProfileSchemaExport}
		exportedSchemas, err := addExportsToIndex(filepath.Join(shared, "schemas", "index.ts"), "// grit:schemas", schemaExports)
		if err != nil {
			return err
		}
		exported = append(exported, exportedSchemas...)

		shown := shared
		if rel, err := filepath.Rel(root, shared); err == nil {
			shown = filepath.ToSlash(rel)
		}
		if len(created) > 0 || len(exported) > 0 {
			fmt.Printf("  ✓ %s: the admin imports model types and the profile schemas from here (%d file(s) added, %d export(s))\n",
				shown, len(created), len(exported))
		}
	}
	return nil
}

// addExportsToIndex applies addMissingExports to an index file on disk.
func addExportsToIndex(path, marker string, exports map[string]string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	crlf := strings.Contains(string(data), "\r\n")
	src := strings.ReplaceAll(string(data), "\r\n", "\n")
	out, added := addMissingExports(src, marker, exports)
	if len(added) == 0 {
		return nil, nil
	}
	if crlf {
		out = strings.ReplaceAll(out, "\n", "\r\n")
	}
	// An index is where grit generate adds exports too, so an edited one is
	// normal: the lines are inserted either way, and the recorded hash moves
	// with them only when the file was Grit's to begin with.
	pristine := manifest.IsUnchanged(path)
	if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
		return nil, fmt.Errorf("writing %s: %w", path, err)
	}
	if pristine {
		manifest.Refresh(path)
	}
	return added, nil
}
