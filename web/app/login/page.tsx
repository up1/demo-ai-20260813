import type { Metadata } from "next";
import LoginForm from "./components/LoginForm";

export const metadata: Metadata = {
  title: "EnterpriseSecure - Sign In",
};

export default function LoginPage() {
  return (
    <main className="flex-grow flex items-center justify-center p-4 md:p-10">
      <LoginForm />
    </main>
  );
}
