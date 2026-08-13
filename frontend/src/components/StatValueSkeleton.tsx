interface StatValueSkeletonProps {
  widthClass: string;
}

export default function StatValueSkeleton({ widthClass }: StatValueSkeletonProps) {
  return <div aria-label="Statistic skeleton" className={`shimmer rounded h-7 ${widthClass}`} role="region" />;
}
