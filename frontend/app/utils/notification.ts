import type { BackendNotification } from '@/types/backend/notification'

/**
 * Groups notifications by their title field.
 * Notifications without a title fall into "Без категории".
 */
export function groupByTitle(list: BackendNotification[]): Map<string, BackendNotification[]> {
  const groups = new Map<string, BackendNotification[]>()
  for (const n of list) {
    const key = n.title || 'Без категории'
    if (!groups.has(key)) groups.set(key, [])
    groups.get(key)!.push(n)
  }
  return groups
}

export type SortField = 'created_at' | 'title'

/**
 * Sorts grouped notification entries by date or title in asc/desc order.
 * Date sorting uses the latest notification within each group.
 */
export function sortGroupEntries(
  groups: Map<string, BackendNotification[]>,
  field: SortField,
  order: 'asc' | 'desc',
): [string, BackendNotification[]][] {
  const entries = Array.from(groups.entries())
  return entries.sort(([, a], [, b]) => {
    const aLatest = Math.max(...a.map((n) => new Date(n.createdAt).getTime()))
    const bLatest = Math.max(...b.map((n) => new Date(n.createdAt).getTime()))

    if (field === 'title') {
      const cmp = a[0].localeCompare(b[0])
      return order === 'asc' ? cmp : -cmp
    }
    return order === 'asc' ? aLatest - bLatest : bLatest - aLatest
  })
}

/**
 * Toggles an item in a Set (immutable — returns a new Set).
 */
export function toggleSetItem<T>(set: Set<T>, item: T): Set<T> {
  const next = new Set(set)
  if (next.has(item)) {
    next.delete(item)
  } else {
    next.add(item)
  }
  return next
}

/**
 * Returns the number of unread notifications in a list.
 */
export function unreadCountIn(items: BackendNotification[]): number {
  return items.filter((n) => !n.isRead).length
}

/**
 * Formats an ISO date string to a locale-friendly format (ru-RU by default).
 */
export function formatDate(iso: string, locale = 'ru-RU'): string {
  return new Date(iso).toLocaleDateString(locale, {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

/**
 * Pluralization helper: returns the right word form for a given count.
 * @example pluralize(1, ["уведомление", "уведомлений"]) // "уведомление"
 */
export function pluralize(count: number, forms: [string, string]): string {
  return count === 1 ? forms[0] : forms[1]
}
