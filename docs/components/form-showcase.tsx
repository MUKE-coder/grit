'use client'

import { useState } from 'react'
import type React from 'react'
import { Upload, Link2, Table2, ListOrdered, WifiOff } from 'lucide-react'

/**
 * What one generate command actually produces, on screen.
 *
 * The pitch "it generates a form for you" is cheap; every scaffolder claims it.
 * These are screenshots of the forms Grit generated from the field definitions
 * printed beside them — same project, same session, nothing hand-edited except
 * the steps block on the multi-step tab, which is the documented way to turn a
 * flat form into a wizard.
 *
 * The `fields` strings are the exact arguments that produced each screenshot.
 * If you change the field syntax, re-run them and re-shoot rather than editing
 * the string to match.
 */

interface Demo {
  key: string
  label: string
  icon: React.ComponentType<{ className?: string }>
  headline: string
  body: string
  /** The command that produced the screenshot, or the config that enabled it. */
  command: string
  commandNote: string
  detail: string
  image: string
}

const DEMOS: Demo[] = [
  {
    key: 'uploads',
    label: 'File uploads',
    icon: Upload,
    headline: 'Four kinds of upload, four words',
    body:
      'Declare a field as file or files and add what it accepts — image, pdf, zip, doc, video, or a bracketed list. You get a dropzone with the right filter, image previews, per-type icons, size caps and progress, on both create and edit.',
    command:
      'grit generate resource Product --fields \\\n  "name:string,category:belongs_to:Category,price:float,\\\n   cover:file:image,gallery:files:image,\\\n   spec_sheet:file:pdf,downloads:files:[zip,doc],\\\n   description:richtext,published:bool"',
    commandNote: 'One command — model, migration, API, types, hooks and this form.',
    detail:
      'Uploads go browser-to-storage through a presigned URL, and the field’s accept list is enforced on the server too — a field declared file:pdf will not take a PNG even if the client asks nicely.',
    image: '/images/forms/uploads.png',
  },
  {
    key: 'relationship',
    label: 'Relationships',
    icon: Link2,
    headline: 'A foreign key you can fill without leaving the form',
    body:
      'A belongs_to field becomes a searchable select fed by the related resource. If the record you need does not exist yet, create it inline — the related resource’s own form opens, and the new record is selected when it closes.',
    command: 'category:belongs_to:Category',
    commandNote: 'One field. The select, the search, and New Category all follow.',
    detail:
      'The inline form is the real one, multi-step included. No half-featured “quick add” dialog that drifts out of sync with the actual create form.',
    image: '/images/forms/relationship.png',
  },
  {
    key: 'line-items',
    label: 'Line items',
    icon: Table2,
    headline: 'Invoices, orders, anything with rows',
    body:
      'Declare a has-many child with --items and the parent form grows an editable line-items table: add and remove rows, per-row totals, and a grand total that keeps itself right.',
    command:
      'grit generate resource Invoice \\\n  --fields "number:string:auto,client_name:string,\\\n            issued_on:date,due_on:date,notes:text" \\\n  --items "InvoiceItem:description:string,qty:int,unit_rate:float"',
    commandNote: 'Parent, child, the join, the form and a PDF endpoint.',
    detail:
      'Children are written in the parent’s transaction, so a half-saved invoice with orphaned rows is not a state you can reach. number:string:auto gives you INV-0001 without a counter to maintain.',
    image: '/images/forms/line-items.png',
  },
  {
    key: 'multi-step',
    label: 'Multi-step',
    icon: ListOrdered,
    headline: 'Long forms, broken into steps',
    body:
      'Group the fields into steps and the same generated form becomes a wizard — on a page or in a modal, with validation per step so someone cannot skip past a required field.',
    command:
      'formView: "page-steps",\nform: {\n  steps: [\n    { title: "Applicant",    fields: ["full_name", "email", "phone"] },\n    { title: "Experience",   fields: ["role_applied", "years_experience", "portfolio_url", "resume"] },\n    { title: "Availability", fields: ["available_from", "expected_salary", "cover_letter", "notes"] },\n  ],\n}',
    commandNote: 'Added to the generated resource file. The fields are unchanged.',
    detail:
      'On edit, each step gets its own Update button that saves only that step, and stays disabled until that step actually changes — so fixing a phone number does not rewrite the whole record.',
    image: '/images/forms/multi-step.png',
  },
  {
    key: 'offline',
    label: 'Offline desktop',
    icon: WifiOff,
    headline: 'The desktop app keeps working with no network',
    body:
      'The Wails app mirrors your data into local SQLite. Reads keep answering, writes queue, and everything pushes automatically when the connection returns — including file uploads, which hold their bytes until storage is reachable.',
    command: 'grit new myapp --triple --desktop',
    commandNote: 'The mirror, the queue and this page come with the app.',
    detail:
      'The state is visible rather than magic: status, last sync, pending count and a Work offline switch for testing the path deliberately instead of by unplugging your router.',
    image: '/images/forms/offline.png',
  },
]

export function FormShowcase() {
  const [active, setActive] = useState(DEMOS[0].key)
  const demo = DEMOS.find((d) => d.key === active) ?? DEMOS[0]

  return (
    <div>
      <div role="tablist" aria-label="Form capability" className="flex flex-wrap gap-2 mb-8">
        {DEMOS.map((d) => {
          const selected = d.key === active
          return (
            <button
              key={d.key}
              type="button"
              role="tab"
              aria-selected={selected}
              onClick={() => setActive(d.key)}
              className={`inline-flex items-center gap-2 rounded-full border px-4 py-2 text-sm font-medium transition-colors ${
                selected
                  ? 'border-primary/40 bg-primary/10 text-foreground'
                  : 'border-border/60 text-muted-foreground hover:text-foreground hover:border-border'
              }`}
            >
              <d.icon className="h-4 w-4" />
              {d.label}
            </button>
          )
        })}
      </div>

      <div className="grid lg:grid-cols-[1fr_20rem] gap-8 lg:gap-10 items-start">
        <div className="rounded-xl overflow-hidden border border-border bg-card/40 shadow-[0_24px_64px_-16px_rgba(2,6,23,0.5)]">
          <div className="flex items-center gap-2 px-3.5 py-2.5 bg-card/70 border-b border-border/60">
            <div className="flex items-center gap-1.5">
              <span className="h-2.5 w-2.5 rounded-full bg-[#ff5f57]" />
              <span className="h-2.5 w-2.5 rounded-full bg-[#febc2e]" />
              <span className="h-2.5 w-2.5 rounded-full bg-[#28c840]" />
            </div>
            <span className="mx-auto text-[11px] font-mono text-muted-foreground">
              {demo.key === 'offline' ? 'myapp — desktop' : 'localhost:3001'}
            </span>
          </div>
          {/* eslint-disable-next-line @next/next/no-img-element */}
          <img
            key={demo.key}
            src={demo.image}
            alt={`${demo.label}: ${demo.headline}`}
            className="w-full h-auto block"
            loading="lazy"
          />
        </div>

        <div>
          <h3 className="text-lg font-semibold mb-3 leading-snug">{demo.headline}</h3>
          <p className="text-sm text-muted-foreground leading-relaxed mb-5">{demo.body}</p>

          <div className="rounded-xl border border-border/60 bg-card/40 overflow-hidden mb-4">
            <div className="px-4 pt-3 pb-2 text-[10.5px] font-mono uppercase tracking-wider text-muted-foreground/70">
              What produced it
            </div>
            {/* Wide command lines scroll inside the card rather than pushing the
                page sideways. */}
            <div className="overflow-x-auto px-4 pb-3">
              <pre className="text-[11.5px] font-mono text-foreground/90 leading-relaxed whitespace-pre">
                {demo.command}
              </pre>
            </div>
            <div className="px-4 py-2.5 border-t border-border/50 bg-card/40">
              <p className="text-[11.5px] text-muted-foreground">{demo.commandNote}</p>
            </div>
          </div>

          <div className="rounded-xl border border-border/60 bg-card/40 p-4">
            <div className="text-[10.5px] font-mono uppercase tracking-wider text-muted-foreground/70 mb-2">
              What to notice
            </div>
            <p className="text-[12.5px] text-foreground/80 leading-relaxed">{demo.detail}</p>
          </div>
        </div>
      </div>
    </div>
  )
}
