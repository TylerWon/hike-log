import { describe, expect, test, vi } from "vitest";

import { readFileAsDataUrl } from "../../../utils/file";

describe("file", () => {
  describe("readFileAsDataUrl", () => {
    const file = new File(["photo"], "photo.png", { type: "image/png" });

    test("throws error when file cannot be read", async () => {
      vi.spyOn(FileReader.prototype, "readAsDataURL").mockImplementation(function (this: FileReader) {
        this.dispatchEvent(new ProgressEvent("error"));
      });

      await expect(readFileAsDataUrl(file)).rejects.toThrow();
    });

    test("returns data URL on success", async () => {
      const url = await readFileAsDataUrl(file);

      const bytes = new Uint8Array(await file.arrayBuffer());
      const binary = Array.from(bytes, (byte) => String.fromCharCode(byte)).join("");
      const base64 = btoa(binary);
      expect(url).toEqual(`data:${file.type};base64,${base64}`);
    });
  });
});
