import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'

export default function Lesson() {
  return (
    <>
      <p>
        k6 is a single Go binary. No Node. No Python. No dependencies.
        Three install methods covered here — pick the one that matches
        your OS.
      </p>

      <h2>macOS — Homebrew</h2>
      <CodeBlock
        language="bash"
        code={`brew install k6`}
      />

      <h2>Linux — apt / yum / pacman</h2>
      <CodeBlock
        language="bash"
        code={`# Debian / Ubuntu
sudo gpg -k
sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
sudo apt-get update
sudo apt-get install k6

# Arch
sudo pacman -S k6`}
      />

      <h2>Windows — Chocolatey or Scoop</h2>
      <CodeBlock
        language="powershell"
        code={`# Chocolatey
choco install k6

# Or Scoop
scoop install k6`}
      />

      <h2>Docker — if you don&apos;t want to install</h2>
      <CodeBlock
        language="bash"
        code={`docker run --rm -i grafana/k6 run - <script.js`}
      />
      <p>
        The <code>-i</code> pipes your local script into the container.
        Useful in CI; clunky in dev. Native install is usually nicer.
      </p>

      <h2>Verify</h2>
      <CodeBlock
        language="bash"
        code={`k6 version
# k6 v0.50.0 (2024-04-18T...)`}
      />
      <p>
        Any 0.50+ works. Most APIs in this course assume features
        from 0.49 or later.
      </p>

      <h2>The k6 binary on its own</h2>
      <p>
        k6 has no UI. It&apos;s pure CLI. You write a JavaScript file
        (k6 has its own JS runtime, NOT Node), then run{' '}
        <code>k6 run script.js</code>. Output goes to your terminal.
        That&apos;s the entire mental model.
      </p>

      <TipBox tone="info">
        <strong>k6 cloud is optional.</strong> Grafana also runs a SaaS
        for k6 — pretty dashboards, team scaling. You don&apos;t need it
        to learn or use k6. Skip it; come back if you need centralised
        results.
      </TipBox>

      <h2>Where to keep your scripts</h2>
      <CodeBlock
        language="text"
        code={`my-grit-app/
└── tests/
    └── k6/
        ├── smoke.js
        ├── load.js
        ├── stress.js
        ├── spike.js
        ├── soak.js
        └── lib/
            ├── auth.js          ← reusable: log in, get a token
            └── thresholds.js    ← shared SLO thresholds`}
      />
      <p>
        Tests live with code. Same repo, version-controlled. Anyone can
        run them.
      </p>

      <Exercise
        prompt={
          <>
            <p>Confirm installation:</p>
            <ol>
              <li>Install k6 via the method for your OS.</li>
              <li>Run <code>k6 version</code> and capture the output.</li>
              <li>
                Make a folder <code>tests/k6/</code> in your Grit project
                — we&apos;ll put real scripts there next lesson.
              </li>
              <li>Paste the version output into <code>notes.md</code>.</li>
            </ol>
          </>
        }
        solution={
          <>
            <p>
              That&apos;s it. Next lesson we&apos;ll write a 20-line
              script that hits your local API.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next — <strong>Your first k6 script</strong>. A real test
        against your running Grit API.
      </p>
    </>
  )
}
