# Rehearsal: resource-definition-stage

The demo in the landing page's **Make it yours** section: a resource
definition above the admin screen it produces, and a column added to one
arriving in the other.

## Gold

`stacked-cards` (two components at one width, a small gap, the lower one off
the bottom).

**Kept:** the geometry, one width for both cards, the 12px gap, the lower card
bleeding off the bottom, the gutter measured on the upper card's sides.

**Changed:** the pair. The gold is a plan over its usage; this is a file over
the screen it describes, so the upper card is code and the lower one is a
table. The look is the site's rather than the kit's default: `theme` and
`accent` follow `next-themes`, so the demo is dark with Grit's violet on the
landing page and light if a reader switches.

**Started from `panel-annotation` and abandoned it.** It is the gold for "a
setting and its effect", which is exactly this story, but its satellite
straddles the hero's right edge, and the right edge of an admin window is
table columns. The rule is that a satellite covers chrome, not content, and
no arrangement of a 940, 760 or 620 stage put the panel over margin without
either covering the Due column or straddling by a tenth. Stacking also reads
in the order the sentence does.

## Beats

1. Rest, 1.1s: three columns in the file, three in the table.
2. The cursor moves to the end of the `columns` array.
3. `{ name: 'priority', label: 'Priority' },` opens in the file, carrying the
   two lines that close it; the PRIORITY column opens in the table at the same
   moment, in width the row already had.
4. Payoff holds 3.2s.

Loop about 5.5s, `play="once"`: the kit's name for plays when first
seen, then holds the end. Poster is the payoff.

It shipped as `play="in-view"` in v3.560.0, which is not one of the
kit's three values. It fell through every branch, so the script ran but
the stage never marked itself done, and the demo replayed each time it
scrolled back into view. `tsc` had been saying so; the docs build does
not read it (`ignoreBuildErrors`), which is why a wrong literal reached
a release.

No leader. The gap is 12px, which is no path for one, and the gold's own note
is that where there is no clear path the shared change carries the cause: the
line and the column arrive together and share the accent.

## Captured

- Wide (620 × 500) and compact (358 × 580), both on the running site.
- Contact sheets at 12 and 14 frames through the whole loop.
- Dark, which is the landing page's theme; the light branch is the same
  component with two token values and was not filmed.

## Gate

```
film: inconclusive; watch the motion and inspect every image before approval
rehearsal failed: 2 never-intentional line(s) above
```

Both lines are the same finding, before and after the file grows:

```
[stage-1] uneven gutters: section.rd-file sits 40px from the left of
div.ss-stage, 40px from the right and 221px from the bottom; bleed it off an
edge, or keep one gutter on every side it doesn't bleed
```

**It stands, and here is the evidence.** The gold produces it too:

```
$ node scripts/film.mjs "...?capture=stacked-cards" --gate
[stage-1] uneven gutters: section.sk-card sits 90px from the left of
div.ss-stage, 90px from the right and 225px from the bottom; ...
rehearsal failed: 1 never-intentional line(s) above
```

The upper card of a stacked pair cannot bleed, and its distance to the stage
floor is the lower card's business. Satisfying the line would mean leaving the
layout. This is parity with the gold, not a pass.

## Film: every line, the frame, the fix

| Finding | Frame | Outcome |
|---|---|---|
| Panel over the Due column | first still | Layout changed to stacked cards |
| Void under a 5-row table in a 500px window | first still | The lower card bleeds; the list runs off the crop |
| Code lines flowing two to a row | first still | `.rd-ln { display: block }` |
| Status pills invisible: two near-white tones on white | rd-04 | Tones from `--ss-text`, not `--ss-bar-*` |
| Half the file card empty at 480 wide | rd-04 | Real `label:` on every column and a `filter:` on one, which is what a definition carries |
| The screen card covering the file's closing lines | rd-07 | Screen pinned below the file's *opened* height |
| Compact: the status line sliced by the card edge | stage-390 | The filter clause is wide-only |
| Compact: the same overlap | stage-390 | Compact screen top 252 → 296 |
| Uneven gutters | every frame | Stands: the gold does the same |

## Assumptions

- The claim is honest: adding an entry to a resource definition's `columns`
  array is what puts a column in the generated table. The demo illustrates a
  real mechanism, and the words on screen are a plausible app's, not a real
  customer's.
- `admin.fieldkit.co` and the two task titles are invented. Nothing here is
  presented as a screenshot of a real deployment, which is why the sections
  that *are* real screenshots (platforms, auth, batteries, deploy) were left
  alone.

## Unresolved

- The light-theme branch is wired and unfilmed.
- The hover and reduced-motion states were not captured; `play="in-view"`
  plays once and holds the payoff, so reduced motion shows the poster, which
  is that same frame.

## Not checked

- Safari and Firefox. Filmed in Chrome only.
- The compact layout between 390 and 580, where the stage is compact but the
  column is wider than 358.
