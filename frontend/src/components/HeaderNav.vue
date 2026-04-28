<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/userStore'
import { useGameStore } from '../stores/gameStore'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const gameStore = useGameStore()

const drawer = ref(false)
const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// Badge counts from backend
const yourTurnCount = ref(0) // Number of challenges waiting for your turn
const incomingReqCount = ref(0) // Number of incoming friend requests

// Load badge counts
async function loadBadgeCounts() {
  if (!userStore.userId) return

  try {
    // Load incoming friend requests count
    const friendsRes = await fetch(`${apiUrl}/api/v1/friends/requests/pending?user_id=${userStore.userId}`)
    if (friendsRes.ok) {
      const requests = await friendsRes.json()
      incomingReqCount.value = requests.length
    }

    // Load sessions where it's your turn
    const sessionsRes = await fetch(`${apiUrl}/api/v1/sessions/my?user_id=${userStore.userId}`)
    if (sessionsRes.ok) {
      const sessions = await sessionsRes.json()
      // Count sessions where user hasn't played yet
      yourTurnCount.value = sessions.filter(s => !s.played).length
    }
  } catch (error) {
    console.error('Failed to load badge counts:', error)
  }
}

// Watch for userId changes and reload counts
watch(() => userStore.userId, (newVal) => {
  if (newVal) {
    loadBadgeCounts()
  } else {
    yourTurnCount.value = 0
    incomingReqCount.value = 0
  }
}, { immediate: true })

// Reload counts on mount if user is authenticated
onMounted(() => {
  if (userStore.userId) {
    loadBadgeCounts()
  }
})

const links = computed(() => [
  { id: '/', label: t('nav.home'), icon: 'home' },
  { id: '/play', label: t('nav.play'), icon: 'play', badge: yourTurnCount.value },
  { id: '/friends', label: t('nav.friends'), icon: 'users', badge: incomingReqCount.value },
  { id: '/leaderboard', label: t('nav.leaderboard'), icon: 'trophy' },
])

// Route matching for active state
const routeMatch = {
  '/play': ['/play', '/match-history', '/results'],
  '/friends': ['/friends'],
  '/leaderboard': ['/leaderboard'],
  '/': ['/']
}

function isActive(linkId) {
  // Special case for home - exact match only
  if (linkId === '/') {
    return route.path === '/'
  }

  const paths = routeMatch[linkId] || [linkId]
  return paths.some(p => route.path.startsWith(p))
}

function goHome() {
  gameStore.endGame()
  router.push('/')
  drawer.value = false
}

function navigateTo(path) {
  router.push(path)
  drawer.value = false
}

function quickPlay() {
  gameStore.startGame(userStore.soloTime, userStore.soloLetters, userStore.soloLang)
  router.push('/game')
  drawer.value = false
}

function signOut() {
  userStore.signOut()
  router.push('/')
  drawer.value = false
}

const userInitial = computed(() => {
  if (!userStore.user?.username) return '?'
  return userStore.user.username[0].toUpperCase()
})
</script>

<template>
  <nav class="nav">
    <div class="shell nav-inner">
      <!-- Brand -->
      <div class="nav-brand" @click="goHome" title="anagrams.">
        <div class="nav-mark">AN</div>
      </div>

      <!-- Quick Play Button -->
      <button class="nav-quickplay" @click="quickPlay" :title="t('nav.quickPlay')">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
          <path d="M13 2L4 14h7l-1 8 9-12h-7l1-8z"/>
        </svg>
        <span>{{ t('nav.quickPlay') }}</span>
      </button>

      <!-- Desktop Links -->
      <div class="nav-links">
        <button
          v-for="link in links"
          :key="link.id"
          class="nav-link"
          :data-active="isActive(link.id)"
          @click="navigateTo(link.id)"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <template v-if="link.icon === 'home'">
              <path d="M3 12l9-9 9 9"/>
              <path d="M5 10v10h14V10"/>
            </template>
            <path v-else-if="link.icon === 'play'" d="M6 4l14 8-14 8z"/>
            <template v-else-if="link.icon === 'users'">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
              <path d="M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75"/>
            </template>
            <template v-else-if="link.icon === 'trophy'">
              <path d="M8 21h8M12 17v4M7 4h10v5a5 5 0 0 1-10 0V4zM17 4h3v3a3 3 0 0 1-3 3M7 4H4v3a3 3 0 0 0 3 3"/>
            </template>
          </svg>
          {{ link.label }}
          <span v-if="link.badge && link.badge > 0" class="nav-badge">{{ link.badge }}</span>
        </button>
      </div>

      <!-- Right Side -->
      <div class="nav-right">
        <!-- Help Button -->
        <button class="nav-icon" :title="$t('nav.help')" @click="userStore.setShowHelp(true)">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"/>
            <path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3M12 17h.01"/>
          </svg>
        </button>

        <!-- Sign In / Avatar -->
        <button
          v-if="!userStore.isAuthenticated"
          class="nav-pill"
          @click="navigateTo('/auth')"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
          {{ $t('nav.signIn') }}
        </button>
        <button
          v-else
          class="nav-avatar"
          :title="$t('nav.settings')"
          @click="navigateTo('/settings')"
        >
          {{ userInitial }}
        </button>

        <!-- Mobile Menu Button -->
        <button class="nav-burger" @click="drawer = !drawer">
          <svg v-if="!drawer" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M3 12h18M3 6h18M3 18h18"/>
          </svg>
          <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 6L6 18M6 6l12 12"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- Mobile Drawer -->
    <div v-if="drawer" class="nav-drawer">
      <!-- Quick Play -->
      <button class="nav-link nav-link--quick" @click="quickPlay">
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
          <path d="M13 2L4 14h7l-1 8 9-12h-7l1-8z"/>
        </svg>
        {{ t('nav.quickPlay') }}
      </button>

      <div class="hr" />

      <!-- Main Links -->
      <button
        v-for="link in links"
        :key="link.id"
        class="nav-link"
        :data-active="isActive(link.id)"
        @click="navigateTo(link.id)"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
          <template v-if="link.icon === 'home'">
            <path d="M3 12l9-9 9 9"/>
            <path d="M5 10v10h14V10"/>
          </template>
          <path v-else-if="link.icon === 'play'" d="M6 4l14 8-14 8z"/>
          <template v-else-if="link.icon === 'users'">
            <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75"/>
          </template>
          <template v-else-if="link.icon === 'trophy'">
            <path d="M8 21h8M12 17v4M7 4h10v5a5 5 0 0 1-10 0V4zM17 4h3v3a3 3 0 0 1-3 3M7 4H4v3a3 3 0 0 0 3 3"/>
          </template>
        </svg>
        {{ link.label }}
        <span v-if="link.badge && link.badge > 0" class="nav-badge">{{ link.badge }}</span>
      </button>

      <div class="hr" />

      <!-- Auth Actions -->
      <template v-if="userStore.isAuthenticated">
        <button class="nav-link" @click="navigateTo('/settings')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="3"/>
            <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
          </svg>
          {{ $t('nav.settings') }}
        </button>
        <button class="nav-link" @click="signOut">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4M16 17l5-5-5-5M21 12H9"/>
          </svg>
          {{ $t('nav.signOut') }}
        </button>
      </template>
      <template v-else>
        <button class="nav-link" @click="navigateTo('/auth')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
          {{ $t('nav.signIn') }}
        </button>
      </template>
    </div>
  </nav>
</template>

<style scoped>
/* Most styles are in app.css, only component-specific overrides here */
.nav {
  position: sticky;
  top: 0;
  z-index: 20;
  background: color-mix(in oklab, var(--milk) 88%, transparent);
  backdrop-filter: blur(14px);
  -webkit-backdrop-filter: blur(14px);
  border-bottom: 1px solid var(--border-hairline);
}

.nav-inner {
  display: flex;
  align-items: center;
  gap: var(--sp-4);
  height: 68px;
}

.nav-brand {
  cursor: pointer;
  user-select: none;
}

.nav-mark {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: var(--navy);
  display: grid;
  place-items: center;
  color: var(--milk);
  font-family: var(--font-display);
  font-weight: 800;
  font-size: 18px;
  letter-spacing: -0.5px;
  box-shadow: var(--shadow-sm);
  transition: transform var(--dur-base) var(--ease-out);
}

.nav-brand:hover .nav-mark {
  transform: rotate(-4deg);
}

/* Quick Play Button */
.nav-quickplay {
  appearance: none;
  border: 0;
  background: var(--accent);
  color: var(--milk);
  padding: 0 16px;
  height: 36px;
  border-radius: 10px;
  font-family: var(--font-body);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: all var(--dur-base) var(--ease-out);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.nav-quickplay:hover {
  background: color-mix(in oklab, var(--accent) 90%, black);
  transform: translateY(-1px);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.12);
}

.nav-quickplay span {
  white-space: nowrap;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 4px;
  flex: 1;
  justify-content: center;
}

.nav-link {
  appearance: none;
  border: 0;
  background: transparent;
  padding: 8px 14px;
  border-radius: 10px;
  font-family: var(--font-body);
  font-size: 13px;
  font-weight: 500;
  color: var(--fg-secondary);
  cursor: pointer;
  transition: all var(--dur-base) var(--ease-std);
  display: inline-flex;
  align-items: center;
  gap: 8px;
  position: relative;
}

.nav-link:hover {
  background: var(--bg-hover);
  color: var(--fg-primary);
}

.nav-link[data-active="true"] {
  background: var(--navy);
  color: var(--milk);
}

/* Badge */
.nav-badge {
  background: var(--danger);
  color: var(--milk);
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 700;
  min-width: 18px;
  height: 18px;
  border-radius: 999px;
  display: inline-grid;
  place-items: center;
  padding: 0 5px;
  margin-left: -2px;
}

.nav-link[data-active="true"] .nav-badge {
  background: var(--milk);
  color: var(--danger);
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 6px;
}

.nav-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: transparent;
  border: 1px solid var(--border-hairline);
  color: var(--fg-muted);
  display: grid;
  place-items: center;
  cursor: pointer;
  transition: all var(--dur-base) var(--ease-std);
}

.nav-icon:hover {
  background: var(--milk-2);
  color: var(--navy);
  border-color: var(--border-default);
}

.nav-pill {
  padding: 0 14px;
  height: 36px;
  border-radius: 10px;
  background: var(--navy);
  color: var(--milk);
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 600;
  border: 0;
  cursor: pointer;
  font-family: var(--font-body);
  transition: transform var(--dur-base) var(--ease-out), background var(--dur-base);
}

.nav-pill:hover {
  background: color-mix(in oklab, var(--navy) 90%, black);
  transform: translateY(-1px);
}

.nav-avatar {
  width: 32px;
  height: 32px;
  border-radius: 999px;
  background: var(--grad-accent);
  color: var(--milk);
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 13px;
  display: grid;
  place-items: center;
  cursor: pointer;
  border: 2px solid var(--milk);
  box-shadow: 0 0 0 1px var(--border-default);
  transition: transform var(--dur-base);
}

.nav-avatar:hover {
  transform: scale(1.05);
}

.nav-burger {
  display: none;
}

/* Horizontal divider */
.hr {
  height: 1px;
  background: var(--border-hairline);
  margin: 8px 0;
}

@media (max-width: 900px) {
  .nav-quickplay span {
    display: none;
  }

  .nav-quickplay {
    width: 36px;
    padding: 0;
    justify-content: center;
  }
}

@media (max-width: 760px) {
  .nav-links {
    display: none;
  }

  .nav-quickplay {
    display: none;
  }

  .nav-burger {
    display: grid;
    place-items: center;
    width: 36px;
    height: 36px;
    border-radius: 10px;
    background: transparent;
    border: 1px solid var(--border-hairline);
    color: var(--navy);
    cursor: pointer;
    transition: all var(--dur-base);
  }

  .nav-burger:hover {
    background: var(--bg-hover);
  }

  .nav-drawer {
    position: absolute;
    top: 68px;
    left: 0;
    right: 0;
    background: var(--milk);
    border-bottom: 1px solid var(--border-subtle);
    padding: 12px 20px 20px;
    display: flex;
    flex-direction: column;
    gap: 4px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  }

  .nav-drawer .nav-link {
    padding: 12px 14px;
    font-size: 15px;
    width: 100%;
    justify-content: flex-start;
  }

  .nav-drawer .nav-link--quick {
    background: var(--accent);
    color: var(--milk);
    font-weight: 600;
  }

  .nav-drawer .nav-link--quick:hover {
    background: color-mix(in oklab, var(--accent) 90%, black);
  }
}
</style>
