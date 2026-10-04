package services

import (
	"time"

	"gorm.io/gorm"
)

// dayExpr is a SQL expression for column's calendar day in UTC, as YYYY-MM-DD,
// in the dialect db speaks. Grouping by it counts rows per day in the database;
// the stats and chart services used to load every row of the last 30 days into
// the API and count them there.
func dayExpr(db *gorm.DB, column string) string {
	switch db.Dialector.Name() {
	case "postgres":
		return "to_char(" + column + " AT TIME ZONE 'UTC', 'YYYY-MM-DD')"
	case "mysql":
		return "DATE_FORMAT(" + column + ", '%Y-%m-%d')"
	default: // sqlite
		return "strftime('%Y-%m-%d', " + column + ")"
	}
}

// lastThirtyDays is the first moment of the 30 UTC days that end today.
func lastThirtyDays() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -29)
}

// asStored is t written the way GORM writes created_at: the same instant, in
// local time. SQLite compares times as text, so a UTC cutoff against local
// timestamps was off by the machine's UTC offset. Postgres compares instants
// and is unaffected.
func asStored(t time.Time) time.Time {
	return t.In(time.Local)
}
