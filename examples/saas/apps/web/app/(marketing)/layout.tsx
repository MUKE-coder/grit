import { Navbar } from "@/components/navbar";
import { Footer } from "@/components/footer";

// The public site. Anything under app/(marketing) gets this chrome; the group
// name is in parentheses, so it is not part of the URL: this file wraps / and
// /blog, and nothing else in the app.
export default function MarketingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      <Navbar />
      <main className="min-h-screen">{children}</main>
      <Footer />
    </>
  );
}
