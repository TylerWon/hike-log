import "../../assets/styles/text.css";
import "../../assets/styles/animation.css";
import "./hike-log-content.css";
import { Fragment } from "react/jsx-runtime";

interface HikeLogContentProps {
  hikeCards: React.ReactElement[];
  overallStatValues: React.ReactElement[];
}

const OVERALL_STAT_LABELS = ["Hikes", "Distance", "Elevation", "Time"];

export default function HikeLogContent({ hikeCards, overallStatValues }: HikeLogContentProps) {
  if (overallStatValues.length !== OVERALL_STAT_LABELS.length) {
    throw new Error("Unexpected number of overall stat values provided");
  }

  return (
    <div className="min-h-screen bg-forest-900 text-cream-100">
      <header className="max-w-3xl mx-auto px-6 pt-16 pb-10">
        <p className="text-xs font-mono tracking-[0.2em] uppercase text-sage-500 mb-3">Field Notes</p>
        <h1 className="text-5xl font-bold font-serif leading-tight mb-6">Hiking Log</h1>

        <div className="flex flex-wrap gap-6">
          {OVERALL_STAT_LABELS.map((label, i) => (
            <Fragment key={i}>
              <div>
                <p className="field-label mb-1">{label}</p>
                {overallStatValues[i]}
              </div>
              {i < OVERALL_STAT_LABELS.length - 1 && <div className="overall-stat-divider" />}
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
