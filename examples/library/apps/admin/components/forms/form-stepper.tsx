"use client";

import { useMemo, useState } from "react";
import { useForm, useWatch, type UseFormTrigger, type FieldValues } from "react-hook-form";
import type { FieldDefinition, FormDefinition } from "@/lib/resource";
import { writableValues } from "@/lib/form-values";
import { FieldRenderer, buildDefaults } from "./form-builder";
import { Button, buttonClasses } from "@/components/ui/button";
import { Check, ChevronLeft, ChevronRight, Loader2 } from "@/lib/icons";

/**
 * Has this value changed from what is stored?
 *
 * Deliberately loose, because "unchanged" has several spellings by the time a
 * value has been through a form control:
 *
 *   - A column that was never set arrives as null, undefined or "" depending on
 *     the field type. Treating those as different from each other lights up the
 *     Update button on a step nobody touched.
 *   - A number input hands back the string "5" for a stored number 5.
 *   - Relationship arrays and line-items are objects, so identity comparison
 *     always reports a change.
 *
 * Getting this wrong in the permissive direction means a permanently-enabled
 * Update button, which trains people to ignore it.
 */
function sameValue(a: unknown, b: unknown): boolean {
  if (a === b) return true;
  const emptyA = a === "" || a === null || a === undefined;
  const emptyB = b === "" || b === null || b === undefined;
  if (emptyA && emptyB) return true;
  if (emptyA !== emptyB) return false;
  if (typeof a === "object" || typeof b === "object") {
    try {
      return JSON.stringify(a) === JSON.stringify(b);
    } catch {
      return false;
    }
  }
  return String(a) === String(b);
}

interface ComputedStep {
  title: string;
  description?: string;
  fields: FieldDefinition[];
}

function computeSteps(form: FormDefinition): ComputedStep[] {
  if (form.steps && form.steps.length > 0) {
    return form.steps.map((step) => ({
      title: step.title,
      description: step.description,
      fields: step.fields
        .map((key) => form.fields.find((f) => f.key === key))
        .filter(Boolean) as FieldDefinition[],
    }));
  }
  const perStep = form.fieldsPerStep ?? 4;
  const chunks: FieldDefinition[][] = [];
  for (let i = 0; i < form.fields.length; i += perStep) {
    chunks.push(form.fields.slice(i, i + perStep));
  }
  return chunks.map((fields, i) => ({
    title: `Step ${i + 1}`,
    fields,
  }));
}

interface FormStepperProps {
  form: FormDefinition;
  defaultValues?: Record<string, unknown>;
  onSubmit: (data: Record<string, unknown>) => void;
  onCancel: () => void;
  isSubmitting?: boolean;
  submitLabel?: string;
  /**
   * EDIT MODE ONLY. Persist just the fields belonging to one step.
   *
   * When supplied, every step gets its own Update button and the whole-form
   * submit disappears — the point being that editing the address on step 3 must
   * not rewrite the twenty fields on steps 1 and 2 with whatever the form
   * happens to be holding.
   *
   * Must reject on failure. The stepper only clears the step's dirty state when
   * this resolves, so a silent catch here would show a saved step that was not.
   */
  onStepSave?: (values: Record<string, unknown>) => Promise<unknown>;
  /** Label for the last step's close button when onStepSave is in play. */
  doneLabel?: string;
}

export function FormStepper({
  form: formDef,
  defaultValues = {},
  onSubmit,
  onCancel,
  isSubmitting,
  submitLabel = "Save",
  onStepSave,
  doneLabel = "Done",
}: FormStepperProps) {
  const [currentStep, setCurrentStep] = useState(0);
  const steps = computeSteps(formDef);
  const isVertical = formDef.stepVariant === "vertical";
  const isTwoColumn = formDef.layout === "two-column";
  const isLastStep = currentStep === steps.length - 1;

  // defaultValues is a fresh object on every render of the parent; keying off
  // its identity would rebuild the baseline constantly and every step would
  // read as clean forever.
  // biome-ignore lint/correctness/useExhaustiveDependencies: keyed on the serialised defaults, see above.
  const initialValues = useMemo(
    () => buildDefaults(formDef.fields, defaultValues),
    [formDef.fields, JSON.stringify(defaultValues ?? {})]
  );

  const {
    control,
    handleSubmit,
    trigger,
    getValues,
    formState: { errors },
  } = useForm({ defaultValues: initialValues });

  // What is currently persisted, per field. Advanced only when a step save
  // actually succeeds, so a failed request leaves the step dirty and the button
  // enabled to retry.
  const [baseline, setBaseline] = useState<Record<string, unknown>>(initialValues);
  const [savingStep, setSavingStep] = useState<number | null>(null);

  const stepKeys = useMemo(
    () => (steps[currentStep]?.fields ?? []).map((f) => f.key),
    [steps, currentStep]
  );

  // Subscribing to just this step's fields, rather than watch()-ing everything,
  // so typing in step 1 does not re-render step 4.
  const watched = useWatch({ control, name: stepKeys as string[] });
  const stepDirty =
    !!onStepSave &&
    stepKeys.some((key, i) => !sameValue((watched as unknown[])?.[i], baseline[key]));

  const handleStepSave = async () => {
    if (!onStepSave) return;
    const valid = await trigger(stepKeys);
    if (!valid) return;

    const all = getValues() as Record<string, unknown>;
    const patch: Record<string, unknown> = {};
    for (const key of stepKeys) patch[key] = all[key];

    setSavingStep(currentStep);
    try {
      await onStepSave(writableValues(formDef.fields, patch));
      setBaseline((prev) => ({ ...prev, ...patch }));
    } catch {
      // The mutation already surfaced the error. Leave the baseline alone so
      // the step stays dirty and the button stays live for another attempt.
    } finally {
      setSavingStep(null);
    }
  };

  const handleNext = async () => {
    const fieldKeys = steps[currentStep].fields.map((f) => f.key);
    const valid = await trigger(fieldKeys);
    if (valid) setCurrentStep((s) => s + 1);
  };

  const handlePrev = () => setCurrentStep((s) => Math.max(0, s - 1));

  const handleFinalSubmit = handleSubmit((data) => onSubmit(writableValues(formDef.fields, data)));

  return (
    <div className={isVertical ? "flex gap-8" : "space-y-6"}>
      {/* Step Indicator */}
      {isVertical ? (
        <VerticalIndicator steps={steps} current={currentStep} onStepClick={setCurrentStep} trigger={trigger} />
      ) : (
        <HorizontalIndicator steps={steps} current={currentStep} onStepClick={setCurrentStep} trigger={trigger} />
      )}

      {/* Step Content */}
      <div className={isVertical ? "flex-1 min-w-0" : ""}>
        <div className="overflow-hidden">
          {steps.map((step, idx) => (
            <div
              key={idx}
              className={idx === currentStep ? "block" : "hidden"}
            >
              {step.description && (
                <p className="text-sm text-text-secondary mb-4">{step.description}</p>
              )}
              <div className={`grid gap-4 ${isTwoColumn ? "grid-cols-1 sm:grid-cols-2" : "grid-cols-1"}`}>
                {step.fields.map((field) => (
                  <div
                    key={field.key}
                    className={field.colSpan === 2 && isTwoColumn ? "sm:col-span-2" : ""}
                  >
                    <FieldRenderer field={field} control={control} errors={errors} getValues={getValues} />
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>

        {/* Progress Bar */}
        <div className="mt-6 mb-4">
          <div className="h-1 w-full rounded-full bg-border">
            <div
              className="h-1 rounded-full bg-accent transition-all duration-300"
              style={{ width: `${((currentStep + 1) / steps.length) * 100}%` }}
            />
          </div>
          <p className="text-xs text-text-muted mt-1.5">
            Step {currentStep + 1} of {steps.length}
            {stepDirty && (
              <span className="text-warning ml-2">Unsaved changes on this step</span>
            )}
          </p>
        </div>

        {/* Navigation */}
        <div className="flex items-center justify-between pt-4 border-t border-border">
          <div>
            {currentStep > 0 ? (
              <Button variant="outline" onClick={handlePrev}>
                <ChevronLeft className="h-4 w-4" />
                Previous
              </Button>
            ) : (
              <Button variant="outline" onClick={onCancel}>
                Cancel
              </Button>
            )}
          </div>
          {/* Editing an existing record saves step by step; creating a new one
              still submits once at the end, because a half-created record has
              nothing to PATCH against. */}
          {onStepSave ? (
            <div className="flex items-center gap-2">
              <Button
                onClick={handleStepSave}
                disabled={!stepDirty || savingStep !== null}
                loading={savingStep === currentStep}
                title={stepDirty ? undefined : "No changes on this step"}
              >
                Update
              </Button>
              {isLastStep ? (
                <button
                  type="button"
                  onClick={onCancel}
                  className="rounded-lg border border-border px-5 py-2 text-sm font-medium text-text-secondary hover:bg-bg-hover transition-colors"
                >
                  {doneLabel}
                </button>
              ) : (
                <button
                  type="button"
                  onClick={handleNext}
                  className="flex items-center gap-1.5 rounded-lg border border-border px-5 py-2 text-sm font-medium text-foreground hover:bg-bg-hover transition-colors"
                >
                  Next
                  <ChevronRight className="h-4 w-4" />
                </button>
              )}
            </div>
          ) : (
            <div>
              {isLastStep ? (
                <button
                  type="button"
                  onClick={handleFinalSubmit}
                  disabled={isSubmitting}
                  className={buttonClasses()}
                >
                  {isSubmitting && <Loader2 className="h-4 w-4 animate-spin" />}
                  {submitLabel}
                </button>
              ) : (
                <button
                  type="button"
                  onClick={handleNext}
                  className={buttonClasses()}
                >
                  Next
                  <ChevronRight className="h-4 w-4" />
                </button>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

/* ─── Horizontal Step Indicator ─────────────────────────────────── */

function HorizontalIndicator({
  steps,
  current,
  onStepClick,
  trigger,
}: {
  steps: ComputedStep[];
  current: number;
  onStepClick: (i: number) => void;
  trigger: UseFormTrigger<FieldValues>;
}) {
  const handleClick = async (idx: number) => {
    if (idx < current) {
      onStepClick(idx);
    } else if (idx === current + 1) {
      const fieldKeys = steps[current].fields.map((f) => f.key);
      const valid = await trigger(fieldKeys);
      if (valid) onStepClick(idx);
    }
  };

  return (
    <div className="flex items-center justify-center">
      {steps.map((step, idx) => {
        const state = idx < current ? "completed" : idx === current ? "active" : "upcoming";
        return (
          <div key={idx} className="flex items-center">
            <button
              type="button"
              onClick={() => handleClick(idx)}
              className="flex flex-col items-center gap-1.5 group"
            >
              <div
                className={`
                  flex h-9 w-9 items-center justify-center rounded-full text-sm font-semibold transition-all
                  ${state === "completed"
                    ? "bg-success text-white"
                    : state === "active"
                    ? "bg-accent text-accent-fg ring-4 ring-accent/20"
                    : "bg-bg-hover text-text-muted border border-border group-hover:border-border/80"}
                `}
              >
                {state === "completed" ? <Check className="h-4 w-4" /> : idx + 1}
              </div>
              <span
                className={`
                  text-xs whitespace-nowrap transition-colors
                  ${state === "active" ? "text-foreground font-medium" : state === "completed" ? "text-foreground" : "text-text-muted"}
                `}
              >
                {step.title}
              </span>
            </button>
            {idx < steps.length - 1 && (
              <div
                className={`
                  h-0.5 w-12 mx-2 rounded-full transition-colors
                  ${idx < current ? "bg-success" : "bg-border"}
                `}
              />
            )}
          </div>
        );
      })}
    </div>
  );
}

/* ─── Vertical Step Indicator ───────────────────────────────────── */

function VerticalIndicator({
  steps,
  current,
  onStepClick,
  trigger,
}: {
  steps: ComputedStep[];
  current: number;
  onStepClick: (i: number) => void;
  trigger: UseFormTrigger<FieldValues>;
}) {
  const handleClick = async (idx: number) => {
    if (idx < current) {
      onStepClick(idx);
    } else if (idx === current + 1) {
      const fieldKeys = steps[current].fields.map((f) => f.key);
      const valid = await trigger(fieldKeys);
      if (valid) onStepClick(idx);
    }
  };

  return (
    <div className="w-52 shrink-0">
      <nav className="space-y-1">
        {steps.map((step, idx) => {
          const state = idx < current ? "completed" : idx === current ? "active" : "upcoming";
          return (
            <div key={idx}>
              <button
                type="button"
                onClick={() => handleClick(idx)}
                className={`
                  flex items-center gap-3 w-full rounded-lg px-3 py-2.5 text-left transition-colors
                  ${state === "active" ? "bg-accent/10" : "hover:bg-bg-hover"}
                `}
              >
                <div
                  className={`
                    flex h-7 w-7 shrink-0 items-center justify-center rounded-full text-xs font-semibold transition-all
                    ${state === "completed"
                      ? "bg-success text-white"
                      : state === "active"
                      ? "bg-accent text-accent-fg ring-2 ring-accent/20"
                      : "bg-bg-hover text-text-muted border border-border"}
                  `}
                >
                  {state === "completed" ? <Check className="h-3 w-3" /> : idx + 1}
                </div>
                <div className="min-w-0">
                  <p
                    className={`
                      text-sm truncate
                      ${state === "active" ? "text-foreground font-medium" : state === "completed" ? "text-foreground" : "text-text-muted"}
                    `}
                  >
                    {step.title}
                  </p>
                  {step.description && (
                    <p className="text-xs text-text-muted truncate">{step.description}</p>
                  )}
                </div>
              </button>
              {idx < steps.length - 1 && (
                <div className="ml-6 py-1">
                  <div
                    className={`
                      w-0.5 h-4 rounded-full transition-colors
                      ${idx < current ? "bg-success" : "bg-border"}
                    `}
                  />
                </div>
              )}
            </div>
          );
        })}
      </nav>
    </div>
  );
}
