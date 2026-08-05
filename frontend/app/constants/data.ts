import type { MostSearched } from '~/types/frontend/search/searches'

export const COUNTRIES_LIST = [
  'Япония',
  'Южная Корея',
  'Китай',
  'США',
  'Россия',
  'Великобритания',
  'Франция',
  'Германия',
  'Канада',
  'Австралия',
]

export const METIER_LIST = ['Автор', 'Художник', 'Сценарист', 'Студия', 'Издатель']

export const NOVELA_TYPES = ['Ранобе', 'Манга', 'Манхва']

export const NOVELA_STATUSES = ['Анонс', 'Завершено', 'Заморожен', 'Продолжается']

export const AGE_RATINGS = ['0+', '6+', '12+', '16+', '18+']

export const TRANSLATION_STATUSES = ['Закончен', 'Заморожено', 'Продолжается', 'Нет переводчика']

export const NOVELA_GENRES = [
  'Фэнтези',
  'Фантастика',
  'Драма',
  'Комедия',
  'Спорт',
  'Детектив',
  'Романтика',
  'Приключения',
  'Мистика',
  'Психология',
]

export const NOVELA_CATEGORIES = [
  'Алхимия',
  'Ангелы',
  'Демоны',
  'Герои',
  'Гильдии',
  'Геймеры',
  'ГГ имба',
  'ГГ женщина',
  'Выживание',
  'Дружба',
]

export const YEAR_MIN = 1900
export const YEAR_MAX = new Date().getFullYear() + 1
export const YEAR_DEFAULT = '2006'
export const ACTIVITY_PAGE_SIZE = 3

// LIMIT FOR FETCH IN MAIN CONTENT
export const NOVELA_FETCH_LIMIT = 10
export const NOVELA_FETCH_LIMIT_NEW = 7
export const NOVELA_FETCH_LIMIT_POPULAR = 7
export const NOVELA_FETCH_LIMIT_LAST_UPDATED = 5
export const NOVELA_FETCH_LIMIT_TRENDING = 5

export const NOVELA_MAX_COMMENT_LENGTH = 5000

// Updaet last user activity timestamp (one time in 30 seconds)
export const USER_ACTIVITY_INTERVAL = 30000

// Idle timeout for user activity (in milliseconds)
export const USER_ACTIVITY_IDLE_TIMEOUT = 60 * 1000

// Catalog
export const NOVELA_CATALOG_PAGE_SIZE = 30

// CatalogNovelaCard
export const NOVELA_CATALOG_CARD_TITLE_MAX_LENGTH = 30
export const NOVELA_CATALOG_CARD_TITLE_DEVIDE = 43

// Header
export const MAX_UNREAD_COUNT = 99

// COMMENTS
export const REPORT_COMMNET_REASONS = [
  { id: 'spam', label: 'Спам и реклама' },
  { id: 'insult', label: 'Оскорбления и агрессия' },
  { id: 'spoiler', label: 'Спойлеры без маскировки' },
  { id: 'inappropriate', label: 'Непристойный контент' },
  { id: 'other', label: 'Другое' },
]

// Most searched values
export const MOST_SEARCHED_DEFAULT: MostSearched[] = [
  {
    id: 1,
    type: 'genres',
    label: 'Детектив',
  },
  {
    id: 2,
    type: 'genres',
    label: 'Романтика',
  },
  {
    id: 3,
    type: 'categories',
    label: 'Перерождение',
  },
  {
    id: 4,
    type: 'categories',
    label: 'Академия',
  },
]
