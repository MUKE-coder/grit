/*
 * What CI actually does to each architecture, per shape.
 *
 * /docs/stability says how proven each subsystem is. It did not say which
 * shapes those subsystems are proven in, and the two are different questions:
 * the Go API is scaffolded, migrated and driven over HTTP on five databases on
 * every push, while two of the six architectures are not scaffolded by any
 * workflow at all. A reader picking `--single` deserves to know that before
 * they build on it, not after (grit#247).
 *
 * Every row here is checked against the workflow files by
 * scripts/check-ci-coverage.py: a shape this claims a workflow scaffolds, and
 * which no `grit new` line in that workflow produces, fails the build. The
 * claim and the thing it claims about are read in one pass, because a table
 * like this is exactly the kind that drifts into flattery.
 *
 * It is deliberately about scaffolding and building, which a workflow can be
 * read for. What a workflow then does with the project is in `detail`, which
 * is prose and is nobody's guarantee but ours.
 */

export type Depth = 'driven' | 'built' | 'generated'

export const depthStyle: Record<Depth, string> = {
  driven: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20',
  built: 'bg-sky-500/10 text-sky-400 border-sky-500/20',
  generated: 'bg-amber-500/10 text-amber-400 border-amber-500/20',
}

export const depthMeaning: Record<Depth, string> = {
  driven:
    'scaffolded, built, migrated and driven over HTTP against a real database, on every push',
  built:
    'scaffolded, installed, built, type-checked and linted, nightly, with no lockfile',
  generated:
    'covered by the generator and scaffolder tests, and by manual passes recorded below. No workflow scaffolds it.',
}

export interface ShapeCoverage {
  /** The flag that selects it, as `grit new` takes it. */
  flag: string
  label: string
  depth: Depth
  /**
   * The workflow files that scaffold this shape. Checked: a name here that
   * no `grit new` line in that file produces fails the build. Empty for a
   * shape no workflow scaffolds, which is the honest answer for three of
   * them.
   */
  workflows: string[]
  detail: string
}

export const SHAPE_COVERAGE: ShapeCoverage[] = [
  {
    flag: '--api',
    label: 'API only: the Go server, no frontend',
    depth: 'driven',
    workflows: ['ci.yml', 'live.yml', 'scan.yml'],
    detail:
      'The deepest coverage by a distance. Scaffolded and built on every push, then migrated and driven over HTTP: Postgres 15, 16 and 17, MySQL 8, SQLite and an in-memory database, reading the database directly where that is the only witness. gosec and Trivy run over the generated project, not only over the CLI.',
  },
  {
    flag: '--full',
    label: 'Everything: API, web, admin, docs, Expo and desktop',
    depth: 'built',
    workflows: ['canary.yml'],
    detail:
      'Scaffolded nightly with no lockfile, because a generated project resolves its dependencies at your install time and not at our release time. Every frontend is installed, built, type-checked and linted, and the docs paths every page names are checked against the project that run produced.',
  },
  {
    flag: '--triple --vite',
    label: 'API, web and admin on Vite and TanStack Router',
    depth: 'built',
    workflows: ['canary.yml'],
    detail:
      'The same nightly treatment. It is its own row because the Vite admin resolves a different dependency graph from the Next one, so a Next-only canary would miss half the exposure.',
  },
  {
    flag: '--triple --next',
    label: 'API, web and admin on Next.js',
    depth: 'built',
    workflows: ['canary.yml'],
    detail:
      'Covered by the --full run, which is a triple project with Expo and a desktop app added. This is the recommended shape and the one every tutorial assumes.',
  },
  {
    flag: '--double',
    label: 'API and web, no admin',
    depth: 'generated',
    workflows: [],
    detail:
      'No workflow scaffolds it. It is the triple shape with the admin left out, so most of what covers triple covers it, and the scaffolder tests assert which files each shape writes. Treat the frontend half as less proven than triple.',
  },
  {
    flag: '--single',
    label: 'One binary: the Go API with the React app embedded',
    depth: 'generated',
    workflows: [],
    detail:
      'No workflow scaffolds it. The embedding, the route groups and the admin section have their own scaffolder tests, and the shape was generated, run and driven through a browser in the eight-shape pass described below, which found a project whose admin could not reach its API in development.',
  },
  {
    flag: '--mobile',
    label: 'API and an Expo app',
    depth: 'generated',
    workflows: [],
    detail:
      'No workflow scaffolds this combination, though the Expo app itself is built nightly as part of --full. The API half is the --api row.',
  },
]
