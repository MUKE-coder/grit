import { useState, useRef, useEffect } from 'react'
import { Search, ChevronDown, Check } from 'lucide-react'
import { cn } from '@/lib/utils'

interface Option {
  value: string
  label: string
  subtitle?: string
}

interface SearchableSelectProps {
  options: Option[]
  value: string
  onChange: (value: string) => void
  placeholder?: string
  label?: string
  required?: boolean
  error?: string
  hint?: string
  disabled?: boolean
}

/**
 * Searchable dropdown — type-to-filter combobox. Used as the default for any
 * picker with more than ~5 options. Visual style matches the new design tokens
 * (h-10 rounded-lg, blue accent, subtle border).
 */
export function SearchableSelect({
  options,
  value,
  onChange,
  placeholder = 'Select...',
  label,
  required,
  error,
  hint,
  disabled,
}: SearchableSelectProps) {
  const [open, setOpen] = useState(false)
  const [search, setSearch] = useState('')
  const ref = useRef<HTMLDivElement>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  const filtered = options.filter(
    (o) =>
      o.label.toLowerCase().includes(search.toLowerCase()) ||
      o.subtitle?.toLowerCase().includes(search.toLowerCase()),
  )

  const selected = options.find((o) => o.value === value)

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false)
        setSearch('')
      }
    }
    document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [])

  useEffect(() => {
    if (open) inputRef.current?.focus()
  }, [open])

  return (
    <div ref={ref} className="relative space-y-1">
      {label && (
        <label className="block text-[12.5px] font-medium text-foreground-secondary">
          {label}
          {required && <span className="text-danger ml-0.5">*</span>}
        </label>
      )}

      <button
        type="button"
        disabled={disabled}
        onClick={() => !disabled && setOpen(!open)}
        className={cn(
          'w-full flex items-center justify-between px-3 h-10 rounded-lg border bg-surface text-[13.5px] transition disabled:bg-surface-2 disabled:text-foreground-muted',
          open
            ? 'border-accent ring-2 ring-accent/15'
            : error
            ? 'border-danger'
            : 'border-border hover:border-foreground-muted',
        )}
      >
        <span className={selected ? 'text-foreground' : 'text-foreground-muted'}>
          {selected ? selected.label : placeholder}
        </span>
        <ChevronDown size={14} className={cn('text-foreground-muted transition', open && 'rotate-180')} />
      </button>

      {/* Hidden input for native form validation */}
      {required && (
        <input
          tabIndex={-1}
          autoComplete="off"
          style={{ position: 'absolute', opacity: 0, height: 0, width: 0 }}
          value={value || ''}
          required={required}
          onChange={() => {}}
        />
      )}

      {error ? (
        <span className="block text-[11.5px] text-danger">{error}</span>
      ) : hint ? (
        <span className="block text-[11.5px] text-foreground-muted">{hint}</span>
      ) : null}

      {open && (
        <div className="absolute z-50 mt-1 w-full bg-surface rounded-lg border border-border shadow-md overflow-hidden">
          <div className="p-1.5 border-b border-border-subtle">
            <div className="relative">
              <Search size={13} className="absolute left-2.5 top-1/2 -translate-y-1/2 text-foreground-muted" />
              <input
                ref={inputRef}
                type="text"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="Search..."
                className="w-full pl-8 pr-3 h-8 rounded-md bg-surface-2 text-[12.5px] placeholder:text-foreground-muted focus:outline-none"
              />
            </div>
          </div>

          <div className="max-h-56 overflow-y-auto py-1">
            {filtered.length === 0 ? (
              <p className="px-3 py-3 text-[12.5px] text-foreground-muted text-center">No results.</p>
            ) : (
              filtered.map((option) => (
                <button
                  key={option.value}
                  type="button"
                  onClick={() => {
                    onChange(option.value)
                    setOpen(false)
                    setSearch('')
                  }}
                  className={cn(
                    'w-full text-left px-3 py-2 text-[13px] flex items-center justify-between hover:bg-surface-hover transition',
                    option.value === value && 'bg-accent-light',
                  )}
                >
                  <div className="min-w-0">
                    <p className={cn('font-medium truncate', option.value === value ? 'text-accent' : 'text-foreground')}>
                      {option.label}
                    </p>
                    {option.subtitle && (
                      <p className="text-[11.5px] text-foreground-muted truncate mt-0.5">{option.subtitle}</p>
                    )}
                  </div>
                  {option.value === value && <Check size={14} className="text-accent shrink-0 ml-2" />}
                </button>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  )
}
