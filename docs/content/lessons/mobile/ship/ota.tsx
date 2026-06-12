import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        OTA = over-the-air updates. Ship a JS bundle change to every user
        without re-submitting to the store. Your apps update overnight while
        users sleep. This is the Expo superpower most teams underuse.
      </p>

      <h2>What can update OTA, what can&apos;t</h2>
      <p>
        Your app has two layers:
      </p>
      <ul>
        <li>
          <strong>Native code</strong> — installed via APK/IPA. Changing it
          requires a new store submission.
        </li>
        <li>
          <strong>JS bundle</strong> — your React Native code, components,
          styles. Updates OTA.
        </li>
      </ul>
      <p>What that means in practice:</p>

      <div className="overflow-x-auto my-5">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b border-border/40">
              <th className="text-left px-3 py-2 font-medium">Change</th>
              <th className="text-left px-3 py-2 font-medium">OTA?</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-border/20">
            <tr><td className="px-3 py-2">UI / styles / new screens</td><td className="text-emerald-500">Yes</td></tr>
            <tr><td className="px-3 py-2">Bug fixes in JS / TS</td><td className="text-emerald-500">Yes</td></tr>
            <tr><td className="px-3 py-2">API contract changes (calling new endpoints)</td><td className="text-emerald-500">Yes</td></tr>
            <tr><td className="px-3 py-2">Adding a new npm dep that&apos;s pure-JS</td><td className="text-emerald-500">Yes</td></tr>
            <tr><td className="px-3 py-2">Adding a native module (camera, BLE, etc.)</td><td className="text-rose-500">No — store re-submit</td></tr>
            <tr><td className="px-3 py-2">Changing app icon or splash screen</td><td className="text-rose-500">No — store re-submit</td></tr>
            <tr><td className="px-3 py-2">Bumping Expo SDK</td><td className="text-rose-500">No — store re-submit</td></tr>
          </tbody>
        </table>
      </div>

      <h2>Shipping an OTA</h2>
      <CodeBlock terminal code={`# from apps/mobile/
eas update --branch production --message "Fix dashboard refresh bug"`} />
      <p>
        Bundles your JS, uploads to Expo. Apps configured for{' '}
        <code>production</code> branch check for updates on next launch and
        download the new bundle.
      </p>

      <h2>Branches map to release channels</h2>
      <p>
        Set up at least two branches:
      </p>
      <ul>
        <li>
          <strong>preview</strong> — internal testers / beta users
        </li>
        <li>
          <strong>production</strong> — store users
        </li>
      </ul>
      <p>
        Push to preview first, validate, then push to production. Same
        idea as canary / staging / prod for backends.
      </p>

      <h2>What the user experiences</h2>
      <ol>
        <li>User opens the app</li>
        <li>Expo checks for new bundle (~200ms HTTP call)</li>
        <li>If new bundle available, downloads it in the background</li>
        <li>Next launch uses the new bundle</li>
      </ol>
      <p>
        First-launch users see the bundled bundle (from the store). On
        their second launch, they see your update.
      </p>

      <h2>Forcing an immediate reload</h2>
      <p>
        For urgent fixes, you can ask Expo to fetch + reload on this
        session, not next launch:
      </p>
      <CodeBlock
        language="ts"
        code={`import * as Updates from 'expo-updates'

useEffect(() => {
  Updates.checkForUpdateAsync().then((update) => {
    if (update.isAvailable) {
      Updates.fetchUpdateAsync().then(() => Updates.reloadAsync())
    }
  })
}, [])`}
      />
      <p>
        Use sparingly — reloading mid-session is jarring. Better for
        important security fixes.
      </p>

      <TipBox tone="warning">
        <strong>Don&apos;t OTA something that needs a native change.</strong>{' '}
        If your JS expects a new native module that&apos;s not in the installed
        binary, the app crashes on first use. Match OTA bundle compatibility
        to the installed runtime version.
      </TipBox>

      <h2>Rolling back</h2>
      <p>
        Pushed a bad update? Republish the previous version&apos;s bundle:
      </p>
      <CodeBlock
        terminal
        code={`eas update --branch production --republish --group <old-group-id>`}
      />
      <p>
        Every <code>eas update</code> creates a &quot;group&quot; you can
        re-publish later. Roll back to any earlier group, users get it on
        next launch.
      </p>

      <KnowledgeCheck
        question="You shipped a UI tweak via `eas update`. A user emails saying their app crashed after opening. Diff: you removed `react-native-camera` from the JS imports, but the binary on their phone still has the native module. Why did it crash?"
        choices={[
          {
            label: 'Removing a native module from a working binary always crashes',
            feedback:
              "Misleading — removing imports doesn't crash. The crash is from something else in the diff.",
          },
          {
            label: 'You probably also OTA\'d a JS change that references a NEW native module the user\'s binary doesn\'t have',
            correct: true,
            feedback:
              "Most likely — JS calling a missing native module crashes. OTA can't add native modules; only the store re-submit can. Audit the diff for any new native deps.",
          },
          {
            label: 'OTA disabled the user\'s permissions',
            feedback:
              "OTA doesn't touch permissions. Permissions live in the installed binary's Info.plist / AndroidManifest.xml.",
          },
          {
            label: 'Expo OTAs always require a binary rebuild',
            feedback:
              "Wrong — OTAs explicitly avoid binary rebuilds for pure JS changes. That's the whole point.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Ship your first OTA:</p>
            <ol>
              <li>
                Make a small visible change to your app (e.g., change the
                home-screen title from &quot;Users&quot; to &quot;Team&quot;).
              </li>
              <li>
                <code>eas update --branch preview --message &quot;rename users to team&quot;</code>
              </li>
              <li>
                On your phone (where you installed the preview build), force-
                quit the app and reopen it.
              </li>
              <li>
                You should see the new title without re-installing.
              </li>
            </ol>
            <p>
              Paste a before/after screenshot in <code>notes.md</code>.
              That&apos;s chapter 5&apos;s assignment done.
            </p>
          </>
        }
        hint={
          <>
            The update is only downloaded on launch. If you don&apos;t see the
            change, the download may not have completed yet. Open and close
            once more — second launch typically picks it up.
          </>
        }
        solution={
          <>
            <p>
              You shipped a code change to a real device without a new
              install. That&apos;s the OTA pipeline.
            </p>
          </>
        }
      />

      <h2>You finished Building Mobile with Go API 🎉</h2>
      <p>
        Five chapters, ~13 lessons. You can now scaffold a mobile + API
        monorepo, sync types end-to-end, build login + secure storage +
        refresh, register and receive push notifications, ship to the
        stores via EAS, and OTA updates over the air.
      </p>
      <p>
        From here: pair this with{' '}
        <code>--triple --mobile</code> in the{' '}
        <a href="/courses/multiplatform" className="text-primary hover:underline">
          Multi-Platform course
        </a>{' '}
        to add web + admin to the mix, or move to{' '}
        <a href="/courses/web-nextjs" className="text-primary hover:underline">
          Building Web with Next.js
        </a>{' '}
        to add a marketing site + dashboard.
      </p>
    </>
  )
}
