import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Support Grit | Grit",
  description:
    "Support the Grit open-source framework with a one-time donation. Help fund new features, documentation, and community support.",
};

export default function DonateLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <>{children}</>;
}
