import type { SocialNetwork } from '@/types/enums/social-network-enum'

export interface NewsSource {
  id: SocialNetwork
  name: string
  icon: string
  styles?: string
  color?: string
}
