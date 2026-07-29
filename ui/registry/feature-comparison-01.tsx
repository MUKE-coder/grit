// feature-comparison-01 — Feature comparison table (Linear pricing page-inspired)
export function FeatureComparison({
  plans = ["Starter", "Pro", "Enterprise"],
  categories = [
    {
      name: "Core",
      features: [
        { label: "Projects", values: ["5", "Unlimited", "Unlimited"] },
        { label: "Team members", values: ["3", "15", "Unlimited"] },
        { label: "Storage", values: ["5 GB", "100 GB", "1 TB"] },
      ],
    },
    {
      name: "Features",
      features: [
        { label: "Analytics", values: [false, true, true] },
        { label: "Custom domain", values: [false, true, true] },
        { label: "Priority support", values: [false, false, true] },
        { label: "SSO / SAML", values: [false, false, true] },
        { label: "Audit log", values: [false, false, true] },
      ],
    },
  ],
}: {
  plans?: string[];
  categories?: Array<{
    name: string;
    features: Array<{ label: string; values: (string | boolean)[] }>;
  }>;
}) {
  const Cell = ({ v }: { v: string | boolean }) => {
    if (typeof v === "boolean") {
      return v ? (
        <svg className="w-5 h-5 text-success mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
        </svg>
      ) : (
        <svg className="w-4 h-4 text-text-muted mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 12H4" />
        </svg>
      );
    }
    return <span className="text-sm text-foreground font-medium">{v}</span>;
  };
  return (
    <div className="w-full overflow-x-auto">
      <table className="w-full">
        <thead>
          <tr className="border-b border-border">
            <th className="text-left py-4 pr-4 text-sm font-medium text-text-muted w-48">Features</th>
            {plans.map((plan) => (
              <th key={plan} className="text-center py-4 px-4 text-sm font-semibold text-foreground">{plan}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {categories.map((cat) => (
            <>
              <tr key={cat.name}>
                <td colSpan={plans.length + 1} className="pt-6 pb-2">
                  <span className="text-xs font-semibold text-text-muted uppercase tracking-wider">{cat.name}</span>
                </td>
              </tr>
              {cat.features.map((feat) => (
                <tr key={feat.label} className="border-b border-border/50 hover:bg-bg-hover transition-colors">
                  <td className="py-3 pr-4 text-sm text-text-secondary">{feat.label}</td>
                  {feat.values.map((v, i) => (
                    <td key={i} className="py-3 px-4 text-center"><Cell v={v} /></td>
                  ))}
                </tr>
              ))}
            </>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default FeatureComparison;
