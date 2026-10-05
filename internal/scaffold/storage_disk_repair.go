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
		"// Outside production, minio that is not answering becomes local, so a new\n" +
		"// project stores uploads before Docker is running instead of answering each\n" +
		"// one with STORAGE_UNAVAILABLE. Production never falls back: a server that\n" +
		"// lost its bucket should say so, not start writing to a disk that is not\n" +
		"// backed up and that the next container will not have.\n" +
		"//\n" +
		"// This used to test for a missing MINIO_ACCESS_KEY, and grit writes that key\n" +
		"// into .env, so the condition could never be true and the fallback never\n" +
		"// fired. Every new project on a machine without Docker answered every upload\n" +
		"// with \"File storage is not configured\" until somebody found STORAGE_DRIVER\n" +
		"// for themselves. Whether the key is set was never the question; whether\n" +
		"// MinIO answers is.\n" +
		"func resolveStorageDriver() string {\n" +
		"	driver := getEnv(\"STORAGE_DRIVER\", \"minio\")\n" +
		"	if driver != \"minio\" || getEnv(\"APP_ENV\", \"production\") == \"production\" {\n" +
		"		return driver\n" +
		"	}\n" +
		"	if getEnv(\"MINIO_ACCESS_KEY\", \"\") == \"\" {\n" +
		"		log.Println(\"MinIO has no credentials (MINIO_ACCESS_KEY), so files are kept on the local disk. Set STORAGE_DRIVER=local to make that the choice, or set the MinIO credentials to use MinIO\")\n" +
		"		return \"local\"\n" +
		"	}\n" +
		"	endpoint := getEnv(\"MINIO_ENDPOINT\", \"http://localhost:9000\")\n" +
		"	if minioAnswers(endpoint) {\n" +
		"		return driver\n" +
		"	}\n" +
		"	log.Printf(\"MinIO at %s is not answering, so files are being kept on the local disk at %s and uploads work. Start it with docker compose up -d minio, or set STORAGE_DRIVER=local to choose the disk deliberately.\", endpoint, getEnv(\"STORAGE_LOCAL_ROOT\", \"storage/app\"))\n" +
		"	return \"local\"\n" +
		"}\n" +
		"\n" +
		"// minioAnswers reports whether something is listening on the endpoint.\n" +
		"//\n" +
		"// A TCP dial, not a bucket call: this runs while the config is being read,\n" +
		"// before there is a client to ask, and the question is only whether the\n" +
		"// container is up. 300ms because the answer is immediate either way on a\n" +
		"// local machine, and a slow boot would be worse than one wrong guess.\n" +
		"func minioAnswers(endpoint string) bool {\n" +
		"	host := endpoint\n" +
		"	if u, err := url.Parse(endpoint); err == nil && u.Host != \"\" {\n" +
		"		host = u.Host\n" +
		"	}\n" +
		"	if !strings.Contains(host, \":\") {\n" +
		"		host += \":80\"\n" +
		"	}\n" +
		"	conn, err := net.DialTimeout(\"tcp\", host, 300*time.Millisecond)\n" +
		"	if err != nil {\n" +
		"		return false\n" +
		"	}\n" +
		"	_ = conn.Close()\n" +
		"	return true\n" +
		"}\n" +
		"\n"

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
	// mainStorageInitC1 is what v3.276.0 wrote. It connected a bucket only when
	// an endpoint was set, so STORAGE_DRIVER=s3 at AWS's default endpoint never
	// connected.
	mainStorageInitC1 = "\t// File storage. STORAGE_DRIVER=local keeps files in a directory on this\n" +
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

	// mainStorageInit opens the default store and every named disk (C3), with
	// storage.Open, which connects AWS S3 without an endpoint.
	mainStorageInit = "\t// File storage. STORAGE_DRIVER=local keeps files in a directory on this\n" +
		"\t// server and serves them from /files; minio, s3, r2 and b2 use a bucket.\n" +
		"\t// STORAGE_DISKS adds named disks, such as a private bucket for backups,\n" +
		"\t// which code reaches through storage.Disks.Get(name).\n" +
		"\tstorage.SetPublicPrefixes(cfg.StoragePublicPrefixes)\n" +
		"\tvar storageService *storage.Storage\n" +
		"\tif s, err := storage.Open(cfg.StorageDriver, cfg.Storage, storage.LocalConfig{\n" +
		"\t\tRoot:      cfg.Storage.LocalRoot,\n" +
		"\t\tPublicURL: cfg.Storage.PublicURL,\n" +
		"\t\tSecret:    cfg.Storage.URLSecret,\n" +
		"\t}); err != nil {\n" +
		"\t\tlog.Printf(\"Warning: Storage unavailable: %v (uploads disabled)\", err)\n" +
		"\t} else {\n" +
		"\t\tstorageService = s\n" +
		"\t\tstorage.Disks.SetDefault(s)\n" +
		"\t\tlog.Printf(\"File storage: %s\", s.Describe())\n" +
		"\t}\n" +
		"\tfor _, named := range cfg.StorageDisks {\n" +
		"\t\ts, err := storage.Open(named.Driver, named.Storage, storage.LocalConfig{\n" +
		"\t\t\tRoot:      named.Storage.LocalRoot,\n" +
		"\t\t\tPublicURL: named.Storage.PublicURL,\n" +
		"\t\t\tSecret:    named.Storage.URLSecret,\n" +
		"\t\t})\n" +
		"\t\tif err != nil {\n" +
		"\t\t\tlog.Printf(\"Warning: storage disk %q unavailable: %v\", named.Name, err)\n" +
		"\t\t\tcontinue\n" +
		"\t\t}\n" +
		"\t\tstorage.Disks.Add(named.Name, s)\n" +
		"\t\tlog.Printf(\"Storage disk %q: %s\", named.Name, s.Describe())\n" +
		"\t}\n"

	// Config gains the named disks and the public prefixes (C3, C4).
	configStorageDisksFieldAnchor = "\tStorage       StorageConfig // Resolved config for the active driver\n"
	configStorageDisksField       = configStorageDisksFieldAnchor + "\n" +
		"\t// StorageDisks are the named disks in STORAGE_DISKS, such as a private bucket\n" +
		"\t// for backups. StoragePublicPrefixes are the key prefixes anyone may read\n" +
		"\t// without a signed link (STORAGE_PUBLIC_PREFIXES).\n" +
		"\tStorageDisks          []StorageDisk\n" +
		"\tStoragePublicPrefixes []string\n"
	configStorageDisksLoadAnchor = "\t\tStorage:       resolveStorage(storageDriver),\n"
	configStorageDisksLoad       = configStorageDisksLoadAnchor + "\n" +
		"\t\tStorageDisks:          resolveStorageDisks(storageDriver),\n" +
		"\t\tStoragePublicPrefixes: resolveStoragePublicPrefixes(),\n"
	configStorageDisksFuncs = "// StorageDisk is one named disk from STORAGE_DISKS.\n" +
		"type StorageDisk struct {\n" +
		"\tName    string\n" +
		"\tDriver  string\n" +
		"\tStorage StorageConfig\n" +
		"}\n\n" +
		"// resolveStorageDisks reads the named disks in STORAGE_DISKS, a comma\n" +
		"// separated list of names.\n" +
		"//\n" +
		"// Each name takes STORAGE_DISK_<NAME>_* settings (DRIVER, BUCKET, ENDPOINT,\n" +
		"// ACCESS_KEY, SECRET_KEY, REGION, PUBLIC_URL, ROOT), and whatever it leaves\n" +
		"// unset comes from its driver's own settings. A private bucket for backups on\n" +
		"// the same provider as the uploads needs one line more than its name:\n" +
		"//\n" +
		"//\tSTORAGE_DISKS=backups\n" +
		"//\tSTORAGE_DISK_BACKUPS_BUCKET=myapp-backups\n" +
		"func resolveStorageDisks(defaultDriver string) []StorageDisk {\n" +
		"\tvar disks []StorageDisk\n" +
		"\tfor _, name := range strings.Split(os.Getenv(\"STORAGE_DISKS\"), \",\") {\n" +
		"\t\tname = strings.ToLower(strings.TrimSpace(name))\n" +
		"\t\tif name == \"\" || name == \"default\" {\n" +
		"\t\t\tcontinue\n" +
		"\t\t}\n" +
		"\t\tif strings.IndexFunc(name, func(r rune) bool {\n" +
		"\t\t\treturn (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '_' && r != '-'\n" +
		"\t\t}) >= 0 {\n" +
		"\t\t\tlog.Printf(\"STORAGE_DISKS: %q is not a disk name (letters, digits, _ and -), so it was skipped\", name)\n" +
		"\t\t\tcontinue\n" +
		"\t\t}\n" +
		"\t\tprefix := \"STORAGE_DISK_\" + strings.ToUpper(strings.ReplaceAll(name, \"-\", \"_\")) + \"_\"\n" +
		"\t\tdriver := getEnv(prefix+\"DRIVER\", defaultDriver)\n" +
		"\t\tdisk := resolveStorage(driver)\n" +
		"\t\tif driver == \"local\" {\n" +
		"\t\t\t// Apart from the default disk's files, and served under its route.\n" +
		"\t\t\tdisk.LocalRoot = filepath.Join(filepath.Dir(disk.LocalRoot), name)\n" +
		"\t\t\tdisk.PublicURL = strings.TrimRight(disk.PublicURL, \"/\") + \"/_disks/\" + name\n" +
		"\t\t}\n" +
		"\t\tdisk.Endpoint = getEnv(prefix+\"ENDPOINT\", disk.Endpoint)\n" +
		"\t\tdisk.AccessKey = getEnv(prefix+\"ACCESS_KEY\", disk.AccessKey)\n" +
		"\t\tdisk.SecretKey = getEnv(prefix+\"SECRET_KEY\", disk.SecretKey)\n" +
		"\t\tdisk.Bucket = getEnv(prefix+\"BUCKET\", disk.Bucket)\n" +
		"\t\tdisk.Region = getEnv(prefix+\"REGION\", disk.Region)\n" +
		"\t\tdisk.PublicURL = getEnv(prefix+\"PUBLIC_URL\", disk.PublicURL)\n" +
		"\t\tdisk.LocalRoot = getEnv(prefix+\"ROOT\", disk.LocalRoot)\n" +
		"\t\tif driver == defaultDriver && driver != \"local\" && disk.Bucket == resolveStorage(defaultDriver).Bucket {\n" +
		"\t\t\tlog.Printf(\"STORAGE_DISKS: %q uses the same bucket as the default disk, set %sBUCKET to keep its files apart\", name, prefix)\n" +
		"\t\t}\n" +
		"\t\tdisks = append(disks, StorageDisk{Name: name, Driver: driver, Storage: disk})\n" +
		"\t}\n" +
		"\treturn disks\n" +
		"}\n\n" +
		"// resolveStoragePublicPrefixes reads STORAGE_PUBLIC_PREFIXES: the key\n" +
		"// prefixes anyone may read without a signed link, comma separated. Every\n" +
		"// other key, backups included, is private.\n" +
		"func resolveStoragePublicPrefixes() []string {\n" +
		"\tvar prefixes []string\n" +
		"\tfor _, prefix := range strings.Split(getEnv(\"STORAGE_PUBLIC_PREFIXES\", \"uploads/,thumbnails/\"), \",\") {\n" +
		"\t\tif prefix = strings.Trim(strings.TrimSpace(prefix), \"/\"); prefix != \"\" {\n" +
		"\t\t\tprefixes = append(prefixes, prefix+\"/\")\n" +
		"\t\t}\n" +
		"\t}\n" +
		"\treturn prefixes\n" +
		"}\n\n"

	// backupCLIConfigMarker is in a config.go that has the named disks, which
	// cmd/backup/main.go reads.
	backupCLIConfigMarker = "func resolveStorageDisks("

	// routes.go serves a stored file through the API (C2).
	routesUploadDownloadAnchor = "\t\tprotected.GET(\"/uploads/:id\", uploadHandler.GetByID)\n"
	routesUploadDownload       = routesUploadDownloadAnchor +
		"\t\t// The file itself, through the API: private files too, on any driver.\n" +
		"\t\tprotected.GET(\"/uploads/:id/download\", uploadHandler.Download)\n"

	// .env.example documents the storage settings. The head as it was before
	// v3.276.0, as v3.276.0 wrote it, and as it is now.
	envStorageHeadPreC1 = "# Storage — Which provider to use: minio, s3, r2, b2\nSTORAGE_DRIVER=minio\n"
	envStorageHeadC1    = "# Storage: which provider to use: local, minio, s3, r2, b2\n" +
		"STORAGE_DRIVER=minio\n" +
		"# local keeps files in STORAGE_LOCAL_ROOT and the API serves them from\n" +
		"# APP_URL/files. Outside production, minio with no MINIO_ACCESS_KEY uses local,\n" +
		"# so uploads work before Docker is running. Production refuses local unless\n" +
		"# ALLOW_LOCAL_STORAGE_IN_PRODUCTION=true (one server, the directory on a volume).\n" +
		"# STORAGE_LOCAL_ROOT=storage/app\n" +
		"# Signs local temporary URLs. Unset, JWT_SECRET is used.\n" +
		"# STORAGE_URL_SECRET=\n"
	envStorageHeadNew = "# Storage: which provider to use: local, minio, s3, r2, b2\n" +
		"STORAGE_DRIVER=minio\n" +
		"# local keeps files in STORAGE_LOCAL_ROOT and the API serves them from\n" +
		"# APP_URL/files. No bucket, no Docker and nothing else to run, so it is the\n" +
		"# one to pick if you have no object storage: the files sit in a folder you\n" +
		"# can open, and uploads, thumbnails and signed links all work.\n" +
		"#\n" +
		"# You may not need to pick it. Outside production, minio that is not\n" +
		"# answering falls back to local on its own and says so in the log, so a new\n" +
		"# project takes uploads before Docker is running. Production never falls\n" +
		"# back, and refuses local unless ALLOW_LOCAL_STORAGE_IN_PRODUCTION=true (one\n" +
		"# server, with the directory on a volume that is backed up).\n" +
		"# STORAGE_LOCAL_ROOT=storage/app\n" +
		"# Signs local temporary URLs. Unset, JWT_SECRET is used.\n" +
		"# STORAGE_URL_SECRET=\n" +
		"# Key prefixes anyone may read without a signed link, comma separated. Every\n" +
		"# other key, backups included, is private: S3 and MinIO enforce it with the\n" +
		"# bucket policy, the local driver with its file route. R2 and B2 have no bucket\n" +
		"# policies, so there a bucket is public or not whatever the key: keep private\n" +
		"# files in a private bucket and serve them with a temporary URL.\n" +
		"# STORAGE_PUBLIC_PREFIXES=uploads/,thumbnails/\n" +
		"# Named disks, such as a private bucket for backups, which use a disk named\n" +
		"# backups when there is one. Each takes STORAGE_DISK_<NAME>_* settings (DRIVER,\n" +
		"# BUCKET, ENDPOINT, ACCESS_KEY, SECRET_KEY, REGION, PUBLIC_URL, ROOT), and\n" +
		"# anything unset comes from that driver's settings below.\n" +
		"# STORAGE_DISKS=backups\n" +
		"# STORAGE_DISK_BACKUPS_BUCKET=\n"

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

// repairStorageFallbackSource brings an existing resolveStorageDriver to the
// version that falls back when MinIO is not answering.
//
// The old one fell back when MINIO_ACCESS_KEY was empty, and grit writes that
// key into .env, so it could never fire: a project on a machine without Docker
// answered every upload with "File storage is not configured". The project that
// needs this most is an existing one, so it cannot be a new-projects-only fix.
func repairStorageFallbackSource(src string) (string, []string, []string) {
	if strings.Contains(src, "func minioAnswers(") {
		return src, nil, nil
	}
	start := strings.Index(src, storageDriverFuncHead)
	if start < 0 {
		return src, nil, []string{"config.go's resolveStorageDriver is not the one Grit wrote: " +
			"fall back to the local disk when MinIO is not answering, not when MINIO_ACCESS_KEY " +
			"is empty, or a project without Docker refuses every upload"}
	}
	// To the end of the function: the first line that is a closing brace in
	// column zero, which in generated Go is always the function's own.
	end := strings.Index(src[start:], "\n}\n")
	if end < 0 {
		return src, nil, []string{"config.go's resolveStorageDriver does not end where expected"}
	}
	end += start + len("\n}\n")

	out := src[:start] + configResolveStorageDriverFunc + src[end:]
	out, ok := withImports(out, "log", "net", "net/url", "strings", "time")
	if !ok {
		return src, nil, []string{"could not add imports to config.go"}
	}
	return out, []string{"uploads fall back to the local disk when MinIO is not answering, " +
		"so a project without Docker takes uploads instead of refusing them"}, nil
}

const (
	appURLFunc = "// resolveAppURL reads APP_URL and says so when it names a different port from\n" +
		"// the one the server is about to listen on.\n" +
		"//\n" +
		"// APP_URL is where the outside world reaches this app, so behind a proxy it is\n" +
		"// https://shop.example.com while the server listens on 8080, and that is\n" +
		"// correct. On localhost it is not: a local APP_URL whose port disagrees with\n" +
		"// APP_PORT is always a mistake, and a silent one. Every URL built from it is\n" +
		"// wrong, which with STORAGE_DRIVER=local means every uploaded file's URL points\n" +
		"// at a port where nothing is listening, and the symptom is an image that does\n" +
		"// not load rather than anything that mentions APP_URL.\n" +
		"func resolveAppURL(port string) string {\n" +
		"	appURL := getEnv(\"APP_URL\", \"http://localhost:8080\")\n" +
		"	u, err := url.Parse(appURL)\n" +
		"	if err != nil || u.Host == \"\" {\n" +
		"		return appURL\n" +
		"	}\n" +
		"	host, urlPort, err := net.SplitHostPort(u.Host)\n" +
		"	if err != nil || urlPort == port {\n" +
		"		return appURL\n" +
		"	}\n" +
		"	if host != \"localhost\" && host != \"127.0.0.1\" && host != \"[::1]\" {\n" +
		"		// A real hostname: a proxy in front, which is the normal case.\n" +
		"		return appURL\n" +
		"	}\n" +
		"	log.Printf(\"WARNING: APP_URL is %s but this server listens on port %s. Every URL the app builds, uploaded files included, will point at port %s where nothing is answering. Set APP_URL=http://%s:%s\",\n" +
		"		appURL, port, urlPort, host, port)\n" +
		"	return appURL\n" +
		"}\n" +
		"\n"
	appURLFuncAnchor = "// warnProviderMismatch says so when DATABASE_URL names a different engine from"
	appURLLoadOld    = "AppURL:             getEnv(\"APP_URL\", \"http://localhost:8080\"),"
	appURLLoadNew    = "AppURL:             resolveAppURL(firstNonEmpty(os.Getenv(\"PORT\"), os.Getenv(\"APP_PORT\"), \"8080\")),"
)

// storageDriverFuncHead is where the old resolveStorageDriver starts, comment
// and all, so the replacement takes its documentation with it.
const storageDriverFuncHead = "// resolveStorageDriver picks the storage driver from STORAGE_DRIVER."

// repairAppURLPortSource routes APP_URL through resolveAppURL, which says when
// it names a different port from the one the server listens on.
//
// Silent until it is not: every URL the app builds comes from APP_URL, so with
// STORAGE_DRIVER=local every uploaded file's URL points at a port where nothing
// is answering, and the symptom is an image that does not load rather than
// anything mentioning APP_URL.
func repairAppURLPortSource(src string) (string, []string, []string) {
	if strings.Contains(src, "func resolveAppURL(") {
		return src, nil, nil
	}
	if strings.Count(src, appURLLoadOld) != 1 || strings.Count(src, appURLFuncAnchor) != 1 {
		return src, nil, []string{"config.go is not the file Grit wrote: read APP_URL through a " +
			"resolveAppURL that compares its port with the one the server listens on, or a " +
			"mismatch serves every uploaded file from an address nothing answers"}
	}
	out := strings.Replace(src, appURLLoadOld, appURLLoadNew, 1)
	out = strings.Replace(out, appURLFuncAnchor, appURLFunc+appURLFuncAnchor, 1)
	out, ok := withImports(out, "log", "net", "net/url")
	if !ok {
		return src, nil, []string{"could not add imports to config.go"}
	}
	return out, []string{"an APP_URL whose port disagrees with the one the server listens on is reported at boot"}, nil
}

// repairStorageConfigSource adds the local driver to config.go.
func repairStorageConfigSource(src string) (string, []string, []string) {
	if strings.Contains(src, "func resolveStorageDriver()") {
		// Already has one, which is every project since the local driver was
		// added. It may still be the version that could never fall back.
		return repairStorageFallbackSource(src)
	}
	if !strings.Contains(src, "func resolveStorage(driver string) StorageConfig") {
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

// repairStorageInitSource builds the stores in main.go with storage.Open: the
// one STORAGE_DRIVER names, local and AWS S3 at its default endpoint included,
// and each named disk.
func repairStorageInitSource(src string) (string, []string, []string) {
	if strings.Contains(src, "storage.Disks.SetDefault(") || !strings.Contains(src, "storage.New(cfg.Storage)") {
		return src, nil, nil
	}
	for _, old := range []string{mainStorageInitC1, mainStorageInitOld, singleStorageInitOld} {
		if strings.Count(src, old) == 1 {
			out := strings.Replace(src, old, mainStorageInit, 1)
			return out, []string{"storage.Open builds the store STORAGE_DRIVER names (local, and AWS S3 without an endpoint) and each disk in STORAGE_DISKS"}, nil
		}
	}
	return src, nil, []string{"main.go is not the file Grit wrote: build the store with storage.Open(cfg.StorageDriver, cfg.Storage, storage.LocalConfig{...}), record it with storage.Disks.SetDefault, and open each of cfg.StorageDisks, or STORAGE_DRIVER=s3 without an endpoint and named disks are never opened"}
}

// repairStorageDisksConfigSource adds the named disks and the public prefixes to
// a config.go that already has the local driver.
func repairStorageDisksConfigSource(src string) (string, []string, []string) {
	if strings.Contains(src, backupCLIConfigMarker) || !strings.Contains(src, "func resolveStorageDriver()") {
		return src, nil, nil
	}
	for _, anchor := range []string{configStorageDisksFieldAnchor, configStorageDisksLoadAnchor, configResolveStorageAnchor} {
		if strings.Count(src, anchor) != 1 {
			return src, nil, []string{"config.go is not the file Grit wrote: add StorageDisks (read from STORAGE_DISKS) and StoragePublicPrefixes (read from STORAGE_PUBLIC_PREFIXES) to Config, as a new project's config.go does"}
		}
	}
	out := strings.Replace(src, configStorageDisksFieldAnchor, configStorageDisksField, 1)
	out = strings.Replace(out, configStorageDisksLoadAnchor, configStorageDisksLoad, 1)
	out = strings.Replace(out, configResolveStorageAnchor, configStorageDisksFuncs+configResolveStorageAnchor, 1)
	// Beside "os" in the standard library group, where a new config.go has it.
	if !strings.Contains(out, "\t\"path/filepath\"\n") && strings.Count(out, "\t\"os\"\n") == 1 {
		out = strings.Replace(out, "\t\"os\"\n", "\t\"os\"\n\t\"path/filepath\"\n", 1)
	}
	out, ok := withImports(out, "log", "os", "path/filepath", "strings")
	if !ok {
		return src, nil, []string{"could not add imports to config.go"}
	}
	return out, []string{"reads STORAGE_DISKS for named disks and STORAGE_PUBLIC_PREFIXES for the keys anyone may read"}, nil
}

// repairUploadDownloadRouteSource mounts GET /uploads/:id/download.
func repairUploadDownloadRouteSource(src string) (string, []string, []string) {
	if strings.Contains(src, "uploadHandler.Download") || !strings.Contains(src, "uploadHandler.GetByID") {
		return src, nil, nil
	}
	if strings.Count(src, routesUploadDownloadAnchor) != 1 {
		return src, nil, []string{"routes.go is not the file Grit wrote: mount GET /uploads/:id/download on uploadHandler.Download to serve stored files through the API"}
	}
	return strings.Replace(src, routesUploadDownloadAnchor, routesUploadDownload, 1),
		[]string{"GET /uploads/:id/download serves a stored file through the API, with Range support"}, nil
}

// repairEnvStorageSource documents the storage settings in .env.example.
func repairEnvStorageSource(src string) (string, []string, []string) {
	if strings.Contains(src, "STORAGE_DISKS") || !strings.Contains(src, "STORAGE_DRIVER=") {
		return src, nil, nil
	}
	for _, old := range []string{envStorageHeadC1, envStorageHeadPreC1} {
		if strings.Count(src, old) == 1 {
			return strings.Replace(src, old, envStorageHeadNew, 1),
				[]string{"documents STORAGE_DRIVER=local, STORAGE_PUBLIC_PREFIXES and STORAGE_DISKS"}, nil
		}
	}
	return src, nil, []string{"the storage block is not the one Grit wrote: see a new project's .env.example for STORAGE_DRIVER=local, STORAGE_PUBLIC_PREFIXES and STORAGE_DISKS"}
}

// writeBackupCLI brings cmd/backup/main.go to storage.Open, once config.go has
// the fields it reads and the backup service has the disk it asks for.
func writeBackupCLI(apiRoot string, opts Options) error {
	cli := filepath.Join(apiRoot, "cmd", "backup", "main.go")
	if !fileExists(cli) || fileContains(cli, "storage.Open(") ||
		!fileContains(filepath.Join(apiRoot, "internal", "config", "config.go"), backupCLIConfigMarker) ||
		!fileContains(filepath.Join(apiRoot, "internal", "backup", "backup.go"), "func (s *Service) Store()") {
		return nil
	}
	content := strings.ReplaceAll(strings.ReplaceAll(backupCmdMainGo(), "~", "`"), "{{MODULE}}", opts.Module())
	if err := writeFile(cli, content); err != nil {
		return err
	}
	if fileContains(cli, "storage.Open(") {
		fmt.Println("  ✓ cmd/backup/main.go: grit backup uses the configured storage, and the backups disk when there is one")
	} else {
		fmt.Println("  ⚠ cmd/backup/main.go has been edited, so it still calls storage.New: open the store with storage.Open as a new project's does, or grit backup ignores STORAGE_DRIVER=local and the backups disk")
	}
	return nil
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
		// Named disks and public prefixes (C3, C4), on top of the local driver.
		if err := repairSourceFile(root, m, configPath, repairStorageDisksConfigSource); err != nil {
			return err
		}
		// An APP_URL that disagrees with the listening port serves every
		// uploaded file from an address nothing answers.
		if err := repairSourceFile(root, m, configPath, repairAppURLPortSource); err != nil {
			return err
		}
	}
	// main.go and cmd/backup read the fields the config repairs add, so they
	// wait for them, and for the storage package that has storage.Open.
	storageOpen := fileContains(filepath.Join(apiRoot, "internal", "storage", "disks.go"), "func Open(")
	if storageOpen && fileContains(configPath, "LocalRoot string") && fileContains(configPath, backupCLIConfigMarker) {
		if err := writeBackupCLI(apiRoot, opts); err != nil {
			return err
		}
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
		if fileContains(filepath.Join(apiRoot, "internal", "handlers", "upload.go"), "func (h *UploadHandler) Download(") {
			if err := repairSourceFile(root, m, routes, repairUploadDownloadRouteSource); err != nil {
				return err
			}
		}
	}
	if env := filepath.Join(root, ".env.example"); fileExists(env) {
		if err := repairTextFile(root, m, env, repairEnvStorageSource); err != nil {
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
