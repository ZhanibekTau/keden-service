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
    const response = await api.get<Subscription[]>('/admin/subscriptions/pending')
    pendingSubscriptions.value = response.data
  }

  async function approveSubscription(id: number, comment: string) {
    await api.post(`/admin/subscriptions/${id}/approve`, { comment })
    pendingSubscriptions.value = pendingSubscriptions.value.filter(s => s.id !== id)
  }

  async function rejectSubscription(id: number, comment: string) {
    await api.post(`/admin/subscriptions/${id}/reject`, { comment })
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
    approveSubscription, rejectSubscription, updateClientStatus, fetchDocuments
  }
})
