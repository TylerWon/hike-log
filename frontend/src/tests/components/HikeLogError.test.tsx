import { describe, expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";

import HikeLogError from "../../components/HikeLogError";

describe("HikeLogError", () => {
  test("displays error message and refresh button", async () => {
    const screen = await render(<HikeLogError />);

    const errorMsg = screen.getByRole("heading", { name: "Couldn't load hikes" });
    await expect.element(errorMsg).toBeInTheDocument();

    const refreshButton = screen.getByRole("button", { name: "Retry button" });
    await expect.element(refreshButton).toBeInTheDocument();
  });

  test("refreshes page when try again button is clicked", async () => {
    const onRetryMock = vi.fn();
    const screen = await render(<HikeLogError onRetry={onRetryMock} />);

    const refreshButton = screen.getByRole("button", { name: "Retry button" });
    await refreshButton.click();

    expect(onRetryMock).toHaveBeenCalledOnce();
  });
});
