import type { GenderSetting } from '../enums/gender-enum'

export interface GetUserDto {
  id: number
  email?: string
  name: string
  picture: string
  banner: string
  role: string
  status: string
  about: string
  gender: GenderSetting.MALE | GenderSetting.FEMALE | GenderSetting.HIDDEN
  birthday: string
  is_public: boolean
  last_login: string
  create_at: string
  active_tag: string
  tags: UserTagDTO[]
  settings: PublicUserSettingsDTO
  user_level: UserLevel
  user_stats: UserStats
}

export interface UserStats {
  total_comments: number
  read_chapters: number
}

export interface UserTagDTO {
  id: number
  tag: string
}

export interface PublicUserSettingsDTO {
  is_show_tag: boolean
  is_show_bookmark: boolean
}

export interface UserLevel {
  level: number
  experience: number
  xp_needed_for_next: number
  level_title: string
}

export interface UserDefinitions {
  level: number
  title: string
  total_xp_required: number
}

export interface UserTagDefinitions {
  tag_id: number
  tag: string
}
