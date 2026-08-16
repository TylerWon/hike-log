/**
 * Converts a JS object to an object that can be serialized with JSON.stringify(). This involves converting BigInts to
 * Numbers as JSON doesn't know how to serialize BigInts.
 */
export function toSerializableJsonObj(value: object): object {
  return JSON.parse(JSON.stringify(value, (_, v) => (typeof v === "bigint" ? Number(v) : v)));
}
