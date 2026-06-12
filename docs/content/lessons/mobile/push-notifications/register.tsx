import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Push notifications go: phone gets a token from Apple / Google →
        mobile sends the token to your API → API saves it on the user
        → later, a job worker sends pushes to that token via Expo Push.
        This lesson covers the first two steps.
      </p>

      <h2>The plumbing</h2>
      <CodeBlock terminal code={`pnpm add expo-notifications expo-device`} />
      <p>Already in the Grit mobile scaffold.</p>

      <h2>Request permission + grab the token</h2>
      <CodeBlock
        language="ts"
        filename="apps/mobile/lib/push.ts"
        code={`import * as Notifications from 'expo-notifications'
import * as Device from 'expo-device'
import Constants from 'expo-constants'
import { Platform } from 'react-native'

export async function registerForPush(): Promise<string | null> {
  if (!Device.isDevice) return null  // simulator can't receive real pushes

  // iOS-specific: must request permission
  const { status: existing } = await Notifications.getPermissionsAsync()
  let final = existing
  if (existing !== 'granted') {
    const { status } = await Notifications.requestPermissionsAsync()
    final = status
  }
  if (final !== 'granted') return null  // user denied

  // Android-specific: must create a channel
  if (Platform.OS === 'android') {
    await Notifications.setNotificationChannelAsync('default', {
      name: 'default',
      importance: Notifications.AndroidImportance.DEFAULT,
    })
  }

  const projectId = Constants.expoConfig?.extra?.eas?.projectId
  const token = (await Notifications.getExpoPushTokenAsync({ projectId })).data
  return token  // ExponentPushToken[xxxxxxxxxxxxxxxxxxxx]
}`}
      />
      <p>
        Returns an Expo Push token — a string starting with{' '}
        <code>ExponentPushToken[</code>. Expo&apos;s push service uses it to
        identify the device when you want to send a push.
      </p>

      <h2>Save it to your API</h2>
      <CodeBlock
        language="ts"
        filename="apps/mobile/hooks/use-auth.ts (after login)"
        code={`async function login(email: string, password: string) {
  // ... existing login logic ...

  // After successful login, register for push and save the token
  const pushToken = await registerForPush()
  if (pushToken) {
    await api.post('/api/users/me/push-token', { push_token: pushToken })
  }
}`}
      />
      <p>
        On the API side, add a handler that updates the user&apos;s{' '}
        <code>push_token</code> column.
      </p>

      <h2>Foreground vs background — what the user sees</h2>
      <ul>
        <li>
          <strong>App in foreground</strong> — Expo doesn&apos;t show a system
          notification by default. Use{' '}
          <code>Notifications.setNotificationHandler</code> to decide
          whether to show a banner.
        </li>
        <li>
          <strong>App backgrounded / killed</strong> — system shows the
          notification automatically. Tapping it deep-links into the app.
        </li>
      </ul>
      <CodeBlock
        language="ts"
        code={`Notifications.setNotificationHandler({
  handleNotification: async () => ({
    shouldShowAlert: true,
    shouldPlaySound: true,
    shouldSetBadge: true,
  }),
})`}
      />

      <TipBox tone="info">
        <strong>Simulator gotcha:</strong>{' '}
        <code>Device.isDevice</code> is <code>false</code> in the iOS
        simulator and Android emulator, so push registration is skipped.
        You must test push on a physical device.
      </TipBox>

      <h2>Handling taps + deep-links</h2>
      <CodeBlock
        language="ts"
        code={`Notifications.addNotificationResponseReceivedListener((response) => {
  const data = response.notification.request.content.data
  if (data?.orderId) {
    router.push(\`/orders/\${data.orderId}\`)
  }
})`}
      />
      <p>
        The <code>data</code> payload you sent from the API arrives here.
        Use it to navigate deep into the app — &quot;new message&quot; opens the
        thread, &quot;order shipped&quot; opens the order.
      </p>

      <h2>Updating the token over time</h2>
      <p>
        Expo push tokens can change — after app reinstall, after iOS
        re-permission cycle. Re-register on every cold start so your API
        always has the current one:
      </p>
      <CodeBlock
        language="ts"
        code={`useEffect(() => {
  if (isAuthenticated) {
    registerForPush().then((token) => {
      if (token) api.post('/api/users/me/push-token', { push_token: token })
    })
  }
}, [isAuthenticated])`}
      />

      <KnowledgeCheck
        question="You're testing push notifications on your iOS simulator. registerForPush() returns null every time. What's wrong?"
        choices={[
          {
            label: 'Permission was denied',
            feedback:
              "Possible but unlikely — the simulator usually shows the prompt. Even if denied, you'd see the prompt first.",
          },
          {
            label: 'iOS simulator can\'t receive real pushes — Device.isDevice is false, so registerForPush short-circuits',
            correct: true,
            feedback:
              "Right — Apple's Push Notification Service (APNs) doesn't deliver to simulators. expo-device's `isDevice` is false in sims; the function returns null intentionally.",
          },
          {
            label: 'expo-notifications is out of date',
            feedback:
              "Unlikely to cause this specific symptom. Update if you want, but the root cause is the simulator limitation.",
          },
          {
            label: 'Push permissions need to be set in app.json',
            feedback:
              "You do need iOS permission strings in app.json for production, but the simulator wouldn't get past Device.isDevice anyway.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Register your device for push:</p>
            <ol>
              <li>
                On a <strong>physical phone</strong> (not the simulator),
                run the app via Expo Go.
              </li>
              <li>Sign in.</li>
              <li>
                Call <code>registerForPush()</code> on mount (or log in).
                Console-log the returned token.
              </li>
              <li>
                The token starts with <code>ExponentPushToken[</code> — copy
                it.
              </li>
              <li>Paste the token in <code>notes.md</code>.</li>
            </ol>
          </>
        }
        hint={
          <>
            If you don&apos;t see the token, check the device console (Metro
            terminal shows console.log from connected devices). On Android,
            you may also need to manually grant the Notifications permission
            in Settings if the prompt didn&apos;t fire.
          </>
        }
        solution={
          <>
            <p>The token looks like:</p>
            <CodeBlock
              language="text"
              code={`ExponentPushToken[xxXxXxxxXXXxxxx-xXxx]`}
            />
            <p>
              Save that token (it&apos;s yours). In the next lesson we send a
              push to it.
            </p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Last lesson of the chapter — actually <strong>send</strong> a push
        from your Grit API to your phone.
      </p>
    </>
  )
}
