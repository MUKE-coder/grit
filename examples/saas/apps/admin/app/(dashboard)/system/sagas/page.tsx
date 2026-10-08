"use client";

import { useState } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { PageHeader } from "@/components/chrome/PageHeader";
import { apiClient } from "@/lib/api-client";
import { useToastedMutation } from "@/hooks/use-toasted-mutation";
import {
  RefreshCw, AlertTriangle, CheckCircle, Loader2, ArrowRight, Workflow,
} from "@/lib/icons";

// The five statuses a run can be in. See internal/saga: only "stuck" needs a
// person, which is why it is the one with a button.
type RunStatus = "running" | "compensating" | "done" | "compensated" | "stuck";

interface SagaStep {
  id: string;
  idx: number;
  name: string;
  status: string;
  attempts: number;
  last_error?: string;
  started_at?: string;
  finished_at?: string;
}

interface SagaRun {
  id: string;
  name: string;
  key?: string;
  status: RunStatus;
  cursor: number;
  attempts: number;
  last_error?: string;
  created_at: string;
  finished_at?: string;
  steps?: SagaStep[];
}

interface RunsResponse {
  data: SagaRun[];
  meta: { total: number; page: number; page_size: number; pages: number };
  counts?: Record<string, number>;
}

const STATUS_TONE: Record<string, string> = {
  running: "border-info/30 bg-info/10 text-info",
  compensating: "border-warning/30 bg-warning/10 text-warning",
  done: "border-success/30 bg-success/10 text-success",
  compensated: "border-border bg-bg-elevated text-text-secondary",
  stuck: "border-danger/30 bg-danger/10 text-danger",
  // Step statuses, which share the palette.
  pending: "border-border bg-bg-elevated text-text-muted",
  failed: "border-danger/30 bg-danger/10 text-danger",
};

// What each status means, in the words an operator needs rather than the words
// the code uses. "compensated" reading as a failure is the misunderstanding
// worth preventing: it is the good outcome of a bad run.
const STATUS_HELP: Record<string, string> = {
  running: "Working through the steps.",
  compensating: "A step failed for good and the completed ones are being undone.",
  done: "Every step completed.",
  compensated: "A step failed and everything before it was undone. The world is back where it started.",
  stuck: "A compensation itself failed. Something happened that could not be taken back: this needs you.",
};

function StatusChip({ status, count }: { status: string; count?: number }) {
  return (
    <span
      title={STATUS_HELP[status]}
      className={
        "inline-flex items-center gap-1.5 rounded-full border px-2.5 py-1 text-xs font-semibold " +
        (STATUS_TONE[status] ?? "border-border bg-bg-elevated text-text-secondary")
      }
    >
      {status}
      {count != null && <span className="font-mono opacity-70">{count}</span>}
    </span>
  );
}

export default function SagasPage() {
  const queryClient = useQueryClient();
  const [status, setStatus] = useState<string>("");
  const [openRun, setOpenRun] = useState<string | null>(null);

  const { data, isLoading, refetch, isFetching } = useQuery<RunsResponse>({
    queryKey: ["system-sagas", status],
    queryFn: async () => {
      const { data } = await apiClient.get<RunsResponse>("/api/admin/sagas", {
        params: { page_size: 50, ...(status ? { "filter[status]": status } : {}) },
      });
      return data;
    },
    refetchInterval: 10_000,
  });

  const { data: detail } = useQuery<{ data: SagaRun }>({
    queryKey: ["system-saga", openRun],
    queryFn: async () => {
      const { data } = await apiClient.get<{ data: SagaRun }>("/api/admin/sagas/" + openRun);
      return data;
    },
    enabled: !!openRun,
  });

  const retry = useToastedMutation<unknown, unknown, string>({
    mutationFn: async (id: string) => {
      const { data } = await apiClient.post("/api/admin/sagas/" + id + "/retry");
      return data;
    },
    successMessage: "The run is compensating again",
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["system-sagas"] });
      queryClient.invalidateQueries({ queryKey: ["system-saga"] });
    },
  });

  const counts = data?.counts ?? {};
  const stuck = counts.stuck ?? 0;
  const runs = data?.data ?? [];
  const open = detail?.data;

  return (
    <div className="space-y-6">
      <PageHeader
        title="Sagas"
        subtitle="Multi-step processes that have to either finish or be undone."
        actions={
          <button
            onClick={() => refetch()}
            className="inline-flex items-center gap-2 rounded-lg border border-border px-3 py-2 text-sm text-text-secondary hover:bg-bg-hover"
          >
            <RefreshCw className={"h-4 w-4 " + (isFetching ? "animate-spin" : "")} /> Refresh
          </button>
        }
      />

      {stuck > 0 && (
        <div className="flex items-start gap-2 rounded-lg border border-danger/30 bg-danger/5 px-4 py-3 text-sm text-danger">
          <AlertTriangle className="mt-0.5 h-4 w-4 shrink-0" />
          <span>
            {stuck} run{stuck === 1 ? "" : "s"} could not be undone. Each one did something that was
            not taken back, so there is a charge, a reservation or a booking still standing. Fix
            what the error names, then retry the run.
          </span>
        </div>
      )}

      <div className="flex flex-wrap items-center gap-2">
        <button
          onClick={() => setStatus("")}
          className={
            "rounded-full border px-2.5 py-1 text-xs font-semibold " +
            (status === "" ? "border-accent/40 bg-accent/10 text-accent" : "border-border text-text-secondary hover:bg-bg-hover")
          }
        >
          all {data?.meta?.total ?? 0}
        </button>
        {(["running", "compensating", "stuck", "done", "compensated"] as const).map((s) => (
          <button key={s} onClick={() => setStatus(status === s ? "" : s)} className={status === s ? "ring-2 ring-accent/40 rounded-full" : ""}>
            <StatusChip status={s} count={counts[s] ?? 0} />
          </button>
        ))}
      </div>

      {isLoading ? (
        <div className="flex items-center gap-2 text-sm text-text-muted">
          <Loader2 className="h-4 w-4 animate-spin" /> Reading the runs
        </div>
      ) : runs.length === 0 ? (
        <div className="rounded-xl border border-border bg-bg-secondary p-8 text-center">
          <Workflow className="mx-auto h-8 w-8 text-text-muted" />
          <p className="mt-3 text-sm font-semibold text-foreground">No runs yet</p>
          <p className="mt-1 text-xs text-text-secondary">
            Generate one with <code className="font-mono">grit generate workflow Checkout --steps &quot;charge,ship&quot;</code>,
            then start it with <code className="font-mono">saga.Start</code>.
          </p>
        </div>
      ) : (
        <div className="overflow-hidden rounded-xl border border-border">
          <table className="w-full text-sm">
            <thead className="bg-bg-elevated text-left text-xs uppercase tracking-wider text-text-muted">
              <tr>
                <th className="px-4 py-3 font-medium">Saga</th>
                <th className="px-4 py-3 font-medium">Status</th>
                <th className="px-4 py-3 font-medium">Step</th>
                <th className="px-4 py-3 font-medium">Started</th>
                <th className="px-4 py-3 font-medium"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-border">
              {runs.map((run) => (
                <tr key={run.id} className="hover:bg-bg-hover">
                  <td className="px-4 py-3">
                    <p className="font-medium text-foreground">{run.name}</p>
                    <p className="font-mono text-[11px] text-text-muted">{run.key ?? run.id}</p>
                  </td>
                  <td className="px-4 py-3">
                    <StatusChip status={run.status} />
                    {run.last_error && (
                      <p className="mt-1 max-w-md truncate text-[11px] text-text-muted" title={run.last_error}>
                        {run.last_error}
                      </p>
                    )}
                  </td>
                  <td className="px-4 py-3 font-mono text-xs text-text-secondary">{run.cursor}</td>
                  <td className="px-4 py-3 text-xs text-text-secondary">
                    {new Date(run.created_at).toLocaleString()}
                  </td>
                  <td className="px-4 py-3 text-right">
                    <div className="flex items-center justify-end gap-2">
                      {run.status === "stuck" && (
                        <button
                          onClick={() => retry.mutate(run.id)}
                          disabled={retry.isPending}
                          className="inline-flex items-center gap-1.5 rounded-lg border border-danger/30 bg-danger/10 px-2.5 py-1.5 text-xs font-semibold text-danger hover:bg-danger/20 disabled:opacity-50"
                        >
                          <RefreshCw className="h-3.5 w-3.5" /> Retry
                        </button>
                      )}
                      <button
                        onClick={() => setOpenRun(openRun === run.id ? null : run.id)}
                        className="inline-flex items-center gap-1 text-xs text-text-secondary hover:text-foreground"
                      >
                        Steps <ArrowRight className="h-3.5 w-3.5" />
                      </button>
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {open && (
        <div className="rounded-xl border border-border bg-bg-secondary p-4">
          <div className="mb-3 flex items-center justify-between">
            <div>
              <p className="text-sm font-semibold text-foreground">{open.name}</p>
              <p className="font-mono text-[11px] text-text-muted">{open.id}</p>
            </div>
            <StatusChip status={open.status} />
          </div>
          <ol className="space-y-2">
            {(open.steps ?? []).map((step) => (
              <li key={step.id} className="flex items-start gap-3 rounded-lg border border-border bg-background px-3 py-2">
                <span className="mt-0.5 font-mono text-[11px] text-text-muted">{step.idx}</span>
                <div className="min-w-0 flex-1">
                  <div className="flex items-center gap-2">
                    <span className="text-sm text-foreground">{step.name}</span>
                    <StatusChip status={step.status} />
                    {step.attempts > 1 && (
                      <span className="text-[11px] text-text-muted">{step.attempts} attempts</span>
                    )}
                  </div>
                  {step.last_error && (
                    <p className="mt-1 break-words text-[11px] text-danger">{step.last_error}</p>
                  )}
                </div>
                {step.status === "done" && <CheckCircle className="h-4 w-4 shrink-0 text-success" />}
              </li>
            ))}
          </ol>
          <p className="mt-3 text-[11px] text-text-muted">
            Steps run top to bottom and are undone bottom to top. A step with no Undo is skipped on
            the way back, which is why anything you cannot take back belongs at the end.
          </p>
        </div>
      )}
    </div>
  );
}
