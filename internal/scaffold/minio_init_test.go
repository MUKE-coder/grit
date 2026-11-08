package scaffold

import (
	"strings"
	"testing"
)

// grit#185: the MinIO image runs as a non-root user but ships /data owned by
// root. Docker copies that ownership onto a fresh named volume, so MinIO cannot
// write to its own backend and exits:
//
//	FATAL Unable to initialize backend: file access denied
//
// An init service hands the directory over first. These tests hold the two
// halves of that together, because each one passes on its own while the pair is
// broken: a chown to a user MinIO does not run as leaves the volume exactly as
// unusable, and nothing would say so until somebody started a stack.

// TestMinioInitClaimsTheUserMinioRunsAs reads both sides. The uid the init
// service chowns to has to be the uid the minio service runs as.
func TestMinioInitClaimsTheUserMinioRunsAs(t *testing.T) {
	for _, c := range []struct {
		name string
		src  string
	}{
		{"docker-compose.yml", dockerCompose(Options{ProjectName: "demo", Architecture: ArchTriple})},
		{"docker-compose.prod.yml", dockerComposeProd(Options{ProjectName: "demo", Architecture: ArchTriple})},
	} {
		initStart, initEnd, ok := serviceBlock(c.src, "minio-init")
		if !ok {
			t.Fatalf("%s: no minio-init service, so a fresh volume stays root-owned and MinIO cannot start", c.name)
		}
		minioStart, minioEnd, ok := serviceBlock(c.src, "minio")
		if !ok {
			t.Fatalf("%s: no minio service", c.name)
		}
		initBlock, minioBlock := c.src[initStart:initEnd], c.src[minioStart:minioEnd]

		// The owner the init service hands /data to.
		chownTo := "\"chown\", \"-R\", \"" + minioRunAsUID + ":" + minioRunAsGID + "\""
		if !strings.Contains(initBlock, chownTo) {
			t.Errorf("%s: minio-init does not chown /data to %s:%s, which is the user the minio service runs as:\n%s",
				c.name, minioRunAsUID, minioRunAsGID, initBlock)
		}
		// And the test that keeps it O(1) on a volume that is already right.
		if !strings.Contains(initBlock, "\"-maxdepth\", \"0\", \"!\", \"-user\", \""+minioRunAsUID+"\"") {
			t.Errorf("%s: minio-init chowns unconditionally, which walks every object in the bucket on every `docker compose up`:\n%s", c.name, initBlock)
		}
		// The uid MinIO runs as, written out rather than left to the image.
		if !strings.Contains(minioBlock, "user: \""+minioRunAsUID+":"+minioRunAsGID+"\"") {
			t.Errorf("%s: the minio service does not pin user %s:%s, so an image that changes its uid would take the chown out from under itself:\n%s",
				c.name, minioRunAsUID, minioRunAsGID, minioBlock)
		}
		// And nothing starts MinIO before the chown has succeeded.
		if !strings.Contains(minioBlock, "condition: service_completed_successfully") {
			t.Errorf("%s: minio does not wait for minio-init, so the chown races its own startup:\n%s", c.name, minioBlock)
		}
		// The init service must be root to be able to chown at all.
		if !strings.Contains(initBlock, "user: \"0:0\"") {
			t.Errorf("%s: minio-init is not root, so it cannot chown a root-owned directory:\n%s", c.name, initBlock)
		}
	}
}

// TestMinioInitCommandNeedsNoShell guards the reason it is written with find.
//
// Both compose files are fmt.Sprintf format strings, so a % in this text
// renders as %!u(MISSING); and the upgrade repair splices the same text in
// directly, so escaping it to %% would be wrong there. A $ would be read by
// Compose as a variable. The obvious `stat -c %u` has both problems.
func TestMinioInitCommandNeedsNoShell(t *testing.T) {
	if strings.ContainsAny(minioInitService, "%$") {
		t.Errorf("minio-init's command carries a %% or a $: the first breaks the compose templates, which are format strings, and the second is read by Compose as a variable:\n%s", minioInitService)
	}
	for _, c := range []struct {
		name string
		src  string
	}{
		{"docker-compose.yml", dockerCompose(Options{ProjectName: "demo", Architecture: ArchTriple})},
		{"docker-compose.prod.yml", dockerComposeProd(Options{ProjectName: "demo", Architecture: ArchTriple})},
	} {
		if strings.Contains(c.src, "%!") {
			t.Errorf("%s: a format verb did not render", c.name)
		}
	}
}

// TestRepairedComposeIsTheGeneratedCompose upgrades a file shaped like the one
// the previous release wrote and asserts it comes out as the current template,
// byte for byte. Without this the repair can insert the right lines in the
// wrong place and only a reader would ever notice.
func TestRepairedComposeIsTheGeneratedCompose(t *testing.T) {
	opts := Options{ProjectName: "demo", Architecture: ArchTriple}
	for _, c := range []struct {
		name string
		want string
	}{
		{"docker-compose.yml", dockerCompose(opts)},
		{"docker-compose.prod.yml", dockerComposeProd(opts)},
	} {
		// The same file as it stood before this change: no init service, and no
		// user or depends_on on minio.
		old := strings.Replace(c.want, minioInitService, "", 1)
		old = strings.Replace(old, minioRunAsUser+minioInitDependsOn, "", 1)
		if old == c.want {
			t.Fatalf("%s: could not build the pre-fix file, so this test is asserting nothing", c.name)
		}

		// The same pass drops container_name on purpose, so that it is Compose's
		// to choose and a second copy of the stack can run beside the first.
		// Everything else has to match.
		want := composeContainerName.ReplaceAllString(c.want, "")

		got, changes, warnings := repairComposeHardeningSource(old)
		if len(warnings) != 0 {
			t.Errorf("%s: unexpected warnings: %v", c.name, warnings)
		}
		if got != want {
			t.Errorf("%s: a repaired file is not the generated file.\n--- repaired ---\n%s\n--- generated ---\n%s",
				c.name, minioRegion(got), minioRegion(want))
		}
		var said bool
		for _, ch := range changes {
			if strings.Contains(ch, "claims its data volume") {
				said = true
			}
		}
		if !said {
			t.Errorf("%s: the repair changed the file without saying what it did: %v", c.name, changes)
		}

		// And again: a second upgrade must find nothing left to do.
		twice, changes, _ := repairComposeHardeningSource(got)
		if twice != got {
			t.Errorf("%s: a second upgrade changed the file again: %v", c.name, changes)
		}
	}
}

// TestRepairLeavesAnEditedMinioAlone: a project that already gave MinIO a user
// of its own has made a decision, and a warning is the whole of Grit's business
// there.
func TestRepairLeavesAnEditedMinioAlone(t *testing.T) {
	opts := Options{ProjectName: "demo", Architecture: ArchTriple}
	edited := strings.Replace(dockerCompose(opts), minioInitService, "", 1)
	edited = strings.Replace(edited, minioRunAsUser+minioInitDependsOn, "    user: \"1000:1000\"\n", 1)

	got, _, _ := repairComposeHardeningSource(edited)
	if strings.Contains(got, "minio-init:") {
		t.Error("the repair added an init service to a minio that already runs as a user somebody chose")
	}
	if !strings.Contains(got, "user: \"1000:1000\"") {
		t.Error("the repair overwrote a user somebody chose")
	}
}

// TestRepairWarnsWhenItCannotPlaceTheInit: no `restart:` line to anchor on
// means the service has been rewritten, and the person has to be told what is
// wrong rather than left with a stack that will not start.
func TestRepairWarnsWhenItCannotPlaceTheInit(t *testing.T) {
	opts := Options{ProjectName: "demo", Architecture: ArchTriple}
	src := strings.Replace(dockerCompose(opts), minioInitService, "", 1)
	src = strings.Replace(src, minioRunAsUser+minioInitDependsOn, "", 1)
	start, end, ok := serviceBlock(src, "minio")
	if !ok {
		t.Fatal("no minio service")
	}
	src = src[:start] + strings.Replace(src[start:end], "    restart: unless-stopped\n", "", 1) + src[end:]

	got, _, warnings := repairComposeHardeningSource(src)
	if strings.Contains(got, "minio-init:") {
		t.Error("the repair guessed at where to put the init service in a file it does not recognise")
	}
	var warned bool
	for _, w := range warnings {
		if strings.Contains(w, "file access denied") {
			warned = true
		}
	}
	if !warned {
		t.Errorf("the repair said nothing about a minio that cannot start on a fresh volume: %v", warnings)
	}
}

// minioRegion trims a compose file to the part these tests are about, so a
// failure is readable.
func minioRegion(src string) string {
	i := strings.Index(src, "  minio")
	if i == -1 {
		return src
	}
	j := strings.Index(src[i:], "\nvolumes:")
	if j == -1 {
		return src[i:]
	}
	return src[i : i+j]
}
