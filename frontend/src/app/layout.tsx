import type { Metadata } from "next";
import Link from "next/link";
import { Space_Grotesk, IBM_Plex_Mono } from "next/font/google";
import "./globals.css";

const displayFont = Space_Grotesk({ subsets: ["latin"], variable: "--font-display" });
const monoFont = IBM_Plex_Mono({ subsets: ["latin"], variable: "--font-mono", weight: ["400", "500"] });

export const metadata: Metadata = {
  title: "NEXUS-MCP Dashboard",
  description: "Operational dashboard for gateway routing, audit, and agent health"
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en" className={`${displayFont.variable} ${monoFont.variable}`}>
      <body className="font-[var(--font-display)] text-[#F8FAFC]">
        <header className="border-b border-[#1E293B]/80 bg-[#0B1020]/70 backdrop-blur-xl">
          <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
            <Link href="/dashboard" className="text-lg font-semibold tracking-tight">
              NEXUS-MCP
            </Link>
            <nav className="flex items-center gap-4 text-sm text-[#94A3B8]">
              <Link href="/dashboard" className="hover:text-[#F8FAFC]">
                Dashboard
              </Link>
              <Link href="/audit" className="hover:text-[#F8FAFC]">
                Audit
              </Link>
              <Link href="/agents" className="hover:text-[#F8FAFC]">
                Agents
              </Link>
            </nav>
          </div>
        </header>
        <main className="mx-auto max-w-6xl px-4 py-6 sm:px-6 lg:px-8">{children}</main>
      </body>
    </html>
  );
}
