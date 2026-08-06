/*
 * Demo imagery for the team blocks.
 *
 * Unsplash, which is free for commercial use with no attribution required, and
 * every URL here was checked to return 200 rather than guessed at. The size and
 * crop are baked into the query string so the browser fetches roughly what it
 * paints instead of a 4000px original scaled down in the layout.
 *
 * These are DEMO photographs of people who do not work at your company. They
 * exist so the preview looks like the real thing; replace them with your own
 * team before shipping. Every block takes its data as a prop for exactly that
 * reason, so swapping them is a value change and not a search and replace.
 */

const P = 'https://images.unsplash.com'

/** Square portrait, sized for a grid cell or an avatar. */
export const portrait = (id: string, size = 480) =>
  `${P}/${id}?w=${size}&h=${size}&fit=crop&crop=faces&q=80`

/** Taller crop, for full-bleed cards where the name sits over the image. */
export const portraitTall = (id: string, w = 480) =>
  `${P}/${id}?w=${w}&h=${Math.round(w * 1.25)}&fit=crop&crop=faces&q=80`

/** Landscape, for collages and polaroids. */
export const scene = (id: string, w = 800) =>
  `${P}/${id}?w=${w}&h=${Math.round(w * 0.75)}&fit=crop&q=80`

/** Portrait ids, verified. */
export const FACES = [
  'photo-1494790108377-be9c29b29330',
  'photo-1472099645785-5658abf4ff4e',
  'photo-1534528741775-53994a69daeb',
  'photo-1506794778202-cad84cf45f1d',
  'photo-1544005313-94ddf0286df2',
  'photo-1500648767791-00dcc994a43e',
  'photo-1438761681033-6461ffad8d80',
  'photo-1507003211169-0a1dd7228f2d',
  'photo-1517841905240-472988babdf9',
  'photo-1520813792240-56fc4a3765a7',
  'photo-1487412720507-e7ab37603c6f',
  'photo-1519244703995-f4e0f30006d5',
] as const

/** Team and workspace ids, verified. */
export const SCENES = [
  'photo-1522071820081-009f0129c71c',
  'photo-1600880292203-757bb62b4baf',
  'photo-1552664730-d307ca884978',
  'photo-1531482615713-2afd69097998',
  'photo-1521737604893-d14cc237f11d',
  'photo-1517048676732-d65bc937f952',
  'photo-1542744173-8e7e53415bb0',
  'photo-1497366216548-37526070297c',
] as const
