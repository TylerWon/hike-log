/**
 * Takes an ISO date string (i.e. YYYY-MM-DD) and outputs it as "[month] [day], [year]" (ex. "Jun 26, 2026").
 */
export function formatDate(date: string): string {
  const d = new Date(date + "T12:00:00");
  return d.toLocaleDateString("en-CA", {
    day: "numeric",
    month: "short",
    year: "numeric",
  });
}

/**
 * Takes a distance (in kilometres) and outputs it as "[distance] km" (ex. "51.2 km"). The distance will have 1 decimal
 * point.
 */
export function formatDistance(distanceKm: number): string {
  return `${distanceKm.toFixed(1)} km`;
}

/**
 * Takes a duration (in minutes) and outputs it as "[hours]h [minutes]m" (ex. "5h 20m") or just "[minutes]m" (ex. "58m")
 * if the duration is less than an hour.
 */
export function formatDuration(durationMins: bigint): string {
  const duration = Number(durationMins);
  const h = Math.floor(duration / 60);
  const m = duration % 60;
  if (h === 0) return `${m}m`;
  if (m === 0) return `${h}h`;
  return `${h}h ${m}m`;
}

/**
 * Takes an elevation (in meters) and outputs it as "[elevation] m" (ex. "6,809 m"). Every thousandth place will be
 * separated with a comma.
 */
export function formatElevation(elevationM: bigint): string {
  return `${elevationM.toLocaleString("en-CA")} m`;
}
