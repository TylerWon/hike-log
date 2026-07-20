interface DifficultyBarProps {
  difficulty: number; // 0–10, 0.5 steps
}

export default function DifficultyBar({ difficulty }: DifficultyBarProps) {
  return (
    <span aria-label={`Difficulty ${difficulty} out of 10`} className="inline-flex items-center gap-0.5">
      {Array.from({ length: 10 }).map((_, i) => {
        const fill = Math.min(Math.max(difficulty - i, 0), 1); // 0, 0.5, or 1

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
            className="inline-block w-[10px] h-[8px] rounded-[2px]"
            key={i}
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
