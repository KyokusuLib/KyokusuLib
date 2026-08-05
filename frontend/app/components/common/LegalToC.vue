<script setup lang="ts">
export interface TocItem {
  id: string
  title: string
  icon: string
}

const props = defineProps<{
  items: TocItem[]
  title?: string
}>()

const activeId = ref<string | null>(null)

function scrollToSection(id: string) {
  const el = document.getElementById(id)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
    activeId.value = id
  }
}

function updateActiveSection() {
  const offsets = props.items
    .map((item) => {
      const el = document.getElementById(item.id)
      if (!el) return { id: item.id, top: Infinity }
      const rect = el.getBoundingClientRect()
      return { id: item.id, top: Math.abs(rect.top) }
    })
    .sort((a, b) => a.top - b.top)

  activeId.value = offsets[0]?.id ?? null
}

onMounted(() => {
  updateActiveSection()
  window.addEventListener('scroll', updateActiveSection, { passive: true })
})

onUnmounted(() => {
  window.removeEventListener('scroll', updateActiveSection)
})
</script>

<template>
  <nav class="w-full" aria-label="Содержание">
    <div
      v-if="title"
      class="text-center text-xs font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500 mb-6"
    >
      {{ title }}
    </div>
    <ul class="space-y-1">
      <li v-for="item in items" :key="item.id">
        <button
          class="w-full flex items-center gap-3 px-3 py-2 cursor-pointer rounded-xl text-sm text-left transition-all duration-200"
          :class="
            activeId === item.id
              ? 'bg-yellow-400/10 text-yellow-600 dark:text-yellow-400 font-semibold shadow-sm'
              : 'text-zinc-500 dark:text-zinc-400 hover:text-zinc-800 dark:hover:text-zinc-200 hover:bg-zinc-100 dark:hover:bg-zinc-800/60'
          "
          @click="scrollToSection(item.id)"
        >
          <span class="shrink-0 w-5 h-5 flex items-center justify-center">
            <Icon :name="item.icon" size="16" />
          </span>
          <span class="leading-snug">{{ item.title }}</span>
        </button>
      </li>
    </ul>
  </nav>
</template>
