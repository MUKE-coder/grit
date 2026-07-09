package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeDesktopGoFiles(root string, opts DesktopOptions) error {
	files := map[string]string{
		filepath.Join(root, "main.go"): desktopMainGo(),
		filepath.Join(root, "app.go"):  desktopAppGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "<MODULE>", opts.ProjectName)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func desktopMainGo() string {
	return `package main

import (
	"context"
	"embed"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"<MODULE>/internal/api"
	"<MODULE>/internal/config"
	"<MODULE>/internal/db"
	"<MODULE>/internal/service"
	"<MODULE>/internal/storage"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	database := db.Connect(cfg)

	authSvc := service.NewAuthService(database)
	blogSvc := service.NewBlogService(database)
	contactSvc := service.NewContactService(database)
	exportSvc := service.NewExportService()
	// grit:service-init

	app := NewApp(authSvc, blogSvc, contactSvc, exportSvc, /* grit:app-args */)

	// Hybrid: alongside the Wails bindings this app runs a real REST API.
	//
	// Uploads land in the OS app-data dir (writable even when the binary sits in
	// Program Files), and the router is mounted twice:
	//
	//   AssetServer.Handler  → the webview fetches /api/... and loads
	//                          <img src="/uploads/x.jpg"> same-origin, so there's
	//                          no port to discover and no CORS to configure.
	//   127.0.0.1:APIPort    → curl, scripts, or a second client hit the same API.
	store, err := storage.New(cfg.AppName)
	if err != nil {
		log.Fatalf("storage: %v", err)
	}
	router := api.NewRouter(&api.Server{DB: database, Storage: store})

	go func() {
		addr := "127.0.0.1:" + cfg.APIPort
		log.Printf("REST API listening on http://%s", addr)
		if err := http.ListenAndServe(addr, router); err != nil {
			// A bind error here almost always means the port is already
			// taken — most commonly by the Wails dev server (34115) if
			// API_PORT was set back to it. The in-webview API still works
			// (it's mounted on the AssetServer handler, port-free); only
			// the loopback listener for curl/other clients is affected.
			log.Printf("WARNING: embedded REST API could not start on %s: %v", addr, err)
			log.Printf("         (in-app features still work; set a free API_PORT to enable curl/external access)")
		}
	}()

	// In-app auto-updater. Best-effort cleanup of the .exe.old left
	// behind by a previous in-app update runs in the background so
	// startup isn't held up by a busy filesystem.
	updater := NewUpdater()
	go updater.CleanupOldOnStartup()

	err = wails.Run(&options.App{
		Title:     cfg.AppName,
		Width:     1280,
		Height:    800,
		MinWidth:  900,
		MinHeight: 600,
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
			// Anything the frontend bundle doesn't resolve (/api/*, /uploads/*)
			// falls through to the REST API — same origin as the webview.
			Handler: router,
		},
		BackgroundColour: &options.RGBA{R: 10, G: 10, B: 15, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			updater.SetContext(ctx)
		},
		Bind: []interface{}{
			app,
			updater,
		},
		CSSDragProperty: "--wails-draggable",
		CSSDragValue:    "drag",
	})
	if err != nil {
		log.Fatal(err)
	}
}
`
}

func desktopAppGo() string {
	return `package main

import (
	"context"
	"fmt"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"<MODULE>/internal/models"
	"<MODULE>/internal/service"
	// grit:imports
)

// fmtDesktopDate renders a *time.Time for CSV/PDF export, tolerating nil (an
// unset optional date) instead of panicking.
func fmtDesktopDate(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

type App struct {
	ctx     context.Context
	auth    *service.AuthService
	blog    *service.BlogService
	contact *service.ContactService
	export  *service.ExportService
	// grit:fields
}

func NewApp(auth *service.AuthService, blog *service.BlogService, contact *service.ContactService, export *service.ExportService, /* grit:constructor-params */) *App {
	return &App{auth: auth, blog: blog, contact: contact, export: export, /* grit:constructor-assign */}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ── Auth ─────────────────────────────────────────────────────────────────────

func (a *App) Login(email, password string) (*models.AuthResponse, error) {
	return a.auth.Login(email, password)
}

func (a *App) Register(name, email, password string) (*models.AuthResponse, error) {
	return a.auth.Register(name, email, password)
}

func (a *App) GetCurrentUser(id string) (*models.User, error) {
	return a.auth.GetUser(id)
}

func (a *App) UpdateUser(id string, name, email string) (*models.User, error) {
	return a.auth.UpdateUser(id, name, email)
}

// ── Blogs ────────────────────────────────────────────────────────────────────

func (a *App) GetBlogs(page, pageSize int, search string) (*models.PaginatedResult, error) {
	return a.blog.List(page, pageSize, search)
}

func (a *App) GetBlog(id string) (*models.Blog, error) {
	return a.blog.GetByID(id)
}

func (a *App) CreateBlog(input models.BlogInput) (*models.Blog, error) {
	return a.blog.Create(input)
}

func (a *App) UpdateBlog(id string, input models.BlogInput) (*models.Blog, error) {
	return a.blog.Update(id, input)
}

func (a *App) DeleteBlog(id string) error {
	return a.blog.Delete(id)
}

// ── Contacts ─────────────────────────────────────────────────────────────────

func (a *App) GetContacts(page, pageSize int, search string) (*models.PaginatedResult, error) {
	return a.contact.List(page, pageSize, search)
}

func (a *App) GetContact(id string) (*models.Contact, error) {
	return a.contact.GetByID(id)
}

func (a *App) CreateContact(input models.ContactInput) (*models.Contact, error) {
	return a.contact.Create(input)
}

func (a *App) UpdateContact(id string, input models.ContactInput) (*models.Contact, error) {
	return a.contact.Update(id, input)
}

func (a *App) DeleteContact(id string) error {
	return a.contact.Delete(id)
}

// ── Export ────────────────────────────────────────────────────────────────────

func (a *App) ExportBlogsPDF() (string, error) {
	blogs, err := a.blog.ListAll()
	if err != nil {
		return "", err
	}
	headers := []string{"Title", "Slug", "Published", "Created"}
	var rows [][]string
	for _, b := range blogs {
		published := "Draft"
		if b.Published {
			published = "Published"
		}
		rows = append(rows, []string{b.Title, b.Slug, published, b.CreatedAt.Format("2006-01-02")})
	}
	return a.export.ExportPDF("Blogs", headers, rows)
}

func (a *App) ExportBlogsExcel() (string, error) {
	blogs, err := a.blog.ListAll()
	if err != nil {
		return "", err
	}
	headers := []string{"Title", "Slug", "Published", "Created"}
	var rows [][]string
	for _, b := range blogs {
		published := "Draft"
		if b.Published {
			published = "Published"
		}
		rows = append(rows, []string{b.Title, b.Slug, published, b.CreatedAt.Format("2006-01-02")})
	}
	return a.export.ExportExcel("Blogs", headers, rows)
}

func (a *App) ExportContactsPDF() (string, error) {
	contacts, err := a.contact.ListAll()
	if err != nil {
		return "", err
	}
	headers := []string{"Name", "Email", "Phone", "Company"}
	var rows [][]string
	for _, c := range contacts {
		rows = append(rows, []string{c.Name, c.Email, c.Phone, c.Company})
	}
	return a.export.ExportPDF("Contacts", headers, rows)
}

func (a *App) ExportContactsExcel() (string, error) {
	contacts, err := a.contact.ListAll()
	if err != nil {
		return "", err
	}
	headers := []string{"Name", "Email", "Phone", "Company"}
	var rows [][]string
	for _, c := range contacts {
		rows = append(rows, []string{c.Name, c.Email, c.Phone, c.Company})
	}
	return a.export.ExportExcel("Contacts", headers, rows)
}

// ── File Dialog ──────────────────────────────────────────────────────────────

func (a *App) OpenFileDialog(title string) (string, error) {
	result, err := wailsRuntime.OpenFileDialog(a.ctx, wailsRuntime.OpenDialogOptions{
		Title: title,
	})
	if err != nil {
		return "", fmt.Errorf("file dialog: %w", err)
	}
	return result, nil
}

// ── Window Controls ──────────────────────────────────────────────────────────

func (a *App) MinimiseWindow() {
	wailsRuntime.WindowMinimise(a.ctx)
}

func (a *App) ToggleMaximise() {
	wailsRuntime.WindowToggleMaximise(a.ctx)
}

func (a *App) CloseApp() {
	wailsRuntime.Quit(a.ctx)
}

// grit:methods
`
}
