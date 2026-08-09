import { describe, expect, test } from "vitest";

import { formatDate, formatDuration } from "../../utils/formatters";

describe("Formatters", () => {
  describe("formatDate", () => {
    test("returns '2026-01-17' as 'Jan 17, 2026'", async () => {
      const date = formatDate("2026-01-17");
      expect(date).toEqual("Jan 17, 2026");
    });
  });

  describe("formatDuration", () => {
    test("returns '0m' when duration is 0 minutes", async () => {
      const duration = formatDuration(BigInt(0));
      expect(duration).toEqual("0m");
    });

    test("returns minutes only when duration is under 60 minutes", async () => {
      const duration = formatDuration(BigInt(59));
      expect(duration).toEqual("59m");
    });

    test("returns hours only when duration is divisible by 60", async () => {
      const duration = formatDuration(BigInt(60));
      expect(duration).toEqual("1h");
    });

    test("returns hours and minutes when duration is over 60 minutes and not divisible by 60", async () => {
      const duration = formatDuration(BigInt(61));
      expect(duration).toEqual("1h 1m");
    });
  });
});
