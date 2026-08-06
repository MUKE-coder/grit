import { FACES, portrait } from './photos'

/*
 * Circular portraits, centred, name and role beneath.
 *
 * The classic leadership layout. A circle crops to the face and away from
 * whatever the background happened to be, which is why it survives a set of
 * photographs taken by eight different people in eight different rooms. It is
 * the most forgiving of these blocks if your team photos are not a matched set.
 *
 * Demo photography is from Unsplash. Replace `members` with your own team.
 */

export interface TeamMember {
  name: string
  role: string
  photo?: string
}

const MEMBERS: TeamMember[] = [
  { name: 'Sarah Mitchell', role: 'Co-Founder, CEO', photo: portrait(FACES[0], 256) },
  { name: 'Marcus Chen', role: 'Co-Founder, CTO', photo: portrait(FACES[2], 256) },
  { name: 'David Thompson', role: 'Chief Operating Officer', photo: portrait(FACES[5], 256) },
  { name: 'James Rodriguez', role: 'VP of Engineering', photo: portrait(FACES[9], 256) },
  { name: 'Emily Watson', role: 'Head of Product', photo: portrait(FACES[6], 256) },
  { name: 'Rachel Kim', role: 'Chief Marketing Officer', photo: portrait(FACES[4], 256) },
  { name: 'Michael Foster', role: 'Chief Financial Officer', photo: portrait(FACES[3], 256) },
  { name: 'Amanda Patel', role: 'Head of Design', photo: portrait(FACES[8], 256) },
]

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

export default function TeamCircularCentered({
  title = 'Our incredible leadership team',
  members = MEMBERS,
}: {
  title?: string
  members?: TeamMember[]
}) {
  return (
    <section className="bg-white py-24 sm:py-32 dark:bg-gray-950">
      <div className="mx-auto max-w-7xl px-6 lg:px-8">
        <h2 className="mx-auto max-w-2xl text-center text-4xl font-semibold tracking-tight text-balance text-gray-900 sm:text-5xl dark:text-white">
          {title}
        </h2>

        <ul
          role="list"
          className="mx-auto mt-20 grid max-w-2xl grid-cols-2 gap-x-8 gap-y-16 lg:max-w-none lg:grid-cols-4"
        >
          {members.map((member) => (
            <li key={member.name} className="text-center">
              {member.photo ? (
                <img
                  src={member.photo}
                  alt=""
                  className="mx-auto size-32 rounded-full object-cover"
                />
              ) : (
                <div
                  aria-hidden="true"
                  className={`mx-auto flex size-32 items-center justify-center rounded-full bg-gradient-to-br ${tintFor(
                    member.name,
                  )}`}
                >
                  <span className="text-3xl font-bold text-white">
                    {initials(member.name)}
                  </span>
                </div>
              )}

              <h3 className="mt-6 text-base font-semibold text-gray-900 dark:text-white">
                {member.name}
              </h3>
              <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{member.role}</p>
            </li>
          ))}
        </ul>
      </div>
    </section>
  )
}
