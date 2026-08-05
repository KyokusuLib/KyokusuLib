<script setup lang="ts">
import FooterApp from '@/components/common/FooterApp.vue'
</script>

<template>
  <div class="layout-auth min-h-screen flex flex-col">
    <main class="grow">
      <slot />
    </main>

    <FooterApp />
  </div>
</template>
