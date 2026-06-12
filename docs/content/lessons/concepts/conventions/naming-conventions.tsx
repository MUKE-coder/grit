import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Grit names things consistently so you never have to guess what to
        call your new file. This is a 5-minute lesson — but it&apos;s the one
        that makes your future code look like it belongs.
      </p>

      <h2>The full table</h2>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Thing</th>
              <th className="text-left px-3 py-2 font-medium">Convention</th>
              <th className="text-left px-3 py-2 font-medium">Example</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr><td className="px-3 py-2">Go files</td><td className="font-mono text-[12px]">snake_case</td><td className="font-mono text-[12px]">user_handler.go</td></tr>
            <tr><td className="px-3 py-2">Go structs</td><td className="font-mono text-[12px]">PascalCase</td><td className="font-mono text-[12px]">type User struct</td></tr>
            <tr><td className="px-3 py-2">Go funcs (exported)</td><td className="font-mono text-[12px]">PascalCase</td><td className="font-mono text-[12px]">GetUsers()</td></tr>
            <tr><td className="px-3 py-2">Go funcs (private)</td><td className="font-mono text-[12px]">camelCase</td><td className="font-mono text-[12px]">parseToken()</td></tr>
            <tr><td className="px-3 py-2">TS files</td><td className="font-mono text-[12px]">kebab-case</td><td className="font-mono text-[12px]">use-users.ts</td></tr>
            <tr><td className="px-3 py-2">React components</td><td className="font-mono text-[12px]">PascalCase</td><td className="font-mono text-[12px]">DataTable.tsx</td></tr>
            <tr><td className="px-3 py-2">API routes</td><td className="font-mono text-[12px]">plural, lowercase</td><td className="font-mono text-[12px]">/api/users</td></tr>
            <tr><td className="px-3 py-2">DB tables</td><td className="font-mono text-[12px]">plural snake_case</td><td className="font-mono text-[12px]">blog_posts</td></tr>
            <tr><td className="px-3 py-2">Zod schemas</td><td className="font-mono text-[12px]">PascalCase + Schema</td><td className="font-mono text-[12px]">UserSchema</td></tr>
          </tbody>
        </table>
      </div>

      <h2>The pattern behind the table</h2>
      <p>Three principles cover everything:</p>
      <ol>
        <li>
          <strong>Go uses snake_case for filenames + PascalCase for types
          and exported funcs.</strong> That&apos;s the Go community convention.
          Grit doesn&apos;t fight it.
        </li>
        <li>
          <strong>TS uses kebab-case for filenames + PascalCase for
          components.</strong> The kebab-case convention is from React
          Router / Next.js — readable in URLs.
        </li>
        <li>
          <strong>API routes and DB tables are plural.</strong>{' '}
          <code>/api/users</code>, <code>users</code> table — a user is a
          single resource; the collection is plural. Avoids the &quot;is it /user
          or /users?&quot; debate.
        </li>
      </ol>

      <TipBox tone="info">
        These are the same conventions {' '}
        <a href="https://google.github.io/styleguide/" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Google&apos;s style guide</a>,{' '}
        <a href="https://airbnb.io/javascript/" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Airbnb&apos;s React guide</a>, and{' '}
        most production Go codebases use. Nothing exotic — just one place that picks all of them
        and applies them consistently.
      </TipBox>

      <KnowledgeCheck
        question="You're adding a React Query hook that fetches invoices. What's the filename?"
        choices={[
          {
            label: 'useInvoices.ts',
            feedback:
              "Wrong — camelCase isn't the TS file convention here. The hook NAME is camelCase (useInvoices), but the FILE is kebab-case.",
          },
          {
            label: 'use-invoices.ts',
            correct: true,
            feedback:
              "Right — TS files are kebab-case. Inside, the export is `useInvoices` (camelCase). File: kebab. Identifier: camel.",
          },
          {
            label: 'UseInvoices.ts',
            feedback:
              "Wrong — that's PascalCase, reserved for React components. Hooks aren't components.",
          },
          {
            label: 'use_invoices.ts',
            feedback:
              "Wrong — snake_case is for Go files only.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Look at your scaffolded project. Find <strong>three
              files</strong> across web, admin, and api — pick one from each.
              In <code>notes.md</code>, write the filename + which convention
              it follows.
            </p>
          </>
        }
        hint={<>Each app has very different conventions. That&apos;s the point.</>}
        solution={
          <>
            <p>Example picks:</p>
            <ul>
              <li>
                <code>apps/api/internal/handlers/user.go</code> — snake_case
                (Go file)
              </li>
              <li>
                <code>apps/web/components/landing-hero.tsx</code> — kebab-case
                (TS file)
              </li>
              <li>
                <code>apps/admin/app/resources/users/page.tsx</code> — plural
                lowercase URL segment (resource name) + kebab-case files
              </li>
            </ul>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Files are named. Next — how Grit&apos;s API responses are shaped. One
        envelope every endpoint follows. Memorise it once, use it forever.
      </p>
    </>
  )
}
