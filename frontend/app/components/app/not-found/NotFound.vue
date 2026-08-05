<script setup lang="ts">
import ArrowBackIcon from '@/assets/images/special/back-arrow.png'
</script>

<template>
  <div
    class="min-h-screen bg-zinc-50 dark:bg-[#0f0f0f] flex items-center justify-center p-4 sm:p-6 relative overflow-hidden transition-colors duration-300"
  >
    <div
      class="absolute inset-0 bg-[radial-gradient(circle_at_50%_50%,#e4e4e7_0%,#fafafa_100%)] dark:bg-[radial-gradient(circle_at_50%_50%,#1a1a1a_0%,#0f0f0f_100%)] transition-colors duration-300"
    ></div>

    <div class="relative z-10 max-w-3xl w-full animate-enter">
      <div
        class="backdrop-blur-xl bg-white/70 dark:bg-[#161616]/70 border border-zinc-200 dark:border-zinc-800/80 rounded-2xl sm:rounded-3xl p-6 sm:p-12 shadow-2xl shadow-zinc-200/50 dark:shadow-black/80 ring-1 ring-zinc-200 dark:ring-zinc-800 hover:ring-zinc-300 dark:hover:ring-zinc-700 transition-all duration-500"
      >
        <div class="relative mb-6 sm:mb-8 flex justify-center select-none">
          <h1
            class="relative text-7xl sm:text-9xl md:text-[10rem] lg:text-[14rem] leading-none font-black tracking-tighter"
          >
            <span
              class="bg-gradient-to-r from-zinc-400 via-zinc-600 to-zinc-400 dark:from-zinc-400 dark:via-zinc-300 dark:to-zinc-400 bg-clip-text text-transparent drop-shadow-2xl"
            >
              404
            </span>

            <span
              class="absolute inset-0 text-zinc-300/60 dark:text-zinc-700/60 blur-sm scale-105 select-none"
              aria-hidden="true"
              >404</span
            >
            <span
              class="absolute inset-0 text-zinc-400/50 dark:text-zinc-600/50 blur-md scale-110 select-none"
              aria-hidden="true"
              >404</span
            >
            <span
              class="absolute inset-0 text-zinc-500/40 blur-xl scale-115 animate-pulse select-none"
              aria-hidden="true"
              >404</span
            >

            <span
              class="absolute inset-0 text-zinc-200 dark:text-zinc-800 blur-sm opacity-0 glitch select-none"
              aria-hidden="true"
              >404</span
            >
            <span
              class="absolute inset-0 text-zinc-400 dark:text-zinc-600 blur-sm opacity-0 glitch-delay select-none"
              aria-hidden="true"
              >404</span
            >
          </h1>
        </div>

        <h2
          class="text-center text-2xl sm:text-4xl md:text-6xl font-bold tracking-wider mb-4 sm:mb-6 px-2"
        >
          <span
            class="bg-gradient-to-r from-zinc-700 to-zinc-500 dark:from-zinc-300 dark:to-zinc-400 bg-clip-text text-transparent drop-shadow-lg shadow-zinc-500/20 dark:shadow-zinc-400/50 block sm:inline"
          >
            СТРАНИЦА
          </span>
          <span
            class="bg-gradient-to-r from-zinc-700 to-zinc-500 dark:from-zinc-300 dark:to-zinc-400 bg-clip-text text-transparent drop-shadow-lg shadow-zinc-500/20 dark:shadow-zinc-400/50 block sm:inline sm:ml-3"
          >
            НЕ НАЙДЕНА
          </span>
        </h2>

        <div
          class="text-center space-y-2 text-zinc-500 dark:text-zinc-400 text-base sm:text-lg md:text-xl font-light tracking-wide px-4"
        >
          <p>Страница не найдена или была перемещена</p>
          <p>Проверьте адрес и попробуйте снова</p>
        </div>

        <div class="flex justify-center mt-8 sm:mt-12">
          <NuxtLink to="/" class="w-full sm:w-auto">
            <button
              class="group relative w-full sm:w-auto px-6 py-4 sm:px-10 sm:py-5 bg-gradient-to-r from-zinc-100 to-zinc-200 dark:from-zinc-800 dark:to-zinc-700 rounded-xl sm:rounded-2xl font-bold text-lg sm:text-xl tracking-wider hover:scale-[1.02] sm:hover:scale-110 active:scale-95 transition-all duration-300 flex items-center justify-center gap-3 sm:gap-4 overflow-hidden cursor-pointer border border-zinc-200 dark:border-transparent"
            >
              <span
                class="absolute inset-0 bg-black/5 dark:bg-white/10 translate-y-full group-hover:translate-y-0 transition-transform duration-500"
              ></span>

              <img
                :src="ArrowBackIcon"
                alt="back"
                class="w-6 h-6 sm:w-8 sm:h-8 invert-0 dark:invert group-hover:-translate-x-2 transition-transform duration-300"
              />

              <span
                class="relative z-10 text-zinc-800 dark:text-zinc-200 drop-shadow-lg whitespace-nowrap"
              >
                ВЕРНУТЬСЯ НАЗАД
              </span>
            </button>
          </NuxtLink>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.animate-enter {
  animation: fadeInUp 1.2s cubic-bezier(0.16, 1, 0.3, 1) forwards;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(40px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@keyframes glitch {
  0%,
  100% {
    transform: translate(0);
  }
  20% {
    transform: translate(-3px, 3px);
  }
  40% {
    transform: translate(-3px, -3px);
  }
  60% {
    transform: translate(3px, 3px);
  }
  80% {
    transform: translate(3px, -3px);
  }
}

.glitch {
  animation: glitch 4s infinite;
}
.glitch-delay {
  animation: glitch 4s infinite 0.4s;
}
</style>
