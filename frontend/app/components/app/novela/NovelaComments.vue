<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useNovelaComments } from '@/composables/api/novela/useNovelaComments'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'
import CommentItem from './CommentItem.vue'
import AuthRequiredModal from '~/components/common/AuthRequiredModal.vue'
import { RichText as BaseRichTextEditor } from '@kyokusu-ui/vue'
import { NOVELA_MAX_COMMENT_LENGTH } from '~/constants/data'
import { useUserActivity } from '~/composables/api/profile/useUserActivity'
import { ACTIVITY_TYPES } from '~/constants/user-activity'
import FilterBar, { type FilterConfig } from '@/components/ui/FilterBar.vue'
import type { NovelaCommentResponse } from '@/types/backend/novela'

const props = defineProps<{
  novelaId: number
  comments: NovelaCommentResponse[]
  isLoading: boolean
  novelaTitle: string
}>()

const route = useRoute()
const { user, isAuthenticated } = useAuthStore()
const { notify } = useNotificationStore()
const novelaCommentsApi = useNovelaComments()

const { createUserActivity } = useUserActivity()

const isAuthModalOpen = ref(false)
const commentText = ref('')
const replyTo = ref<{ id: number; name: string } | null>(null)
const editingComment = ref<{ id: number; content: string } | null>(null)

const sortedComments = ref<NovelaCommentResponse[]>([...props.comments])

const teamFilterConfig: FilterConfig = {
  defaultField: 'created_at',
  sortOptions: [
    {
      field: 'created_at',
      label: 'По новизне',
      icon: 'ph:calendar',
      compare: (a, b) =>
        new Date((a as NovelaCommentResponse).created_at).getTime() -
        new Date((b as NovelaCommentResponse).created_at).getTime(),
    },
    {
      field: 'like_count',
      label: 'По популярности',
      icon: 'ph:heart',
      compare: (a, b) =>
        (a as NovelaCommentResponse).like_count - (b as NovelaCommentResponse).like_count,
    },
  ],
}

onMounted(() => {
  if (route.hash) {
    setTimeout(() => {
      const el = document.querySelector(route.hash)
      if (el) {
        el.scrollIntoView({ behavior: 'smooth', block: 'center' })
        el.classList.add(
          'ring-2',
          'ring-yellow-500',
          'ring-offset-2',
          'ring-offset-white',
          'dark:ring-offset-zinc-900',
        )
        setTimeout(() => {
          el.classList.remove(
            'ring-2',
            'ring-yellow-500',
            'ring-offset-2',
            'ring-offset-white',
            'dark:ring-offset-zinc-900',
          )
        }, 2000)
      }
    }, 300)
  }
})

const scrolltoEditor = () => {
  nextTick(() => {
    document
      .getElementById('comment-editor-section')
      ?.scrollIntoView({ behavior: 'smooth', block: 'center' })
  })
}

const onReplyRequest = (payload: { id: number; name: string }) => {
  editingComment.value = null
  replyTo.value = payload
  scrolltoEditor()
}

const onUpdateRequest = (comment: { id: number; content: string }) => {
  replyTo.value = null
  editingComment.value = comment
  commentText.value = comment.content
  scrolltoEditor()
}

const resetForm = () => {
  editingComment.value = null
  replyTo.value = null
  commentText.value = ''
}

const handleAction = async () => {
  if (!isAuthenticated) {
    isAuthModalOpen.value = true
    return
  }

  if (!commentText.value.trim() || props.isLoading) return

  try {
    if (editingComment.value) {
      await novelaCommentsApi.updateComment(editingComment.value.id, commentText.value)
      notify({ title: 'Успех', content: 'Комментарий обновлен', type: 'success' })
      await createUserActivity({
        user_id: user!.id,
        activity_type: ACTIVITY_TYPES.COMMENT_UPDATE,
        target_id: props.novelaId,
        metadata: {
          novela_title: props.novelaTitle,
          desc: 'Пользователь обновил комментарий',
        },
      })
    } else {
      await novelaCommentsApi.addComment({
        novela_id: props.novelaId,
        content: commentText.value,
        parent_id: replyTo.value?.id || null,
      })
      notify({ title: 'Успех', content: 'Комментарий опубликован', type: 'success' })
      await createUserActivity({
        user_id: user!.id,
        activity_type: ACTIVITY_TYPES.COMMENT_ADDED,
        target_id: props.novelaId,
        metadata: {
          novela_title: props.novelaTitle,
          desc: 'Пользователь добавил комментарий',
        },
      })
    }

    resetForm()
  } catch (e: any) {
    notify({ title: 'Ошибка', content: e.message, type: 'error' })
  }
}

const handleDelete = async (commentId: number) => {
  if (!confirm('Вы уверены, что хотите удалить комментарий?')) return

  try {
    await novelaCommentsApi.deleteComment(commentId)
    novelaCommentsApi.comments.value = novelaCommentsApi.comments.value.filter(
      (c) => c.id !== commentId,
    )
    notify({ title: 'Удалено', content: 'Комментарий успешно удален', type: 'success' })
    await createUserActivity({
      user_id: user!.id,
      activity_type: ACTIVITY_TYPES.COMMENT_REMOVE,
      target_id: props.novelaId,
      metadata: {
        novela_title: props.novelaTitle,
        desc: 'Пользователь удалил комментарий',
      },
    })
  } catch (e: any) {
    notify({ title: 'Ошибка', content: e.message, type: 'error' })
  }
}

const handleLike = async (commentId: number) => {
  if (!isAuthenticated) {
    isAuthModalOpen.value = true
    return
  }

  try {
    await novelaCommentsApi.setCommentLike(commentId)
    novelaCommentsApi.comments.value = novelaCommentsApi.comments.value.map((c) => {
      if (c.id === commentId) {
        c.like_count++
        c.has_like = true
      }
      return c
    })
    await createUserActivity({
      user_id: user!.id,
      activity_type: ACTIVITY_TYPES.COMMENT_LIKE,
      target_id: props.novelaId,
      metadata: {
        novela_title: props.novelaTitle,
        desc: 'Пользователь поставил лайк комментарию',
      },
    })
  } catch (e: any) {
    notify({ title: 'Ошибка', content: e.message, type: 'error' })
  }
}

const handleUnsetLike = async (commentId: number) => {
  if (!isAuthenticated) {
    isAuthModalOpen.value = true
    return
  }

  try {
    await novelaCommentsApi.unsetCommentLike(commentId)
    novelaCommentsApi.comments.value = novelaCommentsApi.comments.value.map((c) => {
      if (c.id === commentId) {
        c.like_count--
        c.has_like = null
      }
      return c
    })
  } catch (e: any) {
    notify({ title: 'Ошибка', content: e.message, type: 'error' })
  }
}

const handleReport = async ({ commentId, reason }: { commentId: number; reason: string }) => {
  if (!isAuthenticated) {
    isAuthModalOpen.value = true
    return
  }

  try {
    await novelaCommentsApi.reportComment(commentId, reason)
    notify({
      title: 'Жалоба отправлена',
      content: 'Модераторы рассмотрят ваше обращение',
      type: 'success',
    })

    await createUserActivity({
      user_id: user!.id,
      activity_type: ACTIVITY_TYPES.COMMENT_REPORT,
      target_id: props.novelaId,
      metadata: {
        novela_title: props.novelaTitle,
        desc: `Пользователь пожаловался на комментарий (ID: ${commentId})`,
      },
    })
    const commentWithReport = novelaCommentsApi.comments.value.find((c) => c.id === commentId)

    if (commentWithReport) {
      commentWithReport.has_report = true
    }
  } catch (e: any) {
    notify({ title: 'Ошибка', content: e.message, type: 'error' })
  }
}

const handleReportCancel = async (id: number) => {
  if (!isAuthenticated) {
    isAuthModalOpen.value = true
    return
  }

  try {
    await novelaCommentsApi.cancelCommentReport(id)
    notify({ title: 'Жалоба отменена', content: 'Жалоба на комментарий отменена', type: 'success' })
  } catch (e: any) {
    notify({ title: 'Ошибка', content: e.message, type: 'error' })
  }

  const commentWithReport = novelaCommentsApi.comments.value.find((c) => c.id === id)

  if (commentWithReport) {
    commentWithReport.has_report = false
  }
}
</script>

<template>
  <div class="space-y-10">
    <section id="comment-editor-section" class="space-y-4">
      <Transition name="fade" mode="out-in">
        <div
          v-if="replyTo || editingComment"
          class="flex items-center justify-between p-3 rounded-2xl border transition-all"
          :class="
            editingComment
              ? 'bg-blue-500/10 border-blue-500/20'
              : 'bg-yellow-500/10 border-yellow-500/20'
          "
        >
          <div class="flex items-center gap-2">
            <Icon
              :name="editingComment ? 'ph:pencil-line-bold' : 'ph:arrows-out-line-horizontal-bold'"
              :class="editingComment ? 'text-blue-500' : 'text-yellow-600'"
            />
            <span
              class="text-[10px] font-black uppercase tracking-widest"
              :class="editingComment ? 'text-blue-500' : 'text-yellow-600'"
            >
              {{ editingComment ? 'Редактирование' : `Ответ: ${replyTo?.name}` }}
            </span>
          </div>

          <button
            @click="resetForm"
            class="cursor-pointer hover:rotate-90 transition-transform p-1"
          >
            <Icon
              name="ph:x-bold"
              size="14"
              :class="editingComment ? 'text-blue-500' : 'text-yellow-600'"
            />
          </button>
        </div>
      </Transition>

      <BaseRichTextEditor
        id="comment-editor"
        v-model="commentText"
        placeholder="Ваше мнение о произведении..."
        :disabled="props.isLoading"
        :maxLength="NOVELA_MAX_COMMENT_LENGTH"
      />

      <div class="flex justify-end gap-3">
        <button
          v-if="editingComment || replyTo"
          @click="resetForm"
          class="px-6 py-3 text-zinc-400 font-bold text-xs uppercase tracking-widest hover:text-zinc-900 dark:hover:text-white transition-colors cursor-pointer"
        >
          Отмена
        </button>
        <button
          @click="handleAction"
          :disabled="!commentText.trim() || props.isLoading"
          class="px-10 py-3 bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 font-black rounded-2xl text-xs uppercase tracking-widest active:scale-95 transition-all flex items-center gap-3 disabled:opacity-50 shadow-xl shadow-zinc-500/10 cursor-pointer"
        >
          <Icon v-if="props.isLoading" name="ph:circle-notch-bold" class="animate-spin" />
          {{ editingComment ? 'Сохранить' : 'Отправить' }}
        </button>
      </div>
    </section>

    <div class="relative min-h-75">
      <div v-if="props.isLoading && !props.comments.length" class="py-20 flex justify-center">
        <Icon name="ph:circle-notch-bold" size="48" class="animate-spin text-zinc-300" />
      </div>

      <div
        v-else-if="!props.comments.length"
        class="py-24 flex flex-col items-center justify-center border-2 border-dashed border-zinc-200 dark:border-zinc-800 rounded-[3rem] text-zinc-400 bg-zinc-50/50 dark:bg-zinc-900/10"
      >
        <Icon name="ph:chat-circle-dots-light" size="64" class="opacity-20 mb-4" />
        <p class="font-bold uppercase tracking-[0.2em] text-[10px]">Здесь пока нет обсуждений</p>
        <p class="text-xs opacity-60 mt-1">Будьте первым, кто поделится мыслями</p>
      </div>

      <TransitionGroup v-else name="list" tag="div" class="space-y-6">
        <FilterBar
          v-if="!isLoading && props.comments.length > 0"
          :config="teamFilterConfig"
          :items="props.comments"
          v-model="sortedComments"
        />
        <CommentItem
          v-for="comment in sortedComments"
          :key="comment.id"
          :comment="comment"
          :user-id="user?.id || 0"
          @reply="onReplyRequest"
          @update="onUpdateRequest"
          @delete="handleDelete"
          @like="handleLike"
          @unset-like="handleUnsetLike"
          @report="handleReport"
          @cancel-report="handleReportCancel"
        />
      </TransitionGroup>
    </div>

    <AuthRequiredModal
      v-if="isAuthModalOpen"
      v-model="isAuthModalOpen"
      action-text="взаимодействовать с комментариями"
    />
  </div>
</template>

<style scoped>
.list-enter-active,
.list-leave-active {
  transition: all 0.4s ease;
}
.list-enter-from {
  opacity: 0;
  transform: translateY(30px);
}
.list-leave-to {
  opacity: 0;
  transform: scale(0.95);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
