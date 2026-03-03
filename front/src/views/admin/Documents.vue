<template>
  <div>
    <h2 style="margin-bottom: 24px">Все документы</h2>
    <el-card>
      <el-table :data="adminStore.documents" style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="Клиент" min-width="180">
          <template #default="{ row }">{{ row.user ? `${row.user.first_name} ${row.user.last_name}` : '-' }}</template>
        </el-table-column>
        <el-table-column prop="original_name" label="Файл" min-width="200" />
        <el-table-column label="Размер" width="100">
          <template #default="{ row }">{{ formatSize(row.file_size) }}</template>
        </el-table-column>
        <el-table-column label="Статус" width="140">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Загружен" width="160">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="Обработан" width="160">
          <template #default="{ row }">{{ formatDate(row.processed_at) }}</template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'

const adminStore = useAdminStore()

onMounted(() => adminStore.fetchDocuments())

function formatDate(d: string | null) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function formatSize(bytes: number) {
  if (!bytes) return '-'
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
