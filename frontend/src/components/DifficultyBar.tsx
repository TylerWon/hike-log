interface DifficultyBarProps {
  difficulty: number; // 0–10, 0.5 steps
}

const MAX_DIFFICULTY = 10;

export default function DifficultyBar({ difficulty }: DifficultyBarProps) {
  // Clamp to [0, MAX_DIFFICULTY] and snap to the nearest 0.5
  const clampedDifficulty = Math.min(Math.max(Math.round(difficulty * 2) / 2, 0), MAX_DIFFICULTY);

  return (
    <span
      aria-label={`${clampedDifficulty} out of 10 difficulty`}
      className="inline-flex items-center gap-0.5"
      role="img"
    >
      {Array.from({ length: MAX_DIFFICULTY }).map((_, i) => {
        const fill = Math.min(Math.max(clampedDifficulty - i, 0), 1); // 0, 0.5, or 1

        let emptyColor, fillColor;
        if (i < 3) {
          // Easy
          fillColor = "var(--color-sage-500)";
          emptyColor = "var(--color-sage-950)";
        } else if (i < 6) {
          // Moderate
          fillColor = "var(--color-amber-500)";
          emptyColor = "var(--color-amber-950)";
        } else {
          // Hard
          fillColor = "var(--color-coral-500)";
          emptyColor = "var(--color-coral-950)";
        }

        return (
          <span
            aria-label={fill == 0 ? "Empty bar" : fill == 1 ? "Filled bar" : "Half bar"}
            className="inline-block w-[10px] h-[8px] rounded-[2px]"
            key={i}
            role="img"
            // Easier to use inline CSS here than Tailwind
            style={{
              background:
                fill === 0
                  ? emptyColor
                  : fill === 1
                    ? fillColor
                    : `linear-gradient(to right, ${fillColor} 50%, ${emptyColor} 50%)`,
            }}
          />
        );
      })}
    </span>
  );
}
