package scaffold

import (
	"os"
	"strings"
	"testing"
)

// POST /uploads/complete recorded a row for any object key the caller named,
// including another user's upload, and deleting that row deleted the object.
// It now records only a key presigned for the caller, and only once.
func TestCompleteUploadRecordsOnlyTheCallersPresignedKey(t *testing.T) {
	raw, err := os.ReadFile("api_storage_files.go")
	if err != nil {
		t.Fatal(err)
	}
	// A Windows checkout with core.autocrlf writes CRLF, and the search for the
	// end of the handler below looks for "\n}\n".
	src := strings.ReplaceAll(string(raw), "\r\n", "\n")

	if !strings.Contains(src, `key := fmt.Sprintf("uploads/%s/%s/%s", userID, time.Now().Format("2006/01"), filename)`) {
		t.Error("presigned keys are not under the caller's own prefix")
	}

	start := strings.Index(src, "func (h *UploadHandler) CompleteUpload(c *gin.Context) {")
	if start < 0 {
		t.Fatal("CompleteUpload is missing from the upload handler template")
	}
	end := strings.Index(src[start:], "\n}\n")
	body := src[start : start+end]

	prefix := strings.Index(body, `strings.HasPrefix(req.Key, "uploads/"+userID+"/")`)
	stat := strings.Index(body, "h.Storage.Stat(")
	switch {
	case prefix < 0:
		t.Error("CompleteUpload does not check that the key is the caller's")
	case stat < 0:
		t.Error("CompleteUpload no longer asks the bucket what it received")
	case prefix > stat:
		t.Error("the bucket is asked before the key is checked, so the answer reveals whether any key exists")
	}
	if !strings.Contains(body, "UPLOAD_ALREADY_RECORDED") {
		t.Error("a key can be recorded more than once")
	}
	if strings.Contains(body, "userID.(string)") {
		t.Error("CompleteUpload still reads the user id twice")
	}
}
