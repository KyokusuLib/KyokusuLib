<script setup lang="ts">
import ModalConfirm from '~/components/common/ModalConfirm.vue'
import TgPostCard from '@/components/app/news/TgPostCard.vue'
import { useTgPosts } from '@/composables/api/news/useTgPosts'
import { useRolePermissions } from '@/composables/api/role/useRolePermissions'
import { NEWS_SOURCES, TG_CHANNEL_NAME } from '@/constants/news'
import { KyokusuAppRole } from '@/types/enums/role-enum'
import { SocialNetwork } from '@/types/enums/social-network-enum'
import type { BackendTgPost } from '@/types/backend/tg-post'

const { hasPermission } = useRolePermissions()
const { posts, fetchTgPosts, connectStream, disconnectStream, deleteTgPost } = useTgPosts()

const source = computed(() =>
  NEWS_SOURCES.find((s) => s.id === SocialNetwork.Telegram),
)

const postUrl = (post: BackendTgPost) =>
  `https://t.me/${TG_CHANNEL_NAME}/${post.messageId}`

const canDelete = computed(() => hasPermission(KyokusuAppRole.MODERATOR))

const postToDelete = ref<BackendTgPost | null>(null)
const deleteError = ref('')

const showDeleteModal = computed({
  get: () => postToDelete.value !== null,
  set: (visible: boolean) => {
    if (!visible) postToDelete.value = null
  },
})

async function confirmDelete() {
  if (!postToDelete.value) return
  deleteError.value = ''
  try {
    await deleteTgPost(postToDelete.value.id)
    postToDelete.value = null
  } catch (err: any) {
    deleteError.value = err?.message ?? 'Не удалось удалить пост'
  }
}

onMounted(() => {
  fetchTgPosts()
  connectStream()
})

onUnmounted(() => {
  disconnectStream()
})
</script>

<template>
  <div class="flex flex-col gap-3">
    <p
      v-if="posts.length === 0"
      class="text-sm text-zinc-500 dark:text-zinc-400 text-center py-6"
    >
      Посты появятся, как только будут опубликованы в канале
    </p>

    <p v-if="deleteError" class="text-sm text-red-500 dark:text-red-400 text-center">
      {{ deleteError }}
    </p>

    <TgPostCard
      v-for="post in posts"
      :key="post.id"
      :post="post"
      :source="source!"
      :post-url="postUrl(post)"
      :can-delete="canDelete"
      @delete="postToDelete = $event"
    />

    <ModalConfirm
      v-model="showDeleteModal"
      title="Удаление поста"
      :description="`Пост #${postToDelete?.messageId ?? ''} будет удалён безвозвратно`"
      confirm-text="Удалить"
      cancel-text="Отмена"
      @confirm="confirmDelete"
    />
  </div>
</template>
