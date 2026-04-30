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
            path: '/play',
            name: 'play',
            component: () => import('../views/PlayHubPage.vue')
        },
        {
            path: '/challenge/:sessionId',
            name: 'challenge',
            component: () => import('../views/ChallengeInfoPage.vue')
        },
        {
            path: '/results/:sessionId',
            name: 'results',
            component: () => import('../views/ResultsPage.vue')
        },
        {
            path: '/friends',
            name: 'friends',
            component: () => import('../views/FriendsPage.vue')
        },
        {
            path: '/match-history',
            name: 'match-history',
            component: () => import('../views/MatchHistoryPage.vue')
        },
        {
            path: '/solo-history',
            name: 'solo-history',
            component: () => import('../views/SoloHistoryPage.vue')
        },
        {
            path: '/privacy',
            name: 'privacy',
            component: () => import('../views/PrivacyPage.vue')
        }
    ]
})

export default router