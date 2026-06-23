import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The admin panel is the operator side. The customer-facing web
        app (<code>apps/web</code>) is the other consumer of your
        generated API. This lesson shows you how to use the
        auto-generated React Query hook, the shared types, and the
        shared Zod schemas to build a list page and a create form in
        the web app without duplicating any types or contracts.
      </p>

      <h2>What the generator gave you for the web</h2>
      <p>
        For every <code>grit generate resource Contact …</code>, the
        web side gets three files dropped into{' '}
        <code>apps/web</code>:
      </p>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">File</th>
              <th className="text-left px-3 py-2 font-medium">What it does</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">apps/web/hooks/use-contacts.ts</td>
              <td className="px-3 py-2 text-[12px]">React Query hooks — useContacts, useContact, useCreateContact, useUpdateContact, useDeleteContact.</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">packages/shared/types/contact.ts</td>
              <td className="px-3 py-2 text-[12px]">TypeScript interface — what a Contact looks like coming back from the API.</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">packages/shared/schemas/contact.ts</td>
              <td className="px-3 py-2 text-[12px]">Zod schemas — CreateContactSchema, UpdateContactSchema. Use them in forms and at API boundaries.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <p>
        All three are <em>importable from the web app</em>. Same source
        of truth as the admin and the Go API.
      </p>

      <h2>1. Listing contacts on a customer page</h2>
      <p>
        The minimum useful list — paginated, with a loading state and
        an empty state:
      </p>

      <CodeBlock
        language="tsx"
        filename="apps/web/app/contacts/page.tsx"
        code={`"use client";

import { useContacts } from "@/hooks/use-contacts";

export default function ContactsPage() {
  const { data, isLoading, error } = useContacts({ page: 1 });

  if (isLoading) return <p>Loading…</p>;
  if (error)     return <p className="text-red-500">Failed to load.</p>;
  if (!data || data.data.length === 0) return <p>No contacts yet.</p>;

  return (
    <ul className="space-y-2">
      {data.data.map((c) => (
        <li key={c.id} className="rounded-lg border p-3">
          <p className="font-medium">{c.name}</p>
          <p className="text-sm text-gray-500">{c.email}</p>
        </li>
      ))}
    </ul>
  );
}`}
      />

      <p>
        Three things to notice:
      </p>
      <ul>
        <li>
          <strong>No <code>fetch</code> call.</strong> The hook owns
          the axios client, the auth cookies, the React Query cache —
          you just call it.
        </li>
        <li>
          <strong>Typed all the way through.</strong>{' '}
          <code>c.name</code> autocompletes, <code>c.xyz</code>{' '}
          errors. The type comes from{' '}
          <code>@repo/shared/types</code> via the hook.
        </li>
        <li>
          <strong>Pagination shape is shared.</strong>{' '}
          <code>data.data</code> is the row array;{' '}
          <code>data.meta</code> has total/page/page_size/pages — same
          shape every endpoint returns.
        </li>
      </ul>

      <h2>2. Searching + paginating</h2>
      <p>
        The hook accepts the same query-string params the API does:
      </p>

      <CodeBlock
        language="tsx"
        code={`"use client";

import { useState } from "react";
import { useContacts } from "@/hooks/use-contacts";

export default function ContactsPage() {
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);

  const { data, isLoading } = useContacts({ search, page, page_size: 20 });

  return (
    <>
      <input
        value={search}
        onChange={(e) => { setSearch(e.target.value); setPage(1); }}
        placeholder="Search by name or email…"
        className="w-full rounded-lg border px-3 py-2"
      />

      {isLoading ? (
        <p>Loading…</p>
      ) : (
        <ul>{data?.data.map((c) => <li key={c.id}>{c.name}</li>)}</ul>
      )}

      <div className="mt-4 flex gap-2">
        <button onClick={() => setPage((p) => Math.max(1, p - 1))} disabled={page === 1}>Prev</button>
        <span>Page {page} of {data?.meta.pages ?? 1}</span>
        <button onClick={() => setPage((p) => p + 1)} disabled={page >= (data?.meta.pages ?? 1)}>Next</button>
      </div>
    </>
  );
}`}
      />

      <h2>3. Loading a single contact</h2>

      <CodeBlock
        language="tsx"
        filename="apps/web/app/contacts/[id]/page.tsx"
        code={`"use client";

import { use } from "react";
import { useContact } from "@/hooks/use-contacts";

export default function ContactDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = use(params);
  const { data, isLoading, error } = useContact(id);

  if (isLoading) return <p>Loading…</p>;
  if (error || !data) return <p>Not found.</p>;

  return (
    <article>
      <h1>{data.name}</h1>
      <p>{data.email}</p>
      <p>{data.phone}</p>
    </article>
  );
}`}
      />

      <h2>4. Creating a contact from a public form</h2>
      <p>
        The shared Zod schema becomes the form&apos;s validator. One
        source of truth — change a field requirement in Go, regenerate
        with <code>grit sync</code>, and the form&apos;s validation
        catches up automatically.
      </p>

      <CodeBlock
        language="tsx"
        filename="apps/web/app/contacts/new/page.tsx"
        code={`"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { CreateContactSchema, type CreateContactInput } from "@repo/shared/schemas";
import { useCreateContact } from "@/hooks/use-contacts";

export default function NewContactPage() {
  const router = useRouter();
  const { mutate: createContact, isPending, error: apiError } = useCreateContact();
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreateContactInput>({
    resolver: zodResolver(CreateContactSchema),
  });

  const onSubmit = (input: CreateContactInput) => {
    createContact(input, {
      onSuccess: (created) => router.push("/contacts/" + created.id),
    });
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <label>
        <span>Name</span>
        <input {...register("name")} className="block w-full rounded-lg border px-3 py-2" />
        {errors.name && <p className="text-red-500 text-sm">{errors.name.message}</p>}
      </label>

      <label>
        <span>Email</span>
        <input type="email" {...register("email")} className="block w-full rounded-lg border px-3 py-2" />
        {errors.email && <p className="text-red-500 text-sm">{errors.email.message}</p>}
      </label>

      <label>
        <span>Phone</span>
        <input {...register("phone")} className="block w-full rounded-lg border px-3 py-2" />
      </label>

      <button
        type="submit"
        disabled={isPending}
        className="rounded-lg bg-accent px-4 py-2 text-white disabled:opacity-50"
      >
        {isPending ? "Creating…" : "Create contact"}
      </button>

      {apiError && (
        <p className="text-red-500">
          {(apiError as { response?: { data?: { error?: { message?: string } } } })
            ?.response?.data?.error?.message ?? "Something went wrong"}
        </p>
      )}
    </form>
  );
}`}
      />

      <TipBox tone="success">
        Notice: zero hand-written types. <code>CreateContactInput</code>{' '}
        and the validator both come from{' '}
        <code>@repo/shared/schemas</code>. If you add a field to Go and
        run <code>grit sync</code>, the form starts demanding it. No
        manual TS drift.
      </TipBox>

      <h2>5. Updating + deleting</h2>

      <CodeBlock
        language="tsx"
        code={`import { useUpdateContact, useDeleteContact } from "@/hooks/use-contacts";

const { mutate: updateContact } = useUpdateContact();
const { mutate: deleteContact } = useDeleteContact();

// Update
updateContact({ id: "01HX…", input: { name: "New name" } });

// Delete (asks the API to soft-delete — soft delete is the default in Grit)
deleteContact("01HX…");`}
      />

      <p>
        Both mutations call <code>queryClient.invalidateQueries({`{ queryKey: ['contacts'] }`})</code>{' '}
        on success — so any list pages currently rendered re-fetch and
        repaint. No manual cache management needed.
      </p>

      <h2>The auto-generated hook file in full</h2>
      <p>
        Curious what&apos;s actually inside{' '}
        <code>apps/web/hooks/use-contacts.ts</code>? Roughly this:
      </p>

      <CodeBlock
        language="ts"
        filename="apps/web/hooks/use-contacts.ts (auto-generated)"
        code={`import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import type { Contact } from "@repo/shared/types";

const KEY = ["contacts"] as const;

type ListParams = { page?: number; page_size?: number; search?: string };
type ListResponse = { data: Contact[]; meta: { total: number; page: number; page_size: number; pages: number } };

export function useContacts(params?: ListParams) {
  return useQuery({
    queryKey: [...KEY, params],
    queryFn: async () => {
      const { data } = await apiClient.get<ListResponse>("/api/contacts", { params });
      return data;
    },
  });
}

export function useContact(id: string) {
  return useQuery({
    queryKey: [...KEY, id],
    queryFn: async () => {
      const { data } = await apiClient.get<{ data: Contact }>("/api/contacts/" + id);
      return data.data;
    },
    enabled: Boolean(id),
  });
}

export function useCreateContact() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (input: Partial<Contact>) => {
      const { data } = await apiClient.post<{ data: Contact }>("/api/contacts", input);
      return data.data;
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: KEY }),
  });
}

export function useUpdateContact() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async ({ id, input }: { id: string; input: Partial<Contact> }) => {
      const { data } = await apiClient.put<{ data: Contact }>("/api/contacts/" + id, input);
      return data.data;
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: KEY }),
  });
}

export function useDeleteContact() {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (id: string) => {
      await apiClient.delete("/api/contacts/" + id);
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: KEY }),
  });
}`}
      />

      <h2>Wait — what about auth?</h2>
      <p>
        The customer-facing web app might run public pages too. The
        generated routes default to <code>middleware.Auth(...)</code>{' '}
        — they require a logged-in user. Two ways to expose a resource
        publicly:
      </p>
      <ul>
        <li>
          <strong>Public read, gated write.</strong> Edit{' '}
          <code>apps/api/internal/routes/routes.go</code> and move the{' '}
          <code>GET</code> routes out of the auth-protected group while
          keeping POST/PUT/DELETE inside it.
        </li>
        <li>
          <strong>Different resource entirely.</strong> Generate a
          public version that filters server-side (e.g. only{' '}
          <code>is_public = true</code>) and keep the auth&apos;d
          version for staff.
        </li>
      </ul>

      <KnowledgeCheck
        question="You added a `salutation` field to Contact in Go, ran `grit migrate` + `grit sync`. The web form using CreateContactSchema doesn't ask for it. What happened?"
        choices={[
          {
            label: 'sync only updates the admin form, not the web',
            feedback:
              "Wrong — sync regenerates both the shared schema (used by web) AND the TS type. Restart your Next.js dev server.",
          },
          {
            label: 'You need to restart `next dev` so the regenerated schema is picked up',
            correct: true,
            feedback:
              "Right. Next.js dev server caches module exports. After grit sync edits packages/shared/schemas/contact.ts, restarting next dev (or saving the importing file to trigger HMR) refreshes the Zod schema.",
          },
          {
            label: 'CreateContactSchema is hand-written and grit sync left it alone',
            feedback:
              "Wrong — for resources you generated with grit, the schema IS regenerated. (User is special-cased and skipped.)",
          },
          {
            label: 'The web app uses a different schema in @repo/shared/schemas/contact-create.ts',
            feedback:
              "Wrong — there's only one schema file per resource. Both admin and web import the same source of truth.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              In your contact-app, build two pages in the web app that
              consume the generated Contact resource:
            </p>
            <ol>
              <li>
                <code>apps/web/app/contacts/page.tsx</code> — list all
                contacts, with a search box.
              </li>
              <li>
                <code>apps/web/app/contacts/new/page.tsx</code> —
                create form using <code>CreateContactSchema</code>{' '}
                from <code>@repo/shared/schemas</code>.
              </li>
            </ol>
            <p>
              Open <code>http://localhost:3000/contacts</code>, create a
              contact via the form, and confirm it appears in the list
              (and in the admin panel — both apps share the API).
            </p>
          </>
        }
        hint={
          <>
            The web app uses HttpOnly cookie auth too — you&apos;ll need
            to log in via{' '}
            <code>http://localhost:3000/login</code> (or the admin&apos;s
            login page) first, since the generated routes default to
            auth-required. To make them public, edit{' '}
            <code>apps/api/internal/routes/routes.go</code> and move
            the contacts <code>GET</code> routes out of the auth
            group.
          </>
        }
        solution={
          <>
            <p>
              The list page is just <code>useContacts()</code> + a map
              over <code>data.data</code>. The form is{' '}
              <code>useCreateContact()</code> +{' '}
              <code>zodResolver(CreateContactSchema)</code>. If the
              creates fail with 401, you forgot to either log in or
              make the route public.
            </p>
            <p>
              The key insight from this exercise: one resource
              generation gave you a working admin panel <em>and</em> a
              working web app, with shared types, shared validators,
              and zero duplicated code.
            </p>
          </>
        }
      />

      <h2>You&apos;ve finished Chapter 4</h2>
      <p>
        Eight lessons in. You can now generate, sync, customise, and
        remove resources, model relationships across all three
        cardinalities, pick between short and long form, and consume
        the generated API from both the admin and the customer web
        app. The rest of Grit assumes this fluency — every lesson from
        here builds on top of <em>resources</em>.
      </p>
    </>
  )
}
