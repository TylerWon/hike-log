import "../../assets/styles/text.css";
import "../../assets/styles/animation.css";
import "./hike-log-content.css";
import { Fragment } from "react/jsx-runtime";

interface HikeLogContentProps {
  hikeCards: React.ReactElement[];
  overallStats: OverallStat[];
}
interface OverallStat {
  label: string;
  value: React.ReactElement;
}

export default function HikeLogContent({ hikeCards, overallStats }: HikeLogContentProps) {
  return (
    <div className="min-h-screen bg-forest-900 text-cream-100">
      <header className="max-w-3xl mx-auto px-6 pt-16 pb-10">
        <p className="text-xs font-mono tracking-[0.2em] uppercase text-sage-500 mb-3">Field Notes</p>
        <h1 className="text-5xl font-bold font-serif leading-tight mb-6">Hiking Log</h1>

        <div className="flex flex-wrap gap-6">
          {overallStats.map((stat, i) => (
            <Fragment key={i}>
              <div>
                <p className="field-label mb-1">{stat.label}</p>
                {stat.value}
              </div>
              {i < overallStats.length - 1 && <div className="overall-stat-divider" />}
            </Fragment>
          ))}
        </div>

        <div className="mt-8 h-px bg-forest-800" />
      </header>

      <main className="max-w-3xl mx-auto px-6 pb-24">
        <ol className="flex flex-col gap-3">
          {hikeCards.map((card, i) => (
            <li key={i}>{card}</li>
          ))}
        </ol>
      </main>
    </div>
  );
}
