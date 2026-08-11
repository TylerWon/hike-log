import "../../assets/styles/text.css";
import { Fragment } from "react/jsx-runtime";

import "./stats-overview.css";

export interface Stat {
  label: string;
  value: React.ReactNode;
}

interface StatsOverviewProps {
  stats: Stat[];
}

export default function StatsOverview({ stats }: StatsOverviewProps) {
  return (
    <div className="flex flex-wrap gap-6">
      {stats.map((stat, i) => (
        <Fragment key={i}>
          <div>
            <p className="field-label mb-1">{stat.label}</p>
            <div className="stats-overview-stat-value">{stat.value}</div>
          </div>
          {i < stats.length - 1 && <div className="stats-overview-stat-divider" />}
        </Fragment>
      ))}
    </div>
  );
}
