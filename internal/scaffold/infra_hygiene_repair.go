package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Contact-app review L10 to L14. The new text lives here once: the templates
// splice it in, and repairInfraHygiene puts it into an existing project,
// anchored on what Grit generated before.

// L10 and L11: images that are still supported, pinned, and a pnpm that reads
// the settings the project writes for it.
const (
	// Alpine 3.19 went out of support in November 2025.
	runtimeAlpineImage = "alpine:3.24"
	// A Node 22 release rather than whatever node:22-alpine means on the day of
	// the build. Dependabot's docker entry moves it.
	nodeImage = "node:22.23-alpine"
	// nginx 1.27 was a mainline branch, closed when 1.28 became stable.
	nginxImage = "nginx:1.30-alpine"
	// pnpmVersion is the pnpm a project installs with everywhere: packageManager
	// in package.json, the Dockerfiles, and CI through packageManager.
	//
	// 10.33.4, not the newest 10.x. From 10.34.0 pnpm refuses a URL tarball whose
	// lockfile entry has no integrity, and pnpm 10 never writes one for the
	// SheetJS CDN tarball the web app's xlsx comes from, so `pnpm install` failed
	// on a project fresh from `grit new`. pnpm 11.9 computes the integrity on
	// download; moving to it is a pnpm 11 migration of its own.
	pnpmVersion = "10.33.4"

	pnpmPinOld = "# Pin pnpm — pnpm@latest started resolving to pnpm 11 which needs Node 22's\n" +
		"# node:sqlite builtin. Pinning here avoids surprise breakage on rebuilds.\n" +
		"RUN corepack enable && corepack prepare pnpm@9.15.0 --activate\n"
	pnpmPinSingleOld = "# Pin pnpm — pnpm@latest started resolving to pnpm 11 which needs Node 22's\n" +
		"# node:sqlite builtin; pinning avoids surprise breakage from that drift.\n" +
		"RUN corepack enable && corepack prepare pnpm@9.15.0 --activate\n"
	pnpmPinNew = "# The pnpm that package.json names in packageManager, which CI installs too.\n" +
		"# The image ran pnpm 9, which does not read onlyBuiltDependencies, so every\n" +
		"# dependency's install script ran during the build. Raise the two together.\n" +
		"RUN corepack enable && corepack prepare pnpm@" + pnpmVersion + " --activate\n"

	packageManagerOld = `"packageManager": "pnpm@10.0.0"`
	packageManagerNew = `"packageManager": "pnpm@` + pnpmVersion + `"`

	// The Next.js runner used to be FROM base, so the image that serves traffic
	// carried corepack, pnpm, npm and yarn, none of which it runs.
	nextRunnerOld = "# Run\nFROM base AS runner\nWORKDIR /app\n"
	nextRunnerNew = "# Run. A fresh Node image rather than the base stage, and without the package\n" +
		"# managers Node ships: the server is node and the standalone output, so pnpm,\n" +
		"# corepack, npm and yarn in the image that serves traffic are only something\n" +
		"# more to patch.\n" +
		"FROM " + nodeImage + " AS runner\n" +
		"RUN rm -rf /usr/local/lib/node_modules/npm /usr/local/lib/node_modules/corepack \\\n" +
		"    /usr/local/bin/npm /usr/local/bin/npx /usr/local/bin/corepack \\\n" +
		"    /usr/local/bin/yarn /usr/local/bin/yarnpkg /opt/yarn-*\n" +
		"WORKDIR /app\n"

	// Docker, Compose and most PaaS read a HEALTHCHECK to restart or stop routing
	// to a container that no longer answers. /api/health answers 200 while the
	// server is serving and reports each dependency in its body.
	apiHealthcheck = "# Unhealthy when the API stops answering. /api/health answers 200 while the\n" +
		"# server is serving, and reports each dependency in its body.\n" +
		"HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 \\\n" +
		"  CMD wget -q -O /dev/null \"http://127.0.0.1:${APP_PORT:-8080}/api/health\" || exit 1\n\n"
	// Any HTTP answer counts: a page that fails because the API is down is the
	// API's health to report, not the frontend's.
	nextHealthcheck = "# Unhealthy when the server stops answering. Any response counts: a page that\n" +
		"# fails because the API is down is the API's health check to report.\n" +
		"HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \\\n" +
		"  CMD node -e \"fetch('http://127.0.0.1:'+(process.env.PORT||3000)+'/').then(()=>process.exit(0),()=>process.exit(1))\"\n\n"
	viteHealthcheck = "# Unhealthy when nginx stops serving the app.\n" +
		"HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \\\n" +
		"  CMD wget -q -O /dev/null http://127.0.0.1:3000/ || exit 1\n\n"
)

// L14: development MinIO answers on this machine only, unless .env says so.
const (
	devMinioPortsOld = "      # Bound to all interfaces (not 127.0.0.1) so a phone/emulator on your LAN\n" +
		"      # can load uploaded images: stored URLs point at this host:9002 and the\n" +
		"      # Expo app rewrites \"localhost\" to your dev IP (apps/expo/lib/images.ts).\n" +
		"      - \"${MINIO_PORT:-9002}:9000\"\n" +
		"      - \"${MINIO_CONSOLE_PORT:-9003}:9001\"\n"
	devMinioPortsNew = "      # 127.0.0.1 like every other port here, unless MINIO_BIND_ADDRESS in .env\n" +
		"      # says otherwise. A phone or emulator on your LAN loads uploaded images from\n" +
		"      # this port (the Expo app rewrites \"localhost\" to your dev IP), so an Expo\n" +
		"      # project sets MINIO_BIND_ADDRESS=0.0.0.0. The console stays on this machine.\n" +
		"      - \"${MINIO_BIND_ADDRESS:-127.0.0.1}:${MINIO_PORT:-9002}:9000\"\n" +
		"      - \"127.0.0.1:${MINIO_CONSOLE_PORT:-9003}:9001\"\n"
	devMinioCredsCommentOld = "    # The root credentials are generated per project into .env. They were\n" +
		"    # minioadmin/minioadmin, on a port bound to every interface so a phone on\n" +
		"    # the LAN can load images, so anyone on the same network owned the bucket.\n"
	devMinioCredsCommentNew = "    # The root credentials are generated per project into .env. They were\n" +
		"    # minioadmin/minioadmin, on a port bound to every interface, so anyone on the\n" +
		"    # same network owned the bucket.\n"
)

// minioBindEnv is the .env setting for the MinIO port. An Expo app on a phone
// needs the LAN; nothing else does.
func minioBindEnv(opts Options) string {
	if opts.ShouldIncludeExpo() {
		return "# MinIO's API port on every interface, so the Expo app on a phone or emulator\n" +
			"# on your LAN can load uploaded images. The root credentials below are the\n" +
			"# only thing between that network and the bucket. 127.0.0.1 keeps it local.\n" +
			"MINIO_BIND_ADDRESS=0.0.0.0\n"
	}
	return "# MinIO's API port answers on this machine only. 0.0.0.0 opens it to your LAN,\n" +
		"# which an app on a phone needs to load uploaded images.\n" +
		"MINIO_BIND_ADDRESS=127.0.0.1\n"
}

// L12: the seeder's weak defaults are for development only.
const (
	seedAdminPasswordOld = "\t// Password resolution: SEED_ADMIN_PASSWORD wins so a real deployment can\n" +
		"\t// seed a strong credential. Otherwise fall back to the docs default\n" +
		"\t// \"admin123\" — but ONLY outside production, so the weak dev password can\n" +
		"\t// never slip into a prod database unnoticed.\n" +
		"\tpassword := os.Getenv(\"SEED_ADMIN_PASSWORD\")\n" +
		"\tif password == \"\" {\n" +
		"\t\tif os.Getenv(\"APP_ENV\") == \"production\" {\n" +
		"\t\t\treturn fmt.Errorf(\"refusing to seed the default admin in production: set SEED_ADMIN_PASSWORD to a strong password, or run the seeder with a non-production APP_ENV\")\n" +
		"\t\t}\n" +
		"\t\tpassword = \"admin123\"\n" +
		"\t}\n"
	seedAdminPasswordNew = "\t// SEED_ADMIN_PASSWORD wins, so a real deployment seeds a strong credential.\n" +
		"\t// The documented \"admin123\" is for APP_ENV=development only. It was the\n" +
		"\t// fallback for every environment except the one spelled \"production\", so a\n" +
		"\t// staging or \"prod\" database got an administrator anyone could guess.\n" +
		"\tpassword := os.Getenv(\"SEED_ADMIN_PASSWORD\")\n" +
		"\tif password == \"\" {\n" +
		"\t\tif !seedDevDefaults() {\n" +
		"\t\t\treturn fmt.Errorf(\"refusing to seed admin@example.com with the development password: APP_ENV is %q, so set SEED_ADMIN_PASSWORD to a password of at least 12 characters\", os.Getenv(\"APP_ENV\"))\n" +
		"\t\t}\n" +
		"\t\tpassword = \"admin123\"\n" +
		"\t} else if !seedDevDefaults() && len(password) < 12 {\n" +
		"\t\treturn fmt.Errorf(\"refusing to seed admin@example.com: SEED_ADMIN_PASSWORD is shorter than 12 characters and APP_ENV is %q\", os.Getenv(\"APP_ENV\"))\n" +
		"\t}\n"
	seedDemoGuardOld = "\t// Demo users are dev fixtures sharing a weak password — never seed them\n" +
		"\t// into a production database.\n" +
		"\tif os.Getenv(\"APP_ENV\") == \"production\" {\n" +
		"\t\tlog.Println(\"Skipping demo users in production\")\n" +
		"\t\treturn nil\n" +
		"\t}\n"
	seedDemoGuardNew = "\t// Demo users are development fixtures sharing a weak password, so they are\n" +
		"\t// seeded with APP_ENV=development and nowhere else.\n" +
		"\tif !seedDevDefaults() {\n" +
		"\t\tlog.Printf(\"Skipping demo users: APP_ENV is %q, not development\", os.Getenv(\"APP_ENV\"))\n" +
		"\t\treturn nil\n" +
		"\t}\n"
	seedDevDefaultsFunc = "\n// seedDevDefaults reports whether the weak development accounts may be seeded:\n" +
		"// only when APP_ENV says development. Unset, staging and prod are not.\n" +
		"func seedDevDefaults() bool {\n" +
		"\treturn os.Getenv(\"APP_ENV\") == \"development\"\n" +
		"}\n"

	// A single seeds itself on first boot, and did so in every environment but
	// production.
	singleAutoSeedOld = `seedEnabled := autoSeed == "true" || (autoSeed != "false" && cfg.AppEnv != "production")`
	singleAutoSeedNew = `seedEnabled := autoSeed == "true" || (autoSeed != "false" && cfg.AppEnv == "development")`

	singleAutoSeedCommentOld = "\t// First-boot seed: only runs when the users table is empty. Off by default\n" +
		"\t// in production unless AUTO_SEED=true is set explicitly.\n"
	singleAutoSeedCommentNew = "\t// First-boot seed: only runs when the users table is empty. On by default\n" +
		"\t// with APP_ENV=development only; AUTO_SEED=true turns it on elsewhere.\n"
)

// L13: the QR code on the 2FA setup screen. skip2/go-qrcode's last release was
// in 2020; boombuler/barcode is maintained and has no dependencies of its own.
const (
	qrModuleOld     = "github.com/skip2/go-qrcode"
	qrModuleNew     = "github.com/boombuler/barcode"
	qrModuleVersion = "v1.1.0"
)

// repairInfraHygiene brings a project up to L10 to L14 of the contact-app
// review. The compose files and users_seeder.go arrive whole through upgrade
// where they are unedited; the Dockerfiles, the root package.json, a single's
// main.go and go.mod do not, so they are repaired here.
func repairInfraHygiene(root string, opts Options) error {
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	apiRoot := opts.APIRoot(root)
	repairs := []struct {
		path string
		fn   func(string) (string, []string, []string)
	}{
		{filepath.Join(apiRoot, "Dockerfile"), repairAPIDockerfileSource},
		{filepath.Join(root, "Dockerfile"), repairSingleDockerfileSource},
		{filepath.Join(root, "apps", "web", "Dockerfile"), repairFrontendDockerfileSource},
		{filepath.Join(root, "apps", "admin", "Dockerfile"), repairFrontendDockerfileSource},
		{filepath.Join(root, "apps", "docs", "Dockerfile"), repairFrontendDockerfileSource},
		{filepath.Join(root, "package.json"), repairPackageManagerSource},
		{filepath.Join(root, "docker-compose.yml"), repairDevMinioBindSource},
		{filepath.Join(apiRoot, "internal", "database", "users_seeder.go"), repairUsersSeederSource},
	}
	if opts.Architecture == ArchSingle {
		// A single's API root is the project root, so its one Dockerfile is the
		// second entry.
		repairs = append(repairs[1:], struct {
			path string
			fn   func(string) (string, []string, []string)
		}{filepath.Join(root, "main.go"), repairSingleAutoSeedSource})
	} else {
		// Outside a single, a Dockerfile at the root is not one Grit wrote.
		repairs = append(repairs[:1], repairs[2:]...)
	}
	for _, r := range repairs {
		if !fileExists(r.path) {
			continue
		}
		if strings.HasSuffix(r.path, ".go") {
			err = repairSourceFile(root, m, r.path, r.fn)
		} else {
			err = repairTextFile(root, m, r.path, r.fn)
		}
		if err != nil {
			return err
		}
	}
	if err := ensureExpoMinioBind(root, opts); err != nil {
		return err
	}
	return replaceQRModule(apiRoot)
}

const expoMinioBindLines = "# 0.0.0.0 so the Expo app on a phone or emulator on your LAN can load uploaded images.\n" +
	"MINIO_BIND_ADDRESS=0.0.0.0\n"

var minioBindKeyRe = regexp.MustCompile(`(?m)^\s*(?:export\s+)?MINIO_BIND_ADDRESS\s*=`)

// ensureExpoMinioBind keeps an Expo project's MinIO on the LAN once compose
// defaults it to 127.0.0.1: the phone loads uploaded images from that port. It
// appends the setting to .env and .env.example when they do not define it, and
// never changes a value already there. Other projects get nothing.
func ensureExpoMinioBind(root string, opts Options) error {
	if !opts.IncludeExpo && !dirExists(filepath.Join(root, "apps", "expo")) {
		return nil
	}
	for _, name := range []string{".env", ".env.example"} {
		path := filepath.Join(root, name)
		raw, err := os.ReadFile(path)
		if err != nil {
			continue // a project that keeps its environment elsewhere
		}
		if minioBindKeyRe.Match(raw) {
			continue
		}
		src := string(raw)
		nl := "\n"
		if strings.Contains(src, "\r\n") {
			nl = "\r\n"
		}
		add := strings.ReplaceAll(expoMinioBindLines, "\n", nl)
		if src != "" && !strings.HasSuffix(src, "\n") {
			add = nl + add
		}
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, []byte(src+add), info.Mode().Perm()); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
		fmt.Printf("  ✓ %s: added MINIO_BIND_ADDRESS=0.0.0.0, so the Expo app still loads images now that MinIO defaults to 127.0.0.1\n", name)
	}
	return nil
}

var (
	alpineRuntimeRe = regexp.MustCompile(`(?m)^FROM alpine:3\.(\d+)\s*$`)
	nodeFloatingRe  = regexp.MustCompile(`(?m)^FROM node:22-alpine\b`)
	pnpmPrepareRe   = regexp.MustCompile(`corepack prepare pnpm@(\d+\.\d+\.\d+) --activate`)
	nginxOldRe      = regexp.MustCompile(`(?m)^FROM nginx:1\.(2[0-9])-alpine\b`)
)

// raiseImages moves the pins every Grit Dockerfile shares: the Alpine runtime,
// the floating Node tag and pnpm.
func raiseImages(src string) (string, []string) {
	out := src
	var fixed []string
	out = alpineRuntimeRe.ReplaceAllStringFunc(out, func(line string) string {
		minor := alpineRuntimeRe.FindStringSubmatch(line)[1]
		if !versionLess("3."+minor+".0", strings.TrimPrefix(runtimeAlpineImage, "alpine:")+".0") {
			return line
		}
		return "FROM " + runtimeAlpineImage
	})
	if out != src {
		fixed = append(fixed, "runtime image "+runtimeAlpineImage)
	}
	beforePnpm := out
	out = strings.Replace(out, pnpmPinOld, pnpmPinNew, 1)
	out = strings.Replace(out, pnpmPinSingleOld, pnpmPinNew, 1)
	out = pnpmPrepareRe.ReplaceAllStringFunc(out, func(match string) string {
		if !versionLess(pnpmPrepareRe.FindStringSubmatch(match)[1], pnpmVersion) {
			return match
		}
		return "corepack prepare pnpm@" + pnpmVersion + " --activate"
	})
	if out != beforePnpm {
		fixed = append(fixed, "pnpm "+pnpmVersion)
	}
	if next := nodeFloatingRe.ReplaceAllString(out, "FROM "+nodeImage); next != out {
		out = next
		fixed = append(fixed, nodeImage)
	}
	return out, fixed
}

// insertHealthcheck puts check above the final CMD line when the file has no
// HEALTHCHECK of its own.
func insertHealthcheck(src, check, cmd string) (string, bool) {
	if strings.Contains(src, "HEALTHCHECK") {
		return src, true
	}
	i := strings.LastIndex(src, "\n"+cmd)
	if i < 0 {
		return src, false
	}
	return src[:i+1] + check + src[i+1:], true
}

func dockerfileRepairResult(src, out string, fixed, warnings []string) (string, []string, []string) {
	if out == src {
		return src, nil, warnings
	}
	return out, fixed, warnings
}

func repairAPIDockerfileSource(src string) (string, []string, []string) {
	out, fixed := raiseImages(src)
	var warnings []string
	if next, ok := insertHealthcheck(out, apiHealthcheck, `CMD ["./server"]`); !ok {
		warnings = append(warnings, "no CMD [\"./server\"] to put a HEALTHCHECK above: add one that requests /api/health")
	} else if next != out {
		out = next
		fixed = append(fixed, "a HEALTHCHECK on /api/health")
	}
	return dockerfileRepairResult(src, out, fixed, warnings)
}

func repairSingleDockerfileSource(src string) (string, []string, []string) {
	return repairAPIDockerfileSource(src)
}

// repairFrontendDockerfileSource covers both frontend images: Next.js, which
// runs node, and Vite, which nginx serves.
func repairFrontendDockerfileSource(src string) (string, []string, []string) {
	out, fixed := raiseImages(src)
	var warnings []string
	if strings.Contains(out, "FROM nginx:") {
		if next := nginxOldRe.ReplaceAllString(out, "FROM "+nginxImage); next != out {
			out = next
			fixed = append(fixed, nginxImage)
		}
		if next, ok := insertHealthcheck(out, viteHealthcheck, `CMD ["nginx"`); ok && next != out {
			out = next
			fixed = append(fixed, "a HEALTHCHECK")
		}
		return dockerfileRepairResult(src, out, fixed, warnings)
	}
	if next := strings.Replace(out, nextRunnerOld, nextRunnerNew, 1); next != out {
		out = next
		fixed = append(fixed, "a runner image without pnpm, corepack, npm or yarn")
	} else if strings.Contains(out, "FROM base AS runner") {
		warnings = append(warnings, "the runner stage is FROM base and not the one Grit wrote: start it FROM "+nodeImage+" so the image that serves traffic carries no package managers")
	}
	if next, ok := insertHealthcheck(out, nextHealthcheck, `CMD ["node"`); !ok {
		warnings = append(warnings, "no CMD [\"node\", ...] to put a HEALTHCHECK above")
	} else if next != out {
		out = next
		fixed = append(fixed, "a HEALTHCHECK")
	}
	return dockerfileRepairResult(src, out, fixed, warnings)
}

var packageManagerRe = regexp.MustCompile(`"packageManager":\s*"pnpm@(\d+\.\d+\.\d+)"`)

// repairPackageManagerSource raises a pnpm 10 packageManager below pnpmVersion.
// Another major is the project's own choice and is left alone.
func repairPackageManagerSource(src string) (string, []string, []string) {
	m := packageManagerRe.FindStringSubmatch(src)
	if m == nil || !strings.HasPrefix(m[1], "10.") || !versionLess(m[1], pnpmVersion) {
		return src, nil, nil
	}
	out := strings.Replace(src, m[0], packageManagerNew, 1)
	return out, []string{"packageManager pnpm@" + pnpmVersion + ", the release the Dockerfiles install"}, nil
}

func repairDevMinioBindSource(src string) (string, []string, []string) {
	if strings.Contains(src, "${MINIO_BIND_ADDRESS:-127.0.0.1}") {
		return src, nil, nil
	}
	if !strings.Contains(src, devMinioPortsOld) {
		if strings.Contains(src, "\"${MINIO_PORT:-9002}:9000\"") {
			return src, nil, []string{"MinIO's ports are not the ones Grit wrote: prefix them with 127.0.0.1: so the bucket is not on your LAN"}
		}
		return src, nil, nil
	}
	out := strings.Replace(src, devMinioPortsOld, devMinioPortsNew, 1)
	out = strings.Replace(out, devMinioCredsCommentOld, devMinioCredsCommentNew, 1)
	return out, []string{"MinIO answers on 127.0.0.1 unless MINIO_BIND_ADDRESS says otherwise"}, nil
}

func repairUsersSeederSource(src string) (string, []string, []string) {
	if strings.Contains(src, "func seedDevDefaults() bool") {
		return src, nil, nil
	}
	if !strings.Contains(src, seedAdminPasswordOld) || !strings.Contains(src, seedDemoGuardOld) {
		if strings.Contains(src, "\"admin123\"") {
			return src, nil, []string{"the seeder is not the one Grit wrote: seed admin123 and the demo users only when APP_ENV=development"}
		}
		return src, nil, nil
	}
	out := strings.Replace(src, seedAdminPasswordOld, seedAdminPasswordNew, 1)
	out = strings.Replace(out, seedDemoGuardOld, seedDemoGuardNew, 1)
	out = strings.TrimRight(out, "\n") + "\n" + seedDevDefaultsFunc
	return out, []string{"admin123 and the demo users are seeded only with APP_ENV=development"}, nil
}

func repairSingleAutoSeedSource(src string) (string, []string, []string) {
	if !strings.Contains(src, singleAutoSeedOld) {
		return src, nil, nil
	}
	out := strings.Replace(src, singleAutoSeedOld, singleAutoSeedNew, 1)
	out = strings.Replace(out, singleAutoSeedCommentOld, singleAutoSeedCommentNew, 1)
	return out,
		[]string{"the first-boot seed runs by default with APP_ENV=development only"}, nil
}

// replaceQRModule swaps skip2/go-qrcode for boombuler/barcode in go.mod once
// handlers/totp.go uses the new one. The old requirement goes only when no Go
// file still imports it, so a device-pairing handler keeps building.
func replaceQRModule(apiRoot string) error {
	gomod := filepath.Join(apiRoot, "go.mod")
	data, err := os.ReadFile(gomod)
	if err != nil {
		return nil
	}
	uses := func(module string) (bool, error) {
		found := false
		err := filepath.WalkDir(apiRoot, func(path string, d os.DirEntry, err error) error {
			if err != nil || found {
				return err
			}
			if d.IsDir() && (d.Name() == "node_modules" || d.Name() == "vendor" || strings.HasPrefix(d.Name(), ".")) && path != apiRoot {
				return filepath.SkipDir
			}
			if !d.IsDir() && strings.HasSuffix(path, ".go") && fileContains(path, "\""+module) {
				found = true
			}
			return nil
		})
		return found, err
	}
	newUsed, err := uses(qrModuleNew)
	if err != nil {
		return err
	}
	if newUsed && !strings.Contains(string(data), qrModuleNew+" ") {
		if err := goGet(apiRoot, qrModuleNew+"@"+qrModuleVersion); err != nil {
			fmt.Printf("  ⚠ could not add %s, so run `go get %s@%s` in %s: %v\n", qrModuleNew, qrModuleNew, qrModuleVersion, apiRoot, err)
			return nil
		}
		fmt.Printf("  ✓ go.mod: %s %s for the 2FA QR code\n", qrModuleNew, qrModuleVersion)
	}
	if !strings.Contains(string(data), qrModuleOld+" ") {
		return nil
	}
	oldUsed, err := uses(qrModuleOld)
	if err != nil || oldUsed {
		return err
	}
	if err := goModEdit(apiRoot, "-droprequire="+qrModuleOld); err != nil {
		fmt.Printf("  ⚠ could not drop %s from go.mod, so run `go mod tidy` in %s: %v\n", qrModuleOld, apiRoot, err)
		return nil
	}
	fmt.Printf("  ✓ go.mod: dropped %s, unmaintained since 2020 and no longer imported\n", qrModuleOld)
	return nil
}
