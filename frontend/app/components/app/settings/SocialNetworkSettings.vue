<script setup lang="ts">
import DiscordIcon from '@/assets/images/social/discord.png'
import GoogleIcon from '@/assets/images/social/google.png'
import { useSocialNetworkSettings } from '@/composables/api/settings/useSocialNetworkSettings'
import SocialNetworkButton from '@/components/app/settings/SocialNetworkButton.vue'

const { socialNetworkData, handleDiscord, handleGoogle } = useSocialNetworkSettings()
</script>

<template>
  <div class="space-y-8 cursor-default transition-colors duration-300">
    <div class="flex flex-col items-center">
      <h2 class="text-3xl font-bold mb-1 text-zinc-900 dark:text-white">Социальные сети</h2>
    </div>
    <div class="space-y-6 text-left">
      <div class="pb-6 border-b space-y-4 border-zinc-200 dark:border-zinc-800">
        <h3 class="text-base leading-xl font-semibold text-zinc-800 dark:text-zinc-200">
          Подключение
        </h3>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 md:gap-6">
        <SocialNetworkButton
          name="Discord"
          :iconSrc="DiscordIcon"
          :isConnected="socialNetworkData.is_discord_connected"
          bgClasses="dark:bg-[#5865f2]/20 dark:hover:bg-[#5865f2]/30 dark:border-zinc-800 dark:hover:border-[#5865f2]/50"
          @click="handleDiscord"
        />

        <SocialNetworkButton
          name="Google"
          :iconSrc="GoogleIcon"
          :isConnected="socialNetworkData.is_google_connected"
          bgClasses="dark:bg-zinc-900 dark:hover:bg-zinc-800 dark:border-zinc-800"
          @click="handleGoogle"
        >
          <template #background>
            <div
              class="absolute inset-0 opacity-0 dark:opacity-10 bg-linear-to-r from-red-500 via-green-500 to-blue-500 transition-opacity group-hover:opacity-20 pointer-events-none"
            ></div>
          </template>
        </SocialNetworkButton>
      </div>
    </div>
  </div>
</template>
