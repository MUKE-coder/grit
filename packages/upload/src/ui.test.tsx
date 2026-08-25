import { describe, it, expect, vi, beforeAll, afterEach } from "vitest";
import { render, screen, waitFor, fireEvent, cleanup } from "@testing-library/react";
import { Dropzone } from "./ui";
import type { FileRef } from "./types";

/**
 * jsdom implements neither of these, and both are load-bearing: the previews
 * are object URLs, and the revoke is what stops every full-size original being
 * pinned in memory for the life of the page.
 */
beforeAll(() => {
  let n = 0;
  URL.createObjectURL = vi.fn(() => `blob:mock/${++n}`);
  URL.revokeObjectURL = vi.fn();
});

/**
 * Explicit, because Testing Library only registers its own afterEach when
 * vitest runs with globals enabled. Without this each render stays in the
 * document and the queries start matching several elements at once, which
 * fails as "found multiple" and looks like a component bug rather than a
 * harness one.
 */
afterEach(cleanup);

function uploaderThatShrinks(to: number) {
  return {
    upload: vi.fn(async (file: Blob, filename: string, options?: any): Promise<FileRef> => {
      options?.onOptimized?.({
        primary: { blob: new Blob([new Uint8Array(to)]), mime: "image/webp" },
      });
      options?.onProgress?.(1);
      return {
        url: `https://cdn.example/${filename}`,
        key: `uploads/${filename}`,
        name: filename,
        mime: "image/webp",
        size: to,
      };
    }),
    getProfiles: vi.fn(),
    profileFor: vi.fn(),
  } as never;
}

function pick(input: HTMLElement, file: File) {
  Object.defineProperty(input, "files", { value: [file], configurable: true });
  fireEvent.change(input);
}

const bigPhoto = () =>
  new File([new Uint8Array(5_200_000)], "photo.jpg", { type: "image/jpeg" });

describe("<Dropzone />", () => {
  it("uses a real file input behind a label, not a div with a click handler", () => {
    const { container } = render(<Dropzone uploader={uploaderThatShrinks(41_000)} />);
    const input = container.querySelector('input[type="file"]');

    // This is what gives keyboard activation, focus and the mobile picker for
    // free. A div would have to reimplement all of it.
    expect(input).toBeTruthy();
    expect(input?.id).toBeTruthy();
    expect(container.querySelector(`label[for="${input?.id}"]`)).toBeTruthy();
  });

  it("shows what the optimisation saved", async () => {
    const { container } = render(<Dropzone uploader={uploaderThatShrinks(41_000)} />);
    pick(container.querySelector('input[type="file"]')!, bigPhoto());

    await screen.findByText("photo.jpg");
    // The whole reason the component exists: without this the only visible
    // effect of the optimisation is that the bar finishes sooner.
    await waitFor(() => {
      expect(screen.getByText(/5\.0 MB -> 40 KB/)).toBeTruthy();
    });
  });

  it("hands back the refs of what uploaded", async () => {
    const onChange = vi.fn();
    const { container } = render(
      <Dropzone uploader={uploaderThatShrinks(41_000)} onChange={onChange} />,
    );
    pick(container.querySelector('input[type="file"]')!, bigPhoto());

    await waitFor(() => expect(onChange).toHaveBeenCalledTimes(1));
    expect(onChange.mock.calls[0][0][0].name).toBe("photo.jpg");
  });

  it("keeps a failed file in the list, with the reason", async () => {
    const failing = {
      upload: vi.fn(async () => {
        throw new Error("Storage rejected the upload");
      }),
    } as never;
    const { container } = render(<Dropzone uploader={failing} />);
    pick(container.querySelector('input[type="file"]')!, bigPhoto());

    // A row that vanishes on failure reads as success.
    await screen.findByText("Storage rejected the upload");
    expect(screen.getByText("photo.jpg")).toBeTruthy();
  });

  it("passes the profile through to the uploader", async () => {
    const uploader = uploaderThatShrinks(41_000);
    const { container } = render(<Dropzone uploader={uploader} profile="product-image" />);
    pick(container.querySelector('input[type="file"]')!, bigPhoto());

    await waitFor(() => {
      expect((uploader as any).upload).toHaveBeenCalledWith(
        expect.anything(),
        "photo.jpg",
        expect.objectContaining({ profile: "product-image" }),
      );
    });
  });

  it("stops accepting once maxFiles is reached", async () => {
    const { container } = render(<Dropzone uploader={uploaderThatShrinks(41_000)} maxFiles={1} />);
    const input = container.querySelector('input[type="file"]') as HTMLInputElement;
    pick(input, bigPhoto());

    await screen.findByText("photo.jpg");
    await waitFor(() => {
      expect(screen.getByText("Maximum of 1 files reached")).toBeTruthy();
      expect(input.disabled).toBe(true);
    });
  });

  it("removes a file, and releases its preview", async () => {
    const { container } = render(<Dropzone uploader={uploaderThatShrinks(41_000)} />);
    pick(container.querySelector('input[type="file"]')!, bigPhoto());

    await screen.findByText("photo.jpg");
    fireEvent.click(screen.getByRole("button", { name: "Remove photo.jpg" }));

    await waitFor(() => expect(screen.queryByText("photo.jpg")).toBeNull());
    expect(URL.revokeObjectURL).toHaveBeenCalled();
  });

  it("exposes progress to assistive tech while it runs", async () => {
    // Held open so the progressbar is still on screen when it is asserted.
    let release: (v: FileRef) => void = () => {};
    const slow = {
      upload: vi.fn((_f: Blob, name: string, options?: any) => {
        options?.onProgress?.(0.4);
        return new Promise<FileRef>((resolve) => {
          release = resolve;
        });
      }),
    } as never;

    const { container } = render(<Dropzone uploader={slow} />);
    pick(container.querySelector('input[type="file"]')!, bigPhoto());

    const bar = await screen.findByRole("progressbar");
    expect(bar.getAttribute("aria-valuenow")).toBe("40");
    expect(bar.getAttribute("aria-label")).toBe("Uploading photo.jpg");

    release({ url: "u", key: "k", name: "photo.jpg", mime: "image/webp", size: 1 });
  });
});
