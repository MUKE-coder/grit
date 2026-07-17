import { redirect } from 'next/navigation'

/**
 * /donate is superseded by /sponsor, which offers the same one-time amounts plus
 * recurring tiers. Kept as a redirect because the URL is in the wild (README,
 * older posts, the header CTA). /donate/success is untouched — Stripe returns
 * there after checkout.
 */
export default function DonatePage() {
  redirect('/sponsor')
}
