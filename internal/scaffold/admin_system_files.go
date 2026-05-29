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
  id: string;
  filename: string;
  original_name: string;
  mime_type: string;
  size: number;
  path: string;
  url: string;
  thumbnail_url?: string;
  user_id: string;
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
    mutationFn: async (id: string) => {
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

import { useEffect, useState } from "react";
import { Shield, ExternalLink, AlertTriangle, Zap, Activity, AlertCircle, Globe } from "@/lib/icons";
import { api } from "@/lib/api";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

interface ScoreCard {
  score?: number
  grade?: string
  factors?: { name: string; score: number }[]
}

interface Summary {
  summary?: { threats_total?: number; threats_24h?: number; blocked?: number; actors?: number }
  score?: ScoreCard
  threats?: { data?: Array<{ id: string; type: string; severity: string; cvss?: number; source_ip?: string; route?: string; created_at?: string }> }
  auth_shield?: { enabled?: boolean; failures?: number; locked?: Array<{ ip: string; until: string }> }
  csp_top?: { data?: Array<{ violated_directive: string; blocked_uri: string; count: number }> }
  performance?: { p50_ms?: number; p95_ms?: number; p99_ms?: number; rps?: number; pipeline_drops?: number }
  trends?: unknown
  top_routes?: { data?: Array<{ route: string; count: number }> }
  _errors?: Record<string, string>
}

export default function SecurityPage() {
  const [data, setData] = useState<Summary | null>(null);
  const [err, setErr] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let alive = true;
    const fetch = async () => {
      try {
        const res = await api.get("/api/admin/security/summary");
        if (alive) setData(res.data?.data || res.data);
        setErr(null);
      } catch (e: any) {
        if (alive) setErr(e?.response?.data?.error?.message || "Failed to load");
      } finally {
        if (alive) setLoading(false);
      }
    };
    fetch();
    const t = setInterval(fetch, 20_000); // refresh every 20s
    return () => { alive = false; clearInterval(t); };
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground flex items-center gap-2.5">
            <Shield className="h-6 w-6 text-accent" /> Security
          </h1>
          <p className="text-sm text-text-secondary mt-1">
            Live summary from Sentinel — score, threats, AuthShield, CSP, performance
          </p>
        </div>
        <a
          href={` + "`${API_URL}/sentinel/ui`" + `}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 rounded-lg border border-border bg-bg-secondary px-4 py-2.5 text-sm font-medium text-foreground hover:bg-bg-hover transition-colors"
        >
          Open full dashboard <ExternalLink className="h-4 w-4" />
        </a>
      </div>

      {err && !loading && (
        <div className="rounded-xl border border-warning/30 bg-warning/5 p-4 text-sm text-warning flex items-start gap-2">
          <AlertTriangle className="h-4 w-4 mt-0.5 shrink-0" />
          <div>
            <strong>Couldn&apos;t reach Sentinel.</strong> Make sure <code className="px-1 py-0.5 rounded bg-bg-secondary text-xs font-mono">SENTINEL_ENABLED=true</code> and the API is running.
            <div className="text-xs text-text-muted mt-1">{err}</div>
          </div>
        </div>
      )}

      {/* Top-line scorecard */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        <KpiCard title="Security score" value={data?.score?.score ?? "—"} sub={data?.score?.grade} tone="success" icon={Shield} />
        <KpiCard title="Threats (24h)" value={data?.summary?.threats_24h ?? "—"} sub={` + "`${data?.summary?.threats_total ?? 0} total`" + `} tone="warning" icon={AlertCircle} />
        <KpiCard title="Blocked" value={data?.summary?.blocked ?? "—"} sub="auto by WAF" tone="info" icon={Zap} />
        <KpiCard title="Actors" value={data?.summary?.actors ?? "—"} sub="unique sources" tone="default" icon={Globe} />
      </div>

      {/* Recent threats */}
      <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
        <div className="flex items-center justify-between px-4 py-3 border-b border-border">
          <h2 className="text-sm font-semibold text-foreground">Recent threats</h2>
          <a href={` + "`${API_URL}/sentinel/ui/threats`" + `} target="_blank" rel="noopener noreferrer" className="text-xs text-accent hover:underline">View all</a>
        </div>
        <div className="divide-y divide-border">
          {(data?.threats?.data ?? []).slice(0, 8).map((t) => (
            <div key={t.id} className="flex items-center gap-3 px-4 py-2.5 text-sm">
              <SeverityChip severity={t.severity} cvss={t.cvss} />
              <span className="text-foreground font-mono text-xs truncate flex-1">{t.type}</span>
              <span className="text-text-muted text-xs hidden md:inline">{t.route}</span>
              <span className="text-text-muted text-xs font-mono">{t.source_ip}</span>
            </div>
          ))}
          {(!data?.threats?.data || data.threats.data.length === 0) && (
            <div className="px-4 py-6 text-center text-sm text-text-muted">No threats in window.</div>
          )}
        </div>
      </div>

      {/* AuthShield + CSP side-by-side */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Panel title="AuthShield">
          {data?.auth_shield?.enabled === false ? (
            <p className="text-sm text-text-muted">AuthShield is disabled in this environment.</p>
          ) : (
            <div className="space-y-3">
              <div className="grid grid-cols-2 gap-2 text-sm">
                <Stat label="Failures (rolling)" value={data?.auth_shield?.failures ?? 0} />
                <Stat label="Locked IPs" value={data?.auth_shield?.locked?.length ?? 0} />
              </div>
              <ul className="space-y-1.5 text-xs font-mono">
                {(data?.auth_shield?.locked ?? []).slice(0, 5).map((l) => (
                  <li key={l.ip} className="flex justify-between gap-2 text-text-secondary">
                    <span>{l.ip}</span>
                    <span className="text-text-muted">until {l.until}</span>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </Panel>
        <Panel title="Top CSP violations">
          <ul className="space-y-1.5 text-xs font-mono">
            {(data?.csp_top?.data ?? []).slice(0, 5).map((v, i) => (
              <li key={i} className="flex justify-between gap-2">
                <span className="truncate text-text-secondary">{v.violated_directive}</span>
                <span className="text-text-muted shrink-0">{v.count}×</span>
              </li>
            ))}
            {(!data?.csp_top?.data || data.csp_top.data.length === 0) && (
              <li className="text-text-muted">No violations.</li>
            )}
          </ul>
        </Panel>
      </div>
    </div>
  );
}

function KpiCard({ title, value, sub, tone, icon: Icon }: { title: string; value: any; sub?: any; tone: "default" | "success" | "warning" | "info" | "danger"; icon: any }) {
  const toneCls = { default: "text-foreground", success: "text-success", warning: "text-warning", info: "text-info", danger: "text-danger" }[tone];
  return (
    <div className="rounded-xl border border-border bg-bg-secondary p-4">
      <div className="flex items-center justify-between mb-1">
        <span className="text-[10px] font-mono uppercase tracking-wider text-text-muted">{title}</span>
        <Icon className={` + "`h-4 w-4 ${toneCls}`" + `} />
      </div>
      <p className={` + "`text-2xl font-semibold tabular-nums ${toneCls}`" + `}>{value}</p>
      {sub && <p className="text-[11px] text-text-muted mt-0.5">{sub}</p>}
    </div>
  );
}

function Panel({ title, children }: { title: string; children: any }) {
  return (
    <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
      <div className="px-4 py-3 border-b border-border"><h2 className="text-sm font-semibold text-foreground">{title}</h2></div>
      <div className="p-4">{children}</div>
    </div>
  );
}

function Stat({ label, value }: { label: string; value: any }) {
  return (
    <div className="rounded-lg border border-border bg-bg-elevated px-3 py-2">
      <p className="text-[10px] uppercase tracking-wider text-text-muted">{label}</p>
      <p className="text-lg font-semibold text-foreground tabular-nums">{value}</p>
    </div>
  );
}

function SeverityChip({ severity, cvss }: { severity?: string; cvss?: number }) {
  const tone =
    severity === "critical" ? "bg-danger/15 text-danger border-danger/30" :
    severity === "high"     ? "bg-warning/15 text-warning border-warning/30" :
    severity === "medium"   ? "bg-info/15 text-info border-info/30" :
                              "bg-bg-elevated text-text-muted border-border";
  return (
    <span className={` + "`inline-flex items-center gap-1 px-2 py-0.5 rounded border text-[10px] font-mono ${tone}`" + `}>
      {severity ?? "info"}
      {cvss !== undefined && cvss > 0 && <span>· {cvss.toFixed(1)}</span>}
    </span>
  );
}
`
}

func adminObservabilityPage() string {
	return `"use client";

import { useEffect, useState } from "react";
import { Activity, ExternalLink, AlertTriangle, Cpu, Database, Zap } from "@/lib/icons";
import { api } from "@/lib/api";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

interface Summary {
  overview?: { p50_ms?: number; p95_ms?: number; p99_ms?: number; rps?: number; error_rate?: number; total_requests?: number }
  slos?: { data?: Array<{ name: string; target: number; current: number; budget_remaining: number; status: string }> }
  use?: { resources?: Array<{ name: string; utilization: { value: number; band: string }; saturation: { value: number; band: string }; errors: { value: number; band: string } }> }
  n1_ranked?: { data?: Array<{ route: string; pattern: string; occurrences: number; avg_queries_per_request: number; impact_score: number }> }
  errors?: { data?: Array<{ id: string; type: string; message: string; route: string; count: number }> }
  runtime?: { heap_alloc_mb?: number; goroutines?: number; gc_pause_ms?: number }
  health_checks?: { data?: Array<{ name: string; status: string; latency_ms: number; error?: string }> }
  alerts?: { data?: Array<{ id: string; name: string; severity: string; message: string }> }
  _errors?: Record<string, string>
}

const BAND_CLS: Record<string, string> = {
  green:   "bg-success/15 text-success border-success/30",
  amber:   "bg-warning/15 text-warning border-warning/30",
  red:     "bg-danger/15 text-danger border-danger/30",
  unknown: "bg-bg-elevated text-text-muted border-border",
};

export default function ObservabilityPage() {
  const [data, setData] = useState<Summary | null>(null);
  const [err, setErr] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    let alive = true;
    const fetch = async () => {
      try {
        const res = await api.get("/api/admin/observability/summary");
        if (alive) setData(res.data?.data || res.data);
        setErr(null);
      } catch (e: any) {
        if (alive) setErr(e?.response?.data?.error?.message || "Failed to load");
      } finally {
        if (alive) setLoading(false);
      }
    };
    fetch();
    const t = setInterval(fetch, 10_000);
    return () => { alive = false; clearInterval(t); };
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-foreground flex items-center gap-2.5">
            <Activity className="h-6 w-6 text-accent" /> Observability
          </h1>
          <p className="text-sm text-text-secondary mt-1">
            Live summary from Pulse — percentile latency, SLOs, USE grid, top N+1, errors, runtime
          </p>
        </div>
        <a
          href={` + "`${API_URL}/pulse/ui`" + `}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 rounded-lg border border-border bg-bg-secondary px-4 py-2.5 text-sm font-medium text-foreground hover:bg-bg-hover transition-colors"
        >
          Open full dashboard <ExternalLink className="h-4 w-4" />
        </a>
      </div>

      {err && !loading && (
        <div className="rounded-xl border border-warning/30 bg-warning/5 p-4 text-sm text-warning flex items-start gap-2">
          <AlertTriangle className="h-4 w-4 mt-0.5 shrink-0" />
          <div>
            <strong>Couldn&apos;t reach Pulse.</strong> Make sure <code className="px-1 py-0.5 rounded bg-bg-secondary text-xs font-mono">PULSE_ENABLED=true</code> and the API is running.
            <div className="text-xs text-text-muted mt-1">{err}</div>
          </div>
        </div>
      )}

      {/* Top-line KPIs */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
        <Kpi label="p95 latency" value={data?.overview?.p95_ms != null ? ` + "`${data.overview.p95_ms.toFixed(0)} ms`" + ` : "—"} tone="info" icon={Zap} />
        <Kpi label="p99 latency" value={data?.overview?.p99_ms != null ? ` + "`${data.overview.p99_ms.toFixed(0)} ms`" + ` : "—"} tone="warning" icon={Zap} />
        <Kpi label="Error rate" value={data?.overview?.error_rate != null ? ` + "`${(data.overview.error_rate * 100).toFixed(2)}%`" + ` : "—"} tone="danger" icon={AlertTriangle} />
        <Kpi label="RPS" value={data?.overview?.rps?.toFixed?.(1) ?? "—"} tone="success" icon={Activity} />
      </div>

      {/* SLO compliance bars */}
      <Panel title="SLOs">
        {(data?.slos?.data ?? []).length === 0 ? (
          <p className="text-sm text-text-muted">No SLOs configured. Add to <code className="px-1 py-0.5 rounded bg-bg-elevated text-xs font-mono">pulse.Config.SLOs</code>.</p>
        ) : (
          <div className="space-y-3">
            {(data?.slos?.data ?? []).map((s) => (
              <div key={s.name} className="space-y-1">
                <div className="flex items-center justify-between text-xs">
                  <span className="font-medium text-foreground">{s.name}</span>
                  <span className={` + "`font-mono ${s.status === \"firing\" ? \"text-danger\" : \"text-success\"}`" + `}>{(s.current * 100).toFixed(2)}% / {(s.target * 100).toFixed(2)}%</span>
                </div>
                <div className="h-2 rounded-full bg-bg-elevated overflow-hidden">
                  <div className={` + "`h-full ${s.status === \"firing\" ? \"bg-danger\" : \"bg-success\"}`" + `} style={{ width: ` + "`${Math.min(s.current * 100, 100)}%`" + ` }} />
                </div>
                <p className="text-[10px] text-text-muted">Budget remaining: {(s.budget_remaining * 100).toFixed(1)}%</p>
              </div>
            ))}
          </div>
        )}
      </Panel>

      {/* USE grid + N+1 side-by-side */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Panel title="USE method">
          <table className="w-full text-xs">
            <thead>
              <tr className="text-text-muted">
                <th className="text-left font-medium py-1">Resource</th>
                <th className="text-center font-medium">U</th>
                <th className="text-center font-medium">S</th>
                <th className="text-center font-medium">E</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {(data?.use?.resources ?? []).map((r) => (
                <tr key={r.name}>
                  <td className="py-1.5 font-medium text-foreground capitalize">{r.name}</td>
                  <td className="py-1.5 text-center"><Cell band={r.utilization?.band} value={r.utilization?.value} /></td>
                  <td className="py-1.5 text-center"><Cell band={r.saturation?.band} value={r.saturation?.value} /></td>
                  <td className="py-1.5 text-center"><Cell band={r.errors?.band} value={r.errors?.value} /></td>
                </tr>
              ))}
              {(!data?.use?.resources || data.use.resources.length === 0) && (
                <tr><td colSpan={4} className="text-center text-text-muted py-3">No USE samples yet.</td></tr>
              )}
            </tbody>
          </table>
        </Panel>

        <Panel title="Top N+1 by impact">
          <ul className="space-y-2 text-xs">
            {(data?.n1_ranked?.data ?? []).slice(0, 6).map((n, i) => (
              <li key={i} className="rounded-lg border border-border bg-bg-elevated p-2.5">
                <div className="flex items-center justify-between mb-0.5">
                  <span className="font-mono text-foreground truncate">{n.route}</span>
                  <span className="text-text-muted shrink-0 ml-2">impact {n.impact_score.toFixed(0)}</span>
                </div>
                <p className="font-mono text-[10px] text-text-muted truncate">{n.pattern}</p>
                <p className="text-[10px] text-text-muted mt-0.5">{n.occurrences} occurrences · ~{n.avg_queries_per_request} queries/req</p>
              </li>
            ))}
            {(!data?.n1_ranked?.data || data.n1_ranked.data.length === 0) && (
              <li className="text-text-muted text-center py-3">No N+1 detections.</li>
            )}
          </ul>
        </Panel>
      </div>

      {/* Errors + Runtime side-by-side */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <Panel title="Recent unresolved errors">
          <ul className="space-y-1.5 text-xs">
            {(data?.errors?.data ?? []).slice(0, 6).map((e) => (
              <li key={e.id} className="flex items-start justify-between gap-2 py-1 border-b border-border last:border-0">
                <div className="min-w-0 flex-1">
                  <p className="text-foreground truncate"><span className="font-mono text-[10px] text-text-muted">{e.type}</span> {e.message}</p>
                  <p className="font-mono text-[10px] text-text-muted truncate">{e.route}</p>
                </div>
                <span className="text-text-muted shrink-0">{e.count}×</span>
              </li>
            ))}
            {(!data?.errors?.data || data.errors.data.length === 0) && <li className="text-text-muted text-center py-3">No unresolved errors.</li>}
          </ul>
        </Panel>

        <Panel title="Go runtime">
          <div className="grid grid-cols-3 gap-2 text-xs">
            <Stat label="Heap alloc" value={data?.runtime?.heap_alloc_mb != null ? ` + "`${data.runtime.heap_alloc_mb.toFixed(0)} MB`" + ` : "—"} />
            <Stat label="Goroutines" value={data?.runtime?.goroutines ?? "—"} />
            <Stat label="GC pause" value={data?.runtime?.gc_pause_ms != null ? ` + "`${data.runtime.gc_pause_ms.toFixed(1)} ms`" + ` : "—"} />
          </div>
          {data?.health_checks?.data && data.health_checks.data.length > 0 && (
            <div className="mt-4 pt-3 border-t border-border space-y-1">
              <p className="text-[10px] uppercase font-mono tracking-wider text-text-muted">Health checks</p>
              {data.health_checks.data.map((h) => (
                <div key={h.name} className="flex items-center justify-between text-xs">
                  <span className="font-mono text-foreground">{h.name}</span>
                  <span className={` + "`text-[10px] ${h.status === \"healthy\" ? \"text-success\" : \"text-danger\"}`" + `}>
                    {h.status} · {h.latency_ms}ms
                  </span>
                </div>
              ))}
            </div>
          )}
        </Panel>
      </div>
    </div>
  );
}

function Kpi({ label, value, tone, icon: Icon }: { label: string; value: any; tone: "default" | "success" | "warning" | "info" | "danger"; icon: any }) {
  const toneCls = { default: "text-foreground", success: "text-success", warning: "text-warning", info: "text-info", danger: "text-danger" }[tone];
  return (
    <div className="rounded-xl border border-border bg-bg-secondary p-4">
      <div className="flex items-center justify-between mb-1">
        <span className="text-[10px] font-mono uppercase tracking-wider text-text-muted">{label}</span>
        <Icon className={` + "`h-4 w-4 ${toneCls}`" + `} />
      </div>
      <p className={` + "`text-2xl font-semibold tabular-nums ${toneCls}`" + `}>{value}</p>
    </div>
  );
}

function Panel({ title, children }: { title: string; children: any }) {
  return (
    <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
      <div className="px-4 py-3 border-b border-border"><h2 className="text-sm font-semibold text-foreground">{title}</h2></div>
      <div className="p-4">{children}</div>
    </div>
  );
}

function Stat({ label, value }: { label: string; value: any }) {
  return (
    <div className="rounded-lg border border-border bg-bg-elevated px-3 py-2">
      <p className="text-[10px] uppercase tracking-wider text-text-muted">{label}</p>
      <p className="text-base font-semibold text-foreground tabular-nums">{value}</p>
    </div>
  );
}

function Cell({ band, value }: { band?: string; value?: number }) {
  const cls = BAND_CLS[band || "unknown"];
  return (
    <span className={` + "`inline-flex items-center justify-center min-w-[44px] px-2 py-0.5 rounded border text-[10px] font-mono ${cls}`" + `}>
      {value != null ? value.toFixed(0) : "—"}
    </span>
  );
}
`
}
