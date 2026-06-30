import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Push Notifications — Grit Mobile Course',
  description: 'Learn Expo push tokens, the grit-notifications plugin, FCM, APNS, local notifications, and handling notification taps in React Native.',
}

export default function NotificationsCourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <Link href="/courses/grit-mobile" className="hover:text-foreground transition-colors">Grit Mobile</Link>
          <span>/</span>
          <span className="text-foreground">Push Notifications</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Course 4 of 5</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">10 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Push Notifications
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Push notifications keep users engaged by delivering timely messages even when the app
            is closed. In this course, you will learn how push notifications work end-to-end — from
            registering device tokens to sending notifications from your Go API and handling them
            in the app.
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ What are Push Notifications? ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What are Push Notifications?</h2>

          <Definition term="Push Notification">
            A message that appears on the user{"'"}s phone even when the app is closed. The message
            shows up in the notification tray (the pull-down area on both iOS and Android). Push
            notifications are sent from your server through Apple (APNS) or Google (FCM) to the
            device.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Push notifications are different from in-app messages. An in-app message only shows
            when the user has the app open. A push notification reaches the user regardless —
            their phone buzzes, the message appears on the lock screen, and tapping it opens your app.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The flow involves three parties:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">1.</span> <strong className="text-foreground">Your app</strong> — registers for notifications and gets a push token</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">2.</span> <strong className="text-foreground">Your Go API</strong> — stores the token and sends notifications when events occur</li>
            <li className="flex gap-2"><span className="text-primary font-mono font-bold">3.</span> <strong className="text-foreground">Apple/Google</strong> — delivers the notification to the specific device</li>
          </ul>

          <Challenge number={1} title="Understand the Flow">
            <p>Explain the push notification flow in your own words. Why can{"'"}t your server send a notification directly to a phone? Why do Apple and Google sit in the middle?</p>
          </Challenge>
        </section>

        {/* ═══ Expo Push Notifications ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Expo Push Notifications</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Expo simplifies push notifications significantly. Instead of configuring APNS and FCM
            separately, you use Expo{"'"}s unified push service. Your server sends a notification
            to Expo, and Expo handles routing it through APNS or FCM depending on the device.
          </p>

          <Definition term="Push Token">
            A unique identifier for a specific device and app combination. Your server uses this
            token to send a notification to that exact device. Each device has a different token,
            and tokens can change over time (e.g., when the user reinstalls the app).
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The complete flow looks like this:
          </p>

          <CodeBlock language="text">{`1. App launches → requests notification permissions
2. App gets Expo push token (e.g., "ExponentPushToken[xxxxxx]")
3. App sends token to your API: POST /api/devices { push_token: token }
4. API stores token in the database (linked to the user)
5. When an event occurs, API sends notification via Expo push service
6. Expo routes it through APNS (iOS) or FCM (Android)
7. Device receives the notification`}</CodeBlock>

          <Challenge number={2} title="Request Permissions">
            <p>What permissions does the app need to request for push notifications? On iOS, the user must explicitly grant permission. On Android 13+, notification permission is also required. What happens if the user denies permission?</p>
          </Challenge>
        </section>

        {/* ═══ Getting the Push Token ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Getting the Push Token</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Use Expo{"'"}s Notifications API to get the device{"'"}s push token:
          </p>

          <CodeBlock language="tsx">{`import * as Notifications from 'expo-notifications'
import * as Device from 'expo-device'
import { Platform } from 'react-native'

async function registerForPushNotifications() {
  // Push notifications only work on physical devices
  if (!Device.isDevice) {
    console.log('Push notifications require a physical device')
    return null
  }

  // Request permission
  const { status: existingStatus } = await Notifications.getPermissionsAsync()
  let finalStatus = existingStatus

  if (existingStatus !== 'granted') {
    const { status } = await Notifications.requestPermissionsAsync()
    finalStatus = status
  }

  if (finalStatus !== 'granted') {
    console.log('Permission not granted for push notifications')
    return null
  }

  // Get the Expo push token
  const token = (await Notifications.getExpoPushTokenAsync()).data
  console.log('Push token:', token)

  // Android needs a notification channel
  if (Platform.OS === 'android') {
    await Notifications.setNotificationChannelAsync('default', {
      name: 'default',
      importance: Notifications.AndroidImportance.MAX,
    })
  }

  return token
}`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            After getting the token, send it to your API so the server can store it:
          </p>

          <CodeBlock language="tsx">{`// Send token to your API
const token = await registerForPushNotifications()
if (token) {
  await api.post('/api/devices', { push_token: token })
}`}</CodeBlock>

          <Note>
            Push notifications do <strong>not</strong> work in the iOS Simulator or Android Emulator.
            You need a physical device to test them. The token registration code should handle
            this gracefully by checking <Code>Device.isDevice</Code> first.
          </Note>

          <Challenge number={3} title="Register a Push Token">
            <p>On a physical device, run the push token registration code. What does the token look like? It should start with {'"'}ExponentPushToken[{'"'}. Send it to your API and verify it{"'"}s stored in the database.</p>
          </Challenge>
        </section>

        {/* ═══ The grit-notifications Plugin ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">The grit-notifications Plugin</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit provides a notifications plugin that handles device token storage, sending
            notifications to one or many devices, and keeping a notification history.
          </p>

          <CodeBlock language="bash">{`go get github.com/MUKE-coder/grit-plugins/grit-notifications`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The plugin provides:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Device model</strong> — stores push tokens linked to user accounts</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Send to one</strong> — send a notification to a single device</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Send to many</strong> — broadcast to all devices of a user, or all users</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Notification history</strong> — log of all sent notifications for debugging</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Token cleanup</strong> — automatically removes invalid tokens</li>
          </ul>

          <Challenge number={4} title="Explore the Plugin">
            <p>Look at the grit-notifications plugin documentation. What functions does it expose? What database tables does it create?</p>
          </Challenge>
        </section>

        {/* ═══ Sending a Notification ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Sending a Notification</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            From your Go API, send a notification using the plugin:
          </p>

          <CodeBlock language="go">{`import notify "github.com/MUKE-coder/grit-plugins/grit-notifications"

// Send a notification to a specific device
err := notify.Send(notify.PushMessage{
    To:    pushToken,
    Title: "New Workout Added",
    Body:  "Your trainer added a new workout for today.",
    Data:  map[string]string{"screen": "workout", "id": "42"},
})

// Send to all devices of a specific user
err = notify.SendToUser(db, userID, notify.PushMessage{
    Title: "Weekly Summary",
    Body:  "You completed 5 workouts this week!",
    Data:  map[string]string{"screen": "stats"},
})`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>Data</Code> field is a key-value map that is delivered to the app but
            <strong className="text-foreground"> not shown to the user</strong>. Use it to pass information
            the app needs to handle the notification — like which screen to navigate to and what
            ID to load.
          </p>

          <Tip>
            Always include a <Code>screen</Code> key in the <Code>Data</Code> field. This tells the
            app which screen to navigate to when the user taps the notification.
          </Tip>

          <Challenge number={5} title="Design a Notification">
            <p>If you were building an order confirmation notification for an e-commerce app, what would the <Code>Title</Code>, <Code>Body</Code>, and <Code>Data</Code> fields contain? Write out the complete <Code>PushMessage</Code> struct.</p>
          </Challenge>
        </section>

        {/* ═══ Handling Notifications ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Handling Notifications</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Notifications can arrive in three different states, and you need to handle each one:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Foreground</strong> — app is open and visible. Show an in-app banner or toast.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Background</strong> — app is in the background. The OS shows the notification in the tray.</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Killed</strong> — app is completely closed. The OS shows the notification. Tapping it opens the app.</li>
          </ul>

          <CodeBlock language="tsx">{`import * as Notifications from 'expo-notifications'
import { useRouter } from 'expo-router'

export function useNotificationListeners() {
  const router = useRouter()

  useEffect(() => {
    // Foreground: notification received while app is open
    const receivedSub = Notifications.addNotificationReceivedListener(
      (notification) => {
        const data = notification.request.content.data
        console.log('Received in foreground:', data)
        // Show an in-app toast or banner
      }
    )

    // User tapped: notification was tapped (from any state)
    const responseSub = Notifications.addNotificationResponseReceivedListener(
      (response) => {
        const data = response.notification.request.content.data
        // Navigate to the relevant screen
        if (data.screen === 'workout' && data.id) {
          router.push('/workouts/' + data.id)
        }
      }
    )

    return () => {
      receivedSub.remove()
      responseSub.remove()
    }
  }, [])
}`}</CodeBlock>

          <Definition term="Notification Channel (Android)">
            On Android 8+, notifications must belong to a channel. Channels let users control
            notification settings per category — they can mute {'"'}Marketing{'"'} notifications while
            keeping {'"'}Order Updates{'"'} enabled. You create channels in your app code.
          </Definition>

          <Challenge number={6} title="Handle Notification Tap">
            <p>What should happen when a user taps a notification about a new workout? The app should open, navigate to the workout detail screen, and show the specific workout. Write the handler code.</p>
          </Challenge>

          <Challenge number={7} title="Foreground Behavior">
            <p>When a notification arrives while the app is open (foreground), the default behavior on iOS is to <strong>not</strong> show it. How would you configure Expo to show a banner even in the foreground? (Hint: look at <Code>setNotificationHandler</Code>)</p>
          </Challenge>
        </section>

        {/* ═══ Local Notifications ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Local Notifications</h2>

          <Definition term="Local Notification">
            A notification triggered by the app itself — not from a server. Local notifications
            are useful for reminders, timers, alarms, and scheduled prompts. They work offline
            and do not require a push token or server infrastructure.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Use Expo{"'"}s <Code>scheduleNotificationAsync</Code> to schedule a local notification:
          </p>

          <CodeBlock language="tsx">{`import * as Notifications from 'expo-notifications'

// Schedule a notification for 1 hour from now
await Notifications.scheduleNotificationAsync({
  content: {
    title: "Reminder",
    body: "Time to work out!",
    data: { screen: "workouts" },
  },
  trigger: {
    seconds: 3600, // 1 hour
  },
})

// Schedule a daily notification at 9:00 AM
await Notifications.scheduleNotificationAsync({
  content: {
    title: "Good morning!",
    body: "Ready for today's workout?",
  },
  trigger: {
    hour: 9,
    minute: 0,
    repeats: true,
  },
})`}</CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Local notifications are perfect for features like:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> Workout reminders ({'"'}Time to exercise!{'"'})</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Timer completions ({'"'}Rest period over — next set!{'"'})</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Daily streaks ({'"'}Don{"'"}t break your 7-day streak!{'"'})</li>
            <li className="flex gap-2"><span className="text-primary">•</span> Meal reminders ({'"'}Time to log your lunch{'"'})</li>
          </ul>

          <Challenge number={8} title="Schedule a Local Notification">
            <p>Schedule a local notification that fires 10 seconds from now with a title and body of your choice. Does it appear on your device? What happens if you close the app before it fires?</p>
          </Challenge>
        </section>

        {/* ═══ Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s what you learned in this course:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-6">
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Push notifications are delivered through Apple (APNS) or Google (FCM)</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Expo simplifies push by providing a unified API for both platforms</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Push tokens uniquely identify a device and must be stored on your server</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Handle three states: foreground, background, and killed</li>
            <li className="flex gap-2"><CheckCircle2 className="h-5 w-5 text-green-400 shrink-0 mt-0.5" /> Local notifications work offline and are great for reminders and timers</li>
          </ul>

          <Challenge number={9} title="Design a Notification System">
            <p>When should your fitness app send push notifications? List 5 scenarios. For each, specify whether it should be a push notification (from the server) or a local notification (from the app).</p>
          </Challenge>

          <Challenge number={10} title="Write the Payloads">
            <p>For each of your 5 scenarios from the previous challenge, write the complete notification payload: <Code>Title</Code>, <Code>Body</Code>, and <Code>Data</Code> (with screen and ID). Make the copy concise and actionable.</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses/grit-mobile/api-offline', label: 'Prev: API Integration & Offline' }}
            next={{ href: '/courses/grit-mobile/build-store', label: 'Next: Build & App Store' }}
          />
        </div>
      </main>
    </div>
  )
}
