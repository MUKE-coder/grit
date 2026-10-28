import {
  DockerLogo,
  GoLogo,
  NextLogo,
  PostgresLogo,
  ReactLogo,
  RedisLogo,
  TailwindLogo,
  TanStackLogo,
  TypeScriptLogo,
} from '@/components/framework-logos'

/**
 * The strip under the hero: the tools a Grit project is actually made of.
 *
 * Only logos the repo already draws. A marquee that mixes real marks with
 * name-only entries reads as a list somebody padded, which is the opposite of
 * the point: every one of these is a dependency you can go and look at.
 *
 * The list is rendered twice because the animation translates the track by
 * -50%; the second copy is what the first scrolls into, and it is hidden from
 * assistive tech so the names are not read out twice.
 */
const TOOLS = [
  { name: 'Go', Logo: GoLogo },
  { name: 'PostgreSQL', Logo: PostgresLogo },
  { name: 'Redis', Logo: RedisLogo },
  { name: 'Docker', Logo: DockerLogo },
  { name: 'Next.js', Logo: NextLogo },
  { name: 'React', Logo: ReactLogo },
  { name: 'TanStack', Logo: TanStackLogo },
  { name: 'TypeScript', Logo: TypeScriptLogo },
  { name: 'Tailwind', Logo: TailwindLogo },
]

function Row({ hidden = false }: { hidden?: boolean }) {
  return (
    <>
      {TOOLS.map((tool) => (
        <span
          key={tool.name}
          className="g-logo inline-flex items-center gap-2.5 whitespace-nowrap px-7 text-[15px] font-medium text-fg"
          aria-hidden={hidden || undefined}
        >
          <tool.Logo className="h-[22px] w-[22px]" />
          {tool.name}
        </span>
      ))}
    </>
  )
}

export function ToolMarquee() {
  return (
    <div className="relative border-t border-border">
      <div className="g-frame py-7">
        <span aria-hidden="true" className="g-cnr g-cnr--tl" />
        <span aria-hidden="true" className="g-cnr g-cnr--tr" />
        <p className="g-eyebrow mb-[18px] text-center">Built on tools you already trust</p>
        <div className="g-marquee-wrap">
          <div className="g-marquee">
            <Row />
            <Row hidden />
          </div>
        </div>
      </div>
    </div>
  )
}
