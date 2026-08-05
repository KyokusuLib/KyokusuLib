<script setup lang="ts">
import { ref } from 'vue'
import { useParse } from '@/composables/api/parse/useParse'
import { isStringInt } from '@/utils/str'

const rhId = ref('')
const { isLoading, parseRanobeHub } = useParse()

const handleParse = async () => {
  if (!rhId.value || !isStringInt(rhId.value)) return
  const success = await parseRanobeHub(rhId.value)
  if (success) {
    rhId.value = ''
  }
}
</script>

<template>
  <div class="space-y-8">
    <div>
      <h2 class="text-xl font-bold text-zinc-900 dark:text-white flex items-center gap-2">
        <Icon name="ph:download-simple-bold" class="text-yellow-500" />
        Импорт контента
      </h2>
      <p class="text-sm text-zinc-500 dark:text-zinc-400 mt-1">
        Парсинг новелл из внешних источников.
      </p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <div
        class="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6 shadow-sm"
      >
        <NuxtLink to="https://ranobehub.org/" class="cursor-pointer" target="_blank">
          <div class="flex items-center gap-3 mb-6 group group-hover:text-amber-500">
            <div
              class="w-12 h-12 rounded-2xl bg-orange-500/10 text-orange-500 flex items-center justify-center"
            >
              <Icon name="ph:globe-bold" size="24" />
            </div>
            <div>
              <h3 class="font-bold text-zinc-900 dark:text-white group-hover:text-amber-500">
                RanobeHub
              </h3>
              <p class="text-xs text-zinc-500">ranobehub.org</p>
            </div>
          </div>
        </NuxtLink>

        <div class="space-y-4">
          <div>
            <label
              class="block text-xs font-black uppercase tracking-widest text-zinc-400 mb-2 ml-1.5"
            >
              ID Новеллы
            </label>
            <input
              v-model="rhId"
              type="text"
              placeholder="Например: 123"
              class="w-full bg-zinc-100 dark:bg-zinc-800 border-none rounded-2xl px-4 py-3 text-sm focus:ring-2 focus:ring-yellow-500/50 outline-none transition-all"
            />
          </div>

          <button
            @click="handleParse"
            :disabled="isLoading || !rhId"
            class="w-full py-2 bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 font-bold rounded-2xl hover:opacity-90 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 cursor-pointer"
          >
            <Icon v-if="isLoading" name="ph:spinner-bold" class="animate-spin" />
            <Icon v-else name="ph:lightning-bold" />
            <span class="flex mr-4">
              {{ isLoading ? 'Парсинг...' : 'Начать импорт' }}
            </span>
          </button>
        </div>

        <div
          class="mt-6 p-4 bg-zinc-50 dark:bg-zinc-800/50 rounded-2xl border border-zinc-100 dark:border-zinc-700/50"
        >
          <p class="text-[10px] text-zinc-500 leading-relaxed">
            <Icon name="ph:info-bold" class="inline mr-1" />
            Введите числовой ID новеллы из URL (например, в
            <code>ranobehub.org/ranobe/123-title</code> ID это 123). Парсинг может занять несколько
            минут.
          </p>
        </div>
      </div>

      <div
        class="border-2 border-dashed border-zinc-200 dark:border-zinc-800 rounded-3xl p-6 flex flex-col items-center justify-center text-zinc-400 gap-3"
      >
        <Icon name="ph:plus-circle-bold" size="32" class="opacity-20" />
        <p class="text-xs font-bold uppercase tracking-widest">Скоро новые источники</p>
      </div>
    </div>
  </div>
</template>
