import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        The React side: a banner that appears when an update is
        available, a modal that shows progress + the &quot;Install&quot; button.
        The scaffold ships these — this lesson explains the components so
        you can customise them.
      </p>

      <h2>The hook that polls</h2>
      <CodeBlock
        language="ts"
        filename="src/hooks/use-update-checker.ts"
        code={`import { useQuery } from '@tanstack/react-query'
import { CheckForUpdates } from '../../wailsjs/go/main/Updater'

export function useUpdateChecker() {
  const q = useQuery({
    queryKey: ['updater', 'latest'],
    queryFn: () => CheckForUpdates(),
    refetchInterval: 6 * 60 * 60 * 1000,  // 6h
    refetchOnWindowFocus: false,
  })

  const isUpdateAvailable = !!q.data?.is_newer
  return { update: q.data, isUpdateAvailable, checkNow: () => q.refetch() }
}`}
      />
      <p>
        Cheap call — one HTTP round-trip to GitHub. Defaulting to 6 hours
        means each running app costs ~4 GitHub requests/day.
      </p>

      <h2>The banner</h2>
      <CodeBlock
        language="tsx"
        filename="src/components/update-banner.tsx"
        code={`export function UpdateBanner() {
  const { update, isUpdateAvailable } = useUpdateChecker()
  const [dismissed, setDismissed] = useState(false)
  const [modalOpen, setModalOpen] = useState(false)

  if (!isUpdateAvailable || dismissed) return null

  return (
    <>
      <div className="flex items-center gap-3 px-4 py-2 border-b bg-primary/5">
        <Download className="h-3.5 w-3.5 text-primary" />
        <span className="flex-1 text-sm">
          Update available: <strong>v{update.version}</strong>
        </span>
        <button onClick={() => setModalOpen(true)} className="text-xs text-primary">
          Update now
        </button>
        <button onClick={() => setDismissed(true)} className="text-muted-foreground">
          <X className="h-3.5 w-3.5" />
        </button>
      </div>
      {modalOpen && <UpdateModal update={update} onClose={() => setModalOpen(false)} />}
    </>
  )
}`}
      />
      <p>
        Mount once in the root layout. Dismissal is per-session — the
        banner reappears on next launch if the update is still pending.
      </p>

      <h2>The modal — download progress + install</h2>
      <CodeBlock
        language="tsx"
        filename="src/components/update-modal.tsx (key parts)"
        code={`export function UpdateModal({ update, onClose }) {
  const [progress, setProgress] = useState({ status: 'idle', bytesDownloaded: 0, bytesTotal: 0 })

  // Auto-start the download on mount
  useEffect(() => {
    StartDownload(update.download_url, update.version)
  }, [])

  // Poll progress every 500ms
  useEffect(() => {
    const id = setInterval(async () => {
      const p = await GetProgress()
      setProgress(p)
    }, 500)
    return () => clearInterval(id)
  }, [])

  const pct = (progress.bytesDownloaded / progress.bytesTotal) * 100

  return (
    <Dialog open onOpenChange={onClose}>
      <DialogContent>
        <h2>Update available</h2>
        <p className="text-sm text-muted-foreground">v{APP_VERSION} → v{update.version}</p>

        {/* Release notes — capped so the install button stays visible */}
        <div className="max-h-[180px] overflow-y-auto text-sm">
          {update.release_notes}
        </div>

        {progress.status === 'downloading' && (
          <>
            <div className="flex justify-between text-xs">
              <span>Downloading…</span>
              <span>{Math.round(pct)}%</span>
            </div>
            <div className="h-2 bg-muted rounded">
              <div className="h-full bg-primary rounded" style={{ width: pct + '%' }} />
            </div>
          </>
        )}

        {progress.status === 'downloaded' && (
          <Button onClick={InstallAndRestart}>Install & restart</Button>
        )}
      </DialogContent>
    </Dialog>
  )
}`}
      />

      <TipBox tone="success">
        <strong>Cap the release-notes height.</strong> Long changelogs
        push the install button below the fold. The scaffold caps at
        180px max-height + scroll. Users always see the action.
      </TipBox>

      <h2>Error UX</h2>
      <p>
        Download fails (network blip, GitHub down). The user shouldn&apos;t
        be stuck. The scaffold handles four states:
      </p>
      <ul>
        <li><code>idle</code> — nothing happening</li>
        <li><code>downloading</code> — progress bar</li>
        <li><code>downloaded</code> — &quot;Install &amp; restart&quot; button</li>
        <li><code>error</code> — error message + &quot;Try again&quot; button</li>
      </ul>

      <h2>Where to mount the banner</h2>
      <p>
        Top of <code>App.tsx</code>, above the routes. Always visible
        when an update is pending — but stays out of the way (one row).
      </p>
      <CodeBlock
        language="tsx"
        code={`<div className="h-screen flex flex-col">
  <TitleBar />
  <UpdateBanner />   {/* ← here */}
  <main className="flex-1 overflow-auto">
    {/* routes */}
  </main>
</div>`}
      />

      <KnowledgeCheck
        question="A user gets the update banner. They click Update now. The download fails halfway with a network error. What should the modal show?"
        choices={[
          {
            label: 'Auto-close — the network will eventually recover',
            feedback:
              "Wrong — the user wanted the update. Failing silently leaves them confused. Show an error.",
          },
          {
            label: 'Show the error message + a "Try again" button',
            correct: true,
            feedback:
              "Right — the scaffolded modal does exactly this. The user retains agency: retry, or close the modal and try later. Per-session dismissal of the banner respects their choice.",
          },
          {
            label: 'Auto-retry every 5 seconds',
            feedback:
              "Wastes their bandwidth and burns CPU. Manual retry is the right control surface.",
          },
          {
            label: 'Roll back the partial download silently',
            feedback:
              "The scaffold does delete the .partial file, but ALSO needs to show the user what happened. Silence is bad UX.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>
              Inspect the three scaffolded components:
            </p>
            <ol>
              <li><code>src/components/update-banner.tsx</code></li>
              <li><code>src/components/update-modal.tsx</code></li>
              <li><code>src/hooks/use-update-checker.ts</code></li>
            </ol>
            <p>
              In <code>notes.md</code>, write down:
            </p>
            <ul>
              <li>The polling interval the hook uses</li>
              <li>The four progress states the modal handles</li>
              <li>Where the dismissal state is stored (component / global)</li>
            </ul>
          </>
        }
        hint={
          <>
            Dismissal is local state in the banner component — re-mounts
            on next launch, so the banner reappears.
          </>
        }
        solution={
          <>
            <p>You should find:</p>
            <ul>
              <li>Polling: 6 hours (refetchInterval)</li>
              <li>States: idle / downloading / downloaded / error (installing is a brief intermediate)</li>
              <li>Dismissal: <code>useState(false)</code> in UpdateBanner — per-session only</li>
            </ul>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last lesson of the chapter — the{' '}
        <strong>release script</strong> that builds, tags, and publishes
        a new version to GitHub so the updater can find it.
      </p>
    </>
  )
}
