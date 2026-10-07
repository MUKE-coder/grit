package scaffold

// The Next.js every scaffolded frontend pins.
//
// Exact, like React and pnpm, and for the same reason: a generated project
// installs what this version of Grit was tested against, not whatever was
// published this morning.
//
// It used to be "^16.1.6", which is any 16.x, and that is not a theoretical
// risk. On 2026-10-07 a project reported a blank page and this in its terminal:
//
//	Uncaught TypeError: Cannot read properties of undefined (reading 'slots')
//	    at OuterLayoutRouter
//
// It had installed 16.4.0, published 2026-10-06 at 18:21, roughly eighteen
// hours earlier. That release rewrote the App Router's layout-router to keep a
// render tree separate from the route tree, and it reads
//
//	parentRenderTree.slots?.get(parallelRouterKey)
//
// guarding `slots` and not `parentRenderTree`. When the client's tree does not
// match the server's, the read is of undefined and it throws from a stack with
// no application frames in it. 16.3.8 has no parentRenderTree at all: the code
// that throws does not exist there.
//
// pnpm's own supply-chain policy rejects 16.4.0 today, because it is inside the
// minimumReleaseAge cutoff. A project that pinned the newest release would fail
// its own install, which is a second reason not to chase the latest.
//
// Pinning does not fix a bug in Next. It makes the version a decision, so a
// move happens when somebody runs the tiers against it rather than when npm
// publishes.
//
// To move it: change this, run the tier tutorials end to end, check the new
// version is past the release-age cutoff, and ship it as its own release so a
// project can date the change.
const nextVersion = "16.3.8"
