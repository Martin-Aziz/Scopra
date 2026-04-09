import { ReactNode } from "react";

type KpiCardProps = {
  label: string;
  value: number;
  icon: ReactNode;
  tone?: "default" | "success" | "warning";
};

export function KpiCard({ label, value, icon, tone = "default" }: KpiCardProps) {
  const borderColor = {
    default: "border-[#334155]",
    success: "border-[#10B981]/40",
    warning: "border-[#F59E0B]/40"
  }[tone];

  return (
    <article
      className={`rounded-xl border ${borderColor} bg-[rgba(15,23,42,0.72)] p-5 shadow-soft backdrop-blur-md transition-transform duration-300 hover:-translate-y-0.5`}
    >
      <div className="flex items-center justify-between text-sm text-[#94A3B8]">
        <span>{label}</span>
        <span>{icon}</span>
      </div>
      <p className="mt-3 text-3xl font-semibold text-[#F8FAFC]">{value.toLocaleString()}</p>
    </article>
  );
}
