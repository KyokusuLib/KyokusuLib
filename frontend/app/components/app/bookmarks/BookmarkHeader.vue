<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    totalBookmarkCounts: number
  }>(),
  {
    totalBookmarkCounts: 0,
  },
)
</script>

<template>
  <div class="flex flex-col">
    <h1 class="text-3xl font-black text-zinc-900 dark:text-white tracking-tight">
      Мои Закладки ({{ totalBookmarkCounts }})
    </h1>

    <svg
      class="mt-2 mb-4 drop-shadow-[0_4px_10px_rgba(234,179,8,0.5)] overflow-visible"
      width="180"
      height="12"
      viewBox="0 0 160 12"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        d="M2 2C50 14 110 14 158 2"
        stroke="url(#arc-gradient)"
        stroke-width="4"
        stroke-linecap="round"
      />
      <defs>
        <linearGradient
          id="arc-gradient"
          x1="0"
          y1="0"
          x2="160"
          y2="0"
          gradientUnits="userSpaceOnUse"
        >
          <stop offset="0%" stop-color="#ca8a04" />
          <stop offset="100%" stop-color="#facc15" />
        </linearGradient>
      </defs>
    </svg>
  </div>
</template>
