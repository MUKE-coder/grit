package database

import "github.com/brianvoe/gofakeit/v7"

// pickID returns a random id from ids, or "" if empty. Faker seeders use it to
// link a belongs_to field to a real existing parent (load the parent ids with
// db.Model(&models.Parent{}).Pluck("id", &ids) first).
func pickID(ids []string) string {
	if len(ids) == 0 {
		return ""
	}
	return ids[gofakeit.Number(0, len(ids)-1)]
}

// firstID returns the first id, or "" if empty — the static example seeders
// use it to link to a parent deterministically.
func firstID(ids []string) string {
	if len(ids) == 0 {
		return ""
	}
	return ids[0]
}
