export interface NovelaShort {
  novelaId: number
  title: string
  description: string
  posterUrl: string
  releaseDate: Date
}

export interface TotalNovelaStatistics {
  novela: NovelaShort
  bookmarkCount: number
  readCount: number
  ratingCount: number
  commentCount: number
  volumeCount: number
  chapterCount: number
  rating: number
}

export interface TotalNovelaStatisticsResponse {
  data: TotalNovelaStatistics[]
  total: number
}

export interface GeneralStatistics {
  bookmarkCount: number
  readCount: number
  ratingCount: number
  commentCount: number
  volumeCount: number
  chapterCount: number
}

export interface MonthlyStatisticsItem {
  month: number
  year: number
  bookmarkCount: number
  readCount: number
  commentCount: number
  ratingCount: number
}

interface ShortNovelaMonthly {
  novelaId: number
  title: string
  posterUrl: string
}

export interface NovelaMonthlySeries {
  novela: ShortNovelaMonthly
  monthlyReads: number[]
}
