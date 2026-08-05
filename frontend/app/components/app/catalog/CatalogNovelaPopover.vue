<script setup lang="ts">
import type { NovelaDetails } from '@/types/backend/novela'
import { roundTo } from '@/utils/num'
import { replaceTags } from '@/utils/str'

const props = defineProps<{
  novela: NovelaDetails
}>()

const novelaDescription = computed(() => {
  return props.novela.description ? replaceTags(props.novela.description) : ''
})

const novelaReleaseYear = computed(() => {
  return props.novela.release_date ? new Date(props.novela.release_date).getFullYear() : 'N/A'
})
</script>

<template>
  <div
    class="absolute z-50 w-72 md:w-80 bottom-full left-1/2 -translate-x-1/2 mb-3 md:left-full md:ml-3 md:top-0 md:bottom-auto md:translate-x-0 md:mb-0 opacity-0 invisible group-hover/popover:opacity-100 group-hover/popover:visible transition-all duration-200 delay-150 pointer-events-none"
  >
    <div
      class="bg-white dark:bg-zinc-800 rounded-xl shadow-2xl border border-zinc-200 dark:border-zinc-700 p-4 max-h-80 overflow-y-auto"
    >
      <h3
        class="font-bold text-zinc-900 dark:text-white text-sm leading-snug mb-1.5 hyphens-auto"
        lang="ru"
      >
        {{ novela.title }}
      </h3>

      <div v-if="novela.alternative_titles?.length" class="mb-2">
        <p
          v-for="(alt, i) in novela.alternative_titles.slice(0, 2)"
          :key="i"
          class="text-[11px] text-zinc-400 dark:text-zinc-500 italic truncate"
        >
          {{ alt }}
        </p>
        <p
          v-if="novela.alternative_titles.length > 2"
          class="text-[11px] text-zinc-400 dark:text-zinc-500 italic"
        >
          +{{ novela.alternative_titles.length - 2 }} more
        </p>
      </div>

      <p class="text-[11px] text-zinc-600 dark:text-zinc-400 leading-relaxed mb-2.5 line-clamp-4">
        {{ novelaDescription }}
      </p>

      <div v-if="novela.genres?.length" class="flex flex-wrap gap-1 mb-2.5">
        <span
          v-for="genre in novela.genres.slice(0, 4)"
          :key="genre"
          class="text-[10px] px-2 py-0.5 rounded-full bg-zinc-100 dark:bg-zinc-700 text-zinc-600 dark:text-zinc-300 font-medium"
        >
          {{ genre }}
        </span>
        <span
          v-if="novela.genres.length > 4"
          class="text-[10px] px-2 py-0.5 rounded-full bg-zinc-100 dark:bg-zinc-700 text-zinc-400 dark:text-zinc-500 font-medium"
        >
          +{{ novela.genres.length - 4 }}
        </span>
      </div>

      <div
        class="flex items-center gap-2 text-[10px] text-zinc-500 dark:text-zinc-400 font-semibold mb-2 flex-wrap"
      >
        <span class="uppercase">{{ novela.type }}</span>
        <span class="w-1 h-1 rounded-full bg-zinc-300 dark:bg-zinc-600 shrink-0" />
        <span class="truncate">{{ novela.status }}</span>
        <span v-if="novela.release_date" class="shrink-0">
          <Icon name="ph:calendar-bold" size="10" class="inline -mt-0.5" />
          {{ novelaReleaseYear }}
        </span>
        <span v-if="novela.country" class="shrink-0">{{ novela.country }}</span>
      </div>

      <div class="flex items-center gap-2 mb-2">
        <span
          v-if="novela.age_rating"
          class="text-[10px] font-bold px-1.5 py-0.5 rounded bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400"
        >
          {{ novela.age_rating }}
        </span>
        <span
          v-if="novela.translation_status"
          class="text-[10px] font-medium px-1.5 py-0.5 rounded bg-emerald-100 dark:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400"
        >
          {{ novela.translation_status }}
        </span>
      </div>

      <div class="flex items-center gap-4">
        <div class="flex items-center gap-1.5">
          <span class="text-amber-500 text-sm font-bold">★</span>
          <span class="text-sm font-bold text-zinc-800 dark:text-zinc-200">
            {{ roundTo(novela.rating_details?.total_rating, 1) }}
          </span>
          <span class="text-[10px] text-zinc-400 dark:text-zinc-500">
            ({{ novela.rating_details?.total ?? 0 }})
          </span>
        </div>
        <div
          v-if="novela.views"
          class="flex items-center gap-1 text-[11px] text-zinc-400 dark:text-zinc-500"
        >
          <Icon name="ph:eye-bold" size="12" />
          {{ novela.views >= 1000 ? `${(novela.views / 1000).toFixed(1)}k` : novela.views }}
        </div>
        <div
          v-if="novela.like_count"
          class="flex items-center gap-1 text-[11px] text-zinc-400 dark:text-zinc-500"
        >
          <Icon name="ph:heart-fill" size="14" class="text-amber-500" />
          {{
            novela.like_count >= 1000
              ? `${(novela.like_count / 1000).toFixed(1)}k`
              : novela.like_count
          }}
        </div>
        <div
          v-if="novela.bookmark_details"
          class="flex items-center gap-1 text-[11px] text-zinc-400 dark:text-zinc-500"
        >
          <Icon name="ph:bookmark-simple-fill" size="14" class="text-amber-500" />
          {{
            novela.bookmark_details.total >= 1000
              ? `${(novela.bookmark_details.total / 1000).toFixed(1)}k`
              : novela.bookmark_details.total
          }}
        </div>
      </div>
    </div>

    <!-- Arrow: mobile (down) -->
    <div
      class="md:hidden absolute left-1/2 -translate-x-1/2 top-full -mt-0.5 w-3 h-3 rotate-45 bg-white dark:bg-zinc-800 border-r border-b border-zinc-200 dark:border-zinc-700"
    />
    <!-- Arrow: desktop (left) -->
    <div
      class="hidden md:block absolute -left-1.5 top-4 w-3 h-3 rotate-45 bg-white dark:bg-zinc-800 border-l border-b border-zinc-200 dark:border-zinc-700"
    />
  </div>
</template>
