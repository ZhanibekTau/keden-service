import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/services/api'
import type { User, AuthResponse, LoginRequest, RegisterRequest } from '@/types'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)

  const isAuthenticated = computed(() => !!accessToken.value)
  const isClient = computed(() => user.value?.role?.name === 'client')
  const isAdmin = computed(() => user.value?.role?.name === 'admin')
  const displayName = computed(() => {
    if (!user.value) return ''
    return `${user.value.first_name} ${user.value.last_name}`
  })

  function init() {
    const token = localStorage.getItem('access_token')
    const savedUser = localStorage.getItem('user')
    if (token && savedUser) {
      accessToken.value = token
      user.value = JSON.parse(savedUser)
    }
  }

  async function login(data: LoginRequest) {
    const response = await api.post<AuthResponse>('/auth/login', data)
    setAuth(response.data)
  }

  async function register(data: RegisterRequest) {
    const response = await api.post<AuthResponse>('/auth/register', data)
    setAuth(response.data)
  }

  function setAuth(auth: AuthResponse) {
    accessToken.value = auth.access_token
    user.value = auth.user
    localStorage.setItem('access_token', auth.access_token)
    localStorage.setItem('refresh_token', auth.refresh_token)
    localStorage.setItem('user', JSON.stringify(auth.user))
  }

  async function logout() {
    const refreshToken = localStorage.getItem('refresh_token')
    if (refreshToken) {
      try {
        await api.post('/auth/logout', { refresh_token: refreshToken })
      } catch { /* ignore */ }
    }
    accessToken.value = null
    user.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')
    router.push('/login')
  }

  init()

  return { user, accessToken, isAuthenticated, isClient, isAdmin, displayName, login, register, logout }
})
