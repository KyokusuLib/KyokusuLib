const YEAR_DEFAULT = 1970

export function parseDateToISO(input: string): string {
  if (!input) return ''
  const date = new Date(input).toISOString().split('T')[0]!
  return date
}

export function parseDateToLocale(input: string): string {
  return new Date(input).toLocaleDateString('ru-RU', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  })
}

export function parseStringToDate(input: string): Date {
  return new Date(input)
}

export function getStingYear(input: string): string {
  return new Date(input).getFullYear().toString()
}

export function fmtDateTime(dateStr: string | null | undefined): string {
  if (!dateStr) return '—'
  return new Date(dateStr).toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function fmtRelativeTime(dateStr: string | null | undefined, now: Date = new Date()): string {
  if (!dateStr) return '—'
  const date = new Date(dateStr)
  const diffMs = now.getTime() - date.getTime()
  if (diffMs < 0) return fmtDateTime(dateStr)

  const minutes = Math.floor(diffMs / 60000)
  if (minutes < 1) return 'только что'
  if (minutes < 60) {
    const m = Math.floor(minutes)
    const label = plural(m, 'минуту', 'минуты', 'минут')
    return `${m} ${label} назад`
  }

  const hours = Math.floor(minutes / 60)
  if (hours < 24) {
    const label = plural(hours, 'час', 'часа', 'часов')
    return `${hours} ${label} назад`
  }

  const days = Math.floor(hours / 24)
  if (days === 1) return 'вчера'
  if (days < 7) {
    const label = plural(days, 'день', 'дня', 'дней')
    return `${days} ${label} назад`
  }

  return fmtDateTime(dateStr)
}

function plural(n: number, one: string, few: string, many: string): string {
  const mod10 = n % 10
  const mod100 = n % 100
  if (mod10 === 1 && mod100 !== 11) return one
  if (mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14)) return few
  return many
}

export function fmtDate(dateStr: string | null | undefined): string {
  if (!dateStr) return 'Не указанно'

  const date = new Date(dateStr)

  if (date.getFullYear() === YEAR_DEFAULT) return 'Не указанно'

  return date.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
  })
}

export function fmtNumber(n: number | null | undefined): string {
  if (n === null) return '—'
  return n.toLocaleString('ru-RU')
}
