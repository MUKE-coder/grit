// The book's design system.
//
// Everything visual lives here so a chapter file is only prose, code and the
// four teaching devices below. If a chapter needs a new kind of box, it goes
// in this file rather than being styled inline, or the book stops looking like
// one book.

#let accent = rgb("#1f6feb") // links, chapter numbers, rules
#let ink = rgb("#14161a") // body text
#let muted = rgb("#5b6672") // captions, asides
#let codebg = rgb("#f6f7f9")
#let codeline = rgb("#e3e6ea")

// The chapter number on a chapter opener.
//
// Not counter(heading): headings here are set numbering: none, and an
// unnumbered heading does not step its counter, so every chapter printed the
// same number. This one is stepped by the level 1 show rule itself. The
// display is the count minus one, because the book opens with a chapter 0,
// the setup chapter, and Part I has to start at 1.
#let chapternum = counter("chapter")

// Back matter opens like a chapter and carries no number, because "59
// Afterword" is not a thing anybody says. The counter still steps; only the
// number is hidden, so nothing later depends on where this was set.
#let backmatter = state("backmatter", false)

// The four teaching devices, each with its own colour so a reader learns the
// shapes rather than reading the label every time.
#let bridgecol = rgb("#f0b429") // yellow: you already know this
#let gotchacol = rgb("#d1343c") // red: this will bite you
#let projectcol = rgb("#1f9d55") // green: build something
#let trycol = rgb("#6b46c1") // purple: your turn

#let callout(fill: gray, title: "", body) = block(
  width: 100%,
  inset: (left: 12pt, right: 12pt, top: 10pt, bottom: 10pt),
  radius: 3pt,
  fill: fill.lighten(90%),
  stroke: (left: 3pt + fill),
  breakable: true,
  [
    #text(fill: fill.darken(25%), weight: "bold", size: 9pt)[#upper(title)]
    #v(-4pt)
    #set text(size: 9.8pt)
    // A callout has 27pt less measure than the body, so its code is scaled
    // by the same ratio to keep one column budget across the whole book.
    #show raw: set text(size: 7.6pt)
    #body
  ],
)

// The device this whole book rests on: take something the reader already does
// in JavaScript and name the Go equivalent before explaining it.
#let bridge(title: "If you know JavaScript", body) = callout(
  fill: bridgecol, title: title, body,
)

#let gotcha(title: "Watch out", body) = callout(
  fill: gotchacol, title: title, body,
)

#let project(title: "Project", body) = callout(
  fill: projectcol, title: title, body,
)

#let tryit(title: "Your turn", body) = callout(
  fill: trycol, title: title, body,
)

// The assignment: the one device in the book that leaves the page.
//
// Every chapter ends with one. It is work the reader pushes to their own
// repository under an exact commit message, which is what the Committed app
// scores for the Gopher League leaderboard. The commit name is part of the
// device rather than prose, because a leaderboard that cannot match the
// message gives no points, and a reader who has to invent the name will
// invent a different one every time.
#let assignment(name: "", body) = block(
  width: 100%,
  inset: (left: 12pt, right: 12pt, top: 10pt, bottom: 11pt),
  radius: 3pt,
  fill: accent.lighten(92%),
  stroke: (left: 3pt + accent),
  breakable: true,
  [
    #text(fill: accent.darken(20%), weight: "bold", size: 9pt)[ASSIGNMENT]
    #v(-4pt)
    #set text(size: 9.8pt)
    #show raw: set text(size: 7.6pt)
    #body
    #v(3pt)
    #block(
      width: 100%,
      fill: rgb("#ffffff"),
      stroke: 0.5pt + accent.lighten(60%),
      radius: 2pt,
      inset: (x: 8pt, y: 6pt),
      [
        #text(size: 8pt, fill: muted, weight: "bold")[COMMIT AS]
        #v(-5pt)
        #raw("git commit -m \"" + name + "\"", lang: "bash", block: true)
      ],
    )
  ],
)

// A JavaScript / Go comparison, stacked rather than side by side.
//
// Two columns was the first attempt and it was wrong at this trim size: half
// of a 5.35in text block is about thirty characters, and indented Go wraps
// into nonsense well before that. Stacked gives each language the full
// measure, which is what the code was written against.
#let compare(js, go, jslabel: "JavaScript", golabel: "Go") = block(
  width: 100%,
  breakable: true,
  [
    #text(size: 8pt, fill: muted, weight: "bold")[#upper(jslabel)]
    #v(-3pt)
    #raw(js, lang: "js", block: true)
    #v(2pt)
    #text(size: 8pt, fill: muted, weight: "bold")[#upper(golabel)]
    #v(-3pt)
    #raw(go, lang: "go", block: true)
  ],
)

#let book(title: "", subtitle: "", author: "", version: "", body) = {
  set document(title: title, author: author)

  set page(
    width: 7in,
    height: 9.25in,
    margin: (inside: 0.9in, outside: 0.75in, top: 0.85in, bottom: 0.8in),
    binding: left,
  )

  set text(font: "Constantia", size: 10.5pt, fill: ink, lang: "en")
  set par(justify: true, leading: 0.68em, first-line-indent: 0pt, spacing: 0.9em)

  show link: it => text(fill: accent, it)

  // Headings. Chapters open a new right-hand page and carry a big number, so
  // flicking through the printed book lands you on a boundary.
  show heading.where(level: 1): it => {
    pagebreak(weak: true, to: "odd")
    chapternum.step()
    block(above: 0pt, below: 22pt)[
      // No hyphenation in a chapter title: "phrase-book" split across two
      // lines reads as a different word, and a title is short enough that a
      // ragged break is always the better one.
      #set text(font: "Cambria", hyphenate: false)
      #context if not backmatter.get() {
        text(size: 42pt, fill: accent.lighten(55%), weight: "bold")[
          #(chapternum.get().first() - 1)
        ]
      }
      #v(-14pt)
      #text(size: 21pt, weight: "bold")[#it.body]
      #v(3pt)
      #line(length: 100%, stroke: 0.8pt + accent)
    ]
  }
  show heading.where(level: 2): it => block(above: 17pt, below: 8pt)[
    #set text(font: "Cambria", size: 13.5pt, weight: "bold")
    #it.body
  ]
  show heading.where(level: 3): it => block(above: 13pt, below: 6pt)[
    #set text(font: "Cambria", size: 11.2pt, weight: "bold", fill: ink.lighten(15%))
    #it.body
  ]

  // Code. Measured rather than calculated: at 8.8pt a 72 column line wrapped
  // on the page, so the font-metric estimate of 80 was wrong by a tenth. At
  // 8.2pt the budget is 75 columns, which every listing in the book fits.
  show raw: set text(font: "Consolas", size: 8.2pt)
  show raw.where(block: true): it => block(
    width: 100%,
    fill: codebg,
    stroke: 0.5pt + codeline,
    radius: 3pt,
    inset: (x: 9pt, y: 8pt),
    breakable: true,
    it,
  )
  show raw.where(block: false): it => box(
    fill: codebg,
    inset: (x: 2.5pt, y: 0pt),
    outset: (y: 3pt),
    radius: 2pt,
    it,
  )

  // Title page
  page(margin: 0.9in, numbering: none, header: none, footer: none)[
    #v(1.6in)
    #line(length: 100%, stroke: 2.5pt + accent)
    #v(10pt)
    #text(font: "Cambria", size: 33pt, weight: "bold")[#title]
    #v(6pt)
    #text(font: "Cambria", size: 14.5pt, fill: muted)[#subtitle]
    #v(10pt)
    #line(length: 100%, stroke: 0.8pt + codeline)
    #v(1fr)
    #text(size: 11pt)[#author]
    #v(-5pt)
    #text(size: 9pt, fill: muted)[#version]
  ]

  // Front matter, numbered in roman so the body can start at page 1.
  set page(
    numbering: "i",
    footer: context align(center, text(size: 8.5pt, fill: muted, counter(page).display("i"))),
  )
  counter(page).update(1)

  block(above: 0pt)[
    #text(font: "Cambria", size: 19pt, weight: "bold")[Contents]
    #v(8pt)
    #line(length: 100%, stroke: 0.8pt + accent)
  ]
  v(8pt)
  show outline.entry.where(level: 1): it => {
    v(6pt, weak: true)
    text(weight: "bold", size: 10.5pt, it)
  }
  outline(title: none, depth: 2, indent: 1.1em)

  // Body
  set page(
    numbering: "1",
    header: context {
      let p = counter(page).get().first()
      if p <= 1 { return }
      let chapters = query(heading.where(level: 1))
      let current = none
      for c in chapters {
        if c.location().page() <= here().page() { current = c }
      }
      if current == none { return }
      set text(size: 8.5pt, fill: muted)
      grid(
        columns: (1fr, auto),
        align(left)[#current.body],
        align(right)[#counter(page).display("1")],
      )
      v(-6pt)
      line(length: 100%, stroke: 0.4pt + codeline)
    },
    footer: none,
  )
  counter(page).update(1)
  chapternum.update(0)

  set heading(numbering: none)
  body
}
