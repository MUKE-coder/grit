import { Star } from "lucide-react";

interface Testimonial {
  quote: string;
  name: string;
  role: string;
  company: string;
  stars: number;
}

const testimonials: Testimonial[] = [
  {
    quote:
      "Grit cut our backend setup time from a week to under an hour. The generated types and handlers are exactly what we would have written by hand — just faster.",
    name: "Sofia Ramirez",
    role: "CTO",
    company: "Stackline",
    stars: 5,
  },
  {
    quote:
      "The admin panel generation alone is worth it. We launched a fully featured dashboard to our clients in two days. Previously that was a month of work.",
    name: "James Whitfield",
    role: "Lead Engineer",
    company: "Orbit AI",
    stars: 5,
  },
  {
    quote:
      "Finally a Go framework that doesn't make me choose between productivity and control. The grit generate command is genuinely magical.",
    name: "Priya Nambiar",
    role: "Senior Developer",
    company: "Nexus Labs",
    stars: 5,
  },
];

function TestimonialCard({ t }: { t: Testimonial }) {
  const initials = t.name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);

  return (
    <div className="flex flex-col gap-5 p-6 rounded-2xl border border-border bg-bg-elevated hover:bg-bg-hover transition-colors duration-200 h-full">
      <div className="flex items-center gap-1">
        {Array.from({ length: 5 }).map((_, i) => (
          <Star
            key={i}
            size={14}
            className={i < t.stars ? "text-warning fill-warning" : "text-border fill-border"}
          />
        ))}
      </div>
      <blockquote className="text-foreground text-sm leading-relaxed flex-1">
        &ldquo;{t.quote}&rdquo;
      </blockquote>
      <div className="flex items-center gap-3 pt-2 border-t border-border">
        <div className="w-9 h-9 rounded-full bg-accent flex items-center justify-center text-white text-xs font-bold shrink-0">
          {initials}
        </div>
        <div className="flex flex-col min-w-0">
          <span className="text-foreground text-sm font-semibold truncate">{t.name}</span>
          <span className="text-text-muted text-xs truncate">
            {t.role} at {t.company}
          </span>
        </div>
      </div>
    </div>
  );
}

export default function TestimonialGrid01() {
  return (
    <section className="bg-background py-24">
      <div className="container mx-auto px-6">
        <div className="text-center mb-16">
          <h2 className="text-4xl font-bold text-foreground tracking-tight">
            Loved by developers
          </h2>
          <p className="mt-4 text-text-secondary text-lg max-w-xl mx-auto">
            Teams across the world trust Grit to ship production-ready apps faster.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {testimonials.map((t) => (
            <TestimonialCard key={t.name} t={t} />
          ))}
        </div>
      </div>
    </section>
  );
}
