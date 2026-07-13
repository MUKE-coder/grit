import Link from 'next/link'
import { CheckCircle2 } from 'lucide-react'
import { SiteHeader } from '@/components/site-header'
import { GridFrame } from '@/components/grid-frame'
import { CodeBlock, Challenge, Note, Tip, Definition, Code, CourseNav, CourseFooter } from '@/components/course-components'
import type { Metadata } from 'next'

export const metadata: Metadata = {
  title: 'Build a Fitness App: Go API + Expo React Native',
  description: 'Build a fitness tracking app with a Go backend and Expo React Native frontend. Track workouts, exercises, sets, reps, and weight with a real mobile app.',
}

export default function MobileFitnessAppCourse() {
  return (
    <div className="relative min-h-screen bg-background isolate">
      <SiteHeader />
      <GridFrame />

      <main className="max-w-3xl mx-auto px-6 py-12">
        {/* Breadcrumb */}
        <div className="flex items-center gap-2 text-sm text-muted-foreground mb-8">
          <Link href="/courses" className="hover:text-foreground transition-colors">Courses</Link>
          <span>/</span>
          <span className="text-foreground">Mobile Fitness App</span>
        </div>

        {/* Header */}
        <div className="mb-10">
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs font-mono font-medium text-primary bg-primary/10 px-2 py-0.5 rounded">Standalone Course</span>
            <span className="text-xs text-muted-foreground">~30 min</span>
            <span className="text-xs text-muted-foreground">•</span>
            <span className="text-xs text-muted-foreground">12 challenges</span>
          </div>
          <h1 className="text-3xl md:text-4xl font-bold text-foreground mb-4">
            Build a Fitness App: Go API + Expo React Native
          </h1>
          <p className="text-lg text-muted-foreground leading-relaxed">
            Build a fitness tracking app with a Go backend and Expo React Native frontend. You{"'"}ll
            design a data model for workouts and exercises, generate API resources, build mobile screens
            with React Native, and create a workout logger you can use on your actual phone.
          </p>
        </div>

        <div className="my-4 rounded-lg border border-primary/20 bg-primary/5 px-4 py-3">
          <p className="text-sm text-muted-foreground">
            <strong className="text-foreground">Reference docs:</strong>{' '}
            <a href="/docs/mobile/first-app" className="text-primary hover:underline">Mobile Apps →</a>
          </p>
        </div>

        <hr className="border-border/40 mb-10" />

        {/* ═══ Section 1: What We're Building ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">What We{"'"}re Building</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            We{"'"}re building a fitness tracking app — the kind of app people actually use every day. The
            Go API stores all the data (exercises, workouts, sets, reps, weight). The Expo app runs on
            your phone and lets you log workouts, browse exercises, and view your history.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The app has 4 core screens:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Exercise Library</strong> — browse all exercises grouped by muscle group (chest, back, legs, shoulders, arms)</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Workout Logger</strong> — create a workout, add exercises, log sets/reps/weight in real-time</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Workout History</strong> — see past workouts with total volume and duration</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Dashboard</strong> — weekly summary, workout streak, total volume lifted</li>
          </ul>

          <Definition term="Expo">
            A framework for building React Native apps with zero native configuration. You write
            TypeScript/React code, and Expo handles compiling it for iOS and Android. You can test on
            your physical phone using the Expo Go app — no Xcode or Android Studio required for
            development.
          </Definition>

          <Definition term="React Native">
            A framework that lets you build native mobile apps using React. Unlike a web view or hybrid
            app, React Native renders actual native components (UIView on iOS, android.view on Android).
            The result feels like a native app because it is one.
          </Definition>

          <Challenge number={1} title="Sketch Your Screens">
            <p>On paper or a whiteboard, sketch 4 screens for a fitness app: Exercise Library, Workout
            Logger, Workout History, and Dashboard. For each screen, write down what data it needs to
            display. What actions can the user take on each screen?</p>
          </Challenge>
        </section>

        {/* ═══ Section 2: Scaffold the Project ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Scaffold the Project</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Grit{"'"}s <Code>--mobile</Code> flag scaffolds a monorepo with a Go API and an Expo React
            Native app:
          </p>

          <CodeBlock filename="Terminal">
{`grit new fitness --mobile`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            This creates a monorepo with three parts:
          </p>

          <CodeBlock filename="Project Structure">
{`fitness/
├── docker-compose.yml          # PostgreSQL, Redis, MinIO, Mailhog
├── turbo.json                  # Monorepo task runner
├── pnpm-workspace.yaml         # Workspace definition
├── apps/
│   ├── api/                    # Go backend (Gin + GORM)
│   │   ├── cmd/server/main.go
│   │   └── internal/           # models, handlers, services, middleware
│   └── expo/                   # Expo React Native app
│       ├── app/                # File-based routing (Expo Router)
│       ├── components/         # Reusable React Native components
│       ├── hooks/              # Custom hooks (useAuth, useQuery)
│       ├── lib/                # API client, storage, utils
│       ├── app.json            # Expo config
│       └── package.json
└── packages/
    └── shared/                 # Shared types and validation (Zod)`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The key difference from a web project: instead of <Code>apps/web</Code> (Next.js), you
            get <Code>apps/expo</Code> (Expo React Native). The Go API is identical — the mobile app
            consumes the same REST endpoints.
          </p>

          <Note>
            The <Code>packages/shared</Code> directory contains TypeScript types and Zod schemas
            shared between the API types and the Expo app. When you generate a resource, Grit updates
            the shared types so both the API and the mobile app stay in sync.
          </Note>

          <Challenge number={2} title="Scaffold and Explore">
            <p>Run <Code>grit new fitness --mobile</Code>. Open the project in your editor. Compare
            the folder structure to a web project (<Code>grit new myapp</Code>). What{"'"}s the same?
            What{"'"}s different? How many files are in <Code>apps/expo/</Code>?</p>
          </Challenge>
        </section>

        {/* ═══ Section 3: Design the Data Model ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Design the Data Model</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Before writing code, let{"'"}s design the data model. A fitness app needs three core entities:
          </p>

          <Definition term="Data Model">
            The structure of your application{"'"}s data — what entities exist, what fields they have,
            and how they relate to each other. A good data model is the foundation of every application.
            Get it right, and everything else is easier. Get it wrong, and you{"'"}ll fight your database
            forever.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Our three resources:
          </p>

          <div className="overflow-x-auto mb-4">
            <table className="w-full text-sm border border-border/40 rounded-lg overflow-hidden">
              <thead>
                <tr className="bg-muted/20">
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Resource</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Fields</th>
                  <th className="text-left px-3 py-2.5 font-semibold text-foreground border-b border-border/40">Relationships</th>
                </tr>
              </thead>
              <tbody className="text-muted-foreground">
                <tr className="border-b border-border/20">
                  <td className="px-3 py-2 font-medium text-foreground">Exercise</td>
                  <td className="px-3 py-2">name, muscle_group, description (optional)</td>
                  <td className="px-3 py-2">Has many WorkoutExercises</td>
                </tr>
                <tr className="border-b border-border/20">
                  <td className="px-3 py-2 font-medium text-foreground">Workout</td>
                  <td className="px-3 py-2">name, date, duration (minutes), notes (optional)</td>
                  <td className="px-3 py-2">Has many WorkoutExercises</td>
                </tr>
                <tr>
                  <td className="px-3 py-2 font-medium text-foreground">WorkoutExercise</td>
                  <td className="px-3 py-2">sets, reps, weight (kg/lbs)</td>
                  <td className="px-3 py-2">Belongs to Workout, belongs to Exercise</td>
                </tr>
              </tbody>
            </table>
          </div>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The relationship is: a <strong className="text-foreground">Workout</strong> contains
            many <strong className="text-foreground">WorkoutExercises</strong>, and
            each <strong className="text-foreground">WorkoutExercise</strong> references
            an <strong className="text-foreground">Exercise</strong>. This is a classic many-to-many
            relationship through a join table (WorkoutExercise).
          </p>

          <CodeBlock filename="Generate All 3 Resources">
{`# Exercise — the library of all possible exercises
grit generate resource Exercise --fields "name:string,muscle_group:string,description:text:optional"

# Workout — a single training session
grit generate resource Workout --fields "name:string,date:date,duration:int,notes:text:optional"

# WorkoutExercise — a specific exercise performed in a workout
grit generate resource WorkoutExercise --fields "sets:int,reps:int,weight:float,workout_id:belongs_to:Workout,exercise_id:belongs_to:Exercise"`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            After generating, restart the API. GORM will auto-migrate the new tables.
          </p>

          <Challenge number={3} title="Generate the Resources">
            <p>Generate all 3 resources using the commands above. Restart the API. Open GORM Studio
            at <Code>/studio</Code> and verify that the <Code>exercises</Code>, <Code>workouts</Code>,
            and <Code>workout_exercises</Code> tables exist. How many columns does
            the <Code>workout_exercises</Code> table have?</p>
          </Challenge>
        </section>

        {/* ═══ Section 4: Seed Sample Data ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Seed Sample Data</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            An empty exercise library is useless. Let{"'"}s populate it with real exercises. You can do this
            through GORM Studio{"'"}s inline editing or through the API.
          </p>

          <Definition term="Seed Data">
            Initial data that your application needs to function. For a fitness app, the exercise library
            is seed data — users don{"'"}t create exercises, they pick from a pre-populated list. Seed data
            is typically loaded once during setup and rarely changes.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here are the exercises to create, organized by muscle group:
          </p>

          <CodeBlock filename="Seed Exercises via API">
{`# Chest
POST /api/exercises  {"name": "Bench Press", "muscle_group": "Chest", "description": "Flat barbell bench press"}
POST /api/exercises  {"name": "Incline Dumbbell Press", "muscle_group": "Chest", "description": "Incline bench with dumbbells"}
POST /api/exercises  {"name": "Cable Fly", "muscle_group": "Chest", "description": "Cable crossover fly movement"}

# Back
POST /api/exercises  {"name": "Deadlift", "muscle_group": "Back", "description": "Conventional barbell deadlift"}
POST /api/exercises  {"name": "Pull-up", "muscle_group": "Back", "description": "Bodyweight pull-up (add weight as needed)"}
POST /api/exercises  {"name": "Barbell Row", "muscle_group": "Back", "description": "Bent-over barbell row"}

# Legs
POST /api/exercises  {"name": "Squat", "muscle_group": "Legs", "description": "Barbell back squat"}
POST /api/exercises  {"name": "Romanian Deadlift", "muscle_group": "Legs", "description": "Stiff-leg deadlift for hamstrings"}
POST /api/exercises  {"name": "Leg Press", "muscle_group": "Legs", "description": "Machine leg press"}

# Shoulders
POST /api/exercises  {"name": "Overhead Press", "muscle_group": "Shoulders", "description": "Standing barbell press"}`}
          </CodeBlock>

          <Tip>
            You can also create these through the interactive API docs at <Code>/docs</Code>. Click
            the POST /api/exercises endpoint, fill in the JSON body, and hit Execute. It{"'"}s faster
            than curl for manual data entry.
          </Tip>

          <Challenge number={4} title="Seed 10 Exercises">
            <p>Create at least 10 exercises across 4 different muscle groups (Chest, Back, Legs,
            Shoulders). Use the API docs or GORM Studio. After creating them, call <Code>GET
            /api/exercises</Code> and verify all 10 are returned. Try filtering
            with <Code>?muscle_group=Chest</Code> — does it work?</p>
          </Challenge>
        </section>

        {/* ═══ Section 5: Start API + Expo ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Start API + Expo</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            With the API running and seed data loaded, it{"'"}s time to start the Expo mobile app.
            You{"'"}ll need two terminals — one for the API and one for Expo.
          </p>

          <CodeBlock filename="Terminal 1 — Start Infrastructure + API">
{`cd fitness
docker compose up -d
cd apps/api && go run cmd/server/main.go`}
          </CodeBlock>

          <CodeBlock filename="Terminal 2 — Start Expo">
{`cd fitness/apps/expo
pnpm install
pnpm start`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Expo will show a QR code in the terminal. To open the app on your phone:
          </p>

          <ul className="space-y-2 text-muted-foreground mb-4">
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">iPhone</strong> — install Expo Go from the App Store, scan the QR code with your camera</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Android</strong> — install Expo Go from Google Play, scan the QR code from the Expo Go app</li>
            <li className="flex gap-2"><span className="text-primary">•</span> <strong className="text-foreground">Emulator</strong> — press <Code>i</Code> for iOS Simulator or <Code>a</Code> for Android Emulator</li>
          </ul>

          <Note>
            Your phone and computer must be on the same Wi-Fi network for Expo to connect. If you{"'"}re
            on a corporate network that blocks local connections, use <Code>pnpm start --tunnel</Code> to
            route through Expo{"'"}s servers instead.
          </Note>

          <Challenge number={5} title="Run the App">
            <p>Start both the API and Expo. Open the app on your phone or emulator. You should see
            the default Grit mobile app with login/register screens. Register a user, log in, and
            confirm you see the main app screen. Is the API responding to requests from the mobile app?</p>
          </Challenge>
        </section>

        {/* ═══ Section 6: Build the Exercise List Screen ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Build the Exercise List Screen</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The exercise list is the first screen users see. It fetches exercises from the API and
            displays them in a scrollable list, grouped by muscle group.
          </p>

          <Definition term="FlatList">
            React Native{"'"}s performant list component. Unlike mapping over an array and rendering all
            items at once (which crashes with large lists), FlatList only renders items that are visible
            on screen. As the user scrolls, new items are rendered and off-screen items are recycled.
          </Definition>

          <Definition term="useQuery (TanStack Query)">
            A React hook for fetching, caching, and synchronizing server data. It handles loading states,
            error states, caching, background refetching, and stale data — all the things you{"'"}d
            normally write manually. You provide a query key and a fetch function, and useQuery handles
            the rest.
          </Definition>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s how to fetch and display exercises:
          </p>

          <CodeBlock filename="app/(tabs)/exercises.tsx">
{`import { View, Text, FlatList, TouchableOpacity } from 'react-native'
import { useQuery } from '@tanstack/react-query'
import { api } from '@/lib/api'

interface Exercise {
  id: number
  name: string
  muscle_group: string
  description?: string
}

export default function ExercisesScreen() {
  const { data, isLoading } = useQuery({
    queryKey: ['exercises'],
    queryFn: () => api.get('/api/exercises'),
  })

  const exercises: Exercise[] = data?.data || []

  // Group by muscle_group
  const grouped = exercises.reduce((acc, ex) => {
    if (!acc[ex.muscle_group]) acc[ex.muscle_group] = []
    acc[ex.muscle_group].push(ex)
    return acc
  }, {} as Record<string, Exercise[]>)

  if (isLoading) return <Text>Loading...</Text>

  return (
    <FlatList
      data={Object.entries(grouped)}
      keyExtractor={([group]) => group}
      renderItem={({ item: [group, items] }) => (
        <View style={{ padding: 16 }}>
          <Text style={{ fontSize: 18, fontWeight: 'bold', marginBottom: 8 }}>
            {group}
          </Text>
          {items.map((ex) => (
            <TouchableOpacity key={ex.id} style={{ padding: 12, marginBottom: 4 }}>
              <Text style={{ fontSize: 16 }}>{ex.name}</Text>
              {ex.description && (
                <Text style={{ color: '#666', fontSize: 13 }}>{ex.description}</Text>
              )}
            </TouchableOpacity>
          ))}
        </View>
      )}
    />
  )
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The pattern is straightforward: <Code>useQuery</Code> fetches the data, we group it by
            muscle group using <Code>reduce</Code>, and <Code>FlatList</Code> renders each group with
            its exercises.
          </p>

          <Challenge number={6} title="Build the Exercise List">
            <p>Create the exercise list screen. Display exercises grouped by muscle group. Each exercise
            should show its name and description. Verify that all the exercises you seeded appear. How
            many muscle groups are displayed?</p>
          </Challenge>
        </section>

        {/* ═══ Section 7: Build the Workout Logger ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Build the Workout Logger</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The workout logger is the core feature of the app. Users create a workout, add exercises
            to it, and log sets with reps and weight for each exercise. This involves creating
            a <Code>Workout</Code> record and then creating <Code>WorkoutExercise</Code> records
            linked to it.
          </p>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The flow is:
          </p>

          <CodeBlock filename="Workout Logging Flow">
{`1. User taps "Start Workout"
   → POST /api/workouts { name: "Push Day", date: "2026-03-27" }
   → Returns workout with ID

2. User picks exercises from the library
   → Each selection adds to a local list

3. For each exercise, user logs sets:
   → POST /api/workout-exercises {
       workout_id: 1,
       exercise_id: 3,
       sets: 4,
       reps: 10,
       weight: 80.0
     }

4. User taps "Finish Workout"
   → PATCH /api/workouts/1 { duration: 55 }
   → Duration calculated from start to finish`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The workout form combines API calls with local state. The workout itself is created on the
            server first (so we get an ID), then each exercise entry is posted as the user logs it:
          </p>

          <CodeBlock filename="Workout Logger Component">
{`import { useState } from 'react'
import { View, Text, TextInput, TouchableOpacity, ScrollView } from 'react-native'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '@/lib/api'

interface SetEntry {
  exercise_id: number
  exercise_name: string
  sets: number
  reps: number
  weight: number
}

export default function WorkoutLoggerScreen() {
  const [workoutId, setWorkoutId] = useState<number | null>(null)
  const [entries, setEntries] = useState<SetEntry[]>([])
  const [startTime] = useState(new Date())
  const queryClient = useQueryClient()

  const createWorkout = useMutation({
    mutationFn: (name: string) =>
      api.post('/api/workouts', {
        name,
        date: new Date().toISOString().split('T')[0],
        duration: 0,
      }),
    onSuccess: (data) => setWorkoutId(data.data.id),
  })

  const logExercise = useMutation({
    mutationFn: (entry: SetEntry) =>
      api.post('/api/workout-exercises', {
        workout_id: workoutId,
        exercise_id: entry.exercise_id,
        sets: entry.sets,
        reps: entry.reps,
        weight: entry.weight,
      }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['workouts'] }),
  })

  // ... render workout form with exercise picker and set inputs
}`}
          </CodeBlock>

          <Tip>
            Use <Code>useMutation</Code> from TanStack Query for POST/PATCH/DELETE operations. It
            handles loading states, error states, and cache invalidation — just
            like <Code>useQuery</Code> does for GET requests.
          </Tip>

          <Challenge number={7} title="Log a Complete Workout">
            <p>Build the workout logger screen. Create a workout called {'"'}Push Day.{'"'} Add 4 exercises
            (e.g., Bench Press, Incline Dumbbell Press, Cable Fly, Overhead Press). For each exercise,
            log sets, reps, and weight. When you{"'"}re done, verify the data exists in the API by
            calling <Code>GET /api/workout-exercises?workout_id=1</Code>.</p>
          </Challenge>
        </section>

        {/* ═══ Section 8: Build the History Screen ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Build the History Screen</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The history screen shows all past workouts with summary data: name, date, duration, and
            total volume. <strong className="text-foreground">Volume</strong> is the key metric in
            fitness tracking — it{"'"}s the total weight lifted in a workout.
          </p>

          <Definition term="Volume">
            The total amount of weight lifted in a workout, calculated as: sets x reps x weight for
            each exercise, summed across all exercises. For example, 4 sets of 10 reps at 80kg = 3,200kg
            volume. Tracking volume over time shows strength progress.
          </Definition>

          <CodeBlock filename="app/(tabs)/history.tsx">
{`import { View, Text, FlatList } from 'react-native'
import { useQuery } from '@tanstack/react-query'
import { api } from '@/lib/api'

interface Workout {
  id: number
  name: string
  date: string
  duration: number
  workout_exercises: {
    sets: number
    reps: number
    weight: number
    exercise: { name: string }
  }[]
}

function calculateVolume(workout: Workout): number {
  return workout.workout_exercises.reduce(
    (total, we) => total + we.sets * we.reps * we.weight,
    0
  )
}

export default function HistoryScreen() {
  const { data } = useQuery({
    queryKey: ['workouts'],
    queryFn: () => api.get('/api/workouts?sort=date&order=desc'),
  })

  const workouts: Workout[] = data?.data || []

  return (
    <FlatList
      data={workouts}
      keyExtractor={(w) => String(w.id)}
      renderItem={({ item: workout }) => (
        <View style={{ padding: 16, borderBottomWidth: 1, borderColor: '#222' }}>
          <Text style={{ fontSize: 18, fontWeight: 'bold' }}>{workout.name}</Text>
          <Text style={{ color: '#666' }}>{workout.date}</Text>
          <View style={{ flexDirection: 'row', gap: 16, marginTop: 8 }}>
            <Text style={{ color: '#999' }}>{workout.duration} min</Text>
            <Text style={{ color: '#6c5ce7' }}>
              {calculateVolume(workout).toLocaleString()} kg volume
            </Text>
          </View>
          <Text style={{ color: '#888', marginTop: 4 }}>
            {workout.workout_exercises.length} exercises
          </Text>
        </View>
      )}
    />
  )
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            The <Code>calculateVolume</Code> function sums <Code>sets * reps * weight</Code> for every
            exercise in the workout. This gives users a clear, numeric measure of how hard they trained.
          </p>

          <Challenge number={8} title="Build the History Screen">
            <p>Create the history screen. Log at least 5 workouts over a {'"'}week{'"'} (you can
            backdate them by setting different dates in the API). Display them in reverse chronological
            order (newest first). Each entry should show: workout name, date, duration, number of
            exercises, and total volume. Which workout had the highest volume?</p>
          </Challenge>
        </section>

        {/* ═══ Section 9: Add Pull-to-Refresh ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Add Pull-to-Refresh</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Mobile users expect to pull down on a list to refresh it. This is a standard mobile pattern
            that React Native supports natively through the <Code>RefreshControl</Code> component.
          </p>

          <Definition term="Pull-to-Refresh">
            A mobile interaction pattern where the user drags a list downward past its top edge to
            trigger a data reload. A spinner appears while the data is refreshing, then disappears
            when the new data is loaded. Nearly every mobile app with lists supports this gesture.
          </Definition>

          <CodeBlock filename="Adding RefreshControl to a FlatList">
{`import { FlatList, RefreshControl } from 'react-native'
import { useQuery } from '@tanstack/react-query'

export default function ExercisesScreen() {
  const { data, isLoading, refetch, isRefetching } = useQuery({
    queryKey: ['exercises'],
    queryFn: () => api.get('/api/exercises'),
  })

  return (
    <FlatList
      data={data?.data || []}
      keyExtractor={(item) => String(item.id)}
      renderItem={({ item }) => (
        // ... render exercise
      )}
      refreshControl={
        <RefreshControl
          refreshing={isRefetching}
          onRefresh={refetch}
          tintColor="#6c5ce7"
        />
      }
    />
  )
}`}
          </CodeBlock>

          <p className="text-muted-foreground leading-relaxed mb-4">
            TanStack Query makes this trivial. The <Code>refetch</Code> function triggers a fresh
            API call, and <Code>isRefetching</Code> tracks whether the refetch is in progress. Wire
            these into <Code>RefreshControl</Code> and you{"'"}re done.
          </p>

          <Challenge number={9} title="Add Pull-to-Refresh">
            <p>Add <Code>RefreshControl</Code> to both the Exercise List and Workout History screens.
            Pull down on each list — does the spinner appear? Does the data reload? Test it by adding
            a new exercise through the API docs and pulling to refresh on the exercise list. Does the
            new exercise appear?</p>
          </Challenge>
        </section>

        {/* ═══ Section 10: Summary ═══ */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold text-foreground mb-4">Summary</h2>

          <p className="text-muted-foreground leading-relaxed mb-4">
            Here{"'"}s everything you learned in this course:
          </p>

          <ul className="space-y-3 mb-6">
            {[
              'Grit --mobile scaffolds a monorepo with Go API + Expo React Native',
              'Data modeling with three resources: Exercise, Workout, WorkoutExercise',
              'belongs_to relationships connect WorkoutExercises to Workouts and Exercises',
              'Seed data provides an exercise library users can pick from',
              'FlatList renders performant scrollable lists in React Native',
              'useQuery (TanStack Query) handles fetching, caching, and loading states',
              'useMutation handles POST/PATCH/DELETE with cache invalidation',
              'Volume (sets x reps x weight) is the key metric for tracking progress',
              'RefreshControl adds pull-to-refresh to any FlatList',
              'Expo Go lets you test on a physical device without native toolchains',
            ].map((item, i) => (
              <li key={i} className="flex gap-3 text-muted-foreground">
                <CheckCircle2 className="w-5 h-5 text-primary shrink-0 mt-0.5" />
                <span>{item}</span>
              </li>
            ))}
          </ul>

          <Challenge number={10} title="Personal Records (Part 1)">
            <p>Design a Personal Records (PR) feature. For each exercise, track the highest weight
            the user has ever lifted. How would you calculate this? You don{"'"}t need a new database
            table — you can derive PRs from existing WorkoutExercise data. Write the API query that
            finds the max weight for each exercise.</p>
          </Challenge>

          <Challenge number={11} title="Personal Records (Part 2)">
            <p>Build a PR display screen. For each exercise the user has performed, show the exercise
            name, the PR weight, and the date it was set. Sort by most recent PR first.</p>
          </Challenge>

          <Challenge number={12} title="Personal Records (Part 3)">
            <p>Add a {'"'}New PR!{'"'} badge to the workout logger. When a user logs a weight that
            exceeds their previous best for that exercise, display a visual indicator (a badge, animation,
            or color change). Compare the current weight against the max weight from all previous
            WorkoutExercise records for that exercise. This is the final touch — congratulate users
            when they beat their records.</p>
          </Challenge>
        </section>

        {/* Footer */}
        <CourseFooter />

        {/* Navigation */}
        <div className="mt-8">
          <CourseNav
            prev={{ href: '/courses', label: 'All Courses' }}
            next={{ href: '/courses', label: 'More Courses' }}
          />
        </div>
      </main>
    </div>
  )
}
