import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { useRouter } from "next/navigation";

interface UpdateProfileData {
  first_name?: string;
  last_name?: string;
  email?: string;
  job_title?: string;
  bio?: string;
  avatar?: string;
  password?: string;
  // Needed to change the email or the password.
  current_password?: string;
}

export function useUpdateProfile() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (data: UpdateProfileData) => {
      const { data: response } = await apiClient.put("/api/profile", data);
      return response;
    },
    onSuccess: (data) => {
      queryClient.setQueryData(["me"], data.data);
    },
  });
}

export function useChangePassword() {
  return useMutation({
    mutationFn: async (data: { password: string; current_password: string }) => {
      const { data: response } = await apiClient.put("/api/profile", data);
      return response;
    },
  });
}

export function useDeleteAccount() {
  const router = useRouter();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      // The API's DELETE handler clears the auth cookies via Set-Cookie
      // max-age=0 as part of the same response. JS doesn't need to touch
      // anything cookie-side.
      await apiClient.delete("/api/profile");
    },
    onSuccess: () => {
      queryClient.clear();
      router.push("/login");
    },
  });
}
