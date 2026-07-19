export interface Photo {
  url: string;
  caption?: string;
}

export interface Hike {
  id: string;
  trail: string;
  date: string; // ISO "YYYY-MM-DD"
  rating: number; // 0–5, 0.5 steps
  difficulty: number; // 0–10, 0.5 steps
  distanceKm: number;
  elevationGainM: number;
  durationMinutes: number;
  notes: string;
  allTrailsUrl: string;
  coverPhoto: Photo | null;
  photos: Photo[];
}

export const mockHikeData: Hike[] = [
  {
    id: "panorama-ridge-2026",
    trail: "Panorama Ridge",
    date: "2026-06-28",
    rating: 5,
    difficulty: 8.5,
    distanceKm: 30.2,
    elevationGainM: 1525,
    durationMinutes: 660,
    notes:
      "One of the best days I've had on a trail. Started at 5am to beat the crowds and caught alpenglow on the Black Tusk. The ridge itself is a narrow knife-edge with jaw-dropping views in every direction. Snow on the descent made trekking poles essential. Would not skip this one.",
    allTrailsUrl: "https://www.alltrails.com/trail/canada/british-columbia/panorama-ridge",
    coverPhoto: null,
    photos: [],
  },
  {
    id: "joffre-lakes-2026",
    trail: "Joffre Lakes",
    date: "2026-05-17",
    rating: 4.5,
    difficulty: 5.5,
    distanceKm: 11.8,
    elevationGainM: 370,
    durationMinutes: 240,
    coverPhoto: {
      url: "https://images.unsplash.com/photo-1780945008220-7edb56992450?w=800&h=600&fit=crop&auto=format",
      caption: "Upper Joffre Lake with the glacier above",
    },
    notes:
      "Arrived at the trailhead at 7am and got one of the last parking spots. All three lakes are stunning but Upper Joffre is the showstopper — glacier calving ice directly into the teal water. Light snow made the upper section slippery.",
    allTrailsUrl: "https://www.alltrails.com/trail/canada/british-columbia/joffre-lakes",
    photos: [
      {
        url: "https://images.unsplash.com/photo-1780945008220-7edb56992450?w=800&h=600&fit=crop&auto=format",
        caption: "Upper Joffre Lake with the glacier above",
      },
      {
        url: "https://images.unsplash.com/photo-1763593862654-52f923fa7d85?w=1200&h=800&fit=crop&auto=format",
        caption: "Middle lake — flat calm in the morning",
      },
      {
        url: "https://images.unsplash.com/photo-1723045278368-1ec8fee8d6b6?w=1200&h=800&fit=crop&auto=format",
        caption: "Upper lake with the Matier Glacier",
      },
    ],
  },
  {
    id: "garibaldi-lake-2025",
    trail: "Garibaldi Lake via Rubble Creek",
    date: "2025-09-06",
    rating: 5,
    difficulty: 6.5,
    distanceKm: 18.0,
    elevationGainM: 820,
    durationMinutes: 390,
    coverPhoto: {
      url: "https://images.unsplash.com/photo-1629401681628-a37c83eb57d9?w=800&h=600&fit=crop&auto=format",
      caption: "The trail winds through old-growth forest before opening to the lake",
    },
    notes:
      "Perfect September day — no bugs, cool air, and the lake was an impossible shade of blue. Camped overnight at Battleship Islands. Watching stars reflect in the lake was something I'll remember for a long time.",
    allTrailsUrl: "https://www.alltrails.com/trail/canada/british-columbia/garibaldi-lake",
    photos: [
      {
        url: "https://images.unsplash.com/photo-1629401681628-a37c83eb57d9?w=800&h=600&fit=crop&auto=format",
        caption: "The trail winds through old-growth forest before opening to the lake",
      },
      {
        url: "https://images.unsplash.com/photo-1501554728187-ce583db33af7?w=1200&h=800&fit=crop&auto=format",
        caption: "Heading down through the meadows",
      },
    ],
  },
  {
    id: "grouse-grind-2025",
    trail: "Grouse Grind",
    date: "2025-08-02",
    rating: 3,
    difficulty: 7,
    distanceKm: 5.6,
    elevationGainM: 853,
    durationMinutes: 58,
    coverPhoto: {
      url: "https://images.unsplash.com/photo-1562593028-1fe2d15bde36?w=800&h=600&fit=crop&auto=format",
      caption: "Dense forest corridor on the Grind",
    },
    notes:
      "Finished in 58 minutes — a new PR. The trail itself is relentlessly steep with no real views until the top, but it's an excellent fitness benchmark. Busy on a Saturday afternoon. Took the gondola down.",
    allTrailsUrl: "https://www.alltrails.com/trail/canada/british-columbia/grouse-grind",
    photos: [
      {
        url: "https://images.unsplash.com/photo-1562593028-1fe2d15bde36?w=800&h=600&fit=crop&auto=format",
        caption: "Dense forest corridor on the Grind",
      },
    ],
  },
  {
    id: "black-tusk-2025",
    trail: "Black Tusk via Taylor Meadows",
    date: "2025-07-19",
    rating: 4.5,
    difficulty: 9,
    distanceKm: 28.8,
    elevationGainM: 1585,
    durationMinutes: 720,
    coverPhoto: {
      url: "https://images.unsplash.com/photo-1676782778930-11b311ec5134?w=800&h=600&fit=crop&auto=format",
      caption: "Approaching the Tusk's volcanic summit scramble",
    },
    notes:
      "The summit scramble was genuinely exposed — loose volcanic rock and a long runout below. Beautiful panorama at the top but not for the faint-hearted. Wildflowers in Taylor Meadows were at peak bloom.",
    allTrailsUrl: "https://www.alltrails.com/trail/canada/british-columbia/black-tusk-via-taylor-meadows",
    photos: [
      {
        url: "https://images.unsplash.com/photo-1676782778930-11b311ec5134?w=800&h=600&fit=crop&auto=format",
        caption: "Approaching the Tusk's volcanic summit scramble",
      },
      {
        url: "https://images.unsplash.com/photo-1692802039917-07ce2958bcb1?w=1200&h=800&fit=crop&auto=format",
        caption: "Taylor Meadows blanketed in wildflowers",
      },
      {
        url: "https://images.unsplash.com/photo-1663524963924-4d84fd7204b5?w=1200&h=800&fit=crop&auto=format",
        caption: "Views north toward the Coast Mountains",
      },
    ],
  },
  {
    id: "waterfalls-loop-2025",
    trail: "Shannon Falls & Sea to Summit Connector",
    date: "2025-05-24",
    rating: 3.5,
    difficulty: 2.5,
    distanceKm: 7.4,
    elevationGainM: 260,
    durationMinutes: 165,
    coverPhoto: {
      url: "https://images.unsplash.com/photo-1621810430423-f60c1da430be?w=800&h=600&fit=crop&auto=format",
      caption: "Shannon Falls in full spring flow",
    },
    notes:
      "Good warm-up hike after a long winter. The falls were roaring with snowmelt. Short but satisfying — great option when you only have half a day.",
    allTrailsUrl: "https://www.alltrails.com/trail/canada/british-columbia/shannon-falls",
    photos: [
      {
        url: "https://images.unsplash.com/photo-1621810430423-f60c1da430be?w=800&h=600&fit=crop&auto=format",
        caption: "Shannon Falls in full spring flow",
      },
      {
        url: "https://images.unsplash.com/photo-1632883324493-9caac5ee5371?w=1200&h=800&fit=crop&auto=format",
        caption: "Sunlight filtering through the forest",
      },
      {
        url: "https://images.unsplash.com/photo-1513956767189-5fcd856c2557?w=1200&h=800&fit=crop&auto=format",
        caption: "Old-growth cedars along the connector trail",
      },
    ],
  },
  {
    id: "golden-ears-2024",
    trail: "Golden Ears Summit",
    date: "2024-09-14",
    rating: 4,
    difficulty: 9.5,
    distanceKm: 33.6,
    elevationGainM: 1770,
    durationMinutes: 870,
    coverPhoto: {
      url: "https://images.unsplash.com/photo-1692802039917-07ce2958bcb1?w=800&h=600&fit=crop&auto=format",
      caption: "Above the clouds on the Golden Ears ridge",
    },
    notes:
      "Two-day trip. The summit day was brutal — a very long ridge walk with significant exposure. Fog rolled in just as we topped out so views were minimal, but the experience was memorable. Camp at Alouette Lake was peaceful.",
    allTrailsUrl: "https://www.alltrails.com/trail/canada/british-columbia/golden-ears-summit",
    photos: [
      {
        url: "https://images.unsplash.com/photo-1692802039917-07ce2958bcb1?w=800&h=600&fit=crop&auto=format",
        caption: "Above the clouds on the Golden Ears ridge",
      },
      {
        url: "https://images.unsplash.com/photo-1533240332313-0db49b459ad6?w=1200&h=800&fit=crop&auto=format",
        caption: "Looking back down the ridge through the fog",
      },
    ],
  },
  {
    id: "cypress-howe-sound-2024",
    trail: "Howe Sound Crest Trail",
    date: "2024-08-03",
    rating: 4,
    difficulty: 8,
    distanceKm: 29.5,
    elevationGainM: 1350,
    durationMinutes: 540,
    coverPhoto: {
      url: "https://images.unsplash.com/photo-1501554728187-ce583db33af7?w=800&h=600&fit=crop&auto=format",
      caption: "The crest trail at dusk, looking south over Howe Sound",
    },
    notes:
      "Completed northbound in a single long day. The section around James Peak is the highlight — sweeping views of Howe Sound, the islands, and on a clear day, Vancouver Island. Trail is well-marked but relentlessly up-and-down.",
    allTrailsUrl: "https://www.alltrails.com/trail/canada/british-columbia/howe-sound-crest-trail",
    photos: [
      {
        url: "https://images.unsplash.com/photo-1501554728187-ce583db33af7?w=800&h=600&fit=crop&auto=format",
        caption: "The crest trail at dusk, looking south over Howe Sound",
      },
    ],
  },
];
