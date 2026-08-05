export enum KyokusuAppRole {
  ADMIN = 'admin',
  MODERATOR = 'moderator',
  PUBLISHER = 'publisher',
  USER = 'user',
}

export const ROLE_WEIGHTS: Record<KyokusuAppRole, number> = {
  [KyokusuAppRole.USER]: 10,
  [KyokusuAppRole.PUBLISHER]: 30,
  [KyokusuAppRole.MODERATOR]: 50,
  [KyokusuAppRole.ADMIN]: 100,
}
