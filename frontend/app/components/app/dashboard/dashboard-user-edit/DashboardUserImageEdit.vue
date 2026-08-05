<script lang="ts" setup>
import { Label, Separator } from '@kyokusu-ui/vue'
import { useImageEdit } from '@/composables/ui/useImageEdit'

const avatar = defineModel<string>('avatar', { required: true })
const banner = defineModel<string>('banner', { required: true })

const {
  avatarInput,
  bannerInput,
  trigger,
  onAvatarChange,
  onBannerChange,
  removeAvatar,
  removeBanner,
} = useImageEdit(avatar, banner)
</script>

<template>
  <div class="flex flex-col gap-4 space-y-4">
    <div class="flex items-center gap-4">
      <div
        class="relative w-24 h-24 rounded-2xl overflow-hidden bg-zinc-200 dark:bg-zinc-800 cursor-pointer group shrink-0 border border-zinc-300 dark:border-zinc-700"
        @click="trigger(avatarInput)"
      >
        <img v-if="avatar" :src="avatar" alt="Avatar" class="w-full h-full object-cover" />
        <div
          class="absolute inset-0 flex items-center justify-center bg-black/0 group-hover:bg-black/40 transition-colors rounded-2xl"
        >
          <Icon
            name="ph:camera"
            size="16"
            class="text-white opacity-0 group-hover:opacity-100 transition-opacity"
          />
        </div>
        <button
          v-if="avatar"
          class="absolute top-1 right-1 w-6 h-6 rounded-full bg-black/50 hover:bg-red-500 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all cursor-pointer"
          title="Удалить аватар"
          @click.stop="removeAvatar"
        >
          <Icon name="ph:trash-bold" size="12" class="text-white" />
        </button>
        <input
          ref="avatarInput"
          type="file"
          accept="image/*"
          class="hidden"
          @change="onAvatarChange"
        />
      </div>
      <div class="flex flex-col">
        <Label label="Аватар" />
        <div class="flex flex-col ml-1">
          <p class="text-xs text-zinc-400">Нажмите на картинку для замены</p>
          <span class="text-xs text-zinc-500 dark:text-stone-500 mt-1">
            Рекомендуется: изображение, мин. 400x400px.
          </span>
        </div>
      </div>
      <Separator orientation="vertical" />
      <div class="flex flex-col">
        <Label label="Баннер" />
        <div class="flex flex-col ml-1">
          <p class="text-xs text-zinc-400">Нажмите на картинку для замены</p>
          <span class="text-xs text-zinc-500 dark:text-stone-500 mt-1">
            Рекомендуется: изображение, мин. 1200x320px.
          </span>
        </div>
      </div>
    </div>

    <div
      class="relative w-full h-32 rounded-xl overflow-hidden bg-zinc-200 dark:bg-zinc-800 cursor-pointer group border border-zinc-300 dark:border-zinc-700"
      @click="trigger(bannerInput)"
    >
      <img v-if="banner" :src="banner" alt="Banner" class="w-full h-full" />
      <div
        class="absolute inset-0 flex items-center justify-center bg-black/0 group-hover:bg-black/40 transition-colors"
      >
        <Icon
          name="ph:camera-bold"
          size="24"
          class="text-white opacity-0 group-hover:opacity-100 transition-opacity"
        />
      </div>
      <button
        v-if="banner"
        class="absolute top-2 right-2 w-7 h-7 rounded-full bg-black/50 hover:bg-red-500 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-all cursor-pointer"
        title="Удалить баннер"
        @click.stop="removeBanner"
      >
        <Icon name="ph:trash-bold" size="14" class="text-white" />
      </button>
      <input
        ref="bannerInput"
        type="file"
        accept="image/*"
        class="hidden"
        @change="onBannerChange"
      />
    </div>
  </div>
</template>
