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
