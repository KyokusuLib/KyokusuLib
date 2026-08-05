<script setup lang="ts">
import { ref } from 'vue'
import { type NovelaAuthorDetails } from '@/types/backend/novela'
import { staticImage } from '@/utils/str'
import { KyokusuAppRole } from '@/types/enums/role-enum'
import { useRolePermissions } from '@/composables/api/role/useRolePermissions'
import AuthorSettings from './AuthorSettings.vue'

const props = defineProps<{
  author: NovelaAuthorDetails
}>()

const emit = defineEmits(['updated'])

const { hasPermission } = useRolePermissions()
const isOpenAuthorSettings = ref(false)

const updatedAuthor = (payload: NovelaAuthorDetails) => {
  Object.assign(props.author, payload)
  emit('updated', payload)
}
</script>

<template>
  <div
    class="bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-3xl overflow-hidden shadow-sm"
  >
    <div
      class="h-32 sm:h-48 w-full bg-linear-to-r from-zinc-800 to-zinc-700 relative overflow-hidden group"
    >
      <div class="absolute inset-0 bg-black/20 z-10 backdrop-blur-sm"></div>
      <img
        v-if="author.picture"
        :src="staticImage(author.picture)"
        alt="Banner Background"
        class="absolute inset-0 w-full h-full object-cover blur-xl scale-125 opacity-50"
      />
      <div class="absolute inset-0 bg-linear-to-t from-black/60 to-transparent z-10"></div>
    </div>

    <div class="px-4 sm:px-8 pb-8 relative z-20">
      <div
        class="flex flex-col sm:flex-row items-center sm:items-end gap-4 sm:gap-6 -mt-16 sm:-mt-20 mb-6"
      >
        <div
          class="w-32 h-32 sm:w-40 sm:h-40 rounded-3xl border-4 border-white dark:border-[#18181b] bg-zinc-100 dark:bg-zinc-800 overflow-hidden shadow-xl shrink-0 flex items-center justify-center"
        >
          <img
            v-if="author.picture"
            :src="staticImage(author.picture)"
            alt="Author picture"
            class="w-full h-full object-cover"
          />
          <Icon v-else name="ph:user-bold" size="64" class="text-zinc-300 dark:text-zinc-600" />
        </div>

        <div
          class="flex-1 w-full text-center sm:text-left pt-2 sm:pt-0 sm:pb-2 flex justify-between items-start"
        >
          <div>
            <h1
              class="text-3xl sm:text-4xl font-black text-zinc-900 dark:text-white tracking-tight mb-2"
            >
              {{ author.name }}
            </h1>

            <div class="flex flex-wrap items-center justify-center sm:justify-start gap-2">
              <span
                v-if="author.country"
                class="inline-flex items-center gap-1 px-3 py-1 rounded-full bg-zinc-100 dark:bg-zinc-800 text-xs font-bold text-zinc-600 dark:text-zinc-300"
              >
                <Icon name="ph:map-pin-bold" size="14" class="text-zinc-400" />
                {{ author.country }}
              </span>

              <span
                v-if="author.metier"
                class="inline-flex items-center gap-1 px-3 py-1 rounded-full bg-yellow-50 dark:bg-yellow-500/10 text-xs font-bold text-yellow-600 dark:text-yellow-500 border border-yellow-200/50 dark:border-yellow-500/20"
              >
                <Icon name="ph:briefcase-bold" size="14" />
                {{ author.metier }}
              </span>
            </div>
          </div>

          <button
            v-if="hasPermission(KyokusuAppRole.MODERATOR)"
            @click="isOpenAuthorSettings = true"
            class="mt-4 cursor-pointer"
          >
            <Icon
              name="ph:gear"
              size="24"
              class="text-zinc-400 hover:text-yellow-500 transition-colors duration-300"
            />
          </button>
        </div>
        <AuthorSettings v-model="isOpenAuthorSettings" :author="author" @updated="updatedAuthor" />
      </div>

      <div class="mt-8">
        <h3 class="text-lg font-bold text-zinc-900 dark:text-white mb-4 flex items-center gap-2">
          <Icon name="ph:info-bold" class="text-yellow-500" />
          Об авторе
        </h3>

        <div
          class="bg-zinc-50 dark:bg-zinc-900/50 rounded-2xl p-5 sm:p-6 border border-zinc-100 dark:border-zinc-800"
        >
          <div
            v-if="author.bio"
            class="prose dark:prose-invert max-w-none text-zinc-700 dark:text-zinc-300 text-sm sm:text-base leading-relaxed whitespace-pre-line"
            v-html="author.bio"
          ></div>
          <p v-else class="text-sm font-medium text-zinc-500 italic text-center py-4">
            Информация об авторе пока не заполнена.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
