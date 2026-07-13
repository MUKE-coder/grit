import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import { LaneFlow } from '@/components/lane-flow'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'CI/CD with GitHub Actions: Automated Testing & Deployment',
  description: 'Set up automated testing and deployment pipelines for Grit projects using GitHub Actions. CI workflows, release automation, branch protection, and environment deploys.',
}

export default function CICDGitHubCourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">CI/CD with GitHub Actions</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            CI/CD with GitHub Actions: Automated Testing & Deployment
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Set up automated testing and deployment for your Grit projects. You{"'"}ll understand the CI/CD
            workflows Grit scaffolds, extend them with frontend tests, configure automated deployments,
            and set up branch protection for a professional development workflow.
          </p>
          <LaneFlow id="c-cicd" lanes={['git push', 'GitHub Actions', 'Result']}
            nodes={[
              { id: 'push', lane: 0, row: 1, title: 'git push', sub: 'PR / main', tone: 'primary' },
              { id: 'test', lane: 1, row: 0, title: 'Test', sub: 'go + vitest', tone: 'cyan' },
              { id: 'build', lane: 1, row: 1, title: 'Build', sub: 'cross-platform', tone: 'blue' },
              { id: 'deploy', lane: 2, row: 1, title: 'Deploy', sub: 'on main', tone: 'green' },
            ]}
            edges={[{ from: 'push', to: 'test', tone: 'cyan' }, { from: 'push', to: 'build', label: 'runs', tone: 'blue' }, { from: 'build', to: 'deploy', tone: 'green' }]}
            legend={[{ tone: 'cyan', label: 'Test + build' }, { tone: 'green', label: 'Auto-deploy' }]}
            caption="Every push runs tests and builds; merges to main deploy automatically" />
        </div>

        <div className="my-4 rounded-lg border border-primary/20 bg-primary/5 px-4 py-3">
          <p className="text-sm text-muted-foreground">
            <strong className="text-foreground">Reference docs:</strong>{' '}
            <a href="/docs/infrastructure/deployment" className="text-primary hover:underline">Deployment →</a>
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What is CI/CD? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What is CI/CD?</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Without CI/CD, shipping code looks like this: you write code, run tests manually (maybe),
            build the project, SSH into a server, pull the code, restart services, and pray nothing
            breaks. With CI/CD, you push code and everything else happens automatically.
          </p>

          <Definition term="CI (Continuous Integration)">
            Automatically testing code every time it{"'"}s pushed to the repository. Every push triggers
            a build and test suite. If tests fail, the team knows immediately — before the code reaches
            production. CI catches bugs early.
          </Definition>

          <Definition term="CD (Continuous Deployment)">
            Automatically deploying code after tests pass. When CI confirms the code is good, CD ships
            it to staging or production without manual intervention. No more SSH-and-pray deployments.
          </Definition>

          <Definition term="Pipeline">
            A series of automated steps that run in sequence: test, build, deploy. Each step must pass
            before the next one runs. If any step fails, the pipeline stops and notifies the team.
            Think of it as an assembly line for code.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            A typical CI/CD pipeline:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">1.</span> Developer pushes code to GitHub</li>
            <li className="flex gap-2"><span className="text-primary">2.</span> CI runs Go tests (<Code>go test -race ./...</Code>)</li>
            <li className="flex gap-2"><span className="text-primary">3.</span> CI runs frontend tests (<Code>pnpm test</Code>)</li>
            <li className="flex gap-2"><span className="text-primary">4.</span> If all tests pass, CD deploys to staging</li>
            <li className="flex gap-2"><span className="text-primary">5.</span> On tag push (v1.0.0), CD deploys to production</li>
          </ul>

          <Challenge number={1} title="Explain CI/CD">
            <p>Explain CI/CD in your own words. Why is it important? What happens when a team doesn{"'"}t
            use CI/CD? Think about: how bugs reach production, how long deployments take, and how
            confident developers feel pushing code on a Friday afternoon.</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: Grit's Scaffolded CI ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Grit{"'"}s Scaffolded CI</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When you run <Code>grit new</Code>, it scaffolds a <Code>.github/workflows/ci.yml</Code> file.
            This is a complete CI workflow that runs automatically on every push and pull request.
          </p>

          <CodeBlock filename=".github/workflows/ci.yml">
{`name: CI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - run: cd apps/api && go test -race ./...`}
          </CodeBlock>

          <Definition term="GitHub Actions">
            GitHub{"'"}s built-in CI/CD platform. It runs jobs on GitHub{"'"}s servers, triggered by events
            like push, pull request, schedule, or manual dispatch. Free for public repos, with generous
            free minutes for private repos. No external CI service needed.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This workflow does 3 things:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Checks out your code</strong> — downloads the repo to the runner</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Sets up Go</strong> — installs the specified Go version</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Runs tests</strong> — executes all Go tests with race detection enabled</li>
          </ul>

          <Tip>
            The <Code>-race</Code> flag enables Go{"'"}s race detector, which catches concurrent access
            bugs. It{"'"}s slower than normal tests but catches real production issues. Always use it in CI.
          </Tip>

          <Challenge number={2} title="Find Your CI Workflow">
            <p>Find the <Code>.github/workflows/</Code> directory in your Grit project. What workflow
            files exist? Open <Code>ci.yml</Code> and read through it. What triggers the workflow?
            What Go version does it use?</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Understanding the Workflow ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Understanding the Workflow</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            YAML workflows can look intimidating. Let{"'"}s break down every section so you can read
            and write them confidently.
          </p>

          <CodeBlock filename="Workflow Anatomy">
{`# name: Human-readable name shown in GitHub UI
name: CI

# on: Events that trigger this workflow
# [push, pull_request] = run on every push AND every PR
on: [push, pull_request]

# jobs: Groups of steps that run on a fresh virtual machine
jobs:
  # "test" is the job ID (you choose the name)
  test:
    # runs-on: The OS for the virtual machine
    runs-on: ubuntu-latest

    # steps: Commands that run sequentially
    steps:
      # uses: Run a pre-built action (like a plugin)
      - uses: actions/checkout@v4

      # with: Configuration for the action
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'

      # run: Execute a shell command
      - run: cd apps/api && go test -race ./...`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Key concepts:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Jobs run in parallel</strong> by default — unless you use <Code>needs</Code> to create dependencies</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Steps run sequentially</strong> within a job — if one fails, the rest are skipped</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Each job gets a fresh VM</strong> — nothing persists between jobs unless you use artifacts</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Actions are reusable</strong> — <Code>actions/checkout@v4</Code> is maintained by GitHub</li>
          </ul>

          <Challenge number={3} title="Read the CI Workflow">
            <p>Read your <Code>ci.yml</Code> file carefully. Answer these questions: (1) What Go version
            does it use? (2) What test flags are set? (3) What events trigger it? (4) What operating
            system does the runner use? (5) How many jobs are defined?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Adding Frontend Tests ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Adding Frontend Tests</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The scaffolded CI only tests the Go API. For a complete pipeline, add a separate job
            that tests the frontend — Vitest for unit tests, Playwright for end-to-end tests.
          </p>

          <CodeBlock filename="Frontend Test Job">
{`  test-frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - run: npm install -g pnpm
      - run: pnpm install
      - run: pnpm test`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This job runs in parallel with the Go test job. Both must pass for the overall CI to
            be green. If either fails, the PR gets a red X.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            For end-to-end tests with Playwright, you need a running API. Use a service container:
          </p>

          <CodeBlock filename="E2E Test Job (Advanced)">
{`  test-e2e:
    runs-on: ubuntu-latest
    needs: [test, test-frontend]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - run: npm install -g pnpm
      - run: pnpm install
      - run: npx playwright install --with-deps
      - run: pnpm test:e2e`}
          </CodeBlock>

          <Note>
            The <Code>needs: [test, test-frontend]</Code> line means E2E tests only run after both
            unit test jobs pass. No point running expensive browser tests if basic tests fail.
          </Note>

          <Challenge number={4} title="Add Frontend Tests to CI">
            <p>Add a <Code>test-frontend</Code> job to your CI workflow. It should set up Node.js 20,
            install pnpm, install dependencies, and run <Code>pnpm test</Code>. Push the change and
            check the Actions tab — do both jobs run?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: The Release Workflow ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The Release Workflow</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit also scaffolds a <Code>release.yml</Code> workflow. Unlike CI (which runs on every push),
            the release workflow only triggers when you push a version tag like <Code>v1.0.0</Code>. It
            builds cross-platform binaries and creates a GitHub Release with downloadable files.
          </p>

          <CodeBlock filename=".github/workflows/release.yml (Simplified)">
{`name: Release
on:
  push:
    tags: ['v*']

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      # Build for Linux, macOS, and Windows
      # Create GitHub Release with binaries attached`}
          </CodeBlock>

          <Definition term="GitHub Release">
            A tagged version of your project with downloadable binaries. Users install your app by
            downloading from the release page. Each release has a version number (v1.0.0), release
            notes, and attached files (binaries, source archives).
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            To create a release:
          </p>

          <CodeBlock filename="Creating a Release">
{`# Tag the commit
git tag v1.0.0

# Push the tag (triggers the release workflow)
git push origin v1.0.0

# GitHub Actions will:
# 1. Build binaries for linux/amd64, darwin/amd64, darwin/arm64, windows/amd64
# 2. Create a GitHub Release
# 3. Attach the binaries as downloadable assets`}
          </CodeBlock>

          <Challenge number={5} title="Read the Release Workflow">
            <p>Open <Code>release.yml</Code> and answer: (1) What event triggers it? (2) What platforms
            does it build for? (3) How would you create a release? Write the git commands you{"'"}d run
            to release version 1.0.0.</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Automated Deployment ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Automated Deployment</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The most powerful CI/CD feature: automatic deployment. After tests pass, the deploy job
            ships your code to production. No SSH, no manual steps, no forgetting to restart the service.
          </p>

          <CodeBlock filename="Deploy Job">
{`  deploy:
    needs: [test]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      # Run: grit deploy --host DEPLOY_HOST --domain DEPLOY_DOMAIN
      # Uses secrets for sensitive values`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Key parts of the deploy job:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>needs: [test]</Code> — only runs after tests pass</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <Code>if: github.ref == {"'"}refs/heads/main{"'"}</Code> — only deploys from the main branch (not PRs or feature branches)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Secrets store sensitive values like server addresses and SSH keys</li>
          </ul>

          <Definition term="GitHub Secrets">
            Encrypted environment variables stored in your repository settings. They{"'"}re never exposed
            in logs, even if a workflow prints them. Use secrets for: API keys, deploy credentials,
            SSH keys, database passwords, and any sensitive configuration.
          </Definition>

          <Tip>
            To add secrets: Repository Settings, then Secrets and variables, then Actions, then
            New repository secret. Name it something clear like <Code>DEPLOY_HOST</Code> or
            <Code>SSH_PRIVATE_KEY</Code>.
          </Tip>

          <Challenge number={6} title="List Your Deploy Secrets">
            <p>What secrets would you need for automated deployment? List them all. Think about:
            server address, domain name, SSH credentials, database connection, and any API keys
            your app needs in production.</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Branch Protection ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Branch Protection</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            CI is only useful if you enforce it. Branch protection rules prevent merging code that
            fails CI. No green checkmark, no merge — period.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            To set up branch protection:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">1.</span> Go to Repository Settings</li>
            <li className="flex gap-2"><span className="text-primary">2.</span> Click Branches in the sidebar</li>
            <li className="flex gap-2"><span className="text-primary">3.</span> Add a branch protection rule for <Code>main</Code></li>
            <li className="flex gap-2"><span className="text-primary">4.</span> Enable {'"'}Require status checks to pass before merging{'"'}</li>
            <li className="flex gap-2"><span className="text-primary">5.</span> Select the CI jobs (test, test-frontend) as required checks</li>
            <li className="flex gap-2"><span className="text-primary">6.</span> Optionally require pull request reviews</li>
          </ul>

          <p className="text-muted-foreground leading-relaxed mb-4">
            With this setup, the development workflow becomes:
          </p>

          <CodeBlock filename="Protected Branch Workflow">
{`# 1. Create a feature branch
git checkout -b feature/add-comments

# 2. Write code, commit, push
git push origin feature/add-comments

# 3. Open a Pull Request on GitHub
# CI runs automatically on the PR

# 4. If CI passes: green checkmark, merge allowed
# If CI fails: red X, merge blocked

# 5. After merge to main: deploy job runs automatically`}
          </CodeBlock>

          <Challenge number={7} title="Set Up Branch Protection">
            <p>Set up branch protection on your main branch. Require the CI test job to pass before
            merging. Create a feature branch, make a change, push it, and open a PR. Does CI run?
            Can you merge before CI passes?</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Environment-Specific Deploys ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Environment-Specific Deploys</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Production apps need a staging environment — a copy of production where you test changes
            before they go live. Different branches deploy to different environments:
          </p>

          <CodeBlock filename="Multi-Environment Deploy">
{`  deploy-staging:
    needs: [test]
    if: github.ref == 'refs/heads/staging'
    runs-on: ubuntu-latest
    steps:
      # Deploy to staging.myapp.com
      # Uses STAGING_HOST and STAGING_DOMAIN secrets

  deploy-production:
    needs: [test]
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    steps:
      # Deploy to myapp.com
      # Uses PRODUCTION_HOST and PRODUCTION_DOMAIN secrets`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The branching strategy:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Feature branches</strong> — where you write code. CI runs tests.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">staging branch</strong> — merge features here first. Auto-deploys to staging.myapp.com for testing.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">main branch</strong> — merge from staging when ready. Auto-deploys to myapp.com (production).</li>
          </ul>

          <Note>
            Staging should mirror production as closely as possible — same OS, same database engine,
            same environment variables (with different values). If it works on staging, it should
            work in production.
          </Note>

          <Challenge number={8} title="Design Your Branching Strategy">
            <p>Design a branching strategy for your project. Which branch deploys where? How does
            code flow from feature to staging to production? Draw the flow: feature branch, then PR
            to staging, then test, then PR to main, then production deploy.</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Notifications ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Notifications</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            When CI fails at 2 AM, you want to know. When a deployment succeeds, the team should
            celebrate. Notifications close the feedback loop — the pipeline tells you what happened.
          </p>

          <CodeBlock filename="Slack Notification on Failure">
{`  # Add this step at the end of any job
  - name: Notify on failure
    if: failure()
    # Send a POST request to Slack webhook
    # Body: "CI failed on branch main - commit abc123"
    # The webhook URL is stored as a secret`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Common notification patterns:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Test failure</strong> — urgent notification to Slack/Discord with the failing branch and commit</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Successful deploy</strong> — info notification: {'"'}v1.2.0 deployed to production{'"'}</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">New release</strong> — announcement: {'"'}Release v1.2.0 published with 5 new features{'"'}</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Deploy failure</strong> — critical alert: {'"'}Production deploy failed — rollback may be needed{'"'}</li>
          </ul>

          <Tip>
            GitHub also sends email notifications by default. For teams, Slack or Discord webhooks
            are more visible. Everyone sees the notification in the team channel, not buried in
            individual inboxes.
          </Tip>

          <Challenge number={9} title="Plan Your Notifications">
            <p>What notifications would you want for these events? (1) Test failure on a PR,
            (2) Successful deploy to staging, (3) Successful deploy to production, (4) New GitHub
            Release created. For each, specify: who should be notified, how (Slack, email, Discord),
            and what the message should say.</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s everything you learned in this course:
          </p>

          <ul className="space-y-3 mb-6">
            {[
              'CI automatically tests code on every push — catching bugs before they reach production',
              'CD automatically deploys code after tests pass — no manual SSH deployments',
              'Grit scaffolds ci.yml (test on push/PR) and release.yml (build on tag push)',
              'GitHub Actions runs jobs on GitHub servers triggered by events (push, PR, tag)',
              'YAML workflows define jobs, steps, triggers, and environment configuration',
              'Frontend tests run as a separate parallel job alongside Go tests',
              'GitHub Secrets store sensitive values like deploy credentials securely',
              'Branch protection prevents merging code that fails CI',
              'Environment-specific deploys: staging branch to staging server, main to production',
              'Notifications keep the team informed of failures, deploys, and releases',
            ].map((item, i) => (
              <li key={i} className="flex gap-3 text-muted-foreground">
                <CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" />
                <span>{item}</span>
              </li>
            ))}
          </ul>

          <Challenge number={10} title="Build the Complete Pipeline (Part 1)">
            <p>Set up a complete CI workflow with two jobs: Go tests (<Code>go test -race ./...</Code>)
            and frontend tests (<Code>pnpm test</Code>). Push it to GitHub and verify both jobs run
            on the Actions tab.</p>
          </Challenge>

          <Challenge number={11} title="Build the Complete Pipeline (Part 2)">
            <p>Add a deploy job that runs after both test jobs pass. It should only trigger on pushes
            to the main branch. Use GitHub Secrets for the deploy host and domain. What secrets did
            you create?</p>
          </Challenge>

          <Challenge number={12} title="Build the Complete Pipeline (Part 3)">
            <p>Set up the full workflow: branch protection on main requiring CI to pass, a staging
            branch with auto-deploy to staging, and tag-triggered releases. Test the complete flow:
            feature branch, then PR, then CI passes, then merge to staging, then staging deploy, then
            merge to main, then production deploy. This is a professional-grade CI/CD setup.</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
