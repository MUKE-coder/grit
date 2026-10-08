import { Database } from "@/lib/icons";
import { useT } from "@/lib/i18n";

export function TableEmptyState() {
  const t = useT();
  return (
    <div className="flex flex-col items-center justify-center py-16 px-4">
      <div className="rounded-full bg-bg-tertiary p-4 mb-4">
        <Database className="h-8 w-8 text-text-muted" />
      </div>
      <h3 className="text-sm font-medium text-foreground mb-1">{t("table.empty", "No records found")}</h3>
      <p className="text-sm text-text-muted">
        {t("table.emptyHint", "Try adjusting your search or filters")}
      </p>
    </div>
  );
}
