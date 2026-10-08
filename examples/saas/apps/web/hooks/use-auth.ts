import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import type {
  User,
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  ApiResponse,
} from "@repo/shared/types";
import { api } from "@/lib/api";
import { setWebSessionMarker, clearWebSessionMarker } from "@/lib/web-session";

// Token storage policy (Grit 3.26+):
//   - The API issues HttpOnly cookies (grit_access + grit_refresh) on
//     login/register/refresh and clears them on logout. The browser
//     handles them automatically — JS never reads or writes the access
//     token, so XSS can't exfiltrate it.
//   - The axios client uses withCredentials: true so the browser
//     attaches the cookies on every request.
//   - There is no Bearer-header path here. Mobile/desktop bearer
//     clients live in their own apps; the web app is cookie-only.
//
// Types live in packages/shared/types/user.ts — they're consumed by
// web, admin, and the Go API (via grit sync). Never inline them.

export function useMe() {
  return useQuery<User | null>({
    queryKey: ["me"],
    queryFn: async () => {
      try {
        const { data } = await api.get<ApiResponse<User>>("/api/auth/me");
        return data.data;
      } catch (err: unknown) {
        // 401 is the canonical "not logged in" — return null instead
        // of throwing so guarded pages can read user === null cleanly.
        const e = err as { response?: { status?: number } };
        if (e.response?.status === 401) return null;
        throw err;
      }
    },
    retry: false,
    staleTime: 10 * 60 * 1000,
  });
}

export function useAuth() {
  const { data: user, isLoading, isError } = useMe();
  return {
    user: user ?? null,
    isAuthenticated: !!user,
    isLoading,
    isError,
  };
}

export function useLogin() {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (credentials: LoginRequest) => {
      // POST is enough — the API sets HttpOnly cookies on the response.
      // The 'tokens' field is still in the JSON body so native bearer
      // clients work too, but the browser ignores it.
      const { data } = await api.post<ApiResponse<AuthResponse>>(
        "/api/auth/login",
        credentials
      );
      return data;
    },
    onSuccess: (data) => {
      queryClient.setQueryData(["me"], data.data.user);
      // v3.31.42: stamp the web-origin marker so the middleware
      // doesn't bounce signed-in web users to /login on the next
      // navigation. See lib/web-session.ts for the rationale.
      setWebSessionMarker();
      router.push("/");
    },
  });
}

export function useRegister() {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: RegisterRequest) => {
      const { data } = await api.post<ApiResponse<AuthResponse>>(
        "/api/auth/register",
        payload
      );
      return data;
    },
    onSuccess: (data) => {
      queryClient.setQueryData(["me"], data.data.user);
      // v3.31.42: same web-origin marker as useLogin.
      setWebSessionMarker();
      router.push("/");
    },
  });
}

export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      try {
        await api.post("/api/auth/logout");
      } catch {
        // The cookies are already cleared by the API's Set-Cookie even
        // on a 4xx; just make sure local state is wiped either way.
      }
    },
    onSettled: () => {
      queryClient.clear();
      // v3.31.42: clear the marker so middleware no longer sees
      // a session on the next request.
      clearWebSessionMarker();
      if (typeof window !== "undefined") {
        window.location.href = "/login";
      }
    },
  });
}
