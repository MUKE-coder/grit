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

/* Demo photography from Unsplash, free for commercial use with no attribution
   required. Every URL was checked to return 200. The width, height and crop are
   in the query string so the browser fetches roughly what it paints rather than
   a 4000px original the layout then scales down.

   These are people and places that have nothing to do with your company.
   Replace them: every block here takes its content as a prop for that reason. */
const MEMBERS: TeamMember[] = [
  { name: 'Sarah Mitchell', role: 'Co-Founder, CEO', photo: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=480&h=600&fit=crop&crop=faces&q=80' },
  { name: 'Marcus Chen', role: 'Co-Founder, CTO', photo: 'https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=480&h=600&fit=crop&crop=faces&q=80' },
  { name: 'David Thompson', role: 'Chief Operating Officer', photo: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=480&h=600&fit=crop&crop=faces&q=80' },
  { name: 'Emily Watson', role: 'Head of Product', photo: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=480&h=600&fit=crop&crop=faces&q=80' },
  { name: 'James Rodriguez', role: 'VP of Engineering', photo: 'https://images.unsplash.com/photo-1520813792240-56fc4a3765a7?w=480&h=600&fit=crop&crop=faces&q=80' },
  { name: 'Rachel Kim', role: 'Chief Marketing Officer', photo: 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?w=480&h=600&fit=crop&crop=faces&q=80' },
  { name: 'Michael Foster', role: 'Chief Financial Officer', photo: 'https://images.unsplash.com/photo-1506794778202-cad84cf45f1d?w=480&h=600&fit=crop&crop=faces&q=80' },
  { name: 'Amanda Patel', role: 'Head of Design', photo: 'https://images.unsplash.com/photo-1517841905240-472988babdf9?w=480&h=600&fit=crop&crop=faces&q=80' },
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
