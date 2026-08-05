<script setup lang="ts">
import type { NovelaChapter } from '~/types/backend/novela'
import { useReadProgress } from '~/composables/api/novela/useReadProgress'
import { useRolePermissions } from '~/composables/api/role/useRolePermissions'
import { KyokusuAppRole } from '~/types/enums/role-enum'
import { useNovela } from '~/composables/api/novela/useNovela'
import ModalConfirm from '~/components/common/ModalConfirm.vue'
import { TeleportedTooltip } from '@kyokusu-ui/vue'

const props = defineProps<{
  chapter: NovelaChapter
  novelaId: number
}>()

const { isAuthenticated } = useAuthStore()
const { isChapterRead, markChapterAsRead } = useReadProgress()
const { hasPermission } = useRolePermissions()
const { deleteChapter } = useNovela()

const read = computed(() => isChapterRead(props.chapter))
const canManage = computed(() => hasPermission(KyokusuAppRole.MODERATOR))

const isDeleteConfirmOpen = ref(false)

async function onDelete() {
  await deleteChapter(props.novelaId, props.chapter.id)
  isDeleteConfirmOpen.value = false
}

async function markAsRead() {
  const isRead = await markChapterAsRead(props.chapter.id)
  props.chapter.is_read = isRead
}
</script>

<template>
  <div
    :class="[
      'flex justify-between items-center p-4 border-b border-zinc-100 dark:border-zinc-800/50 last:border-0 transition-colors group',
      read
        ? 'opacity-50 hover:opacity-80 bg-zinc-50 dark:bg-zinc-900/30'
        : 'hover:bg-zinc-50 dark:hover:bg-zinc-800/50',
    ]"
  >
    <NuxtLink
      :to="`/novela/reader/${novelaId}/${chapter.id}`"
      class="flex grow justify-start items-center gap-3 cursor-pointer"
    >
      <span class="mt-1 text-zinc-400 font-mono text-sm w-8">#{{ chapter.number }}</span>
      <span
        :class="[
          'font-medium transition-colors',
          read
            ? 'text-zinc-500 dark:text-zinc-500'
            : 'text-zinc-900 dark:text-zinc-200 group-hover:text-yellow-600 dark:group-hover:text-yellow-600',
        ]"
      >
        {{ chapter.title }}
      </span>
    </NuxtLink>
    <div v-if="isAuthenticated" class="flex items-center mr-6">
      <TeleportedTooltip text="Отметить как прочитанную">
        <button
          @click="markAsRead"
          class="flex items-center rounded-lg transition-colors cursor-pointer"
          :class="[read ? 'text-zinc-500 dark:text-zinc-500' : 'text-zinc-900 dark:text-zinc-200']"
        >
          <Icon name="ph:eye-bold" size="18" />
        </button>
      </TeleportedTooltip>
    </div>

    <div
      v-if="canManage"
      class="flex items-center gap-6 md:mr-6 opacity-0 group-hover:opacity-100 transition-opacity"
    >
      <TeleportedTooltip text="Редактировать главу">
        <NuxtLink
          :to="`/novela/${novelaId}/edit-chapter/${chapter.id}`"
          class="flex items-center rounded-lg hover:text-zinc-200 dark:hover:text-zinc-700 text-zinc-500 dark:text-zinc-400 transition-colors cursor-pointer"
        >
          <Icon name="ph:pencil-simple-bold" size="18" />
        </NuxtLink>
      </TeleportedTooltip>
      <TeleportedTooltip text="Удалить главу">
        <button
          @click="isDeleteConfirmOpen = true"
          class="flex items-center rounded-lg text-zinc-500 dark:text-zinc-400 hover:text-red-600 dark:hover:text-red-400 transition-colors cursor-pointer"
        >
          <Icon name="ph:trash-bold" size="18" />
        </button>
      </TeleportedTooltip>
    </div>

    <ModalConfirm
      v-model="isDeleteConfirmOpen"
      title="Удаление главы"
      :description="`Вы уверены, что хотите удалить главу №${chapter.number} «${chapter.title}»? Это действие необратимо.`"
      confirm-text="Удалить"
      cancel-text="Отмена"
      @confirm="onDelete"
    />
  </div>
</template>
