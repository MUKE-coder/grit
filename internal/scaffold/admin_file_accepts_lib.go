package scaffold

// admin_file_accepts_lib.go emits lib/file-accepts.ts into the admin
// scaffold. It bridges the CLI accept-alias vocabulary (image, video,
// pdf, doc, ...) used in resource definitions to:
//
//   1. The react-dropzone `accept` map (alias → MIME → extension list).
//   2. The /api/uploads query string (`?accepts=image,pdf&max_size=N`)
//      so the server can enforce per-field validation matching the
//      client's accept set.
//
// Keeping this in a single shared module avoids drift between the
// FileField and the FilesField -- both call the same helpers.

func adminFileAcceptsLib() string {
	return `// File accept-alias bridge for the v3.31.30 file fields. The CLI
// syntax (image:file:image, attachments:files:[pdf,doc,image]) gives
// us aliases like "image" / "pdf"; this module translates them into
// the formats react-dropzone and the /api/uploads endpoint expect.
//
// Keep the alias keys in sync with internal/files/accepts.go on the Go
// side -- both sides resolve aliases the same way, but the server
// re-validates the upload so client tampering is harmless.

import type { Accept } from "react-dropzone";

type AliasMap = { mimes: string[]; exts: string[] };

const ALIASES: Record<string, AliasMap> = {
  image: {
    mimes: ["image/jpeg", "image/png", "image/gif", "image/webp", "image/avif", "image/svg+xml"],
    exts: [".jpg", ".jpeg", ".png", ".gif", ".webp", ".avif", ".svg"],
  },
  video: {
    mimes: ["video/mp4", "video/webm", "video/quicktime", "video/x-msvideo", "video/x-matroska"],
    exts: [".mp4", ".webm", ".mov", ".avi", ".mkv"],
  },
  audio: {
    mimes: ["audio/mpeg", "audio/wav", "audio/ogg", "audio/x-m4a", "audio/webm"],
    exts: [".mp3", ".wav", ".ogg", ".m4a"],
  },
  pdf: {
    mimes: ["application/pdf"],
    exts: [".pdf"],
  },
  doc: {
    mimes: [
      "application/msword",
      "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    ],
    exts: [".doc", ".docx"],
  },
  excel: {
    mimes: [
      "application/vnd.ms-excel",
      "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
    ],
    exts: [".xls", ".xlsx"],
  },
  csv: {
    mimes: ["text/csv", "application/vnd.ms-excel"],
    exts: [".csv"],
  },
  zip: {
    mimes: ["application/zip", "application/x-zip-compressed"],
    exts: [".zip"],
  },
  archive: {
    mimes: [
      "application/zip",
      "application/x-zip-compressed",
      "application/x-tar",
      "application/gzip",
      "application/x-rar-compressed",
      "application/x-7z-compressed",
    ],
    exts: [".zip", ".tar", ".gz", ".tgz", ".rar", ".7z"],
  },
};

// acceptsToReactDropzoneFormat converts our alias list to react-dropzone's
// Accept map (a record of MIME -> extension list). Returning undefined
// (the "all" case) makes react-dropzone accept any file.
export function acceptsToReactDropzoneFormat(accepts: string[]): Accept | undefined {
  if (!accepts || accepts.length === 0) return undefined;
  if (accepts.includes("all")) return undefined;

  const out: Accept = {};
  for (const alias of accepts) {
    const m = ALIASES[alias];
    if (!m) continue;
    for (const mime of m.mimes) {
      out[mime] = out[mime] ? Array.from(new Set([...out[mime], ...m.exts])) : [...m.exts];
    }
  }
  return Object.keys(out).length > 0 ? out : undefined;
}

// buildUploadEndpoint constructs the per-field upload URL with the
// query params the API needs (?accepts=...&max_size=...). The base
// path stays /api/uploads -- this is just attaching field metadata.
export function buildUploadEndpoint(accepts: string[] | undefined, maxBytes: number | undefined): string {
  const params = new URLSearchParams();
  if (accepts && accepts.length > 0) {
    params.set("accepts", accepts.join(","));
  }
  if (maxBytes && maxBytes > 0) {
    params.set("max_size", String(maxBytes));
  }
  const qs = params.toString();
  return qs ? "/api/uploads?" + qs : "/api/uploads";
}

// mimeKind groups a MIME type into a high-level kind for icon /
// preview selection. Used by the dropzone's existing FilePreview
// switch when it needs a fallback icon for a generic upload.
export function mimeKind(mime: string): "image" | "video" | "audio" | "pdf" | "doc" | "excel" | "csv" | "zip" | "other" {
  if (mime.startsWith("image/")) return "image";
  if (mime.startsWith("video/")) return "video";
  if (mime.startsWith("audio/")) return "audio";
  if (mime === "application/pdf") return "pdf";
  if (mime.includes("wordprocessing") || mime === "application/msword") return "doc";
  if (mime.includes("spreadsheet") || mime === "application/vnd.ms-excel") return "excel";
  if (mime === "text/csv") return "csv";
  if (mime.includes("zip") || mime.includes("tar") || mime.includes("gzip") || mime.includes("rar") || mime.includes("7z")) return "zip";
  return "other";
}
`
}
