<script setup lang="ts">
import { useAuthStore } from '~/stores/auth'
const authStore = useAuthStore()

if (authStore.user?.id) {
  await navigateTo(`/profile/${authStore.user.id}`, { replace: true })
} else {
  await navigateTo('/login')
}
</script>

<template>
  <div>Загрузка профиля...</div>
</template>
