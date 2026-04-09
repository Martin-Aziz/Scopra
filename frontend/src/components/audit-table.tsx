export type AuditItem = {
  id: number;
  agentId: string;
  tool: string;
  action: string;
  status: string;
  correlationId: string;
  createdAt: string;
};

type AuditTableProps = {
  events: AuditItem[];
};

export function AuditTable({ events }: AuditTableProps) {
  if (events.length === 0) {
    return (
      <div className="rounded-xl border border-[#334155] bg-[rgba(15,23,42,0.72)] p-6 text-[#94A3B8]">
        No audit events yet. Execute a tool call through the gateway to populate the stream.
      </div>
    );
  }

  return (
    <div className="overflow-hidden rounded-xl border border-[#334155] bg-[rgba(15,23,42,0.72)]">
      <table className="w-full text-left text-sm text-[#E2E8F0]">
        <thead className="bg-[#0F172A] text-xs uppercase text-[#94A3B8]">
          <tr>
            <th className="px-4 py-3">Timestamp</th>
            <th className="px-4 py-3">Agent</th>
            <th className="px-4 py-3">Tool</th>
            <th className="px-4 py-3">Action</th>
            <th className="px-4 py-3">Status</th>
            <th className="px-4 py-3">Correlation</th>
          </tr>
        </thead>
        <tbody>
          {events.map((event) => (
            <tr key={event.id} className="border-t border-[#1E293B] hover:bg-[#1E293B]/40">
              <td className="px-4 py-3">{new Date(event.createdAt).toLocaleString()}</td>
              <td className="px-4 py-3 font-mono text-xs">{event.agentId.slice(0, 12)}...</td>
              <td className="px-4 py-3">{event.tool}</td>
              <td className="px-4 py-3">{event.action}</td>
              <td className="px-4 py-3">
                <span className="rounded bg-[#1E293B] px-2 py-1 text-xs">{event.status}</span>
              </td>
              <td className="px-4 py-3 font-mono text-xs text-[#94A3B8]">{event.correlationId}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
