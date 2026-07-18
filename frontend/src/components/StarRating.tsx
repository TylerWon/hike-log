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
            {/* clipPath controls how much of the filled star is shown (either full, half, or zero width) */}
            <defs>
              <clipPath id={id}>
                <rect x="0" y="0" width={12 * fill} height="12" />
              </clipPath>
            </defs>

            {/* Background: empty star */}
            <polygon
              className="fill-forest-800"
              points="6,1 7.5,4.5 11.5,5 8.5,7.8 9.3,11.5 6,9.5 2.7,11.5 3.5,7.8 0.5,5 4.5,4.5"
            />

            {/* Foreground: filled star (potentially clipped or omitted completely) */}
            {fill > 0 && (
              <polygon
                className="fill-sage-500"
                points="6,1 7.5,4.5 11.5,5 8.5,7.8 9.3,11.5 6,9.5 2.7,11.5 3.5,7.8 0.5,5 4.5,4.5"
                clipPath={`url(#${id})`} // Links this polygon to the clipPath so that its shape is cut according to the clipPath
              />
            )}
          </svg>
        );
      })}
    </span>
  );
}
