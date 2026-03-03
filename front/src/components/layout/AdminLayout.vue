<template>
  <el-container style="min-height: 100vh">
    <el-aside width="220px" style="background: #1a1a2e">
      <div style="padding: 20px; text-align: center">
        <h2 style="color: #26A65B; margin: 0; font-size: 20px; font-weight: 800; letter-spacing: 2px">KEDEN</h2>
        <p style="color: #909399; font-size: 12px; margin: 5px 0 0">Админ-панель</p>
      </div>
      <el-menu
        :default-active="$route.path"
        router
        background-color="#1a1a2e"
        text-color="#bfcbd9"
        active-text-color="#26A65B"
      >
        <el-menu-item index="/admin/dashboard">
          <el-icon><Odometer /></el-icon>
          <span>Дэшборд</span>
        </el-menu-item>
        <el-menu-item index="/admin/companies">
          <el-icon><OfficeBuilding /></el-icon>
          <span>Компании</span>
        </el-menu-item>
        <el-menu-item index="/admin/subscriptions">
          <el-icon><Ticket /></el-icon>
          <span>Подписки</span>
        </el-menu-item>
        <el-menu-item index="/admin/documents">
          <el-icon><Document /></el-icon>
          <span>Документы</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-header style="display: flex; align-items: center; justify-content: flex-end; background: #fff; box-shadow: 0 1px 4px rgba(0,21,41,.08)">
        <span style="margin-right: 16px; color: #606266">Администратор</span>
        <el-button type="danger" text @click="authStore.logout">Выйти</el-button>
      </el-header>
      <el-main style="background: #f5f7f5; padding: 24px">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
const authStore = useAuthStore()
</script>
