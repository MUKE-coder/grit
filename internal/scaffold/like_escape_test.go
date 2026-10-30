package scaffold

import (
	"strings"
	"testing"
)

// Review 2.19: a search term went into a LIKE pattern as typed, and Export
// carried its own copy of the clause builder.
//
// The term is a bound parameter, so this was never injection. It was a wrong
// answer: "50%" matched every row because % is the wildcard, "a_b" matched
// "axb" because _ is any character, and a pattern of nothing but wildcards is
// work the database does for no result.
//
// Checked against the template text rather than by running the generated code,
// which is what every other test of this layer does: the thing that can go
// wrong is the text, and there is no Go here to call.
func TestTheSearchClauseEscapesLikeWildcards(t *testing.T) {
	src := strings.ReplaceAll(apiPaginateGo(), "{{MODULE}}", "example.com/app")

	for _, want := range []string{
		// Exported, so Export can call it instead of copying it.
		"func BuildSearchClause(cols []string, term string) (string, []any) {",
		"func EscapeLike(term string) string {",
		// The pattern declares its escape character, or escaping it does
		// nothing on SQLite, which has no default.
		`clause += "LOWER(" + col + ") LIKE LOWER(?) ESCAPE '" + LikeEscapeChar + "'"`,
		// The three characters that have to be escaped, including the escape
		// character itself.
		"if r == '%' || r == '_' || r == '!' {",
		// Applied to the term before the wildcards are added around it, or the
		// wildcards this clause means would be escaped too.
		`wild := "%" + EscapeLike(term) + "%"`,
	} {
		if !strings.Contains(src, want) {
			t.Errorf("internal/paginate is missing %q", want)
		}
	}

	// Not a backslash: Postgres and MySQL disagree about what '\' means inside
	// a string literal, and one of them does not parse.
	if strings.Contains(src, `ESCAPE '\\'`) || strings.Contains(src, `LikeEscapeChar = "\\"`) {
		t.Error("the escape character is a backslash, which is a syntax error on MySQL")
	}
	if !strings.Contains(src, `const LikeEscapeChar = "!"`) {
		t.Error("the escape character is not the one the comment and the test agree on")
	}

	// The old name still compiles, because a project may call it.
	if !strings.Contains(src, "func buildSearchClause(cols []string, term string) (string, []any) {") {
		t.Error("the unexported name is gone, so a project that calls it stops compiling")
	}
}
