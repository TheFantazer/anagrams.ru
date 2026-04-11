<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/userStore'
import { useGameStore } from '../stores/gameStore'

const router = useRouter()
const userStore = useUserStore()
const gameStore = useGameStore()

function handleSignOut() {
  userStore.signOut()
  router.push('/auth')
}
</script>

<template>
  <div class="settings-page">
    <h2 class="page-title">Settings</h2>

    <div class="settings-card">
      <h3 class="settings-h">Account</h3>
      <div class="info-row">
        <span class="info-label">Username</span>
        <span class="info-value">{{ userStore.username }}</span>
      </div>
      <div class="info-row">
        <span class="info-label">Email</span>
        <span class="info-value">{{ userStore.email }}</span>
      </div>
      <div class="info-row no-border">
        <span class="info-label">Joined</span>
        <span class="info-value">{{ userStore.joinedDate }}</span>
      </div>
    </div>

    <div class="settings-card">
      <h3 class="settings-h">Game defaults</h3>
      <label class="label">Language</label>
      <select v-model="gameStore.settingsLang" class="select">
        <option value="ru">Русский</option>
        <option value="en">English</option>
      </select>

      <label class="label">Default letters</label>
      <select v-model.number="gameStore.settingsLetters" class="select">
        <option v-for="n in [6, 7, 8, 9]" :key="n" :value="n">{{ n }} letters</option>
      </select>
    </div>

    <div class="settings-card">
      <h3 class="settings-h">Stats</h3>
      <div class="info-row">
        <span class="info-label">Games played</span>
        <span class="info-value">{{ userStore.gamesPlayed }}</span>
      </div>
      <div class="info-row">
        <span class="info-label">Best score</span>
        <span class="info-value">{{ userStore.bestScore.toLocaleString() }}</span>
      </div>
      <div class="info-row no-border">
        <span class="info-label">Longest word</span>
        <span class="info-value accent">{{ userStore.longestWord }}</span>
      </div>
    </div>

    <button class="btn-secondary" @click="handleSignOut">
      Sign out
    </button>
  </div>
</template>

<style scoped>
.settings-page {
  max-width: 500px;
  margin: 40px auto;
  padding: 0 24px;
}

.page-title {
  font-family: 'Space Mono', monospace;
  font-size: 22px;
  font-weight: 700;
  margin: 0 0 24px;
  color: #e8e6e1;
}

.settings-card {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 16px;
}

.settings-h {
  font-family: 'Space Mono', monospace;
  font-size: 14px;
  font-weight: 700;
  color: var(--accent);
  margin: 0 0 20px;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  padding: 10px 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);
  font-size: 14px;
}

.info-row.no-border {
  border-bottom: none;
}

.info-label {
  color: #666;
}

.info-value {
  color: #ccc;
  font-weight: 500;
}

.info-value.accent {
  color: var(--accent);
}

.label {
  font-size: 12px;
  color: #888;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin: 0 0 8px;
  display: block;
  font-weight: 500;
}

.select {
  width: 100%;
  padding: 12px 16px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #e8e6e1;
  font-size: 14px;
  outline: none;
  font-family: 'Outfit', sans-serif;
  appearance: none;
  margin-bottom: 20px;
  cursor: pointer;
}

.btn-secondary {
  width: 100%;
  padding: 14px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #e8e6e1;
  font-size: 15px;
  font-weight: 500;
  cursor: pointer;
  font-family: 'Outfit', sans-serif;
  transition: all 0.2s;
  margin-top: 8px;
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.1);
}
</style>
