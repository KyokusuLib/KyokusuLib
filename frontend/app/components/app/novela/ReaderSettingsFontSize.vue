<script setup lang="ts">
import { useReaderSettings } from '~/composables/ui/useReaderSettings'
import RangeSlider from '~/components/ui/RangeSlider.vue'

const { fontSize } = useReaderSettings()
</script>

<template>
  <RangeSlider
    v-model="fontSize"
    :min="12"
    :max="40"
    :step="1"
    label="Размер шрифта"
    suffix="px"
    icon="ph:text-aa-bold"
  />
</template>
