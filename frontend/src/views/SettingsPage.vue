<script setup>
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/userStore'
import { useGameStore } from '../stores/gameStore'

const router = useRouter()
const userStore = useUserStore()
const gameStore = useGameStore()

// Redirect to auth if not logged in
if (!userStore.isAuthenticated) {
  router.push('/auth')
}

const saving = ref(false)
const saveMessage = ref('')
const stats = ref({
  games_played: 0,
  best_score: 0,
  longest_word: '',
  total_words: 0,
  average_score: 0
})
const loadingStats = ref(false)

watch([() => userStore.soloLang, () => userStore.soloLetters, () => userStore.soloTime], async () => {
  if (!userStore.isAuthenticated || !userStore.userId) return

  await saveSettings()
}, { deep: true })

async function saveSettings() {
  if (saving.value) return
  saving.value = true
  saveMessage.value = ''

  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/api/v1/auth/settings?user_id=${userStore.userId}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        letter_count: userStore.soloLetters,
        language: userStore.soloLang,
        time_limit: userStore.soloTime
      })
    })

    if (!response.ok) {
      throw new Error('Failed to save settings')
    }

    const data = await response.json()
    userStore.setUser(data)
    saveMessage.value = 'Settings saved!'
    setTimeout(() => {
      saveMessage.value = ''
    }, 2000)
  } catch (error) {
    console.error('Failed to save settings:', error)
    saveMessage.value = 'Failed to save settings'
  } finally {
    saving.value = false
  }
}

async function loadStats() {
  if (!userStore.isAuthenticated || !userStore.userId) return

  loadingStats.value = true
  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/api/v1/auth/stats?user_id=${userStore.userId}`)

    if (response.ok) {
      stats.value = await response.json()
      if (stats.value.games_played !== undefined) {
        userStore.gamesPlayed = stats.value.games_played
      }
      if (stats.value.best_score !== undefined) {
        userStore.bestScore = stats.value.best_score
      }
      if (stats.value.longest_word) {
        userStore.longestWord = stats.value.longest_word
      }
    }
  } catch (error) {
    console.error('Failed to load stats:', error)
  } finally {
    loadingStats.value = false
  }
}

onMounted(() => {
  loadStats()
})

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
      <select v-model="userStore.soloLang" class="select">
        <option value="ru">Русский</option>
        <option value="en">English</option>
      </select>

      <label class="label">Default letters</label>
      <select v-model.number="userStore.soloLetters" class="select">
        <option v-for="n in [6, 7, 8, 9, 10]" :key="n" :value="n">{{ n }} letters</option>
      </select>

      <label class="label">Time limit (seconds)</label>
      <select v-model.number="userStore.soloTime" class="select">
        <option :value="30">30 seconds</option>
        <option :value="60">60 seconds</option>
        <option :value="90">90 seconds</option>
        <option :value="120">2 minutes</option>
      </select>

      <div v-if="saveMessage" class="save-message">{{ saveMessage }}</div>
    </div>

    <div class="settings-card">
      <h3 class="settings-h">Stats</h3>
      <div v-if="loadingStats" class="loading-text">Loading stats...</div>
      <div v-else>
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

.save-message {
  margin-top: 12px;
  padding: 10px 14px;
  border-radius: 8px;
  background: rgba(99, 230, 190, 0.15);
  border: 1px solid rgba(99, 230, 190, 0.3);
  color: #63e6be;
  font-size: 13px;
  text-align: center;
}

.loading-text {
  color: #999;
  font-size: 14px;
  text-align: center;
  padding: 20px;
}
</style>
