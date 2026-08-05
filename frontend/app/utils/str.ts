import PersonIcon from '@/assets/images/special/user.png'
import { getApiBase } from '@/constants/urls'
import { KyokusuAppRole } from '~/types/enums/role-enum'

export function correctProfileImageLink(link: string) {
  const backendUrl = getApiBase()
  return `${backendUrl}/${link}`
}

export function staticImage(picture: string | null | undefined) {
  if (picture && picture.startsWith('/uploads')) {
    const backendUrl = getApiBase()
    return `${backendUrl}${picture}`
  } else if (picture) {
    return picture
  } else {
    return PersonIcon
  }
}

export function hasSpecialSymbols(str: string): boolean {
  return /[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]/.test(str)
}

export const formatDateUserActivity = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return 'Только что'
  if (minutes < 60) return `${minutes} м. назад`
  if (hours < 24) return `${hours} ч. назад`
  if (days < 7) return `${days} д. назад`

  return date.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
  })
}

type Lang = 'ru' | 'eng'

const ROLE_LABELS: Record<KyokusuAppRole, Record<Lang, string>> = {
  [KyokusuAppRole.ADMIN]: {
    eng: 'admin',
    ru: 'Администратор',
  },
  [KyokusuAppRole.MODERATOR]: {
    eng: 'moderator',
    ru: 'Модератор',
  },
  [KyokusuAppRole.PUBLISHER]: {
    eng: 'publisher',
    ru: 'Публикатор',
  },
  [KyokusuAppRole.USER]: {
    eng: 'user',
    ru: 'Пользователь',
  },
}

export const roleConv = (role: KyokusuAppRole | `${KyokusuAppRole}`, lang: Lang) =>
  ROLE_LABELS[role][lang]

export const replaceTags = (text: string) => {
  return text.replace(/<[^>]*>/g, '')
}

export const textMaxNumValue = (value: number, max: number) => {
  return value > max ? `${max}+` : value
}

export const isStringInt = (str: string) => {
  const trimmed = str.trim()
  if (!trimmed) return false

  const num = Number(trimmed)

  if (isNaN(num) || num < 0) return false

  return Number.isInteger(num)
}
