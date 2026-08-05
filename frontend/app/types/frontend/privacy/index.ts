export interface PrivacyCardItem {
  title: string
  description?: string
  text?: string

  subcards?: PrivacySubCard[]
}

export interface PrivacySubCard {
  title: string
  text: string
}
