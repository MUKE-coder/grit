import Link from "next/link"
import {
  Smartphone,
  Clock,
  BookOpen,
  ArrowLeft,
  Trophy,
  ChevronDown,
  Rocket,
  ShieldCheck,
  Cloud,
  Bell,
  Store,
} from "lucide-react"
import { SiteHeader } from "@/components/site-header"
import type { Metadata } from "next"

export const metadata: Metadata = {
  title: "Grit Mobile Course — Building Mobile Applications with Grit",
  description:
    "Free, self-paced course covering Go API + Expo React Native mobile app development. 5 mini-courses, ~2.5 hours, 5 challenges.",
  openGraph: {
    title: "Grit Mobile Course — Building Mobile Applications with Grit",
    description:
      "Free, self-paced course covering Go API + Expo React Native mobile app development. 5 mini-courses, ~2.5 hours, 5 challenges.",
    url: "https://gritframework.dev/courses/grit-mobile",
  },
}

/* -- Course data ----------------------------------------------------------- */

interface Lesson {
  title: string
  description: string
}

interface MiniCourse {
  number: number
  title: string
  duration: string
  overview: string
  icon: React.ComponentType<{ className?: string }>
  lessons: Lesson[]
  challenge: {
    title: string
    description: string
  }
}

const courses: MiniCourse[] = [
  {
    number: 1,
    title: "Your First Mobile App",
    duration: "30 min",
    overview:
      "Scaffold a mobile project with Grit (API + Expo), set up the development environment, run on a device or emulator.",
    icon: Rocket,
    lessons: [
      {
        title: "Prerequisites",
        description:
          "Go, Node, pnpm, Expo CLI, iOS Simulator or Android Emulator.",
      },
      {
        title: "Scaffolding",
        description:
          "grit new myapp --mobile — creates an API + Expo monorepo.",
      },
      {
        title: "Project Structure",
        description:
          "apps/api (Go), apps/expo (React Native + Expo Router), packages/shared (Zod + types).",
      },
      {
        title: "Running the API",
        description:
          "docker compose up -d, cd apps/api && go run cmd/server/main.go.",
      },
      {
        title: "Running Expo",
        description:
          "cd apps/expo && npx expo start — scan QR code or use emulator.",
      },
      {
        title: "Shared Types",
        description:
          "packages/shared with Zod schemas and TypeScript types used by both API and Expo.",
      },
    ],
    challenge: {
      title: "DIY Challenge",
      description:
        'Scaffold a "fitness" mobile app, start both API and Expo, register a user, and navigate through the default screens.',
    },
  },
  {
    number: 2,
    title: "Mobile Auth & Navigation",
    duration: "30 min",
    overview:
      "Implement authentication flows and navigation patterns in React Native with Expo Router.",
    icon: ShieldCheck,
    lessons: [
      {
        title: "Expo Router",
        description:
          "File-based routing for React Native — same mental model as Next.js App Router.",
      },
      {
        title: "Auth Flow",
        description:
          "Login screen, secure token storage with expo-secure-store.",
      },
      {
        title: "Protected Routes",
        description:
          "Auth context provider, automatic redirect to login for unauthenticated users.",
      },
      {
        title: "Tab Navigation",
        description:
          "Bottom tabs for main sections — the primary navigation pattern on mobile.",
      },
      {
        title: "Stack Navigation",
        description:
          "Nested screens within tabs for drill-down flows.",
      },
      {
        title: "JWT Token Refresh",
        description:
          "Automatic token refresh in the API client — seamless re-authentication.",
      },
    ],
    challenge: {
      title: "DIY Challenge",
      description:
        "Build a login/register flow with secure token storage. Create a bottom tab navigation with Home, Workouts, and Profile tabs.",
    },
  },
  {
    number: 3,
    title: "API Integration & Offline",
    duration: "30 min",
    overview:
      "Connect to the Grit API from React Native, handle data fetching, and implement offline support.",
    icon: Cloud,
    lessons: [
      {
        title: "API Client",
        description:
          "Configured Axios/fetch with base URL and auth headers — works identically to web.",
      },
      {
        title: "TanStack Query",
        description:
          "Data fetching hooks in React Native — same API as the web app.",
      },
      {
        title: "Shared Types",
        description:
          "Using packages/shared types in mobile code for full type safety.",
      },
      {
        title: "Pull-to-Refresh",
        description:
          "RefreshControl with query invalidation for native-feeling data refresh.",
      },
      {
        title: "Infinite Scrolling",
        description:
          "FlatList with useInfiniteQuery for paginated lists.",
      },
      {
        title: "Offline Support",
        description:
          "TanStack Query persistence with AsyncStorage cache for offline-first UX.",
      },
    ],
    challenge: {
      title: "DIY Challenge",
      description:
        "Create a Workout resource on the API, build a workout list with pull-to-refresh and infinite scroll. Test offline mode by turning off the network.",
    },
  },
  {
    number: 4,
    title: "Push Notifications",
    duration: "30 min",
    overview:
      "Set up push notifications using Expo and the grit-notifications plugin.",
    icon: Bell,
    lessons: [
      {
        title: "Expo Push Tokens",
        description:
          "Registering device tokens with the Expo push notification service.",
      },
      {
        title: "grit-notifications Plugin",
        description:
          "Storing tokens on the API, sending notifications from Go.",
      },
      {
        title: "Notification Types",
        description:
          "Local notifications vs push notifications — when to use each.",
      },
      {
        title: "Firebase Cloud Messaging",
        description:
          "Setup for Android push notifications via FCM.",
      },
      {
        title: "Apple Push Notification Service",
        description:
          "Setup for iOS push notifications via APNs.",
      },
      {
        title: "Handling Notifications",
        description:
          "Foreground, background, and killed states — different behaviors for each.",
      },
    ],
    challenge: {
      title: "DIY Challenge",
      description:
        "Set up push notifications. Send a notification when a new workout is created. Handle notification tap to navigate to the workout detail screen.",
    },
  },
  {
    number: 5,
    title: "Build & App Store",
    duration: "30 min",
    overview:
      "Build your mobile app for iOS and Android, prepare for App Store and Play Store submission.",
    icon: Store,
    lessons: [
      {
        title: "EAS Build",
        description:
          "Expo Application Services for cloud builds — no local Xcode or Android Studio needed.",
      },
      {
        title: "iOS Build",
        description:
          "Provisioning profiles, certificates, and TestFlight distribution.",
      },
      {
        title: "Android Build",
        description:
          "Signing keys, APK vs AAB formats, and internal testing tracks.",
      },
      {
        title: "App Store Assets",
        description:
          "App icon, screenshots, description, and keywords for the Apple App Store.",
      },
      {
        title: "Play Store Assets",
        description:
          "Feature graphic, screenshots, and store listing for Google Play.",
      },
      {
        title: "OTA Updates",
        description:
          "Expo Updates for over-the-air JavaScript updates — ship fixes without a new build.",
      },
    ],
    challenge: {
      title: "DIY Challenge",
      description:
        "Build your fitness app for your platform (iOS or Android). Install it on your physical device. Take screenshots for a store listing.",
    },
  },
]

/* -- Page component -------------------------------------------------------- */

export default function GritMobileCoursePage() {
  return (
    <div className="min-h-screen bg-[#0b1120]">
      <SiteHeader />

      <main>
        {/* Back link */}
        <div className="container max-w-screen-xl px-6 pt-8">
          <Link
            href="/courses"
            className="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
          >
            <ArrowLeft className="h-4 w-4" />
            Back to All Courses
          </Link>
        </div>

        {/* Hero */}
        <section className="relative overflow-hidden border-b border-border/30">
          <div className="absolute inset-0 bg-gradient-to-b from-primary/[0.04] to-transparent" />
          <div className="container max-w-screen-xl relative py-16 px-6">
            <div className="max-w-3xl mx-auto text-center">
              <div className="inline-flex items-center gap-2 rounded-full bg-primary/10 border border-primary/20 px-3 py-1 text-xs font-mono font-medium text-primary mb-6">
                <Smartphone className="h-3.5 w-3.5" />
                Mobile Track
              </div>

              <h1 className="text-4xl sm:text-5xl font-bold tracking-tight mb-5 leading-[1.1] text-foreground">
                Building Mobile Applications with Grit
              </h1>
              <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl mx-auto mb-8">
                Learn to build cross-platform mobile apps with{" "}
                <span className="text-foreground font-medium">
                  Go API + Expo + React Native
                </span>
                . From scaffolding to the App Store in five self-paced
                mini-courses.
              </p>

              {/* Stats */}
              <div className="flex items-center justify-center gap-6 text-sm text-muted-foreground">
                <span className="flex items-center gap-1.5">
                  <BookOpen className="h-4 w-4 text-primary/60" />
                  5 courses
                </span>
                <span className="flex items-center gap-1.5">
                  <Clock className="h-4 w-4 text-primary/60" />
                  ~2.5 hours
                </span>
                <span className="flex items-center gap-1.5">
                  <Trophy className="h-4 w-4 text-primary/60" />
                  5 challenges
                </span>
              </div>
            </div>
          </div>
        </section>

        {/* Courses */}
        <section className="py-16">
          <div className="container max-w-screen-xl px-6">
            <div className="max-w-3xl mx-auto space-y-6">
              {courses.map((course) => (
                <details
                  key={course.number}
                  className="group rounded-xl border border-border/40 bg-card/50 overflow-hidden"
                >
                  <summary className="flex items-center gap-4 p-6 cursor-pointer select-none hover:bg-card/70 transition-colors list-none [&::-webkit-details-marker]:hidden">
                    {/* Number badge */}
                    <span className="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-primary/10 border border-primary/15 text-sm font-mono font-semibold text-primary">
                      {course.number}
                    </span>

                    <div className="flex-1 min-w-0">
                      <div className="flex items-center gap-2">
                        <h2 className="text-lg font-semibold tracking-tight text-foreground">
                          {course.title}
                        </h2>
                      </div>
                      <p className="text-sm text-muted-foreground mt-0.5">
                        {course.overview}
                      </p>
                    </div>

                    <div className="flex items-center gap-3 shrink-0">
                      <span className="hidden sm:flex items-center gap-1 text-xs text-muted-foreground">
                        <Clock className="h-3.5 w-3.5" />
                        {course.duration}
                      </span>
                      <ChevronDown className="h-5 w-5 text-muted-foreground transition-transform group-open:rotate-180" />
                    </div>
                  </summary>

                  <div className="px-6 pb-6 pt-2 border-t border-border/20">
                    {/* Lessons */}
                    <h3 className="text-xs font-mono font-medium text-muted-foreground uppercase tracking-wider mb-4">
                      Lessons
                    </h3>
                    <ul className="space-y-3 mb-6">
                      {course.lessons.map((lesson, i) => (
                        <li
                          key={i}
                          className="flex items-start gap-3 text-sm"
                        >
                          <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded bg-primary/10 text-[11px] font-mono font-medium text-primary/70 mt-px">
                            {i + 1}
                          </span>
                          <div>
                            <span className="font-medium text-foreground">
                              {lesson.title}
                            </span>
                            <span className="text-muted-foreground">
                              {" "}
                              &mdash; {lesson.description}
                            </span>
                          </div>
                        </li>
                      ))}
                    </ul>

                    {/* DIY Challenge */}
                    <div className="bg-primary/10 border border-primary/20 rounded-lg p-4">
                      <div className="flex items-center gap-2 mb-2">
                        <Trophy className="h-4 w-4 text-primary" />
                        <h4 className="text-sm font-semibold text-foreground">
                          {course.challenge.title}
                        </h4>
                      </div>
                      <p className="text-sm text-muted-foreground leading-relaxed">
                        {course.challenge.description}
                      </p>
                    </div>
                  </div>
                </details>
              ))}
            </div>
          </div>
        </section>
      </main>
    </div>
  )
}
