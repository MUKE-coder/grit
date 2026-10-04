"use client";

import { useState, useEffect, useRef } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  PersonalInfoSchema,
  ProfessionalInfoSchema,
  ChangePasswordSchema,
  type PersonalInfoInput as PersonalInfoValues,
  type ProfessionalInfoInput as ProfessionalInfoValues,
  type ChangePasswordInput as ChangePasswordValues,
} from "@repo/shared/schemas";
import { useMe } from "@/hooks/use-auth";
import { useUpdateProfile, useChangePassword } from "@/hooks/use-profile";
import { PageHeader } from "@/components/chrome/PageHeader";
import { AvatarCropper } from "@/components/account/avatar-cropper";
import { DeleteAccountDialog } from "@/components/profile/delete-account-dialog";
import { ActiveSessions } from "@/components/profile/active-sessions";
import { TwoFactorCard } from "@/components/profile/two-factor-card";
import { uploadFile } from "@/lib/api-client";
import { buttonClasses } from "@/components/ui/button";
import {
  User as UserIcon, Briefcase, Lock, Trash2, Save, Loader2, Upload, ShieldCheck,
} from "@/lib/icons";

const inputClass =
  "w-full rounded-lg border border-border bg-bg-secondary px-3 py-2.5 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent";
const errorInputClass =
  "w-full rounded-lg border border-danger/50 bg-bg-secondary px-3 py-2.5 text-sm text-foreground placeholder:text-text-muted focus:border-danger focus:outline-none focus:ring-1 focus:ring-danger";

export default function ProfilePage() {
  const { data: user } = useMe();
  const updateProfile = useUpdateProfile();
  const changePassword = useChangePassword();
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [avatarUploading, setAvatarUploading] = useState(false);
  const [avatarError, setAvatarError] = useState("");
  // The picked file waits here while it is positioned.
  const [pickedAvatar, setPickedAvatar] = useState<File | null>(null);
  const avatarInputRef = useRef<HTMLInputElement>(null);

  const personalForm = useForm<PersonalInfoValues>({
    resolver: zodResolver(PersonalInfoSchema),
    defaultValues: { first_name: "", last_name: "", email: "" },
  });
  const professionalForm = useForm<ProfessionalInfoValues>({
    resolver: zodResolver(ProfessionalInfoSchema),
    defaultValues: { job_title: "", bio: "" },
  });
  const passwordForm = useForm<ChangePasswordValues>({
    resolver: zodResolver(ChangePasswordSchema),
    defaultValues: { password: "", confirm_password: "" },
  });

  // biome-ignore lint/correctness/useExhaustiveDependencies: form.reset is stable; the forms refill when the user loads.
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

  // Pick, position, then upload. The file waits in state while the cropper is
  // open, so cancelling sends nothing.
  const handleAvatarUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    // Cleared so picking the same file twice reopens the cropper.
    e.target.value = "";
    if (!file) return;
    setAvatarError("");
    setPickedAvatar(file);
  };

  const handleAvatarCropped = async (file: File) => {
    setAvatarUploading(true);
    setAvatarError("");
    try {
      const result = await uploadFile(file);
      const url = (result.data as Record<string, unknown>)?.url as string;
      if (url) {
        updateProfile.mutate({ avatar: url });
      } else {
        setAvatarError("That upload came back without a URL. Try again.");
      }
    } catch {
      // Said out loud rather than swallowed: a picture that silently does not
      // change is a bug report nobody can describe.
      setAvatarError("That picture did not upload. Try again.");
    } finally {
      setAvatarUploading(false);
      setPickedAvatar(null);
    }
  };

  const onPersonalSubmit = (data: PersonalInfoValues) => updateProfile.mutate(data);
  const onProfessionalSubmit = (data: ProfessionalInfoValues) => updateProfile.mutate(data);
  const onPasswordSubmit = (data: ChangePasswordValues) =>
    changePassword.mutate(
      { password: data.password, current_password: data.current_password },
      { onSuccess: () => passwordForm.reset() },
    );

  if (!user) return null;

  const initials = ((user.first_name?.[0] || "") + (user.last_name?.[0] || "")).toUpperCase() || "U";
  const fullName = [user.first_name, user.last_name].filter(Boolean).join(" ");

  return (
    <div>
      <PageHeader title="Profile" subtitle="Manage your personal details, job profile, and password." />

      {/* Hero card — Profile Picture section with explicit header so the
          user knows what the Upload New button replaces. Layout: gradient
          banner up top, avatar floating half-out, name/email/role next to
          it, prominent Upload New CTA on the right. */}
      <section className="mb-6 overflow-hidden rounded-2xl border border-border bg-bg-elevated">
        <header className="border-b border-border px-6 py-4">
          <p className="text-sm font-semibold text-foreground">Profile Picture</p>
          <p className="mt-0.5 text-xs text-text-muted">
            Replace the image teammates and customers see on your profile.
          </p>
        </header>
        <div className="h-20 bg-gradient-to-r from-accent/30 via-accent/15 to-transparent" />
        <div className="-mt-12 flex flex-col items-start gap-4 px-6 pb-6 sm:flex-row sm:items-end">
          <div className="relative">
            <span className="block h-24 w-24 overflow-hidden rounded-2xl ring-4 ring-bg-elevated bg-bg-secondary">
              {user.avatar ? (
                <img src={user.avatar} alt={fullName || "Avatar"} className="h-full w-full object-cover" />
              ) : (
                <span className="flex h-full w-full items-center justify-center text-2xl font-bold text-foreground">
                  {initials}
                </span>
              )}
            </span>
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-lg font-bold text-foreground truncate">{fullName || "Anonymous"}</p>
            <p className="text-sm text-text-muted truncate">{user.email}</p>
            <span className="mt-1 inline-flex items-center gap-1 rounded-md bg-accent/10 px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide text-accent">
              {user.role}
            </span>
          </div>
          <div className="flex flex-wrap gap-2">
            <input ref={avatarInputRef} type="file" accept="image/*" className="hidden" onChange={handleAvatarUpload} />
            <button
              type="button"
              onClick={() => avatarInputRef.current?.click()}
              disabled={avatarUploading}
              className={buttonClasses({ className: "shadow-sm" })}
            >
              {avatarUploading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Upload className="h-4 w-4" />}
              Upload new
            </button>
            {avatarError && (
              <p role="alert" className="w-full text-xs text-danger">
                {avatarError}
              </p>
            )}
            <AvatarCropper
              file={pickedAvatar}
              saving={avatarUploading}
              onCancel={() => setPickedAvatar(null)}
              onCropped={handleAvatarCropped}
            />
          </div>
        </div>
      </section>

      {/* Personal Information — full-width card, grid form */}
      <ProfileCard
        icon={<UserIcon className="h-4 w-4" />}
        title="Personal information"
        description="Your name and primary email address."
        message={updateProfile.isSuccess ? "Saved" : undefined}
      >
        <form onSubmit={personalForm.handleSubmit(onPersonalSubmit)} className="space-y-4">
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <Field label="First name" error={personalForm.formState.errors.first_name?.message}>
              <input
                {...personalForm.register("first_name")}
                className={personalForm.formState.errors.first_name ? errorInputClass : inputClass}
                placeholder="Jane"
              />
            </Field>
            <Field label="Last name" error={personalForm.formState.errors.last_name?.message}>
              <input
                {...personalForm.register("last_name")}
                className={personalForm.formState.errors.last_name ? errorInputClass : inputClass}
                placeholder="Doe"
              />
            </Field>
          </div>
          <Field label="Email" error={personalForm.formState.errors.email?.message}>
            <input
              type="email"
              {...personalForm.register("email")}
              className={personalForm.formState.errors.email ? errorInputClass : inputClass}
              placeholder="you@example.com"
            />
          </Field>
          <Field label="Current password (only to change your email)">
            <input
              type="password"
              autoComplete="current-password"
              {...personalForm.register("current_password")}
              className={inputClass}
            />
          </Field>
          {updateProfile.isError && (
            <p className="text-sm text-danger">{(updateProfile.error as { response?: { data?: { error?: { message?: string } } } } | null)?.response?.data?.error?.message ?? "That did not save. Try again."}</p>
          )}
          <div className="flex justify-end">
            <SubmitButton pending={updateProfile.isPending}>
              <Save className="h-4 w-4" />
              Save changes
            </SubmitButton>
          </div>
        </form>
      </ProfileCard>

      {/* Professional Information */}
      <ProfileCard
        icon={<Briefcase className="h-4 w-4" />}
        title="Professional information"
        description="What you do, and a short bio teammates and customers can see."
      >
        <form onSubmit={professionalForm.handleSubmit(onProfessionalSubmit)} className="space-y-4">
          <Field label="Job title">
            <input
              {...professionalForm.register("job_title")}
              className={inputClass}
              placeholder="Software Engineer"
            />
          </Field>
          <Field label="Bio">
            <textarea
              {...professionalForm.register("bio")}
              rows={4}
              className={inputClass}
              placeholder="A short bio about yourself..."
            />
          </Field>
          <div className="flex justify-end">
            <SubmitButton pending={updateProfile.isPending}>
              <Save className="h-4 w-4" />
              Save changes
            </SubmitButton>
          </div>
        </form>
      </ProfileCard>

      {/* Password */}
      <ProfileCard
        icon={<Lock className="h-4 w-4" />}
        title="Password"
        description="Choose a new password that's at least 8 characters long."
        message={changePassword.isSuccess ? "Password updated" : undefined}
      >
        <form onSubmit={passwordForm.handleSubmit(onPasswordSubmit)} className="space-y-4">
          <Field label="Current password" error={passwordForm.formState.errors.current_password?.message}>
            <input
              type="password"
              autoComplete="current-password"
              {...passwordForm.register("current_password")}
              className={passwordForm.formState.errors.current_password ? errorInputClass : inputClass}
            />
          </Field>
          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <Field label="New password" error={passwordForm.formState.errors.password?.message}>
              <input
                type="password"
                {...passwordForm.register("password")}
                className={passwordForm.formState.errors.password ? errorInputClass : inputClass}
                placeholder="At least 8 characters"
              />
            </Field>
            <Field label="Confirm password" error={passwordForm.formState.errors.confirm_password?.message}>
              <input
                type="password"
                {...passwordForm.register("confirm_password")}
                className={passwordForm.formState.errors.confirm_password ? errorInputClass : inputClass}
                placeholder="Re-enter your password"
              />
            </Field>
          </div>
          {changePassword.isError && (
            <p className="text-sm text-danger">{(changePassword.error as { response?: { data?: { error?: { message?: string } } } } | null)?.response?.data?.error?.message ?? "That did not save. Try again."}</p>
          )}
          <div className="flex justify-end">
            <SubmitButton pending={changePassword.isPending}>
              <Lock className="h-4 w-4" />
              Update password
            </SubmitButton>
          </div>
        </form>
      </ProfileCard>

      {/* Two-factor — above sessions, because turning it on is what makes
          those sessions harder to steal. Renders its own card chrome, so it
          is deliberately not wrapped in a ProfileCard. */}
      <div className="mt-6">
        <TwoFactorCard />
      </div>

      {/* Active sessions */}
      <ProfileCard
        icon={<ShieldCheck className="h-4 w-4" />}
        title="Active sessions"
        description="Every device signed in to your account. Changing your password signs the others out automatically."
      >
        <ActiveSessions />
      </ProfileCard>

      {/* Danger zone */}
      <section className="mt-6 rounded-2xl border border-danger/30 bg-danger/5 p-6">
        <div className="flex flex-col items-start gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <p className="text-sm font-semibold text-foreground">Delete account</p>
            <p className="mt-1 text-sm text-text-secondary">
              Permanently remove your account and all associated data. This action cannot be undone.
            </p>
          </div>
          <button
            type="button"
            onClick={() => setShowDeleteDialog(true)}
            className="inline-flex items-center gap-2 rounded-lg border border-danger/40 bg-danger/10 px-4 py-2 text-sm font-semibold text-danger hover:bg-danger/20"
          >
            <Trash2 className="h-4 w-4" />
            Delete account
          </button>
        </div>
      </section>

      <DeleteAccountDialog open={showDeleteDialog} onClose={() => setShowDeleteDialog(false)} />
    </div>
  );
}

function ProfileCard({
  icon, title, description, message, children,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
  message?: string;
  children: React.ReactNode;
}) {
  return (
    <section className="mb-6 rounded-2xl border border-border bg-bg-elevated">
      <header className="flex items-start justify-between gap-3 border-b border-border px-6 py-4">
        <div className="flex items-start gap-3 min-w-0">
          <span className="inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg bg-accent/10 text-accent">
            {icon}
          </span>
          <div className="min-w-0">
            <p className="text-sm font-semibold text-foreground">{title}</p>
            <p className="mt-0.5 text-xs text-text-muted">{description}</p>
          </div>
        </div>
        {message && (
          <span className="shrink-0 rounded-md bg-success/10 px-2 py-1 text-xs font-medium text-success">
            {message}
          </span>
        )}
      </header>
      <div className="px-6 py-5">{children}</div>
    </section>
  );
}

function Field({ label, error, children }: { label: string; error?: string; children: React.ReactNode }) {
  return (
    <label className="block">
      <span className="mb-1.5 block text-xs font-semibold uppercase tracking-wide text-text-muted">{label}</span>
      {children}
      {error && <p className="mt-1 text-xs text-danger">{error}</p>}
    </label>
  );
}

function SubmitButton({ pending, children }: { pending: boolean; children: React.ReactNode }) {
  return (
    <button
      type="submit"
      disabled={pending}
      className={buttonClasses()}
    >
      {pending ? <Loader2 className="h-4 w-4 animate-spin" /> : children}
    </button>
  );
}
