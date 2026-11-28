package generate

import (
	"fmt"
	"path/filepath"
	"strings"
)

// The TSX the portal writes.
//
// Three rules hold this file together, and breaking any of them produces a
// screen that compiles and looks wrong.
//
//  1. **The web app's palette, not the admin's.** apps/web has
//     bg-bg-secondary, text-text-secondary and text-accent. The admin's
//     muted / card / primary are not utilities there, and Tailwind renders an
//     unknown utility as nothing at all: no error, no colour.
//  2. **A card list, not a data table.** A customer is looking at their own
//     things, not auditing a table, and the generated screen should look like
//     the former even before anybody edits it.
//  3. **Nothing here is overwritten on a second run.** The pages are the first
//     files anybody edits. A new field reaches the form through the
//     descriptor, which is generated.

// portalFormField is one editable field, resolved once and used by both the
// descriptor and the form.
type portalFormField struct {
	Name     string // Go field name, PascalCase
	JSON     string // the json tag, which is what the API takes
	Label    string
	Kind     string // text, textarea, number, money, checkbox, date, datetime, select, relation, file
	Required bool
	Options  []string // select and radio
	Relation string   // the related resource, plural kebab, for a relation
}

// portalFields is every field a customer may edit, in declaration order.
//
// Deny by default is not possible here without a second declaration, so the
// rule is the next best thing: everything a customer could sensibly edit,
// minus the columns that are not theirs to touch. The owner column above all.
// A customer setting user_id would be giving their row away, and the service
// stamps it from the session anyway.
func (g *Generator) portalFields() []portalFormField {
	owner := g.Definition.OwnerField()
	var out []portalFormField
	for _, f := range g.Definition.Fields {
		if f.IsAuto() || f.ServerSet || f.IsSlug() {
			continue
		}
		if owner != nil && f.Name == owner.Name {
			continue
		}
		// Written and never read back. A form cannot show the current value,
		// so an edit screen would blank it on every save.
		if f.Sensitive {
			continue
		}
		if f.IsManyToMany() {
			// A picker of its own, and the admin's is not portable here yet.
			continue
		}
		out = append(out, portalFormField{
			Name:     toPascalCase(f.Name),
			JSON:     toSnakeCase(f.Name),
			Label:    humanLabel(f.Name),
			Kind:     portalKind(f),
			Required: f.Required,
			Options:  portalOptions(f),
			Relation: portalRelation(f),
		})
	}
	return out
}

// portalRelation is the related resource as the URL spells it, for a
// belongs_to, and empty for everything else.
func portalRelation(f Field) string {
	if !f.IsBelongsTo() {
		return ""
	}
	return MakeNames(f.RelatedModelName()).PluralKebab
}

// portalKind maps a field type onto the input the portal form renders.
//
// Everything not named here falls to a text input, which is wrong for nothing
// and perfect for a few: a field type the portal has not learned yet still
// saves its value.
func portalKind(f Field) string {
	switch {
	case f.IsBelongsTo():
		return "relation"
	case f.IsFile() || f.IsFiles():
		return "file"
	case f.IsSelect() || f.IsRadio():
		return "select"
	case f.IsToggle() || f.IsCheck() || f.GoType() == "bool":
		return "checkbox"
	}
	switch FieldType(f.Type) {
	case FieldText, FieldRichtext:
		return "textarea"
	case FieldMoney:
		return "money"
	case FieldDate:
		return "date"
	case FieldDatetime:
		return "datetime"
	}
	switch f.GoType() {
	case "int", "int64", "uint", "float64":
		return "number"
	}
	return "text"
}

// portalDisplayField is what a card shows as its title: the first text-ish
// column, which is what the admin's display field logic also lands on.
func (g *Generator) portalDisplayField() string {
	for _, f := range g.Definition.Fields {
		if f.IsRelationship() || f.IsFile() || f.IsFiles() {
			continue
		}
		switch portalKind(f) {
		case "text", "textarea":
			return toSnakeCase(f.Name)
		}
	}
	return "id"
}

// portalSubtitleField is the second line on a card, or empty.
func (g *Generator) portalSubtitleField() string {
	title := g.portalDisplayField()
	for _, f := range g.Definition.Fields {
		name := toSnakeCase(f.Name)
		if name == title || f.IsRelationship() || f.IsFile() || f.IsFiles() {
			continue
		}
		switch portalKind(f) {
		case "text", "textarea", "select":
			return name
		}
	}
	return ""
}

// writePortalDescriptor writes what the headless level is for: the shape of
// the resource as data, with no React in it.
func (g *Generator) writePortalDescriptor(names Names) error {
	fields := g.portalFields()
	var rows []string
	for _, f := range fields {
		parts := []string{
			fmt.Sprintf("name: %q", f.JSON),
			fmt.Sprintf("label: %q", f.Label),
			fmt.Sprintf("kind: %q", f.Kind),
		}
		if f.Required {
			parts = append(parts, "required: true")
		}
		if len(f.Options) > 0 {
			quoted := make([]string, 0, len(f.Options))
			for _, o := range f.Options {
				quoted = append(quoted, fmt.Sprintf("%q", o))
			}
			parts = append(parts, "options: ["+strings.Join(quoted, ", ")+"]")
		}
		if f.Relation != "" {
			parts = append(parts, fmt.Sprintf("relation: %q", f.Relation))
		}
		rows = append(rows, "  { "+strings.Join(parts, ", ")+" },")
	}

	scope := "user"
	if g.Definition.TenantOwned {
		scope = "organization"
	}

	content := fmt.Sprintf(`import type { PortalField } from "./types";

// %[1]s, as data.
//
// Generated from the Go model, and rewritten whenever the model changes, so a
// column added with grit generate field turns up here without an edit. Build
// your own screen against this and it keeps up with the model; hard-code the
// fields in a component and it will not.
//
// Nothing in here is a permission check. The API scopes every read and write
// by the %[4]s on the session, which is what makes these rows the
// customer's. This describes what to render, never what is allowed.
export const %[2]sPortal = {
  resource: "%[3]s",
  label: %[5]q,
  path: "/%[3]s",
  api: "/api/%[3]s",
  scope: "%[4]s",
  display: %[6]q,
  subtitle: %[7]q,
  fields: [
%[8]s
  ] satisfies PortalField[],
} as const;

export type %[9]sPortalField = (typeof %[2]sPortal)["fields"][number];
`,
		portalLabel(names),
		names.Camel,
		names.PluralKebab,
		scope,
		portalLabel(names),
		g.portalDisplayField(),
		g.portalSubtitleField(),
		strings.Join(rows, "\n"),
		names.Pascal,
	)

	path := filepath.Join(g.Root, "apps", "web", "lib", "portal", names.PluralKebab+".ts")
	return writeFileWithDirs(path, content)
}

// ensurePortalKit writes the shared pieces, once per project.
//
// Each is skipped when it already exists, so an edited kit survives every
// later generate. That is the same bargain the Expo kit makes.
func (g *Generator) ensurePortalKit() error {
	web := filepath.Join(g.Root, "apps", "web")
	for _, file := range []struct {
		path string
		body string
	}{
		{filepath.Join(web, "lib", "portal", "types.ts"), portalTypesTS},
		{filepath.Join(web, "components", "portal", "portal-page.tsx"), portalPageTSX},
		{filepath.Join(web, "components", "portal", "portal-field.tsx"), portalFieldTSX},
	} {
		if fileExists(file.path) {
			continue
		}
		if err := writeFileWithDirs(file.path, file.body); err != nil {
			return err
		}
	}
	return nil
}

const portalTypesTS = `// The shape of a generated portal descriptor.
//
// One file for every resource, so a screen written against PortalField keeps
// working when a resource gains a column.

export type PortalFieldKind =
  | "text"
  | "textarea"
  | "number"
  | "money"
  | "checkbox"
  | "date"
  | "datetime"
  | "select"
  | "relation"
  | "file";

export interface PortalField {
  /** The column name, which is what the API takes and returns. */
  name: string;
  label: string;
  kind: PortalFieldKind;
  required?: boolean;
  /** The closed set, for a select. */
  options?: readonly string[];
  /** The related resource as its URL spells it, for a relation. */
  relation?: string;
}

export interface PortalDescriptor {
  resource: string;
  label: string;
  path: string;
  api: string;
  /** Who the rows belong to: the signed-in person, or their organization. */
  scope: "user" | "organization";
  display: string;
  subtitle: string;
  fields: readonly PortalField[];
}
`

const portalPageTSX = `"use client";

import Link from "next/link";
import type { ReactNode } from "react";

// The chrome every generated customer screen shares: a heading, an action, and
// the two states a list has before it has anything in it.
//
// Separate from the screens so that changing how they all look is one edit.
// The screens themselves are yours to rewrite; this is the part that is worth
// keeping consistent.

export function PortalPage({
  title,
  description,
  action,
  children,
}: {
  title: string;
  description?: string;
  action?: ReactNode;
  children: ReactNode;
}) {
  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-foreground">{title}</h1>
          {description ? (
            <p className="mt-1 text-sm text-text-secondary">{description}</p>
          ) : null}
        </div>
        {action}
      </div>
      {children}
    </div>
  );
}

export function PortalEmpty({
  title,
  description,
  href,
  cta,
}: {
  title: string;
  description: string;
  href?: string;
  cta?: string;
}) {
  return (
    <div className="rounded-xl border border-dashed border-border bg-background p-12 text-center">
      <h2 className="text-base font-semibold text-foreground">{title}</h2>
      <p className="mx-auto mt-2 max-w-sm text-sm text-text-secondary">{description}</p>
      {href && cta ? (
        <Link
          href={href}
          className="mt-6 inline-flex items-center rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white transition-colors hover:opacity-90"
        >
          {cta}
        </Link>
      ) : null}
    </div>
  );
}

export function PortalSkeleton({ rows = 3 }: { rows?: number }) {
  return (
    <div className="space-y-3">
      {Array.from({ length: rows }).map((_, i) => (
        <div key={i} className="h-20 animate-pulse rounded-xl border border-border bg-bg-secondary" />
      ))}
    </div>
  );
}

export function PortalError({ message }: { message?: string }) {
  return (
    <div className="rounded-xl border border-danger/30 bg-danger/5 p-6 text-sm text-danger">
      {message ?? "That did not load. Try again in a moment."}
    </div>
  );
}
`

const portalFieldTSX = `"use client";

import { useState } from "react";
import type { PortalField } from "@/lib/portal/types";

// One input per field kind, driven by the generated descriptor.
//
// It is deliberately plain. A customer-facing form is where a product earns
// its look, and anything clever here would be in the way of that. What it does
// buy is that a column added to the Go model appears on the form without
// anybody editing a component.

export function PortalInput({
  field,
  value,
  onChange,
}: {
  field: PortalField;
  value: unknown;
  onChange: (value: unknown) => void;
}) {
  const base =
    "w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground " +
    "outline-none transition-colors placeholder:text-text-secondary focus:border-accent";

  if (field.kind === "checkbox") {
    return (
      <label className="flex items-center gap-2 text-sm text-foreground">
        <input
          type="checkbox"
          checked={Boolean(value)}
          onChange={(e) => onChange(e.target.checked)}
          className="h-4 w-4 rounded border-border accent-accent"
        />
        {field.label}
      </label>
    );
  }

  const label = (
    <label htmlFor={field.name} className="mb-1.5 block text-sm font-medium text-foreground">
      {field.label}
      {field.required ? <span className="ml-0.5 text-danger">*</span> : null}
    </label>
  );

  if (field.kind === "textarea") {
    return (
      <div>
        {label}
        <textarea
          id={field.name}
          rows={5}
          required={field.required}
          value={(value as string) ?? ""}
          onChange={(e) => onChange(e.target.value)}
          className={base}
        />
      </div>
    );
  }

  if (field.kind === "select") {
    return (
      <div>
        {label}
        <select
          id={field.name}
          required={field.required}
          value={(value as string) ?? ""}
          onChange={(e) => onChange(e.target.value)}
          className={base}
        >
          <option value="">Choose one</option>
          {(field.options ?? []).map((option) => (
            <option key={option} value={option}>
              {option}
            </option>
          ))}
        </select>
      </div>
    );
  }

  if (field.kind === "file") {
    return (
      <div>
        {label}
        <PortalUpload value={value} onChange={onChange} />
      </div>
    );
  }

  const type =
    field.kind === "number" || field.kind === "money"
      ? "number"
      : field.kind === "date"
        ? "date"
        : field.kind === "datetime"
          ? "datetime-local"
          : "text";

  return (
    <div>
      {label}
      <input
        id={field.name}
        type={type}
        required={field.required}
        step={field.kind === "money" ? "0.01" : undefined}
        value={(value as string | number) ?? ""}
        onChange={(e) =>
          onChange(
            field.kind === "number" || field.kind === "money"
              ? e.target.value === ""
                ? null
                : Number(e.target.value)
              : e.target.value,
          )
        }
        className={base}
      />
    </div>
  );
}

// PortalUpload sends the file to the API and keeps what comes back.
//
// The API answers with the whole FileRef, which is what the column stores: the
// key, the URL, the MIME type, the size, and the renditions when the upload
// was an image. Storing only the URL is how a file loses its renditions.
function PortalUpload({
  value,
  onChange,
}: {
  value: unknown;
  onChange: (value: unknown) => void;
}) {
  const [busy, setBusy] = useState(false);
  const [failed, setFailed] = useState<string | null>(null);
  const current = value as { url?: string; name?: string } | null;

  async function upload(file: File) {
    setBusy(true);
    setFailed(null);
    try {
      const body = new FormData();
      body.append("file", file);
      // No Content-Type header: the browser sets the multipart boundary, and
      // naming the type by hand is how an upload arrives with no body.
      const res = await fetch("/api/uploads", {
        method: "POST",
        body,
        credentials: "include",
      });
      if (!res.ok) {
        throw new Error("Upload failed");
      }
      const json = await res.json();
      onChange(json.data);
    } catch (err) {
      setFailed(err instanceof Error ? err.message : "Upload failed");
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="space-y-2">
      <input
        type="file"
        disabled={busy}
        onChange={(e) => {
          const file = e.target.files?.[0];
          if (file) void upload(file);
        }}
        className="block w-full text-sm text-text-secondary file:mr-3 file:rounded-lg file:border-0 file:bg-bg-secondary file:px-3 file:py-2 file:text-sm file:text-foreground"
      />
      {busy ? <p className="text-xs text-text-secondary">Uploading...</p> : null}
      {failed ? <p className="text-xs text-danger">{failed}</p> : null}
      {current?.url ? (
        <a
          href={current.url}
          target="_blank"
          rel="noreferrer"
          className="inline-block text-xs text-accent underline"
        >
          {current.name ?? "View current file"}
        </a>
      ) : null}
    </div>
  );
}
`

// portalOptions is the closed set of a choice column, values only.
func portalOptions(f Field) []string {
	tag := f.EnumTag()
	if tag == "" {
		return nil
	}
	return strings.Split(tag, "|")
}
