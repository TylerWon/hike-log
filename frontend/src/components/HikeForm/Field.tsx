interface FieldProps {
  children: React.ReactNode;
  error?: string;
  label?: string;
  required?: boolean;
}

export default function Field({ children, error, label, required }: FieldProps) {
  return (
    <div className="w-full">
      {label && (
        <label className="font-mono block text-[10px] uppercase tracking-widest text-forest-700 mb-1.5">
          {label}
          {required && <span className="text-coral-500 ml-1">*</span>}
        </label>
      )}
      {children}
      {error && <p className="text-coral-500 text-xs mt-1">{error}</p>}
    </div>
  );
}
