<script setup>
import { ref, watch, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../stores/userStore'
import { useGameStore } from '../stores/gameStore'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()
const gameStore = useGameStore()

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
let saveMessageTimeout = null

watch([() => userStore.soloLang, () => userStore.soloLetters, () => userStore.soloTime], async () => {
  if (!userStore.isAuthenticated || !userStore.userId) return

  await saveSettings()
}, { deep: true })

async function saveSettings() {
  if (saving.value) return
  saving.value = true

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

    if (!saveMessage.value) {
      saveMessage.value = t('settings.saved')
    }

    if (saveMessageTimeout) {
      clearTimeout(saveMessageTimeout)
    }

    saveMessageTimeout = setTimeout(() => {
      saveMessage.value = ''
      saveMessageTimeout = null
    }, 2000)
  } catch (error) {
    console.error('Failed to save settings:', error)
    saveMessage.value = t('settings.failed')
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

const languageOptions = computed(() => [
  { id: 'en', label: t('settings.gameDefaults.languages.en') },
  { id: 'ru', label: t('settings.gameDefaults.languages.ru') }
])
</script>

<template>
  <div class="page">
    <div class="shell settings-wrap">
      <header class="page-head">
        <div>
          <div class="page-eyebrow">{{ $t('settings.title') }}</div>
          <h1 class="page-title-display">{{ $t('settings.subtitle') }}</h1>
        </div>
        <button class="btn btn--ghost" @click="handleSignOut">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4M16 17l5-5-5-5M21 12H9"/>
          </svg>
          {{ $t('nav.signOut') }}
        </button>
      </header>

      <div class="settings-grid">
        <!-- Account card -->
        <section class="card">
          <h3 style="margin:0 0 16px">{{ $t('settings.account.title') }}</h3>
          <div class="kv">
            <span class="kv-k">{{ $t('settings.account.username') }}</span>
            <span class="kv-v">{{ userStore.username }}</span>
          </div>
          <div class="kv">
            <span class="kv-k">{{ $t('settings.account.email') }}</span>
            <span class="kv-v">{{ userStore.email }}</span>
          </div>
          <div class="kv" style="border-bottom:0">
            <span class="kv-k">{{ $t('settings.account.joined') }}</span>
            <span class="kv-v mono">{{ userStore.joinedDate }}</span>
          </div>
        </section>

        <!-- Game defaults card -->
        <section class="card">
          <h3 style="margin:0 0 16px">{{ $t('settings.gameDefaults.title') }}</h3>
          <div class="field">
            <label class="field-label">{{ $t('settings.gameDefaults.language') }}</label>
            <div class="checkbox-row">
              <button
                v-for="lang in languageOptions"
                :key="lang.id"
                class="chip-toggle"
                :data-active="userStore.soloLang === lang.id"
                @click="userStore.soloLang = lang.id"
              >
                {{ lang.label }}
              </button>
            </div>
          </div>

          <div class="field">
            <label class="field-label">{{ $t('settings.gameDefaults.letters') }}</label>
            <div class="checkbox-row">
              <button
                v-for="n in [6, 7, 8]"
                :key="n"
                class="chip-toggle chip-toggle--mono"
                :data-active="userStore.soloLetters === n"
                @click="userStore.soloLetters = n"
              >
                {{ n }}
              </button>
            </div>
          </div>

          <div class="field" style="margin-bottom:0">
            <label class="field-label">{{ $t('settings.gameDefaults.timeLimit') }}</label>
            <div class="checkbox-row">
              <button
                v-for="time in [{ val: 30, label: '30s' }, { val: 60, label: '1:00' }, { val: 90, label: '1:30' }, { val: 120, label: '2:00' }]"
                :key="time.val"
                class="chip-toggle chip-toggle--mono"
                :data-active="userStore.soloTime === time.val"
                @click="userStore.soloTime = time.val"
              >
                {{ time.label }}
              </button>
            </div>
          </div>

          <div v-if="saveMessage" style="margin-top:12px; padding:10px; background:var(--success-soft); border:1px solid var(--success); border-radius:10px; color:var(--success); font-size:13px; text-align:center">
            {{ saveMessage }}
          </div>
        </section>

        <!-- Stats card -->
        <section class="card" style="grid-column:1 / -1">
          <h3 style="margin:0 0 16px">{{ $t('settings.stats.title') }}</h3>
          <div v-if="loadingStats" style="text-align:center; padding:40px; color:var(--fg-muted)">
            {{ $t('settings.stats.loading') }}
          </div>
          <div v-else class="stats-grid">
            <div class="stat-cell">
              <div class="stat-k">{{ $t('settings.stats.gamesPlayed') }}</div>
              <div class="stat-v">{{ userStore.gamesPlayed }}</div>
            </div>
            <div class="stat-cell">
              <div class="stat-k">{{ $t('settings.stats.bestScore') }}</div>
              <div class="stat-v accent-text">{{ userStore.bestScore.toLocaleString() }}</div>
            </div>
            <div class="stat-cell">
              <div class="stat-k">{{ $t('settings.stats.longestWord') }}</div>
              <div class="stat-v mono" style="font-size:22px">{{ userStore.longestWord || '—' }}</div>
            </div>
            <div class="stat-cell">
              <div class="stat-k">{{ $t('settings.stats.wordsFound') }}</div>
              <div class="stat-v">{{ stats.total_words?.toLocaleString() || '—' }}</div>
            </div>
            <div class="stat-cell">
              <div class="stat-k">{{ $t('settings.stats.avgScore') }}</div>
              <div class="stat-v">{{ Math.round(stats.average_score) || '—' }}</div>
            </div>
            <div class="stat-cell">
              <div class="stat-k">{{ $t('settings.stats.currentStreak') }}</div>
              <div class="stat-v">— <span class="muted" style="font-size:14px">{{ $t('settings.stats.days') }}</span></div>
            </div>
          </div>
        </section>
      </div>

      <div style="margin-top:24px">
        <button class="btn btn--accent btn--lg" @click="userStore.setShowSoloSettings(true)">
          {{ $t('settings.startNewGame') }}
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* All styles are in pages.css and app.css */
</style>
