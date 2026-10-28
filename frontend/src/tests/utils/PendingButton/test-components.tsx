import PendingButton from "../../../components/PendingButton";

/**
 * Wrapper around a PendingButton which provides default values for some props.
 */
export function PendingButtonWrapper({ pending }: { pending: boolean }) {
  return (
    <PendingButton
      buttonProps={{
        "aria-label": "Submit",
        className: "primary-button-default px-5 py-2",
        type: "submit",
      }}
      pending={pending}
      pendingText="Adding..."
      text="Submit"
    />
  );
}
