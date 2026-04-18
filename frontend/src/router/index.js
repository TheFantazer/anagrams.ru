import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: () => import('../views/HomePage.vue')
        },
        {
            path: '/game',
            name: 'game',
            component: () => import('../views/GamePage.vue')
        },
        {
            path: '/settings',
            name: 'settings',
            component: () => import('../views/SettingsPage.vue')
        },
        {
            path: '/auth',
            name: 'auth',
            component: () => import('../views/AuthPage.vue')
        },
        {
            path: '/auth/callback',
            name: 'auth-callback',
            component: () => import('../views/AuthCallbackPage.vue')
        },
        {
            path: '/leaderboard',
            name: 'leaderboard',
            component: () => import('../views/LeaderboardPage.vue')
        },
        {
            path: '/multiplayer',
            name: 'multiplayer',
            component: () => import('../views/MultiplayerPage.vue')
        }
    ]
})

export default router