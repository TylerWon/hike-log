import { useState } from "react";

import type { PhotoData, PhotoErrors } from "../../../components/HikeForm/types";

import Field from "../../../components/HikeForm/Field";
import PhotoField from "../../../components/HikeForm/PhotoField";

// Wrapper around a PhotoField that allows photos to be added, removed, and updated.
export function StatefulPhotoField({ errors }: { errors?: PhotoErrors }) {
  const [photos, setPhotos] = useState<PhotoData[]>([]);

  return <PhotoField errors={errors} photos={photos} setPhotos={setPhotos} />;
}

// Wrapper around a Field that sets some default options.
export function FieldWrapper({ error, required }: { error?: string; required?: boolean }) {
  return (
    <Field error={error} id="name" label="Name" required={required}>
      <input id="name" placeholder="John Doe" type="text" />
    </Field>
  );
}
