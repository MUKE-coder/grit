"use client";

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
