package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// C2: the helpers exist, and the upload handler uses them instead of its own
// sniffing and a key built from the uploaded name.
func TestUploadHandlerUsesTheStorageHelpers(t *testing.T) {
	store := storageStoreGo()
	for _, want := range []string{
		"func Store(ctx context.Context, disk Disk, dir string, file *multipart.FileHeader, opts StoreOptions) (string, error)",
		"func StoreAs(ctx context.Context, disk Disk, dir string, file *multipart.FileHeader, name string, opts StoreOptions) (string, error)",
		"func ServeFile(c *gin.Context, disk Disk, key string, disposition Disposition)",
		"func DetectContentType(r io.ReadSeeker, declared string) (string, error)",
		"uuid.NewString()", "Content-Range", "StatusPartialContent", "Last-Modified",
	} {
		if !strings.Contains(store, want) {
			t.Errorf("store.go is missing %q", want)
		}
	}

	handler := uploadHandlerGo()
	for _, want := range []string{
		"storage.DetectContentType(file, header.Header.Get(\"Content-Type\"))",
		"storage.NewKey(\"uploads\", res.Primary.Ext)",
		"storage.Store(c.Request.Context(), disk, \"uploads\", header,",
		"storage.VisibilityPrivate",
		"func (h *UploadHandler) Download(c *gin.Context) {",
		"storage.ServeFileAs(c, h.Storage.Disk(), upload.Path, disposition, upload.OriginalName)",
	} {
		if !strings.Contains(handler, want) {
			t.Errorf("upload.go is missing %q", want)
		}
	}
	for _, gone := range []string{"stamp.UnixNano()", "http.DetectContentType(sniff"} {
		if strings.Contains(handler, gone) {
			t.Errorf("upload.go still has %q", gone)
		}
	}
	if !strings.Contains(apiRoutesGo(), routesUploadDownload) {
		t.Error("routes.go does not mount the download route")
	}
}

// C3 and C4: the registry, visibility, the backup disk and the env keys.
func TestNamedDisksAndVisibilityTemplates(t *testing.T) {
	disks := storageDisksGo()
	for _, want := range []string{"var Disks = &Registry{}", "func Open(driver string, cfg config.StorageConfig, local LocalConfig) (*Storage, error)", `case "s3":`} {
		if !strings.Contains(disks, want) {
			t.Errorf("disks.go is missing %q", want)
		}
	}
	if !strings.Contains(storageDiskGo(), "Visibility Visibility") || !strings.Contains(storageDiskGo(), "ErrVisibilityMismatch") {
		t.Error("PutOptions has no Visibility")
	}
	for name, src := range map[string]string{"storage.go": storageServiceGo(), "local.go": storageLocalGo()} {
		if !strings.Contains(src, "checkVisibility(key, opts.Visibility)") {
			t.Errorf("%s does not check visibility on Put", name)
		}
	}
	service := storageServiceGo()
	for _, want := range []string{"func SetPublicPrefixes(", "func (d *S3Disk) GetRange(", ".amazonaws.com/", "if cfg.AccessKey != \"\" {"} {
		if !strings.Contains(service, want) {
			t.Errorf("storage.go is missing %q", want)
		}
	}

	backup := backupServiceGo()
	if !strings.Contains(backup, "storage.Disks.Get(DiskName)") || !strings.Contains(backup, `const DiskName = "backups"`) {
		t.Error("backups do not use the backups disk")
	}
	cli := backupCmdMainGo()
	if strings.Contains(cli, "storage.New(cfg.Storage)") || !strings.Contains(cli, "storage.Open(cfg.StorageDriver, cfg.Storage,") || !strings.Contains(cli, "named.Name != backup.DiskName") {
		t.Error("grit backup does not open the configured storage")
	}
	handler := backupHandlerGo()
	if strings.Contains(handler, "h.Storage == nil") || !strings.Contains(handler, "h.svc().Locate(") {
		t.Error("the backup handler does not use the backups disk")
	}
	for name, src := range map[string]string{"backup.go": backupServiceGo(), "handlers/backup.go": handler, "cmd/backup/main.go": cli} {
		mustFormatGo(t, name, strings.ReplaceAll(strings.ReplaceAll(src, "~", "`"), "{{MODULE}}", "example.com/app"))
	}

	env := envExampleFile(Options{ProjectName: "app"})
	if !strings.Contains(env, envStorageHeadNew) {
		t.Error(".env.example does not document STORAGE_DISKS and STORAGE_PUBLIC_PREFIXES")
	}
	if strings.Contains(envStorageHeadNew, "%") {
		t.Error("the storage block would be read as a format string")
	}
}

// C5: the legacy helpers go through media.Transform.
func TestImageHelpersUseTheMediaPipeline(t *testing.T) {
	img := storageImageGo()
	for _, want := range []string{"media.Transform(", "const ThumbnailSize = 400", "media.Fill(ThumbnailSize, ThumbnailSize)"} {
		if !strings.Contains(img, want) {
			t.Errorf("image.go is missing %q", want)
		}
	}
	for _, gone := range []string{"imaging.Fill(", "imaging.Resize(", "image.Decode("} {
		if strings.Contains(img, gone) {
			t.Errorf("image.go still decodes on its own with %q", gone)
		}
	}
	mustFormatGo(t, "image.go", strings.ReplaceAll(img, "{{MODULE}}", "example.com/app"))
	test := storageImageTestGo()
	if !strings.Contains(test, "func TestThumbnailsFollowEXIFOrientation(") {
		t.Error("image_test.go has no EXIF orientation test")
	}
	mustFormatGo(t, "image_test.go", strings.ReplaceAll(test, "{{MODULE}}", "example.com/app"))
}

func TestEnvStorageRepair(t *testing.T) {
	fresh := envExampleFile(Options{ProjectName: "app"})
	for name, old := range map[string]string{
		"v3.276.0":        strings.Replace(fresh, envStorageHeadNew, envStorageHeadC1, 1),
		"before v3.276.0": strings.Replace(fresh, envStorageHeadNew, envStorageHeadPreC1, 1),
	} {
		t.Run(name, func(t *testing.T) {
			if old == fresh {
				t.Fatal("could not reconstruct the old file")
			}
			repaired, fixed, warnings := repairEnvStorageSource(old)
			if len(fixed) == 0 || len(warnings) > 0 || repaired != fresh {
				t.Fatalf("fixed %v, warnings %v, same as the template: %v", fixed, warnings, repaired == fresh)
			}
			if again, fixed, warnings := repairEnvStorageSource(repaired); again != repaired || len(fixed) > 0 || len(warnings) > 0 {
				t.Fatal("a second repair changed something")
			}
		})
	}
	edited := "STORAGE_DRIVER=s3\n"
	if out, fixed, warnings := repairEnvStorageSource(edited); out != edited || len(fixed) > 0 || len(warnings) == 0 {
		t.Error("an .env.example Grit did not write was changed, or not reported")
	}
}

func TestStorageDisksRepairsLeaveUnknownFilesAlone(t *testing.T) {
	config := strings.Replace(withoutStorageDisksConfig(apiConfigGo()), configStorageDisksLoadAnchor, "\t\tStorage: loadMyStorage(),\n", 1)
	if out, fixed, warnings := repairStorageDisksConfigSource(config); out != config || len(fixed) > 0 || len(warnings) == 0 {
		t.Error("a config.go Grit did not write was changed, or not reported")
	}
	routes := strings.Replace(apiRoutesGo(), routesUploadDownload, "\t\tapi.GET(\"/uploads/:id\", uploadHandler.GetByID)\n", 1)
	if out, fixed, warnings := repairUploadDownloadRouteSource(routes); out != routes || len(fixed) > 0 || len(warnings) == 0 {
		t.Error("a routes.go Grit did not write was changed, or not reported")
	}
}

// cmd/backup/main.go reads config fields an older config.go lacks, so upgrade
// holds it back until the config repair has added them, then writes it.
func TestBackupCLIWaitsForTheConfigRepair(t *testing.T) {
	root := t.TempDir()
	opts := Options{ProjectName: "app"}
	apiRoot := opts.APIRoot(root)
	configPath := filepath.Join(apiRoot, "internal", "config", "config.go")
	cli := filepath.Join(apiRoot, "cmd", "backup", "main.go")
	if err := os.MkdirAll(filepath.Dir(configPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(configPath, []byte(withoutStorageDisksConfig(apiConfigGo())), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeBackupFiles(root, opts); err != nil {
		t.Fatal(err)
	}
	if fileExists(cli) {
		t.Fatal("cmd/backup/main.go was written for a config.go without StorageDisks")
	}
	if !fileContains(filepath.Join(apiRoot, "internal", "backup", "backup.go"), "func (s *Service) Store()") {
		t.Fatal("the backup service was not written")
	}

	oldCLI := strings.ReplaceAll(backupCmdMainGo(), "storage.Open(", "storage.New(")
	if err := os.MkdirAll(filepath.Dir(cli), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cli, []byte(oldCLI), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeBackupCLI(apiRoot, opts); err != nil {
		t.Fatal(err)
	}
	if fileContains(cli, "storage.Open(") {
		t.Fatal("the CLI was rewritten before config.go had StorageDisks")
	}

	if err := os.WriteFile(configPath, []byte(apiConfigGo()), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := writeBackupCLI(apiRoot, opts); err != nil {
		t.Fatal(err)
	}
	if !fileContains(cli, "storage.Open(") {
		t.Fatal("the CLI was not rewritten once config.go had StorageDisks")
	}
}
