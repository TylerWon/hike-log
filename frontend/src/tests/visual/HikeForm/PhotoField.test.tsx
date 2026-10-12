import { describe, expect, test } from "vitest";
import { beforeAll } from "vitest";
import { render } from "vitest-browser-react";

import type { PhotoErrors } from "../../../components/HikeForm/types";

import joffreLakes1 from "../../assets/images/joffre_lakes_1.avif";
import joffreLakes2 from "../../assets/images/joffre_lakes_2.avif";
import joffreLakes3 from "../../assets/images/joffre_lakes_3.avif";
import { StatefulPhotoField } from "../../utils/HikeForm/test-components";
import { createPhotoFile, fillFieldInput, uploadPhotos } from "../../utils/HikeForm/test-helpers";

describe("PhotoField", async () => {
  let photos: File[];

  beforeAll(async () => {
    photos = await Promise.all([
      createPhotoFile(joffreLakes1, "joffre_lakes_1.avif"),
      createPhotoFile(joffreLakes2, "joffre_lakes_2.avif"),
      createPhotoFile(joffreLakes3, "joffre_lakes_3.avif"),
    ]);
  });

  test("displays button to upload photos when no photos have been added", async () => {
    const screen = await render(<StatefulPhotoField />);
    const photoField = screen.getByRole("group", { name: "Photo field" });
    await expect(photoField).toMatchScreenshot();
  });

  test("displays uploaded photos", async () => {
    const screen = await render(<StatefulPhotoField />);

    await uploadPhotos(screen, photos);

    const photoField = screen.getByRole("group", { name: "Photo field" });
    await expect(photoField).toMatchScreenshot();
  });

  test("displays caption inputted by user", async () => {
    const screen = await render(<StatefulPhotoField />);

    await uploadPhotos(screen, photos);
    await fillFieldInput(screen, "Caption 1", "TESTING");

    const photoField = screen.getByRole("group", { name: "Photo field" });
    await expect(photoField).toMatchScreenshot();
  });

  test("displays photo errors", async () => {
    const errors: PhotoErrors = {
      errors: [],
      items: [
        { errors: [] },
        {
          errors: [],
          properties: {
            file: { errors: ["Image must be under 10MB"] },
          },
        },
        { errors: [] },
      ],
    };
    const screen = await render(<StatefulPhotoField errors={errors} />);

    await uploadPhotos(screen, photos);

    const photoField = screen.getByRole("group", { name: "Photo field" });
    await expect(photoField).toMatchScreenshot();
  });
});
