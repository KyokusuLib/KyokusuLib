<script setup lang="ts">
import { ref, computed } from 'vue'

import { TeleportedTooltip } from '@kyokusu-ui/vue'

import type { TotalNovelaStatistics } from '@/types/backend/statistics'

import { staticImage } from '@/utils/str'
import { roundTo } from '@/utils/num'

const props = defineProps<{
  items: TotalNovelaStatistics[]
  offset: number
}>()

type SortField =
  | 'rating'
  | 'bookmarkCount'
  | 'readCount'
  | 'ratingCount'
  | 'commentCount'
  | 'volumeCount'
  | 'chapterCount'
  | 'title'

const sortField = ref<SortField | null>(null)
const sortDir = ref<'asc' | 'desc'>('desc')

function setSort(field: SortField) {
  if (sortField.value === field) {
    sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortField.value = field
    sortDir.value = 'desc'
  }
}

const sortedItems = computed(() => {
  if (!sortField.value) return props.items
  const dir = sortDir.value === 'asc' ? 1 : -1
  return [...props.items].sort((a, b) => {
    const f = sortField.value!
    if (f === 'title') {
      return dir * a.novela.title.localeCompare(b.novela.title)
    }
    const aVal = a[f] as number
    const bVal = b[f] as number
    return dir * (aVal - bVal)
  })
})
</script>

<template>
  <div class="overflow-x-auto">
    <table class="w-full text-left border-collapse">
      <thead>
        <tr
          class="border-b border-zinc-200 dark:border-zinc-800 text-sm font-semibold text-zinc-500 dark:text-zinc-400"
        >
          <th class="py-4 px-4 w-12 text-center">#</th>
          <th
            class="py-4 px-4 cursor-pointer select-none hover:text-zinc-700 dark:hover:text-zinc-300 transition-colors"
            @click="setSort('title')"
          >
            <span class="inline-flex items-center gap-1">
              Новелла
              <Icon
                v-if="sortField === 'title'"
                :name="sortDir === 'asc' ? 'ph:caret-up-bold' : 'ph:caret-down-bold'"
                size="12"
              />
            </span>
          </th>
          <th
            class="py-4 px-4 text-center cursor-pointer select-none hover:text-zinc-700 dark:hover:text-zinc-300 transition-colors"
            @click="setSort('rating')"
          >
            <span class="inline-flex items-center gap-1">
              Рейтинг
              <Icon
                v-if="sortField === 'rating'"
                :name="sortDir === 'asc' ? 'ph:caret-up-bold' : 'ph:caret-down-bold'"
                size="12"
              />
            </span>
          </th>
          <th
            class="py-4 px-4 text-center cursor-pointer select-none hover:text-zinc-700 dark:hover:text-zinc-300 transition-colors"
            @click="setSort('readCount')"
          >
            <span class="inline-flex items-center gap-1">
              <TeleportedTooltip text="Прочтения">
                <Icon name="ph:eye-bold" size="18" class="inline-block" />
              </TeleportedTooltip>
            </span>
          </th>
          <th
            class="py-4 px-4 text-center cursor-pointer select-none hover:text-zinc-700 dark:hover:text-zinc-300 transition-colors"
            @click="setSort('bookmarkCount')"
          >
            <span class="inline-flex items-center gap-1">
              <TeleportedTooltip text="Закладки">
                <Icon name="ph:bookmark-simple-bold" size="18" class="inline-block" />
              </TeleportedTooltip>
            </span>
          </th>
          <th
            class="py-4 px-4 text-center cursor-pointer select-none hover:text-zinc-700 dark:hover:text-zinc-300 transition-colors"
            @click="setSort('commentCount')"
          >
            <span class="inline-flex items-center gap-1">
              <TeleportedTooltip text="Комментарии">
                <Icon name="ph:chat-circle-bold" size="18" class="inline-block" />
              </TeleportedTooltip>
            </span>
          </th>
          <th
            class="py-4 px-4 text-center cursor-pointer select-none hover:text-zinc-700 dark:hover:text-zinc-300 transition-colors"
            @click="setSort('ratingCount')"
          >
            <span class="inline-flex items-center gap-1">
              <TeleportedTooltip text="Оценки">
                <Icon name="ph:star-bold" size="18" class="inline-block" />
              </TeleportedTooltip>
            </span>
          </th>
          <th
            class="py-4 px-4 text-center cursor-pointer select-none hover:text-zinc-700 dark:hover:text-zinc-300 transition-colors"
            @click="setSort('volumeCount')"
          >
            <span class="inline-flex items-center gap-1">
              <TeleportedTooltip text="Тома">
                <Icon name="ph:books-bold" size="18" class="inline-block" />
              </TeleportedTooltip>
            </span>
          </th>
          <th
            class="py-4 px-4 text-center cursor-pointer select-none hover:text-zinc-700 dark:hover:text-zinc-300 transition-colors"
            @click="setSort('chapterCount')"
          >
            <span class="inline-flex items-center gap-1">
              <TeleportedTooltip text="Главы">
                <Icon name="ph:file-text-bold" size="18" class="inline-block" />
              </TeleportedTooltip>
            </span>
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(item, index) in sortedItems"
          :key="item.novela.novelaId"
          class="border-b border-zinc-100 dark:border-zinc-800/50 hover:bg-zinc-50 dark:hover:bg-zinc-800/30 transition-colors"
        >
          <td class="py-3 px-4 text-center font-medium text-zinc-500">
            {{ offset + index + 1 }}
          </td>
          <td class="py-3 px-4">
            <NuxtLink :to="`/novela/${item.novela.novelaId}`" class="flex items-center gap-3 group">
              <img
                :src="staticImage(item.novela.posterUrl)"
                :alt="item.novela.title"
                class="w-10 h-14 rounded-lg object-cover shadow-sm group-hover:shadow-md transition-shadow shrink-0"
              />
              <span
                class="font-semibold text-zinc-900 dark:text-zinc-100 group-hover:text-yellow-500 dark:group-hover:text-yellow-400 transition-colors line-clamp-2"
              >
                {{ item.novela.title }}
              </span>
            </NuxtLink>
          </td>
          <td class="py-3 px-4 text-center text-zinc-600 dark:text-zinc-400 font-medium">
            <div class="flex items-center justify-center gap-1">
              <Icon name="ph:star-fill" class="text-yellow-500 shrink-0" size="16" />
              <span>{{ roundTo(item.rating, 1) }}</span>
            </div>
          </td>
          <td class="py-3 px-4 text-center text-zinc-600 dark:text-zinc-400 font-medium">
            {{ item.readCount }}
          </td>
          <td class="py-3 px-4 text-center text-zinc-600 dark:text-zinc-400 font-medium">
            {{ item.bookmarkCount }}
          </td>
          <td class="py-3 px-4 text-center text-zinc-600 dark:text-zinc-400 font-medium">
            {{ item.commentCount }}
          </td>
          <td class="py-3 px-4 text-center text-zinc-600 dark:text-zinc-400 font-medium">
            {{ item.ratingCount }}
          </td>
          <td class="py-3 px-4 text-center text-zinc-600 dark:text-zinc-400 font-medium">
            {{ item.volumeCount }}
          </td>
          <td class="py-3 px-4 text-center text-zinc-600 dark:text-zinc-400 font-medium">
            {{ item.chapterCount }}
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
