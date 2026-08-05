export type NovelaSort = 'asc' | 'desc'

export interface NovelaFilters {
  search?: string
  genres?: string[]
  categories?: string[]
  type?: string
  status?: string
  translationStatus?: string
  chaptersFrom?: number
  chaptersTo?: number
  sort?: NovelaSort
  limit?: number
  offset?: number
}
