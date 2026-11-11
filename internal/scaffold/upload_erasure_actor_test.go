package scaffold

import (
	"strings"
	"testing"
)

// grit#201: the upload allowlist could only ever grow, and PPTX was missing.

// PowerPoint is in the baseline, beside Word and Excel.
//
// It was the only one of the three missing, which reads as an oversight rather
// than a decision. The reporter's product exports a deck as PPTX and could not
// accept one back, while it could accept a RAR archive.
func TestTheUploadBaselineTakesPowerPoint(t *testing.T) {
	src := uploadHandlerGo()

	for _, mime := range []string{
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"application/vnd.ms-powerpoint",
		"image/tiff",
	} {
		if !strings.Contains(src, mime) {
			t.Errorf("the upload baseline does not accept %s", mime)
		}
	}

	// The two that were already there, so this cannot pass by the list having
	// been replaced with something else.
	for _, mime := range []string{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	} {
		if !strings.Contains(src, mime) {
			t.Errorf("the upload baseline no longer accepts %s", mime)
		}
	}
}

// A project can say "exactly these types" and be believed.
//
// UPLOAD_ALLOWED_MIME only ever extended, so a project that accepts four
// document types could not refuse the thirty it did not ask for, archives
// included. The alternative was editing the line in routes.go that builds the
// list, which means owning 1600 lines of central routing for an allowlist.
func TestTheUploadAllowlistCanBeNarrowed(t *testing.T) {
	src := uploadHandlerGo()

	if !strings.Contains(src, `strings.CutPrefix(m, "only:")`) {
		t.Fatal("there is no way to replace the baseline, so a project can only ever " +
			"accept more types than Grit chose")
	}

	at := strings.Index(src, "func UploadMIMEAllowlist(")
	if at == -1 {
		t.Fatal("UploadMIMEAllowlist is gone")
	}
	body := src[at:]
	if end := strings.Index(body, "\n}\n"); end != -1 {
		body = body[:end]
	}

	// The baseline is skipped only when narrowing was asked for.
	if !strings.Contains(body, "if !exact || len(cleaned) == 0 {") {
		t.Error("the baseline is not conditional, so only: cannot replace it")
	}
	// An only: with nothing after it refuses every upload, which nobody meant
	// to type. It has to fall back rather than lock the project out of its own
	// uploads over a trailing comma.
	if !strings.Contains(body, "len(cleaned) == 0") {
		t.Error("an empty only: list would refuse every upload")
	}
}

// grit#202: erasure deleted the rows and left the files in the bucket.
//
// Scrub hard-deletes every registered row and anonymises the user, and had no
// awareness of storage at all, so a right-to-erasure request reported success
// while the documents were still there. The objects were unreferenced by then,
// so nothing could find them to clean up later either.
func TestErasureDeletesTheFilesTheRowsName(t *testing.T) {
	erasure := apiErasureGo()

	if !strings.Contains(erasure, "type Cleanup func(tx *gorm.DB, ownerColumn, userID string) error") {
		t.Fatal("there is no cleanup hook, so a model that owns files has nowhere to say so")
	}
	if !strings.Contains(erasure, "func RegisterWithCleanup(") {
		t.Error("nothing can register a cleanup")
	}

	// Before the delete, or there is nothing left to read the keys off.
	call := strings.Index(erasure, "t.cleanup(tx, t.column, userID)")
	del := strings.Index(erasure, `Delete(t.model).Error; err != nil`)
	if call == -1 || del == -1 {
		t.Fatal("could not find the cleanup call and the delete")
	}
	if call > del {
		t.Error("the cleanup runs after the rows are deleted, so the keys it needs are gone")
	}

	// And Grit's own uploads table uses it, so the common case needs no
	// project code. This is the whole point: the reporter could not close it
	// downstream without owning framework files.
	models := apiErasableModelsGo()
	if !strings.Contains(models, "erasure.RegisterWithCleanup(&Upload{}, \"user_id\", deleteUploadObjects)") {
		t.Error("uploads do not register a cleanup, so erasing a user still leaves their files")
	}
	if !strings.Contains(models, "disk.DeleteMany(ctx, batch)") {
		t.Error("the uploads cleanup does not delete anything from storage")
	}
	// Uploads must not also be in the plain Register loop, or the rows are
	// registered twice and deleted twice: the second pass has no cleanup and
	// counts the same table again.
	loop := models[strings.Index(models, "for _, m := range []interface{}{"):]
	if end := strings.Index(loop, "}"); end != -1 {
		loop = loop[:end]
	}
	if strings.Contains(loop, "&Upload{}") {
		t.Error("Upload is registered twice: once plainly and once with the cleanup")
	}
}

// grit#203: framework handlers attached no actor, so a guard below the handler
// worked for generated resources and silently stood aside for uploads.
//
// Standing aside for a context with no actor is correct, because that is a job
// or a migration acting as the server. The bug was that a request-scoped call
// was indistinguishable from one.
func TestTheAuthMiddlewareAttachesTheActor(t *testing.T) {
	auth := apiAuthMiddlewareGo()
	const attach = "c.Request = c.Request.WithContext(authz.WithActor(c.Request.Context(), authz.ActorOf(c)))"

	// Both of them: Auth, and Identify, which is the one public pages use.
	// Fixing only Auth would leave every signed-in reader of a public page
	// invisible to the same guard.
	if n := strings.Count(auth, attach); n != 2 {
		t.Errorf("the actor is attached %d times; it belongs in both Auth and Identify, "+
			"or a guard below the handler sees a caller on some routes and not others", n)
	}

	// After the role is resolved, because ActorOf reads it to decide Admin.
	role := strings.LastIndex(auth, `c.Set("user_roles", names)`)
	at := strings.Index(auth, attach)
	if role == -1 || at == -1 {
		t.Fatal("could not find the role resolution and the actor attachment")
	}
	if at < strings.Index(auth, `c.Set("user_roles", names)`) {
		t.Error("the actor is attached before the role is resolved, so every caller " +
			"looks like a non-admin to anything reading it")
	}
}
