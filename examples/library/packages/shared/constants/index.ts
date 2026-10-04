export const ROLES = {
  ADMIN: "ADMIN",
  EDITOR: "EDITOR",
  USER: "USER",
  // grit:role-constants
} as const;

export const API_ROUTES = {
  AUTH: {
    LOGIN: "/api/auth/login",
    REGISTER: "/api/auth/register",
    REFRESH: "/api/auth/refresh",
    LOGOUT: "/api/auth/logout",
    ME: "/api/auth/me",
    FORGOT_PASSWORD: "/api/auth/forgot-password",
    RESET_PASSWORD: "/api/auth/reset-password",
    OAUTH: {
      GOOGLE: "/api/auth/oauth/google",
      GITHUB: "/api/auth/oauth/github",
    },
  },
  USERS: {
    LIST: "/api/users",
    GET: (id: string) => `/api/users/${id}`,
    UPDATE: (id: string) => `/api/users/${id}`,
    DELETE: (id: string) => `/api/users/${id}`,
  },
  UPLOADS: {
    CREATE: "/api/uploads",
    LIST: "/api/uploads",
    GET: (id: string) => `/api/uploads/${id}`,
    DELETE: (id: string) => `/api/uploads/${id}`,
  },
  AI: {
    COMPLETE: "/api/ai/complete",
    CHAT: "/api/ai/chat",
    STREAM: "/api/ai/stream",
  },
  ADMIN: {
    JOBS_STATS: "/api/admin/jobs/stats",
    JOBS_LIST: (status: string) => `/api/admin/jobs/${status}`,
    JOBS_RETRY: (id: string) => `/api/admin/jobs/${id}/retry`,
    JOBS_CLEAR: (queue: string) => `/api/admin/jobs/queue/${queue}`,
    CRON_TASKS: "/api/admin/cron/tasks",
  },
  PROFILE: {
    GET: "/api/profile",
    UPDATE: "/api/profile",
    DELETE: "/api/profile",
  },
  BLOGS: {
    LIST: "/api/blogs",
    GET: (slug: string) => `/api/blogs/${slug}`,
    ADMIN_LIST: "/api/admin/blogs",
    CREATE: "/api/admin/blogs",
    UPDATE: (id: string) => `/api/admin/blogs/${id}`,
    DELETE: (id: string) => `/api/admin/blogs/${id}`,
  },
  HEALTH: "/api/health",
  AUTHORS: {
    LIST: "/api/authors",
    GET: (id: number) => `/api/authors/${id}`,
    CREATE: "/api/authors",
    UPDATE: (id: number) => `/api/authors/${id}`,
    DELETE: (id: number) => `/api/authors/${id}`,
  },
  BOOKS: {
    LIST: "/api/books",
    GET: (id: number) => `/api/books/${id}`,
    CREATE: "/api/books",
    UPDATE: (id: number) => `/api/books/${id}`,
    DELETE: (id: number) => `/api/books/${id}`,
  },
  LOANS: {
    LIST: "/api/loans",
    GET: (id: number) => `/api/loans/${id}`,
    CREATE: "/api/loans",
    UPDATE: (id: number) => `/api/loans/${id}`,
    DELETE: (id: number) => `/api/loans/${id}`,
  },
  // grit:api-routes
} as const;
