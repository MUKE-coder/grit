package storage

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"

	"saas/apps/api/internal/config"
)

// ErrNotConfigured is returned, wrapped, by Open for a bucket driver without
// the settings it needs to connect.
var ErrNotConfigured = errors.New("storage: not configured")

// Open builds the store a driver names: a directory for "local", a bucket for
// "minio", "s3", "r2" and "b2".
//
// AWS S3 needs no endpoint, and with no access key the SDK's own credential
// chain supplies one, so an IAM role on EC2, ECS or Lambda works with only
// S3_BUCKET set. MinIO, R2 and B2 need an endpoint and a key.
func Open(driver string, cfg config.StorageConfig, local LocalConfig) (*Storage, error) {
	switch driver {
	case "local":
		return NewLocal(local)
	case "s3":
		return New(cfg)
	case "minio", "r2", "b2", "":
		if cfg.Endpoint == "" || cfg.AccessKey == "" {
			return nil, fmt.Errorf("%w: STORAGE_DRIVER=%s needs an endpoint and an access key", ErrNotConfigured, driver)
		}
		return New(cfg)
	default:
		return nil, fmt.Errorf("storage: unknown STORAGE_DRIVER %q, use local, minio, s3, r2 or b2", driver)
	}
}

// Disks holds every store the app opened: the default one, and each named disk
// in STORAGE_DISKS. main.go fills it at startup, before anything is served.
//
//	backups := storage.Disks.Get("backups") // nil when STORAGE_DISKS has no backups
var Disks = &Registry{}

// Registry is a set of stores by name.
type Registry struct {
	mu    sync.RWMutex
	def   *Storage
	named map[string]*Storage
}

// SetDefault records the store the app uses when it names none.
func (r *Registry) SetDefault(s *Storage) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.def = s
}

// Add records a named store. "" and "default" set the default one.
func (r *Registry) Add(name string, s *Storage) {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" || name == "default" {
		r.SetDefault(s)
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.named == nil {
		r.named = map[string]*Storage{}
	}
	r.named[name] = s
}

// Default is the store the app uses when it names none, or nil.
func (r *Registry) Default() *Storage {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.def
}

// Get returns the store called name, or nil when none is configured. "" and
// "default" are the default store. A missing name does not fall back to the
// default: code asking for a private bucket should not quietly write to the
// public one.
func (r *Registry) Get(name string) *Storage {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" || name == "default" {
		return r.Default()
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.named[name]
}

// Names lists the named stores, sorted.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.named))
	for name := range r.named {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
