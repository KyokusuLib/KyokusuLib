<script setup lang="ts">
import LogoIcon from '@/assets/images/Kyokusu/kyokusulib_logo_2.png'

const lineWidths = [85, 72, 91, 64, 78, 55]
</script>

<template>
  <div class="book-container cursor-default select-none">
    <div class="book">
      <!-- Корешок книги -->
      <div class="book-spine bg-zinc-950 shadow-2xl"></div>

      <!-- ПРАВАЯ СТРАНИЦА (Остается неподвижной, содержит описание) -->
      <div class="page right-side bg-zinc-50 border-l border-zinc-200 shadow-inner">
        <div class="content p-5 flex flex-col h-full relative z-10">
          <!-- Красная линия полей -->
          <div class="absolute left-4 top-0 bottom-0 w-px bg-red-500/10"></div>

          <div class="pl-4 space-y-3">
            <div class="flex items-center gap-2">
              <span class="text-[10px] font-black uppercase tracking-[0.2em] text-yellow-600"
                >KyokusuLib</span
              >
            </div>

            <h4 class="text-[13px] font-black leading-tight text-zinc-900 uppercase italic">
              Портал
            </h4>

            <div class="space-y-2">
              <p class="text-[9px] leading-relaxed text-zinc-600 font-bold">
                Современная платформа для чтения манги и ранобэ с фокусом на производительность и
                удобство.
              </p>
            </div>
          </div>

          <div class="mt-auto pl-4 flex justify-between items-end">
            <span class="text-[8px] font-black text-zinc-300">Страница 01</span>
            <img :src="LogoIcon" class="w-6 h-6 grayscale" />
          </div>
        </div>
        <!-- Тень сгиба -->
        <div
          class="absolute inset-y-0 left-0 w-8 bg-linear-to-r from-black/10 to-transparent"
        ></div>
      </div>

      <!-- Внутренняя страница 2 (Улетает влево) -->
      <div class="page p2 bg-zinc-100">
        <div class="content p-4 flex flex-col items-center justify-center h-full">
          <Icon name="ph:scribble-loop-bold" size="40" class="text-zinc-300 opacity-50" />
          <div class="mt-4 space-y-2 w-full px-4">
            <div class="h-1 w-full bg-zinc-200 rounded-full"></div>
            <div class="h-1 w-2/3 bg-zinc-200 rounded-full"></div>
          </div>
        </div>
      </div>

      <!-- Внутренняя страница 1 (Улетает влево) -->
      <div class="page p1 bg-zinc-50">
        <div class="content p-4">
          <div class="margin-line"></div>
          <div class="mt-6 space-y-2 px-4">
            <div
              v-for="(width, index) in lineWidths"
              :key="index"
              class="h-1 bg-zinc-200 rounded-full"
              :style="{ width: width + '%' }"
            ></div>
          </div>
        </div>
      </div>

      <!-- ПЕРЕДНЯЯ ОБЛОЖКА -->
      <div class="page front-wrap">
        <!-- Внешняя сторона -->
        <div
          class="front-side bg-zinc-900 border border-zinc-800 flex flex-col items-center justify-center"
        >
          <div class="absolute inset-2 border border-yellow-500/10 rounded-sm"></div>
          <img
            :src="LogoIcon"
            class="w-16 h-16 object-contain brightness-0 invert opacity-90 drop-shadow-[0_10px_20px_rgba(0,0,0,0.5)]"
          />
          <div class="mt-4 flex flex-col items-center">
            <div class="h-px w-8 bg-yellow-500/40 mb-1"></div>
            <span class="text-[7px] font-black uppercase tracking-[0.4em] text-zinc-500"
              >KyokusuLib</span
            >
          </div>
        </div>
        <!-- Внутренняя сторона -->
        <div class="inner-side bg-zinc-200 border-r border-zinc-300">
          <div
            class="absolute inset-0 opacity-10 bg-[url('https://www.transparenttextures.com/patterns/paper-fibers.png')]"
          ></div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.book-container {
  perspective: 2000px;
  width: 170px;
  height: 230px;
}

.book {
  width: 100%;
  height: 100%;
  position: relative;
  left: 50px;
  transform-style: preserve-3d;
  transform: rotateY(-20deg) rotateX(5deg);
  animation: bookShow 1.4s cubic-bezier(0.34, 1.56, 0.64, 1) forwards;
}

.page {
  position: absolute;
  width: 100%;
  height: 100%;
  top: 0;
  left: 0;
  transform-origin: left;
  border-radius: 2px 8px 8px 2px;
  overflow: hidden;
}

.right-side {
  z-index: 1;
  transform: translateZ(-2px);
}

.front-wrap {
  z-index: 10;
  transform-style: preserve-3d;
  animation: frontFlip 1.2s 0.5s cubic-bezier(0.645, 0.045, 0.355, 1) forwards;
}

.front-side,
.inner-side {
  position: absolute;
  width: 100%;
  height: 100%;
  backface-visibility: hidden;
}

.inner-side {
  transform: rotateY(180deg);
}

/* Анимация перелистывания */
.p1 {
  z-index: 9;
  animation: pageFlip 1.2s 0.7s cubic-bezier(0.645, 0.045, 0.355, 1) forwards;
}
.p2 {
  z-index: 8;
  animation: pageFlip 1.2s 0.85s cubic-bezier(0.645, 0.045, 0.355, 1) forwards;
}

@keyframes bookShow {
  0% {
    opacity: 0;
    transform: rotateY(0deg) rotateX(10deg) scale(0.8);
  }
  100% {
    opacity: 1;
    transform: rotateY(-25deg) rotateX(5deg) scale(1.1);
  }
}

@keyframes frontFlip {
  100% {
    transform: rotateY(-160deg);
  }
}

@keyframes pageFlip {
  100% {
    transform: rotateY(-155deg);
    box-shadow: -10px 10px 30px rgba(0, 0, 0, 0.2);
  }
}

.book-spine {
  position: absolute;
  width: 36px;
  height: 100%;
  left: -18px;
  transform: rotateY(-90deg);
  transform-origin: right;
  border-radius: 5px 0 0 5px;
  z-index: 5;
  background: linear-gradient(to right, #09090b, #18181b);
}

.margin-line {
  position: absolute;
  left: 30px;
  top: 0;
  bottom: 0;
  width: 1px;
  background: rgba(239, 68, 68, 0.1);
}
</style>
