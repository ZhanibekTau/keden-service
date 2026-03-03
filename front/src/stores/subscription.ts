import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/services/api'
import type { Subscription } from '@/types'

export const useSubscriptionStore = defineStore('subscription', () => {
  const current = ref<Subscription | null>(null)
  const history = ref<Subscription[]>([])
  const loading = ref(false)

  async function fetchCurrent() {
    loading.value = true
    try {
      const response = await api.get('/subscription/current')
      if (response.data.status === 'none') {
        current.value = null
      } else {
        current.value = response.data
      }
    } finally {
      loading.value = false
    }
  }

  async function fetchHistory() {
    const response = await api.get<Subscription[]>('/subscription/history')
    history.value = response.data
  }

  async function requestSubscription() {
    const response = await api.post<Subscription>('/subscription/request')
    current.value = response.data
    return response.data
  }

  return { current, history, loading, fetchCurrent, fetchHistory, requestSubscription }
})
