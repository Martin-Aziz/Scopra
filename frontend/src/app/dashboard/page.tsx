import { Activity, ShieldCheck, Wrench, Ban } from "lucide-react";
import { KpiCard } from "@/components/kpi-card";
import { AuditTable } from "@/components/audit-table";
import { fetchDashboardSummary, fetchRecentAuditEvents } from "@/services/api";

export const dynamic = "force-dynamic";

export default async function DashboardPage() {
  const [summary, events] = await Promise.all([fetchDashboardSummary(), fetchRecentAuditEvents()]);

  return (
    <section className="space-y-6">
      <div>
        <h1 className="text-3xl font-semibold tracking-tight">Gateway Overview</h1>
        <p className="mt-2 text-sm text-[#94A3B8]">Live routing posture, agent activity, and security signal snapshots.</p>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <KpiCard label="Requests (24h)" value={summary.requestsLast24h} icon={<Activity className="h-4 w-4" />} />
        <KpiCard label="Active Agents" value={summary.activeAgents} icon={<ShieldCheck className="h-4 w-4" />} tone="success" />
        <KpiCard label="Connected Tools" value={summary.connectedTools} icon={<Wrench className="h-4 w-4" />} />
        <KpiCard label="Blocked Events" value={summary.blockedEvents} icon={<Ban className="h-4 w-4" />} tone="warning" />
      </div>

      <div className="space-y-2">
        <h2 className="text-xl font-medium">Recent Audit Activity</h2>
        <AuditTable events={events} />
      </div>
    </section>
  );
}
