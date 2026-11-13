package scaffold

import (
	"strings"
)

// RepairGeneratedServiceBatchTx teaches an existing generated service to join
// a batch's transaction.
//
// POST /api/v1/batch replays its requests through the router with the
// transaction on the context, and every read and write a service does goes
// through its db(ctx) helper. A service that does not look there writes
// outside the batch: its rows survive a batch that was rolled back, which is
// the one thing a batch exists to prevent.
//
// Its own repair rather than part of the row-rules one, because the two
// shipped in different releases and a project that took the first already
// looks repaired to it.
//
// Idempotent: a service that already reads dbtx is returned untouched.
func RepairGeneratedServiceBatchTx(src, module string) (string, bool) {
	pascal := serviceTypeName(src)
	if pascal == "" {
		return src, false
	}

	old := "func (s *" + pascal + "Service) db(ctx context.Context) *gorm.DB {\n" +
		"\treturn s.DB.WithContext(ctx)\n" +
		"}"
	if !strings.Contains(src, old) {
		// Already repaired, already generated against it, or a shape this
		// repair does not recognise. All three are "leave it alone".
		return src, false
	}

	out := src

	importLine := "\t\"" + module + "/internal/dbtx\"\n"
	if !strings.Contains(out, importLine) {
		anchor := "\t\"" + module + "/internal/authz\"\n"
		if !strings.Contains(out, anchor) {
			return src, false
		}
		out = strings.Replace(out, anchor, anchor+importLine, 1)
	}

	replacement := "func (s *" + pascal + "Service) db(ctx context.Context) *gorm.DB {\n" +
		"\t// A batch puts its transaction on the context, and everything one of its\n" +
		"\t// sub-requests does has to happen inside it. Every read and write in this\n" +
		"\t// service goes through here, which is why this is the one place that has\n" +
		"\t// to know. See internal/dbtx.\n" +
		"\tif tx := dbtx.TxFrom(ctx); tx != nil {\n" +
		"\t\treturn tx.WithContext(ctx)\n" +
		"\t}\n" +
		"\treturn s.DB.WithContext(ctx)\n" +
		"}"
	out = strings.Replace(out, old, replacement, 1)
	return out, true
}
