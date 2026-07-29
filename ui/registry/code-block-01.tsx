"use client";

// code-block-01 — Code display block (Vercel docs / GitHub-inspired)
import { useState } from "react";
export function CodeBlock({
  code = "grit new my-app --template=saas",
  language = "bash",
  filename,
  showLineNumbers = false,
}: {
  code?: string; language?: string; filename?: string; showLineNumbers?: boolean;
}) {
  const [copied, setCopied] = useState(false);
  const copy = () => {
    navigator.clipboard.writeText(code).catch(() => {});
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };
  const lines = code.split("\n");
  return (
    <div className="bg-background border border-border rounded-xl overflow-hidden">
      <div className="flex items-center justify-between px-4 py-2.5 border-b border-border bg-bg-elevated">
        <div className="flex items-center gap-2">
          {filename ? (
            <span className="text-xs text-text-secondary font-mono">{filename}</span>
          ) : (
            <span className="text-xs text-text-muted uppercase tracking-wider">{language}</span>
          )}
        </div>
        <button onClick={copy} className={"flex items-center gap-1.5 text-xs px-2 py-1 rounded transition-colors " + (copied ? "text-success" : "text-text-muted hover:text-foreground")}>
          {copied ? (
            <><svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" /></svg>Copied</>
          ) : (
            <><svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2" /></svg>Copy</>
          )}
        </button>
      </div>
      <pre className="p-4 overflow-x-auto text-sm font-mono">
        {showLineNumbers ? (
          <table className="w-full border-collapse">
            <tbody>
              {lines.map((line, i) => (
                <tr key={i}>
                  <td className="pr-4 text-right text-text-muted select-none w-8">{i + 1}</td>
                  <td className="text-text-secondary">{line}</td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <code className="text-text-secondary">{code}</code>
        )}
      </pre>
    </div>
  );
}

export default CodeBlock;
