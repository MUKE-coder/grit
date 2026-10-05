// Code generated from Grit's error catalogue. DO NOT EDIT.
//
// Every error this API returns has a code, and a code always arrives with the
// same status. That second half is the part a frontend needs and the part that
// used to be untrue: VALIDATION_ERROR came back as 422 in some handlers and 400
// in others, so a client had to handle whichever pairing it had happened to see.
//
// The same catalogue generates packages/shared/types/errors.ts, so a TypeScript
// switch over these codes can be exhaustive, and the table at /docs/backend/errors.
//
// Regenerate with: grit upgrade
package respond

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

// Code is one of the documented error codes. Use the constants below rather than
// a string literal: a typo in a literal is a code no client has ever heard of.
type Code string

// Category is what kind of problem a code reports, which is what decides how a
// caller should treat it.
type Category string

const (
	CategoryRequest    Category = "request"
	CategoryAuth       Category = "auth"
	CategoryPermission Category = "permission"
	CategoryNotFound   Category = "notfound"
	CategoryConflict   Category = "conflict"
	CategoryState      Category = "state"
	CategoryLimit      Category = "limit"
	CategoryServer     Category = "server"
	CategoryUpstream   Category = "upstream"
	CategoryDisabled   Category = "disabled"
)

const (
	// ACCOUNT_DISABLED — The account has been deactivated.
	CodeAccountDisabled Code = "ACCOUNT_DISABLED"
	// ACCOUNT_LOCKED — Too many wrong two-factor codes, so the account is locked for a while. Password sign-in reports a lock as INVALID_CREDENTIALS.
	CodeAccountLocked Code = "ACCOUNT_LOCKED"
	// AI_ERROR — The provider failed in a way that is not one of the above.
	CodeAIError Code = "AI_ERROR"
	// AI_FORBIDDEN — The provider refused this request, usually its own policy.
	CodeAIForbidden Code = "AI_FORBIDDEN"
	// AI_MODEL_NOT_FOUND — The provider does not know the model that was asked for.
	CodeAIModelNotFound Code = "AI_MODEL_NOT_FOUND"
	// AI_RATE_LIMITED — The provider is rate-limiting this deployment.
	CodeAIRateLimited Code = "AI_RATE_LIMITED"
	// AI_UNAUTHORIZED — The provider rejected the server's API key.
	CodeAIUnauthorized Code = "AI_UNAUTHORIZED"
	// AI_UNAVAILABLE — No AI provider is configured in this deployment.
	CodeAIUnavailable Code = "AI_UNAVAILABLE"
	// ALREADY_VERIFIED — The address on this link is already confirmed.
	CodeAlreadyVerified Code = "ALREADY_VERIFIED"
	// API_KEY_REQUIRED — The route is key-guarded and no key was sent.
	CodeAPIKeyRequired Code = "API_KEY_REQUIRED"
	// BAD_REQUEST — The request could not be understood at all.
	CodeBadRequest Code = "BAD_REQUEST"
	// CANNOT_COMPLETE — The review cannot be completed in its current state.
	CodeCannotComplete Code = "CANNOT_COMPLETE"
	// CANNOT_GENERATE — The combinations could not be generated from the options given.
	CodeCannotGenerate Code = "CANNOT_GENERATE"
	// CHART_FAILED — The chart could not be built from those parameters.
	CodeChartFailed Code = "CHART_FAILED"
	// CLEAR_FAILED — The queue could not be cleared.
	CodeClearFailed Code = "CLEAR_FAILED"
	// CONFLICT — The write collided with the state already there, such as a unique column.
	CodeConflict Code = "CONFLICT"
	// CSRF_INVALID — The CSRF token is missing or does not match the cookie.
	CodeCSRFInvalid Code = "CSRF_INVALID"
	// DB_ERROR — A query behind a dashboard failed.
	CodeDBError Code = "DB_ERROR"
	// EMAIL_EXISTS — An account already has that address.
	CodeEmailExists Code = "EMAIL_EXISTS"
	// EMAIL_NOT_VERIFIED — The account exists and its address has not been confirmed.
	CodeEmailNotVerified Code = "EMAIL_NOT_VERIFIED"
	// ENDPOINT_NOT_ALLOWED — The key is valid and is not allowed to call this endpoint.
	CodeEndpointNotAllowed Code = "ENDPOINT_NOT_ALLOWED"
	// ERASE_FAILED — The erasure did not complete.
	CodeEraseFailed Code = "ERASE_FAILED"
	// EXPORT_FAILED — The subject-access export could not be assembled.
	CodeExportFailed Code = "EXPORT_FAILED"
	// EXTRACT_FAILED — The archive could not be opened or does not hold what a restore needs.
	CodeExtractFailed Code = "EXTRACT_FAILED"
	// FILE_TOO_LARGE — The file is larger than this endpoint accepts.
	CodeFileTooLarge Code = "FILE_TOO_LARGE"
	// FORBIDDEN — The caller is signed in and this action is not theirs to take.
	CodeForbidden Code = "FORBIDDEN"
	// INSUFFICIENT_FUNDS — The balance is lower than the amount, and nothing was moved.
	CodeInsufficientFunds Code = "INSUFFICIENT_FUNDS"
	// INSUFFICIENT_STOCK — The row has less available than the quantity asked for, and nothing was taken.
	CodeInsufficientStock Code = "INSUFFICIENT_STOCK"
	// INTERNAL_ERROR — A fault on the server. The message is deliberately vague; the detail is in the server log.
	CodeInternalError Code = "INTERNAL_ERROR"
	// INVALID_API_KEY — The key is unknown, revoked or expired.
	CodeInvalidAPIKey Code = "INVALID_API_KEY"
	// INVALID_BACKUP_CODE — That backup code is wrong, or has been used.
	CodeInvalidBackupCode Code = "INVALID_BACKUP_CODE"
	// INVALID_BODY — The body was not valid JSON, or was not the shape this endpoint reads.
	CodeInvalidBody Code = "INVALID_BODY"
	// INVALID_CODE — The recovery code is wrong or has expired.
	CodeInvalidCode Code = "INVALID_CODE"
	// INVALID_CREDENTIALS — The email and password do not match an account. A locked account, and one that signs in only with a social provider, get this too until the password is right.
	CodeInvalidCredentials Code = "INVALID_CREDENTIALS"
	// INVALID_CSV — The file is not readable as CSV.
	CodeInvalidCSV Code = "INVALID_CSV"
	// INVALID_DECISION — That is not a decision this review accepts.
	CodeInvalidDecision Code = "INVALID_DECISION"
	// INVALID_FILE — No file was attached, or it could not be read.
	CodeInvalidFile Code = "INVALID_FILE"
	// INVALID_FILE_TYPE — The file's type is not accepted for this field.
	CodeInvalidFileType Code = "INVALID_FILE_TYPE"
	// INVALID_LINK — A one-time link, such as email verification or password reset, is wrong or has expired.
	CodeInvalidLink Code = "INVALID_LINK"
	// INVALID_MOVE — That move would put a node inside its own subtree, or under a parent that cannot hold it.
	CodeInvalidMove Code = "INVALID_MOVE"
	// INVALID_PASSWORD — The password given for a confirmation step is not correct.
	CodeInvalidPassword Code = "INVALID_PASSWORD"
	// INVALID_PENDING_TOKEN — The short-lived token between password and second factor is expired or unknown.
	CodeInvalidPendingToken Code = "INVALID_PENDING_TOKEN"
	// INVALID_RECOVERY_ADDRESS — The recovery email or phone number is not usable.
	CodeInvalidRecoveryAddress Code = "INVALID_RECOVERY_ADDRESS"
	// INVALID_SCHEDULE — The schedule is not a cron expression this API accepts.
	CodeInvalidSchedule Code = "INVALID_SCHEDULE"
	// INVALID_SIGNATURE — The webhook signature does not match the body and the shared secret.
	CodeInvalidSignature Code = "INVALID_SIGNATURE"
	// INVALID_SINCE — The since parameter is not an RFC3339 timestamp.
	CodeInvalidSince Code = "INVALID_SINCE"
	// INVALID_STATUS — That queue state is not one the endpoint accepts.
	CodeInvalidStatus Code = "INVALID_STATUS"
	// INVALID_TOKEN — The token is malformed, expired, or was issued for something else.
	CodeInvalidToken Code = "INVALID_TOKEN"
	// INVALID_TOTP_CODE — The six-digit code is wrong or has expired.
	CodeInvalidTOTPCode Code = "INVALID_TOTP_CODE"
	// INVALID_TRANSITION — That move is not declared in the workflow for this status.
	CodeInvalidTransition Code = "INVALID_TRANSITION"
	// ITEM_LOCKED — That item already has a decision and will not take another.
	CodeItemLocked Code = "ITEM_LOCKED"
	// JOB_ERROR — The import job could not be started.
	CodeJobError Code = "JOB_ERROR"
	// MAIL_FAILED — The mail driver refused or failed to send the verification code.
	CodeMailFailed Code = "MAIL_FAILED"
	// MAINTENANCE — The API is in maintenance mode and is refusing everything.
	CodeMaintenance Code = "MAINTENANCE"
	// MISSING_MODEL — The request did not name a model to sync.
	CodeMissingModel Code = "MISSING_MODEL"
	// MISSING_TOKEN — No token was sent where one is required.
	CodeMissingToken Code = "MISSING_TOKEN"
	// NOT_AVAILABLE — That backup cannot be downloaded: it is not finished, or it is not stored here.
	CodeNotAvailable Code = "NOT_AVAILABLE"
	// NOT_FOUND — No such row, or none this caller is allowed to see.
	CodeNotFound Code = "NOT_FOUND"
	// NOT_SYNCABLE — That model is not exposed to offline sync.
	CodeNotSyncable Code = "NOT_SYNCABLE"
	// NO_ORGANIZATION — The row belongs to an organization and the request has no active one: the caller belongs to none, or to several and named neither.
	CodeNoOrganization Code = "NO_ORGANIZATION"
	// NO_PASSWORD — The account has no password set, so it cannot be confirmed with one.
	CodeNoPassword Code = "NO_PASSWORD"
	// NO_SCOPE — The request did not say which scope to act in.
	CodeNoScope Code = "NO_SCOPE"
	// OPTION_IN_USE — Variants are built on this option, so it cannot be removed.
	CodeOptionInUse Code = "OPTION_IN_USE"
	// ORIGIN_NOT_ALLOWED — The request's Origin is not on the key's allowlist.
	CodeOriginNotAllowed Code = "ORIGIN_NOT_ALLOWED"
	// PASSKEYS_NOT_CONFIGURED — This deployment has no passkey configuration, so the endpoints are inert.
	CodePasskeysNotConfigured Code = "PASSKEYS_NOT_CONFIGURED"
	// PASSKEY_REJECTED — The challenge no longer exists: it expired, or it was already answered.
	CodePasskeyRejected Code = "PASSKEY_REJECTED"
	// PASSWORD_REQUIRED — The shared form is password-protected.
	CodePasswordRequired Code = "PASSWORD_REQUIRED"
	// PAYLOAD_TOO_LARGE — The request body is larger than the server accepts.
	CodePayloadTooLarge Code = "PAYLOAD_TOO_LARGE"
	// PAYMENTS_UNAVAILABLE — No payment provider is configured here: STRIPE_SECRET_KEY is not set.
	CodePaymentsUnavailable Code = "PAYMENTS_UNAVAILABLE"
	// PAYMENT_PROVIDER_ERROR — Stripe refused or failed the request. The message carries what Stripe said.
	CodePaymentProviderError Code = "PAYMENT_PROVIDER_ERROR"
	// PDF_ERROR — The PDF could not be rendered.
	CodePDFError Code = "PDF_ERROR"
	// PERSIST_FAILED — The change was accepted and could not be written.
	CodePersistFailed Code = "PERSIST_FAILED"
	// PRESIGN_FAILED — A presigned upload URL could not be produced.
	CodePresignFailed Code = "PRESIGN_FAILED"
	// PUBLISHABLE_KEY_NOT_ALLOWED — A publishable key was used where only a secret key is accepted.
	CodePublishableKeyNotAllowed Code = "PUBLISHABLE_KEY_NOT_ALLOWED"
	// PULSE_OFF — Pulse is not enabled in this deployment.
	CodePulseOff Code = "PULSE_OFF"
	// PULSE_UNAVAILABLE — Pulse is enabled but answered none of the performance dashboard's calls. The dashboard used to render zeros here, which reads as "no traffic and no errors".
	CodePulseUnavailable Code = "PULSE_UNAVAILABLE"
	// QUERY_FAILED — The audit query failed.
	CodeQueryFailed Code = "QUERY_FAILED"
	// RANGE_NOT_SATISFIABLE — The Range header asks for bytes outside the file. The Content-Range header gives the real size.
	CodeRangeNotSatisfiable Code = "RANGE_NOT_SATISFIABLE"
	// RATE_LIMITED — Too many requests from this caller.
	CodeRateLimited Code = "RATE_LIMITED"
	// READ_BODY_FAILED — The body could not be read to the end.
	CodeReadBodyFailed Code = "READ_BODY_FAILED"
	// REDIS_UNAVAILABLE — Redis is not configured or not reachable, so the queue cannot be read.
	CodeRedisUnavailable Code = "REDIS_UNAVAILABLE"
	// RESEAL_REFUSED — The audit chain was not broken where the reseal claimed, so nothing was resealed.
	CodeResealRefused Code = "RESEAL_REFUSED"
	// RETRY_FAILED — The job could not be re-queued.
	CodeRetryFailed Code = "RETRY_FAILED"
	// REVIEW_CLOSED — The review is closed, so its decisions cannot change.
	CodeReviewClosed Code = "REVIEW_CLOSED"
	// REVIEW_INCOMPLETE — Some items still have no decision, so the review cannot be completed.
	CodeReviewIncomplete Code = "REVIEW_INCOMPLETE"
	// SELF_ERASE — An account cannot erase itself through this endpoint.
	CodeSelfErase Code = "SELF_ERASE"
	// SENTINEL_OFF — Sentinel is not enabled in this deployment, so there is nothing to report.
	CodeSentinelOff Code = "SENTINEL_OFF"
	// SENTINEL_UNAVAILABLE — Sentinel is enabled but answered none of the security dashboard's calls. The dashboard used to render zeros here, which reads as "nothing is attacking you".
	CodeSentinelUnavailable Code = "SENTINEL_UNAVAILABLE"
	// SESSION_REVOKED — The session behind this token was signed out, on this device or another.
	CodeSessionRevoked Code = "SESSION_REVOKED"
	// SETTINGS_UNAVAILABLE — The settings store could not be read or written.
	CodeSettingsUnavailable Code = "SETTINGS_UNAVAILABLE"
	// SETTING_REJECTED — The value did not pass the setting's own validation.
	CodeSettingRejected Code = "SETTING_REJECTED"
	// SMS_FAILED — The SMS provider refused or failed to send.
	CodeSMSFailed Code = "SMS_FAILED"
	// SMS_NOT_CONFIGURED — No SMS provider is configured, so phone recovery cannot run here.
	CodeSMSNotConfigured Code = "SMS_NOT_CONFIGURED"
	// SOCIAL_AUTH_ONLY — The account signs in with a social provider and has no password. Only projects from before sign-in stopped reporting it return this; they now get INVALID_CREDENTIALS.
	CodeSocialAuthOnly Code = "SOCIAL_AUTH_ONLY"
	// STATS_FAILED — The statistics could not be computed. This one deliberately conflates an unknown resource with a failed query, so a dashboard widget can render an error state instead of crashing.
	CodeStatsFailed Code = "STATS_FAILED"
	// STORAGE_UNAVAILABLE — Object storage is not configured here, or is not answering.
	CodeStorageUnavailable Code = "STORAGE_UNAVAILABLE"
	// SUBMISSION_FAILED — The submission was refused: a field, a file or the form's own rules.
	CodeSubmissionFailed Code = "SUBMISSION_FAILED"
	// TEMP_ERROR — The upload could not be buffered to disk before importing.
	CodeTempError Code = "TEMP_ERROR"
	// TOKEN_ERROR — The access and refresh tokens could not be issued.
	CodeTokenError Code = "TOKEN_ERROR"
	// TOO_MANY_CHANGES — A sync push carried more than 500 changes.
	CodeTooManyChanges Code = "TOO_MANY_CHANGES"
	// TOTP_ALREADY_ENABLED — Two-factor authentication is already on for this account.
	CodeTOTPAlreadyEnabled Code = "TOTP_ALREADY_ENABLED"
	// TOTP_ERROR — The secret, QR code or backup codes could not be produced.
	CodeTOTPError Code = "TOTP_ERROR"
	// TOTP_NOT_ENABLED — The account has no second factor, so there is nothing to confirm or turn off.
	CodeTOTPNotEnabled Code = "TOTP_NOT_ENABLED"
	// TRANSITION_REFUSED — A transition hook refused the move, and nothing was written.
	CodeTransitionRefused Code = "TRANSITION_REFUSED"
	// UNAUTHORIZED — The credentials are missing, expired or not accepted.
	CodeUnauthorized Code = "UNAUTHORIZED"
	// UNKNOWN_MODEL — No model is registered under that name.
	CodeUnknownModel Code = "UNKNOWN_MODEL"
	// UNKNOWN_PROVIDER — No social provider is configured under that name.
	CodeUnknownProvider Code = "UNKNOWN_PROVIDER"
	// UNKNOWN_SETTING — No setting is registered under that key.
	CodeUnknownSetting Code = "UNKNOWN_SETTING"
	// UPLOAD_ALREADY_RECORDED — This upload has already been recorded, and a key is recorded only once.
	CodeUploadAlreadyRecorded Code = "UPLOAD_ALREADY_RECORDED"
	// UPLOAD_FAILED — The file reached the API and could not be stored.
	CodeUploadFailed Code = "UPLOAD_FAILED"
	// UPLOAD_KEY_FORBIDDEN — The key being recorded was not presigned for this user. Only keys under the caller's own uploads/<user_id>/ prefix are accepted.
	CodeUploadKeyForbidden Code = "UPLOAD_KEY_FORBIDDEN"
	// UPLOAD_NOT_FOUND — Nothing is stored under that key.
	CodeUploadNotFound Code = "UPLOAD_NOT_FOUND"
	// USER_ERROR — The signed-in account could not be loaded, which means the token outlived its row.
	CodeUserError Code = "USER_ERROR"
	// VALIDATION_ERROR — The body parsed, and a field in it is missing or not acceptable.
	CodeValidationError Code = "VALIDATION_ERROR"
	// VALUE_IN_USE — That value is part of existing variants.
	CodeValueInUse Code = "VALUE_IN_USE"
	// VERIFY_FAILED — The audit chain could not be verified.
	CodeVerifyFailed Code = "VERIFY_FAILED"
	// VERSION_CONFLICT — Somebody else changed the row since the version in your If-Match.
	CodeVersionConflict Code = "VERSION_CONFLICT"
)

// Meaning is everything the catalogue knows about one code.
type Meaning struct {
	Status   int
	Category Category
	// Area is the part of the API that raises it.
	Area string
	// What happened, from the caller's side.
	Meaning string
	// What the caller should do about it.
	Client string
}

// catalogue is the whole set. Generated, so it matches the documentation exactly.
var catalogue = map[Code]Meaning{
	CodeAccountDisabled: {Status: http.StatusForbidden, Category: CategoryPermission, Area: "auth",
		Meaning: "The account has been deactivated.",
		Client:  "Do not retry. This needs an administrator, not a different password."},
	CodeAccountLocked: {Status: http.StatusTooManyRequests, Category: CategoryLimit, Area: "auth",
		Meaning: "Too many wrong two-factor codes, so the account is locked for a while. Password sign-in reports a lock as INVALID_CREDENTIALS.",
		Client:  "Show the wait, and offer password reset. Retrying sooner extends nothing but the lock."},
	CodeAIError: {Status: http.StatusBadGateway, Category: CategoryUpstream, Area: "ai",
		Meaning: "The provider failed in a way that is not one of the above.",
		Client:  "Retry once. The message carries what the provider said."},
	CodeAIForbidden: {Status: http.StatusBadGateway, Category: CategoryUpstream, Area: "ai",
		Meaning: "The provider refused this request, usually its own policy.",
		Client:  "Do not retry the same prompt unchanged."},
	CodeAIModelNotFound: {Status: http.StatusBadGateway, Category: CategoryUpstream, Area: "ai",
		Meaning: "The provider does not know the model that was asked for.",
		Client:  "Pick a model the API lists. A model name can disappear without notice."},
	CodeAIRateLimited: {Status: http.StatusTooManyRequests, Category: CategoryLimit, Area: "ai",
		Meaning: "The provider is rate-limiting this deployment.",
		Client:  "Back off and retry with a delay. Queue rather than loop."},
	CodeAIUnauthorized: {Status: http.StatusBadGateway, Category: CategoryUpstream, Area: "ai",
		Meaning: "The provider rejected the server's API key.",
		Client:  "Nothing for the client to do. The key on the server is wrong or out of credit."},
	CodeAIUnavailable: {Status: http.StatusServiceUnavailable, Category: CategoryDisabled, Area: "ai",
		Meaning: "No AI provider is configured in this deployment.",
		Client:  "Hide AI features unless the API reports one as available."},
	CodeAlreadyVerified: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "auth",
		Meaning: "The address on this link is already confirmed.",
		Client:  "Treat it as success and continue to sign-in."},
	CodeAPIKeyRequired: {Status: http.StatusUnauthorized, Category: CategoryAuth, Area: "apikeys",
		Meaning: "The route is key-guarded and no key was sent.",
		Client:  "Send the key in the documented header. A user token is not a substitute here."},
	CodeBadRequest: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "core",
		Meaning: "The request could not be understood at all.",
		Client:  "Fix the request. Retrying the same one will fail the same way."},
	CodeCannotComplete: {Status: http.StatusBadRequest, Category: CategoryState, Area: "review",
		Meaning: "The review cannot be completed in its current state.",
		Client:  "Re-read it: the message says what is missing."},
	CodeCannotGenerate: {Status: http.StatusUnprocessableEntity, Category: CategoryRequest, Area: "variants",
		Meaning: "The combinations could not be generated from the options given.",
		Client:  "Check that every option has at least one value."},
	CodeChartFailed: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "charts",
		Meaning: "The chart could not be built from those parameters.",
		Client:  "Check the resource and preset against the ones the dashboard offers."},
	CodeClearFailed: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "jobs",
		Meaning: "The queue could not be cleared.",
		Client:  "Retry once, then look at the Redis connection."},
	CodeConflict: {Status: http.StatusConflict, Category: CategoryConflict, Area: "core",
		Meaning: "The write collided with the state already there, such as a unique column.",
		Client:  "Re-read, show what is there, and let the person decide."},
	CodeCSRFInvalid: {Status: http.StatusForbidden, Category: CategoryPermission, Area: "core",
		Meaning: "The CSRF token is missing or does not match the cookie.",
		Client:  "Read the CSRF cookie and send it back in the header on every unsafe request."},
	CodeDBError: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "observability",
		Meaning: "A query behind a dashboard failed.",
		Client:  "Retry once, then report it."},
	CodeEmailExists: {Status: http.StatusConflict, Category: CategoryConflict, Area: "auth",
		Meaning: "An account already has that address.",
		Client:  "Offer sign-in or password reset rather than registration."},
	CodeEmailNotVerified: {Status: http.StatusForbidden, Category: CategoryPermission, Area: "auth",
		Meaning: "The account exists and its address has not been confirmed.",
		Client:  "Send them to the verification flow, and offer to resend the link."},
	CodeEndpointNotAllowed: {Status: http.StatusForbidden, Category: CategoryPermission, Area: "apikeys",
		Meaning: "The key is valid and is not allowed to call this endpoint.",
		Client:  "Widen the key's endpoint list, or call it with one that may."},
	CodeEraseFailed: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "gdpr",
		Meaning: "The erasure did not complete.",
		Client:  "Do not assume anything was erased. Retry, and check the audit log."},
	CodeExportFailed: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "gdpr",
		Meaning: "The subject-access export could not be assembled.",
		Client:  "Retry once, then report it: this is a request with a legal clock on it."},
	CodeExtractFailed: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "backups",
		Meaning: "The archive could not be opened or does not hold what a restore needs.",
		Client:  "Upload an archive this API produced. A re-zipped one usually fails here."},
	CodeFileTooLarge: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "uploads",
		Meaning: "The file is larger than this endpoint accepts.",
		Client:  "Show the limit before the upload starts rather than after it finishes."},
	CodeForbidden: {Status: http.StatusForbidden, Category: CategoryPermission, Area: "core",
		Meaning: "The caller is signed in and this action is not theirs to take.",
		Client:  "Do not retry. Hide the action rather than letting it fail, if the role is known to the client."},
	CodeInsufficientFunds: {Status: http.StatusUnprocessableEntity, Category: CategoryState, Area: "money",
		Meaning: "The balance is lower than the amount, and nothing was moved.",
		Client:  "Show the balance. A retry only helps after money arrives."},
	CodeInsufficientStock: {Status: http.StatusUnprocessableEntity, Category: CategoryState, Area: "stock",
		Meaning: "The row has less available than the quantity asked for, and nothing was taken.",
		Client:  "Show what is left and let the person choose again; do not retry the same quantity."},
	CodeInternalError: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "core",
		Meaning: "A fault on the server. The message is deliberately vague; the detail is in the server log.",
		Client:  "Retry once, then report it. Nothing the client changes will help."},
	CodeInvalidAPIKey: {Status: http.StatusUnauthorized, Category: CategoryAuth, Area: "apikeys",
		Meaning: "The key is unknown, revoked or expired.",
		Client:  "Issue a new key. Do not retry with the same one."},
	CodeInvalidBackupCode: {Status: http.StatusUnauthorized, Category: CategoryAuth, Area: "twofactor",
		Meaning: "That backup code is wrong, or has been used.",
		Client:  "Each code works once. Offer the remaining count the status endpoint reports."},
	CodeInvalidBody: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "core",
		Meaning: "The body was not valid JSON, or was not the shape this endpoint reads.",
		Client:  "Send a JSON body matching the documented request type."},
	CodeInvalidCode: {Status: http.StatusUnprocessableEntity, Category: CategoryRequest, Area: "recovery",
		Meaning: "The recovery code is wrong or has expired.",
		Client:  "Offer to send a new one rather than retrying the same code."},
	CodeInvalidCredentials: {Status: http.StatusUnauthorized, Category: CategoryAuth, Area: "auth",
		Meaning: "The email and password do not match an account. A locked account, and one that signs in only with a social provider, get this too until the password is right.",
		Client:  "Say only that the details are wrong, and offer password reset: which of the two it was, and whether the account exists or is locked, is deliberately not reported."},
	CodeInvalidCSV: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "import",
		Meaning: "The file is not readable as CSV.",
		Client:  "Check the delimiter, the quoting and that the header row matches the template."},
	CodeInvalidDecision: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "review",
		Meaning: "That is not a decision this review accepts.",
		Client:  "Use one of the documented decisions."},
	CodeInvalidFile: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "uploads",
		Meaning: "No file was attached, or it could not be read.",
		Client:  "Send multipart form data with the documented field name."},
	CodeInvalidFileType: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "uploads",
		Meaning: "The file's type is not accepted for this field.",
		Client:  "Check the type client-side before uploading, and say which types are allowed."},
	CodeInvalidLink: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "core",
		Meaning: "A one-time link, such as email verification or password reset, is wrong or has expired.",
		Client:  "Offer to send a new link. This is not a sign-in failure: the caller has no credentials to fix, which is why it is 400 and not 401."},
	CodeInvalidMove: {Status: http.StatusUnprocessableEntity, Category: CategoryRequest, Area: "tree",
		Meaning: "That move would put a node inside its own subtree, or under a parent that cannot hold it.",
		Client:  "Refuse the drop in the UI rather than sending it."},
	CodeInvalidPassword: {Status: http.StatusUnauthorized, Category: CategoryAuth, Area: "auth",
		Meaning: "The password given for a confirmation step is not correct.",
		Client:  "Ask again. This is the re-authentication prompt, not a sign-in."},
	CodeInvalidPendingToken: {Status: http.StatusUnauthorized, Category: CategoryAuth, Area: "twofactor",
		Meaning: "The short-lived token between password and second factor is expired or unknown.",
		Client:  "Start the sign-in again from the password step."},
	CodeInvalidRecoveryAddress: {Status: http.StatusUnprocessableEntity, Category: CategoryRequest, Area: "recovery",
		Meaning: "The recovery email or phone number is not usable.",
		Client:  "Validate the format before sending, and show which one was rejected."},
	CodeInvalidSchedule: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "backups",
		Meaning: "The schedule is not a cron expression this API accepts.",
		Client:  "Validate the expression client-side, or offer fixed choices."},
	CodeInvalidSignature: {Status: http.StatusUnauthorized, Category: CategoryAuth, Area: "webhooks",
		Meaning: "The webhook signature does not match the body and the shared secret.",
		Client:  "Sign the exact bytes sent, with the secret for this endpoint. A reformatted body will not verify."},
	CodeInvalidSince: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "sync",
		Meaning: "The since parameter is not an RFC3339 timestamp.",
		Client:  "Send the cursor the last sync returned, unchanged."},
	CodeInvalidStatus: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "jobs",
		Meaning: "That queue state is not one the endpoint accepts.",
		Client:  "Use one of the documented states."},
	CodeInvalidToken: {Status: http.StatusUnauthorized, Category: CategoryAuth, Area: "core",
		Meaning: "The token is malformed, expired, or was issued for something else.",
		Client:  "Refresh it. A refresh token that fails this way has been used already or revoked: sign in again."},
	CodeInvalidTOTPCode: {Status: http.StatusUnauthorized, Category: CategoryAuth, Area: "twofactor",
		Meaning: "The six-digit code is wrong or has expired.",
		Client:  "Let them try the next code. Clock drift on the device is the usual cause of repeated failures."},
	CodeInvalidTransition: {Status: http.StatusUnprocessableEntity, Category: CategoryState, Area: "workflow",
		Meaning: "That move is not declared in the workflow for this status.",
		Client:  "Offer only the transitions the API lists for the current status."},
	CodeItemLocked: {Status: http.StatusBadRequest, Category: CategoryState, Area: "review",
		Meaning: "That item already has a decision and will not take another.",
		Client:  "Re-read the review before submitting again."},
	CodeJobError: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "import",
		Meaning: "The import job could not be started.",
		Client:  "Retry once. Nothing was imported."},
	CodeMailFailed: {Status: http.StatusBadGateway, Category: CategoryUpstream, Area: "recovery",
		Meaning: "The mail driver refused or failed to send the verification code.",
		Client:  "Retry once. If it keeps failing, the address may be unreachable or mail is misconfigured."},
	CodeMaintenance: {Status: http.StatusServiceUnavailable, Category: CategoryDisabled, Area: "core",
		Meaning: "The API is in maintenance mode and is refusing everything.",
		Client:  "Retry later. Show a maintenance state rather than an error."},
	CodeMissingModel: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "sync",
		Meaning: "The request did not name a model to sync.",
		Client:  "Send the model name the sync manifest lists."},
	CodeMissingToken: {Status: http.StatusUnauthorized, Category: CategoryAuth, Area: "core",
		Meaning: "No token was sent where one is required.",
		Client:  "Send the access token as Authorization: Bearer <token>."},
	CodeNotAvailable: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "backups",
		Meaning: "That backup cannot be downloaded: it is not finished, or it is not stored here.",
		Client:  "Re-read the backup's status before offering a download link."},
	CodeNotFound: {Status: http.StatusNotFound, Category: CategoryNotFound, Area: "core",
		Meaning: "No such row, or none this caller is allowed to see.",
		Client:  "Treat it as absent. On an owned resource this is also the answer for somebody else's row, on purpose: a wrong guess cannot be told from a right one."},
	CodeNotSyncable: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "sync",
		Meaning: "That model is not exposed to offline sync.",
		Client:  "Only sync models the manifest marks as syncable."},
	CodeNoOrganization: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "tenancy",
		Meaning: "The row belongs to an organization and the request has no active one: the caller belongs to none, or to several and named neither.",
		Client:  "Send the active organization as X-Organization-ID. If the caller belongs to no organization, they cannot read this at all: put them in one, or send them somewhere that does not need one."},
	CodeNoPassword: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "auth",
		Meaning: "The account has no password set, so it cannot be confirmed with one.",
		Client:  "Send them through set-a-password first."},
	CodeNoScope: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "settings",
		Meaning: "The request did not say which scope to act in.",
		Client:  "Send the scope the settings list gives for that key."},
	CodeOptionInUse: {Status: http.StatusConflict, Category: CategoryConflict, Area: "variants",
		Meaning: "Variants are built on this option, so it cannot be removed.",
		Client:  "Clear the combinations that use it first, and say so rather than failing silently."},
	CodeOriginNotAllowed: {Status: http.StatusForbidden, Category: CategoryPermission, Area: "apikeys",
		Meaning: "The request's Origin is not on the key's allowlist.",
		Client:  "Add the origin to the key, rather than relaxing CORS for everybody."},
	CodePasskeysNotConfigured: {Status: http.StatusNotImplemented, Category: CategoryDisabled, Area: "passkeys",
		Meaning: "This deployment has no passkey configuration, so the endpoints are inert.",
		Client:  "Hide passkey buttons unless the API reports the feature as available."},
	CodePasskeyRejected: {Status: http.StatusGone, Category: CategoryState, Area: "passkeys",
		Meaning: "The challenge no longer exists: it expired, or it was already answered.",
		Client:  "Start the ceremony again. Do not retry the same assertion."},
	CodePasswordRequired: {Status: http.StatusUnauthorized, Category: CategoryAuth, Area: "forms",
		Meaning: "The shared form is password-protected.",
		Client:  "Prompt for the form's password and send it with the submission."},
	CodePayloadTooLarge: {Status: http.StatusRequestEntityTooLarge, Category: CategoryRequest, Area: "core",
		Meaning: "The request body is larger than the server accepts.",
		Client:  "Send less, or upload the file directly to storage with a presigned URL."},
	CodePaymentsUnavailable: {Status: http.StatusServiceUnavailable, Category: CategoryDisabled, Area: "payments",
		Meaning: "No payment provider is configured here: STRIPE_SECRET_KEY is not set.",
		Client:  "Hide checkout unless payments are configured. Nothing the client sends will fix it."},
	CodePaymentProviderError: {Status: http.StatusBadGateway, Category: CategoryUpstream, Area: "payments",
		Meaning: "Stripe refused or failed the request. The message carries what Stripe said.",
		Client:  "Retry once. Nobody has been charged: a payment is only taken when the customer confirms it."},
	CodePDFError: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "pdf",
		Meaning: "The PDF could not be rendered.",
		Client:  "Retry once. The record itself is unaffected."},
	CodePersistFailed: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "core",
		Meaning: "The change was accepted and could not be written.",
		Client:  "Retry once. Treat the write as not having happened."},
	CodePresignFailed: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "uploads",
		Meaning: "A presigned upload URL could not be produced.",
		Client:  "Retry once, then fall back to uploading through the API."},
	CodePublishableKeyNotAllowed: {Status: http.StatusForbidden, Category: CategoryPermission, Area: "apikeys",
		Meaning: "A publishable key was used where only a secret key is accepted.",
		Client:  "Call this from the server with the secret key. A publishable key is public by design."},
	CodePulseOff: {Status: http.StatusServiceUnavailable, Category: CategoryDisabled, Area: "observability",
		Meaning: "Pulse is not enabled in this deployment.",
		Client:  "Hide the metrics dashboard unless the API says it is on."},
	CodePulseUnavailable: {Status: http.StatusBadGateway, Category: CategoryServer, Area: "observability",
		Meaning: "Pulse is enabled but answered none of the performance dashboard's calls. The dashboard used to render zeros here, which reads as \"no traffic and no errors\".",
		Client:  "Show the dashboard as unavailable rather than empty. details.degraded names the calls that failed."},
	CodeQueryFailed: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "gdpr",
		Meaning: "The audit query failed.",
		Client:  "Retry once, then report it."},
	CodeRangeNotSatisfiable: {Status: 416, Category: CategoryRequest, Area: "uploads",
		Meaning: "The Range header asks for bytes outside the file. The Content-Range header gives the real size.",
		Client:  "Read the size from Content-Range and request a range inside it, or drop the Range header to get the whole file."},
	CodeRateLimited: {Status: http.StatusTooManyRequests, Category: CategoryLimit, Area: "core",
		Meaning: "Too many requests from this caller.",
		Client:  "Back off. Honour Retry-After if it is present rather than retrying immediately."},
	CodeReadBodyFailed: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "core",
		Meaning: "The body could not be read to the end.",
		Client:  "Retry. If it keeps happening, the connection is dropping or the body is larger than the server accepts."},
	CodeRedisUnavailable: {Status: http.StatusServiceUnavailable, Category: CategoryDisabled, Area: "jobs",
		Meaning: "Redis is not configured or not reachable, so the queue cannot be read.",
		Client:  "Retry later. Jobs, cache and cron all depend on it."},
	CodeResealRefused: {Status: http.StatusConflict, Category: CategoryConflict, Area: "gdpr",
		Meaning: "The audit chain was not broken where the reseal claimed, so nothing was resealed.",
		Client:  "Verify the chain again and reseal from the entry the verification names."},
	CodeRetryFailed: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "jobs",
		Meaning: "The job could not be re-queued.",
		Client:  "Retry once. The job is still where it was."},
	CodeReviewClosed: {Status: http.StatusBadRequest, Category: CategoryState, Area: "review",
		Meaning: "The review is closed, so its decisions cannot change.",
		Client:  "Open a new review rather than editing a closed one."},
	CodeReviewIncomplete: {Status: http.StatusBadRequest, Category: CategoryState, Area: "review",
		Meaning: "Some items still have no decision, so the review cannot be completed.",
		Client:  "Show which items are outstanding."},
	CodeSelfErase: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "gdpr",
		Meaning: "An account cannot erase itself through this endpoint.",
		Client:  "Have another administrator run it, so the action has an actor who remains."},
	CodeSentinelOff: {Status: http.StatusServiceUnavailable, Category: CategoryDisabled, Area: "observability",
		Meaning: "Sentinel is not enabled in this deployment, so there is nothing to report.",
		Client:  "Hide the security dashboard unless the API says it is on."},
	CodeSentinelUnavailable: {Status: http.StatusBadGateway, Category: CategoryServer, Area: "observability",
		Meaning: "Sentinel is enabled but answered none of the security dashboard's calls. The dashboard used to render zeros here, which reads as \"nothing is attacking you\".",
		Client:  "Show the dashboard as unavailable rather than empty. details.degraded names the calls that failed."},
	CodeSessionRevoked: {Status: http.StatusUnauthorized, Category: CategoryAuth, Area: "core",
		Meaning: "The session behind this token was signed out, on this device or another.",
		Client:  "Sign in again. Do not retry with the same refresh token: reuse is what revoked it."},
	CodeSettingsUnavailable: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "settings",
		Meaning: "The settings store could not be read or written.",
		Client:  "Retry once, then report it."},
	CodeSettingRejected: {Status: http.StatusUnprocessableEntity, Category: CategoryRequest, Area: "settings",
		Meaning: "The value did not pass the setting's own validation.",
		Client:  "Show the message against the field: it comes from the setting's rule."},
	CodeSMSFailed: {Status: http.StatusBadGateway, Category: CategoryUpstream, Area: "recovery",
		Meaning: "The SMS provider refused or failed to send.",
		Client:  "Retry once, then offer email. The number may be unreachable."},
	CodeSMSNotConfigured: {Status: http.StatusNotImplemented, Category: CategoryDisabled, Area: "recovery",
		Meaning: "No SMS provider is configured, so phone recovery cannot run here.",
		Client:  "Offer email recovery instead."},
	CodeSocialAuthOnly: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "auth",
		Meaning: "The account signs in with a social provider and has no password. Only projects from before sign-in stopped reporting it return this; they now get INVALID_CREDENTIALS.",
		Client:  "Offer the provider button instead of the password form."},
	CodeStatsFailed: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "charts",
		Meaning: "The statistics could not be computed. This one deliberately conflates an unknown resource with a failed query, so a dashboard widget can render an error state instead of crashing.",
		Client:  "Render the widget's error state. The message says which of the two it was."},
	CodeStorageUnavailable: {Status: http.StatusServiceUnavailable, Category: CategoryDisabled, Area: "uploads",
		Meaning: "Object storage is not configured here, or is not answering.",
		Client:  "Retry later. Nothing the client sends will fix it."},
	CodeSubmissionFailed: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "forms",
		Meaning: "The submission was refused: a field, a file or the form's own rules.",
		Client:  "Show the message. It is written for the person filling the form in."},
	CodeTempError: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "import",
		Meaning: "The upload could not be buffered to disk before importing.",
		Client:  "Retry once. Check free disk on the server if it repeats."},
	CodeTokenError: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "auth",
		Meaning: "The access and refresh tokens could not be issued.",
		Client:  "Retry once. The credentials were accepted, so do not ask for them again."},
	CodeTooManyChanges: {Status: http.StatusRequestEntityTooLarge, Category: CategoryLimit, Area: "sync",
		Meaning: "A sync push carried more than 500 changes.",
		Client:  "Send the outbox in pushes of at most 500 changes, as the Grit sync clients do."},
	CodeTOTPAlreadyEnabled: {Status: http.StatusConflict, Category: CategoryConflict, Area: "twofactor",
		Meaning: "Two-factor authentication is already on for this account.",
		Client:  "Show it as enabled rather than offering setup."},
	CodeTOTPError: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "twofactor",
		Meaning: "The secret, QR code or backup codes could not be produced.",
		Client:  "Retry once, then report it. Nothing is half-enabled: setup only counts once confirmed."},
	CodeTOTPNotEnabled: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "twofactor",
		Meaning: "The account has no second factor, so there is nothing to confirm or turn off.",
		Client:  "Offer setup instead."},
	CodeTransitionRefused: {Status: http.StatusUnprocessableEntity, Category: CategoryState, Area: "workflow",
		Meaning: "A transition hook refused the move, and nothing was written.",
		Client:  "Show the message: it is the business rule that said no."},
	CodeUnauthorized: {Status: http.StatusUnauthorized, Category: CategoryAuth, Area: "core",
		Meaning: "The credentials are missing, expired or not accepted.",
		Client:  "Refresh the access token, and sign in again if the refresh is rejected too."},
	CodeUnknownModel: {Status: http.StatusBadRequest, Category: CategoryRequest, Area: "sync",
		Meaning: "No model is registered under that name.",
		Client:  "Read the manifest rather than hard-coding names."},
	CodeUnknownProvider: {Status: http.StatusNotFound, Category: CategoryNotFound, Area: "auth",
		Meaning: "No social provider is configured under that name.",
		Client:  "Only offer the providers the API reports as enabled."},
	CodeUnknownSetting: {Status: http.StatusNotFound, Category: CategoryNotFound, Area: "settings",
		Meaning: "No setting is registered under that key.",
		Client:  "Read the settings list rather than guessing keys."},
	CodeUploadAlreadyRecorded: {Status: http.StatusConflict, Category: CategoryConflict, Area: "uploads",
		Meaning: "This upload has already been recorded, and a key is recorded only once.",
		Client:  "Treat the upload as done: fetch it from the uploads list instead of completing it again."},
	CodeUploadFailed: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "uploads",
		Meaning: "The file reached the API and could not be stored.",
		Client:  "Retry once. Nothing was recorded, so there is no half-uploaded row to clean up."},
	CodeUploadKeyForbidden: {Status: http.StatusForbidden, Category: CategoryPermission, Area: "uploads",
		Meaning: "The key being recorded was not presigned for this user. Only keys under the caller's own uploads/<user_id>/ prefix are accepted.",
		Client:  "Send back the key the presign returned. Another user's key, or one outside the uploads, will not be accepted."},
	CodeUploadNotFound: {Status: http.StatusNotFound, Category: CategoryNotFound, Area: "uploads",
		Meaning: "Nothing is stored under that key.",
		Client:  "Treat it as absent: the key is wrong, or the object was deleted."},
	CodeUserError: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "auth",
		Meaning: "The signed-in account could not be loaded, which means the token outlived its row.",
		Client:  "Sign in again. If it repeats, the account data is inconsistent and needs a look."},
	CodeValidationError: {Status: http.StatusUnprocessableEntity, Category: CategoryRequest, Area: "core",
		Meaning: "The body parsed, and a field in it is missing or not acceptable.",
		Client:  "Read error.details: it maps each field to what is wrong with it. Show those against the inputs."},
	CodeValueInUse: {Status: http.StatusConflict, Category: CategoryConflict, Area: "variants",
		Meaning: "That value is part of existing variants.",
		Client:  "Delete those variants first."},
	CodeVerifyFailed: {Status: http.StatusInternalServerError, Category: CategoryServer, Area: "gdpr",
		Meaning: "The audit chain could not be verified.",
		Client:  "Report it. This is the check that says whether the log has been tampered with."},
	CodeVersionConflict: {Status: http.StatusConflict, Category: CategoryConflict, Area: "core",
		Meaning: "Somebody else changed the row since the version in your If-Match.",
		Client:  "The response carries the current version. Re-read, merge, and send the new ETag."},
}

// StatusOf is the status a code is always returned with, or 0 if the code is not
// one of the documented ones.
func StatusOf(code Code) int {
	return catalogue[code].Status
}

// Lookup returns what the catalogue says about a code.
func Lookup(code Code) (Meaning, bool) {
	meaning, ok := catalogue[code]
	return meaning, ok
}

// Codes lists every documented code, sorted, which is useful in a test that
// asserts your own handlers only return codes a client has been told about.
func Codes() []Code {
	out := make([]Code, 0, len(catalogue))
	for code := range catalogue {
		out = append(out, code)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

// Fail writes a documented error, at the status the catalogue gives it.
//
// Prefer this to writing c.JSON with a status and a code side by side: that is
// how one code ends up with two statuses. An unknown code is written as a 500,
// because a code nothing documents is a bug in the handler, not in the request.
//
//	respond.Fail(c, respond.CodeForbidden, "Only an owner can archive an invoice")
//
// Pass per-field messages for a validation failure:
//
//	respond.Fail(c, respond.CodeValidationError, "Check the highlighted fields",
//	    map[string]string{"email": "That address is already in use"})
func Fail(c *gin.Context, code Code, message string, details ...map[string]string) {
	status := StatusOf(code)
	if status == 0 {
		status = http.StatusInternalServerError
	}
	body := Error{Code: string(code), Message: message}
	if len(details) > 0 && len(details[0]) > 0 {
		body.Details = details[0]
	}
	c.AbortWithStatusJSON(status, gin.H{"error": body})
}
