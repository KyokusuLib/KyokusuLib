<script setup lang="ts">
import { PageLayout } from '@kyokusu-ui/vue';
import NewsSourcesPanel from '@/components/app/news/NewsSourcesPanel.vue';
import NewsFeed from '@/components/app/news/NewsFeed.vue';
</script>

<template>
    <PageLayout>
        <div class="grid grid-cols-2 md:grid-cols-8 gap-6">
            <div class="md:col-span-2">
                <NewsSourcesPanel />
            </div>
            <div class="md:col-span-6">
                <NewsFeed />
            </div>
        </div>
    </PageLayout>
</template>