"use client";

// profile-avatar-01 — User profile avatar with upload (Clerk / GitHub-inspired)
import { useRef, useState } from "react";
export function ProfileAvatar({
  name = "Sarah Chen",
  email = "sarah@example.com",
  avatarUrl,
  onUpload,
}: {
  name?: string; email?: string; avatarUrl?: string;
  onUpload?: (file: File) => void;
}) {
  const [preview, setPreview] = useState(avatarUrl ?? "");
  const inputRef = useRef<HTMLInputElement>(null);
  const initials = name.split(" ").map((n) => n[0]).join("").toUpperCase().slice(0, 2);
  const handleFile = (file: File | undefined) => {
    if (!file) return;
    const url = URL.createObjectURL(file);
    setPreview(url);
    onUpload?.(file);
  };
  return (
    <div className="flex items-center gap-4">
      <div className="relative group cursor-pointer" onClick={() => inputRef.current?.click()}>
        <div className="w-20 h-20 rounded-full overflow-hidden bg-accent/20 flex items-center justify-center border-2 border-border">
          {preview ? (
            <img src={preview} alt={name} className="w-full h-full object-cover" />
          ) : (
            <span className="text-2xl font-bold text-accent">{initials}</span>
          )}
        </div>
        <div className="absolute inset-0 rounded-full bg-background/70 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity">
          <svg className="w-6 h-6 text-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z M15 13a3 3 0 11-6 0 3 3 0 016 0z" /></svg>
        </div>
        <input ref={inputRef} type="file" accept="image/*" className="hidden" onChange={(e) => handleFile(e.target.files?.[0])} />
      </div>
      <div>
        <p className="text-sm font-semibold text-foreground">{name}</p>
        <p className="text-xs text-text-muted">{email}</p>
        <button onClick={() => inputRef.current?.click()} className="text-xs text-accent hover:underline mt-1.5">Upload photo</button>
        {preview && <button onClick={() => { setPreview(""); }} className="text-xs text-text-muted hover:text-danger ml-2">Remove</button>}
      </div>
    </div>
  );
}

export default ProfileAvatar;
