package database

import (
	"log"

	"commerce/apps/api/internal/models"
	"gorm.io/gorm"
)

// SeedPages inserts the shop's static pages.
//
// Hand-edited from the generated one-row seeder. These are the pages every shop
// has and nobody enjoys writing: the footer links to them by handle and the
// storefront renders them from one dynamic route, so adding a fourth needs no
// code at all.
func SeedPages(db *gorm.DB) error {
	var count int64
	db.Model(&models.Page{}).Count(&count)
	if count > 0 {
		log.Println("Pages already seeded, skipping...")
		return nil
	}

	records := []models.Page{
		{
			Title:  "About",
			Handle: "about",
			Body: "<p>We make a small number of things and then keep making them. There " +
				"is no spring collection, because the things we sold last spring are " +
				"still good.</p>" +
				"<p>Everything here is cut from cloth we have used before, by people we " +
				"have worked with before, which is less romantic than it sounds and the " +
				"only reliable way to know how something will wash.</p>",
		},
		{
			Title:  "Shipping & Returns",
			Handle: "shipping-returns",
			Body: "<p>Orders placed before noon go out the same day. Everything else " +
				"goes out the next working day, and you get a tracking number when it " +
				"actually leaves rather than when the label is printed.</p>" +
				"<p>Thirty days to return anything unworn, and we pay the postage both " +
				"ways. If something fails in a way cloth should not fail, tell us " +
				"whenever it happens and we will sort it out.</p>",
		},
		{
			Title:  "Terms & Conditions",
			Handle: "terms-conditions",
			Body: "<p>A demonstration shop, so there is nothing here to agree to. In a " +
				"real one this page is written by somebody who is paid to write it, and " +
				"it lives in exactly this table.</p>" +
				"<p>The point of this page is that it needed no code: it is a row in " +
				"<code>pages</code>, rendered by the same route as the other two.</p>",
		},
	}

	for i := range records {
		if err := db.Create(&records[i]).Error; err != nil {
			log.Printf("Warning: failed to seed page %q: %v", records[i].Handle, err)
		}
	}
	log.Printf("Seeded %d page(s)", len(records))
	return nil
}
