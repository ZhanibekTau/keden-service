<template>
  <div>
    <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px">
      <h2 style="margin: 0">Управление подписками</h2>
      <el-button :loading="loading" @click="reload" size="small" text type="primary">Обновить</el-button>
    </div>

    <el-card>
      <el-tabs v-model="activeTab" type="border-card">
        <el-tab-pane v-for="tab in tabs" :key="tab.status" :name="tab.status">
          <template #label>
            <span style="display: flex; align-items: center; gap: 6px">
              {{ tab.label }}
              <el-badge
                :value="countByStatus(tab.status)"
                :hidden="countByStatus(tab.status) === 0"
                :type="tab.badgeType"
              />
            </span>
          </template>

          <el-table
            :data="byStatus(tab.status)"
            v-loading="loading"
            style="width: 100%"
            :empty-text="`Нет заявок со статусом «${tab.label}»`"
          >
            <el-table-column prop="id" label="ID" width="60" />

            <!-- Client -->
            <el-table-column label="Клиент" min-width="200">
              <template #default="{ row }">
                <div style="font-weight: 600">{{ row.user?.first_name }} {{ row.user?.last_name }}</div>
                <div style="font-size: 12px; color: #909399">{{ row.user?.email }}</div>
                <div style="font-size: 12px; color: #909399">{{ row.user?.phone }}</div>
              </template>
            </el-table-column>

            <!-- Account type + company -->
            <el-table-column label="Тип / Компания" min-width="230">
              <template #default="{ row }">
                <el-tag
                  :type="row.user?.account_type === 'company' ? 'primary' : 'success'"
                  size="small"
                  style="margin-bottom: 4px"
                >
                  {{ row.user?.account_type === 'company' ? 'Компания' : 'Физ. лицо' }}
                </el-tag>
                <template v-if="row.user?.account_type === 'company'">
                  <div style="font-size: 13px; margin-top: 4px">
                    <div><strong>{{ row.company_name || '—' }}</strong></div>
                    <div style="color: #606266; font-size: 12px">{{ row.legal_name }}</div>
                    <div style="color: #409EFF; font-weight: 600; font-size: 13px">БИН: {{ row.bin || '—' }}</div>
                  </div>
                </template>
              </template>
            </el-table-column>

            <!-- Amount -->
            <el-table-column label="Сумма" width="130">
              <template #default="{ row }">{{ row.amount?.toLocaleString('ru-RU') }} ₸</template>
            </el-table-column>

            <!-- Date -->
            <el-table-column label="Дата заявки" width="120">
              <template #default="{ row }">{{ formatDate(row.requested_at) }}</template>
            </el-table-column>

            <!-- Comment -->
            <el-table-column label="Комментарий" min-width="150">
              <template #default="{ row }">
                <span style="font-size: 12px; color: #606266">{{ row.admin_comment || '—' }}</span>
              </template>
            </el-table-column>

            <!-- Actions -->
            <el-table-column label="Действия" width="200" fixed="right">
              <template #default="{ row }">
                <div style="display: flex; flex-direction: column; gap: 6px">
                  <el-button
                    v-if="nextAction(row.status)"
                    :type="nextAction(row.status)!.btnType"
                    size="small"
                    style="width: 100%; margin: 0"
                    @click="doAction(row, nextAction(row.status)!.status)"
                  >
                    {{ nextAction(row.status)!.label }}
                  </el-button>
                  <el-button
                    type="danger"
                    size="small"
                    plain
                    style="width: 100%; margin: 0"
                    @click="doReject(row)"
                  >
                    Отклонить
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Subscription } from '@/types'

const adminStore = useAdminStore()
const loading = ref(false)
const activeTab = ref('pending')

const tabs = [
  { status: 'pending',      label: 'Новые заявки',    badgeType: 'warning' as const },
  { status: 'in_progress',  label: 'В работе',         badgeType: 'primary' as const },
  { status: 'invoice_sent', label: 'Счёт отправлен',   badgeType: 'danger'  as const },
]

onMounted(reload)

async function reload() {
  loading.value = true
  try {
    await adminStore.fetchPendingSubscriptions()
  } finally {
    loading.value = false
  }
}

function byStatus(status: string) {
  return adminStore.pendingSubscriptions.filter(s => s.status === status)
}

function countByStatus(status: string) {
  return adminStore.pendingSubscriptions.filter(s => s.status === status).length
}

function nextAction(status: string): { label: string; status: string; btnType: 'success' | 'primary' | 'warning' } | null {
  const map: Record<string, { label: string; status: string; btnType: 'success' | 'primary' | 'warning' }> = {
    pending:      { label: 'Принять в работу', status: 'in_progress',  btnType: 'primary' },
    in_progress:  { label: 'Отправить счёт',   status: 'invoice_sent', btnType: 'warning' },
    invoice_sent: { label: 'Активировать',     status: 'active',       btnType: 'success' },
  }
  return map[status] ?? null
}

async function doAction(row: Subscription, newStatus: string) {
  const labels: Record<string, string> = {
    in_progress:  'Принять заявку в работу?',
    invoice_sent: 'Подтвердить отправку счёта?',
    active:       'Активировать подписку?',
  }
  try {
    const result = await ElMessageBox.prompt(
      labels[newStatus] ?? 'Изменить статус?',
      'Подтверждение',
      { confirmButtonText: 'Подтвердить', cancelButtonText: 'Отмена', inputPlaceholder: 'Комментарий (необязательно)', inputValue: '' }
    )
    const comment = typeof result === 'object' && 'value' in result ? (result as any).value : ''
    await adminStore.updateSubscriptionStatus(row.id, newStatus, comment || '')
    await adminStore.fetchStats()
    ElMessage.success('Статус обновлён')
    // Auto-switch to the next tab
    const nextTab: Record<string, string> = { in_progress: 'in_progress', invoice_sent: 'invoice_sent', active: 'pending' }
    if (nextTab[newStatus]) activeTab.value = nextTab[newStatus]
  } catch { /* cancelled */ }
}

async function doReject(row: Subscription) {
  try {
    const result = await ElMessageBox.prompt(
      'Причина отклонения:',
      'Отклонение заявки',
      { confirmButtonText: 'Отклонить', cancelButtonText: 'Отмена', type: 'warning', inputPlaceholder: 'Укажите причину' }
    )
    const comment = typeof result === 'object' && 'value' in result ? (result as any).value : ''
    await adminStore.updateSubscriptionStatus(row.id, 'rejected', comment || '')
    await adminStore.fetchStats()
    ElMessage.success('Заявка отклонена')
  } catch { /* cancelled */ }
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric' })
}
</script>
