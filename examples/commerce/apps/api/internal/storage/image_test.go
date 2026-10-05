package storage

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"runtime"
	"testing"

	"commerce/apps/api/internal/media"
)

// pngClaiming is a PNG whose header claims w x h RGBA pixels, with one short
// row of data behind it. Decoding it allocates for every claimed pixel before
// discovering the data is not there.
func pngClaiming(w, h uint32) []byte {
	var b bytes.Buffer
	b.WriteString("\x89PNG\r\n\x1a\n")
	chunk := func(kind string, data []byte) {
		_ = binary.Write(&b, binary.BigEndian, uint32(len(data)))
		b.WriteString(kind)
		b.Write(data)
		_ = binary.Write(&b, binary.BigEndian, crc32.ChecksumIEEE(append([]byte(kind), data...)))
	}
	ihdr := make([]byte, 13)
	binary.BigEndian.PutUint32(ihdr[0:], w)
	binary.BigEndian.PutUint32(ihdr[4:], h)
	ihdr[8], ihdr[9] = 8, 6
	chunk("IHDR", ihdr)
	var z bytes.Buffer
	zw := zlib.NewWriter(&z)
	_, _ = zw.Write(make([]byte, 64))
	_ = zw.Close()
	chunk("IDAT", z.Bytes())
	chunk("IEND", nil)
	return b.Bytes()
}

// A few kilobytes claiming 30000x30000 would allocate 3.6 GB in image.Decode.
// The thumbnail worker ran that in the API process and retried it five times.
func TestAPixelBombIsRefusedBeforeDecoding(t *testing.T) {
	bomb := pngClaiming(30000, 30000)
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	_, err := GenerateThumbnail(bytes.NewReader(bomb), "image/png")
	runtime.ReadMemStats(&after)
	if !errors.Is(err, ErrImageTooLarge) {
		t.Fatalf("want ErrImageTooLarge, got %v", err)
	}
	if grew := after.TotalAlloc - before.TotalAlloc; grew > 64<<20 {
		t.Errorf("refusing the image allocated %d MB", grew>>20)
	}
	if _, err := ProcessImage(bytes.NewReader(bomb), "image/png"); !errors.Is(err, ErrImageTooLarge) {
		t.Errorf("ProcessImage: want ErrImageTooLarge, got %v", err)
	}
}

func TestNotAnImageIsUnreadable(t *testing.T) {
	if _, err := GenerateThumbnail(bytes.NewReader([]byte("not an image")), "image/png"); !errors.Is(err, ErrUnreadableImage) {
		t.Errorf("want ErrUnreadableImage, got %v", err)
	}
}

func TestAnOrdinaryImageStillMakesAThumbnail(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 640, 480))
	for x := 0; x < 640; x++ {
		img.Set(x, x%480, color.RGBA{R: 255, A: 255})
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatal(err)
	}
	thumb, err := GenerateThumbnail(bytes.NewReader(buf.Bytes()), "image/png")
	if err != nil {
		t.Fatal(err)
	}
	got, _, err := image.DecodeConfig(bytes.NewReader(thumb))
	if err != nil || got.Width != ThumbnailSize || got.Height != ThumbnailSize {
		t.Errorf("thumbnail is %dx%d (%v), want %dx%d", got.Width, got.Height, err, ThumbnailSize, ThumbnailSize)
	}
	if thumb := media.DefaultProfile().Renditions["thumb"]; thumb.Width != ThumbnailSize || thumb.Height != ThumbnailSize {
		t.Errorf("ThumbnailSize is %d, the media pipeline's thumb is %dx%d", ThumbnailSize, thumb.Width, thumb.Height)
	}
}

// jpegWithOrientation encodes img as a JPEG carrying an EXIF Orientation tag,
// as a phone writes a portrait photo: the pixels stay landscape and the tag
// says how to turn them.
func jpegWithOrientation(t *testing.T, img image.Image, orientation byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95}); err != nil {
		t.Fatal(err)
	}
	src := buf.Bytes()
	// A big-endian TIFF header and one IFD entry: tag 0x0112, SHORT, count 1.
	tiff := []byte{'M', 'M', 0, 42, 0, 0, 0, 8, 0, 1, 0x01, 0x12, 0, 3, 0, 0, 0, 1, 0, orientation, 0, 0, 0, 0, 0, 0}
	payload := append([]byte("Exif\x00\x00"), tiff...)
	out := append([]byte{}, src[:2]...)
	out = append(out, 0xFF, 0xE1, byte((len(payload)+2)>>8), byte(len(payload)+2))
	out = append(out, payload...)
	return append(out, src[2:]...)
}

func isRed(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return r > 0xa000 && g < 0x6000 && b < 0x6000
}

func isBlue(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return b > 0xa000 && r < 0x6000 && g < 0x6000
}

// A photo whose EXIF says "turn me 90 degrees clockwise" gets a thumbnail the
// right way up. The helpers used to decode without orientation, so a presigned
// upload or a skipped one got a sideways thumbnail.
func TestThumbnailsFollowEXIFOrientation(t *testing.T) {
	// Landscape pixels: red on the left, blue on the right. Turned clockwise,
	// red is on top.
	wide := image.NewRGBA(image.Rect(0, 0, 160, 80))
	for y := 0; y < 80; y++ {
		for x := 0; x < 160; x++ {
			c := color.RGBA{R: 255, A: 255}
			if x >= 80 {
				c = color.RGBA{B: 255, A: 255}
			}
			wide.Set(x, y, c)
		}
	}
	photo := jpegWithOrientation(t, wide, 6)

	thumb, err := GenerateThumbnail(bytes.NewReader(photo), "image/jpeg")
	if err != nil {
		t.Fatal(err)
	}
	img, _, err := image.Decode(bytes.NewReader(thumb))
	if err != nil {
		t.Fatal(err)
	}
	if !isRed(img.At(100, 60)) || !isRed(img.At(300, 60)) || !isBlue(img.At(100, 340)) || !isBlue(img.At(300, 340)) {
		t.Errorf("the thumbnail is not the right way up: top %v %v, bottom %v %v",
			img.At(100, 60), img.At(300, 60), img.At(100, 340), img.At(300, 340))
	}

	processed, err := ProcessImage(bytes.NewReader(photo), "image/jpeg")
	if err != nil {
		t.Fatal(err)
	}
	if cfg, _, err := image.DecodeConfig(bytes.NewReader(processed)); err != nil || cfg.Width != 80 || cfg.Height != 160 {
		t.Errorf("ProcessImage = %dx%d (%v), want the photo turned upright at 80x160", cfg.Width, cfg.Height, err)
	}
}
