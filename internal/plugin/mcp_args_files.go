package plugin

// mcpArgsGo is internal/mcp/args.go: reading what a model sent.
func mcpArgsGo(ctx Context) string {
	return `package mcp

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Reading a model's arguments is not like reading a form's.
//
// A schema says a field is a number and the model sends "5". It says the field
// is an id and the model sends the customer's name, because that is what the
// person typed. Refusing both is correct by the schema and useless in
// practice: the call fails, the model apologises, and the person is no further
// forward.
//
// So these helpers are deliberately forgiving about shape and strict about
// meaning. A string that is a number becomes a number. A name becomes an id,
// or an answer naming the candidates when it is ambiguous, which the model can
// put back to the person. What they never do is guess which record was meant.

// String is a required string argument.
func String(args map[string]any, key string) (string, error) {
	value, present := args[key]
	if !present || value == nil {
		return "", fmt.Errorf("%s is required", key)
	}
	switch typed := value.(type) {
	case string:
		trimmed := strings.TrimSpace(typed)
		if trimmed == "" {
			return "", fmt.Errorf("%s is required", key)
		}
		return trimmed, nil
	case float64:
		// A model sending an id it read out of a previous answer as a number.
		return strconv.FormatFloat(typed, 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(typed), nil
	}
	return "", fmt.Errorf("%s should be text", key)
}

// OptionalString is a string that may be absent, and "" when it is.
func OptionalString(args map[string]any, key string) string {
	value, err := String(args, key)
	if err != nil {
		return ""
	}
	return value
}

// Int is an integer argument, accepting the string a model often sends.
func Int(args map[string]any, key string) (int, error) {
	value, present := args[key]
	if !present || value == nil {
		return 0, fmt.Errorf("%s is required", key)
	}
	switch typed := value.(type) {
	case float64:
		return int(typed), nil
	case int:
		return typed, nil
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(typed))
		if err != nil {
			return 0, fmt.Errorf("%s should be a whole number, not %q", key, typed)
		}
		return parsed, nil
	}
	return 0, fmt.Errorf("%s should be a whole number", key)
}

// OptionalInt is Int with a default.
func OptionalInt(args map[string]any, key string, fallback int) int {
	value, err := Int(args, key)
	if err != nil {
		return fallback
	}
	return value
}

// Float is a number argument. Money arrives here in ordinary units: a tool
// converts to minor units on the way in, because a model asked for "forty
// dollars" will send 40 and not 4000, every time.
func Float(args map[string]any, key string) (float64, error) {
	value, present := args[key]
	if !present || value == nil {
		return 0, fmt.Errorf("%s is required", key)
	}
	switch typed := value.(type) {
	case float64:
		return typed, nil
	case int:
		return float64(typed), nil
	case string:
		cleaned := strings.TrimSpace(strings.NewReplacer(",", "", "$", "", " ", "").Replace(typed))
		parsed, err := strconv.ParseFloat(cleaned, 64)
		if err != nil {
			return 0, fmt.Errorf("%s should be a number, not %q", key, typed)
		}
		return parsed, nil
	}
	return 0, fmt.Errorf("%s should be a number", key)
}

// Bool is a boolean argument, accepting "true" and "yes".
func Bool(args map[string]any, key string, fallback bool) bool {
	value, present := args[key]
	if !present || value == nil {
		return fallback
	}
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		switch strings.ToLower(strings.TrimSpace(typed)) {
		case "true", "yes", "y", "1":
			return true
		case "false", "no", "n", "0":
			return false
		}
	case float64:
		return typed != 0
	}
	return fallback
}

// Strings is a list argument, accepting a single value as a list of one.
//
// A model told a field takes a list will quite often send the bare item, and
// refusing that is pedantry: the intent is unambiguous.
func Strings(args map[string]any, key string) []string {
	value, present := args[key]
	if !present || value == nil {
		return nil
	}
	switch typed := value.(type) {
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			if text, ok := item.(string); ok && strings.TrimSpace(text) != "" {
				out = append(out, strings.TrimSpace(text))
			}
		}
		return out
	case []string:
		return typed
	case string:
		if trimmed := strings.TrimSpace(typed); trimmed != "" {
			return []string{trimmed}
		}
	}
	return nil
}

// Page turns a model's paging arguments into the query the API takes, with
// limits a model cannot talk its way past.
//
// The cap matters: "list every invoice" with no page size is a model asking
// for a hundred thousand rows it will then put in its own context. 100 is
// enough to answer a question and small enough to be affordable.
func Page(args map[string]any) url.Values {
	query := url.Values{}

	page := OptionalInt(args, "page", 1)
	if page < 1 {
		page = 1
	}
	size := OptionalInt(args, "page_size", 20)
	switch {
	case size < 1:
		size = 20
	case size > 100:
		size = 100
	}

	query.Set("page", strconv.Itoa(page))
	query.Set("page_size", strconv.Itoa(size))

	if search := OptionalString(args, "search"); search != "" {
		query.Set("search", search)
	}
	if sort := OptionalString(args, "sort"); sort != "" {
		query.Set("sort", sort)
	}
	return query
}

// Ambiguous is the answer when a name matched more than one record.
//
// An error and not a guess. Picking the first match is the behaviour that
// quietly invoices the wrong customer, and a model handed the candidates can
// ask the person which one they meant, which is the only correct resolution.
type Ambiguous struct {
	What       string
	Term       string
	Candidates []string
}

func (e *Ambiguous) Error() string {
	return fmt.Sprintf("%q matches more than one %s: %s. Ask which one is meant, "+
		"or pass its id.", e.Term, e.What, strings.Join(e.Candidates, ", "))
}

// Resolve turns a name into an id, by asking the app.
//
// Through the router like everything else, so a lookup is scoped to what this
// connection may see: a name that matches a record the person cannot read is
// not found, which is the same answer they would get in the browser.
//
// Returns the value unchanged when it already looks like an id, so a tool can
// pass whatever it was given without deciding first.
func Resolve(ctx context.Context, call *Call, what, path, term string) (string, error) {
	if term == "" {
		return "", fmt.Errorf("no %s given", what)
	}
	if looksLikeID(term) {
		return term, nil
	}

	query := url.Values{}
	query.Set("search", term)
	query.Set("page_size", "10")

	result, err := call.Get(ctx, path, query)
	if err != nil {
		return "", err
	}
	if !result.OK() {
		return "", fmt.Errorf("looking up the %s %q: %s", what, term, result.Message())
	}

	var rows []struct {
		ID    string ` + "`" + `json:"id"` + "`" + `
		Name  string ` + "`" + `json:"name"` + "`" + `
		Title string ` + "`" + `json:"title"` + "`" + `
		Label string ` + "`" + `json:"label"` + "`" + `
	}
	if err := result.Decode(&rows); err != nil {
		return "", fmt.Errorf("looking up the %s %q: %w", what, term, err)
	}

	switch len(rows) {
	case 0:
		return "", fmt.Errorf("no %s called %q", what, term)
	case 1:
		return rows[0].ID, nil
	}

	// More than one row came back, but search is a substring match: an exact
	// match among them is the answer, and only a genuine tie is ambiguous.
	var exact []string
	for _, row := range rows {
		if strings.EqualFold(labelOf(row.Name, row.Title, row.Label), term) {
			exact = append(exact, row.ID)
		}
	}
	if len(exact) == 1 {
		return exact[0], nil
	}

	candidates := make([]string, 0, len(rows))
	for _, row := range rows {
		candidates = append(candidates, labelOf(row.Name, row.Title, row.Label))
	}
	return "", &Ambiguous{What: what, Term: term, Candidates: candidates}
}

func labelOf(options ...string) string {
	for _, option := range options {
		if option != "" {
			return option
		}
	}
	return ""
}

// looksLikeID recognises the UUID every Grit model uses for its primary key.
//
// On the shape and not on a lookup, so passing an id costs no request. A name
// that happens to be 36 characters with hyphens in the right places does not
// exist in practice.
func looksLikeID(value string) bool {
	if len(value) != 36 {
		return false
	}
	for i, ch := range value {
		switch i {
		case 8, 13, 18, 23:
			if ch != '-' {
				return false
			}
		default:
			isHex := (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
			if !isHex {
				return false
			}
		}
	}
	return true
}
`
}
