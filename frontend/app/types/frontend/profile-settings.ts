export interface UserProfileSettings {
  id: number
  email: string
  name: string
  picture: string
  role: string
  status: string
  last_login: string

  about: string
  gender: 'male' | 'female' | 'hidden'
  birthdate: string
}
