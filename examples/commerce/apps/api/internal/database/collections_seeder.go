package database

import (
	"log"

	"commerce/apps/api/internal/files"
	"commerce/apps/api/internal/models"
	"gorm.io/gorm"
)

// SeedCollections inserts the shop's collections.
//
// Hand-edited from what `grit generate resource Collection --seed` wrote, which
// was one row called "Sample Title". The generated seeder exists to be edited:
// it knows the shape of the model and nothing about the shop.
//
// Three is a deliberate number. One collection is indistinguishable from no
// collections on a storefront, because nothing has to be filtered; three is
// enough for the navigation, the collection pages and the related-items strip
// to all have something to show.
func SeedCollections(db *gorm.DB) error {
	var count int64
	db.Model(&models.Collection{}).Count(&count)
	if count > 0 {
		log.Println("Collections already seeded, skipping...")
		return nil
	}

	records := []models.Collection{
		{
			Title:       "Apparel",
			Handle:      "apparel",
			Description: "Tees, shirts and outerwear, cut from heavier cloth than they need to be.",
			Image:       shopImage("collection-apparel", "apparel.jpg"),
		},
		{
			Title:       "Accessories",
			Handle:      "accessories",
			Description: "Bags, caps and small leather goods: the things you lose and replace.",
			Image:       shopImage("collection-accessories", "accessories.jpg"),
		},
		{
			Title:       "Home",
			Handle:      "home",
			Description: "Mugs, candles and objects for a flat you did not choose.",
			Image:       shopImage("collection-home", "home.jpg"),
		},
	}

	for i := range records {
		if err := db.Create(&records[i]).Error; err != nil {
			log.Printf("Warning: failed to seed collection %q: %v", records[i].Handle, err)
		}
	}
	log.Printf("Seeded %d collection(s)", len(records))
	return nil
}

// shopImage builds a stable placeholder.
//
// Seeded by a string we choose rather than a random one, so the same product
// keeps the same picture across a reseed, and a screenshot taken last week
// still matches the shop. A real shop uploads these through the admin and gets
// a FileRef pointing at its own storage; the shape is identical, which is why
// the storefront needs no change when it does.
func shopImage(seed, name string) *files.FileRef {
	return &files.FileRef{
		URL:  "https://picsum.photos/seed/" + seed + "/900/1100",
		Name: name,
		MIME: "image/jpeg",
	}
}
