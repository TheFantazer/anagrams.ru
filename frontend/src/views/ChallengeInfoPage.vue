<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../stores/userStore'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loading = ref(true)
const error = ref(null)
const session = ref(null)

const sessionId = computed(() => route.params.sessionId)

const languageLabel = computed(() => {
  if (!session.value) return ''
  return session.value.language === 'ru'
    ? t('settings.gameDefaults.languages.ru')
    : t('settings.gameDefaults.languages.en')
})

const timeLimitLabel = computed(() => {
  if (!session.value) return ''
  const seconds = session.value.time_limit
  if (seconds < 60) return `${seconds}s`
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return remainingSeconds > 0 ? `${minutes}:${String(remainingSeconds).padStart(2, '0')}` : `${minutes}:00`
})

// Check if current user has played this challenge
const hasUserPlayed = computed(() => {
  if (!session.value || !userStore.userId) return false
  return session.value.results?.some(r => r.user_id === userStore.userId) || false
})

async function loadSession() {
  loading.value = true
  error.value = null

  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/api/v1/sessions/${sessionId.value}`)

    if (!response.ok) {
      if (response.status === 404) {
        error.value = t('challenge.errors.notFound')
      } else {
        error.value = t('challenge.errors.loadFailed')
      }
      return
    }

    session.value = await response.json()
  } catch (err) {
    console.error('Failed to load session:', err)
    error.value = t('challenge.errors.loadFailed')
  } finally {
    loading.value = false
  }
}

function acceptChallenge() {
  if (!userStore.isAuthenticated) {
    // Сохраняем текущий путь для редиректа после авторизации
    sessionStorage.setItem('redirectAfterAuth', route.fullPath)
    router.push('/auth')
    return
  }

  // Переходим на GamePage с session_id
  router.push(`/game?session_id=${sessionId.value}`)
}

onMounted(() => {
  loadSession()
})
</script>

<template>
  <div class="page">
    <div class="shell challenge-wrap">
      <!-- Loading State -->
      <div v-if="loading" class="challenge-loading">
        <div class="spinner"></div>
        <p class="muted">{{ $t('common.loading') }}</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="challenge-error">
        <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="12" cy="12" r="10"/>
          <line x1="12" y1="8" x2="12" y2="12"/>
          <line x1="12" y1="16" x2="12.01" y2="16"/>
        </svg>
        <h2>{{ error }}</h2>
        <button class="btn btn--primary" @click="router.push('/')">
          {{ $t('challenge.backToHome') }}
        </button>
      </div>

      <!-- Challenge Info -->
      <div v-else-if="session" class="challenge-content">
        <div class="challenge-header">
          <div class="challenge-eyebrow">{{ $t('challenge.title') }}</div>
          <h1 class="challenge-title">{{ $t('challenge.subtitle') }}</h1>
        </div>

        <div class="card card--paper challenge-card">
          <div class="challenge-letters">
            <span
              v-for="(letter, i) in session.letters.toUpperCase().split('')"
              :key="i"
              class="challenge-tile"
            >
              {{ hasUserPlayed ? letter : '?' }}
            </span>
          </div>

          <div class="challenge-meta">
            <div class="meta-item">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"/>
                <path d="M12 6v6l4 2"/>
              </svg>
              <span>{{ timeLimitLabel }}</span>
            </div>
            <div class="meta-item">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10 10-4.5 10-10S17.5 2 12 2z"/>
                <path d="M12 6v6l4 2"/>
              </svg>
              <span>{{ languageLabel }}</span>
            </div>
            <div class="meta-item">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M4 7h16M4 12h16M4 17h10"/>
              </svg>
              <span>{{ session.letters.length }} {{ $t('challenge.letters') }}</span>
            </div>
          </div>

          <div class="challenge-description">
            <p>{{ $t('challenge.description') }}</p>
          </div>

          <button class="btn btn--accent btn--lg" @click="acceptChallenge">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M6 4l14 8-14 8z"/>
            </svg>
            {{ userStore.isAuthenticated ? $t('challenge.accept') : $t('challenge.signInToPlay') }}
          </button>
        </div>

        <div class="challenge-footer">
          <button class="btn btn--ghost" @click="router.push('/')">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M19 12H5M12 19l-7-7 7-7"/>
            </svg>
            {{ $t('challenge.backToHome') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.challenge-wrap {
  max-width: 640px;
  margin: 0 auto;
  padding: 40px 20px;
  min-height: 60vh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.challenge-loading,
.challenge-error {
  text-align: center;
  padding: 40px 20px;
}

.challenge-error svg {
  color: var(--danger);
  margin-bottom: 20px;
}

.challenge-error h2 {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 700;
  color: var(--fg-primary);
  margin: 0 0 24px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border-hairline);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.challenge-content {
  width: 100%;
}

.challenge-header {
  text-align: center;
  margin-bottom: 32px;
}

.challenge-eyebrow {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 2px;
  color: var(--fg-muted);
  font-weight: 600;
  margin-bottom: 8px;
}

.challenge-title {
  font-family: var(--font-display);
  font-size: 36px;
  font-weight: 700;
  letter-spacing: -1px;
  color: var(--fg-primary);
  margin: 0;
}

.challenge-card {
  padding: 40px;
  text-align: center;
}

.challenge-letters {
  display: flex;
  justify-content: center;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 32px;
}

.challenge-tile {
  width: 64px;
  height: 72px;
  border-radius: 14px;
  background: var(--navy);
  color: var(--milk);
  display: grid;
  place-items: center;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 26px;
  box-shadow: 0 3px 0 var(--navy-2), var(--shadow-sm);
}

.challenge-meta {
  display: flex;
  justify-content: center;
  gap: 24px;
  flex-wrap: wrap;
  margin-bottom: 24px;
  padding: 20px;
  background: var(--bg-surface);
  border-radius: 12px;
  border: 1px solid var(--border-hairline);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  font-weight: 600;
  color: var(--fg-secondary);
}

.meta-item svg {
  color: var(--accent);
}

.challenge-description {
  margin-bottom: 32px;
  color: var(--fg-secondary);
  font-size: 15px;
  line-height: 1.6;
}

.challenge-footer {
  text-align: center;
  margin-top: 32px;
}

@media (max-width: 600px) {
  .challenge-wrap {
    padding: 20px;
  }

  .challenge-title {
    font-size: 28px;
  }

  .challenge-card {
    padding: 28px 20px;
  }

  .challenge-tile {
    width: 52px;
    height: 60px;
    font-size: 22px;
  }

  .challenge-meta {
    gap: 16px;
    padding: 16px;
  }
}
</style>
