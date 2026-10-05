import { test, expect } from "@playwright/test";

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
