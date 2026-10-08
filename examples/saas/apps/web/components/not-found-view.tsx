import Link from "next/link";

export function NotFoundView() {
  return (
    <div className="flex min-h-[60vh] items-center justify-center px-4">
      <div className="w-full max-w-md text-center">
        <p className="mb-4 text-7xl font-bold text-accent">404</p>
        <h2 className="mb-2 text-2xl font-bold text-foreground">Page not found</h2>
        <p className="mb-8 text-text-muted">
          The page you&apos;re looking for doesn&apos;t exist or has been moved.
        </p>
        <div className="flex gap-3 justify-center">
          <Link
            href="/"
            className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-accent-fg hover:bg-accent-hover transition-colors"
          >
            Home
          </Link>
        </div>
      </div>
    </div>
  );
}
