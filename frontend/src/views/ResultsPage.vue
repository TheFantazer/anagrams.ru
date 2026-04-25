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
const results = ref([])

const sessionId = computed(() => route.params.sessionId)

const myResult = computed(() => {
  if (!results.value || !userStore.userId) return null
  return results.value.find(r => r.user_id === userStore.userId)
})

const opponentResult = computed(() => {
  if (!results.value || !userStore.userId) return null
  return results.value.find(r => r.user_id !== userStore.userId)
})

const gameOutcome = computed(() => {
  if (!myResult.value || !opponentResult.value) return null

  if (myResult.value.score > opponentResult.value.score) {
    return 'won'
  } else if (myResult.value.score < opponentResult.value.score) {
    return 'lost'
  } else {
    return 'tie'
  }
})

const myWords = computed(() => {
  return myResult.value?.found_words?.map(w => w.toUpperCase()) || []
})

const opponentWords = computed(() => {
  return opponentResult.value?.found_words?.map(w => w.toUpperCase()) || []
})

const allValidWords = computed(() => {
  return session.value?.valid_words?.map(w => w.toUpperCase()) || []
})

const sortedWords = computed(() => {
  return [...allValidWords.value].sort((a, b) => {
    if (a.length !== b.length) return b.length - a.length
    return a.localeCompare(b)
  })
})

const wordComparison = computed(() => {
  return sortedWords.value.map(word => ({
    word,
    foundByMe: myWords.value.includes(word),
    foundByOpponent: opponentWords.value.includes(word),
    foundByBoth: myWords.value.includes(word) && opponentWords.value.includes(word)
  }))
})

async function loadResults() {
  loading.value = true
  error.value = null

  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

    // Load session data
    const sessionResponse = await fetch(`${apiUrl}/api/v1/sessions/${sessionId.value}`)
    if (!sessionResponse.ok) {
      throw new Error('Failed to load session')
    }
    session.value = await sessionResponse.json()

    // Load results
    const resultsResponse = await fetch(`${apiUrl}/api/v1/sessions/${sessionId.value}/results`)
    if (!resultsResponse.ok) {
      throw new Error('Failed to load results')
    }
    results.value = await resultsResponse.json()

  } catch (err) {
    console.error('Failed to load results:', err)
    error.value = 'Failed to load results'
  } finally {
    loading.value = false
  }
}

function goBack() {
  router.push('/multiplayer')
}

function playAgain() {
  router.push(`/play/${sessionId.value}`)
}

onMounted(() => {
  loadResults()
})
</script>

<template>
  <div class="page">
    <div class="shell results-wrap">
      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p class="muted">{{ $t('common.loading') }}</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="error-state">
        <p class="muted">{{ error }}</p>
        <button class="btn btn--primary" @click="goBack">
          {{ $t('challenge.backToHome') }}
        </button>
      </div>

      <!-- Results -->
      <div v-else-if="session && results.length > 0">
        <div class="results-header">
          <div class="results-eyebrow">{{ $t('challenge.title') }}</div>
          <h1 class="results-title">{{ $t('game.gameOver.title') }}</h1>
        </div>

        <!-- Victory/Defeat Banner -->
        <div v-if="gameOutcome" class="outcome-banner" :class="`outcome-${gameOutcome}`">
          <svg v-if="gameOutcome === 'won'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
          </svg>
          <svg v-else-if="gameOutcome === 'lost'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M10 15l4-4m0 4l-4-4m13 1a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 2v20M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6"/>
          </svg>
          <span>
            {{ gameOutcome === 'won' ? $t('game.gameOver.youWon') : (gameOutcome === 'lost' ? $t('game.gameOver.youLost') : $t('game.gameOver.tie')) }}
          </span>
        </div>

        <!-- Scores Comparison -->
        <div class="scores-comparison">
          <div class="score-card">
            <div class="score-label">{{ $t('game.gameOver.yourScore') }}</div>
            <div class="score-value">{{ myResult?.score.toLocaleString() || 0 }}</div>
            <div class="score-meta">
              {{ myResult?.found_words?.length || 0 }} {{ $t('game.gameOver.wordsFound') }}
            </div>
          </div>

          <div class="vs-divider">{{ $t('multiplayer.vs') }}</div>

          <div class="score-card">
            <div class="score-label">{{ opponentResult?.player_name || $t('multiplayer.from') }}</div>
            <div class="score-value">{{ opponentResult?.score.toLocaleString() || 0 }}</div>
            <div class="score-meta">
              {{ opponentResult?.found_words?.length || 0 }} {{ $t('game.gameOver.wordsFound') }}
            </div>
          </div>
        </div>

        <!-- Letters -->
        <div class="results-letters">
          <span v-for="(letter, i) in session.letters.split('')" :key="i" class="results-tile">
            {{ letter.toUpperCase() }}
          </span>
        </div>

        <!-- Word Comparison -->
        <div class="word-comparison-section">
          <h3 class="section-title">{{ $t('game.gameOver.allWords') }}</h3>

          <div class="word-grid">
            <div
              v-for="(wordObj, i) in wordComparison"
              :key="i"
              :class="['word-chip-result', {
                'found-me': wordObj.foundByMe && !wordObj.foundByBoth,
                'found-opponent': wordObj.foundByOpponent && !wordObj.foundByBoth,
                'found-both': wordObj.foundByBoth,
                'missed': !wordObj.foundByMe && !wordObj.foundByOpponent
              }]"
            >
              {{ wordObj.word.toLowerCase() }}
            </div>
          </div>

          <div class="legend">
            <div class="legend-item">
              <span class="legend-color found-me"></span>
              <span>{{ $t('game.gameOver.myWords') }}</span>
            </div>
            <div class="legend-item">
              <span class="legend-color found-opponent"></span>
              <span>{{ $t('game.gameOver.opponentWords') }}</span>
            </div>
            <div class="legend-item">
              <span class="legend-color found-both"></span>
              <span>{{ $t('game.gameOver.allWords') }}</span>
            </div>
          </div>
        </div>

        <!-- Actions -->
        <div class="results-actions">
          <button class="btn btn--accent btn--lg" @click="playAgain">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M13 2L4 14h7l-1 8 9-12h-7l1-8z"/>
            </svg>
            {{ $t('game.gameOver.playAgain') }}
          </button>
          <button class="btn btn--primary btn--lg" @click="goBack">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M19 12H5M12 19l-7-7 7-7"/>
            </svg>
            {{ $t('multiplayer.backBtn') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.results-wrap {
  max-width: 800px;
  margin: 0 auto;
  padding: 40px 20px;
}

.loading-state,
.error-state {
  text-align: center;
  padding: 60px 20px;
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

.results-header {
  text-align: center;
  margin-bottom: 32px;
}

.results-eyebrow {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 2px;
  color: var(--fg-muted);
  font-weight: 600;
  margin-bottom: 8px;
}

.results-title {
  font-family: var(--font-display);
  font-size: 40px;
  font-weight: 700;
  letter-spacing: -1.2px;
  color: var(--fg-primary);
  margin: 0;
}

.outcome-banner {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 16px 24px;
  border-radius: 12px;
  margin-bottom: 32px;
  font-weight: 700;
  font-size: 16px;
}

.outcome-won {
  background: var(--success-soft);
  color: var(--success);
}

.outcome-lost {
  background: var(--danger-soft);
  color: var(--danger);
}

.outcome-tie {
  background: var(--bg-surface);
  color: var(--fg-secondary);
  border: 2px solid var(--border-default);
}

.scores-comparison {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 24px;
  align-items: center;
  margin-bottom: 32px;
  padding: 32px;
  background: var(--bg-surface);
  border-radius: 20px;
  border: 1px solid var(--border-subtle);
}

.score-card {
  text-align: center;
}

.score-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--fg-muted);
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 8px;
}

.score-value {
  font-family: var(--font-display);
  font-size: 48px;
  font-weight: 700;
  color: var(--accent);
  line-height: 1;
  margin-bottom: 8px;
}

.score-meta {
  font-size: 13px;
  color: var(--fg-secondary);
}

.vs-divider {
  font-family: var(--font-display);
  font-size: 20px;
  font-weight: 700;
  color: var(--fg-muted);
}

.results-letters {
  display: flex;
  justify-content: center;
  gap: 12px;
  flex-wrap: wrap;
  margin-bottom: 40px;
}

.results-tile {
  width: 56px;
  height: 64px;
  border-radius: 12px;
  background: var(--navy);
  color: var(--milk);
  display: grid;
  place-items: center;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 24px;
  box-shadow: 0 3px 0 var(--navy-2);
}

.word-comparison-section {
  margin-bottom: 40px;
}

.section-title {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 700;
  color: var(--fg-primary);
  margin: 0 0 20px;
  text-align: center;
}

.word-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: 8px;
  margin-bottom: 20px;
}

.word-chip-result {
  padding: 10px 14px;
  border-radius: 8px;
  font-family: var(--font-mono);
  font-size: 13px;
  font-weight: 600;
  text-align: center;
  border: 2px solid transparent;
  transition: all 0.2s;
}

.word-chip-result.found-me {
  background: var(--accent-soft);
  color: var(--accent);
  border-color: var(--accent);
}

.word-chip-result.found-opponent {
  background: var(--danger-soft);
  color: var(--danger);
  border-color: var(--danger);
}

.word-chip-result.found-both {
  background: var(--success-soft);
  color: var(--success);
  border-color: var(--success);
}

.word-chip-result.missed {
  background: var(--bg-surface);
  color: var(--fg-muted);
  border-color: var(--border-hairline);
}

.legend {
  display: flex;
  justify-content: center;
  gap: 24px;
  flex-wrap: wrap;
  padding: 16px;
  background: var(--bg-surface);
  border-radius: 12px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--fg-secondary);
}

.legend-color {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  border: 2px solid transparent;
}

.legend-color.found-me {
  background: var(--accent-soft);
  border-color: var(--accent);
}

.legend-color.found-opponent {
  background: var(--danger-soft);
  border-color: var(--danger);
}

.legend-color.found-both {
  background: var(--success-soft);
  border-color: var(--success);
}

.results-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  flex-wrap: wrap;
}

@media (max-width: 600px) {
  .results-title {
    font-size: 32px;
  }

  .scores-comparison {
    grid-template-columns: 1fr;
    gap: 16px;
    padding: 24px;
  }

  .vs-divider {
    display: none;
  }

  .score-value {
    font-size: 40px;
  }

  .word-grid {
    grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  }
}
</style>
