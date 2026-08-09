import { SocialNetwork } from '@/types/enums/social-network-enum'
import type { NewsSource } from '@/types/frontend/news'

export const TG_CHANNEL_NAME = 'KyokusuLib'

export const NEWS_SOURCES: NewsSource[] = [
  { id: SocialNetwork.Telegram, name: 'Telegram', icon: 'ph:telegram-logo', color: '#2A7B9B', styles: "background: #2A7B9B; background: linear-gradient(140deg,rgba(42, 123, 155, 1) 19%, rgba(125, 174, 194, 1) 82%);" },
  { id: SocialNetwork.Discord, name: 'Discord', icon: 'ph:discord-logo', color: '#7289DA', styles: "background: #7289DA; background: linear-gradient(140deg,rgba(114, 137, 218, 1) 19%, rgba(161, 186, 232, 1) 82%);" },
]
