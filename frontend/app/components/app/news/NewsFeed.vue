<script setup lang="ts">
import { Card } from '@kyokusu-ui/vue'
import type { Component } from 'vue'
import SelectedSourcesBadges from '@/components/app/news/SelectedSourcesBadges.vue'
import TelegramFeed from '@/components/app/news/TelegramFeed.vue'
import { useNewsSources } from '@/composables/news/useNewsSources'
import { SocialNetwork } from '@/types/enums/social-network-enum'
import type { NewsSource } from '@/types/frontend/news'

const { selectedSourceDetails, hasSelectedSources } = useNewsSources()

const feedComponents: Partial<Record<SocialNetwork, Component>> = {
  [SocialNetwork.Telegram]: TelegramFeed,
}

interface ActiveFeed {
  source: NewsSource
  component: Component
}

const activeFeeds = computed<ActiveFeed[]>(() =>
  selectedSourceDetails.value.flatMap((source) => {
    const component = feedComponents[source.id]
    return component ? [{ source, component }] : []
  }),
)
</script>

<template>
  <Card>
    <template #header>
      <div class="flex items-center justify-center gap-4">
        <h2 class="text-xl font-bold flex justify-center items-center">Новости</h2>
        <Icon name="ph:newspaper-bold" size="24" />
      </div>
    </template>

    <div class="flex flex-col gap-4 min-h-30">
      <SelectedSourcesBadges v-if="hasSelectedSources" :sources="selectedSourceDetails" />

      <component
        v-for="feed in activeFeeds"
        :key="feed.source.id"
        :is="feed.component"
      />

      <p
        v-if="!hasSelectedSources"
        class="text-sm text-zinc-500 dark:text-zinc-400 text-center py-10"
      >
        Выберите источники, чтобы видеть новости
      </p>
    </div>
  </Card>
</template>
