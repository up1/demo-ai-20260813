// @ts-check
const { test, expect } = require("@playwright/test");
const { LoginPage } = require("../pages/loginPage");
const { DashboardPage } = require("../pages/dashboardPage");

test.describe("Login", () => {
  test("user can sign in with valid credentials", async ({ page }) => {
    const loginPage = new LoginPage(page);
    const dashboardPage = new DashboardPage(page);

    await loginPage.goto();
    await loginPage.login("somkiat", "12345678");

    await expect(page).toHaveURL(/\/dashboard$/);
    await expect(dashboardPage.welcomeHeading).toHaveText("Welcome, somkiat");
  });
});
