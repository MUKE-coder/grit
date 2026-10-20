#import "../template.typ": *

= The shared package: one definition, four outputs

D3 was "the same type exists four times". A `Blog` is a Go struct, a database
column list, a TypeScript interface and a validation schema, and nothing makes
them agree except the person who remembers to change all four.

This chapter is about the place where three of those four stop being separate
things, and about the one command that keeps them that way.

== What is in it

```
packages/shared/
├── package.json
├── types/          TypeScript interfaces, one file per model
├── schemas/        Zod schemas, one file per model
├── constants/      API routes and shared constants
├── themes.ts       the palette both frontends read
└── brand.config.ts the product name, logo and copy
```

It is a workspace package rather than a published one, so there is no build
step and no version to bump. The frontends import it by subpath:

```json
"exports": {
  "./schemas":   "./schemas/index.ts",
  "./types":     "./types/index.ts",
  "./constants": "./constants/index.ts",
  "./brand":     "./brand.config.ts",
  "./themes":    "./themes.ts"
}
```

```tsx
import { CreateBlogSchema } from "@repo/shared/schemas"
import type { Blog } from "@repo/shared/types"
```

Everything that renders imports from here: the web app, the admin panel, and
the desktop and mobile clients if the project has them. One definition, four
consumers.

== A generated pair, in full

Given this model:

```go
type Blog struct {
    ID          string     `json:"id"          gorm:"primaryKey"`
    Title       string     `json:"title"       binding:"required"`
    Slug        string     `json:"slug"`
    Content     string     `json:"content"     gorm:"type:text"`
    Image       *string    `json:"image"`
    Published   bool       `json:"published"`
    PublishedAt *time.Time `json:"published_at"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}
```

you get `types/blog.ts`:

```ts
export interface Blog {
  id: string;
  title: string;
  slug: string;
  content: string;
  image: string | null;
  published: boolean;
  published_at: string | null;
  created_at: string;
  updated_at: string;
}
```

and `schemas/blog.ts`, which is three schemas rather than one:

```ts
export const BlogSchema = z.object({
  id: z.string(),
  title: z.string(),
  slug: z.string(),
  content: z.string(),
  image: z.string().nullable(),
  published: z.boolean(),
  published_at: z.string().nullable(),
  created_at: z.string(),
  updated_at: z.string(),
});

export const CreateBlogSchema = z.object({
  title: z.string().min(1, "Title is required"),
  content: z.string().optional(),
  image: z.string().optional(),
  published: z.boolean().optional(),
});

export const UpdateBlogSchema = CreateBlogSchema.partial();
```

== Why three, and what each one is for

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`BlogSchema`], [The row as the API returns it. Every field, including the
    ones the server assigns],
  [`CreateBlogSchema`], [What a `POST` accepts. No `id`, no timestamps, and
    `required` where the Go binding tag said so],
  [`UpdateBlogSchema`], [The create schema, `.partial()`. A `PATCH` sends the
    fields that changed and nothing else],
)

`id`, `created_at`, `updated_at` and `deleted_at` are dropped from the create
schema automatically. They are assigned by the server, so a form that offers
them is a form that can lie.

#bridge[
  This is the shape you would have built by hand in a Next.js project, and the
  part that goes wrong by hand is `UpdateSchema`.

  People reuse the create schema for editing, so every `PATCH` has to send
  every field, so an edit form that was missing one input quietly blanks that
  column. `.partial()` is the fix, and having it generated means you cannot
  forget it on the eleventh resource.
]

== nullable against optional, which is not pedantry

`image` is `z.string().nullable()` in the row schema and `z.string()
.optional()` in the create schema, and the difference is the difference between
two states that look the same in a form:

- *nullable*: the server will send `"image": null`. The key is there, the value
  is empty.
- *optional*: the client may leave the key out entirely, and the server applies
  its default.

A Go pointer field (`*string`, `*time.Time`) becomes nullable, because a
pointer is how Go says "this may be absent". A non-required field becomes
optional on create. Get these the wrong way round by hand and you get a form
that refuses to submit because a column the user cannot see is `null`.

== grit sync, and exactly what it does

```bash
grit sync
```

It parses every file in `apps/api/internal/models/` with Go's own parser, not a
regular expression, and for each struct it writes:

#table(
  columns: (auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [`types/<name>.ts`], [The interface, overwritten],
  [`schemas/<name>.ts`], [The three schemas, overwritten],
  [the admin resource], [New fields added to the generated section only],
)

The name is kebab-cased from the struct: `BlogPost` becomes `blog-post.ts`.

`User` is skipped. It has hand-written schemas covering login, registration,
password reset and profile updates, which are not derivable from the struct,
and regenerating them would throw away work the framework cannot reproduce.

The mapping, for the types you will actually hit:

#table(
  columns: (auto, auto, 1fr),
  stroke: 0.4pt + rgb("#e3e6ea"),
  inset: 5pt,
  [*Go*], [*TypeScript*], [*Zod*],
  [`string`], [`string`], [`z.string()`],
  [`int`, `int64`], [`number`], [`z.number().int()`],
  [`uint`], [`number`], [`z.number().int().nonnegative()`],
  [`float64`], [`number`], [`z.number()`],
  [`bool`], [`boolean`], [`z.boolean()`],
  [`time.Time`], [`string`], [`z.string()`],
  [`*time.Time`], [`string | null`], [`z.string().nullable()`],
)

Times cross as strings, because JSON has no date type. Chapter nineteen covers
the format; the thing to hold onto here is that the TypeScript type says
`string` and means RFC 3339, and the parsing is the frontend's job.

#gotcha[
  *`grit sync` overwrites those two files completely.* There is no merge and no
  manifest conflict: the generated type and the generated schema are outputs,
  and an output is rewritten.

  So a refinement you add inside `schemas/blog.ts`, a `.refine()` for a
  cross-field rule or a stricter message, is gone the next time somebody adds a
  column. Put it in a file beside it and compose:

  ```ts
  // schemas/blog-form.ts, which sync does not touch
  export const BlogFormSchema = CreateBlogSchema.extend({
    slug: z.string().regex(/^[a-z0-9-]+$/, "lowercase and dashes"),
  }).refine((v) => !v.published || v.content.length > 0, {
    message: "a published post needs content",
    path: ["content"],
  })
  ```

  The rule for the whole framework, stated once: *edit the inputs, extend the
  outputs*. The model is an input. The shared schema is an output.
]

== The admin resource is updated, not overwritten

The third thing `sync` touches is the admin resource definition, and it is
treated differently because it is yours.

The generated file carries markers around the section the generator owns. New
model fields are added inside those markers, to the columns and the form.
Anything you customised, a renamed label, a `render` function, a reordered
column, sits outside them and survives.

So the usual loop after changing a model is two commands:

```bash
grit migrate
grit sync
```

and the admin has the new field in the table and the form without you opening
the file.

== Where the chain stops, honestly

`sync` keeps three things in step: the type, the schema, and the admin
resource. It does not update the generated *mobile* forms or the *desktop*
columns, because those have no marker scheme yet.

It says so rather than pretending:

```
  ⚠ 2 field(s) never reached apps/mobile:
    reading_notes, verified_at
```

That warning exists because the failure without it is miserable to diagnose: a
correct API, a correct type, a correct admin screen, and a mobile create form
frozen at the three fields the resource was generated with. The field is in the
response and absent from the form, and nothing is broken enough to show an
error.

#gotcha[
  Treat that warning as a task, not as noise. The fix is to open the named
  screen and add the field by hand, which is five minutes, or to regenerate the
  screen if you have not customised it.

  A framework that quietly half-synced would be worse than one that does not
  sync at all, because you would trust the half.
]

== Using it on the frontend

The schemas are not decoration. They are the validation layer of every
generated form:

```tsx
const form = useForm<CreateBlogInput>({
  resolver: zodResolver(CreateBlogSchema),
})
```

The inferred input type comes from the same schema:

```ts
export type CreateBlogInput = z.infer<typeof CreateBlogSchema>;
```

so renaming a field in the Go model and running `grit sync` makes the editor
underline every place the old name is used, before anything runs. That is the
actual payoff of D3: not that the four definitions exist in one place, but that
disagreeing with them is a compile error rather than a support ticket.

The error codes travel the same way. `types/errors.ts` is generated from the Go
error catalogue, so a frontend can match on a code rather than on a message:

```ts
if (apiErrorCode(err) === "EMAIL_EXISTS") {
  form.setError("email", { message: "That email is taken" })
}
```

== Constants, themes and brand

Three smaller files, each solving a small duplication that gets annoying at the
third app:

- `constants/` holds the API route paths, so a URL typo is a type error rather
  than a 404 at runtime.
- `themes.ts` holds the palette, read by the admin, the web app and the desktop
  client, so a theme change lands everywhere at once.
- `brand.config.ts` holds the product name, logo and marketing copy, which is
  why a generated project shows *your* name rather than Grit's.

#project(title: "Project 33: break the chain and watch it fail")[
  The point of this project is to feel what the single definition is worth, by
  removing it.

  ```bash
  grit generate resource Reading \
    --fields "label:string, value:float, taken_at:datetime"
  grit migrate
  ```

  *One: the normal loop.* Add a field to the Go model by hand:

  ```go
  Verified bool `json:"verified"`
  ```

  Run `grit migrate` then `grit sync`. Open `packages/shared/types/
  reading.ts`, the schema, and the admin resource file. Three files changed
  from one edit.

  *Two: skip the sync.* Add a second field and run only `grit migrate`. The
  column exists, the API returns it, and the TypeScript type does not have it.
  Try to read it in a page and watch the compiler stop you. That error is the
  framework working, not failing.

  *Three: the overwrite.* Add a `.refine()` inside the generated
  `schemas/reading.ts`, then run `grit sync` again. It is gone. Now move it to
  a sibling file as the gotcha above shows, and run sync a third time.

  *Four: the drift warning.* If the project has a mobile or desktop client, add
  a third field and read what sync prints. Open the named screen and confirm
  the field really is missing.

  *Five: the payoff.* Rename `label` to `name` in the Go model, migrate, sync,
  and then run `pnpm typecheck`. Every place that used the old name is now an
  error with a file and a line. Count them, and ask how you would have found
  them in a project where the type lived in two languages.
]

#tryit[
  Open `packages/shared/schemas/index.ts` and read the exports. Every generated
  resource adds three schemas and two inferred types to that barrel.

  Then ask the question that decides whether this package earns its place: if
  your application has one frontend and will never have another, is the shared
  package still worth it? The honest answer is that it is worth it the moment
  there is a second consumer of the type, and that a mobile app, a desktop
  client and an admin panel are each a second consumer.
]

#assignment(name: "Shared Package - Assignment One")[
  Make the single definition prove itself on your own project.

  *One.* Add two fields to an existing model, run `grit migrate` then
  `grit sync`, and commit the diff. In `NOTES.md`, list every file that changed
  from those two edits and mark which were generated and which were updated in
  place.

  *Two.* Write a cross-field rule that the generator cannot infer: a published
  post needs content, an end date must follow a start date, a discount cannot
  exceed the total. Put it in a sibling schema file, use it in a form, and then
  run `grit sync` and show that it survived.

  *Three.* Rename a field in the Go model. Migrate, sync, and run the
  TypeScript check. Paste the compiler output into `NOTES.md`: that list is the
  thing D3 buys you, and it is worth seeing written out.

  *Four.* One paragraph on the honest limit: which clients did `sync` not
  update, what did it tell you, and what would you have to do by hand?
]
