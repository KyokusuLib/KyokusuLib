import type { NovelaFilterState } from '@/types/frontend/query/novela-filters'

export interface CatalogFilterPreset {
  id: number
  user_id?: number
  name: string
  filters: NovelaFilterState
  created_at: string
}
