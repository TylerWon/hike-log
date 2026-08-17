import { describe, expect, test } from "vitest";

import { formatDate, formatDistance, formatDuration, formatElevation } from "../../utils/formatters";

describe("formatters", () => {
  describe("formatDate", () => {
    test("returns '2026-01-17' as 'Jan 17, 2026'", async () => {
      const date = formatDate("2026-01-17");
      expect(date).toEqual("Jan 17, 2026");
    });
  });

  describe("formatDistance", () => {
    test("returns distance with 1 decimal point when input has more than 1 decimal point", async () => {
      const distance = formatDistance(2.37);
      expect(distance).toEqual("2.4 km");
    });

    test("returns distance with 1 decmial point when input is a whole number", async () => {
      const distance = formatDistance(16);
      expect(distance).toEqual("16.0 km");
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

  describe("formatElevation", () => {
    test("returns elevation with no commas when the input is below 1000", async () => {
      const elevation = formatElevation(BigInt(999));
      expect(elevation).toEqual("999 m");
    });

    test("returns elevation with every thousandth place separated by a comma when the input is above 1000", async () => {
      const elevation = formatElevation(BigInt(12568));
      expect(elevation).toEqual("12,568 m");
    });
  });
});
