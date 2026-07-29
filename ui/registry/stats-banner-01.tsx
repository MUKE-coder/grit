interface Stat {
  value: string;
  label: string;
}

interface StatsBanner01Props {
  stats?: Stat[];
}

export default function StatsBanner01({ stats }: StatsBanner01Props) {
  const displayStats =
    stats && stats.length > 0
      ? stats
      : [
          { value: "10K+", label: "Developers" },
          { value: "$2M+", label: "Saved in dev hours" },
          { value: "500+", label: "Companies" },
          { value: "99.9%", label: "Uptime" },
        ];

  return (
    <section className="bg-bg-secondary border-y border-border py-16 relative overflow-hidden">
      {/* Subtle radial glow */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          background:
            "radial-gradient(ellipse 60% 100% at 50% 50%, rgba(108,92,231,0.07) 0%, transparent 70%)",
        }}
      />

      <div className="container mx-auto px-6 relative z-10">
        <div className="grid grid-cols-2 lg:grid-cols-4 gap-8 lg:gap-0 lg:divide-x lg:divide-border">
          {displayStats.map((stat, i) => (
            <div
              key={i}
              className="flex flex-col items-center justify-center gap-2 py-4"
            >
              <span className="text-5xl font-bold text-foreground tracking-tight">
                {stat.value}
              </span>
              <span className="text-sm text-text-secondary font-medium text-center">
                {stat.label}
              </span>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
