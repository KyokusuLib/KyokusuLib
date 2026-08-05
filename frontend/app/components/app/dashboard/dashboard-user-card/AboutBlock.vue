<script setup lang="ts">
defineProps<{ about: string }>()
</script>

<template>
  <div class="bg-zinc-50 dark:bg-zinc-800/50 rounded-xl px-4 py-3 border border-zinc-700">
    <div class="flex gap-2">
      <p class="text-[10px] font-bold uppercase tracking-widest text-zinc-400 mb-1.5">О себе</p>
      <Icon name="ph:book-bold" size="16" class="text-zinc-400 shrink-0" />
    </div>
    <p
      class="text-sm font-semibold text-zinc-600 dark:text-zinc-400 leading-relaxed text-wrap wrap-break-word"
    >
      {{ about ? about : 'Пользователь решил не рассказаывать о себе.' }}
    </p>
  </div>
</template>
