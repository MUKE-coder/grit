import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Warn, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'
import { Diagram, DiagramBox } from '@/components/diagram'

export const metadata: Metadata = {
  title: 'GORM relationships: belongs to, has many, many to many',
  description:
    'Each GORM relation drawn, the struct that produces it, the SQL it writes, and the N+1 query problem every one of them can cause.',
  alternates: { canonical: 'https://gritframework.dev/guides/gorm-tutorial-relationships' },
  openGraph: {
    title: 'GORM relationships: belongs to, has many, many to many',
    description: 'Every relation, the struct, the SQL, and how to load them without N+1.',
    url: 'https://gritframework.dev/guides/gorm-tutorial-relationships',
    type: 'article',
  },
}

/** Two tables with the cardinality between them. */
function Relation({
  left,
  leftSub,
  right,
  rightSub,
  label,
}: {
  left: string
  leftSub: string
  right: string
  rightSub: string
  label: string
}) {
  return (
    <Diagram>
      <div className="flex flex-col items-stretch gap-3 sm:flex-row sm:items-center">
        <DiagramBox title={left} sub={leftSub} />
        <div className="shrink-0 px-2 text-center font-mono text-[11px] text-accent">{label}</div>
        <DiagramBox title={right} sub={rightSub} tone="primary" />
      </div>
    </Diagram>
  )
}

export default function Page() {
  return (
    <GuidePage
      slug="gorm-tutorial-relationships"
      lead={
        <>
          <p className="mb-4">
            A relationship in GORM is two things at once: a field on a Go struct, which is how your
            code reads it, and a foreign key column, which is how the database stores it. Most
            confusion with GORM comes from reading one and expecting the other.
          </p>
          <p>
            Each relation below is drawn, then written as a struct, then shown as the SQL GORM
            produces. If you only read one section, read the one on loading them, because that is
            where the performance goes.
          </p>
        </>
      }
    >
      <Outcome>
        <p>The four relation kinds, with the struct tags that make each one work.</p>
        <p>
          The ability to predict which column GORM will look for, which is what the error messages
          are really telling you.
        </p>
        <p>Eager loading, and why a list endpoint got slow after you added an author name.</p>
      </Outcome>

      <Section id="conventions" title="The conventions, before the relations">
        <p>
          GORM guesses a great deal from names. Knowing the guesses turns its error messages from
          mysterious into obvious.
        </p>
        <Checklist
          items={[
            <>
              A struct called <code>Post</code> maps to a table called <code>posts</code>:
              pluralised, snake_cased.
            </>,
            <>
              A field called <code>AuthorID</code> becomes a column called <code>author_id</code>.
            </>,
            <>
              A field named <code>ID</code> is the primary key, and is auto-incrementing if it is
              an integer.
            </>,
            <>
              <code>CreatedAt</code> and <code>UpdatedAt</code> are maintained for you if present.{' '}
              <code>DeletedAt gorm.DeletedAt</code> turns on soft delete for the whole model.
            </>,
          ]}
        />
        <Warn title="One naming trap worth knowing now">
          <p>
            GORM splits names on case boundaries, and it does not know your initialisms. A field
            called <code>AllowIDPInitiated</code> becomes the column{' '}
            <code>allow_id_p_initiated</code>, not <code>allow_idp_initiated</code>. If a column
            name ever looks mangled, pin it explicitly with{' '}
            <code>gorm:&quot;column:allow_idp_initiated&quot;</code> rather than fighting the
            conversion.
          </p>
        </Warn>
      </Section>

      <Section id="belongs-to" title="Belongs to: the many side holds the key">
        <Relation
          left="posts"
          leftSub="id, title, author_id"
          right="users"
          rightSub="id, name"
          label="author_id  →  id"
        />
        <p>
          A post belongs to one user. The foreign key lives on <code>posts</code>, which is the
          side that can have many rows pointing at the same user.
        </p>
        <CodeBlock
          code={`type Post struct {
	ID       uint   \`gorm:"primaryKey"\`
	Title    string
	// The column. This is what the database actually stores.
	AuthorID uint   \`gorm:"index;not null"\`
	// The Go-side field. Empty until you load it; GORM never fills it in
	// by accident, which is the behaviour people find surprising.
	Author   User
}

type User struct {
	ID   uint \`gorm:"primaryKey"\`
	Name string
}`}
        />
        <p>
          The pairing is by name: GORM sees a field of type <code>User</code> called{' '}
          <code>Author</code>, and looks for <code>AuthorID</code>. Name the field{' '}
          <code>Writer</code> and it looks for <code>WriterID</code>. When they do not line up, say
          so:
        </p>
        <CodeBlock
          code={`Author User \`gorm:"foreignKey:CreatedByID;references:ID"\``}
        />
        <Note title="Index the foreign key">
          <p>
            <code>gorm:&quot;index&quot;</code> on <code>AuthorID</code> is not optional in
            practice. Every &quot;posts by this author&quot; query filters on it, and without an
            index that is a full table scan. Postgres creates an index for a primary key and{' '}
            <strong>not</strong> for a foreign key, which catches people who assume otherwise.
          </p>
        </Note>
      </Section>

      <Section id="has-many" title="Has many: the same key, read the other way">
        <Relation
          left="users"
          leftSub="id, name"
          right="posts"
          rightSub="id, title, author_id"
          label="id  →  author_id"
        />
        <p>
          Has-many is belongs-to seen from the other end. There is no new column: the same{' '}
          <code>author_id</code> on <code>posts</code> supports both directions.
        </p>
        <CodeBlock
          code={`type User struct {
	ID    uint \`gorm:"primaryKey"\`
	Name  string
	// A slice, so "many". GORM looks on Post for a UserID, because this
	// struct is called User. Ours is called AuthorID, so say so.
	Posts []Post \`gorm:"foreignKey:AuthorID"\`
}`}
        />
        <p>
          Declaring both sides is optional. Declare the side you will actually read from; a
          relation nobody traverses is a field that confuses the next reader.
        </p>
        <H3>What happens to the children when the parent goes</H3>
        <p>
          Deleting a user with posts leaves orphaned rows pointing at an id that no longer exists,
          unless you say otherwise. Decide this deliberately:
        </p>
        <CodeBlock
          code={`// The database enforces it: deleting a user deletes their posts.
Posts []Post \`gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE"\`

// Or keep the posts and null the author, for "deleted user" content.
AuthorID *uint \`gorm:"index"\`
Posts    []Post \`gorm:"foreignKey:AuthorID;constraint:OnDelete:SET NULL"\``}
        />
        <p>
          Note the pointer in the second version. A nullable foreign key must be{' '}
          <code>*uint</code>, because a plain <code>uint</code> cannot hold null and GORM will
          write 0, which is a row pointing at a user that does not exist.
        </p>
      </Section>

      <Section id="has-one" title="Has one: a belongs-to that happens to be unique">
        <Relation
          left="users"
          leftSub="id, email"
          right="profiles"
          rightSub="id, user_id (unique), bio"
          label="id  →  user_id"
        />
        <CodeBlock
          code={`type User struct {
	ID      uint \`gorm:"primaryKey"\`
	Profile Profile
}

type Profile struct {
	ID     uint   \`gorm:"primaryKey"\`
	// uniqueIndex is what makes this has-ONE rather than has-many. Without
	// it the schema permits two profiles per user and only your code
	// objects, which means it will happen.
	UserID uint   \`gorm:"uniqueIndex;not null"\`
	Bio    string
}`}
        />
        <p>
          Worth asking whether you need it at all. A profile with three fields is usually three
          more columns on <code>users</code>. Split the table when the data is genuinely optional,
          genuinely large, or has a different access pattern: a thousand-character bio you never
          select in a list is a reasonable split.
        </p>
      </Section>

      <Section id="many-to-many" title="Many to many: a third table does the work">
        <Relation
          left="posts"
          leftSub="id, title"
          right="tags"
          rightSub="id, name"
          label="post_tags  ⟷"
        />
        <p>
          A post has many tags and a tag has many posts, so neither table can hold the key. A join
          table does, with one row per pairing.
        </p>
        <CodeBlock
          code={`type Post struct {
	ID   uint  \`gorm:"primaryKey"\`
	Tags []Tag \`gorm:"many2many:post_tags;"\`
}

type Tag struct {
	ID    uint   \`gorm:"primaryKey"\`
	Name  string \`gorm:"uniqueIndex;not null"\`
	Posts []Post \`gorm:"many2many:post_tags;"\`
}`}
        />
        <p>
          GORM creates <code>post_tags</code> with <code>post_id</code> and <code>tag_id</code>,
          and manages it. Adding and removing happens through the association API rather than by
          touching the join table:
        </p>
        <CodeBlock
          code={`db.Model(&post).Association("Tags").Append(&tag)     // add one
db.Model(&post).Association("Tags").Delete(&tag)     // remove one
db.Model(&post).Association("Tags").Replace(tags)    // set the whole list`}
        />
        <H3>When the join needs its own fields</H3>
        <p>
          The moment the relationship itself has attributes, the implicit join table stops being
          enough. &quot;User is a member of Organisation&quot; is many-to-many, until somebody asks
          what their role is, or when they joined.
        </p>
        <CodeBlock
          code={`// A real model, not an implicit table, because the membership has data.
type Membership struct {
	ID             uint      \`gorm:"primaryKey"\`
	UserID         uint      \`gorm:"index;uniqueIndex:idx_user_org;not null"\`
	OrganisationID uint      \`gorm:"index;uniqueIndex:idx_user_org;not null"\`
	Role           string    \`gorm:"not null;default:'member'"\`
	JoinedAt       time.Time
}`}
        />
        <p>
          The shared <code>uniqueIndex:idx_user_org</code> name on both columns makes one composite
          unique index, which is what stops the same person being added twice. Enforce it in the
          database rather than in a handler: a second write path will eventually arrive that does
          not check.
        </p>
      </Section>

      <Section id="loading" title="Loading them, and the query that ruins a list endpoint">
        <p>
          This is the section that matters for performance. A relationship field is empty until you
          ask for it, and <em>how</em> you ask decides whether your endpoint runs one query or five
          hundred.
        </p>
        <H3>Preload: a second query, then stitching</H3>
        <CodeBlock
          code={`var posts []Post
db.Preload("Author").Find(&posts)

// SELECT * FROM posts;
// SELECT * FROM users WHERE id IN (1,2,3,...);   <- one more, not one each`}
        />
        <p>
          Two queries regardless of how many posts come back. Nested relations use a dotted path,
          and a preload can be filtered:
        </p>
        <CodeBlock
          code={`db.Preload("Author").Preload("Comments.User").Find(&posts)

db.Preload("Comments", func(db *gorm.DB) *gorm.DB {
	return db.Where("approved = ?", true).Order("created_at DESC").Limit(5)
}).Find(&posts)`}
        />
        <H3>Joins: one query, for belongs-to and has-one</H3>
        <CodeBlock
          code={`db.Joins("Author").Find(&posts)
// SELECT posts.*, Author.* FROM posts LEFT JOIN users "Author" ON ...`}
        />
        <p>
          One query instead of two, and it lets you filter on the related table. Use it for the
          single-record relations. Do not use it for has-many: joining one post to its forty
          comments returns the post forty times, and the row count multiplies with every
          additional has-many you join.
        </p>
        <H3>The N+1</H3>
        <CodeBlock
          code={`// One query for the list, and then one more per row, in a loop that
// does not look like it is querying anything.
db.Find(&posts)
for i := range posts {
	db.First(&posts[i].Author, posts[i].AuthorID)   // 500 posts, 501 queries
}`}
        />
        <p>
          It is invisible locally, where the database is on the same machine and the dataset is
          small. In production, with a millisecond of network latency per query, five hundred rows
          is half a second of waiting for data you could have fetched in two queries.
        </p>
        <p>
          Watch for the version that hides inside a serialiser or a view helper, where the loop is
          in one file and the query is in another. The fix is always the same: load what you need
          before the loop.
        </p>
        <Note title="Make it visible">
          <p>
            Turn on query logging in development and watch the count while you click around. An
            N+1 announces itself immediately when you can see the queries and is nearly impossible
            to notice when you cannot.
          </p>
          <CodeBlock
            code={`db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
	Logger: logger.Default.LogMode(logger.Info),   // development only
})`}
          />
        </Note>
      </Section>

      <Section id="migrations" title="Getting the tables made">
        <p>
          <code>AutoMigrate</code> compares your structs to the database and adds what is missing:
          tables, columns, indexes and foreign keys.
        </p>
        <CodeBlock
          code={`// Order matters: a table referenced by a foreign key must exist first.
db.AutoMigrate(&User{}, &Profile{}, &Post{}, &Tag{})`}
        />
        <p>What it will not do, and this is the whole reason production needs something else:</p>
        <Checklist
          items={[
            <>It never drops a column. Removing a field from the struct leaves the column there.</>,
            <>
              It never renames. Renaming a field adds a new column and leaves the old one, with
              your data in it.
            </>,
            <>It will not narrow a type, because that could lose data.</>,
            <>
              You cannot review it before it runs. There is no file to read in a pull request and
              no way to roll it back.
            </>,
          ]}
        />
        <p>
          So: <code>AutoMigrate</code> in development, versioned SQL files in production. The{' '}
          <Link href="/docs/backend/migrations" className="text-accent hover:underline">
            migrations reference
          </Link>{' '}
          covers the pattern.
        </p>
      </Section>

      <Section id="with-grit" title="Declaring relations from the command line">
        <p>
          Grit generates these from a field list, which is mostly useful because it writes the
          matching TypeScript, the Zod schema and the admin picker at the same time, so the three
          cannot disagree.
        </p>
        <CodeBlock
          terminal
          code={`grit generate resource Post \\
  --fields "title:string:required,body:text,author:belongs_to:User,tags:many_to_many:Tag"`}
        />
        <p>
          A <code>belongs_to</code> field produces the foreign key, the index, the Go field, a
          searchable relationship picker in the admin showing the related record&apos;s label, and
          the preload on the list endpoint, which is the part people forget and then diagnose as an
          N+1 three months later.
        </p>
      </Section>

      <Faq
        slug="gorm-tutorial-relationships"
        items={[
          {
            q: 'Why is my Author field empty when the post loaded fine?',
            a: (
              <>
                Because GORM does not load relations unless you ask. <code>db.First(&amp;post)</code>{' '}
                fetches the <code>posts</code> row and nothing else, so <code>post.Author</code> is
                a zero-valued struct and <code>post.AuthorID</code> holds the real information. Add{' '}
                <code>.Preload(&quot;Author&quot;)</code> or <code>.Joins(&quot;Author&quot;)</code>
                . This is deliberate: the alternative is an ORM that fetches half your database
                because you touched a field.
              </>
            ),
            text:
              'Because GORM does not load relations unless asked. db.First(&post) fetches the posts row only, so post.Author is a zero-valued struct while post.AuthorID holds the real value. Add .Preload("Author") or .Joins("Author"). The behaviour is deliberate: lazy loading by default fetches half the database when you touch a field.',
          },
          {
            q: 'Preload or Joins?',
            a: (
              <>
                <code>Joins</code> for belongs-to and has-one: one query, and you can filter on the
                related table. <code>Preload</code> for has-many and many-to-many: joining those
                multiplies your rows, and with two has-many relations in one query the row count is
                the product of both. When unsure, <code>Preload</code> is the safer default.
              </>
            ),
            text:
              'Joins for belongs-to and has-one: one query, and you can filter on the joined table. Preload for has-many and many-to-many, because joining those multiplies rows, and two has-many joins in one query produce the product of both. Preload is the safer default when unsure.',
          },
          {
            q: 'Should I use many2many or write the join model myself?',
            a: (
              <>
                <code>many2many</code> while the relationship is only a link. Write the model the
                moment the link has attributes of its own, a role, a quantity, a joined date, or
                needs its own timestamps or soft delete. Migrating later is possible and tedious,
                so if you can imagine the attribute arriving, write the model now.
              </>
            ),
            text:
              'Use many2many while the relationship is only a link. Write an explicit model as soon as the link carries attributes such as a role, a quantity or a joined date, or needs its own timestamps or soft delete. Migrating later is possible but tedious.',
          },
          {
            q: 'Does AutoMigrate delete data?',
            a: (
              <>
                Not by itself: it only adds. The danger is the other direction. It leaves columns
                you removed and duplicates columns you renamed, so the schema quietly drifts from
                the struct, and then one day somebody writes a cleanup migration from the struct
                and drops a column that still has data in it.
              </>
            ),
            text:
              'Not by itself; it only adds. The risk is drift: it leaves columns you removed and duplicates ones you renamed, so the schema diverges from the struct, and a later cleanup written from the struct can drop a column that still holds data.',
          },
          {
            q: 'Why is a nullable foreign key a pointer?',
            a: (
              <>
                Because a plain <code>uint</code> cannot represent null: its zero value is 0, so
                GORM writes 0, and you get a row pointing at a record with id 0 that does not
                exist. <code>*uint</code> distinguishes &quot;no author&quot; from &quot;author
                zero&quot;. The same applies to any optional column: use a pointer, or{' '}
                <code>sql.NullString</code> and friends.
              </>
            ),
            text:
              'Because a plain uint cannot represent null: its zero value is 0, so GORM writes 0 and the row points at a non-existent record with id 0. A *uint distinguishes no author from author zero. The same applies to any optional column.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/first-go-react-app',
            label: 'Build your first full-stack app',
            blurb: 'Where the first model and the first AutoMigrate call come from.',
          },
          {
            href: '/guides/gin-rest-api-tutorial',
            label: 'Build a REST API with Gin',
            blurb: 'The handlers that read and write these models.',
          },
          {
            href: '/docs/backend/models',
            label: 'Models reference',
            blurb: 'Every field type a generated model can carry, and the tags each produces.',
          },
          {
            href: '/docs/backend/migrations',
            label: 'Migrations',
            blurb: 'What replaces AutoMigrate when the database has data you care about.',
          },
        ]}
      />
    </GuidePage>
  )
}
