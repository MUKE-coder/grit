import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Lesson 2 listed thirteen field types in a table. This lesson is
        the field-by-field deep dive — the patterns you reach for once a
        plain <code>string</code> isn&apos;t enough. By the end you can
        generate an Article with a URL slug, a cover image, a YouTube
        video, a gallery of photos, and a list of tags — and know
        exactly why each one was modelled the way it was.
      </p>

      <h2><code>slug</code> — URL-friendly auto-IDs</h2>
      <p>
        A slug is a URL-safe identifier derived from another field. You
        want <code>/blog/my-first-post</code>, not{' '}
        <code>/blog/c8f5a93b-2401-4a7e-9d11</code>. <code>slug</code>{' '}
        gives you that without writing the slugifier yourself.
      </p>

      <h3>Auto-generated from the first string field</h3>
      <CodeBlock
        terminal
        code={`grit generate resource Article \\
  --fields "title:string,slug:slug,content:richtext"`}
      />
      <p>
        When you create an Article with{' '}
        <code>{`{ title: "My First Post" }`}</code>, the{' '}
        <code>BeforeCreate</code> hook fills in{' '}
        <code>slug = &quot;my-first-post-a8f3&quot;</code> — lowercased,
        non-alphanumerics replaced with hyphens, plus a 4-byte hex suffix
        so two posts with the same title don&apos;t collide.
      </p>

      <h3>Customise which field the slug comes from</h3>
      <p>
        Pass the source field as the third colon-separated part:
      </p>
      <CodeBlock
        terminal
        code={`grit generate resource Product \\
  --fields "sku:string:unique,name:string,slug:slug:sku"`}
      />
      <p>
        Now the slug is derived from <code>sku</code> instead of the
        first string column. (Without that hint the generator would have
        used <code>sku</code> anyway because it&apos;s the first string —
        but being explicit beats relying on field order.)
      </p>

      <h3>What the generated hook looks like</h3>
      <CodeBlock
        language="go"
        filename="apps/api/internal/models/article.go (excerpt)"
        code={`func (m *Article) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.Slug == "" {
		m.Slug = slugify(fmt.Sprintf("%v", m.Title))
	}
	return nil
}`}
      />

      <p>
        Two more facts worth remembering:
      </p>
      <ul>
        <li>
          The <code>slug</code> column gets a unique index automatically
          — no <code>:unique</code> modifier needed.
        </li>
        <li>
          Slug fields don&apos;t show up in the admin form (you can
          override by editing the resource page) — they&apos;re intended
          to be derived, not typed by hand.
        </li>
      </ul>

      <h2>Images, avatars, banners — string fields with smart sizing</h2>
      <p>
        Grit doesn&apos;t have an <code>image</code> field type. It has
        something better: it watches the <em>name</em> of a string field
        and upgrades the column to <code>VARCHAR(500)</code> automatically
        if it looks like a URL. So you just use{' '}
        <code>:string</code> and Grit handles the storage:
      </p>

      <CodeBlock
        terminal
        code={`grit generate resource Article \\
  --fields "title:string,cover:string,thumbnail:string,avatar:string"`}
      />

      <p>
        All four columns become <code>VARCHAR(500)</code> (signed S3 URLs
        and UTM-tagged image links eat 255 characters fast). The
        heuristic triggers on:
      </p>

      <CodeBlock
        language="text"
        code={`Exact match:  url, image, avatar, thumbnail, logo, cover, icon, banner, photo
Suffix:       anything_url  (image_url, profile_url, callback_url, …)`}
      />

      <h3>How the upload flow works</h3>
      <p>
        The field stores the <em>URL</em>; the file itself lives on S3
        (or MinIO / R2 / B2). Grit ships an upload handler at{' '}
        <code>POST /api/uploads</code> that:
      </p>
      <ol>
        <li>Receives a multipart file from the admin Form.</li>
        <li>Streams it to your bucket via the storage service.</li>
        <li>
          Returns <code>{`{ data: { url: "https://…/file.jpg" } }`}</code>.
        </li>
      </ol>
      <p>
        The admin <code>FormBuilder</code> spots fields named like
        images (<code>image</code>, <code>avatar</code>, <code>cover</code>,
        …) and renders a drag-and-drop uploader that POSTs to{' '}
        <code>/api/uploads</code> and writes the returned URL back into
        the form. Zero glue code.
      </p>

      <TipBox tone="info">
        Want the upload widget but the URL column is named differently
        (e.g., <code>headshot</code>)? Either rename it to a heuristic
        match (<code>photo</code>, <code>avatar</code>) or open the
        generated admin page and change the form field&apos;s{' '}
        <code>type: &quot;text&quot;</code> to{' '}
        <code>type: &quot;image&quot;</code> by hand.
      </TipBox>

      <h2>Videos — same trick, different field name</h2>
      <p>
        Grit doesn&apos;t care if the URL points at a JPEG or an MP4. A
        video field is just a string column holding the URL:
      </p>
      <CodeBlock
        terminal
        code={`grit generate resource Course \\
  --fields "title:string,video_url:string,duration_seconds:int"`}
      />
      <p>
        The <code>_url</code> suffix triggers the VARCHAR(500) upgrade.
        The frontend decides what to do with the URL — embed a{' '}
        <code>&lt;video&gt;</code> tag, load it into a player, or treat
        it as a YouTube/Vimeo embed URL.
      </p>

      <p>
        For uploaded video files, the same <code>/api/uploads</code>{' '}
        endpoint works — Grit&apos;s upload handler accepts any MIME type
        (and an env var caps the max size).
      </p>

      <h2><code>string_array</code> — galleries and tag lists</h2>
      <p>
        Two-for-one: same type, two common uses depending on whether
        you&apos;re storing URLs or freeform strings.
      </p>

      <h3>Use 1: photo gallery / screenshot list</h3>
      <CodeBlock
        terminal
        code={`grit generate resource Listing \\
  --fields "title:string,description:text,photos:string_array"`}
      />
      <p>
        The Go side becomes:
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/models/listing.go (excerpt)"
        code={`type Listing struct {
	ID          string                       ` + '`' + `gorm:"primarykey;size:36" json:"id"` + '`' + `
	Title       string                       ` + '`' + `gorm:"size:255" json:"title"` + '`' + `
	Description string                       ` + '`' + `gorm:"type:text" json:"description"` + '`' + `
	Photos      datatypes.JSONSlice[string]  ` + '`' + `gorm:"type:json" json:"photos"` + '`' + `
	// …
}`}
      />
      <p>
        Stored as a single JSON column (<code>{`["url1","url2","url3"]`}</code>),
        not a separate table. The TS type is{' '}
        <code>string[]</code>. And — this is the nice bit — the admin
        form renders <code>string_array</code> with a{' '}
        <strong>multi-file image uploader</strong> out of the box. Drag
        in five photos, the form POSTs each to <code>/api/uploads</code>,
        and the URLs land in the array.
      </p>

      <h3>Use 2: freeform tag list</h3>
      <CodeBlock
        terminal
        code={`grit generate resource Post \\
  --fields "title:string,body:richtext,tags:string_array"`}
      />
      <p>
        Same column type — just stores{' '}
        <code>{`["tutorial","go","react"]`}</code> instead of URLs. The
        admin form&apos;s image uploader is appropriate for galleries;
        for freeform tags you&apos;ll usually swap the form field type
        to a chips-input in the generated{' '}
        <code>page.tsx</code>. (Or use the <code>many_to_many</code>{' '}
        relationship covered in the next lesson, which gets you a proper
        Tag table with its own list page.)
      </p>

      <TipBox tone="info">
        <strong>When to pick which:</strong>{' '}
        <code>string_array</code> when the values are <em>strings the
        user types</em> (tags, photo URLs, keywords) and you don&apos;t
        need to query/list them as their own entities.{' '}
        <code>many_to_many:Tag</code> when tags need their own admin
        page, slug, color, usage count — when a Tag is a thing.
      </TipBox>

      <h2><code>text</code> vs <code>richtext</code> — when format matters</h2>
      <p>
        Both store long text in a <code>TEXT</code> column. The
        difference is the editor:
      </p>
      <ul>
        <li>
          <code>text</code> — plain textarea. Good for notes,
          internal-only descriptions, prompts you&apos;ll send to an LLM.
        </li>
        <li>
          <code>richtext</code> — Tiptap Word-style editor with
          bold/italic/underline, headings, lists, links, images, tables.
          Good for blog posts, knowledge-base articles, marketing copy.
        </li>
      </ul>
      <CodeBlock
        terminal
        code={`grit generate resource KnowledgeArticle \\
  --fields "title:string,slug:slug,summary:text,body:richtext"`}
      />

      <h2>Heuristic field names — free upgrades</h2>
      <p>
        Three name patterns that change column storage even though the
        type is plain <code>string</code> or <code>float</code>:
      </p>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Pattern (on type)</th>
              <th className="text-left px-3 py-2 font-medium">Column becomes</th>
              <th className="text-left px-3 py-2 font-medium">Examples</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">URL-shaped (string)</td>
              <td className="px-3 py-2 font-mono text-[12px]">VARCHAR(500)</td>
              <td className="px-3 py-2 font-mono text-[12px]">avatar, logo, photo, banner, *_url</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">Long-text-shaped (string)</td>
              <td className="px-3 py-2 font-mono text-[12px]">TEXT</td>
              <td className="px-3 py-2 font-mono text-[12px]">description, notes, content, body, summary, bio, message</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">Money-shaped (float)</td>
              <td className="px-3 py-2 font-mono text-[12px]">DECIMAL(12,2)</td>
              <td className="px-3 py-2 font-mono text-[12px]">price, amount, total, *_cost, *_fee, *_salary, *_balance</td>
            </tr>
          </tbody>
        </table>
      </div>

      <CodeBlock
        terminal
        code={`grit generate resource Invoice \\
  --fields "number:string:unique,amount:float,description:string,due_date:date"`}
      />
      <p>The generator quietly does the right thing for each column:</p>
      <ul>
        <li>
          <code>amount:float</code> → <code>DECIMAL(12,2)</code>{' '}
          (the name <code>amount</code> triggers the money heuristic).
        </li>
        <li>
          <code>description:string</code> → <code>TEXT</code>{' '}
          (long-text heuristic — even though you wrote{' '}
          <code>string</code>).
        </li>
        <li>
          <code>due_date:date</code> → <code>DATE</code> column, not{' '}
          <code>TIMESTAMP</code> (use <code>:datetime</code> if you need
          the time component).
        </li>
      </ul>

      <h2>Defaults — when you need them, switch to YAML</h2>
      <p>
        Defaults are <em>not</em> available in the inline{' '}
        <code>--fields</code> string (intentional — keeps the syntax
        copy-pasteable). When you need a default, use a YAML definition:
      </p>

      <CodeBlock
        language="yaml"
        filename="task.yaml"
        code={`name: Task
fields:
  - name: title
    type: string
    required: true
  - name: status
    type: string
    default: pending      # GORM "default:pending" → DB-side default
  - name: priority
    type: int
    default: 3
  - name: archived
    type: bool
    default: false`}
      />
      <CodeBlock
        terminal
        code={`grit generate resource Task --from task.yaml`}
      />

      <h2>Common-field cookbook</h2>
      <p>Five recipes you&apos;ll re-use across most resources:</p>

      <CodeBlock
        language="text"
        code={`# Blog post
title:string, slug:slug, excerpt:text, cover:string, body:richtext,
published:bool, published_at:datetime, tags:string_array

# Course / lesson
title:string, slug:slug, description:text, thumbnail:string,
video_url:string, duration_seconds:int, free_preview:bool

# E-commerce product
sku:string:unique, name:string, slug:slug:sku, description:richtext,
price:float, stock_quantity:int, photos:string_array, featured:bool

# Real-estate listing
title:string, slug:slug, address:string, city:string, state:string,
price:float, bedrooms:int, bathrooms:float, square_feet:int,
description:text, photos:string_array, listed_at:datetime

# Calendar event
title:string, location:string, starts_at:datetime, ends_at:datetime,
all_day:bool, notes:text, color:string`}
      />

      <KnowledgeCheck
        question="You're building a /portfolio resource for showcasing client projects. Each project has a title, a hero image, a body of marketing copy with formatting, and 4–10 screenshots. Best field spec?"
        choices={[
          {
            label: 'title:string, image:string, body:text, screenshots:string',
            feedback:
              'Two issues: body:text is plain text (no headings, bold, etc.), and a single string can\'t hold 4–10 screenshots.',
          },
          {
            label: 'title:string, hero:string, body:richtext, screenshots:string_array',
            correct: true,
            feedback:
              "Spot on. \"hero\" hits the URL heuristic → VARCHAR(500); richtext gets the Tiptap editor; string_array stores the multi-image gallery and the admin form will render a multi-file uploader.",
          },
          {
            label: 'title:string, image:string, body:richtext, screenshots:json',
            feedback:
              'Close — but Grit doesn\'t have a json field type. Use string_array for a list of strings, or many_to_many for a relationship.',
          },
          {
            label: 'title:string, image:image, body:richtext, screenshots:images',
            feedback:
              'Grit doesn\'t have :image or :images types. Use :string (heuristic name) for one URL and :string_array for a gallery.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Design a <code>Recipe</code> resource for a cooking app.
              Fields it needs:
            </p>
            <ul>
              <li>A title with a URL slug</li>
              <li>A cover photo</li>
              <li>Prep time in minutes (whole numbers)</li>
              <li>A formatted body with steps (bold, lists, headings)</li>
              <li>A list of 3–10 ingredient strings</li>
              <li>A &quot;published&quot; toggle</li>
            </ul>
            <p>
              Write the <code>grit generate resource Recipe --fields
              &quot;…&quot;</code> command and paste it into{' '}
              <code>notes.md</code>. Then actually run it on your
              project.
            </p>
          </>
        }
        hint={
          <>
            Six fields, all reachable with what you&apos;ve learnt: one
            slug from title, one URL-named string, one int (uint if you
            want non-negative), one richtext, one string_array, one bool.
          </>
        }
        solution={
          <>
            <CodeBlock
              terminal
              code={`grit generate resource Recipe \\
  --fields "title:string,slug:slug,cover:string,prep_minutes:int,body:richtext,ingredients:string_array,published:bool"

grit migrate`}
            />
            <p>
              After migration, the admin page at{' '}
              <code>/resources/recipes</code> renders with: a Title input,
              an auto-hidden slug, a single-image uploader for{' '}
              <code>cover</code>, a number input, a Tiptap editor, a
              multi-image uploader for <code>ingredients</code> (which
              you may want to swap to a chips input — see the tip above),
              and a toggle.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        You can now generate single-table resources with rich fields. The
        next lesson is what makes Grit feel like a framework instead of
        a script: <strong>relationships</strong>. We&apos;ll build a
        Contact + Group pair where each contact belongs to a group and
        every group lists its contacts — using <code>belongs_to</code>{' '}
        and <code>many_to_many</code>.
      </p>
    </>
  )
}
