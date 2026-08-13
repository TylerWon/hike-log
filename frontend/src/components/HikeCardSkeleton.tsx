export default function HikeCardSkeleton() {
  return (
    <div
      aria-label="Hike card skeleton"
      className="rounded-lg border border-forest-800 bg-forest-900 p-4 flex gap-4"
      role="region"
    >
      {/* Thumbnail */}
      <div className="shrink-0 rounded shimmer w-[96px] h-[72px]" />

      <div className="flex-1 min-w-0 flex flex-col justify-between py-0.5">
        <div>
          {/* Trail name */}
          <div className="shimmer rounded h-4 mb-2 w-[55%]" />
          {/* Date */}
          <div className="shimmer rounded h-3 w-[28%]" />
        </div>

        {/* Stats row */}
        <div className="flex items-center gap-4 mt-2.5">
          <div className="shimmer rounded h-3 w-[72px]" />
          <div className="shimmer rounded h-3 w-[48px]" />
          <div className="shimmer rounded h-3 w-[40px]" />
          <div className="shimmer rounded h-3 w-[36px]" />
        </div>
      </div>
    </div>
  );
}
