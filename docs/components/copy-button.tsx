'use client'

import * as React from 'react'
import { Check, Copy } from 'lucide-react'
import { Button, type ButtonProps } from '@/components/ui/button'
import { cn } from '@/lib/utils'

interface CopyButtonProps extends ButtonProps {
  text: string
}

export function CopyButton({ text, className, children, size = 'icon', variant = 'ghost', ...props }: CopyButtonProps) {
  const [copied, setCopied] = React.useState(false)

  const handleCopy = async () => {
    await navigator.clipboard.writeText(text)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  if (children) {
    return (
      <Button
        variant={variant}
        size={size}
        className={className}
        onClick={handleCopy}
        {...props}
      >
        {copied ? (
          <>
            <Check className="h-4 w-4 mr-2" />
            Copied!
          </>
        ) : (
          <>
            <Copy className="h-4 w-4 mr-2" />
            {children}
          </>
        )}
      </Button>
    )
  }

  return (
    <Button
      variant={variant}
      size="icon"
      className={cn('h-8 w-8 shrink-0', className)}
      onClick={handleCopy}
      {...props}
    >
      {copied ? (
        <Check className="h-4 w-4 text-green-500" />
      ) : (
        <Copy className="h-4 w-4" />
      )}
      <span className="sr-only">Copy to clipboard</span>
    </Button>
  )
}
