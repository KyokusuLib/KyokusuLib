<script setup lang="ts">
import { staticImage, roleConv } from '@/utils/str'
import { roundTo } from '@/utils/num'
import type { SearchCategory, SearchResultItem } from '@/composables/api/search/useSearch'

const props = defineProps<{
  item: SearchResultItem
  activeCategory: SearchCategory
}>()

const emit = defineEmits(['select'])

const getIcon = () => {
  switch (props.activeCategory) {
    case 'ranobe':
      return 'ph:book-open-bold'
    case 'users':
      return 'ph:user-bold'
    case 'teams':
      return 'ph:users-three-bold'
    case 'authors':
      return 'ph:pen-nib-bold'
    default:
      return 'ph:image-bold'
  }
}

const getTitle = () => {
  return props.item.title || props.item.name || props.item.label || 'Без названия'
}

const getSubtitle = () => {
  if (props.item.role) {
    return roleConv(props.item.role, 'ru')
  }
  return props.item.type || props.item.description || ''
}
</script>

<template>
  <button
    @click="emit('select', item)"
    class="w-full flex items-center gap-4 p-2.5 sm:p-3 rounded-2xl hover:bg-white dark:hover:bg-[#1f1f23] border border-transparent hover:border-zinc-200 dark:hover:border-zinc-800 hover:shadow-[0_4px_20px_-4px_rgba(0,0,0,0.05)] dark:hover:shadow-[0_4px_20px_-4px_rgba(0,0,0,0.3)] transition-all duration-300 text-left group cursor-pointer active:scale-[0.98]"
  >
    <div
      class="shrink-0 bg-zinc-100 dark:bg-zinc-800 overflow-hidden border border-zinc-200/80 dark:border-zinc-700/80 group-hover:border-yellow-500/30 transition-colors duration-300"
      :class="[
        activeCategory === 'ranobe'
          ? 'w-12 h-16 sm:w-14 sm:h-20 rounded-xl shadow-inner'
          : 'w-12 h-12 sm:w-14 sm:h-14 rounded-full shadow-sm',
      ]"
    >
      <img
        v-if="item.poster_url || item.picture || item.avatar || item.banner || item.avatar_url"
        :src="
          staticImage(
            item.poster_url || item.picture || item.avatar || item.banner || item.avatar_url,
          )
        "
        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
      />
      <div
        v-else
        class="w-full h-full flex items-center justify-center bg-zinc-100 dark:bg-zinc-800 text-zinc-400 group-hover:text-yellow-500 transition-colors duration-300"
      >
        <Icon
          :name="getIcon()"
          size="20"
          class="group-hover:scale-110 transition-transform duration-300"
        />
      </div>
    </div>

    <div class="flex flex-col flex-1 min-w-0 py-0.5 justify-center">
      <div class="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-2 min-w-0">
        <div
          class="flex flex-col sm:flex-row items-center gap-4 text-sm sm:text-base font-bold text-zinc-900 dark:text-zinc-100 line-clamp-2 sm:truncate group-hover:text-yellow-600 dark:group-hover:text-yellow-500 transition-colors duration-300"
        >
          <div class="flex flex-col items-start">
            <div class="flex items-center text-center">{{ getTitle() }}</div>

            <div
              v-if="item.slug"
              class="text-[12px] text-zinc-500 dark:text-zinc-400 group-hover:text-zinc-100 transition-colors"
            >
              @{{ item.slug }}
            </div>
          </div>

          <div
            v-if="item.stats !== null && item.stats !== undefined"
            class="flex flex-col sm:flex-row items-center gap-2 text-[12px] text-zinc-500 dark:text-zinc-400 group-hover:text-zinc-100 transition-colors"
          >
            <div
              class="flex items-center justify-center px-2 py-0.5 mt-0.5 bg-zinc-100 dark:bg-zinc-800 rounded-4xl text-[10px] font-semibold"
            >
              Участников: {{ item.stats.member_count }}
            </div>
            <div
              class="flex items-center justify-center px-2 py-0.5 mt-0.5 bg-zinc-100 dark:bg-zinc-800 rounded-4xl text-[10px] font-semibold"
            >
              Подписчиков: {{ item.stats.subscriber_count ? item.stats.subscriber_count : 0 }}
            </div>

            <div
              v-if="item.is_member"
              class="flex items-center justify-center px-2 py-0.5 mt-0.5 bg-zinc-100 dark:bg-zinc-800 rounded-4xl text-[10px] font-semibold text-yellow-500"
            >
              Участник
            </div>
            <div
              v-if="item.is_subscriber"
              class="flex items-center justify-center px-2 py-0.5 mt-0.5 bg-zinc-100 dark:bg-zinc-800 rounded-4xl text-[10px] font-semibold text-yellow-500"
            >
              Подписчик
            </div>
          </div>

          <div
            v-if="activeCategory !== 'teams' && getSubtitle()"
            class="flex flex-wrap items-center gap-2 mt-1 sm:mt-1.5"
          >
            <span
              class="flex items-center gap-1.5 text-xs text-zinc-500 dark:text-zinc-400 min-w-0 max-w-full"
            >
              <span
                v-if="activeCategory === 'ranobe'"
                class="w-1.5 h-1.5 mb-1.5 shrink-0 rounded-full bg-zinc-300 dark:bg-zinc-700 group-hover:bg-yellow-500/50 transition-colors duration-300"
              ></span>
              <span class="flex mb-1.5 truncate">{{ getSubtitle() }}</span>
            </span>
          </div>

          <div
            v-if="item.rating_details && item.rating_details.total_rating"
            class="flex flex-wrap items-center gap-2 mt-1 sm:mt-1.5"
          >
            <span
              class="flex items-center gap-1 text-[10px] font-bold px-1.5 py-0.5 rounded bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-300 shrink-0"
            >
              <Icon name="ph:star-fill" class="text-yellow-500" />
              {{ roundTo(item.rating_details.total_rating, 2) }}
            </span>
          </div>
        </div>

        <div v-if="item.user_level" class="flex items-center gap-1.5 sm:gap-2 shrink-0">
          <div
            class="relative flex items-center justify-center w-4 h-4 rounded-md bg-zinc-900 dark:bg-zinc-100 text-white dark:text-zinc-900 shadow-sm"
          >
            <span class="text-[10px] font-black italic leading-none">{{
              item.user_level.level
            }}</span>
            <div
              class="absolute inset-0 bg-linear-to-tr from-white/10 to-transparent rounded-md pointer-events-none"
            ></div>
          </div>
          <div
            class="flex shrink-0 px-2 h-4 rounded-full items-center justify-center bg-zinc-900 dark:bg-zinc-100 text-white dark:text-zinc-900 shadow-sm"
          >
            <span class="text-[9px] font-black italic leading-none">{{
              item.user_level.level_title
            }}</span>
          </div>
        </div>
      </div>
      <div
        v-if="item.alternative_titles && item.alternative_titles.length"
        class="flex text-center text-[11px] sm:text-xs text-zinc-500 dark:text-zinc-400 line-clamp-1 sm:truncate mt-0.5 sm:mt-1"
      >
        {{ item.alternative_titles.join(' • ') }}
      </div>

      <div
        v-if="activeCategory === 'teams' && getSubtitle()"
        class="text-xs text-zinc-500 dark:text-zinc-400 line-clamp-2 mt-2 pt-2 border-t border-zinc-100 dark:border-zinc-800"
      >
        {{ getSubtitle() }}
      </div>
    </div>

    <div
      class="flex items-center justify-center w-8 h-8 shrink-0 rounded-full bg-zinc-50 dark:bg-[#121214] border border-zinc-200 dark:border-zinc-800 opacity-0 -translate-x-2 group-hover:opacity-100 group-hover:translate-x-0 transition-all duration-300 sm:mr-1"
    >
      <Icon
        name="ph:arrow-right-bold"
        size="14"
        class="text-zinc-400 group-hover:text-yellow-500 transition-colors"
      />
    </div>
  </button>
</template>
