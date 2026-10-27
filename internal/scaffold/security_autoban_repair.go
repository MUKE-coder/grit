package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// The scheduled ban sweep, delivered to a project that already exists.
//
// The sweep itself is a new file and the handler is refreshed wholesale, so
// neither needs an anchor. What needs one is everything that wires them in:
// the configuration it reads, the routes the dashboard calls, the cron
// registration, and the two call sites of cron.New, which gained a parameter.
//
// Each repair is anchored on what Grit wrote. A file somebody has made their
// own is left alone and named, because a half-applied wiring change is worse
// than none: it does not compile, and the message says so with the edit to
// make by hand.

func repairSecurityAutoBan(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	for _, f := range []struct {
		path string
		fn   func(string) (string, []string, []string)
	}{
		{filepath.Join(apiRoot, "internal", "config", "config.go"), repairAutoBanConfigSource},
		{filepath.Join(apiRoot, "internal", "routes", "routes.go"), repairAutoBanRoutesSource},
		{filepath.Join(apiRoot, "internal", "cron", "cron.go"), repairAutoBanCronSource},
		{filepath.Join(apiRoot, "internal", "cron", "start.go"), repairAutoBanCronStartSource},
		{filepath.Join(apiRoot, "internal", "jobs", "client.go"), repairAutoBanTaskTypeSource},
		{filepath.Join(apiRoot, "internal", "jobs", "workers.go"), repairAutoBanWorkerSource},
		{filepath.Join(apiRoot, "cmd", "server", "main.go"), repairAutoBanMainSource},
	} {
		if !fileExists(f.path) {
			continue
		}
		if err := repairSourceFile(root, m, f.path, f.fn); err != nil {
			return fmt.Errorf("%s: %w", f.path, err)
		}
	}
	return nil
}

// --------------------------------------------------------------------- config

const autoBanConfigFields = `	// The scheduled sweep. Sentinel's own AutoBlock catches a burst, ten
	// high-severity threats inside ten minutes; this catches the scanner
	// patient enough to stay under that. An address over
	// SecurityAutoBanMinThreats high-severity threats inside
	// SecurityAutoBanWindow is banned, every SecurityAutoBanEvery.
	SecurityAutoBanEnabled    bool
	SecurityAutoBanMinThreats int
	SecurityAutoBanWindow     time.Duration
	SecurityAutoBanEvery      time.Duration
	SecurityAutoBanPermanent  bool
`

const autoBanConfigLoads = `		SecurityAutoBanEnabled:    getEnv("SECURITY_AUTO_BAN", "true") == "true",
		SecurityAutoBanMinThreats: getEnvInt("SECURITY_AUTO_BAN_MIN_THREATS", 3),
		SecurityAutoBanWindow:     getEnvDuration("SECURITY_AUTO_BAN_WINDOW_HOURS", 24, time.Hour),
		SecurityAutoBanEvery:      getEnvDuration("AUTO_BAN_INTERVAL_HOURS", 6, time.Hour),
		SecurityAutoBanPermanent:  getEnv("SECURITY_AUTO_BAN_PERMANENT", "true") == "true",
`

// lineWith finds the whole line containing marker, and returns the offsets of
// its start and of the character after its newline.
//
// Anchors here are lines in a gofmt-formatted file, and gofmt aligns a struct's
// field names into columns: "SentinelTurnstileSecret string" in the template is
// "SentinelTurnstileSecret    string" once a longer field joins the block. An
// exact-substring anchor matches the template and not one real project, which
// is the quiet way a repair reports success on every project and changes none
// of them. Matching the line and ignoring its spacing is what survives.
func lineWith(src, marker string) (start, afterNewline int, ok bool) {
	at := strings.Index(src, marker)
	if at < 0 {
		return 0, 0, false
	}
	start = strings.LastIndexByte(src[:at], '\n') + 1
	end := strings.IndexByte(src[at:], '\n')
	if end < 0 {
		return start, len(src), true
	}
	return start, at + end + 1, true
}

func repairAutoBanConfigSource(src string) (string, []string, []string) {
	if strings.Contains(src, "SecurityAutoBanEnabled") {
		return src, nil, nil
	}
	at, _, ok := lineWith(src, "SentinelTurnstileSecret ")
	if !ok {
		return src, nil, []string{"its Config has no SentinelTurnstileSecret field where Grit wrote one, " +
			"so the scheduled ban sweep has nothing to read: add the SecurityAutoBan fields by hand"}
	}
	src = src[:at] + autoBanConfigFields + src[at:]

	_, after, ok := lineWith(src, "SentinelTurnstileSecret:")
	if !ok {
		return src, nil, []string{"its Config never loads SENTINEL_TURNSTILE_SECRET, so the sweep's " +
			"settings are not read from the environment: add the SecurityAutoBan loads by hand"}
	}
	src = src[:after] + autoBanConfigLoads + src[after:]
	return src, []string{"the scheduled ban sweep reads its policy from the environment"}, nil
}

// --------------------------------------------------------------------- routes

const autoBanRoutes = `		staff.GET("/admin/security/trends", middleware.RequireRole("ADMIN", "perm:system.view"), securityHandler.Trends)
		staff.GET("/admin/security/top", middleware.RequireRole("ADMIN", "perm:system.view"), securityHandler.Top)
		staff.GET("/admin/security/threats", middleware.RequireRole("ADMIN", "perm:system.view"), securityHandler.Threats)
		// Banning is a write, and one that can take the site off the internet
		// if it is pointed at the wrong address, so it is ADMIN and nothing
		// weaker. perm:system.view reads; it does not ban.
		staff.POST("/admin/security/ban", middleware.RequireRole("ADMIN"), securityHandler.Ban)
		staff.DELETE("/admin/security/ban/:ip", middleware.RequireRole("ADMIN"), securityHandler.Unban)
		staff.POST("/admin/security/ban-offenders", middleware.RequireRole("ADMIN"), securityHandler.BanOffenders)
`

func repairAutoBanRoutesSource(src string) (string, []string, []string) {
	var fixed, warn []string

	// The handler gained the sweep and the schedule it reports on the page.
	const oldHandler = "\tsecurityHandler := &handlers.SecurityHandler{Bridge: svc.SecObs}"
	if strings.Contains(src, oldHandler) {
		src = strings.Replace(src, oldHandler,
			"\t// The dashboard's button and the scheduled job share one sweep, so the\n"+
				"\t// policy shown on the page is the policy that runs at 3am.\n"+
				"\tsecurityHandler := &handlers.SecurityHandler{\n"+
				"\t\tBridge:       svc.SecObs,\n"+
				"\t\tAutoBan:      services.NewSecurityAutoBan(svc.SecObs, cfg),\n"+
				"\t\tAutoBanEvery: cfg.SecurityAutoBanEvery,\n"+
				"\t}", 1)
		fixed = append(fixed, "the security handler carries the ban sweep")
	} else if !strings.Contains(src, "AutoBan:") {
		warn = append(warn, "its securityHandler is not the one Grit wrote, so the ban sweep is "+
			"not wired to the dashboard: add AutoBan and AutoBanEvery to it")
	}

	if !strings.Contains(src, `"/admin/security/ban-offenders"`) {
		const routeAnchor = "\t\tstaff.GET(\"/admin/security/summary\", middleware.RequireRole(\"ADMIN\", \"perm:system.view\"), securityHandler.Summary)\n"
		at := strings.Index(src, routeAnchor)
		if at < 0 {
			warn = append(warn, "its /admin/security/summary route is not where Grit wrote it, so the "+
				"new security endpoints were not added: the dashboard will show only the summary")
		} else {
			at += len(routeAnchor)
			src = src[:at] + autoBanRoutes + src[at:]
			fixed = append(fixed, "the security dashboard's endpoints are routed")
		}
	}
	return src, fixed, warn
}

// ----------------------------------------------------------------------- cron

const autoBanCronRegistration = `	// The security sweep: bans the addresses that keep attacking.
	//
	// Sentinel's own AutoBlock answers a burst, ten high-severity threats from
	// one address inside ten minutes. This answers patience: four an hour for a
	// week never trips that rule, and it is most of what real scanning looks
	// like.
	autoBanSpec, autoBanWindow := everyHours(autoBanEvery)
	_, err = scheduler.Register(
		autoBanSpec,
		asynq.NewTask("security:auto_ban", nil),
		asynq.Unique(autoBanWindow-time.Minute),
		asynq.Timeout(15*time.Minute),
		asynq.MaxRetry(2),
	)
	if err != nil {
		return nil, fmt.Errorf("registering the security auto-ban sweep: %w", err)
	}
	RegisteredTasks = append(RegisteredTasks, Task{
		Name:     "Ban repeat attackers",
		Schedule: autoBanSpec,
		Type:     "security:auto_ban",
	})

`

const autoBanEveryHours = `
// everyHours turns a duration into an hourly cron spec, and returns the
// duration actually used.
//
// Cron has no "every 90 minutes": the spec is a set of hours, so anything that
// is not a whole number of hours has to become one. Rounding down, because a
// sweep asked to run more often should not run less often than asked.
func everyHours(d time.Duration) (string, time.Duration) {
	hours := int(d / time.Hour)
	if hours < 1 {
		hours = 1
	}
	if hours >= 24 {
		return "0 0 * * *", 24 * time.Hour
	}
	return fmt.Sprintf("0 */%d * * *", hours), time.Duration(hours) * time.Hour
}
`

func repairAutoBanCronSource(src string) (string, []string, []string) {
	if strings.Contains(src, `"security:auto_ban"`) {
		return src, nil, nil
	}
	const sig = "func New(redisURL string) (*Scheduler, error) {"
	if !strings.Contains(src, sig) {
		return src, nil, []string{"its cron.New is not the one Grit wrote, so the ban sweep is not " +
			"scheduled: give it an autoBanEvery parameter and register security:auto_ban"}
	}
	src = strings.Replace(src, sig,
		"// autoBanEvery is how often the security sweep runs. It is a parameter\n"+
			"// rather than a constant because it is the one schedule here an operator\n"+
			"// reasonably changes.\n"+
			"func New(redisURL string, autoBanEvery time.Duration) (*Scheduler, error) {", 1)

	const marker = "\t// grit:cron-tasks\n"
	at := strings.Index(src, marker)
	if at < 0 {
		return src, nil, []string{"its cron.go has no grit:cron-tasks marker, so the ban sweep was not " +
			"registered: add a scheduler.Register for security:auto_ban"}
	}
	src = src[:at] + autoBanCronRegistration + src[at:]

	// everyHours goes beside Start, which every generated cron.go has.
	const startDoc = "// Start runs the scheduler"
	sat := strings.Index(src, startDoc)
	if sat < 0 {
		sat = strings.Index(src, "// Start begins executing scheduled tasks.")
	}
	if sat < 0 {
		return src, nil, []string{"its cron.go has no Start to put everyHours beside"}
	}
	src = src[:sat] + strings.TrimPrefix(autoBanEveryHours, "\n") + "\n" + src[sat:]
	return src, []string{"the ban sweep runs on a schedule"}, nil
}

func repairAutoBanCronStartSource(src string) (string, []string, []string) {
	const old = "\ts, err := New(cfg.RedisURL)\n"
	if !strings.Contains(src, old) {
		return src, nil, nil
	}
	src = strings.Replace(src, old, "\ts, err := New(cfg.RedisURL, cfg.SecurityAutoBanEvery)\n", 1)
	return src, []string{"the scheduler is told how often to sweep"}, nil
}

// ------------------------------------------------------------- the task type

// repairAutoBanTaskTypeSource adds the constant, in the file that declares the
// others. It is a separate repair from the worker's because they are separate
// files: jobs/client.go holds the task types and jobs/workers.go handles them.
func repairAutoBanTaskTypeSource(src string) (string, []string, []string) {
	if strings.Contains(src, "TypeSecurityAutoBan") {
		return src, nil, nil
	}
	_, after, ok := lineWith(src, "TypeTrashPurgeExpired ")
	if !ok {
		return src, nil, []string{"its task-type list is not the one Grit wrote, so security:auto_ban " +
			"has no constant: add TypeSecurityAutoBan = \"security:auto_ban\""}
	}
	src = src[:after] + "\tTypeSecurityAutoBan = \"security:auto_ban\"\n" + src[after:]
	return src, []string{"security:auto_ban has a task type"}, nil
}

// -------------------------------------------------------------------- worker

func repairAutoBanWorkerSource(src string) (string, []string, []string) {
	if strings.Contains(src, "TypeSecurityAutoBan") {
		return src, nil, nil
	}
	var warn []string

	const depsAnchor = "type WorkerDeps struct {\n\tDB      *gorm.DB\n\tMailer  *mail.Mailer\n\tStorage *storage.Storage\n\tCache   *cache.Cache\n}"
	if strings.Contains(src, depsAnchor) {
		src = strings.Replace(src, depsAnchor,
			"type WorkerDeps struct {\n\tDB      *gorm.DB\n\tMailer  *mail.Mailer\n\tStorage *storage.Storage\n\tCache   *cache.Cache\n"+
				"\t// AutoBan is the scheduled security sweep. Nil when Sentinel is off,\n"+
				"\t// and the job says so rather than failing.\n\tAutoBan *services.SecurityAutoBan\n}", 1)
	} else if !strings.Contains(src, "AutoBan *services.SecurityAutoBan") {
		warn = append(warn, "its WorkerDeps is not the one Grit wrote, so the sweep cannot reach the worker")
	}

	if _, after, ok := lineWith(src, "mux.HandleFunc(TypeTrashPurgeExpired"); ok {
		src = src[:after] + "\tmux.HandleFunc(TypeSecurityAutoBan, handleSecurityAutoBan(deps))\n" + src[after:]
	} else {
		warn = append(warn, "its worker mux is not the one Grit wrote, so security:auto_ban has no handler")
	}

	const handlerAnchor = "// handleBackupWeekly takes the scheduled full-database backup"
	if at := strings.Index(src, handlerAnchor); at >= 0 {
		src = src[:at] + autoBanJobHandler + src[at:]
	} else {
		warn = append(warn, "its job handlers are not where Grit wrote them, so handleSecurityAutoBan was not added")
	}
	if len(warn) > 0 {
		return src, nil, warn
	}
	return src, []string{"the ban sweep runs as a background job"}, nil
}

const autoBanJobHandler = `// handleSecurityAutoBan bans the addresses that keep attacking.
//
// High severity only. An actor's total threat count includes every
// low-severity heuristic it ever tripped, and banning on that number is how a
// dashboard starts refusing real users.
func handleSecurityAutoBan(deps WorkerDeps) func(ctx context.Context, task *asynq.Task) error {
	return func(ctx context.Context, task *asynq.Task) error {
		if deps.AutoBan == nil {
			log.Println("Auto-ban sweep: Sentinel is not configured, nothing to sweep")
			return nil
		}
		if !deps.AutoBan.Config.Enabled {
			log.Println("Auto-ban sweep: turned off by SECURITY_AUTO_BAN")
			return nil
		}
		result, err := deps.AutoBan.Run(ctx, false)
		if err != nil {
			return fmt.Errorf("auto-ban sweep: %w", err)
		}
		result.LogResult(deps.AutoBan.Config)
		return nil
	}
}

`

// ---------------------------------------------------------------------- main

func repairAutoBanMainSource(src string) (string, []string, []string) {
	var fixed, warn []string

	const oldDeps = "\t\tstop, err := jobs.StartWorker(cfg.RedisURL, jobs.WorkerDeps{\n\t\t\tDB:      db,\n\t\t\tMailer:  mailer,\n\t\t\tStorage: storageService,\n\t\t\tCache:   cacheService,\n\t\t})"
	if strings.Contains(src, oldDeps) {
		src = strings.Replace(src, oldDeps,
			"\t\tstop, err := jobs.StartWorker(cfg.RedisURL, jobs.WorkerDeps{\n"+
				"\t\t\tDB:      db,\n\t\t\tMailer:  mailer,\n\t\t\tStorage: storageService,\n\t\t\tCache:   cacheService,\n"+
				"\t\t\t// Nil when Sentinel is off, and the sweep says so rather than\n"+
				"\t\t\t// failing: a project with no firewall did not ask for a ban sweep.\n"+
				"\t\t\tAutoBan: services.NewSecurityAutoBan(secObsBridge, cfg),\n\t\t})", 1)
		fixed = append(fixed, "the worker is given the ban sweep")
	} else if !strings.Contains(src, "AutoBan: services.NewSecurityAutoBan") {
		warn = append(warn, "its jobs.StartWorker call is not the one Grit wrote, so the sweep never "+
			"reaches the worker: add AutoBan: services.NewSecurityAutoBan(secObsBridge, cfg)")
	}

	const oldCron = "\t\tcs, err := cron.New(cfg.RedisURL)\n"
	if strings.Contains(src, oldCron) {
		src = strings.Replace(src, oldCron, "\t\tcs, err := cron.New(cfg.RedisURL, cfg.SecurityAutoBanEvery)\n", 1)
		fixed = append(fixed, "the scheduler is told how often to sweep")
	} else if strings.Contains(src, "cron.New(cfg.RedisURL)") {
		warn = append(warn, "its cron.New call is not where Grit wrote it: pass cfg.SecurityAutoBanEvery as its second argument")
	}
	return src, fixed, warn
}
