import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/services/api'
import type { Document, AIData } from '@/types'

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
    triggerDownload(response.data, filename.replace('.pdf', '.xlsx'))
  }

  async function fetchAIData(docId: number): Promise<AIData> {
    const response = await api.get<AIData>(`/documents/${docId}/ai-data`)
    return response.data
  }

  async function updateAIData(docId: number, data: AIData): Promise<void> {
    await api.put(`/documents/${docId}/ai-data`, data)
  }

  async function downloadXML(docId: number, filename: string) {
    const response = await api.get(`/documents/${docId}/download/xml`, {
      responseType: 'blob'
    })
    triggerDownload(response.data, filename.replace('.pdf', '.xml'))
  }

  function triggerDownload(data: Blob, filename: string) {
    const url = window.URL.createObjectURL(new Blob([data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', filename)
    document.body.appendChild(link)
    link.click()
    link.remove()
    window.URL.revokeObjectURL(url)
  }

  return { documents, loading, fetchDocuments, uploadDocument, downloadExcel, fetchAIData, updateAIData, downloadXML }
})
