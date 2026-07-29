"use client";

import { useState } from "react";
import { ChevronDown } from "lucide-react";

interface Faq {
  question: string;
  answer: string;
}

interface FaqAccordion01Props {
  faqs?: Faq[];
}

const defaultFaqs: Faq[] = [
  {
    question: "What is Grit?",
    answer:
      "Grit is a full-stack meta-framework that combines a Go (Gin + GORM) backend with a Next.js frontend. A single CLI scaffolds your entire project, generates resources, and keeps your TypeScript types in sync with your Go models.",
  },
  {
    question: "Do I need to know both Go and TypeScript?",
    answer:
      "Basic Go knowledge is helpful, but Grit is designed to be beginner-friendly. The scaffolded code follows strict conventions so you can learn as you build. The frontend is pure TypeScript/React, so frontend developers can start there.",
  },
  {
    question: "How is Grit different from other Go frameworks?",
    answer:
      "Most Go frameworks handle only the backend. Grit generates the full stack — Go API, Next.js admin panel, web frontend, shared types, Docker setup, and tests — all wired together and production-ready from day one.",
  },
  {
    question: "Is there a self-hosted option?",
    answer:
      "Yes. Grit is open-source and fully self-hostable. The included Docker Compose file starts PostgreSQL, Redis, MinIO, and Mailhog locally. You can deploy to any cloud provider that runs Docker.",
  },
];

export default function FaqAccordion01({ faqs }: FaqAccordion01Props) {
  const items = faqs && faqs.length > 0 ? faqs : defaultFaqs;
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  return (
    <section className="bg-background py-24">
      <div className="container mx-auto px-6 max-w-3xl">
        <div className="text-center mb-14">
          <h2 className="text-4xl font-bold text-foreground tracking-tight">
            Frequently asked questions
          </h2>
          <p className="mt-4 text-text-secondary text-lg">
            Everything you need to know before getting started.
          </p>
        </div>

        <div className="flex flex-col divide-y divide-border rounded-2xl border border-border overflow-hidden">
          {items.map((faq, i) => {
            const isOpen = openIndex === i;
            return (
              <div key={i} className="bg-bg-elevated">
                <button
                  onClick={() => setOpenIndex(isOpen ? null : i)}
                  className="w-full flex items-center justify-between gap-4 px-6 py-5 text-left hover:bg-bg-hover transition-colors duration-150"
                  aria-expanded={isOpen}
                >
                  <span className="text-foreground font-semibold text-sm leading-snug">
                    {faq.question}
                  </span>
                  <ChevronDown
                    size={18}
                    className="text-text-secondary shrink-0 transition-transform duration-200"
                    style={{ transform: isOpen ? "rotate(180deg)" : "rotate(0deg)" }}
                  />
                </button>

                <div
                  className="overflow-hidden transition-all duration-200"
                  style={{ maxHeight: isOpen ? "400px" : "0px" }}
                >
                  <p className="px-6 pb-5 text-text-secondary text-sm leading-relaxed">
                    {faq.answer}
                  </p>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
