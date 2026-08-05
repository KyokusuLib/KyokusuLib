<script setup lang="ts">
import NotFound from '@/components/app/not-found/NotFound.vue'
const event = useRequestEvent()

if (event) {
  setResponseStatus(event, 404)
}
</script>

<template>
  <NotFound />
</template>
