package models

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
)

// slugify generates a URL-friendly slug with a unique suffix.
func slugify(s string) string {
	slug := strings.ToLower(s)
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	b := make([]byte, 4)
	rand.Read(b)
	return slug + "-" + hex.EncodeToString(b)
}
