import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        EAS Build is Expo&apos;s cloud build service. You push a config, EAS
        builds an iOS .ipa or Android .apk / .aab in its cloud, you
        download it. No Xcode on your machine required. This lesson
        covers the setup.
      </p>

      <h2>One-time install</h2>
      <CodeBlock terminal code={`npm install -g eas-cli
eas login   # uses your Expo account`} />

      <h2>The config — eas.json</h2>
      <p>
        Scaffolded into <code>apps/mobile/eas.json</code>:
      </p>
      <CodeBlock
        language="json"
        filename="apps/mobile/eas.json"
        code={`{
  "build": {
    "development": {
      "developmentClient": true,
      "distribution": "internal"
    },
    "preview": {
      "distribution": "internal",
      "android": { "buildType": "apk" }
    },
    "production": {
      "autoIncrement": true
    }
  },
  "submit": {
    "production": {}
  }
}`}
      />
      <p>Three profiles to know:</p>
      <ul>
        <li>
          <strong>development</strong> — dev build, JS hot-reloads live.
          Use during day-to-day dev.
        </li>
        <li>
          <strong>preview</strong> — APK / IPA you send to testers.
          Distributed via internal links.
        </li>
        <li>
          <strong>production</strong> — what goes to the App Store / Play
          Store. Auto-increments build number.
        </li>
      </ul>

      <h2>Your first preview build</h2>
      <CodeBlock terminal code={`# from apps/mobile/
eas build --platform android --profile preview`} />
      <p>
        First run prompts you to:
      </p>
      <ol>
        <li>Configure a project ID (one-time per project)</li>
        <li>
          For iOS: log in to Apple Developer account so EAS can sign with
          your cert
        </li>
        <li>
          For Android: pick auto-managed signing or upload your keystore
        </li>
      </ol>
      <p>
        Then it queues the build. You get a URL to track progress. ~10
        minutes later, downloadable APK / IPA.
      </p>

      <TipBox tone="info">
        <strong>Free tier:</strong> Expo gives you ~30 builds/month free,
        with queue priority on slow lanes. Paid tier ($29/mo) is faster
        queues + more builds. For a solo dev this is plenty.
      </TipBox>

      <h2>Installing the preview build</h2>
      <ul>
        <li>
          <strong>Android APK:</strong> download, copy to phone, install
          (enable &quot;install from unknown sources&quot;). Or use the EAS QR.
        </li>
        <li>
          <strong>iOS:</strong> you need TestFlight (covered in the next
          lesson) — APKs don&apos;t exist on iOS for unsigned distribution.
        </li>
      </ul>

      <h2>Versioning</h2>
      <p>
        Two numbers matter:
      </p>
      <ul>
        <li>
          <strong>version</strong> (in <code>app.json</code>) — what the
          user sees (&quot;v1.2.3&quot;). Bump manually for marketing-visible
          releases.
        </li>
        <li>
          <strong>build number / version code</strong> — internal
          monotonically increasing integer. EAS auto-increments with the
          production profile.
        </li>
      </ul>
      <p>
        The store rejects builds with duplicate build numbers. Let EAS
        manage this.
      </p>

      <h2>What gets built into the binary</h2>
      <p>
        Build-time environment is locked in: <code>app.json</code>,
        <code>app.config.ts</code>, and any <code>process.env</code>{' '}
        baked in via <code>extra</code> in app.json. Runtime config (API
        URL, feature flags) flows through <code>Constants.expoConfig.extra</code>:
      </p>
      <CodeBlock
        language="ts"
        filename="apps/mobile/app.config.ts"
        code={`export default ({ config }) => ({
  ...config,
  extra: {
    apiUrl: process.env.EAS_API_URL ?? 'http://localhost:8080',
    sentryDsn: process.env.EAS_SENTRY_DSN,
  },
})`}
      />
      <p>
        Set the env vars in EAS dashboard or with{' '}
        <code>eas secret:create</code> for per-profile secrets.
      </p>

      <KnowledgeCheck
        question="You ran `eas build --profile production` but you forgot to bump the version. The build succeeds. What happens when you try to upload to TestFlight?"
        choices={[
          {
            label: 'TestFlight accepts it and overwrites the previous build',
            feedback:
              "Wrong — Apple specifically rejects duplicate build numbers. You'd see an upload error.",
          },
          {
            label: 'TestFlight rejects it with \'duplicate CFBundleVersion\'',
            correct: true,
            feedback:
              "Right — Apple's rule is that each upload needs a unique CFBundleVersion (build number). EAS's production profile auto-increments to avoid this; if you bypassed it, you have to bump manually.",
          },
          {
            label: 'TestFlight prompts you to choose which build to keep',
            feedback:
              "Wrong — Apple doesn't ask. It just rejects.",
          },
          {
            label: 'The build expires after 30 days unused',
            feedback:
              "TestFlight builds DO expire after 90 days, but that's unrelated to upload acceptance.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Produce your first preview build:</p>
            <ol>
              <li>
                Run <code>npm install -g eas-cli</code> if you haven&apos;t.
              </li>
              <li><code>eas login</code></li>
              <li><code>eas build:configure</code> (one-time setup)</li>
              <li><code>eas build --platform android --profile preview</code></li>
              <li>Wait ~10 min, download the APK, install on your phone.</li>
              <li>
                Confirm the app launches and looks identical to the Expo Go
                version.
              </li>
            </ol>
            <p>Paste the EAS build URL in <code>notes.md</code>.</p>
          </>
        }
        hint={
          <>
            If the build fails, the EAS dashboard shows the log. Most
            failures are missing config (project ID, app slug) or a native
            module incompatibility — fix the config and re-run.
          </>
        }
        solution={
          <>
            <p>
              The build URL looks like:
            </p>
            <CodeBlock
              language="text"
              code={`https://expo.dev/accounts/<you>/projects/field-app/builds/<id>`}
            />
            <p>
              Click it — you&apos;ll see the queue position, build log, and
              download link. Install the APK and the app should run
              identically to dev mode (just slower to launch since it&apos;s a
              real binary).
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Next lesson — submitting your build to the App Store and Play
        Store. The actual ship.
      </p>
    </>
  )
}
