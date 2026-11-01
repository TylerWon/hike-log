import { useState } from "react";

import type { PhotoData, PhotoErrors } from "../../../components/form/PhotoField/types";

import PhotoField from "../../../components/form/PhotoField/PhotoField";

/**
 * Wrapper around a PhotoField that allows photos to be added, removed, and updated.
 */

export function StatefulPhotoField({ errors }: { errors?: PhotoErrors }) {
  const [photos, setPhotos] = useState<PhotoData[]>([]);

  return <PhotoField errors={errors} photos={photos} setPhotos={setPhotos} />;
}
