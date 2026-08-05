<script setup lang="ts">
import { useRoute } from 'vue-router'
import { useAuthors } from '~/composables/api/authors/useAuthors'

import AuthorCard from '~/components/app/authors/AuthorCard.vue'
import AuthorSameWork from '~/components/app/authors/AuthorSameWork.vue'

const route = useRoute()

const { isSearching, getAuthorById } = useAuthors()

const authorId = computed(() => route.params.id as string)

const { data: author, status } = await useAsyncData(`author-${authorId.value}`, () =>
  getAuthorById(authorId.value),
)

useSeoMeta({
  title: () =>
    author.value?.name ? `${author.value.name} - Авторы - KyokusuLib` : 'Автор - KyokusuLib',
  ogTitle: () => author.value?.name,
  ogImage: () => author.value?.picture,
})
</script>

<template>
  <div
    class="min-h-screen bg-zinc-50 dark:bg-[#0f0f0f] text-zinc-900 dark:text-zinc-200 flex flex-col font-sans transition-colors duration-300"
  >
    <div class="w-full max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12 space-y-6">
      <div v-if="status === 'pending'" class="flex justify-center items-center py-32">
        <Icon name="ph:spinner-gap-bold" size="48" class="animate-spin text-yellow-500 mb-4" />
      </div>
      <template v-else-if="author">
        <AuthorCard :author="author" />
        <AuthorSameWork :authorId="author.id" :authorName="author.name" />
      </template>
      <div
        v-else
        class="flex flex-col items-center justify-center py-32 text-center border-2 border-dashed border-zinc-200 dark:border-zinc-800 rounded-3xl bg-white dark:bg-[#18181b]"
      >
        <Icon name="ph:user-minus-bold" size="64" class="text-zinc-300 dark:text-zinc-700 mb-4" />
        <h2 class="text-xl font-bold text-zinc-900 dark:text-zinc-100 mb-2">Автор не найден</h2>
        <p class="text-sm font-medium text-zinc-500">
          Запрошенный автор не существует или был удален.
        </p>
      </div>
    </div>
  </div>
</template>
