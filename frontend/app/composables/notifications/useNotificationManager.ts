import { ref, computed } from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useNotifications } from '~/composables/api/notifications/useNotifications'
import { groupByTitle, sortGroupEntries, toggleSetItem } from '~/utils/notification'
import type { SortField } from '~/utils/notification'

export function useNotificationManager() {
  const authStore = useAuthStore()
  const { isAuthenticated } = storeToRefs(authStore)

  const { list, unreadCount, fetchList, markRead, markAllRead, remove } = useNotifications()

  // ── UI state ──
  const isLoading = ref(false)
  const selectedIds = ref<Set<number>>(new Set())
  const expandedGroups = ref<Set<string>>(new Set())
  const sortField = ref<SortField>('created_at')
  const sortOrder = ref<'asc' | 'desc'>('desc')

  // ── Fetch on auth ──
  async function load() {
    if (!isAuthenticated.value) return
    isLoading.value = true
    await fetchList(100, 0)
    isLoading.value = false
  }

  // ── Computed: groups ──
  const grouped = computed(() => groupByTitle(list.value))

  const groupEntries = computed(() =>
    sortGroupEntries(grouped.value, sortField.value, sortOrder.value),
  )

  // ── Computed: selection ──
  const allSelected = computed(
    () => list.value.length > 0 && list.value.every((n) => selectedIds.value.has(n.id)),
  )

  const selectedCount = computed(() => selectedIds.value.size)

  function toggleSelectAll() {
    selectedIds.value = allSelected.value ? new Set() : new Set(list.value.map((n) => n.id))
  }

  function toggleSelect(id: number) {
    selectedIds.value = toggleSetItem(selectedIds.value, id)
  }

  // ── Expand / collapse ──
  function toggleExpand(title: string) {
    expandedGroups.value = toggleSetItem(expandedGroups.value, title)
  }

  function expandAll() {
    expandedGroups.value = new Set(Array.from(grouped.value.keys()))
  }

  function collapseAll() {
    expandedGroups.value = new Set()
  }

  // ── Group selection ──

  function isGroupAllSelected(title: string): boolean {
    const group = grouped.value.get(title)
    if (!group || group.length === 0) return false
    return group.every((n) => selectedIds.value.has(n.id))
  }

  function toggleSelectGroup(title: string) {
    const group = grouped.value.get(title)
    if (!group) return
    const allGroupSelected = group.every((n) => selectedIds.value.has(n.id))
    const next = new Set(selectedIds.value)
    for (const n of group) {
      if (allGroupSelected) {
        next.delete(n.id)
      } else {
        next.add(n.id)
      }
    }
    selectedIds.value = next
  }

  // ── Group actions ──

  async function handleGroupMarkRead(title: string) {
    const group = grouped.value.get(title)
    if (!group) return
    await Promise.all(group.filter((n) => !n.isRead).map((n) => markRead(n.id)))
  }

  async function handleGroupDelete(title: string) {
    const group = grouped.value.get(title)
    if (!group) return
    const idsToRemove = new Set(group.map((n) => n.id))
    for (const id of idsToRemove) {
      selectedIds.value.delete(id)
    }
    selectedIds.value = new Set(selectedIds.value)
    await Promise.all(Array.from(idsToRemove).map((id) => remove(id)))
  }

  // ── Single-item actions ──

  async function handleMarkRead(id: number) {
    await markRead(id)
  }

  async function handleDelete(id: number) {
    selectedIds.value.delete(id)
    selectedIds.value = new Set(selectedIds.value)
    await remove(id)
  }

  async function handleMarkAllRead() {
    await markAllRead()
  }

  async function handleDeleteSelected() {
    await Promise.all(Array.from(selectedIds.value).map((id) => remove(id)))
    selectedIds.value = new Set()
  }

  function toggleSortField(field: SortField) {
    sortField.value = field
  }

  function toggleSortOrder() {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  }

  return {
    isLoading,
    selectedIds,
    expandedGroups,
    sortField,
    sortOrder,

    list,
    unreadCount,
    groupEntries,
    allSelected,
    selectedCount,

    load,

    toggleSelectAll,
    toggleSelect,
    isGroupAllSelected,
    toggleSelectGroup,

    toggleExpand,
    expandAll,
    collapseAll,

    handleMarkRead,
    handleDelete,
    handleMarkAllRead,
    handleDeleteSelected,
    handleGroupMarkRead,
    handleGroupDelete,

    toggleSortField,
    toggleSortOrder,
  }
}
