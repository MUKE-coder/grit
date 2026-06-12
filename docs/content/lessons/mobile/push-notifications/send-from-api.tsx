import { CodeBlock } from '@/components/code-block'
import { TipBox } from '@/components/course/tip-box'
import { Exercise } from '@/components/course/exercise'
import { KnowledgeCheck } from '@/components/course/knowledge-check'

export default function Lesson() {
  return (
    <>
      <p>
        Phone is registered, token is in the DB. Now the API sends a push.
        We use Expo&apos;s push service — one HTTP endpoint, no APNs / FCM
        wiring on our side, no certificates. Job worker dispatches; user
        sees a banner.
      </p>

      <h2>The Expo Push endpoint</h2>
      <CodeBlock
        language="text"
        code={`POST https://exp.host/--/api/v2/push/send

{
  "to": "ExponentPushToken[xxx]",
  "title": "Order shipped",
  "body": "Your headphones are on the way",
  "data": { "orderId": "abc123" },
  "sound": "default"
}`}
      />
      <p>
        Pass the user&apos;s saved push token in <code>to</code>, the visible
        text in <code>title</code> + <code>body</code>, and any deep-link
        payload in <code>data</code>.
      </p>

      <h2>From a Grit job worker</h2>
      <CodeBlock
        language="go"
        filename="apps/api/internal/jobs/push.go"
        code={`const TaskTypeSendPush = "push:send"

type SendPushPayload struct {
    UserID string                 \`json:"user_id"\`
    Title  string                 \`json:"title"\`
    Body   string                 \`json:"body"\`
    Data   map[string]any         \`json:"data,omitempty"\`
}

func (j *Jobs) EnqueueSendPush(ctx context.Context, userID, title, body string, data map[string]any) error {
    payload, _ := json.Marshal(SendPushPayload{UserID: userID, Title: title, Body: body, Data: data})
    _, err := j.client.EnqueueContext(ctx, asynq.NewTask(TaskTypeSendPush, payload),
        asynq.MaxRetry(5),
    )
    return err
}

func (j *Jobs) HandleSendPush(ctx context.Context, task *asynq.Task) error {
    var p SendPushPayload
    json.Unmarshal(task.Payload(), &p)

    user, err := j.users.GetByID(ctx, p.UserID)
    if err != nil { return err }
    if user.PushToken == "" {
        return nil  // not registered — nothing to do
    }

    body, _ := json.Marshal(map[string]any{
        "to":    user.PushToken,
        "title": p.Title,
        "body":  p.Body,
        "data":  p.Data,
        "sound": "default",
    })
    resp, err := http.Post("https://exp.host/--/api/v2/push/send",
        "application/json", bytes.NewReader(body))
    if err != nil { return err }
    defer resp.Body.Close()

    if resp.StatusCode >= 300 {
        return fmt.Errorf("expo push: HTTP %d", resp.StatusCode)
    }
    return nil
}`}
      />

      <h2>Triggering it from anywhere</h2>
      <p>
        Any handler can fire a push by enqueuing the job. Example —
        notifying a user when an order ships:
      </p>
      <CodeBlock
        language="go"
        filename="apps/api/internal/handlers/order.go (excerpt)"
        code={`func (h *OrderHandler) MarkShipped(c *gin.Context) {
    order, _ := h.svc.MarkShipped(ctx, id)

    h.jobs.EnqueueSendPush(ctx, order.UserID,
        "Order shipped",
        "Your order is on the way",
        map[string]any{"orderId": order.ID.String()},
    )
    respond.OK(c, order, "Order shipped")
}`}
      />
      <p>
        The handler returns immediately. The worker picks up the job and
        sends the push asynchronously. If Expo&apos;s push API is briefly
        down, Asynq retries automatically.
      </p>

      <h2>Error codes from Expo</h2>
      <p>
        Expo&apos;s response includes per-token results:
      </p>
      <CodeBlock
        language="json"
        code={`{
  "data": [
    { "status": "ok", "id": "receipt-id-1" },
    { "status": "error", "message": "...",
      "details": { "error": "DeviceNotRegistered" } }
  ]
}`}
      />
      <p>
        Worth handling: <code>DeviceNotRegistered</code> means the token is
        dead (app uninstalled, push revoked). Mark the user&apos;s push_token
        as empty so you stop trying.
      </p>

      <TipBox tone="success">
        <strong>Batch by default.</strong> Expo accepts up to 100 messages
        in a single POST. If you&apos;re fanning out to thousands, batch them.
        For the &quot;welcome push&quot; / &quot;order shipped&quot; case, single sends
        are fine.
      </TipBox>

      <h2>Foreground display</h2>
      <p>
        Remember from the last lesson — by default Expo doesn&apos;t show
        notifications when your app is in the foreground. Set the
        notification handler to <code>shouldShowAlert: true</code> if you
        want banners regardless.
      </p>

      <KnowledgeCheck
        question="A push job fails because Expo's API returns 'DeviceNotRegistered'. What should the worker do?"
        choices={[
          {
            label: 'Retry indefinitely — the device might come back',
            feedback:
              "Wrong — DeviceNotRegistered means the token is permanently dead (app uninstalled, push revoked). Retrying is wasted work and floods the dead-letter queue.",
          },
          {
            label: 'Clear the push_token field on the user, return nil to mark the job done',
            correct: true,
            feedback:
              "Right — the token is dead. Update the user row so future pushes skip this user. Return nil so Asynq doesn't retry.",
          },
          {
            label: 'Email the user instead',
            feedback:
              "Sensible fallback for important messages, but doesn't address the immediate problem: the dead token is still in the DB.",
          },
          {
            label: 'Log out the user',
            feedback:
              "Way too aggressive. The user is fine; only their push subscription is dead.",
          },
        ]}
      />

      <Exercise
        prompt={
          <>
            <p>Send yourself a push from your API:</p>
            <ol>
              <li>
                Add the <code>HandleSendPush</code> job + register it with
                your worker mux.
              </li>
              <li>
                Add a debug endpoint <code>POST /api/notify/me</code> that
                enqueues a push to the authenticated user.
              </li>
              <li>
                Make sure your push_token is saved on your user row (from
                last lesson).
              </li>
              <li>Curl the endpoint with your auth token.</li>
              <li>Push should arrive on your phone within ~5 seconds.</li>
            </ol>
            <p>
              That&apos;s chapter 4&apos;s assignment done. Paste a screenshot of
              the notification on your lock screen in{' '}
              <code>notes.md</code>.
            </p>
          </>
        }
        hint={
          <>
            If no push arrives but the job succeeded, double-check
            (a) you used a physical device, not a simulator, and (b) your
            phone has notifications enabled for Expo Go in system
            Settings.
          </>
        }
        solution={
          <>
            <p>You should see a system notification with your title + body. Tapping it deep-links into the app if you wired the response listener.</p>
            <p>You now have a full push pipeline: phone registers, API saves the token, job worker sends pushes. The same pattern handles welcome messages, order updates, chat notifications, anything.</p>
          </>
        }
      />

      <h2>What&apos;s next</h2>
      <p>
        Final chapter — <strong>Ship It</strong>. EAS Build for store-ready
        binaries, App Store + Play Store submission, OTA updates.
      </p>
    </>
  )
}
