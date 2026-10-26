export default function FormError() {
  return (
    <div aria-label="Form error" className="flex items-center gap-2 px-5 py-3 bg-coral-950" role="region">
      <svg
        className="shrink-0"
        fill="none"
        height="14"
        stroke="#c0604a"
        strokeLinecap="round"
        strokeWidth="1.5"
        viewBox="0 0 14 14"
        width="14"
      >
        <circle cx="7" cy="7" r="5.5" />
        <path d="M7 4.5v3M7 9.5v.5" />
      </svg>
      <p className="font-mono flex-1 text-xs text-coral-500">Something went wrong. Please try again.</p>
    </div>
  );
}
