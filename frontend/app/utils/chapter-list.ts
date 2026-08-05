import type { NovelaVolume, NovelaChapter } from '~/types/backend/novela'

export function sortVolumesAscending(volumes: NovelaVolume[]): NovelaVolume[] {
  return [...volumes].sort((a, b) => a.number - b.number)
}

export function sortChaptersAscending(chapters: NovelaChapter[]): NovelaChapter[] {
  return [...chapters].sort((a, b) => a.number - b.number)
}

export function filterChaptersByQuery(chapters: NovelaChapter[], query: string): NovelaChapter[] {
  if (!query) return chapters
  const lower = query.toLowerCase()
  return chapters.filter(
    (c) => c.title.toLowerCase().includes(lower) || c.number.toString().includes(query),
  )
}

export function findVolumeById(volumes: NovelaVolume[], id: string): NovelaVolume | undefined {
  return volumes.find((v) => v.id === id)
}
