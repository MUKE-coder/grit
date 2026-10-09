# What to build next in Grit

Research, October 2026. The question was what tool would make Grit more
adoptable or more enterprise-ready. Every claim is linked, and where I checked
the repo rather than the web it says so.

The first thing to record is that my opening instinct was wrong: I went looking
for the gap Laravel filled with Boost, and **Grit already has an MCP server
with ten tools.** So this is a gap analysis against something that exists, not
a proposal to start one.

---

## 1. What the market did

Laravel shipped [Boost](https://laravel.com/blog/announcing-laravel-boost): a
framework-aware MCP server giving agents "a set of tools, versioned guidelines,
and version-specific documentation so they behave like an experienced Laravel
developer." Its
[15+ tools](https://learnwithmux.com/posts/mastering-the-mcp-toolkit-a-laravel-boost-deep-dive/)
are application info, search docs, Tinker, browser logs, database queries,
database schema, list Artisan commands, last errors, list routes, read
configuration, read log entries and report feedback.

The convention layer standardised while nobody was looking:

- [AGENTS.md](https://blog.buildbetter.ai/agents-md-complete-guide-for-engineering-teams-in-2026/)
  is stewarded by the Agentic AI Foundation, a Linux Foundation project. 28,000
  repositories at launch, **60,000+ by mid-2026**, read natively by Claude
  Code, Codex CLI, Cursor, Aider, Devin, Copilot, Gemini CLI, Windsurf and
  Amazon Q.
- [Agent Skills](https://atlan.com/know/ai-agent/ai-agent-skills/what-are-agent-skills/)
  are supported by 40+ platforms, with around 490,000 skills published.

**I found no Go framework equivalent of Boost.** That is the open position.

It matters more than a feature because of how applications are written now. A
framework an agent uses correctly on the first attempt is worth more than a
framework with more features, and for a project with one maintainer it is
leverage of exactly the right kind: the agent reads the guidelines instead of
asking the maintainer.

---

## 2. What Grit already has

Checked in the repo, not assumed.

`internal/mcp` is about 2,400 lines including tests, and `grit mcp` exposes ten
tools:

| tool | what it answers |
|---|---|
| `grit_project_info` | version, architecture, which apps exist |
| `grit_list_routes` | the route surface |
| `grit_describe_models` | models, from the Go structs |
| `grit_list_resources` | generated resources |
| `grit_file_ownership` | who owns a file, framework or you |
| `grit_doctor` | the health checks |
| `grit_list_permissions` | the permission strings |
| `grit_env_keys` | configuration keys |
| `grit_cli_reference` | the CLI surface |
| `grit_generate_resource` | run the generator |

`grit_file_ownership` has no counterpart in Boost and is the more thoughtful
tool of the two sets: it is the one that stops an agent editing a file that
`grit upgrade` will overwrite.

A generated project also gets `CLAUDE.md` and `AGENTS.md`, from
`scaffold.WriteAgentsDoc`: about 96 lines across eleven sections, and the
release workflow produces an SBOM.

So the foundation is there and is further along than Laravel's in one respect.

---

## 3. The four gaps, in order

### 3.1 Version-pinned documentation search

Boost's headline tool and the one Grit does not have. Agents "use search docs
to find the correct API for the exact package version."

Grit is the framework that needs this most, and the reason is its own release
cadence. It is at v3.393 and ships several versions a week. An agent working in
a project pinned to v3.350 that reads today's docs is confidently wrong, and
the symptom is generated code that references a helper that did not exist yet
or a flag that has since been renamed. Nobody else has this problem as acutely,
which makes it both the largest exposure and the least contested thing to fix.

The docs site already carries a single-source version. The missing piece is
retrieval over MCP, scoped to the version in the project's `grit.json`.

### 3.2 Logs, errors and the running application

Boost reads application logs and browser logs, and has a `last_errors` tool.
Grit has none of these, and the evidence that it matters is this week's work on
a Grit-built application: every protocol bug that mattered was found by reading
a log file by hand. `NoData` in the simple query protocol, the SASL exchange
failing on a parse, the router answering `N` to every TLS request. Each one was
invisible in the exit code and obvious in the log.

The project's own standing rule is "look at the running admin, screenshots
before calling a UI change done." That ritual is entirely manual. Nothing
exposes it to an agent, so an agent cannot follow the one rule the project
considers most important.

`grit_describe_models` reads the Go structs, which is not the same as what the
database actually has after a migration. Grit already embeds GORM Studio, which
answers live schema over HTTP: the capability exists and is not wired to the
agent.

### 3.3 The conventions that ship with nothing

This is the asymmetry worth naming. The Grit repository has excellent scoped
`AGENTS.md` files for people working **on** Grit: one each for
`internal/scaffold`, `internal/generate`, `cmd/grit`, `internal/ui`, `docs` and
`examples`, each covering the boundary rules that are easy to break there.

A project built **with** Grit gets 96 lines.

And the framework's hardest-won knowledge is in neither. There are roughly
sixty accumulated conventions that are not written down anywhere a user or an
agent can read: that an icon in the admin's `iconMap` is not automatically a
named export from `@/lib/icons`; that a change to login needs six edits across
six auth styles; that a framework-owned file needs registering in both the
new-project map and the upgrade writer; that `grit migrate` runs from
`apps/api` while a built binary runs from anywhere, so a relative `.env` path
gives a project two SQLite databases; that `.npmrc`'s `node-linker` is ignored
when a `pnpm-workspace.yaml` exists.

Every one of those was learned by something breaking. They are the most
valuable documentation the project has and they ship with nothing. Promoting
them into the generated `AGENTS.md`, into skills, and where possible into
`grit doctor` checks is pure extraction of an asset that already exists.

### 3.4 The single-maintainer answer

The documented first barrier to enterprise adoption of a project like this is
not features. [Almost 25% of open source projects have one developer
contributing code and 94% have ten or fewer](https://www.atlanticcouncil.org/in-depth-research-reports/report/open-source-software-as-infrastructure/),
the truck factor is assessed in procurement, and XZ Utils made maintainer
burnout a known attack vector rather than a theoretical one. The adoption
question underneath it is the blunt one: ["why should I learn another
one?"](https://dev.to/md8_habibullah/stop-picking-the-wrong-framework-here-are-the-only-10-that-matter-in-2026-32gm)

That cannot be fixed with code, and it can be answered with evidence. A Grit
application is a Go application and a Next.js application. The answer to "what
if Grit stops" is a test that deletes the Grit CLI, builds a generated project,
runs its suite, and passes. Done in CI, on every release, it converts the
largest objection into a checked claim.

It is also the cheapest item on this list.

---

## 4. Also real, and lower

**The EU Cyber Resilience Act.** Vulnerability and incident reporting was due
[11 September 2026](https://bearingpoint.services/foss/en/newsblogs/2026-sbom-minimum-elements-software-bill-of-materials-compliance/),
and SBOM production and maintenance take effect 27 December 2027. The 2026
minimum elements added author signature, component hash and component licence.
Grit signs and produces an SBOM for itself; whether a **generated** application
does is the question worth answering, because that is the artefact a customer's
procurement asks about. A real checkbox with a date, and narrower and less
urgent than the four above.

**What not to chase.** Managed-Postgres complaints are complaints about
providers: Neon's
[500-900 ms cold starts](https://dev.to/thiago_alvarez_a7561753aa/neon-vs-supabase-2026-database-or-backend-the-real-tradeoffs-3ggn),
Supabase's lock-in. A framework cannot fix a cold start.

---

## 5. Recommendation

In this order, because the first two are where the evidence is strongest and
the last is the cheapest:

1. **`grit docs` over MCP, pinned to the project's version.** The one tool
   Boost has that Grit does not, and the one Grit's release cadence makes
   urgent.
2. **Logs, last errors and live schema over MCP.** Cheap, and this week's bug
   hunt is the argument: the failures that mattered lived in log files and
   nothing could read them but a person.
3. **Promote the sixty conventions into shipped, versioned guidelines and
   doctor checks.** Extraction, not invention. The knowledge exists.
4. **The eject proof in CI.** One job that deletes the CLI and proves a
   generated project still builds and passes. Answers the first enterprise
   objection with a green tick instead of a paragraph.

The unifying claim, and the one worth putting on the homepage: **the first Go
framework an agent already knows how to use.** It is nearly true today, which
is the best reason to finish it.
