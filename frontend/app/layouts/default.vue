<script setup lang="ts">
import HeaderApp from '@/components/common/HeaderApp.vue'
import FooterApp from '@/components/common/FooterApp.vue'
</script>

<template>
  <div class="layout-default min-h-screen flex flex-col">
    <HeaderApp />

    <main class="grow">
      <slot />
    </main>

    <FooterApp />
  </div>
</template>
