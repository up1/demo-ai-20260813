import type { LoginRequest, LoginResponse } from "../types";

/**
 * Calls the internal /api/login route, which proxies the request to the
 * authentication backend (see app/api/login/route.ts).
 */
export async function login(credentials: LoginRequest): Promise<{
  ok: boolean;
  data: LoginResponse;
}> {
  const response = await fetch("/api/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(credentials),
  });

  const data = (await response.json()) as LoginResponse;

  return { ok: response.ok, data };
}
