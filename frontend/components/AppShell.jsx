"use client";

import { usePathname } from "next/navigation";
import Sidebar from "./Sidebar";

export default function AppShell({ children }) {
  const pathname = usePathname();
  const isShell = pathname === "/login" || pathname === "/setup";

  if (isShell) return children;

  return (
    <div className="app">
      <Sidebar />
      <main className="main">
        {children}
      </main>
    </div>
  );
}
