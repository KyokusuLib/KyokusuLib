<script setup lang="ts">
import { Tooltip } from '@kyokusu-ui/vue'
import { SORT_OPTIONS, NOVELA_TYPE_OPTIONS } from '@/types/frontend/query/novela-filters'

const socialLinks = [
  { name: 'Telegram', icon: 'ph:telegram-logo', href: 'https://t.me/KyokusuLib' },
  { name: 'Discord', icon: 'ph:discord-logo', href: '#' },
  { name: 'GitHub', icon: 'ph:github-logo', href: 'https://github.com/Lanxre/KyokusuLib' },
]
</script>

<template>
  <footer
    class="w-full border-t-2 border-zinc-200 dark:border-zinc-700 p-6 md:p-12 mt-auto bg-zinc-50 dark:bg-[#0f0f0f] text-zinc-900 dark:text-white transition-colors duration-300"
  >
    <div class="mx-auto max-w-6xl px-12">
      <div class="grid grid-cols-1 gap-12 sm:grid-cols-2 md:grid-cols-4">
        <div class="flex flex-col gap-4 items-center md:items-start text-center md:text-left">
          <NuxtLink to="/" class="group">
            <h1 class="text-xs font-black uppercase tracking-[0.2em] text-zinc-100">KyokusuLib</h1>
          </NuxtLink>
          <div class="text-sm leading-relaxed text-zinc-500 dark:text-zinc-400 cursor-default">
            <p>Есть вопросы?</p>
            <p>Пишите в техподдержку:</p>
            <a
              href="mailto:burtsevvitalij@gmail.com"
              class="block mt-1 font-bold text-zinc-900 dark:text-zinc-200 hover:text-yellow-500 transition-colors"
            >
              burtsevvitalij@gmail.com
            </a>
          </div>
        </div>

        <div class="flex flex-col gap-4 items-center md:items-start text-center md:text-left">
          <h2 class="text-xs font-black uppercase tracking-[0.2em] text-zinc-100">Информация</h2>
          <nav class="flex flex-col gap-2">
            <NuxtLink
              to="/privacy"
              class="text-sm font-medium text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors"
              >Политика конфиденциальности</NuxtLink
            >
            <NuxtLink
              to="/terms"
              class="text-sm font-medium text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors"
              >Условия использования</NuxtLink
            >
          </nav>
        </div>

        <div class="flex flex-col gap-4 items-center md:items-start text-center md:text-left">
          <h2 class="text-xs font-black uppercase tracking-[0.2em] text-zinc-100">Разделы</h2>
          <nav class="flex flex-col gap-2 text-sm font-medium text-zinc-600 dark:text-zinc-400">
            <NuxtLink to="/top" class="hover:text-zinc-900 dark:hover:text-white transition-colors"
              >Популярное</NuxtLink
            >
            <NuxtLink
              :to="{ path: '/catalog', query: { sort: SORT_OPTIONS[0]?.value } }"
              class="hover:text-zinc-900 dark:hover:text-white transition-colors"
              >Новинки</NuxtLink
            >
            <NuxtLink
              :to="{ path: '/catalog', query: { type: NOVELA_TYPE_OPTIONS[0]?.value } }"
              class="hover:text-zinc-900 dark:hover:text-white transition-colors"
              >Каталог ранобе</NuxtLink
            >
          </nav>
        </div>

        <div class="flex flex-col gap-4 full text-center items-center md:items-start md:text-left">
          <h2 class="text-xs font-black uppercase tracking-[0.2em] text-zinc-100">Сообщество</h2>
          <div class="grid grid-cols-2 gap-3 w-fit">
            <a
              v-for="social in socialLinks"
              :key="social.name"
              :href="social.href"
              class="flex items-center justify-center w-11 h-11 bg-zinc-200 hover:bg-yellow-500 dark:bg-zinc-800 dark:hover:bg-yellow-500 rounded-2xl transition-all duration-300 group shadow-sm hover:shadow-yellow-500/20"
              :aria-label="social.name"
              target="_blank"
            >
              <Tooltip :text="social.name" position="left" class="mt-1.5">
                <Icon
                  :name="social.icon"
                  size="24"
                  class="text-zinc-700 dark:text-zinc-200 group-hover:text-white transition-colors"
                />
              </Tooltip>
            </a>
          </div>
        </div>
      </div>

      <div
        class="flex flex-col sm:flex-row justify-between items-center pt-12 mt-12 border-t border-zinc-200 dark:border-zinc-800 gap-4"
      >
        <p class="text-xs font-bold text-zinc-400 uppercase tracking-widest">
          © 2026 KyokusuLib. Все права защищены.
        </p>
      </div>
    </div>
  </footer>
</template>
