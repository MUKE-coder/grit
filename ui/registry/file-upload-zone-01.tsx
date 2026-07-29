"use client";

// file-upload-zone-01 — Drag and drop upload (Linear attachments-inspired)
import { useState, useRef } from "react";
export function FileUploadZone({
  accept = ".png,.jpg,.pdf,.zip",
  maxSizeMB = 10,
  onFiles,
}: {
  accept?: string; maxSizeMB?: number;
  onFiles?: (files: File[]) => void;
}) {
  const [isDragging, setIsDragging] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);
  const handle = (files: FileList | null) => {
    if (!files) return;
    onFiles?.(Array.from(files));
  };
  return (
    <div
      onDragOver={(e) => { e.preventDefault(); setIsDragging(true); }}
      onDragLeave={() => setIsDragging(false)}
      onDrop={(e) => { e.preventDefault(); setIsDragging(false); handle(e.dataTransfer.files); }}
      onClick={() => inputRef.current?.click()}
      className={"flex flex-col items-center justify-center gap-3 rounded-xl border-2 border-dashed p-10 cursor-pointer transition-all " + (isDragging ? "border-accent bg-accent/5 scale-[1.01]" : "border-border hover:border-accent/50 hover:bg-bg-hover")}
    >
      <div className="w-12 h-12 rounded-full bg-bg-elevated border border-border flex items-center justify-center">
        <svg className="w-6 h-6 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
        </svg>
      </div>
      <div className="text-center">
        <p className="text-sm font-medium text-foreground">Drop files here or <span className="text-accent">browse</span></p>
        <p className="text-xs text-text-muted mt-1">{accept} · Max {maxSizeMB}MB</p>
      </div>
      <input ref={inputRef} type="file" multiple accept={accept} className="hidden" onChange={(e) => handle(e.target.files)} />
    </div>
  );
}

export default FileUploadZone;
