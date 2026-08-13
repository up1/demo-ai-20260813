"use client";

import { useState, type SubmitEvent } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/useAuthStore";
import { login } from "../services/authService";
import { validateLogin } from "../services/validateLogin";
import type { LoginFormErrors } from "../types";

export default function LoginForm() {
  const router = useRouter();
  const setUser = useAuthStore((state) => state.setUser);

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [fieldErrors, setFieldErrors] = useState<LoginFormErrors>({});
  const [formError, setFormError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();
    setFormError(null);

    const errors = validateLogin(username, password);
    setFieldErrors(errors);
    if (Object.keys(errors).length > 0) {
      return;
    }

    setIsSubmitting(true);
    try {
      const { ok, data } = await login({ username, password });

      if (ok && data.status === "success") {
        setUser({
          userId: data.data.user_id,
          username: data.data.username,
          token: data.data.token,
        });
        router.push("/dashboard");
        return;
      }

      setFormError(
        data.status === "error" ? data.message : "Unable to sign in. Please try again."
      );
    } catch {
      setFormError("Internal server error. Please try again later.");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <div className="w-full max-w-[440px] bg-surface-container-lowest rounded-[0.5rem] shadow-ambient-large border border-[#F1F5F9] p-10 flex flex-col items-center">
      {/* Brand Logo Placeholder */}
      <div className="mb-8 w-16 h-16 rounded-[0.25rem] bg-surface-container flex items-center justify-center text-primary border border-surface-container-highest">
        <span
          className="material-symbols-outlined text-[32px] font-bold"
          style={{ fontVariationSettings: "'FILL' 1" }}
        >
          shield_lock
        </span>
      </div>

      {/* Heading */}
      <div className="text-center mb-8 w-full">
        <h1 className="font-headline text-[24px] leading-[32px] tracking-[-0.01em] font-semibold text-on-background mb-2">
          Sign In to Your Account
        </h1>
        <p className="font-body text-[16px] leading-[24px] text-on-surface-variant">
          Enter your credentials to access the secure portal.
        </p>
      </div>

      {/* Login Form */}
      <form className="w-full flex flex-col gap-4" onSubmit={handleSubmit} noValidate>
        {formError && (
          <div
            role="alert"
            className="w-full rounded-[0.125rem] border border-error-container bg-error-container px-4 py-3 text-[14px] leading-[20px] text-on-error-container"
          >
            {formError}
          </div>
        )}

        {/* Username Input Group */}
        <div className="flex flex-col gap-2">
          <label
            className="font-body text-[14px] leading-[20px] font-semibold text-on-surface"
            htmlFor="username"
          >
            Username or Email
          </label>
          <div
            className={`relative flex items-center w-full border rounded-[0.125rem] bg-surface-container-lowest overflow-hidden transition-all duration-200 input-focused ${
              fieldErrors.username ? "border-error" : "border-outline-variant"
            }`}
          >
            <span className="material-symbols-outlined text-outline ml-4 absolute pointer-events-none">
              person
            </span>
            <input
              data-testid="demo_username"
              className="w-full pl-12 pr-4 py-3 bg-transparent border-none focus:ring-0 font-body text-[16px] leading-[24px] text-on-surface placeholder-on-surface-variant/50 outline-none"
              id="username"
              name="username"
              placeholder="name@company.com"
              type="text"
              autoComplete="username"
              value={username}
              onChange={(event) => setUsername(event.target.value)}
              aria-invalid={Boolean(fieldErrors.username)}
              aria-describedby={fieldErrors.username ? "username-error" : undefined}
            />
          </div>
          {fieldErrors.username && (
            <p id="username-error" className="text-[12px] leading-[16px] text-error">
              {fieldErrors.username}
            </p>
          )}
        </div>

        {/* Password Input Group */}
        <div className="flex flex-col gap-2">
          <label
            className="font-body text-[14px] leading-[20px] font-semibold text-on-surface"
            htmlFor="password"
          >
            Password
          </label>
          <div
            className={`relative flex items-center w-full border rounded-[0.125rem] bg-surface-container-lowest overflow-hidden transition-all duration-200 input-focused ${
              fieldErrors.password ? "border-error" : "border-outline-variant"
            }`}
          >
            <span className="material-symbols-outlined text-outline ml-4 absolute pointer-events-none">
              lock
            </span>
            <input
              data-testid="demo_password"
              className="w-full pl-12 pr-12 py-3 bg-transparent border-none focus:ring-0 font-body text-[16px] leading-[24px] text-on-surface placeholder-on-surface-variant/50 outline-none"
              id="password"
              name="password"
              placeholder="••••••••"
              type={showPassword ? "text" : "password"}
              autoComplete="current-password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              aria-invalid={Boolean(fieldErrors.password)}
              aria-describedby={fieldErrors.password ? "password-error" : undefined}
            />
            <button
              aria-label="Toggle password visibility"
              className="absolute right-4 text-outline hover:text-primary transition-colors focus:outline-none"
              type="button"
              onClick={() => setShowPassword((value) => !value)}
            >
              <span className="material-symbols-outlined">
                {showPassword ? "visibility" : "visibility_off"}
              </span>
            </button>
          </div>
          {fieldErrors.password && (
            <p id="password-error" className="text-[12px] leading-[16px] text-error">
              {fieldErrors.password}
            </p>
          )}
        </div>

        {/* Utilities Row */}
        <div className="flex items-center justify-between mt-2 mb-4 w-full">
          <label className="flex items-center gap-2 cursor-pointer group">
            <input
              className="w-4 h-4 rounded border-outline-variant text-primary focus:ring-primary focus:ring-offset-surface-container-lowest transition-colors"
              type="checkbox"
            />
            <span className="font-body text-[12px] leading-[16px] font-medium text-on-surface-variant group-hover:text-on-surface transition-colors">
              Remember me
            </span>
          </label>
          <a
            className="font-body text-[12px] leading-[16px] font-medium text-primary hover:text-primary-container transition-colors underline-offset-2 hover:underline"
            href="#"
          >
            Forgot password?
          </a>
        </div>

        {/* Primary Action */}
        <button
          data-testid="demo_submit"
          className="w-full bg-primary text-on-primary font-body text-[14px] leading-[20px] font-semibold py-3 px-6 rounded-[0.125rem] transition-all duration-200 btn-hover btn-active flex items-center justify-center gap-2 disabled:opacity-60 disabled:cursor-not-allowed"
          type="submit"
          disabled={isSubmitting}
        >
          {isSubmitting ? "Signing In..." : "Sign In"}
          {!isSubmitting && (
            <span className="material-symbols-outlined text-[18px]">arrow_forward</span>
          )}
        </button>
      </form>

      {/* Secondary Action */}
      <div className="mt-8 text-center w-full pt-6 border-t border-surface-container-high">
        <span className="font-body text-[16px] leading-[24px] text-on-surface-variant">
          New to EnterpriseSecure?
        </span>
        <a
          className="font-body text-[14px] leading-[20px] font-semibold text-primary ml-1 hover:text-primary-container transition-colors hover:underline underline-offset-4"
          href="#"
        >
          Sign up
        </a>
      </div>
    </div>
  );
}
