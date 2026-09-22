import Link from 'next/link'
import { ArrowLeft, ArrowRight, Bell } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { SiteHeader } from '@/components/site-header'
import { DocsSidebar } from '@/components/docs-sidebar'
import { CodeBlock } from '@/components/code-block'

export const metadata = {
  title: 'Push notifications plugin | Grit',
  description:
    'Push notifications to the Expo app on iOS and Android, through Expo’s push service: device tokens, sending, and cleaning up tokens for deleted apps.',
  alternates: { canonical: 'https://gritframework.dev/docs/plugins/push' },
}

export default function PushPluginPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="mx-auto max-w-3xl px-6 py-12">
          <div className="mb-3 flex items-center gap-2">
            <Bell className="h-5 w-5 text-primary" />
            <span className="font-mono text-xs uppercase tracking-wider text-muted-foreground">Plugin</span>
          </div>

          <h1 className="mb-4 font-display text-4xl font-bold tracking-tight">Push notifications</h1>
          <p className="mb-8 text-lg leading-relaxed text-muted-foreground">
            Send a notification to a user&apos;s phone from Go, in one line. The API talks to Expo&apos;s push
            service, which relays to Apple and Google, so there is no APNs certificate or Firebase key to manage.
          </p>

          <CodeBlock language="bash" code={`grit plugin add push
grit migrate
pnpm install`} />

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">Send one</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            <code>Go</code> sends in the background, so the request that caused it never waits on the push service.
            Every device the users registered gets it.
          </p>
          <CodeBlock language="go" code={`push := services.NewPush(db)

push.Go([]string{order.CustomerID}, services.PushMessage{
    Title: "Your order shipped",
    Body:  "Order " + order.Number + " is on its way.",
    Sound: "default",
    // Arrives with the notification: what the app opens when it is tapped.
    Data: map[string]string{"order_id": order.ID},
})`} />
          <p className="mb-4 mt-4 leading-relaxed text-muted-foreground">
            <code>Send</code> does the same and waits, returning how many devices Expo accepted: for a background
            job that wants the result.
          </p>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">In the Expo app</h2>
          <p className="mb-4 leading-relaxed text-muted-foreground">
            <code>apps/expo/lib/push.ts</code> asks for permission, gets the device&apos;s token and registers it.
            Call it after sign-in and on every launch while signed in; unregister before signing out.
          </p>
          <CodeBlock language="tsx" code={`import { onNotificationTap, registerForPush, unregisterPush } from "@/lib/push";

useEffect(() => {
  if (user) void registerForPush();
}, [user?.id]);

// The tap that opened the app counts too.
useEffect(() => onNotificationTap((data) => {
  if (data.order_id) router.push(\`/orders/\${data.order_id}\`);
}), []);

// In logout, before the session ends:
await unregisterPush();`} />

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">What you get</h2>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-border text-left">
                  <th className="py-2 pr-4 font-medium text-foreground">Endpoint</th>
                  <th className="py-2 font-medium text-foreground">Does</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/50">
                  <td className="py-2 pr-4 font-mono text-xs">POST /api/v1/push/tokens</td>
                  <td className="py-2">Registers the signed-in device. Anything but an Expo push token is refused.</td>
                </tr>
                <tr className="border-b border-border/50">
                  <td className="py-2 pr-4 font-mono text-xs">POST /api/v1/push/tokens/remove</td>
                  <td className="py-2">Unregisters it, if it is yours. The answer is the same either way.</td>
                </tr>
                <tr>
                  <td className="py-2 pr-4 font-mono text-xs">POST /api/v1/push/test</td>
                  <td className="py-2">Sends yourself a test notification, to check the set-up end to end.</td>
                </tr>
              </tbody>
            </table>
          </div>

          <h2 className="mb-4 mt-12 text-2xl font-semibold tracking-tight">The details that matter</h2>
          <ul className="mb-8 list-disc space-y-3 pl-6 leading-relaxed text-muted-foreground">
            <li>
              <strong className="text-foreground">Dead tokens clean themselves up.</strong> When an app is deleted or
              notifications are turned off, Expo answers <code>DeviceNotRegistered</code> for that token, and it is
              removed on the spot instead of being tried with every message from then on.
            </li>
            <li>
              <strong className="text-foreground">A phone belongs to whoever signed in last.</strong> A token is
              unique, and registering it moves it to the current user, so a shared or handed-down phone stops getting
              the previous person&apos;s notifications.
            </li>
            <li>
              <strong className="text-foreground">Batches of 100.</strong> Expo takes at most 100 messages a request;
              sending to more users makes as many requests as it needs.
            </li>
            <li>
              <strong className="text-foreground">Ten devices a user.</strong> The oldest are dropped first, since a
              phone that was replaced stops re-registering.
            </li>
            <li>
              <strong className="text-foreground">Real devices only.</strong> Simulators and the web cannot receive
              push; <code>registerForPush</code> returns <code>null</code> there and the app carries on. Store builds
              need an EAS project id; set <code>EXPO_ACCESS_TOKEN</code> if the project uses enhanced push security.
            </li>
          </ul>

          <p className="mb-8 leading-relaxed text-muted-foreground">
            Checked against Expo&apos;s live service while building the WhatsApp blueprint: a message to someone with
            a registered token went out through Expo, and a token no device owned came back as not registered and
            was removed.
          </p>

          <div className="mt-16 flex items-center justify-between border-t border-border/50 pt-8">
            <Link href="/docs/plugins/device-pairing">
              <Button variant="ghost" className="gap-2">
                <ArrowLeft className="h-4 w-4" />
                Device pairing
              </Button>
            </Link>
            <Link href="/docs/plugins/video">
              <Button variant="ghost" className="gap-2">
                Video
                <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
          </div>
        </div>
      </main>
    </div>
  )
}
