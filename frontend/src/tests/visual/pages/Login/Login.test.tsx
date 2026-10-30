import { describe, expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";

import { login } from "../../../../api/auth/auth";
import { fillField } from "../../../utils/Field/test-helpers";
import { LoginWrapper } from "../../../utils/Login/test-components";
import { clickLoginButton } from "../../../utils/Login/test-helpers";

vi.mock(import("../../../../api/auth/auth"), async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    login: vi.fn(),
  };
});

describe("Login", () => {
  test("displays password field and login button", async () => {
    const screen = await render(<LoginWrapper />);
    await expect(screen.container).toMatchScreenshot();
  });

  test("displays field error when password is submitted with invalid value", async () => {
    const screen = await render(<LoginWrapper />);
    await clickLoginButton(screen);
    await expect(screen.container).toMatchScreenshot();
  });

  test("displays error message when login fails unexpectedly", async () => {
    vi.mocked(login).mockRejectedValue(new Error("Something went wrong"));

    const screen = await render(<LoginWrapper />);

    await fillField(screen, "Password", "admin");
    await clickLoginButton(screen);
    await expect(screen.container).toMatchScreenshot();
  });

  test("displays disabled login button while login is pending", async () => {
    // Have login sleep for a few seconds so we can check button is disabled
    vi.mocked(login).mockImplementation(() => new Promise((resolve) => setTimeout(resolve, 2000)));

    const screen = await render(<LoginWrapper />);

    await fillField(screen, "Password", "admin");
    await clickLoginButton(screen);
    await expect(screen.container).toMatchScreenshot();
  });
});
