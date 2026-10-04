// Package sync owns the model registry used by the offline-first
// /api/sync/push and /api/sync/pull endpoints.
//
// Every model that should be syncable from a desktop client must:
//  1. Have an ID string (UUID) primary key.
//  2. Have a Version int field.
//  3. Have CreatedAt / UpdatedAt timestamps.
//  4. Have a BeforeUpdate hook that increments Version.
//  5. Be registered with Register("table_name", &models.X{}).
//
// The handler uses reflection to decode push payloads into the
// registered struct type, run a versioned update, and detect conflicts
// when the client's version doesn't match what's on disk.
package sync

import (
	"fmt"
	"reflect"
	"sync"
)

// Registry holds the syncable model types keyed by their plural snake_case
// name (e.g. "buildings"). Population happens at app boot from routes.Setup.
type Registry struct {
	mu     sync.RWMutex
	models map[string]reflect.Type
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{models: make(map[string]reflect.Type)}
}

// Register adds a model under its plural-snake table name. proto must be
// a pointer to a zero-value struct (e.g. &models.Building{}).
func (r *Registry) Register(table string, proto interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t := reflect.TypeOf(proto)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	r.models[table] = t
}

// New returns a new pointer to a zero-value model struct for the given
// table, or an error if the table isn't registered.
func (r *Registry) New(table string) (interface{}, error) {
	r.mu.RLock()
	t, ok := r.models[table]
	r.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("sync: unknown table %q", table)
	}
	return reflect.New(t).Interface(), nil
}

// Tables lists every registered table name. Used by /api/sync/pull when
// the client asks for the full set of types.
func (r *Registry) Tables() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.models))
	for k := range r.models {
		out = append(out, k)
	}
	return out
}
