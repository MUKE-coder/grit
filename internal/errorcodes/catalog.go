// Package errorcodes is the catalogue of every error a Grit API can return.
//
// A client needs two things from an error: a code it can branch on, and the
// knowledge that the code always arrives with the same status. Grit had the
// first and not the second. VALIDATION_ERROR came back as 422 in thirty-eight
// places and 400 in twenty-five, INVALID_TOKEN as both 401 and 400, and nothing
// anywhere listed what the codes were, so a frontend had to discover each one by
// triggering it.
//
// This is the list, and it is the only list. Three things are generated from it,
// so they cannot disagree:
//
//   - internal/respond/codes.go in every project, with a typed constant and the
//     status for each code;
//   - packages/shared/types/errors.ts, so a TypeScript switch over codes can be
//     exhaustive;
//   - the table on /docs/backend/errors.
//
// And a test walks the templates: a code emitted anywhere that is not in here,
// or that is emitted with a status other than the one here, fails the build.
// That test is the reason the statuses above are consistent now.
package errorcodes

import (
	"net/http"
	"sort"
)

// Category is what kind of problem a code reports, which is what decides how a
// client should treat it: fix the request, sign in, wait, or give up and report.
type Category string

const (
	// CategoryRequest: the request itself is wrong. Fixing it is the client's job.
	CategoryRequest Category = "request"
	// CategoryAuth: nobody is signed in, or the credentials are not accepted.
	CategoryAuth Category = "auth"
	// CategoryPermission: signed in, but not allowed to do this.
	CategoryPermission Category = "permission"
	// CategoryNotFound: no such thing, or nothing this caller may see.
	CategoryNotFound Category = "notfound"
	// CategoryConflict: the state of the server says no, for now.
	CategoryConflict Category = "conflict"
	// CategoryState: the thing exists but is not in a state that allows this.
	CategoryState Category = "state"
	// CategoryLimit: too much, too fast.
	CategoryLimit Category = "limit"
	// CategoryServer: a fault on our side. The client did nothing wrong.
	CategoryServer Category = "server"
	// CategoryUpstream: a service the API depends on failed or refused.
	CategoryUpstream Category = "upstream"
	// CategoryDisabled: the feature is not configured or not running here.
	CategoryDisabled Category = "disabled"
)

// Entry is one code, and everything a caller needs to know about it.
type Entry struct {
	Code     string
	Status   int
	Category Category
	// Area is the part of the API that raises it, for grouping in the docs.
	Area string
	// Meaning is what happened, in one sentence, from the caller's side.
	Meaning string
	// Client is what the caller should do about it. This is the field a frontend
	// developer actually wants, and the reason a bare status is not enough.
	Client string
}

// catalog is every code a generated project can emit. Keep it sorted by area,
// then code, and keep one status per code: the second status is the bug.
var catalog = []Entry{
	// ── The envelope every handler shares ──────────────────────────────────────
	{"BAD_REQUEST", http.StatusBadRequest, CategoryRequest, "core",
		"The request could not be understood at all.",
		"Fix the request. Retrying the same one will fail the same way."},
	{"INVALID_BODY", http.StatusBadRequest, CategoryRequest, "core",
		"The body was not valid JSON, or was not the shape this endpoint reads.",
		"Send a JSON body matching the documented request type."},
	{"READ_BODY_FAILED", http.StatusBadRequest, CategoryRequest, "core",
		"The body could not be read to the end.",
		"Retry. If it keeps happening, the connection is dropping or the body is larger than the server accepts."},
	{"VALIDATION_ERROR", http.StatusUnprocessableEntity, CategoryRequest, "core",
		"The body parsed, and a field in it is missing or not acceptable.",
		"Read error.details: it maps each field to what is wrong with it. Show those against the inputs."},
	{"PAYLOAD_TOO_LARGE", http.StatusRequestEntityTooLarge, CategoryRequest, "core",
		"The request body is larger than the server accepts.",
		"Send less, or upload the file directly to storage with a presigned URL."},
	{"UNAUTHORIZED", http.StatusUnauthorized, CategoryAuth, "core",
		"The credentials are missing, expired or not accepted.",
		"Refresh the access token, and sign in again if the refresh is rejected too."},
	{"MISSING_TOKEN", http.StatusUnauthorized, CategoryAuth, "core",
		"No token was sent where one is required.",
		"Send the access token as Authorization: Bearer <token>."},
	{"INVALID_TOKEN", http.StatusUnauthorized, CategoryAuth, "core",
		"The token is malformed, expired, or was issued for something else.",
		"Refresh it. A refresh token that fails this way has been used already or revoked: sign in again."},
	{"INVALID_LINK", http.StatusBadRequest, CategoryRequest, "core",
		"A one-time link, such as email verification or password reset, is wrong or has expired.",
		"Offer to send a new link. This is not a sign-in failure: the caller has no credentials to fix, which is why it is 400 and not 401."},
	{"SESSION_REVOKED", http.StatusUnauthorized, CategoryAuth, "core",
		"The session behind this token was signed out, on this device or another.",
		"Sign in again. Do not retry with the same refresh token: reuse is what revoked it."},
	{"CSRF_INVALID", http.StatusForbidden, CategoryPermission, "core",
		"The CSRF token is missing or does not match the cookie.",
		"Read the CSRF cookie and send it back in the header on every unsafe request."},
	{"FORBIDDEN", http.StatusForbidden, CategoryPermission, "core",
		"The caller is signed in and this action is not theirs to take.",
		"Do not retry. Hide the action rather than letting it fail, if the role is known to the client."},
	{"NOT_FOUND", http.StatusNotFound, CategoryNotFound, "core",
		"No such row, or none this caller is allowed to see.",
		"Treat it as absent. On an owned resource this is also the answer for somebody else's row, on purpose: a wrong guess cannot be told from a right one."},
	{"CONFLICT", http.StatusConflict, CategoryConflict, "core",
		"The write collided with the state already there, such as a unique column.",
		"Re-read, show what is there, and let the person decide."},
	{"VERSION_CONFLICT", http.StatusConflict, CategoryConflict, "core",
		"Somebody else changed the row since the version in your If-Match.",
		"The response carries the current version. Re-read, merge, and send the new ETag."},
	{"RATE_LIMITED", http.StatusTooManyRequests, CategoryLimit, "core",
		"Too many requests from this caller.",
		"Back off. Honour Retry-After if it is present rather than retrying immediately."},
	{"INTERNAL_ERROR", http.StatusInternalServerError, CategoryServer, "core",
		"A fault on the server. The message is deliberately vague; the detail is in the server log.",
		"Retry once, then report it. Nothing the client changes will help."},
	{"MAINTENANCE", http.StatusServiceUnavailable, CategoryDisabled, "core",
		"The API is in maintenance mode and is refusing everything.",
		"Retry later. Show a maintenance state rather than an error."},
	{"PERSIST_FAILED", http.StatusInternalServerError, CategoryServer, "core",
		"The change was accepted and could not be written.",
		"Retry once. Treat the write as not having happened."},

	// ── Organizations (the multitenant plugin) ─────────────────────────────────
	{"NO_ORGANIZATION", http.StatusBadRequest, CategoryRequest, "tenancy",
		"The row belongs to an organization and the request has no active one: the caller belongs to none, or to several and named neither.",
		"Send the active organization as X-Organization-ID. If the caller belongs to no organization, they cannot read this at all: put them in one, or send them somewhere that does not need one."},

	// ── Sign-in and accounts ───────────────────────────────────────────────────
	{"INVALID_CREDENTIALS", http.StatusUnauthorized, CategoryAuth, "auth",
		"The email and password do not match an account.",
		"Say only that the details are wrong: which of the two it was is deliberately not reported."},
	{"INVALID_PASSWORD", http.StatusUnauthorized, CategoryAuth, "auth",
		"The password given for a confirmation step is not correct.",
		"Ask again. This is the re-authentication prompt, not a sign-in."},
	{"EMAIL_EXISTS", http.StatusConflict, CategoryConflict, "auth",
		"An account already has that address.",
		"Offer sign-in or password reset rather than registration."},
	{"EMAIL_NOT_VERIFIED", http.StatusForbidden, CategoryPermission, "auth",
		"The account exists and its address has not been confirmed.",
		"Send them to the verification flow, and offer to resend the link."},
	{"ACCOUNT_DISABLED", http.StatusForbidden, CategoryPermission, "auth",
		"The account has been deactivated.",
		"Do not retry. This needs an administrator, not a different password."},
	{"ACCOUNT_LOCKED", http.StatusTooManyRequests, CategoryLimit, "auth",
		"Too many failed attempts, so the account is locked for a while.",
		"Show the wait, and offer password reset. Retrying sooner extends nothing but the lock."},
	{"ALREADY_VERIFIED", http.StatusBadRequest, CategoryRequest, "auth",
		"The address on this link is already confirmed.",
		"Treat it as success and continue to sign-in."},
	{"SOCIAL_AUTH_ONLY", http.StatusBadRequest, CategoryRequest, "auth",
		"The account signs in with a social provider and has no password.",
		"Offer the provider button instead of the password form."},
	{"NO_PASSWORD", http.StatusBadRequest, CategoryRequest, "auth",
		"The account has no password set, so it cannot be confirmed with one.",
		"Send them through set-a-password first."},
	{"UNKNOWN_PROVIDER", http.StatusNotFound, CategoryNotFound, "auth",
		"No social provider is configured under that name.",
		"Only offer the providers the API reports as enabled."},
	{"TOKEN_ERROR", http.StatusInternalServerError, CategoryServer, "auth",
		"The access and refresh tokens could not be issued.",
		"Retry once. The credentials were accepted, so do not ask for them again."},
	{"USER_ERROR", http.StatusInternalServerError, CategoryServer, "auth",
		"The signed-in account could not be loaded, which means the token outlived its row.",
		"Sign in again. If it repeats, the account data is inconsistent and needs a look."},
	{"INVALID_SIGNATURE", http.StatusUnauthorized, CategoryAuth, "webhooks",
		"The webhook signature does not match the body and the shared secret.",
		"Sign the exact bytes sent, with the secret for this endpoint. A reformatted body will not verify."},

	// ── Two-factor authentication ──────────────────────────────────────────────
	{"INVALID_TOTP_CODE", http.StatusUnauthorized, CategoryAuth, "twofactor",
		"The six-digit code is wrong or has expired.",
		"Let them try the next code. Clock drift on the device is the usual cause of repeated failures."},
	{"INVALID_BACKUP_CODE", http.StatusUnauthorized, CategoryAuth, "twofactor",
		"That backup code is wrong, or has been used.",
		"Each code works once. Offer the remaining count the status endpoint reports."},
	{"INVALID_PENDING_TOKEN", http.StatusUnauthorized, CategoryAuth, "twofactor",
		"The short-lived token between password and second factor is expired or unknown.",
		"Start the sign-in again from the password step."},
	{"TOTP_ALREADY_ENABLED", http.StatusConflict, CategoryConflict, "twofactor",
		"Two-factor authentication is already on for this account.",
		"Show it as enabled rather than offering setup."},
	{"TOTP_NOT_ENABLED", http.StatusBadRequest, CategoryRequest, "twofactor",
		"The account has no second factor, so there is nothing to confirm or turn off.",
		"Offer setup instead."},
	{"TOTP_ERROR", http.StatusInternalServerError, CategoryServer, "twofactor",
		"The secret, QR code or backup codes could not be produced.",
		"Retry once, then report it. Nothing is half-enabled: setup only counts once confirmed."},

	// ── Passkeys ───────────────────────────────────────────────────────────────
	{"PASSKEYS_NOT_CONFIGURED", http.StatusNotImplemented, CategoryDisabled, "passkeys",
		"This deployment has no passkey configuration, so the endpoints are inert.",
		"Hide passkey buttons unless the API reports the feature as available."},
	{"PASSKEY_REJECTED", http.StatusGone, CategoryState, "passkeys",
		"The challenge no longer exists: it expired, or it was already answered.",
		"Start the ceremony again. Do not retry the same assertion."},

	// ── Account recovery ───────────────────────────────────────────────────────
	{"INVALID_RECOVERY_ADDRESS", http.StatusUnprocessableEntity, CategoryRequest, "recovery",
		"The recovery email or phone number is not usable.",
		"Validate the format before sending, and show which one was rejected."},
	{"INVALID_CODE", http.StatusUnprocessableEntity, CategoryRequest, "recovery",
		"The recovery code is wrong or has expired.",
		"Offer to send a new one rather than retrying the same code."},
	{"SMS_NOT_CONFIGURED", http.StatusNotImplemented, CategoryDisabled, "recovery",
		"No SMS provider is configured, so phone recovery cannot run here.",
		"Offer email recovery instead."},
	{"SMS_FAILED", http.StatusBadGateway, CategoryUpstream, "recovery",
		"The SMS provider refused or failed to send.",
		"Retry once, then offer email. The number may be unreachable."},

	// ── API keys ───────────────────────────────────────────────────────────────
	{"API_KEY_REQUIRED", http.StatusUnauthorized, CategoryAuth, "apikeys",
		"The route is key-guarded and no key was sent.",
		"Send the key in the documented header. A user token is not a substitute here."},
	{"INVALID_API_KEY", http.StatusUnauthorized, CategoryAuth, "apikeys",
		"The key is unknown, revoked or expired.",
		"Issue a new key. Do not retry with the same one."},
	{"ENDPOINT_NOT_ALLOWED", http.StatusForbidden, CategoryPermission, "apikeys",
		"The key is valid and is not allowed to call this endpoint.",
		"Widen the key's endpoint list, or call it with one that may."},
	{"ORIGIN_NOT_ALLOWED", http.StatusForbidden, CategoryPermission, "apikeys",
		"The request's Origin is not on the key's allowlist.",
		"Add the origin to the key, rather than relaxing CORS for everybody."},
	{"PUBLISHABLE_KEY_NOT_ALLOWED", http.StatusForbidden, CategoryPermission, "apikeys",
		"A publishable key was used where only a secret key is accepted.",
		"Call this from the server with the secret key. A publishable key is public by design."},

	// ── Uploads and storage ────────────────────────────────────────────────────
	{"INVALID_FILE", http.StatusBadRequest, CategoryRequest, "uploads",
		"No file was attached, or it could not be read.",
		"Send multipart form data with the documented field name."},
	{"INVALID_FILE_TYPE", http.StatusBadRequest, CategoryRequest, "uploads",
		"The file's type is not accepted for this field.",
		"Check the type client-side before uploading, and say which types are allowed."},
	{"FILE_TOO_LARGE", http.StatusBadRequest, CategoryRequest, "uploads",
		"The file is larger than this endpoint accepts.",
		"Show the limit before the upload starts rather than after it finishes."},
	{"UPLOAD_FAILED", http.StatusInternalServerError, CategoryServer, "uploads",
		"The file reached the API and could not be stored.",
		"Retry once. Nothing was recorded, so there is no half-uploaded row to clean up."},
	{"UPLOAD_NOT_FOUND", http.StatusNotFound, CategoryNotFound, "uploads",
		"Nothing is stored under that key.",
		"Treat it as absent: the key is wrong, or the object was deleted."},
	{"PRESIGN_FAILED", http.StatusInternalServerError, CategoryServer, "uploads",
		"A presigned upload URL could not be produced.",
		"Retry once, then fall back to uploading through the API."},
	{"STORAGE_UNAVAILABLE", http.StatusServiceUnavailable, CategoryDisabled, "uploads",
		"Object storage is not configured here, or is not answering.",
		"Retry later. Nothing the client sends will fix it."},

	// ── AI gateway ─────────────────────────────────────────────────────────────
	{"AI_UNAVAILABLE", http.StatusServiceUnavailable, CategoryDisabled, "ai",
		"No AI provider is configured in this deployment.",
		"Hide AI features unless the API reports one as available."},
	{"AI_UNAUTHORIZED", http.StatusBadGateway, CategoryUpstream, "ai",
		"The provider rejected the server's API key.",
		"Nothing for the client to do. The key on the server is wrong or out of credit."},
	{"AI_FORBIDDEN", http.StatusBadGateway, CategoryUpstream, "ai",
		"The provider refused this request, usually its own policy.",
		"Do not retry the same prompt unchanged."},
	{"AI_RATE_LIMITED", http.StatusTooManyRequests, CategoryLimit, "ai",
		"The provider is rate-limiting this deployment.",
		"Back off and retry with a delay. Queue rather than loop."},
	{"AI_MODEL_NOT_FOUND", http.StatusBadGateway, CategoryUpstream, "ai",
		"The provider does not know the model that was asked for.",
		"Pick a model the API lists. A model name can disappear without notice."},
	{"AI_ERROR", http.StatusBadGateway, CategoryUpstream, "ai",
		"The provider failed in a way that is not one of the above.",
		"Retry once. The message carries what the provider said."},

	// ── Background jobs ────────────────────────────────────────────────────────
	{"REDIS_UNAVAILABLE", http.StatusServiceUnavailable, CategoryDisabled, "jobs",
		"Redis is not configured or not reachable, so the queue cannot be read.",
		"Retry later. Jobs, cache and cron all depend on it."},
	{"INVALID_STATUS", http.StatusBadRequest, CategoryRequest, "jobs",
		"That queue state is not one the endpoint accepts.",
		"Use one of the documented states."},
	{"RETRY_FAILED", http.StatusInternalServerError, CategoryServer, "jobs",
		"The job could not be re-queued.",
		"Retry once. The job is still where it was."},
	{"CLEAR_FAILED", http.StatusInternalServerError, CategoryServer, "jobs",
		"The queue could not be cleared.",
		"Retry once, then look at the Redis connection."},
	{"JOB_ERROR", http.StatusInternalServerError, CategoryServer, "import",
		"The import job could not be started.",
		"Retry once. Nothing was imported."},
	{"TEMP_ERROR", http.StatusInternalServerError, CategoryServer, "import",
		"The upload could not be buffered to disk before importing.",
		"Retry once. Check free disk on the server if it repeats."},
	{"INVALID_CSV", http.StatusBadRequest, CategoryRequest, "import",
		"The file is not readable as CSV.",
		"Check the delimiter, the quoting and that the header row matches the template."},

	// ── Backups ────────────────────────────────────────────────────────────────
	{"NOT_AVAILABLE", http.StatusBadRequest, CategoryRequest, "backups",
		"That backup cannot be downloaded: it is not finished, or it is not stored here.",
		"Re-read the backup's status before offering a download link."},
	{"INVALID_SCHEDULE", http.StatusBadRequest, CategoryRequest, "backups",
		"The schedule is not a cron expression this API accepts.",
		"Validate the expression client-side, or offer fixed choices."},
	{"EXTRACT_FAILED", http.StatusBadRequest, CategoryRequest, "backups",
		"The archive could not be opened or does not hold what a restore needs.",
		"Upload an archive this API produced. A re-zipped one usually fails here."},

	// ── GDPR ───────────────────────────────────────────────────────────────────
	{"EXPORT_FAILED", http.StatusInternalServerError, CategoryServer, "gdpr",
		"The subject-access export could not be assembled.",
		"Retry once, then report it: this is a request with a legal clock on it."},
	{"ERASE_FAILED", http.StatusInternalServerError, CategoryServer, "gdpr",
		"The erasure did not complete.",
		"Do not assume anything was erased. Retry, and check the audit log."},
	{"SELF_ERASE", http.StatusBadRequest, CategoryRequest, "gdpr",
		"An account cannot erase itself through this endpoint.",
		"Have another administrator run it, so the action has an actor who remains."},
	{"QUERY_FAILED", http.StatusInternalServerError, CategoryServer, "gdpr",
		"The audit query failed.",
		"Retry once, then report it."},
	{"VERIFY_FAILED", http.StatusInternalServerError, CategoryServer, "gdpr",
		"The audit chain could not be verified.",
		"Report it. This is the check that says whether the log has been tampered with."},
	{"RESEAL_REFUSED", http.StatusConflict, CategoryConflict, "gdpr",
		"The audit chain was not broken where the reseal claimed, so nothing was resealed.",
		"Verify the chain again and reseal from the entry the verification names."},

	// ── Settings ───────────────────────────────────────────────────────────────
	{"UNKNOWN_SETTING", http.StatusNotFound, CategoryNotFound, "settings",
		"No setting is registered under that key.",
		"Read the settings list rather than guessing keys."},
	{"SETTING_REJECTED", http.StatusUnprocessableEntity, CategoryRequest, "settings",
		"The value did not pass the setting's own validation.",
		"Show the message against the field: it comes from the setting's rule."},
	{"SETTINGS_UNAVAILABLE", http.StatusInternalServerError, CategoryServer, "settings",
		"The settings store could not be read or written.",
		"Retry once, then report it."},
	{"NO_SCOPE", http.StatusBadRequest, CategoryRequest, "settings",
		"The request did not say which scope to act in.",
		"Send the scope the settings list gives for that key."},

	// ── Variants ───────────────────────────────────────────────────────────────
	{"OPTION_IN_USE", http.StatusConflict, CategoryConflict, "variants",
		"Variants are built on this option, so it cannot be removed.",
		"Clear the combinations that use it first, and say so rather than failing silently."},
	{"VALUE_IN_USE", http.StatusConflict, CategoryConflict, "variants",
		"That value is part of existing variants.",
		"Delete those variants first."},
	{"CANNOT_GENERATE", http.StatusUnprocessableEntity, CategoryRequest, "variants",
		"The combinations could not be generated from the options given.",
		"Check that every option has at least one value."},

	// ── Access reviews ─────────────────────────────────────────────────────────
	{"REVIEW_CLOSED", http.StatusBadRequest, CategoryState, "review",
		"The review is closed, so its decisions cannot change.",
		"Open a new review rather than editing a closed one."},
	{"REVIEW_INCOMPLETE", http.StatusBadRequest, CategoryState, "review",
		"Some items still have no decision, so the review cannot be completed.",
		"Show which items are outstanding."},
	{"ITEM_LOCKED", http.StatusBadRequest, CategoryState, "review",
		"That item already has a decision and will not take another.",
		"Re-read the review before submitting again."},
	{"INVALID_DECISION", http.StatusBadRequest, CategoryRequest, "review",
		"That is not a decision this review accepts.",
		"Use one of the documented decisions."},
	{"CANNOT_COMPLETE", http.StatusBadRequest, CategoryState, "review",
		"The review cannot be completed in its current state.",
		"Re-read it: the message says what is missing."},

	// ── Public forms ───────────────────────────────────────────────────────────
	{"PASSWORD_REQUIRED", http.StatusUnauthorized, CategoryAuth, "forms",
		"The shared form is password-protected.",
		"Prompt for the form's password and send it with the submission."},
	{"SUBMISSION_FAILED", http.StatusBadRequest, CategoryRequest, "forms",
		"The submission was refused: a field, a file or the form's own rules.",
		"Show the message. It is written for the person filling the form in."},

	// ── Offline sync ───────────────────────────────────────────────────────────
	{"MISSING_MODEL", http.StatusBadRequest, CategoryRequest, "sync",
		"The request did not name a model to sync.",
		"Send the model name the sync manifest lists."},
	{"UNKNOWN_MODEL", http.StatusBadRequest, CategoryRequest, "sync",
		"No model is registered under that name.",
		"Read the manifest rather than hard-coding names."},
	{"NOT_SYNCABLE", http.StatusBadRequest, CategoryRequest, "sync",
		"That model is not exposed to offline sync.",
		"Only sync models the manifest marks as syncable."},
	{"INVALID_SINCE", http.StatusBadRequest, CategoryRequest, "sync",
		"The since parameter is not an RFC3339 timestamp.",
		"Send the cursor the last sync returned, unchanged."},

	// ── Workflows and trees ────────────────────────────────────────────────────
	{"INVALID_TRANSITION", http.StatusUnprocessableEntity, CategoryState, "workflow",
		"That move is not declared in the workflow for this status.",
		"Offer only the transitions the API lists for the current status."},
	{"TRANSITION_REFUSED", http.StatusUnprocessableEntity, CategoryState, "workflow",
		"A transition hook refused the move, and nothing was written.",
		"Show the message: it is the business rule that said no."},
	{"INVALID_MOVE", http.StatusUnprocessableEntity, CategoryRequest, "tree",
		"That move would put a node inside its own subtree, or under a parent that cannot hold it.",
		"Refuse the drop in the UI rather than sending it."},

	// ── Reporting and observability ────────────────────────────────────────────
	{"CHART_FAILED", http.StatusBadRequest, CategoryRequest, "charts",
		"The chart could not be built from those parameters.",
		"Check the resource and preset against the ones the dashboard offers."},
	{"STATS_FAILED", http.StatusBadRequest, CategoryRequest, "charts",
		"The statistics could not be computed. This one deliberately conflates an unknown resource with a failed query, so a dashboard widget can render an error state instead of crashing.",
		"Render the widget's error state. The message says which of the two it was."},
	{"PDF_ERROR", http.StatusInternalServerError, CategoryServer, "pdf",
		"The PDF could not be rendered.",
		"Retry once. The record itself is unaffected."},
	{"DB_ERROR", http.StatusInternalServerError, CategoryServer, "observability",
		"A query behind a dashboard failed.",
		"Retry once, then report it."},
	{"SENTINEL_OFF", http.StatusServiceUnavailable, CategoryDisabled, "observability",
		"Sentinel is not enabled in this deployment, so there is nothing to report.",
		"Hide the security dashboard unless the API says it is on."},
	{"PULSE_OFF", http.StatusServiceUnavailable, CategoryDisabled, "observability",
		"Pulse is not enabled in this deployment.",
		"Hide the metrics dashboard unless the API says it is on."},
}

// areaLabels name the areas for a reader. The key is what the entries carry; the
// label is what the documentation shows.
var areaLabels = map[string]string{
	"core":          "Every endpoint",
	"auth":          "Sign-in and accounts",
	"twofactor":     "Two-factor authentication",
	"passkeys":      "Passkeys",
	"recovery":      "Account recovery",
	"apikeys":       "API keys",
	"uploads":       "Uploads and storage",
	"ai":            "AI gateway",
	"jobs":          "Background jobs",
	"import":        "CSV import",
	"backups":       "Backups and restore",
	"gdpr":          "GDPR and the audit log",
	"settings":      "Settings",
	"variants":      "Product variants",
	"review":        "Access reviews",
	"forms":         "Public forms",
	"sync":          "Offline sync",
	"workflow":      "Workflows",
	"tree":          "Trees",
	"charts":        "Charts and statistics",
	"pdf":           "PDF rendering",
	"observability": "Security and metrics dashboards",
	"webhooks":      "Incoming webhooks",
	"tenancy":       "Organizations (the multitenant plugin)",
}

// AreaLabel is the heading for an area. An area with no label is returned as it
// is, so a new one shows up looking unfinished rather than disappearing.
func AreaLabel(area string) string {
	if label, ok := areaLabels[area]; ok {
		return label
	}
	return area
}

// All returns the catalogue, sorted by code.
func All() []Entry {
	out := make([]Entry, len(catalog))
	copy(out, catalog)
	sort.Slice(out, func(i, j int) bool { return out[i].Code < out[j].Code })
	return out
}

// ByArea returns the catalogue grouped by area, and the area names in the order
// the catalogue declares them, which is roughly the order a reader meets them.
func ByArea() ([]string, map[string][]Entry) {
	var areas []string
	grouped := map[string][]Entry{}
	for _, entry := range catalog {
		if _, seen := grouped[entry.Area]; !seen {
			areas = append(areas, entry.Area)
		}
		grouped[entry.Area] = append(grouped[entry.Area], entry)
	}
	return areas, grouped
}

// Lookup finds one code.
func Lookup(code string) (Entry, bool) {
	for _, entry := range catalog {
		if entry.Code == code {
			return entry, true
		}
	}
	return Entry{}, false
}

// StatusOf is the status a code is always returned with, or 0 if it is not a
// catalogued code.
func StatusOf(code string) int {
	if entry, ok := Lookup(code); ok {
		return entry.Status
	}
	return 0
}
