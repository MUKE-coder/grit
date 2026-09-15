package scaffold

import (
	"strings"
	"testing"
)

func TestMailQueueAndPreviewInTemplates(t *testing.T) {
	if !strings.Contains(jobsClientGo(), jobsEnqueuePayloadFunc) {
		t.Error("jobs.Client has no EnqueuePayload for mail.Queue")
	}
	if !strings.Contains(jobsWorkersGo(), jobsWorkerMailNew) {
		t.Error("the email:send worker does not send a queued mail.Message")
	}
	if !strings.Contains(apiRoutesGo(), routesMailPreviewAnchor+routesMailPreview) {
		t.Error("routes.go does not mount the Mail Preview endpoints")
	}
	if !strings.Contains(apiDocsRoutesGo(), "GET /api/v1/admin/mail/preview/:template") {
		t.Error("the API reference does not describe the Mail Preview endpoints")
	}
	ticket := ticketMailGo()
	if !strings.Contains(ticket, ticketMailFuncOpen+ticketMailSettingCheck) || !strings.Contains(ticket, "\t\"{{MODULE}}/internal/settings\"\n") {
		t.Error("the ticket email does not honour notifications.email_enabled")
	}
	mustFormatGo(t, "ticket_mail.go", strings.ReplaceAll(ticket, "{{MODULE}}", "example.com/app"))

	recovery := recoveryHandlerGo()
	if strings.Contains(recovery, "_ = h.Mail.SendRaw") || !strings.Contains(recovery, `"MAIL_FAILED"`) {
		t.Error("the recovery handler still ignores SendRaw's error")
	}
	mustFormatGo(t, "recovery.go", strings.ReplaceAll(recovery, "{{MODULE}}", "example.com/app"))
	mustFormatGo(t, "mail_preview.go", strings.ReplaceAll(mailPreviewHandlerGo(), "{{MODULE}}", "example.com/app"))

	page := adminMailPage()
	for _, want := range []string{"/api/admin/mail/templates", "/api/admin/mail/preview/", "srcDoc", `sandbox=""`} {
		if !strings.Contains(page, want) {
			t.Errorf("the Mail Preview page does not use %s", want)
		}
	}
	if strings.Contains(page, "sampleData") || strings.Contains(page, "Reset Your Password") {
		t.Error("the Mail Preview page still carries JSX copies of the templates")
	}
	for name, src := range map[string]string{"admin health": adminSystemHealthPage(), "desktop health": desktopClientSystemHealthPage()} {
		if strings.Contains(src, "Email (Resend)") || !strings.Contains(src, "data.email.driver") {
			t.Errorf("the %s page labels email as Resend instead of showing the driver", name)
		}
	}
	for rel, src := range map[string]string{"page": page, "queue.go": mailQueueGo(), "preview.go": mailPreviewGo(), "handler": mailPreviewHandlerGo(), "repair": routesMailPreview + ticketMailSettingCheck + jobsWorkerMailNew + jobsEnqueuePayloadFunc} {
		if strings.Contains(src, "—") {
			t.Errorf("%s contains an em dash", rel)
		}
	}
}

func TestMailQueueAndPreviewRepairs(t *testing.T) {
	cases := []struct {
		name  string
		fresh string
		old   string
		fn    func(string) (string, []string, []string)
	}{
		{"jobs client", jobsClientGo(), strings.Replace(jobsClientGo(), jobsEnqueuePayloadFunc, "", 1), repairJobsEnqueuePayloadSource},
		{"worker", jobsWorkersGo(), strings.Replace(jobsWorkersGo(), jobsWorkerMailNew, jobsWorkerMailOld, 1), repairJobsWorkerMailSource},
		{"routes", apiRoutesGo(), strings.Replace(apiRoutesGo(), routesMailPreviewAnchor+routesMailPreview, routesMailPreviewAnchor, 1), repairRoutesMailPreviewSource},
		{"ticket mail", ticketMailGo(), strings.Replace(strings.Replace(ticketMailGo(), ticketMailSettingCheck, "", 1), "\t\"{{MODULE}}/internal/settings\"\n", "", 1), repairTicketMailSettingSource},
	}
	for _, tc := range cases {
		if tc.old == tc.fresh {
			t.Fatalf("%s: could not rebuild the old file", tc.name)
		}
		out, changes, warnings := tc.fn(tc.old)
		if out != tc.fresh || len(changes) != 1 || len(warnings) != 0 {
			t.Errorf("%s: repairing the old file did not give the template (changes %v, warnings %v)", tc.name, changes, warnings)
		}
		if again, changes, warnings := tc.fn(out); again != out || len(changes) != 0 || len(warnings) != 0 {
			t.Errorf("%s: the repair is not idempotent", tc.name)
		}
	}

	// An old client with the method appended at the end still compiles.
	client := strings.Replace(jobsClientGo(), jobsEnqueuePayloadFunc, "", 1)
	client = strings.Replace(client, "\n// EnqueueProcessImage enqueues an image processing job.", "\n// EnqueueProcessImage enqueues an image job.", 1)
	if out, changes, _ := repairJobsEnqueuePayloadSource(client); len(changes) != 1 || !strings.HasSuffix(out, jobsEnqueuePayloadFunc) {
		t.Error("a client without the usual neighbour did not get EnqueuePayload at the end")
	}

	edited := map[string]func(string) (string, []string, []string){
		"package jobs\n\ntype Client struct{}\n":                       repairJobsEnqueuePayloadSource,
		"package jobs\n\nfunc handleEmailSend() {}\n":                  repairJobsWorkerMailSource,
		"package routes\n\nfunc Setup() { staff := v1.Group(\"\") }\n": repairRoutesMailPreviewSource,
		"package services\n\nfunc SendTicketCreatedEmail() {}\n":       repairTicketMailSettingSource,
	}
	for src, fn := range edited {
		if out, _, warnings := fn(src); out != src || len(warnings) != 1 {
			t.Errorf("an edited file got no note or was changed: %q", src)
		}
	}
}
