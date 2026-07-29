import { Star } from "lucide-react";

interface TestimonialCard01Props {
  quote?: string;
  name?: string;
  role?: string;
  company?: string;
  avatar?: string;
  stars?: number;
}

export default function TestimonialCard01({
  quote = "We replaced three weeks of setup with one command. The generated code is ordinary Go and TypeScript, so onboarding a new engineer took an afternoon instead of a sprint.",
  name = "Amara Okonkwo",
  role = "Staff Engineer",
  company = "Northwind",
  avatar,
  stars = 5,
}: TestimonialCard01Props) {
  const initials = name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);

  return (
    <div className="flex flex-col gap-5 p-6 rounded-2xl border border-border bg-bg-elevated hover:bg-bg-hover transition-colors duration-200">
      {/* Stars */}
      <div className="flex items-center gap-1">
        {Array.from({ length: 5 }).map((_, i) => (
          <Star
            key={i}
            size={14}
            className={i < stars ? "text-warning fill-warning" : "text-border fill-border"}
          />
        ))}
      </div>

      {/* Quote */}
      <blockquote className="text-foreground text-sm leading-relaxed">
        &ldquo;{quote}&rdquo;
      </blockquote>

      {/* Author */}
      <div className="flex items-center gap-3 pt-2 border-t border-border">
        {avatar ? (
          <img
            src={avatar}
            alt={name}
            className="w-9 h-9 rounded-full object-cover border border-border"
          />
        ) : (
          <div className="w-9 h-9 rounded-full bg-accent flex items-center justify-center text-white text-xs font-bold shrink-0">
            {initials}
          </div>
        )}
        <div className="flex flex-col min-w-0">
          <span className="text-foreground text-sm font-semibold truncate">{name}</span>
          <span className="text-text-muted text-xs truncate">
            {role} at {company}
          </span>
        </div>
      </div>
    </div>
  );
}
