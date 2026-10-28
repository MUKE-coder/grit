#!/usr/bin/env python3
"""Copy the blueprint screenshots into the website, cropped and resized.

Every blueprint repository ships its own mockups under docs/assets/mockups:
real screenshots of the running app, made for its README. The website wants
the same pictures, smaller, and without the promo headline that the phone ones
carry baked in at the top.

    python scripts/blueprint-shots.py

Reads D:/PROJECTS/grit-blueprints/<name>/docs/assets/mockups and writes
public/blueprints/<name>/NN-<label>.webp. The source repositories are not part
of this one, so this cannot run in CI: the OUTPUT is committed, and this script
is how it is regenerated when a blueprint gets new screenshots. With the
repositories missing it says so and exits 0 rather than failing a build that
has nothing to do with it.

The phone mockups are 860x1864 with a headline above the device. Cropping the
top 27% away leaves the phone, which is the part worth showing at card size;
the headline is unreadable there and reads as noise.
"""
import os
import sys

try:
    from PIL import Image
except ImportError:
    print("Pillow is not installed: pip install Pillow")
    sys.exit(0)

HERE = os.path.dirname(os.path.abspath(__file__))
DOCS = os.path.dirname(HERE)
SOURCE = os.environ.get("BLUEPRINT_REPOS", "D:/PROJECTS/grit-blueprints")
OUT = os.path.join(DOCS, "public", "blueprints")

# blueprint -> (kind, [(source file, output name)]). Three each: the card
# shows all three and the detail page leads with the first.
SHOTS = {
    "whatsapp": (
        "phone",
        [
            ("01-inbox.webp", "01-inbox.webp"),
            ("02-chat.webp", "02-chat.webp"),
            ("03-group.webp", "03-group.webp"),
        ],
    ),
    "instagram": (
        "phone",
        [
            ("01-feed.webp", "01-feed.webp"),
            ("02-carousel.webp", "02-carousel.webp"),
            ("03-profile.webp", "03-profile.webp"),
        ],
    ),
    "storefront": (
        "screen",
        [
            ("03-product.webp", "01-product.webp"),
            ("02-shop.webp", "02-shop.webp"),
            ("06-seller.webp", "03-seller.webp"),
        ],
    ),
}

# Width to write, by kind. A card renders these at about 380 CSS pixels, so
# these are comfortably past 2x without being a megabyte each.
WIDTH = {"phone": 480, "screen": 1280}

# How much of a phone mockup is headline rather than phone, and how much of
# each side is the promo background the device sits on. Cropping both leaves
# the screen itself, which is what fans nicely three to a card.
PHONE_TOP_CROP = 0.27
PHONE_SIDE_CROP = 0.135


def main():
    if not os.path.isdir(SOURCE):
        print("%s is not on this machine, so there is nothing to copy." % SOURCE)
        print("The committed images under public/blueprints are unchanged.")
        return 0

    written = 0
    for name, (kind, files) in SHOTS.items():
        src_dir = os.path.join(SOURCE, name, "docs", "assets", "mockups")
        if not os.path.isdir(src_dir):
            print("skipped %-12s no mockups at %s" % (name, src_dir))
            continue

        dest_dir = os.path.join(OUT, name)
        os.makedirs(dest_dir, exist_ok=True)

        for src_file, out_file in files:
            src = os.path.join(src_dir, src_file)
            if not os.path.exists(src):
                print("skipped %-12s %s is missing" % (name, src_file))
                continue

            im = Image.open(src).convert("RGB")
            if kind == "phone":
                top = int(im.height * PHONE_TOP_CROP)
                side = int(im.width * PHONE_SIDE_CROP)
                im = im.crop((side, top, im.width - side, im.height))

            width = WIDTH[kind]
            if im.width > width:
                height = round(im.height * width / im.width)
                im = im.resize((width, height), Image.LANCZOS)

            dest = os.path.join(dest_dir, out_file)
            im.save(dest, "WEBP", quality=82, method=6)
            written += 1
            print("%-12s %-20s %dx%d  %dKB"
                  % (name, out_file, im.width, im.height,
                     os.path.getsize(dest) // 1024))

    print("\nwrote %d images into public/blueprints" % written)
    return 0


if __name__ == "__main__":
    sys.exit(main())
