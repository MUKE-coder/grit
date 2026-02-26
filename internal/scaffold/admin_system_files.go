package scaffold

func adminUseSystem() string {
	return `import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient, uploadFile } from "@/lib/api-client";

// ── Jobs ────────────────────────────────────────────────────────

interface QueueStats {
  queue: string;
  size: number;
  active: number;
  pending: number;
  completed: number;
  failed: number;
  retry: number;
  scheduled: number;
  processed: number;
}

interface Job {
  id: string;
  type: string;
  queue: string;
  max_retry: number;
  retried: number;
  last_error: string;
}

export function useJobStats() {
  return useQuery<QueueStats[]>({
    queryKey: ["admin", "jobs", "stats"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/admin/jobs/stats");
      return data.data;
    },
    refetchInterval: 5000,
  });
}

export function useJobsByStatus(status: string, queue = "default") {
  return useQuery<Job[]>({
    queryKey: ["admin", "jobs", status, queue],
    queryFn: async () => {
      const { data } = await apiClient.get(` + "`" + `/api/admin/jobs/${status}?queue=${queue}` + "`" + `);
      return data.data;
    },
    refetchInterval: 5000,
  });
}

export function useRetryJob() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ id, queue }: { id: string; queue?: string }) => {
      await apiClient.post(` + "`" + `/api/admin/jobs/${id}/retry?queue=${queue || "default"}` + "`" + `);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "jobs"] });
    },
  });
}

export function useClearQueue() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (queue: string) => {
      await apiClient.delete(` + "`" + `/api/admin/jobs/queue/${queue}` + "`" + `);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "jobs"] });
    },
  });
}

// ── Files ───────────────────────────────────────────────────────

interface Upload {
  id: number;
  filename: string;
  original_name: string;
  mime_type: string;
  size: number;
  path: string;
  url: string;
  thumbnail_url?: string;
  user_id: number;
  created_at: string;
  updated_at: string;
}

interface UploadListResponse {
  data: Upload[];
  meta: {
    total: number;
    page: number;
    page_size: number;
    pages: number;
  };
}

export function useUploads(page = 1, pageSize = 20) {
  return useQuery<UploadListResponse>({
    queryKey: ["admin", "uploads", page, pageSize],
    queryFn: async () => {
      const { data } = await apiClient.get(` + "`" + `/api/uploads?page=${page}&page_size=${pageSize}` + "`" + `);
      return data;
    },
  });
}

export function useUploadFile() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (file: File) => {
      const result = await uploadFile(file);
      return result.data as Upload;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "uploads"] });
    },
  });
}

export function useDeleteUpload() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (id: number) => {
      await apiClient.delete(` + "`" + `/api/uploads/${id}` + "`" + `);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "uploads"] });
    },
  });
}

// ── Cron ────────────────────────────────────────────────────────

interface CronTask {
  name: string;
  schedule: string;
  type: string;
}

export function useCronTasks() {
  return useQuery<CronTask[]>({
    queryKey: ["admin", "cron", "tasks"],
    queryFn: async () => {
      const { data } = await apiClient.get("/api/admin/cron/tasks");
      return data.data;
    },
  });
}
`
}

func adminJobsPage() string {
	return `"use client";

import { useState } from "react";
import { useJobStats, useJobsByStatus, useRetryJob, useClearQueue } from "@/hooks/use-system";
import { Briefcase, RefreshCw, Trash2, Loader2, AlertCircle } from "@/lib/icons";

const statuses = ["active", "pending", "completed", "failed", "retry"] as const;

export default function JobsPage() {
  const [activeTab, setActiveTab] = useState<string>("active");
  const { data: stats, isLoading: statsLoading } = useJobStats();
  const { data: jobs, isLoading: jobsLoading } = useJobsByStatus(activeTab);
  const retryJob = useRetryJob();
  const clearQueue = useClearQueue();

  const totalStats = stats?.[0];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Background Jobs</h1>
          <p className="text-sm text-text-secondary mt-1">Monitor and manage background job queues</p>
        </div>
      </div>

      {/* Stats cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4">
        {[
          { label: "Active", value: totalStats?.active ?? 0, color: "text-info" },
          { label: "Pending", value: totalStats?.pending ?? 0, color: "text-warning" },
          { label: "Completed", value: totalStats?.completed ?? 0, color: "text-success" },
          { label: "Failed", value: totalStats?.failed ?? 0, color: "text-danger" },
          { label: "Retry", value: totalStats?.retry ?? 0, color: "text-accent" },
          { label: "Processed", value: totalStats?.processed ?? 0, color: "text-foreground" },
        ].map((stat) => (
          <div key={stat.label} className="rounded-xl border border-border bg-bg-secondary p-4">
            <p className="text-xs text-text-muted uppercase tracking-wider">{stat.label}</p>
            <p className={` + "`" + `text-2xl font-bold mt-1 ${stat.color}` + "`" + `}>
              {statsLoading ? "—" : stat.value.toLocaleString()}
            </p>
          </div>
        ))}
      </div>

      {/* Tab navigation */}
      <div className="flex items-center gap-1 border-b border-border">
        {statuses.map((status) => (
          <button
            key={status}
            onClick={() => setActiveTab(status)}
            className={` + "`" + `px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${
              activeTab === status
                ? "border-accent text-accent"
                : "border-transparent text-text-secondary hover:text-foreground"
            }` + "`" + `}
          >
            {status.charAt(0).toUpperCase() + status.slice(1)}
          </button>
        ))}

        <div className="ml-auto">
          <button
            onClick={() => clearQueue.mutate("default")}
            className="flex items-center gap-2 px-3 py-1.5 text-sm text-text-secondary hover:text-danger transition-colors"
          >
            <Trash2 className="h-3.5 w-3.5" />
            Clear completed
          </button>
        </div>
      </div>

      {/* Job list */}
      <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="border-b border-border">
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">ID</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">Type</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">Queue</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">Retries</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">Error</th>
              <th className="px-4 py-3 text-right text-xs font-medium text-text-muted uppercase">Actions</th>
            </tr>
          </thead>
          <tbody>
            {jobsLoading ? (
              <tr>
                <td colSpan={6} className="px-4 py-12 text-center">
                  <Loader2 className="h-6 w-6 animate-spin text-accent mx-auto" />
                </td>
              </tr>
            ) : jobs && jobs.length > 0 ? (
              jobs.map((job) => (
                <tr key={job.id} className="border-b border-border last:border-0 hover:bg-bg-hover transition-colors">
                  <td className="px-4 py-3 text-sm font-mono text-text-secondary">{job.id.slice(0, 8)}...</td>
                  <td className="px-4 py-3">
                    <span className="rounded-md bg-accent/10 px-2 py-0.5 text-xs font-medium text-accent">
                      {job.type}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-sm text-text-secondary">{job.queue}</td>
                  <td className="px-4 py-3 text-sm text-text-secondary">{job.retried}/{job.max_retry}</td>
                  <td className="px-4 py-3 text-sm text-danger max-w-[200px] truncate">{job.last_error || "—"}</td>
                  <td className="px-4 py-3 text-right">
                    {(activeTab === "failed" || activeTab === "retry") && (
                      <button
                        onClick={() => retryJob.mutate({ id: job.id, queue: job.queue })}
                        className="inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-medium text-accent hover:bg-accent/10 rounded-md transition-colors"
                      >
                        <RefreshCw className="h-3 w-3" />
                        Retry
                      </button>
                    )}
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={6} className="px-4 py-12 text-center">
                  <Briefcase className="h-8 w-8 text-text-muted mx-auto mb-2" />
                  <p className="text-sm text-text-secondary">No {activeTab} jobs</p>
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
`
}

func adminFilesPage() string {
	return `"use client";

import { useState, useRef } from "react";
import { useUploads, useUploadFile, useDeleteUpload } from "@/hooks/use-system";
import { FolderOpen, Upload, Trash2, Loader2, X } from "@/lib/icons";

function formatFileSize(bytes: number): string {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + " " + sizes[i];
}

function formatDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString("en-US", {
    month: "short",
    day: "numeric",
    year: "numeric",
  });
}

export default function FilesPage() {
  const [page, setPage] = useState(1);
  const { data, isLoading } = useUploads(page);
  const uploadFile = useUploadFile();
  const deleteUpload = useDeleteUpload();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [dragOver, setDragOver] = useState(false);

  const uploads = data?.data ?? [];
  const meta = data?.meta;

  const handleFiles = (files: FileList | null) => {
    if (!files) return;
    Array.from(files).forEach((file) => uploadFile.mutate(file));
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setDragOver(false);
    handleFiles(e.dataTransfer.files);
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">File Storage</h1>
          <p className="text-sm text-text-secondary mt-1">Manage uploaded files and images</p>
        </div>
        <button
          onClick={() => fileInputRef.current?.click()}
          className="flex items-center gap-2 rounded-lg bg-accent px-4 py-2.5 text-sm font-medium text-white hover:bg-accent-hover transition-colors"
        >
          <Upload className="h-4 w-4" />
          Upload File
        </button>
        <input
          ref={fileInputRef}
          type="file"
          className="hidden"
          multiple
          onChange={(e) => handleFiles(e.target.files)}
        />
      </div>

      {/* Drop zone */}
      <div
        onDragOver={(e) => { e.preventDefault(); setDragOver(true); }}
        onDragLeave={() => setDragOver(false)}
        onDrop={handleDrop}
        className={` + "`" + `rounded-xl border-2 border-dashed p-8 text-center transition-colors ${
          dragOver ? "border-accent bg-accent/5" : "border-border"
        }` + "`" + `}
      >
        <Upload className="h-8 w-8 text-text-muted mx-auto mb-2" />
        <p className="text-sm text-text-secondary">
          Drag & drop files here, or{" "}
          <button
            onClick={() => fileInputRef.current?.click()}
            className="text-accent hover:underline"
          >
            browse
          </button>
        </p>
        <p className="text-xs text-text-muted mt-1">Max 10 MB per file</p>
      </div>

      {/* Upload progress */}
      {uploadFile.isPending && (
        <div className="flex items-center gap-3 rounded-lg border border-border bg-bg-secondary p-4">
          <Loader2 className="h-5 w-5 animate-spin text-accent" />
          <span className="text-sm text-text-secondary">Uploading...</span>
        </div>
      )}

      {/* File grid */}
      {isLoading ? (
        <div className="flex justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin text-accent" />
        </div>
      ) : uploads.length === 0 ? (
        <div className="flex flex-col items-center py-16">
          <FolderOpen className="h-12 w-12 text-text-muted mb-3" />
          <p className="text-sm text-text-secondary">No files uploaded yet</p>
        </div>
      ) : (
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4">
          {uploads.map((upload) => (
            <div
              key={upload.id}
              className="group relative rounded-xl border border-border bg-bg-secondary overflow-hidden hover:border-accent/50 transition-colors"
            >
              {/* Preview */}
              <div className="aspect-square bg-bg-tertiary flex items-center justify-center">
                {upload.mime_type.startsWith("image/") ? (
                  <img
                    src={upload.thumbnail_url || upload.url}
                    alt={upload.original_name}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <div className="text-center p-4">
                    <FolderOpen className="h-8 w-8 text-text-muted mx-auto mb-1" />
                    <p className="text-xs text-text-muted uppercase">
                      {upload.mime_type.split("/")[1]?.slice(0, 4)}
                    </p>
                  </div>
                )}
              </div>

              {/* Info */}
              <div className="p-3">
                <p className="text-sm font-medium text-foreground truncate">{upload.original_name}</p>
                <div className="flex items-center justify-between mt-1">
                  <span className="text-xs text-text-muted">{formatFileSize(upload.size)}</span>
                  <span className="text-xs text-text-muted">{formatDate(upload.created_at)}</span>
                </div>
              </div>

              {/* Delete overlay */}
              <button
                onClick={() => deleteUpload.mutate(upload.id)}
                className="absolute top-2 right-2 p-1.5 rounded-lg bg-black/50 text-white opacity-0 group-hover:opacity-100 hover:bg-danger transition-all"
              >
                <X className="h-3.5 w-3.5" />
              </button>
            </div>
          ))}
        </div>
      )}

      {/* Pagination */}
      {meta && meta.pages > 1 && (
        <div className="flex items-center justify-center gap-2">
          <button
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page === 1}
            className="px-3 py-1.5 text-sm rounded-lg border border-border text-text-secondary hover:bg-bg-hover disabled:opacity-50 transition-colors"
          >
            Previous
          </button>
          <span className="text-sm text-text-secondary">
            Page {page} of {meta.pages}
          </span>
          <button
            onClick={() => setPage((p) => Math.min(meta.pages, p + 1))}
            disabled={page === meta.pages}
            className="px-3 py-1.5 text-sm rounded-lg border border-border text-text-secondary hover:bg-bg-hover disabled:opacity-50 transition-colors"
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}
`
}

func adminCronPage() string {
	return `"use client";

import { useCronTasks } from "@/hooks/use-system";
import { Calendar, Loader2 } from "@/lib/icons";

export default function CronPage() {
  const { data: tasks, isLoading } = useCronTasks();

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-foreground">Cron Scheduler</h1>
        <p className="text-sm text-text-secondary mt-1">View registered scheduled tasks</p>
      </div>

      <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
        <table className="w-full">
          <thead>
            <tr className="border-b border-border">
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">Task</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">Schedule</th>
              <th className="px-4 py-3 text-left text-xs font-medium text-text-muted uppercase">Type</th>
            </tr>
          </thead>
          <tbody>
            {isLoading ? (
              <tr>
                <td colSpan={3} className="px-4 py-12 text-center">
                  <Loader2 className="h-6 w-6 animate-spin text-accent mx-auto" />
                </td>
              </tr>
            ) : tasks && tasks.length > 0 ? (
              tasks.map((task, i) => (
                <tr key={i} className="border-b border-border last:border-0 hover:bg-bg-hover transition-colors">
                  <td className="px-4 py-3">
                    <span className="text-sm font-medium text-foreground">{task.name}</span>
                  </td>
                  <td className="px-4 py-3">
                    <code className="rounded-md bg-bg-tertiary px-2 py-1 text-xs font-mono text-accent">
                      {task.schedule}
                    </code>
                  </td>
                  <td className="px-4 py-3">
                    <span className="rounded-md bg-accent/10 px-2 py-0.5 text-xs font-medium text-accent">
                      {task.type}
                    </span>
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan={3} className="px-4 py-12 text-center">
                  <Calendar className="h-8 w-8 text-text-muted mx-auto mb-2" />
                  <p className="text-sm text-text-secondary">No cron tasks registered</p>
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
`
}

func adminMailPage() string {
	return `"use client";

import { useState } from "react";
import { Mail } from "@/lib/icons";

const templates = [
  {
    name: "welcome",
    label: "Welcome Email",
    description: "Sent when a new user registers",
    sampleData: {
      AppName: "MyApp",
      Name: "John Doe",
      DashboardURL: "http://localhost:3000/dashboard",
      Year: new Date().getFullYear(),
    },
  },
  {
    name: "password-reset",
    label: "Password Reset",
    description: "Sent when a user requests a password reset",
    sampleData: {
      AppName: "MyApp",
      ResetURL: "http://localhost:3000/reset-password?token=abc123",
      Year: new Date().getFullYear(),
    },
  },
  {
    name: "email-verification",
    label: "Email Verification",
    description: "Sent to verify a user's email address",
    sampleData: {
      AppName: "MyApp",
      VerifyURL: "http://localhost:3000/verify?token=abc123",
      Year: new Date().getFullYear(),
    },
  },
  {
    name: "notification",
    label: "Notification",
    description: "General purpose notification email",
    sampleData: {
      AppName: "MyApp",
      Title: "New Activity",
      Message: "Someone commented on your post. Check it out!",
      ActionURL: "http://localhost:3000/activity",
      ActionText: "View Activity",
      Year: new Date().getFullYear(),
    },
  },
];

export default function MailPage() {
  const [selected, setSelected] = useState(templates[0]);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-foreground">Email Templates</h1>
        <p className="text-sm text-text-secondary mt-1">Preview email templates with sample data</p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Template selector */}
        <div className="space-y-2">
          {templates.map((t) => (
            <button
              key={t.name}
              onClick={() => setSelected(t)}
              className={` + "`" + `w-full text-left rounded-xl border p-4 transition-colors ${
                selected.name === t.name
                  ? "border-accent bg-accent/5"
                  : "border-border bg-bg-secondary hover:border-accent/30"
              }` + "`" + `}
            >
              <div className="flex items-center gap-3">
                <Mail className={` + "`" + `h-4 w-4 ${selected.name === t.name ? "text-accent" : "text-text-muted"}` + "`" + `} />
                <div>
                  <p className="text-sm font-medium text-foreground">{t.label}</p>
                  <p className="text-xs text-text-muted mt-0.5">{t.description}</p>
                </div>
              </div>
            </button>
          ))}
        </div>

        {/* Preview */}
        <div className="lg:col-span-3">
          <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
            <div className="flex items-center justify-between border-b border-border px-4 py-3">
              <div>
                <p className="text-sm font-medium text-foreground">{selected.label}</p>
                <p className="text-xs text-text-muted">Template: {selected.name}</p>
              </div>
              <span className="rounded-md bg-accent/10 px-2 py-0.5 text-xs font-medium text-accent">
                Preview
              </span>
            </div>
            <div className="p-4">
              <div className="rounded-lg border border-border bg-[#0a0a0f] p-6">
                <div className="max-w-md mx-auto">
                  <div className="rounded-xl border border-[#2a2a3a] bg-[#111118] p-8">
                    <div className="text-center mb-6">
                      <span className="text-2xl font-bold text-[#6c5ce7]">
                        {selected.sampleData.AppName}
                      </span>
                    </div>
                    <h2 className="text-lg font-semibold text-[#e8e8f0] mb-3">
                      {selected.name === "welcome" && ` + "`" + `Welcome, ${selected.sampleData.Name}!` + "`" + `}
                      {selected.name === "password-reset" && "Reset Your Password"}
                      {selected.name === "email-verification" && "Verify Your Email"}
                      {selected.name === "notification" && String((selected.sampleData as Record<string, unknown>).Title ?? "")}
                    </h2>
                    <p className="text-sm text-[#9090a8] mb-4 leading-relaxed">
                      {selected.name === "welcome" && "Thanks for signing up. Your account is ready to use."}
                      {selected.name === "password-reset" && "We received a request to reset your password. Click the button below to set a new one."}
                      {selected.name === "email-verification" && "Please verify your email address by clicking the button below."}
                      {selected.name === "notification" && String((selected.sampleData as Record<string, unknown>).Message ?? "")}
                    </p>
                    <div className="text-center">
                      <span className="inline-block rounded-lg bg-[#6c5ce7] px-6 py-2.5 text-sm font-semibold text-white">
                        {selected.name === "welcome" && "Go to Dashboard"}
                        {selected.name === "password-reset" && "Reset Password"}
                        {selected.name === "email-verification" && "Verify Email"}
                        {selected.name === "notification" && String((selected.sampleData as Record<string, unknown>).ActionText ?? "View")}
                      </span>
                    </div>
                  </div>
                  <p className="text-center text-xs text-[#606078] mt-4">
                    &copy; {selected.sampleData.Year} {selected.sampleData.AppName}. All rights reserved.
                  </p>
                </div>
              </div>
            </div>

            {/* Sample data */}
            <div className="border-t border-border px-4 py-3">
              <p className="text-xs font-medium text-text-muted uppercase mb-2">Template Data</p>
              <pre className="text-xs text-text-secondary font-mono bg-bg-tertiary rounded-lg p-3 overflow-x-auto">
                {JSON.stringify(selected.sampleData, null, 2)}
              </pre>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
`
}

func adminSecurityPage() string {
	return `"use client";

import { useState, useEffect } from "react";
import { Shield, ExternalLink, AlertTriangle } from "@/lib/icons";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export default function SecurityPage() {
  const [loaded, setLoaded] = useState(false);
  const [error, setError] = useState(false);

  useEffect(() => {
    // Check if sentinel is reachable
    fetch(` + "`" + `${API_URL}/sentinel/api/auth/verify` + "`" + `, { method: "GET" })
      .then(() => setLoaded(true))
      .catch(() => setError(true));
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground">Security Dashboard</h1>
          <p className="text-sm text-text-secondary mt-1">
            WAF, rate limiting, threat detection, and anomaly monitoring powered by Sentinel
          </p>
        </div>
        <a
          href={` + "`" + `${API_URL}/sentinel/ui` + "`" + `}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 rounded-lg border border-border bg-bg-secondary px-4 py-2.5 text-sm font-medium text-foreground hover:bg-bg-hover transition-colors"
        >
          Open in new tab
          <ExternalLink className="h-4 w-4" />
        </a>
      </div>

      {error ? (
        <div className="rounded-xl border border-warning/20 bg-warning/5 p-8 text-center">
          <AlertTriangle className="h-10 w-10 text-warning mx-auto mb-4" />
          <h2 className="text-lg font-semibold text-foreground mb-2">Sentinel Not Available</h2>
          <p className="text-sm text-text-secondary max-w-md mx-auto">
            The Sentinel security dashboard could not be reached. Make sure your API server is running
            and <code className="text-xs font-mono bg-bg-secondary px-1.5 py-0.5 rounded">SENTINEL_ENABLED=true</code> is set in your <code className="text-xs font-mono bg-bg-secondary px-1.5 py-0.5 rounded">.env</code> file.
          </p>
        </div>
      ) : (
        <div className="rounded-xl border border-border overflow-hidden" style={{ height: "calc(100vh - 200px)" }}>
          <iframe
            src={` + "`" + `${API_URL}/sentinel/ui` + "`" + `}
            className="w-full h-full border-0"
            title="Sentinel Security Dashboard"
            onLoad={() => setLoaded(true)}
          />
        </div>
      )}
    </div>
  );
}
`
}
