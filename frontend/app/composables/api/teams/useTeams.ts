import { ref } from 'vue'
import { useNotificationStore } from '~/stores/notification'
import { useApi, $api } from '../useApi'
import type { Team, TeamMember, TeamSubscriber } from '~/types/frontend/teams'

export function useTeam() {
  const { notify } = useNotificationStore()
  const team = ref<Team | null>(null)
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const getTeamBySlug = async (slug: string) => {
    isLoading.value = true
    error.value = null
    try {
      const data = await $api<Team>(`/api/teams/${slug}`)
      return data
    } catch (e: any) {
      error.value = e.message
      return null
    } finally {
      isLoading.value = false
    }
  }

  const fetchTeam = async (slug: string) => {
    isLoading.value = true
    error.value = null
    try {
      const { data, error: fetchError } = await useApi<Team>(`/api/teams/${slug}`)

      if (fetchError.value) throw new Error(fetchError.value.message)
      team.value = data.value!
    } catch (e: any) {
      error.value = e.message
      team.value = null
    } finally {
      isLoading.value = false
    }
  }

  const createTeam = async (form: { name: string; slug: string; description: string }) => {
    isLoading.value = true
    try {
      const data = await $api<Team>('/api/teams', {
        method: 'POST',
        body: form,
      })

      notify({ type: 'success', title: 'Успешно', content: 'Команда создана' })
      return data!
    } catch (e: any) {
      console.log(e)
      notify({ type: 'error', title: 'Ошибка', content: e.message })
      return null
    } finally {
      isLoading.value = false
    }
  }

  const updateTeam = async (slug: string, form: FormData) => {
    isLoading.value = true
    try {
      const data = await $api<Team>(`/api/teams/${slug}`, {
        method: 'PATCH',
        body: form,
      })

      if (team.value && data) {
        Object.assign(team.value, data)
      }

      notify({
        type: 'success',
        title: 'Успешно',
        content: 'Команда обновлена',
      })
      return data
    } catch (e: any) {
      notify({ type: 'error', title: 'Ошибка', content: e.message })
      return null
    } finally {
      isLoading.value = false
    }
  }

  const joinTeam = async (slug: string) => {
    isLoading.value = true
    try {
      const data = await $api<{ status: string }>(`/api/teams/${slug}/join`, {
        method: 'POST',
      })

      if (data?.status === 'request_sent') {
        notify({
          type: 'success',
          title: 'Заявка отправлена',
          content: 'Ваша заявка на вступление отправлена на рассмотрение',
        })
        return 'request_sent'
      }

      notify({
        type: 'success',
        title: 'Успешно',
        content: 'Вы вступили в команду',
      })
      return 'joined'
    } catch (e: any) {
      notify({ type: 'error', title: 'Ошибка', content: e.message })
      return false
    } finally {
      isLoading.value = false
    }
  }

  const leaveTeam = async (slug: string) => {
    isLoading.value = true
    try {
      await $api(`/api/teams/${slug}/leave`, {
        method: 'POST',
      })

      notify({
        type: 'success',
        title: 'Успешно',
        content: 'Вы покинули команду',
      })
      return true
    } catch (e: any) {
      notify({ type: 'error', title: 'Ошибка', content: e.message })
      return false
    } finally {
      isLoading.value = false
    }
  }

  const getUserTeams = async (userId: number) => {
    isLoading.value = true
    error.value = null
    try {
      const data = await $api<Team[]>(`/api/teams`, {
        query: { user_id: userId },
      })
      return data || []
    } catch (e: any) {
      error.value = e.message
      return []
    } finally {
      isLoading.value = false
    }
  }

  const getTeamMembers = async (slug: string, limit: number = 20, offset: number = 0) => {
    isLoading.value = true
    error.value = null
    try {
      const data = await $api<TeamMember[]>(`/api/teams/${slug}/members`, {
        query: { limit, offset },
      })
      return data || []
    } catch (e: any) {
      error.value = e.message
      return []
    } finally {
      isLoading.value = false
    }
  }

  const getTeamSubscribers = async (slug: string, limit: number = 20, offset: number = 0) => {
    isLoading.value = true
    error.value = null
    try {
      const data = await $api<TeamSubscriber[]>(`/api/teams/${slug}/subscribers`, {
        query: { limit, offset },
      })
      return data || []
    } catch (e: any) {
      error.value = e.message
      return []
    } finally {
      isLoading.value = false
    }
  }

  const subscribeTeam = async (teamId: number) => {
    try {
      await $api(`/api/teams/subscribe?team_id=${teamId}`, {
        method: 'POST',
      })
    } catch (e: any) {
      notify({ type: 'error', title: 'Ошибка', content: e.message })
      throw e
    }
  }

  const unsubscribeTeam = async (teamId: number) => {
    try {
      await $api(`/api/teams/unsubscribe?team_id=${teamId}`, {
        method: 'POST',
      })
    } catch (e: any) {
      notify({ type: 'error', title: 'Ошибка', content: e.message })
      throw e
    }
  }

  const getTeamJoinRequests = async (slug: string, limit = 20, offset = 0) => {
    try {
      const data = await $api<any[]>(`/api/teams/${slug}/requests?limit=${limit}&offset=${offset}`)
      return data || []
    } catch (e: any) {
      console.error(e)
      return []
    }
  }

  const acceptJoinRequest = async (slug: string, userId: number) => {
    isLoading.value = true
    try {
      await $api(`/api/teams/${slug}/requests/accept`, {
        method: 'POST',
        body: { user_id: userId },
      })
      notify({
        type: 'success',
        title: 'Успешно',
        content: 'Заявка принята',
      })
      return true
    } catch (e: any) {
      notify({ type: 'error', title: 'Ошибка', content: e.message })
      return false
    } finally {
      isLoading.value = false
    }
  }

  const rejectJoinRequest = async (slug: string, userId: number) => {
    isLoading.value = true
    try {
      await $api(`/api/teams/${slug}/requests/reject`, {
        method: 'POST',
        body: { user_id: userId },
      })
      notify({
        type: 'success',
        title: 'Успешно',
        content: 'Заявка отклонена',
      })
      return true
    } catch (e: any) {
      notify({ type: 'error', title: 'Ошибка', content: e.message })
      return false
    } finally {
      isLoading.value = false
    }
  }

  const updateTeamMember = async (
    slug: string,
    userId: number,
    data: { role?: string; custom_role_name?: string },
  ) => {
    isLoading.value = true
    try {
      await $api(`/api/teams/${slug}/members/${userId}`, {
        method: 'PATCH',
        body: data,
      })
      notify({
        type: 'success',
        title: 'Успешно',
        content: 'Данные участника обновлены',
      })
      return true
    } catch (e: any) {
      notify({ type: 'error', title: 'Ошибка', content: e.message })
      return false
    } finally {
      isLoading.value = false
    }
  }

  return {
    team,
    isLoading,
    error,
    getTeamBySlug,
    fetchTeam,
    createTeam,
    updateTeam,
    joinTeam,
    leaveTeam,
    getUserTeams,
    getTeamMembers,
    getTeamSubscribers,
    subscribeTeam,
    unsubscribeTeam,
    getTeamJoinRequests,
    acceptJoinRequest,
    rejectJoinRequest,
    updateTeamMember,
  }
}
