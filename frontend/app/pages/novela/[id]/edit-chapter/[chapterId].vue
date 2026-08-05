<script setup lang="ts">
import { useNovela } from '~/composables/api/novela/useNovela'
import { $api } from '~/composables/api/useApi'
import { parseChapterContent, buildChapterContent } from '~/utils/chapterContent'
import { uid } from '~/utils/chapter'
import ChapterEditorForm from '~/components/app/novela/ChapterEditorForm.vue'
import type { ChapterReaderResponse } from '~/types/backend/novela'
import type { UploadedImage } from '~/composables/ui/useImageUpload'

const route = useRoute()
const router = useRouter()
const novelaId = route.params.id as string
const chapterId = route.params.chapterId as string

const { novela, fetchNovela, updateChapter, deleteChapterImages, addChapterImage } = useNovela()

await useAsyncData(`novela-edit-chapter-${novelaId}`, () => fetchNovela(novelaId))

const { data: chapterData } = await useAsyncData(`chapter-${chapterId}`, () =>
  $api<ChapterReaderResponse>(`/api/novela/chapters/${chapterId}`),
)

const chapter = chapterData.value
const parsed = chapter ? parseChapterContent(chapter.content) : null

// Convert gallery images (chapter.images) into editor blocks + UploadedImage[]
let galleryBlocks = null
let galleryImages: UploadedImage[] | null = null
if (chapter?.images?.length) {
  galleryImages = chapter.images.map((img) => ({
    id: `gallery-${img.id}`,
    sourceUrl: img.image_url,
    previewUrl: img.image_url,
    caption: img.caption,
    isGallery: true,
  }))
  galleryBlocks = galleryImages.map((gImg) => ({
    id: uid('block'),
    type: 'image' as const,
    content: gImg.id,
  }))
}

// Combine parsed inline blocks + gallery blocks
const combinedBlocks = parsed
  ? galleryBlocks
    ? [...parsed.blocks, ...galleryBlocks]
    : parsed.blocks
  : null
const combinedImages = parsed
  ? galleryImages
    ? [...parsed.images, ...galleryImages]
    : parsed.images
  : null

const selectedVolumeId = ref(chapter?.volume_id ?? '')
const chapterNumber = ref(chapter?.number ?? 1)
const chapterTitle = ref(chapter?.title ?? '')

const volumes = computed(() =>
  (novela.value?.volumes ?? []).map((v) => ({
    value: v.id,
    label: `Том ${v.number}${v.title ? ` — ${v.title}` : ''}`,
  })),
)

const isSubmitting = ref(false)
const formRef = ref()

async function onSubmit(_content: string) {
  isSubmitting.value = true
  try {
    const form = (formRef.value as any)?.form
    if (!form) return

    const allBlocks = form.blocks
    const allImages = form.images

    // Separate gallery vs inline
    const inlineImages = allImages.filter((img) => !img.isGallery)
    const galleryImgs: UploadedImage[] = []

    // Collect gallery images in block order (not image array order)
    for (const block of allBlocks) {
      if (block.type === 'image') {
        const img = allImages.find((i) => i.id === block.content)
        if (img?.isGallery) galleryImgs.push(img)
      }
    }

    // Build content from inline blocks only (gallery blocks → "" → filter(Boolean) removes them)
    const content = buildChapterContent(allBlocks, inlineImages)

    // Save chapter content (inline only)
    await updateChapter(
      Number(novelaId),
      chapterId,
      chapterNumber.value,
      chapterTitle.value || `Глава ${chapterNumber.value}`,
      content,
    )

    // Save gallery images separately: delete old, re-create
    if (galleryImgs.length > 0) {
      await deleteChapterImages(chapterId)
      await Promise.all(
        galleryImgs.map((img, i) => addChapterImage(chapterId, img.sourceUrl, img.caption, i)),
      )
    }

    router.push(`/novela/${novelaId}`)
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <div
    v-if="!chapter"
    class="flex flex-col items-center justify-center min-h-screen gap-4 text-zinc-500"
  >
    <Icon name="ph:warning-circle-bold" size="48" />
    <p>Глава не найдена</p>
    <button
      class="px-4 py-2 rounded-xl text-sm font-medium bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900 cursor-pointer"
      @click="router.push(`/novela/${novelaId}`)"
    >
      Вернуться к новелле
    </button>
  </div>

  <ChapterEditorForm
    v-else
    ref="formRef"
    :volumes="volumes"
    :selected-volume-id="selectedVolumeId"
    :chapter-number="chapterNumber"
    :chapter-title="chapterTitle"
    :poster-url="novela?.poster_url"
    :initial-blocks="combinedBlocks ?? undefined"
    :initial-images="combinedImages ?? undefined"
    :is-submitting="isSubmitting"
    title="Редактировать главу"
    :subtitle="`${novela?.title} — Глава ${chapter.number}`"
    submit-label="Сохранить изменения"
    submitting-label="Сохранение..."
    @update:selected-volume-id="selectedVolumeId = $event"
    @update:chapter-number="chapterNumber = $event"
    @update:chapter-title="chapterTitle = $event"
    @submit="onSubmit"
  >
    <template #back>
      <button
        class="flex items-center p-2 rounded-xl bg-white/50 dark:bg-zinc-900/50 backdrop-blur-md border border-zinc-200/50 dark:border-zinc-800/50 hover:bg-zinc-200 dark:hover:bg-zinc-800 transition-colors cursor-pointer"
        @click="router.push(`/novela/reader/${novelaId}/${chapterId}`)"
      >
        <Icon name="ph:arrow-left-bold" size="20" />
      </button>
    </template>
    <template #actions>
      <button
        type="button"
        class="flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium border border-zinc-200 dark:border-zinc-700 text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 bg-zinc-800/70 cursor-pointer"
        @click="router.push(`/novela/${novelaId}`)"
      >
        <Icon name="ph:x-bold" size="16" />
        Отмена
      </button>
    </template>
  </ChapterEditorForm>
</template>
