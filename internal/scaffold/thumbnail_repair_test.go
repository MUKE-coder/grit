package scaffold

import (
	"go/format"
	"strings"
	"testing"
)

func TestRepairThumbnailRetry(t *testing.T) {
	src := "package jobs\n\nimport (\n\t\"fmt\"\n)\n\nfunc handleImageProcess() {\n\tfunc() error {\n" + thumbnailCallOld + "\t\treturn nil\n\t}()\n}\n"
	out, fixed, warn := repairThumbnailRetrySource(src)
	if len(warn) > 0 || len(fixed) != 1 {
		t.Fatalf("fixed %v, warned %v", fixed, warn)
	}
	for _, want := range []string{"asynq.SkipRetry", "storage.ErrImageTooLarge", "\"errors\""} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %s:\n%s", want, out)
		}
	}
	if _, err := format.Source([]byte(out)); err != nil {
		t.Fatalf("not valid Go: %v", err)
	}
	if again, fixed, _ := repairThumbnailRetrySource(out); again != out || len(fixed) > 0 {
		t.Error("a second upgrade changed workers.go again")
	}
}

func TestFreshTemplatesNeedNoThumbnailRepair(t *testing.T) {
	src := jobsWorkersGo()
	if out, fixed, warn := repairThumbnailRetrySource(src); out != src || len(fixed) > 0 || len(warn) > 0 {
		t.Errorf("the worker template still needs the repair (fixed %v, warn %v)", fixed, warn)
	}
	img := storageImageGo()
	for _, want := range []string{"image.DecodeConfig", "ErrImageTooLarge", "media.Get(\"\").MaxPixels"} {
		if !strings.Contains(img, want) {
			t.Errorf("storage/image.go is missing %s", want)
		}
	}
	if strings.Contains(img, "img, _, err := image.Decode(reader)") {
		t.Error("an image helper still decodes without checking the header")
	}
	if _, err := format.Source([]byte(storageImageTestGo())); err != nil {
		t.Errorf("image_test.go: %v", err)
	}
}
