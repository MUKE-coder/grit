import type { Config } from 'tailwindcss'

const config: Config = {
  darkMode: ['class'],
  content: [
    './pages/**/*.{js,ts,jsx,tsx,mdx}',
    './components/**/*.{js,ts,jsx,tsx,mdx}',
    './app/**/*.{js,ts,jsx,tsx,mdx}',
    '*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      fontFamily: {
        // One family pair, per the design system: Geist and Geist Mono.
        // `display` is kept as a name because the docs use it in dozens of
        // places, but it resolves to Geist now rather than Ubuntu.
        sans: ['var(--font-geist-sans)', 'ui-sans-serif', 'system-ui', 'sans-serif'],
        display: ['var(--font-geist-sans)', 'ui-sans-serif', 'system-ui', 'sans-serif'],
        mono: ['var(--font-geist-mono)', 'ui-monospace', 'SFMono-Regular', 'Menlo', 'monospace'],
      },
      /*
       * The twelve named type styles, from the design system.
       *
       * They have names here because the contract is "use one of the twelve",
       * and a rule like that is only checkable if the twelve are things you
       * can name. Before this block, `text-lead` was written in two
       * components and existed in no stylesheet: Tailwind emits nothing for
       * a class it does not know, so every section lead on the landing page
       * rendered at the inherited 16px instead of 20, silently, for as long
       * as it took somebody to measure it.
       *
       * Each entry carries its line height and weight, so `text-h2` is the
       * whole style and not just a size somebody still has to pair.
       */
      fontSize: {
        'display-xl': ['80px', { lineHeight: '84px', fontWeight: '600', letterSpacing: '-0.045em' }],
        display: ['56px', { lineHeight: '60px', fontWeight: '600', letterSpacing: '-0.04em' }],
        stat: ['64px', { lineHeight: '64px', fontWeight: '500', letterSpacing: '-0.04em' }],
        h1: ['40px', { lineHeight: '48px', fontWeight: '600', letterSpacing: '-0.035em' }],
        h2: ['28px', { lineHeight: '36px', fontWeight: '600', letterSpacing: '-0.025em' }],
        h3: ['20px', { lineHeight: '28px', fontWeight: '600', letterSpacing: '-0.02em' }],
        h4: ['16px', { lineHeight: '24px', fontWeight: '600', letterSpacing: '-0.01em' }],
        lead: ['20px', { lineHeight: '32px', fontWeight: '400' }],
        body: ['16px', { lineHeight: '28px', fontWeight: '400' }],
        'body-sm': ['14px', { lineHeight: '22px', fontWeight: '400' }],
        caption: ['13px', { lineHeight: '20px', fontWeight: '400' }],
        eyebrow: ['12px', { lineHeight: '16px', fontWeight: '500', letterSpacing: '0.08em' }],
        code: ['13.5px', { lineHeight: '24px', fontWeight: '400' }],
      },
      colors: {
        // Grit design system. The shadcn names below are aliases kept so
        // every existing components/ui/ file keeps working.
        bg: 'rgb(from var(--bg) r g b / <alpha-value>)',
        surface: {
          1: 'rgb(from var(--surface-1) r g b / <alpha-value>)',
          2: 'rgb(from var(--surface-2) r g b / <alpha-value>)',
          3: 'rgb(from var(--surface-3) r g b / <alpha-value>)',
        },
        fg: {
          DEFAULT: 'rgb(from var(--fg) r g b / <alpha-value>)',
          muted: 'rgb(from var(--fg-muted) r g b / <alpha-value>)',
          subtle: 'rgb(from var(--fg-subtle) r g b / <alpha-value>)',
        },
        accent: {
          DEFAULT: 'rgb(from var(--accent) r g b / <alpha-value>)',
          fg: 'rgb(from var(--accent-fg) r g b / <alpha-value>)',
          soft: 'var(--accent-soft)',
          line: 'var(--accent-line)',
        },
        // shadcn's `accent` means hover surface, not brand colour. It keeps
        // its own name so `bg-accent` is unambiguously the violet.
        'accent-surface': {
          DEFAULT: 'rgb(from var(--accent-surface) r g b / <alpha-value>)',
          foreground: 'rgb(from var(--accent-foreground) r g b / <alpha-value>)',
        },
        success: 'rgb(from var(--success) r g b / <alpha-value>)',
        warning: 'rgb(from var(--warning) r g b / <alpha-value>)',
        danger: 'rgb(from var(--danger) r g b / <alpha-value>)',
        // The code palette: the syntax roles, so a hand-written snippet can
        // be coloured from tokens rather than from hex.
        //
        // There is deliberately no DEFAULT here. `text-code` is a font size
        // in the scale above, and a colour of the same name would make one
        // class quietly set both a size and a colour. The panel background
        // is `bg-code-bg`.
        code: {
          bg: 'rgb(from var(--code-bg) r g b / <alpha-value>)',
          fg: 'rgb(from var(--code-fg) r g b / <alpha-value>)',
          comment: 'rgb(from var(--code-comment) r g b / <alpha-value>)',
          keyword: 'rgb(from var(--code-keyword) r g b / <alpha-value>)',
          string: 'rgb(from var(--code-string) r g b / <alpha-value>)',
          fn: 'rgb(from var(--code-fn) r g b / <alpha-value>)',
          number: 'rgb(from var(--code-number) r g b / <alpha-value>)',
        },
        chart: {
          // Grit carries colour; every rival is one neutral grey, so a chart
          // never reads as a traffic light.
          grit: 'rgb(from var(--chart-grit) r g b / <alpha-value>)',
          rival: 'rgb(from var(--chart-rival) r g b / <alpha-value>)',
          '1': 'rgb(from var(--chart-1) r g b / <alpha-value>)',
          '2': 'rgb(from var(--chart-2) r g b / <alpha-value>)',
          '3': 'rgb(from var(--chart-3) r g b / <alpha-value>)',
          '4': 'rgb(from var(--chart-4) r g b / <alpha-value>)',
          '5': 'rgb(from var(--chart-5) r g b / <alpha-value>)',
        },
        background: 'rgb(from var(--background) r g b / <alpha-value>)',
        foreground: 'rgb(from var(--foreground) r g b / <alpha-value>)',
        card: {
          DEFAULT: 'rgb(from var(--card) r g b / <alpha-value>)',
          foreground: 'rgb(from var(--card-foreground) r g b / <alpha-value>)',
        },
        popover: {
          DEFAULT: 'rgb(from var(--popover) r g b / <alpha-value>)',
          foreground: 'rgb(from var(--popover-foreground) r g b / <alpha-value>)',
        },
        primary: {
          DEFAULT: 'rgb(from var(--primary) r g b / <alpha-value>)',
          foreground: 'rgb(from var(--primary-foreground) r g b / <alpha-value>)',
        },
        secondary: {
          DEFAULT: 'rgb(from var(--secondary) r g b / <alpha-value>)',
          foreground: 'rgb(from var(--secondary-foreground) r g b / <alpha-value>)',
        },
        muted: {
          DEFAULT: 'rgb(from var(--muted) r g b / <alpha-value>)',
          foreground: 'rgb(from var(--muted-foreground) r g b / <alpha-value>)',
        },
        destructive: {
          DEFAULT: 'rgb(from var(--destructive) r g b / <alpha-value>)',
          foreground: 'rgb(from var(--destructive-foreground) r g b / <alpha-value>)',
        },
        // `border-border` is the hairline, `border-border-strong` the one a
        // control draws so it reads as interactive rather than as a panel.
        border: {
          DEFAULT: 'rgb(from var(--border) r g b / <alpha-value>)',
          strong: 'rgb(from var(--border-strong) r g b / <alpha-value>)',
        },
        input: 'rgb(from var(--input) r g b / <alpha-value>)',
        ring: 'rgb(from var(--ring) r g b / <alpha-value>)',
        sidebar: {
          DEFAULT: 'rgb(from var(--sidebar-background) r g b / <alpha-value>)',
          foreground: 'rgb(from var(--sidebar-foreground) r g b / <alpha-value>)',
          primary: 'rgb(from var(--sidebar-primary) r g b / <alpha-value>)',
          'primary-foreground': 'rgb(from var(--sidebar-primary-foreground) r g b / <alpha-value>)',
          accent: 'rgb(from var(--sidebar-accent) r g b / <alpha-value>)',
          'accent-foreground': 'rgb(from var(--sidebar-accent-foreground) r g b / <alpha-value>)',
          border: 'rgb(from var(--sidebar-border) r g b / <alpha-value>)',
          ring: 'rgb(from var(--sidebar-ring) r g b / <alpha-value>)',
        },
      },
      // DESIGN.md section 4 allows three radii plus full: 6, 10 and 16.
      //
      // Tailwind ships five names, and the old mapping generated 6, 8 and 10
      // from one --radius, so `rounded-md` produced an 8px that is on no
      // scale. Collapsing all five names onto the three allowed values makes
      // the budget unbreakable through a utility: there is no class that can
      // produce an off-scale corner.
      // Every name reads a token from globals.css, so there is no arithmetic
      // here that can invent a value. The old mapping derived 6, 8 and 10
      // from one --radius, and the 8px it produced was on no scale and on
      // every button.
      //
      // Names outnumber values on purpose. `lg` stays at 10 and bare
      // `rounded` stays at 6 because between them they are 2,600 corners on
      // this site, and moving a radius 2px in 2,600 places is a large diff
      // that nobody asked for and nobody can see. The two button radii have
      // names of their own rather than a Tailwind size, because they belong
      // to a control rather than to a step on the scale.
      borderRadius: {
        DEFAULT: 'var(--radius-sm)', // 6
        xs: 'var(--radius-xs)', // 4
        sm: 'var(--radius-sm)', // 6
        'btn-sm': 'var(--radius-btn-sm)', // 8
        md: 'var(--radius-md)', // 10
        lg: 'var(--radius-md)', // 10
        btn: 'var(--radius-btn-lg)', // 12
        xl: 'var(--radius-lg)', // 16
        '2xl': 'var(--radius-xl)', // 20
        '3xl': 'var(--radius-xl)', // 20
      },
      keyframes: {
        'accordion-down': {
          from: { height: '0' },
          to: { height: 'var(--radix-accordion-content-height)' },
        },
        'accordion-up': {
          from: { height: 'var(--radix-accordion-content-height)' },
          to: { height: '0' },
        },
        'fade-in': {
          from: { opacity: '0', transform: 'translateY(8px)' },
          to: { opacity: '1', transform: 'translateY(0)' },
        },
        'pulse-glow': {
          '0%, 100%': { opacity: '0.4' },
          '50%': { opacity: '0.8' },
        },
      },
      animation: {
        'accordion-down': 'accordion-down 0.2s ease-out',
        'accordion-up': 'accordion-up 0.2s ease-out',
        'fade-in': 'fade-in 0.5s ease-out forwards',
        'pulse-glow': 'pulse-glow 3s ease-in-out infinite',
      },
    },
  },
  plugins: [require('tailwindcss-animate')],
}
export default config
