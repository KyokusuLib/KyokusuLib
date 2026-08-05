<script setup lang="ts">
import { useNovela } from '@/composables/api/novela/useNovela'
import { staticImage } from '@/utils/str'
import ModalConfirm from '~/components/common/ModalConfirm.vue'
import { useReadProgress } from '~/composables/api/novela/useReadProgress'
import { useUserActivity } from '~/composables/api/profile/useUserActivity'
import { useRolePermissions } from '~/composables/api/role/useRolePermissions'
import { $api } from '~/composables/api/useApi'
import { ACTIVITY_TYPES } from '~/constants/user-activity'
import { useNotificationStore } from '~/stores/notification'
import type { NovelaCommentResponse, NovelaDetails, NovelaVolume } from '~/types/backend/novela'
import { KyokusuAppRole } from '~/types/enums/role-enum'
import { type BookmarkCategory, getBookmarkCategoryLabel } from '~/types/frontend/bookmarks'
import {
  convToRu,
  NOVELA_ACTIVE_TABS,
  type NovelaActiveTabs,
  NovelaActiveTabsEnum,
} from '~/types/frontend/novela/novela-active-tabs'
import ChapterList from './ChapterList.vue'
import NovelaBookmarkButton from './NovelaBookmarkButton.vue'
import NovelaBookmarkStats from './NovelaBookmarkStats.vue'
import NovelaComments from './NovelaComments.vue'
import NovelaLikeButton from './NovelaLikeButton.vue'
import NovelaRating from './NovelaRating.vue'
import NovelaRatingStats from './NovelaRatingStats.vue'
import NovelaSettings from './NovelaSettings.vue'
import NovelaStats from './NovelaStats.vue'

import { useNovelaComments } from '@/composables/api/novela/useNovelaComments'

const route = useRoute()

const { user } = useAuthStore()
const { novela, fetchNovela } = useNovela()
const chaptersKey = ref(0)

const { createUserActivity } = useUserActivity()
const { hasPermission } = useRolePermissions()
const { getContinueReadingUrl } = useReadProgress()
const { comments, totalComments, fetchComments, isLoading: commentsLoading } = useNovelaComments()

const novelaId = computed(() => route.params.id as string)

const continueReadingUrl = computed(() =>
  novela.value ? getContinueReadingUrl(novela.value) : null,
)

const lastReadChapterLabel = computed(() => {
  if (!novela.value?.last_readed?.chapter_id) return null
  return `Глава ${novela.value.last_readed.chapter_number}`
})

const activeTab = ref<NovelaActiveTabs>(
  route.query.tab === NovelaActiveTabsEnum.COMMENTS
    ? NovelaActiveTabsEnum.COMMENTS
    : NovelaActiveTabsEnum.ABOUT,
)

await useAsyncData(`novela-${novelaId.value}`, () => fetchNovela(novelaId.value))

await useAsyncData(`novela-comments-${novelaId.value}`, () => fetchComments(Number(novelaId.value)))

const bookmarkInitial = ref(Boolean(novela.value?.bookmark))
const currentBookmarkCategory = ref(novela.value?.bookmark || null)
const isOpenNovelaSettings = ref(false)
const isOpenDeleteConfirm = ref(false)

const { notify } = useNotificationStore()

const deleteNovela = async () => {
  if (!novela.value) return
  try {
    await $api(`/api/novela/${novela.value.id}`, { method: 'DELETE' })
    notify({
      type: 'success',
      title: 'Новела удалена',
      content: `«${novela.value.title}» успешно удалена`,
    })
    await navigateTo('/')
  } catch (e: any) {
    notify({
      type: 'error',
      title: 'Ошибка',
      content: e?.message || 'Не удалось удалить новелу',
    })
  }
}

const openedSections = ref<string[]>([])

const totalChapters = computed(() => {
  if (!novela.value) return 0
  return novela.value.volumes.reduce((acc, v) => acc + (v.chapters?.length || 0), 0)
})

const novelaInfo = computed(() => [
  { label: 'Статус', value: novela.value?.status },
  { label: 'Перевод', value: novela.value?.translation_status },
  { label: 'Страна', value: novela.value?.country },
])

const updateCountBookmarks = async (newCategory: BookmarkCategory | null) => {
  if (!novela.value || !user?.id) return

  const prevCategory = currentBookmarkCategory.value

  if (prevCategory === newCategory) return

  if (newCategory) {
    const item = novela.value.bookmark_details.nc_items.find((i) => i.value === newCategory)
    if (item) item.count++
  }

  if (prevCategory) {
    const item = novela.value.bookmark_details.nc_items.find((i) => i.value === prevCategory)
    if (item) item.count--
  }

  if (prevCategory === null && newCategory !== null) {
    novela.value.bookmark_details.total++

    await createUserActivity({
      user_id: user.id,
      activity_type: ACTIVITY_TYPES.NOVELA_BOOKMARK,
      target_id: novela.value.id,
      metadata: {
        novela_title: novela.value.title,
        desc: 'Пользователь добавил сюжет в закладки',
      },
    })
  } else if (prevCategory !== null && newCategory !== null) {
    await createUserActivity({
      user_id: user.id,
      activity_type: ACTIVITY_TYPES.NOVELA_BOOKMARK,
      target_id: novela.value.id,
      metadata: {
        novela_title: novela.value.title,
        desc: `Пользователь изменил категорию на "${getBookmarkCategoryLabel(newCategory)}"`,
      },
    })
  } else if (prevCategory !== null && newCategory === null) {
    novela.value.bookmark_details.total--

    await createUserActivity({
      user_id: user.id,
      activity_type: ACTIVITY_TYPES.NOVELA_BOOKMARK_REMOVE,
      target_id: novela.value.id,
      metadata: {
        novela_title: novela.value.title,
        desc: 'Пользователь удалил сюжет из закладок',
      },
    })
  }

  currentBookmarkCategory.value = newCategory
  bookmarkInitial.value = Boolean(newCategory)
}

const updateCountLike = async (has_liked: boolean) => {
  if (novela.value && has_liked) {
    if (user?.id) {
      await createUserActivity({
        user_id: user.id,
        activity_type: ACTIVITY_TYPES.USER_NOVELA_LIKE,
        target_id: novela.value.id,
        metadata: {
          name: novela.value.title,
          desc: 'Пользователь поставил лайк сюжету',
        },
      })
    }
    novela.value.like_count = (novela.value.like_count || 0) + 1
  } else if (novela.value && !has_liked) {
    if (user?.id) {
      await createUserActivity({
        user_id: user.id,
        activity_type: ACTIVITY_TYPES.USER_NOVELA_LIKE_REMOVE,
        target_id: novela.value.id,
        metadata: {
          name: novela.value.title,
          desc: 'Пользователь убрал лайк',
        },
      })
    }
    novela.value.like_count = (novela.value.like_count || 0) - 1
  }
}

const updateRating = (rating: number) => {
  if (!novela.value) return

  const oldCount = novela.value.rating_details.total || 0
  const oldAverage = novela.value.rating_details.total_rating || 0
  const previousUserRating = novela.value.user_rating || 0

  if (previousUserRating === rating) return

  const item = novela.value.rating_details.nc_items.find((i) => i.value === rating)
  if (item) item.count++

  if (previousUserRating !== 0) {
    const item2 = novela.value.rating_details.nc_items.find((i) => i.value === previousUserRating)
    if (item2) item2.count--
  }

  let newCount = oldCount
  let newAverage = oldAverage

  if (previousUserRating === 0) {
    newCount = oldCount + 1
    newAverage = (oldAverage * oldCount + rating) / newCount
  } else {
    const oldSum = oldAverage * oldCount
    newAverage = (oldSum - previousUserRating + rating) / oldCount
  }

  novela.value.rating_details.total = newCount
  novela.value.rating_details.total_rating = roundTo(newAverage, 1)
  novela.value.user_rating = rating

  if (user?.id) {
    createUserActivity({
      user_id: user.id,
      activity_type: ACTIVITY_TYPES.USER_NOVELA_RATING,
      target_id: novela.value.id,
      metadata: {
        name: novela.value.title,
        rating: rating,
        desc: 'Пользователь оценил произведение',
      },
    })
  }
}

const updatedNovela = (payload: NovelaDetails) => {
  if (novela.value !== null) {
    Object.assign(novela.value, payload)
  }
}

const toggleSection = (name: string) => {
  if (openedSections.value.includes(name)) {
    openedSections.value = openedSections.value.filter((s) => s !== name)
  } else {
    openedSections.value.push(name)
  }
}

function onVolumeDeleted(volumeId: string) {
  chaptersKey.value++
  if (novela.value) {
    Object.assign(novela.value, {
      volumes: novela.value.volumes.filter((v) => v.id !== volumeId),
    })
  }
}

function onVolumeAdded(volumeData?: {
  id: string
  title: string
  number: number
  status?: string
}) {
  chaptersKey.value++
  if (novela.value && volumeData) {
    const newVolume: NovelaVolume = {
      id: volumeData.id,
      title: volumeData.title,
      number: volumeData.number,
      status: volumeData.status,
      chapters: [],
    }
    novela.value.volumes.push(newVolume)
  }
}

watch(
  () => route.query.tab,
  (newTab) => {
    if (newTab === NovelaActiveTabsEnum.COMMENTS) activeTab.value = NovelaActiveTabsEnum.COMMENTS
  },
)
</script>

<template>
  <div
    class="min-h-screen bg-zinc-50 dark:bg-[#0f0f0f] text-zinc-900 dark:text-zinc-200 flex flex-col font-sans transition-colors duration-300"
  >
    <main class="grow" v-if="novela">
      <div class="relative h-75 md:h-112.5 w-full overflow-hidden">
        <div class="absolute inset-0 bg-black/40 z-10 backdrop-blur-sm"></div>
        <div
          class="absolute inset-0 bg-linear-to-t from-zinc-50 dark:from-[#0f0f0f] via-transparent to-transparent z-20"
        ></div>
        <img
          :src="staticImage(novela.poster_url || '')"
          class="w-full h-full object-cover blur-md scale-110"
          alt="Backdrop"
        />
      </div>

      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 -mt-64 md:-mt-106 relative z-30 pb-12">
        <div class="flex flex-col md:flex-row gap-8">
          <div class="w-full md:w-75 shrink-0 flex flex-col gap-6">
            <div
              class="relative rounded-2xl overflow-hidden shadow-2xl aspect-2/3 border border-white/10 group"
            >
              <img
                :src="staticImage(novela.poster_url || '')"
                class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
                :alt="novela.title"
              />
              <div
                class="absolute top-3 left-3 bg-zinc-900/80 backdrop-blur px-2 py-1 rounded-lg text-xs font-bold text-white"
              >
                {{ novela.age_rating }}
              </div>
              <div
                class="absolute top-3 right-3 bg-zinc-900/80 backdrop-blur px-2 py-1 rounded-lg text-xs font-bold text-white"
              >
                {{ novela.type }}
              </div>
            </div>

            <div class="flex flex-col gap-3">
              <NuxtLink
                v-if="continueReadingUrl"
                :to="{
                  path: continueReadingUrl,
                  query: { scroll: novela.last_readed!.chapter_scroll },
                }"
                class="w-full py-1.5 text-center cursor-pointer rounded-xl font-bold shadow-lg transition-all active:scale-95 bg-yellow-500 text-white hover:bg-yellow-600 flex flex-col items-center leading-tight"
              >
                <span>{{ novela.last_readed ? 'Продолжить' : 'Читать' }}</span>
                <span v-if="lastReadChapterLabel" class="text-[10px] font-medium opacity-80">{{
                  lastReadChapterLabel
                }}</span>
              </NuxtLink>
              <button
                v-else
                class="w-full py-1.5 cursor-not-allowed opacity-50 rounded-xl font-bold shadow-lg bg-white text-zinc-900"
              >
                Нет глав
              </button>
              <div class="grid grid-cols-2 gap-3">
                <NovelaBookmarkButton
                  v-model="novela.bookmark"
                  :novela-id="novela.id"
                  @update:model-value="updateCountBookmarks"
                />
                <NovelaLikeButton
                  v-model="novela.has_liked"
                  :novela-id="novela.id"
                  @update:model-value="updateCountLike"
                />
              </div>
            </div>

            <div
              class="hidden md:flex flex-col gap-4 bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-zinc-800 p-6 rounded-2xl backdrop-blur-sm"
            >
              <div v-for="info in novelaInfo" :key="info.label">
                <div class="flex justify-between items-center text-sm">
                  <span class="text-zinc-500">{{ info.label }}</span>
                  <span
                    :class="[
                      'font-semibold cursor-default',
                      info.value ? 'text-zinc-900 dark:text-zinc-200' : 'text-zinc-500',
                    ]"
                  >
                    {{ info.value || '—' }}
                  </span>
                </div>
              </div>
              <div class="flex justify-between items-start text-sm">
                <span class="text-zinc-500">{{ novela.authors?.length ? 'Авторы' : 'Автор' }}</span>
                <div
                  :class="[
                    'font-semibold cursor-default',
                    !!novela.authors ? 'text-zinc-900 dark:text-zinc-200' : 'text-zinc-500',
                  ]"
                  class="flex flex-col items-end gap-2"
                >
                  <span v-for="author in novela.authors" :key="author.id">
                    {{ author.name }}
                  </span>
                </div>
              </div>
              <div class="h-px bg-zinc-200 dark:bg-zinc-800 my-1"></div>
              <div class="flex flex-col justify-center flex-wrap gap-2">
                <h3 class="text-xl text-center font-black mb-4 uppercase tracking-tighter">
                  {{ novela.genres.length >= 2 ? 'Жанры' : 'Жанр' }}
                </h3>
                <div class="flex justify-center flex-wrap gap-2">
                  <span
                    v-for="genre in novela.genres"
                    :key="genre"
                    class="px-3 py-1 rounded-lg bg-zinc-100 dark:bg-zinc-800 text-[11px] font-bold uppercase tracking-wider text-zinc-600 dark:text-zinc-400 hover:bg-zinc-200 dark:hover:bg-zinc-700 cursor-pointer transition-colors"
                  >
                    {{ genre }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div class="flex-1 min-w-0 md:pt-12">
            <div class="md:mb-48">
              <div class="flex flex-row">
                <div class="flex flex-col gap-3 mb-12">
                  <h1
                    class="text-4xl md:text-4xl shadow-text text-white font-black tracking-tight leading-none"
                  >
                    {{ novela.title }}
                  </h1>
                  <h2
                    v-if="novela.alternative_titles?.length"
                    class="text-lg md:text-xl text-zinc-300 dark:text-zinc-300/80 font-medium italic"
                  >
                    {{ novela.alternative_titles.join(' • ') }}
                  </h2>
                </div>
                <div
                  v-if="hasPermission(KyokusuAppRole.MODERATOR)"
                  class="flex items-start mt-4.5 w-auto gap-6 md:ml-12"
                >
                  <button
                    class="p-2 -m-4 cursor-pointer group outline-none"
                    @click="isOpenNovelaSettings = true"
                  >
                    <Icon
                      name="ph:gear"
                      size="24"
                      class="text-zinc-200 group-hover:text-zinc-400 transition-colors duration-300"
                    />
                  </button>
                  <button
                    v-if="hasPermission(KyokusuAppRole.ADMIN)"
                    class="p-2 -m-4 cursor-pointer group outline-none"
                    @click="isOpenDeleteConfirm = true"
                  >
                    <Icon
                      name="ph:trash"
                      size="24"
                      class="text-zinc-200 group-hover:text-zinc-400 transition-colors duration-300"
                    />
                  </button>
                </div>
              </div>

              <div class="flex flex-wrap items-baseline-last gap-6 ml-2">
                <NovelaStats :novela="novela" />
                <NovelaRating
                  :novela-id="novela.id"
                  :rating="novela.rating_details.total_rating"
                  :count="novela.rating_details.total"
                  :user-rating="novela.user_rating || 1"
                  @update:rated="updateRating"
                />
              </div>
              <div class="flex flex-row gap-6 mb-2 md:flex-row items-start">
                <NovelaRatingStats
                  :is-expanded="openedSections.includes('rating')"
                  @toggle="toggleSection('rating')"
                  :rating-details="novela.rating_details"
                />
                <NovelaBookmarkStats
                  :is-expanded="openedSections.includes('bookmark')"
                  @toggle="toggleSection('bookmark')"
                  :bookmark-details="novela.bookmark_details"
                />
              </div>
            </div>

            <div
              class="md:hidden mb-8 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 p-5 rounded-2xl space-y-4"
            >
              <div
                v-for="info in novelaInfo"
                :key="info.label"
                class="flex justify-between text-sm"
              >
                <span class="text-zinc-500">{{ info.label }}</span>
                <span class="font-bold">{{ info.value || '—' }}</span>
              </div>
              <div class="flex justify-between items-start text-sm">
                <span class="text-zinc-500">{{ novela.authors?.length ? 'Авторы' : 'Автор' }}</span>
                <div
                  :class="[
                    'font-semibold cursor-default',
                    !!novela.authors ? 'text-zinc-900 dark:text-zinc-200' : 'text-zinc-500',
                  ]"
                  class="flex flex-col items-end gap-2"
                >
                  <span v-for="author in novela.authors" :key="author.id">
                    {{ author.name }}
                  </span>
                </div>
              </div>
              <div class="flex flex-wrap gap-2">
                <span
                  v-for="genre in novela.genres"
                  :key="genre"
                  class="px-2 py-1 rounded bg-zinc-100 dark:bg-zinc-800 text-[10px] font-bold text-zinc-500"
                >
                  {{ genre }}
                </span>
              </div>
            </div>

            <div
              class="overflow-x-auto no-scrollbar mb-8 border-b border-zinc-200 dark:border-zinc-800 min-w-0 max-w-full w-full"
            >
              <div class="inline-flex gap-8 whitespace-nowrap px-4 md:px-0">
                <button
                  v-for="tab in NOVELA_ACTIVE_TABS"
                  :key="tab"
                  class="flex flex-row items-center pb-4 text-lg font-bold transition-all relative cursor-pointer capitalize shrink-0 min-h-11"
                  :class="activeTab === tab ? 'text-zinc-900 dark:text-white' : 'text-zinc-400'"
                  @click="activeTab = tab as any"
                >
                  {{ convToRu(tab) }}
                  <span v-if="tab === NovelaActiveTabsEnum.CHAPTERS" class="ml-2 text-lg opacity-50"
                    >({{ totalChapters }})</span
                  >
                  <span v-if="tab === NovelaActiveTabsEnum.COMMENTS" class="ml-2 text-lg opacity-50"
                    >({{ totalComments }})</span
                  >
                  <div
                    v-if="activeTab === tab"
                    class="absolute bottom-0 left-0 w-full h-1 bg-yellow-500 rounded-t-full"
                  ></div>
                </button>
              </div>
            </div>

            <div class="min-h-75">
              <transition name="fade" mode="out-in">
                <div
                  v-if="activeTab === NovelaActiveTabsEnum.ABOUT"
                  :key="NovelaActiveTabsEnum.ABOUT"
                  class="space-y-10"
                >
                  <div
                    class="prose dark:prose-invert max-w-none text-zinc-700 dark:text-zinc-300 leading-relaxed text-base md:text-lg whitespace-pre-line"
                    v-html="novela.description"
                  ></div>

                  <div v-if="novela.categories?.length">
                    <h3 class="text-xl font-black mb-4 uppercase tracking-tighter">Категории</h3>
                    <div class="flex flex-wrap gap-3">
                      <span
                        v-for="cat in novela.categories"
                        :key="cat"
                        class="px-4 py-2 rounded-xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900/50 text-sm font-medium hover:border-yellow-500 cursor-pointer transition-all"
                      >
                        #{{ cat }}
                      </span>
                    </div>
                  </div>
                </div>
                <div
                  v-else-if="activeTab === NovelaActiveTabsEnum.CHAPTERS"
                  :key="NovelaActiveTabsEnum.CHAPTERS"
                >
                  <ChapterList
                    :key="`chapters-${chaptersKey}`"
                    :volumes="novela.volumes"
                    :can-manage="!!user"
                    :novela-id="novela.id"
                    @volume-deleted="onVolumeDeleted"
                    @volume-added="onVolumeAdded"
                  />
                </div>
                <div
                  v-else-if="activeTab === NovelaActiveTabsEnum.COMMENTS"
                  :key="NovelaActiveTabsEnum.COMMENTS"
                >
                  <NovelaComments
                    v-if="comments !== null"
                    :comments="comments"
                    :novela-id="novela.id"
                    :is-loading="commentsLoading"
                    :novela-title="novela.title"
                    :total-comments="totalComments"
                  />
                </div>
              </transition>
            </div>
          </div>
        </div>
      </div>
    </main>

    <div v-else class="grow flex items-center justify-center">
      <div
        class="w-12 h-12 border-4 border-zinc-300 border-t-yellow-500 rounded-full animate-spin"
      ></div>
    </div>

    <NovelaSettings
      v-if="hasPermission(KyokusuAppRole.MODERATOR)"
      v-model="isOpenNovelaSettings"
      :novela="novela!"
      @updated="updatedNovela"
    />

    <ModalConfirm
      v-if="novela"
      v-model="isOpenDeleteConfirm"
      title="Удаление новелы"
      :description="`Вы уверены, что хотите удалить «${novela.title}»? Это действие необратимо — все тома, главы, изображения и комментарии будут удалены.`"
      confirm-text="Удалить"
      cancel-text="Отмена"
      @confirm="deleteNovela"
    />
  </div>
</template>

<style scoped>
.shadow-text {
  text-shadow: 4px 4px 2px rgba(0, 0, 0, 0.6);
}
.fade-enter-active,
.fade-leave-active {
  transition: all 0.3s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

.no-scrollbar::-webkit-scrollbar {
  display: none;
}

.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
