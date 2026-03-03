<template>
  <div>
    <h2 style="margin-bottom: 24px">Дэшборд</h2>
    <el-row :gutter="20" style="margin-bottom: 24px">
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header><span>Статус подписки</span></template>
          <el-tag :type="subscriptionTagType" size="large">{{ subscriptionStatusText }}</el-tag>
          <p v-if="subStore.current?.end_date" style="margin-top: 8px; color: #909399; font-size: 13px">
            Действует до: {{ formatDate(subStore.current.end_date) }}
          </p>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header><span>Профиль</span></template>
          <p><strong>{{ authStore.displayName }}</strong></p>
          <p style="color: #909399; font-size: 13px">{{ authStore.user?.email }}</p>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header><span>Документы</span></template>
          <p style="font-size: 28px; font-weight: bold; color: #26A65B">{{ docStore.documents.length }}</p>
          <p style="color: #909399; font-size: 13px">Всего загружено</p>
        </el-card>
      </el-col>
    </el-row>

    <el-card v-if="!subStore.current || subStore.current.status !== 'active'">
      <el-empty description="Для работы с документами необходима активная подписка">
        <router-link to="/broker/subscription">
          <el-button type="primary">Перейти к подписке</el-button>
        </router-link>
      </el-empty>
    </el-card>

    <el-card v-else>
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center">
          <span>Последние документы</span>
          <router-link to="/broker/documents"><el-button type="primary" size="small">Все документы</el-button></router-link>
        </div>
      </template>
      <el-table :data="docStore.documents.slice(0, 5)" style="width: 100%">
        <el-table-column prop="original_name" label="Файл" />
        <el-table-column prop="status" label="Статус" width="150">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="Дата" width="180">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useSubscriptionStore } from '@/stores/subscription'
import { useDocumentsStore } from '@/stores/documents'

const authStore = useAuthStore()
const subStore = useSubscriptionStore()
const docStore = useDocumentsStore()

onMounted(async () => {
  await subStore.fetchCurrent()
  if (subStore.current?.status === 'active') {
    await docStore.fetchDocuments()
  }
})

const subscriptionStatusText = computed(() => {
  if (!subStore.current) return 'Нет подписки'
  const map: Record<string, string> = { pending: 'Заявка отправлена', active: 'Активна', expired: 'Истекла', rejected: 'Отклонена' }
  return map[subStore.current.status] || subStore.current.status
})

const subscriptionTagType = computed(() => {
  if (!subStore.current) return 'info'
  const map: Record<string, string> = { pending: 'warning', active: 'success', expired: 'danger', rejected: 'danger' }
  return (map[subStore.current.status] || 'info') as any
})

function formatDate(d: string | null) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

function getStatusType(s: string) {
  const map: Record<string, string> = { uploaded: '', queued: 'warning', processing: 'warning', completed: 'success', error: 'danger' }
  return (map[s] || 'info') as any
}

function getStatusText(s: string) {
  const map: Record<string, string> = { uploaded: 'Загружен', queued: 'В очереди', processing: 'Обработка', completed: 'Готово', error: 'Ошибка' }
  return map[s] || s
}
</script>
