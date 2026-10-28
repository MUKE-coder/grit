import React from "react";
import type { Metadata } from "next";
import Script from "next/script";
import { GeistSans } from "geist/font/sans";
import { GeistMono } from "geist/font/mono";
import { ThemeProvider } from "@/components/theme-provider";
import { OrganizationSchema } from "@/components/structured-data";
import { NexoraChat } from "@/components/nexora-chat";
import { Analytics } from "zenith-analytics/next";
import { siteConfig } from "@/config/site";
import { ZENITH_PUBLIC } from "@/config/zenith";

import "./globals.css";

// One family pair: Geist for everything, Geist Mono for anything a developer
// would type. This replaces Ubuntu headings plus JetBrains Mono code, which
// was three voices for one brand.
//
// From the `geist` package rather than next/font/google: it ships the variable
// fonts locally, so there is no request to Google and no layout shift while
// one resolves.

export const metadata: Metadata = {
  metadataBase: new URL(siteConfig.url),
  title: {
    default: siteConfig.title,
    template: `%s | ${siteConfig.name}`,
  },
  description: siteConfig.description,
  keywords: siteConfig.keywords,
  authors: [{ name: siteConfig.author, url: siteConfig.website }],
  creator: siteConfig.author,
  publisher: siteConfig.author,
  openGraph: {
    type: "website",
    locale: "en_US",
    url: siteConfig.url,
    title: siteConfig.title,
    description: siteConfig.description,
    siteName: siteConfig.name,
    images: [
      {
        url: siteConfig.ogImage,
        width: 1200,
        height: 630,
        alt: "Grit — Go + React Full-Stack Framework",
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title: siteConfig.title,
    description: siteConfig.description,
    images: [siteConfig.ogImage],
    creator: "@jaboronics",
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      "max-video-preview": -1,
      "max-image-preview": "large",
      "max-snippet": -1,
    },
  },
  alternates: {
    canonical: siteConfig.url,
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${GeistSans.variable} ${GeistMono.variable} font-sans antialiased`}
      >
        <ThemeProvider attribute="class" defaultTheme="dark" enableSystem>
          <OrganizationSchema />
          {children}
          <NexoraChat />
        </ThemeProvider>
        {/* Zenith analytics — server-rendered so the ~1KB snippet is inlined
            into the HTML. Cookieless, so there is nothing to consent to. */}
        <Analytics config={ZENITH_PUBLIC} required />
        {process.env.NEXT_PUBLIC_UMAMI_WEBSITE_ID && (
          <Script
            defer
            src={`${process.env.NEXT_PUBLIC_UMAMI_URL || "https://analytics.gritframework.dev"}/script.js`}
            data-website-id={process.env.NEXT_PUBLIC_UMAMI_WEBSITE_ID}
            strategy="afterInteractive"
          />
        )}
      </body>
    </html>
  );
}
