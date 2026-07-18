interface StarRatingProps {
  rating: number; // 0–5, 0.5 steps
  max?: number;
}

export default function StarRating({ rating, max = 5 }: StarRatingProps) {
  return (
    <span className="inline-flex items-center gap-0.5" aria-label={`${rating} out of ${max} stars`}>
      {Array.from({ length: max }).map((_, i) => {
        const fill = Math.min(Math.max(rating - i, 0), 1); // 0, 0.5, or 1
        const id = `star-clip-${i}-${rating}`;
        return (
          <svg key={i} width="12" height="12" viewBox="0 0 12 12" fill="none" xmlns="http://www.w3.org/2000/svg">
            <defs>
              <clipPath id={id}>
                <rect x="0" y="0" width={12 * fill} height="12" />
              </clipPath>
            </defs>
            {/* Empty star */}
            <polygon
              points="6,1 7.5,4.5 11.5,5 8.5,7.8 9.3,11.5 6,9.5 2.7,11.5 3.5,7.8 0.5,5 4.5,4.5"
              fill="var(--color-forest-800)"
            />
            {/* Filled star (clipped) */}
            {fill > 0 && (
              <polygon
                points="6,1 7.5,4.5 11.5,5 8.5,7.8 9.3,11.5 6,9.5 2.7,11.5 3.5,7.8 0.5,5 4.5,4.5"
                fill="var(--color-sage-500)"
                clipPath={`url(#${id})`}
              />
            )}
          </svg>
        );
      })}
    </span>
  );
}
