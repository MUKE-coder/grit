package generate

import (
	"fmt"
	"path/filepath"
	"strings"
)

// The three pages and the form.
//
// No JavaScript template literals anywhere in here. The templates are Go raw
// strings, a raw string cannot contain a backtick, and the escape dance that
// gets one in is how a generated file ends up with a stray quote nobody sees
// until the app will not build. String concatenation reads the same in the
// output and cannot go wrong in the input.

// writePortalForm writes the create-and-edit form for one resource.
//
// Driven by the generated descriptor rather than a field list baked into the
// component, so grit generate field reaches it without an edit here.
func (g *Generator) writePortalForm(names Names) error {
	path := filepath.Join(g.Root, "apps", "web", "components", "portal", names.Kebab+"-form.tsx")
	if fileExists(path) {
		return nil
	}
	content := fmt.Sprintf(`"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { %[1]sPortal } from "@/lib/portal/%[2]s";
import { PortalInput } from "@/components/portal/portal-field";
import type { %[3]s } from "@repo/shared/types";

// The customer's form for a %[4]s.
//
// Written once and never rewritten: this file is yours. The fields come from
// the descriptor, so a column added to the Go model turns up here on its own,
// while the layout, the copy and the order are edits you make and keep.
export function %[3]sForm({
  initial,
  onSubmit,
  submitting,
  error,
  submitLabel,
}: {
  initial?: Partial<%[3]s>;
  onSubmit: (values: Record<string, unknown>) => void;
  submitting?: boolean;
  error?: string | null;
  submitLabel: string;
}) {
  const router = useRouter();
  const [values, setValues] = useState<Record<string, unknown>>(
    () => ({ ...(initial ?? {}) }) as Record<string, unknown>,
  );

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        onSubmit(values);
      }}
      className="max-w-2xl space-y-5 rounded-xl border border-border bg-background p-6"
    >
      {%[1]sPortal.fields.map((field) => (
        <PortalInput
          key={field.name}
          field={field}
          value={values[field.name]}
          onChange={(value) => setValues((v) => ({ ...v, [field.name]: value }))}
        />
      ))}

      {error ? <p className="text-sm text-danger">{error}</p> : null}

      <div className="flex items-center gap-3 pt-2">
        <button
          type="submit"
          disabled={submitting}
          className="rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white transition-colors hover:opacity-90 disabled:opacity-60"
        >
          {submitting ? "Saving..." : submitLabel}
        </button>
        <button
          type="button"
          onClick={() => router.back()}
          className="rounded-lg px-4 py-2 text-sm text-text-secondary transition-colors hover:text-foreground"
        >
          Cancel
        </button>
      </div>
    </form>
  );
}
`,
		names.Camel,
		names.PluralKebab,
		names.Pascal,
		names.Lower,
	)
	return writeFileWithDirs(path, content)
}

// portalListPage is the customer's list of their own rows.
//
// Cards rather than a table, deliberately. The admin has a table because an
// operator audits every row of a shared dataset. A customer is looking at
// their own handful of things, and the generated screen should look like that
// before anybody touches it.
func (g *Generator) portalListPage(names Names) string {
	return fmt.Sprintf(`"use client";

import Link from "next/link";
import { use%[5]s } from "@/hooks/use-%[2]s";
import { %[1]sPortal } from "@/lib/portal/%[2]s";
import { PortalPage, PortalEmpty, PortalSkeleton, PortalError } from "@/components/portal/portal-page";

// %[6]s, as the customer sees them.
//
// Every row here is theirs: the API scopes the list by the %[7]s on the
// session, so this page asks for everything and receives only what it may
// have. There is no filter in this file doing that job, and there should not
// be: a check that lives in a page is a check an API client walks past.
export default function %[5]sPage() {
  const { data, isLoading, error } = use%[5]s();
  const rows = data?.data ?? [];

  return (
    <PortalPage
      title={%[1]sPortal.label}
      action={
        <Link
          href="/%[2]s/new"
          className="rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white transition-colors hover:opacity-90"
        >
          New
        </Link>
      }
    >
      {isLoading ? <PortalSkeleton /> : null}
      {error ? <PortalError /> : null}
      {!isLoading && !error && rows.length === 0 ? (
        <PortalEmpty
          title=%[8]q
          description="When you create one it will show up here."
          href="/%[2]s/new"
          cta="Create the first one"
        />
      ) : null}

      {rows.length > 0 ? (
        <ul className="grid gap-3 sm:grid-cols-2">
          {rows.map((row) => (
            <li key={row.id}>
              <Link
                href={"/%[2]s/" + row.id}
                className="block rounded-xl border border-border bg-background p-4 transition-colors hover:border-accent"
              >
                <p className="truncate font-medium text-foreground">
                  {String((row as unknown as Record<string, unknown>)[%[1]sPortal.display] ?? "Untitled")}
                </p>%[9]s
              </Link>
            </li>
          ))}
        </ul>
      ) : null}
    </PortalPage>
  );
}
`,
		names.Camel,
		names.PluralKebab,
		names.Pascal,
		names.Lower,
		names.PluralPascal,
		portalLabel(names),
		g.portalScopeWord(),
		"No "+strings.ToLower(portalLabel(names))+" yet",
		g.portalSubtitleMarkup(names),
	)
}

// portalSubtitleMarkup is the second line on a card, or nothing when the
// resource has no second text column worth showing.
func (g *Generator) portalSubtitleMarkup(names Names) string {
	if g.portalSubtitleField() == "" {
		return ""
	}
	return fmt.Sprintf(`
                <p className="mt-1 line-clamp-2 text-sm text-text-secondary">
                  {String((row as unknown as Record<string, unknown>)[%sPortal.subtitle] ?? "")}
                </p>`, names.Camel)
}

// portalScopeWord is who the rows belong to, for the comment at the top of a
// generated page. The distinction is the whole reason the page needs no filter.
func (g *Generator) portalScopeWord() string {
	if g.Definition.TenantOwned {
		return "organization"
	}
	return "user"
}

// portalNewPage is the create screen.
func (g *Generator) portalNewPage(names Names) string {
	return fmt.Sprintf(`"use client";

import { useRouter } from "next/navigation";
import type { Create%[3]sInput } from "@repo/shared/schemas";
import { useCreate%[3]s } from "@/hooks/use-%[2]s";
import { %[3]sForm } from "@/components/portal/%[4]s-form";
import { PortalPage } from "@/components/portal/portal-page";

export default function New%[3]sPage() {
  const router = useRouter();
  const create = useCreate%[3]s();

  return (
    <PortalPage title="New %[5]s">
      <%[3]sForm
        submitLabel="Create"
        submitting={create.isPending}
        error={create.error ? "That did not save. Check the fields and try again." : null}
        onSubmit={(values) => {
          // The form builds its payload from the generated descriptor, so the
          // shape is right at runtime and unprovable at compile time. One
          // cast, here, rather than an any loose in the form.
          create.mutate(values as Create%[3]sInput, {
            onSuccess: () => router.push("/%[2]s"),
          });
        }}
      />
    </PortalPage>
  );
}
`,
		names.Camel,
		names.PluralKebab,
		names.Pascal,
		names.Kebab,
		names.Lower,
	)
}

// portalDetailPage is the edit-and-delete screen for one row.
func (g *Generator) portalDetailPage(names Names) string {
	return fmt.Sprintf(`"use client";

import { useParams, useRouter } from "next/navigation";
import type { Update%[3]sInput } from "@repo/shared/schemas";
import { useGet%[3]s, useUpdate%[3]s, useDelete%[3]s } from "@/hooks/use-%[2]s";
import { %[3]sForm } from "@/components/portal/%[4]s-form";
import { PortalPage, PortalSkeleton, PortalError } from "@/components/portal/portal-page";

export default function %[3]sDetailPage() {
  const params = useParams<{ id: string }>();
  const id = params?.id ?? "";
  const router = useRouter();
  const { data: row, isLoading, error } = useGet%[3]s(id);
  const update = useUpdate%[3]s();
  const remove = useDelete%[3]s();

  if (isLoading) return <PortalSkeleton rows={1} />;
  // A row that is not yours answers 404, the same answer as a row that does
  // not exist. That is deliberate on the API side: anything else would tell a
  // stranger which ids are real.
  if (error || !row) return <PortalError message="That is not here." />;

  return (
    <PortalPage
      title="Edit %[5]s"
      action={
        <button
          type="button"
          onClick={() => {
            if (!confirm("Delete this? That cannot be undone.")) return;
            remove.mutate(id, { onSuccess: () => router.push("/%[2]s") });
          }}
          className="rounded-lg border border-danger/30 px-4 py-2 text-sm text-danger transition-colors hover:bg-danger/10"
        >
          Delete
        </button>
      }
    >
      <%[3]sForm
        initial={row}
        submitLabel="Save changes"
        submitting={update.isPending}
        error={update.error ? "That did not save. Check the fields and try again." : null}
        onSubmit={(values) => {
          // version is what the row was at when this page loaded. It goes out
          // as If-Match, so an edit saved by somebody else in the meantime
          // answers 409 instead of being silently overwritten.
          update.mutate(
            { id, version: row.version, ...values } as {
              id: string;
              version: number;
            } & Update%[3]sInput,
            { onSuccess: () => router.push("/%[2]s") },
          );
        }}
      />
    </PortalPage>
  );
}
`,
		names.Camel,
		names.PluralKebab,
		names.Pascal,
		names.Kebab,
		names.Lower,
	)
}
