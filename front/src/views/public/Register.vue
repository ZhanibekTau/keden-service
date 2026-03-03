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
        Регистрация
      </h2>

      <div v-if="accountType" style="text-align: center; margin-bottom: 16px">
        <el-tag :type="accountType === 'company' ? 'primary' : 'success'" size="large">
          {{ accountType === 'company' ? 'Компания' : 'Физическое лицо' }}
        </el-tag>
      </div>

      <el-form ref="formRef" :model="form" :rules="computedRules" label-position="top" @submit.prevent="handleRegister">
        <el-form-item label="Email" prop="email">
          <el-input v-model="form.email" placeholder="example@company.kz" size="large" />
        </el-form-item>
        <el-form-item label="Пароль" prop="password">
          <el-input v-model="form.password" type="password" placeholder="Минимум 6 символов" size="large" show-password />
        </el-form-item>
        <el-form-item label="Подтверждение пароля" prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" placeholder="Повторите пароль" size="large" show-password />
        </el-form-item>

        <el-divider>{{ accountType === 'company' ? 'Данные контактного лица' : 'Персональные данные' }}</el-divider>

        <el-form-item label="Имя" prop="first_name">
          <el-input v-model="form.first_name" placeholder="Имя" size="large" />
        </el-form-item>
        <el-form-item label="Фамилия" prop="last_name">
          <el-input v-model="form.last_name" placeholder="Фамилия" size="large" />
        </el-form-item>
        <el-form-item label="Телефон" prop="phone">
          <el-input v-model="form.phone" placeholder="+7 (XXX) XXX-XX-XX" size="large" />
        </el-form-item>

        <template v-if="accountType === 'company'">
          <el-divider>Данные компании</el-divider>
          <el-form-item label="Название компании" prop="company_name">
            <el-input v-model="form.company_name" placeholder="ТОО «Ваша Компания»" size="large" />
          </el-form-item>
          <el-form-item label="Юридическое название" prop="legal_name">
            <el-input v-model="form.legal_name" placeholder="Полное юридическое название" size="large" />
          </el-form-item>
          <el-form-item label="БИН" prop="bin">
            <el-input v-model="form.bin" placeholder="12 цифр" maxlength="12" size="large" />
          </el-form-item>
          <el-form-item label="Контактное лицо" prop="contact_person">
            <el-input v-model="form.contact_person" placeholder="ФИО контактного лица" size="large" />
          </el-form-item>
        </template>

        <el-form-item>
          <el-button type="primary" size="large" style="width: 100%" :loading="loading" @click="handleRegister">
            Зарегистрироваться
          </el-button>
        </el-form-item>
      </el-form>
      <div style="text-align: center">
        <span style="color: #909399">Уже есть аккаунт?</span>
        <router-link to="/login" style="margin-left: 8px">Войти</router-link>
      </div>
    </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage, type FormInstance } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const accountType = ref<'individual' | 'company'>('individual')

onMounted(() => {
  const type = route.query.type as string
  if (type === 'company' || type === 'individual') {
    accountType.value = type
  }
})

const form = reactive({
  email: '', password: '', confirmPassword: '',
  first_name: '', last_name: '', phone: '',
  company_name: '', legal_name: '', bin: '',
  contact_person: ''
})

const validateConfirm = (_rule: any, value: string, callback: any) => {
  if (value !== form.password) callback(new Error('Пароли не совпадают'))
  else callback()
}

const validateBIN = (_rule: any, value: string, callback: any) => {
  if (!/^\d{12}$/.test(value)) callback(new Error('БИН должен содержать 12 цифр'))
  else callback()
}

const baseRules = {
  email: [{ required: true, message: 'Введите email', trigger: 'blur' }, { type: 'email' as const, message: 'Неверный формат', trigger: 'blur' }],
  password: [{ required: true, message: 'Введите пароль', trigger: 'blur' }, { min: 6, message: 'Минимум 6 символов', trigger: 'blur' }],
  confirmPassword: [{ required: true, message: 'Подтвердите пароль', trigger: 'blur' }, { validator: validateConfirm, trigger: 'blur' }],
  first_name: [{ required: true, message: 'Введите имя', trigger: 'blur' }],
  last_name: [{ required: true, message: 'Введите фамилию', trigger: 'blur' }],
  phone: [{ required: true, message: 'Введите телефон', trigger: 'blur' }]
}

const companyRules = {
  company_name: [{ required: true, message: 'Введите название компании', trigger: 'blur' }],
  legal_name: [{ required: true, message: 'Введите юридическое название', trigger: 'blur' }],
  bin: [{ required: true, message: 'Введите БИН', trigger: 'blur' }, { validator: validateBIN, trigger: 'blur' }],
  contact_person: [{ required: true, message: 'Введите контактное лицо', trigger: 'blur' }]
}

const computedRules = computed(() => {
  if (accountType.value === 'company') {
    return { ...baseRules, ...companyRules }
  }
  return baseRules
})

async function handleRegister() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const data: any = {
      email: form.email,
      password: form.password,
      first_name: form.first_name,
      last_name: form.last_name,
      phone: form.phone,
      account_type: accountType.value
    }

    if (accountType.value === 'company') {
      data.company_name = form.company_name
      data.legal_name = form.legal_name
      data.bin = form.bin
      data.contact_person = form.contact_person
    }

    await authStore.register(data)
    ElMessage.success('Регистрация успешна!')
    await router.push('/broker/dashboard').catch(() => {
      window.location.href = '/broker/dashboard'
    })
  } catch (err: any) {
    const serverError = err.response?.data?.error || ''
    const errorMap: Record<string, string> = {
      'email already registered': 'Email уже зарегистрирован',
      'BIN already registered': 'Компания с таким БИН уже зарегистрирована',
      'failed to create company record, please try again': 'Ошибка создания компании, попробуйте ещё раз',
      'company_name, legal_name and bin are required for company account type': 'Заполните все данные компании',
    }
    ElMessage.error(errorMap[serverError] || serverError || 'Ошибка регистрации')
  } finally {
    loading.value = false
  }
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
.auth-card { width: 500px; padding: 20px; border-radius: 12px; }
</style>
