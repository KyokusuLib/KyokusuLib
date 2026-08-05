export interface DashboardUser {
  id: number
  email: string
  name: string
  picture: string
  banner: string
  role: string
  status: string
  about: string
  birthday: string | null
  gender: string
  is_public: boolean
  last_login: string
  create_at: string
  active_tag: string
  tags: DashboardUserTag[]
  settings: DashboardUserSettings
  user_level: DashboardUserLevel | null
  user_stats: DashboardUserStats | null
}

export interface DashboardUserLevel {
  level: number
  experience: number
  level_title: string
  xp_needed_for_next: number
}

export interface DashboardUserTag {
  tag_id: number
  tag: string
}

export interface DashboardUserSettings {
  is_show_tag: boolean
  is_show_bookmark: boolean
}

export interface DashboardUserStats {
  total_comments: number
  read_chapters: number
}
