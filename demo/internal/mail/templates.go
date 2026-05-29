package mail

// EmailTemplates contains all available email templates.
var EmailTemplates = map[string]string{
	"welcome":            welcomeTemplate,
	"password-reset":     passwordResetTemplate,
	"email-verification": emailVerificationTemplate,
	"notification":       notificationTemplate,
	"invitation":         invitationTemplate,
	"low-stock":          lowStockTemplate,
}

const lowStockTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { margin: 0; padding: 0; background-color: #F8FAFC; color: #0F172A; font-family: 'Inter', -apple-system, sans-serif; }
    .container { max-width: 600px; margin: 0 auto; padding: 40px 20px; }
    .card { background-color: #FFFFFF; border: 1px solid #E2E8F0; border-radius: 16px; padding: 32px; }
    .logo { text-align: center; margin-bottom: 24px; }
    .logo-icon { display: inline-block; width: 48px; height: 48px; line-height: 48px; border-radius: 12px; background-color: #16A34A; color: #fff; font-size: 24px; font-weight: 700; text-align: center; }
    h1 { font-size: 20px; margin: 0 0 8px; color: #0F172A; }
    p { font-size: 14px; line-height: 1.6; color: #64748B; margin: 0 0 16px; }
    .alert-box { background: #FEF2F2; border: 1px solid #FECACA; border-radius: 12px; padding: 16px; text-align: center; margin: 16px 0; }
    .alert-stock { font-size: 32px; font-weight: 800; color: #EF4444; }
    .alert-label { font-size: 12px; color: #EF4444; font-weight: 600; margin-top: 4px; }
    .info-row { display: flex; justify-content: space-between; padding: 8px 0; border-bottom: 1px solid #F1F5F9; font-size: 14px; }
    .info-label { color: #64748B; }
    .info-value { color: #0F172A; font-weight: 600; }
    .footer { text-align: center; margin-top: 24px; font-size: 12px; color: #94A3B8; }
  </style>
</head>
<body>
  <div class="container">
    <div class="card">
      <div class="logo"><span class="logo-icon" style="font-size: 14px;">KM</span></div>
      <h1>Low Stock Alert</h1>
      <p>A product in your inventory has dropped below the low stock threshold and needs attention.</p>
      <div class="alert-box">
        <div class="alert-stock">{{.CurrentStock}} units</div>
        <div class="alert-label">remaining (threshold: {{.Threshold}})</div>
      </div>
      <div style="margin-top: 16px;">
        <div class="info-row"><span class="info-label">Product</span><span class="info-value">{{.ProductTitle}}</span></div>
        <div class="info-row"><span class="info-label">Branch</span><span class="info-value">{{.BranchName}}</span></div>
      </div>
      <p style="margin-top: 20px; font-size: 13px;">Log in to Grit Motors to restock this product or transfer stock from another branch.</p>
    </div>
    <div class="footer">
      <p>Grit Motors — Motorcycles. Spares. Loans.</p>
    </div>
  </div>
</body>
</html>`

const invitationTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    body { margin: 0; padding: 0; background-color: #F8FAFC; color: #0F172A; font-family: 'Inter', -apple-system, sans-serif; }
    .container { max-width: 600px; margin: 0 auto; padding: 40px 20px; }
    .card { background-color: #FFFFFF; border: 1px solid #E2E8F0; border-radius: 16px; padding: 32px; }
    .logo { text-align: center; margin-bottom: 24px; }
    .logo-icon { display: inline-block; width: 48px; height: 48px; line-height: 48px; border-radius: 12px; background-color: #16A34A; color: #fff; font-size: 24px; font-weight: 700; text-align: center; }
    h1 { font-size: 20px; margin: 0 0 8px; color: #0F172A; }
    p { font-size: 14px; line-height: 1.6; color: #64748B; margin: 0 0 16px; }
    .role-badge { display: inline-block; padding: 4px 12px; border-radius: 20px; font-size: 13px; font-weight: 600; background: #F0FDF4; color: #16A34A; }
    .btn { display: inline-block; background-color: #16A34A; color: #ffffff; text-decoration: none; padding: 12px 32px; border-radius: 10px; font-weight: 600; font-size: 14px; }
    .footer { text-align: center; margin-top: 24px; font-size: 12px; color: #94A3B8; }
  </style>
</head>
<body>
  <div class="container">
    <div class="card">
      <div class="logo"><span class="logo-icon" style="font-size: 14px;">KM</span></div>
      <h1>You're invited to join {{.BusinessName}}</h1>
      <p>You've been invited to join <strong>{{.BusinessName}}</strong> on Grit Motors as a <span class="role-badge">{{.Role}}</span>.</p>
      <p>Click the button below to accept the invitation and get started:</p>
      <p style="text-align: center; margin-top: 24px;">
        <a href="{{.InviteURL}}" class="btn">Accept Invitation</a>
      </p>
      <p style="font-size: 12px; color: #94A3B8; margin-top: 24px;">This invitation expires in 7 days. If you didn't expect this, you can safely ignore it.</p>
    </div>
    <div class="footer">
      <p>Grit Motors — Motorcycles. Spares. Loans.</p>
    </div>
  </div>
</body>
</html>`

const baseLayout = `<!DOCTYPE html>
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
</html>`

const welcomeTemplate = `<!DOCTYPE html>
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
</html>`

const passwordResetTemplate = `<!DOCTYPE html>
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
</html>`

const emailVerificationTemplate = `<!DOCTYPE html>
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
</html>`

const notificationTemplate = `<!DOCTYPE html>
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
</html>`
