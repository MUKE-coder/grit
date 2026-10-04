import type { Metadata } from "next";
import { Inter, JetBrains_Mono } from "next/font/google";

const inter = Inter({ subsets: ["latin"], variable: "--font-display" });
const jetbrainsMono = JetBrains_Mono({ subsets: ["latin"], variable: "--font-mono", preload: false });
import "./globals.css";
import { Providers } from "@/components/shared/providers";

export const metadata: Metadata = {
  title: "library Admin",
  description: "Admin panel, built with Grit",
  // <meta name="generator" content="Grit">, the convention every static site
  // generator and CMS uses and the one technology scanners read. No version:
  // knowing an app is a Grit app is fine, and knowing which release it is
  // running is a list of that release's advisories. Delete the line to opt out.
  generator: "Grit",
};

// Applied synchronously, before the first paint, from the <head> below.
// Writes all three signals the stylesheets key off: data-theme-mode for the
// CSS variable cascade, .dark for Tailwind's darkMode: "class", and
// color-scheme so the browser's own scrollbars and inputs match.
const themeScript =
  "(function(){try{var m=localStorage.getItem('grit-theme-mode');if(m!=='dark'&&m!=='light'){m=window.matchMedia&&window.matchMedia('(prefers-color-scheme: dark)').matches?'dark':'light';}var r=document.documentElement;r.setAttribute('data-theme-mode',m);r.classList.toggle('dark',m==='dark');r.style.colorScheme=m;}catch(e){}})();";

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  // v3.28.1: data-theme drives the CSS variable cascade in globals.css.
  // Reading process.env at render time means changing THEME in .env +
  // restarting the dev server re-paints the dashboard without code edits.
  // Defaults to "atlas" so a missing env doesn't blank the surface.
  const dataTheme = process.env.NEXT_PUBLIC_THEME || "atlas";

  return (
    // suppressHydrationWarning: the theme script below and DarkModeToggle both
    // mutate html.classList + data-theme-mode + style. Without this React would
    // log a noisy mismatch on the first paint even though the behaviour is
    // intentional.
    <html lang="en" data-theme={dataTheme} suppressHydrationWarning>
      <head>
        {/*
          The stored theme is applied before the browser paints anything. It used
          to be applied in an effect inside DarkModeToggle, which mounts only
          after /auth/me answers, so every load of a dark dashboard showed a
          light one first for as long as that request took.

          dangerouslySetInnerHTML is the only way to get a synchronous inline
          script into the document: next/script defers it, and a deferred script
          is the flash again.
        */}
        {/* biome-ignore lint/security/noDangerouslySetInnerHtml: a constant script of Grit's, no user input. */}
        <script dangerouslySetInnerHTML={{ __html: themeScript }} />
      </head>
      <body className={`${inter.variable} ${jetbrainsMono.variable} min-h-screen font-sans antialiased`}>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
