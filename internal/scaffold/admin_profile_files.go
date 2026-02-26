package scaffold

// adminProfilePage returns the profile page component using React Hook Form + Zod.
func adminProfilePage() string {
	return `"use client";

import { useState, useEffect, useRef } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useMe } from "@/hooks/use-auth";
import { useUpdateProfile, useChangePassword } from "@/hooks/use-profile";
import { User, Briefcase, Lock, Trash2, Save, Loader2, Upload } from "@/lib/icons";
import { DeleteAccountDialog } from "@/components/profile/delete-account-dialog";
import { uploadFile } from "@/lib/api-client";

const PersonalInfoSchema = z.object({
  first_name: z.string().min(2, "First name must be at least 2 characters"),
  last_name: z.string().min(2, "Last name must be at least 2 characters"),
  email: z.string().email("Please enter a valid email"),
});
type PersonalInfoValues = z.infer<typeof PersonalInfoSchema>;

const ProfessionalInfoSchema = z.object({
  job_title: z.string().optional().default(""),
  bio: z.string().optional().default(""),
});
type ProfessionalInfoValues = z.infer<typeof ProfessionalInfoSchema>;

const ChangePasswordSchema = z
  .object({
    password: z.string().min(8, "Password must be at least 8 characters"),
    confirm_password: z.string().min(1, "Please confirm your password"),
  })
  .refine((data) => data.password === data.confirm_password, {
    message: "Passwords do not match",
    path: ["confirm_password"],
  });
type ChangePasswordValues = z.infer<typeof ChangePasswordSchema>;

const inputClass =
  "w-full rounded-lg border border-border bg-bg-tertiary px-3 py-2 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent";
const errorInputClass =
  "w-full rounded-lg border border-danger/50 bg-bg-tertiary px-3 py-2 text-sm text-foreground placeholder:text-text-muted focus:border-danger focus:outline-none focus:ring-1 focus:ring-danger";

export default function ProfilePage() {
  const { data: user } = useMe();
  const updateProfile = useUpdateProfile();
  const changePassword = useChangePassword();
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [avatarUploading, setAvatarUploading] = useState(false);
  const avatarInputRef = useRef<HTMLInputElement>(null);

  const handleAvatarUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setAvatarUploading(true);
    try {
      const result = await uploadFile(file);
      const url = (result.data as Record<string, unknown>)?.url as string;
      if (url) {
        updateProfile.mutate({ avatar: url });
      }
    } catch {
      // Upload failed silently
    } finally {
      setAvatarUploading(false);
      if (avatarInputRef.current) avatarInputRef.current.value = "";
    }
  };

  // Personal info form
  const personalForm = useForm<PersonalInfoValues>({
    resolver: zodResolver(PersonalInfoSchema),
    defaultValues: {
      first_name: user?.first_name || "",
      last_name: user?.last_name || "",
      email: user?.email || "",
    },
  });

  // Professional info form
  const professionalForm = useForm<ProfessionalInfoValues>({
    resolver: zodResolver(ProfessionalInfoSchema),
    defaultValues: {
      job_title: user?.job_title || "",
      bio: user?.bio || "",
    },
  });

  // Password form
  const passwordForm = useForm<ChangePasswordValues>({
    resolver: zodResolver(ChangePasswordSchema),
    defaultValues: { password: "", confirm_password: "" },
  });

  // Reset form defaults when user data loads
  useEffect(() => {
    if (user) {
      personalForm.reset({
        first_name: user.first_name,
        last_name: user.last_name,
        email: user.email,
      });
      professionalForm.reset({
        job_title: user.job_title || "",
        bio: user.bio || "",
      });
    }
  }, [user]);

  const onPersonalSubmit = (data: PersonalInfoValues) => {
    updateProfile.mutate(data);
  };

  const onProfessionalSubmit = (data: ProfessionalInfoValues) => {
    updateProfile.mutate(data);
  };

  const onPasswordSubmit = (data: ChangePasswordValues) => {
    changePassword.mutate(
      { password: data.password },
      { onSuccess: () => passwordForm.reset() }
    );
  };

  if (!user) return null;

  return (
    <div className="space-y-6">
      {/* Page header */}
      <div>
        <h1 className="text-2xl font-bold text-foreground">Profile</h1>
        <p className="mt-1 text-sm text-text-secondary">
          Manage your personal information, professional details, and account settings.
        </p>
      </div>

      {/* Profile overview card */}
      <div className="rounded-xl border border-border bg-bg-secondary p-6">
        <div className="flex items-center gap-4">
          <div className="relative group">
            <div className="flex h-16 w-16 items-center justify-center rounded-full bg-accent/15 border border-accent/20 overflow-hidden">
              {user.avatar ? (
                <img src={user.avatar} alt="" className="h-16 w-16 rounded-full object-cover" />
              ) : (
                <span className="text-2xl font-bold text-accent">
                  {user.first_name?.[0]}{user.last_name?.[0]}
                </span>
              )}
            </div>
            <button
              type="button"
              onClick={() => avatarInputRef.current?.click()}
              disabled={avatarUploading}
              className="absolute inset-0 flex items-center justify-center rounded-full bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity cursor-pointer"
            >
              {avatarUploading ? (
                <Loader2 className="h-5 w-5 text-white animate-spin" />
              ) : (
                <Upload className="h-5 w-5 text-white" />
              )}
            </button>
            <input
              ref={avatarInputRef}
              type="file"
              accept="image/*"
              onChange={handleAvatarUpload}
              className="hidden"
            />
          </div>
          <div>
            <h2 className="text-lg font-semibold text-foreground">
              {user.first_name} {user.last_name}
            </h2>
            <p className="text-sm text-text-secondary">{user.email}</p>
            <div className="mt-1 flex items-center gap-2">
              <span className="inline-flex items-center rounded-full bg-accent/10 px-2.5 py-0.5 text-xs font-medium text-accent">
                {user.role}
              </span>
              {user.job_title && (
                <span className="text-xs text-text-muted">{user.job_title}</span>
              )}
            </div>
          </div>
        </div>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        {/* Personal Information */}
        <div className="rounded-xl border border-border bg-bg-secondary">
          <div className="border-b border-border px-6 py-4">
            <div className="flex items-center gap-2">
              <User className="h-4 w-4 text-accent" />
              <h3 className="font-semibold text-foreground">Personal Information</h3>
            </div>
            <p className="mt-1 text-xs text-text-muted">Update your name and email address.</p>
          </div>
          <form onSubmit={personalForm.handleSubmit(onPersonalSubmit)} className="space-y-4 p-6">
            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-1.5">
                <label className="text-sm font-medium text-text-secondary">First name</label>
                <input
                  {...personalForm.register("first_name")}
                  className={personalForm.formState.errors.first_name ? errorInputClass : inputClass}
                />
                {personalForm.formState.errors.first_name && (
                  <p className="text-xs text-danger">{personalForm.formState.errors.first_name.message}</p>
                )}
              </div>
              <div className="space-y-1.5">
                <label className="text-sm font-medium text-text-secondary">Last name</label>
                <input
                  {...personalForm.register("last_name")}
                  className={personalForm.formState.errors.last_name ? errorInputClass : inputClass}
                />
                {personalForm.formState.errors.last_name && (
                  <p className="text-xs text-danger">{personalForm.formState.errors.last_name.message}</p>
                )}
              </div>
            </div>
            <div className="space-y-1.5">
              <label className="text-sm font-medium text-text-secondary">Email</label>
              <input
                type="email"
                {...personalForm.register("email")}
                className={personalForm.formState.errors.email ? errorInputClass : inputClass}
              />
              {personalForm.formState.errors.email && (
                <p className="text-xs text-danger">{personalForm.formState.errors.email.message}</p>
              )}
            </div>
            <div className="flex items-center justify-between">
              {updateProfile.isSuccess && (
                <p className="text-sm text-success">Profile updated successfully.</p>
              )}
              <div className="ml-auto">
                <button
                  type="submit"
                  disabled={updateProfile.isPending}
                  className="flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
                >
                  {updateProfile.isPending ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <Save className="h-4 w-4" />
                  )}
                  Save changes
                </button>
              </div>
            </div>
          </form>
        </div>

        {/* Professional Information */}
        <div className="rounded-xl border border-border bg-bg-secondary">
          <div className="border-b border-border px-6 py-4">
            <div className="flex items-center gap-2">
              <Briefcase className="h-4 w-4 text-accent" />
              <h3 className="font-semibold text-foreground">Professional Information</h3>
            </div>
            <p className="mt-1 text-xs text-text-muted">Add your job title and a short bio.</p>
          </div>
          <form onSubmit={professionalForm.handleSubmit(onProfessionalSubmit)} className="space-y-4 p-6">
            <div className="space-y-1.5">
              <label className="text-sm font-medium text-text-secondary">Job title</label>
              <input
                {...professionalForm.register("job_title")}
                placeholder="e.g. Software Engineer"
                className={inputClass}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-sm font-medium text-text-secondary">Bio</label>
              <textarea
                {...professionalForm.register("bio")}
                rows={4}
                placeholder="Write a short bio about yourself..."
                className={inputClass + " resize-none"}
              />
            </div>
            <div className="flex justify-end">
              <button
                type="submit"
                disabled={updateProfile.isPending}
                className="flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
              >
                {updateProfile.isPending ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <Save className="h-4 w-4" />
                )}
                Save changes
              </button>
            </div>
          </form>
        </div>

        {/* Change Password */}
        <div className="rounded-xl border border-border bg-bg-secondary">
          <div className="border-b border-border px-6 py-4">
            <div className="flex items-center gap-2">
              <Lock className="h-4 w-4 text-accent" />
              <h3 className="font-semibold text-foreground">Change Password</h3>
            </div>
            <p className="mt-1 text-xs text-text-muted">Update your password to keep your account secure.</p>
          </div>
          <form onSubmit={passwordForm.handleSubmit(onPasswordSubmit)} className="space-y-4 p-6">
            <div className="space-y-1.5">
              <label className="text-sm font-medium text-text-secondary">New password</label>
              <input
                type="password"
                {...passwordForm.register("password")}
                placeholder="Min. 8 characters"
                className={passwordForm.formState.errors.password ? errorInputClass : inputClass}
              />
              {passwordForm.formState.errors.password && (
                <p className="text-xs text-danger">{passwordForm.formState.errors.password.message}</p>
              )}
            </div>
            <div className="space-y-1.5">
              <label className="text-sm font-medium text-text-secondary">Confirm new password</label>
              <input
                type="password"
                {...passwordForm.register("confirm_password")}
                placeholder="Repeat your new password"
                className={passwordForm.formState.errors.confirm_password ? errorInputClass : inputClass}
              />
              {passwordForm.formState.errors.confirm_password && (
                <p className="text-xs text-danger">{passwordForm.formState.errors.confirm_password.message}</p>
              )}
            </div>
            <div className="flex items-center justify-between">
              {changePassword.isSuccess && (
                <p className="text-sm text-success">Password updated successfully.</p>
              )}
              <div className="ml-auto">
                <button
                  type="submit"
                  disabled={changePassword.isPending}
                  className="flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white hover:bg-accent-hover disabled:opacity-50 transition-colors"
                >
                  {changePassword.isPending ? (
                    <Loader2 className="h-4 w-4 animate-spin" />
                  ) : (
                    <Lock className="h-4 w-4" />
                  )}
                  Update password
                </button>
              </div>
            </div>
          </form>
        </div>

        {/* Danger Zone */}
        <div className="rounded-xl border border-danger/30 bg-bg-secondary">
          <div className="border-b border-danger/20 px-6 py-4">
            <div className="flex items-center gap-2">
              <Trash2 className="h-4 w-4 text-danger" />
              <h3 className="font-semibold text-danger">Danger Zone</h3>
            </div>
            <p className="mt-1 text-xs text-text-muted">Irreversible and destructive actions.</p>
          </div>
          <div className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <h4 className="text-sm font-medium text-foreground">Delete account</h4>
                <p className="mt-0.5 text-xs text-text-muted">
                  Permanently delete your account and all associated data. This action cannot be undone.
                </p>
              </div>
              <button
                onClick={() => setShowDeleteDialog(true)}
                className="rounded-lg border border-danger/30 bg-danger/10 px-4 py-2 text-sm font-medium text-danger hover:bg-danger/20 transition-colors"
              >
                Delete account
              </button>
            </div>
          </div>
        </div>
      </div>

      <DeleteAccountDialog
        open={showDeleteDialog}
        onClose={() => setShowDeleteDialog(false)}
      />
    </div>
  );
}
`
}

// adminUseProfile returns the profile React Query hooks.
func adminUseProfile() string {
	return `import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";
import { useRouter } from "next/navigation";
import Cookies from "js-cookie";

interface UpdateProfileData {
  first_name?: string;
  last_name?: string;
  email?: string;
  job_title?: string;
  bio?: string;
  avatar?: string;
  password?: string;
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
    mutationFn: async (data: { password: string }) => {
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
      await apiClient.delete("/api/profile");
    },
    onSuccess: () => {
      Cookies.remove("access_token");
      Cookies.remove("refresh_token");
      queryClient.clear();
      router.push("/login");
    },
  });
}
`
}

// adminDeleteAccountDialog returns the delete account confirmation dialog.
func adminDeleteAccountDialog() string {
	return `"use client";

import { useState } from "react";
import { useDeleteAccount } from "@/hooks/use-profile";
import { AlertTriangle, Loader2, X } from "@/lib/icons";

interface DeleteAccountDialogProps {
  open: boolean;
  onClose: () => void;
}

export function DeleteAccountDialog({ open, onClose }: DeleteAccountDialogProps) {
  const [confirmation, setConfirmation] = useState("");
  const deleteAccount = useDeleteAccount();

  if (!open) return null;

  const handleDelete = () => {
    if (confirmation !== "DELETE") return;
    deleteAccount.mutate();
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div className="fixed inset-0 bg-black/60 backdrop-blur-sm" onClick={onClose} />
      <div className="relative w-full max-w-md rounded-xl border border-border bg-bg-secondary p-6 shadow-2xl">
        <button
          onClick={onClose}
          className="absolute right-4 top-4 text-text-muted hover:text-foreground transition-colors"
        >
          <X className="h-4 w-4" />
        </button>

        <div className="flex items-center gap-3 text-danger">
          <div className="flex h-10 w-10 items-center justify-center rounded-full bg-danger/10">
            <AlertTriangle className="h-5 w-5" />
          </div>
          <div>
            <h3 className="font-semibold text-foreground">Delete Account</h3>
            <p className="text-xs text-text-muted">This action is permanent</p>
          </div>
        </div>

        <div className="mt-4 space-y-3">
          <p className="text-sm text-text-secondary">
            This will permanently delete your account, including all your data, settings, and
            access. This action <strong className="text-foreground">cannot be undone</strong>.
          </p>

          <div className="space-y-1.5">
            <label className="text-sm font-medium text-text-secondary">
              Type <span className="font-mono text-danger">DELETE</span> to confirm
            </label>
            <input
              type="text"
              value={confirmation}
              onChange={(e) => setConfirmation(e.target.value)}
              placeholder="DELETE"
              className="w-full rounded-lg border border-border bg-bg-tertiary px-3 py-2 text-sm text-foreground placeholder:text-text-muted focus:border-danger focus:outline-none focus:ring-1 focus:ring-danger"
            />
          </div>
        </div>

        <div className="mt-6 flex justify-end gap-3">
          <button
            onClick={onClose}
            className="rounded-lg border border-border px-4 py-2 text-sm font-medium text-foreground hover:bg-bg-hover transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={handleDelete}
            disabled={confirmation !== "DELETE" || deleteAccount.isPending}
            className="flex items-center gap-2 rounded-lg bg-danger px-4 py-2 text-sm font-medium text-white hover:bg-danger/90 disabled:opacity-50 transition-colors"
          >
            {deleteAccount.isPending ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <AlertTriangle className="h-4 w-4" />
            )}
            Delete my account
          </button>
        </div>
      </div>
    </div>
  );
}
`
}
