package scaffold

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairMailQueueAndPreview gives an existing project what mail.Queue and the
// Mail Preview need outside internal/mail, which writeMailFiles has already
// delivered: jobs.Client.EnqueuePayload, the worker sending a queued Message,
// the two preview routes, and the ticket email honouring
// notifications.email_enabled. Called from repairMailDrivers once the package
// is in place, so it reaches single projects with it.
func repairMailQueueAndPreview(root string, opts Options, m *manifest.Manifest) error {
	apiRoot := opts.APIRoot(root)
	mailDir := filepath.Join(apiRoot, "internal", "mail")
	if !fileContains(filepath.Join(mailDir, "queue.go"), "type QueuedMessage struct") ||
		!fileContains(filepath.Join(mailDir, "preview.go"), "func LookupTemplate(") {
		return nil
	}

	steps := []struct {
		path string
		when bool
		fn   func(string) (string, []string, []string)
	}{
		{filepath.Join(apiRoot, "internal", "jobs", "client.go"), true, repairJobsEnqueuePayloadSource},
		{filepath.Join(apiRoot, "internal", "jobs", "workers.go"), true, repairJobsWorkerMailSource},
		{filepath.Join(apiRoot, "internal", "routes", "routes.go"),
			fileContains(filepath.Join(apiRoot, "internal", "handlers", "mail_preview.go"), "type MailPreviewHandler struct"),
			repairRoutesMailPreviewSource},
		// Only where the setting is declared: settings.Bool on a key the project
		// never defined is false, which would silence the email instead.
		{filepath.Join(apiRoot, "internal", "services", "ticket_mail.go"),
			settingsHave(filepath.Join(apiRoot, "internal", "settings"), "func Bool(ctx context.Context, key string) bool") &&
				settingsHave(filepath.Join(apiRoot, "internal", "settings"), `"notifications.email_enabled"`),
			repairTicketMailSettingSource},
	}
	for _, s := range steps {
		if !s.when || !fileExists(s.path) {
			continue
		}
		if err := repairSourceFile(root, m, s.path, s.fn); err != nil {
			return err
		}
	}
	return nil
}

// settingsHave reports whether any file of the project's settings package
// contains needle.
func settingsHave(dir, needle string) bool {
	files, err := goSources(dir)
	if err != nil {
		return false
	}
	for _, f := range files {
		if fileContains(f, needle) {
			return true
		}
	}
	return false
}

func repairJobsEnqueuePayloadSource(src string) (string, []string, []string) {
	if strings.Contains(src, "func (c *Client) EnqueuePayload(") || !strings.Contains(src, "type Client struct") {
		return src, nil, nil
	}
	if !strings.Contains(src, "func (c *Client) Enqueue(ctx context.Context, taskType string, payload any") ||
		!strings.Contains(src, "ErrDuplicateTask") || !strings.Contains(src, "\t\"errors\"\n") || !strings.Contains(src, "\t\"context\"\n") {
		return src, nil, []string{"jobs.Client is not the one Grit wrote: add an EnqueuePayload(ctx, taskType string, payload []byte) error method, which mail.Queue sends through"}
	}
	const next = "}\n\n// EnqueueProcessImage enqueues an image processing job."
	change := []string{"has EnqueuePayload, so mail.Queue can put a message on the email:send job"}
	if strings.Count(src, next) == 1 {
		return strings.Replace(src, next, "}\n"+jobsEnqueuePayloadFunc+"\n// EnqueueProcessImage enqueues an image processing job.", 1), change, nil
	}
	return strings.TrimRight(src, "\n") + "\n" + jobsEnqueuePayloadFunc, change, nil
}

func repairJobsWorkerMailSource(src string) (string, []string, []string) {
	if strings.Contains(src, "mail.QueuedMessage") || !strings.Contains(src, "func handleEmailSend(") {
		return src, nil, nil
	}
	if strings.Count(src, jobsWorkerMailOld) != 1 || !strings.Contains(src, "\t\"encoding/json\"\n") ||
		!strings.Contains(src, "\t\"github.com/hibiken/asynq\"\n") || !strings.Contains(src, "/internal/mail\"\n") {
		return src, nil, []string{"handleEmailSend is not the one Grit wrote: decode a mail.QueuedMessage first and send its Message with deps.Mailer.SendMessage, or mail.Queue's messages are sent as empty templates"}
	}
	return strings.Replace(src, jobsWorkerMailOld, jobsWorkerMailNew, 1),
		[]string{"the email:send worker sends the messages mail.Queue queues"}, nil
}

func repairRoutesMailPreviewSource(src string) (string, []string, []string) {
	if strings.Contains(src, "MailPreviewHandler") || !strings.Contains(src, "staff := v1.Group(") {
		return src, nil, nil
	}
	if strings.Count(src, routesMailPreviewAnchor) != 1 {
		return src, nil, []string{"the staff routes are not the ones Grit wrote: mount handlers.MailPreviewHandler at GET /admin/mail/templates and /admin/mail/preview/:template, or the admin's Mail Preview has nothing to show"}
	}
	return strings.Replace(src, routesMailPreviewAnchor, routesMailPreviewAnchor+routesMailPreview, 1),
		[]string{"mounts /api/admin/mail/templates and /api/admin/mail/preview/:template for the admin's Mail Preview"}, nil
}

var ticketMailModelsImport = regexp.MustCompile(`\n\t"([^"\n]+)/internal/models"\n\)\n`)

func repairTicketMailSettingSource(src string) (string, []string, []string) {
	// The check has to be inside SendTicketCreatedEmail specifically: the same
	// setting is read by QueueTicketCreatedEmail further down the file, and
	// looking for the setting name alone found that one and stopped.
	if strings.Contains(src, ticketMailFuncOpen+ticketMailSettingCheck) || !strings.Contains(src, "func SendTicketCreatedEmail(") {
		return src, nil, nil
	}
	loc := ticketMailModelsImport.FindAllStringSubmatchIndex(src, -1)
	if strings.Count(src, ticketMailFuncOpen) != 1 || len(loc) != 1 ||
		!strings.Contains(src, "\t\"context\"\n") || !strings.Contains(src, "\t\"log\"\n") {
		return src, nil, []string{"SendTicketCreatedEmail is not the one Grit wrote: return early when settings.Bool(ctx, \"notifications.email_enabled\") is false, or turning notification emails off does not stop it"}
	}
	module := src[loc[0][2]:loc[0][3]]
	importLine := "\n\t\"" + module + "/internal/models\"\n\t\"" + module + "/internal/settings\"\n)\n"
	out := src[:loc[0][0]] + importLine + src[loc[0][1]:]
	out = strings.Replace(out, ticketMailFuncOpen, ticketMailFuncOpen+ticketMailSettingCheck, 1)
	return out, []string{"the new-ticket email honours the notifications.email_enabled setting"}, nil
}
