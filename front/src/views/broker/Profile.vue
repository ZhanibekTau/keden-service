<template>
  <div>
    <h2 style="margin-bottom: 24px">Профиль</h2>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>Персональные данные</span>
              <el-button type="primary" size="small" @click="editing = !editing">
                {{ editing ? 'Отмена' : 'Редактировать' }}
              </el-button>
            </div>
          </template>

          <el-form v-if="editing" :model="form" label-position="top">
            <el-form-item label="Имя">
              <el-input v-model="form.first_name" />
            </el-form-item>
            <el-form-item label="Фамилия">
              <el-input v-model="form.last_name" />
            </el-form-item>
            <el-form-item label="Телефон">
              <el-input v-model="form.phone" />
            </el-form-item>
            <el-button type="primary" :loading="saving" @click="saveProfile">Сохранить</el-button>
          </el-form>

          <el-descriptions v-else :column="1" border>
            <el-descriptions-item label="Имя">{{ authStore.user?.first_name }}</el-descriptions-item>
            <el-descriptions-item label="Фамилия">{{ authStore.user?.last_name }}</el-descriptions-item>
            <el-descriptions-item label="Email">{{ authStore.user?.email }}</el-descriptions-item>
            <el-descriptions-item label="Телефон">{{ authStore.user?.phone }}</el-descriptions-item>
            <el-descriptions-item label="Тип аккаунта">
              {{ authStore.user?.account_type === 'company' ? 'Компания' : 'Физическое лицо' }}
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </el-col>

      <el-col :span="12">
        <el-card>
          <template #header><span>Смена пароля</span></template>
          <el-form :model="pwForm" label-position="top">
            <el-form-item label="Текущий пароль">
              <el-input v-model="pwForm.current_password" type="password" show-password />
            </el-form-item>
            <el-form-item label="Новый пароль">
              <el-input v-model="pwForm.new_password" type="password" show-password />
            </el-form-item>
            <el-button type="primary" :loading="changingPw" @click="changePassword">Сменить пароль</el-button>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'
import { ElMessage } from 'element-plus'

const authStore = useAuthStore()
const editing = ref(false)
const saving = ref(false)
const changingPw = ref(false)

const form = reactive({
  first_name: authStore.user?.first_name || '',
  last_name: authStore.user?.last_name || '',
  phone: authStore.user?.phone || ''
})

const pwForm = reactive({ current_password: '', new_password: '' })

async function saveProfile() {
  saving.value = true
  try {
    const response = await api.put('/company/profile', form)
    if (authStore.user) {
      Object.assign(authStore.user, response.data)
      localStorage.setItem('user', JSON.stringify(authStore.user))
    }
    editing.value = false
    ElMessage.success('Профиль обновлен')
  } catch (err: any) {
    ElMessage.error(err.response?.data?.error || 'Ошибка')
  } finally {
    saving.value = false
  }
}

async function changePassword() {
  if (!pwForm.current_password || !pwForm.new_password) {
    ElMessage.warning('Заполните все поля')
    return
  }
  changingPw.value = true
  try {
    await api.put('/company/password', pwForm)
    pwForm.current_password = ''
    pwForm.new_password = ''
    ElMessage.success('Пароль изменен')
  } catch (err: any) {
    ElMessage.error(err.response?.data?.error || 'Ошибка')
  } finally {
    changingPw.value = false
  }
}
</script>
