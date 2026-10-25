import "../../assets/styles/dialog.css";
import "../../assets/styles/button.css";

import { useMutation } from "@tanstack/react-query";

import { classnames } from "../../utils/classnames";
import FormError from "./FormError";
import TextWithSpinner from "./TextWithSpinner";

type ActionType = "default" | "destructive";

interface ConfirmationDialogProps {
  acceptButtonPendingText: string;
  acceptButtonText: string;
  actionType: ActionType;
  description: React.ReactElement;
  onAccept: () => Promise<void>;
  onCancel: () => void;
  title: string;
}

export default function ConfirmationDialog({
  acceptButtonPendingText,
  acceptButtonText,
  actionType,
  description,
  onAccept,
  onCancel,
  title,
}: ConfirmationDialogProps) {
  const mutation = useMutation({
    mutationFn: onAccept,
    onError: (error) => console.error(error),
  });

  const acceptPending = mutation.isPending;
  const acceptError = mutation.isError;

  const handleCancel = () => {
    // Don't allow dialog to be cancelled while accept is in progress
    if (!acceptPending) {
      onCancel();
    }
  };

  return (
    <div aria-label="Confirmation dialog backdrop" className="dialog-backdrop" onClick={handleCancel} role="region">
      <div
        aria-label="Confirmation dialog"
        aria-modal="true"
        className="w-full max-w-sm rounded-lg border border-forest-800 bg-forest-900 overflow-hidden"
        onClick={(e) => e.stopPropagation()} // Prevents dialog from closing when it is clicked
        role="alertdialog"
      >
        <div className="flex items-center justify-between px-5 py-4 border-b border-forest-800">
          <h2 className="font-serif text-lg font-semibold text-cream-100">{title}</h2>
        </div>

        <div className="px-6 py-5 flex gap-4 items-start">
          <p className="text-sm text-cream-100 leading-relaxed pt-1.5">{description}</p>
        </div>

        <div className="border-t border-forest-800">
          {acceptError && <FormError />}

          <div className="flex items-center justify-end gap-2 px-5 py-4">
            <button
              aria-label="Cancel"
              className="secondary-button px-4 py-2"
              disabled={acceptPending}
              onClick={handleCancel}
            >
              Cancel
            </button>
            <button
              aria-label="Accept"
              className={classnames("px-5 py-2", {
                "danger-button": actionType === "destructive",
                "primary-button": actionType === "default",
              })}
              disabled={acceptPending}
              onClick={() => mutation.mutate()}
            >
              {acceptPending ? <TextWithSpinner text={acceptButtonPendingText} /> : acceptButtonText}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
