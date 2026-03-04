package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeDesktopFrontendUIFiles(root string, opts DesktopOptions) error {
	files := map[string]string{
		filepath.Join(root, "frontend", "src", "components", "ui", "button.tsx"):        desktopUIButton(),
		filepath.Join(root, "frontend", "src", "components", "ui", "input.tsx"):         desktopUIInput(),
		filepath.Join(root, "frontend", "src", "components", "ui", "label.tsx"):         desktopUILabel(),
		filepath.Join(root, "frontend", "src", "components", "ui", "card.tsx"):          desktopUICard(),
		filepath.Join(root, "frontend", "src", "components", "ui", "badge.tsx"):         desktopUIBadge(),
		filepath.Join(root, "frontend", "src", "components", "ui", "dialog.tsx"):        desktopUIDialog(),
		filepath.Join(root, "frontend", "src", "components", "ui", "dropdown-menu.tsx"): desktopUIDropdownMenu(),
		filepath.Join(root, "frontend", "src", "components", "ui", "select.tsx"):        desktopUISelect(),
		filepath.Join(root, "frontend", "src", "components", "ui", "tabs.tsx"):          desktopUITabs(),
		filepath.Join(root, "frontend", "src", "components", "ui", "avatar.tsx"):        desktopUIAvatar(),
		filepath.Join(root, "frontend", "src", "components", "ui", "toast.tsx"):         desktopUIToast(),
		filepath.Join(root, "frontend", "src", "components", "ui", "scroll-area.tsx"):   desktopUIScrollArea(),
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}

	return nil
}

func desktopUIButton() string {
	return `import { ButtonHTMLAttributes, ReactNode } from "react";
import { cn } from "@/lib/utils";

type Variant = "default" | "secondary" | "outline" | "ghost" | "destructive";
type Size = "sm" | "default" | "lg";

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: Variant;
  size?: Size;
  children: ReactNode;
}

const variantClasses: Record<Variant, string> = {
  default: "bg-accent text-white hover:bg-accent-hover",
  secondary: "bg-bg-elevated text-foreground hover:bg-bg-hover",
  outline: "border border-border bg-transparent text-foreground hover:bg-bg-elevated",
  ghost: "bg-transparent text-foreground hover:bg-bg-elevated",
  destructive: "bg-danger text-white hover:bg-danger/80",
};

const sizeClasses: Record<Size, string> = {
  sm: "h-8 px-3 text-xs",
  default: "h-10 px-4 text-sm",
  lg: "h-12 px-6 text-base",
};

export default function Button({
  variant = "default",
  size = "default",
  className,
  children,
  disabled,
  ...props
}: ButtonProps) {
  return (
    <button
      className={cn(
        "inline-flex items-center justify-center rounded-lg font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-accent/50 disabled:opacity-50 disabled:pointer-events-none",
        variantClasses[variant],
        sizeClasses[size],
        className
      )}
      disabled={disabled}
      {...props}
    >
      {children}
    </button>
  );
}
`
}

func desktopUIInput() string {
	return `import { forwardRef, InputHTMLAttributes } from "react";
import { cn } from "@/lib/utils";

const Input = forwardRef<HTMLInputElement, InputHTMLAttributes<HTMLInputElement>>(
  ({ className, ...props }, ref) => {
    return (
      <input
        ref={ref}
        className={cn(
          "w-full bg-bg-elevated border border-border rounded-lg px-3 py-2 text-sm text-foreground placeholder:text-text-muted focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors disabled:opacity-50 disabled:cursor-not-allowed",
          className
        )}
        {...props}
      />
    );
  }
);

Input.displayName = "Input";

export default Input;
`
}

func desktopUILabel() string {
	return `import { LabelHTMLAttributes } from "react";
import { cn } from "@/lib/utils";

interface LabelProps extends LabelHTMLAttributes<HTMLLabelElement> {}

export default function Label({ className, children, ...props }: LabelProps) {
  return (
    <label
      className={cn("text-sm font-medium text-foreground", className)}
      {...props}
    >
      {children}
    </label>
  );
}
`
}

func desktopUICard() string {
	return `import { HTMLAttributes } from "react";
import { cn } from "@/lib/utils";

export function Card({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn("rounded-xl border border-border bg-bg-elevated", className)}
      {...props}
    />
  );
}

export function CardHeader({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn("flex flex-col space-y-1.5 p-6", className)}
      {...props}
    />
  );
}

export function CardTitle({ className, ...props }: HTMLAttributes<HTMLHeadingElement>) {
  return (
    <h3
      className={cn("text-lg font-semibold text-foreground", className)}
      {...props}
    />
  );
}

export function CardDescription({ className, ...props }: HTMLAttributes<HTMLParagraphElement>) {
  return (
    <p
      className={cn("text-sm text-text-secondary", className)}
      {...props}
    />
  );
}

export function CardContent({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn("p-6 pt-0", className)}
      {...props}
    />
  );
}

export function CardFooter({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn("flex items-center p-6 pt-0", className)}
      {...props}
    />
  );
}
`
}

func desktopUIBadge() string {
	return `import { HTMLAttributes } from "react";
import { cn } from "@/lib/utils";

type Variant = "default" | "secondary" | "success" | "danger" | "outline";

interface BadgeProps extends HTMLAttributes<HTMLSpanElement> {
  variant?: Variant;
}

const variantClasses: Record<Variant, string> = {
  default: "bg-accent/20 text-accent",
  secondary: "bg-bg-hover text-text-secondary",
  success: "bg-success/20 text-success",
  danger: "bg-danger/20 text-danger",
  outline: "border border-border text-text-secondary bg-transparent",
};

export default function Badge({ variant = "default", className, ...props }: BadgeProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium",
        variantClasses[variant],
        className
      )}
      {...props}
    />
  );
}
`
}

func desktopUIDialog() string {
	return `import { ReactNode, useEffect, useCallback } from "react";
import { createPortal } from "react-dom";
import { cn } from "@/lib/utils";

type DialogSize = "sm" | "md" | "lg";

interface DialogProps {
  isOpen: boolean;
  onClose: () => void;
  title?: string;
  children: ReactNode;
  size?: DialogSize;
}

const sizeClasses: Record<DialogSize, string> = {
  sm: "max-w-sm",
  md: "max-w-md",
  lg: "max-w-lg",
};

export default function Dialog({
  isOpen,
  onClose,
  title,
  children,
  size = "md",
}: DialogProps) {
  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    },
    [onClose]
  );

  useEffect(() => {
    if (isOpen) {
      document.addEventListener("keydown", handleKeyDown);
      return () => document.removeEventListener("keydown", handleKeyDown);
    }
  }, [isOpen, handleKeyDown]);

  if (!isOpen) return null;

  return createPortal(
    <div
      className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose();
      }}
    >
      <div
        className={cn(
          "w-full rounded-xl border border-border bg-bg-elevated shadow-lg",
          sizeClasses[size]
        )}
      >
        {title && (
          <div className="flex items-center justify-between border-b border-border px-6 py-4">
            <h2 className="text-lg font-semibold text-foreground">{title}</h2>
            <button
              onClick={onClose}
              className="text-text-muted hover:text-foreground transition-colors"
            >
              &#x2715;
            </button>
          </div>
        )}
        <div className="p-6">{children}</div>
      </div>
    </div>,
    document.body
  );
}
`
}

func desktopUIDropdownMenu() string {
	return `import { ReactNode, useState, useRef, useEffect } from "react";
import { cn } from "@/lib/utils";

interface DropdownMenuProps {
  trigger: ReactNode;
  children: ReactNode;
  align?: "left" | "right";
}

export default function DropdownMenu({
  trigger,
  children,
  align = "left",
}: DropdownMenuProps) {
  const [isOpen, setIsOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

  return (
    <div ref={ref} className="relative inline-block">
      <div onClick={() => setIsOpen((prev) => !prev)}>{trigger}</div>
      {isOpen && (
        <div
          className={cn(
            "absolute z-50 mt-1 min-w-[180px] rounded-lg border border-border bg-bg-elevated py-1 shadow-lg",
            align === "right" && "right-0",
            align === "left" && "left-0"
          )}
        >
          {children}
        </div>
      )}
    </div>
  );
}

interface DropdownMenuItemProps {
  children: ReactNode;
  onClick?: () => void;
  className?: string;
  destructive?: boolean;
}

export function DropdownMenuItem({
  children,
  onClick,
  className,
  destructive = false,
}: DropdownMenuItemProps) {
  return (
    <button
      className={cn(
        "flex w-full items-center px-3 py-2 text-sm transition-colors",
        destructive
          ? "text-danger hover:bg-danger/10"
          : "text-text-secondary hover:bg-bg-hover hover:text-foreground",
        className
      )}
      onClick={onClick}
    >
      {children}
    </button>
  );
}
`
}

func desktopUISelect() string {
	return `import { forwardRef, SelectHTMLAttributes } from "react";
import { cn } from "@/lib/utils";

const Select = forwardRef<HTMLSelectElement, SelectHTMLAttributes<HTMLSelectElement>>(
  ({ className, children, ...props }, ref) => {
    return (
      <select
        ref={ref}
        className={cn(
          "w-full appearance-none bg-bg-elevated border border-border rounded-lg px-3 py-2 pr-8 text-sm text-foreground focus:border-accent focus:ring-1 focus:ring-accent outline-none transition-colors disabled:opacity-50 disabled:cursor-not-allowed",
          "bg-[url('data:image/svg+xml;charset=utf-8,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20width%3D%2216%22%20height%3D%2216%22%20viewBox%3D%220%200%2024%2024%22%20fill%3D%22none%22%20stroke%3D%22%239090a8%22%20stroke-width%3D%222%22%3E%3Cpath%20d%3D%22m6%209%206%206%206-6%22%2F%3E%3C%2Fsvg%3E')] bg-[length:16px] bg-[right_8px_center] bg-no-repeat",
          className
        )}
        {...props}
      >
        {children}
      </select>
    );
  }
);

Select.displayName = "Select";

export default Select;
`
}

func desktopUITabs() string {
	return `import { createContext, useContext, useState, ReactNode } from "react";
import { cn } from "@/lib/utils";

interface TabsContextType {
  activeTab: string;
  setActiveTab: (tab: string) => void;
}

const TabsContext = createContext<TabsContextType | null>(null);

function useTabs() {
  const ctx = useContext(TabsContext);
  if (!ctx) throw new Error("Tabs components must be used within <Tabs>");
  return ctx;
}

interface TabsProps {
  defaultValue: string;
  children: ReactNode;
  className?: string;
}

export function Tabs({ defaultValue, children, className }: TabsProps) {
  const [activeTab, setActiveTab] = useState(defaultValue);

  return (
    <TabsContext.Provider value={{ activeTab, setActiveTab }}>
      <div className={cn("w-full", className)}>{children}</div>
    </TabsContext.Provider>
  );
}

interface TabsListProps {
  children: ReactNode;
  className?: string;
}

export function TabsList({ children, className }: TabsListProps) {
  return (
    <div
      className={cn(
        "flex border-b border-border",
        className
      )}
    >
      {children}
    </div>
  );
}

interface TabsTriggerProps {
  value: string;
  children: ReactNode;
  className?: string;
}

export function TabsTrigger({ value, children, className }: TabsTriggerProps) {
  const { activeTab, setActiveTab } = useTabs();
  const isActive = activeTab === value;

  return (
    <button
      className={cn(
        "px-4 py-2 text-sm font-medium transition-colors -mb-px",
        isActive
          ? "text-accent border-b-2 border-accent"
          : "text-text-secondary hover:text-foreground",
        className
      )}
      onClick={() => setActiveTab(value)}
    >
      {children}
    </button>
  );
}

interface TabsContentProps {
  value: string;
  children: ReactNode;
  className?: string;
}

export function TabsContent({ value, children, className }: TabsContentProps) {
  const { activeTab } = useTabs();

  if (activeTab !== value) return null;

  return (
    <div className={cn("pt-4", className)}>
      {children}
    </div>
  );
}
`
}

func desktopUIAvatar() string {
	return `import { cn } from "@/lib/utils";

type AvatarSize = "sm" | "md" | "lg";

interface AvatarProps {
  src?: string;
  name: string;
  size?: AvatarSize;
  className?: string;
}

const sizeClasses: Record<AvatarSize, string> = {
  sm: "h-8 w-8 text-xs",
  md: "h-10 w-10 text-sm",
  lg: "h-14 w-14 text-lg",
};

function getInitials(name: string): string {
  const parts = name.trim().split(" ");
  if (parts.length >= 2) {
    return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
  }
  return name.slice(0, 2).toUpperCase();
}

export default function Avatar({ src, name, size = "md", className }: AvatarProps) {
  if (src) {
    return (
      <img
        src={src}
        alt={name}
        className={cn(
          "rounded-full object-cover",
          sizeClasses[size],
          className
        )}
      />
    );
  }

  return (
    <div
      className={cn(
        "inline-flex items-center justify-center rounded-full bg-accent/20 text-accent font-medium",
        sizeClasses[size],
        className
      )}
    >
      {getInitials(name)}
    </div>
  );
}
`
}

func desktopUIToast() string {
	return `export { toast } from "sonner";
`
}

func desktopUIScrollArea() string {
	return `import { HTMLAttributes } from "react";
import { cn } from "@/lib/utils";

interface ScrollAreaProps extends HTMLAttributes<HTMLDivElement> {}

export default function ScrollArea({ className, children, ...props }: ScrollAreaProps) {
  return (
    <div
      className={cn("overflow-auto scrollbar", className)}
      {...props}
    >
      {children}
    </div>
  );
}
`
}
