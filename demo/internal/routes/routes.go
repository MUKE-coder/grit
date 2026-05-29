package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/MUKE-coder/pulse/pulse"
	"github.com/MUKE-coder/sentinel"

	"gritdemo/internal/cache"
	"gritdemo/internal/config"
	"gritdemo/internal/dgateway"
	"gritdemo/internal/handlers"
	"gritdemo/internal/mail"
	"gritdemo/internal/middleware"
	"gritdemo/internal/services"
	"gritdemo/internal/storage"
)

// Services holds all injectable services.
type Services struct {
	Cache   *cache.Cache
	Storage *storage.Storage
	Mailer  *mail.Mailer
	// SecObsBridge talks to the locally-mounted Sentinel/Pulse APIs over
	// loopback so the in-app Security/Observability dashboards can show
	// summary cards without an iframe. Nil when both are disabled.
	SecObs  *services.SecObsBridge
}

// Setup configures all routes and returns the Gin engine.
func Setup(db *gorm.DB, cfg *config.Config, svc *Services) *gin.Engine {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(cfg.CORSOrigins))
	r.Use(middleware.Logger())

	// ─── Sentinel security suite (mounts /sentinel/ui + /sentinel/api) ───
	if cfg.SentinelEnabled {
		isDev := cfg.AppEnv != "production"
		if err := sentinel.MountE(r, db, sentinel.Config{
			Dashboard: sentinel.DashboardConfig{
				Username:              cfg.SentinelUsername,
				Password:              cfg.SentinelPassword,
				SecretKey:             cfg.SentinelSecretKey,
				AllowInsecureDefaults: isDev,
			},
			WAF: sentinel.WAFConfig{
				Enabled: true,
				Mode: func() sentinel.WAFMode {
					if isDev {
						return sentinel.ModeLog
					}
					return sentinel.ModeBlock
				}(),
				TrustedProxies:      cfg.SentinelTrustedProxies,
				MaxBodyBytes:        64 * 1024,
				RejectOversizedBody: true,
			},
			RateLimit: sentinel.RateLimitConfig{
				Enabled: !isDev,
				ByIP:    &sentinel.Limit{Requests: 100, Window: time.Minute},
				ByRoute: map[string]sentinel.Limit{
					"/api/auth/login": {Requests: 5, Window: 15 * time.Minute},
				},
			},
			AuthShield: sentinel.AuthShieldConfig{Enabled: !isDev, LoginRoute: "/api/auth/login"},
			Anomaly:    sentinel.AnomalyConfig{Enabled: !isDev},
			Geo:        sentinel.GeoConfig{Enabled: !isDev},
		}); err != nil {
			log.Printf("Warning: Sentinel mount failed: %v", err)
		} else {
			log.Println("Sentinel mounted at /sentinel")
		}
	}

	// ─── Pulse observability (mounts /pulse/ui + /pulse/api) ──────────────
	if cfg.PulseEnabled {
		opts := []pulse.Option{
			pulse.WithAppName(cfg.AppName),
			pulse.WithCredentials(cfg.PulseUsername, cfg.PulsePassword),
			pulse.WithExcludePaths("/studio/*", "/sentinel/*", "/docs/*", "/pulse/*"),
			pulse.WithPrometheus(),
		}
		if cfg.AppEnv != "production" {
			opts = append(opts, pulse.WithDevMode())
		}
		if cfg.PulseStorage == "sqlite" && cfg.PulseStorageDSN != "" {
			opts = append(opts, pulse.WithSQLite(cfg.PulseStorageDSN))
		}
		pulse.Mount(context.Background(), r, db, opts...)
		log.Println("Pulse mounted at /pulse")
	}

	// Auth service
	authService := &services.AuthService{
		Secret:        cfg.JWTSecret,
		AccessExpiry:  cfg.JWTAccessExpiry,
		RefreshExpiry: cfg.JWTRefreshExpiry,
	}

	// Audit service — single instance, shared across handlers. Async writes;
	// failures log to stdout but don't break the request.
	auditService := services.NewAuditService(db)

	// Handlers
	authHandler := &handlers.AuthHandler{DB: db, AuthService: authService, Config: cfg, Audit: auditService}
	businessHandler := &handlers.BusinessHandler{DB: db}
	branchHandler := &handlers.BranchHandler{DB: db, Audit: auditService}
	staffHandler := &handlers.StaffHandler{DB: db, Audit: auditService}
	invitationHandler := &handlers.InvitationHandler{
		DB: db, Mailer: svc.Mailer, AuthService: authService, Config: cfg, Audit: auditService,
	}
	categoryHandler := &handlers.CategoryHandler{DB: db}
	productHandler := &handlers.ProductHandler{DB: db, Storage: svc.Storage, Audit: auditService}
	stockHandler := &handlers.StockHandler{DB: db, Mailer: svc.Mailer, Audit: auditService}
	saleHandler := &handlers.SaleHandler{DB: db, Mailer: svc.Mailer, Audit: auditService}
	reportHandler := &handlers.ReportHandler{DB: db}
	segmentReportHandler := &handlers.SegmentReportHandler{DB: db}
	activityHandler := &handlers.ActivityHandler{DB: db}

	// Motorcycle / Loan / Daily Boda handlers (Grit Motors).
	motorcycleHandler := &handlers.MotorcycleHandler{DB: db, Storage: svc.Storage, Audit: auditService}
	loanProductHandler := &handlers.LoanProductHandler{DB: db, Audit: auditService}
	borrowerHandler := &handlers.BorrowerHandler{DB: db, Storage: svc.Storage, Audit: auditService}
	loanHandler := &handlers.LoanHandler{DB: db, Audit: auditService}
	repaymentHandler := &handlers.RepaymentHandler{DB: db, Audit: auditService}
	cashSaleHandler := &handlers.CashSaleHandler{DB: db, Audit: auditService}
	dailyBodaHandler := &handlers.DailyBodaHandler{DB: db, Audit: auditService}

	adminHandler := &handlers.AdminHandler{DB: db}

	// In-app Security + Observability dashboards (read from local Sentinel/
	// Pulse APIs) + the notification bell.
	notificationHandler := &handlers.NotificationHandler{DB: db}
	securityHandler := &handlers.SecurityHandler{Bridge: svc.SecObs}
	observabilityHandler := &handlers.ObservabilityHandler{Bridge: svc.SecObs}

	// DGateway client (mobile money). nil if not configured — handlers return 503.
	var dgClient *dgateway.Client
	if cfg.DGatewayAPIKey != "" {
		dgClient = dgateway.New(cfg.DGatewayAPIKey, cfg.DGatewayBaseURL)
	}
	dgatewayHandler := &handlers.DGatewayHandler{
		DB:                 db,
		DG:                 dgClient,
		DefaultProvider:    cfg.DGatewayDefaultProvider,
		WebhookCallbackURL: cfg.AppURL + "/api/dgateway/webhook",
	}

	// Health check (no API key required)
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "app": "grit-motors", "version": "0.1.0"})
	})

	// DGateway webhook — public endpoint. DGateway can't carry our X-Internal-Key,
	// so we accept the inbound POST and re-verify with DGateway before applying any state.
	r.POST("/api/dgateway/webhook", dgatewayHandler.Webhook)

	// All /api routes require X-Internal-Key
	api := r.Group("/api")
	api.Use(middleware.InternalKey(cfg.InternalAPIKey))

	// Public auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
	}

	// Public invite routes (no JWT, just API key)
	api.GET("/invite/:token", invitationHandler.GetInvite)
	api.POST("/invite/:token/accept", invitationHandler.AcceptInvite)

	// Protected routes (API key + JWT)
	protected := api.Group("")
	protected.Use(middleware.Auth(db, authService))
	{
		protected.GET("/auth/me", authHandler.Me)
		protected.POST("/auth/logout", authHandler.Logout)

		// Businesses (user-level, no business scope needed)
		protected.GET("/businesses", businessHandler.List)
		protected.POST("/businesses", businessHandler.Create)

		// Notification bell (any authenticated user; admins also see broadcast)
		protected.GET("/notifications", notificationHandler.List)
		protected.POST("/notifications/:id/read", notificationHandler.MarkRead)
		protected.POST("/notifications/read-all", notificationHandler.MarkAllRead)

		// In-app Security + Observability summary endpoints. Both proxy
		// loopback calls into the locally-mounted Sentinel/Pulse APIs and
		// return a single envelope so the React page does one round-trip.
		// (Auth middleware is enough — these are admin-only in the UI.)
		protected.GET("/admin/security/summary", securityHandler.Summary)
		protected.GET("/admin/observability/summary", observabilityHandler.Summary)
	}

	// Business-scoped routes (API key + JWT + X-Business-ID)
	scoped := api.Group("")
	scoped.Use(middleware.Auth(db, authService))
	scoped.Use(middleware.BusinessAccess(db))
	{
		// Business settings
		scoped.PUT("/businesses/:id", businessHandler.Update)

		// Branches
		scoped.GET("/branches", branchHandler.List)
		scoped.POST("/branches", branchHandler.Create)
		scoped.PUT("/branches/:id", branchHandler.Update)
		scoped.DELETE("/branches/:id", branchHandler.Delete)

		// Staff
		scoped.GET("/staff", staffHandler.List)
		scoped.POST("/staff", staffHandler.Create) // direct create (admin only) — replaces invite as primary path
		scoped.PUT("/staff/:userId/role", staffHandler.UpdateRole)
		scoped.DELETE("/staff/:userId", staffHandler.Remove)

		// Invitations (kept for backwards compatibility; UI prefers /staff direct-create)
		scoped.POST("/invitations", invitationHandler.Create)
		scoped.GET("/invitations", invitationHandler.List)
		scoped.DELETE("/invitations/:id", invitationHandler.Revoke)

		// Activities (audit log, admin only)
		scoped.GET("/activities", activityHandler.List)

		// Categories
		scoped.GET("/categories", categoryHandler.List)
		scoped.POST("/categories", categoryHandler.Create)
		scoped.PUT("/categories/:id", categoryHandler.Update)
		scoped.DELETE("/categories/:id", categoryHandler.Delete)

		// Products
		scoped.GET("/products", productHandler.List)
		scoped.GET("/products/pos", productHandler.POS)
		scoped.GET("/products/:id", productHandler.Get)
		scoped.POST("/products", productHandler.Create)
		scoped.POST("/products/import", productHandler.Import)
		scoped.PUT("/products/:id", productHandler.Update)
		scoped.DELETE("/products/:id", productHandler.Delete)
		scoped.POST("/products/:id/image", productHandler.UploadImage)
		scoped.DELETE("/products/:id/image", productHandler.DeleteImage)

		// Stock
		scoped.GET("/stock", stockHandler.Levels)
		scoped.POST("/stock/in", stockHandler.StockIn)
		scoped.POST("/stock/transfer", stockHandler.Transfer)
		scoped.GET("/stock/movements", stockHandler.Movements)
		scoped.GET("/stock/transfers", stockHandler.Transfers)

		// Sales
		scoped.POST("/sales", saleHandler.Create)
		scoped.GET("/sales", saleHandler.List)
		scoped.GET("/sales/:id", saleHandler.Get)
		scoped.POST("/sales/:id/returns", saleHandler.CreateReturn)
		scoped.GET("/sales/:id/returns", saleHandler.ListReturns)

		// Reports — spares
		scoped.GET("/reports/dashboard", reportHandler.Dashboard)
		scoped.GET("/reports/daily", reportHandler.Daily)
		scoped.GET("/reports/stock", reportHandler.StockReport)
		scoped.GET("/reports/pnl", reportHandler.PnL)
		scoped.GET("/reports/cashiers", reportHandler.Cashiers)

		// Reports — segments + combined home
		scoped.GET("/reports/main", segmentReportHandler.Main)
		scoped.GET("/reports/loans", segmentReportHandler.LoansReport)
		scoped.GET("/reports/collections", segmentReportHandler.CollectionsReport)
		scoped.GET("/reports/motorcycles", segmentReportHandler.MotorcyclesReport)
		scoped.GET("/reports/daily-boda", segmentReportHandler.DailyBodaReport)

		// Motorcycles inventory
		scoped.GET("/motorcycles", motorcycleHandler.List)
		scoped.POST("/motorcycles", motorcycleHandler.Create)
		scoped.GET("/motorcycles/:id", motorcycleHandler.Get)
		scoped.PUT("/motorcycles/:id", motorcycleHandler.Update)
		scoped.DELETE("/motorcycles/:id", motorcycleHandler.Delete)
		scoped.POST("/motorcycles/:id/transfer", motorcycleHandler.Transfer)
		scoped.POST("/motorcycles/:id/image", motorcycleHandler.UploadImage)

		// Loan products
		scoped.GET("/loan-products", loanProductHandler.List)
		scoped.POST("/loan-products", loanProductHandler.Create)
		scoped.GET("/loan-products/:id", loanProductHandler.Get)
		scoped.PUT("/loan-products/:id", loanProductHandler.Update)
		scoped.DELETE("/loan-products/:id", loanProductHandler.Delete)

		// Borrowers
		scoped.GET("/borrowers", borrowerHandler.List)
		scoped.POST("/borrowers", borrowerHandler.Create)
		scoped.GET("/borrowers/:id", borrowerHandler.Get)
		scoped.PUT("/borrowers/:id", borrowerHandler.Update)
		scoped.DELETE("/borrowers/:id", borrowerHandler.Delete)
		scoped.POST("/borrowers/:id/document", borrowerHandler.UploadDocument)

		// Loans
		scoped.GET("/loans", loanHandler.List)
		scoped.POST("/loans", loanHandler.Create)
		scoped.GET("/loans/:id", loanHandler.Get)
		scoped.GET("/loans/:id/schedule", loanHandler.Schedule)
		scoped.POST("/loans/:id/approve", loanHandler.Approve)
		scoped.POST("/loans/:id/reject", loanHandler.Reject)
		scoped.POST("/loans/:id/disburse", loanHandler.Disburse)

		// Repayments
		scoped.GET("/repayments", repaymentHandler.List)
		scoped.POST("/repayments", repaymentHandler.Create)
		scoped.GET("/repayments/:id", repaymentHandler.Get)
		scoped.POST("/repayments/:id/approve", repaymentHandler.Approve)
		scoped.POST("/repayments/:id/reject", repaymentHandler.Reject)

		// Cash sales (motorcycles)
		scoped.GET("/cash-sales", cashSaleHandler.List)
		scoped.POST("/cash-sales", cashSaleHandler.Create)
		scoped.GET("/cash-sales/:id", cashSaleHandler.Get)

		// Daily Boda — drivers
		scoped.GET("/daily-boda/drivers", dailyBodaHandler.ListDrivers)
		scoped.POST("/daily-boda/drivers", dailyBodaHandler.CreateDriver)
		scoped.GET("/daily-boda/drivers/:id", dailyBodaHandler.GetDriver)
		scoped.PUT("/daily-boda/drivers/:id", dailyBodaHandler.UpdateDriver)

		// Daily Boda — motorcycles
		scoped.GET("/daily-boda/motorcycles", dailyBodaHandler.ListMotorcycles)
		scoped.POST("/daily-boda/motorcycles", dailyBodaHandler.CreateMotorcycle)
		scoped.POST("/daily-boda/motorcycles/:id/assign", dailyBodaHandler.AssignDriver)
		scoped.POST("/daily-boda/motorcycles/:id/return", dailyBodaHandler.ReturnMotorcycle)

		// Daily Boda — payments
		scoped.GET("/daily-boda/payments", dailyBodaHandler.ListPayments)
		scoped.POST("/daily-boda/payments", dailyBodaHandler.CreatePayment)
		scoped.POST("/daily-boda/payments/:id/verify", dailyBodaHandler.VerifyPayment)

		// DGateway mobile money (collect = initiate, verify = poll status)
		scoped.POST("/dgateway/collect", dgatewayHandler.Collect)
		scoped.POST("/dgateway/verify/:id", dgatewayHandler.Verify)

		// Admin one-shot ops (TEMPORARY — remove once first-deploy demo data is cleaned)
		scoped.POST("/admin/wipe-seed", adminHandler.WipeSeedData)
	}

	return r
}
