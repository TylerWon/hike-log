interface TextWithSpinnerProps {
  text: string;
}

export default function TextWithSpinner({ text }: TextWithSpinnerProps) {
  return (
    <div className="inline-flex items-center gap-2">
      <svg
        className="animate-spin"
        fill="none"
        height="12"
        stroke="currentColor"
        strokeLinecap="round"
        strokeWidth="1.5"
        viewBox="0 0 12 12"
        width="12"
      >
        <path d="M6 1v2M6 9v2M1 6h2M9 6h2M2.5 2.5l1.4 1.4M8.1 8.1l1.4 1.4M9.5 2.5L8.1 3.9M3.9 8.1L2.5 9.5" />
      </svg>
      {text}
    </div>
  );
}
