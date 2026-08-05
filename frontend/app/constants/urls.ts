import type { UserActivity } from '~/types/backend/user_activity'
import { ACTIVITY_TYPES } from './user-activity'

export function getApiBase() {
  try {
    const config = useRuntimeConfig()
    return config.public.apiBase
  } catch (_e) {
    return ''
  }
}

export function naviagateFromActivity(activity: UserActivity) {
  switch (activity.activity_type) {
    case ACTIVITY_TYPES.NOVELA_BOOKMARK:
      navigateTo(`/novela/${activity.target_id}`)
      break
    case ACTIVITY_TYPES.NOVELA_BOOKMARK_REMOVE:
      navigateTo(`/novela/${activity.target_id}`)
      break
    case ACTIVITY_TYPES.RANOBE_ADD:
      navigateTo(`/ranobe/${activity.target_id}`)
      break
    case ACTIVITY_TYPES.USER_NOVELA_LIKE:
      navigateTo(`/novela/${activity.target_id}`)
      break
    case ACTIVITY_TYPES.USER_NOVELA_LIKE_REMOVE:
      navigateTo(`/novela/${activity.target_id}`)
      break
    case ACTIVITY_TYPES.COMMENT_ADDED:
      navigateTo(`/novela/${activity.target_id}`)
      break
    case ACTIVITY_TYPES.COMMENT_REMOVE:
      navigateTo(`/novela/${activity.target_id}`)
      break
    case ACTIVITY_TYPES.COMMENT_UPDATE:
      navigateTo(`/novela/${activity.target_id}`)
      break
    case ACTIVITY_TYPES.COMMENT_LIKE:
      navigateTo(`/novela/${activity.target_id}`)
      break
    case ACTIVITY_TYPES.COMMENT_REPORT:
      navigateTo(`/novela/${activity.target_id}`)
      break
  }
}
