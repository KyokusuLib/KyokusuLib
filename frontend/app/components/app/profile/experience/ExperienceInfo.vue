<script setup lang="ts">
type Rank = {
  level: number
  experience: number
  name: string
}

const ranks: Rank[] = [
  { level: 0, experience: 0, name: 'Новичок' },
  { level: 1, experience: 1000, name: 'Ученик' },
  { level: 2, experience: 1500, name: 'Читатель' },
  { level: 3, experience: 2000, name: 'Практик' },
  { level: 4, experience: 2500, name: 'Эксперт' },
  { level: 5, experience: 3000, name: 'Мастер' },
  { level: 6, experience: 3500, name: 'Глава' },
  { level: 7, experience: 4000, name: 'Величие' },
]
</script>

<template>
  <div class="flex flex-col gap-4 font-medium text-foreground">
    <p class="mt-2">
      На сайте действует уникальная система рангов, которая отражает активность читателей через
      количество просмотров. Каждый уровень символизирует этап развития — от первых шагов до
      величия.
    </p>

    <p>Ниже представлен список доступных уровней и необходимое количество просмотров:</p>

    <ol class="list-decimal list-inside">
      <li v-for="rank in ranks" :key="rank.level" class="font-bold text-xl">
        {{ rank.name }} — {{ rank.experience }} опыта
      </li>
    </ol>

    <p>
      Этот путь подчёркивает вклад каждого читателя в развитие сообщества. Какой уровень сможете
      достичь вы? 🚀
    </p>
  </div>
</template>
