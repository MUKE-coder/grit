import { CapabilityMatrix } from '@/components/capability-matrix'
import { readyCount, TOTAL_ROWS } from '@/config/capabilities'

/*
 * Sits directly under the benchmark. The two answer different halves of the
 * same question: the benchmark says how fast the request path is, this says how
 * much of the product you did not have to assemble to get one.
 */

export function HomepageCapabilities() {
  return (
    <section className="py-16 md:py-20 px-6 border-b border-border/40">
      <div className="max-w-5xl mx-auto">
        <div className="text-center mb-8">
          <h2 className="text-2xl md:text-3xl font-bold tracking-tight text-foreground mb-3">
            Speed is the easy half
          </h2>
          <p className="text-sm md:text-base text-muted-foreground max-w-2xl mx-auto leading-relaxed">
            The benchmark above measures a request. This measures a Monday: how many of the{' '}
            {TOTAL_ROWS} things a real product needs are already running the first time you open
            the app. Grit has {readyCount('grit')}. The {TOTAL_ROWS - readyCount('grit')} it does
            not have are in the table too, because a comparison where one column wins every row is
            not a comparison.
          </p>
        </div>

        <CapabilityMatrix />
      </div>
    </section>
  )
}
