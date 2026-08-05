import type { UserActivity } from '@/types/backend/user_activity'
import type { ActivityUIConfig } from '@/types/frontend/user-activity'
import { DEFAULT_CONFIG, STRATEGIES, type UserActivityType } from '@/constants/user-activity'

export function useActivityConfig() {
  const getActivityConfig = (activity: UserActivity): ActivityUIConfig => {
    const type = activity.activity_type as UserActivityType
    const strategy = STRATEGIES[type]
    if (!strategy) {
      return DEFAULT_CONFIG
    }
    return (strategy as (a: UserActivity) => ActivityUIConfig)(activity)
  }

  return {
    getActivityConfig,
  }
}
