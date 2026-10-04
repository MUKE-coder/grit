import type { ReactNode } from "react";

interface StackedCellProps {
  top: string;
  bottom?: string;
  topClassName?: string;
  bottomClassName?: string;
}

// Renders a two-line cell: bold primary text on top, muted secondary
// below. Designed to be called directly from a resource definition's
// cell: callback — no JSX needed at the call site, so resource files
// stay .ts instead of being forced to .tsx.
export function StackedCell({
  top,
  bottom,
  topClassName,
  bottomClassName,
}: StackedCellProps): ReactNode {
  return (
    <div className="flex flex-col">
      <span className={topClassName ?? "font-medium text-foreground"}>{top}</span>
      {bottom && (
        <span className={bottomClassName ?? "text-xs text-text-muted"}>
          {bottom}
        </span>
      )}
    </div>
  );
}
