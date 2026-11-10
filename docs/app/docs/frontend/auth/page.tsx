import Link from 'next/link'
import { ArrowRight, ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { DocsShell } from '@/components/docs/shell'
import { CodeBlock } from '@/components/code-block'
import { getDocMetadata } from '@/config/docs-metadata'

export const metadata = getDocMetadata('/docs/frontend/auth')

export default function WebAuthPage() {
  return (
    <DocsShell>
      {/* Header */}
      <div className="mb-10">
        <span className="tag-mono text-accent/80 mb-3 block">Frontend</span>
        <h1 className="text-4xl font-bold tracking-tight mb-4">
          Auth and protected pages
        </h1>
        <p className="text-lg text-muted-foreground leading-relaxed">
          The API already does authentication. What the web app needs is a
          sign-in page, a sign-up page, a hook that knows who is signed in, and
          a gate on the pages that are private. Four files, and one idea that is
          worth getting right: the gate in front of a page is not the thing that
          protects it.
        </p>
      </div>

      <div className="prose-grit">
        <h2 id="what-you-already-have">What the API already does</h2>
        <p>
          Everything on the server side ships with the project. The endpoints
          exist, the password hashing exists, the sessions table exists, the
          refresh rotation and replay detection exist. There is nothing to
          install:
        </p>
        <div className="rounded-lg border border-border/30 bg-card/30 overflow-hidden mb-6">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-border/30 bg-accent/20">
                <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Endpoint</th>
                <th className="text-left px-4 py-2.5 font-medium text-foreground/80">Does</th>
              </tr>
            </thead>
            <tbody>
              <tr className="border-b border-border/20">
                <td className="px-4 py-2.5 font-mono text-xs">POST /api/auth/register</td>
                <td className="px-4 py-2.5 text-muted-foreground">Creates the account and signs it in</td>
              </tr>
              <tr className="border-b border-border/20">
                <td className="px-4 py-2.5 font-mono text-xs">POST /api/auth/login</td>
                <td className="px-4 py-2.5 text-muted-foreground">Signs in, or asks for a 2FA code</td>
              </tr>
              <tr className="border-b border-border/20">
                <td className="px-4 py-2.5 font-mono text-xs">GET /api/auth/me</td>
                <td className="px-4 py-2.5 text-muted-foreground">Who this is, or 401</td>
              </tr>
              <tr className="border-b border-border/20">
                <td className="px-4 py-2.5 font-mono text-xs">POST /api/auth/refresh</td>
                <td className="px-4 py-2.5 text-muted-foreground">A new access token from the refresh cookie</td>
              </tr>
              <tr>
                <td className="px-4 py-2.5 font-mono text-xs">POST /api/auth/logout</td>
                <td className="px-4 py-2.5 text-muted-foreground">Clears the cookies, ends the session row</td>
              </tr>
            </tbody>
          </table>
        </div>
        <p>
          A successful sign-in or sign-up answers with cookies, and the shape of
          them is the thing to understand before writing any of the frontend:
        </p>
        <CodeBlock
          language="text"
          code={`HTTP/1.1 201 Created
Set-Cookie: grit_access=...;    Path=/;             Max-Age=900;    HttpOnly; SameSite=Lax
Set-Cookie: grit_refresh=...;   Path=/api/v1/auth;  Max-Age=604800; HttpOnly; SameSite=Lax
Set-Cookie: grit_signed_in=1;   Path=/;             Max-Age=604800; HttpOnly; SameSite=Lax`}
        />
        <p>
          All three are <code>HttpOnly</code>, which means JavaScript cannot read
          them, which is the point. The access token lasts fifteen minutes. The
          refresh token lasts a week and is scoped to the refresh route, so it is
          not sent with every unrelated request.
        </p>
        <p>
          The response body also carries an access token. That is for native
          clients, the Expo app and the desktop app, which have somewhere better
          than a cookie to keep it. <strong>A browser ignores it.</strong>{' '}
          Copying it into <code>localStorage</code> undoes the reason the cookies
          are <code>HttpOnly</code>: anything readable from JavaScript is
          reachable by an XSS bug, a compromised dependency or a browser
          extension.
        </p>

        <h2 id="what-ships">What the web app ships with, and what it does not</h2>
        <p>
          <code>apps/web/lib/api.ts</code> is already there, and it already does
          the hard part. It sends cookies with every request, and when a request
          comes back <code>401</code> it spends the refresh cookie on a new
          access token and sends the request again. Requests that fail together
          share one refresh, because the server rotates refresh tokens and treats
          a second use of the old one as theft.
        </p>
        <p>
          What a triple project&apos;s <code>apps/web</code> does{' '}
          <em>not</em> ship with is any auth UI. That is on purpose: a storefront,
          a marketing site and a customer portal want different sign-in pages,
          and a generated one would be deleted. The admin panel has six of them
          and they are a reasonable thing to read for reference.
        </p>
        <p>So there are four files to write. Here they are.</p>

        <h2 id="marker-cookie">Step 1: a marker cookie, so the edge can tell</h2>
        <p>
          The session cookies are <code>HttpOnly</code> and set on the{' '}
          <em>API&apos;s</em> origin. In development that is{' '}
          <code>localhost:8080</code>, and the browser does send them to{' '}
          <code>localhost:3000</code>. In production, with the site on one domain
          and the API on another, it does not.
        </p>
        <p>
          So the frontend writes a marker of its own, on its own origin, whenever
          the API hands back a session. It is not a credential and it is not
          checked against anything:
        </p>
        <CodeBlock
          language="ts"
          code={`// apps/web/lib/web-session.ts
export const WEB_SESSION_COOKIE = "shop_session";

// Set by the API beside its session cookies, and only when the API shares this
// app's host. The gate accepts it too, so a browser that is already signed in
// is not sent to the sign-in page the first time it loads a page that has just
// gained a gate.
export const API_SIGNED_IN_COOKIE = "grit_signed_in";

const MAX_AGE_SECONDS = 7 * 24 * 60 * 60;   // the refresh window

export function setWebSession(): void {
  if (typeof document === "undefined") return;
  const secure =
    typeof location !== "undefined" && location.protocol === "https:" ? "; Secure" : "";
  document.cookie =
    WEB_SESSION_COOKIE + "=1; Path=/; Max-Age=" + MAX_AGE_SECONDS +
    "; SameSite=Lax" + secure;
}

export function clearWebSession(): void {
  if (typeof document === "undefined") return;
  document.cookie = WEB_SESSION_COOKIE + "=; Path=/; Max-Age=0; SameSite=Lax";
}`}
        />
        <p>
          Anyone can write this cookie from the browser console. All it buys them
          is the HTML of a page whose every request the API refuses, which is the
          correct amount of damage for a cookie that holds no secret. It answers
          one question, &quot;has this browser signed in here?&quot;, and that is
          enough to decide whether to send a page or a redirect.
        </p>

        <h2 id="hooks">Step 2: the hooks</h2>
        <p>
          Four of them. The <code>User</code> type comes from{' '}
          <code>packages/shared</code>, which is the same shape the Go API
          serialises, so it cannot drift from the API the way a hand-written
          interface does.
        </p>
        <CodeBlock
          language="ts"
          code={`// apps/web/hooks/use-auth.ts
"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import type { ApiResponse, AuthResponse, LoginRequest, RegisterRequest, User }
  from "@repo/shared/types";
import { api } from "@/lib/api";
import { clearWebSession, setWebSession } from "@/lib/web-session";

// null rather than a throw on 401, so a guard can read \`user === null\`
// without a try/catch. Anything else is a real failure and surfaces.
export function useMe() {
  return useQuery<User | null>({
    queryKey: ["me"],
    queryFn: async () => {
      try {
        const { data } = await api.get<ApiResponse<User>>("/api/auth/me");
        return data.data;
      } catch (err: unknown) {
        const status = (err as { response?: { status?: number } }).response?.status;
        if (status === 401) {
          clearWebSession();
          return null;
        }
        throw err;
      }
    },
    retry: false,
    staleTime: 10 * 60 * 1000,
  });
}

export function useLogin() {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (credentials: LoginRequest) => {
      const { data } = await api.post<ApiResponse<AuthResponse>>(
        "/api/auth/login", credentials
      );
      return data.data;
    },
    onSuccess: (auth) => {
      setWebSession();
      queryClient.setQueryData(["me"], auth.user);
      router.push("/account");
    },
  });
}

export function useRegister() {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (body: RegisterRequest) => {
      const { data } = await api.post<ApiResponse<AuthResponse>>(
        "/api/auth/register", body
      );
      return data.data;
    },
    onSuccess: (auth) => {
      setWebSession();
      queryClient.setQueryData(["me"], auth.user);
      router.push("/account");
    },
  });
}

export function useLogout() {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      await api.post("/api/auth/logout");
    },
    // onSettled, not onSuccess. A logout that fails on the network still has
    // to end with this browser signed out locally, or the person is left
    // looking at a page they believe is private.
    onSettled: () => {
      clearWebSession();
      queryClient.setQueryData(["me"], null);
      queryClient.clear();
      router.push("/");
    },
  });
}`}
        />
        <p>
          <code>queryClient.clear()</code> on logout rather than just dropping{' '}
          <code>me</code>: every other cached query was fetched as that person,
          and leaving an order list in the cache means the next person to sign in
          on the same browser sees it for a frame.
        </p>

        <h2 id="sign-in">Step 3: the sign-in page</h2>
        <CodeBlock
          language="tsx"
          code={`// apps/web/app/(auth)/login/page.tsx
"use client";

import { useState } from "react";
import { useLogin } from "@/hooks/use-auth";

export default function LoginPage() {
  const login = useLogin();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  // The API's message, not one of ours. It distinguishes a wrong password from
  // a locked account from a rate limit, and replacing all three with "Login
  // failed" throws away the only thing that tells somebody what to do next.
  const error = login.error
    ? ((login.error as { response?: { data?: { error?: { message?: string } } } })
        .response?.data?.error?.message ?? "Could not sign in. Try again.")
    : null;

  return (
    <form
      onSubmit={(e) => {
        e.preventDefault();
        login.mutate({ email, password });
      }}
    >
      <label htmlFor="email">Email</label>
      <input id="email" type="email" autoComplete="email" required
        value={email} onChange={(e) => setEmail(e.target.value)} />

      <label htmlFor="password">Password</label>
      <input id="password" type="password" autoComplete="current-password" required
        value={password} onChange={(e) => setPassword(e.target.value)} />

      {error ? <p role="alert">{error}</p> : null}

      <button type="submit" disabled={login.isPending}>
        {login.isPending ? "Signing in..." : "Sign in"}
      </button>
    </form>
  );
}`}
        />
        <p>
          Three small things that are easy to leave out and annoying to live
          without: <code>htmlFor</code> on every label, so the field has an
          accessible name and a tap on the label focuses it;{' '}
          <code>autoComplete</code>, so a password manager fills it; and{' '}
          <code>role=&quot;alert&quot;</code> on the error, so a screen reader
          announces it instead of leaving somebody waiting.
        </p>

        <h2 id="sign-up">Step 4: the sign-up page</h2>
        <p>
          One thing to get right. The API&apos;s{' '}
          <code>RegisterRequest</code> is <code>first_name</code> and{' '}
          <code>last_name</code>, not <code>name</code>. Posting{' '}
          <code>{'{ name }'}</code> gets this, which is at least explicit about
          it:
        </p>
        <CodeBlock
          language="json"
          code={`{"error":{"code":"VALIDATION_ERROR","message":"Key: 'RegisterRequest.first_name' Error:Field validation for 'first_name' failed on the 'required' tag\\nKey: 'RegisterRequest.last_name' Error:Field validation for 'last_name' failed on the 'required' tag"}}`}
        />
        <p>
          The shape is in <code>packages/shared/types/user.ts</code>, which is
          the file to read rather than guessing:
        </p>
        <CodeBlock
          language="ts"
          code={`export interface RegisterRequest {
  first_name: string;
  last_name: string;
  email: string;
  password: string;
  mac_address?: string;   // optional, for native clients that can get it
}`}
        />
        <CodeBlock
          language="tsx"
          code={`// apps/web/app/(auth)/sign-up/page.tsx
const register = useRegister();
const [form, setForm] = useState({
  first_name: "", last_name: "", email: "", password: "",
});

const set = (key: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement>) =>
  setForm((f) => ({ ...f, [key]: e.target.value }));

// ...then one input per field, and:
register.mutate(form);`}
        />
        <p>
          The project also ships a <code>PasswordStrength</code> component in{' '}
          <code>apps/web/components/password-strength.tsx</code> and the rules it
          measures against in <code>lib/password-rules.ts</code>, which are the
          same rules the API enforces. Use them rather than writing a second set
          that disagrees.
        </p>

        <h2 id="the-gate">Step 5: the gate, which is proxy.ts</h2>
        <p>
          A Next 16 project uses <code>proxy.ts</code> at the app root. Not{' '}
          <code>middleware.ts</code>: that is the old name, and a project that
          has both runs neither in the way you expect.
        </p>
        <CodeBlock
          language="ts"
          code={`// apps/web/proxy.ts
import { NextResponse, type NextRequest } from "next/server";
import { API_SIGNED_IN_COOKIE, WEB_SESSION_COOKIE } from "@/lib/web-session";

// Private by prefix, rather than public by prefix. The storefront is public,
// so listing what is private is the shorter and safer list: a new marketing
// page is reachable by default, and a new private page is not reachable by
// accident.
const PRIVATE_PREFIXES = ["/account", "/orders"];

function isPrivate(pathname: string): boolean {
  return PRIVATE_PREFIXES.some((p) => pathname === p || pathname.startsWith(p + "/"));
}

export function proxy(request: NextRequest) {
  const { pathname } = request.nextUrl;
  if (!isPrivate(pathname)) return NextResponse.next();

  if (
    request.cookies.has(WEB_SESSION_COOKIE) ||
    request.cookies.has(API_SIGNED_IN_COOKIE)
  ) {
    return NextResponse.next();
  }

  const url = request.nextUrl.clone();
  url.pathname = "/login";
  url.search = "?next=" + encodeURIComponent(pathname);   // where they were going
  return NextResponse.redirect(url);
}

export const config = {
  // Everything but Next's own output and anything with a file extension.
  // [.] rather than a backslash: this string is read as JavaScript before it
  // is read as a regex, so a single backslash is eaten and the pattern then
  // matches any character where it means a dot.
  matcher: ["/((?!_next/|[^?]*[.][A-Za-z0-9]+$).*)"],
};
`}
        />
        <p>
          It decides one thing: whether to send the HTML of a private page to
          this browser at all. <strong>It is not authorisation.</strong> It reads
          a cookie that anybody can write, and it never asks the API anything.
          Its whole job is that a signed-out visitor who types{' '}
          <code>/account</code> gets the sign-in page instead of a flash of an
          empty dashboard.
        </p>

        <h2 id="page-guard">Step 6: the half that is authoritative</h2>
        <p>
          The page asks the API, because the API is the only thing that knows
          whether the session is still good. It may have been revoked from
          another device, or expired, or the account may have been deactivated
          since the marker cookie was written.
        </p>
        <CodeBlock
          language="tsx"
          code={`// apps/web/app/account/page.tsx
"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { useLogout, useMe } from "@/hooks/use-auth";

export default function AccountPage() {
  const router = useRouter();
  const { data: user, isLoading } = useMe();
  const logout = useLogout();

  useEffect(() => {
    if (!isLoading && user === null) router.replace("/login?next=/account");
  }, [isLoading, user, router]);

  if (isLoading) {
    return <div className="h-8 w-40 animate-pulse rounded bg-bg-hover" />;
  }

  // The redirect is on its way. Rendering nothing beats rendering an empty
  // shell that looks like the account page with no data in it.
  if (!user) return null;

  return (
    <div className="mx-auto max-w-2xl px-6 py-16">
      <h1 className="text-3xl font-bold tracking-tight">
        {user.first_name} {user.last_name}
      </h1>
      <p className="mt-2 text-text-secondary">{user.email}</p>

      <button onClick={() => logout.mutate()} disabled={logout.isPending}>
        {logout.isPending ? "Signing out..." : "Sign out"}
      </button>
    </div>
  );
}`}
        />
        <p>
          Both halves, and neither is redundant. Without the gate, a signed-out
          visitor downloads the page and watches it redirect, which looks broken
          and leaks the shape of your private UI. Without the page&apos;s own
          check, a browser holding a stale marker cookie renders a logged-in
          page for somebody with no session, and every request on it fails
          with a 401 nothing is handling.
        </p>

        <h2 id="verify">Verifying it, rather than assuming it</h2>
        <p>
          Four requests, and each one asserts something different. Run them
          against a project with the API on 8080 and the web app on 3000:
        </p>
        <CodeBlock
          language="bash"
          code={`# Signed out, a private page redirects and remembers where you were going.
$ curl -s -o /dev/null -w "%{http_code} -> %{redirect_url}\n" localhost:3000/account
307 -> http://localhost:3000/login?next=%2Faccount

# A public page is untouched by the gate.
$ curl -s -o /dev/null -w "%{http_code}\n" localhost:3000/products
200

# A marker cookie opens the gate, which is all it does.
$ curl -s -o /dev/null -w "%{http_code}\n" -b "shop_session=1" localhost:3000/account
200

# And the API is not fooled by it.
$ curl -s -o /dev/null -w "%{http_code}\n" -b "shop_session=1" localhost:8080/api/auth/me
401`}
        />
        <p>
          The last two together are the point of the design. The cookie gets a
          stranger the HTML and nothing else.
        </p>
        <p>Then the real flow, with a cookie jar:</p>
        <CodeBlock
          language="bash"
          code={`$ curl -s -c jar.txt -X POST localhost:8080/api/auth/register \
    -H "Content-Type: application/json" \
    -d '{"first_name":"Ada","last_name":"Lovelace","email":"ada@example.com","password":"Str0ng-Passw0rd!"}'

$ curl -s -b jar.txt localhost:8080/api/auth/me
{"data":{"first_name":"Ada","email":"ada@example.com","role":"USER", ...}}`}
        />

        <h2 id="roles">Gating on a role</h2>
        <p>
          <code>user.role</code> is <code>ADMIN</code>, <code>EDITOR</code> or{' '}
          <code>USER</code>. Hiding a link from a role is a UI decision and
          belongs in the component:
        </p>
        <CodeBlock
          language="tsx"
          code={`const { data: user } = useMe();

{user?.role === "ADMIN" ? (
  <a href={process.env.NEXT_PUBLIC_ADMIN_URL}>Admin panel</a>
) : null}`}
        />
        <p>
          <strong>
            Deciding what a role may <em>do</em> is not a UI decision.
          </strong>{' '}
          A hidden button is not a permission: the endpoint is still there and
          anybody can call it. Role and per-resource permissions are enforced in
          the service layer on the Go side, which is where they hold for every
          route, every job and every test. See{' '}
          <Link href="/docs/security/authorization">Authorization</Link>, and{' '}
          <code>--owned-by</code> and <code>--roles</code> in{' '}
          <Link href="/docs/concepts/field-types">the generator flags</Link> for
          the two cases you get for free.
        </p>

        <h2 id="mistakes">Four ways this goes wrong</h2>
        <ul>
          <li>
            <strong>Reading the token in JavaScript.</strong> It is{' '}
            <code>HttpOnly</code>, so you cannot, and the workaround, copying the
            one in the response body into <code>localStorage</code>, hands it to
            the first XSS bug or bad dependency you ship. The cookies are the
            session. Let them be.
          </li>
          <li>
            <strong>Treating the gate as security.</strong>{' '}
            <code>proxy.ts</code> decides what HTML to send. The API decides what
            data to return. Confusing the two produces a site that looks locked
            and is not.
          </li>
          <li>
            <strong>Writing your own <code>User</code> interface.</strong> It
            drifts from the API the first time a column is added, and nothing
            tells you. Import it from <code>@repo/shared/types</code>.
          </li>
          <li>
            <strong>Replacing the API&apos;s error message.</strong> &quot;Login
            failed&quot; is the same string for a wrong password, a locked
            account and a rate limit. The API already distinguishes them, and
            that distinction is what tells somebody whether to try again or wait.
          </li>
        </ul>

        <div className="p-4 rounded-lg border border-accent/20 bg-accent/5">
          <p className="text-sm text-foreground/80">
            <strong>One sentence to keep.</strong> The marker cookie answers
            &quot;has this browser signed in here?&quot; and decides what page to
            send. <code>/api/auth/me</code> answers &quot;is this session still
            valid?&quot; and decides what to render. Only the second one is a
            security boundary, and it lives on the server.
          </p>
        </div>
      </div>

      {/* Nav */}
      <div className="flex items-center justify-between pt-6 border-t border-border/30">
        <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
          <Link href="/docs/frontend/fetching-data" className="gap-1.5">
            <ArrowLeft className="h-3.5 w-3.5" />
            Fetching Data
          </Link>
        </Button>
        <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
          <Link href="/docs/frontend/tanstack-router" className="gap-1.5">
            TanStack Router
            <ArrowRight className="h-3.5 w-3.5" />
          </Link>
        </Button>
      </div>
    </DocsShell>
  )
}
