<template>
  <div>
    <h2 style="margin-bottom: 24px">Подписка</h2>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header><span>Текущая подписка</span></template>

          <div v-if="!subStore.current || subStore.current.status === 'expired' || subStore.current.status === 'rejected'">
            <el-empty description="Подписка отсутствует" :image-size="80">
              <el-button type="primary" :loading="requesting" @click="requestSub">
                Подать заявку на подписку
              </el-button>
            </el-empty>
          </div>

          <div v-else-if="subStore.current.status === 'pending'">
            <el-result icon="info" title="Заявка отправлена" sub-title="Ожидает рассмотрения администратором" />
          </div>

          <div v-else-if="subStore.current.status === 'in_progress'">
            <el-result icon="info" title="Принято в работу" sub-title="Менеджер обрабатывает вашу заявку" />
          </div>

          <div v-else-if="subStore.current.status === 'invoice_sent'">
            <el-result icon="warning" title="Счёт отправлен">
              <template #sub-title>
                <p>Счёт на оплату отправлен. После получения оплаты подписка будет активирована.</p>
                <p v-if="subStore.current.admin_comment" style="margin-top: 8px; color: #606266">
                  Комментарий: {{ subStore.current.admin_comment }}
                </p>
              </template>
            </el-result>
          </div>

          <div v-else-if="subStore.current.status === 'active'">
            <el-result icon="success" title="Подписка активна">
              <template #sub-title>
                <p>Действует до: <strong>{{ formatDate(subStore.current.end_date) }}</strong></p>
              </template>
            </el-result>
          </div>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header><span>Информация</span></template>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="Стоимость">12 990 &#8376; / месяц</el-descriptions-item>
            <el-descriptions-item label="Срок">1 месяц</el-descriptions-item>
            <el-descriptions-item label="Включено">Неограниченная загрузка PDF, AI-обработка, генерация Excel</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>
    </el-row>

    <el-card style="margin-top: 24px">
      <template #header><span>История заявок</span></template>
      <el-table :data="subStore.history" style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="Статус" width="150">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Дата заявки" width="180">
          <template #default="{ row }">{{ formatDate(row.requested_at) }}</template>
        </el-table-column>
        <el-table-column label="Начало" width="150">
          <template #default="{ row }">{{ formatDate(row.start_date) }}</template>
        </el-table-column>
        <el-table-column label="Окончание" width="150">
          <template #default="{ row }">{{ formatDate(row.end_date) }}</template>
        </el-table-column>
        <el-table-column prop="admin_comment" label="Комментарий" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useSubscriptionStore } from '@/stores/subscription'
import { ElMessage } from 'element-plus'

const subStore = useSubscriptionStore()
const requesting = ref(false)

onMounted(async () => {
  try {
    await subStore.fetchCurrent()
    await subStore.fetchHistory()
  } catch {
    ElMessage.error('Ошибка загрузки данных подписки')
  }
})

async function requestSub() {
  requesting.value = true
  try {
    await subStore.requestSubscription()
    await subStore.fetchHistory()
    ElMessage.success('Заявка на подписку отправлена')
  } catch (err: any) {
    const serverError = err.response?.data?.error || ''
    const errorMap: Record<string, string> = {
      'you already have a pending subscription request': 'У вас уже есть активная заявка на рассмотрении',
      'you already have an active subscription': 'У вас уже есть активная подписка',
    }
    ElMessage.error(errorMap[serverError] || serverError || 'Ошибка при создании заявки')
    // Обновляем состояние на случай рассинхронизации
    await subStore.fetchCurrent().catch(() => {})
    await subStore.fetchHistory().catch(() => {})
  } finally {
    requesting.value = false
  }
}

function formatDate(d: string | null) {
  if (!d) return '-'
  return new Date(d).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric' })
}

function getStatusType(s: string) {
  const map: Record<string, any> = {
    pending: 'warning', in_progress: 'primary', invoice_sent: '',
    active: 'success', expired: 'info', rejected: 'danger'
  }
  return map[s] || 'info'
}

function getStatusText(s: string) {
  const map: Record<string, string> = {
    pending: 'Ожидает', in_progress: 'В работе', invoice_sent: 'Счёт отправлен',
    active: 'Активна', expired: 'Истекла', rejected: 'Отклонена'
  }
  return map[s] || s
}
</script>
