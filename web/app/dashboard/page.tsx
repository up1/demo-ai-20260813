"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/useAuthStore";

export default function DashboardPage() {
  const router = useRouter();
  const user = useAuthStore((state) => state.user);
  const logout = useAuthStore((state) => state.logout);

  useEffect(() => {
    if (!user) {
      router.replace("/login");
    }
  }, [user, router]);

  if (!user) {
    return null;
  }

  return (
    <main className="flex-grow flex items-center justify-center p-4 md:p-10">
      <div className="w-full max-w-[440px] bg-surface-container-lowest rounded-[0.5rem] shadow-ambient-large border border-[#F1F5F9] p-10 flex flex-col items-center text-center gap-4">
        <h1
          data-testid="demo_welcome_heading"
          className="font-headline text-[24px] leading-[32px] font-semibold text-on-background"
        >
          Welcome, {user.username}
        </h1>
        <p className="font-body text-[16px] leading-[24px] text-on-surface-variant">
          You have successfully signed in to EnterpriseSecure.
        </p>
        <button
          className="mt-4 bg-primary text-on-primary font-body text-[14px] leading-[20px] font-semibold py-3 px-6 rounded-[0.125rem] transition-all duration-200 hover:opacity-90"
          onClick={() => {
            logout();
            router.push("/login");
          }}
        >
          Sign Out
        </button>
      </div>
    </main>
  );
}
