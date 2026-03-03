<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px">
      <h2 style="margin: 0">Документы</h2>
      <el-button type="primary" @click="showUpload = true">
        <el-icon style="margin-right: 4px"><Upload /></el-icon> Загрузить PDF
      </el-button>
    </div>

    <!-- Upload dialog -->
    <el-dialog v-model="showUpload" title="Загрузка документа" width="500px">
      <el-upload
        ref="uploadRef"
        drag
        :auto-upload="false"
        :limit="1"
        accept=".pdf"
        :on-change="handleFileChange"
      >
        <el-icon :size="48"><Upload /></el-icon>
        <div class="el-upload__text">Перетащите PDF-файл сюда или <em>нажмите для выбора</em></div>
        <template #tip>
          <div class="el-upload__tip">Только PDF файлы, максимум 50 МБ</div>
        </template>
      </el-upload>
      <template #footer>
        <el-button @click="showUpload = false">Отмена</el-button>
        <el-button type="primary" :loading="uploading" :disabled="!selectedFile" @click="handleUpload">
          Загрузить
        </el-button>
      </template>
    </el-dialog>

    <!-- Documents table -->
    <el-card>
      <el-table :data="docStore.documents" v-loading="docStore.loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="original_name" label="Файл" min-width="200" />
        <el-table-column label="Размер" width="120">
          <template #default="{ row }">{{ formatSize(row.file_size) }}</template>
        </el-table-column>
        <el-table-column label="Статус" width="150">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Дата" width="180">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="Действия" width="150" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'completed'"
              type="success"
              size="small"
              @click="downloadExcel(row)"
            >
              <el-icon><Download /></el-icon> Excel
            </el-button>
            <el-tooltip v-if="row.status === 'error'" :content="row.error_message" placement="top">
              <el-button type="danger" size="small" text>
                <el-icon><Warning /></el-icon>
              </el-button>
            </el-tooltip>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useDocumentsStore } from '@/stores/documents'
import { ElMessage } from 'element-plus'
import { Upload, Download, Warning } from '@element-plus/icons-vue'
import type { Document } from '@/types'
import type { UploadFile } from 'element-plus'

const docStore = useDocumentsStore()
const showUpload = ref(false)
const uploading = ref(false)
const selectedFile = ref<File | null>(null)

onMounted(() => docStore.fetchDocuments())

function handleFileChange(file: UploadFile) {
  selectedFile.value = file.raw || null
}

async function handleUpload() {
  if (!selectedFile.value) return
  uploading.value = true
  try {
    await docStore.uploadDocument(selectedFile.value)
    ElMessage.success('Документ загружен успешно')
    showUpload.value = false
    selectedFile.value = null
  } catch (err: any) {
    ElMessage.error(err.response?.data?.error || 'Ошибка загрузки')
  } finally {
    uploading.value = false
  }
}

async function downloadExcel(doc: Document) {
  try {
    await docStore.downloadExcel(doc.id, doc.original_name)
    ElMessage.success('Файл скачивается')
  } catch {
    ElMessage.error('Ошибка скачивания')
  }
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function formatSize(bytes: number) {
  if (bytes < 1024) return bytes + ' Б'
  if (bytes < 1048576) return (bytes / 1024).toFixed(1) + ' КБ'
  return (bytes / 1048576).toFixed(1) + ' МБ'
}

function getStatusType(s: string) {
  const map: Record<string, any> = { uploaded: '', queued: 'warning', processing: 'warning', completed: 'success', error: 'danger' }
  return map[s] || 'info'
}

function getStatusText(s: string) {
  const map: Record<string, string> = { uploaded: 'Загружен', queued: 'В очереди', processing: 'Обработка', completed: 'Готово', error: 'Ошибка' }
  return map[s] || s
}
</script>
