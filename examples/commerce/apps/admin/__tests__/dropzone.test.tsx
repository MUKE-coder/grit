import { describe, it, expect, vi, beforeAll, beforeEach } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

// The upload goes to the API; here it answers the way /api/uploads does.
vi.mock("@/lib/api-client", () => ({
  uploadFile: vi.fn(async (file: File) => ({
    data: { url: "https://cdn.test/" + file.name, key: "uploads/" + file.name, name: file.name, size: file.size, mime: file.type },
  })),
}));

import {
  AvatarDropzone,
  BoxDropzone,
  CompactDropzone,
  Dropzone,
  InlineDropzone,
  MinimalDropzone,
  type UploadedFile,
} from "@/components/ui/dropzone";

beforeAll(() => {
  // jsdom has no object URLs, and a dropzone that does not upload makes one
  // per local preview.
  URL.createObjectURL = vi.fn(() => "blob:preview");
  URL.revokeObjectURL = vi.fn();
});

let onFilesChange: ReturnType<typeof vi.fn>;
beforeEach(() => {
  onFilesChange = vi.fn();
});

function picker(container: HTMLElement): HTMLInputElement {
  const input = container.querySelector('input[type="file"]');
  if (!(input instanceof HTMLInputElement)) throw new Error("no file input");
  return input;
}

const png = (name: string) => new File(["x"], name, { type: "image/png" });
const stored = (name: string, type = "image/png"): UploadedFile => ({
  url: "https://cdn.test/" + name,
  name,
  size: 2048,
  type,
});

describe("AvatarDropzone", () => {
  it("uploads one picture and shows it in the circle", async () => {
    const user = userEvent.setup();
    const { container } = render(
      <AvatarDropzone label="Photo" accept={{ "image/*": [".png"] }} onFilesChange={onFilesChange} />
    );
    expect(screen.getByText("Photo")).toBeInTheDocument();
    expect(screen.getByText("Upload")).toBeInTheDocument();

    await user.upload(picker(container), png("me.png"));

    await waitFor(() =>
      expect(onFilesChange).toHaveBeenCalledWith([
        expect.objectContaining({ url: "https://cdn.test/me.png", key: "uploads/me.png" }),
      ])
    );
    expect(screen.getByAltText("me.png")).toHaveAttribute("src", "https://cdn.test/me.png");
  });

  it("removes the current picture", async () => {
    const user = userEvent.setup();
    render(<AvatarDropzone value={[stored("current.png")]} onFilesChange={onFilesChange} />);
    expect(screen.getByAltText("current.png")).toBeInTheDocument();

    await user.click(screen.getByRole("button"));

    expect(onFilesChange).toHaveBeenCalledWith([]);
    expect(screen.queryByAltText("current.png")).not.toBeInTheDocument();
  });
});

describe("InlineDropzone", () => {
  it("lists the files it holds and keeps a local pick without uploading", async () => {
    const user = userEvent.setup();
    const { container } = render(
      <InlineDropzone autoUpload={false} maxFiles={3} value={[stored("a.pdf", "application/pdf")]} onFilesChange={onFilesChange} />
    );
    expect(screen.getByText("Choose files")).toBeInTheDocument();
    expect(screen.getByText("a.pdf")).toBeInTheDocument();
    expect(screen.getByText("2 KB")).toBeInTheDocument();

    await user.upload(picker(container), png("b.png"));

    await waitFor(() =>
      expect(onFilesChange).toHaveBeenCalledWith([
        expect.objectContaining({ name: "a.pdf" }),
        expect.objectContaining({ name: "b.png", url: "blob:preview" }),
      ])
    );
    // Inline cards have no reorder buttons.
    expect(screen.queryByTitle("Move up")).not.toBeInTheDocument();
  });
});

describe("BoxDropzone", () => {
  it("reorders its files when reorderable", async () => {
    const user = userEvent.setup();
    render(
      <BoxDropzone maxFiles={5} reorderable value={[stored("one.png"), stored("two.png")]} onFilesChange={onFilesChange} />
    );
    expect(screen.getByText("Click to upload or drag and drop")).toBeInTheDocument();

    await user.click(screen.getAllByTitle("Move down")[0]);

    expect(onFilesChange).toHaveBeenCalledWith([
      expect.objectContaining({ name: "two.png" }),
      expect.objectContaining({ name: "one.png" }),
    ]);
  });

  it("shows the error in place of the description", () => {
    render(<BoxDropzone description="PNG only" error="An image is required" />);
    expect(screen.getByText("An image is required")).toBeInTheDocument();
    expect(screen.queryByText("PNG only")).not.toBeInTheDocument();
  });
});

describe("CompactDropzone", () => {
  it("shows a chip per file and removes one", async () => {
    const user = userEvent.setup();
    render(<CompactDropzone maxFiles={2} value={[stored("clip.mp4", "video/mp4")]} onFilesChange={onFilesChange} />);
    expect(screen.getByText(/Drop files or click to browse/)).toBeInTheDocument();
    expect(screen.getByText("clip.mp4")).toBeInTheDocument();

    await user.click(screen.getByRole("button"));

    expect(onFilesChange).toHaveBeenCalledWith([]);
  });
});

describe("MinimalDropzone", () => {
  it("is an upload button that lists file names", () => {
    render(<MinimalDropzone value={[stored("notes.txt", "text/plain")]} />);
    expect(screen.getByText("Upload file")).toBeInTheDocument();
    expect(screen.getByText("notes.txt")).toBeInTheDocument();
  });
});

describe("Dropzone", () => {
  it("still picks the look from the variant prop", () => {
    const { unmount } = render(<Dropzone variant="avatar" value={[stored("face.png")]} />);
    expect(screen.getByAltText("face.png")).toBeInTheDocument();
    unmount();

    render(<Dropzone variant="inline" provider="minio" />);
    expect(screen.getByText("Choose files")).toBeInTheDocument();
  });

  it("builds a custom look from its parts", () => {
    render(
      <Dropzone.Root label="Receipts" maxFiles={4} value={[stored("receipt.png")]}>
        <Dropzone.Target className="receipts">Drop receipts here</Dropzone.Target>
        <Dropzone.Progress />
        <Dropzone.FileList />
      </Dropzone.Root>
    );
    expect(screen.getByText("Receipts")).toBeInTheDocument();
    expect(screen.getByText("Drop receipts here")).toHaveClass("receipts");
    expect(screen.getByText("receipt.png")).toBeInTheDocument();
    expect(screen.queryByText(/Uploading/)).not.toBeInTheDocument();
  });
});
