package storage

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"library/apps/api/internal/config"
)

// The same checks run against every driver, so code written against Disk
// behaves the same on the local disk in development and on a bucket in
// production.

func TestLocalDisk(t *testing.T) {
	var disk *LocalDisk
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		disk.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)

	root := filepath.Join(t.TempDir(), "files")
	d, err := NewLocalDisk(LocalConfig{Root: root, PublicURL: server.URL + "/files", Secret: "a-test-secret-long-enough-to-sign-with"})
	if err != nil {
		t.Fatal(err)
	}
	disk = d

	runDiskSuite(t, d, func(key string) string {
		return d.signedURL(key, time.Now().Add(-time.Minute))
	})

	t.Run("a refused key writes nothing outside the root", func(t *testing.T) {
		_ = d.Put(context.Background(), "../escape.txt", strings.NewReader("x"), PutOptions{})
		if _, err := os.Stat(filepath.Join(filepath.Dir(root), "escape.txt")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("a file was written outside the root: %v", err)
		}
	})

	t.Run("a signature for one key does not open another", func(t *testing.T) {
		ctx := context.Background()
		for _, key := range []string{"private/a.txt", "private/b.txt"} {
			if err := d.Put(ctx, key, strings.NewReader(key), PutOptions{}); err != nil {
				t.Fatal(err)
			}
		}
		signed, err := d.TemporaryURL(ctx, "private/a.txt", time.Minute)
		if err != nil {
			t.Fatal(err)
		}
		forged := strings.Replace(signed, "private/a.txt", "private/b.txt", 1)
		if status, _ := fetch(t, forged); status != http.StatusForbidden {
			t.Fatalf("a signature for a.txt opened b.txt: status %d", status)
		}
	})
}

// TestS3Disk runs the same suite against a bucket, when one is named:
// STORAGE_TEST_S3_ENDPOINT, STORAGE_TEST_S3_ACCESS_KEY, STORAGE_TEST_S3_SECRET_KEY
// and STORAGE_TEST_S3_BUCKET. A local MinIO works.
func TestS3Disk(t *testing.T) {
	endpoint := os.Getenv("STORAGE_TEST_S3_ENDPOINT")
	if endpoint == "" {
		t.Skip("set STORAGE_TEST_S3_ENDPOINT, STORAGE_TEST_S3_ACCESS_KEY, STORAGE_TEST_S3_SECRET_KEY and STORAGE_TEST_S3_BUCKET to run against a bucket")
	}
	d, err := NewS3Disk(config.StorageConfig{
		Endpoint:  endpoint,
		AccessKey: os.Getenv("STORAGE_TEST_S3_ACCESS_KEY"),
		SecretKey: os.Getenv("STORAGE_TEST_S3_SECRET_KEY"),
		Bucket:    os.Getenv("STORAGE_TEST_S3_BUCKET"),
		Region:    "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	runDiskSuite(t, d, func(key string) string {
		signed, err := d.TemporaryURL(context.Background(), key, time.Second)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(2 * time.Second)
		return signed
	})
}

// The methods handlers have always called still work, over any driver.
func TestStorageKeepsItsMethodNames(t *testing.T) {
	ctx := context.Background()
	d, err := NewLocalDisk(LocalConfig{Root: t.TempDir(), PublicURL: "http://localhost:8080/files", Secret: "secret"})
	if err != nil {
		t.Fatal(err)
	}
	s := Wrap(d)

	if err := s.Upload(ctx, "uploads/a.txt", strings.NewReader("hello"), "text/plain"); err != nil {
		t.Fatal(err)
	}
	r, err := s.Download(ctx, "uploads/a.txt")
	if err != nil {
		t.Fatal(err)
	}
	body, err := io.ReadAll(r)
	r.Close()
	if err != nil || string(body) != "hello" {
		t.Fatalf("Download = %q, %v", body, err)
	}
	if size, ctype, err := s.Stat(ctx, "uploads/a.txt"); err != nil || size != 5 || ctype != "text/plain" {
		t.Fatalf("Stat = %d, %q, %v", size, ctype, err)
	}
	if got := s.GetURL("uploads/a.txt"); got != "http://localhost:8080/files/uploads/a.txt" {
		t.Fatalf("GetURL = %q", got)
	}
	if _, err := s.GetSignedURL(ctx, "uploads/a.txt", time.Minute); err != nil {
		t.Fatal(err)
	}
	if _, err := s.PresignPutURL(ctx, "uploads/b.txt", "text/plain", 5); !errors.Is(err, ErrPresignUnsupported) {
		t.Fatalf("PresignPutURL on the local disk = %v, want ErrPresignUnsupported", err)
	}
	if s.FileServer() == nil {
		t.Fatal("the local disk has no file server")
	}
	if err := s.DeleteMany(ctx, []string{"uploads/a.txt", "uploads/missing.txt"}); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(ctx, "uploads/a.txt"); err != nil {
		t.Fatal(err)
	}
}

func runDiskSuite(t *testing.T, disk Disk, expiredURL func(key string) string) {
	ctx := context.Background()
	run := strconv.FormatInt(time.Now().UnixNano(), 36)
	public := "uploads/disk-suite-" + run + "/"
	private := "private/disk-suite-" + run + "/"
	t.Cleanup(func() {
		for _, prefix := range []string{public, private} {
			objects, err := disk.List(ctx, prefix)
			if err != nil {
				t.Errorf("cleaning up %s: %v", prefix, err)
				continue
			}
			keys := make([]string, 0, len(objects))
			for _, obj := range objects {
				keys = append(keys, obj.Key)
			}
			if len(keys) > 0 {
				if err := disk.Delete(ctx, keys...); err != nil {
					t.Errorf("cleaning up %s: %v", prefix, err)
				}
			}
		}
	})

	put := func(t *testing.T, key, body string) {
		t.Helper()
		if err := disk.Put(ctx, key, strings.NewReader(body), PutOptions{ContentType: "text/plain"}); err != nil {
			t.Fatalf("Put(%q): %v", key, err)
		}
	}
	read := func(t *testing.T, key string) string {
		t.Helper()
		r, err := disk.Get(ctx, key)
		if err != nil {
			t.Fatalf("Get(%q): %v", key, err)
		}
		defer r.Close()
		body, err := io.ReadAll(r)
		if err != nil {
			t.Fatalf("reading %q: %v", key, err)
		}
		return string(body)
	}
	present := func(t *testing.T, key string) bool {
		t.Helper()
		ok, err := disk.Exists(ctx, key)
		if err != nil {
			t.Fatalf("Exists(%q): %v", key, err)
		}
		return ok
	}

	t.Run("put, get, exists and stat", func(t *testing.T) {
		key := public + "hello.txt"
		put(t, key, "hello")
		if got := read(t, key); got != "hello" {
			t.Fatalf("Get = %q, want hello", got)
		}
		if !present(t, key) || present(t, public+"missing.txt") {
			t.Fatal("Exists is wrong")
		}
		obj, err := disk.Stat(ctx, key)
		if err != nil {
			t.Fatal(err)
		}
		if obj.Size != 5 || obj.ContentType != "text/plain" || obj.LastModified.IsZero() {
			t.Fatalf("Stat = %+v", obj)
		}
		if _, err := disk.Stat(ctx, public+"missing.txt"); !errors.Is(err, ErrNotFound) {
			t.Fatalf("Stat of a missing file = %v, want ErrNotFound", err)
		}
		if _, err := disk.Get(ctx, public+"missing.txt"); !errors.Is(err, ErrNotFound) {
			t.Fatalf("Get of a missing file = %v, want ErrNotFound", err)
		}
		put(t, key, "replaced")
		if got := read(t, key); got != "replaced" {
			t.Fatalf("a second Put left %q", got)
		}
	})

	t.Run("copy and move", func(t *testing.T) {
		src := public + "copy/source.txt"
		put(t, src, "copied")
		if err := disk.Copy(ctx, src, public+"copy/target.txt"); err != nil {
			t.Fatal(err)
		}
		if read(t, public+"copy/target.txt") != "copied" || !present(t, src) {
			t.Fatal("Copy did not leave two files")
		}
		if err := disk.Move(ctx, public+"copy/target.txt", private+"moved/file.txt"); err != nil {
			t.Fatal(err)
		}
		if read(t, private+"moved/file.txt") != "copied" || present(t, public+"copy/target.txt") {
			t.Fatal("Move did not rename the file")
		}
		if err := disk.Move(ctx, public+"copy/missing.txt", public+"copy/x.txt"); !errors.Is(err, ErrNotFound) {
			t.Fatalf("Move of a missing file = %v, want ErrNotFound", err)
		}
	})

	t.Run("list", func(t *testing.T) {
		dir := public + "list/"
		for _, name := range []string{"a.txt", "b.txt", "nested/c.txt"} {
			put(t, dir+name, name)
		}
		put(t, public+"list-sibling.txt", "not in the directory")
		objects, err := disk.List(ctx, dir)
		if err != nil {
			t.Fatal(err)
		}
		var keys []string
		for _, obj := range objects {
			keys = append(keys, obj.Key)
			if obj.Size == 0 {
				t.Errorf("%s has no size", obj.Key)
			}
		}
		sort.Strings(keys)
		want := []string{dir + "a.txt", dir + "b.txt", dir + "nested/c.txt"}
		if strings.Join(keys, ",") != strings.Join(want, ",") {
			t.Fatalf("List = %v, want %v", keys, want)
		}
		if missing, err := disk.List(ctx, public+"nothing-here/"); err != nil || len(missing) != 0 {
			t.Fatalf("List of an empty prefix = %v, %v", missing, err)
		}
	})

	t.Run("delete many", func(t *testing.T) {
		keys := []string{public + "del/1.txt", public + "del/2.txt", public + "del/3.txt"}
		for _, key := range keys {
			put(t, key, key)
		}
		if err := disk.Delete(ctx, append(keys, public+"del/never-there.txt")...); err != nil {
			t.Fatal(err)
		}
		for _, key := range keys {
			if present(t, key) {
				t.Fatalf("%s survived Delete", key)
			}
		}
	})

	t.Run("URL serves a public file", func(t *testing.T) {
		key := public + "page one.txt"
		put(t, key, "public body")
		link := disk.URL(key)
		if !strings.Contains(link, "page%20one.txt") {
			t.Fatalf("URL = %q, want the space escaped", link)
		}
		if status, body := fetch(t, link); status != http.StatusOK || body != "public body" {
			t.Fatalf("GET %s = %d %q", link, status, body)
		}
	})

	t.Run("ServeFile streams a range", func(t *testing.T) {
		key := private + "range.txt"
		put(t, key, "0123456789")
		rec := serve(t, disk, key, Attachment, http.MethodGet, map[string]string{"Range": "bytes=3-6"})
		if rec.Code != http.StatusPartialContent || rec.Body.String() != "3456" || rec.Header().Get("Content-Range") != "bytes 3-6/10" {
			t.Fatalf("Range 3-6 = %d %q %v", rec.Code, rec.Body.String(), rec.Header())
		}
		if full := serve(t, disk, key, Attachment, http.MethodGet, nil); full.Code != http.StatusOK || full.Body.String() != "0123456789" || full.Header().Get("Content-Length") != "10" {
			t.Fatalf("GET = %d %q %v", full.Code, full.Body.String(), full.Header())
		}
	})

	t.Run("temporary URL opens a private file until it expires", func(t *testing.T) {
		key := private + "secret.txt"
		put(t, key, "private body")
		if status, _ := fetch(t, disk.URL(key)); status == http.StatusOK {
			t.Fatal("a private file is served without a signature")
		}
		signed, err := disk.TemporaryURL(ctx, key, time.Minute)
		if err != nil {
			t.Fatal(err)
		}
		if status, body := fetch(t, signed); status != http.StatusOK || body != "private body" {
			t.Fatalf("GET signed URL = %d %q", status, body)
		}
		if status, _ := fetch(t, expiredURL(key)); status != http.StatusForbidden {
			t.Fatalf("GET expired URL = %d, want 403", status)
		}
	})

	t.Run("keys that climb out of the store are refused", func(t *testing.T) {
		for _, key := range []string{"../escape.txt", "uploads/../../escape.txt", "/etc/passwd", ""} {
			if err := disk.Put(ctx, key, strings.NewReader("x"), PutOptions{}); !errors.Is(err, ErrInvalidKey) {
				t.Errorf("Put(%q) = %v, want ErrInvalidKey", key, err)
			}
			if _, err := disk.Get(ctx, key); !errors.Is(err, ErrInvalidKey) {
				t.Errorf("Get(%q) = %v, want ErrInvalidKey", key, err)
			}
			if err := disk.Delete(ctx, key); !errors.Is(err, ErrInvalidKey) {
				t.Errorf("Delete(%q) = %v, want ErrInvalidKey", key, err)
			}
		}
	})
}

func fetch(t *testing.T, link string) (int, string) {
	t.Helper()
	res, err := http.Get(link)
	if err != nil {
		t.Fatalf("GET %s: %v", link, err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("reading %s: %v", link, err)
	}
	return res.StatusCode, string(body)
}
