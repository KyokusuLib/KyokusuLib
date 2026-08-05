import { useLocalStorage } from '@vueuse/core'

export const READER_FONTS = [
  { name: 'Системный', value: 'font-sans' },
  { name: 'С засечками', value: 'font-serif' },
  { name: 'Моноширинный', value: 'font-mono' },
]

export function useReaderSettings() {
  const fontSize = useLocalStorage('reader-font-size', 18)
  const fontFamily = useLocalStorage('reader-font-family', 'font-sans')
  const lineWeight = useLocalStorage('reader-line-weight', 1.6)
  const isAutoScrollEnabled = useLocalStorage('reader-auto-scroll', false)

  const increaseFontSize = () => {
    if (fontSize.value < 40) fontSize.value += 1
  }

  const decreaseFontSize = () => {
    if (fontSize.value > 12) fontSize.value -= 1
  }

  return {
    fontSize,
    fontFamily,
    lineWeight,
    isAutoScrollEnabled,
    increaseFontSize,
    decreaseFontSize,
  }
}
