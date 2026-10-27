import TextField from "../../../components/form/TextField";

/**
 * Wrapper around a TextField which provides default values for some props.
 */
export function TextFieldWrapper({ error, required }: { error?: string; required?: boolean }) {
  return (
    <TextField
      error={error}
      id="name"
      inputProps={{ placeholder: "e.g. John Doe", type: "text" }}
      label="Name"
      required={required}
    />
  );
}
