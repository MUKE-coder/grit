
const DOCS_URL = "https://gritframework.dev/docs";

export function Footer() {
  return (
    <footer className="border-t border-border/50 bg-background">
      <div className="mx-auto max-w-5xl px-6 py-8">
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2 text-sm text-text-muted">
            <span className="font-semibold text-text-secondary">library</span>
            <span className="text-border">·</span>
            <span>Built with Grit</span>
          </div>
          <div className="flex items-center gap-6 text-sm text-text-muted">
            <a
              href="https://github.com/MUKE-coder/grit"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-foreground transition-colors"
            >
              GitHub
            </a>
            <a
              href={DOCS_URL}
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-foreground transition-colors"
            >
              Documentation
            </a>
          </div>
        </div>
        <div className="mt-6 pt-6 border-t border-border/30 text-center">
          <p className="text-xs text-text-muted">
            &copy; {new Date().getFullYear()} library. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
}
