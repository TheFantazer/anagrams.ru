<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/userStore'
import { useGameStore } from '../stores/gameStore'

const router = useRouter()
const userStore = useUserStore()
const gameStore = useGameStore()

function goHome() {
  gameStore.endGame()
  router.push('/')
}
</script>

<template>
  <header class="header">
    <div class="logo-box" @click="goHome">
      <span class="logo-text">
        <img src="../../public/icon.png" alt="icon" style="border-radius: 8px">
      </span>
    </div>
    <div class="header-right">
      <div class="icon-btn" @click="userStore.setShowHelp(true)" title="How to play">
        ?
      </div>
      <div class="icon-btn" @click="router.push('/leaderboard')" title="Leaderboard">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M8 21V12h8v9M3 21V15h5v6M16 21V9h5v12"/>
        </svg>
      </div>
      <div v-if="!userStore.isAuthenticated" class="icon-btn account-btn" @click="router.push('/auth')" title="Sign In">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M15 3h4a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2h-4M10 17l5-5-5-5M15 12H3"/>
        </svg>
      </div>
      <div v-else class="icon-btn account-btn" @click="router.push('/settings')" title="Account">
        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <circle cx="12" cy="8" r="4"/>
          <path d="M4 21v-1a6 6 0 0112 0v1"/>
        </svg>
      </div>
    </div>
  </header>
</template>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 24px;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  position: relative;
  z-index: 10;
  backdrop-filter: blur(12px);
  background: rgba(10, 10, 15, 0.8);
}

.logo-box {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgb(255, 255, 255);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.logo-box:hover {
  background: rgba(99, 230, 190, 0.15);
  border-color: rgba(99, 230, 190, 0.3);
}

.logo-text {
  font-family: 'Space Mono', monospace;
  font-weight: 700;
  font-size: 20px;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text img {
  width: 32px;
  height: 32px;
  object-fit: cover;
  border-radius: 6px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.icon-btn {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #999;
  transition: all 0.2s;
  font-size: 16px;
}

.icon-btn.account-btn {
  width: auto;
  padding: 0 12px;
  gap: 6px;
}

.icon-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #e8e6e1;
}
</style>
