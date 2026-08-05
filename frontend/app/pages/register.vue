<script setup lang="ts">
import RegisterPage from '~/components/app/auth/RegisterPage.vue'

definePageMeta({
  layout: 'auth',
  middleware: ['guest'],
})
</script>

<template>
  <RegisterPage />
</template>
