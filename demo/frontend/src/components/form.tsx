import { forwardRef, useEffect, useRef, useState } from 'react'
import { Controller, type Control, type FieldPath, type FieldValues } from 'react-hook-form'
import { cn } from '@/lib/utils'
import { SearchableSelect } from '@/components/ui/SearchableSelect'

/**
 * Form primitives: TextField, TextAreaField, SelectField, FormGrid, FormActions.
 *
 * Designed to compose with react-hook-form via `register('name')`. Spread the
 * register result directly:
 *   <TextField label="Name" {...register('first_name', { required: true })} />
 *
 * Validation errors come from `formState.errors[name]?.message` — pass to
 * `error` prop. Hints appear when there's no error.
 */

function FieldWrap({
  label,
  hint,
  error,
  required,
  className,
  htmlFor,
  children,
}: {
  label?: string
  hint?: string
  error?: string
  required?: boolean
  className?: string
  htmlFor?: string
  children: React.ReactNode
}) {
  return (
    <div className={cn('block space-y-1', className)}>
      {label && (
        <label htmlFor={htmlFor} className="block text-[12.5px] font-medium text-foreground-secondary">
          {label}
          {required && <span className="text-danger ml-0.5">*</span>}
        </label>
      )}
      {children}
      {error ? (
        <span className="block text-[11.5px] text-danger">{error}</span>
      ) : hint ? (
        <span className="block text-[11.5px] text-foreground-muted">{hint}</span>
      ) : null}
    </div>
  )
}

const baseInput =
  'w-full h-10 px-3 rounded-lg border border-border bg-surface text-[13.5px] text-foreground placeholder:text-foreground-muted focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/15 disabled:bg-surface-2 disabled:text-foreground-muted'

interface InputProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'> {
  label?: string
  hint?: string
  error?: string
}

export const TextField = forwardRef<HTMLInputElement, InputProps>(
  ({ label, hint, error, required, className, id, ...props }, ref) => (
    <FieldWrap label={label} hint={hint} error={error} required={required} htmlFor={id}>
      <input
        id={id}
        ref={ref}
        required={required}
        className={cn(baseInput, error && 'border-danger focus:border-danger focus:ring-danger/15', className)}
        {...props}
      />
    </FieldWrap>
  ),
)
TextField.displayName = 'TextField'

interface TextAreaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement> {
  label?: string
  hint?: string
  error?: string
}

export const TextAreaField = forwardRef<HTMLTextAreaElement, TextAreaProps>(
  ({ label, hint, error, required, className, id, ...props }, ref) => (
    <FieldWrap label={label} hint={hint} error={error} required={required} htmlFor={id}>
      <textarea
        id={id}
        ref={ref}
        required={required}
        className={cn(baseInput, 'h-auto min-h-[90px] py-2 resize-y', error && 'border-danger', className)}
        {...props}
      />
    </FieldWrap>
  ),
)
TextAreaField.displayName = 'TextAreaField'

interface SelectOpt {
  value: string
  label: string
}
interface SelectProps extends React.SelectHTMLAttributes<HTMLSelectElement> {
  label?: string
  hint?: string
  error?: string
  options: SelectOpt[]
  placeholder?: string
}

/**
 * SelectField now renders the SearchableSelect combobox under the hood while
 * keeping the same register-friendly API (spread `{...register('field')}`).
 *
 * Mechanism: a hidden native `<select>` carries the ref + name so RHF can
 * register it; the visible UI is the SearchableSelect, and picking an option
 * updates the hidden select via the native value setter + dispatches a real
 * 'change' event so RHF + native form validation both fire correctly.
 */
export const SelectField = forwardRef<HTMLSelectElement, SelectProps>(
  ({ label, hint, error, required, className, options, placeholder, id, name, onChange, onBlur, defaultValue, value, disabled, ...rest }, ref) => {
    const hiddenRef = useRef<HTMLSelectElement | null>(null)
    const [internal, setInternal] = useState<string>(
      (defaultValue as string | undefined) ?? (value as string | undefined) ?? ''
    )

    // Forward the hidden ref to the consumer's ref (RHF's register).
    const setRefs = (el: HTMLSelectElement | null) => {
      hiddenRef.current = el
      if (typeof ref === 'function') ref(el)
      else if (ref) (ref as React.MutableRefObject<HTMLSelectElement | null>).current = el
    }

    // If the consumer drives `value` (controlled usage), keep our display in sync.
    useEffect(() => {
      if (value !== undefined && String(value) !== internal) setInternal(String(value))
    }, [value])

    // RHF's setValue / form reset writes to the hidden input directly. Mirror it
    // back to our combobox state so the UI stays correct after a reset.
    useEffect(() => {
      const el = hiddenRef.current
      if (!el) return
      const observer = new MutationObserver(() => {
        if (el.value !== internal) setInternal(el.value)
      })
      observer.observe(el, { attributes: true, attributeFilter: ['value'] })
      return () => observer.disconnect()
    }, [internal])

    const pick = (next: string) => {
      setInternal(next)
      const el = hiddenRef.current
      if (el) {
        const setter = Object.getOwnPropertyDescriptor(window.HTMLSelectElement.prototype, 'value')?.set
        if (setter) setter.call(el, next)
        el.dispatchEvent(new Event('change', { bubbles: true }))
      }
    }

    return (
      <FieldWrap label={label} hint={hint} error={error} required={required} htmlFor={id}>
        <SearchableSelect
          options={options}
          value={internal}
          onChange={pick}
          placeholder={placeholder}
          required={required}
          error={error}
          disabled={disabled}
        />
        {/* Hidden native select for RHF register + form-validation. */}
        <select
          id={id}
          ref={setRefs}
          name={name}
          required={required}
          disabled={disabled}
          onChange={onChange}
          onBlur={onBlur}
          tabIndex={-1}
          aria-hidden="true"
          style={{ position: 'absolute', opacity: 0, height: 0, width: 0, pointerEvents: 'none' }}
          {...rest}
        >
          <option value=""></option>
          {options.map((o) => (
            <option key={o.value} value={o.value}>{o.label}</option>
          ))}
        </select>
      </FieldWrap>
    )
  },
)
SelectField.displayName = 'SelectField'

export function FormGrid({
  columns = 2,
  children,
}: {
  columns?: 1 | 2 | 3
  children: React.ReactNode
}) {
  const cls =
    columns === 1
      ? 'grid grid-cols-1 gap-4'
      : columns === 3
      ? 'grid grid-cols-1 sm:grid-cols-3 gap-4'
      : 'grid grid-cols-1 sm:grid-cols-2 gap-4'
  return <div className={cls}>{children}</div>
}

/**
 * SearchableSelectField — react-hook-form Controller wrapper around the
 * SearchableSelect combobox. Use this instead of SelectField for any picker
 * with more than ~5 options (branches, motorcycles, drivers, borrowers).
 *
 *   <SearchableSelectField
 *     control={control}
 *     name="branch_id"
 *     label="Branch"
 *     required
 *     options={branches.map(b => ({ value: String(b.id), label: b.name }))}
 *   />
 */
interface SSOpt {
  value: string
  label: string
  subtitle?: string
}
interface SearchableSelectFieldProps<T extends FieldValues> {
  control: Control<T>
  name: FieldPath<T>
  label?: string
  hint?: string
  placeholder?: string
  required?: boolean
  options: SSOpt[]
  /** When true, the form value is stored as a number (Number(value)). Defaults to true if name ends in "_id" or "Id". */
  asNumber?: boolean
}
export function SearchableSelectField<T extends FieldValues>({
  control, name, label, hint, placeholder, required, options, asNumber,
}: SearchableSelectFieldProps<T>) {
  const inferNumber = asNumber ?? (name.endsWith('_id') || name.endsWith('Id'))
  return (
    <Controller
      control={control}
      name={name}
      rules={required ? { validate: (v) => v !== undefined && v !== null && v !== '' || 'Required' } : undefined}
      render={({ field, fieldState }) => (
        <SearchableSelect
          label={label}
          hint={hint}
          placeholder={placeholder}
          required={required}
          error={fieldState.error?.message}
          options={options}
          value={field.value === undefined || field.value === null ? '' : String(field.value)}
          onChange={(v) => field.onChange(inferNumber ? (v === '' ? undefined : Number(v)) : v)}
        />
      )}
    />
  )
}

/** Sticky footer with Save (primary) + Cancel. Place at the end of a form. */
export function FormActions({
  onCancel,
  submitLabel = 'Save',
  isPending,
  extra,
}: {
  onCancel?: () => void
  submitLabel?: string
  isPending?: boolean
  extra?: React.ReactNode
}) {
  return (
    <div className="flex items-center justify-between pt-4 border-t border-border-subtle">
      <div>{extra}</div>
      <div className="flex gap-2">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            className="h-9 px-3.5 rounded-lg border border-border bg-surface text-[13px] font-medium text-foreground-secondary hover:bg-surface-hover transition"
          >
            Cancel
          </button>
        )}
        <button
          type="submit"
          disabled={isPending}
          className="h-9 px-3.5 rounded-lg bg-accent text-white text-[13px] font-medium hover:bg-accent-hover disabled:opacity-60 transition"
        >
          {isPending ? 'Saving…' : submitLabel}
        </button>
      </div>
    </div>
  )
}

/** Section header inside a form (visual divider with optional description). */
export function FormSection({
  title,
  description,
  children,
}: {
  title: string
  description?: string
  children: React.ReactNode
}) {
  return (
    <div className="space-y-3">
      <div>
        <h3 className="text-[12.5px] font-semibold uppercase tracking-wider text-foreground-muted">
          {title}
        </h3>
        {description && (
          <p className="text-[12px] text-foreground-muted mt-0.5">{description}</p>
        )}
      </div>
      {children}
    </div>
  )
}
