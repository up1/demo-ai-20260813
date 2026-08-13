// @ts-check

class LoginPage {
  /**
   * @param {import('@playwright/test').Page} page
   */
  constructor(page) {
    this.page = page;
    this.usernameInput = page.getByTestId("demo_username");
    this.passwordInput = page.getByTestId("demo_password");
    this.submitButton = page.getByTestId("demo_submit");
    this.formError = page.getByRole("alert");
  }

  async goto() {
    await this.page.goto("/login");
  }

  async login(username, password) {
    await this.usernameInput.fill(username);
    await this.passwordInput.fill(password);
    await this.submitButton.click();
  }
}

module.exports = { LoginPage };
