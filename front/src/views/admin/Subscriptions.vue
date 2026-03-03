<template>
  <div>
    <h2 style="margin-bottom: 24px">Управление подписками</h2>

    <el-card style="margin-bottom: 24px">
      <template #header><span>Ожидающие подтверждения</span></template>
      <el-table :data="adminStore.pendingSubscriptions" style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="Клиент" min-width="180">
          <template #default="{ row }">
            <div>
              <strong>{{ row.user?.first_name }} {{ row.user?.last_name }}</strong>
              <p style="margin: 0; font-size: 12px; color: #909399">{{ row.user?.email }}</p>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="Тип" width="140">
          <template #default="{ row }">
            <el-tag :type="row.user?.account_type === 'company' ? 'primary' : 'success'" size="small">
              {{ row.user?.account_type === 'company' ? 'Компания' : 'Физ. лицо' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="Телефон" width="160">
          <template #default="{ row }">{{ row.user?.phone }}</template>
        </el-table-column>
        <el-table-column label="Сумма" width="120">
          <template #default="{ row }">{{ row.amount?.toLocaleString() }} &#8376;</template>
        </el-table-column>
        <el-table-column label="Дата заявки" width="140">
          <template #default="{ row }">{{ formatDate(row.requested_at) }}</template>
        </el-table-column>
        <el-table-column label="Действия" width="240" fixed="right">
          <template #default="{ row }">
            <el-button type="success" size="small" @click="handleApprove(row.id)">Одобрить</el-button>
            <el-button type="danger" size="small" @click="handleReject(row.id)">Отклонить</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="!adminStore.pendingSubscriptions.length" description="Нет ожидающих заявок" :image-size="60" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'
import { ElMessage, ElMessageBox } from 'element-plus'

const adminStore = useAdminStore()

onMounted(() => adminStore.fetchPendingSubscriptions())

async function handleApprove(id: number) {
  try {
    const result = await ElMessageBox.prompt('Комментарий (способ оплаты):', 'Одобрение подписки', {
      confirmButtonText: 'Одобрить',
      cancelButtonText: 'Отмена',
      inputPlaceholder: 'Например: Kaspi перевод'
    })
    const comment = typeof result === 'object' && 'value' in result ? (result as any).value : ''
    await adminStore.approveSubscription(id, comment || '')
    ElMessage.success('Подписка одобрена')
  } catch {
    // cancelled
  }
}

async function handleReject(id: number) {
  try {
    const result = await ElMessageBox.prompt('Причина отклонения:', 'Отклонение заявки', {
      confirmButtonText: 'Отклонить',
      cancelButtonText: 'Отмена',
      type: 'warning'
    })
    const comment = typeof result === 'object' && 'value' in result ? (result as any).value : ''
    await adminStore.rejectSubscription(id, comment || '')
    ElMessage.success('Заявка отклонена')
  } catch {
    // cancelled
  }
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric' })
}
</script>
