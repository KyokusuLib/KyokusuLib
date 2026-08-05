import { defineStore, storeToRefs } from 'pinia'
import { computed, watch } from 'vue'
import { useIdle, useWindowFocus, useOnline, useIntervalFn, useEventListener } from '@vueuse/core'
import { useAuthStore } from '@/stores/auth'
import { OnlineStatusEnum } from '@/types/enums/online-status-enum'
import { $api } from '@/composables/api/useApi'
import { USER_ACTIVITY_INTERVAL, USER_ACTIVITY_IDLE_TIMEOUT } from '@/constants/data'
import { DashboardRowUserStatus } from '~/types/enums/dashboard-table'

export const useActivityStore = defineStore('activity', () => {
  const authStore = useAuthStore()
  const { isAuthenticated, user } = storeToRefs(authStore)

  const { idle, lastActive } = useIdle(USER_ACTIVITY_IDLE_TIMEOUT)
  const focused = useWindowFocus()
  const online = useOnline()

  const isUserActive = computed(() => {
    const status = online.value && focused.value && !idle.value
    if (status) user.value!.status = DashboardRowUserStatus.online
    return status
  })

  const sendHeartbeat = async () => {
    if (!isAuthenticated.value || import.meta.server) return

    try {
      await $api('/api/user/activity', {
        method: 'POST',
        body: {
          status: isUserActive.value ? OnlineStatusEnum.ONLINE : OnlineStatusEnum.OFFLINE,
          last_active: lastActive.value,
        },
      })
      if (user.value) {
        const now = new Date().toISOString()
        user.value.last_login = now
      }
    } catch {
      console.warn('Heartbeat failed')
    }
  }

  const { pause, resume } = useIntervalFn(
    () => {
      if (isUserActive.value) {
        sendHeartbeat()
      }
    },
    USER_ACTIVITY_INTERVAL,
    { immediate: false },
  )

  const initActivityTracking = () => {
    if (import.meta.server) return

    const sendOffline = () => {
      if (!isAuthenticated.value) return
      navigator.sendBeacon(
        '/api/user/activity',
        new Blob(
          [
            JSON.stringify({
              status: OnlineStatusEnum.OFFLINE,
              last_active: Date.now(),
            }),
          ],
          { type: 'application/json' },
        ),
      )
    }

    useEventListener(window, 'beforeunload', sendOffline)

    watch(
      isAuthenticated,
      (isAuth) => {
        if (isAuth) {
          sendHeartbeat()
          resume()
        } else {
          pause()
        }
      },
      { immediate: true },
    )

    watch(isUserActive, (active) => {
      if (active && isAuthenticated.value) {
        sendHeartbeat()
      }
    })
  }

  return {
    isUserActive,
    initActivityTracking,
  }
})
