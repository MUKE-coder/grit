import { createFileRoute, Outlet, Navigate } from '@tanstack/react-router'
import { useAuthStore, useAuthHydrated } from '@/stores/auth.store'
import { TrendingUp, ShoppingCart, BarChart3 } from 'lucide-react'
import { useState, useEffect } from 'react'

export const Route = createFileRoute('/_auth')({
  component: AuthLayout,
})

const slides = [
  {
    badge: 'Motorcycle Sales',
    title: 'Cash Or Loan',
    description: 'Sell motorcycles by cash or finance them on flexible loan plans. Every unit tracked by number plate.',
  },
  {
    badge: 'Spares Inventory',
    title: 'Sell Parts In Seconds',
    description: 'A POS built for speed across every branch — barcode scanning, quick search, one-tap checkout.',
  },
  {
    badge: 'Loans & Repayments',
    title: 'Collect On Time',
    description: 'Auto-generated repayment schedules, mobile-money collections, and overdue tracking — all in one place.',
  },
]

function AuthLayout() {
  const hydrated = useAuthHydrated()
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated())
  const [currentSlide, setCurrentSlide] = useState(0)

  useEffect(() => {
    const timer = setInterval(() => {
      setCurrentSlide((prev) => (prev + 1) % slides.length)
    }, 5000)
    return () => clearInterval(timer)
  }, [])

  if (!hydrated) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background">
        <div className="w-10 h-10 rounded-xl bg-accent flex items-center justify-center text-white font-bold text-sm animate-pulse">
          KM
        </div>
      </div>
    )
  }

  if (isAuthenticated) {
    return <Navigate to="/" />
  }

  const slide = slides[currentSlide]
  const icons = [TrendingUp, ShoppingCart, BarChart3]
  const SlideIcon = icons[currentSlide]

  return (
    <div className="min-h-screen flex bg-background">
      {/* Left panel — branding (dark, accent-tinted) */}
      <div className="hidden lg:flex lg:w-[480px] xl:w-[520px] shrink-0 relative overflow-hidden flex-col bg-foreground">
        {/* Subtle dot pattern */}
        <div className="absolute inset-0 opacity-10">
          <svg width="100%" height="100%">
            <defs>
              <pattern id="dots" x="0" y="0" width="30" height="30" patternUnits="userSpaceOnUse">
                <circle cx="2" cy="2" r="1.5" fill="white" />
              </pattern>
            </defs>
            <rect width="100%" height="100%" fill="url(#dots)" />
          </svg>
        </div>

        {/* Content */}
        <div className="relative z-10 flex flex-col h-full p-10">
          {/* Logo */}
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-accent flex items-center justify-center text-white font-bold text-[13px] shrink-0">
              KM
            </div>
            <span className="text-white text-[18px] font-semibold tracking-tight">Grit Motors</span>
          </div>

          {/* Slide content */}
          <div className="mt-auto mb-auto">
            <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-white/10 text-white text-[12.5px] font-medium mb-5">
              <SlideIcon size={14} />
              {slide.badge}
            </div>

            <h2 className="text-[36px] xl:text-[40px] font-semibold text-white leading-tight mb-3 tracking-tight">
              {slide.title}
            </h2>

            <p className="text-white/70 text-[15px] leading-relaxed max-w-sm">
              {slide.description}
            </p>

            {/* Slide indicators */}
            <div className="flex gap-2 mt-7">
              {slides.map((_, i) => (
                <button
                  key={i}
                  type="button"
                  onClick={() => setCurrentSlide(i)}
                  className={`h-1.5 rounded-full transition-all duration-300 ${
                    i === currentSlide ? 'w-7 bg-white' : 'w-1.5 bg-white/30'
                  }`}
                  aria-label={`Slide ${i + 1}`}
                />
              ))}
            </div>
          </div>

          {/* Tagline */}
          <p className="text-white/40 text-[12.5px]">
            Motorcycles. Spares. Loans. One platform.
          </p>
        </div>

        {/* Decorative blue accent blob */}
        <div className="absolute -bottom-24 -right-24 w-72 h-72 rounded-full bg-accent/30 blur-3xl pointer-events-none" />
      </div>

      {/* Right panel — form */}
      <div className="flex-1 flex items-center justify-center px-6 py-12">
        <div className="w-full max-w-md">
          {/* Mobile logo (hidden on desktop where the left panel handles it) */}
          <div className="lg:hidden text-center mb-8">
            <div className="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-accent text-white text-[14px] font-bold shadow-xs mb-3">
              KM
            </div>
            <h1 className="text-[18px] font-semibold text-foreground">Grit Motors</h1>
          </div>

          <Outlet />
        </div>
      </div>
    </div>
  )
}
