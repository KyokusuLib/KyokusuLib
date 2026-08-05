import type { UserLevel } from '@/types/backend/user'

export interface DashboardUserTableRow {
  id: number
  name: string
  role: string
  level: UserLevel
  levelDisplay: string
  registered: string
}
