import { Github } from "lucide-react";

interface SocialLoginButtons01Props {
  googleHref?: string;
  githubHref?: string;
}

export default function SocialLoginButtons01({
  googleHref = "/auth/google",
  githubHref = "/auth/github",
}: SocialLoginButtons01Props) {
  return (
    <div className="flex flex-col gap-4 w-full">
      {/* Divider */}
      <div className="relative flex items-center gap-4">
        <div className="flex-1 h-px bg-border" />
        <span className="text-text-muted text-xs font-medium whitespace-nowrap">
          or continue with
        </span>
        <div className="flex-1 h-px bg-border" />
      </div>

      {/* Buttons */}
      <a
        href={googleHref}
        className="w-full flex items-center justify-center gap-3 px-4 py-2.5 rounded-lg border border-border bg-bg-elevated hover:bg-bg-hover text-foreground text-sm font-medium transition-colors duration-150"
      >
        <svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true">
          <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
          <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
          <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
          <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
        </svg>
        Continue with Google
      </a>

      <a
        href={githubHref}
        className="w-full flex items-center justify-center gap-3 px-4 py-2.5 rounded-lg border border-border bg-bg-elevated hover:bg-bg-hover text-foreground text-sm font-medium transition-colors duration-150"
      >
        <Github size={18} />
        Continue with GitHub
      </a>
    </div>
  );
}
