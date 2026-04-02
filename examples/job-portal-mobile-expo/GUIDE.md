# Job Portal -- Mobile (Expo) Build Guide

A step-by-step guide to building a mobile job portal with Grit.
Go API on the backend, Expo React Native on the frontend, running on your phone.

## Step 1: Scaffold the Project

```bash
grit new job-portal --mobile
cd job-portal
```

The `--mobile` flag tells Grit to scaffold a Go API plus an Expo React Native app:

```
job-portal/
├── apps/
│   ├── api/                # Go backend
│   └── expo/               # Expo React Native app
│       ├── app/            # Expo Router (file-based navigation)
│       ├── components/     # Reusable components
│       ├── hooks/          # Custom hooks (auth, data)
│       ├── lib/            # API client, SecureStore
│       └── app.json        # Expo config
├── packages/shared/        # Shared types + Zod schemas
├── turbo.json
├── pnpm-workspace.yaml
└── docker-compose.yml
```

## Step 2: Start Infrastructure

```bash
docker compose up -d
```

This starts PostgreSQL, Redis, MinIO, and Mailhog.

## Step 3: Install Dependencies

```bash
pnpm install
```

This installs dependencies for both the Expo app and the shared package.

## Step 4: Generate Resources

```bash
grit generate resource Category \
  --fields "name:string:unique,slug:slug:name"

grit generate resource Company \
  --fields "name:string,slug:slug:name,description:richtext,logo:string,website:string:optional,email:string"

grit generate resource Job \
  --fields "title:string,description:richtext,company_id:belongs_to:Company,location:string,salary_min:float,salary_max:float,type:string,experience:string,active:bool,deadline:date"

grit generate resource Application \
  --fields "job_id:belongs_to:Job,applicant_name:string,applicant_email:string,resume:string,cover_letter:text:optional,status:string"
```

With `--mobile`, each `grit generate resource` creates:
- **Go files**: model, service, handler (same as all architectures)
- **Shared types**: TypeScript interface in `packages/shared/types/`
- **Shared schemas**: Zod validation in `packages/shared/schemas/`
- **Expo hook**: TanStack Query hook in `apps/expo/hooks/`

## Step 5: Run Migrations and Seed

```bash
grit migrate
grit seed
```

Seeds an admin user: `admin@example.com` / `password`.

## Step 6: Start the API

```bash
cd apps/api && go run cmd/server/main.go
```

Verify it works: http://localhost:8080/docs

## Step 7: Start Expo

In a second terminal:

```bash
cd apps/expo && npx expo start
```

You will see a QR code in your terminal. Open the Expo Go app on your phone
and scan it.

### Important: Local IP for Physical Devices

If you are running on a physical phone (not a simulator), the app cannot reach
`localhost` because that refers to the phone itself. Find your computer's local
IP:

```bash
# macOS/Linux
ifconfig | grep "inet " | grep -v 127.0.0.1

# Windows
ipconfig | findstr "IPv4"
```

Update `EXPO_PUBLIC_API_URL` in `.env`:
```
EXPO_PUBLIC_API_URL=http://192.168.1.100:8080
```

Restart Expo after changing the env.

## Step 8: Understand Expo Router

Expo Router uses file-based routing, similar to Next.js App Router. The file
structure maps directly to screens:

```
app/
├── _layout.tsx              # Root layout (QueryClientProvider, fonts, auth)
├── (auth)/                  # Auth group (no tab bar)
│   ├── _layout.tsx          # Stack navigator for auth screens
│   ├── login.tsx            # Login screen
│   └── register.tsx         # Register screen
└── (tabs)/                  # Tab group (tab bar at bottom)
    ├── _layout.tsx          # Tab navigator config
    ├── index.tsx            # Home tab
    ├── jobs.tsx             # Jobs tab
    ├── applications.tsx     # My Applications tab
    └── profile.tsx          # Profile tab
```

### Root Layout (`app/_layout.tsx`)

The root layout wraps the entire app with providers:

```typescript
import { QueryClientProvider } from "@tanstack/react-query";
import { AuthProvider } from "../hooks/useAuth";
import { Stack } from "expo-router";

export default function RootLayout() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <Stack screenOptions={{ headerShown: false }}>
          <Stack.Screen name="(auth)" />
          <Stack.Screen name="(tabs)" />
        </Stack>
      </AuthProvider>
    </QueryClientProvider>
  );
}
```

### Tab Layout (`app/(tabs)/_layout.tsx`)

The tab layout configures the bottom tab bar:

```typescript
import { Tabs } from "expo-router";
import { Home, Briefcase, FileText, User } from "lucide-react-native";

export default function TabLayout() {
  return (
    <Tabs screenOptions={{ tabBarActiveTintColor: "#6c5ce7" }}>
      <Tabs.Screen
        name="index"
        options={{ title: "Home", tabBarIcon: ({ color }) => <Home color={color} /> }}
      />
      <Tabs.Screen
        name="jobs"
        options={{ title: "Jobs", tabBarIcon: ({ color }) => <Briefcase color={color} /> }}
      />
      <Tabs.Screen
        name="applications"
        options={{ title: "Applications", tabBarIcon: ({ color }) => <FileText color={color} /> }}
      />
      <Tabs.Screen
        name="profile"
        options={{ title: "Profile", tabBarIcon: ({ color }) => <User color={color} /> }}
      />
    </Tabs>
  );
}
```

## Step 9: Build the Job Listing Screen

The Jobs tab uses `FlatList` with TanStack Query's `useInfiniteQuery` for
paginated loading:

```typescript
// app/(tabs)/jobs.tsx
import { FlatList, RefreshControl, View, Text } from "react-native";
import { useInfiniteQuery } from "@tanstack/react-query";
import { fetchJobs } from "../../lib/api";
import { JobCard } from "../../components/JobCard";
import { SearchBar } from "../../components/SearchBar";

export default function JobsScreen() {
  const [search, setSearch] = useState("");

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isRefetching,
    refetch,
  } = useInfiniteQuery({
    queryKey: ["jobs", search],
    queryFn: ({ pageParam = 1 }) =>
      fetchJobs({ page: pageParam, search }),
    getNextPageParam: (lastPage) =>
      lastPage.meta.page < lastPage.meta.pages
        ? lastPage.meta.page + 1
        : undefined,
  });

  const jobs = data?.pages.flatMap((page) => page.data) ?? [];

  return (
    <View style={{ flex: 1 }}>
      <SearchBar value={search} onChangeText={setSearch} />
      <FlatList
        data={jobs}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => <JobCard job={item} />}
        onEndReached={() => hasNextPage && fetchNextPage()}
        onEndReachedThreshold={0.5}
        refreshControl={
          <RefreshControl refreshing={isRefetching} onRefresh={refetch} />
        }
        ListFooterComponent={
          isFetchingNextPage ? <Text>Loading more...</Text> : null
        }
      />
    </View>
  );
}
```

Pull down to refresh. Scroll to the bottom to load more. Works offline with
cached data from TanStack Query.

## Step 10: Build the Job Detail Screen

Create a dynamic route for job details:

```typescript
// app/jobs/[id].tsx
import { useLocalSearchParams } from "expo-router";
import { useQuery } from "@tanstack/react-query";
import { ScrollView, View, Text, Pressable } from "react-native";
import { fetchJob } from "../../lib/api";
import { router } from "expo-router";

export default function JobDetailScreen() {
  const { id } = useLocalSearchParams<{ id: string }>();
  const { data: job, isLoading } = useQuery({
    queryKey: ["jobs", id],
    queryFn: () => fetchJob(id),
  });

  if (isLoading) return <Text>Loading...</Text>;

  return (
    <ScrollView style={{ flex: 1, padding: 16 }}>
      <Text style={{ fontSize: 24, fontWeight: "bold" }}>{job.title}</Text>
      <Text style={{ color: "#9090a8" }}>{job.company?.name}</Text>
      <Text>{job.location}</Text>
      <Text>
        ${job.salary_min.toLocaleString()} - ${job.salary_max.toLocaleString()}
      </Text>

      <Pressable
        style={{ backgroundColor: "#6c5ce7", padding: 16, borderRadius: 8, marginTop: 24 }}
        onPress={() => router.push(`/apply/${job.id}`)}
      >
        <Text style={{ color: "white", textAlign: "center", fontWeight: "600" }}>
          Apply Now
        </Text>
      </Pressable>
    </ScrollView>
  );
}
```

## Step 11: Build the Application Form with Resume Upload

The apply screen lets users upload a resume from their device:

```typescript
// app/apply/[jobId].tsx
import * as DocumentPicker from "expo-document-picker";
import { useMutation } from "@tanstack/react-query";
import { uploadFile, submitApplication } from "../../lib/api";
import { useAuth } from "../../hooks/useAuth";

export default function ApplyScreen() {
  const { jobId } = useLocalSearchParams<{ jobId: string }>();
  const { user } = useAuth();
  const [resumeUri, setResumeUri] = useState<string | null>(null);
  const [coverLetter, setCoverLetter] = useState("");

  const pickResume = async () => {
    const result = await DocumentPicker.getDocumentAsync({
      type: ["application/pdf", "application/msword"],
    });
    if (!result.canceled) {
      setResumeUri(result.assets[0].uri);
    }
  };

  const applyMutation = useMutation({
    mutationFn: async () => {
      // 1. Get presigned URL
      // 2. Upload file
      const fileUrl = await uploadFile(resumeUri!);
      // 3. Submit application
      return submitApplication({
        job_id: jobId,
        applicant_name: user.name,
        applicant_email: user.email,
        resume: fileUrl,
        cover_letter: coverLetter,
      });
    },
    onSuccess: () => router.replace("/(tabs)/applications"),
  });

  return (
    <ScrollView style={{ flex: 1, padding: 16 }}>
      <Pressable onPress={pickResume} style={styles.uploadButton}>
        <Text>{resumeUri ? "Resume selected" : "Pick resume (PDF/DOC)"}</Text>
      </Pressable>

      <TextInput
        placeholder="Cover letter (optional)"
        value={coverLetter}
        onChangeText={setCoverLetter}
        multiline
        numberOfLines={6}
        style={styles.textArea}
      />

      <Pressable
        onPress={() => applyMutation.mutate()}
        disabled={!resumeUri || applyMutation.isPending}
        style={styles.submitButton}
      >
        <Text style={{ color: "white" }}>
          {applyMutation.isPending ? "Submitting..." : "Submit Application"}
        </Text>
      </Pressable>
    </ScrollView>
  );
}
```

## Step 12: SecureStore Auth Flow

The auth hook manages tokens in SecureStore, not AsyncStorage:

```typescript
// hooks/useAuth.ts
import * as SecureStore from "expo-secure-store";
import { createContext, useContext, useEffect, useState } from "react";
import { api } from "../lib/api";
import { router } from "expo-router";

interface AuthState {
  user: User | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
}

export function AuthProvider({ children }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // On mount, check for stored tokens
  useEffect(() => {
    (async () => {
      const token = await SecureStore.getItemAsync("access_token");
      if (token) {
        api.defaults.headers.common["Authorization"] = `Bearer ${token}`;
        try {
          const { data } = await api.get("/api/auth/me");
          setUser(data.data);
        } catch {
          await SecureStore.deleteItemAsync("access_token");
          await SecureStore.deleteItemAsync("refresh_token");
        }
      }
      setIsLoading(false);
    })();
  }, []);

  const login = async (email: string, password: string) => {
    const { data } = await api.post("/api/auth/login", { email, password });
    await SecureStore.setItemAsync("access_token", data.data.access_token);
    await SecureStore.setItemAsync("refresh_token", data.data.refresh_token);
    api.defaults.headers.common["Authorization"] = `Bearer ${data.data.access_token}`;
    setUser(data.data.user);
    router.replace("/(tabs)");
  };

  const logout = async () => {
    await SecureStore.deleteItemAsync("access_token");
    await SecureStore.deleteItemAsync("refresh_token");
    delete api.defaults.headers.common["Authorization"];
    setUser(null);
    router.replace("/(auth)/login");
  };

  // ... register similar to login

  return (
    <AuthContext.Provider value={{ user, isLoading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}
```

The API client includes an interceptor that refreshes expired tokens automatically:

```typescript
// lib/api.ts
import axios from "axios";
import * as SecureStore from "expo-secure-store";

export const api = axios.create({
  baseURL: process.env.EXPO_PUBLIC_API_URL,
});

api.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      const refreshToken = await SecureStore.getItemAsync("refresh_token");
      if (refreshToken) {
        const { data } = await api.post("/api/auth/refresh", {
          refresh_token: refreshToken,
        });
        await SecureStore.setItemAsync("access_token", data.data.access_token);
        api.defaults.headers.common["Authorization"] =
          `Bearer ${data.data.access_token}`;
        return api.request(error.config);
      }
    }
    return Promise.reject(error);
  }
);
```

## Step 13: Build the Profile Screen

```typescript
// app/(tabs)/profile.tsx
import { View, Text, Pressable, Switch } from "react-native";
import { useAuth } from "../../hooks/useAuth";

export default function ProfileScreen() {
  const { user, logout } = useAuth();

  if (!user) {
    return (
      <View style={{ flex: 1, justifyContent: "center", alignItems: "center" }}>
        <Text>Sign in to view your profile</Text>
        <Pressable onPress={() => router.push("/(auth)/login")}>
          <Text style={{ color: "#6c5ce7" }}>Sign In</Text>
        </Pressable>
      </View>
    );
  }

  return (
    <View style={{ flex: 1, padding: 16 }}>
      <Text style={{ fontSize: 24, fontWeight: "bold" }}>{user.name}</Text>
      <Text style={{ color: "#9090a8" }}>{user.email}</Text>
      <Text style={{ color: "#9090a8" }}>Role: {user.role}</Text>

      <View style={{ marginTop: 32, gap: 16 }}>
        <Pressable onPress={() => router.push("/settings/2fa")}>
          <Text>Two-Factor Authentication</Text>
        </Pressable>

        <Pressable onPress={() => router.push("/settings/notifications")}>
          <Text>Push Notifications</Text>
        </Pressable>

        <Pressable onPress={logout}>
          <Text style={{ color: "#ff6b6b" }}>Sign Out</Text>
        </Pressable>
      </View>
    </View>
  );
}
```

## Step 14: Add Push Notifications

Register for push notifications after login:

```typescript
// hooks/usePushNotifications.ts
import * as Notifications from "expo-notifications";
import * as Device from "expo-device";
import { api } from "../lib/api";
import { Platform } from "react-native";

export async function registerPushNotifications() {
  if (!Device.isDevice) return; // Push only works on real devices

  const { status } = await Notifications.requestPermissionsAsync();
  if (status !== "granted") return;

  const token = await Notifications.getExpoPushTokenAsync({
    projectId: "your-expo-project-id",
  });

  // Send token to the API so it can send push notifications
  await api.post("/api/users/push-token", {
    token: token.data,
    platform: Platform.OS,
  });
}

// Configure notification handling
Notifications.setNotificationHandler({
  handleNotification: async () => ({
    shouldShowAlert: true,
    shouldPlaySound: true,
    shouldSetBadge: true,
  }),
});
```

Call `registerPushNotifications()` in the auth flow after a successful login.

The API sends push notifications via Expo's push API when application status
changes. Users get a notification like "Your application for Senior Go Developer
has been moved to Interview stage."

## Step 15: Deploy

### Deploy the API

Same as any Grit project:

```bash
grit deploy --host deploy@server.com --domain api.jobs.example.com
```

Or with Docker:

```bash
docker compose -f docker-compose.prod.yml up -d
```

### Build the Mobile App

```bash
cd apps/expo

# Install EAS CLI
npm install -g eas-cli

# Login to Expo
eas login

# Configure EAS Build
eas build:configure

# Set production API URL
# In .env: EXPO_PUBLIC_API_URL=https://api.jobs.example.com

# Build for both platforms
eas build --platform ios
eas build --platform android

# Submit to App Store and Google Play
eas submit --platform ios
eas submit --platform android
```

### Over-the-Air Updates

After the initial store release, push JS bundle updates without going through
app review:

```bash
eas update --branch production --message "Fix job search filters"
```

This updates the app for all users within minutes, no store review needed.
Only JS/asset changes can be pushed OTA -- native code changes require a new
build.

## Summary

What we built:
- A Go API with 4 resources (Job, Company, Application, Category)
- An Expo React Native app with tab navigation
- JWT authentication stored in SecureStore
- Infinite scrolling job listings with pull-to-refresh
- Resume upload via document picker and presigned URLs
- Push notifications for application status updates
- Offline support via TanStack Query cache
- Production deployment via EAS Build

The API is identical to every other Grit architecture. The mobile app is a
native client that talks to it over HTTP, with mobile-specific patterns for
auth, navigation, and file handling.
