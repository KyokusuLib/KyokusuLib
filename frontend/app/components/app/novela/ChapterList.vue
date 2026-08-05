<script setup lang="ts">
import { computed, ref } from 'vue'
import ModalConfirm from '~/components/common/ModalConfirm.vue'
import { useNovela } from '~/composables/api/novela/useNovela'
import { useChapterList } from '~/composables/ui/useChapterList'
import type { NovelaVolume } from '~/types/backend/novela'
import ChapterListContent from './ChapterListContent.vue'
import ChapterSearchHeader from './ChapterSearchHeader.vue'
import NovelaAddVolumeModal from './NovelaAddVolumeModal.vue'
import VolumeTabBar from './VolumeTabBar.vue'

const props = defineProps<{
  volumes: NovelaVolume[]
  canManage?: boolean
  novelaId: number
}>()

const emit = defineEmits<{
  'volume-deleted': [volumeId: string]
  'volume-added': [
    volume: { id: string; title: string; number: number; status?: string } | undefined,
  ]
}>()

const router = useRouter()

const { activeVolumeId, searchQuery, sortedVolumes, filteredChapters, setActiveVolume } =
  useChapterList(computed(() => props.volumes))

const { deleteVolume } = useNovela()
const showAddVolumeModal = ref(false)
const showDeleteVolumeConfirm = ref(false)

const activeVolume = computed(() => sortedVolumes.value.find((v) => v.id === activeVolumeId.value))

function goToAddChapter() {
  router.push(`/novela/${props.novelaId}/add-chapter`)
}

async function confirmDeleteVolume() {
  if (!activeVolume.value) return
  emit('volume-deleted', activeVolume.value.id)
  await deleteVolume(props.novelaId, activeVolume.value.id)
}
</script>

<template>
  <div class="space-y-4">
    <div class="flex flex-col sm:flex-row gap-4 justify-between items-start sm:items-center">
      <VolumeTabBar
        :volumes="sortedVolumes"
        :active-volume-id="activeVolumeId"
        @select="setActiveVolume"
      />

      <ChapterSearchHeader
        v-model="searchQuery"
        :can-manage="canManage"
        @add-volume="showAddVolumeModal = true"
        @add-chapter="goToAddChapter"
      />
    </div>

    <div v-if="canManage && activeVolume" class="flex items-center gap-2 text-xs text-zinc-500">
      <button
        type="button"
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg border border-red-200 dark:border-red-800 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors cursor-pointer"
        @click="showDeleteVolumeConfirm = true"
      >
        <Icon name="ph:trash-bold" size="14" />
        Удалить том {{ activeVolume.number }}
      </button>
    </div>

    <ChapterListContent :chapters="filteredChapters" :novela-id="novelaId" />

    <NovelaAddVolumeModal
      v-if="canManage"
      v-model="showAddVolumeModal"
      :novela-id="novelaId"
      :volumes="sortedVolumes"
      @volume-added="(volume: any) => emit('volume-added', volume)"
    />

    <ModalConfirm
      v-if="activeVolume"
      v-model="showDeleteVolumeConfirm"
      title="Удаление тома"
      :description="`Вы уверены, что хотите удалить том ${activeVolume.number} «${activeVolume.title || 'Без названия'}»? Все главы и изображения этого тома будут безвозвратно удалены.`"
      confirm-text="Удалить"
      cancel-text="Отмена"
      @confirm="confirmDeleteVolume"
    />
  </div>
</template>
