# Job Portal -- Mobile Architecture (API + Expo)

A mobile-first job portal built with Grit using the **Mobile** architecture:
- `apps/api` -- Go API (Gin + GORM + PostgreSQL)
- `apps/expo` -- Expo React Native mobile app (iOS + Android)
- `packages/shared` -- Shared TypeScript types and validation schemas

## Architecture

This configuration pairs a Go API with an Expo React Native mobile app. Users
browse jobs, apply, and track applications from their phone. The monorepo
structure keeps API types and mobile code in sync via `packages/shared`.

Best for: mobile-first apps, on-the-go workflows, apps where push notifications
matter, teams building for iOS and Android simultaneously.

## Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+
- pnpm 8+
- Docker (for PostgreSQL, Redis, MinIO, Mailhog)
- Expo CLI (`npx expo --version` to verify)
- Expo Go app on your phone ([iOS](https://apps.apple.com/app/expo-go/id982107779) / [Android](https://play.google.com/store/apps/details?id=host.exp.exponent))

### Setup

1. Install Grit CLI:
   ```bash
   go install github.com/MUKE-coder/grit/v3/cmd/grit@latest
   ```

2. Create the project:
   ```bash
   grit new job-portal --mobile
   ```

3. Start Docker services:
   ```bash
   cd job-portal
   docker compose up -d
   ```

4. Generate resources:
   ```bash
   grit generate resource Category --fields "name:string:unique,slug:slug:name"
   grit generate resource Company --fields "name:string,slug:slug:name,description:richtext,logo:string,website:string:optional,email:string"
   grit generate resource Job --fields "title:string,description:richtext,company_id:belongs_to:Company,location:string,salary_min:float,salary_max:float,type:string,experience:string,active:bool,deadline:date"
   grit generate resource Application --fields "job_id:belongs_to:Job,applicant_name:string,applicant_email:string,resume:string,cover_letter:text:optional,status:string"
   ```

5. Run migrations and seed:
   ```bash
   grit migrate
   grit seed
   ```

6. Install dependencies and start:
   ```bash
   pnpm install

   # Terminal 1: API
   cd apps/api && go run cmd/server/main.go

   # Terminal 2: Expo
   cd apps/expo && npx expo start
   ```

7. Open:
   - Scan the QR code in your terminal with Expo Go
   - API Docs: http://localhost:8080/docs
   - GORM Studio: http://localhost:8080/studio
   - Mailhog: http://localhost:8025

## Connecting from a Physical Device

When running on a physical phone, `localhost` points to the phone itself, not
your computer. You need to use your computer's local IP address.

```bash
# Find your local IP
# macOS/Linux:
ifconfig | grep "inet " | grep -v 127.0.0.1
# Windows:
ipconfig | findstr "IPv4"
```

Set `EXPO_PUBLIC_API_URL` in your `.env`:
```
EXPO_PUBLIC_API_URL=http://192.168.1.100:8080
```

Make sure `CORS_ORIGINS` in the API `.env` includes this address or is set to `*`
for development.

## Project Structure

```
job-portal/
├── apps/
│   ├── api/                      # Go backend
│   │   ├── cmd/server/main.go    # Entry point
│   │   └── internal/
│   │       ├── models/           # Job, Company, Application, Category, User
│   │       ├── handlers/         # REST handlers
│   │       ├── services/         # Business logic
│   │       ├── middleware/        # Auth, CORS, cache, rate limiting
│   │       └── routes/           # Route registration
│   └── expo/                     # Expo React Native app
│       ├── app/                  # Expo Router (file-based routing)
│       │   ├── _layout.tsx       # Root layout (providers, fonts)
│       │   ├── (auth)/           # Auth screens (login, register)
│       │   │   ├── login.tsx
│       │   │   └── register.tsx
│       │   └── (tabs)/           # Tab navigation
│       │       ├── _layout.tsx   # Tab bar configuration
│       │       ├── index.tsx     # Home tab (featured jobs)
│       │       ├── jobs.tsx      # Jobs tab (search + list)
│       │       ├── applications.tsx  # My Applications tab
│       │       └── profile.tsx   # Profile tab
│       ├── components/           # Reusable React Native components
│       │   ├── JobCard.tsx       # Job listing card
│       │   ├── ApplicationCard.tsx
│       │   ├── SearchBar.tsx
│       │   └── StatusBadge.tsx
│       ├── hooks/                # Custom hooks
│       │   ├── useAuth.ts        # Auth state + SecureStore
│       │   ├── useJobs.ts        # TanStack Query for jobs
│       │   └── useApplications.ts
│       ├── lib/                  # Utilities
│       │   ├── api.ts            # Axios instance with auth interceptor
│       │   └── secure-store.ts   # SecureStore wrapper
│       └── app.json              # Expo configuration
├── packages/shared/              # Shared types + schemas
│   ├── schemas/                  # Zod validation
│   └── types/                    # TypeScript interfaces
├── .env                          # Environment variables
├── .env.example                  # Template
├── docker-compose.yml            # Dev services
├── docker-compose.prod.yml       # Production (API only)
├── turbo.json                    # Turborepo config
└── pnpm-workspace.yaml           # Monorepo workspaces
```

## Features Demonstrated

### Authentication
- Email/password registration and login
- OAuth2 social login (Google, GitHub) via AuthSession
- JWT access tokens (15 min) + refresh tokens (7 days)
- TOTP two-factor authentication with backup codes
- Tokens stored in SecureStore (encrypted, not AsyncStorage)
- Role-based access: ADMIN (full access), EDITOR (manage jobs), USER (apply)

### Mobile-Specific Features

#### Tab Navigation
Four tabs at the bottom of the screen:
- **Home** -- Featured jobs, search bar, category quick links
- **Jobs** -- Full job listing with search, filters, infinite scroll
- **Applications** -- User's submitted applications with status tracking
- **Profile** -- User info, settings, logout

#### SecureStore for Tokens
Unlike web apps that use `localStorage`, the mobile app stores JWT tokens in
Expo SecureStore -- encrypted storage backed by Keychain (iOS) and EncryptedSharedPreferences (Android).

```typescript
import * as SecureStore from "expo-secure-store";

await SecureStore.setItemAsync("access_token", token);
const token = await SecureStore.getItemAsync("access_token");
```

#### Pull-to-Refresh
All list screens support pull-to-refresh via `RefreshControl`:

```typescript
<FlatList
  data={jobs}
  refreshControl={
    <RefreshControl refreshing={isRefreshing} onRefresh={refetch} />
  }
/>
```

#### Infinite Scrolling
Job listings use `onEndReached` with TanStack Query's `useInfiniteQuery`:

```typescript
const { data, fetchNextPage, hasNextPage } = useInfiniteQuery({
  queryKey: ["jobs"],
  queryFn: ({ pageParam }) => fetchJobs({ page: pageParam }),
  getNextPageParam: (lastPage) => lastPage.meta.page < lastPage.meta.pages
    ? lastPage.meta.page + 1
    : undefined,
});
```

#### Push Notifications
Application status updates trigger push notifications via Expo Notifications:

```typescript
import * as Notifications from "expo-notifications";

// Register for push notifications on login
const { data: token } = await Notifications.getExpoPushTokenAsync();
// Send token to API: POST /api/users/push-token
```

The API sends notifications when application status changes (Reviewed, Interview,
Accepted, Rejected).

#### Camera/Gallery for Resume Upload
Users can take a photo or pick a file from their device to upload as a resume:

```typescript
import * as ImagePicker from "expo-image-picker";
import * as DocumentPicker from "expo-document-picker";

// Pick a document (PDF, DOC)
const result = await DocumentPicker.getDocumentAsync({
  type: ["application/pdf", "application/msword"],
});
```

The file is uploaded via presigned URL, same as the web flow.

#### Offline Support
TanStack Query's built-in cache provides offline read support. Previously loaded
job listings and application statuses are available without a network connection.

### Emails
- Welcome email on registration
- Application confirmation email
- Application status update notifications

### Dashboard Stats
- `GET /api/stats` -- Total jobs, companies, applications, users
- `GET /api/stats/applications` -- Applications by status breakdown

### Deployment
See the Deployment section below.

## API Reference

The API is identical to the API Only architecture. See the full reference at
http://localhost:8080/docs when running locally.

### Authentication

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/auth/register` | Register new user | -- |
| `POST` | `/api/auth/login` | Login (returns JWT) | -- |
| `POST` | `/api/auth/refresh` | Refresh access token | Refresh token |
| `GET` | `/api/auth/me` | Get current user | Bearer |
| `POST` | `/api/auth/totp/setup` | Enable 2FA | Bearer |
| `POST` | `/api/auth/totp/verify` | Verify TOTP code | Bearer |

### Jobs

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/jobs` | List jobs (paginated) | -- |
| `POST` | `/api/jobs` | Create job | ADMIN, EDITOR |
| `GET` | `/api/jobs/:id` | Get job by ID | -- |
| `PUT` | `/api/jobs/:id` | Update job | ADMIN, EDITOR |
| `DELETE` | `/api/jobs/:id` | Delete job | ADMIN |

### Companies

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/companies` | List companies | -- |
| `POST` | `/api/companies` | Create company | ADMIN, EDITOR |
| `GET` | `/api/companies/:id` | Get company by ID | -- |
| `PUT` | `/api/companies/:id` | Update company | ADMIN, EDITOR |

### Applications

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/applications` | Submit application | Bearer |
| `GET` | `/api/applications` | List applications | ADMIN |
| `GET` | `/api/applications/:id` | Get application | ADMIN, Owner |
| `PUT` | `/api/applications/:id` | Update status | ADMIN |

### Categories

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `GET` | `/api/categories` | List categories | -- |
| `POST` | `/api/categories` | Create category | ADMIN |

### Uploads

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/uploads/presign` | Get presigned upload URL | Bearer |

### Push Notifications

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| `POST` | `/api/users/push-token` | Register Expo push token | Bearer |

## Environment Variables

See `.env.example` for the full list.

Key variables:
- `DATABASE_URL` -- PostgreSQL connection string
- `JWT_SECRET` -- Must be changed in production
- `STORAGE_ENDPOINT` -- MinIO (dev) or S3/R2 (prod)
- `RESEND_API_KEY` -- For sending emails
- `EXPO_PUBLIC_API_URL` -- API URL the mobile app connects to
- `GOOGLE_CLIENT_ID` -- For OAuth (optional)
- `GITHUB_CLIENT_ID` -- For OAuth (optional)

## Deployment

### API Deployment

The API deploys exactly like any Grit API:

#### Option 1: grit deploy (VPS)
```bash
grit deploy --host deploy@server.com --domain api.jobs.example.com
```

#### Option 2: Docker Compose
```bash
docker compose -f docker-compose.prod.yml up -d
```

### Mobile App Deployment

The Expo app is built and submitted separately using EAS (Expo Application Services):

```bash
# Install EAS CLI
npm install -g eas-cli

# Configure EAS
cd apps/expo
eas build:configure

# Build for iOS
eas build --platform ios

# Build for Android
eas build --platform android

# Submit to stores
eas submit --platform ios
eas submit --platform android
```

Before building, update `EXPO_PUBLIC_API_URL` in `.env` to your production API URL:
```
EXPO_PUBLIC_API_URL=https://api.jobs.example.com
```

### Over-the-Air Updates

After the initial store release, push updates without going through app review:

```bash
eas update --branch production --message "Fix job search filters"
```

## LLM Instructions

If you are an AI assistant working on this project:
1. Read `.claude/skills/grit/SKILL.md` for Grit conventions
2. Use `grit generate resource` for new resources
3. Never delete `// grit:*` marker comments
4. Follow the API response format: `{ data, message }` / `{ data, meta }` / `{ error: { code, message } }`
5. Use TanStack Query for data fetching, Zod for validation
6. Use Expo Router for navigation (file-based, not React Navigation manually)
7. Use SecureStore for tokens, never AsyncStorage
8. Report bugs at https://github.com/MUKE-coder/grit/issues

## License

MIT
