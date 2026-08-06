import { FACES, portraitTall } from './photos'

/*
 * Full-bleed portraits with the name and role sitting on the image.
 *
 * The overlay is a gradient rather than a flat scrim: a flat panel needs to be
 * dark enough for the worst photograph anyone will ever pass in, which means it
 * is too dark for all the others. A gradient from transparent to near-black
 * only darkens the strip the text occupies.
 *
 * Demo photography is from Unsplash. Swap `members` for your own team; a
 * member with no `photo` falls back to a monogram.
 */

export interface TeamMember {
  name: string
  role: string
  photo?: string
}

const MEMBERS: TeamMember[] = [
  { name: 'Sarah Mitchell', role: 'Co-Founder, CEO', photo: portraitTall(FACES[0]) },
  { name: 'Marcus Chen', role: 'Co-Founder, CTO', photo: portraitTall(FACES[2]) },
  { name: 'David Thompson', role: 'Chief Operating Officer', photo: portraitTall(FACES[5]) },
  { name: 'Emily Watson', role: 'Head of Product', photo: portraitTall(FACES[6]) },
  { name: 'James Rodriguez', role: 'VP of Engineering', photo: portraitTall(FACES[9]) },
  { name: 'Rachel Kim', role: 'Chief Marketing Officer', photo: portraitTall(FACES[4]) },
  { name: 'Michael Foster', role: 'Chief Financial Officer', photo: portraitTall(FACES[3]) },
  { name: 'Amanda Patel', role: 'Head of Design', photo: portraitTall(FACES[8]) },
]

const TINTS = [
  'from-violet-500 to-sky-500',
  'from-sky-500 to-emerald-500',
  'from-amber-500 to-rose-500',
  'from-rose-500 to-violet-500',
  'from-emerald-500 to-teal-500',
  'from-indigo-500 to-fuchsia-500',
]

function tintFor(name: string) {
  let sum = 0
  for (let i = 0; i < name.length; i++) sum += name.charCodeAt(i)
  return TINTS[sum % TINTS.length]
}

function initials(name: string) {
  return name
    .split(' ')
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() ?? '')
    .join('')
}

export default function TeamOverlayGrid({
  title = 'Our team',
  subtitle = 'A diverse group of passionate individuals working together to build something great.',
  members = MEMBERS,
}: {
  title?: string
  subtitle?: string
  members?: TeamMember[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <div className="max-w-2xl">
          <h2 className="text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
            {title}
          </h2>
          <p className="mt-6 text-lg/8 text-gray-600 dark:text-gray-400">{subtitle}</p>
        </div>

        <ul
          role="list"
          className="mt-16 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4"
        >
          {members.map((member) => (
            <li
              key={member.name}
              className="group relative isolate overflow-hidden rounded-2xl"
            >
              {member.photo ? (
                <img
                  src={member.photo}
                  alt=""
                  className="aspect-[4/5] w-full object-cover"
                />
              ) : (
                <div
                  aria-hidden="true"
                  className={`flex aspect-[4/5] w-full items-center justify-center bg-gradient-to-br ${tintFor(
                    member.name,
                  )}`}
                >
                  <span className="text-6xl font-bold text-white/90">
                    {initials(member.name)}
                  </span>
                </div>
              )}

              {/* Only the bottom third is darkened, which keeps the face lit. */}
              <div
                aria-hidden="true"
                className="absolute inset-x-0 bottom-0 h-2/5 bg-gradient-to-t from-black/80 via-black/40 to-transparent"
              />

              <div className="absolute inset-x-0 bottom-0 p-4">
                <h3 className="text-sm font-semibold text-white">{member.name}</h3>
                <p className="mt-0.5 text-sm text-white/70">{member.role}</p>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
