package scaffold

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

/*
 * The MinIO image in docker-compose.yml has to be one a stranger can pull.
 *
 * It has broken twice. minio/minio went 404 when MinIO took it off Docker
 * Hub; quay.io/minio/minio then went 401, including with an anonymous pull
 * token, so `docker compose up` on a fresh machine failed with a registry
 * error that says nothing about which part of the project is wrong
 * (grit#109).
 *
 * Both times the only thing that would have caught it before a user did is
 * asking the registry. So this asks, and is skipped when there is no network
 * or when -short is set, because a unit test that needs the internet must
 * never be the reason a build fails offline.
 */

// TestComposeUsesAnImageTheWorldCanPull resolves the reference the scaffold
// writes, anonymously, exactly as Docker would.
func TestComposeUsesAnImageTheWorldCanPull(t *testing.T) {
	if testing.Short() {
		t.Skip("needs the network")
	}
	if os.Getenv("GRIT_OFFLINE") != "" {
		t.Skip("GRIT_OFFLINE is set")
	}

	ref := strings.TrimSpace(strings.TrimPrefix(minioImageNew, "image: "))
	repo, tag, ok := strings.Cut(ref, ":")
	if !ok {
		t.Fatalf("the MinIO image %q has no tag; an unpinned image in a compose "+
			"file people run for months is a moving target", ref)
	}
	if strings.Contains(repo, "/") && strings.Count(repo, "/") > 1 {
		t.Skipf("%s is not on Docker Hub; this test only knows Docker Hub's auth", repo)
	}

	client := &http.Client{Timeout: 20 * time.Second}

	// Docker Hub hands out an anonymous pull token for a public repository.
	tokenURL := "https://auth.docker.io/token?service=registry.docker.io&scope=repository:" +
		repo + ":pull"
	resp, err := client.Get(tokenURL)
	if err != nil {
		t.Skipf("no network: %v", err)
	}
	defer resp.Body.Close()
	var tok struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil || tok.Token == "" {
		t.Skipf("could not get an anonymous token for %s: %v", repo, err)
	}

	req, err := http.NewRequest(http.MethodHead,
		"https://registry-1.docker.io/v2/"+repo+"/manifests/"+tag, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tok.Token)
	req.Header.Set("Accept", strings.Join([]string{
		"application/vnd.oci.image.index.v1+json",
		"application/vnd.docker.distribution.manifest.list.v2+json",
		"application/vnd.docker.distribution.manifest.v2+json",
	}, ","))

	got, err := client.Do(req)
	if err != nil {
		t.Skipf("no network: %v", err)
	}
	defer got.Body.Close()

	if got.StatusCode != http.StatusOK {
		t.Errorf("docker-compose.yml points at %s, and an anonymous pull of it "+
			"answers HTTP %d. `docker compose up` fails for anybody without a "+
			"login to that registry, which is everybody running `grit new` for "+
			"the first time.", ref, got.StatusCode)
	}
}

// TestComposeMinioImageIsPinned is the half that needs no network.
func TestComposeMinioImageIsPinned(t *testing.T) {
	ref := strings.TrimSpace(strings.TrimPrefix(minioImageNew, "image: "))
	_, tag, ok := strings.Cut(ref, ":")
	if !ok || tag == "latest" || tag == "" {
		t.Errorf("the MinIO image is %q; it has to carry a real tag, because a "+
			"compose file is run for months and `latest` makes the stack "+
			"different on every machine", ref)
	}
}
