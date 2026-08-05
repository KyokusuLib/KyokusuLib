export type NovelaActiveTabs = 'about' | 'chapters' | 'comments'

export enum NovelaActiveTabsEnum {
  ABOUT = 'about',
  CHAPTERS = 'chapters',
  COMMENTS = 'comments',
}

export const NOVELA_ACTIVE_TABS: NovelaActiveTabs[] = ['about', 'chapters', 'comments']

export function convToRu(tab: NovelaActiveTabs): string {
  switch (tab) {
    case 'about':
      return 'Описание'
    case 'chapters':
      return 'Главы'
    case 'comments':
      return 'Комментарии'
    default:
      return ''
  }
}
