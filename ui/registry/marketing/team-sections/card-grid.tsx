import { FACES, portrait } from './photos'

/*
 * A team grid: portrait, name, role. The workhorse layout.
 *
 * Demo photography comes from Unsplash and is of people who do not work at your
 * company. Replace `members` with your own; the monogram fallback below is what
 * renders for anyone whose `photo` you leave out, so a half-filled directory
 * still looks deliberate rather than broken.
 */

export interface TeamMember {
  name: string
  role: string
  /** A real photograph. Omit it and a monogram is drawn instead. */
  photo?: string
}

const MEMBERS: TeamMember[] = [
  { name: 'Alexander Chee', role: 'Co-Founder, CEO', photo: portrait(FACES[1]) },
  { name: 'Sarah Johnson', role: 'Co-Founder, CTO', photo: portrait(FACES[0]) },
  { name: 'Michael Chen', role: 'Head of Engineering', photo: portrait(FACES[3]) },
  { name: 'Emily Rodriguez', role: 'Head of Design', photo: portrait(FACES[2]) },
  { name: 'David Kim', role: 'Head of Product', photo: portrait(FACES[7]) },
  { name: 'Lisa Wang', role: 'Head of Marketing', photo: portrait(FACES[4]) },
]

/* Deterministic tint per person, so a monogram is not a grey box and the same
   name is always the same colour between renders and between machines. */
const TINTS = [
  'from-violet-400 to-sky-400',
  'from-sky-400 to-emerald-400',
  'from-amber-400 to-rose-400',
  'from-rose-400 to-violet-400',
  'from-emerald-400 to-teal-400',
  'from-indigo-400 to-fuchsia-400',
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

export default function TeamCardGrid({
  title = 'Our incredible leadership team',
  members = MEMBERS,
}: {
  title?: string
  members?: TeamMember[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <h2 className="max-w-md text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
          {title}
        </h2>

        <ul
          role="list"
          className="mt-16 grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3"
        >
          {members.map((member) => (
            <li
              key={member.name}
              className="rounded-2xl border border-gray-200 bg-white p-3 shadow-sm dark:border-white/10 dark:bg-white/5"
            >
              {member.photo ? (
                <img
                  src={member.photo}
                  // The name is already in the caption below, so repeating it
                  // here makes a screen reader say it twice.
                  alt=""
                  className="aspect-square w-full rounded-xl object-cover"
                />
              ) : (
                <div
                  aria-hidden="true"
                  className={`flex aspect-square w-full items-center justify-center rounded-xl bg-gradient-to-br ${tintFor(
                    member.name,
                  )}`}
                >
                  <span className="text-5xl font-bold text-white">
                    {initials(member.name)}
                  </span>
                </div>
              )}

              <div className="px-2 pt-4 pb-2">
                <h3 className="text-base font-semibold text-gray-900 dark:text-white">
                  {member.name}
                </h3>
                <p className="mt-0.5 text-sm text-gray-500 dark:text-gray-400">
                  {member.role}
                </p>
              </div>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
