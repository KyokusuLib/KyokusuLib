<script setup lang="ts">
import AnimatedBook from '~/components/features/AnimatedBook/AnimatedBook.vue'
</script>

<template>
  <section
    class="relative w-full h-[60vh] min-h-[500px] overflow-hidden flex items-center cursor-default"
  >
    <div class="container mx-auto px-4 flex justify-center items-center relative z-10">
      <div class="flex flex-col md:flex-row items-center justify-center gap-12 md:gap-32">
        <div
          class="flex flex-col items-center md:items-start text-center md:text-left space-y-6 select-none"
        >
          <div class="space-y-2">
            <h1
              class="text-6xl md:text-9xl font-black tracking-tighter text-zinc-700/80 dark:text-white uppercase italic drop-shadow-2xl"
            >
              KyokusuLib
            </h1>
            <div
              class="h-1.5 w-64 md:w-136 bg-linear-to-r from-yellow-600 to-yellow-400 rounded-full shadow-[0_0_20px_rgba(234,179,8,0.4)]"
            ></div>
          </div>

          <p class="text-xl md:text-3xl font-medium text-zinc-500 leading-tight">
            Cайт для чтения
            <span class="text-zinc-700 dark:text-zinc-200 italic">манги</span> и
            <span class="text-zinc-700 dark:text-zinc-200 italic">ранобе</span>
          </p>

          <div class="pt-4 flex gap-4">
            <NuxtLink
              to="/catalog"
              class="flex items-center gap-4 px-8 py-4 bg-zinc-400 dark:bg-white text-black font-black rounded-2xl hover:bg-yellow-500 hover:shadow-[0_0_30px_rgba(255,255,255,0.1)] transition-all duration-500 uppercase text-xs tracking-widest"
            >
              Открыть библиотеку
              <Icon name="ph:arrow-right" size="16" />
            </NuxtLink>
          </div>
        </div>

        <div class="hidden lg:block relative shrink-0">
          <AnimatedBook />
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.italic-shadow {
  filter: drop-shadow(0 10px 20px rgba(0, 0, 0, 0.8));
}

section::before {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 50%;
  height: 100%;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.03) 0%, transparent 70%);
  transform: translate(-50%, -50%);
  pointer-events: none;
}
</style>
