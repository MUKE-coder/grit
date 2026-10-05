"use client";

import { useState } from "react";
import { ChevronLeft, ChevronRight } from "lucide-react";

import type { FileRef } from "@repo/shared/types";

// The product's pictures.
//
// A big one and a row of thumbnails, with arrows for anyone who would rather
// not aim at a 72px target. The thumbnails are real buttons in a list, so the
// whole gallery is reachable with Tab and arrow keys announce "image 2 of 4"
// rather than nothing.
export function Gallery({
  images,
  title,
}: {
  images: FileRef[];
  title: string;
}) {
  const [index, setIndex] = useState(0);

  if (images.length === 0) {
    return (
      <div className="flex aspect-[9/11] items-center justify-center rounded-xl border border-dashed border-border bg-bg-secondary text-sm text-text-muted">
        No pictures yet
      </div>
    );
  }

  const current = images[Math.min(index, images.length - 1)];
  const step = (by: number) =>
    setIndex((i) => (i + by + images.length) % images.length);

  return (
    <div>
      <div className="group relative overflow-hidden rounded-xl border border-border bg-bg-secondary">
        <img
          src={current.url}
          alt={`${title}, picture ${index + 1} of ${images.length}`}
          width={900}
          height={1100}
          // The first picture is the largest thing on the page and is on screen
          // immediately, so it is fetched eagerly and at high priority.
          loading="eager"
          fetchPriority="high"
          className="aspect-[9/11] w-full object-cover"
        />

        {images.length > 1 && (
          <>
            <button
              type="button"
              onClick={() => step(-1)}
              aria-label="Previous picture"
              className="absolute left-3 top-1/2 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full border border-border bg-background/80 text-foreground opacity-0 backdrop-blur transition-opacity hover:bg-background focus-visible:opacity-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent group-hover:opacity-100"
            >
              <ChevronLeft className="h-5 w-5" aria-hidden="true" />
            </button>
            <button
              type="button"
              onClick={() => step(1)}
              aria-label="Next picture"
              className="absolute right-3 top-1/2 flex h-10 w-10 -translate-y-1/2 items-center justify-center rounded-full border border-border bg-background/80 text-foreground opacity-0 backdrop-blur transition-opacity hover:bg-background focus-visible:opacity-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent group-hover:opacity-100"
            >
              <ChevronRight className="h-5 w-5" aria-hidden="true" />
            </button>
          </>
        )}
      </div>

      {images.length > 1 && (
        <ul className="mt-3 flex gap-3">
          {images.map((image, i) => (
            <li key={image.url + i}>
              <button
                type="button"
                onClick={() => setIndex(i)}
                aria-label={`Show picture ${i + 1} of ${images.length}`}
                aria-current={i === index ? "true" : undefined}
                className={[
                  "overflow-hidden rounded-lg border transition-colors",
                  "focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent",
                  i === index
                    ? "border-accent"
                    : "border-border hover:border-accent/50",
                ].join(" ")}
              >
                <img
                  src={image.url}
                  alt=""
                  width={72}
                  height={88}
                  loading="lazy"
                  className="h-[88px] w-[72px] object-cover"
                />
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
