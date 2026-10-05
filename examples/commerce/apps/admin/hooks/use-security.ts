import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

/**
 * What the server says about this account and this deployment.
 *
 * Both halves matter. The account half is state; the deployment half is
 * capability, and it is why sms_provider_configured is here: phone recovery
 * only exists if somebody wired a provider, and the page leaves the card out
 * rather than rendering a control that cannot work.
 */
export interface SecurityOverview {
  email: string;
  email_verified: boolean;
  has_password: boolean;
  provider: string;
  /** Masked, e.g. "b****p@example.com". The server never returns it in full. */
  recovery_email: string;
  recovery_email_verified: boolean;
  recovery_phone: string;
  recovery_phone_verified: boolean;
  sms_provider_configured: boolean;
}

export function useSecurityOverview() {
  return useQuery<SecurityOverview>({
    queryKey: ["security-overview"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/auth/security");
      return data.data;
    },
  });
}

type Kind = "email" | "phone";

/** Starts adding a recovery contact. Sends a code; stores nothing yet. */
export function useSetRecoveryContact(kind: Kind) {
  return useMutation({
    mutationFn: async (input: { password: string; value: string }) => {
      const body =
        kind === "email"
          ? { password: input.password, email: input.value }
          : { password: input.password, phone: input.value };
      const { data } = await apiClient.post("/api/auth/recovery/" + kind, body);
      return data;
    },
  });
}

export function useVerifyRecoveryContact(kind: Kind) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (code: string) => {
      const { data } = await apiClient.post("/api/auth/recovery/" + kind + "/verify", { code });
      return data;
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["security-overview"] }),
  });
}

export function useClearRecoveryContact(kind: Kind) {
  const qc = useQueryClient();
  return useMutation({
    mutationFn: async (password: string) => {
      // A body on DELETE, because removing a recovery contact is as sensitive
      // as adding one and takes the same password.
      const { data } = await apiClient.delete("/api/auth/recovery/" + kind, { data: { password } });
      return data;
    },
    onSuccess: () => qc.invalidateQueries({ queryKey: ["security-overview"] }),
  });
}
