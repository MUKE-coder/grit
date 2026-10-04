// Code generated from Grit's error catalogue. DO NOT EDIT.
//
// Every error the API returns has a code, and each code always arrives with the
// same HTTP status. Switch on the code, not on the status: the status groups
// errors, the code says which one it is.
//
//   import { API_ERRORS, type ApiErrorCode } from '@/shared/types/errors'
//
//   function explain(code: ApiErrorCode) {
//     return API_ERRORS[code].client   // what the person should do about it
//   }
//
// Regenerate with: grit upgrade

export type ApiErrorCategory =
  | 'request'
  | 'auth'
  | 'permission'
  | 'notfound'
  | 'conflict'
  | 'state'
  | 'limit'
  | 'server'
  | 'upstream'
  | 'disabled'

export type ApiErrorCode =
  | 'ACCOUNT_DISABLED'
  | 'ACCOUNT_LOCKED'
  | 'AI_ERROR'
  | 'AI_FORBIDDEN'
  | 'AI_MODEL_NOT_FOUND'
  | 'AI_RATE_LIMITED'
  | 'AI_UNAUTHORIZED'
  | 'AI_UNAVAILABLE'
  | 'ALREADY_VERIFIED'
  | 'API_KEY_REQUIRED'
  | 'BAD_REQUEST'
  | 'CANNOT_COMPLETE'
  | 'CANNOT_GENERATE'
  | 'CHART_FAILED'
  | 'CLEAR_FAILED'
  | 'CONFLICT'
  | 'CSRF_INVALID'
  | 'DB_ERROR'
  | 'EMAIL_EXISTS'
  | 'EMAIL_NOT_VERIFIED'
  | 'ENDPOINT_NOT_ALLOWED'
  | 'ERASE_FAILED'
  | 'EXPORT_FAILED'
  | 'EXTRACT_FAILED'
  | 'FILE_TOO_LARGE'
  | 'FORBIDDEN'
  | 'INSUFFICIENT_FUNDS'
  | 'INSUFFICIENT_STOCK'
  | 'INTERNAL_ERROR'
  | 'INVALID_API_KEY'
  | 'INVALID_BACKUP_CODE'
  | 'INVALID_BODY'
  | 'INVALID_CODE'
  | 'INVALID_CREDENTIALS'
  | 'INVALID_CSV'
  | 'INVALID_DECISION'
  | 'INVALID_FILE'
  | 'INVALID_FILE_TYPE'
  | 'INVALID_LINK'
  | 'INVALID_MOVE'
  | 'INVALID_PASSWORD'
  | 'INVALID_PENDING_TOKEN'
  | 'INVALID_RECOVERY_ADDRESS'
  | 'INVALID_SCHEDULE'
  | 'INVALID_SIGNATURE'
  | 'INVALID_SINCE'
  | 'INVALID_STATUS'
  | 'INVALID_TOKEN'
  | 'INVALID_TOTP_CODE'
  | 'INVALID_TRANSITION'
  | 'ITEM_LOCKED'
  | 'JOB_ERROR'
  | 'MAIL_FAILED'
  | 'MAINTENANCE'
  | 'MISSING_MODEL'
  | 'MISSING_TOKEN'
  | 'NOT_AVAILABLE'
  | 'NOT_FOUND'
  | 'NOT_SYNCABLE'
  | 'NO_ORGANIZATION'
  | 'NO_PASSWORD'
  | 'NO_SCOPE'
  | 'OPTION_IN_USE'
  | 'ORIGIN_NOT_ALLOWED'
  | 'PASSKEYS_NOT_CONFIGURED'
  | 'PASSKEY_REJECTED'
  | 'PASSWORD_REQUIRED'
  | 'PAYLOAD_TOO_LARGE'
  | 'PAYMENTS_UNAVAILABLE'
  | 'PAYMENT_PROVIDER_ERROR'
  | 'PDF_ERROR'
  | 'PERSIST_FAILED'
  | 'PRESIGN_FAILED'
  | 'PUBLISHABLE_KEY_NOT_ALLOWED'
  | 'PULSE_OFF'
  | 'PULSE_UNAVAILABLE'
  | 'QUERY_FAILED'
  | 'RANGE_NOT_SATISFIABLE'
  | 'RATE_LIMITED'
  | 'READ_BODY_FAILED'
  | 'REDIS_UNAVAILABLE'
  | 'RESEAL_REFUSED'
  | 'RETRY_FAILED'
  | 'REVIEW_CLOSED'
  | 'REVIEW_INCOMPLETE'
  | 'SELF_ERASE'
  | 'SENTINEL_OFF'
  | 'SENTINEL_UNAVAILABLE'
  | 'SESSION_REVOKED'
  | 'SETTINGS_UNAVAILABLE'
  | 'SETTING_REJECTED'
  | 'SMS_FAILED'
  | 'SMS_NOT_CONFIGURED'
  | 'SOCIAL_AUTH_ONLY'
  | 'STATS_FAILED'
  | 'STORAGE_UNAVAILABLE'
  | 'SUBMISSION_FAILED'
  | 'TEMP_ERROR'
  | 'TOKEN_ERROR'
  | 'TOO_MANY_CHANGES'
  | 'TOTP_ALREADY_ENABLED'
  | 'TOTP_ERROR'
  | 'TOTP_NOT_ENABLED'
  | 'TRANSITION_REFUSED'
  | 'UNAUTHORIZED'
  | 'UNKNOWN_MODEL'
  | 'UNKNOWN_PROVIDER'
  | 'UNKNOWN_SETTING'
  | 'UPLOAD_ALREADY_RECORDED'
  | 'UPLOAD_FAILED'
  | 'UPLOAD_KEY_FORBIDDEN'
  | 'UPLOAD_NOT_FOUND'
  | 'USER_ERROR'
  | 'VALIDATION_ERROR'
  | 'VALUE_IN_USE'
  | 'VERIFY_FAILED'
  | 'VERSION_CONFLICT'

export interface ApiErrorInfo {
  /** The status this code always arrives with. */
  status: number
  category: ApiErrorCategory
  /** Which part of the API raises it. */
  area: string
  /** What happened, from the caller's side. */
  meaning: string
  /** What to do about it. */
  client: string
}

export const API_ERRORS: Record<ApiErrorCode, ApiErrorInfo> = {
  ACCOUNT_DISABLED: {
    status: 403,
    category: 'permission',
    area: 'auth',
    meaning: 'The account has been deactivated.',
    client: 'Do not retry. This needs an administrator, not a different password.',
  },
  ACCOUNT_LOCKED: {
    status: 429,
    category: 'limit',
    area: 'auth',
    meaning: 'Too many wrong two-factor codes, so the account is locked for a while. Password sign-in reports a lock as INVALID_CREDENTIALS.',
    client: 'Show the wait, and offer password reset. Retrying sooner extends nothing but the lock.',
  },
  AI_ERROR: {
    status: 502,
    category: 'upstream',
    area: 'ai',
    meaning: 'The provider failed in a way that is not one of the above.',
    client: 'Retry once. The message carries what the provider said.',
  },
  AI_FORBIDDEN: {
    status: 502,
    category: 'upstream',
    area: 'ai',
    meaning: 'The provider refused this request, usually its own policy.',
    client: 'Do not retry the same prompt unchanged.',
  },
  AI_MODEL_NOT_FOUND: {
    status: 502,
    category: 'upstream',
    area: 'ai',
    meaning: 'The provider does not know the model that was asked for.',
    client: 'Pick a model the API lists. A model name can disappear without notice.',
  },
  AI_RATE_LIMITED: {
    status: 429,
    category: 'limit',
    area: 'ai',
    meaning: 'The provider is rate-limiting this deployment.',
    client: 'Back off and retry with a delay. Queue rather than loop.',
  },
  AI_UNAUTHORIZED: {
    status: 502,
    category: 'upstream',
    area: 'ai',
    meaning: 'The provider rejected the server\'s API key.',
    client: 'Nothing for the client to do. The key on the server is wrong or out of credit.',
  },
  AI_UNAVAILABLE: {
    status: 503,
    category: 'disabled',
    area: 'ai',
    meaning: 'No AI provider is configured in this deployment.',
    client: 'Hide AI features unless the API reports one as available.',
  },
  ALREADY_VERIFIED: {
    status: 400,
    category: 'request',
    area: 'auth',
    meaning: 'The address on this link is already confirmed.',
    client: 'Treat it as success and continue to sign-in.',
  },
  API_KEY_REQUIRED: {
    status: 401,
    category: 'auth',
    area: 'apikeys',
    meaning: 'The route is key-guarded and no key was sent.',
    client: 'Send the key in the documented header. A user token is not a substitute here.',
  },
  BAD_REQUEST: {
    status: 400,
    category: 'request',
    area: 'core',
    meaning: 'The request could not be understood at all.',
    client: 'Fix the request. Retrying the same one will fail the same way.',
  },
  CANNOT_COMPLETE: {
    status: 400,
    category: 'state',
    area: 'review',
    meaning: 'The review cannot be completed in its current state.',
    client: 'Re-read it: the message says what is missing.',
  },
  CANNOT_GENERATE: {
    status: 422,
    category: 'request',
    area: 'variants',
    meaning: 'The combinations could not be generated from the options given.',
    client: 'Check that every option has at least one value.',
  },
  CHART_FAILED: {
    status: 400,
    category: 'request',
    area: 'charts',
    meaning: 'The chart could not be built from those parameters.',
    client: 'Check the resource and preset against the ones the dashboard offers.',
  },
  CLEAR_FAILED: {
    status: 500,
    category: 'server',
    area: 'jobs',
    meaning: 'The queue could not be cleared.',
    client: 'Retry once, then look at the Redis connection.',
  },
  CONFLICT: {
    status: 409,
    category: 'conflict',
    area: 'core',
    meaning: 'The write collided with the state already there, such as a unique column.',
    client: 'Re-read, show what is there, and let the person decide.',
  },
  CSRF_INVALID: {
    status: 403,
    category: 'permission',
    area: 'core',
    meaning: 'The CSRF token is missing or does not match the cookie.',
    client: 'Read the CSRF cookie and send it back in the header on every unsafe request.',
  },
  DB_ERROR: {
    status: 500,
    category: 'server',
    area: 'observability',
    meaning: 'A query behind a dashboard failed.',
    client: 'Retry once, then report it.',
  },
  EMAIL_EXISTS: {
    status: 409,
    category: 'conflict',
    area: 'auth',
    meaning: 'An account already has that address.',
    client: 'Offer sign-in or password reset rather than registration.',
  },
  EMAIL_NOT_VERIFIED: {
    status: 403,
    category: 'permission',
    area: 'auth',
    meaning: 'The account exists and its address has not been confirmed.',
    client: 'Send them to the verification flow, and offer to resend the link.',
  },
  ENDPOINT_NOT_ALLOWED: {
    status: 403,
    category: 'permission',
    area: 'apikeys',
    meaning: 'The key is valid and is not allowed to call this endpoint.',
    client: 'Widen the key\'s endpoint list, or call it with one that may.',
  },
  ERASE_FAILED: {
    status: 500,
    category: 'server',
    area: 'gdpr',
    meaning: 'The erasure did not complete.',
    client: 'Do not assume anything was erased. Retry, and check the audit log.',
  },
  EXPORT_FAILED: {
    status: 500,
    category: 'server',
    area: 'gdpr',
    meaning: 'The subject-access export could not be assembled.',
    client: 'Retry once, then report it: this is a request with a legal clock on it.',
  },
  EXTRACT_FAILED: {
    status: 400,
    category: 'request',
    area: 'backups',
    meaning: 'The archive could not be opened or does not hold what a restore needs.',
    client: 'Upload an archive this API produced. A re-zipped one usually fails here.',
  },
  FILE_TOO_LARGE: {
    status: 400,
    category: 'request',
    area: 'uploads',
    meaning: 'The file is larger than this endpoint accepts.',
    client: 'Show the limit before the upload starts rather than after it finishes.',
  },
  FORBIDDEN: {
    status: 403,
    category: 'permission',
    area: 'core',
    meaning: 'The caller is signed in and this action is not theirs to take.',
    client: 'Do not retry. Hide the action rather than letting it fail, if the role is known to the client.',
  },
  INSUFFICIENT_FUNDS: {
    status: 422,
    category: 'state',
    area: 'money',
    meaning: 'The balance is lower than the amount, and nothing was moved.',
    client: 'Show the balance. A retry only helps after money arrives.',
  },
  INSUFFICIENT_STOCK: {
    status: 422,
    category: 'state',
    area: 'stock',
    meaning: 'The row has less available than the quantity asked for, and nothing was taken.',
    client: 'Show what is left and let the person choose again; do not retry the same quantity.',
  },
  INTERNAL_ERROR: {
    status: 500,
    category: 'server',
    area: 'core',
    meaning: 'A fault on the server. The message is deliberately vague; the detail is in the server log.',
    client: 'Retry once, then report it. Nothing the client changes will help.',
  },
  INVALID_API_KEY: {
    status: 401,
    category: 'auth',
    area: 'apikeys',
    meaning: 'The key is unknown, revoked or expired.',
    client: 'Issue a new key. Do not retry with the same one.',
  },
  INVALID_BACKUP_CODE: {
    status: 401,
    category: 'auth',
    area: 'twofactor',
    meaning: 'That backup code is wrong, or has been used.',
    client: 'Each code works once. Offer the remaining count the status endpoint reports.',
  },
  INVALID_BODY: {
    status: 400,
    category: 'request',
    area: 'core',
    meaning: 'The body was not valid JSON, or was not the shape this endpoint reads.',
    client: 'Send a JSON body matching the documented request type.',
  },
  INVALID_CODE: {
    status: 422,
    category: 'request',
    area: 'recovery',
    meaning: 'The recovery code is wrong or has expired.',
    client: 'Offer to send a new one rather than retrying the same code.',
  },
  INVALID_CREDENTIALS: {
    status: 401,
    category: 'auth',
    area: 'auth',
    meaning: 'The email and password do not match an account. A locked account, and one that signs in only with a social provider, get this too until the password is right.',
    client: 'Say only that the details are wrong, and offer password reset: which of the two it was, and whether the account exists or is locked, is deliberately not reported.',
  },
  INVALID_CSV: {
    status: 400,
    category: 'request',
    area: 'import',
    meaning: 'The file is not readable as CSV.',
    client: 'Check the delimiter, the quoting and that the header row matches the template.',
  },
  INVALID_DECISION: {
    status: 400,
    category: 'request',
    area: 'review',
    meaning: 'That is not a decision this review accepts.',
    client: 'Use one of the documented decisions.',
  },
  INVALID_FILE: {
    status: 400,
    category: 'request',
    area: 'uploads',
    meaning: 'No file was attached, or it could not be read.',
    client: 'Send multipart form data with the documented field name.',
  },
  INVALID_FILE_TYPE: {
    status: 400,
    category: 'request',
    area: 'uploads',
    meaning: 'The file\'s type is not accepted for this field.',
    client: 'Check the type client-side before uploading, and say which types are allowed.',
  },
  INVALID_LINK: {
    status: 400,
    category: 'request',
    area: 'core',
    meaning: 'A one-time link, such as email verification or password reset, is wrong or has expired.',
    client: 'Offer to send a new link. This is not a sign-in failure: the caller has no credentials to fix, which is why it is 400 and not 401.',
  },
  INVALID_MOVE: {
    status: 422,
    category: 'request',
    area: 'tree',
    meaning: 'That move would put a node inside its own subtree, or under a parent that cannot hold it.',
    client: 'Refuse the drop in the UI rather than sending it.',
  },
  INVALID_PASSWORD: {
    status: 401,
    category: 'auth',
    area: 'auth',
    meaning: 'The password given for a confirmation step is not correct.',
    client: 'Ask again. This is the re-authentication prompt, not a sign-in.',
  },
  INVALID_PENDING_TOKEN: {
    status: 401,
    category: 'auth',
    area: 'twofactor',
    meaning: 'The short-lived token between password and second factor is expired or unknown.',
    client: 'Start the sign-in again from the password step.',
  },
  INVALID_RECOVERY_ADDRESS: {
    status: 422,
    category: 'request',
    area: 'recovery',
    meaning: 'The recovery email or phone number is not usable.',
    client: 'Validate the format before sending, and show which one was rejected.',
  },
  INVALID_SCHEDULE: {
    status: 400,
    category: 'request',
    area: 'backups',
    meaning: 'The schedule is not a cron expression this API accepts.',
    client: 'Validate the expression client-side, or offer fixed choices.',
  },
  INVALID_SIGNATURE: {
    status: 401,
    category: 'auth',
    area: 'webhooks',
    meaning: 'The webhook signature does not match the body and the shared secret.',
    client: 'Sign the exact bytes sent, with the secret for this endpoint. A reformatted body will not verify.',
  },
  INVALID_SINCE: {
    status: 400,
    category: 'request',
    area: 'sync',
    meaning: 'The since parameter is not an RFC3339 timestamp.',
    client: 'Send the cursor the last sync returned, unchanged.',
  },
  INVALID_STATUS: {
    status: 400,
    category: 'request',
    area: 'jobs',
    meaning: 'That queue state is not one the endpoint accepts.',
    client: 'Use one of the documented states.',
  },
  INVALID_TOKEN: {
    status: 401,
    category: 'auth',
    area: 'core',
    meaning: 'The token is malformed, expired, or was issued for something else.',
    client: 'Refresh it. A refresh token that fails this way has been used already or revoked: sign in again.',
  },
  INVALID_TOTP_CODE: {
    status: 401,
    category: 'auth',
    area: 'twofactor',
    meaning: 'The six-digit code is wrong or has expired.',
    client: 'Let them try the next code. Clock drift on the device is the usual cause of repeated failures.',
  },
  INVALID_TRANSITION: {
    status: 422,
    category: 'state',
    area: 'workflow',
    meaning: 'That move is not declared in the workflow for this status.',
    client: 'Offer only the transitions the API lists for the current status.',
  },
  ITEM_LOCKED: {
    status: 400,
    category: 'state',
    area: 'review',
    meaning: 'That item already has a decision and will not take another.',
    client: 'Re-read the review before submitting again.',
  },
  JOB_ERROR: {
    status: 500,
    category: 'server',
    area: 'import',
    meaning: 'The import job could not be started.',
    client: 'Retry once. Nothing was imported.',
  },
  MAIL_FAILED: {
    status: 502,
    category: 'upstream',
    area: 'recovery',
    meaning: 'The mail driver refused or failed to send the verification code.',
    client: 'Retry once. If it keeps failing, the address may be unreachable or mail is misconfigured.',
  },
  MAINTENANCE: {
    status: 503,
    category: 'disabled',
    area: 'core',
    meaning: 'The API is in maintenance mode and is refusing everything.',
    client: 'Retry later. Show a maintenance state rather than an error.',
  },
  MISSING_MODEL: {
    status: 400,
    category: 'request',
    area: 'sync',
    meaning: 'The request did not name a model to sync.',
    client: 'Send the model name the sync manifest lists.',
  },
  MISSING_TOKEN: {
    status: 401,
    category: 'auth',
    area: 'core',
    meaning: 'No token was sent where one is required.',
    client: 'Send the access token as Authorization: Bearer <token>.',
  },
  NOT_AVAILABLE: {
    status: 400,
    category: 'request',
    area: 'backups',
    meaning: 'That backup cannot be downloaded: it is not finished, or it is not stored here.',
    client: 'Re-read the backup\'s status before offering a download link.',
  },
  NOT_FOUND: {
    status: 404,
    category: 'notfound',
    area: 'core',
    meaning: 'No such row, or none this caller is allowed to see.',
    client: 'Treat it as absent. On an owned resource this is also the answer for somebody else\'s row, on purpose: a wrong guess cannot be told from a right one.',
  },
  NOT_SYNCABLE: {
    status: 400,
    category: 'request',
    area: 'sync',
    meaning: 'That model is not exposed to offline sync.',
    client: 'Only sync models the manifest marks as syncable.',
  },
  NO_ORGANIZATION: {
    status: 400,
    category: 'request',
    area: 'tenancy',
    meaning: 'The row belongs to an organization and the request has no active one: the caller belongs to none, or to several and named neither.',
    client: 'Send the active organization as X-Organization-ID. If the caller belongs to no organization, they cannot read this at all: put them in one, or send them somewhere that does not need one.',
  },
  NO_PASSWORD: {
    status: 400,
    category: 'request',
    area: 'auth',
    meaning: 'The account has no password set, so it cannot be confirmed with one.',
    client: 'Send them through set-a-password first.',
  },
  NO_SCOPE: {
    status: 400,
    category: 'request',
    area: 'settings',
    meaning: 'The request did not say which scope to act in.',
    client: 'Send the scope the settings list gives for that key.',
  },
  OPTION_IN_USE: {
    status: 409,
    category: 'conflict',
    area: 'variants',
    meaning: 'Variants are built on this option, so it cannot be removed.',
    client: 'Clear the combinations that use it first, and say so rather than failing silently.',
  },
  ORIGIN_NOT_ALLOWED: {
    status: 403,
    category: 'permission',
    area: 'apikeys',
    meaning: 'The request\'s Origin is not on the key\'s allowlist.',
    client: 'Add the origin to the key, rather than relaxing CORS for everybody.',
  },
  PASSKEYS_NOT_CONFIGURED: {
    status: 501,
    category: 'disabled',
    area: 'passkeys',
    meaning: 'This deployment has no passkey configuration, so the endpoints are inert.',
    client: 'Hide passkey buttons unless the API reports the feature as available.',
  },
  PASSKEY_REJECTED: {
    status: 410,
    category: 'state',
    area: 'passkeys',
    meaning: 'The challenge no longer exists: it expired, or it was already answered.',
    client: 'Start the ceremony again. Do not retry the same assertion.',
  },
  PASSWORD_REQUIRED: {
    status: 401,
    category: 'auth',
    area: 'forms',
    meaning: 'The shared form is password-protected.',
    client: 'Prompt for the form\'s password and send it with the submission.',
  },
  PAYLOAD_TOO_LARGE: {
    status: 413,
    category: 'request',
    area: 'core',
    meaning: 'The request body is larger than the server accepts.',
    client: 'Send less, or upload the file directly to storage with a presigned URL.',
  },
  PAYMENTS_UNAVAILABLE: {
    status: 503,
    category: 'disabled',
    area: 'payments',
    meaning: 'No payment provider is configured here: STRIPE_SECRET_KEY is not set.',
    client: 'Hide checkout unless payments are configured. Nothing the client sends will fix it.',
  },
  PAYMENT_PROVIDER_ERROR: {
    status: 502,
    category: 'upstream',
    area: 'payments',
    meaning: 'Stripe refused or failed the request. The message carries what Stripe said.',
    client: 'Retry once. Nobody has been charged: a payment is only taken when the customer confirms it.',
  },
  PDF_ERROR: {
    status: 500,
    category: 'server',
    area: 'pdf',
    meaning: 'The PDF could not be rendered.',
    client: 'Retry once. The record itself is unaffected.',
  },
  PERSIST_FAILED: {
    status: 500,
    category: 'server',
    area: 'core',
    meaning: 'The change was accepted and could not be written.',
    client: 'Retry once. Treat the write as not having happened.',
  },
  PRESIGN_FAILED: {
    status: 500,
    category: 'server',
    area: 'uploads',
    meaning: 'A presigned upload URL could not be produced.',
    client: 'Retry once, then fall back to uploading through the API.',
  },
  PUBLISHABLE_KEY_NOT_ALLOWED: {
    status: 403,
    category: 'permission',
    area: 'apikeys',
    meaning: 'A publishable key was used where only a secret key is accepted.',
    client: 'Call this from the server with the secret key. A publishable key is public by design.',
  },
  PULSE_OFF: {
    status: 503,
    category: 'disabled',
    area: 'observability',
    meaning: 'Pulse is not enabled in this deployment.',
    client: 'Hide the metrics dashboard unless the API says it is on.',
  },
  PULSE_UNAVAILABLE: {
    status: 502,
    category: 'server',
    area: 'observability',
    meaning: 'Pulse is enabled but answered none of the performance dashboard\'s calls. The dashboard used to render zeros here, which reads as "no traffic and no errors".',
    client: 'Show the dashboard as unavailable rather than empty. details.degraded names the calls that failed.',
  },
  QUERY_FAILED: {
    status: 500,
    category: 'server',
    area: 'gdpr',
    meaning: 'The audit query failed.',
    client: 'Retry once, then report it.',
  },
  RANGE_NOT_SATISFIABLE: {
    status: 416,
    category: 'request',
    area: 'uploads',
    meaning: 'The Range header asks for bytes outside the file. The Content-Range header gives the real size.',
    client: 'Read the size from Content-Range and request a range inside it, or drop the Range header to get the whole file.',
  },
  RATE_LIMITED: {
    status: 429,
    category: 'limit',
    area: 'core',
    meaning: 'Too many requests from this caller.',
    client: 'Back off. Honour Retry-After if it is present rather than retrying immediately.',
  },
  READ_BODY_FAILED: {
    status: 400,
    category: 'request',
    area: 'core',
    meaning: 'The body could not be read to the end.',
    client: 'Retry. If it keeps happening, the connection is dropping or the body is larger than the server accepts.',
  },
  REDIS_UNAVAILABLE: {
    status: 503,
    category: 'disabled',
    area: 'jobs',
    meaning: 'Redis is not configured or not reachable, so the queue cannot be read.',
    client: 'Retry later. Jobs, cache and cron all depend on it.',
  },
  RESEAL_REFUSED: {
    status: 409,
    category: 'conflict',
    area: 'gdpr',
    meaning: 'The audit chain was not broken where the reseal claimed, so nothing was resealed.',
    client: 'Verify the chain again and reseal from the entry the verification names.',
  },
  RETRY_FAILED: {
    status: 500,
    category: 'server',
    area: 'jobs',
    meaning: 'The job could not be re-queued.',
    client: 'Retry once. The job is still where it was.',
  },
  REVIEW_CLOSED: {
    status: 400,
    category: 'state',
    area: 'review',
    meaning: 'The review is closed, so its decisions cannot change.',
    client: 'Open a new review rather than editing a closed one.',
  },
  REVIEW_INCOMPLETE: {
    status: 400,
    category: 'state',
    area: 'review',
    meaning: 'Some items still have no decision, so the review cannot be completed.',
    client: 'Show which items are outstanding.',
  },
  SELF_ERASE: {
    status: 400,
    category: 'request',
    area: 'gdpr',
    meaning: 'An account cannot erase itself through this endpoint.',
    client: 'Have another administrator run it, so the action has an actor who remains.',
  },
  SENTINEL_OFF: {
    status: 503,
    category: 'disabled',
    area: 'observability',
    meaning: 'Sentinel is not enabled in this deployment, so there is nothing to report.',
    client: 'Hide the security dashboard unless the API says it is on.',
  },
  SENTINEL_UNAVAILABLE: {
    status: 502,
    category: 'server',
    area: 'observability',
    meaning: 'Sentinel is enabled but answered none of the security dashboard\'s calls. The dashboard used to render zeros here, which reads as "nothing is attacking you".',
    client: 'Show the dashboard as unavailable rather than empty. details.degraded names the calls that failed.',
  },
  SESSION_REVOKED: {
    status: 401,
    category: 'auth',
    area: 'core',
    meaning: 'The session behind this token was signed out, on this device or another.',
    client: 'Sign in again. Do not retry with the same refresh token: reuse is what revoked it.',
  },
  SETTINGS_UNAVAILABLE: {
    status: 500,
    category: 'server',
    area: 'settings',
    meaning: 'The settings store could not be read or written.',
    client: 'Retry once, then report it.',
  },
  SETTING_REJECTED: {
    status: 422,
    category: 'request',
    area: 'settings',
    meaning: 'The value did not pass the setting\'s own validation.',
    client: 'Show the message against the field: it comes from the setting\'s rule.',
  },
  SMS_FAILED: {
    status: 502,
    category: 'upstream',
    area: 'recovery',
    meaning: 'The SMS provider refused or failed to send.',
    client: 'Retry once, then offer email. The number may be unreachable.',
  },
  SMS_NOT_CONFIGURED: {
    status: 501,
    category: 'disabled',
    area: 'recovery',
    meaning: 'No SMS provider is configured, so phone recovery cannot run here.',
    client: 'Offer email recovery instead.',
  },
  SOCIAL_AUTH_ONLY: {
    status: 400,
    category: 'request',
    area: 'auth',
    meaning: 'The account signs in with a social provider and has no password. Only projects from before sign-in stopped reporting it return this; they now get INVALID_CREDENTIALS.',
    client: 'Offer the provider button instead of the password form.',
  },
  STATS_FAILED: {
    status: 400,
    category: 'request',
    area: 'charts',
    meaning: 'The statistics could not be computed. This one deliberately conflates an unknown resource with a failed query, so a dashboard widget can render an error state instead of crashing.',
    client: 'Render the widget\'s error state. The message says which of the two it was.',
  },
  STORAGE_UNAVAILABLE: {
    status: 503,
    category: 'disabled',
    area: 'uploads',
    meaning: 'Object storage is not configured here, or is not answering.',
    client: 'Retry later. Nothing the client sends will fix it.',
  },
  SUBMISSION_FAILED: {
    status: 400,
    category: 'request',
    area: 'forms',
    meaning: 'The submission was refused: a field, a file or the form\'s own rules.',
    client: 'Show the message. It is written for the person filling the form in.',
  },
  TEMP_ERROR: {
    status: 500,
    category: 'server',
    area: 'import',
    meaning: 'The upload could not be buffered to disk before importing.',
    client: 'Retry once. Check free disk on the server if it repeats.',
  },
  TOKEN_ERROR: {
    status: 500,
    category: 'server',
    area: 'auth',
    meaning: 'The access and refresh tokens could not be issued.',
    client: 'Retry once. The credentials were accepted, so do not ask for them again.',
  },
  TOO_MANY_CHANGES: {
    status: 413,
    category: 'limit',
    area: 'sync',
    meaning: 'A sync push carried more than 500 changes.',
    client: 'Send the outbox in pushes of at most 500 changes, as the Grit sync clients do.',
  },
  TOTP_ALREADY_ENABLED: {
    status: 409,
    category: 'conflict',
    area: 'twofactor',
    meaning: 'Two-factor authentication is already on for this account.',
    client: 'Show it as enabled rather than offering setup.',
  },
  TOTP_ERROR: {
    status: 500,
    category: 'server',
    area: 'twofactor',
    meaning: 'The secret, QR code or backup codes could not be produced.',
    client: 'Retry once, then report it. Nothing is half-enabled: setup only counts once confirmed.',
  },
  TOTP_NOT_ENABLED: {
    status: 400,
    category: 'request',
    area: 'twofactor',
    meaning: 'The account has no second factor, so there is nothing to confirm or turn off.',
    client: 'Offer setup instead.',
  },
  TRANSITION_REFUSED: {
    status: 422,
    category: 'state',
    area: 'workflow',
    meaning: 'A transition hook refused the move, and nothing was written.',
    client: 'Show the message: it is the business rule that said no.',
  },
  UNAUTHORIZED: {
    status: 401,
    category: 'auth',
    area: 'core',
    meaning: 'The credentials are missing, expired or not accepted.',
    client: 'Refresh the access token, and sign in again if the refresh is rejected too.',
  },
  UNKNOWN_MODEL: {
    status: 400,
    category: 'request',
    area: 'sync',
    meaning: 'No model is registered under that name.',
    client: 'Read the manifest rather than hard-coding names.',
  },
  UNKNOWN_PROVIDER: {
    status: 404,
    category: 'notfound',
    area: 'auth',
    meaning: 'No social provider is configured under that name.',
    client: 'Only offer the providers the API reports as enabled.',
  },
  UNKNOWN_SETTING: {
    status: 404,
    category: 'notfound',
    area: 'settings',
    meaning: 'No setting is registered under that key.',
    client: 'Read the settings list rather than guessing keys.',
  },
  UPLOAD_ALREADY_RECORDED: {
    status: 409,
    category: 'conflict',
    area: 'uploads',
    meaning: 'This upload has already been recorded, and a key is recorded only once.',
    client: 'Treat the upload as done: fetch it from the uploads list instead of completing it again.',
  },
  UPLOAD_FAILED: {
    status: 500,
    category: 'server',
    area: 'uploads',
    meaning: 'The file reached the API and could not be stored.',
    client: 'Retry once. Nothing was recorded, so there is no half-uploaded row to clean up.',
  },
  UPLOAD_KEY_FORBIDDEN: {
    status: 403,
    category: 'permission',
    area: 'uploads',
    meaning: 'The key being recorded was not presigned for this user. Only keys under the caller\'s own uploads/<user_id>/ prefix are accepted.',
    client: 'Send back the key the presign returned. Another user\'s key, or one outside the uploads, will not be accepted.',
  },
  UPLOAD_NOT_FOUND: {
    status: 404,
    category: 'notfound',
    area: 'uploads',
    meaning: 'Nothing is stored under that key.',
    client: 'Treat it as absent: the key is wrong, or the object was deleted.',
  },
  USER_ERROR: {
    status: 500,
    category: 'server',
    area: 'auth',
    meaning: 'The signed-in account could not be loaded, which means the token outlived its row.',
    client: 'Sign in again. If it repeats, the account data is inconsistent and needs a look.',
  },
  VALIDATION_ERROR: {
    status: 422,
    category: 'request',
    area: 'core',
    meaning: 'The body parsed, and a field in it is missing or not acceptable.',
    client: 'Read error.details: it maps each field to what is wrong with it. Show those against the inputs.',
  },
  VALUE_IN_USE: {
    status: 409,
    category: 'conflict',
    area: 'variants',
    meaning: 'That value is part of existing variants.',
    client: 'Delete those variants first.',
  },
  VERIFY_FAILED: {
    status: 500,
    category: 'server',
    area: 'gdpr',
    meaning: 'The audit chain could not be verified.',
    client: 'Report it. This is the check that says whether the log has been tampered with.',
  },
  VERSION_CONFLICT: {
    status: 409,
    category: 'conflict',
    area: 'core',
    meaning: 'Somebody else changed the row since the version in your If-Match.',
    client: 'The response carries the current version. Re-read, merge, and send the new ETag.',
  },
}

/** The shape of the error envelope every endpoint uses. */
export interface ApiErrorEnvelope {
  error: {
    code: string
    message: string
    /** Per-field messages, on a validation failure. */
    details?: Record<string, string>
  }
}

/** Narrows an unknown string to a documented code. */
export function isApiErrorCode(value: unknown): value is ApiErrorCode {
  return typeof value === 'string' && value in API_ERRORS
}

/**
 * The documented code behind a failure, narrowed to the union.
 *
 * Takes either the response body or the error a client library threw, because
 * both shapes turn up in practice. Returns null when the code is not one of the
 * documented ones, which covers two cases worth keeping apart from a known code:
 * a proxy or gateway answered instead of the API, or one of your own handlers
 * returned a code of its own. For the raw string, including your own codes, use
 * apiErrorCode from ./api.
 */
export function documentedErrorCode(value: unknown): ApiErrorCode | null {
  const asError = value as { response?: { data?: ApiErrorEnvelope } } | undefined
  const code = asError?.response?.data?.error?.code ?? (value as ApiErrorEnvelope | undefined)?.error?.code
  return isApiErrorCode(code) ? code : null
}

/** What the person in front of the screen should do about it. */
export function apiErrorAdvice(code: ApiErrorCode): string {
  return API_ERRORS[code].client
}

/** Codes that are worth retrying unchanged. Everything else needs a decision. */
export function isRetryable(code: ApiErrorCode): boolean {
  const category = API_ERRORS[code].category
  return category === 'server' || category === 'upstream' || category === 'limit' || category === 'disabled'
}

/** Every documented code, sorted, for a table or a test. */
export const API_ERROR_CODES = Object.keys(API_ERRORS).sort() as ApiErrorCode[]
