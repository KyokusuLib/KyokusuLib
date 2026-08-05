import type { UseFetchOptions } from 'nuxt/app'
import { defu } from 'defu'

export function useApi<T>(url: string | (() => string), options: UseFetchOptions<T> = {}) {
  const defaults: UseFetchOptions<T> = {
    baseURL: '/api',
    headers: useRequestHeaders(['cookie']),
  }

  return useFetch(url, defu(options, defaults))
}

export async function $api<T>(url: string, options: any = {}): Promise<T> {
  const defaults = {
    baseURL: '/api',
    headers: useRequestHeaders(['cookie']),
  }

  const params = defu(options, defaults)

  try {
    return await $fetch<T>(url, params)
  } catch (err: any) {
    const fetchError = err.data

    const errorMessage =
      fetchError?.error ||
      fetchError?.message ||
      err.statusMessage ||
      'Произошла непредвиденная ошибка'

    const error = new Error(errorMessage)

    ;(error as any).statusCode = err.response?.status

    throw error
  }
}
