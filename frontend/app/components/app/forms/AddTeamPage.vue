<script setup lang="ts">
import { reactive, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useTeam } from '@/composables/api/teams/useTeams'
import { Input as BaseInput, Label } from '@kyokusu-ui/vue'

const router = useRouter()
const { createTeam, isLoading } = useTeam()

const form = reactive({
  name: '',
  slug: '',
  description: '',
})

const generateSlug = (val: string) => {
  return val
    .toLowerCase()
    .trim()
    .replace(/[^\w\s-]/g, '')
    .replace(/[\s_]+/g, '_')
    .replace(/^[-_]+|[-_]+$/g, '')
}

watch(
  () => form.name,
  (newName) => {
    form.slug = generateSlug(newName)
  },
)

const handleSubmit = async () => {
  const result = await createTeam({
    name: form.name,
    slug: form.slug,
    description: form.description,
  })

  if (result) {
    router.push(`/team/${result.slug}`)
  }
}
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <main
      class="flex justify-center items-center grow bg-white dark:bg-radial-[at_center] dark:from-zinc-800 dark:to-zinc-950 dark:to-90% text-zinc-900 dark:text-zinc-200 transition-colors duration-300"
    >
      <div
        class="max-w-2xl mx-auto m-6 md:m-12 p-6 bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-sm transition-colors duration-300"
      >
        <div class="mb-2 p-4 flex flex-col items-center gap-2">
          <h2 class="text-4xl font-bold text-zinc-900 dark:text-white">Создание команды</h2>
          <p class="text-sm text-zinc-500 mt-1">
            Создайте команду для совместной работы над переводами.
          </p>
        </div>

        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <Label for="team-name" label="Название команды" />
              <BaseInput
                id="team-name"
                v-model="form.name"
                placeholder="Например: Imperial Scan"
                required
                :disabled="isLoading"
              />
            </div>

            <div>
              <Label for="team-slug" label="Короткое название" />
              <BaseInput
                id="team-slug"
                v-model="form.slug"
                placeholder="imperial-scan"
                required
                :disabled="isLoading"
              />
            </div>
          </div>

          <div class="space-y-1.5">
            <Label for="team-desc" label="Описание команды" />
            <textarea
              id="team-desc"
              v-model="form.description"
              rows="4"
              class="w-full bg-zinc-50 dark:bg-zinc-950 cursor-text border border-zinc-200 dark:border-zinc-800 rounded-xl p-4 text-sm outline-none focus:border-yellow-500/50 focus:ring-4 focus:ring-yellow-500/5 transition-all resize-none"
              placeholder="Расскажите о вашей команде..."
              :disabled="isLoading"
            ></textarea>
          </div>

          <div class="flex justify-end pt-4">
            <button
              type="submit"
              :disabled="isLoading || !form.name || !form.slug"
              class="w-full md:w-auto px-8 py-2 bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 font-bold rounded-xl active:scale-[0.98] transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
            >
              <div
                v-if="isLoading"
                class="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin"
              ></div>
              <span>{{ isLoading ? 'Создание...' : 'Создать команду' }}</span>
            </button>
          </div>
        </form>
      </div>
    </main>
  </div>
</template>
