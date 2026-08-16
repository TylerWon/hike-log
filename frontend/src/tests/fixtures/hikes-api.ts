import { toSerializableJsonObj } from "../../utils/json";
import { HIKE_FIXTURE_1, HIKE_FIXTURE_2 } from "./hike";

export const HIKES_API_RESPONSE_FIXTURE = [
  toSerializableJsonObj(HIKE_FIXTURE_1),
  toSerializableJsonObj(HIKE_FIXTURE_2),
];
