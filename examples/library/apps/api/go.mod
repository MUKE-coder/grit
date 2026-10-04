module library/apps/api

// Go 1.26.6, and not an earlier patch. Go 1.26.4 carries eight standard library
// vulnerabilities this code reaches, among them net/http and the encoding/xml
// behind the SAML sign-in path, before anyone has signed in. The go directive is
// also what CI, setup-go and the release workflow install, and an older Go
// downloads this one, so raising it here raises it everywhere.
go 1.26.6

require (
	// Pure-Go WebP, so image optimisation does not cost the static
	// cross-compiled binary. Lossless (VP8L) only: there is no pure-Go
	// lossy WebP encoder, which is why lossy compression targets JPEG.
	github.com/HugoSmits86/nativewebp v1.3.0
	github.com/MUKE-coder/gin-docs v0.0.0-20260222113017-4d647cb4e7aa
	github.com/MUKE-coder/gorm-studio v1.1.0
	github.com/MUKE-coder/pulse v1.0.0
	// Sentinel now ships a proper /v2 module path, so we track real tags.
	// v2.1.1 is the minimum safe release for WAF.Mode = ModeBlock: v2.1.0
	// fixed the SSRF rule matching "0.0.0.0" inside a Chrome User-Agent
	// (403'ing every Chrome 140/130/120/110 user), and v2.1.1 fixed
	// SQLi_Basic matching a bare "--" inside JWT cookies (roughly one
	// session in ten 403'd at random). v2.2.0 adds ValidateConfig, which
	// Mount runs at startup so dead config shows up in the boot log rather
	// than as a 403 weeks later. Do not downgrade below v2.1.1.
	// v2.2.2 is a security release (client-IP spoofing behind a proxy, a
	// sort_by SQL injection, SSRF bypasses), and v2.5.0 keeps rate limits and
	// lockouts in Redis, so every replica counts against the same numbers.
	github.com/MUKE-coder/sentinel/v2 v2.5.0
	github.com/aws/aws-sdk-go-v2 v1.43.0
	github.com/aws/aws-sdk-go-v2/config v1.32.31
	github.com/aws/aws-sdk-go-v2/credentials v1.19.30
	github.com/aws/aws-sdk-go-v2/service/s3 v1.106.0
	// The QR code on the 2FA setup screen. Replaced skip2/go-qrcode, last
	// released in 2020.
	github.com/boombuler/barcode v1.1.0
	github.com/brianvoe/gofakeit/v7 v7.15.0
	// SAML 2.0 service provider for enterprise SSO. OIDC covers every modern
	// IdP and needs no library, but SAML is still what a lot of enterprise
	// procurement asks for, and it cannot be hand-rolled safely — assertion
	// signature verification, audience restriction and clock-skew handling are
	// exactly the places a DIY implementation becomes an auth bypass.
	github.com/crewjam/saml v0.5.1
	github.com/disintegration/imaging v1.6.2
	github.com/gin-gonic/gin v1.11.0
	github.com/go-pdf/fpdf v1.4.3
	// Pure Go WebAuthn, so passkeys do not cost the static binary.
	github.com/go-webauthn/webauthn v0.18.0
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	// Linked in through goth's gothic, which still asks for v1.6.2 from 2018.
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/gorilla/sessions v1.4.0
	github.com/gorilla/websocket v1.5.3
	github.com/hibiken/asynq v0.24.1
	github.com/joho/godotenv v1.5.1
	github.com/markbates/goth v1.80.0
	// The HTML sanitiser behind internal/sanitize. Rich text is rendered as HTML
	// by the blog and by anything built on a richtext field, so it is cleaned on
	// the way in rather than trusted on the way out.
	github.com/microcosm-cc/bluemonday v1.0.27
	github.com/redis/go-redis/v9 v9.22.0
	// CVE-2026-54063 and CVE-2026-59161: the worksheet and streaming
	// parsers could be made to allocate without bound, and this is what
	// the CSV/XLSX importer hands user uploads to.
	github.com/xuri/excelize/v2 v2.11.0
	// OpenTelemetry, for internal/tracing. Four modules and no more: the
	// OTLP/HTTP exporter rather than both transports, and a hand-written gin
	// middleware rather than otelgin, because the span has to be named after
	// the route pattern and carry this project's request id, which a
	// general-purpose middleware cannot know about.
	//
	// Inert until OTEL_EXPORTER_OTLP_ENDPOINT is set: with no endpoint the
	// global provider is a no-op and a span costs an allocation.
	go.opentelemetry.io/otel v1.45.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.45.0
	go.opentelemetry.io/otel/sdk v1.45.0
	go.opentelemetry.io/otel/trace v1.45.0
	golang.org/x/crypto v0.57.0
	// singleflight, which collapses concurrent cache misses in middleware/cache.go.
	golang.org/x/sync v0.23.0
	// Pulled in by the OTLP exporter even on the HTTP transport: the
	// collector proto package depends on it. Pinned rather than left to
	// whatever the otel modules ask for, which was v1.77.0 and carries two
	// reachable advisories, and v1.83.1 carries CVE-2026-84445.
	google.golang.org/grpc v1.83.2 // indirect
	gorm.io/datatypes v1.2.7
	gorm.io/driver/mysql v1.6.0
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.1
	gorm.io/plugin/dbresolver v1.6.2
)

require (
	github.com/davidbyttow/govips/v2 v2.19.0
	github.com/glebarez/sqlite v1.11.0
	github.com/stretchr/testify v1.12.1
	golang.org/x/net v0.58.0
	gorm.io/driver/sqlite v1.6.0
)

require (
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.14 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.31 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.31 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.31 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.4.32 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.13 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.9.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.31 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.19.32 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.5.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.33.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.38.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.45.0 // indirect
	github.com/aws/smithy-go v1.27.3 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beevik/etree v1.6.0 // indirect
	github.com/bytedance/sonic v1.14.0 // indirect
	github.com/bytedance/sonic/loader v0.3.0 // indirect
	github.com/cenkalti/backoff/v5 v5.0.3 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/ebitengine/purego v0.10.0 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/fxamacker/cbor/v2 v2.9.3 // indirect
	github.com/gabriel-vasile/mimetype v1.4.8 // indirect
	github.com/gin-contrib/cors v1.7.2 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/glebarez/go-sqlite v1.21.2 // indirect
	github.com/go-logr/logr v1.4.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/go-webauthn/x v0.3.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/goccy/go-yaml v1.18.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.2 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/go-tpm v0.9.8 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.29.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/mattermost/xml-roundtrip-validator v0.1.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/richardlehane/mscfb v1.0.7 // indirect
	github.com/richardlehane/msoleps v1.0.6 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/shirou/gopsutil/v4 v4.26.4 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/tiendc/go-deepcopy v1.7.2 // indirect
	github.com/tinylib/msgp v1.6.4 // indirect
	github.com/tklauser/go-sysconf v0.3.16 // indirect
	github.com/tklauser/numcpus v0.11.0 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	github.com/xuri/efp v0.0.1 // indirect
	github.com/xuri/nfp v0.0.2-0.20250530014748-2ddeb826f9a9 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.45.0 // indirect
	go.opentelemetry.io/otel/metric v1.45.0 // indirect
	go.opentelemetry.io/proto/otlp v1.11.0 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/arch v0.20.0 // indirect
	golang.org/x/exp v0.0.0-20251023183803-a4bb9ffd2546 // indirect
	golang.org/x/sys v0.48.0 // indirect
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20260803160001-6ac0973c030d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260803160001-6ac0973c030d // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	modernc.org/libc v1.67.6 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
	modernc.org/sqlite v1.46.1 // indirect
)

// Security floors for transitive dependencies. These are not imported directly;
// they are pinned because a dependency pulls in a version with a known,
// reachable vulnerability and MVS would otherwise settle on it. Each one was
// confirmed with govulncheck against a freshly scaffolded project.
//
// Raise these, never lower them. Re-check with:
//   go install golang.org/x/vuln/cmd/govulncheck@latest && govulncheck ./...
require (
	// Pulled in by the MySQL driver; govulncheck flags releases before v1.1.1.
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/jackc/pgx/v5 v5.9.2 // indirect; GO-2026-5004
	github.com/quic-go/quic-go v0.59.1 // indirect; GO-2026-5676, GO-2025-4233
	// CVE-2026-33487: XML Digital Signature validation could be bypassed. On the
	// SAML assertion path that is an authentication bypass, and crewjam/saml
	// v0.5.1 is its newest release and still asks for the vulnerable version.
	github.com/russellhaering/goxmldsig v1.6.0 // CVE-2026-33487
	golang.org/x/image v0.45.0 // indirect; GO-2026-5066, -5062, -5032, -5031, -4815, CVE-2026-46603
	golang.org/x/oauth2 v0.36.0 // indirect; CVE-2025-22868
	golang.org/x/text v0.42.0 // indirect; GO-2026-5970
)
