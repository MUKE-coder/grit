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

/* Demo photography from Unsplash, free for commercial use with no attribution
   required. Every URL was checked to return 200. The width, height and crop are
   in the query string so the browser fetches roughly what it paints rather than
   a 4000px original the layout then scales down.

   These are people and places that have nothing to do with your company.
   Replace them: every block here takes its content as a prop for that reason. */
const MEMBERS: TeamMember[] = [
  { name: 'Sarah Mitchell', role: 'Co-Founder, CEO', photo: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=256&h=256&fit=crop&crop=faces&q=80' },
  { name: 'Marcus Chen', role: 'Co-Founder, CTO', photo: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=256&h=256&fit=crop&crop=faces&q=80' },
  { name: 'David Thompson', role: 'Chief Operating Officer', photo: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=256&h=256&fit=crop&crop=faces&q=80' },
  { name: 'James Rodriguez', role: 'VP of Engineering', photo: 'https://images.unsplash.com/photo-1520813792240-56fc4a3765a7?w=256&h=256&fit=crop&crop=faces&q=80' },
  { name: 'Emily Watson', role: 'Head of Product', photo: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=256&h=256&fit=crop&crop=faces&q=80' },
  { name: 'Rachel Kim', role: 'Chief Marketing Officer', photo: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=256&h=256&fit=crop&crop=faces&q=80' },
  { name: 'Michael Foster', role: 'Chief Financial Officer', photo: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=256&h=256&fit=crop&crop=faces&q=80' },
  { name: 'Amanda Patel', role: 'Head of Design', photo: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?w=256&h=256&fit=crop&crop=faces&q=80' },
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
