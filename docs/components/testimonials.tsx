import Image from 'next/image'
import Link from 'next/link'
import { ArrowRight, Github, Quote } from 'lucide-react'
import { TESTIMONIALS, TESTIMONIAL_ISSUE_URL, hasTestimonials } from '@/config/testimonials'

/**
 * The testimonials section.
 *
 * Two states, and the empty one is the default rather than an accident. It
 * says plainly that there is nothing here and that nothing will be invented,
 * which is a stronger claim to a sceptical reader than three anonymous quotes
 * about how intuitive everything is.
 *
 * Once data/testimonials.json has entries, this renders them instead. No flag
 * to flip: the presence of approved rows is the switch.
 */
export function Testimonials() {
  return (
    <section className="relative overflow-hidden border-t border-border/40 px-6 py-24">
      <div className="absolute inset-0 -z-10 bg-gradient-to-b from-card/40 via-background to-background" />
      <div className="absolute left-1/4 top-40 -z-10 h-[500px] w-[500px] rounded-full bg-primary/[0.04] blur-[120px]" />

      {hasTestimonials ? (
        <div className="mx-auto max-w-5xl">
          <div className="mx-auto max-w-2xl text-center">
            <h2 className="text-3xl font-bold leading-tight tracking-tight text-foreground md:text-4xl">
              What people built with it
            </h2>
            <p className="mt-4 text-muted-foreground leading-relaxed">
              Every one of these was submitted by the person quoted, with their permission, and
              links back to them so you can check they exist.
            </p>
          </div>

          <ul className="mt-12 grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            {TESTIMONIALS.map((t) => (
              <li
                key={t.name + t.quote.slice(0, 24)}
                className="flex flex-col rounded-xl border border-border/60 bg-card/50 p-6"
              >
                <Quote className="h-5 w-5 shrink-0 text-primary/40" aria-hidden="true" />
                <blockquote className="mt-3 flex-1 text-sm leading-relaxed text-foreground">
                  {t.quote}
                </blockquote>

                <figcaption className="mt-5 flex items-center gap-3 border-t border-border/40 pt-4">
                  {/* The photo is required by the issue template and committed
                      to the repo, so there is no missing-image branch to design
                      around here. */}
                  <Image
                    src={t.photo}
                    alt=""
                    width={40}
                    height={40}
                    className="h-10 w-10 shrink-0 rounded-full object-cover"
                  />
                  <div className="min-w-0">
                    <div className="truncate text-sm font-semibold text-foreground">
                      {t.link ? (
                        <Link
                          href={t.link}
                          target="_blank"
                          rel="noreferrer"
                          className="hover:text-primary"
                        >
                          {t.name}
                        </Link>
                      ) : (
                        t.name
                      )}
                    </div>
                    {t.role && (
                      <div className="truncate text-xs text-muted-foreground">{t.role}</div>
                    )}
                  </div>
                </figcaption>
              </li>
            ))}
          </ul>

          <div className="mt-10 text-center">
            <Link
              href={TESTIMONIAL_ISSUE_URL}
              target="_blank"
              rel="noreferrer"
              className="inline-flex h-11 items-center gap-2 rounded-lg border border-border/60 px-6 text-sm font-semibold text-foreground transition-colors hover:bg-accent/50"
            >
              <Github className="h-4 w-4" aria-hidden="true" />
              Add yours
            </Link>
          </div>
        </div>
      ) : (
        <div className="mx-auto max-w-3xl text-center">
          <h2 className="text-3xl font-bold leading-tight tracking-tight text-foreground md:text-5xl">
            Built something with Grit?
          </h2>
          <p className="mt-5 text-lg leading-relaxed text-muted-foreground">
            There are no testimonials here yet, and there will not be any invented ones. If Grit
            is running something of yours in production, tell us about it, the good and the parts
            that hurt, and it goes on this page with your name, your photo and a link back to you.
          </p>
          <p className="mt-4 text-sm leading-relaxed text-muted-foreground/80">
            Critical ones get published too. A page of uniform praise tells a reader nothing.
          </p>
          <div className="mt-9 flex flex-col items-center justify-center gap-3 sm:flex-row">
            <Link
              href={TESTIMONIAL_ISSUE_URL}
              target="_blank"
              rel="noreferrer"
              className="inline-flex h-11 items-center gap-2 rounded-lg bg-primary px-6 text-sm font-semibold text-primary-foreground transition-colors hover:bg-primary/90"
            >
              <Github className="h-4 w-4" aria-hidden="true" />
              Share your experience
            </Link>
            <Link
              href="/showcase"
              className="inline-flex h-11 items-center gap-2 rounded-lg border border-border/60 px-6 text-sm font-semibold text-foreground transition-colors hover:bg-accent/50"
            >
              See what people are building
              <ArrowRight className="h-4 w-4" aria-hidden="true" />
            </Link>
          </div>
        </div>
      )}
    </section>
  )
}
