<script setup lang="ts">
import {
  STATUS_OPTIONS,
  TRANSLATION_STATUS_OPTIONS,
  NOVELA_TYPE_OPTIONS,
  AGE_RATING_OPTIONS,
  SORT_OPTIONS,
} from '~/types/frontend/query/novela-filters'
import { Select, Label } from '@kyokusu-ui/vue'

const status = defineModel<string>('status', { default: '' })
const translationStatus = defineModel<string>('translationStatus', { default: '' })
const type = defineModel<string>('type', { default: '' })
const ageRating = defineModel<string>('ageRating', { default: '' })
const sort = defineModel<string>('sort', { default: '' })
</script>

<template>
  <div class="space-y-3">
    <div class="space-y-1.5">
      <Label label="Тип" />
      <Select v-model="type" :options="NOVELA_TYPE_OPTIONS" placeholder="Все типы" />
    </div>

    <div class="space-y-1.5">
      <Label label="Статус" />
      <Select v-model="status" :options="STATUS_OPTIONS" placeholder="Все статусы" />
    </div>

    <div class="space-y-1.5">
      <Label label="Статус перевода" />
      <Select
        v-model="translationStatus"
        :options="TRANSLATION_STATUS_OPTIONS"
        placeholder="Все статусы"
      />
    </div>

    <div class="space-y-1.5">
      <Label label="Рейтинг" />
      <Select v-model="ageRating" :options="AGE_RATING_OPTIONS" placeholder="Возврастной рейтинг" />
    </div>

    <div class="space-y-1.5">
      <Label label="Сортировать по" />
      <Select v-model="sort" :options="SORT_OPTIONS" placeholder="Сортировать..." />
    </div>
  </div>
</template>
