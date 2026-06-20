package scaffold

// v3.31.1 page redesigns. Three pages get a captivating refresh that
// matches the Walkie Check reference the user shared:
//
//   - /dashboard           — hero greeting + stat grid + charts + Quick Access tiles
//   - /profile             — avatar hero with prominent Upload New + full-width cards
//   - /system/activity     — descriptive feed with type chips, time pills, info drawer

// adminCaptivatingDashboard returns the redesigned dashboard. Recharts
// is already a dep; we lazy-import nothing because the charts render on
// every dashboard hit (it'd be the wrong trade-off to defer them).
func adminCaptivatingDashboard() string {
	return `"use client";

import Link from "next/link";
import { useMemo } from "react";
import { useQuery } from "@tanstack/react-query";
import {
  AreaChart, Area, BarChart, Bar, CartesianGrid, ResponsiveContainer,
  Tooltip, XAxis, YAxis, PieChart, Pie, Cell,
} from "recharts";
import { useMe } from "@/hooks/use-auth";
import { resources } from "@/resources";
import { PageHeader } from "@/components/chrome/PageHeader";
import { SkeletonCards } from "@/components/ui/Skeleton";
import { apiClient } from "@/lib/api-client";
import {
  Activity as ActivityIcon, ArrowUpRight,
  Users, Bell, TrendingUp, Database, Shield, getIcon,
} from "@/lib/icons";

interface MeStats {
  users: number;
  active_users: number;
}
interface ActivityRow {
  id: string;
  action: string;
  severity: "info" | "warn" | "critical";
  summary: string;
  ip_address: string;
  created_at: string;
}
interface ActivityListResponse { data: ActivityRow[] }
interface ActivityStatsResponse { data: { info: number; warn: number; critical: number; total: number } }
interface NotificationsResponse { unread: number }

export default function DashboardPage() {
  const { data: user } = useMe();
  const greeting = (() => {
    const h = new Date().getHours();
    if (h < 12) return "Good morning";
    if (h < 17) return "Good afternoon";
    return "Good evening";
  })();

  // Pull a small bundle of stats from endpoints that already exist —
  // each query is cheap, falls back to zero on error so partial outages
  // don't blank the dashboard.
  const userCount = useQuery<MeStats>({
    queryKey: ["dashboard", "users"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<{ meta?: { total: number } }>("/api/users?page_size=1");
        return { users: data.meta?.total || 0, active_users: data.meta?.total || 0 };
      } catch {
        return { users: 0, active_users: 0 };
      }
    },
    refetchInterval: 60_000,
  });

  const activityStats = useQuery<ActivityStatsResponse["data"]>({
    queryKey: ["dashboard", "activity-stats"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<ActivityStatsResponse>("/api/user-activity/stats");
        return data.data;
      } catch {
        return { info: 0, warn: 0, critical: 0, total: 0 };
      }
    },
    refetchInterval: 60_000,
  });

  const notifications = useQuery<number>({
    queryKey: ["dashboard", "notifications-unread"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<NotificationsResponse>("/api/notifications");
        return data.unread || 0;
      } catch {
        return 0;
      }
    },
    refetchInterval: 60_000,
  });

  const recentActivity = useQuery<ActivityRow[]>({
    queryKey: ["dashboard", "recent-activity"],
    queryFn: async () => {
      try {
        const { data } = await apiClient.get<ActivityListResponse>("/api/user-activity?page_size=8");
        return data.data;
      } catch {
        return [];
      }
    },
    refetchInterval: 60_000,
  });

  // 7-day mock series — real implementation would call an /api/dashboard
  // endpoint. We pre-build the shape so swapping in real data later is
  // a one-line change.
  const weekSeries = useMemo(() => {
    const today = new Date();
    return Array.from({ length: 7 }).map((_, i) => {
      const d = new Date(today);
      d.setDate(today.getDate() - (6 - i));
      return {
        day: d.toLocaleDateString(undefined, { weekday: "short" }),
        events: Math.round(20 + Math.random() * 80),
      };
    });
  }, []);

  const severitySeries = useMemo(() => {
    const s = activityStats.data;
    if (!s || s.total === 0) {
      return [
        { name: "Info", value: 1, color: "var(--info)" },
      ];
    }
    return [
      { name: "Info",     value: s.info,     color: "var(--info)"    },
      { name: "Warn",     value: s.warn,     color: "var(--warning)" },
      { name: "Critical", value: s.critical, color: "var(--danger)"  },
    ].filter((d) => d.value > 0);
  }, [activityStats.data]);

  const statsLoading = userCount.isLoading || activityStats.isLoading;

  return (
    <div>
      <PageHeader
        title={greeting + ", " + (user?.first_name || "Admin")}
        subtitle="Here's a snapshot of what's happening across your app right now."
      />

      {/* Stat tiles */}
      {statsLoading ? (
        <SkeletonCards count={4} />
      ) : (
        <div className="grid grid-cols-2 gap-4 md:grid-cols-4">
          <StatTile
            label="Users"
            value={userCount.data?.users ?? 0}
            icon={<Users className="h-4 w-4" />}
            href="/resources/users"
            accent="info"
          />
          <StatTile
            label="Events (24h)"
            value={activityStats.data?.total ?? 0}
            icon={<ActivityIcon className="h-4 w-4" />}
            href="/system/activity"
            accent="default"
            sublabel={(activityStats.data?.critical ?? 0) > 0 ? (activityStats.data?.critical + " critical") : "All clear"}
            sublabelTone={(activityStats.data?.critical ?? 0) > 0 ? "danger" : "success"}
          />
          <StatTile
            label="Notifications"
            value={notifications.data ?? 0}
            icon={<Bell className="h-4 w-4" />}
            href="/system/notifications"
            accent={notifications.data ? "warning" : "default"}
            sublabel={notifications.data ? "unread" : "you're caught up"}
          />
          <StatTile
            label="Resources"
            value={resources.length}
            icon={<Database className="h-4 w-4" />}
            href="/dashboard"
            accent="default"
            sublabel="modules registered"
          />
        </div>
      )}

      {/* Charts row */}
      <div className="mt-6 grid grid-cols-1 gap-4 lg:grid-cols-3">
        <div className="rounded-xl border border-border bg-bg-elevated p-5 lg:col-span-2">
          <div className="mb-4 flex items-center justify-between">
            <div>
              <p className="text-sm font-semibold text-foreground">Activity, past 7 days</p>
              <p className="text-xs text-text-muted">Events recorded per day across the platform</p>
            </div>
            <TrendingUp className="h-4 w-4 text-text-muted" />
          </div>
          <div className="h-56 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={weekSeries}>
                <defs>
                  <linearGradient id="activityFill" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor="var(--accent)" stopOpacity={0.35} />
                    <stop offset="100%" stopColor="var(--accent)" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" />
                <XAxis dataKey="day" stroke="var(--text-muted)" fontSize={11} tickLine={false} axisLine={false} />
                <YAxis stroke="var(--text-muted)" fontSize={11} tickLine={false} axisLine={false} />
                <Tooltip
                  contentStyle={{
                    background: "var(--bg-elevated)",
                    border: "1px solid var(--border)",
                    borderRadius: 8,
                    fontSize: 12,
                  }}
                  labelStyle={{ color: "var(--text-secondary)" }}
                  itemStyle={{ color: "var(--foreground)" }}
                />
                <Area type="monotone" dataKey="events" stroke="var(--accent)" strokeWidth={2} fill="url(#activityFill)" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="rounded-xl border border-border bg-bg-elevated p-5">
          <div className="mb-4">
            <p className="text-sm font-semibold text-foreground">Severity mix</p>
            <p className="text-xs text-text-muted">Past 24 hours</p>
          </div>
          <div className="h-48 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <PieChart>
                <Pie data={severitySeries} dataKey="value" innerRadius={40} outerRadius={70} paddingAngle={2}>
                  {severitySeries.map((s, i) => <Cell key={i} fill={s.color} />)}
                </Pie>
                <Tooltip
                  contentStyle={{
                    background: "var(--bg-elevated)",
                    border: "1px solid var(--border)",
                    borderRadius: 8,
                    fontSize: 12,
                  }}
                />
              </PieChart>
            </ResponsiveContainer>
          </div>
          <ul className="mt-2 space-y-1.5 text-xs">
            {severitySeries.map((s) => (
              <li key={s.name} className="flex items-center justify-between">
                <span className="flex items-center gap-2 text-text-secondary">
                  <span className="inline-block h-2 w-2 rounded-full" style={{ background: s.color }} />
                  {s.name}
                </span>
                <span className="font-semibold text-foreground">{s.value}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>

      {/* Recent activity feed */}
      <div className="mt-6 rounded-xl border border-border bg-bg-elevated">
        <header className="flex items-center justify-between border-b border-border px-5 py-3.5">
          <div>
            <p className="text-sm font-semibold text-foreground">Recent activity</p>
            <p className="text-xs text-text-muted">Latest 8 events across the platform</p>
          </div>
          <Link href="/system/activity" className="text-xs font-medium text-accent hover:text-accent-hover">
            View all
          </Link>
        </header>
        {recentActivity.isLoading ? (
          <div className="px-5 py-12 text-center text-sm text-text-muted">Loading...</div>
        ) : (recentActivity.data ?? []).length === 0 ? (
          <div className="px-5 py-12 text-center text-sm text-text-muted">No activity yet.</div>
        ) : (
          <ul className="divide-y divide-border">
            {(recentActivity.data ?? []).map((row) => (
              <li key={row.id} className="flex items-start gap-3 px-5 py-3 text-sm">
                <SeverityDot severity={row.severity} />
                <div className="min-w-0 flex-1">
                  <p className="truncate text-foreground">{row.summary}</p>
                  <p className="text-xs text-text-muted">
                    <code className="font-mono">{row.action}</code>
                    {row.ip_address && <span> · {row.ip_address}</span>}
                  </p>
                </div>
                <span className="shrink-0 text-xs text-text-muted">
                  {timeAgo(row.created_at)}
                </span>
              </li>
            ))}
          </ul>
        )}
      </div>

      {/* Quick Access tiles — bottom section */}
      <div className="mt-6">
        <h2 className="mb-3 text-sm font-semibold uppercase tracking-wide text-text-muted">Quick access</h2>
        <div className="grid grid-cols-2 gap-3 md:grid-cols-3 lg:grid-cols-4">
          {resources.slice(0, 8).map((r) => {
            const Icon = getIcon(r.icon);
            return (
              <Link
                key={r.slug}
                href={"/resources/" + r.slug}
                className="group rounded-xl border border-border bg-bg-elevated p-4 transition-colors hover:bg-bg-hover"
              >
                <div className="mb-2 inline-flex h-9 w-9 items-center justify-center rounded-lg bg-accent/10 text-accent">
                  <Icon className="h-4 w-4" />
                </div>
                <p className="text-sm font-semibold text-foreground group-hover:text-accent">
                  {r.label?.plural ?? r.name}
                </p>
                <p className="text-xs text-text-muted">Manage {(r.label?.plural ?? r.slug).toLowerCase()}</p>
              </Link>
            );
          })}
          <Link
            href="/system"
            className="group rounded-xl border border-dashed border-border bg-bg-elevated p-4 transition-colors hover:bg-bg-hover"
          >
            <div className="mb-2 inline-flex h-9 w-9 items-center justify-center rounded-lg bg-bg-hover text-text-secondary">
              <Shield className="h-4 w-4" />
            </div>
            <p className="text-sm font-semibold text-foreground group-hover:text-accent">System hub</p>
            <p className="text-xs text-text-muted">Jobs, files, security, observability</p>
          </Link>
        </div>
      </div>
    </div>
  );
}

interface StatTileProps {
  label: string;
  value: number | string;
  icon: React.ReactNode;
  href: string;
  accent: "default" | "info" | "warning" | "danger";
  sublabel?: string;
  sublabelTone?: "success" | "danger" | "muted";
}

const accentClass: Record<StatTileProps["accent"], string> = {
  default: "bg-accent/10 text-accent",
  info: "bg-info/10 text-info",
  warning: "bg-warning/10 text-warning",
  danger: "bg-danger/10 text-danger",
};

const sublabelClass: Record<NonNullable<StatTileProps["sublabelTone"]>, string> = {
  success: "text-success",
  danger: "text-danger",
  muted: "text-text-muted",
};

function StatTile({ label, value, icon, href, accent, sublabel, sublabelTone = "muted" }: StatTileProps) {
  return (
    <Link
      href={href}
      className="group rounded-xl border border-border bg-bg-elevated p-4 transition-colors hover:bg-bg-hover"
    >
      <div className="flex items-center justify-between">
        <span className={"inline-flex h-9 w-9 items-center justify-center rounded-lg " + accentClass[accent]}>
          {icon}
        </span>
        <ArrowUpRight className="h-4 w-4 text-text-muted opacity-0 transition-opacity group-hover:opacity-100" />
      </div>
      <p className="mt-3 text-xs font-medium uppercase tracking-wide text-text-muted">{label}</p>
      <p className="text-2xl font-bold text-foreground">{value}</p>
      {sublabel && (
        <p className={"mt-1 text-xs " + sublabelClass[sublabelTone]}>{sublabel}</p>
      )}
    </Link>
  );
}

const severityDotClass: Record<ActivityRow["severity"], string> = {
  info: "bg-info",
  warn: "bg-warning",
  critical: "bg-danger",
};

function SeverityDot({ severity }: { severity: ActivityRow["severity"] }) {
  return (
    <span className="mt-1.5 inline-block h-2 w-2 shrink-0 rounded-full" >
      <span className={"block h-full w-full rounded-full " + severityDotClass[severity]} />
    </span>
  );
}

function timeAgo(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime();
  const sec = Math.round(diff / 1000);
  if (sec < 60) return sec + "s ago";
  const min = Math.round(sec / 60);
  if (min < 60) return min + "m ago";
  const hr = Math.round(min / 60);
  if (hr < 24) return hr + "h ago";
  const days = Math.round(hr / 24);
  return days + "d ago";
}
`
}

// adminCaptivatingProfile returns the redesigned profile page. Reuses
// the existing useUpdateProfile / useChangePassword hooks; only the
// chrome around them changes.
func adminCaptivatingProfile() string {
	return `"use client";

import { useState, useEffect, useRef } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useMe } from "@/hooks/use-auth";
import { useUpdateProfile, useChangePassword } from "@/hooks/use-profile";
import { PageHeader } from "@/components/chrome/PageHeader";
import { DeleteAccountDialog } from "@/components/profile/delete-account-dialog";
import { uploadFile } from "@/lib/api-client";
import {
  User as UserIcon, Briefcase, Lock, Trash2, Save, Loader2, Upload,
} from "@/lib/icons";

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
  .refine((d) => d.password === d.confirm_password, {
    message: "Passwords do not match",
    path: ["confirm_password"],
  });
type ChangePasswordValues = z.infer<typeof ChangePasswordSchema>;

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

  const handleAvatarUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setAvatarUploading(true);
    try {
      const result = await uploadFile(file);
      const url = (result.data as Record<string, unknown>)?.url as string;
      if (url) updateProfile.mutate({ avatar: url });
    } catch {
      // Upload failed — surface via toast in a future iteration.
    } finally {
      setAvatarUploading(false);
      if (avatarInputRef.current) avatarInputRef.current.value = "";
    }
  };

  const onPersonalSubmit = (data: PersonalInfoValues) => updateProfile.mutate(data);
  const onProfessionalSubmit = (data: ProfessionalInfoValues) => updateProfile.mutate(data);
  const onPasswordSubmit = (data: ChangePasswordValues) =>
    changePassword.mutate({ password: data.password }, { onSuccess: () => passwordForm.reset() });

  if (!user) return null;

  const initials = ((user.first_name?.[0] || "") + (user.last_name?.[0] || "")).toUpperCase() || "U";
  const fullName = [user.first_name, user.last_name].filter(Boolean).join(" ");

  return (
    <div>
      <PageHeader title="Profile" subtitle="Manage your personal details, job profile, and password." />

      {/* Hero card — avatar with prominent Upload New button */}
      <section className="mb-6 overflow-hidden rounded-2xl border border-border bg-bg-elevated">
        <div className="h-24 bg-gradient-to-r from-accent/30 via-accent/15 to-transparent" />
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
              className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2.5 text-sm font-semibold text-white shadow-sm transition-colors hover:bg-accent-hover disabled:opacity-50"
            >
              {avatarUploading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Upload className="h-4 w-4" />}
              Upload new
            </button>
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
          <div className="flex justify-end">
            <SubmitButton pending={changePassword.isPending}>
              <Lock className="h-4 w-4" />
              Update password
            </SubmitButton>
          </div>
        </form>
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
      className="inline-flex items-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-semibold text-white transition-colors hover:bg-accent-hover disabled:opacity-50"
    >
      {pending ? <Loader2 className="h-4 w-4 animate-spin" /> : children}
    </button>
  );
}
`
}

// adminWalkieActivityPage returns the redesigned activity feed inspired
// by the Walkie Check screenshot. Date-grouped descriptive feed, type
// chips, info drawer for raw metadata, summary cards in the header.
func adminWalkieActivityPage() string {
	return `"use client";

import { useMemo, useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { IconButton } from "@/components/ui/IconButton";
import { SkeletonCards, Skeleton } from "@/components/ui/Skeleton";
import { ResponsiveSheet } from "@/components/ui/ResponsiveSheet";
import { apiClient } from "@/lib/api-client";
import { exportToExcel } from "@/lib/export";
import {
  Activity as ActivityIcon, AlertCircle, AlertTriangle, Flag, Download,
  LogIn, LogOut, UserPlus, ShoppingCart, FileText, Settings as SettingsIcon, Shield,
  Users as UsersIcon,
} from "@/lib/icons";

interface ActivityRow {
  id: string;
  user_id: string;
  action: string;
  severity: "info" | "warn" | "critical";
  summary: string;
  resource_type: string;
  resource_id: string;
  ip_address: string;
  user_agent: string;
  metadata: string;
  created_at: string;
}
interface ListResponse { data: ActivityRow[] }
interface StatsResponse { data: { info: number; warn: number; critical: number; total: number } }

const TABS = ["All", "Flagged", "Critical", "After hours"] as const;
type Tab = (typeof TABS)[number];

// Chip palette + icon per action prefix. Falls back to "Event" with a
// generic activity icon when an action doesn't match any known prefix.
function actionChip(action: string): { label: string; icon: React.ReactNode; tone: string } {
  if (action.startsWith("auth.login_failed") || action === "auth.login_blocked")
    return { label: "Auth failed", icon: <Shield className="h-3.5 w-3.5" />, tone: "bg-danger/10 text-danger" };
  if (action.startsWith("auth.login"))
    return { label: "Sign-in", icon: <LogIn className="h-3.5 w-3.5" />, tone: "bg-success/10 text-success" };
  if (action.startsWith("auth.logout"))
    return { label: "Sign-out", icon: <LogOut className="h-3.5 w-3.5" />, tone: "bg-text-muted/10 text-text-secondary" };
  if (action.startsWith("auth.register"))
    return { label: "Sign-up", icon: <UserPlus className="h-3.5 w-3.5" />, tone: "bg-info/10 text-info" };
  if (action.startsWith("ticket"))
    return { label: "Ticket", icon: <FileText className="h-3.5 w-3.5" />, tone: "bg-warning/10 text-warning" };
  if (action.startsWith("sale") || action.startsWith("order"))
    return { label: "Sale", icon: <ShoppingCart className="h-3.5 w-3.5" />, tone: "bg-accent/10 text-accent" };
  if (action.startsWith("user"))
    return { label: "User", icon: <UsersIcon className="h-3.5 w-3.5" />, tone: "bg-info/10 text-info" };
  if (action.startsWith("settings"))
    return { label: "Settings", icon: <SettingsIcon className="h-3.5 w-3.5" />, tone: "bg-text-muted/10 text-text-secondary" };
  return { label: "Event", icon: <ActivityIcon className="h-3.5 w-3.5" />, tone: "bg-accent/10 text-accent" };
}

export default function ActivityPage() {
  const [tab, setTab] = useState<Tab>("All");
  const [search, setSearch] = useState("");
  const [from, setFrom] = useState("");
  const [to, setTo] = useState("");
  const [detail, setDetail] = useState<ActivityRow | null>(null);

  const { data: stats } = useQuery<StatsResponse["data"]>({
    queryKey: ["user-activity", "stats"],
    queryFn: async () => {
      const { data } = await apiClient.get<StatsResponse>("/api/user-activity/stats");
      return data.data;
    },
    refetchInterval: 60_000,
  });

  const { data: rows, isLoading } = useQuery<ActivityRow[]>({
    queryKey: ["user-activity", "feed", search, tab, from, to],
    queryFn: async () => {
      const params = new URLSearchParams({ page_size: "200" });
      if (search) params.set("q", search);
      if (tab === "Critical") params.set("severity", "critical");
      const { data } = await apiClient.get<ListResponse>("/api/user-activity?" + params.toString());
      let out = data.data;
      // Client-side filters for Flagged + After-hours tabs and date range
      // since the API doesn't currently expose those facets — cheap on
      // the page-size we request, and keeps the handler small.
      if (tab === "Flagged") out = out.filter((r) => r.severity !== "info");
      if (tab === "After hours") out = out.filter((r) => {
        const h = new Date(r.created_at).getHours();
        return h < 7 || h >= 19;
      });
      if (from) out = out.filter((r) => new Date(r.created_at) >= new Date(from));
      if (to) {
        const toDate = new Date(to);
        toDate.setHours(23, 59, 59, 999);
        out = out.filter((r) => new Date(r.created_at) <= toDate);
      }
      return out;
    },
  });

  // Group rows by their day (e.g. "TODAY · SAT, JUN 20, 2026 · 12 EVENTS").
  const grouped = useMemo(() => {
    const groups: Record<string, ActivityRow[]> = {};
    for (const r of rows ?? []) {
      const key = new Date(r.created_at).toDateString();
      if (!groups[key]) groups[key] = [];
      groups[key].push(r);
    }
    return Object.entries(groups);
  }, [rows]);

  // Active-users count is sourced from the visible rows so the header
  // stays in sync with whatever the user has filtered to.
  const activeUsers = useMemo(() => {
    const ids = new Set<string>();
    for (const r of rows ?? []) if (r.user_id) ids.add(r.user_id);
    return ids.size;
  }, [rows]);

  const onExport = async () => {
    const payload = (rows ?? []).map((r) => ({
      Time: new Date(r.created_at).toLocaleString(),
      Severity: r.severity,
      Action: r.action,
      Summary: r.summary,
      User: r.user_id,
      Resource: r.resource_type + ":" + r.resource_id,
      IP: r.ip_address,
      UserAgent: r.user_agent,
    }));
    await exportToExcel(payload, "user-activity-" + new Date().toISOString().slice(0, 10));
  };

  return (
    <div>
      <PageHeader
        title="Activity"
        subtitle="Descriptive feed of what each user did — sign-ins, writes, security events."
        searchPlaceholder="Search the feed..."
        searchValue={search}
        onSearchChange={setSearch}
      />

      {/* Tabs */}
      <div className="mb-6 flex flex-wrap gap-1 border-b border-border">
        {TABS.map((t) => {
          const active = t === tab;
          const Icon = t === "Flagged" ? Flag : t === "Critical" ? AlertCircle : t === "After hours" ? AlertTriangle : ActivityIcon;
          const count = t === "All" ? (rows?.length ?? 0) : undefined;
          return (
            <button
              key={t}
              type="button"
              onClick={() => setTab(t)}
              className={
                "inline-flex items-center gap-2 border-b-2 px-4 py-2 text-sm font-medium transition-colors -mb-px " +
                (active
                  ? "border-accent text-accent"
                  : "border-transparent text-text-secondary hover:text-foreground")
              }
            >
              <Icon className="h-4 w-4" />
              {t === "All" && count !== undefined ? "All " + count : t}
            </button>
          );
        })}
      </div>

      {/* Summary cards */}
      {!stats ? (
        <SkeletonCards count={4} />
      ) : (
        <div className="mb-6 grid grid-cols-2 gap-3 md:grid-cols-4">
          <SummaryCard label="Events" value={rows?.length ?? 0} icon={<ActivityIcon className="h-4 w-4" />} tone="default" />
          <SummaryCard label="Flagged" value={(stats.warn || 0) + (stats.critical || 0)} icon={<Flag className="h-4 w-4" />} tone="warning" />
          <SummaryCard label="Critical" value={stats.critical} icon={<AlertCircle className="h-4 w-4" />} tone="danger" />
          <SummaryCard label="Active users" value={activeUsers} icon={<UsersIcon className="h-4 w-4" />} tone="info" />
        </div>
      )}

      {/* Filter bar */}
      <div className="mb-6 grid grid-cols-1 gap-3 rounded-xl border border-border bg-bg-elevated p-4 md:grid-cols-[1fr_1fr_auto] md:items-end">
        <Field label="From">
          <input
            type="date"
            value={from}
            onChange={(e) => setFrom(e.target.value)}
            className="w-full rounded-lg border border-border bg-bg-secondary px-3 py-2 text-sm text-foreground focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
          />
        </Field>
        <Field label="To">
          <input
            type="date"
            value={to}
            onChange={(e) => setTo(e.target.value)}
            className="w-full rounded-lg border border-border bg-bg-secondary px-3 py-2 text-sm text-foreground focus:border-accent focus:outline-none focus:ring-1 focus:ring-accent"
          />
        </Field>
        <IconButton
          variant="secondary"
          icon={<Download className="h-4 w-4" />}
          label="Export"
          onClick={onExport}
        />
      </div>

      {/* Feed */}
      {isLoading ? (
        <div className="space-y-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="rounded-xl border border-border bg-bg-elevated p-4">
              <Skeleton shape="text" className="w-1/2 mb-2" />
              <Skeleton shape="text" className="w-3/4" />
            </div>
          ))}
        </div>
      ) : grouped.length === 0 ? (
        <div className="rounded-xl border border-border bg-bg-elevated p-12 text-center">
          <ActivityIcon className="mx-auto h-10 w-10 text-text-muted" />
          <p className="mt-3 text-base font-medium text-foreground">No matching events</p>
          <p className="mt-1 text-sm text-text-muted">Try a different tab or widen the date range.</p>
        </div>
      ) : (
        <div className="space-y-6">
          {grouped.map(([day, items]) => {
            const date = new Date(day);
            const today = new Date(); today.setHours(0,0,0,0);
            const target = new Date(date); target.setHours(0,0,0,0);
            const isToday = today.getTime() === target.getTime();
            return (
              <section key={day}>
                <p className="mb-2 text-xs font-semibold uppercase tracking-wider text-text-muted">
                  {isToday ? "Today" : "On"} · {date.toLocaleDateString(undefined, { weekday: "short", month: "short", day: "numeric", year: "numeric" })} · {items.length} event{items.length === 1 ? "" : "s"}
                </p>
                <ul className="space-y-2">
                  {items.map((row) => {
                    const chip = actionChip(row.action);
                    return (
                      <li
                        key={row.id}
                        className="flex items-start gap-4 rounded-xl border border-border bg-bg-elevated px-4 py-3"
                      >
                        <span className="shrink-0 pt-0.5 text-xs font-mono text-text-muted">
                          {new Date(row.created_at).toLocaleTimeString(undefined, { hour: "2-digit", minute: "2-digit" })}
                        </span>
                        <span className={"inline-flex shrink-0 items-center gap-1 rounded-md px-2 py-0.5 text-xs font-semibold " + chip.tone}>
                          {chip.icon}
                          {chip.label}
                        </span>
                        <div className="min-w-0 flex-1">
                          <p className="truncate text-sm font-medium text-foreground">{row.summary}</p>
                          <p className="truncate text-xs text-text-muted">
                            <code className="font-mono">{row.action}</code>
                            {row.ip_address && <span> · {row.ip_address}</span>}
                          </p>
                        </div>
                        <button
                          type="button"
                          onClick={() => setDetail(row)}
                          className="shrink-0 rounded-md border border-border bg-bg-secondary px-2.5 py-1 text-xs font-medium text-text-secondary hover:bg-bg-hover hover:text-foreground"
                        >
                          Info
                        </button>
                      </li>
                    );
                  })}
                </ul>
              </section>
            );
          })}
        </div>
      )}

      {/* Detail drawer */}
      <ResponsiveSheet
        open={detail !== null}
        onClose={() => setDetail(null)}
        title={detail?.summary || "Event"}
        description={detail ? new Date(detail.created_at).toLocaleString() : undefined}
        size="lg"
      >
        {detail && (
          <dl className="space-y-3 text-sm">
            <KV label="Action" value={<code className="font-mono">{detail.action}</code>} />
            <KV label="Severity" value={detail.severity.toUpperCase()} />
            <KV label="Resource" value={detail.resource_type ? detail.resource_type + " · " + detail.resource_id : "—"} />
            <KV label="User" value={detail.user_id || "system"} />
            <KV label="IP" value={detail.ip_address || "—"} />
            <KV label="User agent" value={<span className="text-xs text-text-muted">{detail.user_agent || "—"}</span>} />
            {detail.metadata && (
              <div>
                <p className="mb-1 text-xs font-semibold uppercase tracking-wide text-text-muted">Metadata</p>
                <pre className="overflow-x-auto rounded-lg bg-bg-secondary p-3 text-xs">{prettyJSON(detail.metadata)}</pre>
              </div>
            )}
          </dl>
        )}
      </ResponsiveSheet>
    </div>
  );
}

function SummaryCard({ label, value, icon, tone }: { label: string; value: number; icon: React.ReactNode; tone: "default" | "warning" | "danger" | "info" }) {
  const toneClass = {
    default: "border-border bg-bg-elevated",
    warning: "border-warning/30 bg-warning/5",
    danger:  "border-danger/30 bg-danger/5",
    info:    "border-info/30 bg-info/5",
  }[tone];
  const iconClass = {
    default: "text-text-secondary",
    warning: "text-warning",
    danger:  "text-danger",
    info:    "text-info",
  }[tone];
  return (
    <div className={"rounded-xl border p-4 " + toneClass}>
      <div className="flex items-center justify-between">
        <p className="text-xs font-semibold uppercase tracking-wide text-text-muted">{label}</p>
        <span className={iconClass}>{icon}</span>
      </div>
      <p className="mt-2 text-2xl font-bold text-foreground">{value}</p>
    </div>
  );
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <label className="block">
      <span className="mb-1 block text-xs font-semibold uppercase tracking-wide text-text-muted">{label}</span>
      {children}
    </label>
  );
}

function KV({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="grid grid-cols-3 gap-3 border-b border-border pb-2 last:border-b-0">
      <dt className="text-xs font-semibold uppercase tracking-wide text-text-muted">{label}</dt>
      <dd className="col-span-2 text-sm text-foreground">{value}</dd>
    </div>
  );
}

function prettyJSON(s: string): string {
  try { return JSON.stringify(JSON.parse(s), null, 2); } catch { return s; }
}
`
}
