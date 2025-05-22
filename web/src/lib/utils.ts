import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

/**
 * Merges and normalizes class names using `clsx` and `tailwind-merge`.
 *
 * This utility function helps combine multiple class names while ensuring
 * Tailwind utility classes are correctly merged to prevent conflicts.
 *
 * @param inputs - A list of class values that can be strings, arrays, or objects.
 * @returns A merged and optimized class string.
 *
 * @example
 * ```tsx
 * cn("bg-red-500", isActive && "text-white", "p-4", "p-2") 
 * // Output: "bg-red-500 text-white p-2" (merges duplicate classes)
 * ```
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
