import Link from 'next/link'
import type { Metadata } from 'next'
import { GuidePage, Section, H3, Outcome, Note, Warn, Checklist, Faq, Related } from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Webhooks in Go: HMAC verification, idempotency and replay protection',
  description:
    'Verifying Stripe and GitHub signatures in Go with constant-time compare, timestamp tolerance, a dedupe table, and sending webhooks that survive a bad receiver.',
  alternates: { canonical: 'https://gritframework.dev/guides/go-webhooks-hmac-idempotency' },
  openGraph: {
    title: 'Webhooks in Go: HMAC verification, idempotency and replay protection',
    description: 'Receiving them safely, and sending them so a bad receiver is not your problem.',
    url: 'https://gritframework.dev/guides/go-webhooks-hmac-idempotency',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="go-webhooks-hmac-idempotency"
      lead={
        <>
          <p className="mb-4">
            A webhook endpoint is a public, unauthenticated route that makes changes to your data
            on behalf of somebody else. That sentence should be uncomfortable, and the signature
            is what makes it acceptable.
          </p>
          <p>
            Three things have to be right: the signature must be verified correctly, the request
            must not be replayable, and the handler must survive being called twice with the same
            event, because it will be.
          </p>
        </>
      }
    >
      <Outcome>
        <p>Stripe and GitHub signature verification in Go, with the mistakes named.</p>
        <p>A dedupe table that makes a repeated delivery a no-op rather than a second charge.</p>
        <p>An outgoing webhook sender that a slow receiver cannot take down.</p>
      </Outcome>

      <Section id="receiving" title="Receiving: the raw body is the thing you verify">
        <p>
          The signature is computed over the exact bytes that were sent. Decode the JSON first and
          re-encode it and the bytes differ, by key order or whitespace, and nothing will ever
          verify.
        </p>
        <CodeBlock
          filename="internal/handlers/webhook.go"
          code={`func (h *WebhookHandler) Stripe(c *gin.Context) {
	// Read the raw body FIRST, before anything parses it. Cap it: this
	// is a public endpoint and the default body size is unlimited.
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, 1<<20))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := verifyStripe(body, c.GetHeader("Stripe-Signature"), h.secret); err != nil {
		// 400, and log it. A signature failure is either a
		// misconfiguration or somebody probing; both are worth knowing.
		log.Warn().Err(err).Str("ip", c.ClientIP()).Msg("webhook signature rejected")
		c.Status(http.StatusBadRequest)
		return
	}

	// Only now is it safe to look at the contents.
	var event StripeEvent
	if err := json.Unmarshal(body, &event); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	...
}`}
        />
        <Warn title="Gin middleware that reads the body breaks this">
          <p>
            A request body is a stream and can be read once. A logging or tracing middleware that
            reads it leaves an empty reader for your handler, and the signature check then fails
            on every request with no obvious cause. If you have such middleware, exclude the
            webhook routes or have it restore the body.
          </p>
        </Warn>
      </Section>

      <Section id="stripe" title="Stripe: a timestamp inside the signed value">
        <CodeBlock
          code={`// Stripe-Signature: t=1763750400,v1=5d41402abc...,v1=aabbcc...
func verifyStripe(body []byte, header, secret string) error {
	var timestamp string
	var signatures []string
	for _, part := range strings.Split(header, ",") {
		k, v, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		switch k {
		case "t":
			timestamp = v
		case "v1":
			signatures = append(signatures, v)
		}
	}
	if timestamp == "" || len(signatures) == 0 {
		return errors.New("malformed signature header")
	}

	// The timestamp is part of what was signed, so it cannot be altered
	// without invalidating the signature. Rejecting old ones is what
	// stops a captured request being replayed next week.
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("bad timestamp: %w", err)
	}
	if d := time.Since(time.Unix(ts, 0)); d > 5*time.Minute || d < -5*time.Minute {
		return fmt.Errorf("timestamp outside tolerance: %s", d)
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp + "." + string(body)))
	expected := mac.Sum(nil)

	// Several v1 values can be present during a secret rotation, so any
	// match is acceptance. hmac.Equal is constant time.
	for _, s := range signatures {
		got, err := hex.DecodeString(s)
		if err != nil {
			continue
		}
		if hmac.Equal(got, expected) {
			return nil
		}
	}
	return errors.New("no matching signature")
}`}
        />
        <p>
          Negative tolerance matters as much as positive: a receiver whose clock is behind sees
          every legitimate delivery as being from the future. Check the absolute difference, and
          keep the server&apos;s clock synchronised.
        </p>
      </Section>

      <Section id="github" title="GitHub: simpler, and the same two rules">
        <CodeBlock
          code={`// X-Hub-Signature-256: sha256=<hex>
func verifyGitHub(body []byte, header, secret string) error {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))

	// hmac.Equal and not ==. A normal string comparison returns as soon
	// as two bytes differ, so the time it takes reveals how much of the
	// prefix was right, which is enough to recover a signature.
	if !hmac.Equal([]byte(header), []byte(expected)) {
		return errors.New("signature mismatch")
	}
	return nil
}`}
        />
        <p>
          GitHub does not include a timestamp, so replay protection has to come from the delivery
          id instead, which the next section needs anyway.
        </p>
        <Note title="The three mistakes, named">
          <p>
            <strong>Comparing with <code>==</code>.</strong> Use <code>hmac.Equal</code>.
          </p>
          <p>
            <strong>Verifying a re-encoded body.</strong> Verify the bytes you received.
          </p>
          <p>
            <strong>Skipping verification in development.</strong> The flag that disables it gets
            deployed. Use the provider&apos;s CLI, which signs properly against a local secret.
          </p>
        </Note>
      </Section>

      <Section id="idempotency" title="It will arrive twice">
        <p>
          Every provider retries, and a retry happens whenever your response was slow, lost, or a
          500. From the sender&apos;s side an unacknowledged delivery and a failed one are
          identical, so a handler that succeeded and timed out gets called again with the same
          event.
        </p>
        <CodeBlock
          language="sql"
          code={`CREATE TABLE webhook_events (
  id            bigserial PRIMARY KEY,
  provider      text        NOT NULL,
  -- The provider's own id. The unique index is what arbitrates, so two
  -- concurrent deliveries cannot both proceed.
  external_id   text        NOT NULL,
  event_type    text        NOT NULL,
  payload       jsonb       NOT NULL,
  received_at   timestamptz NOT NULL DEFAULT now(),
  processed_at  timestamptz,
  UNIQUE (provider, external_id)
);`}
        />
        <CodeBlock
          code={`func (h *WebhookHandler) handle(ctx context.Context, provider string, e Event) error {
	return h.db.Transaction(func(tx *gorm.DB) error {
		rec := models.WebhookEvent{
			Provider: provider, ExternalID: e.ID,
			EventType: e.Type, Payload: e.Raw,
		}
		// The database decides, not a SELECT followed by an INSERT: two
		// deliveries arriving in the same millisecond would both pass a
		// check-then-act.
		err := tx.Create(&rec).Error
		if isUniqueViolation(err) {
			return nil   // already have it
		}
		if err != nil {
			return fmt.Errorf("record event: %w", err)
		}

		// Process in the SAME transaction as the claim. Claim first and
		// work after and a crash between them loses the event, with the
		// dedupe row preventing the retry from helping.
		if err := h.apply(ctx, tx, e); err != nil {
			return err
		}
		return tx.Model(&rec).Update("processed_at", time.Now()).Error
	})
}`}
        />
        <Warn title="Answer 200 quickly, then do the work">
          <p>
            Most providers time out in a handful of seconds and treat a slow response as a
            failure. If handling takes real time, store the event, return 200, and process from a
            queue. Doing the work inline is how a slow third-party call turns into a stream of
            duplicate deliveries.
          </p>
          <p>
            And return 200 for an event you have decided to ignore. A 400 or 404 for an event
            type you do not handle makes the provider retry it, and eventually disable your
            endpoint.
          </p>
        </Warn>
      </Section>

      <Section id="ordering" title="Events arrive out of order">
        <p>
          Retries mean a <code>subscription.updated</code> from 10:00 can land after the one from
          10:05. A handler that applies whatever it is given last will happily move a subscription
          backwards.
        </p>
        <CodeBlock
          code={`// Either refuse to apply an event older than what you have:
if sub.UpdatedFromEventAt.After(e.CreatedAt) {
	return nil   // a later state is already applied
}

// Or, better for anything that matters, treat the webhook as a
// notification rather than as data, and fetch the current state:
current, err := h.stripe.Subscriptions.Get(ctx, e.SubscriptionID)
if err != nil {
	return fmt.Errorf("fetch subscription: %w", err)
}
return h.applyState(ctx, tx, current)`}
        />
        <p>
          The second shape is more robust and costs an API call. For anything financial it is the
          right trade: the payload tells you something changed, and the provider tells you what
          is true now.
        </p>
      </Section>

      <Section id="sending" title="Sending webhooks">
        <p>
          The other direction, where your receivers are slow, down or malicious, and none of that
          may become your problem.
        </p>
        <Checklist
          items={[
            <>
              <strong>Deliver from a queue, never from the request.</strong> A receiver that takes
              thirty seconds must not hold your handler.
            </>,
            <>
              <strong>Sign the way you would want to verify:</strong> a timestamp and the body,
              HMAC-SHA256, in a documented header. The Standard Webhooks specification is a
              reasonable thing to copy rather than invent.
            </>,
            <>
              <strong>Send a unique delivery id,</strong> so your receivers can deduplicate. You
              are the provider in their version of the section above.
            </>,
            <>
              <strong>Retry with exponential backoff and jitter,</strong> then disable the
              endpoint and tell the customer. An endpoint retried forever is capacity spent on
              somebody who is not listening.
            </>,
            <>
              <strong>Time out aggressively,</strong> ten seconds at most, and cap redirects.
            </>,
            <>
              <strong>Treat the target URL as hostile.</strong> A customer-supplied URL pointed at{' '}
              <code>169.254.169.254</code> or <code>localhost</code> turns your sender into a
              request forgery tool. Resolve the hostname and refuse private addresses, on every
              redirect hop.
            </>,
          ]}
        />
        <CodeBlock
          code={`// Refuse private destinations at dial time, which is after DNS, so a
// public hostname resolving to a private address is still caught.
func safeDial(ctx context.Context, network, addr string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		if ip.IP.IsLoopback() || ip.IP.IsPrivate() || ip.IP.IsLinkLocalUnicast() {
			return nil, fmt.Errorf("refusing %s", ip.IP)
		}
	}
	return (&net.Dialer{Timeout: 5 * time.Second}).DialContext(ctx, network, addr)
}`}
        />
      </Section>

      <Section id="testing" title="Testing them">
        <CodeBlock
          terminal
          code={`# Stripe's CLI signs with your local secret and replays real events.
stripe listen --forward-to localhost:8080/webhooks/stripe
stripe trigger payment_intent.succeeded`}
        />
        <p>Three tests worth having, beyond the happy path:</p>
        <Checklist
          items={[
            <>
              <strong>A wrong signature is rejected.</strong> Then break the verification on
              purpose and confirm this test goes red, because a test that posts an unsigned body
              to a handler that never verifies also passes.
            </>,
            <>
              <strong>The same event twice has the same effect as once.</strong> Assert on the
              side effect, not on the response.
            </>,
            <>
              <strong>An old timestamp is rejected.</strong> The replay case, which is the one
              nobody writes.
            </>,
          ]}
        />
      </Section>

      <Section id="with-grit" title="Generated">
        <CodeBlock
          terminal
          code={`grit plugin add webhooks    # sending, to the Standard Webhooks spec
grit plugin add stripe      # receiving, with verification and dedupe`}
        />
        <p>
          The sending side gives you endpoint management per tenant, signing, the delivery queue
          with backoff, automatic disabling after repeated failure, and a delivery log with the
          request, the response and a replay button.
        </p>
        <p>
          The delivery log is what you will actually use. Every webhook integration eventually
          becomes a conversation that begins &quot;we never received it&quot;, and the only
          useful reply is the response code and body from the attempt.
        </p>
        <p>
          <Link href="/docs/plugins/webhooks" className="text-accent hover:underline">
            The webhooks reference
          </Link>{' '}
          covers both directions.
        </p>
      </Section>

      <Faq
        slug="go-webhooks-hmac-idempotency"
        items={[
          {
            q: 'Why does my signature verification always fail?',
            a: (
              <>
                Almost always the body. The signature covers the exact bytes received, so
                anything that parses and re-encodes the JSON breaks it, including a logging
                middleware that reads the body and a framework that binds it before your handler
                runs. Read the raw body first, verify, then parse.
              </>
            ),
            text:
              'Almost always the body. The signature covers the exact bytes received, so anything that parses and re-encodes the JSON breaks it, including a logging middleware that reads the body first. Read the raw body, verify, then parse.',
          },
          {
            q: 'Should I process the event inline or queue it?',
            a: (
              <>
                Store and queue, unless the work is genuinely a few milliseconds. Providers time
                out in seconds and treat a slow response as a failure, so inline processing of
                anything slow produces duplicate deliveries, which is the problem you were trying
                to avoid.
              </>
            ),
            text:
              'Store and queue unless the work takes a few milliseconds. Providers time out in seconds and treat slow responses as failures, so inline processing of slow work produces duplicate deliveries.',
          },
          {
            q: 'What status code should I return for an event I do not handle?',
            a: (
              <>
                200. A 400 or 404 makes the provider retry, and repeated failures eventually get
                your endpoint disabled, which takes the events you <em>do</em> care about with
                it. Acknowledge everything you have received and signed, and ignore what is not
                interesting.
              </>
            ),
            text:
              '200. A 400 or 404 causes retries, and repeated failures can get your endpoint disabled, taking the events you do care about with it. Acknowledge every verified delivery and ignore the uninteresting ones.',
          },
          {
            q: 'Is an unguessable URL an alternative to signing?',
            a: (
              <>
                No. A URL appears in logs, in proxies, in browser history and in whatever the
                customer pasted it into. It is also permanent: rotating it means coordinating with
                the sender. A signature is checkable, rotatable and does not leak by being
                written down.
              </>
            ),
            text:
              'No. A URL appears in logs, proxies and browser history, and rotating it requires coordinating with the sender. A signature is checkable, rotatable and does not leak by being written down.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/go-transactional-outbox',
            label: 'The transactional outbox',
            blurb: 'The same at-least-once problem, from the sending side.',
          },
          {
            href: '/guides/go-background-jobs-asynq',
            label: 'Background jobs in Go',
            blurb: 'Where the work goes after you have answered 200.',
          },
          {
            href: '/docs/plugins/webhooks',
            label: 'Webhooks reference',
            blurb: 'Endpoint management, signing, the delivery log and replay.',
          },
          {
            href: '/docs/plugins/stripe',
            label: 'Stripe',
            blurb: 'Receiving, verifying and deduplicating payment events.',
          },
        ]}
      />
    </GuidePage>
  )
}
