package files

import "strings"

// AcceptsToMIMEs translates the high-level accept aliases used by the
// CLI (image:file:image, attachment:file:[pdf,doc]) to concrete MIME
// types + filename extensions. The upload handler uses this to
// validate per-field uploads at request time.
//
// "all" is a sentinel meaning the field accepts anything — the upload
// handler skips MIME checking and only enforces the size cap.
func AcceptsToMIMEs(accepts []string) (mimes []string, exts []string, acceptAll bool) {
	mimeSet := map[string]bool{}
	extSet := map[string]bool{}
	for _, a := range accepts {
		switch strings.ToLower(strings.TrimSpace(a)) {
		case "all":
			acceptAll = true
		case "image":
			for _, m := range []string{"image/jpeg", "image/png", "image/gif", "image/webp", "image/avif", "image/svg+xml"} {
				mimeSet[m] = true
			}
			for _, e := range []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".avif", ".svg"} {
				extSet[e] = true
			}
		case "video":
			for _, m := range []string{"video/mp4", "video/webm", "video/quicktime", "video/x-msvideo", "video/x-matroska"} {
				mimeSet[m] = true
			}
			for _, e := range []string{".mp4", ".webm", ".mov", ".avi", ".mkv"} {
				extSet[e] = true
			}
		case "audio":
			for _, m := range []string{"audio/mpeg", "audio/wav", "audio/ogg", "audio/x-m4a", "audio/webm"} {
				mimeSet[m] = true
			}
			for _, e := range []string{".mp3", ".wav", ".ogg", ".m4a"} {
				extSet[e] = true
			}
		case "pdf":
			mimeSet["application/pdf"] = true
			extSet[".pdf"] = true
		case "doc":
			for _, m := range []string{"application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document"} {
				mimeSet[m] = true
			}
			for _, e := range []string{".doc", ".docx"} {
				extSet[e] = true
			}
		case "excel":
			for _, m := range []string{"application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"} {
				mimeSet[m] = true
			}
			for _, e := range []string{".xls", ".xlsx"} {
				extSet[e] = true
			}
		case "csv":
			for _, m := range []string{"text/csv", "application/vnd.ms-excel"} {
				mimeSet[m] = true
			}
			extSet[".csv"] = true
		case "zip":
			for _, m := range []string{"application/zip", "application/x-zip-compressed"} {
				mimeSet[m] = true
			}
			extSet[".zip"] = true
		case "archive":
			for _, m := range []string{"application/zip", "application/x-zip-compressed", "application/x-tar", "application/gzip", "application/x-rar-compressed", "application/x-7z-compressed"} {
				mimeSet[m] = true
			}
			for _, e := range []string{".zip", ".tar", ".gz", ".tgz", ".rar", ".7z"} {
				extSet[e] = true
			}
		}
	}
	for m := range mimeSet {
		mimes = append(mimes, m)
	}
	for e := range extSet {
		exts = append(exts, e)
	}
	return mimes, exts, acceptAll
}

// AllowsMIME returns true if the given mime type is acceptable for the
// given accept-aliases. "all" short-circuits to true.
func AllowsMIME(accepts []string, mime string) bool {
	mimes, _, all := AcceptsToMIMEs(accepts)
	if all {
		return true
	}
	mime = strings.ToLower(strings.TrimSpace(mime))
	for _, m := range mimes {
		if strings.EqualFold(m, mime) {
			return true
		}
	}
	return false
}

// DefaultMaxSizeBytes returns the sensible default cap for the given
// accept set. Video-heavy fields get a much larger cap (300MB) because
// even a short 1080p clip is dozens of megabytes; everything else
// defaults to 5MB.
func DefaultMaxSizeBytes(accepts []string) int64 {
	for _, a := range accepts {
		if strings.EqualFold(a, "video") {
			return 300 << 20 // 300 MB
		}
	}
	return 5 << 20 // 5 MB
}
