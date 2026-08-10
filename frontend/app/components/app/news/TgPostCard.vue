<script setup lang="ts">
import { Card, Button, TeleportedTooltip } from '@kyokusu-ui/vue'
import PostImageSlider from '@/components/app/news/PostImageSlider.vue'
import type { BackendTgPost } from '@/types/backend/tg-post'
import type { NewsSource } from '@/types/frontend/news'
import { fmtDateTime } from '@/utils/date'
import { sanitizeHtml } from '@/utils/sanitize'
import { TG_CHANNEL_NAME } from '@/constants/news'

const props = defineProps<{
  post: BackendTgPost
  source: NewsSource
  postUrl: string
  canDelete: boolean
}>()

const emit = defineEmits<{
  delete: [post: BackendTgPost]
}>()

const safeText = computed(() => sanitizeHtml(props.post.text))
const formattedDate = computed(() => fmtDateTime(props.post.createdAt))
</script>

<template>
  <Card variant="outline" padding="sm" shadow class="post-card transition-shadow hover:shadow-lg">
    <div class="flex flex-col gap-3">
      <div class="flex items-start justify-between gap-2">
        <NuxtLink :to="postUrl" target="_blank">
          <div class="flex items-center gap-2.5 min-w-0">
            <div
              class="flex h-9 w-9 shrink-0 items-center justify-center rounded-full text-white"
              :style="source.styles"
            >
              <Icon :name="source.icon" size="18" />
            </div>
            <div class="min-w-0">
              <div class="font-extrabold leading-tight">{{ TG_CHANNEL_NAME }}</div>
              <p class="flex items-center gap-1 text-xs text-zinc-500 dark:text-zinc-400">
                <Icon name="ph:clock" size="14" />
                <span>
                  {{ formattedDate }}
                </span>
              </p>
            </div>
          </div>
        </NuxtLink>

        <div class="flex items-center gap-1">
          <TeleportedTooltip text="Открыть в Telegram">
            <a
              :href="postUrl"
              target="_blank"
              rel="noopener noreferrer"
              aria-label="Открыть пост в Telegram"
              class="flex h-8 w-8 items-center justify-center rounded-lg text-zinc-400 transition-colors hover:text-zinc-700 dark:hover:text-zinc-200"
            >
              <Icon name="ph:arrow-up-right-bold" size="16" />
            </a>
          </TeleportedTooltip>
          <Button
            v-if="canDelete"
            variant="ghost"
            size="icon"
            aria-label="Удалить пост"
            @click="emit('delete', post)"
          >
            <Icon name="ph:trash-bold" size="16" />
          </Button>
        </div>
      </div>

      <div
        v-if="post.text"
        class="tg-content rounded-xl border-2 border-zinc-500/40 bg-zinc-500/10 p-3 text-sm font-medium whitespace-pre-line wrap-break-words dark:bg-zinc-800/40"
        v-html="safeText"
      />
      <PostImageSlider v-if="post.imageUrls.length > 0" :images="post.imageUrls" />
    </div>
  </Card>
</template>

<style scoped>
.post-card {
  content-visibility: auto;
  contain-intrinsic-size: auto 220px;
  animation: post-in 0.3s ease-out both;
}

@keyframes post-in {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.tg-content :deep(a) {
  color: var(--k-link, #3b82f6);
  text-decoration: underline;
}

.tg-content :deep(code) {
  background: rgb(113 113 122 / 0.15);
  border-radius: 4px;
  padding: 1px 4px;
  font-size: 0.875em;
}

.tg-content :deep(pre) {
  background: rgb(113 113 122 / 0.15);
  border-radius: 8px;
  padding: 10px 12px;
  overflow-x: auto;
  white-space: pre;
}

.tg-content :deep(pre code) {
  background: none;
  padding: 0;
}

.tg-content :deep(blockquote) {
  margin: 0;
  padding-left: 10px;
  border-left: 3px solid rgb(113 113 122 / 0.4);
  color: var(--k-text-secondary, #71717a);
}

.tg-content :deep(.tg-spoiler) {
  background: currentColor;
  border-radius: 4px;
  transition: background 0.2s;
}

.tg-content :deep(.tg-spoiler:hover) {
  background: transparent;
}
</style>
