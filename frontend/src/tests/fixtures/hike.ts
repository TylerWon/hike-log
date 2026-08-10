import type { Hike } from "../../schemas/hike";

export const HIKE_FIXTURE: Hike = {
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
      srcUrl: "https://images.unsplash.com/photo-1780945008220-7edb56992450?w=800&h=600&fit=crop&auto=format",
    },
    {
      caption: "Middle lake — flat calm in the morning",
      hikeId: BigInt(1),
      id: BigInt(2),
      srcUrl: "https://images.unsplash.com/photo-1763593862654-52f923fa7d85?w=1200&h=800&fit=crop&auto=format",
    },
    {
      caption: "Upper lake with the Matier Glacier",
      hikeId: BigInt(1),
      id: BigInt(3),
      srcUrl: "https://images.unsplash.com/photo-1723045278368-1ec8fee8d6b6?w=1200&h=800&fit=crop&auto=format",
    },
  ],
  rating: 4.5,
  trailName: "Joffre Lakes",
};
