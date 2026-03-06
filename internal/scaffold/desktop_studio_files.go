package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"
)

func writeDesktopStudioFiles(root string, opts DesktopOptions) error {
	files := map[string]string{
		filepath.Join(root, "cmd", "studio", "main.go"): desktopStudioMainGo(),
	}

	for path, content := range files {
		content = strings.ReplaceAll(content, "<MODULE>", opts.ProjectName)
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func desktopStudioMainGo() string {
	return `package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/MUKE-coder/gorm-studio/studio"

	"<MODULE>/internal/config"
	"<MODULE>/internal/db"
	"<MODULE>/internal/models"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()
	database := db.Connect(cfg)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	studio.Mount(r, database, []interface{}{
		&models.User{},
		&models.Blog{},
		&models.Contact{},
		// grit:studio-models
	}, studio.Config{
		Prefix: "/studio",
	})

	// Auto-open browser after server starts
	go func() {
		time.Sleep(500 * time.Millisecond)
		url := "http://localhost:8080/studio"
		switch runtime.GOOS {
		case "windows":
			exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		case "darwin":
			exec.Command("open", url).Start()
		default:
			exec.Command("xdg-open", url).Start()
		}
	}()

	fmt.Println("GORM Studio running at http://localhost:8080/studio")
	fmt.Println("Press Ctrl+C to stop")
	log.Fatal(r.Run(":8080"))
}
`
}
