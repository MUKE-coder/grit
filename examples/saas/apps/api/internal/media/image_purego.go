//go:build !vips

package media

import "image"

// encodeImage writes a chain's pixels with the encoders transform.go uses.
// WebP here is lossless: quality has no effect on it.
func encodeImage(img image.Image, f Format, quality int) ([]byte, error) {
	return encode(img, f, float64(quality)/100)
}

// autoFormat is what an Image with no To call encodes as: lossless WebP when
// it has transparency, JPEG otherwise, the same choice a profile's Auto makes.
func autoFormat(img image.Image) Format {
	return resolveFormat(Auto, img)
}
