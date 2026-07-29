"use client";

import { useState } from "react";
import { Check, Zap } from "lucide-react";

interface Plan {
  name: string;
  monthlyPrice: number;
  annualPrice: number;
  description: string;
  features: string[];
  cta: string;
  popular: boolean;
  enterprise: boolean;
}

const plans: Plan[] = [
  {
    name: "Starter",
    monthlyPrice: 0,
    annualPrice: 0,
    description: "For solo developers and side projects.",
    features: [
      "Up to 3 projects",
      "grit generate (all resource types)",
      "Docker Compose setup",
      "Community support",
      "1 team member",
    ],
    cta: "Start for free",
    popular: false,
    enterprise: false,
  },
  {
    name: "Pro",
    monthlyPrice: 29,
    annualPrice: 23,
    description: "For professional developers and small teams.",
    features: [
      "Unlimited projects",
      "Everything in Starter",
      "Grit UI component registry",
      "Priority email support",
      "Up to 10 team members",
      "Audit logs",
    ],
    cta: "Get started",
    popular: true,
    enterprise: false,
  },
  {
    name: "Enterprise",
    monthlyPrice: 99,
    annualPrice: 79,
    description: "For growing teams with advanced needs.",
    features: [
      "Everything in Pro",
      "SSO (SAML / OIDC)",
      "Custom deployment support",
      "SLA + dedicated Slack",
      "Unlimited team members",
      "Custom billing",
    ],
    cta: "Contact sales",
    popular: false,
    enterprise: true,
  },
];

export default function PricingToggle01() {
  const [annual, setAnnual] = useState(false);

  return (
    <section className="bg-background py-24">
      <div className="container mx-auto px-6">
        {/* Header */}
        <div className="text-center mb-12">
          <h2 className="text-4xl font-bold text-foreground tracking-tight">
            Simple, transparent pricing
          </h2>
          <p className="mt-4 text-text-secondary text-lg max-w-xl mx-auto">
            Start free, upgrade when you&rsquo;re ready. No hidden fees.
          </p>

          {/* Toggle */}
          <div className="mt-8 inline-flex items-center gap-3 p-1 rounded-xl border border-border bg-bg-elevated">
            <button
              onClick={() => setAnnual(false)}
              className={"px-5 py-2 rounded-lg text-sm font-semibold transition-colors duration-150 " +
                (!annual
                  ? "bg-accent text-white"
                  : "text-text-secondary hover:text-foreground")}
            >
              Monthly
            </button>
            <button
              onClick={() => setAnnual(true)}
              className={"px-5 py-2 rounded-lg text-sm font-semibold transition-colors duration-150 flex items-center gap-2 " +
                (annual
                  ? "bg-accent text-white"
                  : "text-text-secondary hover:text-foreground")}
            >
              Annual
              <span className="text-xs px-2 py-0.5 rounded-full bg-success text-background font-bold">
                Save 20%
              </span>
            </button>
          </div>
        </div>

        {/* Cards */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-5xl mx-auto">
          {plans.map((plan) => (
            <div
              key={plan.name}
              className={"relative flex flex-col rounded-2xl border p-7 transition-all duration-200 " +
                (plan.popular
                  ? "border-accent bg-bg-elevated ring-2 ring-accent ring-offset-2 ring-offset-background scale-105"
                  : "border-border bg-bg-elevated hover:bg-bg-hover")}
            >
              {plan.popular && (
                <div className="absolute -top-4 left-1/2 -translate-x-1/2">
                  <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-accent text-white text-xs font-bold tracking-wide shadow-lg">
                    <Zap size={11} />
                    Most popular
                  </span>
                </div>
              )}

              <div className="mb-6">
                <h3 className="text-foreground font-bold text-lg">{plan.name}</h3>
                <p className="text-text-muted text-sm mt-1">{plan.description}</p>
              </div>

              <div className="mb-6">
                <div className="flex items-end gap-1">
                  <span className="text-5xl font-bold text-foreground tracking-tight">
                    {plan.enterprise
                      ? "Custom"
                      : (annual ? "$" + plan.annualPrice : "$" + plan.monthlyPrice)}
                  </span>
                  {!plan.enterprise && (
                    <span className="text-text-muted text-sm mb-2">/mo</span>
                  )}
                </div>
                {annual && !plan.enterprise && plan.annualPrice > 0 && (
                  <p className="text-text-muted text-xs mt-1">
                    Billed annually (${plan.annualPrice * 12}/yr)
                  </p>
                )}
              </div>

              <ul className="flex flex-col gap-3 mb-8 flex-1">
                {plan.features.map((feature, i) => (
                  <li key={i} className="flex items-center gap-2.5">
                    <Check size={15} className="text-success shrink-0" />
                    <span className="text-text-secondary text-sm">{feature}</span>
                  </li>
                ))}
              </ul>

              <a
                href={plan.enterprise ? "/contact" : "/signup"}
                className={"w-full inline-flex items-center justify-center py-3 rounded-lg font-semibold text-sm transition-colors duration-150 " +
                  (plan.popular
                    ? "bg-accent hover:bg-accent-hover text-white"
                    : "border border-border bg-bg-secondary hover:bg-bg-hover text-foreground")}
              >
                {plan.cta}
              </a>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
