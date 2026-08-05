<script setup lang="ts">
const props = defineProps<{
  userID: number
}>()
</script>

<template>
  <div>
    <p
      class="flex justify-center items-center mt-2 px-4 rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900/50 text-zinc-900 dark:text-zinc-200 shadow-sm transition-all duration-300 hover:border-yellow-500/30 hover:shadow-lg hover:shadow-yellow-500/5 cursor-default h-8 sm:w-auto font-bold"
    >
      ID: #{{ userID }}
    </p>
  </div>
</template>
