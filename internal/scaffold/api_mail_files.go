package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeMailFiles(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	module := opts.Module()

	// The package is framework-owned: mailer.go, the drivers, FromConfig, their
	// tests and the mailtest fake. templates.go is written with them.
	files := map[string]string{
		filepath.Join(apiRoot, "internal", "mail", "templates.go"): mailTemplatesGo(),
	}
	for rel, content := range mailPackageFiles() {
		files[filepath.Join(apiRoot, "internal", "mail", filepath.FromSlash(rel))] = content
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "{{MODULE}}", module)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func mailTemplatesGo() string {
	return `package mail

// EmailTemplates contains all available email templates.
var EmailTemplates = map[string]string{
	"welcome":            welcomeTemplate,
	"password-reset":     passwordResetTemplate,
	"email-verification": emailVerificationTemplate,
	"notification":       notificationTemplate,
}

const baseLayout = ` + "`" + `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { margin: 0; padding: 0; background-color: #0a0a0f; color: #e8e8f0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
    .container { max-width: 600px; margin: 0 auto; padding: 40px 20px; }
    .card { background-color: #111118; border: 1px solid #2a2a3a; border-radius: 12px; padding: 32px; }
    .logo { text-align: center; margin-bottom: 24px; font-size: 24px; font-weight: 700; color: #6c5ce7; }
    h1 { font-size: 20px; margin: 0 0 16px; color: #e8e8f0; }
    p { font-size: 14px; line-height: 1.6; color: #9090a8; margin: 0 0 16px; }
    .btn { display: inline-block; background-color: #6c5ce7; color: #ffffff; text-decoration: none; padding: 12px 24px; border-radius: 8px; font-weight: 600; font-size: 14px; }
    .btn:hover { background-color: #7c6cf7; }
    .footer { text-align: center; margin-top: 24px; font-size: 12px; color: #606078; }
    .code { background-color: #1a1a24; border: 1px solid #2a2a3a; border-radius: 8px; padding: 16px; text-align: center; font-size: 28px; letter-spacing: 4px; font-weight: 700; color: #6c5ce7; margin: 16px 0; }
  </style>
</head>
<body>
  <div class="container">
    <div class="card">
      <div class="logo">{{.AppName}}</div>
      {{.Content}}
    </div>
    <div class="footer">
      <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
    </div>
  </div>
</body>
</html>` + "`" + `

const welcomeTemplate = ` + "`" + `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { margin: 0; padding: 0; background-color: #0a0a0f; color: #e8e8f0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
    .container { max-width: 600px; margin: 0 auto; padding: 40px 20px; }
    .card { background-color: #111118; border: 1px solid #2a2a3a; border-radius: 12px; padding: 32px; }
    .logo { text-align: center; margin-bottom: 24px; font-size: 24px; font-weight: 700; color: #6c5ce7; }
    h1 { font-size: 20px; margin: 0 0 16px; color: #e8e8f0; }
    p { font-size: 14px; line-height: 1.6; color: #9090a8; margin: 0 0 16px; }
    .btn { display: inline-block; background-color: #6c5ce7; color: #ffffff; text-decoration: none; padding: 12px 24px; border-radius: 8px; font-weight: 600; font-size: 14px; }
    .footer { text-align: center; margin-top: 24px; font-size: 12px; color: #606078; }
  </style>
</head>
<body>
  <div class="container">
    <div class="card">
      <div class="logo">{{.AppName}}</div>
      <h1>Welcome, {{.Name}}!</h1>
      <p>Thanks for signing up. Your account is ready to use.</p>
      <p>Get started by exploring the dashboard:</p>
      <p style="text-align: center; margin-top: 24px;">
        <a href="{{.DashboardURL}}" class="btn">Go to Dashboard</a>
      </p>
    </div>
    <div class="footer">
      <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
    </div>
  </div>
</body>
</html>` + "`" + `

// #nosec G101 -- an email about a password, not a credential in the source.
const passwordResetTemplate = ` + "`" + `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { margin: 0; padding: 0; background-color: #0a0a0f; color: #e8e8f0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
    .container { max-width: 600px; margin: 0 auto; padding: 40px 20px; }
    .card { background-color: #111118; border: 1px solid #2a2a3a; border-radius: 12px; padding: 32px; }
    .logo { text-align: center; margin-bottom: 24px; font-size: 24px; font-weight: 700; color: #6c5ce7; }
    h1 { font-size: 20px; margin: 0 0 16px; color: #e8e8f0; }
    p { font-size: 14px; line-height: 1.6; color: #9090a8; margin: 0 0 16px; }
    .btn { display: inline-block; background-color: #6c5ce7; color: #ffffff; text-decoration: none; padding: 12px 24px; border-radius: 8px; font-weight: 600; font-size: 14px; }
    .footer { text-align: center; margin-top: 24px; font-size: 12px; color: #606078; }
  </style>
</head>
<body>
  <div class="container">
    <div class="card">
      <div class="logo">{{.AppName}}</div>
      <h1>Reset Your Password</h1>
      <p>We received a request to reset your password. Click the button below to set a new one:</p>
      <p style="text-align: center; margin-top: 24px;">
        <a href="{{.ResetURL}}" class="btn">Reset Password</a>
      </p>
      <p>This link expires in 1 hour. If you didn't request this, you can safely ignore this email.</p>
    </div>
    <div class="footer">
      <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
    </div>
  </div>
</body>
</html>` + "`" + `

const emailVerificationTemplate = ` + "`" + `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { margin: 0; padding: 0; background-color: #0a0a0f; color: #e8e8f0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
    .container { max-width: 600px; margin: 0 auto; padding: 40px 20px; }
    .card { background-color: #111118; border: 1px solid #2a2a3a; border-radius: 12px; padding: 32px; }
    .logo { text-align: center; margin-bottom: 24px; font-size: 24px; font-weight: 700; color: #6c5ce7; }
    h1 { font-size: 20px; margin: 0 0 16px; color: #e8e8f0; }
    p { font-size: 14px; line-height: 1.6; color: #9090a8; margin: 0 0 16px; }
    .btn { display: inline-block; background-color: #6c5ce7; color: #ffffff; text-decoration: none; padding: 12px 24px; border-radius: 8px; font-weight: 600; font-size: 14px; }
    .footer { text-align: center; margin-top: 24px; font-size: 12px; color: #606078; }
  </style>
</head>
<body>
  <div class="container">
    <div class="card">
      <div class="logo">{{.AppName}}</div>
      <h1>Verify Your Email</h1>
      <p>Please verify your email address by clicking the button below:</p>
      <p style="text-align: center; margin-top: 24px;">
        <a href="{{.VerifyURL}}" class="btn">Verify Email</a>
      </p>
      <p>If you didn't create an account, you can safely ignore this email.</p>
    </div>
    <div class="footer">
      <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
    </div>
  </div>
</body>
</html>` + "`" + `

const notificationTemplate = ` + "`" + `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { margin: 0; padding: 0; background-color: #0a0a0f; color: #e8e8f0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
    .container { max-width: 600px; margin: 0 auto; padding: 40px 20px; }
    .card { background-color: #111118; border: 1px solid #2a2a3a; border-radius: 12px; padding: 32px; }
    .logo { text-align: center; margin-bottom: 24px; font-size: 24px; font-weight: 700; color: #6c5ce7; }
    h1 { font-size: 20px; margin: 0 0 16px; color: #e8e8f0; }
    p { font-size: 14px; line-height: 1.6; color: #9090a8; margin: 0 0 16px; }
    .btn { display: inline-block; background-color: #6c5ce7; color: #ffffff; text-decoration: none; padding: 12px 24px; border-radius: 8px; font-weight: 600; font-size: 14px; }
    .footer { text-align: center; margin-top: 24px; font-size: 12px; color: #606078; }
  </style>
</head>
<body>
  <div class="container">
    <div class="card">
      <div class="logo">{{.AppName}}</div>
      <h1>{{.Title}}</h1>
      <p>{{.Message}}</p>
      {{if .ActionURL}}
      <p style="text-align: center; margin-top: 24px;">
        <a href="{{.ActionURL}}" class="btn">{{.ActionText}}</a>
      </p>
      {{end}}
    </div>
    <div class="footer">
      <p>&copy; {{.Year}} {{.AppName}}. All rights reserved.</p>
    </div>
  </div>
</body>
</html>` + "`" + `
`
}
