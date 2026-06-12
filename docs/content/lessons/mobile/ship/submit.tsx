import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Store submission. Apple takes ~24 hours; Google takes ~2 hours.
        Both have their own gauntlets. <code>eas submit</code> automates the
        upload; the human review is still on you.
      </p>

      <h2>One command per store</h2>
      <CodeBlock
        terminal
        code={`eas build --platform all --profile production    # builds both
eas submit --platform ios                       # upload to App Store Connect
eas submit --platform android                   # upload to Play Console`}
      />
      <p>
        EAS handles the upload. The store dashboards (App Store Connect /
        Play Console) handle the metadata + review.
      </p>

      <h2>Apple App Store — the bureaucracy</h2>
      <p>You need:</p>
      <ul>
        <li>
          <strong>Apple Developer Program</strong> — $99/year. Required for
          any App Store distribution.
        </li>
        <li>
          <strong>App Store Connect</strong> account set up; a new app
          record matching your bundle ID.
        </li>
        <li>
          <strong>Screenshots</strong> — at least 6.5&apos;&apos; iPhone (1290 ×
          2796) and iPad if your app supports it.
        </li>
        <li>
          <strong>App description, keywords, support URL, privacy URL</strong>
        </li>
        <li>
          <strong>Privacy policy</strong> — Apple is strict. You need one
          even for free apps with no data collection.
        </li>
        <li>
          <strong>Data-collection disclosure</strong> — every SDK that
          collects anything (analytics, crash reporting, ads) must be
          declared.
        </li>
      </ul>

      <h2>Apple review — the gotchas</h2>
      <ul>
        <li>
          <strong>Reject reasons</strong> — &quot;app crashes on launch&quot;,
          &quot;feature doesn&apos;t work as described&quot;, &quot;mentions a competitor&quot;,
          &quot;requires login but doesn&apos;t describe what users can do without
          one&quot;. Read the rejection email carefully; resubmit.
        </li>
        <li>
          <strong>Submit a test account</strong> if your app requires login.
          Otherwise reviewers see only the login screen and reject for
          &quot;incomplete experience&quot;.
        </li>
        <li>
          <strong>Use Sign in with Apple</strong> if you offer other social
          login. This is enforced.
        </li>
      </ul>

      <TipBox tone="warning">
        <strong>TestFlight is your friend.</strong> Submit a build to
        TestFlight first (no review required for internal testing). Get 5
        users to install and use it. Fix the bugs. THEN submit to public
        review. You skip a lot of rejections.
      </TipBox>

      <h2>Google Play Store</h2>
      <p>Less strict, faster, but still has rules:</p>
      <ul>
        <li>
          <strong>Play Developer account</strong> — $25 one-time.
        </li>
        <li>
          <strong>Target SDK</strong> — Google forces you to target a
          recent SDK every year. Update the target in <code>app.json</code>{' '}
          when EAS warns.
        </li>
        <li>
          <strong>Data-safety disclosure</strong> — similar to Apple&apos;s
          privacy disclosure.
        </li>
        <li>
          <strong>Internal testing → Closed → Open → Production</strong> —
          the staged rollout track. Use Internal first.
        </li>
      </ul>

      <h2>Subscribe to store-rejection emails</h2>
      <p>
        Both stores email when they accept, reject, or have a question. Set
        up the right recipient (your support email, not your dev email);
        a missed rejection means days of delay.
      </p>

      <h2>Versioning across stores</h2>
      <p>
        Same code, same version, same build number — push the production
        build to BOTH stores at once. iOS and Android users see the same
        update on the same day.
      </p>

      <KnowledgeCheck
        question="Apple rejected your app saying 'reviewers couldn't sign in'. Your login works on every device you've tested. What's the most likely cause?"
        choices={[
          {
            label: 'A bug in your API',
            feedback:
              "Possible but rare — if it worked on your devices it'd work for reviewers. More likely you didn't give them a test account.",
          },
          {
            label: 'You didn\'t include a demo account in App Store Connect\'s reviewer notes',
            correct: true,
            feedback:
              "Right — Apple's reviewers need credentials. Add them in App Store Connect → App Information → App Review Information. Include email + password + any 2FA bypass instructions.",
          },
          {
            label: 'Reviewers are using a different country\'s App Store',
            feedback:
              "Plays a role for region-locked services, but the reviewer's email contains specific failures. Test account is the bigger lever.",
          },
          {
            label: 'Reviewers can\'t use Sign in with Apple',
            feedback:
              "They can. The issue is they don't have a working account at all.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              For the chapter assignment, you don&apos;t need to fully submit
              to a store — that takes days. Instead, write a{' '}
              <strong>submission readiness checklist</strong> for your app
              in <code>notes.md</code>:
            </p>
            <ul>
              <li>App icon (1024×1024) ready? ☐</li>
              <li>Screenshots (at least 6.5&apos;&apos; iPhone) ready? ☐</li>
              <li>App description (max 4000 chars) written? ☐</li>
              <li>Privacy policy URL live? ☐</li>
              <li>Test account credentials prepared for reviewers? ☐</li>
              <li>Bundle ID matches App Store Connect entry? ☐</li>
              <li>EAS production build succeeded? ☐</li>
            </ul>
            <p>
              Tick each off. Anything red means you&apos;d be rejected if you
              submitted today.
            </p>
          </>
        }
        hint={
          <>
            For solo devs, the missing piece is usually the privacy policy.
            Use <a href="https://www.iubenda.com" target="_blank" rel="noopener noreferrer" className="text-primary hover:underline">Iubenda</a> or write a basic one yourself.
          </>
        }
        solution={
          <>
            <p>
              No exact &quot;right&quot; solution here — this exercise is about
              checking each item against your real project. If you got all
              boxes ticked, run{' '}
              <code>eas submit --platform ios</code> and start the clock.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Final lesson — <strong>OTA updates</strong>. Push JS changes to all
        users without a store re-submission.
      </p>
    </>
  )
}
