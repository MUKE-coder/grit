package media

// Image optimisation profiles for this project.
//
// This file is yours. It is written once and never regenerated, so anything
// you add here survives grit generate and grit upgrade.
//
// You do not have to define anything. Every image field without a profile gets
// DefaultProfile(): fit inside 1600x1600, quality 0.82, format chosen per
// image, a 400x400 thumbnail, the original kept privately for reprocessing,
// and EXIF stripped. On a 6 MB phone photo that is about 150 KB out.
//
// Define a profile when a particular field wants something different, then
// name it from the upload: POST /api/v1/uploads?profile=product-image
//
// Fields you leave zero keep the default, so a profile that only wants a
// different size says only that.

func init() {
	// A product photograph on a storefront. Smaller than the default, because
	// a catalogue page shows a dozen of them at once.
	Define("product-image", Profile{
		Max:     Fit(1000, 1000),
		Quality: 0.8,
		Renditions: map[string]Size{
			"thumb": Fill(300, 300),
			"card":  Fit(600, 600),
		},
	})

	// An avatar is always square and always small, so it crops rather than
	// fits: a portrait photo shown in a round frame should not letterbox.
	Define("avatar", Profile{
		Max:        Fill(400, 400),
		Quality:    0.85,
		Renditions: map[string]Size{"thumb": Fill(80, 80)},
		// Nothing to reprocess later: the source is a user's selfie, not
		// artwork you will re-crop, and keeping every original costs storage
		// for a file nobody will ask for again.
		DiscardOriginal: true,
	})

	// A hero or cover image, where width matters more than total pixels.
	Define("cover", Profile{
		Max:        Fit(2400, 1200),
		Quality:    0.82,
		Renditions: map[string]Size{"thumb": Fill(600, 300)},
	})
}
