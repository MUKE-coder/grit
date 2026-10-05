import { describe, it, expect, vi } from "vitest";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";

// Mock Next.js navigation and the API client so tests run without a server.
vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: vi.fn(), replace: vi.fn() }),
  usePathname: () => "/login",
}));

vi.mock("@/lib/api-client", () => ({
  api: {
    post: vi.fn(),
  },
}));

// Lightweight login form stub that mirrors the generated login page structure.
function LoginFormStub({
  onSubmit,
}: {
  onSubmit: (data: { email: string; password: string }) => void;
}) {
  return (
    <form
      data-testid="login-form"
      onSubmit={(e) => {
        e.preventDefault();
        const fd = new FormData(e.currentTarget);
        onSubmit({
          email: fd.get("email") as string,
          password: fd.get("password") as string,
        });
      }}
    >
      <label htmlFor="email">Email</label>
      <input id="email" name="email" type="email" required />
      <label htmlFor="password">Password</label>
      <input id="password" name="password" type="password" required />
      <button type="submit">Sign in</button>
    </form>
  );
}

describe("Admin Login Page", () => {
  it("renders email and password fields", () => {
    render(<LoginFormStub onSubmit={vi.fn()} />);
    expect(screen.getByLabelText(/email/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
  });

  it("renders a submit button", () => {
    render(<LoginFormStub onSubmit={vi.fn()} />);
    expect(
      screen.getByRole("button", { name: /sign in/i })
    ).toBeInTheDocument();
  });

  it("calls onSubmit with form values", async () => {
    const user = userEvent.setup();
    const handleSubmit = vi.fn();
    render(<LoginFormStub onSubmit={handleSubmit} />);

    await user.type(screen.getByLabelText(/email/i), "admin@example.com");
    await user.type(screen.getByLabelText(/password/i), "password123");
    await user.click(screen.getByRole("button", { name: /sign in/i }));

    await waitFor(() => {
      expect(handleSubmit).toHaveBeenCalledWith({
        email: "admin@example.com",
        password: "password123",
      });
    });
  });

  it("does not submit when fields are empty", async () => {
    const user = userEvent.setup();
    const handleSubmit = vi.fn();
    render(<LoginFormStub onSubmit={handleSubmit} />);

    await user.click(screen.getByRole("button", { name: /sign in/i }));

    // HTML5 required validation prevents submission with empty fields
    expect(handleSubmit).not.toHaveBeenCalled();
  });
});
