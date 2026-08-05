<script setup lang="ts">
import LoginPage from '~/components/app/auth/LoginPage.vue'

definePageMeta({
  layout: 'auth',
  middleware: ['guest'],
})
</script>

<template>
  <LoginPage />
</template>
