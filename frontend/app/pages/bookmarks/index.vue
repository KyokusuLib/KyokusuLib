<script setup lang="ts">
import Bookmarks from '~/components/app/bookmarks/Bookmarks.vue'

definePageMeta({ middleware: ['auth'] })
</script>

<template>
  <Bookmarks />
</template>
