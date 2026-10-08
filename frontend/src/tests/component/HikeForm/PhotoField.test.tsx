import { describe, expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";

import type { PhotoErrors } from "../../../components/HikeForm/types";

import { StatefulPhotoField } from "../../utils/HikeForm/photo-field-components";
import { uploadPhotos } from "../../utils/HikeForm/photo-field-helpers";

describe("PhotoField", () => {
  const photo1 = new File(["photo1"], "summit.png", { type: "image/png" });
  const photo2 = new File(["photo2"], "trailhead.jpg", { type: "image/jpeg" });
  const photo3 = new File(["photo3"], "waterfall.jpg", { type: "image/jpeg" });

  test("displays button to upload photos when no photos have been added", async () => {
    const screen = await render(<StatefulPhotoField />);

    const uploadButton = screen.getByRole("button", { name: "Upload photos" });
    await expect.element(uploadButton).toBeInTheDocument();
  });

  test("shows file selector dialog when upload button is clicked", async () => {
    const screen = await render(<StatefulPhotoField />);

    // Use click on file input as indication that file selector dialog will open
    const fileInputClickSpy = vi.spyOn(HTMLInputElement.prototype, "click").mockImplementation(() => {});

    const uploadButton = screen.getByRole("button", { name: "Upload photos" });
    await uploadButton.click();

    expect(fileInputClickSpy).toHaveBeenCalledOnce();
  });

  test("shows file selector dialog when add more photos button is clicked", async () => {
    const screen = await render(<StatefulPhotoField />);

    await uploadPhotos(screen, [photo1]);

    // Use click on file input as indication that file selector dialog will open
    const fileInputClickSpy = vi.spyOn(HTMLInputElement.prototype, "click").mockImplementation(() => {});

    const addMorePhotosButton = screen.getByRole("button", { name: "Add more photos" });
    await addMorePhotosButton.click();

    expect(fileInputClickSpy).toHaveBeenCalledOnce();
  });

  test("displays image preview, caption, and ordering controls for an uploaded photo", async () => {
    const screen = await render(<StatefulPhotoField />);

    await uploadPhotos(screen, [photo1]);

    const previewImage = screen.getByRole("img", { name: "Photo 1" });
    await expect.element(previewImage).toBeInTheDocument();

    const captionInput = screen.getByRole("textbox", { name: "Caption input 1" });
    await expect.element(captionInput).toBeInTheDocument();

    const removeButton = screen.getByRole("button", { name: "Remove photo 1" });
    await expect.element(removeButton).toBeInTheDocument();

    const moveUpButton = screen.getByRole("button", { name: "Move up photo 1" });
    await expect.element(moveUpButton).toBeInTheDocument();

    const moveDownButton = screen.getByRole("button", { name: "Move down photo 1" });
    await expect.element(moveDownButton).toBeInTheDocument();
  });

  test("allows photo caption to be set", async () => {
    const screen = await render(<StatefulPhotoField />);

    await uploadPhotos(screen, [photo1]);

    const captionInput = screen.getByRole("textbox", { name: "Caption input 1" });
    await captionInput.fill("TESTING");
    await expect.element(captionInput).toHaveValue("TESTING");
  });

  test("disabled up button for the first photo and down button for the last photo", async () => {
    const screen = await render(<StatefulPhotoField />);

    await uploadPhotos(screen, [photo1, photo2]);

    const moveUpButtonFirstPhoto = screen.getByRole("button", { name: "Move up photo 1" });
    await expect.element(moveUpButtonFirstPhoto).toBeDisabled();

    const moveDownButtonFirstPhoto = screen.getByRole("button", { name: "Move down photo 1" });
    await expect.element(moveDownButtonFirstPhoto).not.toBeDisabled();

    const moveUpButtonLastPhoto = screen.getByRole("button", { name: "Move up photo 2" });
    await expect.element(moveUpButtonLastPhoto).not.toBeDisabled();

    const moveDownButtonLastPhoto = screen.getByRole("button", { name: "Move down photo 2" });
    await expect.element(moveDownButtonLastPhoto).toBeDisabled();
  });

  test("moves photo up in the order when up button is clicked", async () => {
    const screen = await render(<StatefulPhotoField />);

    await uploadPhotos(screen, [photo1, photo2]);

    // Set caption on photo being moved so it can be used to determine if order changed later
    let captionInput = screen.getByRole("textbox", { name: "Caption input 2" });
    await captionInput.fill("Originally second photo");

    // Move second photo up
    const moveUpButton = screen.getByRole("button", { name: "Move up photo 2" });
    await moveUpButton.click();

    captionInput = screen.getByRole("textbox", { name: "Caption input 1" });
    await expect.element(captionInput).toHaveValue("Originally second photo");
  });

  test("moves photo down in the order when down button is clicked", async () => {
    const screen = await render(<StatefulPhotoField />);

    await uploadPhotos(screen, [photo1, photo2]);

    // Set caption on photo being moved so it can be used to determine if order changed later
    let captionInput = screen.getByRole("textbox", { name: "Caption input 1" });
    await captionInput.fill("Originally first photo");

    // Move first photo down
    const moveDownButton = screen.getByRole("button", { name: "Move down photo 1" });
    await moveDownButton.click();

    captionInput = screen.getByRole("textbox", { name: "Caption input 2" });
    await expect.element(captionInput).toHaveValue("Originally first photo");
  });

  test("removes photo when the delete button is clicked and adjusts relative order of remaining photos", async () => {
    const screen = await render(<StatefulPhotoField />);

    await uploadPhotos(screen, [photo1, photo2, photo3]);

    // Set captions on photos not being deleted so they can be used to determine if order changed later
    let captionInput = screen.getByRole("textbox", { name: "Caption input 2" });
    await captionInput.fill("Originally second photo");

    captionInput = screen.getByRole("textbox", { name: "Caption input 3" });
    await captionInput.fill("Originally third photo");

    // Remove first photo
    const removeButton = screen.getByRole("button", { name: "Remove photo 1" });
    await removeButton.click();

    captionInput = screen.getByRole("textbox", { name: "Caption input 1" });
    await expect.element(captionInput).toHaveValue("Originally second photo");

    captionInput = screen.getByRole("textbox", { name: "Caption input 2" });
    await expect.element(captionInput).toHaveValue("Originally third photo");
  });

  test("removes last photo when the delete button is clicked and shows upload photo buton since no more photos exist", async () => {
    const screen = await render(<StatefulPhotoField />);

    await uploadPhotos(screen, [photo1]);

    const removeButton = screen.getByRole("button", { name: "Remove photo 1" });
    await removeButton.click();

    const uploadButton = screen.getByRole("button", { name: "Upload photos" });
    await expect.element(uploadButton).toBeInTheDocument();
  });

  test("shows errors for correct photos", async () => {
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
      ],
    };

    const screen = await render(<StatefulPhotoField errors={errors} />);

    await uploadPhotos(screen, [photo1, photo2]);

    // Error is displayed as part of the photo's Caption field
    const photo1CaptionField = screen.getByRole("group", { name: "Caption field 1" });
    const photo1Error = photo1CaptionField.getByText("Image must be under 10MB");
    await expect.element(photo1Error).not.toBeInTheDocument();

    const photo2CaptionField = screen.getByRole("group", { name: "Caption field 2" });
    const photo2Error = photo2CaptionField.getByText("Image must be under 10MB");
    await expect.element(photo2Error).toBeInTheDocument();
  });
});
