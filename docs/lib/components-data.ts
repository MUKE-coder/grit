export type Platform = "Next.js" | "Expo" | "Go" | "All";

export interface Component {
  id: string;
  name: string;
  description: string;
  thumbnail: string;
  platform: Platform[];
  isPaid: boolean;
  price?: number;
  functionalities: string[];
  installCommand: string;
  category: string;
  docsUrl: string;
  usageSteps: {
    title: string;
    description: string;
    code?: string;
  }[];
}

export const components: Component[] = [
  {
    id: "jb-better-auth-ui",
    name: "JB Better Auth UI",
    description:
      "Complete authentication system with sign-in, sign-up, email verification, password reset, social auth (Google & GitHub), and session management. Built with Better Auth, Prisma, and Resend.",
    thumbnail: "/api/placeholder/400/300",
    platform: ["Next.js"],
    isPaid: false,
    functionalities: [
      "Sign In & Sign Up Forms",
      "Email OTP Verification",
      "Forgot & Reset Password",
      "Change Password",
      "Social Auth (Google & GitHub)",
      "Session Management",
      "Protected Routes Middleware",
      "User Profile Component",
      "Logout Button",
      "Dark Mode Support",
    ],
    installCommand:
      "pnpm dlx shadcn@latest add https://better-auth-ui.desishub.com/r/auth-components.json",
    category: "Authentication",
    docsUrl: "https://better-auth-ui.desishub.com",
    usageSteps: [
      {
        title: "Install the component",
        description:
          "Run the shadcn CLI to add the auth components to your project.",
        code: "pnpm dlx shadcn@latest add https://better-auth-ui.desishub.com/r/auth-components.json",
      },
      {
        title: "Configure environment variables",
        description:
          "Add your database connection, Better Auth secret, email (Resend), and OAuth provider credentials.",
        code: `DATABASE_URL="postgresql://..."
BETTER_AUTH_SECRET="your-secret"
RESEND_API_KEY="re_..."
GOOGLE_CLIENT_ID="..."
GOOGLE_CLIENT_SECRET="..."
GITHUB_CLIENT_ID="..."
GITHUB_CLIENT_SECRET="..."`,
      },
      {
        title: "Set up Prisma & database",
        description:
          "Run migrations to set up the auth tables in your database.",
        code: `npx prisma migrate dev --name init
npx prisma generate`,
      },
      {
        title: "Use auth components",
        description: "Import and use the auth components in your pages.",
        code: `import { SignIn } from '@/components/auth/sign-in'
import { SignUp } from '@/components/auth/sign-up'

// Sign In page
export default function SignInPage() {
  return <SignIn />
}`,
      },
    ],
  },
  {
    id: "stripe-ui-component",
    name: "Stripe UI Component",
    description:
      "Full-featured Stripe payment integration with checkout flow, order management, product grid, and cart powered by Zustand. Includes server-side API routes and auth-protected checkout.",
    thumbnail: "/api/placeholder/400/300",
    platform: ["Next.js"],
    isPaid: false,
    functionalities: [
      "Stripe Payment Element",
      "Checkout Flow with Intent",
      "Order Management & History",
      "Product Grid Display",
      "Zustand Shopping Cart",
      "Server-side API Routes",
      "Auth-Protected Checkout",
      "Webhook Handling",
    ],
    installCommand:
      "pnpm dlx shadcn@latest add https://stripe-ui-component.desishub.com/r/stripe-ui-component.json",
    category: "Payments",
    docsUrl: "https://stripe-ui-component.desishub.com",
    usageSteps: [
      {
        title: "Install the component",
        description: "Run the shadcn CLI to add the Stripe UI component.",
        code: "pnpm dlx shadcn@latest add https://stripe-ui-component.desishub.com/r/stripe-ui-component.json",
      },
      {
        title: "Configure Stripe keys",
        description:
          "Add your Stripe API keys and webhook secret to environment variables.",
        code: `STRIPE_SECRET_KEY="sk_test_..."
NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY="pk_test_..."
STRIPE_WEBHOOK_SECRET="whsec_..."`,
      },
      {
        title: "Set up database",
        description:
          "Run Prisma migrations to create product and order tables.",
        code: `npx prisma migrate dev --name add-stripe
npx prisma generate`,
      },
      {
        title: "Use the components",
        description: "Import the product grid and checkout components.",
        code: `import { ProductGrid } from '@/components/stripe/product-grid'
import { CheckoutForm } from '@/components/stripe/checkout-form'

// Products page
export default function ProductsPage() {
  return <ProductGrid />
}`,
      },
    ],
  },
  {
    id: "file-storage-ui",
    name: "File Storage UI",
    description:
      "Multi-provider file storage with drag-and-drop upload (5 variants), progress tracking, presigned URLs, file dashboard, and category management. Supports AWS S3 and Cloudflare R2.",
    thumbnail: "/api/placeholder/400/300",
    platform: ["Next.js"],
    isPaid: false,
    functionalities: [
      "Multi-Provider (AWS S3 & Cloudflare R2)",
      "Drag-and-Drop Upload (5 Variants)",
      "Upload Progress Tracking",
      "Presigned URL Generation",
      "File Dashboard & Gallery",
      "Category Management",
      "Storage Statistics",
      "File Type Validation",
    ],
    installCommand:
      "pnpm dlx shadcn@latest add https://file-storage-registry.vercel.app/r/file-storage.json",
    category: "Media",
    docsUrl: "https://file-storage-registry.vercel.app",
    usageSteps: [
      {
        title: "Install the component",
        description: "Run the shadcn CLI to add the file storage component.",
        code: "pnpm dlx shadcn@latest add https://file-storage-registry.vercel.app/r/file-storage.json",
      },
      {
        title: "Configure storage provider",
        description: "Set up your S3 or Cloudflare R2 credentials.",
        code: `# AWS S3
AWS_ACCESS_KEY_ID="..."
AWS_SECRET_ACCESS_KEY="..."
AWS_S3_BUCKET="your-bucket"
AWS_REGION="us-east-1"

# Or Cloudflare R2
R2_ACCESS_KEY_ID="..."
R2_SECRET_ACCESS_KEY="..."
R2_BUCKET="your-bucket"
R2_ENDPOINT="https://..."`,
      },
      {
        title: "Set up database",
        description: "Run migrations to create file and category tables.",
        code: `npx prisma migrate dev --name add-file-storage
npx prisma generate`,
      },
      {
        title: "Use the file upload",
        description: "Import and render the file upload component.",
        code: `import { FileUpload } from '@/components/file-storage/file-upload'
import { FileDashboard } from '@/components/file-storage/file-dashboard'

// Upload page
export default function UploadPage() {
  return <FileUpload variant="default" />
}`,
      },
    ],
  },
  {
    id: "multi-step-form",
    name: "Multi Step Form",
    description:
      "Multi-step form wizard with Zustand state management, step navigation, form validation with Zod, and progress tracking. Perfect for onboarding flows and complex forms.",
    thumbnail: "/api/placeholder/400/300",
    platform: ["Next.js"],
    isPaid: false,
    functionalities: [
      "Multi-Step Form Wizard",
      "Zustand State Management",
      "Step Navigation (Next/Prev)",
      "Form Validation with Zod",
      "Progress Tracking",
      "Responsive Design",
    ],
    installCommand:
      "pnpm dlx shadcn@latest add https://jb.desishub.com/r/multi-step-form.json",
    category: "Forms",
    docsUrl: "https://jb.desishub.com/docs/components/multi-step-form",
    usageSteps: [
      {
        title: "Install the component",
        description: "Run the shadcn CLI to add the multi-step form.",
        code: "pnpm dlx shadcn@latest add https://jb.desishub.com/r/multi-step-form.json",
      },
      {
        title: "Configure your form steps",
        description: "Define the steps and their validation schemas.",
        code: `import { useFormStore } from '@/store/form-store'

// The form store manages step navigation
// and form data persistence across steps`,
      },
      {
        title: "Use the MultiStepForm component",
        description: "Import and render the multi-step form in your page.",
        code: `import { MultiStepForm } from '@/components/multi-step-form'

export default function OnboardingPage() {
  return <MultiStepForm />
}`,
      },
    ],
  },
  {
    id: "zustand-cart",
    name: "Zustand Cart",
    description:
      "Reusable shopping cart powered by Zustand with localStorage persistence, product grid, quantity management, price calculations, and wishlist functionality.",
    thumbnail: "/api/placeholder/400/300",
    platform: ["Next.js"],
    isPaid: false,
    functionalities: [
      "Zustand Cart State",
      "LocalStorage Persistence",
      "Product Grid Display",
      "Add / Increment / Decrement Items",
      "Remove from Cart",
      "Automatic Price Calculations",
      "Wishlist UI",
    ],
    installCommand:
      "pnpm dlx shadcn@latest add https://jb.desishub.com/r/zustand-cart.json",
    category: "E-Commerce",
    docsUrl: "https://jb.desishub.com/docs/components/zustand-cart",
    usageSteps: [
      {
        title: "Install the component",
        description: "Run the shadcn CLI to add the Zustand cart.",
        code: "pnpm dlx shadcn@latest add https://jb.desishub.com/r/zustand-cart.json",
      },
      {
        title: "Set up the cart store",
        description:
          "The cart store is automatically created with Zustand and localStorage persistence.",
        code: `import { useCartStore } from '@/store/cart-store'

// Access cart state anywhere
const { items, addItem, removeItem, total } = useCartStore()`,
      },
      {
        title: "Use the cart components",
        description: "Import the product grid and cart components.",
        code: `import { ProductGrid } from '@/components/cart/product-grid'
import { CartDrawer } from '@/components/cart/cart-drawer'

export default function ShopPage() {
  return (
    <>
      <CartDrawer />
      <ProductGrid />
    </>
  )
}`,
      },
    ],
  },
  {
    id: "work-experience",
    name: "Work Experience",
    description:
      "Professional work experience timeline component with support for multiple companies, positions, employment periods, skills tags, and expandable/collapsible sections.",
    thumbnail: "/api/placeholder/400/300",
    platform: ["Next.js"],
    isPaid: false,
    functionalities: [
      "Multiple Companies with Logos",
      "Multiple Positions per Company",
      "Employment Period & Type",
      "Skills Tags per Position",
      "Expandable/Collapsible Sections",
      "Current Employer Badge",
      "Responsive Timeline Layout",
    ],
    installCommand:
      "pnpm dlx shadcn@latest add https://jb.desishub.com/r/work-experience.json",
    category: "UI Components",
    docsUrl: "https://jb.desishub.com/docs/components/work-experience",
    usageSteps: [
      {
        title: "Install the component",
        description: "Run the shadcn CLI to add the work experience component.",
        code: "pnpm dlx shadcn@latest add https://jb.desishub.com/r/work-experience.json",
      },
      {
        title: "Install dependencies",
        description: "Add the required Tailwind typography plugin.",
        code: "pnpm add @tailwindcss/typography",
      },
      {
        title: "Configure your experience data",
        description:
          "Define your work experience entries with companies, positions, and skills.",
        code: `const experiences = [
  {
    company: "Acme Corp",
    logo: "/logos/acme.png",
    positions: [
      {
        title: "Senior Developer",
        period: "2023 - Present",
        type: "Full-time",
        skills: ["React", "TypeScript", "Node.js"],
        description: "Leading frontend development..."
      }
    ]
  }
]`,
      },
      {
        title: "Use the component",
        description: "Import and render the work experience timeline.",
        code: `import { WorkExperience } from '@/components/work-experience'

export default function ResumePage() {
  return <WorkExperience experiences={experiences} />
}`,
      },
    ],
  },
  {
    id: "offline-sync",
    name: "Offline Sync",
    description:
      "Offline-first data sync architecture with IndexedDB (Dexie) for instant local writes, Prisma for server sync, conflict detection with version tracking, smart retry with exponential backoff, soft deletes, and PWA support. One command adds 12 files for a complete offline-first data layer.",
    thumbnail: "/api/placeholder/400/300",
    platform: ["Next.js"],
    isPaid: false,
    functionalities: [
      "IndexedDB Local Storage (Dexie)",
      "Automatic Background Sync",
      "Conflict Detection (Version Tracking)",
      "Smart Retry with Exponential Backoff",
      "Zod-Validated Server Actions",
      "PWA Ready (Service Worker + Manifest)",
      "UUID Primary Keys",
      "Soft Deletes",
      "Delta Sync (Only Changed Records)",
      "Force Offline Toggle",
      "Pending Changes Badge & Sync Status",
      "Reactive Queries (useLiveQuery)",
    ],
    installCommand:
      "pnpm dlx shadcn@latest add https://offline-sync.desishub.com/r/offline-sync.json",
    category: "UI Components",
    docsUrl: "https://offline-sync.desishub.com/docs",
    usageSteps: [
      {
        title: "Install the component",
        description:
          "Run the shadcn CLI to add offline-sync. This installs 12 files, 5 shadcn dependencies, and 4 npm packages.",
        code: "pnpm dlx shadcn@latest add https://offline-sync.desishub.com/r/offline-sync.json",
      },
      {
        title: "Set up Prisma schema",
        description:
          "Merge the installed prisma/schema.offline-sync.prisma into your schema.prisma. Every sync-enabled model needs id (UUID), version, isDeleted, updatedAt, and createdAt fields.",
        code: `npx prisma db push
npx prisma generate`,
      },
      {
        title: "Wire up providers",
        description:
          "Wrap your app with OnlineProvider and add PWARegister in your root layout.",
        code: `import { OnlineProvider } from "@/components/online-provider"
import { PWARegister } from "@/components/pwa-register"
import { Toaster } from "@/components/ui/sonner"

export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body>
        <PWARegister />
        <OnlineProvider>{children}</OnlineProvider>
        <Toaster />
      </body>
    </html>
  )
}`,
      },
      {
        title: "Use in your pages",
        description:
          "Use useLiveQuery for reactive UI updates and the repository functions for CRUD operations.",
        code: `import { useLiveQuery } from "dexie-react-hooks"
import { db } from "@/lib/offline-sync-db"
import * as repo from "@/lib/offline-sync-repository"
import { OnlineToggle } from "@/components/online-toggle"

export default function Home() {
  const categories = useLiveQuery(
    () => db.categories.filter((c) => !c.isDeleted).toArray(), []
  )

  return (
    <main>
      <OnlineToggle />
      {categories?.map((cat) => (
        <div key={cat.id}>{cat.name}</div>
      ))}
    </main>
  )
}`,
      },
    ],
  },
];

export const categories = [
  "All",
  "Authentication",
  "Payments",
  "UI Components",
  "Media",
  "Forms",
  "E-Commerce",
  "Offline & Sync",
] as const;

export function getComponentById(id: string): Component | undefined {
  return components.find((c) => c.id === id);
}

export function getComponentsByCategory(category: string): Component[] {
  if (category === "All") return components;
  return components.filter((c) => c.category === category);
}

export function getComponentsByPlatform(platform: Platform): Component[] {
  return components.filter(
    (c) => c.platform.includes(platform) || c.platform.includes("All"),
  );
}
