interface FieldProps {
  children: React.ReactNode;
  error?: string;
  id: string; // should match id on input
  label: string;
  required?: boolean;
}

export default function Field({ children, error, id, label, required }: FieldProps) {
  return (
    <div aria-label={`${label} field`} className="w-full" role="group">
      <div className="flex font-mono text-[10px] uppercase tracking-widest mb-1.5">
        <label className="text-forest-700" htmlFor={id}>
          {label}
        </label>
        {required && <span className="text-coral-500 ml-1">*</span>}
      </div>

      {children}

      {error && <p className="text-coral-500 text-xs mt-1">{error}</p>}
    </div>
  );
}
