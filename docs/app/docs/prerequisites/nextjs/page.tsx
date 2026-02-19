import Link from "next/link";
import { ArrowRight, ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { SiteHeader } from "@/components/site-header";
import { DocsSidebar } from "@/components/docs-sidebar";

const reactBasicsCode = `// A React component is just a function that returns JSX
interface GreetingProps {
  name: string;
  role?: string;   // optional prop
  children?: React.ReactNode;
}

function Greeting({ name, role = "user", children }: GreetingProps) {
  return (
    <div className="p-4 rounded-lg border">
      <h2>Hello, {name}!</h2>
      {role === "admin" && <span className="text-red-500">Admin</span>}
      {children}
    </div>
  );
}

// Usage
<Greeting name="Alice" role="admin">
  <p>Welcome back.</p>
</Greeting>

// Event handling
function Counter() {
  const handleClick = (e: React.MouseEvent) => {
    console.log("Clicked!", e.target);
  };
  return <button onClick={handleClick}>Click me</button>;
}`;

const typescriptCode = `// Interface for object shapes
interface User {
  id: number;
  name: string;
  email: string;
  role: "ADMIN" | "EDITOR" | "USER";  // union type
  avatar?: string;                      // optional
}

// Type for unions and computed types
type Status = "active" | "inactive" | "banned";
type UserRole = User["role"];  // extracts "ADMIN" | "EDITOR" | "USER"

// Generics
interface ApiResponse<T> {
  data: T;
  message: string;
}

// Component props with TypeScript
interface UserCardProps {
  user: User;
  onEdit: (id: number) => void;
  variant?: "compact" | "full";
}

function UserCard({ user, onEdit, variant = "full" }: UserCardProps) {
  return (
    <div onClick={() => onEdit(user.id)}>
      <span>{user.name}</span>
      {variant === "full" && <span>{user.email}</span>}
    </div>
  );
}

// Record type for key-value maps
const statusColors: Record<Status, string> = {
  active: "text-green-500",
  inactive: "text-gray-500",
  banned: "text-red-500",
};`;

const useStateEffectCode = `"use client";

import { useState, useEffect } from "react";

function UserProfile({ userId }: { userId: number }) {
  // useState: declare reactive state
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [count, setCount] = useState(0);

  // useEffect: run side effects (API calls, subscriptions, etc.)
  useEffect(() => {
    // This runs when userId changes
    setLoading(true);
    fetch(\`/api/users/\${userId}\`)
      .then((res) => res.json())
      .then((data) => {
        setUser(data.data);
        setLoading(false);
      });

    // Cleanup function (runs before re-run or unmount)
    return () => {
      console.log("Cleaning up previous effect");
    };
  }, [userId]);  // dependency array: re-run when userId changes
                 // [] = run once on mount
                 // no array = run every render (avoid this)

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      <h1>{user?.name}</h1>
      <p>Count: {count}</p>
      <button onClick={() => setCount((prev) => prev + 1)}>
        Increment
      </button>
    </div>
  );
}`;

const appRouterCode = `# Next.js App Router maps folders to URL routes:
#
# app/
# \u251C\u2500\u2500 page.tsx                    \u2192  /
# \u251C\u2500\u2500 layout.tsx                  \u2192  wraps all pages
# \u251C\u2500\u2500 loading.tsx                 \u2192  loading UI for /
# \u251C\u2500\u2500 not-found.tsx               \u2192  404 page
# \u251C\u2500\u2500 (auth)/                     \u2192  route group (no URL segment)
# \u2502   \u251C\u2500\u2500 login/page.tsx          \u2192  /login
# \u2502   \u2514\u2500\u2500 register/page.tsx       \u2192  /register
# \u251C\u2500\u2500 (dashboard)/                \u2192  route group with its own layout
# \u2502   \u251C\u2500\u2500 layout.tsx              \u2192  dashboard layout (sidebar, navbar)
# \u2502   \u251C\u2500\u2500 page.tsx                \u2192  /  (dashboard home)
# \u2502   \u2514\u2500\u2500 users/
# \u2502       \u251C\u2500\u2500 page.tsx            \u2192  /users
# \u2502       \u2514\u2500\u2500 [id]/
# \u2502           \u2514\u2500\u2500 page.tsx        \u2192  /users/123 (dynamic route)
# \u2514\u2500\u2500 api/                        \u2192  API routes (not used in Grit)

# Key files:
# - page.tsx    \u2192  the UI for a route (required to make a route accessible)
# - layout.tsx  \u2192  shared UI that wraps child pages (persists across navigation)
# - loading.tsx \u2192  loading state shown while page data loads
# - error.tsx   \u2192  error boundary for a route segment`;

const serverClientCode = `// SERVER COMPONENT (default - no directive needed)
// Can: fetch data, access DB, read files, use async/await
// Cannot: use useState, useEffect, event handlers, browser APIs
async function UsersPage() {
  const res = await fetch("http://localhost:8080/api/users");
  const data = await res.json();

  return (
    <div>
      <h1>Users</h1>
      <UserList users={data.data} />
    </div>
  );
}

// CLIENT COMPONENT (needs "use client" directive)
// Can: use hooks, event handlers, browser APIs, interactivity
// Cannot: be async, directly access server resources
"use client";

import { useState } from "react";

function UserList({ users }: { users: User[] }) {
  const [search, setSearch] = useState("");

  const filtered = users.filter((u) =>
    u.name.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div>
      <input
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        placeholder="Search users..."
      />
      {filtered.map((user) => (
        <div key={user.id}>{user.name}</div>
      ))}
    </div>
  );
}`;

const tailwindCode = `// Tailwind uses utility classes directly in JSX
function ProductCard({ product }: { product: Product }) {
  return (
    <div className="rounded-xl border border-border bg-card p-6
                    hover:bg-accent/10 transition-colors
                    sm:p-4 md:p-6 lg:p-8">
      {/* Responsive: sm/md/lg prefixes */}
      <h3 className="text-lg font-semibold text-foreground
                     sm:text-base md:text-lg">
        {product.name}
      </h3>
      <p className="text-sm text-muted-foreground mt-2">
        {product.description}
      </p>
      <span className="text-primary font-bold text-xl mt-4 block">
        \${product.price}
      </span>
    </div>
  );
}

// shadcn/ui components (pre-built, customizable)
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Dialog, DialogTrigger, DialogContent } from "@/components/ui/dialog";

function Example() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Create User</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        <Input placeholder="Name" />
        <Input placeholder="Email" type="email" />
        <Button variant="default">Save</Button>
        <Button variant="outline">Cancel</Button>
        <Button variant="destructive">Delete</Button>
      </CardContent>
    </Card>
  );
}`;

const reactQueryCode = `"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

// Fetching data with useQuery
function usePosts(page: number = 1) {
  return useQuery({
    queryKey: ["posts", { page }],   // unique cache key
    queryFn: async () => {
      const { data } = await apiClient.get(\`/api/posts?page=\${page}\`);
      return data;                    // { data: Post[], meta: {...} }
    },
    staleTime: 30 * 1000,            // data is fresh for 30 seconds
  });
}

// Creating data with useMutation
function useCreatePost() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (newPost: CreatePostInput) => {
      const { data } = await apiClient.post("/api/posts", newPost);
      return data;
    },
    onSuccess: () => {
      // Invalidate the cache so the list refetches
      queryClient.invalidateQueries({ queryKey: ["posts"] });
    },
  });
}

// Using in a component
function PostsPage() {
  const { data, isLoading, error } = usePosts(1);
  const createPost = useCreatePost();

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return (
    <div>
      <button onClick={() => createPost.mutate({ title: "New Post" })}>
        Add Post
      </button>
      {data?.data.map((post: Post) => (
        <div key={post.id}>{post.title}</div>
      ))}
    </div>
  );
}`;

const zodCode = `import { z } from "zod";

// Define a schema
const CreateUserSchema = z.object({
  name: z.string().min(2, "Name must be at least 2 characters"),
  email: z.string().email("Invalid email address"),
  age: z.number().min(18, "Must be 18 or older").optional(),
  role: z.enum(["ADMIN", "EDITOR", "USER"]),
  isActive: z.boolean().default(true),
  tags: z.array(z.string()).optional(),
});

// Extract the TypeScript type from the schema
type CreateUserInput = z.infer<typeof CreateUserSchema>;
// Result:
// {
//   name: string;
//   email: string;
//   age?: number;
//   role: "ADMIN" | "EDITOR" | "USER";
//   isActive: boolean;
//   tags?: string[];
// }

// Validate data
const result = CreateUserSchema.safeParse({
  name: "Alice",
  email: "alice@example.com",
  role: "ADMIN",
});

if (result.success) {
  console.log(result.data);  // typed as CreateUserInput
} else {
  console.log(result.error.issues);
  // [{ path: ["email"], message: "Invalid email address" }]
}`;

const reactHookFormCode = `"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";

const LoginSchema = z.object({
  email: z.string().email("Invalid email"),
  password: z.string().min(8, "Password must be at least 8 characters"),
});

type LoginInput = z.infer<typeof LoginSchema>;

function LoginForm() {
  const {
    register,       // connects inputs to the form
    handleSubmit,    // wraps your submit handler with validation
    formState: { errors, isSubmitting },
  } = useForm<LoginInput>({
    resolver: zodResolver(LoginSchema),  // Zod validates on submit
  });

  const onSubmit = async (data: LoginInput) => {
    // data is fully validated and typed
    const response = await fetch("/api/auth/login", {
      method: "POST",
      body: JSON.stringify(data),
    });
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <div>
        <label>Email</label>
        <input {...register("email")} type="email" className="input" />
        {errors.email && (
          <span className="text-red-500 text-sm">
            {errors.email.message}
          </span>
        )}
      </div>

      <div>
        <label>Password</label>
        <input {...register("password")} type="password" className="input" />
        {errors.password && (
          <span className="text-red-500 text-sm">
            {errors.password.message}
          </span>
        )}
      </div>

      <button type="submit" disabled={isSubmitting}>
        {isSubmitting ? "Signing in..." : "Sign In"}
      </button>
    </form>
  );
}`;

const customHooksCode = `// hooks/use-users.ts
"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { User, CreateUserInput } from "@shared/types";

// Encapsulate all user-related data fetching in one hook file
export function useUsers(params?: {
  page?: number;
  search?: string;
  role?: string;
}) {
  return useQuery({
    queryKey: ["users", params],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/users", { params });
      return data;
    },
  });
}

export function useUser(id: number) {
  return useQuery({
    queryKey: ["users", id],
    queryFn: async () => {
      const { data } = await apiClient.get(\`/api/users/\${id}\`);
      return data.data as User;
    },
    enabled: !!id,  // only fetch if id exists
  });
}

export function useCreateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (input: CreateUserInput) => {
      const { data } = await apiClient.post("/api/users", input);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
}

// When to make a custom hook:
// 1. Logic is used in 2+ components
// 2. A component has complex state/effect logic
// 3. You want to separate data fetching from UI`;

const apiClientCode = `// lib/api-client.ts
import axios from "axios";

export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080",
  headers: { "Content-Type": "application/json" },
});

// Automatically attach JWT token to every request
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = \`Bearer \${token}\`;
  }
  return config;
});

// Handle 401 responses (expired token)
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("token");
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

// Folder structure overview:
//
// apps/web/ (or apps/admin/)
// \u251C\u2500\u2500 app/                  \u2192  Routes (pages, layouts, loading states)
// \u2502   \u251C\u2500\u2500 (auth)/           \u2192  Auth pages (login, register)
// \u2502   \u251C\u2500\u2500 (dashboard)/      \u2192  Dashboard pages with shared layout
// \u2502   \u2514\u2500\u2500 layout.tsx        \u2192  Root layout (providers, fonts)
// \u251C\u2500\u2500 components/           \u2192  Reusable UI components
// \u2502   \u251C\u2500\u2500 ui/               \u2192  shadcn/ui primitives (Button, Input, etc.)
// \u2502   \u2514\u2500\u2500 layout/           \u2192  Sidebar, Navbar, etc.
// \u251C\u2500\u2500 hooks/                \u2192  React Query hooks (use-users.ts, etc.)
// \u251C\u2500\u2500 lib/                  \u2192  Utilities (api-client.ts, utils.ts)
// \u2514\u2500\u2500 public/               \u2192  Static assets (images, fonts)`;

export default function NextjsForGritPage() {
  return (
    <div className="min-h-screen bg-background">
      <SiteHeader />
      <DocsSidebar />

      <main className="lg:pl-64">
        <div className="container max-w-screen-xl py-10 px-6">
          <div className="max-w-3xl">
            {/* Header */}
            <div className="mb-10">
              <span className="tag-mono text-primary/80 mb-3 block">Prerequisites</span>
              <h1 className="text-4xl font-bold tracking-tight mb-4">Next.js &amp; React for Grit Developers</h1>
              <p className="text-xl text-muted-foreground leading-relaxed">
                A focused guide to the React and Next.js concepts you need to build with Grit.
                If you know another language or framework but are new to the React ecosystem,
                this page covers everything from components and hooks to data fetching and validation
                -- the building blocks behind Grit&apos;s admin panel and web app.
              </p>
            </div>

            <div className="prose-grit">
              {/* 1. React Basics */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  React Basics
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  React components are JavaScript functions that return JSX -- a syntax that looks like HTML but
                  compiles to JavaScript. Every piece of UI in a React app is a component. Components accept
                  <strong> props</strong> (input data) and can render other components as <strong>children</strong>.
                  You handle events with inline functions and conditionally render elements using JavaScript
                  expressions like <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">&amp;&amp;</code> and
                  ternaries.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">components/Greeting.tsx</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{reactBasicsCode}</pre>
                </div>

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
                  <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
                  <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                    Every file in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">app/</code> and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">components/</code> is
                    a React component. Pages like <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">page.tsx</code> export
                    a default component for the route. Reusable UI pieces live
                    in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">components/</code> and
                    are imported wherever needed. Grit generates both page components and shared components
                    when you scaffold a project.
                  </p>
                </div>
              </div>

              {/* 2. TypeScript Essentials */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  TypeScript Essentials
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  TypeScript adds static types to JavaScript, catching bugs at compile time instead of runtime.
                  Use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">interface</code> to
                  describe object shapes and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">type</code> for
                  unions, intersections, and computed types. Generics let you write reusable typed code
                  (like <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">{"ApiResponse<T>"}</code>).
                  TypeScript infers types wherever possible, so you don&apos;t need to annotate everything --
                  but explicit types on function parameters and return values make your code self-documenting.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">types.ts</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{typescriptCode}</pre>
                </div>

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
                  <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
                  <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                    All Grit code is TypeScript. Shared types live
                    in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">packages/shared/src/types/</code> and
                    are generated from your Go models by <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit sync</code>.
                    Both the admin panel and web app import these types, ensuring your frontend
                    always matches your backend data structures.
                  </p>
                </div>
              </div>

              {/* 3. useState & useEffect */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  useState &amp; useEffect
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  React hooks are functions that let components have state and side effects. <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useState</code> declares
                  a piece of reactive state -- when you call the setter function, React re-renders the component
                  with the new value. <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useEffect</code> runs
                  side effects like API calls, subscriptions, or timers. The dependency array controls when the
                  effect re-runs, and the cleanup function handles teardown.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">components/UserProfile.tsx</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{useStateEffectCode}</pre>
                </div>

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
                  <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
                  <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                    Grit uses <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useState</code> throughout
                    the admin panel for local UI state (search filters, pagination, modal visibility) and
                    in the web app for form inputs and toggles. You rarely need raw <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useEffect</code> for
                    data fetching because React Query handles that, but it&apos;s still used for things like
                    keyboard shortcuts and resize listeners.
                  </p>
                </div>
              </div>

              {/* 4. Next.js App Router */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Next.js App Router
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Next.js uses file-based routing: folders inside <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">app/</code> become
                  URL routes automatically. A <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">page.tsx</code> file
                  makes a route accessible, <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">layout.tsx</code> wraps
                  child pages with shared UI (like a sidebar), and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">loading.tsx</code> shows
                  a loading state. Route groups with parentheses
                  like <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">(auth)</code> organize
                  routes without adding a URL segment. Dynamic routes use square brackets
                  like <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">[id]</code>.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">app/ folder structure</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{appRouterCode}</pre>
                </div>

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
                  <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
                  <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                    Both <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/web</code> and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">apps/admin</code> use
                    the App Router exclusively. The admin app uses a <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">(dashboard)</code> route
                    group with a shared layout containing the sidebar and navbar. Auth pages sit in
                    an <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">(auth)</code> group
                    with a minimal layout. When you run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>,
                    it creates new page files inside the appropriate route group.
                  </p>
                </div>
              </div>

              {/* 5. Server vs Client Components */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Server vs Client Components
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  In the App Router, components are <strong>server components</strong> by default. They run on the server,
                  can be <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">async</code>, and
                  can fetch data directly -- but they cannot use hooks or browser APIs.
                  Add <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">&quot;use client&quot;</code> at
                  the top of a file to make it a <strong>client component</strong>, which runs in the browser
                  and can use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useState</code>, <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useEffect</code>,
                  event handlers, and everything interactive.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">server-vs-client.tsx</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{serverClientCode}</pre>
                </div>

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
                  <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
                  <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                    In Grit, page-level components are typically server components that render the page shell.
                    Interactive parts like DataTable, FormBuilder, and any component using React Query hooks
                    are marked <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">&quot;use client&quot;</code>.
                    The general rule: keep server components for layout and structure, push interactivity down
                    to the smallest client component possible.
                  </p>
                </div>
              </div>

              {/* 6. Tailwind CSS & shadcn/ui */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Tailwind CSS &amp; shadcn/ui
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Tailwind CSS is a utility-first framework where you style elements by composing small CSS classes
                  directly in your JSX. Instead of writing <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">margin-top: 1rem</code> in
                  a CSS file, you write <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">mt-4</code> on
                  the element. Responsive design uses breakpoint
                  prefixes (<code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">sm:</code>, <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">md:</code>, <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">lg:</code>).
                  shadcn/ui provides pre-built, accessible components (Button, Card, Dialog, Input) that use Tailwind
                  under the hood and can be customized by editing the source files directly.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">components/ProductCard.tsx</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{tailwindCode}</pre>
                </div>

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
                  <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
                  <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                    All styling in Grit is Tailwind -- there are no custom CSS files except for the
                    base <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">globals.css</code> which
                    defines CSS variables for the dark theme. UI components come from shadcn/ui and are
                    scaffolded into <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">components/ui/</code>.
                    You can customize any shadcn component by editing it directly since it&apos;s your own code,
                    not a node_modules dependency.
                  </p>
                </div>
              </div>

              {/* 7. React Query */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  React Query (TanStack Query)
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  React Query is the data fetching layer in Grit. <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useQuery</code> fetches
                  and caches server data with automatic background refetching, stale-time management, and
                  loading/error states. <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useMutation</code> handles
                  create, update, and delete operations with callbacks
                  like <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">onSuccess</code> to
                  invalidate caches. The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">QueryClientProvider</code> at
                  the root of the app provides the cache to all components.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">hooks/use-posts.ts</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{reactQueryCode}</pre>
                </div>

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
                  <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
                  <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                    All data fetching in Grit goes through React Query hooks stored
                    in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">hooks/</code>. When you
                    run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>,
                    it creates a complete hook file with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useList</code>, <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useGet</code>, <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useCreate</code>, <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useUpdate</code>,
                    and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useDelete</code> hooks.
                    Components never call <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">fetch</code> or <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">axios</code> directly.
                  </p>
                </div>
              </div>

              {/* 8. Zod Validation */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Zod Validation
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Zod is a TypeScript-first schema validation library. You define the shape and constraints of your
                  data with a fluent API, then use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">z.infer</code> to
                  extract the TypeScript type from the schema -- so your validation rules and types are always
                  in sync. Zod supports strings, numbers, booleans, enums, arrays, nested objects, optional fields,
                  and custom validation messages. Use <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">safeParse</code> for
                  error handling without exceptions.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">schemas/user.ts</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{zodCode}</pre>
                </div>

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
                  <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
                  <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                    Zod schemas live in <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">packages/shared/src/schemas/</code> and
                    are shared across the admin panel and web app. When you run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit generate resource</code>,
                    it creates Zod schemas matching your Go model fields.
                    The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit sync</code> command
                    can regenerate schemas from Go types to keep everything aligned.
                  </p>
                </div>
              </div>

              {/* 9. React Hook Form */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  React Hook Form
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  React Hook Form manages form state, validation, and submission with minimal re-renders.
                  The <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">useForm</code> hook
                  returns <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">register</code> (connects
                  inputs), <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">handleSubmit</code> (validates
                  before calling your handler), and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">formState</code> (errors,
                  submission state). Combined with the Zod resolver, your form validation uses the same schemas
                  you already defined -- no duplicate validation logic.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">components/LoginForm.tsx</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{reactHookFormCode}</pre>
                </div>

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
                  <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
                  <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                    Grit&apos;s admin panel includes a FormBuilder component that wraps React Hook Form with Zod
                    validation, supporting 8 field types (text, textarea, number, select, date, toggle, checkbox,
                    radio). For simple forms or custom requirements, you can use React Hook Form directly. The
                    profile page and auth pages both use raw React Hook Form with Zod resolvers.
                  </p>
                </div>
              </div>

              {/* 10. Custom Hooks */}
              <div className="mb-12">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  Custom Hooks
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  Custom hooks are functions starting with <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">use</code> that
                  extract reusable logic from components. They can call other hooks, manage state, handle
                  side effects, and return any values. Create a custom hook when: the same data-fetching logic
                  appears in multiple components, a component&apos;s hook logic becomes complex, or you want to separate
                  data concerns from UI rendering. Name files with the <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">use-</code> prefix
                  in kebab-case.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">hooks/use-users.ts</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{customHooksCode}</pre>
                </div>

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
                  <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
                  <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                    Grit generates resource-specific hooks
                    like <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">hooks/use-users.ts</code> and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">hooks/use-posts.ts</code>.
                    The admin panel also has <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">hooks/use-resource.ts</code> (a
                    generic hook powered by resource definitions) and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">hooks/use-system.ts</code> for
                    system pages. Following the convention of one hook file per data domain keeps your codebase organized.
                  </p>
                </div>
              </div>

              {/* 11. API Client & Folder Structure */}
              <div className="mb-10">
                <h2 className="text-2xl font-semibold tracking-tight mb-4">
                  API Client &amp; Folder Structure
                </h2>
                <p className="text-muted-foreground leading-relaxed mb-4">
                  The API client is an Axios instance with a base URL pointing to the Go backend and interceptors
                  that automatically attach JWT tokens to every request and redirect to the login page on 401
                  responses. All hooks import this shared client, so authentication is handled in one place.
                  The frontend follows a clear folder structure:
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">app/</code> for routes,
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded"> components/</code> for UI,
                  <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded"> hooks/</code> for data logic,
                  and <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded"> lib/</code> for utilities.
                </p>

                <div className="rounded-xl border border-border/40 bg-card/80 overflow-hidden mb-6">
                  <div className="flex items-center gap-2 px-4 py-2.5 border-b border-border/30 bg-accent/30">
                    <span className="text-[11px] font-mono text-muted-foreground/40">lib/api-client.ts</span>
                  </div>
                  <pre className="p-5 text-sm font-mono text-foreground/80 overflow-x-auto">{apiClientCode}</pre>
                </div>

                <div className="rounded-xl border border-primary/20 bg-primary/5 p-5 mb-8">
                  <h4 className="text-sm font-semibold text-primary/80 uppercase tracking-wider mb-2">In Grit</h4>
                  <p className="text-[13px] text-muted-foreground/80 leading-relaxed">
                    When you run <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">grit new</code>,
                    the API client is scaffolded
                    at <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">lib/api-client.ts</code> in
                    both the admin and web apps. It reads <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">NEXT_PUBLIC_API_URL</code> from
                    your <code className="text-xs font-mono bg-accent/50 px-1.5 py-0.5 rounded">.env</code> file.
                    All generated hooks import from this single client, so you never need to configure
                    authentication headers or base URLs in individual components.
                  </p>
                </div>
              </div>
            </div>

            {/* Navigation footer */}
            <div className="flex items-center justify-between pt-6 border-t border-border/30">
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/prerequisites/golang" className="gap-1.5">
                  <ArrowLeft className="h-3.5 w-3.5" />
                  Go for Grit Developers
                </Link>
              </Button>
              <Button variant="ghost" size="sm" asChild className="text-muted-foreground/60 hover:text-foreground">
                <Link href="/docs/prerequisites/docker" className="gap-1.5">
                  Docker for Grit Developers
                  <ArrowRight className="h-3.5 w-3.5" />
                </Link>
              </Button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}
