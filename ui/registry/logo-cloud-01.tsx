import React from "react";

const logos = [
  "Acme Corp",
  "Veritas",
  "Nexus Labs",
  "Orbit AI",
  "Stackline",
  "Pulsar Co",
  "Zenith",
  "Forge",
  "Altus",
  "Meridian",
  "Apex IO",
  "Crest",
];

export default function LogoCloud01() {
  const doubled = [...logos, ...logos];

  return (
    <section className="bg-background border-y border-border py-12 overflow-hidden">
      <div className="container mx-auto px-6 mb-8 text-center">
        <p className="text-sm text-text-muted font-medium tracking-widest uppercase">
          Trusted by 500+ teams worldwide
        </p>
      </div>

      <div className="relative flex overflow-hidden">
        {/* Fade masks */}
        <div
          className="absolute left-0 top-0 bottom-0 w-24 z-10 pointer-events-none"
          style={{ background: "linear-gradient(to right, var(--color-background, #0a0a0f), transparent)" }}
        />
        <div
          className="absolute right-0 top-0 bottom-0 w-24 z-10 pointer-events-none"
          style={{ background: "linear-gradient(to left, var(--color-background, #0a0a0f), transparent)" }}
        />

        {/* Scrolling track */}
        <div
          className="flex gap-16 items-center whitespace-nowrap"
          style={{
            animation: "marquee 28s linear infinite",
          }}
        >
          {doubled.map((name, i) => (
            <span
              key={i}
              className="text-text-muted font-semibold text-lg tracking-tight opacity-50 hover:opacity-80 transition-opacity duration-200 select-none"
            >
              {name}
            </span>
          ))}
        </div>
      </div>

      <style dangerouslySetInnerHTML={{ __html: "@keyframes marquee { from { transform: translateX(0); } to { transform: translateX(-50%); } }" }} />
    </section>
  );
}
