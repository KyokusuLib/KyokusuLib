export interface TeamRoleNames {
  owner: string
  moderator: string
  member: string
}

export interface Team {
  id: number
  owner_id: number
  name: string
  slug: string
  description: string
  avatar_url: string
  banner_url?: string
  role_names: TeamRoleNames
  created_at: string
  member_count: number
  subscribers_count: number
  is_member: boolean
  is_subscriber: boolean
  has_requested: boolean
  team_type: string
  stats: {
    member_count: number
    subscribers_count: number
  }
  team_type: string
}

export interface TeamUserMember {
  id: number
  name: string
  picture: string
  active_tag: string
  user_level: {
    level: number
    level_title: string
  }
}

export interface TeamMember {
  user: TeamUserMember
  role: string
  role_name: string
  joined_at: string
}

export interface TeamSubscriber {
  user: TeamUserMember
  subscribed_at: string
}

export interface CustomRoleNames {
  label: string
  value: string
}
