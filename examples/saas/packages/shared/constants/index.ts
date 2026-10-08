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
  PLANS: {
    LIST: "/api/plans",
    GET: (id: number) => `/api/plans/${id}`,
    CREATE: "/api/plans",
    UPDATE: (id: number) => `/api/plans/${id}`,
    DELETE: (id: number) => `/api/plans/${id}`,
  },
  SUBSCRIPTIONS: {
    LIST: "/api/subscriptions",
    GET: (id: number) => `/api/subscriptions/${id}`,
    CREATE: "/api/subscriptions",
    UPDATE: (id: number) => `/api/subscriptions/${id}`,
    DELETE: (id: number) => `/api/subscriptions/${id}`,
  },
  INVOICES: {
    LIST: "/api/invoices",
    GET: (id: number) => `/api/invoices/${id}`,
    CREATE: "/api/invoices",
    UPDATE: (id: number) => `/api/invoices/${id}`,
    DELETE: (id: number) => `/api/invoices/${id}`,
  },
  USAGE_RECORDS: {
    LIST: "/api/usage_records",
    GET: (id: number) => `/api/usage_records/${id}`,
    CREATE: "/api/usage_records",
    UPDATE: (id: number) => `/api/usage_records/${id}`,
    DELETE: (id: number) => `/api/usage_records/${id}`,
  },
  // grit:api-routes
} as const;
