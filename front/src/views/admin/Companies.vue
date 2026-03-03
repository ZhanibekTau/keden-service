<template>
  <div>
    <h2 style="margin-bottom: 24px">Клиенты</h2>
    <el-card>
      <el-table :data="adminStore.clients" v-loading="adminStore.loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column label="Имя" min-width="180">
          <template #default="{ row }">{{ row.first_name }} {{ row.last_name }}</template>
        </el-table-column>
        <el-table-column label="Тип" width="130">
          <template #default="{ row }">
            <el-tag :type="row.account_type === 'company' ? 'primary' : 'success'" size="small">
              {{ row.account_type === 'company' ? 'Компания' : 'Физ. лицо' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="Телефон" width="160" />
        <el-table-column prop="email" label="Email" width="220" />
        <el-table-column label="Статус" width="120">
          <template #default="{ row }">
            <el-switch
              :model-value="row.is_active"
              active-text="Акт."
              inactive-text="Выкл."
              @change="(val: boolean) => toggleStatus(row.id, val)"
            />
          </template>
        </el-table-column>
        <el-table-column label="Регистрация" width="140">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAdminStore } from '@/stores/admin'
import { ElMessage } from 'element-plus'

const adminStore = useAdminStore()

onMounted(() => adminStore.fetchClients())

async function toggleStatus(id: number, isActive: boolean) {
  try {
    await adminStore.updateClientStatus(id, isActive)
    ElMessage.success('Статус обновлен')
  } catch { ElMessage.error('Ошибка') }
}

function formatDate(d: string) {
  return new Date(d).toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit', year: 'numeric' })
}
</script>
