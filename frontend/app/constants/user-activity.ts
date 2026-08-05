import type { ActivityUIConfig } from '@/types/frontend/user-activity'
import type { UserActivity } from '@/types/backend/user_activity'

export const ACTIVITY_TYPES = {
  NOVELA_BOOKMARK: 'novela_added_to_bookmarks',
  NOVELA_BOOKMARK_REMOVE: 'novela_removed_from_bookmarks',
  RANOBE_ADD: 'ranobe_added',
  USER_NOVELA_LIKE: 'user_novela_like',
  USER_NOVELA_LIKE_REMOVE: 'user_novela_like_remove',
  USER_NOVELA_RATING: 'user_novela_rating',
  ACHIEVEMENT_EARNED: 'achievement_earned',
  COMMENT_ADDED: 'comment_added',
  COMMENT_REMOVE: 'comment_remove',
  COMMENT_UPDATE: 'comment_update',
  COMMENT_LIKE: 'comment_like',
  COMMENT_REPORT: 'comment_report',
  DEFAULT: 'default',
} as const

export type UserActivityType = (typeof ACTIVITY_TYPES)[keyof typeof ACTIVITY_TYPES]

export const DEFAULT_CONFIG: ActivityUIConfig = {
  icon: 'activity',
  color: 'bg-zinc-100 text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400',
  title: 'Активность',
  description: 'Действие пользователя',
}

export const STRATEGIES: {
  [K in UserActivityType]?: (
    activity: Extract<UserActivity, { activity_type: K }>,
  ) => ActivityUIConfig
} = {
  [ACTIVITY_TYPES.RANOBE_ADD]: (activity) => ({
    icon: 'book',
    color: 'bg-blue-100 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400',
    title: 'Добавлено ранобэ',
    description: `"${activity.metadata.ranobe_title || 'Без названия'}" от ${activity.metadata.author || 'Неизвестный'}`,
  }),

  [ACTIVITY_TYPES.ACHIEVEMENT_EARNED]: (activity) => ({
    icon: 'trophy',
    color: 'bg-yellow-100 text-yellow-600 dark:bg-yellow-900/30 dark:text-yellow-400',
    title: 'Получено достижение',
    description: `${activity.metadata.name || 'Достижение'} — ${activity.metadata.desc || ''}`,
  }),

  [ACTIVITY_TYPES.NOVELA_BOOKMARK]: (activity) => ({
    icon: 'bookmark',
    color: 'bg-pink-100 text-pink-600 dark:bg-pink-900/30 dark:text-pink-400',
    title: 'Добавлено в закладки',
    description: `${activity.metadata.desc || 'Без описания'}: "${activity.metadata.novela_title || 'Без названия'}"`,
  }),

  [ACTIVITY_TYPES.NOVELA_BOOKMARK_REMOVE]: (activity) => ({
    icon: 'bookmark',
    color: 'bg-pink-100 text-pink-600 dark:bg-pink-900/30 dark:text-pink-400',
    title: 'Удалено из закладок',
    description: `${activity.metadata.desc || 'Без описания'}: "${activity.metadata.novela_title || 'Без названия'}"`,
  }),

  [ACTIVITY_TYPES.USER_NOVELA_LIKE]: (activity) => ({
    icon: 'heart',
    color: 'bg-red-100 text-red-600 dark:bg-red-900/30 dark:text-red-400',
    title: 'Понравилось',
    description: `${activity.metadata.desc || 'Без описания'}: "${activity.metadata.name || 'Без названия'}"`,
  }),

  [ACTIVITY_TYPES.USER_NOVELA_LIKE_REMOVE]: (activity) => ({
    icon: 'heart',
    color: 'bg-red-100 text-red-600 dark:bg-red-900/30 dark:text-red-400',
    title: 'Удалено из понравившихся',
    description: `${activity.metadata.desc || 'Без описания'}: "${activity.metadata.name || 'Без названия'}"`,
  }),

  [ACTIVITY_TYPES.USER_NOVELA_RATING]: (activity) => ({
    icon: 'star',
    color: 'bg-yellow-100 text-yellow-600 dark:bg-yellow-900/30 dark:text-yellow-400',
    title: 'Оценено',
    description: `${activity.metadata.desc || 'Без описания'}: "${activity.metadata.name || 'Без названия'}" на ${activity.metadata.rating} баллов`,
  }),

  [ACTIVITY_TYPES.COMMENT_ADDED]: (activity) => ({
    icon: 'comment',
    color: 'bg-cyan-100 text-cyan-600 dark:bg-cyan-900/30 dark:text-cyan-400',
    title: 'Добавлен комментарий',
    description: `${activity.metadata.desc || 'Без описания'} под: "${activity.metadata.novela_title || 'Без названия'}"`,
  }),

  [ACTIVITY_TYPES.COMMENT_REMOVE]: (activity) => ({
    icon: 'comment',
    color: 'bg-cyan-100 text-cyan-600 dark:bg-cyan-900/30 dark:text-cyan-400',
    title: 'Удален комментарий',
    description: `${activity.metadata.desc || 'Без описания'} под: "${activity.metadata.novela_title || 'Без названия'}"`,
  }),

  [ACTIVITY_TYPES.COMMENT_LIKE]: (activity) => ({
    icon: 'heart',
    color: 'bg-red-100 text-red-600 dark:bg-red-900/30 dark:text-red-400',
    title: 'Понравилось',
    description: `${activity.metadata.desc || 'Без описания'} под: "${activity.metadata.novela_title || 'Без названия'}"`,
  }),

  [ACTIVITY_TYPES.COMMENT_UPDATE]: (activity) => ({
    icon: 'comment',
    color: 'bg-cyan-100 text-cyan-600 dark:bg-cyan-900/30 dark:text-cyan-400',
    title: 'Обновлен комментарий',
    description: `${activity.metadata.desc || 'Без описания'} под: "${activity.metadata.novela_title || 'Без названия'}"`,
  }),

  [ACTIVITY_TYPES.COMMENT_REPORT]: (activity) => ({
    icon: 'flag',
    color: 'bg-amber-100 text-amber-600 dark:bg-amber-900/30 dark:text-amber-400',
    title: 'Пожаловался',
    description: `${activity.metadata.desc || 'Без описания'} под: "${activity.metadata.novela_title || 'Без названия'}"`,
  }),
}
