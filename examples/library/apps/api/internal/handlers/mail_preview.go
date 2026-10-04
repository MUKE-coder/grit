package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"library/apps/api/internal/mail"
	// For its init functions: every email grit generate mail writes registers
	// itself for this preview.
	_ "library/apps/api/internal/mail/templates"
	"library/apps/api/internal/respond"
)

// MailPreviewHandler serves the admin's Mail Preview: the registered email
// templates, rendered with sample data by the same code that sends them.
type MailPreviewHandler struct {
	Mailer *mail.Mailer
}

// MailTemplateView is one template in the preview's list.
type MailTemplateView struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Subject     string `json:"subject"`
	HasText     bool   `json:"has_text"`
}

// Templates handles GET /api/admin/mail/templates. driver names the transport
// mail goes out through, empty when none is configured.
func (h *MailPreviewHandler) Templates(c *gin.Context) {
	registered := mail.Templates()
	out := make([]MailTemplateView, 0, len(registered))
	for _, t := range registered {
		view := MailTemplateView{Name: t.Name, Description: t.Description}
		if msg, err := t.Render(); err == nil {
			view.Subject = msg.Subject
			view.HasText = msg.Text != ""
		}
		out = append(out, view)
	}
	driver := ""
	if h.Mailer != nil {
		driver = h.Mailer.Driver()
	}
	c.JSON(http.StatusOK, gin.H{"data": out, "driver": driver})
}

// Preview handles GET /api/admin/mail/preview/:template: the rendered HTML, or
// the text part with ?part=text.
func (h *MailPreviewHandler) Preview(c *gin.Context) {
	t, ok := mail.LookupTemplate(c.Param("template"))
	if !ok {
		respond.Fail(c, respond.CodeNotFound, "No mail template by that name")
		return
	}
	msg, err := t.Render()
	if err != nil {
		log.Printf("mail preview: rendering %s: %v", t.Name, err)
		respond.Fail(c, respond.CodeInternalError, "The template could not be rendered")
		return
	}
	// The admin shows this in a sandboxed iframe. Opened on its own, the page
	// still runs no script and loads nothing but images.
	c.Header("Content-Security-Policy", "default-src 'none'; img-src data: https:; style-src 'unsafe-inline'; sandbox")
	c.Header("Cache-Control", "no-store")
	if c.Query("part") == "text" {
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(msg.Text))
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(msg.HTML))
}
