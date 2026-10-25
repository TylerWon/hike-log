import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";

import type { PhotoData, PhotoErrors } from "../../../pages/HikeLog/HikeForm/types";

import Field from "../../../pages/HikeLog/HikeForm/Field";
import HikeForm, { type HikeFormProps } from "../../../pages/HikeLog/HikeForm/HikeForm";
import PhotoField from "../../../pages/HikeLog/HikeForm/PhotoField";

/**
 * Wrapper around a PhotoField that allows photos to be added, removed, and updated.
 */
export function StatefulPhotoField({ errors }: { errors?: PhotoErrors }) {
  const [photos, setPhotos] = useState<PhotoData[]>([]);

  return <PhotoField errors={errors} photos={photos} setPhotos={setPhotos} />;
}

/**
 * Wrapper around a Field that sets some default options.
 */
export function FieldWrapper({ error, required }: { error?: string; required?: boolean }) {
  return (
    <Field error={error} id="name" label="Name" required={required}>
      <input id="name" placeholder="John Doe" type="text" />
    </Field>
  );
}

/**
 * Wraps a HikeForm with a Tanstack Query client for testing.
 */
export function HikeFormWrapper({ onCancel, onSubmit, onSuccess }: HikeFormProps) {
  const [queryClient] = useState(() => new QueryClient());

  return (
    <QueryClientProvider client={queryClient}>
      <HikeForm onCancel={onCancel} onSubmit={onSubmit} onSuccess={onSuccess} />
    </QueryClientProvider>
  );
}
