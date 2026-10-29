package generate

import (
	"strings"
	"testing"
)

// grit#132: a sequence with Reset "never" issued one number and then failed
// forever.
//
// Its bucket is "", the zero value of the string half of the counter's
// primary key, and GORM's Save treats a zero primary key as a new row: the
// second call INSERTed a duplicate instead of updating. Monthly and yearly
// counters have a non-empty bucket, which is why only "never" broke.
func TestSequenceUpdatesRatherThanSaving(t *testing.T) {
	src := sequencePackageGo()

	if strings.Contains(src, "tx.Save(&c)") {
		t.Error("the counter is still written with Save, which inserts a " +
			"duplicate row when the bucket is empty, so every number after " +
			"the first fails on the unique constraint (grit#132)")
	}
	if !strings.Contains(src, `Where("name = ? AND bucket = ?", cfg.Name, bucket).`) {
		t.Error("the counter is no longer updated by its two key columns; " +
			"check the fix for grit#132 still holds")
	}
}
