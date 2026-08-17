package scaffold

// apiSettingsHandlerGo emits internal/handlers/settings.go.
func apiSettingsHandlerGo() string {
	return `package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"{{MODULE}}/internal/authz"
	"{{MODULE}}/internal/events"
	"{{MODULE}}/internal/settings"
)

// SettingsHandler serves the declared settings and their current values.
type SettingsHandler struct {
	DB *gorm.DB
}

// settingView is one setting as the admin sees it: the declaration plus the
// resolved value.
type settingView struct {
	Key        string             ` + "`" + `json:"key"` + "`" + `
	Type       settings.Kind      ` + "`" + `json:"type"` + "`" + `
	Label      string             ` + "`" + `json:"label"` + "`" + `
	Help       string             ` + "`" + `json:"help,omitempty"` + "`" + `
	Group      string             ` + "`" + `json:"group"` + "`" + `
	Value      string             ` + "`" + `json:"value"` + "`" + `
	Default    string             ` + "`" + `json:"default"` + "`" + `
	Options    []settings.Option  ` + "`" + `json:"options,omitempty"` + "`" + `
	Scope      settings.Scope     ` + "`" + `json:"scope"` + "`" + `
	Permission string             ` + "`" + `json:"permission,omitempty"` + "`" + `
	// FromEnv marks a value that is currently coming from an environment
	// variable because nothing is stored. The admin shows where it came from,
	// and saving still works: the stored value then wins, which is the
	// precedence the store resolves with.
	FromEnv bool ` + "`" + `json:"from_env"` + "`" + `
	// EnvVar names the variable, so the admin can say "currently from
	// APP_NAME" rather than leaving the reader to guess.
	EnvVar   string ` + "`" + `json:"env_var,omitempty"` + "`" + `
	Editable bool   ` + "`" + `json:"editable"` + "`" + `
}

// List handles GET /api/settings.
func (h *SettingsHandler) List(c *gin.Context) {
	ctx := settingsContext(c)
	grants := grantsOf(c)

	declared := settings.All()
	out := make([]settingView, 0, len(declared))
	for _, s := range declared {
		value := settings.String(ctx, s.Key)
		// Only "from the environment" when that is actually where this value
		// came from: an env var that a stored value is already overriding is
		// not what the reader is looking at.
		envValue, hasEnv := s.EnvOverride()
		fromEnv := hasEnv && value == envValue
		if s.Type == settings.TypeSecret && value != "" {
			// Never send a secret back. It can be replaced, not read.
			value = "********"
		}
		out = append(out, settingView{
			Key: s.Key, Type: s.Type, Label: s.Label, Help: s.Help,
			Group: s.Group, Value: value, Default: s.Default,
			Options: s.Options, Scope: s.Scope, Permission: s.Permission,
			FromEnv:  fromEnv,
			EnvVar:   envVarIf(fromEnv, s.Env),
			Editable: s.Permission == "" || authz.Granted(grants, s.Permission),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": out, "groups": settings.Groups()})
}

// UpdateSettingsRequest is a batch of changes, so an admin page saves once.
type UpdateSettingsRequest struct {
	Values map[string]string ` + "`" + `json:"values" binding:"required"` + "`" + `
	// Scope defaults to global. "user" or "tenant" write an override for the
	// caller's own user or organisation.
	Scope string ` + "`" + `json:"scope,omitempty"` + "`" + `
}

// Update handles PUT /api/settings.
//
// Validates everything before writing anything. A batch that half-applies
// leaves the admin page showing a mix of saved and rejected values with no
// way to tell which is which.
func (h *SettingsHandler) Update(c *gin.Context) {
	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code": "INVALID_BODY", "message": err.Error(),
		}})
		return
	}

	store := settings.Default()
	if store == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "SETTINGS_UNAVAILABLE", "message": "the settings store is not initialised",
		}})
		return
	}

	scope := settings.Scope(req.Scope)
	if scope == "" {
		scope = settings.Global
	}
	scopeID := ""
	switch scope {
	case settings.User:
		scopeID = stringClaim(c, "user_id")
	case settings.Tenant:
		scopeID = stringClaim(c, "tenant_id")
	}
	if scope != settings.Global && scopeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"code":    "NO_SCOPE",
			"message": "cannot write a " + string(scope) + " override without a " + string(scope),
		}})
		return
	}

	grants := grantsOf(c)
	details := map[string]string{}

	for key, value := range req.Values {
		declared, ok := settings.Get(key)
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{
				"code": "UNKNOWN_SETTING", "message": "no setting named " + key + " is declared",
			}})
			return
		}
		if declared.Permission != "" && !authz.Granted(grants, declared.Permission) {
			c.JSON(http.StatusForbidden, gin.H{"error": gin.H{
				"code":    "FORBIDDEN",
				"message": "you do not have permission to change " + declared.Label,
			}})
			return
		}
		if err := declared.Parse(value); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{
				"code": "VALIDATION_ERROR", "message": err.Error(),
				"details": gin.H{key: err.Error()},
			}})
			return
		}
		details[key] = value
	}

	actor := stringClaim(c, "user_id")
	for key, value := range details {
		if err := store.Set(key, scope, scopeID, value, actor); err != nil {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{
				"code": "SETTING_REJECTED", "message": err.Error(),
			}})
			return
		}
	}

	// Settings changes are the kind of thing an audit trail exists for: they
	// alter behaviour everywhere and leave no other trace. Secret values are
	// named but never recorded.
	changed := make([]string, 0, len(details))
	for key := range details {
		changed = append(changed, key)
	}
	events.Emitted(c, "settings", "Setting", "updated", string(scope), string(scope),
		"changed "+joinKeys(changed), nil, nil)

	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"updated": len(details), "scope": scope},
		"message": "Settings saved",
	})
}

// Reset handles DELETE /api/settings/:key, removing an override so the level
// below applies again.
func (h *SettingsHandler) Reset(c *gin.Context) {
	key := c.Param("key")
	declared, ok := settings.Get(key)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
			"code": "UNKNOWN_SETTING", "message": "no setting named " + key + " is declared",
		}})
		return
	}
	if declared.Permission != "" && !authz.Granted(grantsOf(c), declared.Permission) {
		c.JSON(http.StatusForbidden, gin.H{"error": gin.H{
			"code": "FORBIDDEN", "message": "you do not have permission to change " + declared.Label,
		}})
		return
	}

	store := settings.Default()
	if store == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "SETTINGS_UNAVAILABLE", "message": "the settings store is not initialised",
		}})
		return
	}

	scope := settings.Scope(c.DefaultQuery("scope", string(settings.Global)))
	scopeID := ""
	switch scope {
	case settings.User:
		scopeID = stringClaim(c, "user_id")
	case settings.Tenant:
		scopeID = stringClaim(c, "tenant_id")
	}

	if err := store.Unset(key, scope, scopeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"code": "INTERNAL_ERROR", "message": err.Error(),
		}})
		return
	}

	events.Emitted(c, "settings", "Setting", "reset", key, declared.Label, "", nil, nil)
	c.JSON(http.StatusOK, gin.H{
		"data":    gin.H{"key": key, "value": settings.String(settingsContext(c), key)},
		"message": declared.Label + " reset to its default",
	})
}

// settingsContext carries the caller's user and tenant into resolution, so a
// read returns what this person should see rather than the global value.
func settingsContext(c *gin.Context) context.Context {
	ctx := c.Request.Context()
	if uid := stringClaim(c, "user_id"); uid != "" {
		ctx = settings.WithUser(ctx, uid)
	}
	if tid := stringClaim(c, "tenant_id"); tid != "" {
		ctx = settings.WithTenant(ctx, tid)
	}
	return ctx
}

func stringClaim(c *gin.Context, key string) string {
	if v, ok := c.Get(key); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func grantsOf(c *gin.Context) []string {
	if v, ok := c.Get("user_grants"); ok {
		if list, ok := v.([]string); ok {
			return list
		}
	}
	return nil
}

func envVarIf(fromEnv bool, name string) string {
	if fromEnv {
		return name
	}
	return ""
}

func joinKeys(keys []string) string {
	out := ""
	for i, k := range keys {
		if i > 0 {
			out += ", "
		}
		out += k
	}
	return out
}
`
}

// apiSettingsDefaultsGo emits internal/settings/defaults.go: the declarations
// every project starts with.
//
// Real ones rather than a commented-out example, because a settings page with
// nothing on it teaches nobody the shape, and these four are wanted by
// essentially every application.
func apiSettingsDefaultsGo() string {
	return `package settings

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
		Default:  "{{PROJECT}}",
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
`
}
