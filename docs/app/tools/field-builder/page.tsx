import Link from 'next/link'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { FieldBuilder } from '@/components/tools/field-builder'

export const metadata = {
  // The root layout appends the site name, so the title is the page's own.
  title: 'Field builder',
  description:
    'Build a grit generate resource command visually: pick from all 32 field types, set relationships, choices and modifiers, and copy the --fields string.',
  alternates: { canonical: 'https://gritframework.dev/tools/field-builder' },
}

export default function Page() {
  return (
    <div className="relative min-h-screen bg-background">
      <SiteHeader />
      <GridFrame />

      <main className="mx-auto max-w-6xl px-6 py-16">
        <nav aria-label="Breadcrumb" className="mb-6 text-xs text-muted-foreground">
          <Link href="/" className="hover:text-foreground">
            Home
          </Link>
          <span className="mx-2 text-muted-foreground/50">/</span>
          <Link href="/tools" className="hover:text-foreground">
            Tools
          </Link>
        </nav>

        <h1 className="mb-4 font-display text-4xl font-bold tracking-tight md:text-5xl">
          Field builder
        </h1>
        <p className="mb-10 max-w-2xl text-lg leading-relaxed text-muted-foreground">
          The <code className="font-mono text-base">--fields</code> string has a compact syntax
          with a few rules that are easy to forget: a relationship names its model, a vector needs
          its dimension, and choices go in brackets. Build it here instead.
        </p>

        <FieldBuilder />

        <section className="mt-16 border-t border-border/60 pt-10">
          <h2 className="mb-4 font-display text-xl font-semibold tracking-tight">
            The syntax, if you would rather type it
          </h2>
          <div className="space-y-3 text-sm leading-relaxed text-muted-foreground">
            <p>
              A field is <code>name:type</code>, with optional extras after it, and fields are
              separated by commas.
            </p>
            <ul className="space-y-2">
              <li>
                <code>title:string:required</code> adds a modifier. The modifiers are{' '}
                <code>required</code>, <code>unique</code> and <code>optional</code>.
              </li>
              <li>
                <code>author:belongs_to:User</code> names the related model in the third
                position. That is not a modifier, so modifiers come after it.
              </li>
              <li>
                <code>status:select[draft, sent, paid]</code> lists the choices. Write{' '}
                <code>draft=Draft in progress</code> when the label is not just the value
                capitalised.
              </li>
              <li>
                <code>embedding:vector[1536]</code> gives the dimension, which is part of the
                column type rather than a setting.
              </li>
              <li>
                <code>photo:file:image</code> restricts what may be uploaded.
              </li>
              <li>
                <code>slug:slug:title</code> says which field the slug is built from.
              </li>
              <li>
                <code>notes:text if status == disputed</code> shows the field only when another
                field has a value. The same thing can be written{' '}
                <code>notes:text:when=status=disputed</code>.
              </li>
            </ul>
            <p className="pt-2">
              The full list of types, with what each one produces across the stack, is in{' '}
              <Link href="/docs/concepts/field-types" className="text-accent hover:underline">
                the field types reference
              </Link>
              .
            </p>
          </div>
        </section>
      </main>
    </div>
  )
}
