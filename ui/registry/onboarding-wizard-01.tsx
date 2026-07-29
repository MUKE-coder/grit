"use client";

import { useState } from "react";
import { ArrowRight, ArrowLeft, Check, User, Building2, Zap } from "lucide-react";

interface OnboardingWizard01Props {
  onComplete?: () => void;
}

const roles = ["Developer", "Designer", "Product Manager", "Founder", "Other"];
const teamSizes = ["Just me", "2–5", "6–15", "16–50", "50+"];

const plans = [
  {
    name: "Starter",
    price: "Free",
    description: "Perfect to explore Grit",
    features: ["3 projects", "Community support"],
  },
  {
    name: "Pro",
    price: "$29/mo",
    description: "For professionals",
    features: ["Unlimited projects", "Priority support", "Grit UI registry"],
    popular: true,
  },
  {
    name: "Enterprise",
    price: "Custom",
    description: "For larger teams",
    features: ["SSO", "SLA", "Dedicated support"],
  },
];

export default function OnboardingWizard01({ onComplete }: OnboardingWizard01Props) {
  const [step, setStep] = useState(0);
  const [fullName, setFullName] = useState("");
  const [role, setRole] = useState("");
  const [company, setCompany] = useState("");
  const [teamSize, setTeamSize] = useState("");
  const [selectedPlan, setSelectedPlan] = useState("Pro");

  const totalSteps = 3;

  const stepMeta = [
    { icon: User, title: "Tell us about you", subtitle: "Help us personalize your experience." },
    { icon: Building2, title: "Set up your workspace", subtitle: "A few details about where you work." },
    { icon: Zap, title: "Choose your plan", subtitle: "Start free, upgrade anytime." },
  ];

  const canAdvance = () => {
    if (step === 0) return fullName.trim() && role;
    if (step === 1) return company.trim() && teamSize;
    return true;
  };

  const handleNext = () => {
    if (step < totalSteps - 1) {
      setStep(step + 1);
    } else {
      onComplete?.();
    }
  };

  const Icon = stepMeta[step].icon;

  return (
    <div className="min-h-screen bg-background flex items-center justify-center px-4 py-12">
      <div className="w-full max-w-lg">
        {/* Step dots */}
        <div className="flex items-center justify-center gap-2 mb-10">
          {Array.from({ length: totalSteps }).map((_, i) => (
            <div
              key={i}
              className={"transition-all duration-300 rounded-full " +
                (i === step
                  ? "w-6 h-2 bg-accent"
                  : i < step
                  ? "w-2 h-2 bg-accent opacity-60"
                  : "w-2 h-2 bg-border")}
            />
          ))}
        </div>

        <div className="rounded-2xl border border-border bg-bg-elevated p-8">
          {/* Header */}
          <div className="flex items-center gap-4 mb-8">
            <div className="w-11 h-11 rounded-xl bg-accent flex items-center justify-center shrink-0">
              <Icon size={20} className="text-white" />
            </div>
            <div>
              <h2 className="text-xl font-bold text-foreground tracking-tight">
                {stepMeta[step].title}
              </h2>
              <p className="text-text-secondary text-sm mt-0.5">{stepMeta[step].subtitle}</p>
            </div>
          </div>

          {/* Step content */}
          {step === 0 && (
            <div className="flex flex-col gap-5">
              <div className="flex flex-col gap-1.5">
                <label className="text-sm font-medium text-foreground">Full name</label>
                <input
                  type="text"
                  value={fullName}
                  onChange={(e) => setFullName(e.target.value)}
                  placeholder="Alex Johnson"
                  className="px-3 py-2.5 rounded-lg border border-border bg-bg-secondary text-foreground text-sm placeholder:text-text-muted focus:outline-none focus:border-accent transition-colors"
                />
              </div>
              <div className="flex flex-col gap-2">
                <label className="text-sm font-medium text-foreground">Your role</label>
                <div className="grid grid-cols-2 sm:grid-cols-3 gap-2">
                  {roles.map((r) => (
                    <button
                      key={r}
                      type="button"
                      onClick={() => setRole(r)}
                      className={"px-3 py-2 rounded-lg border text-sm font-medium transition-colors duration-150 " +
                        (role === r
                          ? "border-accent bg-accent text-white"
                          : "border-border bg-bg-secondary text-text-secondary hover:bg-bg-hover hover:text-foreground")}
                    >
                      {r}
                    </button>
                  ))}
                </div>
              </div>
            </div>
          )}

          {step === 1 && (
            <div className="flex flex-col gap-5">
              <div className="flex flex-col gap-1.5">
                <label className="text-sm font-medium text-foreground">Company name</label>
                <input
                  type="text"
                  value={company}
                  onChange={(e) => setCompany(e.target.value)}
                  placeholder="Acme Corp"
                  className="px-3 py-2.5 rounded-lg border border-border bg-bg-secondary text-foreground text-sm placeholder:text-text-muted focus:outline-none focus:border-accent transition-colors"
                />
              </div>
              <div className="flex flex-col gap-2">
                <label className="text-sm font-medium text-foreground">Team size</label>
                <div className="grid grid-cols-3 sm:grid-cols-5 gap-2">
                  {teamSizes.map((size) => (
                    <button
                      key={size}
                      type="button"
                      onClick={() => setTeamSize(size)}
                      className={"px-3 py-2 rounded-lg border text-sm font-medium transition-colors duration-150 " +
                        (teamSize === size
                          ? "border-accent bg-accent text-white"
                          : "border-border bg-bg-secondary text-text-secondary hover:bg-bg-hover hover:text-foreground")}
                    >
                      {size}
                    </button>
                  ))}
                </div>
              </div>
            </div>
          )}

          {step === 2 && (
            <div className="flex flex-col gap-3">
              {plans.map((plan) => (
                <button
                  key={plan.name}
                  type="button"
                  onClick={() => setSelectedPlan(plan.name)}
                  className={"relative w-full flex items-start gap-4 p-4 rounded-xl border text-left transition-all duration-150 " +
                    (selectedPlan === plan.name
                      ? "border-accent bg-bg-secondary ring-1 ring-accent"
                      : "border-border bg-bg-secondary hover:bg-bg-hover")}
                >
                  <div className={"w-5 h-5 rounded-full border-2 flex items-center justify-center shrink-0 mt-0.5 transition-colors " +
                    (selectedPlan === plan.name ? "border-accent bg-accent" : "border-border")}
                  >
                    {selectedPlan === plan.name && <Check size={11} className="text-white" />}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 flex-wrap">
                      <span className="text-foreground font-semibold text-sm">{plan.name}</span>
                      {(plan as { popular?: boolean }).popular && (
                        <span className="px-2 py-0.5 rounded-full bg-accent text-white text-xs font-bold">
                          Popular
                        </span>
                      )}
                      <span className="ml-auto text-foreground font-bold text-sm">{plan.price}</span>
                    </div>
                    <p className="text-text-muted text-xs mt-0.5">{plan.description}</p>
                    <ul className="mt-2 flex flex-wrap gap-x-4 gap-y-1">
                      {plan.features.map((f) => (
                        <li key={f} className="flex items-center gap-1.5 text-text-secondary text-xs">
                          <Check size={11} className="text-success" />
                          {f}
                        </li>
                      ))}
                    </ul>
                  </div>
                </button>
              ))}
            </div>
          )}

          {/* Navigation */}
          <div className="flex items-center justify-between mt-8 pt-6 border-t border-border">
            <button
              type="button"
              onClick={() => setStep(step - 1)}
              disabled={step === 0}
              className="inline-flex items-center gap-2 px-4 py-2.5 rounded-lg border border-border bg-bg-secondary hover:bg-bg-hover text-foreground text-sm font-medium transition-colors duration-150 disabled:opacity-0 disabled:pointer-events-none"
            >
              <ArrowLeft size={15} />
              Back
            </button>

            <div className="text-text-muted text-xs">
              Step {step + 1} of {totalSteps}
            </div>

            <button
              type="button"
              onClick={handleNext}
              disabled={!canAdvance()}
              className="inline-flex items-center gap-2 px-5 py-2.5 rounded-lg bg-accent hover:bg-accent-hover text-white font-semibold text-sm transition-colors duration-150 disabled:opacity-50"
            >
              {step < totalSteps - 1 ? (
                <>Next <ArrowRight size={15} /></>
              ) : (
                <>Get started <Zap size={15} /></>
              )}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
