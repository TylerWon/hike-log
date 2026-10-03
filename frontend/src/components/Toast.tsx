import { useEffect } from "react";

interface ToastProps {
  message: string;
  onDismiss: () => void;
  duration?: number; // ms, default 5000
}

export default function Toast({ message, onDismiss, duration = 10000 }: ToastProps) {
  useEffect(() => {
    const timer = setTimeout(onDismiss, duration);
    return () => clearTimeout(timer);
  }, [onDismiss, duration]);

  return (
    <div
      className="fixed bottom-6 left-1/2 -translate-x-1/2 z-50 flex items-center gap-3 px-4 py-3 rounded-lg border border-amber-700 bg-amber-950 shadow-lg min-w-[280px] max-w-[420px]"
      role="status"
      aria-live="polite"
    >
      <svg
        width="16"
        height="16"
        viewBox="0 0 16 16"
        fill="none"
        stroke="#c09a4a"
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
        className="flex-shrink-0 mt-px"
      >
        <path d="M8 2L14.5 13H1.5L8 2z" />
        <path d="M8 6.5v3M8 11v.5" />
      </svg>

      <p className="font-mono flex-1 text-xs leading-relaxed text-amber-400">
        {message}
      </p>

      <button
        onClick={onDismiss}
        className="flex-shrink-0 text-amber-700 hover:text-amber-500 transition-colors focus:outline-none mt-px"
        aria-label="Dismiss"
      >
        <svg width="12" height="12" viewBox="0 0 12 12" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round">
          <path d="M2 2l8 8M10 2L2 10" />
        </svg>
      </button>
    </div>
  );
}
