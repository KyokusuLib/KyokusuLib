export function syncArrays<T>(arr1: Array<T>, arr2: Array<T>): boolean {
  const set1 = new Set(arr1)
  const set2 = new Set(arr2)

  if (set1.size !== set2.size) {
    arr1.length = 0
    arr1.push(...arr2)
    return true
  }

  for (const item of set2) {
    if (!set1.has(item)) {
      arr1.length = 0
      arr1.push(...arr2)
      return true
    }
  }

  return false
}
