import { NextResponse } from "next/server";

const API_HOST = process.env.API_HOST ?? "http://localhost:8000";

export async function POST(request: Request) {
  const body = await request.json();

  try {
    const backendResponse = await fetch(`${API_HOST}/api/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });

    const data = await backendResponse.json();
    return NextResponse.json(data, { status: backendResponse.status });
  } catch {
    return NextResponse.json(
      { status: "error", message: "Internal server error" },
      { status: 500 }
    );
  }
}
