"use client";

import { useState, useRef, useEffect, KeyboardEvent, ClipboardEvent } from "react";
import { RefreshCw } from "lucide-react";

interface OtpInput01Props {
  email?: string;
  onComplete?: (code: string) => void;
}

export default function OtpInput01({ email = "you@company.com", onComplete }: OtpInput01Props) {
  const [digits, setDigits] = useState<string[]>(Array(6).fill(""));
  const [countdown, setCountdown] = useState(60);
  const [canResend, setCanResend] = useState(false);
  const inputRefs = useRef<(HTMLInputElement | null)[]>([]);

  useEffect(() => {
    if (countdown <= 0) {
      setCanResend(true);
      return;
    }
    const t = setTimeout(() => setCountdown((c) => c - 1), 1000);
    return () => clearTimeout(t);
  }, [countdown]);

  const handleResend = () => {
    setCountdown(60);
    setCanResend(false);
    setDigits(Array(6).fill(""));
    inputRefs.current[0]?.focus();
  };

  const handleChange = (index: number, value: string) => {
    const char = value.replace(/\D/g, "").slice(-1);
    const next = [...digits];
    next[index] = char;
    setDigits(next);
    if (char && index < 5) {
      inputRefs.current[index + 1]?.focus();
    }
    if (next.every((d) => d !== "")) {
      onComplete?.(next.join(""));
    }
  };

  const handleKeyDown = (index: number, e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Backspace") {
      if (digits[index]) {
        const next = [...digits];
        next[index] = "";
        setDigits(next);
      } else if (index > 0) {
        inputRefs.current[index - 1]?.focus();
      }
    }
    if (e.key === "ArrowLeft" && index > 0) inputRefs.current[index - 1]?.focus();
    if (e.key === "ArrowRight" && index < 5) inputRefs.current[index + 1]?.focus();
  };

  const handlePaste = (e: ClipboardEvent<HTMLInputElement>) => {
    e.preventDefault();
    const pasted = e.clipboardData.getData("text").replace(/\D/g, "").slice(0, 6);
    if (!pasted) return;
    const next = Array(6).fill("");
    pasted.split("").forEach((char, i) => { next[i] = char; });
    setDigits(next);
    const lastFilled = Math.min(pasted.length, 5);
    inputRefs.current[lastFilled]?.focus();
    if (pasted.length === 6) onComplete?.(pasted);
  };

  return (
    <div className="min-h-screen bg-background flex items-center justify-center px-4">
      <div className="w-full max-w-sm text-center">
        <div className="flex justify-center mb-8">
          <div className="w-12 h-12 rounded-2xl bg-accent flex items-center justify-center">
            <span className="text-white font-black text-xl tracking-tight">G</span>
          </div>
        </div>

        <h1 className="text-2xl font-bold text-foreground tracking-tight mb-2">
          Check your email
        </h1>
        <p className="text-text-secondary text-sm mb-1">
          We sent a 6-digit code to
        </p>
        <p className="text-foreground text-sm font-semibold mb-8">{email}</p>

        {/* OTP boxes */}
        <div className="flex justify-center gap-3 mb-8">
          {digits.map((digit, i) => (
            <input
              key={i}
              ref={(el) => { inputRefs.current[i] = el; }}
              type="text"
              inputMode="numeric"
              maxLength={1}
              value={digit}
              onChange={(e) => handleChange(i, e.target.value)}
              onKeyDown={(e) => handleKeyDown(i, e)}
              onPaste={handlePaste}
              className={"w-12 h-14 text-center text-xl font-bold rounded-xl border transition-colors duration-150 bg-bg-elevated text-foreground focus:outline-none " +
                (digit
                  ? "border-accent text-accent"
                  : "border-border focus:border-accent")}
            />
          ))}
        </div>

        {/* Resend */}
        <div className="flex items-center justify-center gap-2 text-sm">
          {canResend ? (
            <button
              onClick={handleResend}
              className="inline-flex items-center gap-1.5 text-accent hover:underline font-medium"
            >
              <RefreshCw size={14} />
              Resend code
            </button>
          ) : (
            <p className="text-text-muted">
              Resend code in{" "}
              <span className="text-text-secondary font-semibold">{countdown}s</span>
            </p>
          )}
        </div>

        <p className="mt-6 text-xs text-text-muted">
          Didn&rsquo;t get the email? Check your spam folder.
        </p>
      </div>
    </div>
  );
}
