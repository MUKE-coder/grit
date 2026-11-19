package scaffold

// internal/notify, which runs something at most once for a key.
//
// From grit#179: a 278-line billing sweep whose every job did the same four
// things, one of which was "send the email at most once, through a dedup row".
// Getting that wrong emails a customer the same overdue reminder on every run
// of the sweep, which is a worse failure than not sending it at all, and it is
// worth getting right once here rather than in every project that needs it.

// apiNotifyGo emits internal/notify/notify.go.
func apiNotifyGo() string { return tmpl("api/notify/notify.go") }

// apiNotifyTestGo emits internal/notify/notify_test.go.
func apiNotifyTestGo() string { return tmpl("api/notify/notify_test.go") }

// apiNotificationClaimModelGo emits internal/models/notification_claim.go.
func apiNotificationClaimModelGo() string { return tmpl("api/models/notification_claim.go") }

// notifyModelRegistration is the claims table in the AutoMigrate list.
//
// Registered like every other table, so grit migrate creates it and a backup
// includes it. A backup without this table restores a deployment that has
// forgotten every notification it has ever sent, and the first sweep after the
// restore sends all of them again.
const notifyModelRegistration = `		// The record that a notification has already been sent, keyed by
		// whatever the caller claimed. Here like every other table, so
		// AutoMigrate creates it and a backup includes it: restoring without
		// it means the next sweep re-sends everything.
		&NotificationClaim{},
`

// notifyModelAnchor is the line the registration goes before.
const notifyModelAnchor = "\t\t// grit:models\n"
