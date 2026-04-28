<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../stores/userStore'
import { useGameStore } from '../stores/gameStore'
import Pagination from '../components/Pagination.vue'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()
const gameStore = useGameStore()

const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// Filters
const typeFilter = ref('all') // all, daily, practice

// Pagination
const page = ref(1)
const pageSize = 15

// Data
const allGames = ref([])
const loading = ref(false)

// Stats (calculated from all games, not just filtered)
const stats = computed(() => {
  const totalGames = allGames.value.length
  const bestScore = allGames.value.length > 0
    ? Math.max(...allGames.value.map(g => g.score))
    : 0
  const totalWords = allGames.value.reduce((sum, g) => sum + g.wordsFound, 0)
  const totalPoints = allGames.value.reduce((sum, g) => sum + g.score, 0)
  return { totalGames, bestScore, totalWords, totalPoints }
})

// Filtered games
const filteredGames = computed(() => {
  let filtered = allGames.value

  if (typeFilter.value === 'daily') {
    filtered = filtered.filter(g => g.isDaily)
  } else if (typeFilter.value === 'practice') {
    filtered = filtered.filter(g => !g.isDaily)
  }

  return filtered
})

// Paginated slice
const paginatedGames = computed(() => {
  const start = (page.value - 1) * pageSize
  const end = start + pageSize
  return filteredGames.value.slice(start, end)
})

const totalPages = computed(() => {
  return Math.max(1, Math.ceil(filteredGames.value.length / pageSize))
})

// Fetch games from backend
async function fetchGames() {
  if (!userStore.userId) {
    loading.value = false
    return
  }

  loading.value = true
  try {
    const params = new URLSearchParams({
      user_id: userStore.userId,
      page: '1',
      per_page: '100' // Get all games for client-side filtering
    })

    const response = await fetch(`${apiUrl}/api/v1/sessions/all?${params}`)

    if (!response.ok) {
      throw new Error('Failed to fetch solo history')
    }

    const data = await response.json()

    // Filter only solo sessions (single result with current user)
    const soloSessions = (data.sessions || []).filter(s =>
      s.results && s.results.length === 1 && s.results[0].user_id === userStore.userId
    )

    // Convert to game format
    allGames.value = soloSessions.map(session => {
      const result = session.results[0]

      return {
        id: session.id,
        letters: session.letters,
        language: session.language,
        isDaily: false, // TODO: Backend doesn't support this yet
        score: result.score || 0,
        wordsFound: result.word_count || 0,
        timeLimit: session.time_limit,
        createdAt: result.played_at || session.created_at
      }
    })

  } catch (error) {
    console.error('Error fetching solo history:', error)
    userStore.showToast('Failed to load solo history', 'error')
    allGames.value = []
  } finally {
    loading.value = false
  }
}

// Format date helper
function formatDate(dateStr) {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now - date
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 1) return 'just now'
  if (diffMins < 60) return `${diffMins}m ago`
  if (diffHours < 24) return `${diffHours}h ago`
  if (diffDays === 1) return 'yesterday'
  if (diffDays < 7) return `${diffDays}d ago`
  if (diffDays < 30) return `${Math.floor(diffDays / 7)}w ago`
  return date.toLocaleDateString()
}

// Replay a game with same letters
function replayGame(game) {
  gameStore.startGame(game.timeLimit, game.letters.length, game.language)
  router.push('/game')
}

// Reset page when filter changes
watch(typeFilter, () => {
  page.value = 1
})

onMounted(() => {
  fetchGames()
})
</script>

<template>
  <div class="page">
    <div class="shell hist-wrap">
      <!-- Header -->
      <header class="page-head">
        <div>
          <button class="btn btn--ghost btn--sm" @click="router.push('/play')" style="margin-bottom: 8px">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M19 12H5M12 19l-7-7 7-7"/>
            </svg>
            {{ t('history.backToPlay') }}
          </button>
          <div class="page-eyebrow">{{ t('history.solo.title') }}</div>
          <h1 class="page-title-display">{{ t('history.solo.title') }}.</h1>
        </div>
        <div class="hist-stats">
          <div class="hist-stat">
            <span class="mono">{{ stats.totalGames }}</span>
            <span class="lbl">{{ t('history.solo.games') }}</span>
          </div>
          <div class="hist-stat">
            <span class="mono">{{ stats.bestScore.toLocaleString() }}</span>
            <span class="lbl">{{ t('history.solo.best') }}</span>
          </div>
          <div class="hist-stat">
            <span class="mono">{{ stats.totalWords }}</span>
            <span class="lbl">{{ t('history.solo.words') }}</span>
          </div>
          <div class="hist-stat">
            <span class="mono">{{ stats.totalPoints.toLocaleString() }}</span>
            <span class="lbl">{{ t('history.solo.totalPts') }}</span>
          </div>
        </div>
      </header>

      <!-- Filters -->
      <div class="hist-filters">
        <div class="hist-filter-grp">
          <span class="hist-filter-lbl muted">{{ t('history.solo.type') }}</span>
          <div class="checkbox-row">
            <button
              class="chip-toggle"
              :data-active="typeFilter === 'all'"
              @click="typeFilter = 'all'"
            >
              {{ t('history.solo.all') }}
            </button>
            <button
              class="chip-toggle"
              :data-active="typeFilter === 'daily'"
              @click="typeFilter = 'daily'"
            >
              {{ t('history.solo.dailyOnly') }}
            </button>
            <button
              class="chip-toggle"
              :data-active="typeFilter === 'practice'"
              @click="typeFilter = 'practice'"
            >
              {{ t('history.solo.practiceOnly') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Loading state -->
      <div v-if="loading" class="hist-empty">
        <div class="muted">{{ t('history.solo.loading') }}</div>
      </div>

      <!-- Empty state -->
      <div v-else-if="filteredGames.length === 0" class="hist-empty">
        <div class="muted">
          {{ allGames.length === 0 ? t('history.solo.noGames') : t('history.solo.noGamesFiltered') }}
        </div>
        <button
          v-if="typeFilter !== 'all'"
          class="btn btn--soft btn--sm"
          @click="typeFilter = 'all'"
        >
          {{ t('history.solo.clearFilter') }}
        </button>
      </div>

      <!-- Table -->
      <div v-else class="hist-table">
        <!-- Header -->
        <div class="hist-row hist-row--head">
          <span class="hr-col-letters">{{ t('history.solo.letters') }}</span>
          <span class="hr-col-with">{{ t('history.solo.mode') }}</span>
          <span class="hr-col-score">{{ t('history.solo.score') }}</span>
          <span class="hr-col-delta">{{ t('history.solo.words') }}</span>
          <span class="hr-col-date">{{ t('history.solo.when') }}</span>
          <span class="hr-col-cta" />
        </div>

        <!-- Rows -->
        <div
          v-for="g in paginatedGames"
          :key="g.id"
          class="hist-row"
        >
          <span class="hr-col-letters mono">
            {{ g.letters }}
            <span v-if="g.isDaily" class="hr-daily-tag">DAILY</span>
          </span>
          <span class="hr-col-with muted">
            {{ g.letters.length }}L · {{ g.timeLimit }}s · {{ g.language.toUpperCase() }}
          </span>
          <span class="hr-col-score mono">{{ g.score.toLocaleString() }}</span>
          <span class="hr-col-delta mono">{{ g.wordsFound }}</span>
          <span class="hr-col-date muted">{{ formatDate(g.createdAt) }}</span>
          <span class="hr-col-cta">
            <button class="btn btn--soft btn--sm" @click="replayGame(g)">
              <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M13 2L4 14h7l-1 8 9-12h-7l1-8z"/>
              </svg>
              {{ t('history.solo.replay') }}
            </button>
          </span>
        </div>
      </div>

      <!-- Pagination -->
      <Pagination
        v-if="filteredGames.length > 0"
        :page="page"
        :total-pages="totalPages"
        :showing="paginatedGames.length"
        :total="filteredGames.length"
        @update:page="page = $event"
      />
    </div>
  </div>
</template>

<style scoped>
/* Import history styles */
.hist-wrap {
  padding-top: var(--sp-8);
  padding-bottom: var(--sp-12);
}

/* ===== Header stats ==================================================== */
.hist-stats {
  display: flex;
  gap: var(--sp-3);
  flex-wrap: wrap;
}

.hist-stat {
  display: flex;
  flex-direction: column;
  padding: var(--sp-3) var(--sp-4);
  background: var(--bg-surface);
  border: 1px solid var(--border-hairline);
  border-radius: var(--radius-md);
  min-width: 80px;
}

.hist-stat .mono {
  font-family: var(--font-mono);
  font-weight: 700;
  font-size: 18px;
  color: var(--fg-primary);
}

.hist-stat .lbl {
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  color: var(--fg-muted);
  font-weight: 600;
  margin-top: 2px;
}

/* ===== Filters ======================================================== */
.hist-filters {
  display: flex;
  gap: var(--sp-5);
  align-items: flex-end;
  flex-wrap: wrap;
  margin-bottom: var(--sp-4);
  padding: var(--sp-3) var(--sp-4);
  background: var(--bg-card);
  border-radius: var(--radius-md);
}

.hist-filter-grp {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}

.hist-filter-lbl {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  font-weight: 600;
}

/* ===== Empty state ==================================================== */
.hist-empty {
  padding: var(--sp-8);
  text-align: center;
  background: var(--bg-card);
  border-radius: var(--radius-md);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--sp-3);
  margin-bottom: var(--sp-5);
}

/* ===== Table ========================================================== */
.hist-table {
  background: var(--bg-surface);
  border: 1px solid var(--border-hairline);
  border-radius: var(--radius-lg);
  overflow: hidden;
  margin-bottom: var(--sp-5);
}

.hist-row {
  display: grid;
  grid-template-columns: 140px 160px 100px 80px 1fr 100px;
  gap: var(--sp-3);
  align-items: center;
  padding: 14px var(--sp-4);
  border-bottom: 1px solid var(--border-hairline);
  font-size: 13px;
  transition: background var(--dur-fast);
}

.hist-row:last-child {
  border-bottom: 0;
}

.hist-row:not(.hist-row--head):hover {
  background: var(--bg-card);
}

.hist-row--head {
  background: var(--bg-card);
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  font-weight: 700;
  color: var(--fg-muted);
  padding: 10px var(--sp-4);
}

.hr-col-letters {
  font-weight: 700;
  letter-spacing: 1px;
  font-size: 13px;
  color: var(--fg-primary);
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.hr-daily-tag {
  font-family: var(--font-body);
  font-size: 9px;
  font-weight: 700;
  letter-spacing: 1px;
  padding: 2px 5px;
  border-radius: 4px;
  background: var(--accent);
  color: var(--milk);
}

.hr-col-with {
  color: var(--fg-secondary);
  font-size: 12px;
}

.hr-col-score {
  color: var(--fg-primary);
  font-weight: 600;
}

.hr-col-delta {
  font-weight: 700;
  color: var(--fg-primary);
}

.hr-col-date {
  font-size: 12px;
}

.hr-col-cta {
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 820px) {
  .hist-row {
    grid-template-columns: 1fr auto;
    grid-template-areas:
      "letters cta"
      "mode score"
      "delta date";
    row-gap: 4px;
    padding: var(--sp-3);
  }

  .hist-row--head {
    display: none;
  }

  .hr-col-letters {
    grid-area: letters;
    font-size: 12px;
  }

  .hr-col-with {
    grid-area: mode;
    font-size: 11px;
  }

  .hr-col-score {
    grid-area: score;
    text-align: right;
    font-size: 13px;
  }

  .hr-col-delta {
    grid-area: delta;
    font-size: 12px;
  }

  .hr-col-date {
    grid-area: date;
    text-align: right;
  }

  .hr-col-cta {
    grid-area: cta;
  }
}
</style>
