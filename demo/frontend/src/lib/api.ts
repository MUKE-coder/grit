import axios from 'axios'
import { useAuthStore } from '@/stores/auth.store'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor: attach headers
api.interceptors.request.use((config) => {
  const { tokens, currentBusinessID } = useAuthStore.getState()

  // Always attach internal API key
  const apiKey = import.meta.env.VITE_INTERNAL_API_KEY
  if (apiKey) {
    config.headers['X-Internal-Key'] = apiKey
  }

  // Attach JWT if authenticated
  if (tokens?.access_token) {
    config.headers['Authorization'] = `Bearer ${tokens.access_token}`
  }

  // Attach business context
  if (currentBusinessID) {
    config.headers['X-Business-ID'] = String(currentBusinessID)
  }

  return config
})

// Response interceptor: handle 401 (but NOT on auth routes)
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      const url = error.config?.url || ''
      // Don't redirect for login/register/invite — let the component handle the error
      const isAuthRoute = url.includes('/auth/login') || url.includes('/auth/register') || url.includes('/invite/')
      if (!isAuthRoute) {
        const { logout } = useAuthStore.getState()
        logout()
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

export default api
