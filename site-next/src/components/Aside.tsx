interface AsideProps {
  type?: "note" | "tip" | "caution" | "danger";
  title?: string;
  children: React.ReactNode;
}

const styles: Record<string, { border: string; bg: string; icon: string }> = {
  note: {
    border: "border-accent/30",
    bg: "bg-accent/5",
    icon: "text-accent",
  },
  tip: {
    border: "border-[#4ade80]/30",
    bg: "bg-[#4ade80]/5",
    icon: "text-[#4ade80]",
  },
  caution: {
    border: "border-[#fbbf24]/30",
    bg: "bg-[#fbbf24]/5",
    icon: "text-[#fbbf24]",
  },
  danger: {
    border: "border-[#f87171]/30",
    bg: "bg-[#f87171]/5",
    icon: "text-[#f87171]",
  },
};

export default function Aside({
  type = "note",
  title,
  children,
}: AsideProps) {
  const style = styles[type] || styles.note;
  const label = title || type.charAt(0).toUpperCase() + type.slice(1);

  return (
    <div
      className={`my-4 rounded-lg border-l-4 ${style.border} ${style.bg} p-4`}
    >
      <p className={`text-sm font-semibold mb-1 ${style.icon}`}>{label}</p>
      <div className="text-sm text-gray-2">{children}</div>
    </div>
  );
}
