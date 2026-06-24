import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Real apps have entities that point at each other. A User has one
        Profile; a Contact belongs to a Group; a Post has many Tags.
        This lesson covers all three relationship cardinalities Grit
        supports — <strong>one-to-one</strong>,{' '}
        <strong>one-to-many</strong>, and{' '}
        <strong>many-to-many</strong> — using just two field types:{' '}
        <code>belongs_to</code> (with or without a unique constraint)
        and <code>many_to_many</code>. Each section is a complete,
        runnable example you can drop into a fresh Grit project.
      </p>

      <h2>The three cardinalities at a glance</h2>
      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Cardinality</th>
              <th className="text-left px-3 py-2 font-medium">Example</th>
              <th className="text-left px-3 py-2 font-medium">Field spec</th>
              <th className="text-left px-3 py-2 font-medium">DB shape</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr>
              <td className="px-3 py-2 text-[12px]">one-to-one</td>
              <td className="px-3 py-2 text-[12px]">User ↔ Profile</td>
              <td className="px-3 py-2 font-mono text-[12px]">belongs_to + unique</td>
              <td className="px-3 py-2 text-[12px]">FK on child with <code>uniqueIndex</code></td>
            </tr>
            <tr>
              <td className="px-3 py-2 text-[12px]">one-to-many</td>
              <td className="px-3 py-2 text-[12px]">Group → Contacts</td>
              <td className="px-3 py-2 font-mono text-[12px]">belongs_to</td>
              <td className="px-3 py-2 text-[12px]">FK on child, regular <code>index</code></td>
            </tr>
            <tr>
              <td className="px-3 py-2 text-[12px]">many-to-many</td>
              <td className="px-3 py-2 text-[12px]">Post ↔ Tags</td>
              <td className="px-3 py-2 font-mono text-[12px]">many_to_many</td>
              <td className="px-3 py-2 text-[12px]">Auto-generated join table</td>
            </tr>
          </tbody>
        </table>
      </div>

      <h2>Example 1 — One-to-many: Group → many Contacts</h2>
      <p>
        Imagine a CRM. Every Contact belongs to exactly one Group
        (Clients, Leads, Vendors), and every Group has many Contacts.
        That&apos;s the textbook one-to-many — modelled in Grit with{' '}
        <strong>one <code>belongs_to</code> field on the child</strong>{' '}
        (Contact), not two declarations.
      </p>

      <CodeBlock
        language="text"
        code={`Group           Contact
─────           ───────────────
id (UUID)  ←─── group_id (FK to Group)
name            name
description     email
                phone`}
      />

      <h3>Step 1 — Generate the parent first (Group)</h3>
      <p>
        Order matters: generate Group <em>before</em> Contact, so the
        Contact resource&apos;s <code>belongs_to</code> field has
        something to point at.
      </p>

      <CodeBlock
        terminal
        code={`grit generate resource Group \\
  --fields "name:string:unique,description:text:optional"

grit migrate`}
      />

      <p>That writes the same eight files you saw in Lesson 3 — model, service, handler, schema, type, hook, admin definition, admin page — plus injections into routes.go and the registry. Now we add the child.</p>

      <h3>Step 2 — Generate the child with belongs_to</h3>
      <CodeBlock
        terminal
        code={`grit generate resource Contact \\
  --fields "name:string,email:string:unique,phone:string:optional,group:belongs_to:Group"

grit migrate`}
      />

      <p>
        Three colons in the relationship spec, each meaningful:
      </p>

      <CodeBlock
        language="text"
        code={`group  :  belongs_to  :  Group
└──┬─┘    └────┬────┘    └──┬──┘
   │           │            │
   │           │            └── Related model. PascalCase, singular. Must already exist.
   │           │
   │           └── Type. Tells Grit to make this a foreign key, not a plain string.
   │
   └── Field name. Convention: name it after the relation, NOT "group_id".
       The generator appends "_id" automatically when writing the column.`}
      />

      <TipBox tone="info">
        Shortcut: if the field name matches the model name lowercased,
        you can drop the third part. <code>group:belongs_to</code> alone
        infers <code>:Group</code>. Be explicit when the field name
        diverges (e.g. <code>parent:belongs_to:Group</code>).
      </TipBox>

      <h3>What the generator produced — Contact side</h3>
      <p>
        Open <code>apps/api/internal/models/contact.go</code> after
        regenerating and you&apos;ll see a new column plus an eager-load
        association:
      </p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/models/contact.go (excerpt)"
        code={`type Contact struct {
	ID      string ` + '`' + `gorm:"primarykey;size:36" json:"id"` + '`' + `
	Name    string ` + '`' + `gorm:"size:255" json:"name" binding:"required"` + '`' + `
	Email   string ` + '`' + `gorm:"size:255;uniqueIndex" json:"email" binding:"required"` + '`' + `
	Phone   string ` + '`' + `gorm:"size:255" json:"phone"` + '`' + `

	// Foreign key column — what GORM stores in the contacts table.
	GroupID string ` + '`' + `gorm:"size:36;index" json:"group_id" binding:"required"` + '`' + `

	// Association — populated when you call .Preload("Group").
	Group   *Group ` + '`' + `json:"group,omitempty"` + '`' + `

	Version   int            ` + '`' + `gorm:"not null;default:1" json:"version"` + '`' + `
	CreatedAt time.Time      ` + '`' + `json:"created_at"` + '`' + `
	UpdatedAt time.Time      ` + '`' + `json:"updated_at"` + '`' + `
	DeletedAt gorm.DeletedAt ` + '`' + `gorm:"index" json:"-"` + '`' + `
}`}
      />

      <p>Two fields for the same relationship — that&apos;s the GORM idiom:</p>
      <ul>
        <li>
          <strong><code>GroupID string</code></strong> — the actual
          column on the <code>contacts</code> table. UUID string,
          indexed. This is what Create/Update endpoints accept and what
          the database stores.
        </li>
        <li>
          <strong><code>Group *Group</code></strong> — the populated
          association. Empty by default; only loaded when the service
          does <code>db.Preload(&quot;Group&quot;)</code>. The
          generator&apos;s <code>List</code> and <code>GetByID</code>{' '}
          methods preload it for you.
        </li>
      </ul>

      <h3>The frontend types stay clean</h3>
      <CodeBlock
        language="ts"
        filename="packages/shared/src/types/contact.ts"
        code={`export interface Contact {
  id: string;
  name: string;
  email: string;
  phone: string;
  group_id: string;
  group?: Group;          // present when preloaded by the API
  version: number;
  created_at: string;
  updated_at: string;
}`}
      />

      <h3>The admin form gets a dropdown</h3>
      <p>
        Open the regenerated{' '}
        <code>apps/admin/resources/contacts.ts</code> and you&apos;ll
        see:
      </p>
      <CodeBlock
        language="ts"
        filename="apps/admin/resources/contacts.ts (excerpt)"
        code={`form: {
  fields: [
    { key: "name",  type: "text", label: "Name",  required: true },
    { key: "email", type: "text", label: "Email", required: true },
    { key: "phone", type: "text", label: "Phone" },
    {
      key: "group_id",
      type: "relationship-select",
      label: "Group",
      relatedEndpoint: "/api/groups",
      displayField: "name",
      relationshipKey: "group_id",
      required: true,
    },
  ],
},`}
      />
      <p>
        The <code>relationship-select</code> field is a server-backed
        dropdown — it hits the <code>relatedEndpoint</code>, lists
        each row by the <code>displayField</code>, and writes the
        selected row&apos;s UUID into the column named in{' '}
        <code>relationshipKey</code>. No glue needed.
      </p>

      <h3>The list page shows the parent name</h3>
      <CodeBlock
        language="ts"
        filename="apps/admin/resources/contacts.ts (excerpt)"
        code={`table: {
  columns: [
    { key: "name",       label: "Name",    sortable: true, searchable: true },
    { key: "email",      label: "Email",   sortable: true, format: "email" },
    { key: "group.name", label: "Group" }, // dotted path reads from the preloaded association
    { key: "created_at", label: "Created", format: "relative" },
  ],
},`}
      />

      <h3>Showing &quot;has many&quot; on the parent — a manual addition</h3>
      <p>
        The generator only writes the <em>child</em> side of{' '}
        <code>belongs_to</code>. Group doesn&apos;t automatically know
        about its contacts. That&apos;s deliberate — listing contacts on
        the Group is usually a separate API call (paginated, filtered,
        searchable), not a preload.
      </p>
      <p>
        Two ways to add the &quot;has many&quot; side:
      </p>

      <h3>Option A — Add a <code>Contacts</code> association to Group manually</h3>
      <p>
        If you want <code>GET /api/groups/:id</code> to embed the
        contact list, add the field by hand:
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/models/group.go"
        code={`type Group struct {
	ID          string ` + '`' + `gorm:"primarykey;size:36" json:"id"` + '`' + `
	Name        string ` + '`' + `gorm:"size:255;uniqueIndex" json:"name"` + '`' + `
	Description string ` + '`' + `gorm:"type:text" json:"description"` + '`' + `

	// Added manually — has-many. GORM infers the FK is "group_id" on Contact.
	Contacts []Contact ` + '`' + `json:"contacts,omitempty"` + '`' + `

	// auto-generated boilerplate stays the same
}`}
      />
      <p>
        Then in the service for <code>GetByID</code>:{' '}
        <code>db.Preload(&quot;Contacts&quot;).First(&amp;group, &quot;id = ?&quot;, id)</code>.
        Run <code>grit sync</code> and the TS type picks up{' '}
        <code>contacts?: Contact[]</code>.
      </p>

      <h3>Option B — Filter the existing endpoint (recommended)</h3>
      <p>
        Use the auto-generated <code>GET /api/contacts?group_id=:id</code>{' '}
        endpoint. The generator wires a search/filter clause on every
        column, including FKs:
      </p>
      <CodeBlock
        terminal
        code={`curl "http://localhost:8080/api/contacts?group_id=01HX...&page=1&page_size=20" \\
  -H "Authorization: Bearer $TOKEN"`}
      />
      <p>
        Cleaner separation — Group endpoints handle Groups, Contact
        endpoints handle Contacts. Easier to paginate, easier to cache.
      </p>

      <h2>Example 2 — One-to-one: User ↔ Profile</h2>
      <p>
        One-to-one is just one-to-many <em>with a unique constraint</em>
        on the foreign key. Pick the side that&apos;s &quot;optional /
        added later&quot; (Profile, Settings, KycRecord) as the child
        and put the FK there. Each Profile points at one User; the
        unique index guarantees no User has more than one Profile.
      </p>

      <p>
        Because the inline <code>--fields</code> syntax doesn&apos;t
        accept <code>:unique</code> on relationship fields, this case is
        a great fit for the YAML long-form (covered in the next lesson):
      </p>

      <CodeBlock
        language="yaml"
        filename="profile.yaml"
        code={`name: Profile
fields:
  - name: user
    type: belongs_to
    related_model: User
    required: true
    unique: true              # ← this is what turns it into one-to-one
  - name: bio
    type: text
  - name: avatar
    type: string              # URL field — auto-becomes VARCHAR(500)
  - name: twitter_handle
    type: string
    required: false`}
      />
      <CodeBlock
        terminal
        code={`grit generate resource Profile --from profile.yaml
grit migrate`}
      />

      <p>The generated Go model is almost identical to the one-to-many case — but the GORM tag is different:</p>

      <CodeBlock
        language="go"
        filename="apps/api/internal/models/profile.go (excerpt)"
        code={`type Profile struct {
	ID     string ` + '`' + `gorm:"primarykey;size:36" json:"id"` + '`' + `
	Bio    string ` + '`' + `gorm:"type:text" json:"bio"` + '`' + `
	Avatar string ` + '`' + `gorm:"size:500" json:"avatar"` + '`' + `

	// Note the uniqueIndex — that's what makes it 1-to-1.
	UserID string ` + '`' + `gorm:"size:36;index;uniqueIndex" json:"user_id" binding:"required"` + '`' + `
	User   *User  ` + '`' + `json:"user,omitempty"` + '`' + `

	// auto-generated boilerplate
}`}
      />

      <p>
        Try inserting two profiles for the same user — the second insert
        fails with a <code>UNIQUE constraint failed</code> error.
        That&apos;s the database guaranteeing the relationship rather
        than your app code having to check.
      </p>

      <TipBox tone="info">
        <strong>Which side gets the FK?</strong> The side that&apos;s{' '}
        <em>optional</em> or <em>added later</em>. A User is created on
        sign-up; a Profile is filled in afterwards — so Profile holds
        the FK. Same logic for User → Settings, Order → Receipt, etc.
        The other side (User) doesn&apos;t need any changes — it stays
        oblivious to the optional companion.
      </TipBox>

      <h3>Loading the pair</h3>
      <p>
        To get a User with their Profile embedded, query from the
        Profile side and preload, or add a <code>HasOne</code>{' '}
        association on User manually:
      </p>
      <CodeBlock
        language="go"
        code={`// In services/user.go — add a method that joins the optional Profile.
func (s *UserService) GetWithProfile(id string) (*models.User, error) {
	var u models.User
	err := s.db.Preload("Profile").First(&u, "id = ?", id).Error
	return &u, err
}`}
      />
      <p>
        And add the association field to <code>models/user.go</code>:
      </p>
      <CodeBlock
        language="go"
        code={`type User struct {
	// ...existing fields...
	Profile *Profile ` + '`' + `gorm:"foreignKey:UserID" json:"profile,omitempty"` + '`' + `
}`}
      />
      <p>
        Run <code>grit sync</code> after the manual addition — the TS
        type picks up the optional <code>profile?: Profile</code>.
      </p>

      <h2>Example 3 — Many-to-many: Post ↔ Tags</h2>
      <p>
        Use <code>many_to_many</code> when an entity has <em>many</em>{' '}
        related items <em>and</em> each related item can belong to{' '}
        <em>many</em> of these entities. The classic case: a Post has
        many Tags, and each Tag is on many Posts.
      </p>

      <CodeBlock
        terminal
        code={`grit generate resource Tag \\
  --fields "name:string:unique,slug:slug:name,color:string:optional"

grit generate resource Post \\
  --fields "title:string,slug:slug,body:richtext,tags:many_to_many:Tag"

grit migrate`}
      />

      <p>The field spec works the same:</p>
      <CodeBlock
        language="text"
        code={`tags  :  many_to_many  :  Tag
└─┬─┘   └─────┬──────┘   └─┬─┘
  │           │            │
  │           │            └── Related model name. REQUIRED for many_to_many (unlike belongs_to).
  │           │
  │           └── Type. Creates a join table behind the scenes.
  │
  └── Field name. Usually plural. Becomes a []string of UUIDs in Go.`}
      />

      <p>What you get:</p>
      <ul>
        <li>
          A join table named <code>post_tags</code> with{' '}
          <code>post_id</code> + <code>tag_id</code> columns, indexed
          both ways.
        </li>
        <li>
          <code>Post.Tags</code> is <code>[]string</code> on the Go side
          (the array of related UUIDs) and{' '}
          <code>string[]</code> in TypeScript.
        </li>
        <li>
          The admin Form gets a{' '}
          <code>multi-relationship-select</code> — a multi-select
          dropdown backed by <code>GET /api/tags</code>. The list
          page renders attached tags as a stack of colored pills
          using the badge format.
        </li>
      </ul>

      <h2>Picking the right relationship</h2>
      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">If…</th>
              <th className="text-left px-3 py-2 font-medium">Use</th>
              <th className="text-left px-3 py-2 font-medium">Where to put it</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr>
              <td className="px-3 py-2 text-[12px]">Each parent has at most one child (User ↔ Profile, Order ↔ Receipt)</td>
              <td className="px-3 py-2 font-mono text-[12px]">belongs_to + unique</td>
              <td className="px-3 py-2 text-[12px]">On the &quot;optional&quot; side — uses YAML (inline can&apos;t set unique on belongs_to)</td>
            </tr>
            <tr>
              <td className="px-3 py-2 text-[12px]">Child has exactly one parent (Contact → Group, Order → Customer)</td>
              <td className="px-3 py-2 font-mono text-[12px]">belongs_to</td>
              <td className="px-3 py-2 text-[12px]">On the child</td>
            </tr>
            <tr>
              <td className="px-3 py-2 text-[12px]">Parent owns N children, children are <em>not</em> shared (BlogPost → Comment)</td>
              <td className="px-3 py-2 font-mono text-[12px]">belongs_to on the child</td>
              <td className="px-3 py-2 text-[12px]">On the child + filter <code>?post_id=</code></td>
            </tr>
            <tr>
              <td className="px-3 py-2 text-[12px]">Both sides can have many of each other (Post ↔ Tag, User ↔ Role)</td>
              <td className="px-3 py-2 font-mono text-[12px]">many_to_many</td>
              <td className="px-3 py-2 text-[12px]">On whichever side reads it more — usually both</td>
            </tr>
            <tr>
              <td className="px-3 py-2 text-[12px]">Loose freeform list (5–10 tags users type by hand, not managed centrally)</td>
              <td className="px-3 py-2 font-mono text-[12px]">string_array</td>
              <td className="px-3 py-2 text-[12px]">A column on the parent (no join table)</td>
            </tr>
          </tbody>
        </table>
      </div>

      <KnowledgeCheck
        question="You're modelling a school app. A Student has one HomeRoom; a Student takes many Classes; each Class has many Students. How do you wire it?"
        choices={[
          {
            label: 'home_room:belongs_to:HomeRoom on Student, classes:many_to_many:Class on Student',
            correct: true,
            feedback:
              "Right. Student → HomeRoom is one-to-many so it's a belongs_to on Student. Student ↔ Class is many-to-many on both sides — declare it on one (usually Student) and the join table covers the reverse.",
          },
          {
            label: 'home_room:many_to_many:HomeRoom and classes:belongs_to:Class — both on Student',
            feedback:
              'Backwards. HomeRoom is a single relation (belongs_to). Classes is multi (many_to_many).',
          },
          {
            label: 'Skip relationships — store home_room_name as a string and classes as a string_array',
            feedback:
              'Tempting but fragile. Renaming a HomeRoom would require an update sweep, and you lose the ability to list "all students in this class" from the Class side.',
          },
          {
            label: 'belongs_to for both — Student has one HomeRoom and belongs_to one Class',
            feedback:
              'Wrong for Class — a Student can be in several Classes. belongs_to forces "exactly one".',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Try all three cardinalities on your machine. Each part is
              self-contained — feel free to do them one at a time.
            </p>

            <p><strong>Part A — one-to-many (Group → Contacts)</strong></p>
            <CodeBlock
              terminal
              code={`grit generate resource Group \\
  --fields "name:string:unique,description:text:optional"

grit generate resource Contact \\
  --fields "name:string,email:string:unique,phone:string:optional,group:belongs_to:Group"

grit migrate`}
            />
            <p>
              In the admin, create one Group then two Contacts pointing
              at it. Hit{' '}
              <code>GET /api/contacts?group_id=&lt;the-group-uuid&gt;</code>{' '}
              — both contacts come back, each with the embedded{' '}
              <code>group</code> object.
            </p>

            <p><strong>Part B — one-to-one (User → Profile)</strong></p>
            <p>
              Drop this YAML in <code>profile.yaml</code> at the project
              root and generate from it:
            </p>
            <CodeBlock
              language="yaml"
              filename="profile.yaml"
              code={`name: Profile
fields:
  - name: user
    type: belongs_to
    related_model: User
    required: true
    unique: true
  - name: bio
    type: text
  - name: avatar
    type: string`}
            />
            <CodeBlock
              terminal
              code={`grit generate resource Profile --from profile.yaml
grit migrate`}
            />
            <p>
              Create one Profile for your admin user. Then try to create
              a second Profile for the same user — the API should return
              a 500 with a unique-constraint error.
            </p>

            <p><strong>Part C — many-to-many (Post ↔ Tag)</strong></p>
            <CodeBlock
              terminal
              code={`grit generate resource Tag \\
  --fields "name:string:unique,slug:slug:name,color:string:optional"

grit generate resource Post \\
  --fields "title:string,slug:slug,body:richtext,tags:many_to_many:Tag"

grit migrate`}
            />
            <p>
              In the admin, create 3 Tags then a Post with all three
              attached via the multi-select. Hit{' '}
              <code>GET /api/posts/&lt;id&gt;</code> — the{' '}
              <code>tags</code> field comes back as an array of UUIDs
              (or full Tag objects if you tweak the service to preload).
            </p>

            <p>
              Paste the three API responses into <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            The Group / Tag UUIDs are visible in the admin panel URL
            after you create them (<code>/resources/groups/&lt;id&gt;</code>),
            or on the list page. If <code>grit migrate</code> errors with
            &quot;table already exists&quot;, that means a previous
            generate created the same resource — rename or drop it
            first.
          </>
        }
        solution={
          <>
            <p>Part A — paginated list with embedded parent object:</p>
            <CodeBlock
              language="json"
              code={`{
  "data": [
    {
      "id": "01HXP…",
      "name": "Alice",
      "email": "alice@example.com",
      "group_id": "01HXG…",
      "group": { "id": "01HXG…", "name": "Clients" }
    },
    { "id": "01HXP…", "name": "Bob", "group_id": "01HXG…", "group": {…} }
  ],
  "meta": { "total": 2, "page": 1 }
}`}
            />
            <p>Part B — the second insert returns:</p>
            <CodeBlock
              language="json"
              code={`{
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "creating profile: UNIQUE constraint failed: profiles.user_id"
  }
}`}
            />
            <p>
              The database — not your app code — is the one enforcing
              one-to-one. That&apos;s exactly what you want.
            </p>
            <p>Part C — Post with attached tags:</p>
            <CodeBlock
              language="json"
              code={`{
  "data": {
    "id": "01HXP…",
    "title": "Getting started with Grit",
    "slug": "getting-started-with-grit-a8f3",
    "body": "<p>…</p>",
    "tags": ["01HXT-uuid-1…", "01HXT-uuid-2…", "01HXT-uuid-3…"]
  }
}`}
            />
            <p>
              You&apos;ve now exercised every relationship cardinality
              the Grit generator supports — with two field types
              (<code>belongs_to</code>, <code>many_to_many</code>) and
              one modifier (<code>unique</code>) that turns one-to-many
              into one-to-one.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        You&apos;ve seen the inline form for every field type Grit
        supports. The next lesson is about choosing between the inline{' '}
        <code>--fields</code> string and the long-form YAML file — when
        each shines, and what extra knobs YAML gives you (defaults, the
        single source of truth you can re-run, anything beyond five
        fields).
      </p>
    </>
  )
}
