<template>
  <div class="auth-layout">
    <header class="auth-header">
      <div class="auth-header-content">
        <router-link to="/" class="logo-group">
          <span class="logo">KEDEN</span>
          <span class="logo-divider"></span>
          <span class="logo-sub">AI-powered</span>
        </router-link>
        <router-link to="/" class="back-link">
          <el-icon><ArrowLeft /></el-icon>
          На главную
        </router-link>
      </div>
    </header>

    <div class="auth-page">
      <el-card class="auth-card" shadow="always">
        <h2 style="text-align: center; margin-bottom: 24px; color: #1a1a2e">
          Вход в систему
        </h2>
        <el-form ref="formRef" :model="form" :rules="rules" label-position="top" @submit.prevent="handleLogin">
          <el-form-item label="Email" prop="email">
            <el-input v-model="form.email" placeholder="example@company.kz" size="large" />
          </el-form-item>
          <el-form-item label="Пароль" prop="password">
            <el-input v-model="form.password" type="password" placeholder="Введите пароль" size="large" show-password />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="large" style="width: 100%" :loading="loading" @click="handleLogin">
              Войти
            </el-button>
          </el-form-item>
        </el-form>
        <div style="text-align: center">
          <span style="color: #909399">Нет аккаунта?</span>
          <el-button text type="primary" style="margin-left: 4px; padding: 0" @click="showRegisterModal = true">
            Зарегистрироваться
          </el-button>
        </div>
      </el-card>
    </div>

    <!-- Модальное окно выбора типа регистрации -->
    <el-dialog
      v-model="showRegisterModal"
      title="Выберите тип регистрации"
      width="500px"
      :close-on-click-modal="true"
      align-center
    >
      <div class="register-type-options">
        <div class="register-type-card" @click="goToRegister('individual')">
          <el-icon :size="48" color="#26A65B"><User /></el-icon>
          <h3>Физическое лицо</h3>
          <p>Регистрация как индивидуальный брокер без привязки к компании</p>
        </div>
        <div class="register-type-card" @click="goToRegister('company')">
          <el-icon :size="48" color="#26A65B"><OfficeBuilding /></el-icon>
          <h3>Компания</h3>
          <p>Регистрация как юридическое лицо с данными компании</p>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage, type FormInstance } from 'element-plus'
import { ArrowLeft, User, OfficeBuilding } from '@element-plus/icons-vue'

const router = useRouter()
const authStore = useAuthStore()
const formRef = ref<FormInstance>()
const loading = ref(false)
const showRegisterModal = ref(false)

const form = reactive({ email: '', password: '' })

const rules = {
  email: [{ required: true, message: 'Введите email', trigger: 'blur' }, { type: 'email' as const, message: 'Неверный формат email', trigger: 'blur' }],
  password: [{ required: true, message: 'Введите пароль', trigger: 'blur' }]
}

async function handleLogin() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    await authStore.login(form)
    const dest = authStore.isAdmin ? '/admin/dashboard' : '/broker/dashboard'
    await router.push(dest).catch(() => { window.location.href = dest })
    ElMessage.success('Вы успешно вошли в систему')
  } catch (err: any) {
    const serverError = err.response?.data?.error || ''
    const errorMap: Record<string, string> = {
      'invalid email or password': 'Неверный email или пароль',
      'account is inactive': 'Аккаунт неактивен',
    }
    ElMessage.error(errorMap[serverError] || serverError || 'Ошибка входа')
  } finally {
    loading.value = false
  }
}

function goToRegister(type: string) {
  showRegisterModal.value = false
  router.push({ path: '/register', query: { type } })
}
</script>

<style scoped>
.auth-layout { min-height: 100vh; display: flex; flex-direction: column; background: linear-gradient(135deg, #f0faf4 0%, #fff 60%); }

.auth-header {
  background: #fff;
  border-bottom: 1px solid #ebeef5;
  padding: 14px 0;
  position: sticky;
  top: 0;
  z-index: 10;
}
.auth-header-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.logo-group { display: flex; align-items: center; gap: 12px; text-decoration: none; }
.logo { font-size: 24px; color: #26A65B; font-weight: 800; letter-spacing: 2px; }
.logo-divider { width: 1px; height: 20px; background: #dcdfe6; }
.logo-sub { color: #909399; font-size: 13px; }
.back-link {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #606266;
  text-decoration: none;
  font-size: 14px;
  transition: color 0.2s;
}
.back-link:hover { color: #26A65B; }

.auth-page { flex: 1; display: flex; align-items: center; justify-content: center; padding: 40px 20px; }
.auth-card { width: 420px; padding: 20px; border-radius: 12px; }

.register-type-options {
  display: flex;
  gap: 20px;
  padding: 8px 0 16px;
}
.register-type-card {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 28px 16px;
  border: 2px solid #ebeef5;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s;
  text-align: center;
}
.register-type-card:hover {
  border-color: #26A65B;
  background: rgba(38, 166, 91, 0.05);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(38, 166, 91, 0.15);
}
.register-type-card h3 { margin: 12px 0 8px; color: #1a1a2e; font-size: 16px; }
.register-type-card p { margin: 0; color: #909399; font-size: 13px; line-height: 1.5; }
</style>
