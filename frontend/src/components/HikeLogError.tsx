interface HikeLogErrorProps {
  onRetry?: () => void; // Helpful for testing retry button
}

export default function HikeLogError({ onRetry = () => window.location.reload() }: HikeLogErrorProps) {
  return (
    <div
      aria-label="Hike log error"
      className="min-h-screen bg-forest-950 text-cream-100 flex items-center justify-center"
      role="region"
    >
      <div className="flex flex-col items-center justify-center py-24 px-6 text-center">
        <div className="mb-4 relative">
          <svg fill="none" height="48" viewBox="0 0 48 48" width="48" xmlns="http://www.w3.org/2000/svg">
            {/* Mountain */}
            <path
              className="fill-forest-900 stroke-forest-800"
              d="M6 38L18 16l7 10 5-6 12 18H6z"
              strokeLinejoin="round"
              strokeWidth="1.5"
            />
            {/* Crack in mountain */}
            <path
              className="stroke-coral-500"
              d="M22 30l2-3 2 2 2-4"
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="1.5"
            />
            {/* X circle */}
            <circle className="fill-forest-900 stroke-forest-800" cx="37" cy="13" r="7" strokeWidth="1.5" />
            <path d="M34.5 10.5l5 5M39.5 10.5l-5 5" stroke="#c0604a" strokeLinecap="round" strokeWidth="1.5" />
          </svg>
        </div>

        <h2 className="text-xl font-semibold font-serif text-cream-100 mb-2">Couldn&apos;t load hikes</h2>

        <p className="text-sm text-forest-600 max-w-xs leading-relaxed mb-6">
          Something went wrong while fetching the hike log. Check your connection and try again.
        </p>

        <button
          aria-label="Retry button"
          className="inline-flex items-center gap-2 px-5 py-2.5 font-mono text-[12px] text-sage-500 rounded border border-forest-700 hover:bg-forest-850 hover:border-sage-500 hover:text-sage-400 transition-all duration-150 cursor-pointer"
          onClick={onRetry}
        >
          <svg fill="currentColor" height="24" viewBox="0 0 24 24" width="24">
            <path d="M12 20c-2.21665 0 -4.10415 -0.77915 -5.6625 -2.3375C4.779165 16.10415 4 14.21665 4 12c0 -2.21665 0.779165 -4.10415 2.3375 -5.6625C7.89585 4.779165 9.78335 4 12 4c1.41665 0 2.65835 0.2875 3.725 0.8625 1.06665 0.575 1.99165 1.3625 2.775 2.3625V4h1.5v6.35H13.65v-1.5h4.2c-0.63335 -1 -1.44165 -1.80835 -2.425 -2.425C14.44165 5.80835 13.3 5.5 12 5.5c-1.81665 0 -3.35415 0.62915 -4.6125 1.8875C6.12915 8.64585 5.5 10.18335 5.5 12c0 1.81665 0.62915 3.35415 1.8875 4.6125C8.64585 17.87085 10.18335 18.5 12 18.5c1.38335 0 2.65 -0.39585 3.8 -1.1875s1.95 -1.8375 2.4 -3.1375h1.55c-0.48335 1.75 -1.44165 3.15835 -2.875 4.225C15.44165 19.46665 13.81665 20 12 20Z" />
          </svg>
          Try again
        </button>
      </div>
    </div>
  );
}
