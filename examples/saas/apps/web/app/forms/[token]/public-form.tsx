"use client";

import { useState } from "react";
import { api } from "@/lib/api";

export type PublicFieldType =
  | "text"
  | "email"
  | "tel"
  | "textarea"
  | "number"
  | "checkbox"
  | "date"
  | "datetime"
  | "file";

export interface PublicField {
  key: string;
  label: string;
  type: PublicFieldType;
  required: boolean;
}

export interface ShareInfo {
  resource_name: string;
  has_password: boolean;
  label: string;
  custom_title: string;
  custom_description: string;
  fields: PublicField[];
}

interface PublicFormProps {
  token: string;
  info: ShareInfo;
}

export function PublicForm({ token, info }: PublicFormProps) {
  const [password, setPassword] = useState("");
  // Mixed value types so checkbox + number fields can survive the
  // round-trip without coercion ceremony at submit time.
  const [fields, setFields] = useState<Record<string, string | number | boolean>>(() => {
    const initial: Record<string, string | number | boolean> = {};
    for (const f of info.fields) {
      initial[f.key] = f.type === "checkbox" ? false : "";
    }
    return initial;
  });
  const [submitting, setSubmitting] = useState(false);
  const [done, setDone] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const update = (key: string, value: string | number | boolean) => {
    setFields((prev) => ({ ...prev, [key]: value }));
  };

  const onSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      // Strip file fields — not supported on public shares yet
      // (auth-gated /api/uploads endpoint). Sending them would
      // confuse the dispatcher's typed unmarshal.
      const payload: Record<string, string | number | boolean> = {};
      for (const f of info.fields) {
        if (f.type === "file") continue;
        payload[f.key] = fields[f.key];
      }
      await api.post("/api/public/forms/" + token + "/submit", {
        _password: password,
        fields: payload,
      });
      setDone(true);
    } catch (err) {
      const e = err as { response?: { data?: { error?: { message?: string } } } };
      setError(e?.response?.data?.error?.message || "Submission failed");
    } finally {
      setSubmitting(false);
    }
  };

  if (done) {
    return (
      <div className="mt-6 rounded-xl border border-emerald-200 bg-emerald-50 p-4 text-sm text-emerald-700">
        <p className="font-medium">Thank you</p>
        <p className="mt-1">Your submission was received.</p>
      </div>
    );
  }

  return (
    <form onSubmit={onSubmit} className="mt-6 space-y-4">
      {info.has_password && (
        <Field
          field={{ key: "_password", label: "Password", type: "text", required: true }}
          inputType="password"
          value={password}
          onChange={(v) => setPassword(String(v))}
          hint="This form is password-protected — ask whoever shared the link."
        />
      )}

      {info.fields.length === 0 && (
        <div className="rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-sm text-amber-700">
          This resource has no public-form fields registered on the API.
          Ask the operator to re-generate the resource so the form-share
          dispatcher picks up the field schema.
        </div>
      )}

      {info.fields.map((f) => (
        <Field
          key={f.key}
          field={f}
          value={fields[f.key] ?? (f.type === "checkbox" ? false : "")}
          onChange={(v) => update(f.key, v)}
        />
      ))}

      {error && (
        <div className="rounded-lg border border-red-200 bg-red-50 px-3 py-2 text-sm text-red-700">
          {error}
        </div>
      )}

      <button
        type="submit"
        disabled={submitting}
        className="w-full rounded-lg bg-slate-900 px-4 py-2.5 text-sm font-semibold text-white hover:bg-slate-800 disabled:opacity-50"
      >
        {submitting ? "Sending…" : "Submit"}
      </button>
    </form>
  );
}

interface FieldProps {
  field: PublicField;
  value: string | number | boolean;
  onChange: (v: string | number | boolean) => void;
  inputType?: string;
  hint?: string;
}

function Field({ field, value, onChange, inputType, hint }: FieldProps) {
  const labelClass = "block text-sm font-medium text-slate-700";
  const inputClass =
    "block w-full rounded-lg border border-slate-300 bg-white px-3 py-2 text-sm text-slate-900 placeholder-slate-400 focus:border-slate-900 focus:outline-none focus:ring-2 focus:ring-slate-900/10";
  const fieldId = "public-field-" + field.key;

  if (field.type === "file") {
    return (
      <div className="space-y-1.5">
        <label className={labelClass} htmlFor={fieldId}>{field.label}</label>
        <div id={fieldId} className="rounded-lg border border-dashed border-slate-300 bg-slate-50 px-3 py-3 text-xs text-slate-500">
          File uploads aren&apos;t supported on public-share forms.
          {field.required
            ? " The operator must collect this file through a different channel."
            : " You can leave this blank."}
        </div>
      </div>
    );
  }

  if (field.type === "textarea") {
    return (
      <div className="space-y-1.5">
        <label className={labelClass} htmlFor={fieldId}>
          {field.label}
          {field.required && <span className="ml-1 text-red-500">*</span>}
        </label>
        <textarea
          id={fieldId}
          value={String(value)}
          onChange={(e) => onChange(e.target.value)}
          required={field.required}
          rows={4}
          className={inputClass}
        />
        {hint && <p className="text-xs text-slate-500">{hint}</p>}
      </div>
    );
  }

  if (field.type === "checkbox") {
    return (
      <div className="space-y-1.5">
        <label className="flex items-center gap-2 text-sm font-medium text-slate-700" htmlFor={fieldId}>
          <input
            id={fieldId}
            type="checkbox"
            checked={Boolean(value)}
            onChange={(e) => onChange(e.target.checked)}
            className="h-4 w-4 rounded border-slate-300"
          />
          {field.label}
          {field.required && <span className="ml-1 text-red-500">*</span>}
        </label>
        {hint && <p className="text-xs text-slate-500">{hint}</p>}
      </div>
    );
  }

  if (field.type === "number") {
    return (
      <div className="space-y-1.5">
        <label className={labelClass} htmlFor={fieldId}>
          {field.label}
          {field.required && <span className="ml-1 text-red-500">*</span>}
        </label>
        <input
          id={fieldId}
          type="number"
          value={value === "" ? "" : String(value)}
          onChange={(e) => onChange(e.target.value === "" ? "" : Number(e.target.value))}
          required={field.required}
          className={inputClass}
        />
        {hint && <p className="text-xs text-slate-500">{hint}</p>}
      </div>
    );
  }

  const htmlType =
    inputType ??
    (field.type === "email"
      ? "email"
      : field.type === "tel"
        ? "tel"
        : field.type === "date"
          ? "date"
          : field.type === "datetime"
            ? "datetime-local"
            : "text");

  return (
    <div className="space-y-1.5">
      <label className={labelClass} htmlFor={fieldId}>
        {field.label}
        {field.required && <span className="ml-1 text-red-500">*</span>}
      </label>
      <input
        id={fieldId}
        type={htmlType}
        value={String(value)}
        onChange={(e) => onChange(e.target.value)}
        required={field.required}
        className={inputClass}
      />
      {hint && <p className="text-xs text-slate-500">{hint}</p>}
    </div>
  );
}
