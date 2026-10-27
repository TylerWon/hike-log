import "./form.css";
import { classnames } from "../../utils/classnames";
import Field from "./Field";

interface TextFieldProps {
  error?: string;
  id: string;
  inputProps: React.ComponentProps<"input">;
  label: string;
  required?: boolean;
}

export default function TextField({ error, id, inputProps, label, required }: TextFieldProps) {
  const { className: inputClassName, ...restOfInputProps } = inputProps;

  return (
    <Field error={error} id={id} label={label} required={required}>
      <input className={classnames("form-input", inputClassName)} id={id} {...restOfInputProps} />
    </Field>
  );
}
