package scaffold

// admin_import_modal.go scaffolds the components/tables/import-modal.tsx
// component: the three-stage Excel import flow used by every resource
// list page. Stages: file pick -> validation preview -> submit with
// progress. The whole thing runs client-side via SheetJS + the
// existing apiClient -- no /import endpoint required.

func adminImportModal() string {
	return `"use client";

import { useCallback, useRef, useState } from "react";
import { useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import {
  AlertCircle,
  CheckCircle,
  Download,
  Loader2,
  Upload,
  X,
} from "@/lib/icons";
import type { ResourceDefinition } from "@/lib/resource";
import {
  downloadImportTemplate,
  parseImportFile,
  submitImport,
  type ParsedImport,
  type ParsedImportRow,
} from "@/lib/excel-utils";

interface ImportModalProps {
  resource: ResourceDefinition;
  onClose: () => void;
}

type Stage = "pick" | "preview" | "submitting" | "done";

interface SubmitSummary {
  succeeded: number;
  failed: { rowNumber: number; message: string }[];
}

// ImportModal walks the user from "drop a file" through validation
// to a per-row submit summary. The whole thing runs in the browser:
// SheetJS parses, we POST each valid row to the resource endpoint,
// React Query invalidates so the underlying list refreshes.
export function ImportModal({ resource, onClose }: ImportModalProps) {
  const queryClient = useQueryClient();
  const [stage, setStage] = useState<Stage>("pick");
  const [fileName, setFileName] = useState<string>("");
  const [parsed, setParsed] = useState<ParsedImport | null>(null);
  const [progress, setProgress] = useState({ done: 0, total: 0 });
  const [summary, setSummary] = useState<SubmitSummary | null>(null);
  const [dragOver, setDragOver] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  const importCfg = resource.table.import;
  const allowedFields =
    importCfg && importCfg !== false ? importCfg.fields : undefined;

  const handleFile = useCallback(
    async (file: File) => {
      setFileName(file.name);
      try {
        const result = await parseImportFile(file, resource, allowedFields);
        setParsed(result);
        setStage("preview");
        if (result.rows.length === 0) {
          toast.error("No rows found in file");
        }
      } catch (err) {
        toast.error("Failed to read file: " + (err as Error).message);
        setStage("pick");
      }
    },
    [resource, allowedFields]
  );

  const onFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) handleFile(file);
  };

  const onDrop = (e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    setDragOver(false);
    const file = e.dataTransfer.files?.[0];
    if (file) handleFile(file);
  };

  const handleSubmit = async () => {
    if (!parsed) return;
    setStage("submitting");
    setProgress({ done: 0, total: parsed.rows.filter((r) => r.errors.length === 0).length });
    const result = await submitImport(
      resource.endpoint,
      parsed.rows,
      (done, total) => setProgress({ done, total })
    );
    setSummary(result);
    setStage("done");
    if (result.succeeded > 0) {
      queryClient.invalidateQueries({ queryKey: [resource.endpoint] });
      toast.success("Imported " + result.succeeded + " row" + (result.succeeded === 1 ? "" : "s"));
    }
    if (result.failed.length > 0) {
      toast.error(result.failed.length + " row" + (result.failed.length === 1 ? "" : "s") + " failed");
    }
  };

  const validCount = parsed?.rows.filter((r) => r.errors.length === 0).length ?? 0;
  const invalidCount = parsed ? parsed.rows.length - validCount : 0;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4">
      <div className="w-full max-w-3xl max-h-[85vh] rounded-xl border border-border bg-bg-secondary shadow-2xl flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between border-b border-border px-6 py-4">
          <div>
            <h2 className="text-lg font-semibold text-foreground">
              Import {resource.label?.plural ?? resource.slug}
            </h2>
            <p className="text-xs text-text-muted mt-0.5">
              Upload an .xlsx or .csv file. Validation happens in your browser.
            </p>
          </div>
          <button
            onClick={onClose}
            disabled={stage === "submitting"}
            className="rounded-md p-1.5 text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors disabled:opacity-40"
          >
            <X className="h-4 w-4" />
          </button>
        </div>

        {/* Body */}
        <div className="flex-1 overflow-y-auto p-6">
          {stage === "pick" && (
            <PickStage
              dragOver={dragOver}
              setDragOver={setDragOver}
              onDrop={onDrop}
              onPick={() => inputRef.current?.click()}
              onDownloadTemplate={() => downloadImportTemplate(resource, allowedFields)}
            />
          )}

          {stage === "preview" && parsed && (
            <PreviewStage
              parsed={parsed}
              fileName={fileName}
              validCount={validCount}
              invalidCount={invalidCount}
            />
          )}

          {stage === "submitting" && (
            <SubmittingStage progress={progress} />
          )}

          {stage === "done" && summary && (
            <DoneStage summary={summary} />
          )}

          <input
            ref={inputRef}
            type="file"
            accept=".xlsx,.xls,.csv"
            className="hidden"
            onChange={onFileInput}
          />
        </div>

        {/* Footer */}
        <div className="flex items-center justify-end gap-2 border-t border-border px-6 py-4">
          {stage === "preview" && (
            <>
              <button
                onClick={() => {
                  setStage("pick");
                  setParsed(null);
                  setFileName("");
                }}
                className="rounded-lg border border-border px-4 py-1.5 text-sm text-text-secondary hover:bg-bg-hover transition-colors"
              >
                Back
              </button>
              <button
                onClick={handleSubmit}
                disabled={validCount === 0}
                className="rounded-lg bg-accent px-4 py-1.5 text-sm font-medium text-white hover:bg-accent-hover transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
              >
                Import {validCount} row{validCount === 1 ? "" : "s"}
              </button>
            </>
          )}
          {(stage === "pick" || stage === "done") && (
            <button
              onClick={onClose}
              className="rounded-lg border border-border px-4 py-1.5 text-sm text-text-secondary hover:bg-bg-hover transition-colors"
            >
              {stage === "done" ? "Done" : "Cancel"}
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

interface PickStageProps {
  dragOver: boolean;
  setDragOver: (v: boolean) => void;
  onDrop: (e: React.DragEvent<HTMLDivElement>) => void;
  onPick: () => void;
  onDownloadTemplate: () => void;
}

function PickStage({ dragOver, setDragOver, onDrop, onPick, onDownloadTemplate }: PickStageProps) {
  return (
    <div className="space-y-4">
      <div
        onDragOver={(e) => {
          e.preventDefault();
          setDragOver(true);
        }}
        onDragLeave={() => setDragOver(false)}
        onDrop={onDrop}
        onClick={onPick}
        className={
          "flex flex-col items-center justify-center rounded-lg border-2 border-dashed px-6 py-12 cursor-pointer transition-colors " +
          (dragOver
            ? "border-accent bg-accent/5"
            : "border-border bg-bg-tertiary hover:bg-bg-hover")
        }
      >
        <div className="rounded-full bg-bg-elevated p-3 mb-3">
          <Upload className="h-6 w-6 text-text-secondary" />
        </div>
        <p className="text-sm font-medium text-foreground mb-1">
          Drop a file here, or click to browse
        </p>
        <p className="text-xs text-text-muted">
          Accepts .xlsx, .xls, or .csv
        </p>
      </div>

      <div className="flex items-center justify-between rounded-lg border border-border bg-bg-tertiary px-4 py-3">
        <div>
          <p className="text-sm font-medium text-foreground">Not sure where to start?</p>
          <p className="text-xs text-text-muted">
            Grab a blank template with the right column headers.
          </p>
        </div>
        <button
          onClick={onDownloadTemplate}
          className="flex items-center gap-1.5 rounded-lg border border-border bg-bg-secondary px-3 py-1.5 text-xs text-text-secondary hover:text-foreground hover:bg-bg-hover transition-colors"
        >
          <Download className="h-3.5 w-3.5" />
          Download template
        </button>
      </div>
    </div>
  );
}

function PreviewStage({
  parsed,
  fileName,
  validCount,
  invalidCount,
}: {
  parsed: ParsedImport;
  fileName: string;
  validCount: number;
  invalidCount: number;
}) {
  const showRows = parsed.rows.slice(0, 25);
  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center gap-3 rounded-lg border border-border bg-bg-tertiary px-4 py-3">
        <p className="text-sm font-medium text-foreground flex-1 truncate">
          {fileName}
        </p>
        <span className="inline-flex items-center gap-1 rounded-full bg-success/10 px-2.5 py-0.5 text-xs font-medium text-success">
          <CheckCircle className="h-3 w-3" />
          {validCount} valid
        </span>
        {invalidCount > 0 && (
          <span className="inline-flex items-center gap-1 rounded-full bg-danger/10 px-2.5 py-0.5 text-xs font-medium text-danger">
            <AlertCircle className="h-3 w-3" />
            {invalidCount} invalid
          </span>
        )}
      </div>

      {parsed.unknownHeaders.length > 0 && (
        <div className="rounded-lg border border-warning/40 bg-warning/10 px-4 py-3">
          <p className="text-xs font-semibold text-warning mb-1">
            Unknown columns ignored
          </p>
          <p className="text-xs text-text-secondary">
            {parsed.unknownHeaders.join(", ")}
          </p>
        </div>
      )}

      <div className="rounded-lg border border-border overflow-hidden">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="bg-bg-tertiary">
              <tr>
                <th className="px-3 py-2 text-left text-[11px] font-semibold uppercase tracking-wider text-text-muted">
                  Row
                </th>
                <th className="px-3 py-2 text-left text-[11px] font-semibold uppercase tracking-wider text-text-muted">
                  Status
                </th>
                {parsed.fields.slice(0, 4).map((f) => (
                  <th
                    key={f.key}
                    className="px-3 py-2 text-left text-[11px] font-semibold uppercase tracking-wider text-text-muted"
                  >
                    {f.label}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {showRows.map((row) => (
                <PreviewRow
                  key={row.rowNumber}
                  row={row}
                  fieldKeys={parsed.fields.slice(0, 4).map((f) => f.key)}
                />
              ))}
            </tbody>
          </table>
        </div>
        {parsed.rows.length > showRows.length && (
          <div className="border-t border-border px-3 py-2 text-center text-xs text-text-muted">
            Showing {showRows.length} of {parsed.rows.length} rows. All rows will be imported.
          </div>
        )}
      </div>
    </div>
  );
}

function PreviewRow({ row, fieldKeys }: { row: ParsedImportRow; fieldKeys: string[] }) {
  const hasError = row.errors.length > 0;
  return (
    <tr className={"border-t border-border/50 " + (hasError ? "bg-danger/5" : "")}>
      <td className="px-3 py-2 text-xs text-text-muted font-mono">
        {row.rowNumber}
      </td>
      <td className="px-3 py-2">
        {hasError ? (
          <span title={row.errors.join("; ")} className="inline-flex items-center gap-1 text-xs text-danger">
            <AlertCircle className="h-3 w-3" />
            {row.errors[0]}
          </span>
        ) : (
          <CheckCircle className="h-3.5 w-3.5 text-success" />
        )}
      </td>
      {fieldKeys.map((k) => (
        <td key={k} className="px-3 py-2 text-xs text-foreground truncate max-w-[160px]">
          {row.values[k] === undefined || row.values[k] === ""
            ? "—"
            : String(row.values[k])}
        </td>
      ))}
    </tr>
  );
}

function SubmittingStage({ progress }: { progress: { done: number; total: number } }) {
  const pct = progress.total > 0 ? Math.round((progress.done / progress.total) * 100) : 0;
  return (
    <div className="flex flex-col items-center justify-center py-10 space-y-4">
      <Loader2 className="h-8 w-8 animate-spin text-accent" />
      <p className="text-sm font-medium text-foreground">
        Importing {progress.done} of {progress.total}…
      </p>
      <div className="w-full max-w-sm h-2 rounded-full bg-bg-tertiary overflow-hidden">
        <div
          className="h-full bg-accent transition-all"
          style={{ width: pct + "%" }}
        />
      </div>
    </div>
  );
}

function DoneStage({ summary }: { summary: SubmitSummary }) {
  return (
    <div className="space-y-4">
      <div className="flex items-center gap-3 rounded-lg border border-success/30 bg-success/10 px-4 py-3">
        <CheckCircle className="h-5 w-5 text-success" />
        <p className="text-sm text-foreground">
          <span className="font-semibold">{summary.succeeded}</span> row{summary.succeeded === 1 ? "" : "s"} imported successfully.
        </p>
      </div>
      {summary.failed.length > 0 && (
        <div className="rounded-lg border border-danger/30 bg-danger/5 p-4">
          <p className="text-sm font-semibold text-danger mb-2">
            {summary.failed.length} failed
          </p>
          <ul className="space-y-1 max-h-48 overflow-y-auto">
            {summary.failed.map((f) => (
              <li key={f.rowNumber} className="text-xs text-text-secondary">
                <span className="font-mono text-text-muted">Row {f.rowNumber}:</span> {f.message}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}
`
}
