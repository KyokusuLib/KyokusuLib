export enum NovelaSort {
  New = 'new',
  Popular = 'popular',
  Updated = 'updated',
  Trending = 'trending',
}

export interface NovelsQueryParams {
  limit?: number
  offset?: number
  order_by?: string
  order?: string
  search?: string
  sort?: NovelaSort
  author_id?: number
}
