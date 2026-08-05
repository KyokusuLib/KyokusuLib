<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { MultiSelect } from '@kyokusu-ui/vue'
import type { MultiSelectOption } from '@kyokusu-ui/vue'
import { useUserTags } from '@/composables/api/user/useUserTags'
import type { DashboardUserTag } from '@/types/frontend/dashboard/users'
import type { UserTagDefinitions } from '@/types/backend/user'

const props = defineProps<{
  tags: DashboardUserTag[]
}>()

const emit = defineEmits<{
  'update:tags': [value: DashboardUserTag[]]
}>()

const { getDefinitions } = useUserTags()

const definitions = ref<UserTagDefinitions[]>([])
const loading = ref(true)

onMounted(async () => {
  loading.value = true
  const raw = await getDefinitions()
  definitions.value = Array.isArray(raw) ? raw : raw ? [raw] : []
  loading.value = false
})

const multiOptions = computed<MultiSelectOption[]>(() =>
  definitions.value.map((d) => ({ id: d.tag_id, label: d.tag })),
)

const selectedIds = computed(() => props.tags.map((t) => t.tag_id))

function onUpdate(value: (string | number | MultiSelectOption)[]) {
  const ids = new Set(value.map((v) => (typeof v === 'object' ? v.id : Number(v))))

  const newTags: DashboardUserTag[] = definitions.value
    .filter((d) => ids.has(d.tag_id))
    .map((d) => ({ tag_id: d.tag_id, tag: d.tag }))

  emit('update:tags', newTags)
}
</script>

<template>
  <MultiSelect
    v-if="!loading"
    id="user-tags"
    :model-value="selectedIds"
    :options="multiOptions"
    placeholder="Выберите теги..."
    @update:model-value="onUpdate"
  />
</template>
