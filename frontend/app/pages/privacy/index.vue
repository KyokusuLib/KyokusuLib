<script setup lang="ts">
import PageLayout from '@/components/features/PageLayout/PageLayout.vue'
import { Separator, Table } from '@kyokusu-ui/vue'
import PrivacyCard from '@/components/common/Privacy/PrivacyCard.vue'
import type { TableHeader } from '@kyokusu-ui/vue'

useSeoMeta({
  title: 'Политика конфиденциальности - KyokusuLib',
  description:
    'Соглашение о конфиденциальности и защите персональных данных пользователей KyokusuLib',
})

interface PrivacyPositionOne {
  category: string
  type: string
  timeRange: string
}

const headers: TableHeader<keyof PrivacyPositionOne>[] = [
  { key: 'category', label: 'Категория персональных данных', align: 'left' },
  { key: 'type', label: 'Тип', align: 'center' },
  { key: 'timeRange', label: 'Временной диапазон', align: 'left' },
]

const table_data: PrivacyPositionOne[] = [
  {
    category: 'Адрес электронной почты',
    type: 'Конфиденциальный',
    timeRange:
      'До удаления аккаунта пользователем или отзыва согласия, но не более 2 лет с даты последней активности',
  },
  {
    category: 'Дата рождения',
    type: 'Конфиденциальный',
    timeRange: 'До удаления аккаунта пользователем или отзыва согласия',
  },
  {
    category: 'Публичные данные профиля пользователя',
    type: 'Общедоступные',
    timeRange: 'До удаления пользователем или модератором; по умолчанию — бессрочно',
  },
  {
    category: 'Комментарии и посты',
    type: 'Общедоступные',
    timeRange: 'До удаления пользователем или модератором; по умолчанию — бессрочно',
  },
]

const lastUpdated = computed(() =>
  new Date('2026-06-17').toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  }),
)
const sections = [
  {
    title: 'Общие положения',
    description:
      'Настоящая Политика конфиденциальности определяет порядок обработки и защиты персональной информации пользователей платформы KyokusuLib.',
    subcards: [
      {
        title: 'Статус документа',
        text: 'Использование сервиса означает безоговорочное согласие пользователя с настоящей Политикой и указанными в ней условиями обработки его персональной информации.',
      },
      {
        title: 'Изменение политики',
        text: 'Администрация KyokusuLib имеет право вносить изменения в настоящую Политику конфиденциальности. Новая редакция вступает в силу с момента ее размещения.',
      },
    ],
  },
  {
    title: 'Собираемая информация',
    description:
      'Мы собираем данные, которые необходимы для обеспечения работы функций платформы и улучшения пользовательского опыта.',
    subcards: [
      {
        title: 'Данные профиля',
        text: 'При регистрации мы сохраняем ваш email, имя пользователя и пароль (в зашифрованном виде). При использовании OAuth (Google/Discord) мы получаем только публичный идентификатор и email.',
      },
      {
        title: 'Технические данные',
        text: 'IP-адрес, тип устройства, версия браузера, файлы cookie и данные об активности (прогресс чтения, закладки, лайки).',
      },
      {
        title: 'Контент пользователя',
        text: 'Комментарии, рецензии и сообщения, которые вы добровольно размещаете на платформе.',
      },
    ],
  },
  {
    title: 'Цели обработки данных',
    description:
      'Персональная информация используется исключительно в рамках предоставления услуг сервиса.',
    subcards: [
      {
        title: 'Идентификация',
        text: 'Установление личности пользователя для доступа к личному кабинету и персональным настройкам (закладки, история чтения).',
      },
      {
        title: 'Коммуникация',
        text: 'Отправка уведомлений о новых главах, системных оповещений или ответов на ваши запросы в техподдержку.',
      },
      {
        title: 'Аналитика',
        text: 'Обезличенная статистика для понимания популярности произведений и улучшения интерфейса сайта.',
      },
    ],
  },
  {
    title: 'Защита информации',
    description:
      'Мы принимаем необходимые организационные и технические меры для защиты персональной информации пользователя от неправомерного доступа.',
    subcards: [
      {
        title: 'Шифрование',
        text: 'Все данные передаются по зашифрованному протоколу HTTPS. Пароли хранятся с использованием современных алгоритмов хеширования.',
      },
      {
        title: 'Доступ к данным',
        text: 'Доступ к личным данным пользователей имеют только авторизованные сотрудники администрации в рамках выполнения своих обязанностей.',
      },
    ],
  },
  {
    title: 'Права пользователя',
    description:
      'Вы обладаете полным контролем над своими данными в рамках действующего законодательства.',
    subcards: [
      {
        title: 'Доступ и исправление',
        text: 'Вы можете изменить свои данные в настройках профиля в любое время.',
      },
      {
        title: 'Удаление данных',
        text: 'Пользователь имеет право запросить полное удаление своего аккаунта и связанных с ним данных через службу поддержки или специальную функцию в настройках.',
      },
    ],
  },
]
</script>

<template>
  <PageLayout>
    <div class="w-full max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
      <div class="mb-12">
        <div
          class="flex flex-col sm:flex-row justify-start gap-2 sm:gap-4 items-start sm:items-center text-2xl sm:text-3xl font-black tracking-tight text-zinc-900 dark:text-white mb-6"
        >
          <span class="text-yellow-500">KyokusuLib</span>
          <Separator orientation="vertical" class="hidden sm:block h-8" />
          <span class="opacity-60">Пользовательское соглашение</span>
        </div>
        <h1
          class="text-4xl sm:text-5xl font-black uppercase tracking-tighter text-zinc-900 dark:text-white mb-4 leading-none"
        >
          Политика конфиденциальности
        </h1>
        <p class="text-zinc-500 dark:text-zinc-400 font-medium">
          Последнее обновление: {{ lastUpdated }}
        </p>
        <Separator class="mt-8" />
      </div>

      <div class="space-y-4">
        <PrivacyCard
          v-for="(section, index) in sections"
          :key="index"
          :index="index + 1"
          :card-title="section.title"
          :description="section.description"
          :subcards="section.subcards"
        />
      </div>
      <div class="h-full overflow-y-auto">
        <Table :headers="headers" :items="table_data" :searchable="false" :showActions="false">
        </Table>
      </div>
    </div>
  </PageLayout>
</template>
