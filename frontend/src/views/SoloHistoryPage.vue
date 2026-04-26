<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../stores/userStore'
import Pagination from '../components/Pagination.vue'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()

// Filters
const typeFilter = ref('all') // all, daily, practice

// Pagination
const page = ref(1)
const pageSize = 15

// Data
const games = ref([])
const totalGames = ref(0)
const loading = ref(false)

// Stats
const stats = ref({
  totalGames: 0,
  bestScore: 0,
  totalWords: 0,
  totalPoints: 0
})

// Fetch solo games from backend
async function fetchGames() {
  if (!userStore.userId) {
    loading.value = false
    return
  }

  loading.value = true
  try {
    const params = new URLSearchParams({
      user_id: userStore.userId,
      page: page.value.toString(),
      per_page: pageSize.toString()
    })

    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/api/v1/sessions/all?${params}`)

    if (!response.ok) {
      throw new Error('Failed to fetch solo history')
    }

    const data = await response.json()

    // Filter only solo sessions (with single result)
    const soloSessions = (data.sessions || []).filter(s => {
      return s.results && s.results.length === 1 && s.results[0].user_id === userStore.userId
    })

    // Convert to game format
    games.value = soloSessions.map(session => {
      const result = session.results[0]

      return {
        id: session.id,
        letters: session.letters,
        language: session.language,
        isDaily: false, // TODO: Determine if it's a daily puzzle
        score: result.score || 0,
        wordsFound: result.word_count || 0,
        timeLimit: session.time_limit,
        createdAt: result.played_at || session.created_at
      }
    })

    // Apply type filter
    let filtered = games.value
    if (typeFilter.value !== 'all') {
      filtered = filtered.filter(g => {
        if (typeFilter.value === 'daily') return g.isDaily
        if (typeFilter.value === 'practice') return !g.isDaily
        return true
      })
    }

    games.value = filtered
    totalGames.value = filtered.length

    // Calculate stats
    calculateMockStats()
  } catch (error) {
    console.error('Error fetching solo history:', error)
    userStore.showToast('Failed to load solo history', 'error')

    // Fallback to mock data
    games.value = generateMockGames()
    totalGames.value = games.value.length
    calculateMockStats()
  } finally {
    loading.value = false
  }
}

// Mock data generator (fallback)
function generateMockGames() {
  const letterSets = [
    'АБВГДЕЖ', 'TESTING', 'EXAMPLE', 'WORDGAM', 'PLAYNOW',
    'АНАГРАМ', 'СЛОВАРИ', 'РЕКЛАМА', 'ЖУРНАЛ', 'ПОБЕДА'
  ]
  const languages = ['en', 'ru']

  const mockData = []
  for (let i = 0; i < 35; i++) {
    const isDaily = Math.random() > 0.7
    const score = 600 + Math.floor(Math.random() * 1500)
    const words = 8 + Math.floor(Math.random() * 25)
    const lang = languages[Math.floor(Math.random() * languages.length)]
    const letters = lang === 'ru'
      ? letterSets.filter(s => /[А-Я]/.test(s))[Math.floor(Math.random() * 5)]
      : letterSets.filter(s => /[A-Z]/.test(s))[Math.floor(Math.random() * 5)]

    mockData.push({
      id: `game-${i}`,
      letters,
      language: lang,
      isDaily,
      score,
      wordsFound: words,
      timeLimit: 180,
      createdAt: new Date(Date.now() - Math.random() * 60 * 24 * 60 * 60 * 1000).toISOString()
    })
  }

  return mockData.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
}

function calculateMockStats() {
  stats.value = {
    totalGames: games.value.length,
    bestScore: Math.max(...games.value.map(g => g.score), 0),
    totalWords: games.value.reduce((sum, g) => sum + (g.wordsFound || 0), 0),
    totalPoints: games.value.reduce((sum, g) => sum + (g.score || 0), 0)
  }
}

// Format date relative (e.g., "2h ago", "3d ago")
function formatRelativeTime(isoDate) {
  const now = new Date()
  const date = new Date(isoDate)
  const diffMs = now - date
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 60) return `${diffMins}m ${t('multiplayer.ago')}`
  if (diffHours < 24) return `${diffHours}h ${t('multiplayer.ago')}`
  if (diffDays === 1) return t('multiplayer.yesterday')
  if (diffDays < 7) return `${diffDays}d ${t('multiplayer.ago')}`

  return date.toLocaleDateString()
}

// Replay game with same settings
function replayGame(game) {
  // Navigate to game page with settings
  router.push({
    path: '/game',
    query: {
      letters: game.letters.length,
      time: game.timeLimit,
      lang: game.language
    }
  })
}

// Computed values
const totalPages = computed(() => Math.max(1, Math.ceil(totalGames.value / pageSize)))
const showing = computed(() => Math.min(pageSize, totalGames.value - (page.value - 1) * pageSize))

// Filter games by type (applied to mock data only, backend handles this)
const filteredGames = computed(() => {
  if (typeFilter.value === 'all') return games.value
  if (typeFilter.value === 'daily') return games.value.filter(g => g.isDaily)
  if (typeFilter.value === 'practice') return games.value.filter(g => !g.isDaily)
  return games.value
})

// Watch filters and reset to page 1
watch(typeFilter, () => {
  page.value = 1
  fetchGames()
})

// Watch page changes
watch(page, () => {
  fetchGames()
})

onMounted(() => {
  fetchGames()
})
</script>

<template>
  <div class="hist-wrap">
    <div class="hist-head">
      <div>
        <div class="page-eyebrow muted">History</div>
        <h1 class="display">Solo History</h1>
        <p class="muted">Your practice games and daily puzzles.</p>
      </div>
    </div>

    <!-- Stats Overview -->
    <div class="hist-stats">
      <div class="hist-stat-card">
        <div class="hist-stat-num accent-text">{{ stats.totalGames }}</div>
        <div class="hist-stat-lbl muted">Total games</div>
      </div>
      <div class="hist-stat-card">
        <div class="hist-stat-num mono">{{ stats.bestScore.toLocaleString() }}</div>
        <div class="hist-stat-lbl muted">Best score</div>
      </div>
      <div class="hist-stat-card">
        <div class="hist-stat-num mono">{{ stats.totalWords.toLocaleString() }}</div>
        <div class="hist-stat-lbl muted">Total words</div>
      </div>
      <div class="hist-stat-card">
        <div class="hist-stat-num mono">{{ stats.totalPoints.toLocaleString() }}</div>
        <div class="hist-stat-lbl muted">Total points</div>
      </div>
    </div>

    <!-- Filters -->
    <div class="hist-filters">
      <div class="hist-filter-group">
        <label class="hist-filter-label muted">Type</label>
        <div class="chip-toggle-group">
          <button
            :class="['chip-toggle', { 'is-active': typeFilter === 'all' }]"
            @click="typeFilter = 'all'"
          >
            All
          </button>
          <button
            :class="['chip-toggle', { 'is-active': typeFilter === 'daily' }]"
            @click="typeFilter = 'daily'"
          >
            Daily only
          </button>
          <button
            :class="['chip-toggle', { 'is-active': typeFilter === 'practice' }]"
            @click="typeFilter = 'practice'"
          >
            Practice only
          </button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="hist-loading">
      <p class="muted">{{ t('common.loading') }}</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredGames.length === 0" class="hist-empty">
      <p class="muted">No games found. Start playing to build your history!</p>
    </div>

    <!-- Table -->
    <div v-else class="hist-table-wrap">
      <table class="hist-table">
        <thead>
          <tr>
            <th>Letters</th>
            <th>Mode</th>
            <th class="hist-th-score">Score</th>
            <th class="hist-th-score">Words</th>
            <th>When</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="game in filteredGames" :key="game.id" class="hist-row">
            <td>
              <div class="hist-letters-cell">
                <span class="mono hist-letters">{{ game.letters }}</span>
                <span class="hist-lang-badge muted">{{ game.language?.toUpperCase() }}</span>
              </div>
            </td>
            <td>
              <span v-if="game.isDaily" class="hist-daily-badge">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M12 2v4m0 12v4M4.93 4.93l2.83 2.83m8.48 8.48l2.83 2.83M2 12h4m12 0h4M4.93 19.07l2.83-2.83m8.48-8.48l2.83-2.83"/>
                </svg>
                Daily
              </span>
              <span v-else class="muted hist-practice">Practice</span>
            </td>
            <td class="hist-score">
              <span class="mono">{{ game.score?.toLocaleString() || '—' }}</span>
            </td>
            <td class="hist-score">
              <span class="mono">{{ game.wordsFound || '—' }}</span>
            </td>
            <td>
              <span class="muted hist-when">{{ formatRelativeTime(game.createdAt) }}</span>
            </td>
            <td>
              <button class="btn btn-sm btn-ghost" @click="replayGame(game)">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/>
                  <path d="M21 3v5h-5"/>
                  <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/>
                  <path d="M3 21v-5h5"/>
                </svg>
                Replay
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <Pagination
      v-if="!loading && filteredGames.length > 0"
      :page="page"
      :total-pages="totalPages"
      :showing="showing"
      :total="totalGames"
      :page-size="pageSize"
      @update:page="page = $event"
    />
  </div>
</template>

<style scoped>
.hist-wrap {
  max-width: 1120px;
  margin: 0 auto;
  padding: var(--sp-8) var(--sp-5);
}

.hist-head {
  margin-bottom: var(--sp-6);
}

.page-eyebrow {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 1.2px;
  font-weight: 600;
  margin-bottom: var(--sp-2);
}

h1.display {
  font-family: var(--font-display);
  font-size: 42px;
  font-weight: 700;
  margin-bottom: var(--sp-2);
  color: var(--navy);
}

/* Stats Grid */
.hist-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: var(--sp-4);
  margin-bottom: var(--sp-6);
}

.hist-stat-card {
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 16px;
  padding: var(--sp-4);
  text-align: center;
}

.hist-stat-num {
  font-family: var(--font-mono);
  font-size: 32px;
  font-weight: 700;
  line-height: 1;
  margin-bottom: var(--sp-2);
}

.hist-stat-lbl {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 600;
}

/* Filters */
.hist-filters {
  display: flex;
  gap: var(--sp-4);
  margin-bottom: var(--sp-5);
  flex-wrap: wrap;
}

.hist-filter-group {
  display: flex;
  flex-direction: column;
  gap: var(--sp-2);
  flex: 1;
}

.hist-filter-label {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 600;
}

/* Table */
.hist-table-wrap {
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 20px;
  overflow: hidden;
  margin-bottom: var(--sp-4);
}

.hist-table {
  width: 100%;
  border-collapse: collapse;
}

.hist-table thead {
  background: var(--bg-card);
  border-bottom: 1px solid var(--border-hairline);
}

.hist-table th {
  padding: var(--sp-3) var(--sp-4);
  text-align: left;
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 600;
  color: var(--fg-muted);
}

.hist-table th.hist-th-score {
  text-align: right;
}

.hist-table tbody tr {
  border-bottom: 1px solid var(--border-hairline);
  transition: background var(--dur-fast);
}

.hist-table tbody tr:last-child {
  border-bottom: none;
}

.hist-table tbody tr:hover {
  background: color-mix(in oklab, var(--bg-card) 50%, transparent);
}

.hist-table td {
  padding: var(--sp-3) var(--sp-4);
  font-size: 14px;
  color: var(--fg-primary);
}

.hist-score {
  text-align: right;
}

/* Letters Cell */
.hist-letters-cell {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}

.hist-letters {
  font-size: 15px;
  font-weight: 700;
  letter-spacing: 1.5px;
  color: var(--navy);
}

.hist-lang-badge {
  font-size: 10px;
  font-weight: 600;
  font-family: var(--font-mono);
  padding: 2px 6px;
  background: var(--bg-card);
  border-radius: 4px;
}

/* Daily Badge */
.hist-daily-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
  background: linear-gradient(135deg, var(--accent) 0%, var(--cocoa, #7a4a2b) 100%);
  color: var(--milk);
  border: 1px solid color-mix(in oklab, var(--accent) 80%, transparent);
}

.hist-daily-badge svg {
  opacity: 0.9;
}

.hist-practice {
  font-size: 12px;
  font-weight: 600;
}

.hist-when {
  font-size: 12px;
}

/* Loading/Empty States */
.hist-loading,
.hist-empty {
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 20px;
  padding: var(--sp-8);
  text-align: center;
  margin-bottom: var(--sp-4);
}

/* Responsive */
@media (max-width: 820px) {
  .hist-table-wrap {
    overflow-x: auto;
  }

  .hist-table {
    min-width: 600px;
  }
}

@media (max-width: 720px) {
  h1.display {
    font-size: 32px;
  }

  .hist-stats {
    grid-template-columns: repeat(2, 1fr);
  }

  .hist-filters {
    flex-direction: column;
  }
}
</style>
