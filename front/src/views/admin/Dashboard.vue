<template>
  <div>
    <h2 style="margin-bottom: 24px">Дэшборд администратора</h2>
    <el-row :gutter="20" style="margin-bottom: 24px">
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="Компании" :value="adminStore.stats?.total_companies || 0">
            <template #prefix><el-icon><OfficeBuilding /></el-icon></template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="Активные подписки" :value="adminStore.stats?.active_subscriptions || 0">
            <template #prefix><el-icon style="color: #67C23A"><CircleCheck /></el-icon></template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="Ожидают подтверждения" :value="adminStore.stats?.pending_subscriptions || 0">
            <template #prefix><el-icon style="color: #E6A23C"><Clock /></el-icon></template>
          </el-statistic>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover">
          <el-statistic title="Обработано документов" :value="adminStore.stats?.completed_documents || 0">
            <template #suffix>
              <span style="font-size: 14px; color: #909399"> / {{ adminStore.stats?.total_documents || 0 }}</span>
            </template>
            <template #prefix><el-icon style="color: #26A65B"><Document /></el-icon></template>
          </el-statistic>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>Ожидающие заявки</span>
              <router-link to="/admin/subscriptions"><el-button size="small" text type="primary">Все</el-button></router-link>
            </div>
          </template>
          <el-table :data="adminStore.pendingSubscriptions.slice(0, 5)" style="width: 100%">
            <el-table-column label="Клиент">
              <template #default="{ row }">{{ row.user ? `${row.user.first_name} ${row.user.last_name}` : '-' }}</template>
            </el-table-column>
            <el-table-column label="Статус" width="160">
              <template #default="{ row }">
                <el-tag :type="statusType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="Дата" width="110">
              <template #default="{ row }">{{ formatDate(row.requested_at) }}</template>
            </el-table-column>
          </el-table>
          <el-empty v-if="!adminStore.pendingSubscriptions.length" description="Нет ожидающих заявок" :image-size="60" />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header><span>Быстрые ссылки</span></template>
          <el-space direction="vertical" :size="12" style="width: 100%">
            <router-link to="/admin/companies" style="display: block"><el-button style="width: 100%" text type="primary">Управление компаниями</el-button></router-link>
            <router-link to="/admin/subscriptions" style="display: block"><el-button style="width: 100%" text type="primary">Управление подписками</el-button></router-link>
            <router-link to="/admin/documents" style="display: block"><el-button style="width: 100%" text type="primary">Все документы</el-button></router-link>
          </el-space>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'

const adminStore = useAdminStore()

onMounted(async () => {
  await Promise.all([adminStore.fetchStats(), adminStore.fetchPendingSubscriptions()])
})

function statusLabel(s: string) {
  const map: Record<string, string> = {
    pending: 'Новая', in_progress: 'В работе', invoice_sent: 'Счёт отправлен'
  }
  return map[s] ?? s
}

function statusType(s: string) {
  const map: Record<string, string> = {
    pending: 'warning', in_progress: 'primary', invoice_sent: ''
  }
  return map[s] ?? 'info'
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric' })
}
</script>
