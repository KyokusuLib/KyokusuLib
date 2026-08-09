<script setup lang="ts">
import { computed } from 'vue'
import { Card, Button, Badge } from '@kyokusu-ui/vue'
import type { PendingVolume, PendingChapter } from '@/composables/api/moderation/useModeration'

const props = defineProps<{
  kind: 'volume' | 'chapter'
  item: PendingVolume | PendingChapter
  loading?: boolean
}>()

const emit = defineEmits<{
  approve: []
  reject: []
  edit: []
}>()

const isVolume = computed(() => props.kind === 'volume')
const isChapter = computed(() => props.kind === 'chapter')

const volumeItem = computed<PendingVolume | null>(() =>
  isVolume.value ? (props.item as PendingVolume) : null,
)
const chapterItem = computed<PendingChapter | null>(() =>
  isChapter.value ? (props.item as PendingChapter) : null,
)

const formatChapterNumber = (num: number) => {
  return Number.isInteger(num) ? num : num.toFixed(1)
}

const novelaTitle = computed(() => props.item.novela_title || 'Без названия')
const author = computed(() => props.item.created_by_name || `#${props.item.created_by}`)
const itemBadge = computed(() =>
  isVolume.value
    ? `Том ${formatChapterNumber(props.item.number)}`
    : `Глава ${formatChapterNumber(props.item.number)}`,
)

// Текст главы без HTML и <figure>-блоков (SSR-safe, без DOM)
const chapterText = computed(() => {
  const raw = chapterItem.value?.content || ''
  if (!raw) return ''
  return raw
    .replace(/<figure[\s\S]*?<\/figure>/gi, ' ')
    .replace(/<[^>]+>/g, ' ')
    .replace(/&nbsp;/g, ' ')
    .replace(/&amp;/g, '&')
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>')
    .replace(/&quot;/g, '"')
    .replace(/&#0*39;/g, "'")
    .replace(/\s+/g, ' ')
    .trim()
})

// Количество изображений в контенте главы
const imageCount = computed(() => {
  const raw = chapterItem.value?.content || ''
  const figures = (raw.match(/<figure[\s\S]*?<\/figure>/gi) || []).length
  const imgs = (raw.match(/<img\b/gi) || []).length
  return Math.max(figures, imgs)
})

const imageBadge = computed(() => {
  if (imageCount.value === 0) return ''
  return imageCount.value === 1 ? 'С изображением' : 'С изображениями'
})

const accentClass = computed(() =>
  isVolume.value ? 'hover:border-yellow-500/40' : 'hover:border-blue-500/40',
)
</script>

<template>
  <Card
    variant="default"
    padding="md"
    shadow
    class="cursor-pointer transition-all hover:shadow-md group min-w-0"
    :class="accentClass"
    @click="emit('edit')"
  >
    <template #header>
      <div class="flex justify-between items-start gap-3">
        <div class="min-w-0">
          <p class="text-xs font-bold uppercase tracking-wider text-amber-500 dark:text-amber-400">
            Новелла
          </p>
          <h3 class="font-bold text-zinc-900 dark:text-white truncate">
            {{ novelaTitle }}
          </h3>
        </div>
        <Badge variant="secondary" size="sm" :text="itemBadge" />
      </div>
    </template>

    <!-- Volume body -->
    <template v-if="volumeItem">
      <div class="flex flex-col gap-3">
        <p v-if="volumeItem.title" class="text-sm text-zinc-600 dark:text-zinc-300">
          {{ volumeItem.title }}
        </p>
        <div class="flex items-center gap-2">
          <Icon name="ph:user-bold" size="14" class="text-zinc-400" />
          <span class="text-xs text-zinc-500 dark:text-zinc-400 truncate">{{ author }}</span>
        </div>
        <div class="flex items-center gap-2 flex-wrap">
          <Badge variant="outline" size="sm" text="Ожидает" />
          <span class="text-[10px] text-zinc-400">нажмите для редактирования</span>
        </div>
      </div>
    </template>

    <!-- Chapter body -->
    <template v-else-if="chapterItem">
      <div class="flex flex-col gap-3">
        <div class="flex items-center gap-2 text-sm">
          <Icon name="ph:books-bold" size="14" class="text-zinc-400" />
          <span class="font-semibold text-zinc-700 dark:text-zinc-300">
            Том {{ formatChapterNumber(chapterItem.volume_number) }}
          </span>
          <span v-if="chapterItem.title" class="text-zinc-500 dark:text-zinc-400 truncate">
            — {{ chapterItem.title }}
          </span>
        </div>

        <p
          class="text-xs text-zinc-500 dark:text-zinc-400 bg-zinc-50 dark:bg-zinc-800/50 p-3 rounded-xl line-clamp-3 wrap-break-words"
        >
          {{ chapterText || 'Без текста' }}
        </p>

        <div class="flex items-center gap-2 flex-wrap">
          <Badge variant="outline" size="sm" text="Ожидает" />
          <Badge v-if="imageBadge" variant="outline" bg="#f59e0b" size="sm" style="color: black" :text="imageBadge" />
          <span class="text-[10px] text-zinc-400">Нажмите для редактирования</span>
        </div>

        <div class="flex items-center gap-2">
          <Icon name="ph:user-bold" size="14" class="text-zinc-400" />
          <span class="text-xs text-zinc-500 dark:text-zinc-400 truncate">{{ author }}</span>
        </div>
      </div>
    </template>

    <template #footer>
      <div class="grid grid-cols-2 gap-1.5 pt-2 border-t border-zinc-100 dark:border-zinc-800">
        <button
          class="flex items-center justify-center py-3 gap-2 h-4 rounded-2xl bg-zinc-100 hover:bg-zinc-300 text-black dark:text-zinc-700 cursor-pointer"
          @click.stop="emit('approve')"
        >
          <Icon name="ph:check-bold" size="12" class="text-black dark:text-zinc-700 mt-0.5" />
          <span class="text-xs font-bold">Одобрить</span>
        </button>

        <button
          class="flex items-center justify-center py-3 gap-2 h-4 rounded-2xl bg-red-500 hover:bg-red-600 text-black dark:text-white cursor-pointer"
          @click.stop="emit('reject')"
        >
          <Icon name="ph:x-bold" size="12" class="text-black dark:text-zinc-200 mt-0.5" />
          <span class="text-xs font-bold">Отклонить</span>
        </button>
        
      </div>
    </template>
  </Card>
</template>

<style scoped>
:deep(.k-button-content) {
  gap: 6px;
}
</style>
