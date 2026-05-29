import { useState, useEffect } from 'react'

interface MoneyInputProps {
  value: number | ''
  onChange: (value: number) => void
  placeholder?: string
  className?: string
  min?: number
  required?: boolean
}

function formatWithCommas(val: string): string {
  const num = val.replace(/[^0-9]/g, '')
  if (!num) return ''
  return Number(num).toLocaleString('en-US')
}

function parseFromCommas(val: string): number {
  return Number(val.replace(/[^0-9]/g, '')) || 0
}

export function MoneyInput({ value, onChange, placeholder, className = '', min, required }: MoneyInputProps) {
  const [display, setDisplay] = useState(() =>
    value !== '' && value > 0 ? Number(value).toLocaleString('en-US') : ''
  )

  // Sync external value changes
  useEffect(() => {
    if (value === '' || value === 0) {
      setDisplay('')
    } else {
      const current = parseFromCommas(display)
      if (current !== value) {
        setDisplay(Number(value).toLocaleString('en-US'))
      }
    }
  }, [value])

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const raw = e.target.value
    const formatted = formatWithCommas(raw)
    setDisplay(formatted)
    onChange(parseFromCommas(raw))
  }

  return (
    <input
      type="text"
      inputMode="numeric"
      value={display}
      onChange={handleChange}
      placeholder={placeholder}
      className={className}
      min={min}
      required={required}
    />
  )
}
