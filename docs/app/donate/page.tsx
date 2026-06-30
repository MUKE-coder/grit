"use client";

import { useState } from "react";
import Link from "next/link";
import { Heart, Sparkles, Loader2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { SiteHeader } from "@/components/site-header";

const PRESET_AMOUNTS = [
  { value: 500, label: "$5", description: "Buy a coffee" },
  { value: 1000, label: "$10", description: "Support a feature" },
  { value: 2500, label: "$25", description: "Fuel a sprint" },
  { value: 5000, label: "$50", description: "Champion the project" },
];

export default function DonatePage() {
  const [selectedAmount, setSelectedAmount] = useState<number>(1000);
  const [customAmount, setCustomAmount] = useState("");
  const [isCustom, setIsCustom] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");

  const effectiveAmount = isCustom
    ? Math.round(parseFloat(customAmount) * 100)
    : selectedAmount;

  async function handleDonate() {
    setError("");

    if (
      isCustom &&
      (!customAmount ||
        isNaN(parseFloat(customAmount)) ||
        parseFloat(customAmount) < 1)
    ) {
      setError("Please enter a valid amount ($1 minimum).");
      return;
    }

    setIsLoading(true);

    try {
      const res = await fetch("/api/checkout", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ amount: effectiveAmount }),
      });

      const data = await res.json();

      if (!res.ok) {
        setError(data.error || "Something went wrong.");
        return;
      }

      window.location.href = data.url;
    } catch {
      setError("Failed to connect. Please try again.");
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <div className="min-h-screen bg-background isolate">
      <SiteHeader />

      {/* Hero */}
      <section className="relative overflow-hidden">
        <div className="absolute inset-0 -z-10">
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 h-[500px] w-[500px] rounded-full bg-pink-500/[0.06] blur-[120px]" />
        </div>
        <div className="container max-w-screen-xl py-24 md:py-36 px-6">
          <div className="max-w-3xl mx-auto text-center">
            <div className="flex items-center justify-center gap-2 rounded-full border border-border/60 bg-pink-500/10 px-4 py-1.5 mb-6 mx-auto w-fit">
              <Heart className="h-3.5 w-3.5 text-pink-500" />
              <span className="text-xs font-mono uppercase tracking-wider text-muted-foreground">
                Support Open Source
              </span>
            </div>
            <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold tracking-tight mb-6">
              <span className="gradient-text">Support Grit</span>
            </h1>
            <p className="text-lg md:text-xl text-muted-foreground leading-relaxed mb-8 max-w-2xl mx-auto">
              Grit is free and open source. Your donation helps fund continued
              development, new features, documentation, and community support.
            </p>
          </div>
        </div>
      </section>

      {/* Donation Form */}
      <section className="border-t border-border/30">
        <div className="container max-w-screen-xl py-20 md:py-28 px-6">
          <div className="max-w-xl mx-auto">
            {/* Preset amount cards */}
            <div className="grid grid-cols-2 gap-3 mb-6">
              {PRESET_AMOUNTS.map((preset) => (
                <button
                  key={preset.value}
                  onClick={() => {
                    setSelectedAmount(preset.value);
                    setIsCustom(false);
                    setError("");
                  }}
                  className={cn(
                    "group rounded-xl border p-4 text-left transition-all duration-200",
                    !isCustom && selectedAmount === preset.value
                      ? "border-primary/40 bg-primary/[0.08] ring-1 ring-primary/20"
                      : "border-border/40 bg-card/50 hover:border-primary/20 hover:bg-card/80"
                  )}
                >
                  <span className="text-2xl font-bold text-foreground">
                    {preset.label}
                  </span>
                  <span className="text-xs text-muted-foreground/60 block mt-1">
                    {preset.description}
                  </span>
                </button>
              ))}
            </div>

            {/* Custom amount */}
            <button
              onClick={() => {
                setIsCustom(true);
                setError("");
              }}
              className={cn(
                "w-full rounded-xl border p-4 text-left transition-all duration-200 mb-6",
                isCustom
                  ? "border-primary/40 bg-primary/[0.08] ring-1 ring-primary/20"
                  : "border-border/40 bg-card/50 hover:border-primary/20 hover:bg-card/80"
              )}
            >
              <span className="text-sm font-semibold text-foreground">
                Custom Amount
              </span>
              {isCustom && (
                <div className="mt-3 flex items-center gap-2">
                  <span className="text-lg font-bold text-muted-foreground">
                    $
                  </span>
                  <Input
                    type="number"
                    min="1"
                    max="999"
                    step="1"
                    placeholder="Enter amount"
                    value={customAmount}
                    onChange={(e) => setCustomAmount(e.target.value)}
                    onClick={(e) => e.stopPropagation()}
                    className="bg-background/50 border-border/40"
                    autoFocus
                  />
                </div>
              )}
            </button>

            {/* Error message */}
            {error && (
              <p className="text-sm text-destructive mb-4">{error}</p>
            )}

            {/* Donate button */}
            <Button
              size="lg"
              className="w-full glow-purple-sm h-12 text-base"
              onClick={handleDonate}
              disabled={isLoading}
            >
              {isLoading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  Redirecting to Stripe...
                </>
              ) : (
                <>
                  <Heart className="mr-2 h-4 w-4" />
                  Donate{" "}
                  {isCustom && customAmount
                    ? `$${customAmount}`
                    : `$${(selectedAmount / 100).toFixed(0)}`}
                </>
              )}
            </Button>

            <p className="text-xs text-muted-foreground/40 text-center mt-4">
              Payments are processed securely by Stripe. No card details are
              stored on our servers.
            </p>
          </div>
        </div>
      </section>

      {/* Where Your Support Goes */}
      <section className="border-t border-border/30 bg-accent/20">
        <div className="container max-w-screen-xl py-20 md:py-28 px-6">
          <div className="text-center mb-14">
            <h2 className="text-3xl md:text-4xl font-bold tracking-tight mb-4">
              Where Your Support Goes
            </h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              Every contribution directly funds the tools and resources that make
              Grit better for everyone.
            </p>
          </div>
          <div className="grid gap-4 sm:grid-cols-2 max-w-3xl mx-auto">
            {[
              "New features and framework improvements",
              "Documentation, tutorials, and video courses",
              "Bug fixes and security patches",
              "Community support and issue triage",
              "Infrastructure and hosting costs",
              "Open-source sustainability",
            ].map((point) => (
              <div
                key={point}
                className="flex items-start gap-3 rounded-lg border border-border/30 bg-card/30 px-4 py-3"
              >
                <Sparkles className="h-4 w-4 text-primary mt-0.5 shrink-0" />
                <span className="text-sm text-foreground/80">{point}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-border/30 py-8 px-6">
        <div className="container max-w-screen-xl flex flex-col sm:flex-row items-center justify-between gap-4 text-sm text-muted-foreground/50">
          <span>Grit Framework &mdash; Go + React. Built with Grit.</span>
          <div className="flex items-center gap-4">
            <Link
              href="/docs"
              className="hover:text-foreground transition-colors"
            >
              Docs
            </Link>
            <Link
              href="/hire"
              className="hover:text-foreground transition-colors"
            >
              Hire Us
            </Link>
            <Link
              href="https://github.com/MUKE-coder/grit"
              target="_blank"
              rel="noreferrer"
              className="hover:text-foreground transition-colors"
            >
              GitHub
            </Link>
          </div>
        </div>
      </footer>
    </div>
  );
}
