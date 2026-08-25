import { describe, it, expect, vi } from "vitest";
import { createUploader } from "./uploader";
import type { MediaProfile, OptimizedImage, UploadTransport } from "./types";

const profile: MediaProfile = {
  name: "default",
  max_width: 1600,
  max_height: 1600,
  crop: false,
  quality: 0.82,
  format: "auto",
  max_pixels: 50000000,
  renditions: { thumb: { width: 400, height: 400, crop: true } },
};

function optimizedTo(primary: number, thumb: number): OptimizedImage {
  return {
    primary: {
      name: "primary",
      blob: new Blob([new Uint8Array(primary)], { type: "image/webp" }),
      width: 1600,
      height: 1200,
      mime: "image/webp",
      ext: ".webp",
    },
    extra: [
      {
        name: "thumb",
        blob: new Blob([new Uint8Array(thumb)], { type: "image/webp" }),
        width: 400,
        height: 400,
        mime: "image/webp",
        ext: ".webp",
      },
    ],
    originalWidth: 4000,
    originalHeight: 3000,
    originalSize: 6000000,
  };
}

interface Recorder extends UploadTransport {
  puts: Array<{ url: string; size: number }>;
  posts: Array<{ path: string; body: Record<string, unknown> }>;
}

function recorder(): Recorder {
  const puts: Array<{ url: string; size: number }> = [];
  const posts: Array<{ path: string; body: Record<string, unknown> }> = [];
  return {
    puts,
    posts,
    get: async () => ({ data: { profiles: [profile] } }) as never,
    post: async (path: string, body: Record<string, unknown>) => {
      posts.push({ path, body });
      if (path === "/uploads/presign") {
        return {
          data: {
            presigned_url: "https://s3.example/" + String(body.filename) + "?sig=x",
            key: "uploads/" + String(body.filename),
            public_url: "https://cdn.example/uploads/" + String(body.filename),
          },
        } as never;
      }
      return {
        data: {
          url: "https://cdn.example/x",
          key: body.key,
          name: body.filename,
          mime: body.content_type,
          size: body.size,
        },
      } as never;
    },
    put: async (url: string, blob: Blob) => {
      puts.push({ url, size: blob.size });
    },
  };
}

describe("uploader", () => {
  it("sends every rendition straight to storage, not to the API", async () => {
    const t = recorder();
    const up = createUploader({
      transport: t,
      optimize: async () => optimizedTo(40000, 9000),
      profiles: [profile],
    });

    const ref = await up.upload(
      new Blob([new Uint8Array(10)], { type: "image/jpeg" }),
      "photo.jpg",
    );

    expect(t.puts).toHaveLength(2);
    expect(t.puts.every((p) => p.url.startsWith("https://s3.example/"))).toBe(true);
    expect(ref.renditions?.thumb).toBeDefined();
    expect(ref.thumbnail_url).toContain("thumb");
    expect(ref.optimised).toBe(true);
  });

  it("presigns with the optimised byte count, not the original", async () => {
    const t = recorder();
    const up = createUploader({
      transport: t,
      optimize: async () => optimizedTo(40000, 9000),
      profiles: [profile],
    });

    await up.upload(
      new Blob([new Uint8Array(6000000)], { type: "image/jpeg" }),
      "photo.jpg",
    );

    // The server signs this number into the URL, so it has to be exactly what
    // gets sent. Reporting the original 6 MB would make S3 reject every PUT.
    const presigns = t.posts.filter((p) => p.path === "/uploads/presign");
    expect(presigns[0].body.file_size).toBe(40000);
    expect(presigns[1].body.file_size).toBe(9000);
  });

  it("uploads a non-image untouched", async () => {
    const t = recorder();
    const optimize = vi.fn();
    const up = createUploader({
      transport: t,
      optimize: optimize as never,
      profiles: [profile],
    });

    const ref = await up.upload(
      new Blob([new Uint8Array(100)], { type: "application/pdf" }),
      "doc.pdf",
    );

    expect(optimize).not.toHaveBeenCalled();
    expect(t.puts).toHaveLength(1);
    expect(ref.optimised).toBe(false);
  });

  it("leaves an animated GIF alone", async () => {
    const t = recorder();
    const optimize = vi.fn();
    const up = createUploader({
      transport: t,
      optimize: optimize as never,
      profiles: [profile],
    });

    await up.upload(new Blob([new Uint8Array(100)], { type: "image/gif" }), "loop.gif");

    // Drawing a GIF to a canvas keeps the first frame only, so optimising one
    // would silently throw the animation away.
    expect(optimize).not.toHaveBeenCalled();
  });

  it("refuses an oversized image instead of uploading it", async () => {
    const t = recorder();
    const up = createUploader({
      transport: t,
      optimize: async () => {
        throw new Error("That image is 12000x12000, larger than this app accepts (50 megapixels).");
      },
      profiles: [profile],
    });

    await expect(
      up.upload(new Blob([new Uint8Array(10)], { type: "image/png" }), "bomb.png"),
    ).rejects.toThrow("megapixels");
    expect(t.puts).toHaveLength(0);
  });

  it("uploads the original when the browser cannot decode it", async () => {
    const t = recorder();
    const up = createUploader({
      transport: t,
      optimize: async () => {
        throw new Error("could not decode image");
      },
      profiles: [profile],
    });

    const ref = await up.upload(
      new Blob([new Uint8Array(500)], { type: "image/jpeg" }),
      "odd.jpg",
    );

    // Losing somebody's file because a codec choked is worse than storing it
    // exactly as it arrived.
    expect(t.puts).toHaveLength(1);
    expect(ref.optimised).toBe(false);
  });
});
