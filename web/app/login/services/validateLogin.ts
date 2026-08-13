import type { LoginFormErrors } from "../types";

export function validateLogin(
  username: string,
  password: string
): LoginFormErrors {
  const errors: LoginFormErrors = {};

  if (!username.trim()) {
    errors.username = "Please enter a valid email address.";
  }

  if (!password || password.length < 8) {
    errors.password = "Password must be at least 8 characters long.";
  }

  return errors;
}
