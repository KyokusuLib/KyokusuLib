<script setup lang="ts">
interface Props {
  maxWidth?: string
}

withDefaults(defineProps<Props>(), {
  maxWidth: '72rem',
})
</script>

<template>
  <div class="k-page-layout">
    <div class="k-page-layout-inner" :style="{ maxWidth }">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.k-page-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #fafafa;
  transition: background-color 0.3s ease;
}

.dark .k-page-layout {
  background: radial-gradient(at center, #27272a 0%, #09090b 90%);
}

.k-page-layout-inner {
  width: 100%;
  margin-left: auto;
  margin-right: auto;
  padding: 2rem 1rem;
  box-sizing: border-box;
}

@media (min-width: 640px) {
  .k-page-layout-inner {
    padding-left: 1.5rem;
    padding-right: 1.5rem;
  }
}
</style>
