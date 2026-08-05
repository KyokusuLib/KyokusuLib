<script setup lang="ts">
import { useProfile } from '~/composables/api/profile/useProfile'
import ProfileSelfContent from '~/components/app/profile/ProfileSelfContent.vue'
import ProfilePublic from '~/components/app/profile/ProfilePublic.vue'
import ProfilePrivate from '~/components/app/profile/ProfilePrivate.vue'

const route = useRoute()
const { profileData, isSelfProfile, isPublicAccount, init, profileError } = useProfile()

definePageMeta({
  middleware: (to) => {
    if (!/^\d+$/.test(to.params.id as string)) {
      return navigateTo('/')
    }
  },
})

const userId = computed(() => Number(route.params.id))

const { status } = await useAsyncData(`profile-${userId.value}`, () => init(userId.value), {
  watch: [userId],
})

if (profileError.value) {
  throw createError({ statusCode: 404 })
}

useSeoMeta({
  title: () => (profileData.value?.name ? `${profileData.value.name} - Профиль` : 'Профиль'),
  ogTitle: () => profileData.value?.name,
  ogImage: () => profileData.value?.picture,
})
</script>

<template>
  <div class="min-h-screen bg-white dark:bg-[#0f0f0f]">
    <div v-if="status === 'pending'" class="flex justify-center py-20 items-center min-h-screen">
      <svg
        class="animate-spin h-10 w-10 text-zinc-500"
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
      >
        <circle
          class="opacity-25"
          cx="12"
          cy="12"
          r="10"
          stroke="currentColor"
          stroke-width="4"
        ></circle>
        <path
          class="opacity-75"
          fill="currentColor"
          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
        ></path>
      </svg>
    </div>

    <div v-else-if="profileData">
      <ProfileSelfContent v-if="isSelfProfile" />
      <ProfilePublic v-else-if="isPublicAccount" :profileData="profileData" />
      <ProfilePrivate v-else :profileData="profileData" />
    </div>
  </div>
</template>
