import { AuditItem } from "@/components/audit-table";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

type DashboardSummary = {
  requestsLast24h: number;
  activeAgents: number;
  connectedTools: number;
  blockedEvents: number;
};

const fallbackSummary: DashboardSummary = {
  requestsLast24h: 0,
  activeAgents: 0,
  connectedTools: 0,
  blockedEvents: 0
};

export async function fetchDashboardSummary(): Promise<DashboardSummary> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/v1/dashboard/summary`, {
      next: { revalidate: 5 }
    });

    if (!response.ok) {
      return fallbackSummary;
    }

    return (await response.json()) as DashboardSummary;
  } catch {
    return fallbackSummary;
  }
}

export async function fetchRecentAuditEvents(): Promise<AuditItem[]> {
  try {
    const response = await fetch(`${API_BASE_URL}/api/v1/audit/recent`, {
      next: { revalidate: 5 }
    });

    if (!response.ok) {
      return [];
    }

    return (await response.json()) as AuditItem[];
  } catch {
    return [];
  }
}
