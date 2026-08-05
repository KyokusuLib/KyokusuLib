import { COUNTRIES_LIST, AGE_RATINGS } from '@/constants/data'

export interface NovelaFilterState {
  search: string
  genres: string[]
  categories: string[]
  status: string
  translationStatus: string
  type: string
  chaptersFrom: number | null
  chaptersTo: number | null
  yearFrom: number | null
  yearTo: number | null
  country: string[]
  ageRating: string
  sort: string
}

export const DEFAULT_FILTER_STATE: NovelaFilterState = {
  search: '',
  genres: [],
  categories: [],
  status: '',
  translationStatus: '',
  type: '',
  chaptersFrom: null,
  chaptersTo: null,
  yearFrom: null,
  yearTo: null,
  country: [],
  ageRating: '',
  sort: '',
}

export interface NovelaFilterOption {
  value: string
  label: string
}

export const STATUS_OPTIONS: NovelaFilterOption[] = [
  { value: 'Продолжается', label: 'Продолжается' },
  { value: 'Завершено', label: 'Завершён' },
  { value: 'Заморожен', label: 'Заморожен' },
  { value: 'Анонс', label: 'Анонс' },
  { value: 'Отменён', label: 'Отменён' },
]

export const TRANSLATION_STATUS_OPTIONS: NovelaFilterOption[] = [
  { value: 'Продолжается', label: 'Переводится' },
  { value: 'Заморожено', label: 'Заморожено' },
  { value: 'Закончен', label: 'Закончен' },
  { value: 'Лицензирован', label: 'Лицензирован' },
  { value: 'Нет переводчика', label: 'Нет переводчика' },
]

export const CATEGORY_OPTIONS: NovelaFilterOption[] = [
  { value: 'Алхимия', label: 'Алхимия' },
  { value: 'Ангелы', label: 'Ангелы' },
  { value: 'Демоны', label: 'Демоны' },
  { value: 'Драконы', label: 'Драконы' },
  { value: 'Зомби', label: 'Зомби' },
  { value: 'Вампиры', label: 'Вампиры' },
  { value: 'Геймеры', label: 'Геймеры' },
  { value: 'Герои', label: 'Герои' },
  { value: 'Гильдии', label: 'Гильдии' },
  { value: 'ГГ имба', label: 'ГГ имба' },
  { value: 'ГГ женщина', label: 'ГГ женщина' },
  { value: 'ГГ мужчина', label: 'ГГ мужчина' },
  { value: 'Выживание', label: 'Выживание' },
  { value: 'Дружба', label: 'Дружба' },
  { value: 'Империя', label: 'Империя' },
  { value: 'Королевство', label: 'Королевство' },
  { value: 'Магия', label: 'Магия' },
]

export const GENRE_OPTIONS: NovelaFilterOption[] = [
  { value: 'Фэнтези', label: 'Фэнтези' },
  { value: 'Романтика', label: 'Романтика' },
  { value: 'Драма', label: 'Драма' },
  { value: 'Приключения', label: 'Приключения' },
  { value: 'Детектив', label: 'Детектив' },
  { value: 'Фантастика', label: 'Фантастика' },
  { value: 'Боевик', label: 'Боевик' },
  { value: 'Триллер', label: 'Триллер' },
  { value: 'Ужасы', label: 'Ужасы' },
  { value: 'Комедия', label: 'Комедия' },
  { value: 'Мистика', label: 'Мистика' },
  { value: 'Психология', label: 'Психология' },
]

export const NOVELA_TYPE_OPTIONS: NovelaFilterOption[] = [
  { value: 'Ранобе', label: 'Ранобе' },
  { value: 'Веб-новелла', label: 'Веб-новелла' },
  { value: 'Манга', label: 'Манга' },
  { value: 'Манхва', label: 'Манхва' },
]

export const COUNTRY_OPTIONS: NovelaFilterOption[] = COUNTRIES_LIST.map((nameCountry) => {
  return {
    label: nameCountry,
    value: nameCountry,
  }
})

export const AGE_RATING_OPTIONS: NovelaFilterOption[] = AGE_RATINGS.map((ageRating) => {
  return {
    label: ageRating,
    value: ageRating,
  }
})

export const SORT_OPTIONS: NovelaFilterOption[] = [
  { value: 'new', label: 'Новизне' },
  { value: 'popular', label: 'Популярности' },
  { value: 'rating', label: 'Рейтингу' },
  { value: 'updated', label: 'Дате добавления' },
  { value: 'trending', label: 'Трендам' },
]
