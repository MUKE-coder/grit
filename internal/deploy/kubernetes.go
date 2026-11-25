package deploy

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Kubernetes output: the manifests a Grit project needs to run on a cluster.
//
// ── Why plain manifests and a kustomization, not a Helm chart ────────────
//
// grit#259 asked for "a Helm chart (or Kustomize base)". This writes the
// second, for a reason worth stating: every file it produces is valid YAML, so
// the generator can be checked by parsing its own output with a real YAML
// parser rather than by matching strings. A Helm chart cannot be checked that
// way, because `{{- if .Values.x }}` is not YAML, and a template that produces
// broken manifests looks exactly like one that does not until somebody runs
// helm against a cluster.
//
// `kubectl apply -k deploy/k8s` works with no extra tool installed. A team
// that wants Helm can wrap these, and if the demand is there a chart is a
// separate, additive piece of work.
//
// ── What it does not write, and why ──────────────────────────────────────
//
// Postgres, Redis and object storage. On a cluster these are nearly always a
// managed service or somebody's own operator, and a generator that guessed
// would produce a single-replica database with no backups that looks
// production-ready. They are named through the Secret instead.
//
// A worker Deployment, which the issue asked for. A Grit API process starts
// its own asynq worker whenever Redis is reachable, and there is no
// worker-only binary to put in a separate Deployment. Generating one would
// produce a chart that does not deploy. The consequence is written into the
// manifests rather than hidden: every API replica is also a queue consumer.

// K8sConfig is what the manifests are built from.
type K8sConfig struct {
	// AppName names every object and the namespace's default.
	AppName string
	// Namespace to deploy into.
	Namespace string
	// Image is the API image, with its tag. The tag matters: it names the
	// migration Job, so a second deploy of the same tag does not collide
	// with the first one's completed Job.
	Image string
	// WebImage and AdminImage are the frontend images, empty when the
	// project has no such app.
	WebImage   string
	AdminImage string
	// Domain is the host an Ingress routes. Empty writes no Ingress.
	Domain string
	// OutDir is where the manifests are written.
	OutDir string
	// Replicas is the API's starting replica count.
	Replicas int
}

// KubernetesFiles renders the manifests as path -> contents, relative to
// OutDir.
//
// Separate from writing them so a test can read what would be written without
// a temporary directory, and so `--dry-run` can list the files.
func KubernetesFiles(cfg K8sConfig) map[string]string {
	cfg = withK8sDefaults(cfg)

	files := map[string]string{
		"README.md":           k8sReadme(cfg),
		"namespace.yaml":      k8sNamespace(cfg),
		"secret.example.yaml": k8sSecretExample(cfg),
		"configmap.yaml":      k8sConfigMap(cfg),
		"api-deployment.yaml": k8sAPIDeployment(cfg),
		"api-service.yaml":    k8sService(cfg, "api", 8080),
		"api-hpa.yaml":        k8sHPA(cfg),
		"migrate-job.yaml":    k8sMigrateJob(cfg),
	}

	if cfg.WebImage != "" {
		files["web-deployment.yaml"] = k8sFrontendDeployment(cfg, "web", cfg.WebImage, 3000)
		files["web-service.yaml"] = k8sService(cfg, "web", 3000)
	}
	if cfg.AdminImage != "" {
		files["admin-deployment.yaml"] = k8sFrontendDeployment(cfg, "admin", cfg.AdminImage, 3001)
		files["admin-service.yaml"] = k8sService(cfg, "admin", 3001)
	}
	if cfg.Domain != "" {
		files["ingress.yaml"] = k8sIngress(cfg)
	}

	files["kustomization.yaml"] = k8sKustomization(cfg, files)
	return files
}

// Kubernetes writes the manifests to cfg.OutDir.
func Kubernetes(cfg K8sConfig) error {
	cfg = withK8sDefaults(cfg)
	files := KubernetesFiles(cfg)

	if err := os.MkdirAll(cfg.OutDir, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", cfg.OutDir, err)
	}
	for _, name := range sortedKeys(files) {
		path := filepath.Join(cfg.OutDir, name)
		if err := os.WriteFile(path, []byte(files[name]), 0o644); err != nil { // #nosec G306 -- manifests are not secret
			return fmt.Errorf("writing %s: %w", path, err)
		}
		fmt.Printf("  ✓ %s\n", filepath.ToSlash(filepath.Join(filepath.Base(cfg.OutDir), name)))
	}

	// Printed relative to where the person is standing, because an absolute
	// Windows path inside a kubectl line is not something anybody pastes.
	dir := filepath.ToSlash(cfg.OutDir)
	if wd, err := os.Getwd(); err == nil {
		if rel, err := filepath.Rel(wd, cfg.OutDir); err == nil && !strings.HasPrefix(rel, "..") {
			dir = filepath.ToSlash(rel)
		}
	}

	fmt.Println()
	fmt.Println("  Next:")
	fmt.Println("    1. Copy secret.example.yaml, fill it in, and apply it. It is NOT in")
	fmt.Println("       kustomization.yaml, so `apply -k` cannot ship the placeholders.")
	fmt.Printf("    2. kubectl apply -f %s/migrate-job.yaml\n", dir)
	fmt.Printf("       kubectl wait --for=condition=complete --timeout=5m job/%s-migrate-%s -n %s\n",
		cfg.AppName, imageTag(cfg.Image), cfg.Namespace)
	fmt.Printf("    3. kubectl apply -k %s\n", dir)
	return nil
}

func withK8sDefaults(cfg K8sConfig) K8sConfig {
	if cfg.AppName == "" {
		cfg.AppName = "grit-app"
	}
	if cfg.Namespace == "" {
		cfg.Namespace = cfg.AppName
	}
	if cfg.Image == "" {
		cfg.Image = cfg.AppName + "-api:latest"
	}
	if cfg.Replicas < 1 {
		cfg.Replicas = 2
	}
	if cfg.OutDir == "" {
		cfg.OutDir = filepath.Join("deploy", "k8s")
	}
	return cfg
}

// imageTag is the part after the last colon, which names the migration Job.
//
// Sanitised to what a Kubernetes name allows: a tag may contain characters an
// object name may not, and a Job whose name a cluster refuses is a deploy that
// fails at the first step.
func imageTag(image string) string {
	tag := "latest"
	if i := strings.LastIndex(image, ":"); i >= 0 && !strings.Contains(image[i:], "/") {
		tag = image[i+1:]
	}
	var b strings.Builder
	for _, r := range strings.ToLower(tag) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-', r == '.':
			b.WriteRune('-')
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "latest"
	}
	if len(out) > 40 {
		out = strings.Trim(out[:40], "-")
	}
	return out
}

func sortedKeys(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func k8sNamespace(cfg K8sConfig) string {
	return `apiVersion: v1
kind: Namespace
metadata:
  name: ` + cfg.Namespace + `
  labels:
    app.kubernetes.io/name: ` + cfg.AppName + `
`
}

// k8sSecretExample is a template, and is deliberately left out of the
// kustomization.
//
// A generated Secret containing a placeholder signing key is the one file in
// this set that is dangerous if it is applied by accident: the application
// would start, mint tokens anybody can forge, and look healthy. So `apply -k`
// cannot reach it, and the only way it gets into a cluster is somebody opening
// it, filling it in and applying it on purpose.
func k8sSecretExample(cfg K8sConfig) string {
	return `# Copy this file, fill in every value, and apply it BEFORE the rest.
#
#   cp secret.example.yaml secret.yaml   # secret.yaml is gitignored
#   $EDITOR secret.yaml
#   kubectl apply -f secret.yaml
#
# It is not listed in kustomization.yaml on purpose. These are placeholders,
# and an application that starts with a placeholder signing key mints tokens
# anybody can forge while looking perfectly healthy.
#
# Generate the two keys with:
#   openssl rand -base64 48
apiVersion: v1
kind: Secret
metadata:
  name: ` + cfg.AppName + `-secrets
  namespace: ` + cfg.Namespace + `
type: Opaque
stringData:
  # Signs and verifies every token this API issues.
  JWT_SECRET: "CHANGE-ME-generate-with-openssl-rand-base64-48"
  # Encrypts columns declared :encrypted. Losing it loses the data in them.
  ENCRYPTION_KEY: "CHANGE-ME-generate-with-openssl-rand-base64-48"

  # Postgres. A managed instance or your own operator: these manifests do not
  # install a database, because a generated single-replica Postgres with no
  # backups looks production-ready and is not.
  POSTGRES_HOST: "postgres.example.internal"
  POSTGRES_PORT: "5432"
  POSTGRES_USER: "` + cfg.AppName + `"
  POSTGRES_PASSWORD: "CHANGE-ME"
  POSTGRES_DB: "` + cfg.AppName + `"

  # Redis. Required: the API starts its asynq worker only when Redis answers,
  # so without it no background job ever runs.
  REDIS_URL: "redis://:CHANGE-ME@redis.example.internal:6379"

  # Object storage for uploads. Any S3-compatible endpoint.
  STORAGE_DRIVER: "s3"
  S3_ENDPOINT: "https://s3.example.com"
  S3_BUCKET: "` + cfg.AppName + `"
  S3_ACCESS_KEY: "CHANGE-ME"
  S3_SECRET_KEY: "CHANGE-ME"
`
}

func k8sConfigMap(cfg K8sConfig) string {
	appURL := "https://" + cfg.Domain
	if cfg.Domain == "" {
		appURL = "http://" + cfg.AppName + "-api." + cfg.Namespace + ".svc.cluster.local:8080"
	}
	return `# Everything that is not a credential. Changing a value here and rolling the
# Deployment is a normal config change; the Secret is where anything that must
# not be read goes.
apiVersion: v1
kind: ConfigMap
metadata:
  name: ` + cfg.AppName + `-config
  namespace: ` + cfg.Namespace + `
data:
  APP_ENV: "production"
  APP_PORT: "8080"
  # Whatever DB_PROVIDER a development .env says. SQLite in a container is a
  # file in that container, so every replica would have its own database and a
  # restart would start from an empty one.
  DB_PROVIDER: "postgres"
  APP_URL: "` + appURL + `"
`
}

// k8sAPIDeployment is the API, which is also every queue consumer.
func k8sAPIDeployment(cfg K8sConfig) string {
	return `# The API.
#
# It is also the worker. A Grit API process starts its own asynq worker
# whenever Redis is reachable, so there is no separate worker Deployment here
# and every replica consumes the queue. Two consequences worth knowing:
#
#   - scaling for HTTP load scales queue consumers too, which is usually what
#     you want and is occasionally not,
#   - a job that is slow or memory-hungry is slow and memory-hungry inside a
#     pod that is also serving requests.
#
# Separating them needs a worker-only mode in the binary, which does not exist
# yet. Ask on the issue tracker if you need one.
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ` + cfg.AppName + `-api
  namespace: ` + cfg.Namespace + `
  labels:
    app.kubernetes.io/name: ` + cfg.AppName + `
    app.kubernetes.io/component: api
spec:
  replicas: ` + fmt.Sprint(cfg.Replicas) + `
  selector:
    matchLabels:
      app.kubernetes.io/name: ` + cfg.AppName + `
      app.kubernetes.io/component: api
  template:
    metadata:
      labels:
        app.kubernetes.io/name: ` + cfg.AppName + `
        app.kubernetes.io/component: api
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 65532
        fsGroup: 65532
      containers:
        - name: api
          image: ` + cfg.Image + `
          imagePullPolicy: IfNotPresent
          ports:
            - name: http
              containerPort: 8080
          envFrom:
            - configMapRef:
                name: ` + cfg.AppName + `-config
            - secretRef:
                name: ` + cfg.AppName + `-secrets
          # /healthz and /readyz, not /api/health.
          #
          # The application serves all three. /api/health probes every
          # dependency it has, which is the right answer for the System Health
          # page and the wrong one for a probe: a slow object store would fail
          # readiness and take every pod out of the load balancer over
          # something that does not stop it serving requests. Liveness probes
          # nothing on purpose, and readiness gates on the database alone.
          livenessProbe:
            httpGet:
              path: /healthz
              port: http
            initialDelaySeconds: 10
            periodSeconds: 20
            timeoutSeconds: 3
            failureThreshold: 3
          readinessProbe:
            httpGet:
              path: /readyz
              port: http
            initialDelaySeconds: 3
            periodSeconds: 10
            timeoutSeconds: 3
            failureThreshold: 3
          # A startup probe as well, so a cold start with a slow database does
          # not trip the liveness probe and restart the pod forever.
          startupProbe:
            httpGet:
              path: /healthz
              port: http
            periodSeconds: 5
            failureThreshold: 30
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              memory: 512Mi
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities:
              drop:
                - ALL
          volumeMounts:
            - name: tmp
              mountPath: /tmp
      volumes:
        - name: tmp
          emptyDir: {}
`
}

func k8sFrontendDeployment(cfg K8sConfig, name, image string, port int) string {
	return `apiVersion: apps/v1
kind: Deployment
metadata:
  name: ` + cfg.AppName + `-` + name + `
  namespace: ` + cfg.Namespace + `
  labels:
    app.kubernetes.io/name: ` + cfg.AppName + `
    app.kubernetes.io/component: ` + name + `
spec:
  replicas: 2
  selector:
    matchLabels:
      app.kubernetes.io/name: ` + cfg.AppName + `
      app.kubernetes.io/component: ` + name + `
  template:
    metadata:
      labels:
        app.kubernetes.io/name: ` + cfg.AppName + `
        app.kubernetes.io/component: ` + name + `
    spec:
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
      containers:
        - name: ` + name + `
          image: ` + image + `
          imagePullPolicy: IfNotPresent
          ports:
            - name: http
              containerPort: ` + fmt.Sprint(port) + `
          envFrom:
            - configMapRef:
                name: ` + cfg.AppName + `-config
          readinessProbe:
            httpGet:
              path: /
              port: http
            initialDelaySeconds: 5
            periodSeconds: 10
          resources:
            requests:
              cpu: 50m
              memory: 128Mi
            limits:
              memory: 512Mi
          securityContext:
            allowPrivilegeEscalation: false
            capabilities:
              drop:
                - ALL
`
}

func k8sService(cfg K8sConfig, component string, port int) string {
	return `apiVersion: v1
kind: Service
metadata:
  name: ` + cfg.AppName + `-` + component + `
  namespace: ` + cfg.Namespace + `
  labels:
    app.kubernetes.io/name: ` + cfg.AppName + `
    app.kubernetes.io/component: ` + component + `
spec:
  type: ClusterIP
  selector:
    app.kubernetes.io/name: ` + cfg.AppName + `
    app.kubernetes.io/component: ` + component + `
  ports:
    - name: http
      port: ` + fmt.Sprint(port) + `
      targetPort: http
`
}

// k8sMigrateJob runs the schema migration once per image tag.
//
// Named for the tag rather than fixed, because a Job's pod template is
// immutable: applying a fixed-name Job a second time with a new image is
// rejected by the API server, so a deploy would fail on its first step with a
// message about an immutable field.
func k8sMigrateJob(cfg K8sConfig) string {
	return `# Run this BEFORE applying the rest, and wait for it:
#
#   kubectl apply -f migrate-job.yaml
#   kubectl wait --for=condition=complete --timeout=5m \
#     job/` + cfg.AppName + `-migrate-` + imageTag(cfg.Image) + ` -n ` + cfg.Namespace + `
#
# An API pod that starts against an un-migrated database answers "relation does
# not exist" to every request, so the order matters. It is not an initContainer
# on the API Deployment because that would run once per replica, concurrently.
#
# The name carries the image tag: a Job's pod template is immutable, so a
# fixed name would make the second deploy fail on an immutable field.
apiVersion: batch/v1
kind: Job
metadata:
  name: ` + cfg.AppName + `-migrate-` + imageTag(cfg.Image) + `
  namespace: ` + cfg.Namespace + `
  labels:
    app.kubernetes.io/name: ` + cfg.AppName + `
    app.kubernetes.io/component: migrate
spec:
  backoffLimit: 0
  ttlSecondsAfterFinished: 604800
  template:
    metadata:
      labels:
        app.kubernetes.io/name: ` + cfg.AppName + `
        app.kubernetes.io/component: migrate
    spec:
      restartPolicy: Never
      securityContext:
        runAsNonRoot: true
        runAsUser: 65532
      containers:
        - name: migrate
          image: ` + cfg.Image + `
          command: ["./migrate"]
          envFrom:
            - configMapRef:
                name: ` + cfg.AppName + `-config
            - secretRef:
                name: ` + cfg.AppName + `-secrets
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              memory: 512Mi
          securityContext:
            allowPrivilegeEscalation: false
            readOnlyRootFilesystem: true
            capabilities:
              drop:
                - ALL
`
}

func k8sHPA(cfg K8sConfig) string {
	return `# Scales the API on CPU.
#
# Remember that every API replica is also a queue consumer, so this scales
# background workers as a side effect of HTTP load. If that is wrong for your
# workload, pin replicas instead and delete this file.
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: ` + cfg.AppName + `-api
  namespace: ` + cfg.Namespace + `
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: ` + cfg.AppName + `-api
  minReplicas: ` + fmt.Sprint(cfg.Replicas) + `
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
`
}

func k8sIngress(cfg K8sConfig) string {
	rules := `    - host: ` + cfg.Domain + `
      http:
        paths:
          - path: /api
            pathType: Prefix
            backend:
              service:
                name: ` + cfg.AppName + `-api
                port:
                  name: http
`
	// The API owns more than /api: the websocket, the probes and the JWKS
	// document all sit outside it, and routing only /api leaves a frontend
	// talking to itself for those.
	for _, path := range []string{"/healthz", "/readyz", "/.well-known/jwks.json"} {
		rules += `          - path: ` + path + `
            pathType: Prefix
            backend:
              service:
                name: ` + cfg.AppName + `-api
                port:
                  name: http
`
	}

	root := cfg.AppName + "-api"
	if cfg.WebImage != "" {
		root = cfg.AppName + "-web"
	}
	rules += `          - path: /
            pathType: Prefix
            backend:
              service:
                name: ` + root + `
                port:
                  name: http
`

	admin := ""
	if cfg.AdminImage != "" {
		admin = `    - host: admin.` + cfg.Domain + `
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: ` + cfg.AppName + `-admin
                port:
                  name: http
`
	}

	return `# Assumes an ingress controller and cert-manager. Change the class and the
# issuer to whatever your cluster runs, or delete this file and route traffic
# the way you already do.
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ` + cfg.AppName + `
  namespace: ` + cfg.Namespace + `
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  ingressClassName: nginx
  tls:
    - hosts:
        - ` + cfg.Domain + `
      secretName: ` + cfg.AppName + `-tls
  rules:
` + rules + admin
}

func k8sKustomization(cfg K8sConfig, files map[string]string) string {
	var resources []string
	for _, name := range sortedKeys(files) {
		switch name {
		case "README.md", "kustomization.yaml", "secret.example.yaml", "migrate-job.yaml":
			// The Secret is a template nobody should apply by accident, and
			// the migration Job runs on its own, before the rest.
			continue
		}
		resources = append(resources, "  - "+name)
	}

	return `# kubectl apply -k .
#
# secret.example.yaml and migrate-job.yaml are deliberately absent:
#
#   - the Secret holds placeholders, and applying it would start the API with
#     a signing key anybody can forge, looking perfectly healthy,
#   - the Job has to run and finish before the API rolls, so it is applied on
#     its own and waited for.
#
# The README beside this file has the order.
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: ` + cfg.Namespace + `
# labels rather than the deprecated commonLabels, and with
# includeSelectors false on purpose. commonLabels writes into
# spec.selector.matchLabels, which is immutable: renaming the app would
# then make every later apply fail on a field that cannot be changed.
# These Deployments already declare their own selectors.
labels:
  - includeSelectors: false
    pairs:
      app.kubernetes.io/name: ` + cfg.AppName + `
      app.kubernetes.io/managed-by: grit
resources:
` + strings.Join(resources, "\n") + `
`
}

func k8sReadme(cfg K8sConfig) string {
	tag := imageTag(cfg.Image)
	return `# ` + cfg.AppName + ` on Kubernetes

Generated by ` + "`grit deploy --k8s`" + `. Plain manifests and a kustomization, so
` + "`kubectl apply -k`" + ` works with nothing else installed.

## Before the first deploy

1. **Build and push the images.** The Dockerfiles are the ones already in the
   repository, under ` + "`apps/`" + `.

2. **Create the secret.** ` + "`secret.example.yaml`" + ` is a template and is not in
   ` + "`kustomization.yaml`" + `, so ` + "`apply -k`" + ` cannot ship its placeholders.

   ` + "```" + `
   cp secret.example.yaml secret.yaml
   $EDITOR secret.yaml
   kubectl apply -f secret.yaml
   ` + "```" + `

3. **Bring your own Postgres and Redis.** These manifests do not install
   either. A generated single-replica database with no backups looks
   production-ready and is not, so the Secret names yours instead.

   Redis is not optional: the API starts its asynq worker only when Redis
   answers, so without it no background job ever runs.

## Every deploy

` + "```" + `
kubectl apply -f migrate-job.yaml
kubectl wait --for=condition=complete --timeout=5m \
  job/` + cfg.AppName + `-migrate-` + tag + ` -n ` + cfg.Namespace + `
kubectl apply -k .
` + "```" + `

The migration runs first and is waited for, because an API pod that starts
against an un-migrated database answers "relation does not exist" to every
request. It is not an initContainer, which would run once per replica at the
same time as its siblings.

The Job's name carries the image tag. A Job's pod template is immutable, so a
fixed name makes the second deploy fail on an immutable field rather than
running.

## The API is also the worker

There is no worker Deployment here. A Grit API process starts its own asynq
worker whenever Redis is reachable, so every replica consumes the queue.
Scaling for HTTP load scales queue consumers with it, and a heavy job runs
inside a pod that is also serving requests. Separating them needs a
worker-only mode in the binary, which does not exist yet.

## Probes

Liveness is ` + "`/healthz`" + ` and readiness is ` + "`/readyz`" + `, not ` + "`/api/health`" + `.
All three are served. ` + "`/api/health`" + ` probes every dependency, which is right
for the System Health page and wrong for a probe: a slow object store would
fail readiness and pull every pod out of the load balancer over something that
does not stop it serving requests.

## What to change

- ` + "`ingress.yaml`" + ` assumes an nginx ingress controller and cert-manager.
- Resource requests and limits are a starting point, not a measurement.
- ` + "`api-hpa.yaml`" + ` scales on CPU. Delete it to pin replicas.
`
}
