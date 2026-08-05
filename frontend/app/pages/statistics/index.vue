<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useTotalStatistics } from '~/composables/api/statistics/total-stats'
import { StatisticsPeriodEnum } from '@/types/enums/statistics'
import { staticImage } from '@/utils/str'
import { Separator } from '@kyokusu-ui/vue'

useSeoMeta({ title: 'Статистика - KyokusuLib' })

const { novelaTotalStats, total, fetchTotalStatistics } = useTotalStatistics()

const currentPeriod = ref<StatisticsPeriodEnum | undefined>(undefined)
const currentPage = ref(1)
const limit = 20
const isLoading = ref(true)

const offset = computed(() => (currentPage.value - 1) * limit)
const totalPages = computed(() => Math.ceil(total.value / limit) || 1)

const loadData = async () => {
  isLoading.value = true
  await fetchTotalStatistics(currentPeriod.value, limit, offset.value)
  isLoading.value = false
}

onMounted(() => {
  loadData()
})

const setPeriod = (period: StatisticsPeriodEnum | undefined) => {
  if (currentPeriod.value === period) return
  currentPeriod.value = period
  currentPage.value = 1
  loadData()
}

const nextPage = () => {
  if (currentPage.value >= totalPages.value) return
  currentPage.value++
  loadData()
}

const prevPage = () => {
  if (currentPage.value <= 1) return
  currentPage.value--
  loadData()
}

const roundTo = (num: number, decimals: number = 1) => {
  return Math.round(num * Math.pow(10, decimals)) / Math.pow(10, decimals)
}
</script>

<template>
  <div
    class="min-h-screen bg-zinc-50 dark:bg-[#0f0f0f] text-zinc-900 dark:text-zinc-200 transition-colors duration-300"
  >
    <div class="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
      <div class="mb-8">
        <h1
          class="text-3xl sm:text-4xl font-black uppercase tracking-tight text-zinc-900 dark:text-white mb-4"
        >
          Статистика
        </h1>
        <Separator />
      </div>

      <div
        class="bg-white dark:bg-zinc-900/40 rounded-[2.5rem] border border-zinc-200 dark:border-zinc-800 p-6 shadow-sm mb-8"
      >
        <div class="flex flex-wrap gap-2 mb-6">
          <button
            @click="setPeriod(undefined)"
            :class="[
              'px-4 py-2 rounded-xl text-sm font-medium transition-colors',
              currentPeriod === undefined
                ? 'bg-zinc-200 dark:bg-zinc-700 text-zinc-900 dark:text-white ring-1 ring-zinc-300 dark:ring-zinc-600'
                : 'text-zinc-500 hover:bg-zinc-100 dark:text-zinc-400 dark:hover:bg-zinc-800',
            ]"
          >
            За всё время
          </button>
          <button
            @click="setPeriod(StatisticsPeriodEnum.Day)"
            :class="[
              'px-4 py-2 rounded-xl text-sm font-medium transition-colors',
              currentPeriod === StatisticsPeriodEnum.Day
                ? 'bg-zinc-200 dark:bg-zinc-700 text-zinc-900 dark:text-white ring-1 ring-zinc-300 dark:ring-zinc-600'
                : 'text-zinc-500 hover:bg-zinc-100 dark:text-zinc-400 dark:hover:bg-zinc-800',
            ]"
          >
            За день
          </button>
          <button
            @click="setPeriod(StatisticsPeriodEnum.Week)"
            :class="[
              'px-4 py-2 rounded-xl text-sm font-medium transition-colors',
              currentPeriod === StatisticsPeriodEnum.Week
                ? 'bg-zinc-200 dark:bg-zinc-700 text-zinc-900 dark:text-white ring-1 ring-zinc-300 dark:ring-zinc-600'
                : 'text-zinc-500 hover:bg-zinc-100 dark:text-zinc-400 dark:hover:bg-zinc-800',
            ]"
          >
            За неделю
          </button>
          <button
            @click="setPeriod(StatisticsPeriodEnum.Month)"
            :class="[
              'px-4 py-2 rounded-xl text-sm font-medium transition-colors',
              currentPeriod === StatisticsPeriodEnum.Month
                ? 'bg-zinc-200 dark:bg-zinc-700 text-zinc-900 dark:text-white ring-1 ring-zinc-300 dark:ring-zinc-600'
                : 'text-zinc-500 hover:bg-zinc-100 dark:text-zinc-400 dark:hover:bg-zinc-800',
            ]"
          >
            За месяц
          </button>
        </div>

        <div v-if="isLoading" class="flex justify-center items-center py-12">
          <Icon name="ph:spinner-gap-bold" size="32" class="animate-spin text-zinc-400" />
          <span class="ml-3 text-zinc-500 font-medium">Загрузка...</span>
        </div>

        <div
          v-else-if="novelaTotalStats.length === 0"
          class="flex flex-col items-center justify-center py-16 text-center"
        >
          <Icon name="ph:chart-bar-bold" size="48" class="text-zinc-300 dark:text-zinc-700 mb-4" />
          <p class="text-zinc-500 dark:text-zinc-400 font-medium">Нет данных для отображения</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead>
              <tr
                class="border-b border-zinc-200 dark:border-zinc-800 text-sm font-semibold text-zinc-500 dark:text-zinc-400"
              >
                <th class="py-4 px-4 w-12 text-center">#</th>
                <th class="py-4 px-4">Новелла</th>
                <th class="py-4 px-4 text-center">Рейтинг</th>
                <th class="py-4 px-4 text-center">
                  <Icon name="ph:bookmark-simple-bold" size="18" class="inline-block" />
                </th>
                <th class="py-4 px-4 text-center">
                  <Icon name="ph:eye-bold" size="18" class="inline-block" />
                </th>
                <th class="py-4 px-4 text-center">
                  <Icon name="ph:star-bold" size="18" class="inline-block" />
                </th>
                <th class="py-4 px-4 text-center">
                  <Icon name="ph:chat-circle-bold" size="18" class="inline-block" />
                </th>
                <th class="py-4 px-4 text-center">
                  <Icon name="ph:books-bold" size="18" class="inline-block" />
                </th>
                <th class="py-4 px-4 text-center">
                  <Icon name="ph:file-text-bold" size="18" class="inline-block" />
                </th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(stat, index) in novelaTotalStats"
                :key="stat.novela.novelaId"
                class="border-b border-zinc-100 dark:border-zinc-800/50 hover:bg-zinc-50 dark:hover:bg-zinc-800/30 transition-colors"
              >
                <td class="py-3 px-4 text-center font-medium text-zinc-500">
                  {{ offset + index + 1 }}
                </td>
                <td class="py-3 px-4">
                  <NuxtLink
                    :to="`/novela/${stat.novela.novelaId}`"
                    class="flex items-center gap-3 group"
                  >
                    <img
                      :src="staticImage(stat.novela.posterUrl)"
                      :alt="stat.novela.title"
                      class="w-10 h-14 rounded-lg object-cover shadow-sm group-hover:shadow-md transition-shadow"
                    />
                    <span
                      class="font-semibold text-zinc-900 dark:text-zinc-100 group-hover:text-yellow-500 dark:group-hover:text-yellow-400 transition-colors line-clamp-2"
                    >
                      {{ stat.novela.title }}
                    </span>
                  </NuxtLink>
                </td>
                <td class="py-3 px-4 text-center">
                  <div
                    class="flex items-center justify-center gap-1 font-bold text-zinc-900 dark:text-zinc-100"
                  >
                    <Icon name="ph:star-fill" size="16" class="text-yellow-500" />
                    {{ roundTo(stat.rating, 1) }}
                  </div>
                </td>
                <td class="py-3 px-4 text-center text-zinc-600 dark:text-zinc-400 font-medium">
                  {{ stat.bookmarkCount }}
                </td>
                <td class="py-3 px-4 text-center text-zinc-600 dark:text-zinc-400 font-medium">
                  {{ stat.readCount }}
                </td>
                <td class="py-3 px-4 text-center text-zinc-600 dark:text-zinc-400 font-medium">
                  {{ stat.ratingCount }}
                </td>
                <td class="py-3 px-4 text-center text-zinc-600 dark:text-zinc-400 font-medium">
                  {{ stat.commentCount }}
                </td>
                <td class="py-3 px-4 text-center text-zinc-600 dark:text-zinc-400 font-medium">
                  {{ stat.volumeCount }}
                </td>
                <td class="py-3 px-4 text-center text-zinc-600 dark:text-zinc-400 font-medium">
                  {{ stat.chapterCount }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <div
          class="flex items-center justify-between mt-6 pt-6 border-t border-zinc-200 dark:border-zinc-800"
        >
          <button
            @click="prevPage"
            :disabled="currentPage <= 1 || isLoading"
            class="px-4 py-2 rounded-xl text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed bg-zinc-100 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 text-zinc-900 dark:text-zinc-100"
          >
            < Назад
          </button>

          <span class="text-sm font-medium text-zinc-500 dark:text-zinc-400">
            {{ currentPage }} / {{ totalPages }}
          </span>

          <button
            @click="nextPage"
            :disabled="currentPage >= totalPages || isLoading"
            class="px-4 py-2 rounded-xl text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed bg-zinc-100 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 text-zinc-900 dark:text-zinc-100"
          >
            Вперёд >
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
