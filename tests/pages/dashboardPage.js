// @ts-check

class DashboardPage {
  /**
   * @param {import('@playwright/test').Page} page
   */
  constructor(page) {
    this.page = page;
    this.welcomeHeading = page.getByTestId("demo_welcome_heading");
  }
}

module.exports = { DashboardPage };
