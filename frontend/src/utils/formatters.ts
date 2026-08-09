/**
 * Takes an ISO date string and outputs it as "[month] [day], [year]" (ex. "Jun 26, 2026").
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
