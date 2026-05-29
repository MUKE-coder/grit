import { useState, useRef, useEffect } from 'react'
import { Controller, type Control, type FieldPath, type FieldValues } from 'react-hook-form'
import { cn } from '@/lib/utils'

/**
 * CurrencyInput / CurrencyField
 *
 * Auto-formats numbers with thousands separators while typing.
 *   3000     -> 3,000
 *   1500000  -> 1,500,000
 *   1500.5   -> 1,500.5
 *
 * The bound form value is a `number | undefined` — never the formatted string.
 * Empty input becomes `undefined` so optional money fields stay optional.
 */

function formatWithCommas(raw: string): string {
  if (!raw) return ''
  const negative = raw.startsWith('-')
  const cleaned = raw.replace(/[^\d.]/g, '')
  if (!cleaned) return ''
  const parts = cleaned.split('.')
  const intPart = parts[0]?.replace(/\B(?=(\d{3})+(?!\d))/g, ',') ?? ''
  const decimalPart = parts.length > 1 ? '.' + parts.slice(1).join('') : ''
  return (negative ? '-' : '') + intPart + decimalPart
}

function parseToNumber(formatted: string): number | undefined {
  if (!formatted) return undefined
  const cleaned = formatted.replace(/,/g, '')
  const n = parseFloat(cleaned)
  return Number.isFinite(n) ? n : undefined
}

interface CurrencyInputProps {
  value?: number
  onChange: (value: number | undefined) => void
  onBlur?: () => void
  placeholder?: string
  disabled?: boolean
  prefix?: string
  className?: string
  id?: string
}

export function CurrencyInput({
  value,
  onChange,
  onBlur,
  placeholder = '0',
  disabled,
  prefix = 'UGX',
  className,
  id,
}: CurrencyInputProps) {
  const [display, setDisplay] = useState<string>(
    value === undefined || value === null ? '' : formatWithCommas(String(value)),
  )
  const lastValueRef = useRef<number | undefined>(value)

  // Keep display in sync if external value changes (form reset, prefilled).
  useEffect(() => {
    if (value !== lastValueRef.current) {
      setDisplay(value === undefined || value === null ? '' : formatWithCommas(String(value)))
      lastValueRef.current = value
    }
  }, [value])

  return (
    <div className={cn('relative', className)}>
      {prefix && (
        <span className="absolute left-3 top-1/2 -translate-y-1/2 text-[12.5px] font-medium text-foreground-muted pointer-events-none">
          {prefix}
        </span>
      )}
      <input
        id={id}
        type="text"
        inputMode="decimal"
        value={display}
        disabled={disabled}
        onBlur={onBlur}
        placeholder={placeholder}
        onChange={(e) => {
          const raw = e.target.value
          const formatted = formatWithCommas(raw)
          setDisplay(formatted)
          const n = parseToNumber(formatted)
          lastValueRef.current = n
          onChange(n)
        }}
        className={cn(
          'w-full h-10 rounded-lg border border-border bg-surface text-[13.5px] text-foreground placeholder:text-foreground-muted',
          'focus:border-accent focus:outline-none focus:ring-2 focus:ring-accent/15',
          'disabled:bg-surface-2 disabled:text-foreground-muted',
          prefix ? 'pl-12 pr-3' : 'px-3',
          'tabular-nums',
        )}
      />
    </div>
  )
}

interface CurrencyFieldProps<T extends FieldValues> {
  control: Control<T>
  name: FieldPath<T>
  label?: string
  hint?: string
  required?: boolean
  prefix?: string
  placeholder?: string
  disabled?: boolean
}

/** react-hook-form Controller wrapper. Use when you have a `control` from `useForm()`. */
export function CurrencyField<T extends FieldValues>({
  control,
  name,
  label,
  hint,
  required,
  prefix,
  placeholder,
  disabled,
}: CurrencyFieldProps<T>) {
  return (
    <Controller
      control={control}
      name={name}
      render={({ field, fieldState }) => (
        <label className="block space-y-1">
          {label && (
            <span className="block text-[12.5px] font-medium text-foreground-secondary">
              {label}
              {required && <span className="text-danger ml-0.5">*</span>}
            </span>
          )}
          <CurrencyInput
            value={field.value as number | undefined}
            onChange={field.onChange}
            onBlur={field.onBlur}
            prefix={prefix}
            placeholder={placeholder}
            disabled={disabled}
          />
          {fieldState.error?.message ? (
            <span className="block text-[11.5px] text-danger">{fieldState.error.message}</span>
          ) : hint ? (
            <span className="block text-[11.5px] text-foreground-muted">{hint}</span>
          ) : null}
        </label>
      )}
    />
  )
}
