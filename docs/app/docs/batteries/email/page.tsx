import Link from "next/link";
import { ArrowRight, ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";
import { CodeBlock } from "@/components/code-block";
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/batteries/email')

export default function EmailPage() {
  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">
                Batteries
              </span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">
                Email System
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed">
                Grit integrates with Resend for transactional email delivery
                using raw{" "}
                <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                  net/http
                </code>{" "}
                calls (no SDK dependency). It includes a Mailer service, Go
                template rendering, and four built-in email templates styled
                with the Grit dark theme.
              </p>
            </div>

            <div className="prose-grit">
              {/* Configuration */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Configuration
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Email sending requires a Resend API key and a verified sender
                  address. For local development, Mailhog is included in Docker
                  Compose for testing without sending real emails.
                </p>

                <CodeBlock language="bash" filename=".env" code={`# Email Configuration
RESEND_API_KEY=re_xxxxxxxxxxxxxxxxxxxx
MAIL_FROM=noreply@yourdomain.com`} />
              </div>

              {/* Mailer Service */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Mailer Service
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The Mailer service at{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    internal/mail/mailer.go
                  </code>{" "}
                  provides two methods for sending emails:{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    Send()
                  </code>{" "}
                  for template-based emails and{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    SendRaw()
                  </code>{" "}
                  for raw HTML emails.
                </p>

                <CodeBlock language="go" filename="internal/mail/mailer.go" code={`// Mailer sends emails via the Resend API.
type Mailer struct {
    apiKey string
    from   string
    client *http.Client
}

// New creates a new Mailer instance.
func New(apiKey, from string) *Mailer

// SendOptions configures an email to send.
type SendOptions struct {
    To       string
    Subject  string
    Template string                 // Template name (e.g., "welcome")
    Data     map[string]interface{} // Template variables
}

// Send renders a template and sends the email via Resend.
func (m *Mailer) Send(ctx context.Context, opts SendOptions) error

// SendRaw sends an email with raw HTML content (no template).
func (m *Mailer) SendRaw(ctx context.Context, to, subject, htmlBody string) error`} />

                <h3 className="text-lg font-semibold tracking-tight mb-3 mt-6">
                  Usage Examples
                </h3>

                <CodeBlock language="go" filename="send-welcome-email.go" code={`// Send a welcome email using a template
err := mailer.Send(ctx, mail.SendOptions{
    To:       "user@example.com",
    Subject:  "Welcome to MyApp!",
    Template: "welcome",
    Data: map[string]interface{}{
        "AppName":      "MyApp",
        "Name":         "John",
        "DashboardURL": "https://myapp.com/dashboard",
        "Year":         time.Now().Year(),
    },
})

// Send a password reset email
err = mailer.Send(ctx, mail.SendOptions{
    To:       "user@example.com",
    Subject:  "Reset Your Password",
    Template: "password-reset",
    Data: map[string]interface{}{
        "AppName":  "MyApp",
        "ResetURL": "https://myapp.com/reset?token=abc123",
        "Year":     time.Now().Year(),
    },
})

// Send a raw HTML email
err = mailer.SendRaw(ctx, "user@example.com", "Custom Email", "<h1>Hello!</h1>")`} />
              </div>

              {/* Template Rendering */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Template Rendering
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Email templates use Go&apos;s{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    html/template
                  </code>{" "}
                  package. Templates are stored as Go constants in{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    internal/mail/templates.go
                  </code>{" "}
                  and compiled at render time. Variables use the standard Go
                  template syntax:{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">{`{{.VariableName}}`}</code>
                  .
                </p>

                <CodeBlock language="go" filename="internal/mail/templates.go (registry)" code={`// EmailTemplates contains all available email templates.
var EmailTemplates = map[string]string{
    "welcome":            welcomeTemplate,
    "password-reset":     passwordResetTemplate,
    "email-verification": emailVerificationTemplate,
    "notification":       notificationTemplate,
}`} />
              </div>

              {/* Built-in Templates */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Built-in Templates
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Grit ships four pre-built email templates, all styled with the
                  dark theme (background: #0a0a0f, accent: #6c5ce7). Each
                  template is a complete HTML document with inline CSS for
                  maximum email client compatibility.
                </p>

                <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
                  <table className="w-full text-sm">
                    <thead>
                      <tr className="border-b border-border/30 bg-accent/20">
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">
                          Template
                        </th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">
                          Variables
                        </th>
                        <th className="text-left px-4 py-2.5 font-medium text-foreground/80">
                          Use Case
                        </th>
                      </tr>
                    </thead>
                    <tbody className="text-muted-foreground">
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">
                          welcome
                        </td>
                        <td className="px-4 py-2.5 text-xs">
                          AppName, Name, DashboardURL, Year
                        </td>
                        <td className="px-4 py-2.5">After user registration</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">
                          password-reset
                        </td>
                        <td className="px-4 py-2.5 text-xs">
                          AppName, ResetURL, Year
                        </td>
                        <td className="px-4 py-2.5">Password reset request</td>
                      </tr>
                      <tr className="border-b border-border/20">
                        <td className="px-4 py-2.5 font-mono text-xs">
                          email-verification
                        </td>
                        <td className="px-4 py-2.5 text-xs">
                          AppName, VerifyURL, Year
                        </td>
                        <td className="px-4 py-2.5">
                          Email address verification
                        </td>
                      </tr>
                      <tr>
                        <td className="px-4 py-2.5 font-mono text-xs">
                          notification
                        </td>
                        <td className="px-4 py-2.5 text-xs">
                          AppName, Title, Message, ActionURL, ActionText, Year
                        </td>
                        <td className="px-4 py-2.5">Generic notifications</td>
                      </tr>
                    </tbody>
                  </table>
                </div>

                <h3 className="text-lg font-semibold tracking-tight mb-3">
                  Template Structure
                </h3>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Each template follows a consistent structure: dark background,
                  card container with logo, heading, body text, optional CTA
                  button, and footer. Here is the welcome template:
                </p>

                <CodeBlock language="markup" filename="welcome template (simplified)" code={`<body style="background-color: #0a0a0f; color: #e8e8f0;">
  <div class="container">
    <div class="card">
      <div class="logo">{{.AppName}}</div>
      <h1>Welcome, {{.Name}}!</h1>
      <p>Thanks for signing up. Your account is ready to use.</p>
      <p>Get started by exploring the dashboard:</p>
      <p style="text-align: center;">
        <a href="{{.DashboardURL}}" class="btn">Go to Dashboard</a>
      </p>
    </div>
    <div class="footer">
      <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
    </div>
  </div>
</body>`} />
              </div>

              {/* Adding Custom Templates */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Adding Custom Templates
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  To add a new email template, define it as a Go constant in{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    templates.go
                  </code>{" "}
                  and register it in the{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    EmailTemplates
                  </code>{" "}
                  map.
                </p>

                <CodeBlock language="go" filename="internal/mail/templates.go" code={`// Add to the EmailTemplates map:
var EmailTemplates = map[string]string{
    "welcome":            welcomeTemplate,
    "password-reset":     passwordResetTemplate,
    "email-verification": emailVerificationTemplate,
    "notification":       notificationTemplate,
    "invoice":            invoiceTemplate,  // Your custom template
}

// Define the template:
const invoiceTemplate = \`<!DOCTYPE html>
<html>
<head>...</head>
<body>
  <div class="container">
    <div class="card">
      <div class="logo">{{.AppName}}</div>
      <h1>Invoice #{{.InvoiceNumber}}</h1>
      <p>Amount: {{.Amount}}</p>
      <p>Due: {{.DueDate}}</p>
      <a href="{{.PaymentURL}}" class="btn">Pay Now</a>
    </div>
  </div>
</body>
</html>\`

// Send it:
mailer.Send(ctx, mail.SendOptions{
    To:       "customer@example.com",
    Subject:  "Invoice #INV-001",
    Template: "invoice",
    Data: map[string]interface{}{
        "AppName":       "MyApp",
        "InvoiceNumber": "INV-001",
        "Amount":        "99.00",
        "DueDate":       "Feb 28, 2026",
        "PaymentURL":    "https://myapp.com/pay/inv-001",
    },
})`} />
              </div>

              {/* Background Email Sending */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Background Email Sending
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  For better API response times, emails should be sent via
                  background jobs rather than inline. Grit&apos;s job queue
                  client provides a{" "}
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">
                    EnqueueSendEmail()
                  </code>{" "}
                  method that sends the email asynchronously.
                </p>

                <CodeBlock language="go" filename="async-email.go" code={`// Instead of sending synchronously:
// mailer.Send(ctx, opts)

// Enqueue a background job (returns immediately):
err := jobsClient.EnqueueSendEmail(
    "user@example.com",
    "Welcome to MyApp!",
    "welcome",
    map[string]interface{}{
        "AppName":      "MyApp",
        "Name":         "John",
        "DashboardURL": "https://myapp.com/dashboard",
        "Year":         2026,
    },
)
// The asynq worker picks it up and calls mailer.Send()`} />

                <div className="p-4 rounded-lg border border-primary/20 bg-primary/5">
                  <p className="text-sm text-foreground/80">
                    <strong>Admin email preview:</strong> The admin panel
                    includes an email template preview page where you can view
                    how each template renders with sample data before sending.
                    See the{" "}
                    <Link
                      href="/docs/batteries/jobs"
                      className="text-primary hover:underline"
                    >
                      Background Jobs
                    </Link>{" "}
                    page for the full job queue documentation.
                  </p>
                </div>
              </div>
            </div>

            {/* Nav */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/batteries/storage" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  File Storage
                </Link>
              </Button>
              <Button
                variant="ghost"
                size="sm"
                asChild
                className="text-muted-foreground/60 hover:text-foreground"
              >
                <Link href="/docs/batteries/jobs" className="gap-1.5">
                  Background Jobs
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
