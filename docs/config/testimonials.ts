import data from '@/data/testimonials.json'

/**
 * Approved testimonials.
 *
 * The pipeline, and the reason it has one:
 *
 *   1. Somebody files the "Share your experience" issue. The template requires
 *      a quote, a name and a photo, and a checkbox giving permission to publish
 *      all three.
 *   2. A maintainer reads it and, if it is real, copies the fields into
 *      data/testimonials.json and commits the photo under public/testimonials/.
 *   3. It appears on the homepage on the next deploy.
 *
 * Nothing is rendered straight from an issue. An issue is a text box on the
 * internet, and this page carries named human beings vouching for the project:
 * the gap between those two things is the approval step.
 *
 * This section previously shipped with quotes attributed to real people and
 * institutions who never gave them, placeholder copy from a design mock. That
 * is why the empty state here is deliberate and why the array starts empty. For
 * a framework asking to be trusted with auth and audit logs, being caught
 * inventing an endorsement costs more than an empty section ever could.
 */

export interface Testimonial {
  /** The line that appears on the card. Required by the issue template. */
  quote: string
  /** Exactly as they asked to be credited. */
  name: string
  /** Optional: "Engineering lead, Acme Logistics". */
  role?: string
  /**
   * Path under /public, e.g. "/testimonials/amina.jpg".
   *
   * Required, and committed to the repo rather than hotlinked from the issue.
   * GitHub's user-content URLs are not a CDN and have gone dead before; a
   * broken face on a testimonial reads worse than no testimonial.
   */
  photo: string
  /** Their site, GitHub or LinkedIn, so a reader can check they exist. */
  link?: string
  /** The issue this came from, so the permission trail is one click away. */
  source?: string
}

export const TESTIMONIALS: Testimonial[] = data.approved as Testimonial[]

export const hasTestimonials = TESTIMONIALS.length > 0

/** Where the "share yours" buttons point. */
export const TESTIMONIAL_ISSUE_URL =
  'https://github.com/MUKE-coder/grit/issues/new?template=testimonial.yml'
