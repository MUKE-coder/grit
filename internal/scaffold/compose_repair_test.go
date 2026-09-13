package scaffold

import (
	"strings"
	"testing"
)

// The compose templates carry no default MinIO credentials, and the production
// side containers take no env_file.
func TestComposeTemplatesHoldNoDefaultSecrets(t *testing.T) {
	opts := Options{ProjectName: "shop"}
	for name, src := range map[string]string{"docker-compose.yml": dockerCompose(opts), "docker-compose.prod.yml": dockerComposeProd(opts)} {
		// The comments explain what minioadmin was, so look for it as a value.
		if strings.Contains(src, ": minioadmin") || strings.Contains(src, ":-minioadmin") {
			t.Errorf("%s still defaults to minioadmin", name)
		}
	}
	prod := dockerComposeProd(opts)
	for _, svc := range []string{"postgres", "minio"} {
		start, end, ok := serviceBlock(prod, svc)
		if !ok {
			t.Fatalf("no %s service in the production compose file", svc)
		}
		if strings.Contains(prod[start:end], "    env_file:") {
			t.Errorf("the production %s service still takes all of .env", svc)
		}
	}
	if start, end, ok := serviceBlock(prod, "api"); !ok || !strings.Contains(prod[start:end], "    env_file:") {
		t.Error("the api service lost its env_file")
	}
	if strings.Contains(prod, "POSTGRES_PASSWORD:-grit") {
		t.Error("the production Postgres password still falls back to grit")
	}
}

func TestEnvGeneratesMinioCredentials(t *testing.T) {
	a, b := envFile(Options{ProjectName: "shop"}), envFile(Options{ProjectName: "shop"})
	if strings.Contains(a, "minioadmin") || strings.Contains(a, "{{MINIO_") {
		t.Error(".env still carries minioadmin or an unfilled placeholder")
	}
	secret := func(s string) string {
		for _, line := range strings.Split(s, "\n") {
			if strings.HasPrefix(line, "MINIO_SECRET_KEY=") {
				return strings.TrimPrefix(line, "MINIO_SECRET_KEY=")
			}
		}
		return ""
	}
	if secret(a) == "" || len(secret(a)) < 32 || secret(a) == secret(b) {
		t.Errorf("MINIO_SECRET_KEY is not generated per project: %d characters", len(secret(a)))
	}
	if strings.Contains(envExampleFile(Options{ProjectName: "shop"}), secret(a)) {
		t.Error("the generated MinIO secret leaked into .env.example")
	}
}

func TestRepairComposeFiles(t *testing.T) {
	dev := "services:\n  minio:\n    image: minio/minio\n    environment:\n" + devMinioCredentials + "    volumes:\n      - minio-data:/data\n"
	out, fixed, _ := repairDevComposeSource(dev)
	if strings.Contains(out, "minioadmin") || len(fixed) != 1 {
		t.Errorf("dev compose not repaired:\n%s", out)
	}

	prod := "services:\n  api:\n    env_file:\n      - .env\n" +
		"  postgres:\n    image: postgres:16-alpine\n    env_file:\n      - .env\n    environment:\n      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:-grit}\n" +
		"  minio:\n    image: minio/minio\n    env_file:\n      - .env\n    environment:\n      MINIO_ROOT_USER: ${MINIO_ACCESS_KEY:-minioadmin}\n      MINIO_ROOT_PASSWORD: ${MINIO_SECRET_KEY:-minioadmin}\n" +
		"\nnetworks:\n  shop:\n"
	out, fixed, _ = repairProdComposeSource(prod)
	if strings.Count(out, "env_file") != 1 {
		t.Errorf("want only the api's env_file left:\n%s", out)
	}
	if strings.Contains(out, "minioadmin") || strings.Contains(out, ":-grit") || len(fixed) != 3 {
		t.Errorf("prod compose not repaired (%v):\n%s", fixed, out)
	}
	if again, fixed, _ := repairProdComposeSource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed the production compose file again")
	}
}

// config.go from v3.243.0 to v3.247.0 has checkSecrets without MinIO.
func TestRepairConfigAddsMinioToCheckSecrets(t *testing.T) {
	current := apiConfigGo()
	old := strings.Replace(strings.Replace(current, checkSecretsMinioLine, "", 1), "\t\"minioadmin\": true,\n", "", 1)
	if strings.Contains(old, `{"MINIO_SECRET_KEY", cfg.Storage.SecretKey`) {
		t.Fatal("could not reconstruct the old config")
	}
	out, fixed, warn := repairConfigMinioSecretSource(old)
	if len(warn) > 0 || len(fixed) != 1 || out != current {
		t.Errorf("repair did not reproduce the template (fixed %v, warn %v)", fixed, warn)
	}
	if out, fixed, _ := repairConfigMinioSecretSource(current); out != current || len(fixed) > 0 {
		t.Error("the scaffold's config still needs the repair")
	}
}
