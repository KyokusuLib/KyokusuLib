import { KyokusuAppRole } from '@/types/enums/role-enum'
import { DashboardRowUserStatus, DashboardRowUserGender } from '@/types/enums/dashboard-table'
import type { DashboardUser } from '@/types/frontend/dashboard/users'

import { parseDateToLocale } from '@/utils/date'
import { roleConv } from '@/utils/str'

export function mapToRows(user: DashboardUser) {
  return {
    id: user.id,
    name: user.name,
    role: user.role,
    status: user.status,
    level: user.user_level?.level ?? null,
    levelDisplay: user.user_level?.level,
    registered: parseDateToLocale(user.create_at),
  }
}

const ROLE_COLORS: Record<KyokusuAppRole, string> = {
  [KyokusuAppRole.ADMIN]: 'bg-rose-500/10 text-rose-600 dark:text-rose-400',
  [KyokusuAppRole.MODERATOR]: 'bg-blue-500/10 text-blue-600 dark:text-blue-400',
  [KyokusuAppRole.PUBLISHER]: 'bg-amber-500/10 text-amber-600 dark:text-amber-400',
  [KyokusuAppRole.USER]: 'bg-zinc-500/10 text-zinc-600 dark:text-zinc-400',
}

const STATUS_COLORS: Record<DashboardRowUserStatus, string> = {
  [DashboardRowUserStatus.online]: 'bg-emerald-500',
  [DashboardRowUserStatus.offline]: 'bg-zinc-400',
  [DashboardRowUserStatus.idle]: 'bg-yellow-500',
  [DashboardRowUserStatus.ban]: 'bg-red-500',
}

const STATUS_LABELS: Record<DashboardRowUserStatus, string> = {
  [DashboardRowUserStatus.online]: 'Онлайн',
  [DashboardRowUserStatus.offline]: 'Офлайн',
  [DashboardRowUserStatus.idle]: 'Отошёл',
  [DashboardRowUserStatus.ban]: 'Заблокирован',
}

const GENDER_LABELS: Record<DashboardRowUserGender, string> = {
  [DashboardRowUserGender.Male]: 'Мужской',
  [DashboardRowUserGender.Female]: 'Женской',
  [DashboardRowUserGender.Hidden]: 'Не указан',
}

const LEVEL_TIERS = [
  { max: 3, bg: 'bg-zinc-500', text: 'text-amber-500' },
  { max: 7, bg: 'bg-amber-600', text: 'text-white' },
  { max: 12, bg: 'bg-stone-400', text: 'text-white' },
  { max: 18, bg: 'bg-yellow-500', text: 'text-black' },
  { max: 25, bg: 'bg-emerald-500', text: 'text-white' },
  { max: Infinity, bg: 'bg-rose-500', text: 'text-white' },
] as const

export const getUserRoleColor = (role: string) =>
  ROLE_COLORS[role as KyokusuAppRole | `${KyokusuAppRole}`]
export const getUserStatusColor = (status: string) =>
  STATUS_COLORS[status as DashboardRowUserStatus | `${DashboardRowUserStatus}`]

export const formatRole = (role: string): string =>
  roleConv(role as KyokusuAppRole | `${KyokusuAppRole}`, 'ru')
export const formatStatus = (
  status: DashboardRowUserStatus | `${DashboardRowUserStatus}`,
): string => STATUS_LABELS[status]
export const formatGender = (gender: string): string =>
  GENDER_LABELS[gender as DashboardRowUserGender | `${DashboardRowUserGender}`]

export function getLevelColor(level: number | null | undefined): string {
  if (level === null) return 'bg-zinc-300 dark:bg-zinc-600 text-zinc-500 dark:text-zinc-400'
  const tier = LEVEL_TIERS.find((t) => level <= t.max)
  return tier
    ? `${tier.bg} ${tier.text}`
    : `${LEVEL_TIERS[LEVEL_TIERS.length - 1]?.bg} ${LEVEL_TIERS[LEVEL_TIERS.length - 1]?.text}`
}
