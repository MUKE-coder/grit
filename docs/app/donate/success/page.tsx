import Link from "next/link";
import type { Metadata } from "next";
import { CheckCircle2, ArrowRight, Heart } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";

export const metadata: Metadata = {
  title: "Thank You | Grit",
  description: "Thank you for supporting the Grit framework!",
};

export default function DonateSuccessPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />

      <section className="relative overflow-hidden">
        <div className="absolute inset-0 -z-10">
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 h-[500px] w-[500px] rounded-full bg-primary/[0.06] blur-[120px]" />
        </div>
        <div className="container max-w-screen-xl py-24 md:py-36 px-6">
          <div className="max-w-2xl mx-auto text-center">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-emerald-500/10 border border-emerald-500/20 mx-auto mb-6">
              <CheckCircle2 className="h-8 w-8 text-emerald-500" />
            </div>
            <h1 className="text-4xl md:text-5xl font-bold tracking-tight mb-6">
              <span className="gradient-text">Thank You!</span>
            </h1>
            <p className="text-lg md:text-xl text-muted-foreground leading-relaxed mb-4 max-w-xl mx-auto">
              Your support means the world. Every contribution helps keep Grit
              free, open source, and actively maintained.
            </p>
            <p className="text-sm text-muted-foreground/60 mb-8">
              A receipt has been sent to your email by Stripe.
            </p>
            <div className="flex flex-col sm:flex-row items-center justify-center gap-3">
              <Button size="lg" className="glow-purple-sm" asChild>
                <Link href="/docs">
                  Back to Docs
                  <ArrowRight className="ml-2 h-4 w-4" />
                </Link>
              </Button>
              <Button
                variant="outline"
                size="lg"
                className="border-border/60 bg-transparent hover:bg-accent/50"
                asChild
              >
                <Link href="/donate">
                  <Heart className="mr-2 h-4 w-4" />
                  Donate Again
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </section>

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
          </div>
        </div>
      </footer>
    </div>
  );
}
