<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import RangeSlider from '@/components/ui/RangeSlider.vue'
import CircleBlock from '@/components/ui/CircleBlock.vue'
import { Label } from '@kyokusu-ui/vue'
import { useUserExperiance } from '@/composables/api/user/useUserExperiance'
import { getLevelColor } from '@/utils/dashboard'
import type { UserDefinitions } from '@/types/backend/user'

const props = defineProps<{
  level: number
  experience: number
}>()

const emit = defineEmits<{
  'update:level': [value: number]
  'update:experience': [value: number]
}>()

const { getDefinitions } = useUserExperiance()

const { data: definitions, status } = useAsyncData<UserDefinitions[]>(
  'user-level-definitions',
  async () => {
    const raw = await getDefinitions()
    if (!raw) return []
    return Array.isArray(raw) ? raw : [raw]
  },
)

const selectedLevel = computed({
  get: () => props.level,
  set: (v) => emit('update:level', v),
})

const selectedExperience = computed({
  get: () => props.experience,
  set: (v) => emit('update:experience', v),
})

const originalLevel = ref(props.level)
const originalExperience = ref(props.experience)

const canReset = computed(
  () =>
    selectedLevel.value !== originalLevel.value ||
    selectedExperience.value !== originalExperience.value,
)

function resetLevel() {
  emit('update:level', originalLevel.value)
  emit('update:experience', originalExperience.value)
}

const defs = computed(() => definitions.value ?? [])

const currentDef = computed(() => defs.value.find((d) => d.level === selectedLevel.value))

const currentDefIndex = computed(() => defs.value.findIndex((d) => d.level === selectedLevel.value))

const nextDef = computed(() => {
  const i = currentDefIndex.value
  return i >= 0 && i < defs.value.length - 1 ? defs.value[i + 1] : null
})

const maxXp = computed(
  () => nextDef.value?.total_xp_required ?? currentDef.value?.total_xp_required ?? 0,
)

watch(selectedLevel, () => {
  if (selectedExperience.value > maxXp.value) {
    selectedExperience.value = maxXp.value
  }
})
</script>

<template>
  <div class="flex flex-col gap-5 w-full">
    <div class="flex flex-col gap-2">
      <div class="flex items-center justify-between">
        <Label label="Уровень" />
        <button
          v-if="canReset"
          class="text-xs font-semibold text-zinc-400 hover:text-yellow-500 transition-colors cursor-pointer flex items-center gap-1"
          @click="resetLevel"
        >
          <Icon name="ph:arrow-counter-clockwise-bold" size="14" />
          Сбросить
        </button>
      </div>
      <div v-if="status === 'pending'" class="text-sm text-zinc-400">Загрузка...</div>
      <div v-else class="relative">
        <select
          v-model.number="selectedLevel"
          class="w-full h-10 rounded-lg border pl-12 pr-10 text-sm outline-none appearance-none transition focus:ring-2"
          :style="{
            backgroundColor: 'var(--k-editor-bg)',
            borderColor: 'var(--k-editor-border)',
            color: 'var(--k-editor-text)',
          }"
        >
          <option v-for="def in defs" :key="def.level" :value="def.level">
            {{ def.title }}
          </option>
        </select>

        <div class="absolute left-2 top-1/2 -translate-y-1/2 pointer-events-none">
          <CircleBlock
            :text="selectedLevel"
            :colorClass="getLevelColor(selectedLevel)"
            border
            size="w-7 h-7"
          />
        </div>

        <Icon
          name="ph:caret-down-bold"
          size="14"
          class="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-400 pointer-events-none"
        />
      </div>
    </div>

    <RangeSlider
      v-if="maxXp > 0"
      v-model="selectedExperience"
      :min="0"
      :max="maxXp"
      :step="1"
      label="Текущий Опыт"
      suffix="XP"
      icon="ph:lightning-bold"
    />
    <div class="flex items-center justify-between text-xs text-zinc-500 dark:text-zinc-400 px-2">
      <span v-if="currentDef" class="font-semibold">
        Титул: <strong>{{ currentDef.title }}</strong>
      </span>
      <span v-if="nextDef" class="font-semibold">
        До уровня {{ nextDef.level }}:
        <strong>{{ Math.max(maxXp - selectedExperience, 0) }} XP</strong>
      </span>
      <span v-else-if="currentDef" class="font-semibold"> Максимальный уровень </span>
    </div>
  </div>
</template>
