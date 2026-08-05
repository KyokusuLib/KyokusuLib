<script setup lang="ts">
import { useReaderSettings } from '~/composables/ui/useReaderSettings'
import RangeSlider from '~/components/ui/RangeSlider.vue'

const { lineWeight } = useReaderSettings()

function formatWeight(v: number) {
  return v.toFixed(1)
}
</script>

<template>
  <RangeSlider
    v-model="lineWeight"
    :min="1.0"
    :max="2.5"
    :step="0.1"
    label="Межстрочный интервал"
    icon="ph:text-align-center-bold"
    :format="formatWeight"
  />
</template>
