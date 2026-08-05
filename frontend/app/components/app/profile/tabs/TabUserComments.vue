<script setup lang="ts">
import { ref } from 'vue'
import { parseDateToLocale } from '@/utils/date'
import { useUserComments } from '~/composables/api/profile/useUserComments'
import FilterBar from '~/components/ui/FilterBar.vue'
import type { FilterConfig } from '~/components/ui/FilterBar.vue'
import type { UserCommentResponse } from '~/types/backend/comments'

const props = defineProps<{ userId: number }>()
const { comments, fetchUserComments } = useUserComments()
const sortedComments = ref<UserCommentResponse[]>([])

const commentFilterConfig: FilterConfig = {
  defaultField: 'created_at',
  defaultOrder: 'desc',
  sortOptions: [
    {
      field: 'created_at',
      label: 'По дате',
      icon: 'ph:calendar',
      compare: (a, b) =>
        new Date((a as UserCommentResponse).created_at).getTime() -
        new Date((b as UserCommentResponse).created_at).getTime(),
    },
    {
      field: 'novela_title',
      label: 'По новелле',
      icon: 'ph:book-open-bold',
      compare: (a, b) =>
        (a as UserCommentResponse).novela_title.localeCompare(
          (b as UserCommentResponse).novela_title,
        ),
    },
    {
      field: 'content_length',
      label: 'По длине',
      icon: 'ph:text-align-left',
      compare: (a, b) =>
        (a as UserCommentResponse).content.length - (b as UserCommentResponse).content.length,
    },
  ],
}

await useAsyncData(`user-comments-${props.userId}`, () => fetchUserComments(props.userId))
</script>

<template>
  <div class="space-y-6">
    <FilterBar
      v-if="!!comments.length"
      :items="comments"
      v-model="sortedComments"
      :config="commentFilterConfig"
    />

    <div v-if="!!comments.length" class="flex flex-col gap-4">
      <div
        v-for="comment in sortedComments"
        :key="comment.id"
        @click="navigateTo(`/novela/${comment.novela_id}?tab=comments#comment-${comment.id}`)"
        class="group relative bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-5 transition-all hover:border-yellow-500/30 hover:shadow-lg hover:shadow-yellow-500/5 cursor-pointer"
      >
        <div class="flex items-center justify-between mb-3 gap-8">
          <NuxtLink
            :to="`/novela/${comment.novela_id}`"
            @click.stop
            class="flex items-center gap-2 group/link"
          >
            <div
              class="w-8 h-8 rounded-lg bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center group-hover/link:bg-yellow-500 transition-colors"
            >
              <Icon
                name="ph:book-open-bold"
                size="16"
                class="text-zinc-500 dark:text-zinc-400 group-hover/link:text-zinc-900"
              />
            </div>
            <span
              class="text-sm font-black text-zinc-900 dark:text-zinc-100 group-hover/link:text-yellow-600 transition-colors truncate max-w-50 sm:max-w-md"
            >
              {{ comment.novela_title }}
            </span>
          </NuxtLink>

          <span class="text-[10px] font-bold text-zinc-400 uppercase tracking-widest">
            {{
              comment.created_at < comment.updated_at
                ? `Изменено: ${parseDateToLocale(comment.updated_at)}`
                : parseDateToLocale(comment.created_at)
            }}
          </span>
        </div>

        <div class="relative">
          <Icon
            name="ph:quotes-fill"
            size="32"
            class="absolute -top-2 -left-2 text-zinc-100 dark:text-zinc-800 -z-10 opacity-50"
          />
          <div
            class="text-[14px] ml-1 leading-relaxed text-zinc-600 dark:text-zinc-300 prose dark:prose-invert max-w-none line-clamp-3 group-hover:line-clamp-none transition-all antialiased"
            v-html="comment.content"
          ></div>
        </div>

        <div
          class="flex justify-end mt-4 pt-4 border-t border-zinc-100 dark:border-zinc-800 opacity-0 group-hover:opacity-100 transition-opacity"
        >
          <NuxtLink
            :to="`/novela/${comment.novela_id}?tab=comments#comment-${comment.id}`"
            @click.stop
            class="text-[10px] font-black uppercase tracking-widest text-yellow-600 hover:text-yellow-500 flex items-center gap-1.5"
          >
            Перейти к обсуждению
            <Icon name="ph:arrow-right-bold" />
          </NuxtLink>
        </div>
      </div>
    </div>

    <div
      v-else
      class="py-24 flex flex-col items-center justify-center border-2 border-dashed border-zinc-200 dark:border-zinc-800 rounded-[3rem] text-zinc-400 bg-zinc-50/50 dark:bg-zinc-900/10"
    >
      <Icon name="ph:chat-teardrop-dots-light" size="64" class="opacity-20 mb-4" />
      <p class="font-bold uppercase tracking-[0.2em] text-[10px]">История комментариев пуста</p>
    </div>
  </div>
</template>
