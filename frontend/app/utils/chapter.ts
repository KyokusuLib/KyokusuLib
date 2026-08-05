import type { NovelaVolume } from '~/types/backend/novela'

export function getLastChapterNumberForVolume(volumes: NovelaVolume[], volumeId: string): number {
  const vol = volumes.find((v) => v.id === volumeId)
  if (!vol?.chapters?.length) return 1
  let max = 0
  for (const ch of vol.chapters) {
    if (ch.number > max) max = ch.number
  }
  return max + 1
}

export function getLastChapterNumber(volumes: NovelaVolume[]): number {
  if (!volumes.length) return 1
  let max = 0
  for (const v of volumes) {
    if (!v.chapters?.length) continue
    for (const ch of v.chapters) {
      if (ch.number > max) max = ch.number
    }
  }
  return max + 1
}

let idCounter = 0
export function uid(prefix: string): string {
  return `${prefix}-${Date.now()}-${++idCounter}`
}
