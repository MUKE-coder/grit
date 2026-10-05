"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Mail, AlertTriangle } from "@/lib/icons";
import { apiClient } from "@/lib/api-client";

interface MailTemplate {
  name: string;
  description: string;
  subject: string;
  has_text: boolean;
}

interface MailTemplatesResponse {
  data: MailTemplate[];
  driver: string;
}

type Part = "html" | "text";

export default function MailPage() {
  const [selected, setSelected] = useState<string | null>(null);
  const [part, setPart] = useState<Part>("html");

  const templates = useQuery<MailTemplatesResponse>({
    queryKey: ["mail-templates"],
    queryFn: async () => {
      const { data } = await apiClient.get<MailTemplatesResponse>("/api/admin/mail/templates");
      return data;
    },
  });

  const list = templates.data?.data ?? [];
  const current = list.find((t) => t.name === selected) ?? list[0];
  const shownPart: Part = part === "text" && current?.has_text ? "text" : "html";

  // The HTML is rendered by the API with sample data and shown in a sandboxed
  // iframe, so the page shows exactly what the mailer sends.
  const preview = useQuery<string>({
    queryKey: ["mail-preview", current?.name ?? "", shownPart],
    enabled: Boolean(current),
    queryFn: async () => {
      const name = current ? current.name : "";
      const { data } = await apiClient.get<string>(
        "/api/admin/mail/preview/" + encodeURIComponent(name) + (shownPart === "text" ? "?part=text" : ""),
        { responseType: "text", transformResponse: [(body: string) => body] },
      );
      return data;
    },
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-foreground">Email Templates</h1>
        <p className="text-sm text-text-secondary mt-1">
          Rendered by the API with sample data.{" "}
          {templates.data && (templates.data.driver ? "Mail is sent with " + templates.data.driver + "." : "No mail driver is configured.")}
        </p>
      </div>

      {templates.isError && (
        <div className="rounded-xl border border-warning/30 bg-warning/5 p-4 text-sm text-warning flex items-start gap-2">
          <AlertTriangle className="h-4 w-4 mt-0.5 shrink-0" />
          <span>The templates could not be loaded from the API.</span>
        </div>
      )}

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        <div className="space-y-2">
          {list.map((t) => {
            const active = current?.name === t.name;
            return (
              <button
                key={t.name}
                type="button"
                onClick={() => setSelected(t.name)}
                aria-pressed={active}
                className={
                  "w-full text-left rounded-xl border p-4 transition-colors " +
                  (active ? "border-accent bg-accent/5" : "border-border bg-bg-secondary hover:border-accent/30")
                }
              >
                <div className="flex items-center gap-3">
                  <Mail className={"h-4 w-4 " + (active ? "text-accent" : "text-text-muted")} />
                  <div className="min-w-0">
                    <p className="text-sm font-medium text-foreground truncate">{t.name}</p>
                    <p className="text-xs text-text-muted mt-0.5">{t.description}</p>
                  </div>
                </div>
              </button>
            );
          })}
        </div>

        <div className="lg:col-span-3">
          {current && (
            <div className="rounded-xl border border-border bg-bg-secondary overflow-hidden">
              <div className="flex items-center justify-between gap-3 border-b border-border px-4 py-3">
                <div className="min-w-0">
                  <p className="text-sm font-medium text-foreground truncate">{current.subject || current.name}</p>
                  <p className="text-xs text-text-muted">Template: {current.name}</p>
                </div>
                <div className="flex items-center gap-1 rounded-lg border border-border p-0.5">
                  {(["html", "text"] as Part[]).map((p) => (
                    <button
                      key={p}
                      type="button"
                      disabled={p === "text" && !current.has_text}
                      onClick={() => setPart(p)}
                      aria-pressed={shownPart === p}
                      className={
                        "rounded-md px-2.5 py-1 text-xs font-medium disabled:opacity-40 " +
                        (shownPart === p ? "bg-accent/10 text-accent" : "text-text-secondary hover:text-foreground")
                      }
                    >
                      {p === "html" ? "HTML" : "Text"}
                    </button>
                  ))}
                </div>
              </div>
              <div className="p-4">
                {preview.isError ? (
                  <p className="text-sm text-warning">The preview could not be rendered.</p>
                ) : shownPart === "text" ? (
                  <pre className="text-xs text-text-secondary font-mono bg-bg-tertiary rounded-lg p-4 whitespace-pre-wrap">
                    {preview.data ?? ""}
                  </pre>
                ) : (
                  <iframe
                    title={"Preview of " + current.name}
                    sandbox=""
                    srcDoc={preview.data ?? ""}
                    className="w-full h-[640px] rounded-lg border border-border bg-white"
                  />
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
