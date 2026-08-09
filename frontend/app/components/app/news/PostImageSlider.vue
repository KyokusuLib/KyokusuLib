<script setup lang="ts">
const props = defineProps<{
  images: string[]
}>()

const current = ref(0)

const showControls = computed(() => props.images.length > 1)

const currentImage = computed(() => props.images[current.value])

function next() {
  current.value = (current.value + 1) % props.images.length
}

function prev() {
  current.value = (current.value - 1 + props.images.length) % props.images.length
}
</script>

<template>
  <div class="relative">
    <img
      :src="currentImage"
      alt="Изображение поста"
      loading="lazy"
      decoding="async"
      class="max-h-96 w-full rounded-lg object-contain"
    />

    <div v-if="showControls" class="absolute inset-y-0 left-0 flex items-center">
      <button
        type="button"
        aria-label="Предыдущее изображение"
        class="ml-2 flex h-9 w-9 cursor-pointer items-center justify-center rounded-full bg-black/40 text-white backdrop-blur transition hover:bg-black/60"
        @click="prev"
      >
        <Icon name="ph:caret-left-bold" size="18" />
      </button>
    </div>

    <div v-if="showControls" class="absolute inset-y-0 right-0 flex items-center">
      <button
        type="button"
        aria-label="Следующее изображение"
        class="mr-2 flex h-9 w-9 cursor-pointer items-center justify-center rounded-full bg-black/40 text-white backdrop-blur transition hover:bg-black/60"
        @click="next"
      >
        <Icon name="ph:caret-right-bold" size="18" />
      </button>
    </div>

    <div v-if="showControls" class="absolute inset-x-0 bottom-3 flex justify-center gap-1.5">
      <button
        v-for="(_, index) in images"
        :key="index"
        type="button"
        :aria-label="`Изображение ${index + 1}`"
        class="h-1.5 rounded-full transition-all"
        :class="index === current ? 'w-5 bg-white' : 'w-1.5 bg-white/50'"
        @click="current = index"
      />
    </div>
  </div>
</template>
