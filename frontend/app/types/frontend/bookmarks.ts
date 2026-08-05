export type BookmarkCategory = string | number

export interface Bookmark {
  user_id: number
  novela_id: number
  category: BookmarkCategory
  category_id?: number
  created_at: string
  updated_at: string
}

export interface BookmarkCategoryTab {
  id: BookmarkCategory
  name: string
  label: string
  icon: string
  user_id: number
  count: number
  visible: boolean
}

export function getBookmarkCategoryLabel(category: BookmarkCategory): string {
  switch (category) {
    case 'planned':
      return 'В планах'
    case 'reading':
      return 'Читаю'
    case 'completed':
      return 'Прочитано'
    case 'on_hold':
      return 'Отложено'
    case 'dropped':
      return 'Брошено'
    default:
      return category.toString()
  }
}
