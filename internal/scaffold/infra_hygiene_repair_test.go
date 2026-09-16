package scaffold

import (
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// checkRepair asserts that fn turns old into want, and that want needs nothing.
func checkRepair(t *testing.T, name string, fn func(string) (string, []string, []string), old, want string) {
	t.Helper()
	if old == want {
		t.Fatalf("%s: the reconstructed old file is the template, so the test proves nothing", name)
	}
	out, fixed, warnings := fn(old)
	if out != want {
		t.Errorf("%s: the repair does not produce the template.\n--- got ---\n%s\n--- want ---\n%s", name, out, want)
	}
	if len(fixed) == 0 || len(warnings) != 0 {
		t.Errorf("%s: fixed %v, warnings %v", name, fixed, warnings)
	}
	if again, fixed, warnings := fn(want); again != want || len(fixed) != 0 || len(warnings) != 0 {
		t.Errorf("%s: the repair is not idempotent: %v %v", name, fixed, warnings)
	}
}

func TestDockerfilesUseSupportedPinnedImages(t *testing.T) {
	files := map[string]string{
		"api":    dockerfileAPI(),
		"single": dockerfileSingle(),
		"next":   dockerfileNextJS("web"),
		"vite":   dockerfileVite("admin"),
	}
	for name, src := range files {
		for _, stale := range []string{"alpine:3.19", "node:22-alpine", "pnpm@9", "nginx:1.27"} {
			if strings.Contains(src, stale) {
				t.Errorf("%s Dockerfile still has %q", name, stale)
			}
		}
		if !strings.Contains(src, "HEALTHCHECK") {
			t.Errorf("%s Dockerfile has no HEALTHCHECK", name)
		}
		if name != "api" && !strings.Contains(src, "pnpm@"+pnpmVersion) {
			t.Errorf("%s Dockerfile does not install pnpm %s", name, pnpmVersion)
		}
	}
	for _, name := range []string{"api", "single"} {
		if !strings.Contains(files[name], "FROM "+runtimeAlpineImage) || !strings.Contains(files[name], "/api/health") {
			t.Errorf("%s Dockerfile: want %s and a health check on /api/health", name, runtimeAlpineImage)
		}
	}
	if strings.Contains(files["next"], "FROM base AS runner") || !strings.Contains(files["next"], "FROM "+nodeImage+" AS runner") {
		t.Error("the Next.js runner still inherits pnpm from the base stage")
	}
	if !strings.Contains(files["vite"], "FROM "+nginxImage) {
		t.Errorf("the Vite runner is not %s", nginxImage)
	}
	if strings.Contains(files["next"], "%!") {
		t.Error("the Next.js Dockerfile has a formatting verb error")
	}
	if !strings.Contains(rootPackageJSON(Options{ProjectName: "app"}), packageManagerNew) {
		t.Error("package.json does not name the pnpm the Dockerfiles install")
	}
}

func TestDockerfileRepairs(t *testing.T) {
	api := dockerfileAPI()
	checkRepair(t, "api", repairAPIDockerfileSource,
		strings.Replace(strings.Replace(api, "FROM "+runtimeAlpineImage, "FROM alpine:3.19", 1), apiHealthcheck, "", 1), api)

	single := dockerfileSingle()
	oldSingle := strings.Replace(single, "FROM "+runtimeAlpineImage, "FROM alpine:3.19", 1)
	oldSingle = strings.Replace(oldSingle, "FROM "+nodeImage, "FROM node:22-alpine", 1)
	oldSingle = strings.Replace(oldSingle, pnpmPinNew, pnpmPinSingleOld, 1)
	oldSingle = strings.Replace(oldSingle, apiHealthcheck, "", 1)
	checkRepair(t, "single", repairSingleDockerfileSource, oldSingle, single)

	for _, app := range []string{"web", "admin", "docs"} {
		next := dockerfileNextJS(app)
		old := strings.Replace(next, nextRunnerNew, nextRunnerOld, 1)
		old = strings.Replace(old, "FROM "+nodeImage+" AS base", "FROM node:22-alpine AS base", 1)
		old = strings.Replace(old, pnpmPinNew, pnpmPinOld, 1)
		old = strings.Replace(old, nextHealthcheck, "", 1)
		checkRepair(t, "next "+app, repairFrontendDockerfileSource, old, next)
	}

	vite := dockerfileVite("web")
	oldVite := strings.Replace(vite, "FROM "+nodeImage, "FROM node:22-alpine", 1)
	oldVite = strings.Replace(oldVite, pnpmPinNew, pnpmPinOld, 1)
	oldVite = strings.Replace(oldVite, "FROM "+nginxImage, "FROM nginx:1.27-alpine", 1)
	oldVite = strings.Replace(oldVite, viteHealthcheck, "", 1)
	checkRepair(t, "vite", repairFrontendDockerfileSource, oldVite, vite)

	// An edited Dockerfile still gets its pins raised, and a newer pin stays.
	edited := "FROM node:22-alpine AS base\nRUN corepack enable && corepack prepare pnpm@9.1.0 --activate\nFROM base AS runner\nCMD [\"node\", \"server.js\"]\n"
	out, _, warnings := repairFrontendDockerfileSource(edited)
	if !strings.Contains(out, "pnpm@"+pnpmVersion) || !strings.Contains(out, "FROM "+nodeImage) || !strings.Contains(out, "HEALTHCHECK") {
		t.Errorf("an edited Dockerfile was not raised:\n%s", out)
	}
	if len(warnings) != 1 {
		t.Errorf("an edited runner stage should be named, got %v", warnings)
	}
	newer := "FROM alpine:3.30\nRUN corepack prepare pnpm@11.0.0 --activate\nHEALTHCHECK CMD true\nCMD [\"./server\"]\n"
	if out, fixed, _ := repairAPIDockerfileSource(newer); out != newer || len(fixed) != 0 {
		t.Errorf("newer pins were changed: %v", fixed)
	}
}

func TestPackageManagerRepair(t *testing.T) {
	fresh := rootPackageJSON(Options{ProjectName: "app"})
	checkRepair(t, "package.json", repairPackageManagerSource, strings.Replace(fresh, packageManagerNew, packageManagerOld, 1), fresh)
	other := `{"packageManager": "pnpm@9.15.0"}`
	if out, _, _ := repairPackageManagerSource(other); out != other {
		t.Error("another pnpm major is the project's choice and should be left alone")
	}
}

func TestDevMinioBindsLocally(t *testing.T) {
	compose := dockerCompose(Options{ProjectName: "app"})
	if strings.Contains(compose, "- \"${MINIO_PORT:-9002}:9000\"") || !strings.Contains(compose, "${MINIO_BIND_ADDRESS:-127.0.0.1}:${MINIO_PORT:-9002}:9000") ||
		!strings.Contains(compose, "127.0.0.1:${MINIO_CONSOLE_PORT:-9003}:9001") {
		t.Error("development MinIO is still published on every interface")
	}
	old := strings.Replace(compose, devMinioPortsNew, devMinioPortsOld, 1)
	old = strings.Replace(old, devMinioCredsCommentNew, devMinioCredsCommentOld, 1)
	checkRepair(t, "docker-compose.yml", repairDevMinioBindSource, old, compose)

	if _, _, warnings := repairDevMinioBindSource("    ports:\n      - \"${MINIO_PORT:-9002}:9000\"\n"); len(warnings) != 1 {
		t.Error("an edited MinIO port block should be named")
	}

	if env := envFile(Options{ProjectName: "app", Theme: "atlas"}); !strings.Contains(env, "\nMINIO_BIND_ADDRESS=127.0.0.1\n") {
		t.Error(".env without Expo should keep MinIO on 127.0.0.1")
	}
	if env := envFile(Options{ProjectName: "app", Theme: "atlas", IncludeExpo: true}); !strings.Contains(env, "\nMINIO_BIND_ADDRESS=0.0.0.0\n") {
		t.Error(".env with Expo should open MinIO to the LAN")
	}
}

func TestExpoMinioBindIsAppended(t *testing.T) {
	setup := func(t *testing.T, expo bool, env string) string {
		t.Helper()
		root := t.TempDir()
		if expo {
			if err := os.MkdirAll(filepath.Join(root, "apps", "expo"), 0o755); err != nil {
				t.Fatal(err)
			}
		}
		for _, name := range []string{".env", ".env.example"} {
			if err := os.WriteFile(filepath.Join(root, name), []byte(env), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return root
	}
	readEnv := func(t *testing.T, root, name string) string {
		t.Helper()
		b, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}

	t.Run("expo without the key gets it once", func(t *testing.T) {
		root := setup(t, true, "APP_ENV=development\nMINIO_PORT=9002")
		for i := 0; i < 2; i++ {
			if err := ensureExpoMinioBind(root, Options{}); err != nil {
				t.Fatal(err)
			}
		}
		for _, name := range []string{".env", ".env.example"} {
			got := readEnv(t, root, name)
			if want := "APP_ENV=development\nMINIO_PORT=9002\n" + expoMinioBindLines; got != want {
				t.Errorf("%s:\n%q\nwant\n%q", name, got, want)
			}
		}
	})

	t.Run("expo with the key keeps its value", func(t *testing.T) {
		for _, value := range []string{"127.0.0.1", "192.168.1.20", ""} {
			env := "MINIO_BIND_ADDRESS=" + value + "\n"
			root := setup(t, true, env)
			if err := ensureExpoMinioBind(root, Options{}); err != nil {
				t.Fatal(err)
			}
			if got := readEnv(t, root, ".env"); got != env {
				t.Errorf("a value the developer set was changed: %q", got)
			}
		}
	})

	t.Run("a project without expo is untouched", func(t *testing.T) {
		env := "APP_ENV=development\n"
		root := setup(t, false, env)
		if err := ensureExpoMinioBind(root, Options{}); err != nil {
			t.Fatal(err)
		}
		if readEnv(t, root, ".env") != env || readEnv(t, root, ".env.example") != env {
			t.Error("a project without Expo gained MINIO_BIND_ADDRESS")
		}
	})
}

func TestSeedDefaultsAreDevelopmentOnly(t *testing.T) {
	seeder := apiUsersSeederGo()
	if strings.Contains(seeder, `os.Getenv("APP_ENV") == "production"`) || !strings.Contains(seeder, seedDevDefaultsFunc) {
		t.Error("the seeder still seeds admin123 outside production")
	}
	if _, err := format.Source([]byte(strings.ReplaceAll(seeder, "{{MODULE}}", "example.com/app"))); err != nil {
		t.Fatalf("the seeder template is not valid Go: %v", err)
	}
	old := strings.Replace(seeder, seedAdminPasswordNew, seedAdminPasswordOld, 1)
	old = strings.Replace(old, seedDemoGuardNew, seedDemoGuardOld, 1)
	old = strings.Replace(old, seedDevDefaultsFunc, "", 1)
	checkRepair(t, "users_seeder.go", repairUsersSeederSource, old, seeder)

	if _, _, warnings := repairUsersSeederSource("package database\nvar p = \"admin123\"\n"); len(warnings) != 1 {
		t.Error("an edited seeder should be named")
	}

	main := singleMainGo(Options{ProjectName: "app", Architecture: ArchSingle})
	if !strings.Contains(main, singleAutoSeedNew) {
		t.Fatal("a single still seeds itself outside development")
	}
	oldMain := strings.Replace(main, singleAutoSeedNew, singleAutoSeedOld, 1)
	oldMain = strings.Replace(oldMain, singleAutoSeedCommentNew, singleAutoSeedCommentOld, 1)
	checkRepair(t, "main.go", repairSingleAutoSeedSource, oldMain, main)
}

func TestQRCodeLibrary(t *testing.T) {
	handler := totpHandlerGo()
	if strings.Contains(handler, qrModuleOld) || !strings.Contains(handler, qrModuleNew+"/qr") {
		t.Error("the 2FA handler still uses skip2/go-qrcode")
	}
	if _, err := format.Source([]byte(strings.ReplaceAll(handler, "{{MODULE}}", "example.com/app"))); err != nil {
		t.Errorf("the 2FA handler template is not valid Go: %v", err)
	}
	gomod := apiGoMod(Options{ProjectName: "app"})
	if strings.Contains(gomod, qrModuleOld) || !strings.Contains(gomod, qrModuleNew+" "+qrModuleVersion) {
		t.Error("go.mod still requires skip2/go-qrcode")
	}
}

func TestReplaceQRModule(t *testing.T) {
	dir := t.TempDir()
	write := func(rel, body string) {
		t.Helper()
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	write("go.mod", "module app\n\nrequire (\n\t"+qrModuleOld+" v0.0.0-20200617195104-da1b6568686e\n)\n")
	write("internal/handlers/totp.go", "package handlers\n\nimport \""+qrModuleNew+"/qr\"\n")

	var got, edits []string
	oldGet, oldEdit := goGet, goModEdit
	t.Cleanup(func() { goGet, goModEdit = oldGet, oldEdit })
	goGet = func(_ string, specs ...string) error { got = append(got, specs...); return nil }
	goModEdit = func(_ string, args ...string) error { edits = append(edits, args...); return nil }

	if err := replaceQRModule(dir); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != qrModuleNew+"@"+qrModuleVersion {
		t.Errorf("go get: %v", got)
	}
	if len(edits) != 1 || edits[0] != "-droprequire="+qrModuleOld {
		t.Errorf("go mod edit: %v", edits)
	}

	// A file still importing the old module keeps its requirement.
	got, edits = nil, nil
	write("internal/handlers/pairing.go", "package handlers\n\nimport qrcode \""+qrModuleOld+"\"\n")
	if err := replaceQRModule(dir); err != nil {
		t.Fatal(err)
	}
	if len(edits) != 0 {
		t.Errorf("the old module was dropped while still imported: %v", edits)
	}
}

func TestDependabotWatchesDockerfiles(t *testing.T) {
	triple := dependabotYAML(Options{ProjectName: "app", Architecture: ArchTriple, Frontend: FrontendNext})
	if !strings.Contains(triple, "package-ecosystem: docker") || !strings.Contains(triple, "\"/apps/api\"") {
		t.Errorf("dependabot.yml does not watch the Dockerfiles:\n%s", triple)
	}
	single := dependabotYAML(Options{ProjectName: "app", Architecture: ArchSingle})
	if !strings.Contains(single, "package-ecosystem: docker\n    directories:\n      - \"/\"\n") {
		t.Errorf("a single's Dockerfile is at the root:\n%s", single)
	}
}
