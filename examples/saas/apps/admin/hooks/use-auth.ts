import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import type {
  User,
  LoginRequest,
  RegisterRequest,
  AuthResponse,
  ApiResponse,
} from "@repo/shared/types";
import { apiClient } from "@/lib/api-client";
import {
  encodeAssertion,
  toRequestOptions,
  type RequestOptionsJSON,
} from "@/lib/webauthn";

// Auth flow (Grit 3.27+):
//   - User shape lives in packages/shared/types/user.ts (the SAME type
//     the Go API serialises and the web app imports). Never inline it.
//   - Login / register / OAuth set HttpOnly grit_access + grit_refresh
//     cookies via Set-Cookie. The browser stores them; JS never reads.
//   - apiClient has withCredentials: true so the cookies are attached
//     automatically on every request.
//   - useMe returns null on 401 instead of throwing so guard components
//     can read user === null cleanly without a try/catch.
//   - useLogout posts to /api/auth/logout — the API clears the cookies
//     via Set-Cookie max-age=0; we just clear the query cache + redirect.

export function useMe() {
  return useQuery<User | null>({
    queryKey: ["me"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<ApiResponse<User>>("/api/auth/me");
        return data.data;
      } catch (err: unknown) {
        const e = err as { response?: { status?: number } };
        if (e.response?.status === 401) return null;
        throw err;
      }
    },
    retry: false,
    staleTime: 10 * 60 * 1000,
  });
}

// A login can end in two places: signed in, or holding a short-lived pending
// token because the account has 2FA on and this device is not trusted.
export interface TOTPChallenge {
  totp_required: true;
  pending_token: string;
}

function isChallenge(d: AuthResponse | TOTPChallenge): d is TOTPChallenge {
  return (d as TOTPChallenge).totp_required === true;
}

export function useLogin() {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (credentials: LoginRequest) => {
      // POST is enough — the API responds with Set-Cookie headers for
      // grit_access + grit_refresh. The 'tokens' field in the JSON body
      // is preserved for native bearer clients (mobile/desktop); the
      // browser ignores it.
      const { data } = await apiClient.post<ApiResponse<AuthResponse | TOTPChallenge>>(
        "/api/auth/login",
        credentials
      );
      return data;
    },
    onSuccess: (data) => {
      // Reading data.data.user here without the guard is what used to break:
      // on a 2FA account it is undefined, and the redirect threw instead of
      // showing the code prompt.
      if (isChallenge(data.data)) return;
      queryClient.setQueryData(["me"], data.data.user);
      router.push(data.data.user.role === "USER" ? "/account" : "/dashboard");
    },
  });
}

// Signs in with a passkey, no email and no password.
//
// BeginLogin is called without a user, so the challenge carries no
// allowCredentials and the browser offers whichever passkeys it holds for this
// site: the discoverable-credential flow, which is the one that lets somebody
// tap a fingerprint reader and be in.
//
// The server answers with the same token pair a password sign-in issues and
// records the same session row, so everything downstream, including Active
// Sessions, is unchanged.
export function usePasskeyLogin() {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      const { data: begun } = await apiClient.post<
        ApiResponse<{ session_id: string; options: RequestOptionsJSON }>
      >("/api/auth/passkeys/login/begin", {});

      const assertion = (await navigator.credentials.get({
        publicKey: toRequestOptions(begun.data.options),
      })) as PublicKeyCredential | null;
      if (!assertion) throw new Error("No passkey was chosen.");

      // The session id goes in the query string because the body is the raw
      // assertion: the server hands it to the WebAuthn library unparsed, and
      // wrapping it in an envelope would mean unwrapping it there.
      const { data } = await apiClient.post<ApiResponse<AuthResponse>>(
        "/api/auth/passkeys/login/finish?session=" + encodeURIComponent(begun.data.session_id),
        encodeAssertion(assertion),
      );
      return data;
    },
    onSuccess: (data) => {
      queryClient.setQueryData(["me"], data.data.user);
      router.push(data.data.user.role === "USER" ? "/account" : "/dashboard");
    },
  });
}

// Completes a login that stopped at the 2FA prompt. Same shape for an
// authenticator code and a backup code — only the endpoint differs — so the
// form can offer both without branching.
export function useVerifyTOTP() {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: {
      pending_token: string;
      code: string;
      trust_device?: boolean;
      backup?: boolean;
    }) => {
      const path = payload.backup
        ? "/api/auth/totp/backup-codes/verify"
        : "/api/auth/totp/verify";
      const { data } = await apiClient.post<ApiResponse<AuthResponse>>(path, {
        pending_token: payload.pending_token,
        code: payload.code.trim(),
        trust_device: payload.trust_device ?? false,
      });
      return data;
    },
    onSuccess: (data) => {
      queryClient.setQueryData(["me"], data.data.user);
      router.push(data.data.user.role === "USER" ? "/account" : "/dashboard");
    },
  });
}

/* ── Two-factor management ─────────────────────────────────────────── */

export interface TOTPStatus {
  enabled: boolean;
  backup_codes_remaining: number;
  trusted_devices: number;
}

export interface TrustedDevice {
  id: string;
  user_agent: string;
  ip_address: string;
  created_at: string;
  expires_at: string;
  current: boolean;
}

export function useTOTPStatus() {
  return useQuery({
    queryKey: ["totp-status"],
    queryFn: async () => {
      const { data } = await apiClient.get<ApiResponse<TOTPStatus>>("/api/auth/totp/status");
      return data.data;
    },
  });
}

export function useTrustedDevices(enabled: boolean) {
  return useQuery({
    queryKey: ["trusted-devices"],
    // Pointless call when 2FA is off — there can be no trusted devices.
    enabled,
    queryFn: async () => {
      const { data } = await apiClient.get<ApiResponse<TrustedDevice[]>>(
        "/api/auth/totp/trusted-devices"
      );
      return data.data;
    },
  });
}

export function useTOTPSetup() {
  return useMutation({
    mutationFn: async () => {
      const { data } = await apiClient.post<
        ApiResponse<{ secret: string; uri: string; qr_code: string }>
      >("/api/auth/totp/setup", {});
      return data.data;
    },
  });
}

export function useEnableTOTP() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (payload: { secret: string; code: string }) => {
      const { data } = await apiClient.post<ApiResponse<{ backup_codes: string[] }>>(
        "/api/auth/totp/enable",
        { secret: payload.secret, code: payload.code.trim() }
      );
      return data.data;
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["totp-status"] }),
  });
}

export function useDisableTOTP() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (password: string) => {
      const { data } = await apiClient.post("/api/auth/totp/disable", { password });
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["totp-status"] });
      queryClient.invalidateQueries({ queryKey: ["trusted-devices"] });
    },
  });
}

export function useRegenerateBackupCodes() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async () => {
      const { data } = await apiClient.post<ApiResponse<{ backup_codes: string[] }>>(
        "/api/auth/totp/backup-codes",
        {}
      );
      return data.data;
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["totp-status"] }),
  });
}

export function useRevokeTrustedDevice() {
  const queryClient = useQueryClient();
  return useMutation({
    // No id revokes every device; an id revokes just that one.
    mutationFn: async (id?: string) => {
      const path = id
        ? "/api/auth/totp/trusted-devices/" + id
        : "/api/auth/totp/trusted-devices";
      const { data } = await apiClient.delete(path);
      return data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["trusted-devices"] });
      queryClient.invalidateQueries({ queryKey: ["totp-status"] });
    },
  });
}

export function useRegister() {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (payload: RegisterRequest) => {
      const { data } = await apiClient.post<ApiResponse<AuthResponse>>(
        "/api/auth/register",
        payload
      );
      return data;
    },
    onSuccess: (data) => {
      queryClient.setQueryData(["me"], data.data.user);
      router.push(data.data.user.role === "USER" ? "/account" : "/dashboard");
    },
  });
}

export function useLogout() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      try {
        await apiClient.post("/api/auth/logout");
      } catch {
        // The API clears the cookies via Set-Cookie max-age=0 even on
        // a 4xx, so falling through here is safe — local state still
        // gets wiped by onSettled.
      }
    },
    onSettled: () => {
      queryClient.clear();
      if (typeof window !== "undefined") {
        window.location.href = "/login";
      }
    },
  });
}
