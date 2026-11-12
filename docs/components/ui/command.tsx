'use client'

import * as React from 'react'
import { type DialogProps } from '@radix-ui/react-dialog'
import { Command as CommandPrimitive } from 'cmdk'
import { Search } from 'lucide-react'

import { cn } from '@/lib/utils'
import { Dialog, DialogContent } from '@/components/ui/dialog'

const Command = React.forwardRef<
  React.ElementRef<typeof CommandPrimitive>,
  React.ComponentPropsWithoutRef<typeof CommandPrimitive>
>(({ className, ...props }, ref) => (
  <CommandPrimitive
    ref={ref}
    className={cn(
      'flex h-full w-full flex-col overflow-hidden rounded-md bg-popover text-popover-foreground',
      className,
    )}
    {...props}
  />
))
Command.displayName = CommandPrimitive.displayName

/**
 * Props that belong to the inner Command (filter, shouldFilter, loop, label)
 * are forwarded rather than swallowed. Without this the dialog is stuck with
 * cmdk's default filter, which scores by subsequence: every character of the
 * query only has to appear somewhere in order, so "deplo" matches "errors fix
 * debug issues problems help broken".
 */
type CommandDialogProps = DialogProps &
  Pick<
    React.ComponentPropsWithoutRef<typeof CommandPrimitive>,
    'filter' | 'shouldFilter' | 'loop' | 'label'
  >

const CommandDialog = ({
  children,
  filter,
  shouldFilter,
  loop,
  label,
  ...props
}: CommandDialogProps) => {
  return (
    <Dialog {...props}>
      {/* Anchored near the top rather than centred: the list grows downward as
          you type, and a vertically centred palette moves its own first result
          out from under the cursor every keystroke. */}
      <DialogContent
        hideClose
        className="top-[12vh] max-w-2xl translate-y-0 gap-0 overflow-hidden rounded-xl border-border bg-surface-1 p-0 shadow-2xl"
      >
        <Command
          filter={filter}
          shouldFilter={shouldFilter}
          loop={loop}
          label={label}
          className="bg-transparent [&_[cmdk-group-heading]]:px-3 [&_[cmdk-group-heading]]:text-eyebrow [&_[cmdk-group-heading]]:uppercase [&_[cmdk-group-heading]]:text-fg-subtle [&_[cmdk-group]:not([hidden])_~[cmdk-group]]:pt-0 [&_[cmdk-group]]:px-2 [&_[cmdk-input-wrapper]_svg]:h-4 [&_[cmdk-input-wrapper]_svg]:w-4 [&_[cmdk-item]_svg]:h-4 [&_[cmdk-item]_svg]:w-4"
        >
          {children}
        </Command>
      </DialogContent>
    </Dialog>
  )
}

const CommandInput = React.forwardRef<
  React.ElementRef<typeof CommandPrimitive.Input>,
  React.ComponentPropsWithoutRef<typeof CommandPrimitive.Input>
>(({ className, ...props }, ref) => (
  <div
    className="flex items-center gap-2.5 border-b border-border px-4"
    cmdk-input-wrapper=""
  >
    <Search className="h-4 w-4 shrink-0 text-fg-subtle" />
    <CommandPrimitive.Input
      ref={ref}
      className={cn(
        'flex h-14 w-full bg-transparent text-[15px] text-fg outline-none placeholder:text-fg-subtle disabled:cursor-not-allowed disabled:opacity-50',
        // See the note in search-dialog.tsx: the global focus ring is a
        // floating violet rectangle, and this input is focused for as long as
        // the palette is open. The underline is the indicator instead.
        'focus-visible:shadow-[inset_0_-2px_0_0_var(--accent)]',
        className,
      )}
      {...props}
    />
    <kbd className="hidden shrink-0 select-none rounded-xs border border-border bg-surface-2 px-1.5 py-0.5 font-mono text-[10px] text-fg-subtle sm:inline-block">
      esc
    </kbd>
  </div>
))

CommandInput.displayName = CommandPrimitive.Input.displayName

const CommandList = React.forwardRef<
  React.ElementRef<typeof CommandPrimitive.List>,
  React.ComponentPropsWithoutRef<typeof CommandPrimitive.List>
>(({ className, ...props }, ref) => (
  <CommandPrimitive.List
    ref={ref}
    className={cn(
      'max-h-[min(60vh,28rem)] overflow-y-auto overflow-x-hidden overscroll-contain p-2',
      className,
    )}
    {...props}
  />
))

CommandList.displayName = CommandPrimitive.List.displayName

const CommandEmpty = React.forwardRef<
  React.ElementRef<typeof CommandPrimitive.Empty>,
  React.ComponentPropsWithoutRef<typeof CommandPrimitive.Empty>
>((props, ref) => (
  <CommandPrimitive.Empty
    ref={ref}
    className="py-6 text-center text-sm"
    {...props}
  />
))

CommandEmpty.displayName = CommandPrimitive.Empty.displayName

const CommandGroup = React.forwardRef<
  React.ElementRef<typeof CommandPrimitive.Group>,
  React.ComponentPropsWithoutRef<typeof CommandPrimitive.Group>
>(({ className, ...props }, ref) => (
  <CommandPrimitive.Group
    ref={ref}
    className={cn(
      'overflow-hidden p-1 text-foreground [&_[cmdk-group-heading]]:px-2 [&_[cmdk-group-heading]]:py-1.5 [&_[cmdk-group-heading]]:text-xs [&_[cmdk-group-heading]]:font-medium [&_[cmdk-group-heading]]:text-muted-foreground',
      className,
    )}
    {...props}
  />
))

CommandGroup.displayName = CommandPrimitive.Group.displayName

const CommandSeparator = React.forwardRef<
  React.ElementRef<typeof CommandPrimitive.Separator>,
  React.ComponentPropsWithoutRef<typeof CommandPrimitive.Separator>
>(({ className, ...props }, ref) => (
  <CommandPrimitive.Separator
    ref={ref}
    className={cn('-mx-1 h-px bg-border', className)}
    {...props}
  />
))
CommandSeparator.displayName = CommandPrimitive.Separator.displayName

const CommandItem = React.forwardRef<
  React.ElementRef<typeof CommandPrimitive.Item>,
  React.ComponentPropsWithoutRef<typeof CommandPrimitive.Item>
>(({ className, ...props }, ref) => (
  <CommandPrimitive.Item
    ref={ref}
    className={cn(
      // bg-accent here was the brand violet, and text-accent-foreground is
      // not a utility this config defines, so the selected row was a violet
      // bar with unstyled text on it. The selected state is the accent tint
      // and its hairline instead, which keeps every row's text at full
      // contrast in both themes.
      'relative flex cursor-default select-none items-center gap-2.5 rounded-btn-sm px-3 py-2.5 text-sm text-fg outline-none transition-colors',
      "data-[selected='true']:bg-accent-soft data-[selected=true]:ring-1 data-[selected=true]:ring-inset data-[selected=true]:ring-accent-line",
      'data-[disabled=true]:pointer-events-none data-[disabled=true]:opacity-50',
      '[&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0',
      className,
    )}
    {...props}
  />
))

CommandItem.displayName = CommandPrimitive.Item.displayName

const CommandShortcut = ({
  className,
  ...props
}: React.HTMLAttributes<HTMLSpanElement>) => {
  return (
    <span
      className={cn(
        'ml-auto text-xs tracking-widest text-muted-foreground',
        className,
      )}
      {...props}
    />
  )
}
CommandShortcut.displayName = 'CommandShortcut'

export {
  Command,
  CommandDialog,
  CommandInput,
  CommandList,
  CommandEmpty,
  CommandGroup,
  CommandItem,
  CommandShortcut,
  CommandSeparator,
}
