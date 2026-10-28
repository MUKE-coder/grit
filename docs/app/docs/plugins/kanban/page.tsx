import Link from 'next/link'
import { ArrowRight, SquareKanban } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Kanban boards — Grit',
  description:
    'grit plugin add kanban: boards, columns and cards with drag-and-drop by pointer, touch or keyboard. Reordering writes one row, using fractional indexing rather than a position column.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/kanban' },
}

export default function KanbanPluginPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <SquareKanban className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">
              Plugin
            </span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">Kanban boards</h1>
          <p className="mb-8 text-lg leading-relaxed text-muted-foreground">
            Boards, columns and cards in the admin, with dragging by pointer, by finger and by
            keyboard. It is also the plugin the{' '}
            <Link href="/docs/plugins/authoring" className="text-primary hover:underline">
              authoring guide
            </Link>{' '}
            is written around, because a board is the smallest feature that needs every part of the
            plugin system at once.
          </p>

          <CodeBlock
            language="bash"
            code={`grit plugin add kanban
pnpm install
cd apps/api && go run cmd/migrate/main.go`}
          />

          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            Boards then appears in the admin sidebar, under Content. A new board starts with To do,
            In progress and Done, because landing on an empty screen with one button is a worse
            first thirty seconds than landing on a board.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">What you get</h2>
          <ul className="mb-6 space-y-2 leading-relaxed text-muted-foreground">
            <li>
              <code>Board</code>, <code>BoardColumn</code> and <code>BoardCard</code> models, with
              migrations
            </li>
            <li>CRUD routes for all three, plus a move endpoint for cards and one for columns</li>
            <li>A boards list page and a board page in the admin, with React Query hooks</li>
            <li>An optional work-in-progress limit per column, enforced on the server</li>
            <li>
              A <code>boards.card_moved</code> event on the domain bus, so webhooks, the audit log
              and realtime subscribers all hear about it without the handler calling them
            </li>
          </ul>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Boards are visible to every signed-in user and record who created them. If boards need
            to be per-tenant, install the{' '}
            <Link href="/docs/plugins/multitenant" className="text-primary hover:underline">
              multitenant plugin
            </Link>
            , which scopes every model including these.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">
            The ordering, and why it is a string
          </h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            This is the part worth reading even if you never install the plugin, because every
            ordered list has the same problem: a playlist, a priority queue, a sortable table.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            The obvious design is an integer <code>position</code> column. Moving a card then means:
          </p>
          <CodeBlock
            language="sql"
            code={`UPDATE cards SET position = position + 1
WHERE column_id = ? AND position >= ?`}
          />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            That rewrites every card below the drop. With one person on the board it works. With
            two, their rewrites interleave and the column ends up with duplicate or skipped
            positions, and on Postgres the two statements can deadlock outright because they touch
            the same rows in opposite orders. Gap numbering (100, 200, 300) only delays it: the gap
            closes after a few drags between the same pair and you are renumbering again.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            So a rank here is a <strong>string</strong>, compared lexicographically, and there is
            always a string strictly between any two others. A move is:
          </p>
          <CodeBlock
            language="sql"
            code={`UPDATE board_cards SET column_id = ?, rank = ? WHERE id = ?`}
          />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            One row. No renumbering, no transaction spanning cards somebody else is dragging, and no
            lock held while a user decides where to drop. The same <code>ORDER BY</code> works
            unchanged on Postgres, MySQL and SQLite, which matters because Grit targets all three
            and a scheme needing dialect-specific SQL is three bugs.
          </p>

          <h3 className="mb-3 mt-8 text-xl font-semibold tracking-tight">What the ranks look like</h3>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Append three cards to an empty column and they get <code>a0</code>, <code>a1</code>,{' '}
            <code>a2</code>. Drag the last one to the top and exactly one row changes:
          </p>
          <CodeBlock
            language="text"
            code={`before                        after
  a0  Write the launch post       Zz  Email the list      <- the only row written
  a1  Record the demo             a0  Write the launch post
  a2  Email the list              a1  Record the demo`}
          />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            The keys stay short because a rank has an integer part that increments rather than a
            string that grows. Ten thousand appends leave the key four characters long. Only
            repeated insertion between the same two neighbours lengthens one, by a character at a
            time, and it never runs out of room.
          </p>

          <h3 className="mb-3 mt-8 text-xl font-semibold tracking-tight">Two people, one slot</h3>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Two people who drop a card into the same slot at the same moment compute the same rank.
            That is a tie, not corruption. Every read orders by <code>(rank, id)</code>, so the tie
            breaks the same way for both of them and the board is stable.
          </p>

          <h3 className="mb-3 mt-8 text-xl font-semibold tracking-tight">Stale neighbours</h3>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            The move request names the two cards the dragged card was dropped between, not an index.
            An index is stale the instant anybody else moves a card, and a stale index silently puts
            the card in the wrong place. Neighbour ids survive that: if one of them has moved away,
            the server notices and appends rather than computing a position from a row that is no
            longer where the browser thought it was.
          </p>
          <CodeBlock
            language="http"
            code={`POST /api/v1/board-cards/:id/move

{
  "column_id": "01a1...",
  "prev_id": "01a2...",   // the card it now sits below, empty for the top
  "next_id": "01a3..."    // the card it now sits above, empty for the bottom
}`}
          />

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Keyboard and screen reader</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            Every card has a drag handle that is a real <code>button</code>, so it is in the tab
            order. Space picks a card up, the arrow keys move it within and between columns, Space
            drops it and Escape cancels. A live region announces the column by name: &quot;Picked up
            Email the list from To do&quot;, &quot;Over In progress&quot;, &quot;Dropped Email the
            list in In progress&quot;.
          </p>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            This is the entire reason the plugin takes a dependency on{' '}
            <code>@dnd-kit</code> rather than using HTML5 drag events. Those do not fire on touch
            devices and cannot be driven from the keyboard at all, so a hand-rolled board works in a
            demo and excludes two groups of people who are not in the room when it is demonstrated.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Work-in-progress limits</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            A column can cap how many cards it holds. The limit is checked on the server, and a move
            that would exceed it is answered with <code>409</code> and the card goes back where it
            came from. Reordering inside a column that is already at its limit is allowed: a limit
            is about how much work is in progress, not about whether it can be rearranged.
          </p>
          <CodeBlock
            language="bash"
            code={`curl -X PATCH /api/v1/board-columns/:id -d '{"wip_limit": 3}'`}
          />

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Removing it</h2>
          <CodeBlock language="bash" code={`grit plugin remove kanban`} />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            Removal replays the install backwards from the lockfile: the files it wrote are deleted,
            the lines it injected are taken back out, and the project compiles again. Your board
            data stays in the database until you drop the tables.
          </p>

          <div className="mt-12 flex flex-col gap-3 sm:flex-row">
            <Link
              href="/docs/plugins/authoring"
              className="group flex flex-1 items-center gap-3 rounded-xl border border-border p-4 transition-colors hover:border-primary"
            >
              <span className="flex-1">
                <span className="block font-medium">Write a plugin like this one</span>
                <span className="text-sm text-muted-foreground">The authoring guide</span>
              </span>
              <ArrowRight className="h-4 w-4 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
            </Link>
            <Link
              href="/docs/extending"
              className="group flex flex-1 items-center gap-3 rounded-xl border border-border p-4 transition-colors hover:border-primary"
            >
              <span className="flex-1">
                <span className="block font-medium">Plugin, resource, block or blueprint</span>
                <span className="text-sm text-muted-foreground">Which one a job needs</span>
              </span>
              <ArrowRight className="h-4 w-4 text-muted-foreground transition-transform group-hover:translate-x-0.5" />
            </Link>
          </div>
        </div>
      </main>
    </div>
  )
}
