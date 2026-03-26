package deploy

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Config holds deployment configuration.
type Config struct {
	Host       string // SSH host (e.g. user@server.com)
	Port       string // SSH port (default: 22)
	KeyFile    string // Path to SSH private key
	AppName    string // Application name
	RemotePath string // Remote deployment directory (default: /opt/{app-name})
	Domain     string // Domain name for Caddy reverse proxy
	AppPort    string // Port the app runs on (default: 8080)
}

// Run executes the full deployment pipeline.
func Run(cfg Config) error {
	if cfg.Host == "" {
		return fmt.Errorf("deploy host is required (set DEPLOY_HOST in .env or use --host)")
	}

	if cfg.Port == "" {
		cfg.Port = "22"
	}
	if cfg.RemotePath == "" {
		cfg.RemotePath = "/opt/" + cfg.AppName
	}
	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}

	steps := []struct {
		name string
		fn   func() error
	}{
		{"Building application", func() error { return buildBinary(cfg) }},
		{"Uploading binary", func() error { return uploadBinary(cfg) }},
		{"Setting up systemd service", func() error { return setupSystemd(cfg) }},
		{"Restarting service", func() error { return restartService(cfg) }},
	}

	// Optional: Caddy setup if domain is provided
	if cfg.Domain != "" {
		steps = append(steps, struct {
			name string
			fn   func() error
		}{"Configuring Caddy reverse proxy", func() error { return setupCaddy(cfg) }})
	}

	for _, step := range steps {
		fmt.Printf("  → %s...\n", step.name)
		if err := step.fn(); err != nil {
			return fmt.Errorf("%s: %w", step.name, err)
		}
	}

	return nil
}

// buildBinary cross-compiles the Go binary for Linux.
func buildBinary(cfg Config) error {
	binaryName := cfg.AppName
	outputPath := filepath.Join("bin", binaryName)

	// Check if frontend exists and build it first
	if _, err := os.Stat("frontend"); err == nil {
		fmt.Println("    Building frontend...")
		feCmd := exec.Command("pnpm", "build")
		feCmd.Dir = "frontend"
		feCmd.Stdout = os.Stdout
		feCmd.Stderr = os.Stderr
		if err := feCmd.Run(); err != nil {
			return fmt.Errorf("building frontend: %w", err)
		}
	}

	// Cross-compile for Linux
	cmd := exec.Command("go", "build", "-o", outputPath, "./cmd/server")
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64", "CGO_ENABLED=0")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// uploadBinary sends the binary to the remote server via SCP.
func uploadBinary(cfg Config) error {
	binaryName := cfg.AppName
	localPath := filepath.Join("bin", binaryName)
	remotePath := cfg.RemotePath + "/" + binaryName

	// Create remote directory
	sshCmd(cfg, fmt.Sprintf("mkdir -p %s", cfg.RemotePath))

	// Upload via SCP
	args := []string{"-P", cfg.Port}
	if cfg.KeyFile != "" {
		args = append(args, "-i", cfg.KeyFile)
	}
	args = append(args, localPath, cfg.Host+":"+remotePath)

	cmd := exec.Command("scp", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("uploading binary: %w", err)
	}

	// Make executable
	return sshCmd(cfg, fmt.Sprintf("chmod +x %s/%s", cfg.RemotePath, binaryName))
}

// setupSystemd creates a systemd service unit file on the remote server.
func setupSystemd(cfg Config) error {
	unit := SystemdUnit(cfg)
	serviceName := cfg.AppName + ".service"

	// Write unit file via SSH
	escaped := strings.ReplaceAll(unit, "'", "'\\''")
	return sshCmd(cfg, fmt.Sprintf("echo '%s' | sudo tee /etc/systemd/system/%s > /dev/null && sudo systemctl daemon-reload && sudo systemctl enable %s", escaped, serviceName, serviceName))
}

// restartService restarts the systemd service.
func restartService(cfg Config) error {
	return sshCmd(cfg, fmt.Sprintf("sudo systemctl restart %s", cfg.AppName))
}

// setupCaddy configures Caddy as a reverse proxy with auto-TLS.
func setupCaddy(cfg Config) error {
	caddyfile := Caddyfile(cfg)
	escaped := strings.ReplaceAll(caddyfile, "'", "'\\''")
	return sshCmd(cfg, fmt.Sprintf("echo '%s' | sudo tee /etc/caddy/Caddyfile > /dev/null && sudo systemctl reload caddy", escaped))
}

// sshCmd executes a command on the remote server.
func sshCmd(cfg Config, command string) error {
	args := []string{"-p", cfg.Port}
	if cfg.KeyFile != "" {
		args = append(args, "-i", cfg.KeyFile)
	}
	args = append(args, cfg.Host, command)

	cmd := exec.Command("ssh", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SystemdUnit generates a systemd service unit file.
func SystemdUnit(cfg Config) string {
	return fmt.Sprintf(`[Unit]
Description=%s
After=network.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=%s
ExecStart=%s/%s
Restart=on-failure
RestartSec=5
Environment=APP_ENV=production
Environment=APP_PORT=%s
EnvironmentFile=%s/.env

# Security hardening
NoNewPrivileges=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=%s

[Install]
WantedBy=multi-user.target
`, cfg.AppName, cfg.RemotePath, cfg.RemotePath, cfg.AppName, cfg.AppPort, cfg.RemotePath, cfg.RemotePath)
}

// Caddyfile generates a Caddy reverse proxy configuration with auto-TLS.
func Caddyfile(cfg Config) string {
	return fmt.Sprintf(`%s {
    reverse_proxy localhost:%s

    encode gzip

    header {
        X-Frame-Options "DENY"
        X-Content-Type-Options "nosniff"
        Referrer-Policy "strict-origin-when-cross-origin"
        -Server
    }

    log {
        output file /var/log/caddy/%s.log {
            roll_size 10mb
            roll_keep 5
        }
    }
}
`, cfg.Domain, cfg.AppPort, cfg.AppName)
}

// Timestamp returns a deployment timestamp for logging.
func Timestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
