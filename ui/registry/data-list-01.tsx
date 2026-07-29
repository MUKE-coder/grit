// data-list-01 — Key-value data list (Stripe details panel-inspired)
export function DataList({
  items = [
    { label: "Status", value: "Active", valueClass: "text-success" },
    { label: "Plan", value: "Pro Monthly" },
    { label: "Next billing", value: "Apr 1, 2025" },
    { label: "Amount", value: "$29.00" },
  ],
}: {
  items?: Array<{ label: string; value: string; valueClass?: string }>;
}) {
  return (
    <dl className="divide-y divide-border">
      {items.map((item) => (
        <div key={item.label} className="flex items-center justify-between py-3 gap-4">
          <dt className="text-sm text-text-secondary">{item.label}</dt>
          <dd className={"text-sm font-medium " + (item.valueClass ?? "text-foreground")}>{item.value}</dd>
        </div>
      ))}
    </dl>
  );
}

export default DataList;
