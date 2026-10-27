import TextAreaField from "../../../components/form/TextAreaField";

/**
 * Wrapper around a TextAreaField which provides default values for some props.
 */
export function TextAreaFieldWrapper({ error, required }: { error?: string; required?: boolean }) {
  return (
    <TextAreaField
      error={error}
      id="notes"
      label="Notes"
      required={required}
      textAreaProps={{ placeholder: "How was the hike?", rows: 3 }}
    />
  );
}
