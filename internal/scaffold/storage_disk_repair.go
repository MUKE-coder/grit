package scaffold

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// Issue 003, C1: the storage Disk interface and STORAGE_DRIVER=local.
//
// internal/storage arrives whole from writeStorageFiles. config.go, main.go and
// routes.go are the developer's, so the driver choice, the local settings, the
// construction and the file route are added on the exact text Grit wrote.

const (
	// StorageConfig gains the local driver's settings.
	configStorageFieldsAnchor = "\tPublicURL string\n}\n"
	configStorageFields       = "\tPublicURL string\n\n" +
		"\t// LocalRoot is the directory STORAGE_DRIVER=local keeps files in\n" +
		"\t// (STORAGE_LOCAL_ROOT). URLSecret signs its temporary URLs\n" +
		"\t// (STORAGE_URL_SECRET, or JWT_SECRET when that is unset).\n" +
		"\tLocalRoot string\n" +
		"\tURLSecret string\n" +
		"}\n"

	configStorageDriverFieldOld = "\tStorageDriver string        // \"minio\", \"s3\", \"r2\", or \"b2\"\n"
	configStorageDriverFieldNew = "\tStorageDriver string        // \"local\", \"minio\", \"s3\", \"r2\", or \"b2\"\n"

	configStorageDriverOld = "\tstorageDriver := getEnv(\"STORAGE_DRIVER\", \"minio\")\n"
	configStorageDriverNew = "\tstorageDriver := resolveStorageDriver()\n"

	// Production refuses the local driver unless told otherwise, as it refuses SQLite.
	configLocalStorageCheckAnchor = "\t// DatabaseURL is always populated by resolveDatabaseURL()"
	configLocalStorageCheck       = "\t// The local driver keeps files on this server's disk: a second replica cannot\n" +
		"\t// see them, and a redeploy that replaces the container loses them. A\n" +
		"\t// single-server deployment with the directory on a volume it backs up says so.\n" +
		"\tif cfg.StorageDriver == \"local\" && cfg.AppEnv == \"production\" && getEnv(\"ALLOW_LOCAL_STORAGE_IN_PRODUCTION\", \"false\") != \"true\" {\n" +
		"\t\treturn nil, fmt.Errorf(\"APP_ENV=production is using STORAGE_DRIVER=local: use minio, s3, r2 or b2, or set ALLOW_LOCAL_STORAGE_IN_PRODUCTION=true if STORAGE_LOCAL_ROOT is on a volume you back up\")\n" +
		"\t}\n\n"

	configResolveStorageAnchor     = "// resolveStorage returns the StorageConfig for the active driver.\n"
	configResolveStorageDriverFunc = "// resolveStorageDriver picks the storage driver from STORAGE_DRIVER.\n" +
		"//\n" +
		"// Outside production, minio with no MINIO_ACCESS_KEY becomes local, so a new\n" +
		"// project stores uploads before Docker is running instead of answering each\n" +
		"// one with STORAGE_UNAVAILABLE. Production never falls back: a server that\n" +
		"// lost its credentials should say so, not start writing to its own disk.\n" +
		"func resolveStorageDriver() string {\n" +
		"\tdriver := getEnv(\"STORAGE_DRIVER\", \"minio\")\n" +
		"\tif driver == \"minio\" && getEnv(\"MINIO_ACCESS_KEY\", \"\") == \"\" && getEnv(\"APP_ENV\", \"production\") != \"production\" {\n" +
		"\t\tlog.Println(\"MinIO has no credentials (MINIO_ACCESS_KEY), so files are kept on the local disk. Set STORAGE_DRIVER=local to make that the choice, or set the MinIO credentials to use MinIO\")\n" +
		"\t\treturn \"local\"\n" +
		"\t}\n" +
		"\treturn driver\n" +
		"}\n\n"

	configLocalCaseAnchor = "\tdefault: // minio\n"
	configLocalCase       = "\tcase \"local\":\n" +
		"\t\treturn StorageConfig{\n" +
		"\t\t\tLocalRoot: getEnv(\"STORAGE_LOCAL_ROOT\", \"storage/app\"),\n" +
		"\t\t\tURLSecret: firstNonEmpty(os.Getenv(\"STORAGE_URL_SECRET\"), os.Getenv(\"JWT_SECRET\")),\n" +
		"\t\t\t// The API serves these files itself, from the route routes.Setup mounts.\n" +
		"\t\t\tPublicURL: firstNonEmpty(os.Getenv(\"STORAGE_PUBLIC_URL\"), strings.TrimRight(getEnv(\"APP_URL\", \"http://localhost:8080\"), \"/\")+\"/files\"),\n" +
		"\t\t}\n"

	// main.go builds the store the driver names.
	mainStorageInitOld = "\t// File storage (S3-compatible)\n" +
		"\tvar storageService *storage.Storage\n" +
		"\tif cfg.Storage.Endpoint != \"\" && cfg.Storage.AccessKey != \"\" {\n" +
		"\t\ts, err := storage.New(cfg.Storage)\n" +
		"\t\tif err != nil {\n" +
		"\t\t\tlog.Printf(\"Warning: Storage unavailable: %v (uploads disabled)\", err)\n" +
		"\t\t} else {\n" +
		"\t\t\tstorageService = s\n" +
		"\t\t\tlog.Println(\"File storage connected\")\n" +
		"\t\t}\n" +
		"\t}\n"
	singleStorageInitOld = "\t// S3-compatible storage\n" +
		"\tvar storageService *storage.Storage\n" +
		"\ts, err := storage.New(cfg.Storage)\n" +
		"\tif err != nil {\n" +
		"\t\tlog.Printf(\"Warning: Storage unavailable: %v\", err)\n" +
		"\t} else {\n" +
		"\t\tstorageService = s\n" +
		"\t\tlog.Println(\"Storage configured\")\n" +
		"\t}\n"
	mainStorageInit = "\t// File storage. STORAGE_DRIVER=local keeps files in a directory on this\n" +
		"\t// server and serves them from /files; minio, s3, r2 and b2 use a bucket.\n" +
		"\tvar storageService *storage.Storage\n" +
		"\tif cfg.StorageDriver == \"local\" {\n" +
		"\t\ts, err := storage.NewLocal(storage.LocalConfig{\n" +
		"\t\t\tRoot:      cfg.Storage.LocalRoot,\n" +
		"\t\t\tPublicURL: cfg.Storage.PublicURL,\n" +
		"\t\t\tSecret:    cfg.Storage.URLSecret,\n" +
		"\t\t})\n" +
		"\t\tif err != nil {\n" +
		"\t\t\tlog.Printf(\"Warning: Storage unavailable: %v (uploads disabled)\", err)\n" +
		"\t\t} else {\n" +
		"\t\t\tstorageService = s\n" +
		"\t\t\tlog.Printf(\"File storage: local disk at %s, served from %s\", cfg.Storage.LocalRoot, cfg.Storage.PublicURL)\n" +
		"\t\t}\n" +
		"\t} else if cfg.Storage.Endpoint != \"\" && cfg.Storage.AccessKey != \"\" {\n" +
		"\t\ts, err := storage.New(cfg.Storage)\n" +
		"\t\tif err != nil {\n" +
		"\t\t\tlog.Printf(\"Warning: Storage unavailable: %v (uploads disabled)\", err)\n" +
		"\t\t} else {\n" +
		"\t\t\tstorageService = s\n" +
		"\t\t\tlog.Println(\"File storage connected\")\n" +
		"\t\t}\n" +
		"\t}\n"

	// routes.go serves the local driver's files.
	routesLocalFilesAnchor = "\t// Health check\n\t// /api/health probes every infrastructure dependency"
	routesLocalFilesBlock  = "\t// Files kept by STORAGE_DRIVER=local. Keys under storage.PublicPrefixes are\n" +
		"\t// served to anyone, every other key only through a signed temporary URL. A\n" +
		"\t// bucket serves its own files, so nothing is mounted for one.\n" +
		"\tif svc.Storage != nil {\n" +
		"\t\tif fileServer := svc.Storage.FileServer(); fileServer != nil {\n" +
		"\t\t\tr.GET(storage.LocalFilesRoute, gin.WrapH(fileServer))\n" +
		"\t\t\tr.HEAD(storage.LocalFilesRoute, gin.WrapH(fileServer))\n" +
		"\t\t}\n" +
		"\t}\n\n"

	// The local root is not source.
	apiGitignoreLocalStorage = "\n# Files kept by STORAGE_DRIVER=local (STORAGE_LOCAL_ROOT)\n/storage/\n"

	// The admin and web apps load local files from the API origin.
	cspImgSrcOld = "\"img-src 'self' data: blob: https: \" + STORAGE_ORIGIN,"
	cspImgSrcNew = "\"img-src 'self' data: blob: https: \" + API_ORIGIN + \" \" + STORAGE_ORIGIN,"
)

// repairStorageConfigSource adds the local driver to config.go.
func repairStorageConfigSource(src string) (string, []string, []string) {
	if strings.Contains(src, "func resolveStorageDriver()") || !strings.Contains(src, "func resolveStorage(driver string) StorageConfig") {
		return src, nil, nil
	}
	warn := []string{"config.go is not the file Grit wrote: add STORAGE_DRIVER=local (LocalRoot and URLSecret on StorageConfig, a \"local\" case in resolveStorage) to keep files on the local disk"}
	for _, anchor := range []string{configStorageFieldsAnchor, configStorageDriverOld, configLocalStorageCheckAnchor, configResolveStorageAnchor, configLocalCaseAnchor} {
		if strings.Count(src, anchor) != 1 {
			return src, nil, warn
		}
	}
	out := strings.Replace(src, configStorageFieldsAnchor, configStorageFields, 1)
	out = strings.Replace(out, configStorageDriverFieldOld, configStorageDriverFieldNew, 1)
	out = strings.Replace(out, configStorageDriverOld, configStorageDriverNew, 1)
	out = strings.Replace(out, configLocalStorageCheckAnchor, configLocalStorageCheck+configLocalStorageCheckAnchor, 1)
	out = strings.Replace(out, configResolveStorageAnchor, configResolveStorageDriverFunc+configResolveStorageAnchor, 1)
	out = strings.Replace(out, configLocalCaseAnchor, configLocalCase+configLocalCaseAnchor, 1)
	out, ok := withImports(out, "fmt", "log", "os", "strings")
	if !ok {
		return src, nil, []string{"could not add imports to config.go"}
	}
	return out, []string{"STORAGE_DRIVER=local keeps files on this machine, and is the default outside production when MinIO has no credentials"}, nil
}

// repairStorageInitSource builds the store STORAGE_DRIVER names in main.go.
func repairStorageInitSource(src string) (string, []string, []string) {
	if strings.Contains(src, "storage.NewLocal(") || !strings.Contains(src, "storage.New(cfg.Storage)") {
		return src, nil, nil
	}
	for _, old := range []string{mainStorageInitOld, singleStorageInitOld} {
		if strings.Count(src, old) == 1 {
			out := strings.Replace(src, old, mainStorageInit, 1)
			return out, []string{"STORAGE_DRIVER=local builds a local disk instead of a bucket client"}, nil
		}
	}
	return src, nil, []string{"main.go is not the file Grit wrote: build storage.NewLocal(storage.LocalConfig{...}) when cfg.StorageDriver is \"local\", or that driver stores nothing"}
}

// repairLocalFilesRouteSource mounts the local driver's file route in routes.go.
func repairLocalFilesRouteSource(src string) (string, []string, []string) {
	if strings.Contains(src, "storage.LocalFilesRoute") || !strings.Contains(src, "func Setup(") {
		return src, nil, nil
	}
	// Above the queue snapshot when the health repair has added it, which it
	// puts directly above the health check.
	anchor := routesLocalFilesAnchor
	if strings.Contains(src, healthQueueStatsSetup) {
		anchor = healthQueueStatsSetup
	}
	if strings.Count(src, anchor) != 1 {
		return src, nil, []string{"routes.go is not the file Grit wrote: mount svc.Storage.FileServer() at storage.LocalFilesRoute, or files kept by STORAGE_DRIVER=local are never served"}
	}
	out := strings.Replace(src, anchor, routesLocalFilesBlock+anchor, 1)
	return out, []string{"files kept by STORAGE_DRIVER=local are served from /files"}, nil
}

// repairLocalStorageGitignoreSource keeps the local root out of git.
func repairLocalStorageGitignoreSource(src string) (string, []string, []string) {
	if strings.Contains(src, "/storage/") {
		return src, nil, nil
	}
	out := strings.TrimRight(src, "\n") + "\n" + apiGitignoreLocalStorage
	return out, []string{"files kept by STORAGE_DRIVER=local are not committed"}, nil
}

// repairCSPImageSource lets a frontend load images from the API origin.
func repairCSPImageSource(src string) (string, []string, []string) {
	if strings.Contains(src, cspImgSrcNew) || strings.Count(src, cspImgSrcOld) != 1 {
		return src, nil, nil
	}
	return strings.Replace(src, cspImgSrcOld, cspImgSrcNew, 1), []string{"images kept by STORAGE_DRIVER=local load from the API origin"}, nil
}

// repairStorageDisk brings the files an existing project keeps up to C1. The
// storage package itself arrives whole, so nothing here runs unless it did.
func repairStorageDisk(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	if !fileContains(filepath.Join(apiRoot, "internal", "storage", "storage.go"), "func NewLocal(") {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}

	configPath := filepath.Join(apiRoot, "internal", "config", "config.go")
	if fileExists(configPath) {
		if err := repairSourceFile(root, m, configPath, repairStorageConfigSource); err != nil {
			return err
		}
	}
	// main.go reads the fields the config repair adds, so it waits for them.
	if fileContains(configPath, "LocalRoot string") && fileContains(configPath, "resolveStorageDriver()") {
		for _, server := range []string{filepath.Join(apiRoot, "cmd", "server", "main.go"), filepath.Join(root, "main.go")} {
			if !fileExists(server) {
				continue
			}
			if err := repairSourceFile(root, m, server, repairStorageInitSource); err != nil {
				return err
			}
		}
	}
	if routes := filepath.Join(apiRoot, "internal", "routes", "routes.go"); fileExists(routes) {
		if err := repairSourceFile(root, m, routes, repairLocalFilesRouteSource); err != nil {
			return err
		}
	}
	if ignore := filepath.Join(apiRoot, ".gitignore"); fileExists(ignore) {
		if err := repairTextFile(root, m, ignore, repairLocalStorageGitignoreSource); err != nil {
			return err
		}
	}
	for _, app := range []string{"admin", "web"} {
		for _, name := range []string{"next.config.ts", "next.config.mjs", "next.config.js", "vite.config.ts"} {
			path := filepath.Join(root, "apps", app, name)
			if !fileExists(path) {
				continue
			}
			if err := repairTextFile(root, m, path, repairCSPImageSource); err != nil {
				return fmt.Errorf("repairing %s: %w", path, err)
			}
		}
	}
	return nil
}
