import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'landing',
      component: () => import('@/views/public/Landing.vue')
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/public/Login.vue')
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/public/Register.vue')
    },
    {
      path: '/broker',
      component: () => import('@/components/layout/BrokerLayout.vue'),
      meta: { requiresAuth: true, role: 'client' },
      children: [
        { path: '', redirect: '/broker/dashboard' },
        { path: 'dashboard', name: 'broker-dashboard', component: () => import('@/views/broker/Dashboard.vue') },
        { path: 'documents', name: 'broker-documents', component: () => import('@/views/broker/Documents.vue') },
        { path: 'subscription', name: 'broker-subscription', component: () => import('@/views/broker/Subscription.vue') },
        { path: 'profile', name: 'broker-profile', component: () => import('@/views/broker/Profile.vue') }
      ]
    },
    {
      path: '/admin',
      component: () => import('@/components/layout/AdminLayout.vue'),
      meta: { requiresAuth: true, role: 'admin' },
      children: [
        { path: '', redirect: '/admin/dashboard' },
        { path: 'dashboard', name: 'admin-dashboard', component: () => import('@/views/admin/Dashboard.vue') },
        { path: 'companies', name: 'admin-companies', component: () => import('@/views/admin/Companies.vue') },
        { path: 'subscriptions', name: 'admin-subscriptions', component: () => import('@/views/admin/Subscriptions.vue') },
        { path: 'documents', name: 'admin-documents', component: () => import('@/views/admin/Documents.vue') }
      ]
    }
  ]
})

router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('access_token')
  const userStr = localStorage.getItem('user')

  if (to.meta.requiresAuth && !token) {
    return next('/login')
  }

  if (to.meta.role && userStr) {
    const user = JSON.parse(userStr)
    const roleName = user.role?.name || ''
    if (roleName !== to.meta.role) {
      return next(roleName === 'admin' ? '/admin/dashboard' : '/broker/dashboard')
    }
  }

  if ((to.name === 'login' || to.name === 'register') && token && userStr) {
    const user = JSON.parse(userStr)
    const roleName = user.role?.name || ''
    return next(roleName === 'admin' ? '/admin/dashboard' : '/broker/dashboard')
  }

  next()
})

export default router
