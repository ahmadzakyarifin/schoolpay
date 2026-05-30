import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../store/auth'
import LoginView from '../views/auth/LoginView.vue'
import DashboardLayout from '../layouts/DashboardLayout.vue'

const router = createRouter({
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            name: 'login',
            component: LoginView,
            meta: { guestOnly: true }
        },
        {
            path: '/forgot-password',
            name: 'forgot-password',
            component: () => import('../views/auth/ForgotPasswordView.vue'),
            meta: { guestOnly: true }
        },
        {
            path: '/reset-password',
            name: 'reset-password',
            component: () => import('../views/auth/ResetPasswordView.vue')
        },
        {
            path: '/activate',
            name: 'activate',
            component: () => import('../views/auth/ActivationView.vue')
        },
        {
            path: '/receipt/:id',
            name: 'receipt-print',
            component: () => import('../views/shared/ReceiptView.vue')
        },
        {
            path: '/dashboard',
            component: DashboardLayout,
            meta: { requiresAuth: true },
            children: [
                {
                    path: '',
                    name: 'dashboard',
                    component: () => import('../views/admin/DashboardView.vue'),
                    meta: { role: 'admin' }
                },
                // Admin Routes
                {
                    path: '/users',
                    name: 'user-management',
                    component: () => import('../views/users/UserManagementView.vue'),
                    meta: { role: 'admin' }
                },
                {
                    path: '/reports',
                    name: 'reports',
                    component: () => import('../views/admin/ReportsView.vue'),
                    meta: { role: 'admin' }
                },
                {
                    path: '/users/:id',
                    name: 'user-details',
                    component: () => import('../views/users/UserDetailsView.vue'),
                    meta: { role: 'admin' }
                },
                {
                    path: '/students',
                    name: 'student-management',
                    component: () => import('../views/students/StudentManagementView.vue'),
                    meta: { role: 'admin' }
                },
                {
                    path: '/students/:id',
                    name: 'student-details',
                    component: () => import('../views/students/StudentDetailsView.vue'),
                    meta: { role: 'admin' }
                },
                // Academic (Admin Only for management)
                {
                    path: '/academic/major',
                    name: 'major-management',
                    component: () => import('../views/academic/MajorManagementView.vue'),
                    meta: { role: 'admin' }
                },
                {
                    path: '/academic/class',
                    name: 'class-management',
                    component: () => import('../views/academic/ClassManagementView.vue'),
                    meta: { role: 'admin' }
                },
                {
                    path: '/academic/years',
                    name: 'academic-year-management',
                    component: () => import('../views/academic/AcademicYearManagementView.vue'),
                    meta: { role: 'admin' }
                },
                // Finance (Admin Only for management)
                {
                    path: '/finance/bill-types',
                    name: 'bill-types',
                    component: () => import('../views/finance/BillTypeManagementView.vue'),
                    meta: { role: 'admin' }
                },
                {
                    path: '/finance/bills',
                    name: 'all-bills',
                    component: () => import('../views/finance/FinanceBillManagementView.vue'),
                    meta: { role: 'admin' }
                },
                {
                    path: '/finance/rules',
                    name: 'billing-rules',
                    component: () => import('../views/finance/BillingRuleManagementView.vue'),
                    meta: { role: 'admin' }
                },
                {
                    path: '/notifications',
                    name: 'notification-logs',
                    component: () => import('../views/admin/NotificationLogsView.vue'),
                    meta: { role: 'admin' }
                },
                {
                    path: '/support/chat',
                    name: 'support-chat',
                    component: () => import('../views/admin/SupportChatView.vue'),
                    meta: { role: 'admin' }
                },
                {
                    path: '/audit-logs',
                    name: 'audit-logs',
                    component: () => import('../views/admin/AuditLogView.vue'),
                    meta: { role: 'admin' }
                },
                {
                    path: '/profile',
                    name: 'profile',
                    component: () => import('../views/shared/ProfileView.vue')
                },

                // Parent Routes
                {
                    path: '/parent/dashboard',
                    name: 'parent-dashboard',
                    component: () => import('../views/parent/ParentDashboard.vue'),
                    meta: { role: 'parent' }
                },
                {
                    path: '/parent/bills',
                    name: 'parent-bills',
                    component: () => import('../views/parent/ParentDashboard.vue'),
                    meta: { role: 'parent' }
                },
                {
                    path: '/parent/history',
                    name: 'parent-history',
                    component: () => import('../views/parent/ParentDashboard.vue'),
                    meta: { role: 'parent' }
                }
            ]
        }
    ]
})

router.beforeEach(async (to, from, next) => {
    const authStore = useAuthStore()
    if (!authStore.isInitialized) {
        await authStore.initializeAuth()
    }
    const isAuthenticated = authStore.isAuthenticated
    const userRole = authStore.userRole

    if (to.meta.guestOnly && isAuthenticated) {
        if (userRole === 'admin') return next({ name: 'dashboard' })
        return next({ name: 'parent-dashboard' })
    }

    if (to.meta.requiresAuth && !isAuthenticated) {
        return next({ name: 'login' })
    }

    if (to.meta.role && to.meta.role !== userRole) {
        return next({ name: 'login' })
    }

    next()
})

export default router
