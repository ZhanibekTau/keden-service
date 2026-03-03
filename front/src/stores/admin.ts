import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/services/api'
import type { AdminStats, User, Subscription, Document } from '@/types'

export const useAdminStore = defineStore('admin', () => {
  const stats = ref<AdminStats | null>(null)
  const clients = ref<User[]>([])
  const pendingSubscriptions = ref<Subscription[]>([])
  const documents = ref<Document[]>([])
  const loading = ref(false)

  async function fetchStats() {
    const response = await api.get<AdminStats>('/admin/stats')
    stats.value = response.data
  }

  async function fetchClients() {
    loading.value = true
    try {
      const response = await api.get<User[]>('/admin/companies')
      clients.value = response.data
    } finally {
      loading.value = false
    }
  }

  async function fetchPendingSubscriptions() {
    const response = await api.get<Subscription[]>('/admin/subscriptions')
    pendingSubscriptions.value = response.data
  }

  async function updateSubscriptionStatus(id: number, status: string, comment = '') {
    await api.put(`/admin/subscriptions/${id}/status`, { status, comment })
    pendingSubscriptions.value = pendingSubscriptions.value.filter(s => s.id !== id)
  }

  async function updateClientStatus(id: number, isActive: boolean) {
    await api.put(`/admin/companies/${id}/status`, { is_active: isActive })
    const idx = clients.value.findIndex(c => c.id === id)
    if (idx !== -1) clients.value[idx].is_active = isActive
  }

  async function fetchDocuments() {
    const response = await api.get<Document[]>('/admin/documents')
    documents.value = response.data
  }

  return {
    stats, clients, pendingSubscriptions, documents, loading,
    fetchStats, fetchClients, fetchPendingSubscriptions,
    updateSubscriptionStatus, updateClientStatus, fetchDocuments
  }
})
