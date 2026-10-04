package mail

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// awsCredentials are what a SigV4 signature is made with.
type awsCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

// signV4 signs req with AWS Signature Version 4, following
// https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_sigv-create-signed-request.html
// Every header already on the request is signed, with host. payload must be
// the request body.
func signV4(req *http.Request, payload []byte, creds awsCredentials, region, service string, now time.Time) {
	amzDate := now.UTC().Format("20060102T150405Z")
	req.Header.Set("X-Amz-Date", amzDate)
	if creds.SessionToken != "" {
		req.Header.Set("X-Amz-Security-Token", creds.SessionToken)
	}
	canonical, signedHeaders := canonicalRequest(req, payload)
	scope := amzDate[:8] + "/" + region + "/" + service + "/aws4_request"
	key := signingKey(creds.SecretAccessKey, amzDate[:8], region, service)
	signature := hex.EncodeToString(hmacSHA256(key, stringToSign(amzDate, scope, canonical)))
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential="+creds.AccessKeyID+"/"+scope+
		", SignedHeaders="+signedHeaders+", Signature="+signature)
}

// canonicalRequest returns the canonical request and its signed-headers list.
func canonicalRequest(req *http.Request, payload []byte) (string, string) {
	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	values := map[string][]string{"host": {host}}
	for name, v := range req.Header {
		lower := strings.ToLower(name)
		if lower == "authorization" {
			continue
		}
		values[lower] = v
	}
	names := make([]string, 0, len(values))
	for name := range values {
		names = append(names, name)
	}
	sort.Strings(names)

	var headers strings.Builder
	for _, name := range names {
		trimmed := make([]string, 0, len(values[name]))
		for _, v := range values[name] {
			trimmed = append(trimmed, strings.Join(strings.Fields(v), " "))
		}
		headers.WriteString(name + ":" + strings.Join(trimmed, ",") + "\n")
	}
	signedHeaders := strings.Join(names, ";")

	path := req.URL.EscapedPath()
	if path == "" {
		path = "/"
	}
	canonical := strings.Join([]string{
		req.Method,
		path,
		canonicalQuery(req.URL),
		headers.String(),
		signedHeaders,
		sha256Hex(payload),
	}, "\n")
	return canonical, signedHeaders
}

func canonicalQuery(u *url.URL) string {
	query := u.Query()
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var pairs []string
	for _, k := range keys {
		vals := append([]string(nil), query[k]...)
		sort.Strings(vals)
		for _, v := range vals {
			pairs = append(pairs, awsEscape(k)+"="+awsEscape(v))
		}
	}
	return strings.Join(pairs, "&")
}

// awsEscape percent-encodes everything but A-Z a-z 0-9 - _ . ~, with a space
// as %20 rather than +.
func awsEscape(s string) string {
	return strings.ReplaceAll(url.QueryEscape(s), "+", "%20")
}

func stringToSign(amzDate, scope, canonical string) string {
	return "AWS4-HMAC-SHA256\n" + amzDate + "\n" + scope + "\n" + sha256Hex([]byte(canonical))
}

func signingKey(secret, date, region, service string) []byte {
	k := hmacSHA256([]byte("AWS4"+secret), date)
	k = hmacSHA256(k, region)
	k = hmacSHA256(k, service)
	return hmacSHA256(k, "aws4_request")
}

func hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}
