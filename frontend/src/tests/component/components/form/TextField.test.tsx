import { describe, expect, test } from "vitest";
import { render } from "vitest-browser-react";

import { TextFieldWrapper } from "../../../utils/TextField/test-components";

describe("TextField", () => {
  test("displays label and input", async () => {
    const screen = await render(<TextFieldWrapper />);

    const label = screen.getByText("Name");
    await expect.element(label).toBeInTheDocument();

    const input = screen.getByLabelText("Name");
    await expect.element(input).toBeInTheDocument();

    const requiredSymbol = screen.getByText("*");
    await expect.element(requiredSymbol).not.toBeInTheDocument();

    const error = screen.getByRole("paragraph");
    await expect.element(error).not.toBeInTheDocument();
  });

  test("displays error when provided", async () => {
    const screen = await render(<TextFieldWrapper error="Invalid" />);

    const label = screen.getByText("Name");
    await expect.element(label).toBeInTheDocument();

    const input = screen.getByLabelText("Name");
    await expect.element(input).toBeInTheDocument();

    const requiredSymbol = screen.getByText("*");
    await expect.element(requiredSymbol).not.toBeInTheDocument();

    const error = screen.getByText("Invalid");
    await expect.element(error).toBeInTheDocument();
  });

  test("displays required symbol when specified", async () => {
    const screen = await render(<TextFieldWrapper required />);

    const label = screen.getByText("Name");
    await expect.element(label).toBeInTheDocument();

    const input = screen.getByLabelText("Name");
    await expect.element(input).toBeInTheDocument();

    const requiredSymbol = screen.getByText("*");
    await expect.element(requiredSymbol).toBeInTheDocument();

    const error = screen.getByRole("paragraph");
    await expect.element(error).not.toBeInTheDocument();
  });
});
