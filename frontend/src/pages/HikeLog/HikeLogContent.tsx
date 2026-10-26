import StatsOverview, { type Stat } from "./StatsOverview/StatsOverview";

interface HikeLogContentProps {
  children: React.ReactNode;
  overallStats: Stat[];
}

export default function HikeLogContent({ children, overallStats }: HikeLogContentProps) {
  return (
    <div aria-label="Hike log content" className="min-h-screen bg-forest-900 text-cream-100" role="region">
      {/* Title, stats overview */}
      <header className="max-w-3xl mx-auto px-6 pt-16 pb-10">
        <h1 className="text-5xl font-bold font-serif leading-tight mb-6">Hike Log</h1>
        <StatsOverview stats={overallStats} />
        <div className="mt-8 h-px bg-forest-800" />
      </header>

      {/* Body (hike cards) */}
      <main className="max-w-3xl mx-auto px-6 pb-24">
        <ol className="flex flex-col gap-3">{children}</ol>
      </main>
    </div>
  );
}
