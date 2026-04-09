import { AuditTable } from "@/components/audit-table";
import { fetchRecentAuditEvents } from "@/services/api";

export const dynamic = "force-dynamic";

export default async function AuditPage() {
  const events = await fetchRecentAuditEvents();

  return (
    <section className="space-y-4">
      <h1 className="text-3xl font-semibold tracking-tight">Audit Ledger</h1>
      <p className="text-sm text-[#94A3B8]">Cryptographically chained event stream with correlation IDs for incident forensics.</p>
      <AuditTable events={events} />
    </section>
  );
}
