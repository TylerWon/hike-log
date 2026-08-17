import { describe, expect, test } from "vitest";

import { toSerializableJsonObj } from "../../../utils/json";

describe("json", () => {
  describe("toSerializableJsonObj", () => {
    test("returns a copy of original object when it is already serializable", async () => {
      const input = {
        category: "Sports Equipment",
        id: 1,
        manufacturer: {
          id: 100,
          name: "Wilson",
          origin: "Chicago, Illinois",
        },
        name: "Baseball Glove",
        outOfStock: false,
        price: 99.99,
        relatedItems: ["Baseball Bat", "Baseball", "Baseball Helmet"],
      };

      const result = toSerializableJsonObj(input);

      expect(result).toEqual(input);
    });

    test("converts BigInts to Numbers", async () => {
      const input = {
        artist: {
          id: BigInt(200),
          name: "Michael Jackson",
        },
        id: BigInt(100),
        releaseDate: "January 2, 1983",
        streams: BigInt(1000000000),
        title: "Billie Jean",
      };
      const expected = {
        artist: {
          id: 200,
          name: "Michael Jackson",
        },
        id: 100,
        releaseDate: "January 2, 1983",
        streams: 1000000000,
        title: "Billie Jean",
      };

      const result = toSerializableJsonObj(input);

      expect(result).toEqual(expected);
    });
  });
});
