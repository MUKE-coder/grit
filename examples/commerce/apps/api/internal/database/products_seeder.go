package database

import (
	"fmt"
	"log"

	"commerce/apps/api/internal/files"
	"commerce/apps/api/internal/models"
	"commerce/apps/api/internal/money"
	"gorm.io/gorm"
)

// SeedProducts inserts the catalogue.
//
// Hand-edited from the one-row seeder `grit generate resource Product --seed`
// wrote. The prices are integers in minor units, because that is what
// money.Money holds: 2800 USD is $28.00. Nothing in the shop ever stores a
// price as a float, which is the whole reason the money type exists.
//
// The first six rows are the ones SeedProductVariants attaches a Colour x Size
// matrix to, so they are deliberately the apparel: a mug does not come in a
// medium. The rest stay plain, and the storefront has to render both.
func SeedProducts(db *gorm.DB) error {
	var count int64
	db.Model(&models.Product{}).Count(&count)
	if count > 0 {
		log.Println("Products already seeded, skipping...")
		return nil
	}

	// The collections, by handle, so a product names the one it belongs to
	// rather than taking whichever row happened to be inserted first.
	byHandle := map[string]string{}
	var collections []models.Collection
	if err := db.Find(&collections).Error; err != nil {
		return fmt.Errorf("loading collections to attach products to: %w", err)
	}
	for _, c := range collections {
		byHandle[c.Handle] = c.ID
	}
	if len(byHandle) == 0 {
		return fmt.Errorf("cannot seed products: no collections exist yet, so run grit seed after the collections seeder has run")
	}

	type entry struct {
		title, handle, collection, description string
		cents                                  int64
	}

	// Apparel first: these six get the variant matrix.
	catalogue := []entry{
		{
			title: "Heavyweight Cotton Tee", handle: "heavyweight-cotton-tee",
			collection: "apparel", cents: 2800,
			description: "<p>Eight ounces of combed cotton, boxy through the body, with a " +
				"collar that keeps its shape after the wash that ruins everything else. " +
				"Pre-shrunk, so the size you order is the size you keep.</p>",
		},
		{
			title: "Oxford Button-Down", handle: "oxford-button-down",
			collection: "apparel", cents: 7500,
			description: "<p>A proper oxford cloth, woven thick enough to hold a crease and " +
				"soft enough to wear on a Sunday. Unlined collar, single patch pocket, " +
				"and a back pleat so you can reach for something without the shirt arguing.</p>",
		},
		{
			title: "Merino Crew Knit", handle: "merino-crew-knit",
			collection: "apparel", cents: 11000,
			description: "<p>Fine-gauge merino, knitted in a mill that has been doing it " +
				"since before anyone asked it to be sustainable. Warm without bulk, and " +
				"it travels folded without resenting you for it.</p>",
		},
		{
			title: "Chore Jacket", handle: "chore-jacket",
			collection: "apparel", cents: 16500,
			description: "<p>Cotton canvas that starts stiff and ends up yours. Three " +
				"pockets, all of them big enough for a phone and one of them big enough " +
				"for a paperback. Unlined, so it works in three seasons out of four.</p>",
		},
		{
			title: "Fleece Half-Zip", handle: "fleece-half-zip",
			collection: "apparel", cents: 9800,
			description: "<p>Brushed back, high collar, and a zip that stops where you " +
				"want it to. The layer you put on in the house and then forget to take " +
				"off when you leave it.</p>",
		},
		{
			title: "Pique Polo", handle: "pique-polo",
			collection: "apparel", cents: 5400,
			description: "<p>Textured cotton pique with a ribbed collar that stands up on " +
				"its own. Two buttons, no logo, and a hem cut straight so it looks " +
				"deliberate untucked.</p>",
		},

		// Accessories and home: no options, which is most of a real shop.
		{
			title: "Canvas Tote", handle: "canvas-tote",
			collection: "accessories", cents: 4200,
			description: "<p>Twenty-ounce canvas, flat-bottomed so it stands up on a " +
				"kitchen floor, with webbing handles long enough to go over a shoulder " +
				"in a coat. One interior pocket for the keys you will still lose.</p>",
		},
		{
			title: "Six-Panel Cap", handle: "six-panel-cap",
			collection: "accessories", cents: 3200,
			description: "<p>Washed cotton twill with a soft, unstructured crown and a " +
				"brass slider at the back. Takes the shape of your head in about a week " +
				"and keeps it.</p>",
		},
		{
			title: "Leather Card Holder", handle: "leather-card-holder",
			collection: "accessories", cents: 5800,
			description: "<p>Vegetable-tanned leather, four slots, no stitching across the " +
				"spine so it folds flat in a front pocket. Arrives pale and ends up the " +
				"colour of a cricket ball.</p>",
		},
		{
			title: "Stoneware Mug", handle: "stoneware-mug",
			collection: "home", cents: 2400,
			description: "<p>Twelve ounces, thick-walled, with a handle sized for a whole " +
				"finger rather than the tip of one. Reactive glaze, so no two are the " +
				"same and none of them are a mistake.</p>",
		},
		{
			title: "Cedar Soy Candle", handle: "cedar-soy-candle",
			collection: "home", cents: 3800,
			description: "<p>Cedar, a little vetiver, and nothing that announces itself " +
				"from the next room. Forty hours in a glass you will keep for pencils " +
				"afterwards.</p>",
		},
		{
			title: "Linen Tea Towel Set", handle: "linen-tea-towel-set",
			collection: "home", cents: 2900,
			description: "<p>Two towels in washed European linen, which dries a glass " +
				"without leaving anything behind on it. Stiff for the first week, then " +
				"better than cotton for the next decade.</p>",
		},
	}

	records := make([]models.Product, 0, len(catalogue))
	for _, e := range catalogue {
		collectionID, ok := byHandle[e.collection]
		if !ok {
			log.Printf("Warning: no %q collection, so %q is being seeded without one", e.collection, e.handle)
		}
		records = append(records, models.Product{
			Title:       e.title,
			Handle:      e.handle,
			Description: e.description,
			Price:       money.New(e.cents, "USD"),
			// Four pictures each, seeded by handle so they are stable across a
			// reseed. FeaturedImage is the card; Images is the gallery, and the
			// detail page shows the featured one first followed by these.
			FeaturedImage: shopImage(e.handle, e.handle+".jpg"),
			Images: files.FileRefs{
				*shopImage(e.handle+"-2", e.handle+"-2.jpg"),
				*shopImage(e.handle+"-3", e.handle+"-3.jpg"),
				*shopImage(e.handle+"-4", e.handle+"-4.jpg"),
			},
			Available:    true,
			CollectionID: collectionID,
		})
	}

	for i := range records {
		if err := db.Create(&records[i]).Error; err != nil {
			log.Printf("Warning: failed to seed product %q: %v", records[i].Handle, err)
		}
	}
	log.Printf("Seeded %d product(s)", len(records))
	return nil
}
