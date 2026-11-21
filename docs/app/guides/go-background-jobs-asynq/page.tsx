import Link from 'next/link'
import type { Metadata } from 'next'
import {
  GuidePage,
  Section,
  Step,
  H3,
  Outcome,
  Note,
  Warn,
  Mistake,
  Checklist,
  Faq,
  Related,
} from '@/components/guide'
import { CodeBlock } from '@/components/code-block'

export const metadata: Metadata = {
  title: 'Background jobs in Go: send email without slowing the API',
  description:
    'Why a queue exists, writing a worker with asynq and Redis, retries and idempotency, scheduled jobs, and testing mail locally before it is real.',
  alternates: { canonical: 'https://gritframework.dev/guides/go-background-jobs-asynq' },
  openGraph: {
    title: 'Background jobs in Go: send email without slowing the API',
    description: 'A queue, a worker, retries that are safe, and local mail you can read.',
    url: 'https://gritframework.dev/guides/go-background-jobs-asynq',
    type: 'article',
  },
}

export default function Page() {
  return (
    <GuidePage
      slug="go-background-jobs-asynq"
      lead={
        <>
          <p className="mb-4">
            Your registration endpoint creates a user in four milliseconds and then spends eight
            hundred waiting for an email provider. The user watches a spinner for the provider,
            not for you, and when the provider has a bad minute your sign-up returns 500 for an
            account that was created perfectly.
          </p>
          <p>
            A queue separates those. Save the user, write down that an email is owed, answer the
            request. Something else sends it. This is how to do that in Go, and the parts that are
            harder than they look are retries and the second delivery nobody planned for.
          </p>
        </>
      }
    >
      <Outcome>
        <p>A worker process, a queue on Redis, and a registration endpoint that returns at once.</p>
        <p>Retries with backoff, and jobs written so a retry cannot do the work twice.</p>
        <p>Email you can read locally without sending anything to a real address.</p>
      </Outcome>

      <Section id="why" title="What belongs in a job">
        <p>Three tests. If any is yes, it belongs behind a queue.</p>
        <Checklist
          items={[
            <>
              <strong>It talks to somebody else&apos;s service.</strong> Email, SMS, payment
              webhooks, a CRM. Their downtime becomes your downtime otherwise.
            </>,
            <>
              <strong>It takes longer than the user should wait.</strong> Resizing images,
              generating a PDF, importing a spreadsheet, rebuilding a report.
            </>,
            <>
              <strong>It can fail without the request failing.</strong> A welcome email that does
              not arrive is worth retrying. It is not worth failing a registration over.
            </>,
          ]}
        />
        <p>
          And the opposite, because queues get overused: anything the user must see the result of
          belongs in the request. A job that the frontend then polls for is usually a sign the
          work was quick enough to do inline.
        </p>
      </Section>

      <Section id="naive" title="The version without a queue, and why it fails">
        <CodeBlock
          code={`func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*models.User, error) {
	user, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// 800ms, and a 500 for the whole registration if the provider is down.
	if err := s.mail.SendWelcome(ctx, user); err != nil {
		return nil, fmt.Errorf("send welcome: %w", err)
	}
	return user, nil
}`}
        />
        <p>
          The usual first fix is <code>go s.mail.SendWelcome(...)</code>, and it is worse than it
          looks. A goroutine is lost on deploy, lost on a crash, retries nothing, reports nothing,
          and the request context it captured is cancelled the moment the response is written, so
          the call often fails immediately. It is a queue with no durability, no retry and no
          visibility.
        </p>
      </Section>

      <Step n={1} title="Define the task">
        <p>
          A task is a name and a payload. Keep the payload small and keep it to identifiers: the
          worker reads the current row rather than acting on a snapshot that may be minutes old.
        </p>
        <CodeBlock
          filename="internal/jobs/tasks.go"
          code={`package jobs

const (
	TypeWelcomeEmail = "email:welcome"
	TypeResizeImage  = "image:resize"
	TypeWeeklyDigest = "email:digest"
)

type WelcomeEmailPayload struct {
	UserID uint \`json:"user_id"\`
}

func NewWelcomeEmail(userID uint) (*asynq.Task, error) {
	payload, err := json.Marshal(WelcomeEmailPayload{UserID: userID})
	if err != nil {
		return nil, fmt.Errorf("marshal welcome payload: %w", err)
	}
	return asynq.NewTask(TypeWelcomeEmail, payload,
		asynq.MaxRetry(5),
		asynq.Timeout(30*time.Second),
		asynq.Queue("default"),
	), nil
}`}
        />
        <Note title="Identifiers, not objects">
          <p>
            Putting the whole user in the payload means the worker emails a name that may have
            changed, or an address that was corrected two minutes ago. It also means the payload
            grows with the model and old jobs in the queue carry a shape your code no longer
            expects. An id is stable and small.
          </p>
        </Note>
      </Step>

      <Step n={2} title="Enqueue it">
        <CodeBlock
          code={`func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*models.User, error) {
	user, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	task, err := jobs.NewWelcomeEmail(user.ID)
	if err != nil {
		return nil, fmt.Errorf("build welcome task: %w", err)
	}
	if _, err := s.queue.EnqueueContext(ctx, task); err != nil {
		// The user exists and that is what was asked for. Log it loudly
		// and return success: failing the registration because a welcome
		// email could not be queued is the wrong trade.
		log.Error().Err(err).Uint("user_id", user.ID).Msg("queue welcome email")
	}

	return user, nil   // back in milliseconds
}`}
        />
        <p>
          Scheduling later is the same call with one option, which is how a trial reminder or a
          delayed nudge is written:
        </p>
        <CodeBlock
          code={`s.queue.Enqueue(task, asynq.ProcessIn(24*time.Hour))
s.queue.Enqueue(task, asynq.ProcessAt(trialEnd.Add(-48*time.Hour)))`}
        />
      </Step>

      <Step n={3} title="Write the worker">
        <p>
          A separate process. It must be separate: you will want to restart the API without
          interrupting a long job, and to run three workers against one API.
        </p>
        <CodeBlock
          filename="cmd/worker/main.go"
          code={`func main() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: cfg.RedisAddr},
		asynq.Config{
			// How many jobs run at once. Too high and the workers fight
			// the API for database connections; the pool is shared.
			Concurrency: 10,

			// Weights, not priorities: critical gets six of every ten
			// slots, and low still gets one, so nothing starves.
			Queues: map[string]int{"critical": 6, "default": 3, "low": 1},

			RetryDelayFunc: func(n int, err error, t *asynq.Task) time.Duration {
				// Exponential, with a cap. Retrying every second against a
				// service that is down is a denial of service you are
				// running against somebody else.
				d := time.Duration(math.Pow(2, float64(n))) * time.Second
				if d > 10*time.Minute {
					d = 10 * time.Minute
				}
				return d
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.HandleFunc(jobs.TypeWelcomeEmail, handlers.HandleWelcomeEmail)
	mux.HandleFunc(jobs.TypeResizeImage, handlers.HandleResizeImage)

	if err := srv.Run(mux); err != nil {
		log.Fatal().Err(err).Msg("worker stopped")
	}
}`}
        />
        <CodeBlock
          filename="internal/jobs/handlers/email.go"
          code={`func HandleWelcomeEmail(ctx context.Context, t *asynq.Task) error {
	var p jobs.WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		// A payload that cannot be parsed will never parse. Saying so
		// stops it consuming five retries on its way to the dead queue.
		return fmt.Errorf("unmarshal: %v: %w", err, asynq.SkipRetry)
	}

	user, err := userRepo.ByID(ctx, p.UserID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Deleted between enqueue and processing. Not a failure.
		return nil
	}
	if err != nil {
		return fmt.Errorf("load user %d: %w", p.UserID, err)   // retry
	}

	if err := mailer.SendWelcome(ctx, user); err != nil {
		return fmt.Errorf("send welcome to %d: %w", p.UserID, err)   // retry
	}
	return nil
}`}
        />
        <p>
          The return value is the whole protocol: <code>nil</code> means done,{' '}
          <code>asynq.SkipRetry</code> means never going to work, and any other error means try
          again later. Getting that distinction right is most of what makes a queue pleasant to
          operate.
        </p>
      </Step>

      <Section id="idempotency" title="A job will run twice. Design for it.">
        <p>
          Not a rare edge case. A worker that completes the work and is killed before recording
          success leaves the job to be retried, and the second run repeats whatever the first one
          did. Across enough jobs this happens weekly.
        </p>
        <p>
          A duplicate log line is nothing. A duplicate email is embarrassing. A duplicate charge is
          a refund and an apology.
        </p>
        <H3>Record that you did it, in the same transaction</H3>
        <CodeBlock
          code={`func HandleWelcomeEmail(ctx context.Context, t *asynq.Task) error {
	var p jobs.WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("unmarshal: %v: %w", err, asynq.SkipRetry)
	}

	// A unique index on (user_id, kind) makes the second insert fail
	// rather than the second email send. The database is the only thing
	// that can arbitrate this correctly when two workers race.
	sent, err := notify.Once(ctx, db, p.UserID, "welcome")
	if err != nil {
		return fmt.Errorf("claim welcome: %w", err)
	}
	if !sent {
		return nil   // somebody already did it
	}

	return mailer.SendWelcome(ctx, p.UserID)
}`}
        />
        <p>
          For anything that costs money, the external service usually offers an idempotency key.
          Use it, and derive the key from the thing being paid for rather than generating a new
          one per attempt, which defeats the purpose entirely.
        </p>
        <Warn title="Enqueue after the commit, not inside it">
          <p>
            Queue a job inside a database transaction and the worker can pick it up before the
            transaction commits, then fail to find the row it was told about. Redis knows nothing
            about your transaction. Either enqueue after the commit, or write the job to a table
            in the same transaction and have a separate process move it to the queue, which is the
            outbox pattern.
          </p>
        </Warn>
      </Section>

      <Section id="scheduled" title="Jobs on a schedule">
        <CodeBlock
          code={`scheduler := asynq.NewScheduler(redisOpt, &asynq.SchedulerOpts{
	// Pin the zone. "Every day at 09:00" on a UTC server is not what
	// anybody meant, and it moves twice a year.
	Location: time.FixedZone("Europe/London", 0),
})

scheduler.Register("0 9 * * MON", jobs.MustNewWeeklyDigest())
scheduler.Register("@every 15m",  jobs.MustNewSyncInventory())
scheduler.Register("0 3 * * *",   jobs.MustNewCleanupExpired())`}
        />
        <Warn title="Exactly one scheduler, however many workers">
          <p>
            Run the scheduler inside the worker binary and then scale to three workers, and the
            nightly digest goes out three times. Either run the scheduler as its own single
            process, or elect a leader among the workers. This is one of the most common
            production surprises when a service is scaled for the first time.
          </p>
        </Warn>
      </Section>

      <Section id="mail" title="Reading your own email, locally">
        <p>
          Do not point development at a real provider. You will eventually send a test to a real
          customer. Mailhog is a fake SMTP server with a web inbox: everything your app sends
          lands there and goes nowhere.
        </p>
        <CodeBlock
          language="yaml"
          filename="docker-compose.yml"
          code={`  mailhog:
    image: mailhog/mailhog
    ports:
      - "127.0.0.1:1025:1025"   # SMTP, what the app connects to
      - "127.0.0.1:8025:8025"   # the inbox you open in a browser`}
        />
        <CodeBlock
          language="bash"
          code={`# development
MAIL_DRIVER=smtp
MAIL_HOST=localhost
MAIL_PORT=1025

# production
MAIL_DRIVER=resend
RESEND_API_KEY=re_...`}
        />
        <p>
          Open <code>localhost:8025</code> and read the mail, including the HTML as it will
          actually render, which catches the template problems that a logged string never shows.
        </p>
      </Section>

      <Section id="operating" title="Operating it">
        <p>
          A queue is a system with its own failure modes, and all of them are invisible until you
          look.
        </p>
        <Checklist
          items={[
            <>
              <strong>Watch the depth.</strong> A queue that grows steadily means workers are
              slower than producers, and the first symptom users notice is email arriving hours
              late.
            </>,
            <>
              <strong>Read the dead letter queue.</strong> Jobs that exhausted their retries land
              there, and an unmonitored dead queue is a silent feature outage.
            </>,
            <>
              <strong>Alert on the oldest pending job,</strong> not just the count. Ten thousand
              jobs enqueued a second ago is healthy; one job pending for six hours is a stuck
              worker.
            </>,
            <>
              <strong>Separate the queues by urgency.</strong> A password reset must not sit
              behind a bulk import of forty thousand rows.
            </>,
            <>
              <strong>Shut down gracefully.</strong> On deploy, stop accepting new jobs and let
              the running ones finish. asynq does this on SIGTERM; make sure your container runner
              actually sends one and waits.
            </>,
          ]}
        />
      </Section>

      <Section id="mistakes" title="When it goes wrong">
        <Mistake error="dial tcp 127.0.0.1:6379: connect: connection refused">
          <p>
            Redis is not running, or is not where the worker thinks. In Docker it is the service
            name and the container port, <code>redis:6379</code>, not{' '}
            <code>localhost:6380</code>, which is the host mapping.
          </p>
        </Mistake>

        <Mistake error="The job never runs, and nothing is logged">
          <p>
            Three usual causes. The worker process is not running at all, which is easy when it is
            a second binary you forgot to start. The task type string does not match between
            enqueue and <code>HandleFunc</code>, and a typo there is silent. Or the job was
            enqueued to a queue name the worker&apos;s <code>Queues</code> map does not list, in
            which case it sits there forever.
          </p>
        </Mistake>

        <Mistake error="Emails arrive twice">
          <p>
            Either the job is not idempotent and was retried, or two workers both have a
            scheduler. Check the scheduler first: it is the cause when the duplicates arrive at a
            round time of day.
          </p>
        </Mistake>

        <Mistake error="The worker is slow and the API gets slower too">
          <p>
            They share the database connection pool. Ten concurrent workers each holding a
            connection, against a pool of ten, leaves the API waiting. Give the worker its own
            pool with its own limit, and size both against what the database actually allows.
          </p>
        </Mistake>

        <Mistake error="context deadline exceeded inside a job">
          <p>
            The task timeout fired. Either the work genuinely takes longer, in which case raise{' '}
            <code>asynq.Timeout</code> and consider splitting it, or something downstream is
            hanging with no timeout of its own. An HTTP client with no timeout will wait forever
            by default, which is how one slow provider blocks a whole worker.
          </p>
        </Mistake>
      </Section>

      <Section id="with-grit" title="Already wired">
        <p>
          A Grit project ships asynq configured, a worker binary, email and image workers, the
          scheduler, Mailhog in the development compose file, Resend in production, and an admin
          screen showing pending, scheduled, retrying and dead jobs with a retry button.
        </p>
        <CodeBlock
          terminal
          code={`grit generate job SendInvoiceReminder --payload "invoice_id:uint"
grit start   # api, worker and scheduler together`}
        />
        <p>
          The <code>notify.Once</code> helper used above is part of it, because &quot;send this
          exactly once&quot; turned out to be the thing every project rewrote.{' '}
          <Link href="/docs/batteries/jobs" className="text-accent hover:underline">
            The background jobs reference
          </Link>{' '}
          covers the queues, the scheduler and the admin screens.
        </p>
      </Section>

      <Faq
        slug="go-background-jobs-asynq"
        items={[
          {
            q: 'Why not just use a goroutine?',
            a: (
              <>
                Because a goroutine is lost on deploy, lost on a crash, retries nothing and reports
                nothing. It also captures the request context, which is cancelled as soon as the
                response is written, so the work frequently fails immediately for a reason nobody
                sees. A queue is a goroutine that survives, retries and can be inspected.
              </>
            ),
            text:
              'A goroutine is lost on deploy and on crash, retries nothing and reports nothing. It also captures the request context, which is cancelled once the response is written, so the work often fails immediately and invisibly. A queue is the same work made durable, retryable and inspectable.',
          },
          {
            q: 'asynq, River, Machinery, or a database table?',
            a: (
              <>
                asynq if you already run Redis: mature, good defaults, a usable dashboard. River if
                you would rather not add Redis, since it uses Postgres and can therefore enqueue
                inside the same transaction as your write, which removes a whole class of bug. A
                plain table polled every few seconds is entirely respectable for low volume and is
                less code than people expect.
              </>
            ),
            text:
              'asynq if you already run Redis: mature, sensible defaults and a usable dashboard. River if you would rather not add Redis, since it uses Postgres and can enqueue inside the same transaction as the write, removing a class of bug. A plain polled table is respectable at low volume.',
          },
          {
            q: 'How many workers should I run?',
            a: (
              <>
                Start with one process and a concurrency of ten, then measure. The binding
                constraint is usually the database connection pool or the rate limit of whatever
                you are calling, not CPU. Raising concurrency past the pool size makes jobs wait
                on connections and slows the API at the same time.
              </>
            ),
            text:
              'Start with one process at concurrency ten and measure. The binding constraint is usually the database connection pool or an external rate limit rather than CPU, and raising concurrency past the pool size makes jobs queue on connections while slowing the API.',
          },
          {
            q: 'What happens to jobs that fail every time?',
            a: (
              <>
                After the retry limit they go to a dead letter queue and stay there. That is the
                correct behaviour and it is also a trap: nobody looks. Put the dead queue count on
                a dashboard, and alert on it, or a feature stops working and the only evidence is
                a number nobody reads.
              </>
            ),
            text:
              'After the retry limit they move to a dead letter queue and remain there. That is correct, and also a trap, because nobody looks. Put the dead queue count on a dashboard and alert on it, or a broken feature leaves no visible evidence.',
          },
        ]}
      />

      <Related
        links={[
          {
            href: '/guides/go-file-upload-s3-presigned',
            label: 'File uploads and image renditions',
            blurb: 'The most common first real job: resizing an image after it arrives.',
          },
          {
            href: '/guides/docker-compose-postgres-redis',
            label: 'Redis and Mailhog with Docker Compose',
            blurb: 'The queue and the fake inbox, running locally.',
          },
          {
            href: '/docs/batteries/jobs',
            label: 'Background jobs reference',
            blurb: 'The generated workers, the scheduler and the admin screens.',
          },
          {
            href: '/docs/backend/outbox',
            label: 'The transactional outbox',
            blurb: 'Enqueuing in the same transaction as the write, when losing a job is not an option.',
          },
        ]}
      />
    </GuidePage>
  )
}
