import Link from 'next/link'
import { ArrowRight } from 'lucide-react'

import { Button, ButtonArrow } from '@/components/grit/button'
import { BlueprintShots } from '@/components/blueprints/shots'
import { BLUEPRINTS } from '@/config/blueprints'
import { PLANS, usd } from '@/config/blueprint-plans'

/*
 * The blueprints, on the home page.
 *
 * A compact band rather than a second gallery: three screenshots each and the
 * one line that says what the thing is, then out to /blueprints. Somebody on
 * the home page is deciding whether this framework is worth their evening,
 * and "here are whole applications built with it" is the strongest answer
 * available, which is why it sits just before the showcase rather than after
 * the FAQ where nobody reaches.
 */

export function BlueprintsBand() {
  const built = BLUEPRINTS.filter((b) => b.access !== 'planned')
  const planned = BLUEPRINTS.length - built.length

  return (
    <>
      <ul className="mt-12 grid list-none grid-cols-1 gap-4 p-0 md:grid-cols-3">
        {built.map((blueprint) => (
          <li key={blueprint.name}>
            <Link
              href={`/blueprints/${blueprint.name}`}
              className="g-card g-card--interactive h-full gap-0 p-0"
            >
              <BlueprintShots
                blueprint={blueprint}
                className="rounded-b-none border-0 border-b"
              />
              <div className="flex flex-1 flex-col gap-2 p-5">
                <h3 className="g-card__title">{blueprint.title}</h3>
                <p className="g-card__body">{blueprint.summary}</p>
                <span className="g-card__more mt-auto pt-4 text-body-sm text-accent">
                  Look at it
                  <ArrowRight className="ml-1.5 inline h-3.5 w-3.5" aria-hidden="true" />
                </span>
              </div>
            </Link>
          </li>
        ))}
      </ul>

      <div className="mt-10 flex flex-col items-center justify-center gap-4 text-center">
        <p className="m-0 max-w-[var(--measure)] text-body-sm text-fg-muted">
          {planned} more are on the list and named rather than hinted at. WhatsApp is
          open source; the rest are {usd(PLANS.single.priceUsd)} each, or{' '}
          {usd(PLANS.monthly.priceUsd)} a month for every one of them.
        </p>
        <Button href="/blueprints" variant="secondary">
          The whole catalogue
          <ButtonArrow />
        </Button>
      </div>
    </>
  )
}
