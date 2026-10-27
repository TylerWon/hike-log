import "./form.css";
import { classnames } from "../../utils/classnames";
import Field from "./Field";

interface TextAreaProps {
  error?: string;
  id: string;
  label: string;
  required?: boolean;
  textAreaProps: React.ComponentProps<"textarea">;
}

export default function TextAreaField({ error, id, label, required, textAreaProps }: TextAreaProps) {
  const { className: textAreaClassName, ...restOfTextAreaProps } = textAreaProps;

  return (
    <Field error={error} id={id} label={label} required={required}>
      <textarea className={classnames("form-input", textAreaClassName)} id={id} {...restOfTextAreaProps} />
    </Field>
  );
}
