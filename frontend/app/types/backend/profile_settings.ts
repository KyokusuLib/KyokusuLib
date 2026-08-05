export interface UserInterfaceSettings {
  theme: string
  is_show_tag: boolean
  is_show_bookmark: boolean
}

export interface UserInterfaceSettingsPatch {
  theme?: string
  is_show_tag?: boolean
  is_show_bookmark?: boolean
}

export interface UserNotifySettings {
  is_app_notify: boolean
  is_new_published_notify: boolean
}
