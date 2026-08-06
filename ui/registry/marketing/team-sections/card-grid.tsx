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

/* Demo photography from Unsplash, free for commercial use with no attribution
   required. Every URL was checked to return 200. The width, height and crop are
   in the query string so the browser fetches roughly what it paints rather than
   a 4000px original the layout then scales down.

   These are people and places that have nothing to do with your company.
   Replace them: every block here takes its content as a prop for that reason. */
const MEMBERS: TeamMember[] = [
  { name: 'Alexander Chee', role: 'Co-Founder, CEO', photo: 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=480&h=480&fit=crop&crop=faces&q=80' },
  { name: 'Sarah Johnson', role: 'Co-Founder, CTO', photo: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=480&h=480&fit=crop&crop=faces&q=80' },
  { name: 'Michael Chen', role: 'Head of Engineering', photo: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=480&h=480&fit=crop&crop=faces&q=80' },
  { name: 'Emily Rodriguez', role: 'Head of Design', photo: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=480&h=480&fit=crop&crop=faces&q=80' },
  { name: 'David Kim', role: 'Head of Product', photo: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=480&h=480&fit=crop&crop=faces&q=80' },
  { name: 'Lisa Wang', role: 'Head of Marketing', photo: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=480&h=480&fit=crop&crop=faces&q=80' },
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
