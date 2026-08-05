<script setup lang="ts">
import { ref, computed } from 'vue'
import type { NovelaCommentResponse } from '@/types/backend/novela'
import { staticImage } from '@/utils/str'
import { parseDateToLocale } from '@/utils/date'
import { Tooltip as BaseToolTip } from '@kyokusu-ui/vue'
import CommentReportMenu from './CommentReportMenu.vue'

interface Props {
  comment: NovelaCommentResponse
  isReply?: boolean
  userId?: number
}

const props = defineProps<Props>()

const emit = defineEmits([
  'reply',
  'delete',
  'update',
  'like',
  'unset-like',
  'report',
  'cancel-report',
])
const isReportMenuOpen = ref(false)
const isExpanded = ref(false)

const totalReplyCount = computed(() => {
  const getCount = (list: any[]): number =>
    list.reduce((total, item) => total + 1 + getCount(item.replies || []), 0)

  return getCount(props.comment.replies || [])
})

const handleLike = () => {
  if (!props.comment.has_like) {
    emit('like', props.comment.id)
  } else {
    emit('unset-like', props.comment.id)
  }
}

const onReportSubmit = (payload: any) => {
  emit('report', payload)
}

const onReportCancel = (id: number) => {
  emit('cancel-report', id)
}

const openReportMenu = () => {
  isReportMenuOpen.value = !isReportMenuOpen.value
}
</script>

<template>
  <div
    :id="`comment-${comment.id}`"
    :class="[
      'flex gap-4 group transition-all',
      isReply
        ? 'mt-4'
        : 'bg-white dark:bg-zinc-900/50 p-4 border border-zinc-200 dark:border-zinc-800 rounded-3xl shadow-sm',
    ]"
  >
    <img
      :src="staticImage(comment.user.picture)"
      class="rounded-full shrink-0 object-cover bg-zinc-200 dark:bg-zinc-700 border-2 border-transparent hover:border-zinc-400 dark:hover:border-zinc-700 shadow-sm cursor-pointer transition-transform hover:scale-105"
      :class="isReply ? 'w-8 h-8' : 'w-11 h-11'"
      @click="navigateTo(`/profile/${comment.user.id}`)"
    />

    <div class="flex-1 min-w-0 space-y-1.5">
      <div class="flex flex-row justify-between items-start">
        <div class="flex items-center gap-3">
          <span
            class="font-black text-zinc-900 dark:text-zinc-100 cursor-pointer text-sm hover:text-yellow-600 transition-colors"
            @click="navigateTo(`/profile/${comment.user.id}`)"
          >
            {{ comment.user.name }}
          </span>
          <span class="text-[9px] font-bold text-zinc-400 uppercase tracking-widest cursor-default">
            {{ parseDateToLocale(comment.created_at) }}
          </span>
        </div>

        <div class="flex flex-row gap-3">
          <div v-if="comment.user.id === userId" class="flex flex-row gap-3">
            <BaseToolTip text="Редактировать" position="left">
              <Icon
                name="ph:pencil-bold"
                size="16"
                class="cursor-pointer text-zinc-400 hover:text-blue-500 transition-colors"
                @click="emit('update', { id: comment.id, content: comment.content })"
              />
            </BaseToolTip>
            <BaseToolTip text="Удалить" position="left">
              <Icon
                name="ph:trash-bold"
                size="16"
                class="cursor-pointer text-zinc-400 hover:text-red-500 transition-colors"
                @click="emit('delete', comment.id)"
              />
            </BaseToolTip>
          </div>
          <div class="flex flex-row gap-3">
            <div class="flex flex-row justify-center gap-1">
              <BaseToolTip text="Лайк" position="left">
                <Icon
                  :name="comment.has_like ? 'ph:heart-fill' : 'ph:heart-straight-bold'"
                  size="16"
                  :class="[
                    'transition-all duration-300 cursor-pointer',
                    comment.has_like ? 'text-yellow-500 scale-110' : 'text-zinc-400',
                  ]"
                  @click="handleLike"
                />
              </BaseToolTip>
              <span class="text-[10px] font-bold text-zinc-400 cursor-default mt-0.5">
                {{ comment.like_count }}
              </span>
            </div>
            <div
              v-if="comment.user.id !== userId"
              class="flex flex-row justify-center gap-1 relative"
            >
              <BaseToolTip
                :text="comment.has_report ? 'Жалоба отправлена' : 'Пожаловаться'"
                position="left"
              >
                <Icon
                  :name="comment.has_report ? 'ph:flag-fill' : 'ph:flag-bold'"
                  size="16"
                  class="cursor-pointer transition-colors report-toggle-btn"
                  :class="[
                    isReportMenuOpen ? 'text-red-500' : '',
                    comment.has_report !== null && comment.has_report === true
                      ? 'text-red-500'
                      : 'text-zinc-400 hover:text-red-500',
                  ]"
                  @click="openReportMenu"
                />
              </BaseToolTip>

              <CommentReportMenu
                v-if="isReportMenuOpen"
                :comment-id="comment.id"
                :reported="Boolean(comment.has_report)"
                @close="isReportMenuOpen = false"
                @submit="onReportSubmit"
                @cancel="onReportCancel"
              />
            </div>
          </div>
        </div>
      </div>

      <div
        class="text-[14px] leading-relaxed text-zinc-700 dark:text-zinc-300 prose dark:prose-invert max-w-none break-words"
        v-html="comment.content"
      ></div>

      <div class="flex items-center gap-4">
        <button
          @click="emit('reply', { id: comment.id, name: comment.user.name })"
          class="flex items-center gap-1.5 text-[10px] font-black uppercase tracking-widest text-zinc-400 hover:text-yellow-600 transition-colors cursor-pointer"
        >
          Ответить
          <Icon name="ph:arrow-bend-up-left-bold" size="14" />
        </button>

        <button
          v-if="comment.replies?.length"
          @click="isExpanded = !isExpanded"
          class="flex items-center gap-1.5 text-[10px] font-black uppercase tracking-widest text-yellow-600/80 hover:text-yellow-500 transition-colors cursor-pointer"
        >
          <Icon :name="isExpanded ? 'ph:caret-up-bold' : 'ph:caret-down-bold'" size="14" />
          {{ isExpanded ? 'Скрыть' : `Ответы (${totalReplyCount})` }}
        </button>
      </div>

      <div
        class="grid transition-[grid-template-rows] duration-300 ease-in-out"
        :class="isExpanded ? 'grid-template-rows-[1fr]' : 'grid-template-rows-[0fr]'"
      >
        <div class="overflow-hidden">
          <div class="mt-2 border-l-2 border-zinc-100 dark:border-zinc-800 pl-4 space-y-2">
            <CommentItem
              v-for="reply in comment.replies"
              :key="reply.id"
              :comment="reply"
              :user-id="userId"
              is-reply
              @reply="(val) => emit('reply', val)"
              @delete="(id) => emit('delete', id)"
              @update="(payload) => emit('update', payload)"
              @like="(id) => emit('like', id)"
              @unset-like="(id) => emit('unset-like', id)"
              @report="(val) => emit('report', val)"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.grid-template-rows-\[0fr\] {
  grid-template-rows: 0fr;
}
.grid-template-rows-\[1fr\] {
  grid-template-rows: 1fr;
}
</style>
