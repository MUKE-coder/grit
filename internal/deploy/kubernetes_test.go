package deploy

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// Every file this writes is parsed by a real YAML parser.
//
// That is the whole reason these are plain manifests and not a Helm chart: a
// chart's templates are not YAML, so a generator that produces broken output
// looks exactly like one that does not until somebody runs helm against a
// cluster. Asserting with somebody else's parser is the only check here that
// cannot agree with whatever this package happened to emit.
func TestEveryManifestIsValidYAML(t *testing.T) {
	files := KubernetesFiles(K8sConfig{
		AppName:    "contacts",
		Image:      "registry.example.com/contacts-api:v1.4.2",
		WebImage:   "registry.example.com/contacts-web:v1.4.2",
		AdminImage: "registry.example.com/contacts-admin:v1.4.2",
		Domain:     "contacts.example.com",
	})

	if len(files) < 10 {
		t.Fatalf("a triple project produced %d files, which is too few: %v", len(files), sortedKeys(files))
	}

	for _, name := range sortedKeys(files) {
		if !strings.HasSuffix(name, ".yaml") {
			continue
		}
		var doc map[string]any
		if err := yaml.Unmarshal([]byte(files[name]), &doc); err != nil {
			t.Errorf("%s is not valid YAML: %v", name, err)
			continue
		}
		// kustomization.yaml has no kind/apiVersion requirement worth
		// asserting beyond what it declares, and every other file is an
		// object the API server has to accept.
		if name == "kustomization.yaml" {
			continue
		}
		for _, key := range []string{"apiVersion", "kind", "metadata"} {
			if _, ok := doc[key]; !ok {
				t.Errorf("%s has no %s, so kubectl would refuse it", name, key)
			}
		}
	}
}

// An api-only project has no web or admin image, and must not be handed
// Deployments for containers that do not exist.
func TestAPIOnlyProjectGetsNoFrontendManifests(t *testing.T) {
	files := KubernetesFiles(K8sConfig{AppName: "svc", Image: "svc-api:v1"})

	for _, name := range []string{"web-deployment.yaml", "web-service.yaml", "admin-deployment.yaml", "admin-service.yaml"} {
		if _, ok := files[name]; ok {
			t.Errorf("%s was written for a project with no such app", name)
		}
	}
	for _, name := range []string{"api-deployment.yaml", "api-service.yaml", "migrate-job.yaml"} {
		if _, ok := files[name]; !ok {
			t.Errorf("%s is missing", name)
		}
	}
	// No domain, no Ingress: an Ingress with an empty host matches every
	// request that reaches the controller.
	if _, ok := files["ingress.yaml"]; ok {
		t.Error("an Ingress was written for a project with no domain")
	}
}

// The placeholder Secret must not be reachable by `kubectl apply -k`.
//
// It is the one file here that is dangerous when applied by accident: the API
// would start, mint tokens anybody can forge, and report itself healthy.
func TestKustomizationCannotShipThePlaceholderSecret(t *testing.T) {
	files := KubernetesFiles(K8sConfig{AppName: "contacts", Image: "contacts-api:v1", Domain: "x.example.com"})

	var k struct {
		Resources []string `yaml:"resources"`
	}
	if err := yaml.Unmarshal([]byte(files["kustomization.yaml"]), &k); err != nil {
		t.Fatalf("kustomization.yaml: %v", err)
	}
	for _, r := range k.Resources {
		if r == "secret.example.yaml" {
			t.Error("kustomization.yaml lists the placeholder secret, so `apply -k` would install a forgeable signing key")
		}
		if r == "migrate-job.yaml" {
			t.Error("kustomization.yaml lists the migration Job, which has to run and finish before the API rolls")
		}
		if _, ok := files[r]; !ok {
			t.Errorf("kustomization.yaml lists %q, which this generator does not write", r)
		}
	}
	if !strings.Contains(files["secret.example.yaml"], "CHANGE-ME") {
		t.Error("the secret template has no obvious placeholder, so a real-looking value could be applied unchanged")
	}
}

// Probes gate on the endpoints meant for an orchestrator.
//
// /api/health probes every dependency the application has, which is right for
// the System Health page and wrong here: a slow object store would fail
// readiness and pull every pod out of the load balancer over something that
// does not stop it serving requests.
func TestProbesUseTheOrchestratorEndpoints(t *testing.T) {
	files := KubernetesFiles(K8sConfig{AppName: "contacts", Image: "contacts-api:v1"})

	var dep struct {
		Spec struct {
			Template struct {
				Spec struct {
					Containers []struct {
						LivenessProbe  map[string]any `yaml:"livenessProbe"`
						ReadinessProbe map[string]any `yaml:"readinessProbe"`
						StartupProbe   map[string]any `yaml:"startupProbe"`
					} `yaml:"containers"`
				} `yaml:"spec"`
			} `yaml:"template"`
		} `yaml:"spec"`
	}
	if err := yaml.Unmarshal([]byte(files["api-deployment.yaml"]), &dep); err != nil {
		t.Fatalf("api-deployment.yaml: %v", err)
	}
	containers := dep.Spec.Template.Spec.Containers
	if len(containers) != 1 {
		t.Fatalf("the API deployment has %d containers, want 1", len(containers))
	}

	path := func(probe map[string]any) string {
		get, _ := probe["httpGet"].(map[string]any)
		p, _ := get["path"].(string)
		return p
	}
	if got := path(containers[0].LivenessProbe); got != "/healthz" {
		t.Errorf("liveness probes %q, want /healthz", got)
	}
	if got := path(containers[0].ReadinessProbe); got != "/readyz" {
		t.Errorf("readiness probes %q, want /readyz", got)
	}
	if containers[0].StartupProbe == nil {
		t.Error("no startup probe, so a cold start against a slow database restarts the pod forever")
	}
}

// A Job's pod template is immutable, so a fixed name makes the second deploy
// fail on an immutable field instead of running.
func TestMigrationJobNameFollowsTheImageTag(t *testing.T) {
	first := KubernetesFiles(K8sConfig{AppName: "contacts", Image: "contacts-api:v1.4.2"})
	second := KubernetesFiles(K8sConfig{AppName: "contacts", Image: "contacts-api:v1.5.0"})

	name := func(src string) string {
		var job struct {
			Metadata struct {
				Name string `yaml:"name"`
			} `yaml:"metadata"`
		}
		if err := yaml.Unmarshal([]byte(src), &job); err != nil {
			t.Fatalf("migrate-job.yaml: %v", err)
		}
		return job.Metadata.Name
	}

	a, b := name(first["migrate-job.yaml"]), name(second["migrate-job.yaml"])
	if a == b {
		t.Errorf("both releases produce the Job name %q, so the second apply is rejected", a)
	}
	if !strings.Contains(a, "v1-4-2") {
		t.Errorf("the Job name %q does not carry the image tag", a)
	}
}

// A tag may hold characters a Kubernetes object name may not, and a name the
// cluster refuses is a deploy that fails on its first step.
func TestImageTagIsSanitisedIntoAName(t *testing.T) {
	for image, want := range map[string]string{
		"app:v1.4.2":                      "v1-4-2",
		"registry.example.com/app:latest": "latest",
		"registry.example.com:5000/app":   "latest",
		"app:FEATURE_Branch":              "featurebranch",
		"app:":                            "latest",
		"app":                             "latest",
		"app:sha256.abcdef":               "sha256-abcdef",
	} {
		if got := imageTag(image); got != want {
			t.Errorf("imageTag(%q) = %q, want %q", image, got, want)
		}
	}
}

// The API container reads the ConfigMap and the Secret this set writes, by the
// names this set gives them. Two things that have to agree, and a typo in
// either reads as an empty environment at runtime rather than as an error.
func TestTheAPIReadsTheConfigAndSecretThisWrites(t *testing.T) {
	files := KubernetesFiles(K8sConfig{AppName: "contacts", Image: "contacts-api:v1"})

	nameOf := func(file string) string {
		var doc struct {
			Metadata struct {
				Name string `yaml:"name"`
			} `yaml:"metadata"`
		}
		if err := yaml.Unmarshal([]byte(files[file]), &doc); err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		return doc.Metadata.Name
	}

	config, secret := nameOf("configmap.yaml"), nameOf("secret.example.yaml")

	for _, file := range []string{"api-deployment.yaml", "migrate-job.yaml"} {
		src := files[file]
		if !strings.Contains(src, "name: "+config) {
			t.Errorf("%s does not read the ConfigMap %q this set writes", file, config)
		}
		if !strings.Contains(src, "name: "+secret) {
			t.Errorf("%s does not read the Secret %q this set writes", file, secret)
		}
	}
}
