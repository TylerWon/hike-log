interface StarRatingProps {
  rating: number;
}

export default function StarRating({ rating }: StarRatingProps) {
  // Clamp to [0, 5] and snap to the nearest 0.5
  const clampedRating = Math.min(Math.max(Math.round(rating * 2) / 2, 0), 5);

  return (
    <span aria-label={`${clampedRating} out of 5 stars`} className="inline-flex items-center gap-0.5" role="img">
      {Array.from({ length: 5 }).map((_, i) => {
        const fill = Math.min(Math.max(clampedRating - i, 0), 1); // 0, 0.5, or 1
        const id = `star-${i}-${clampedRating}`;

        return (
          <svg
            aria-label={fill === 1 ? "Filled star" : fill === 0.5 ? "Half star" : "Empty star"}
            fill="none"
            height="12"
            key={i}
            role="img"
            viewBox="0 0 12 12"
            width="12"
          >
            {/* clipPath controls how much of the filled star is shown (either all, half, or none) */}
            <defs>
              <clipPath id={id}>
                <rect height="12" width={12 * fill} x="0" y="0" />
              </clipPath>
            </defs>

            {/* Background: empty star */}
            <polygon
              className="fill-forest-800"
              points="6,1 7.5,4.5 11.5,5 8.5,7.8 9.3,11.5 6,9.5 2.7,11.5 3.5,7.8 0.5,5 4.5,4.5"
            />

            {/* Foreground: filled star */}
            {fill > 0 && (
              <polygon
                className="fill-sage-500"
                clipPath={`url(#${id})`} // Links this polygon to the clipPath so that its shape is cut according to the clipPath
                points="6,1 7.5,4.5 11.5,5 8.5,7.8 9.3,11.5 6,9.5 2.7,11.5 3.5,7.8 0.5,5 4.5,4.5"
              />
            )}
          </svg>
        );
      })}
    </span>
  );
}
