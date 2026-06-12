import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Welcome. By the end of this lesson you&apos;ll be able to explain — in
        about 30 seconds — what Grit is, who built it for what, and whether your
        next project is a fit. That&apos;s the only goal for the next 6 minutes.
      </p>

      <h2>The one-paragraph answer</h2>
      <p>
        Grit is a <strong>full-stack meta-framework</strong> that scaffolds
        production-ready apps in one command. It pairs a <strong>Go (Gin + GORM)</strong>{' '}
        backend with React on the front and ships every &quot;battery&quot; — auth, jobs,
        mail, file storage, AI, observability, security WAF, admin panel —
        out of the box. Instead of gluing 12 libraries together for a week,
        you run <code>grit new my-app</code> and start shipping product on day
        one.
      </p>

      <h2>What &quot;meta-framework&quot; actually means</h2>
      <p>
        A framework gives you patterns and primitives. A meta-framework gives
        you a framework + the conventions + the scaffolder + the day-one
        wiring. Three concrete differences:
      </p>

      <ul>
        <li>
          <strong>Opinionated folder structure.</strong> You don&apos;t decide where
          handlers go vs. services vs. models — there&apos;s one place each, and
          every Grit project looks the same. That&apos;s a feature: every team
          member, every AI agent, every future-you can find what they need
          without thinking.
        </li>
        <li>
          <strong>Code generation.</strong>{' '}
          <code>grit generate resource Product</code> writes the Go model,
          handler, service, Zod schema, TypeScript type, React Query hook, and
          admin panel resource — 7 files, one command, all wired correctly.
        </li>
        <li>
          <strong>Batteries shipped, not described.</strong> Most frameworks
          tell you to &quot;pick an auth library&quot;. Grit ships JWT + OAuth + 2FA
          out of the box, with the trusted-device flow already wired, and you
          can swap in WorkOS or Auth0 if you outgrow it.
        </li>
      </ul>

      <TipBox tone="info">
        Grit is <strong>not</strong> a no-code tool, not a CMS, and not a
        Backend-as-a-Service. You&apos;re writing real Go and real React. Grit just
        does the scaffolding-and-conventions part for you.
      </TipBox>

      <h2>What you get day one</h2>
      <p>One <code>grit new my-app</code> later, you have:</p>

      <CodeBlock
        language="text"
        filename="my-app/"
        code={`my-app/
├── apps/
│   ├── api/             ← Go API (Gin + GORM)
│   │   ├── cmd/server/  ← main.go entry point
│   │   └── internal/    ← models, handlers, services, middleware
│   ├── web/             ← Next.js public site + product surface
│   └── admin/           ← Filament-style admin panel (Next.js)
├── packages/shared/     ← Zod schemas + TS types (used everywhere)
├── docker-compose.yml   ← Postgres, Redis, MinIO, Mailhog
├── grit.json            ← project metadata
└── .env                 ← random-generated secrets, ready for prod`}
      />

      <p>
        Boot Docker, run the API + web + admin, log into the admin panel as
        the seeded user, and you have a working product skeleton. We&apos;ll do
        exactly that in chapter 2.
      </p>

      <h2>The tagline you&apos;ll see everywhere</h2>
      <blockquote className="border-l-2 border-primary pl-4 italic text-foreground/80 my-4">
        Go + React. Built with Grit.
      </blockquote>
      <p>
        That&apos;s the elevator pitch — <em>Go on the back, React on the front,
        meta-framework gluing them together</em>. Memorise it. You&apos;ll repeat it
        to teammates, in interviews, and at meetups.
      </p>

      <KnowledgeCheck
        question="Which of these would NOT be in a fresh `grit new my-app` output?"
        choices={[
          {
            label: 'A Postgres docker-compose service ready to start',
            feedback:
              'Wrong direction — docker-compose with Postgres ships in every project by default.',
          },
          {
            label: 'A pre-wired admin panel at apps/admin',
            feedback:
              'Wrong direction — the Filament-style admin panel is one of the batteries that ships out of the box.',
          },
          {
            label: 'An empty .env file you need to fill in line by line',
            correct: true,
            feedback:
              "Right — Grit generates random secrets (JWT_SECRET, SENTINEL_SECRET_KEY, etc.) at scaffold time so APP_ENV=production works on first boot.",
          },
          {
            label: 'A Go module at apps/api with main.go',
            feedback:
              'Wrong direction — apps/api/cmd/server/main.go is the API entry point and ships in every project.',
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              In your own words (one paragraph, ~30 seconds when read aloud):
              explain to a friend what Grit is and who would use it. Write the
              answer in <code>notes.md</code> — the file you&apos;ll grow throughout
              this course.
            </p>
          </>
        }
        hint={
          <>
            Cover three things: (1) what it is (meta-framework, Go + React),
            (2) what problem it solves (scaffolding + glue work), (3) who
            it&apos;s for (devs who want to ship products, not wire frameworks).
          </>
        }
        solution={
          <>
            <p>One acceptable answer:</p>
            <blockquote className="border-l-2 border-primary pl-4 italic">
              Grit is an open-source meta-framework for building production
              apps with Go on the backend and React on the front. Instead of
              spending a week wiring 12 libraries — auth, jobs, mail, storage,
              admin panel — Grit ships them already wired and gives you a
              code generator that writes 7 files every time you add a new
              resource. It&apos;s for solo devs and small teams who want to ship
              real products quickly without compromising on the stack.
            </blockquote>
            <p>
              If yours says the same three things — what it is, what it
              replaces, who it&apos;s for — you&apos;ve got it.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        You now have the elevator pitch. The next lesson covers the more
        useful question: <strong>when should you reach for Grit, and when
        should you not?</strong> Trade-offs vs. plain Next.js, Rails, Laravel,
        and others — honestly compared.
      </p>
    </>
  )
}
