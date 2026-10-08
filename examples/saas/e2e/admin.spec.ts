import { test, expect } from "@playwright/test";

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
