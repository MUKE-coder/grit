import type { FlowEdge, FlowNode } from '@/components/system-flow'

/**
 * The shape of a system design page.
 *
 * Grit is not one system. It is roughly thirty of them sharing a process: the
 * thing that decides who you are, the thing that decides what you may touch,
 * the thing that holds a queue, the thing that keeps a cache honest. The
 * reference documentation explains how to use each one. These pages explain
 * how each one is built and why, in the order a system design interview asks:
 * the problem, the requirements, the numbers, the architecture, the data, the
 * API, the classes, the scaling, and what breaks first.
 *
 * One data file and one template rather than thirty hand-written pages,
 * following config/deployment-guides.ts: the value of a set like this is that
 * every entry answers the SAME questions in the same order. Hand-written pages
 * drift, and by the sixth one a missing section is ambiguous between "does not
 * apply here" and "nobody got round to it".
 *
 * WHAT GOES IN HERE: what the system actually does, as built. Numbers that
 * were measured or can be derived from something stated. The components that
 * exist, by the names they have in a generated project.
 *
 * WHAT DOES NOT: aspiration. If a page describes a queue that is not there,
 * the page is worse than no page, because somebody will plan around it.
 */

/** One row of a two-column table. */
export type Pair = [string, string]

/** A worked calculation: the sum, then what it means. */
export interface Estimate {
  label: string
  /** The arithmetic, shown as given. */
  working: string[]
  note?: string
}

export interface ApiRow {
  method: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
  path: string
  what: string
}

export interface ApiGroup {
  title: string
  rows: ApiRow[]
}

export interface ClassSpec {
  name: string
  /** Where it lives in a generated project. */
  file?: string
  what: string
  /** The methods worth naming. */
  methods?: string[]
}

export interface Labelled {
  label: string
  text: string
}

export interface SystemDesign {
  /** URL segment under /docs/systems. */
  slug: string
  /** Page title, without the "System Design" suffix the template adds. */
  name: string
  /** One line for the index card and the meta description. */
  tagline: string
  /** Which family it belongs to, for grouping on the index. */
  group:
    | 'Identity'
    | 'Access'
    | 'Data'
    | 'Delivery'
    | 'Operations'
    | 'Platform'
    | 'Embedded'
  /** The Go package or packages in a generated project. */
  packages: string[]

  /**
   * The interview questions this page answers, shown before the problem
   * statement.
   *
   * These pages are read by two people. One wants to know how Grit works. The
   * other is preparing for an interview and wants to know whether this page
   * covers the thing they will be asked. The second reader cannot tell from a
   * table of contents, so the questions go at the top with a one-paragraph
   * answer each and a link to the section that works it through properly.
   *
   * The rule that keeps this honest: a question only belongs here if the page
   * genuinely answers it. A question with no section behind it is a promise
   * the page does not keep, which is worse than not listing it.
   */
  interview?: {
    intro?: string
    questions: {
      q: string
      /** The short answer, in full, for somebody who reads nothing else. */
      a: string
      /** The id of the section that covers it, for the "read more" link. */
      see?: string
    }[]
  }

  /** 1. The problem, in prose, then what the system must be able to do. */
  problem: { text: string[]; capabilities: string[] }

  /** 2. Requirements. */
  functional: string[]
  nonFunctional: Labelled[]

  /** 3. Capacity estimation. Assumptions, then the sums. */
  capacity?: {
    assumptions: Pair[]
    estimates: Estimate[]
  }

  /** 4. High level design. */
  highLevel: {
    intro?: string
    components: Labelled[]
    flow: {
      title: string
      nodes: FlowNode[]
      edges: FlowEdge[]
      /** The numbered steps, matching the badges on the diagram. */
      steps: string[]
    }
    dataFlow?: string[]
  }

  /** 5. What it is built from. */
  stack: Pair[]

  /** 6. The tables and what they hold. */
  dataModel?: {
    intro?: string
    entities: { name: string; fields: Pair[]; note?: string }[]
    storage?: string[]
  }

  /** 7. The HTTP surface. */
  api?: {
    groups: ApiGroup[]
    samples?: { title: string; language: string; code: string }[]
  }

  /** 8. Low level design. */
  lowLevel?: {
    classes: ClassSpec[]
    principles?: Labelled[]
    patterns?: Pair[]
  }

  /** 9. How it holds up. */
  scaling: string[]

  /** 10. What breaks first, and what to do about it. */
  bottlenecks: { problems: Labelled[]; improvements: Labelled[] }

  /** Where the reference documentation for using it lives. */
  seeAlso?: { title: string; href: string }[]
}
