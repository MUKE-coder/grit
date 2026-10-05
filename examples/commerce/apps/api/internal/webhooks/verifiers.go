package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// HMACVerifier returns a VerifyFunc that validates a hex-encoded
// HMAC-SHA256 signature found in the named header. Most simple
// providers (custom partners, self-rolled webhooks) use this scheme.
//
//	webhooks.Register("partner", webhooks.Provider{
//	    SecretEnv: "PARTNER_WEBHOOK_SECRET",
//	    Verify:    webhooks.HMACVerifier("X-Signature"),
//	    Extract:   webhooks.JSONFieldExtractor("type", "id"),
//	})
func HMACVerifier(header string) VerifyFunc {
	return func(secret string, body []byte, headers map[string]string) error {
		if secret == "" {
			return fmt.Errorf("webhooks: signing secret is empty")
		}
		got := headers[header]
		if got == "" {
			got = headers[strings.ToLower(header)]
		}
		if got == "" {
			return fmt.Errorf("webhooks: missing signature header %q", header)
		}
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expected := hex.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(got), []byte(expected)) {
			return fmt.Errorf("webhooks: signature mismatch")
		}
		return nil
	}
}

// StripeVerifier validates Stripe's "Stripe-Signature" header, which
// has the form "t=<unix>,v1=<hex>" where v1 = HMAC-SHA256 of
// "<timestamp>.<payload>" using the webhook signing secret. Tolerance
// of 5 minutes guards against replay.
//
// See https://stripe.com/docs/webhooks/signatures
func StripeVerifier(secret string, body []byte, headers map[string]string) error {
	const tolerance = 5 * time.Minute
	if secret == "" {
		return fmt.Errorf("webhooks: stripe secret is empty")
	}
	header := headers["Stripe-Signature"]
	if header == "" {
		header = headers["stripe-signature"]
	}
	if header == "" {
		return fmt.Errorf("webhooks: missing Stripe-Signature header")
	}

	var ts int64
	var sigs []string
	for _, part := range strings.Split(header, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "t":
			ts, _ = strconv.ParseInt(kv[1], 10, 64)
		case "v1":
			sigs = append(sigs, kv[1])
		}
	}
	if ts == 0 || len(sigs) == 0 {
		return fmt.Errorf("webhooks: malformed Stripe-Signature header")
	}
	if time.Since(time.Unix(ts, 0)) > tolerance {
		return fmt.Errorf("webhooks: stripe timestamp outside tolerance")
	}

	signed := strconv.FormatInt(ts, 10) + "." + string(body)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signed))
	expected := hex.EncodeToString(mac.Sum(nil))
	for _, s := range sigs {
		if hmac.Equal([]byte(s), []byte(expected)) {
			return nil
		}
	}
	return fmt.Errorf("webhooks: stripe signature mismatch")
}

// GitHubVerifier validates GitHub's "X-Hub-Signature-256" header,
// which is "sha256=<hex>" — HMAC-SHA256 of the raw body using the
// webhook secret.
func GitHubVerifier(secret string, body []byte, headers map[string]string) error {
	if secret == "" {
		return fmt.Errorf("webhooks: github secret is empty")
	}
	header := headers["X-Hub-Signature-256"]
	if header == "" {
		header = headers["x-hub-signature-256"]
	}
	if header == "" {
		return fmt.Errorf("webhooks: missing X-Hub-Signature-256 header")
	}
	prefix := "sha256="
	if !strings.HasPrefix(header, prefix) {
		return fmt.Errorf("webhooks: unexpected X-Hub-Signature-256 format")
	}
	got := header[len(prefix):]
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(got), []byte(expected)) {
		return fmt.Errorf("webhooks: github signature mismatch")
	}
	return nil
}

// JSONFieldExtractor returns an ExtractFunc that pulls type + id from
// top-level JSON fields in the body. Stripe-style payloads use
// JSONFieldExtractor("type", "id") — the most common shape.
func JSONFieldExtractor(typeField, idField string) ExtractFunc {
	return func(body []byte, headers map[string]string) (string, string, error) {
		var raw map[string]interface{}
		if err := json.Unmarshal(body, &raw); err != nil {
			return "", "", fmt.Errorf("decoding payload: %w", err)
		}
		t, _ := raw[typeField].(string)
		id, _ := raw[idField].(string)
		return t, id, nil
	}
}

// StripeExtractor pulls (type, id) from Stripe's standard
// { "type": "...", "id": "evt_..." } envelope.
var StripeExtractor = JSONFieldExtractor("type", "id")

// GitHubExtractor reads the event type from the "X-GitHub-Event"
// header and the delivery ID from "X-GitHub-Delivery".
func GitHubExtractor(body []byte, headers map[string]string) (string, string, error) {
	t := headers["X-GitHub-Event"]
	if t == "" {
		t = headers["x-github-event"]
	}
	id := headers["X-GitHub-Delivery"]
	if id == "" {
		id = headers["x-github-delivery"]
	}
	return t, id, nil
}
