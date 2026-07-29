// card-grid-01 — Responsive card grid wrapper with header
export function CardGrid({
  title,
  description,
  action,
  columns = 3,
  children,
}: {
  title?: string; description?: string; action?: React.ReactNode;
  columns?: 2 | 3 | 4; children?: React.ReactNode;
}) {
  const colClass = { 2: "sm:grid-cols-2", 3: "sm:grid-cols-2 lg:grid-cols-3", 4: "sm:grid-cols-2 lg:grid-cols-4" };
  return (
    <section className="flex flex-col gap-6">
      {(title || action) && (
        <div className="flex items-start justify-between gap-4">
          <div>
            {title && <h2 className="text-lg font-semibold text-foreground">{title}</h2>}
            {description && <p className="text-sm text-text-secondary mt-0.5">{description}</p>}
          </div>
          {action}
        </div>
      )}
      <div className={"grid grid-cols-1 gap-4 " + colClass[columns]}>{children}</div>
    </section>
  );
}

export default CardGrid;
