package scaffold

import (
	"path/filepath"
	"strings"

	"github.com/MUKE-coder/grit/v3/internal/manifest"
)

// repairThumbnailWorker brings a project scaffolded before v3.250.0 up to the
// fix for H12 in the contact-app review: the thumbnail job decoded any upload
// without a pixel limit, in the API process, and retried the out-of-memory crash
// until its retries ran out.
//
// The pixel check is in internal/storage, which upgrade delivers. The job is in
// internal/jobs/workers.go, which it does not, so the one change there anchors on
// what Grit wrote, and waits until storage has the errors it names.
func repairThumbnailWorker(root string, opts Options) error {
	apiRoot := opts.APIRoot(root)
	workers := filepath.Join(apiRoot, "internal", "jobs", "workers.go")
	if !fileExists(workers) || !fileContains(filepath.Join(apiRoot, "internal", "storage", "image.go"), "ErrImageTooLarge") {
		return nil
	}
	m, err := manifest.Load(root)
	if err != nil {
		return err
	}
	return repairSourceFile(root, m, workers, repairThumbnailRetrySource)
}

const (
	thumbnailCallOld = "\t\tthumbBytes, err := storage.GenerateThumbnail(reader, payload.MimeType)\n\t\tif err != nil {\n\t\t\treturn fmt.Errorf(\"generating thumbnail: %w\", err)\n\t\t}\n"
	thumbnailCallNew = "\t\tthumbBytes, err := storage.GenerateThumbnail(reader, payload.MimeType)\n" +
		"\t\t// An image refused as too large, or one that cannot be decoded, fails the\n" +
		"\t\t// same way on every attempt, so it is not retried.\n" +
		"\t\tif errors.Is(err, storage.ErrImageTooLarge) || errors.Is(err, storage.ErrUnreadableImage) {\n" +
		"\t\t\treturn fmt.Errorf(\"generating thumbnail for %s: %v: %w\", payload.Key, err, asynq.SkipRetry)\n\t\t}\n" +
		"\t\tif err != nil {\n\t\t\treturn fmt.Errorf(\"generating thumbnail: %w\", err)\n\t\t}\n"
)

func repairThumbnailRetrySource(src string) (string, []string, []string) {
	if !strings.Contains(src, "func handleImageProcess(") || strings.Contains(src, "storage.ErrImageTooLarge") {
		return src, nil, nil
	}
	if strings.Count(src, thumbnailCallOld) != 1 {
		return src, nil, []string{"the thumbnail job is not the one Grit wrote: return asynq.SkipRetry for storage.ErrImageTooLarge and storage.ErrUnreadableImage, or a decompression bomb is retried"}
	}
	out := strings.Replace(src, thumbnailCallOld, thumbnailCallNew, 1)
	if !strings.Contains(out, "\t\"errors\"\n") {
		var ok bool
		if out, ok = addImportGroup(out, "errors"); !ok {
			return src, nil, []string{"could not add the errors import to jobs/workers.go"}
		}
	}
	return out, []string{"an image refused before decoding, or undecodable, is not retried"}, nil
}
