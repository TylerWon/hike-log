interface DifficultyBarProps {
  difficulty: number; // 0–10, 0.5 steps
}

function segmentColor(index: number): { filled: string; empty: string } {
  // index 0–2: easy (green), 3–5: moderate (amber), 6–9: hard (red)
  if (index < 3) return { filled: "var(--color-sage-500)", empty: "#1f2d1a" };
  if (index < 6) return { filled: "var(--color-amber-500)", empty: "#2d2510" };
  return { filled: "var(--color-coral-500)", empty: "#2d1510" };
}

export default function DifficultyBar({ difficulty }: DifficultyBarProps) {
  return (
    <span
      className="inline-flex items-center gap-0.5"
      aria-label={`Difficulty ${difficulty} out of 10`}
    >
      {Array.from({ length: 10 }).map((_, i) => {
        const fill = Math.min(Math.max(difficulty - i, 0), 1);
        const { filled, empty } = segmentColor(i);
        const fillPct = fill * 100;
        return (
          <span
            key={i}
            style={{
              width: 10,
              height: 8,
              borderRadius: 2,
              background:
                fill === 0
                  ? empty
                  : fill === 1
                    ? filled
                    : `linear-gradient(to right, ${filled} ${fillPct}%, ${empty} ${fillPct}%)`,
              display: "inline-block",
            }}
          />
        );
      })}
    </span>
  );
}
