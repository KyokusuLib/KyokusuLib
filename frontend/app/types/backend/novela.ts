export interface NovelaChapterImage {
  id: number
  image_url: string
  caption: string
  position: number
}

export interface ChapterReaderResponse {
  id: string
  title: string
  number: number
  content: string
  images: NovelaChapterImage[]
  novela_id: number
  novela_title: string
  volume_number: number
  volume_id: string
  prev_chapter_id: string | null
  next_chapter_id: string | null
}

export interface NovelaChapter {
  id: string
  title: string
  number: number
  is_read?: boolean
  status?: string
  created_by?: number
  content: string
  images: NovelaChapterImage[]
}

export interface NovelaVolume {
  id: string
  title: string
  number: number
  status?: string
  created_by?: number
  chapters: NovelaChapter[]
}

export interface NovelaAuthor {
  id: number
  name: string
  role: string
}

export interface NovelaAuthorDetails {
  id: number
  name: string
  country: string
  picture?: string
  metier: string
  bio: string
}

export interface NCItem {
  value: number
  count: number
}

export interface NovelaRatingCategory {
  total: number
  total_rating: number
  nc_items: NCItem[]
}

export interface NovelaBookmarkCategory {
  total: number
  nc_items: BookmarkNCItem[]
}

export interface BookmarkNCItem {
  value: 'planned' | 'reading' | 'completed' | 'on_hold' | 'dropped'
  count: number
}

export interface NovelaDetails {
  id: number
  title: string
  alternative_titles: string[]
  description: string
  type: string
  age_rating: string
  release_date: string
  status: string
  poster_url?: string
  translation_status: string
  views: number
  authors?: NovelaAuthor[]
  genres: string[]
  categories: string[]
  volumes: NovelaVolume[]
  bookmark?: string
  has_liked: boolean
  like_count: number
  user_rating: number
  rating_details: NovelaRatingCategory
  bookmark_details: NovelaBookmarkCategory
  country: string
  last_readed?: NovelaLastReaded
}

export interface NovelaLastReaded {
  chapter_id: string
  chapter_number: number
  chapter_scroll: number
}

interface NovelaCommentUserResponse {
  id: number
  name: string
  picture: string
}

export interface NovelaCommentResponse {
  id: number
  parent_id: number
  content: string
  created_at: string
  has_like: boolean | null
  has_report: boolean | null
  like_count: number
  user: NovelaCommentUserResponse
  replies: NovelaCommentResponse[]
}
