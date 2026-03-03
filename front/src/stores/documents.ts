import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/services/api'
import type { Document } from '@/types'

export const useDocumentsStore = defineStore('documents', () => {
  const documents = ref<Document[]>([])
  const loading = ref(false)

  async function fetchDocuments() {
    loading.value = true
    try {
      const response = await api.get<Document[]>('/documents')
      documents.value = response.data
    } finally {
      loading.value = false
    }
  }

  async function uploadDocument(file: File) {
    const formData = new FormData()
    formData.append('file', file)
    const response = await api.post<Document>('/documents/upload', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    documents.value.unshift(response.data)
    return response.data
  }

  async function downloadExcel(docId: number, filename: string) {
    const response = await api.get(`/documents/${docId}/download`, {
      responseType: 'blob'
    })
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', filename.replace('.pdf', '.xlsx'))
    document.body.appendChild(link)
    link.click()
    link.remove()
    window.URL.revokeObjectURL(url)
  }

  return { documents, loading, fetchDocuments, uploadDocument, downloadExcel }
})
