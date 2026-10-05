"use client";

import { useEffect, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { Loader2, Trash2, Upload, User } from "@/lib/icons";
import { inputClasses } from "@/components/ui/input";
import { useMe } from "@/hooks/use-auth";
import { useUpdateProfile } from "@/hooks/use-profile";
import { DeleteAccountDialog } from "@/components/profile/delete-account-dialog";
import { buttonClasses } from "@/components/ui/button";
import { AvatarCropper } from "@/components/account/avatar-cropper";
import { uploadFile } from "@/lib/api-client";

interface Values {
  first_name: string;
  last_name: string;
  email: string;
  job_title: string;
  bio: string;
  current_password: string;
}

export function ProfileForm() {
  const { data: user } = useMe();
  const updateProfile = useUpdateProfile();
  const [uploading, setUploading] = useState(false);
  const [uploadError, setUploadError] = useState("");
  // The picked file waits here while it is positioned. Nothing is uploaded until
  // the cropper hands back a square, so cancelling sends nothing at all.
  const [picked, setPicked] = useState<File | null>(null);
  const avatarInput = useRef<HTMLInputElement>(null);

  const { register, handleSubmit, reset, watch } = useForm<Values>({
    defaultValues: { first_name: "", last_name: "", email: "", job_title: "", bio: "", current_password: "" },
  });

  useEffect(() => {
    if (!user) return;
    const extra = user as { job_title?: string; bio?: string };
    reset({
      first_name: user.first_name ?? "",
      last_name: user.last_name ?? "",
      email: user.email ?? "",
      job_title: extra.job_title ?? "",
      bio: extra.bio ?? "",
      current_password: "",
    });
  }, [user, reset]);

  const emailChanged = (watch("email") ?? "") !== (user?.email ?? "");

  function onAvatarPicked(event: React.ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    // Cleared so picking the same file twice still opens the cropper: the input
    // fires no change event when the value has not changed.
    event.target.value = "";
    if (!file) return;
    setUploadError("");
    setPicked(file);
  }

  async function onAvatarCropped(file: File) {
    setUploading(true);
    setUploadError("");
    try {
      const result = await uploadFile(file);
      const url = (result.data as Record<string, unknown>)?.url as string | undefined;
      if (url) {
        updateProfile.mutate({ avatar: url });
      } else {
        setUploadError("That upload came back without a URL. Try again.");
      }
    } catch {
      // Said out loud rather than swallowed: a picture that silently does not
      // change is a bug report nobody can describe.
      setUploadError("That picture did not upload. Try again.");
    } finally {
      setUploading(false);
      setPicked(null);
    }
  }

  function onSubmit(values: Values) {
    const payload: Record<string, string> = {
      first_name: values.first_name,
      last_name: values.last_name,
      job_title: values.job_title,
      bio: values.bio,
    };
    if (emailChanged) {
      payload.email = values.email;
      payload.current_password = values.current_password;
    }
    updateProfile.mutate(payload);
  }

  const field = inputClasses();

  return (
    <section className="overflow-hidden rounded-xl border border-border bg-bg-elevated">
      <div className="flex items-start gap-3 border-b border-border px-6 py-4">
        <span className="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-accent/10 text-accent">
          <User className="h-4 w-4" aria-hidden="true" />
        </span>
        <div className="min-w-0">
          <h2 className="text-base font-semibold text-foreground">Profile</h2>
          <p className="mt-1 text-sm leading-relaxed text-text-muted">
            Your name and picture are what other people in this app see.
          </p>
        </div>
      </div>

      <div className="p-6">
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-5">
            <div className="flex items-center gap-4">
              <span className="inline-flex h-16 w-16 shrink-0 items-center justify-center overflow-hidden rounded-full bg-bg-tertiary text-lg font-semibold text-text-secondary">
                {user?.avatar ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img src={user.avatar} alt="" className="h-full w-full object-cover" />
                ) : (
                  (user?.first_name?.[0] ?? "") + (user?.last_name?.[0] ?? "")
                )}
              </span>
              <div>
                <input
                  ref={avatarInput}
                  id="account-avatar"
                  type="file"
                  accept="image/*"
                  className="sr-only"
                  onChange={onAvatarPicked}
                />
                <button
                  type="button"
                  onClick={() => avatarInput.current?.click()}
                  disabled={uploading}
                  className={buttonClasses({ variant: "outline", size: "sm" })}
                >
                  {uploading ? (
                    <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />
                  ) : (
                    <Upload className="h-4 w-4" aria-hidden="true" />
                  )}
                  Change picture
                </button>
                {uploadError && (
                  <p role="alert" className="mt-1.5 text-xs text-danger">
                    {uploadError}
                  </p>
                )}
                <AvatarCropper
                  file={picked}
                  saving={uploading}
                  onCancel={() => setPicked(null)}
                  onCropped={onAvatarCropped}
                />
              </div>
            </div>

            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-1.5">
                <label htmlFor="account-first-name" className="block text-sm font-medium text-text-secondary">
                  First name
                </label>
                <input id="account-first-name" className={field} {...register("first_name")} />
              </div>
              <div className="space-y-1.5">
                <label htmlFor="account-last-name" className="block text-sm font-medium text-text-secondary">
                  Last name
                </label>
                <input id="account-last-name" className={field} {...register("last_name")} />
              </div>
            </div>

            <div className="space-y-1.5">
              <label htmlFor="account-job-title" className="block text-sm font-medium text-text-secondary">
                Job title
              </label>
              <input id="account-job-title" className={field} {...register("job_title")} />
            </div>

            <div className="space-y-1.5">
              <label htmlFor="account-bio" className="block text-sm font-medium text-text-secondary">
                Bio
              </label>
              <textarea id="account-bio" rows={3} className={field} {...register("bio")} />
            </div>

            <div className="space-y-1.5">
              <label htmlFor="account-email" className="block text-sm font-medium text-text-secondary">
                Email
              </label>
              <input id="account-email" type="email" autoComplete="email" className={field} {...register("email")} />
            </div>

            {emailChanged && (
              <div className="space-y-1.5 rounded-lg border border-warning/40 bg-warning/5 p-4">
                <label htmlFor="account-email-password" className="block text-sm font-medium text-text-secondary">
                  Current password
                </label>
                <input
                  id="account-email-password"
                  type="password"
                  autoComplete="current-password"
                  className={field}
                  {...register("current_password")}
                />
                <p className="text-xs text-text-muted">
                  Changing your email changes where a password reset goes, so it asks who you are first.
                </p>
              </div>
            )}

            {updateProfile.isError && (
              <p role="alert" className="text-sm text-danger">
                {(updateProfile.error as { response?: { data?: { error?: { message?: string } } } } | null)
                  ?.response?.data?.error?.message ?? "That did not save. Try again."}
              </p>
            )}
            {updateProfile.isSuccess && (
              <p role="status" className="text-sm text-success">
                Saved.
              </p>
            )}

            <button type="submit" disabled={updateProfile.isPending} className={buttonClasses({ variant: "primary" })}>
              {updateProfile.isPending && <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />}
              Save changes
            </button>
          </form>
      </div>
    </section>
  );
}

/**
 * Closing the account, in a card of its own.
 *
 * Split out of ProfileForm so the page can put it where it belongs, which is
 * last. While it lived inside the component that draws the first card, the
 * only destructive control on the screen sat in the middle of it.
 */
export function CloseAccountCard() {
  const [confirmDelete, setConfirmDelete] = useState(false);

  return (
    <section className="overflow-hidden rounded-xl border border-danger/40 bg-danger/5">
      <div className="flex items-start gap-3 border-b border-danger/20 px-6 py-4">
        <span className="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-danger/10 text-danger">
          <Trash2 className="h-4 w-4" aria-hidden="true" />
        </span>
        <div className="min-w-0">
          <h2 className="text-base font-semibold text-foreground">Close this account</h2>
          <p className="mt-1 text-sm leading-relaxed text-text-muted">
            Your account and the data attached to it. This cannot be undone.
          </p>
        </div>
      </div>

      <div className="p-6">
        <button
          type="button"
          onClick={() => setConfirmDelete(true)}
          className={buttonClasses({ variant: "danger" })}
        >
          Close my account
        </button>
      </div>

      <DeleteAccountDialog open={confirmDelete} onClose={() => setConfirmDelete(false)} />
    </section>
  );
}
