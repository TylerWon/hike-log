import { useEffect } from "react";

interface ToastProps {
  duration?: number; // ms, default 5000
  message: string;
  onDismiss: () => void;
}

export default function Toast({ duration = 10000, message, onDismiss }: ToastProps) {
  useEffect(() => {
    const timer = setTimeout(onDismiss, duration);
    return () => clearTimeout(timer);
  }, [onDismiss, duration]);

  return (
    <div
      aria-live="polite"
      className="fixed bottom-6 left-1/2 -translate-x-1/2 z-50 flex items-center gap-3 px-4 py-3 rounded-lg border border-amber-700 bg-amber-950 shadow-lg min-w-[280px] max-w-[420px]"
      role="status"
    >
      <svg
        className="flex-shrink-0 mt-px"
        fill="none"
        height="16"
        stroke="#c09a4a"
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="1.5"
        viewBox="0 0 16 16"
        width="16"
      >
        <path d="M8 2L14.5 13H1.5L8 2z" />
        <path d="M8 6.5v3M8 11v.5" />
      </svg>

      <p className="font-mono flex-1 text-xs leading-relaxed text-amber-400">{message}</p>

      <button
        aria-label="Dismiss"
        className="flex-shrink-0 text-amber-700 hover:text-amber-500 transition-colors focus:outline-none mt-px"
        onClick={onDismiss}
      >
        <svg
          fill="none"
          height="12"
          stroke="currentColor"
          strokeLinecap="round"
          strokeWidth="1.5"
          viewBox="0 0 12 12"
          width="12"
        >
          <path d="M2 2l8 8M10 2L2 10" />
        </svg>
      </button>
    </div>
  );
}
