package settings

// RegisterDefaults declares the settings every Grit project starts with.
//
// Called once at boot, before Init. Add your own beside these, or in the
// package that reads them, which is usually the better home: a setting
// declared next to its only caller cannot be orphaned by a refactor.
func RegisterDefaults() {
	Define(Setting{
		Key:      "app.name",
		Type:     TypeString,
		Label:    "Application name",
		Help:     "Shown in the browser title, emails and the admin header.",
		Group:    "General",
		Default:  "library",
		Validate: NotEmpty(),
		Order:    1,
	})

	Define(Setting{
		Key:     "app.support_email",
		Type:    TypeString,
		Label:   "Support email",
		Help:    "Where users are told to write when something goes wrong.",
		Group:   "General",
		Default: "",
		Order:   2,
	})

	Define(Setting{
		Key:   "app.timezone",
		Type:  TypeString,
		Label: "Default timezone",
		Help:  "An IANA name, for example Africa/Kampala. Used when a user has not chosen one.",
		Group: "General",
		// Per user, because a distributed team is the normal case and a single
		// global timezone quietly shows everybody else the wrong times.
		Scope:   User,
		Default: "UTC",
		Order:   3,
	})

	Define(Setting{
		Key:   "cors.origins",
		Type:  TypeText,
		Label: "Allowed browser origins",
		Help: "One per line. Sites whose browser JavaScript may call this API. " +
			"Add your storefront's domain here. Does not affect servers, mobile " +
			"apps or curl, none of which enforce CORS.",
		Group: "Security",
		// Empty by default, because the value that ships is whatever
		// CORS_ORIGINS says and the resolution order already prefers a stored
		// value over the environment. Writing the env value in as a default
		// here would mean a fresh install had a stored copy that then stopped
		// tracking the environment.
		Default: "",
		Order:   1,
	})

	Define(Setting{
		Key:   "cache.public_ttl_seconds",
		Type:  TypeNumber,
		Label: "Public response cache (seconds)",
		Help: "How long a public endpoint's response is reused. Read at startup, " +
			"so a change here needs a restart. 0 uses the 60 second default.",
		Group:    "Performance",
		Default:  "60",
		Validate: Between(0, 86400),
		Order:    1,
	})

	Define(Setting{
		Key:     "notifications.email_enabled",
		Type:    TypeBool,
		Label:   "Send notification emails",
		Help:    "Turn off to suppress all outbound notification mail without changing any code.",
		Group:   "Notifications",
		Default: "true",
		Order:   1,
	})
}
