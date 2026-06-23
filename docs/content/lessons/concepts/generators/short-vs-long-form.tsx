import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        You&apos;ve been mixing them all chapter — the inline{' '}
        <code>--fields &quot;…&quot;</code> form for quick resources, the{' '}
        <code>--from x.yaml</code> form when you needed{' '}
        <code>default:</code> or <code>unique:</code> on a belongs_to.
        This lesson is the explicit comparison: when each one earns its
        keep, and which one you reach for in which situation.
      </p>

      <h2>The same resource, both ways</h2>
      <p>
        Here&apos;s a realistic Article resource — title with slug, hero
        image, formatted body, status with a default — written in both
        forms side-by-side.
      </p>

      <h3>Short form (inline <code>--fields</code>)</h3>
      <CodeBlock
        terminal
        code={`grit generate resource Article \\
  --fields "title:string,slug:slug,cover:string,body:richtext,published:bool"`}
      />

      <h3>Long form (YAML)</h3>
      <CodeBlock
        language="yaml"
        filename="article.yaml"
        code={`name: Article
fields:
  - name: title
    type: string
    required: true
  - name: slug
    type: slug
    slug_source: title
  - name: cover
    type: string                # URL heuristic → VARCHAR(500)
  - name: body
    type: richtext
  - name: status
    type: string
    default: draft              # ← only YAML can set defaults
  - name: published
    type: bool
    default: false`}
      />
      <CodeBlock
        terminal
        code={`grit generate resource Article --from article.yaml`}
      />

      <p>
        Both produce the same 8 generated files <em>except</em> the YAML
        version also sets the database defaults
        (<code>status = &apos;draft&apos;</code>, <code>published = false</code>)
        and is committed into the repo — which means the next teammate
        who pulls can regenerate the resource from scratch by re-running
        the same command.
      </p>

      <h2>What each form supports</h2>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Feature</th>
              <th className="text-left px-3 py-2 font-medium">Inline <code>--fields</code></th>
              <th className="text-left px-3 py-2 font-medium">YAML <code>--from</code></th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr><td className="px-3 py-2 text-[12px]">All 13 field types</td><td className="text-center text-[12px]">✓</td><td className="text-center text-[12px]">✓</td></tr>
            <tr><td className="px-3 py-2 text-[12px]"><code>required</code> / <code>optional</code> / <code>unique</code></td><td className="text-center text-[12px]">✓</td><td className="text-center text-[12px]">✓</td></tr>
            <tr><td className="px-3 py-2 text-[12px]"><code>default</code> values</td><td className="text-center text-[12px]">✗</td><td className="text-center text-[12px]">✓</td></tr>
            <tr><td className="px-3 py-2 text-[12px]"><code>unique</code> on a <code>belongs_to</code> (one-to-one)</td><td className="text-center text-[12px]">✗</td><td className="text-center text-[12px]">✓</td></tr>
            <tr><td className="px-3 py-2 text-[12px]">Slug <code>slug_source</code> override</td><td className="text-center text-[12px]">✓ (third colon)</td><td className="text-center text-[12px]">✓ (explicit key)</td></tr>
            <tr><td className="px-3 py-2 text-[12px]">Re-runnable / version-controllable</td><td className="text-center text-[12px]">Sort of (commit message)</td><td className="text-center text-[12px]">✓ (committed file)</td></tr>
            <tr><td className="px-3 py-2 text-[12px]">Comments next to each field</td><td className="text-center text-[12px]">✗</td><td className="text-center text-[12px]">✓</td></tr>
            <tr><td className="px-3 py-2 text-[12px]">Reads cleanly past 5 fields</td><td className="text-center text-[12px]">✗ (one long line)</td><td className="text-center text-[12px]">✓</td></tr>
          </tbody>
        </table>
      </div>

      <h2>YAML field reference — the full spec</h2>
      <p>
        Each field entry under <code>fields:</code> accepts seven keys.
        Most have inline equivalents — two don&apos;t:
      </p>

      <CodeBlock
        language="yaml"
        code={`- name: <field_name>              # camelCase or snake_case — required
  type: <field_type>              # one of the 13 types — required
  required: true|false            # default depends on type (string=true, others=false)
  unique: true|false              # adds DB unique index. Works on belongs_to (1-to-1!).
  default: <value>                # DB default. YAML-only.
  slug_source: <field_name>       # for type: slug — which field to slugify from
  related_model: <ModelName>      # for type: belongs_to or many_to_many — what it points at`}
      />

      <p>
        That&apos;s the whole long-form vocabulary. The inline form
        compresses the first four into colons (<code>name:type:modifier</code>)
        and uses the third colon-separated slot for{' '}
        <code>slug_source</code> / <code>related_model</code>:
      </p>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Inline</th>
              <th className="text-left px-3 py-2 font-medium">YAML equivalent</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">name:string</td>
              <td className="px-3 py-2 font-mono text-[12px]">{`{ name: name, type: string, required: true }`}</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">phone:string:optional</td>
              <td className="px-3 py-2 font-mono text-[12px]">{`{ name: phone, type: string, required: false }`}</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">email:string:unique</td>
              <td className="px-3 py-2 font-mono text-[12px]">{`{ name: email, type: string, required: true, unique: true }`}</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">slug:slug:sku</td>
              <td className="px-3 py-2 font-mono text-[12px]">{`{ name: slug, type: slug, slug_source: sku }`}</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">group:belongs_to:Group</td>
              <td className="px-3 py-2 font-mono text-[12px]">{`{ name: group, type: belongs_to, related_model: Group }`}</td>
            </tr>
            <tr>
              <td className="px-3 py-2 font-mono text-[12px]">tags:many_to_many:Tag</td>
              <td className="px-3 py-2 font-mono text-[12px]">{`{ name: tags, type: many_to_many, related_model: Tag }`}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <h2>When to reach for which</h2>

      <h3>Use the short form when…</h3>
      <ul>
        <li>You have 1–5 fields and no defaults.</li>
        <li>
          You&apos;re prototyping — the field list will probably change
          before this resource is real.
        </li>
        <li>
          You want the command to fit in a commit message or a chat
          message to a teammate.
        </li>
        <li>You&apos;re pairing live and reading the spec out loud.</li>
      </ul>

      <h3>Use the long form (YAML) when…</h3>
      <ul>
        <li>
          You need <code>default</code> values, or you need{' '}
          <code>unique</code> on a <code>belongs_to</code> (one-to-one).
        </li>
        <li>You have more than 5 fields and the inline string stops fitting on one line.</li>
        <li>
          You want the resource spec checked into git so anyone can
          re-generate from scratch later. Drop <code>article.yaml</code>{' '}
          next to <code>grit.config.ts</code>.
        </li>
        <li>
          You&apos;re generating multiple resources at once and want
          them all in <code>resources/*.yaml</code> so a script can
          regenerate the lot.
        </li>
        <li>
          The resource spec is reviewed in a PR before code is written —
          the YAML is the design doc <em>and</em> the source.
        </li>
      </ul>

      <h2>The third option — interactive</h2>
      <p>
        Don&apos;t forget about <code>-i</code> when you genuinely
        don&apos;t know the fields yet:
      </p>
      <CodeBlock
        terminal
        code={`grit generate resource Article -i`}
      />
      <p>
        The CLI walks you through each field one prompt at a time, and
        at the end it shows you the equivalent <code>--fields</code>{' '}
        string. Great for pairing or for when you&apos;re still
        designing the model in your head.
      </p>

      <TipBox tone="info">
        Many teams keep <strong>both</strong>: a{' '}
        <code>resources/</code> directory full of committed YAML files
        for the &quot;real&quot; resources, and ad-hoc inline commands
        for prototypes and experiments. Once a prototype settles, copy
        the inline command into a YAML and commit it.
      </TipBox>

      <KnowledgeCheck
        question="A teammate asks you to generate an Invoice resource with: number (unique string), amount (money), description (text), due_date (date), and a default status of 'pending'. Which form fits?"
        choices={[
          {
            label: 'Short form: grit generate resource Invoice --fields "number:string:unique,amount:float,description:text,due_date:date,status:string:default=pending"',
            feedback:
              "Almost — the syntax is right but `:default=pending` isn't valid in the inline form (only `required`, `optional`, `unique` modifiers are accepted). The default has to come from YAML.",
          },
          {
            label: 'Long form YAML with default: pending on the status field',
            correct: true,
            feedback:
              "Right — inline modifiers don't support defaults. Drop a tiny invoice.yaml at the project root and run `grit generate resource Invoice --from invoice.yaml`.",
          },
          {
            label: 'Short form without the default, then edit the model.go by hand to add `gorm:"default:pending"`',
            feedback:
              "Works but leaves you with a hand-edited model that diverges from any committed YAML. If you ever re-run the generator, the default is lost. Better to keep the spec in YAML.",
          },
          {
            label: 'Interactive mode (`-i`) — the prompt will ask for defaults',
            feedback:
              "Interactive mode goes through fields but doesn't currently surface a default prompt. Use YAML for that.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Take this five-field resource and write it{' '}
              <strong>both ways</strong> in your <code>notes.md</code>,
              then actually generate it once using YAML:
            </p>
            <ul>
              <li>Resource: <strong>Project</strong></li>
              <li><code>name</code> — string, required, unique</li>
              <li><code>slug</code> — auto-generated from <code>name</code></li>
              <li><code>description</code> — text, optional</li>
              <li><code>status</code> — string with default <code>active</code></li>
              <li><code>archived</code> — bool with default <code>false</code></li>
            </ul>
            <p>
              Save the YAML as <code>project.yaml</code>, run{' '}
              <code>grit generate resource Project --from project.yaml</code>,
              then <code>grit migrate</code>. Confirm in the admin panel
              that a new Project starts with <code>status = active</code>{' '}
              and <code>archived = false</code> without you typing them.
            </p>
          </>
        }
        hint={
          <>
            Inline form needs 5 colon-separated specs in one quoted
            string. Two of those will be impossible to fit inline
            (defaults) — that&apos;s the whole point of the exercise.
            Write the inline version <em>without</em> defaults to show
            the gap, then write YAML <em>with</em> defaults to show
            what&apos;s gained.
          </>
        }
        solution={
          <>
            <p>Inline (incomplete — can&apos;t express defaults):</p>
            <CodeBlock
              terminal
              code={`grit generate resource Project \\
  --fields "name:string:unique,slug:slug,description:text:optional,status:string,archived:bool"`}
            />
            <p>YAML (complete):</p>
            <CodeBlock
              language="yaml"
              filename="project.yaml"
              code={`name: Project
fields:
  - name: name
    type: string
    required: true
    unique: true
  - name: slug
    type: slug
    slug_source: name
  - name: description
    type: text
    required: false
  - name: status
    type: string
    default: active
  - name: archived
    type: bool
    default: false`}
            />
            <p>
              Notice the YAML version is longer but says <em>more</em>{' '}
              — it&apos;s explicit about defaults and self-describing.
              Once it&apos;s committed, anyone can rebuild the resource
              identically.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        You can now write any field shape Grit supports in either form.
        One last topic for the chapter:{' '}
        <code>grit sync</code> — what to run after you edit a Go struct
        by hand so the TypeScript types catch up.
      </p>
    </>
  )
}
