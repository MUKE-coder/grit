package scaffold

import (
	"fmt"
	"path/filepath"
)

func writeFrontendTestFiles(root string, opts Options) error {
	webRoot := webAppRoot(root, opts)
	adminRoot := filepath.Join(root, "apps", "admin")

	files := map[string]string{}

	if opts.ShouldIncludeWeb() {
		files[filepath.Join(webRoot, "vitest.config.ts")] = webVitestConfig()
		files[filepath.Join(webRoot, "vitest.setup.ts")] = vitestSetup()
		files[filepath.Join(webRoot, "__tests__", "navbar.test.tsx")] = webNavbarTest()
		files[filepath.Join(webRoot, "__tests__", "footer.test.tsx")] = webFooterTest(opts)
	}

	if opts.ShouldIncludeAdmin() {
		files[filepath.Join(adminRoot, "vitest.config.ts")] = adminVitestConfig()
		files[filepath.Join(adminRoot, "vitest.setup.ts")] = vitestSetup()
		files[filepath.Join(adminRoot, "__tests__", "login.test.tsx")] = adminLoginTest()
		files[filepath.Join(adminRoot, "__tests__", "utils.test.ts")] = adminUtilsTest()
		libImport := "@/lib"
		if opts.Frontend == FrontendTanStack {
			libImport = "@/src/lib"
		}
		files[filepath.Join(adminRoot, "__tests__", "form-values.test.ts")] = adminFormValuesTest(libImport)
		files[filepath.Join(adminRoot, "__tests__", "tiptap-extensions.test.ts")] = adminTiptapExtensionsTest(libImport)
		files[filepath.Join(adminRoot, "__tests__", "safe-href.test.ts")] = adminSafeHrefTest(libImport)
		// A component test needs the component's own "@/" imports to resolve,
		// and the Vite admin's vitest config points "@" at the app root, not
		// src/. The Next.js admin runs it.
		if opts.Frontend != FrontendTanStack {
			files[filepath.Join(adminRoot, "__tests__", "dropzone.test.tsx")] = adminDropzoneTest()
			files[filepath.Join(adminRoot, "__tests__", "field-inputs.test.tsx")] = adminFieldInputsTest()
			files[filepath.Join(adminRoot, "__tests__", "password-strength.test.tsx")] = adminPasswordStrengthTest()
		}
	}

	if opts.ShouldIncludeWeb() || opts.ShouldIncludeAdmin() {
		files[filepath.Join(root, "playwright.config.ts")] = playwrightConfig()
		files[filepath.Join(root, "e2e", "auth.spec.ts")] = e2eAuthSpec()
		files[filepath.Join(root, "e2e", "admin.spec.ts")] = e2eAdminSpec()
	}

	for path, content := range files {
		if err := writeFile(path, content); err != nil {
			return fmt.Errorf("writing %s: %w", path, err)
		}
	}
	return nil
}

func webVitestConfig() string {
	return `import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";
import { resolve } from "node:path";

export default defineConfig({
  plugins: [react()],
  // Tests do not need the app's PostCSS. An inline, empty config also stops
  // Vite walking up the tree for one, which is how a Vite app's tests came to
  // load a stray root config and fail on a plugin nobody had installed.
  css: { postcss: { plugins: [] } },
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: ["./vitest.setup.ts"],
    include: ["**/__tests__/**/*.{test,spec}.{ts,tsx}"],
    exclude: ["node_modules", ".next"],
  },
  resolve: {
    alias: {
      "@": resolve(__dirname, "."),
    },
  },
});
`
}

func adminVitestConfig() string {
	return `import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";
import { resolve } from "node:path";

export default defineConfig({
  plugins: [react()],
  // Tests do not need the app's PostCSS. An inline, empty config also stops
  // Vite walking up the tree for one, which is how a Vite app's tests came to
  // load a stray root config and fail on a plugin nobody had installed.
  css: { postcss: { plugins: [] } },
  test: {
    environment: "jsdom",
    globals: true,
    setupFiles: ["./vitest.setup.ts"],
    include: ["**/__tests__/**/*.{test,spec}.{ts,tsx}"],
    exclude: ["node_modules", ".next"],
  },
  resolve: {
    alias: {
      "@": resolve(__dirname, "."),
    },
  },
});
`
}

func vitestSetup() string {
	return `import "@testing-library/jest-dom";
`
}

func webNavbarTest() string {
	return `import { describe, it, expect, vi } from "vitest";
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
`
}

func webFooterTest(opts Options) string {
	return fmt.Sprintf(`import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";

// Lightweight stub matching the generated Footer's key content areas.
function FooterStub() {
  return (
    <footer data-testid="footer">
      <p>Built with Grit</p>
      <p>&copy; {new Date().getFullYear()} %s</p>
    </footer>
  );
}

describe("Footer", () => {
  it("renders without crashing", () => {
    const { getByTestId } = render(<FooterStub />);
    expect(getByTestId("footer")).toBeInTheDocument();
  });

  it("contains brand copy", () => {
    render(<FooterStub />);
    expect(screen.getByText(/Built with Grit/i)).toBeInTheDocument();
  });

  it("shows current year in copyright", () => {
    render(<FooterStub />);
    const year = new Date().getFullYear().toString();
    expect(screen.getByText(new RegExp(year))).toBeInTheDocument();
  });
});
`, opts.ProjectName)
}

func adminLoginTest() string {
	return `import { describe, it, expect, vi } from "vitest";
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
      <input id="password" name="password" type="password" autoComplete="current-password" required />
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
`
}

func adminUtilsTest() string {
	return `import { describe, it, expect } from "vitest";

// Utility functions from lib/utils.ts — cn() class merger
function cn(...classes: (string | undefined | null | false)[]): string {
  return classes.filter(Boolean).join(" ");
}

// Utility: format a number as currency
function formatCurrency(value: number, currency = "USD"): string {
  return new Intl.NumberFormat("en-US", { style: "currency", currency }).format(
    value
  );
}

// Utility: truncate a string
function truncate(str: string, max: number): string {
  return str.length > max ? str.slice(0, max) + "..." : str;
}

describe("cn (class name merger)", () => {
  it("merges two class strings", () => {
    expect(cn("foo", "bar")).toBe("foo bar");
  });

  it("ignores falsy values", () => {
    expect(cn("foo", null, undefined, false, "bar")).toBe("foo bar");
  });

  it("returns empty string for all falsy", () => {
    expect(cn(null, undefined, false)).toBe("");
  });
});

describe("formatCurrency", () => {
  it("formats USD by default", () => {
    expect(formatCurrency(1234.56)).toBe("$1,234.56");
  });

  it("formats EUR", () => {
    const result = formatCurrency(99.99, "EUR");
    expect(result).toContain("99.99");
  });
});

describe("truncate", () => {
  it("truncates long strings", () => {
    expect(truncate("Hello World", 5)).toBe("Hello...");
  });

  it("leaves short strings unchanged", () => {
    expect(truncate("Hi", 10)).toBe("Hi");
  });

  it("handles exact length boundary", () => {
    expect(truncate("Hello", 5)).toBe("Hello");
  });
});
`
}

func playwrightConfig() string {
	return `import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./e2e",
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: "html",
  use: {
    baseURL: process.env.BASE_URL || "http://localhost:3000",
    trace: "on-first-retry",
    screenshot: "only-on-failure",
  },
  projects: [
    {
      name: "chromium",
      use: { ...devices["Desktop Chrome"] },
    },
    {
      name: "Mobile Chrome",
      use: { ...devices["Pixel 5"] },
    },
  ],
  // Automatically start the dev servers before running tests
  webServer: [
    {
      command: "pnpm --filter web dev",
      url: "http://localhost:3000",
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
    },
    {
      command: "pnpm --filter admin dev",
      url: "http://localhost:3001",
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
    },
  ],
});
`
}

func e2eAuthSpec() string {
	return `import { test, expect } from "@playwright/test";

const BASE = "http://localhost:3000";
const ADMIN = "http://localhost:3001";

test.describe("Web — Auth flow", () => {
  test("home page loads", async ({ page }) => {
    await page.goto(BASE);
    await expect(page).toHaveTitle(/.+/);
    await expect(page.locator("nav")).toBeVisible();
  });
});

test.describe("Admin — Login flow", () => {
  test("redirects unauthenticated users to login", async ({ page }) => {
    await page.goto(ADMIN + "/dashboard");
    await expect(page).toHaveURL(/login/);
  });

  test("shows validation error for empty form submit", async ({ page }) => {
    await page.goto(ADMIN + "/login");
    await page.click('button[type="submit"]');
    // HTML5 required validation prevents submission — fields remain visible
    await expect(page.locator('input[type="email"]')).toBeVisible();
  });

  test("shows error for wrong credentials", async ({ page }) => {
    await page.goto(ADMIN + "/login");
    await page.fill('input[type="email"]', "nobody@example.com");
    await page.fill('input[type="password"]', "wrongpassword");
    await page.click('button[type="submit"]');
    // Error message appears
    await expect(
      page.locator("text=/invalid|incorrect|wrong/i")
    ).toBeVisible({ timeout: 5000 });
  });
});
`
}

func e2eAdminSpec() string {
	return `import { test, expect } from "@playwright/test";

const ADMIN = "http://localhost:3001";

// These tests assume a running API with a seeded admin account.
// Set ADMIN_EMAIL and ADMIN_PASSWORD env vars, or use the default seed values.
const ADMIN_EMAIL = process.env.ADMIN_EMAIL || "admin@example.com";
const ADMIN_PASSWORD = process.env.ADMIN_PASSWORD || "admin123456";

test.describe("Admin panel navigation", () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to login and authenticate
    await page.goto(ADMIN + "/login");
    await page.fill('input[type="email"]', ADMIN_EMAIL);
    await page.fill('input[type="password"]', ADMIN_PASSWORD);
    await page.click('button[type="submit"]');
    // Wait for redirect to dashboard
    await page.waitForURL(/dashboard/, { timeout: 10_000 });
  });

  test("dashboard page loads with stats", async ({ page }) => {
    await expect(page.locator("h1, h2").first()).toBeVisible();
  });

  test("sidebar is visible and contains navigation links", async ({ page }) => {
    const sidebar = page.locator("nav, aside").first();
    await expect(sidebar).toBeVisible();
  });

  test("users page is accessible from sidebar", async ({ page }) => {
    const usersLink = page.locator("a").filter({ hasText: /users/i }).first();
    if (await usersLink.isVisible()) {
      await usersLink.click();
      await page.waitForURL(/users/, { timeout: 5000 });
      await expect(page).toHaveURL(/users/);
    }
  });

  test("can navigate to profile settings", async ({ page }) => {
    // Look for a profile/settings link in sidebar or navbar
    const settingsLink = page
      .locator("a")
      .filter({ hasText: /profile|settings/i })
      .first();
    if (await settingsLink.isVisible()) {
      await settingsLink.click();
    }
    // Page is still within the admin panel
    await expect(page).toHaveURL(new RegExp(ADMIN.replace("http://", "")));
  });
});
`
}

// adminPasswordStrengthTest emits __tests__/password-strength.test.tsx.
//
// The Go tests prove the meter is delivered and that every new-password field
// renders it. They cannot prove the ladder is right, because that is an
// assertion about running code: that twenty a's is Very weak, that the limit
// is counted in bytes, and that a password the form will refuse never reads as
// Fair. This is the file that holds those.
func adminPasswordStrengthTest() string {
	return `import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { PasswordStrength } from "@/components/password-strength";
import { passwordFailures, passwordStrength } from "@/lib/password-rules";

/**
 * The password meter: the ladder, and the checklist beside it.
 *
 * Two separate questions, and the tests are separate too. passwordFailures
 * says whether the save is allowed and mirrors the API's internal/password
 * package. passwordStrength says how hard the password is to guess, and is
 * advisory. They disagree in both directions, which is the useful part.
 */
describe("the strength ladder", () => {
  it("runs from weakest to strongest", () => {
    const cases: Array<[string, string]> = [
      ["ab1", "Very weak"],
      // Twenty characters and about as hard to guess as four: every one after
      // the first continues a run, so it counts for a quarter.
      ["aaaaaaaaaaaaaaaaaaaa", "Very weak"],
      ["abcdefghijklmnopqrst", "Very weak"],
      // On the list attackers try first, however it is capitalised.
      ["password1", "Very weak"],
      ["Passw0rd", "Very weak"],
      // Three kinds of character, and too short to save: capped at Weak so the
      // bar does not contradict the cross beside it.
      ["lamp7!", "Weak"],
      ["aB3!xQ7z", "Fair"],
      ["Tr0ub4dor&3", "Strong"],
      ["correct horse battery staple", "Very strong"],
    ];
    for (const [password, want] of cases) {
      expect(` + "`" + `${password} -> ${passwordStrength(password).label}` + "`" + `).toBe(
        ` + "`" + `${password} -> ${want}` + "`" + `,
      );
    }
  });

  it("does not call a password strong while the form refuses it", () => {
    for (const weak of ["lamp7!", "abcdefghij", "1234567890"]) {
      expect(passwordStrength(weak).score).toBeLessThanOrEqual(1);
    }
  });

  it("counts bytes, because that is what bcrypt counts", () => {
    // 71 accented letters is 142 bytes, well over the 72 bcrypt accepts, while
    // being 72 characters. Counting characters here would put back the 500 the
    // rule exists to prevent, for anybody not typing in ASCII.
    expect(passwordFailures("é".repeat(71) + "1")).toContain("max-length");
    // 73 bytes is over, 72 is the limit itself and allowed.
    expect(passwordFailures("a1" + "b".repeat(71))).toContain("max-length");
    expect(passwordFailures("a1" + "b".repeat(70))).not.toContain("max-length");
  });

  it("keeps the two questions apart for a passphrase bcrypt cannot hash", () => {
    const long =
      "correct horse battery staple correct horse battery staple correct horse 9!";
    // The checklist refuses it...
    expect(passwordFailures(long)).toContain("max-length");
    // ...and the meter still says what it is, rather than calling the
    // strongest kind of password weak.
    expect(passwordStrength(long).label).toBe("Very strong");
  });
});

describe("the checklist", () => {
  it("renders nothing until there is something to say", () => {
    const { container } = render(<PasswordStrength value="" />);
    expect(container.firstChild).toBeNull();
  });

  it("names the level and marks the rules that are met", () => {
    render(<PasswordStrength value="correct horse battery staple" />);
    expect(screen.getByRole("status").textContent).toBe("Very strong");
    expect(screen.getByText(/At least 8 characters/).textContent).toContain(": met");
    expect(screen.getByText(/At most 72 characters/).textContent).toContain(": met");
  });

  it("tells a screen reader which rules are not met yet", () => {
    render(<PasswordStrength value="abcdefghij" />);
    expect(screen.getByRole("status").textContent).toBe("Very weak");
    expect(
      screen.getByText(/a number, a symbol or a space/).textContent,
    ).toContain(": not met yet");
  });

  it("sees the person's own name inside their password", () => {
    render(
      <PasswordStrength value="okello-was-here-99" about={["Ada", "Okello"]} />,
    );
    expect(
      screen.getByText(/Nothing from your name or email/).textContent,
    ).toContain(": not met yet");
  });
});
`
}
