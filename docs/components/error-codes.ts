// Generated from internal/errorcodes by tools/errorcodes. DO NOT EDIT.
//
// Run: go run ./tools/errorcodes
//
// The same catalogue generates internal/respond/codes.go and
// packages/shared/types/errors.ts in every project, which is what makes this
// table the documentation of what the API actually returns rather than a list
// somebody kept up by hand.

export interface ErrorCodeRow {
  code: string
  status: number
  category: string
  /** What happened, from the caller's side. */
  meaning: string
  /** What the caller should do about it. */
  client: string
}

export interface ErrorCodeArea {
  /** The key the API uses. */
  area: string
  /** The heading for a reader. */
  label: string
  codes: ErrorCodeRow[]
}

export const errorCodeAreas: ErrorCodeArea[] = [
  {
    area: 'core',
    label: 'Every endpoint',
    codes: [
      {
        code: 'BAD_REQUEST',
        status: 400,
        category: 'request',
        meaning: 'The request could not be understood at all.',
        client: 'Fix the request. Retrying the same one will fail the same way.',
      },
      {
        code: 'INVALID_BODY',
        status: 400,
        category: 'request',
        meaning: 'The body was not valid JSON, or was not the shape this endpoint reads.',
        client: 'Send a JSON body matching the documented request type.',
      },
      {
        code: 'READ_BODY_FAILED',
        status: 400,
        category: 'request',
        meaning: 'The body could not be read to the end.',
        client: 'Retry. If it keeps happening, the connection is dropping or the body is larger than the server accepts.',
      },
      {
        code: 'VALIDATION_ERROR',
        status: 422,
        category: 'request',
        meaning: 'The body parsed, and a field in it is missing or not acceptable.',
        client: 'Read error.details: it maps each field to what is wrong with it. Show those against the inputs.',
      },
      {
        code: 'PAYLOAD_TOO_LARGE',
        status: 413,
        category: 'request',
        meaning: 'The request body is larger than the server accepts.',
        client: 'Send less, or upload the file directly to storage with a presigned URL.',
      },
      {
        code: 'UNAUTHORIZED',
        status: 401,
        category: 'auth',
        meaning: 'The credentials are missing, expired or not accepted.',
        client: 'Refresh the access token, and sign in again if the refresh is rejected too.',
      },
      {
        code: 'MISSING_TOKEN',
        status: 401,
        category: 'auth',
        meaning: 'No token was sent where one is required.',
        client: 'Send the access token as Authorization: Bearer <token>.',
      },
      {
        code: 'INVALID_TOKEN',
        status: 401,
        category: 'auth',
        meaning: 'The token is malformed, expired, or was issued for something else.',
        client: 'Refresh it. A refresh token that fails this way has been used already or revoked: sign in again.',
      },
      {
        code: 'INVALID_LINK',
        status: 400,
        category: 'request',
        meaning: 'A one-time link, such as email verification or password reset, is wrong or has expired.',
        client: 'Offer to send a new link. This is not a sign-in failure: the caller has no credentials to fix, which is why it is 400 and not 401.',
      },
      {
        code: 'SESSION_REVOKED',
        status: 401,
        category: 'auth',
        meaning: 'The session behind this token was signed out, on this device or another.',
        client: 'Sign in again. Do not retry with the same refresh token: reuse is what revoked it.',
      },
      {
        code: 'CSRF_INVALID',
        status: 403,
        category: 'permission',
        meaning: 'The CSRF token is missing or does not match the cookie.',
        client: 'Read the CSRF cookie and send it back in the header on every unsafe request.',
      },
      {
        code: 'FORBIDDEN',
        status: 403,
        category: 'permission',
        meaning: 'The caller is signed in and this action is not theirs to take.',
        client: 'Do not retry. Hide the action rather than letting it fail, if the role is known to the client.',
      },
      {
        code: 'NOT_FOUND',
        status: 404,
        category: 'notfound',
        meaning: 'No such row, or none this caller is allowed to see.',
        client: 'Treat it as absent. On an owned resource this is also the answer for somebody else\'s row, on purpose: a wrong guess cannot be told from a right one.',
      },
      {
        code: 'CONFLICT',
        status: 409,
        category: 'conflict',
        meaning: 'The write collided with the state already there, such as a unique column.',
        client: 'Re-read, show what is there, and let the person decide.',
      },
      {
        code: 'VERSION_CONFLICT',
        status: 409,
        category: 'conflict',
        meaning: 'Somebody else changed the row since the version in your If-Match.',
        client: 'The response carries the current version. Re-read, merge, and send the new ETag.',
      },
      {
        code: 'RATE_LIMITED',
        status: 429,
        category: 'limit',
        meaning: 'Too many requests from this caller.',
        client: 'Back off. Honour Retry-After if it is present rather than retrying immediately.',
      },
      {
        code: 'INTERNAL_ERROR',
        status: 500,
        category: 'server',
        meaning: 'A fault on the server. The message is deliberately vague; the detail is in the server log.',
        client: 'Retry once, then report it. Nothing the client changes will help.',
      },
      {
        code: 'MAINTENANCE',
        status: 503,
        category: 'disabled',
        meaning: 'The API is in maintenance mode and is refusing everything.',
        client: 'Retry later. Show a maintenance state rather than an error.',
      },
      {
        code: 'PERSIST_FAILED',
        status: 500,
        category: 'server',
        meaning: 'The change was accepted and could not be written.',
        client: 'Retry once. Treat the write as not having happened.',
      },
    ],
  },
  {
    area: 'tenancy',
    label: 'Organizations (the multitenant plugin)',
    codes: [
      {
        code: 'NO_ORGANIZATION',
        status: 400,
        category: 'request',
        meaning: 'The row belongs to an organization and the request has no active one: the caller belongs to none, or to several and named neither.',
        client: 'Send the active organization as X-Organization-ID. If the caller belongs to no organization, they cannot read this at all: put them in one, or send them somewhere that does not need one.',
      },
    ],
  },
  {
    area: 'auth',
    label: 'Sign-in and accounts',
    codes: [
      {
        code: 'INVALID_CREDENTIALS',
        status: 401,
        category: 'auth',
        meaning: 'The email and password do not match an account.',
        client: 'Say only that the details are wrong: which of the two it was is deliberately not reported.',
      },
      {
        code: 'INVALID_PASSWORD',
        status: 401,
        category: 'auth',
        meaning: 'The password given for a confirmation step is not correct.',
        client: 'Ask again. This is the re-authentication prompt, not a sign-in.',
      },
      {
        code: 'EMAIL_EXISTS',
        status: 409,
        category: 'conflict',
        meaning: 'An account already has that address.',
        client: 'Offer sign-in or password reset rather than registration.',
      },
      {
        code: 'EMAIL_NOT_VERIFIED',
        status: 403,
        category: 'permission',
        meaning: 'The account exists and its address has not been confirmed.',
        client: 'Send them to the verification flow, and offer to resend the link.',
      },
      {
        code: 'ACCOUNT_DISABLED',
        status: 403,
        category: 'permission',
        meaning: 'The account has been deactivated.',
        client: 'Do not retry. This needs an administrator, not a different password.',
      },
      {
        code: 'ACCOUNT_LOCKED',
        status: 429,
        category: 'limit',
        meaning: 'Too many failed attempts, so the account is locked for a while.',
        client: 'Show the wait, and offer password reset. Retrying sooner extends nothing but the lock.',
      },
      {
        code: 'ALREADY_VERIFIED',
        status: 400,
        category: 'request',
        meaning: 'The address on this link is already confirmed.',
        client: 'Treat it as success and continue to sign-in.',
      },
      {
        code: 'SOCIAL_AUTH_ONLY',
        status: 400,
        category: 'request',
        meaning: 'The account signs in with a social provider and has no password.',
        client: 'Offer the provider button instead of the password form.',
      },
      {
        code: 'NO_PASSWORD',
        status: 400,
        category: 'request',
        meaning: 'The account has no password set, so it cannot be confirmed with one.',
        client: 'Send them through set-a-password first.',
      },
      {
        code: 'UNKNOWN_PROVIDER',
        status: 404,
        category: 'notfound',
        meaning: 'No social provider is configured under that name.',
        client: 'Only offer the providers the API reports as enabled.',
      },
      {
        code: 'TOKEN_ERROR',
        status: 500,
        category: 'server',
        meaning: 'The access and refresh tokens could not be issued.',
        client: 'Retry once. The credentials were accepted, so do not ask for them again.',
      },
      {
        code: 'USER_ERROR',
        status: 500,
        category: 'server',
        meaning: 'The signed-in account could not be loaded, which means the token outlived its row.',
        client: 'Sign in again. If it repeats, the account data is inconsistent and needs a look.',
      },
    ],
  },
  {
    area: 'webhooks',
    label: 'Incoming webhooks',
    codes: [
      {
        code: 'INVALID_SIGNATURE',
        status: 401,
        category: 'auth',
        meaning: 'The webhook signature does not match the body and the shared secret.',
        client: 'Sign the exact bytes sent, with the secret for this endpoint. A reformatted body will not verify.',
      },
    ],
  },
  {
    area: 'twofactor',
    label: 'Two-factor authentication',
    codes: [
      {
        code: 'INVALID_TOTP_CODE',
        status: 401,
        category: 'auth',
        meaning: 'The six-digit code is wrong or has expired.',
        client: 'Let them try the next code. Clock drift on the device is the usual cause of repeated failures.',
      },
      {
        code: 'INVALID_BACKUP_CODE',
        status: 401,
        category: 'auth',
        meaning: 'That backup code is wrong, or has been used.',
        client: 'Each code works once. Offer the remaining count the status endpoint reports.',
      },
      {
        code: 'INVALID_PENDING_TOKEN',
        status: 401,
        category: 'auth',
        meaning: 'The short-lived token between password and second factor is expired or unknown.',
        client: 'Start the sign-in again from the password step.',
      },
      {
        code: 'TOTP_ALREADY_ENABLED',
        status: 409,
        category: 'conflict',
        meaning: 'Two-factor authentication is already on for this account.',
        client: 'Show it as enabled rather than offering setup.',
      },
      {
        code: 'TOTP_NOT_ENABLED',
        status: 400,
        category: 'request',
        meaning: 'The account has no second factor, so there is nothing to confirm or turn off.',
        client: 'Offer setup instead.',
      },
      {
        code: 'TOTP_ERROR',
        status: 500,
        category: 'server',
        meaning: 'The secret, QR code or backup codes could not be produced.',
        client: 'Retry once, then report it. Nothing is half-enabled: setup only counts once confirmed.',
      },
    ],
  },
  {
    area: 'passkeys',
    label: 'Passkeys',
    codes: [
      {
        code: 'PASSKEYS_NOT_CONFIGURED',
        status: 501,
        category: 'disabled',
        meaning: 'This deployment has no passkey configuration, so the endpoints are inert.',
        client: 'Hide passkey buttons unless the API reports the feature as available.',
      },
      {
        code: 'PASSKEY_REJECTED',
        status: 410,
        category: 'state',
        meaning: 'The challenge no longer exists: it expired, or it was already answered.',
        client: 'Start the ceremony again. Do not retry the same assertion.',
      },
    ],
  },
  {
    area: 'recovery',
    label: 'Account recovery',
    codes: [
      {
        code: 'INVALID_RECOVERY_ADDRESS',
        status: 422,
        category: 'request',
        meaning: 'The recovery email or phone number is not usable.',
        client: 'Validate the format before sending, and show which one was rejected.',
      },
      {
        code: 'INVALID_CODE',
        status: 422,
        category: 'request',
        meaning: 'The recovery code is wrong or has expired.',
        client: 'Offer to send a new one rather than retrying the same code.',
      },
      {
        code: 'SMS_NOT_CONFIGURED',
        status: 501,
        category: 'disabled',
        meaning: 'No SMS provider is configured, so phone recovery cannot run here.',
        client: 'Offer email recovery instead.',
      },
      {
        code: 'SMS_FAILED',
        status: 502,
        category: 'upstream',
        meaning: 'The SMS provider refused or failed to send.',
        client: 'Retry once, then offer email. The number may be unreachable.',
      },
    ],
  },
  {
    area: 'apikeys',
    label: 'API keys',
    codes: [
      {
        code: 'API_KEY_REQUIRED',
        status: 401,
        category: 'auth',
        meaning: 'The route is key-guarded and no key was sent.',
        client: 'Send the key in the documented header. A user token is not a substitute here.',
      },
      {
        code: 'INVALID_API_KEY',
        status: 401,
        category: 'auth',
        meaning: 'The key is unknown, revoked or expired.',
        client: 'Issue a new key. Do not retry with the same one.',
      },
      {
        code: 'ENDPOINT_NOT_ALLOWED',
        status: 403,
        category: 'permission',
        meaning: 'The key is valid and is not allowed to call this endpoint.',
        client: 'Widen the key\'s endpoint list, or call it with one that may.',
      },
      {
        code: 'ORIGIN_NOT_ALLOWED',
        status: 403,
        category: 'permission',
        meaning: 'The request\'s Origin is not on the key\'s allowlist.',
        client: 'Add the origin to the key, rather than relaxing CORS for everybody.',
      },
      {
        code: 'PUBLISHABLE_KEY_NOT_ALLOWED',
        status: 403,
        category: 'permission',
        meaning: 'A publishable key was used where only a secret key is accepted.',
        client: 'Call this from the server with the secret key. A publishable key is public by design.',
      },
    ],
  },
  {
    area: 'uploads',
    label: 'Uploads and storage',
    codes: [
      {
        code: 'INVALID_FILE',
        status: 400,
        category: 'request',
        meaning: 'No file was attached, or it could not be read.',
        client: 'Send multipart form data with the documented field name.',
      },
      {
        code: 'INVALID_FILE_TYPE',
        status: 400,
        category: 'request',
        meaning: 'The file\'s type is not accepted for this field.',
        client: 'Check the type client-side before uploading, and say which types are allowed.',
      },
      {
        code: 'FILE_TOO_LARGE',
        status: 400,
        category: 'request',
        meaning: 'The file is larger than this endpoint accepts.',
        client: 'Show the limit before the upload starts rather than after it finishes.',
      },
      {
        code: 'UPLOAD_FAILED',
        status: 500,
        category: 'server',
        meaning: 'The file reached the API and could not be stored.',
        client: 'Retry once. Nothing was recorded, so there is no half-uploaded row to clean up.',
      },
      {
        code: 'UPLOAD_NOT_FOUND',
        status: 404,
        category: 'notfound',
        meaning: 'Nothing is stored under that key.',
        client: 'Treat it as absent: the key is wrong, or the object was deleted.',
      },
      {
        code: 'PRESIGN_FAILED',
        status: 500,
        category: 'server',
        meaning: 'A presigned upload URL could not be produced.',
        client: 'Retry once, then fall back to uploading through the API.',
      },
      {
        code: 'STORAGE_UNAVAILABLE',
        status: 503,
        category: 'disabled',
        meaning: 'Object storage is not configured here, or is not answering.',
        client: 'Retry later. Nothing the client sends will fix it.',
      },
    ],
  },
  {
    area: 'ai',
    label: 'AI gateway',
    codes: [
      {
        code: 'AI_UNAVAILABLE',
        status: 503,
        category: 'disabled',
        meaning: 'No AI provider is configured in this deployment.',
        client: 'Hide AI features unless the API reports one as available.',
      },
      {
        code: 'AI_UNAUTHORIZED',
        status: 502,
        category: 'upstream',
        meaning: 'The provider rejected the server\'s API key.',
        client: 'Nothing for the client to do. The key on the server is wrong or out of credit.',
      },
      {
        code: 'AI_FORBIDDEN',
        status: 502,
        category: 'upstream',
        meaning: 'The provider refused this request, usually its own policy.',
        client: 'Do not retry the same prompt unchanged.',
      },
      {
        code: 'AI_RATE_LIMITED',
        status: 429,
        category: 'limit',
        meaning: 'The provider is rate-limiting this deployment.',
        client: 'Back off and retry with a delay. Queue rather than loop.',
      },
      {
        code: 'AI_MODEL_NOT_FOUND',
        status: 502,
        category: 'upstream',
        meaning: 'The provider does not know the model that was asked for.',
        client: 'Pick a model the API lists. A model name can disappear without notice.',
      },
      {
        code: 'AI_ERROR',
        status: 502,
        category: 'upstream',
        meaning: 'The provider failed in a way that is not one of the above.',
        client: 'Retry once. The message carries what the provider said.',
      },
    ],
  },
  {
    area: 'jobs',
    label: 'Background jobs',
    codes: [
      {
        code: 'REDIS_UNAVAILABLE',
        status: 503,
        category: 'disabled',
        meaning: 'Redis is not configured or not reachable, so the queue cannot be read.',
        client: 'Retry later. Jobs, cache and cron all depend on it.',
      },
      {
        code: 'INVALID_STATUS',
        status: 400,
        category: 'request',
        meaning: 'That queue state is not one the endpoint accepts.',
        client: 'Use one of the documented states.',
      },
      {
        code: 'RETRY_FAILED',
        status: 500,
        category: 'server',
        meaning: 'The job could not be re-queued.',
        client: 'Retry once. The job is still where it was.',
      },
      {
        code: 'CLEAR_FAILED',
        status: 500,
        category: 'server',
        meaning: 'The queue could not be cleared.',
        client: 'Retry once, then look at the Redis connection.',
      },
    ],
  },
  {
    area: 'import',
    label: 'CSV import',
    codes: [
      {
        code: 'JOB_ERROR',
        status: 500,
        category: 'server',
        meaning: 'The import job could not be started.',
        client: 'Retry once. Nothing was imported.',
      },
      {
        code: 'TEMP_ERROR',
        status: 500,
        category: 'server',
        meaning: 'The upload could not be buffered to disk before importing.',
        client: 'Retry once. Check free disk on the server if it repeats.',
      },
      {
        code: 'INVALID_CSV',
        status: 400,
        category: 'request',
        meaning: 'The file is not readable as CSV.',
        client: 'Check the delimiter, the quoting and that the header row matches the template.',
      },
    ],
  },
  {
    area: 'backups',
    label: 'Backups and restore',
    codes: [
      {
        code: 'NOT_AVAILABLE',
        status: 400,
        category: 'request',
        meaning: 'That backup cannot be downloaded: it is not finished, or it is not stored here.',
        client: 'Re-read the backup\'s status before offering a download link.',
      },
      {
        code: 'INVALID_SCHEDULE',
        status: 400,
        category: 'request',
        meaning: 'The schedule is not a cron expression this API accepts.',
        client: 'Validate the expression client-side, or offer fixed choices.',
      },
      {
        code: 'EXTRACT_FAILED',
        status: 400,
        category: 'request',
        meaning: 'The archive could not be opened or does not hold what a restore needs.',
        client: 'Upload an archive this API produced. A re-zipped one usually fails here.',
      },
    ],
  },
  {
    area: 'gdpr',
    label: 'GDPR and the audit log',
    codes: [
      {
        code: 'EXPORT_FAILED',
        status: 500,
        category: 'server',
        meaning: 'The subject-access export could not be assembled.',
        client: 'Retry once, then report it: this is a request with a legal clock on it.',
      },
      {
        code: 'ERASE_FAILED',
        status: 500,
        category: 'server',
        meaning: 'The erasure did not complete.',
        client: 'Do not assume anything was erased. Retry, and check the audit log.',
      },
      {
        code: 'SELF_ERASE',
        status: 400,
        category: 'request',
        meaning: 'An account cannot erase itself through this endpoint.',
        client: 'Have another administrator run it, so the action has an actor who remains.',
      },
      {
        code: 'QUERY_FAILED',
        status: 500,
        category: 'server',
        meaning: 'The audit query failed.',
        client: 'Retry once, then report it.',
      },
      {
        code: 'VERIFY_FAILED',
        status: 500,
        category: 'server',
        meaning: 'The audit chain could not be verified.',
        client: 'Report it. This is the check that says whether the log has been tampered with.',
      },
      {
        code: 'RESEAL_REFUSED',
        status: 409,
        category: 'conflict',
        meaning: 'The audit chain was not broken where the reseal claimed, so nothing was resealed.',
        client: 'Verify the chain again and reseal from the entry the verification names.',
      },
    ],
  },
  {
    area: 'settings',
    label: 'Settings',
    codes: [
      {
        code: 'UNKNOWN_SETTING',
        status: 404,
        category: 'notfound',
        meaning: 'No setting is registered under that key.',
        client: 'Read the settings list rather than guessing keys.',
      },
      {
        code: 'SETTING_REJECTED',
        status: 422,
        category: 'request',
        meaning: 'The value did not pass the setting\'s own validation.',
        client: 'Show the message against the field: it comes from the setting\'s rule.',
      },
      {
        code: 'SETTINGS_UNAVAILABLE',
        status: 500,
        category: 'server',
        meaning: 'The settings store could not be read or written.',
        client: 'Retry once, then report it.',
      },
      {
        code: 'NO_SCOPE',
        status: 400,
        category: 'request',
        meaning: 'The request did not say which scope to act in.',
        client: 'Send the scope the settings list gives for that key.',
      },
    ],
  },
  {
    area: 'variants',
    label: 'Product variants',
    codes: [
      {
        code: 'OPTION_IN_USE',
        status: 409,
        category: 'conflict',
        meaning: 'Variants are built on this option, so it cannot be removed.',
        client: 'Clear the combinations that use it first, and say so rather than failing silently.',
      },
      {
        code: 'VALUE_IN_USE',
        status: 409,
        category: 'conflict',
        meaning: 'That value is part of existing variants.',
        client: 'Delete those variants first.',
      },
      {
        code: 'CANNOT_GENERATE',
        status: 422,
        category: 'request',
        meaning: 'The combinations could not be generated from the options given.',
        client: 'Check that every option has at least one value.',
      },
    ],
  },
  {
    area: 'review',
    label: 'Access reviews',
    codes: [
      {
        code: 'REVIEW_CLOSED',
        status: 400,
        category: 'state',
        meaning: 'The review is closed, so its decisions cannot change.',
        client: 'Open a new review rather than editing a closed one.',
      },
      {
        code: 'REVIEW_INCOMPLETE',
        status: 400,
        category: 'state',
        meaning: 'Some items still have no decision, so the review cannot be completed.',
        client: 'Show which items are outstanding.',
      },
      {
        code: 'ITEM_LOCKED',
        status: 400,
        category: 'state',
        meaning: 'That item already has a decision and will not take another.',
        client: 'Re-read the review before submitting again.',
      },
      {
        code: 'INVALID_DECISION',
        status: 400,
        category: 'request',
        meaning: 'That is not a decision this review accepts.',
        client: 'Use one of the documented decisions.',
      },
      {
        code: 'CANNOT_COMPLETE',
        status: 400,
        category: 'state',
        meaning: 'The review cannot be completed in its current state.',
        client: 'Re-read it: the message says what is missing.',
      },
    ],
  },
  {
    area: 'forms',
    label: 'Public forms',
    codes: [
      {
        code: 'PASSWORD_REQUIRED',
        status: 401,
        category: 'auth',
        meaning: 'The shared form is password-protected.',
        client: 'Prompt for the form\'s password and send it with the submission.',
      },
      {
        code: 'SUBMISSION_FAILED',
        status: 400,
        category: 'request',
        meaning: 'The submission was refused: a field, a file or the form\'s own rules.',
        client: 'Show the message. It is written for the person filling the form in.',
      },
    ],
  },
  {
    area: 'sync',
    label: 'Offline sync',
    codes: [
      {
        code: 'MISSING_MODEL',
        status: 400,
        category: 'request',
        meaning: 'The request did not name a model to sync.',
        client: 'Send the model name the sync manifest lists.',
      },
      {
        code: 'UNKNOWN_MODEL',
        status: 400,
        category: 'request',
        meaning: 'No model is registered under that name.',
        client: 'Read the manifest rather than hard-coding names.',
      },
      {
        code: 'NOT_SYNCABLE',
        status: 400,
        category: 'request',
        meaning: 'That model is not exposed to offline sync.',
        client: 'Only sync models the manifest marks as syncable.',
      },
      {
        code: 'INVALID_SINCE',
        status: 400,
        category: 'request',
        meaning: 'The since parameter is not an RFC3339 timestamp.',
        client: 'Send the cursor the last sync returned, unchanged.',
      },
    ],
  },
  {
    area: 'workflow',
    label: 'Workflows',
    codes: [
      {
        code: 'INVALID_TRANSITION',
        status: 422,
        category: 'state',
        meaning: 'That move is not declared in the workflow for this status.',
        client: 'Offer only the transitions the API lists for the current status.',
      },
      {
        code: 'TRANSITION_REFUSED',
        status: 422,
        category: 'state',
        meaning: 'A transition hook refused the move, and nothing was written.',
        client: 'Show the message: it is the business rule that said no.',
      },
    ],
  },
  {
    area: 'tree',
    label: 'Trees',
    codes: [
      {
        code: 'INVALID_MOVE',
        status: 422,
        category: 'request',
        meaning: 'That move would put a node inside its own subtree, or under a parent that cannot hold it.',
        client: 'Refuse the drop in the UI rather than sending it.',
      },
    ],
  },
  {
    area: 'charts',
    label: 'Charts and statistics',
    codes: [
      {
        code: 'CHART_FAILED',
        status: 400,
        category: 'request',
        meaning: 'The chart could not be built from those parameters.',
        client: 'Check the resource and preset against the ones the dashboard offers.',
      },
      {
        code: 'STATS_FAILED',
        status: 400,
        category: 'request',
        meaning: 'The statistics could not be computed. This one deliberately conflates an unknown resource with a failed query, so a dashboard widget can render an error state instead of crashing.',
        client: 'Render the widget\'s error state. The message says which of the two it was.',
      },
    ],
  },
  {
    area: 'pdf',
    label: 'PDF rendering',
    codes: [
      {
        code: 'PDF_ERROR',
        status: 500,
        category: 'server',
        meaning: 'The PDF could not be rendered.',
        client: 'Retry once. The record itself is unaffected.',
      },
    ],
  },
  {
    area: 'observability',
    label: 'Security and metrics dashboards',
    codes: [
      {
        code: 'DB_ERROR',
        status: 500,
        category: 'server',
        meaning: 'A query behind a dashboard failed.',
        client: 'Retry once, then report it.',
      },
      {
        code: 'SENTINEL_OFF',
        status: 503,
        category: 'disabled',
        meaning: 'Sentinel is not enabled in this deployment, so there is nothing to report.',
        client: 'Hide the security dashboard unless the API says it is on.',
      },
      {
        code: 'PULSE_OFF',
        status: 503,
        category: 'disabled',
        meaning: 'Pulse is not enabled in this deployment.',
        client: 'Hide the metrics dashboard unless the API says it is on.',
      },
    ],
  },
]

/** How many codes the API documents. Shown on the page, so it cannot be stale. */
export const errorCodeCount = 106

/** Every row, flattened, for searching and for a test that checks coverage. */
export const errorCodes: ErrorCodeRow[] = errorCodeAreas.flatMap((area) => area.codes)
