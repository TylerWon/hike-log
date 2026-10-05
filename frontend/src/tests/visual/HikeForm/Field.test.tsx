import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import Field from "../../../components/HikeForm/Field";

describe("Field", () => {
  test("displays field", async () => {
    const screen = await render(
      <Field label="Name" required>
        <input name="Name" placeholder="John Doe" type="text" />
      </Field>,
    );
    const field = screen.getByRole("group", { name: "Name field" });
    await expect(field).toMatchScreenshot();
  });

  test("displays field with error", async () => {
    const screen = await render(
      <Field error="Invalid name" label="Name" required>
        <input name="Name" placeholder="John Doe" type="text" />
      </Field>,
    );
    const field = screen.getByRole("group", { name: "Name field" });
    await expect(field).toMatchScreenshot();
  });
});
