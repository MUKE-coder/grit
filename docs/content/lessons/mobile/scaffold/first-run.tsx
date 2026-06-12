import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Time to actually run the app. Three modes — Expo Go (fastest),
        simulator (Mac iOS or Android emulator), dev build (most realistic).
        This lesson covers when to use which.
      </p>

      <h2>1. The first run on Expo Go</h2>
      <CodeBlock terminal code={`cd apps/mobile
pnpm dev`} />
      <p>You&apos;ll see a QR code in the terminal:</p>
      <CodeBlock
        language="text"
        code={`Metro Bundler ready

› Scan the QR code above with Expo Go (Android) or the Camera app (iOS)
› Press a │ open Android
› Press i │ open iOS simulator
› Press w │ open web`}
      />
      <ol>
        <li>Install the <strong>Expo Go</strong> app on your phone</li>
        <li>Make sure your phone is on the SAME Wi-Fi as your computer</li>
        <li>Scan the QR — your app loads on the phone</li>
      </ol>

      <h2>2. Simulator / emulator</h2>
      <p>
        Press <code>i</code> in the Metro terminal (Mac only — opens iOS
        simulator) or <code>a</code> (Android emulator, needs Android Studio
        installed and an AVD). The app loads in a window on your computer.
      </p>

      <h2>3. Dev build</h2>
      <p>
        Expo Go runs your JS but uses Expo&apos;s pre-bundled native modules.
        If you add a library that needs a custom native module (e.g.,
        Bluetooth, BLE, camera with custom config), Expo Go can&apos;t load it
        — you need a <strong>dev build</strong>.
      </p>
      <CodeBlock terminal code={`eas build --profile development --platform ios`} />
      <p>
        Same QR-code workflow, but you scan with your dev build app instead
        of Expo Go. Once you&apos;ve made a dev build, JS updates still flow
        live; only native changes require a rebuild.
      </p>

      <TipBox tone="info">
        <strong>When to graduate from Expo Go to a dev build:</strong> the
        moment you add a native module. Most apps spend their first few
        weeks in Expo Go and graduate when they need camera customisation,
        local notifications with custom sounds, BLE, or in-app purchases.
      </TipBox>

      <h2>Same network, different APIs</h2>
      <p>
        When you scan with Expo Go, your phone runs the JS but tries to talk
        to your API. <code>localhost</code> won&apos;t work — the phone&apos;s
        localhost is itself, not your computer.
      </p>
      <CodeBlock
        language="ts"
        filename="apps/mobile/lib/api.ts"
        code={`// WRONG — won't work on a physical device
export const API_URL = 'http://localhost:8080'

// RIGHT — your computer's LAN IP, reachable from the phone
export const API_URL = 'http://192.168.1.42:8080'`}
      />
      <p>
        Find your LAN IP: <code>ifconfig</code> on Mac/Linux,{' '}
        <code>ipconfig</code> on Windows. Look for the 192.168.x.x or
        10.0.x.x address. The simulator can use{' '}
        <code>http://localhost:8080</code> because it shares the host&apos;s
        network.
      </p>

      <h2>The dev menu</h2>
      <p>
        Shake your phone (or press <code>m</code> in Metro) to open the dev
        menu — reload bundle, toggle perf monitor, enable network inspector.
        You&apos;ll use this constantly.
      </p>

      <KnowledgeCheck
        question="You scan the QR code on your phone. The app starts but every API call fails with 'network request failed'. What's the most likely cause?"
        choices={[
          {
            label: 'The API isn\'t running',
            feedback:
              "Possible — check `docker compose ps`. But if you're seeing 'network request failed' specifically (not 401 or 500), the connection isn't even being made.",
          },
          {
            label: 'API_URL is set to localhost — the phone is trying to connect to itself',
            correct: true,
            feedback:
              "Right — the phone's localhost is its own loopback, not your computer's. Use your computer's LAN IP (192.168.x.x) instead.",
          },
          {
            label: 'Expo Go doesn\'t support HTTP — needs HTTPS',
            feedback:
              "Wrong — Expo Go is HTTP-friendly in dev mode. Production iOS builds require HTTPS by default (configurable), but Expo Go is fine.",
          },
          {
            label: 'CORS isn\'t configured',
            feedback:
              "CORS errors are different — they'd show as a CORS-specific message. 'network request failed' is a connectivity issue.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Run your scaffolded app on BOTH a simulator/emulator AND a physical device:</p>
            <ol>
              <li><code>cd apps/mobile && pnpm dev</code></li>
              <li>Press <code>i</code> (Mac) or <code>a</code> (Windows/Linux with AVD) to open simulator</li>
              <li>Install Expo Go on your phone (App Store / Play Store)</li>
              <li>Scan the QR code with your phone</li>
              <li>
                Once both are running: take a screenshot of each. Paste both
                in <code>notes.md</code>.
              </li>
            </ol>
            <p>That completes chapter 1&apos;s assignment.</p>
          </>
        }
        hint={
          <>
            If Expo Go on your phone shows a network error connecting to the
            Metro server, your phone and computer aren&apos;t on the same Wi-Fi
            (or your network blocks AP isolation). Try tethering your
            computer&apos;s Wi-Fi to your phone&apos;s hotspot for a quick fix.
          </>
        }
        solution={
          <>
            <p>
              Two screenshots — the same Expo home screen running on two
              surfaces. The simulator gives you fast iteration; the physical
              device gives you reality. You&apos;ll use both regularly.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Chapter 2 — <strong>Shared types + API client</strong>. We&apos;ll wire{' '}
        <code>grit sync</code> so your TypeScript types match your Go
        structs end to end.
      </p>
    </>
  )
}
