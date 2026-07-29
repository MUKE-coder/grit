"use client";

// error-page-01 — Error page (404 / 500) (Vercel / Linear-inspired)
export function ErrorPage({
  code = "404",
  title = "Page not found",
  description = "The page you're looking for doesn't exist or has been moved.",
  homeHref = "/",
  backLabel = "Go back",
}: {
  code?: string; title?: string; description?: string;
  homeHref?: string; backLabel?: string;
}) {
  return (
    <div className="min-h-[60vh] flex flex-col items-center justify-center px-4 text-center">
      <div className="text-8xl font-black text-bg-hover mb-4 select-none">{code}</div>
      <h1 className="text-2xl font-bold text-foreground mb-2">{title}</h1>
      <p className="text-text-secondary mb-8 max-w-md">{description}</p>
      <div className="flex items-center gap-3">
        <button onClick={() => window.history.back()} className="px-4 py-2 rounded-lg border border-border text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors text-sm">
          {backLabel}
        </button>
        <a href={homeHref} className="px-4 py-2 rounded-lg bg-accent hover:bg-accent/90 text-white text-sm font-semibold transition-colors">
          Back to home
        </a>
      </div>
    </div>
  );
}

export default ErrorPage;
