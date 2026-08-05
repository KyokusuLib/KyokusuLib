import { ref, computed, watch, type Ref } from 'vue'
import type { NovelaVolume, NovelaChapter } from '~/types/backend/novela'
import {
  sortVolumesAscending,
  sortChaptersAscending,
  filterChaptersByQuery,
  findVolumeById,
} from '~/utils/chapter-list'

export function useChapterList(volumes: Ref<NovelaVolume[]>) {
  const activeVolumeId = ref(volumes.value[0]?.id ?? '')
  const searchQuery = ref('')

  const sortedVolumes = computed(() => sortVolumesAscending(volumes.value))

  const activeVolume = computed(() => findVolumeById(volumes.value, activeVolumeId.value))

  const filteredChapters = computed<NovelaChapter[]>(() => {
    if (!activeVolume.value) return []
    const sorted = sortChaptersAscending(activeVolume.value.chapters)
    return filterChaptersByQuery(sorted, searchQuery.value)
  })

  const chapterCount = computed(() => {
    return volumes.value.reduce((acc, v) => acc + (v.chapters?.length || 0), 0)
  })

  watch(volumes, (list) => {
    if (!list.length) return
    const stillExists = list.some((v) => v.id === activeVolumeId.value)
    if (!stillExists) {
      activeVolumeId.value = list[0].id
    }
  })

  function setActiveVolume(id: string) {
    activeVolumeId.value = id
  }

  function setSearchQuery(query: string) {
    searchQuery.value = query
  }

  if (volumes.value.length > 0 && !activeVolumeId.value) {
    activeVolumeId.value = volumes.value[0].id
  }

  return {
    activeVolumeId,
    searchQuery,
    sortedVolumes,
    activeVolume,
    filteredChapters,
    chapterCount,
    setActiveVolume,
    setSearchQuery,
  }
}
