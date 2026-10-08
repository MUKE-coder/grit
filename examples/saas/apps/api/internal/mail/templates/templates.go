// Package templates holds the application's own emails. grit generate mail
// writes one file here per email: a typed data struct, an html/template body
// inside the shared layout, a text alternative, and Send and Queue helpers.
//
// Each file registers its email with mail.Register, so the admin's Mail
// Preview lists it. internal/handlers/mail_preview.go imports this package for
// that reason.
package templates
