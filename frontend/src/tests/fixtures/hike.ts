import type { Hike } from "../../schemas/hike";

import grouse_grind_1 from "../assets/images/grouse_grind_1.avif";
import joffre_lakes_1 from "../assets/images/joffre_lakes_1.avif";
import joffre_lakes_2 from "../assets/images/joffre_lakes_2.avif";
import joffre_lakes_3 from "../assets/images/joffre_lakes_3.avif";

export const HIKE_FIXTURE_1: Hike = {
  allTrailsUrl: "https://www.alltrails.com/trail/canada/british-columbia/joffre-lakes",
  date: "2026-05-17",
  difficulty: 5.5,
  distance: 11.8,
  duration: BigInt(240),
  elevationGain: BigInt(370),
  id: BigInt(1),
  notes:
    "Arrived at the trailhead at 7am and got one of the last parking spots. All three lakes are stunning but Upper Joffre is the showstopper — glacier calving ice directly into the teal water. Light snow made the upper section slippery.",
  photos: [
    {
      caption: "Upper Joffre Lake with the glacier above",
      hikeId: BigInt(1),
      id: BigInt(1),
      srcUrl: new URL(joffre_lakes_1, import.meta.url).href,
    },
    {
      caption: "Middle lake — flat calm in the morning",
      hikeId: BigInt(1),
      id: BigInt(2),
      srcUrl: new URL(joffre_lakes_2, import.meta.url).href,
    },
    {
      caption: "Upper lake with the Matier Glacier",
      hikeId: BigInt(1),
      id: BigInt(3),
      srcUrl: new URL(joffre_lakes_3, import.meta.url).href,
    },
  ],
  rating: 4.5,
  trailName: "Joffre Lakes",
};

export const HIKE_FIXTURE_2: Hike = {
  allTrailsUrl: "https://www.alltrails.com/trail/canada/british-columbia/grouse-grind",
  date: "2025-08-02",
  difficulty: 7,
  distance: 5.6,
  duration: BigInt(58),
  elevationGain: BigInt(853),
  id: BigInt(24),
  notes:
    "Finished in 58 minutes — a new PR. The trail itself is relentlessly steep with no real views until the top, but it's an excellent fitness benchmark. Busy on a Saturday afternoon. Took the gondola down.",
  photos: [
    {
      caption: "Dense forest corridor on the Grind",
      hikeId: BigInt(2),
      id: BigInt(2),
      srcUrl: new URL(grouse_grind_1, import.meta.url).href,
    },
  ],
  rating: 3,
  trailName: "Grouse Grind",
};
