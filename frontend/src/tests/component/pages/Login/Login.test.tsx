import { describe, expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";

import { login } from "../../../../network/api/auth/auth";
import { AuthError } from "../../../../network/api/auth/errors";
import { clickButton } from "../../../utils/button";
import { checkFieldError, checkNoFieldError, fillField } from "../../../utils/Field/test-helpers";
import { LoginWrapper } from "../../../utils/Login/test-components";

vi.mock(import("../../../../network/api/auth/auth"), async (importOriginal) => {
  const actual = await importOriginal();
  return {
    ...actual,
    login: vi.fn(),
  };
});

describe("Login", () => {
  test("displays password field and login button", async () => {
    const screen = await render(<LoginWrapper />);

    const password = screen.getByRole("group", { name: "Password field" });
    await expect.element(password).toBeInTheDocument();

    const button = screen.getByRole("button", { name: "Login" });
    await expect.element(button).toBeInTheDocument();
  });

  test("displays field error when empty password is submitted", async () => {
    const screen = await render(<LoginWrapper />);
    await clickButton(screen, "Login");
    await checkFieldError(screen, "Password", "Required");
  });

  test("displays field error when incorrect password is submitted", async () => {
    vi.mocked(login).mockRejectedValue(new AuthError("Incorrect credentials"));

    const screen = await render(<LoginWrapper />);

    await fillField(screen, "Password", "abc");
    await clickButton(screen, "Login");

    await checkFieldError(screen, "Password", "Incorrect password");
  });

  test("displays error message when login fails unexpectedly", async () => {
    vi.mocked(login).mockRejectedValue(new Error("Something went wrong"));

    const screen = await render(<LoginWrapper />);

    await fillField(screen, "Password", "admin");
    await clickButton(screen, "Login");

    await checkNoFieldError(screen, "Password");

    const errorMessage = screen.getByText("Something went wrong. Please try again.");
    await expect.element(errorMessage).toBeInTheDocument();
  });

  test("disables login button while login is pending", async () => {
    // Have login sleep for a few seconds so we can check button is disabled
    vi.mocked(login).mockImplementation(() => new Promise((resolve) => setTimeout(resolve, 2000)));

    const screen = await render(<LoginWrapper />);

    await fillField(screen, "Password", "admin");
    await clickButton(screen, "Login");

    const button = screen.getByRole("button", { name: "Login" });
    await expect.element(button).toBeDisabled();
  });

  test("redirects to home page when login succeeds", async () => {
    vi.mocked(login).mockResolvedValue(undefined);
    const navigateMock = vi.fn();

    const screen = await render(<LoginWrapper navigateFn={navigateMock} />);

    await fillField(screen, "Password", "admin");
    await clickButton(screen, "Login");

    expect(navigateMock).toHaveBeenCalledWith("/");
  });
});
