<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { onClickOutside } from '@vueuse/core'
import InnerBlock from '~/components/ui/InnerBlock.vue'
import { Tooltip } from '@kyokusu-ui/vue'

export interface SortOption {
  field: string
  label: string
  icon: string
  compare: (a: unknown, b: unknown) => number
}

export interface FilterConfig {
  sortOptions: SortOption[]
  defaultField?: string
  defaultOrder?: 'asc' | 'desc'
}

const props = defineProps<{
  items: unknown[]
  modelValue: unknown[]
  config: FilterConfig
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: unknown[]): void
}>()

const sortOptions = computed(() => props.config.sortOptions)

const isSortOpen = ref(false)
const sortContainerRef = ref<HTMLElement | null>(null)
const sortField = ref<string>(props.config.defaultField ?? props.config.sortOptions[0]?.field ?? '')
const sortIcon = ref<string>(props.config.sortOptions[0]?.icon ?? '')
const sortOrder = ref<'asc' | 'desc'>(props.config.defaultOrder ?? 'asc')

onClickOutside(sortContainerRef, () => (isSortOpen.value = false))

const currentOption = computed(() => sortOptions.value.find((o) => o.field === sortField.value))

const toggleOrder = () => {
  sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
}

const selectField = (field: string, icon: string) => {
  if (sortField.value === field) {
    toggleOrder()
  } else {
    sortField.value = field
    sortOrder.value = 'asc'
    sortIcon.value = icon
  }
  isSortOpen.value = false
}

const applySort = () => {
  const option = sortOptions.value.find((o) => o.field === sortField.value)
  if (!option) {
    emit('update:modelValue', [...props.items])
    return
  }
  const sorted = [...props.items].sort((a, b) => {
    const comparison = option.compare(a, b)
    return sortOrder.value === 'asc' ? comparison : -comparison
  })
  emit('update:modelValue', sorted)
}

watch([() => props.items, sortField, sortOrder], applySort, { immediate: true })
</script>

<template>
  <div>
    <h3
      class="text-md shadow-2xs pl-1 mb-2 font-semibold text-zinc-900 dark:text-white cursor-default"
    >
      Фильтры
    </h3>
    <InnerBlock class="flex justify-start h-12 gap-2">
      <div class="relative inline-block" ref="sortContainerRef">
        <Tooltip text="Сортировать">
          <div
            @click="isSortOpen = !isSortOpen"
            class="flex items-center gap-2 p-2 rounded-md cursor-pointer border bg-white dark:bg-zinc-900/50 dark:border-zinc-800 hover:border-yellow-500/50 transition-colors"
          >
            <div
              class="flex items-center h-4 gap-2 px-2 text-xs font-semibold text-zinc-600 dark:text-zinc-400"
            >
              <Icon
                :name="sortIcon"
                size="18"
                class="text-zinc-400 transition-transform duration-300"
              />
              {{ currentOption?.label }}
            </div>
            <Icon
              name="ph:caret-down-bold"
              size="10"
              class="text-zinc-400 transition-transform duration-300"
              :class="{ 'rotate-180': isSortOpen }"
            />
          </div>
        </Tooltip>

        <Transition name="dropdown">
          <div
            v-if="isSortOpen"
            class="absolute left-0 mt-2 w-56 bg-white/90 dark:bg-zinc-900/90 backdrop-blur-xl border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-2xl z-50 overflow-hidden shadow-black/20"
          >
            <div class="px-4 py-3 border-b border-zinc-100 dark:border-zinc-800/50">
              <p
                class="text-[10px] font-black uppercase tracking-[0.2em] text-zinc-400 text-center"
              >
                Сортировка
              </p>
            </div>

            <div class="p-1.5 space-y-1">
              <button
                v-for="option in sortOptions"
                :key="option.field"
                @click="selectField(option.field, option.icon)"
                class="w-full flex items-center justify-between px-3 py-2.5 rounded-xl text-sm transition-all duration-200 cursor-pointer"
                :class="[
                  option.field === sortField
                    ? 'bg-yellow-500/10 text-yellow-600 dark:text-yellow-500 font-bold'
                    : 'text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 hover:text-zinc-900 dark:hover:text-zinc-100',
                ]"
              >
                <div class="flex items-center gap-2.5">
                  <Icon :name="option.icon" size="16" />
                  <span>{{ option.label }}</span>
                </div>

                <div v-if="option.field === sortField" class="flex items-center gap-1">
                  <Icon
                    :name="sortOrder === 'asc' ? 'ph:caret-up-bold' : 'ph:caret-down-bold'"
                    size="14"
                  />
                </div>
              </button>
            </div>

            <div
              class="p-2 bg-zinc-50 dark:bg-zinc-800/40 border-t border-zinc-100 dark:border-zinc-800/50"
            >
              <div
                class="flex items-center gap-2 mr-2 justify-center text-[12px] text-center text-zinc-400"
              >
                <Icon
                  :name="sortOrder === 'asc' ? 'ph:sort-ascending' : 'ph:sort-descending'"
                  size="14"
                />
                {{ sortOrder === 'asc' ? 'По возрастанию' : 'По убыванию' }}
              </div>
            </div>
          </div>
        </Transition>
      </div>

      <Tooltip :text="sortOrder === 'asc' ? 'По убыванию' : 'По возрастанию'">
        <div
          @click="toggleOrder"
          class="flex items-center h-8 p-2 rounded-md cursor-pointer border bg-white dark:bg-zinc-900/50 dark:border-zinc-800 hover:border-yellow-500/50 transition-colors"
        >
          <Icon
            :name="sortOrder === 'asc' ? 'ph:sort-ascending' : 'ph:sort-descending'"
            size="16"
            class="text-zinc-500 dark:text-zinc-400"
          />
        </div>
      </Tooltip>
    </InnerBlock>
  </div>
</template>

<style scoped>
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.3s cubic-bezier(0.23, 1, 0.32, 1);
}
.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.98);
}
</style>
