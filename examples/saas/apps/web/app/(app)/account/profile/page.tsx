"use client";

import { useEffect, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { UserCog } from "lucide-react";
import { api } from "@/lib/api";
import { useMe } from "@/hooks/use-auth";
import { PasswordStrength } from "@/components/password-strength";

// One endpoint does all of this: PUT /api/profile takes the name fields, the
// email, and a password when the customer wants to change it. An empty password
// is left out of the payload rather than sent as an empty string, which the API
// would hash.
export default function AccountProfilePage() {
  const { data: user } = useMe();
  const queryClient = useQueryClient();

  const [form, setForm] = useState({
    first_name: "",
    last_name: "",
    email: "",
    job_title: "",
    password: "",
    current_password: "",
  });
  const [saved, setSaved] = useState(false);

  useEffect(() => {
    if (!user) return;
    setForm({
      first_name: user.first_name ?? "",
      last_name: user.last_name ?? "",
      email: user.email ?? "",
      job_title: user.job_title ?? "",
      password: "",
      current_password: "",
    });
  }, [user]);

  const save = useMutation({
    mutationFn: async () => {
      const payload: Record<string, string> = {
        first_name: form.first_name,
        last_name: form.last_name,
        email: form.email,
        job_title: form.job_title,
      };
      if (form.password) payload.password = form.password;
      // The API asks for it to change the email or the password.
      if (form.current_password) payload.current_password = form.current_password;
      const { data } = await api.put("/api/profile", payload);
      return data;
    },
    onSuccess: () => {
      setSaved(true);
      setForm((f) => ({ ...f, password: "", current_password: "" }));
      queryClient.invalidateQueries({ queryKey: ["me"] });
      window.setTimeout(() => setSaved(false), 4000);
    },
  });

  const field = (key: keyof typeof form) => ({
    value: form[key],
    onChange: (e: React.ChangeEvent<HTMLInputElement>) =>
      setForm((f) => ({ ...f, [key]: e.target.value })),
    className:
      "w-full rounded-lg border border-border bg-background px-3 py-2 text-sm text-foreground outline-none transition-colors focus:border-accent",
  });

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4">
        <span className="flex h-11 w-11 items-center justify-center rounded-xl bg-accent/10 text-accent">
          <UserCog className="h-5 w-5" />
        </span>
        <div>
          <h1 className="text-xl font-bold tracking-tight text-foreground">Profile</h1>
          <p className="text-sm text-text-secondary">
            This is what the rest of the app knows about you
          </p>
        </div>
      </div>

      <form
        className="space-y-5 rounded-xl border border-border bg-background p-6"
        onSubmit={(e) => {
          e.preventDefault();
          save.mutate();
        }}
      >
        <div className="grid gap-5 sm:grid-cols-2">
          <label className="block space-y-1.5">
            <span className="text-sm font-medium text-foreground">First name</span>
            <input {...field("first_name")} />
          </label>
          <label className="block space-y-1.5">
            <span className="text-sm font-medium text-foreground">Last name</span>
            <input {...field("last_name")} />
          </label>
        </div>

        <label className="block space-y-1.5">
          <span className="text-sm font-medium text-foreground">Email</span>
          <input type="email" {...field("email")} />
        </label>

        <label className="block space-y-1.5">
          <span className="text-sm font-medium text-foreground">Job title</span>
          <input {...field("job_title")} />
        </label>

        <label className="block space-y-1.5">
          <span className="text-sm font-medium text-foreground">New password</span>
          <input
            type="password"
            autoComplete="new-password"
            aria-describedby="password-strength"
            {...field("password")}
          />
          <span className="block text-xs text-text-secondary">
            Leave this empty to keep the password you have.
          </span>
          <PasswordStrength
            value={form.password}
            about={[form.email, form.first_name, form.last_name]}
          />
        </label>

        <label className="block space-y-1.5">
          <span className="text-sm font-medium text-foreground">Current password</span>
          <input type="password" autoComplete="current-password" {...field("current_password")} />
          <span className="block text-xs text-text-secondary">
            Needed only to change your email or your password.
          </span>
        </label>

        <div className="flex items-center gap-3 pt-1">
          <button
            type="submit"
            disabled={save.isPending}
            className="inline-flex items-center rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white transition-colors hover:bg-accent-hover disabled:opacity-60"
          >
            {save.isPending ? "Saving..." : "Save changes"}
          </button>
          {saved ? <span className="text-sm text-text-secondary">Saved.</span> : null}
          {save.isError ? (
            <span className="text-sm text-danger">
              {(save.error as { response?: { data?: { error?: { message?: string } } } } | null)?.response?.data?.error
                ?.message ?? "That did not save. Check the fields and try again."}
            </span>
          ) : null}
        </div>
      </form>
    </div>
  );
}
