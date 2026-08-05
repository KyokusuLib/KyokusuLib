<script setup lang="ts">
import { ref, computed } from 'vue'

export interface UiTableColumn {
  key: string
  label: string
  sortable?: boolean
  align?: 'left' | 'center' | 'right'
}

const props = withDefaults(
  defineProps<{
    columns: UiTableColumn[]
    items: Record<string, any>[]
    loading?: boolean
    emptyText?: string
  }>(),
  { loading: false, emptyText: 'Нет данных' },
)

const search = defineModel<string>('search', { default: '' })

const sortKey = ref<string | null>(null)
const sortDir = ref<'asc' | 'desc'>('asc')

function toggleSort(key: string) {
  if (!key) return
  if (sortKey.value === key) {
    sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = key
    sortDir.value = 'asc'
  }
}

const filteredItems = computed(() => {
  let result = props.items
  if (search.value) {
    const q = search.value.toLowerCase()
    result = result.filter((item) =>
      props.columns.some((col) => {
        const val = item[col.key]
        return val !== null && String(val).toLowerCase().includes(q)
      }),
    )
  }
  if (sortKey.value) {
    const dir = sortDir.value === 'asc' ? 1 : -1
    result = [...result].sort((a, b) => {
      const aVal = a[sortKey.value!]
      const bVal = b[sortKey.value!]
      if (typeof aVal === 'string') return dir * aVal.localeCompare(bVal)
      return dir * ((aVal ?? 0) - (bVal ?? 0))
    })
  }
  return result
})

const slots = defineSlots<{
  [key: `cell-${string}`]: (props: { item: Record<string, any> }) => any
  actions: (props: { item: Record<string, any> }) => any
}>()
</script>

<template>
  <div>
    <!-- loading -->
    <div v-if="loading" class="py-24 flex flex-col items-center justify-center text-zinc-400">
      <Icon name="ph:spinner-bold" size="32" class="animate-spin mb-4" />
      <p class="text-xs font-bold uppercase tracking-[0.2em]">Загрузка...</p>
    </div>

    <!-- empty -->
    <div
      v-else-if="filteredItems.length === 0"
      class="py-24 flex flex-col items-center justify-center border-2 border-dashed border-zinc-200 dark:border-zinc-800 rounded-[3rem] text-zinc-400 bg-zinc-50/50 dark:bg-zinc-900/10"
    >
      <Icon name="ph:users-bold" size="64" class="opacity-20 mb-4" />
      <p class="font-bold uppercase tracking-[0.2em] text-[10px]">
        {{ emptyText }}
      </p>
    </div>

    <!-- table -->
    <div v-else class="overflow-x-auto rounded-2xl border border-zinc-200 dark:border-zinc-800">
      <table class="w-full text-left border-collapse">
        <thead>
          <tr
            class="border-b border-zinc-200 dark:border-zinc-800 text-sm font-semibold text-zinc-500 dark:text-zinc-400 bg-zinc-50 dark:bg-zinc-900/50"
          >
            <th
              v-for="col in columns"
              :key="col.key"
              :class="[
                'py-4 px-4 whitespace-nowrap',
                col.align === 'center'
                  ? 'text-center'
                  : col.align === 'right'
                    ? 'text-right'
                    : 'text-left',
                col.sortable
                  ? 'cursor-pointer select-none hover:text-zinc-700 dark:hover:text-zinc-300 transition-colors'
                  : '',
              ]"
              @click="col.sortable && toggleSort(col.key)"
            >
              <span class="inline-flex items-center gap-1.5">
                {{ col.label }}
                <Icon
                  v-if="col.sortable && sortKey === col.key"
                  :name="sortDir === 'asc' ? 'ph:caret-up-bold' : 'ph:caret-down-bold'"
                  size="12"
                />
              </span>
            </th>
            <th v-if="$slots.actions" class="py-4 px-4 w-24 text-center whitespace-nowrap">
              Действия
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(item, idx) in filteredItems"
            :key="item.id ?? idx"
            class="border-b border-zinc-100 dark:border-zinc-800/50 hover:bg-zinc-50 dark:hover:bg-zinc-800/30 transition-colors"
          >
            <td
              v-for="col in columns"
              :key="col.key"
              :class="[
                'py-3.5 px-4 text-sm text-zinc-700 dark:text-zinc-300',
                col.align === 'center' ? 'text-center' : col.align === 'right' ? 'text-right' : '',
              ]"
            >
              <slot :name="`cell-${col.key}`" :item="item">
                {{ item[col.key] ?? '—' }}
              </slot>
            </td>
            <td v-if="$slots.actions" class="py-3.5 px-4 text-center">
              <slot name="actions" :item="item" />
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
