import { Button } from "@/components/ui/button";

export default function AgentsPage() {
  return (
    <section className="space-y-6">
      <div>
        <h1 className="text-3xl font-semibold tracking-tight">Agent Onboarding</h1>
        <p className="mt-2 text-sm text-[#94A3B8]">
          Connect an AI framework to your gateway endpoint and start routing tool calls with scoped credentials.
        </p>
      </div>

      <div className="rounded-xl border border-[#334155] bg-[rgba(15,23,42,0.72)] p-6">
        <h2 className="text-lg font-medium">Quick Start Snippet</h2>
        <pre className="mt-3 overflow-x-auto rounded-md bg-[#020617] p-4 font-[var(--font-mono)] text-xs text-[#CBD5E1]">
{`# Example: point your agent runtime to the gateway
export NEXUS_GATEWAY_URL=http://localhost:8080
export NEXUS_API_TOKEN=<access-token>
`}
        </pre>
        <div className="mt-4 flex gap-3">
          <Button variant="default">Copy LangChain Snippet</Button>
          <Button variant="secondary">Copy CrewAI Snippet</Button>
        </div>
      </div>
    </section>
  );
}
