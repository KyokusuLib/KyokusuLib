export type RecentSearchType = 'query' | 'genre' | 'category' | 'author' | 'user' | 'team'

export interface MostSearched {
  id: number
  type: string
  label: string
}

export interface RecentSearch {
  text: string
  type: RecentSearchType
  id?: number
  category?: string
}
