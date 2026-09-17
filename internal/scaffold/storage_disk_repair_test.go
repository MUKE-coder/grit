package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

func gofmtOrFatal(t *testing.T, name, src string) string {
	t.Helper()
	out, err := format.Source([]byte(src))
	if err != nil {
		t.Fatalf("%s is not valid Go: %v", name, err)
	}
	return string(out)
}

// The storage package ships the Disk interface, the local driver and one suite
// for both drivers, and keeps the names handlers already call.
func TestStorageTemplatesShipTheDiskInterfaceAndLocalDriver(t *testing.T) {
	disk := storageDiskGo()
	for _, want := range []string{
		"type Disk interface {",
		"Put(ctx context.Context, key string, r io.Reader, opts PutOptions) error",
		"Delete(ctx context.Context, keys ...string) error",
		"TemporaryURL(ctx context.Context, key string, ttl time.Duration) (string, error)",
		"ErrPresignUnsupported",
	} {
		if !strings.Contains(disk, want) {
			t.Errorf("disk.go is missing %q", want)
		}
	}
	local := storageLocalGo()
	for _, want := range []string{"os.CreateTemp(", "os.Rename(tmp.Name(), target)", "hmac.Equal(", "IsPublicKey(key)", "sandbox"} {
		if !strings.Contains(local, want) {
			t.Errorf("local.go is missing %q", want)
		}
	}
	service := storageServiceGo()
	for _, want := range []string{
		"func (s *Storage) Upload(ctx context.Context, key string, reader io.Reader, contentType string) error",
		"func (s *Storage) Delete(ctx context.Context, key string) error",
		"func (s *Storage) DeleteMany(ctx context.Context, keys []string) error",
		"func (s *Storage) Stat(ctx context.Context, key string) (int64, string, error)",
		"func (s *Storage) GetSignedURL(",
		"func (s *Storage) PresignPutURL(",
		"func New(cfg config.StorageConfig) (*Storage, error)",
		"func NewLocal(cfg LocalConfig) (*Storage, error)",
	} {
		if !strings.Contains(service, want) {
			t.Errorf("storage.go is missing %q", want)
		}
	}
	// storage.go arrives before config.go is repaired, so it must not read the
	// fields the repair adds.
	for _, field := range []string{"LocalRoot", "URLSecret"} {
		if strings.Contains(service, "cfg."+field) || strings.Contains(local, "cfg."+field) {
			t.Errorf("the storage package reads config.StorageConfig.%s, which an unrepaired config.go does not have", field)
		}
	}
	// Nor may the helpers and the registry, which arrive with it.
	for name, src := range map[string]string{"store.go": storageStoreGo(), "disks.go": storageDisksGo()} {
		for _, field := range []string{"LocalRoot", "URLSecret", "StorageDisks", "StoragePublicPrefixes"} {
			if strings.Contains(src, "."+field) {
				t.Errorf("%s reads %s, which an unrepaired config.go does not have", name, field)
			}
		}
	}
	for name, src := range map[string]string{
		"storage.go": service, "disk.go": disk, "local.go": local, "disk_test.go": storageDiskTestGo(),
		"store.go": storageStoreGo(), "store_test.go": storageStoreTestGo(), "disks.go": storageDisksGo(),
	} {
		gofmtOrFatal(t, name, strings.ReplaceAll(src, "{{MODULE}}", "example.com/app"))
	}
	if !strings.Contains(storageDiskTestGo(), "func TestS3Disk(t *testing.T)") || !strings.Contains(storageDiskTestGo(), "runDiskSuite(t, d,") {
		t.Error("the Disk suite does not run against both drivers")
	}
}

func TestStorageWiringTemplates(t *testing.T) {
	opts := Options{ProjectName: "app"}
	for name, src := range map[string]string{"main.go": apiMainGo(opts), "single main.go": singleMainGo(opts)} {
		if !strings.Contains(src, mainStorageInit) {
			t.Errorf("%s does not build the local driver", name)
		}
	}
	config := apiConfigGo()
	for _, want := range []string{configStorageFields, configStorageDriverNew, configLocalStorageCheck, configResolveStorageDriverFunc, configLocalCase,
		configStorageDisksField, configStorageDisksLoad, configStorageDisksFuncs} {
		if !strings.Contains(config, want) {
			t.Errorf("config.go is missing %q", want)
		}
	}
	if !strings.Contains(apiRoutesGo(), routesLocalFilesBlock) {
		t.Error("routes.go does not mount the local file route")
	}
	if !strings.Contains(apiGitignore(), "/storage/") {
		t.Error("the API .gitignore does not ignore the local root")
	}
	if strings.Count(nextSecurityHeaders(), cspImgSrcTight) != 1 || strings.Count(viteSecurityHeaders(), cspImgSrcTight) != 1 {
		t.Error("the frontend CSP does not allow images from the API origin")
	}
	handler := uploadHandlerGo()
	if !strings.Contains(handler, "errors.Is(err, storage.ErrPresignUnsupported)") || !strings.Contains(handler, `"method": "multipart"`) {
		t.Error("presign does not send the client to the multipart upload")
	}
}

// A config.go, main.go and routes.go as Grit wrote them before C1 are brought
// to the template, and a second pass changes nothing.
func TestStorageDiskRepairsReproduceTheTemplate(t *testing.T) {
	opts := Options{ProjectName: "app"}
	cases := []struct {
		name   string
		fresh  string
		old    func(string) string
		repair func(string) (string, []string, []string)
	}{
		{
			name:  "config.go from before v3.276.0",
			fresh: apiConfigGo(),
			old: func(s string) string {
				s = withoutStorageDisksConfig(s)
				s = strings.Replace(s, configStorageFields, configStorageFieldsAnchor, 1)
				s = strings.Replace(s, configStorageDriverFieldNew, configStorageDriverFieldOld, 1)
				s = strings.Replace(s, configStorageDriverNew, configStorageDriverOld, 1)
				s = strings.Replace(s, configLocalStorageCheck, "", 1)
				s = strings.Replace(s, configResolveStorageDriverFunc, "", 1)
				return strings.Replace(s, configLocalCase, "", 1)
			},
			repair: func(s string) (string, []string, []string) {
				out, fixed, warn := repairStorageConfigSource(s)
				more, fixedMore, warnMore := repairStorageDisksConfigSource(out)
				return more, append(fixed, fixedMore...), append(warn, warnMore...)
			},
		},
		{
			name:   "config.go from v3.276.0",
			fresh:  apiConfigGo(),
			old:    withoutStorageDisksConfig,
			repair: repairStorageDisksConfigSource,
		},
		{
			name:   "main.go from v3.276.0",
			fresh:  apiMainGo(opts),
			old:    func(s string) string { return strings.Replace(s, mainStorageInit, mainStorageInitC1, 1) },
			repair: repairStorageInitSource,
		},
		{
			name:   "main.go",
			fresh:  apiMainGo(opts),
			old:    func(s string) string { return strings.Replace(s, mainStorageInit, mainStorageInitOld, 1) },
			repair: repairStorageInitSource,
		},
		{
			name:   "routes.go without the download route",
			fresh:  apiRoutesGo(),
			old:    func(s string) string { return strings.Replace(s, routesUploadDownload, routesUploadDownloadAnchor, 1) },
			repair: repairUploadDownloadRouteSource,
		},
		{
			name:   "single main.go",
			fresh:  singleMainGo(opts),
			old:    func(s string) string { return strings.Replace(s, mainStorageInit, singleStorageInitOld, 1) },
			repair: repairStorageInitSource,
		},
		{
			name:   "routes.go",
			fresh:  apiRoutesGo(),
			old:    func(s string) string { return strings.Replace(s, routesLocalFilesBlock, "", 1) },
			repair: repairLocalFilesRouteSource,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			old := c.old(c.fresh)
			if old == c.fresh {
				t.Fatal("could not reconstruct the old file")
			}
			repaired, fixed, warnings := c.repair(old)
			if len(warnings) > 0 || len(fixed) == 0 {
				t.Fatalf("repair: fixed %v, warnings %v", fixed, warnings)
			}
			if gofmtOrFatal(t, "repaired", repaired) != gofmtOrFatal(t, "template", c.fresh) {
				t.Fatal("the repaired file is not the template")
			}
			again, fixed, warnings := c.repair(repaired)
			if again != repaired || len(fixed) > 0 || len(warnings) > 0 {
				t.Fatalf("a second repair changed something: fixed %v, warnings %v", fixed, warnings)
			}
		})
	}
}

// withoutStorageDisksConfig is config.go as v3.276.0 wrote it.
func withoutStorageDisksConfig(s string) string {
	s = strings.Replace(s, configStorageDisksField, configStorageDisksFieldAnchor, 1)
	s = strings.Replace(s, configStorageDisksLoad, configStorageDisksLoadAnchor, 1)
	s = strings.Replace(s, configStorageDisksFuncs, "", 1)
	return strings.Replace(s, "\t\"path/filepath\"\n", "", 1)
}

func TestStorageDiskRepairsLeaveUnknownFilesAlone(t *testing.T) {
	config := strings.Replace(apiConfigGo(), configStorageDriverNew, "\tstorageDriver := os.Getenv(\"DRIVER\")\n", 1)
	config = strings.Replace(config, configResolveStorageDriverFunc, "", 1)
	if out, fixed, warnings := repairStorageConfigSource(config); out != config || len(fixed) > 0 || len(warnings) == 0 {
		t.Error("a config.go Grit did not write was changed, or not reported")
	}
	main := strings.Replace(apiMainGo(Options{ProjectName: "app"}), mainStorageInit, "\tstorageService, _ := storage.New(cfg.Storage)\n", 1)
	if out, fixed, warnings := repairStorageInitSource(main); out != main || len(fixed) > 0 || len(warnings) == 0 {
		t.Error("a main.go Grit did not write was changed, or not reported")
	}

	ignore := "tmp/\n"
	repaired, _, _ := repairLocalStorageGitignoreSource(ignore)
	if !strings.HasSuffix(repaired, "/storage/\n") {
		t.Errorf(".gitignore = %q", repaired)
	}
	if again, fixed, _ := repairLocalStorageGitignoreSource(repaired); again != repaired || len(fixed) > 0 {
		t.Error("a second .gitignore repair changed something")
	}
	csp := "  " + cspImgSrcOld + "\n"
	repaired, _, _ = repairCSPImageSource(csp)
	if repaired != "  "+cspImgSrcNew+"\n" {
		t.Errorf("CSP = %q", repaired)
	}
	if again, fixed, _ := repairCSPImageSource(repaired); again != repaired || len(fixed) > 0 {
		t.Error("a second CSP repair changed something")
	}
}
