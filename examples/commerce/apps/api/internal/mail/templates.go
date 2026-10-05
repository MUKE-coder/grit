package mail

// The six emails, as fragments. Layout is the shell they are rendered into.
//
// Fragments rather than whole documents, because six copies of one shell is how
// they drifted apart last time: three different footers and two different greys
// between them.
var EmailTemplates = map[string]string{
	"welcome":            welcomeTemplate,
	"password-reset":     passwordResetTemplate,
	"email-verification": emailVerificationTemplate,
	"notification":       notificationTemplate,
	"two-factor-code":    twoFactorCodeTemplate,
	"magic-link":         magicLinkTemplate,
}

// Layout wraps every fragment.
//
// Tables and style attributes, which is not a stylistic choice. Outlook on
// Windows renders mail through Word, which ignores most of a <style> block and
// all of flexbox; Gmail drops <style> when it clips a long message or when
// somebody forwards it. Anything that has to survive goes on the element.
//
// The style strings the fragments use ({{.H1}}, {{.P}} and the rest) are
// supplied by the renderer, so there is one definition of what a heading looks
// like rather than one per template.
const Layout = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
<meta name="viewport" content="width=device-width, initial-scale=1" />
<meta name="x-apple-disable-message-reformatting" />
<meta name="color-scheme" content="light" />
<meta name="supported-color-schemes" content="light" />
<title>{{.Subject}}</title>
</head>
<body style="margin:0;padding:0;background-color:{{.Canvas}};">
<div style="display:none;max-height:0;overflow:hidden;mso-hide:all;font-size:1px;line-height:1px;color:{{.Canvas}};">{{.Preheader}}&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;&#8203;</div>
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="background-color:{{.Canvas}};">
<tr>
<td align="center" style="padding:32px 12px;">
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="max-width:560px;width:100%;">
<tr>
<td align="left" style="padding:0 4px 18px;font-family:{{.Font}};font-size:17px;font-weight:700;letter-spacing:-0.01em;color:{{.Accent}};">{{.AppName}}</td>
</tr>
<tr>
<td style="background-color:{{.Card}};border:1px solid {{.Border}};border-radius:14px;padding:32px;">
{{.Content}}
</td>
</tr>
<tr>
<td align="left" style="padding:18px 4px 0;font-family:{{.Font}};font-size:12px;line-height:1.6;color:{{.Muted}};">
&copy; {{.Year}} {{.AppName}}. Sent to you because this address was used on {{.AppName}}.
</td>
</tr>
</table>
</td>
</tr>
</table>
</body>
</html>`

// Fragments: the inside of the card, and nothing else.
//
// Divs with explicit margins rather than h1/p, because Outlook adds margins of
// its own to block elements that no stylesheet can take back.

const welcomeTemplate = `<div style="{{.H1}}">Welcome to {{.AppName}}</div>
<div style="{{.P}}">Hi {{.Name}}, your account is ready. There is nothing else to set up.</div>
{{if .ActionURL}}<table role="presentation" cellpadding="0" cellspacing="0" border="0" style="margin:0 0 20px;"><tr><td align="center" style="{{.BtnCell}}"><a href="{{.ActionURL}}" style="{{.BtnLink}}">{{.ActionText}}</a></td></tr></table>{{end}}
<div style="{{.PLast}}">If you did not create this account you can ignore this message, and nothing further will happen.</div>`

const passwordResetTemplate = `<div style="{{.H1}}">Reset your password</div>
<div style="{{.P}}">Somebody asked to reset the password for this address. Use the button below to choose a new one.</div>
<table role="presentation" cellpadding="0" cellspacing="0" border="0" style="margin:0 0 20px;"><tr><td align="center" style="{{.BtnCell}}"><a href="{{.ResetURL}}" style="{{.BtnLink}}">Choose a new password</a></td></tr></table>
<div style="{{.PLast}}">The link works once and expires in an hour. If this was not you, nothing has changed and you can ignore this message.</div>`

const emailVerificationTemplate = `<div style="{{.H1}}">{{.Title}}</div>
<div style="{{.P}}">{{.Message}}</div>
<table role="presentation" cellpadding="0" cellspacing="0" border="0" style="margin:0 0 20px;"><tr><td align="center" style="{{.BtnCell}}"><a href="{{.VerifyURL}}" style="{{.BtnLink}}">Confirm this address</a></td></tr></table>
<div style="{{.PLast}}">If the button does not work, paste this into your browser:<br /><span style="word-break:break-all;color:{{.Accent}};">{{.VerifyURL}}</span></div>`

const notificationTemplate = `<div style="{{.H1}}">{{.Title}}</div>
<div style="{{.P}}">{{.Message}}</div>
{{if .ActionURL}}<table role="presentation" cellpadding="0" cellspacing="0" border="0" style="margin:0;"><tr><td align="center" style="{{.BtnCell}}"><a href="{{.ActionURL}}" style="{{.BtnLink}}">{{.ActionText}}</a></td></tr></table>{{end}}`

const magicLinkTemplate = `<div style="{{.H1}}">Your sign-in link</div>
<div style="{{.P}}">Use the button below to sign in. There is no password to type.</div>
<table role="presentation" cellpadding="0" cellspacing="0" border="0" style="margin:0 0 20px;"><tr><td align="center" style="{{.BtnCell}}"><a href="{{.MagicURL}}" style="{{.BtnLink}}">Sign in</a></td></tr></table>
<div style="{{.PLast}}">The link lasts {{.Minutes}} minutes and works once. If you did not ask for it, somebody knows your email address, which is not a secret: they cannot sign in without this mailbox.</div>`

const twoFactorCodeTemplate = `<div style="{{.H1}}">Your sign-in code</div>
<div style="{{.P}}">Type this into the page you already have open.</div>
<table role="presentation" width="100%" cellpadding="0" cellspacing="0" border="0" style="margin:0 0 18px;">
<tr><td align="center" style="background-color:{{.Canvas}};border:1px solid {{.Border}};border-radius:12px;padding:20px 12px;font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;font-size:32px;font-weight:700;letter-spacing:10px;color:{{.Ink}};">{{.Code}}</td></tr>
</table>
<div style="{{.PLast}}">It expires in {{.Minutes}} minutes and works once. If you were not signing in, somebody has your password: change it. This code alone lets nobody in.</div>`
