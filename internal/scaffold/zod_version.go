package scaffold

// zodVersion is the zod every generated package.json pins.
//
// One constant because there were eight, at four different 3.x minors, in
// eight files. A dependency the generator writes into every frontend is one
// fact, and eight copies of it is how a project ends up with two zods in a
// hoisted node_modules (grit#175).
//
// # Why 4
//
// zod 3 has a HIGH advisory with no fix inside 3.x: GHSA-mrhx-c7xc-m3w5,
// memory exhaustion while validating oversized input. It was the last
// outstanding advisory in a freshly generated project, and the only patched
// version is 4.
//
// # What the move cost, which was less than expected
//
// Checked against zod 4.6.7 rather than assumed: every 3.x spelling the
// generator emitted still type-checks AND still validates at runtime.
// z.string().email(), z.string().url() and z.string().uuid("message") are
// deprecated in favour of z.email(), z.url() and z.uuid(), and they work.
//
// So this was a pin bump, and the call sites were modernised because they
// should be rather than because they had to be. The same is true of a
// project's own hand-written schemas: they keep working, which is why
// grit upgrade does not need to scan for them or refuse to move the pin.
//
// The input-vs-output divergence around .default() is real in zod 4 and does
// not reach Grit's own forms: LoginSchema, RegisterSchema and
// ForgotPasswordSchema declare no defaults, so z.infer and z.input are the
// same type for every schema handed to useForm.
const zodVersion = "^4.6.7"
