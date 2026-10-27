import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import { TextAreaFieldWrapper } from "../../../utils/TextAreaField/test-components";

describe("TextAreaField", () => {
  test("displays label and input", async () => {
    const screen = await render(<TextAreaFieldWrapper />);

    const label = screen.getByText("Notes");
    await expect.element(label).toBeInTheDocument();

    const input = screen.getByLabelText("Notes");
    await expect.element(input).toBeInTheDocument();

    const requiredSymbol = screen.getByText("*");
    await expect.element(requiredSymbol).not.toBeInTheDocument();

    const error = screen.getByRole("paragraph");
    await expect.element(error).not.toBeInTheDocument();
  });

  test("displays error when provided", async () => {
    const screen = await render(<TextAreaFieldWrapper error="Invalid" />);

    const label = screen.getByText("Notes");
    await expect.element(label).toBeInTheDocument();

    const input = screen.getByLabelText("Notes");
    await expect.element(input).toBeInTheDocument();

    const requiredSymbol = screen.getByText("*");
    await expect.element(requiredSymbol).not.toBeInTheDocument();

    const error = screen.getByText("Invalid");
    await expect.element(error).toBeInTheDocument();
  });

  test("displays required symbol when specified", async () => {
    const screen = await render(<TextAreaFieldWrapper required />);

    const label = screen.getByText("Notes");
    await expect.element(label).toBeInTheDocument();

    const input = screen.getByLabelText("Notes");
    await expect.element(input).toBeInTheDocument();

    const requiredSymbol = screen.getByText("*");
    await expect.element(requiredSymbol).toBeInTheDocument();

    const error = screen.getByRole("paragraph");
    await expect.element(error).not.toBeInTheDocument();
  });
});
