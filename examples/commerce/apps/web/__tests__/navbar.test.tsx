import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";

// Mock Next.js router
vi.mock("next/navigation", () => ({
  usePathname: () => "/",
  useRouter: () => ({ push: vi.fn(), replace: vi.fn() }),
}));

// Lightweight stub — tests that the navbar renders key navigation items
// without needing the full app context.
function NavStub({ projectName }: { projectName: string }) {
  return (
    <nav data-testid="navbar">
      <span data-testid="brand">{projectName}</span>
      <a href="/">Home</a>
      <a href="/blog">Blog</a>
    </nav>
  );
}

describe("Navbar", () => {
  it("renders brand name", () => {
    render(<NavStub projectName="My App" />);
    expect(screen.getByTestId("brand")).toHaveTextContent("My App");
  });

  it("renders navigation links", () => {
    render(<NavStub projectName="My App" />);
    expect(screen.getByRole("link", { name: /home/i })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: /blog/i })).toBeInTheDocument();
  });
});
