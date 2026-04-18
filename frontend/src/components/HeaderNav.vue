<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '../stores/userStore'
import { useGameStore } from '../stores/gameStore'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const gameStore = useGameStore()

const drawer = ref(false)

const links = [
  { id: '/', label: 'Play', icon: 'play' },
  { id: '/multiplayer', label: 'Multiplayer', icon: 'users' },
  { id: '/leaderboard', label: 'Leaderboard', icon: 'trophy' },
]

function goHome() {
  gameStore.endGame()
  router.push('/')
  drawer.value = false
}

function navigateTo(path) {
  router.push(path)
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
      <div class="nav-brand" @click="goHome">
        <div class="nav-mark">AN</div>
        <span class="nav-wordmark">anagrams<span class="dot">.</span></span>
      </div>

      <!-- Desktop Links -->
      <div class="nav-links">
        <button
          v-for="link in links"
          :key="link.id"
          class="nav-link"
          :data-active="route.path === link.id"
          @click="navigateTo(link.id)"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path v-if="link.icon === 'play'" d="M6 4l14 8-14 8z"/>
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
        </button>
      </div>

      <!-- Right Side -->
      <div class="nav-right">
        <!-- Help Button -->
        <button class="nav-icon" title="How to play" @click="userStore.setShowHelp(true)">
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
          Sign in
        </button>
        <button
          v-else
          class="nav-avatar"
          title="Settings"
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
      <button
        v-for="link in links"
        :key="link.id"
        class="nav-link"
        :data-active="route.path === link.id"
        @click="navigateTo(link.id)"
      >
        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
          <path v-if="link.icon === 'play'" d="M6 4l14 8-14 8z"/>
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
      </button>
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
  justify-content: space-between;
  height: 68px;
}

.nav-brand {
  display: flex;
  align-items: center;
  gap: 12px;
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

.nav-wordmark {
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 15px;
  color: var(--fg-primary);
  letter-spacing: -0.2px;
}

.nav-wordmark .dot {
  color: var(--accent);
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 4px;
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
}

.nav-link:hover {
  background: var(--bg-hover);
  color: var(--fg-primary);
}

.nav-link[data-active="true"] {
  background: var(--navy);
  color: var(--milk);
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
  background: var(--navy-2);
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
}

.nav-burger {
  display: none;
}

@media (max-width: 760px) {
  .nav-links {
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
  }

  .nav-drawer .nav-link {
    padding: 12px 14px;
    font-size: 15px;
    width: 100%;
    justify-content: flex-start;
  }
}
</style>
