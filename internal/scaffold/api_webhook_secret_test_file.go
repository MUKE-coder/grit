package scaffold

// apiWebhookSecretTestGo emits internal/webhooks/secret_test.go.
//
// A provider whose signing secret differs per organization had no way into the
// receiver before v3.446.0: Provider carried SecretEnv, one env var for the
// whole deployment, and the receiver read it with os.Getenv. The only recourse
// was a hand-written route, which loses the stored deliveries, the dedupe and
// the admin replay (grit#138).
//
// These are shipped rather than kept in the framework because the thing worth
// protecting is in the project: a resolver that falls back to a default secret
// when it does not recognise a tenant verifies one organization's deliveries
// against another's key.
func apiWebhookSecretTestGo() string {
	return `package webhooks

import (
	"context"
	"errors"
	"os"
	"testing"
)

// Provider.Secret decides, and SecretEnv is the fallback. A provider that has
// both must not read the environment: that is one secret for the deployment,
// which is the thing per-tenant verification exists to escape.
func TestSecretForPrefersTheResolver(t *testing.T) {
	t.Setenv("TEST_WEBHOOK_SECRET", "deployment-wide")

	p := Provider{
		SecretEnv: "TEST_WEBHOOK_SECRET",
		Secret: func(_ context.Context, params, _ map[string]string) (string, error) {
			return "secret-for-" + params["tenant"], nil
		},
	}
	got, err := SecretFor(context.Background(), p, map[string]string{"tenant": "org-a"}, nil)
	if err != nil {
		t.Fatalf("resolving: %v", err)
	}
	if got != "secret-for-org-a" {
		t.Errorf("got %q, want the resolver's answer, not the environment's", got)
	}
}

// Without a resolver, the environment, which is right for Stripe and GitHub.
func TestSecretForFallsBackToTheEnvironment(t *testing.T) {
	t.Setenv("TEST_WEBHOOK_SECRET", "from-the-env")

	p := Provider{SecretEnv: "TEST_WEBHOOK_SECRET"}
	got, err := SecretFor(context.Background(), p, nil, nil)
	if err != nil {
		t.Fatalf("resolving: %v", err)
	}
	if got != "from-the-env" {
		t.Errorf("got %q, want %q", got, os.Getenv("TEST_WEBHOOK_SECRET"))
	}
}

// A resolver that does not recognise the tenant must return an error, and the
// error must reach the caller rather than becoming an empty secret. An empty
// secret is not "no verification": it is verification against "", which some
// HMAC verifiers will happily accept for a payload signed with "".
func TestSecretForReturnsTheResolversError(t *testing.T) {
	sentinel := errors.New("no such tenant")
	p := Provider{
		SecretEnv: "TEST_WEBHOOK_SECRET",
		Secret: func(context.Context, map[string]string, map[string]string) (string, error) {
			return "", sentinel
		},
	}
	if _, err := SecretFor(context.Background(), p, nil, nil); !errors.Is(err, sentinel) {
		t.Errorf("got %v, want the resolver's error: a delivery whose secret cannot "+
			"be resolved must be refused, not verified against an empty string", err)
	}
}
`
}
