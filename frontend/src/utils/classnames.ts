import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

/**
 * Used to build conditional class strings.
 *
 * clsx handles conditional logic and class composition.
 * tailwind-merge removes conflicting Tailwind classes.
 */
export function classnames(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
