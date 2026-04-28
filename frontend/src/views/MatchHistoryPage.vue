<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../stores/userStore'
import Pagination from '../components/Pagination.vue'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()

const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// Filters
const resultFilter = ref('all') // all, won, lost, tie
const opponentFilter = ref('all')

// Pagination
const page = ref(1)
const pageSize = 15

// Data
const allMatches = ref([])
const loading = ref(false)

// Stats (calculated from all matches, not just filtered)
const stats = computed(() => {
  const wins = allMatches.value.filter(m => m.result === 'won').length
  const ties = allMatches.value.filter(m => m.result === 'tie').length
  const losses = allMatches.value.filter(m => m.result === 'lost').length
  const totalPoints = allMatches.value.reduce((sum, m) => sum + m.myScore, 0)
  return { wins, ties, losses, totalPoints }
})

// Unique opponents for filter dropdown
const opponents = computed(() => {
  const unique = new Set()
  allMatches.value.forEach(m => {
    if (m.opponentUsername) unique.add(m.opponentUsername)
  })
  return ['all', ...Array.from(unique).sort()]
})

// Filtered matches
const filteredMatches = computed(() => {
  let filtered = allMatches.value

  if (resultFilter.value !== 'all') {
    filtered = filtered.filter(m => m.result === resultFilter.value)
  }

  if (opponentFilter.value !== 'all') {
    filtered = filtered.filter(m => m.opponentUsername === opponentFilter.value)
  }

  return filtered
})

// Paginated slice
const paginatedMatches = computed(() => {
  const start = (page.value - 1) * pageSize
  const end = start + pageSize
  return filteredMatches.value.slice(start, end)
})

const totalPages = computed(() => {
  return Math.max(1, Math.ceil(filteredMatches.value.length / pageSize))
})

// Fetch matches from backend
async function fetchMatches() {
  if (!userStore.userId) {
    loading.value = false
    return
  }

  loading.value = true
  try {
    const params = new URLSearchParams({
      user_id: userStore.userId,
      page: '1',
      per_page: '100' // Get all matches for client-side filtering
    })

    const response = await fetch(`${apiUrl}/api/v1/sessions/all?${params}`)

    if (!response.ok) {
      throw new Error('Failed to fetch match history')
    }

    const data = await response.json()

    // Filter only completed multiplayer sessions
    const multiplayerSessions = (data.sessions || []).filter(s =>
      s.results && s.results.length >= 2
    )

    // Convert to match format
    allMatches.value = multiplayerSessions.map(session => {
      const myResult = session.results.find(r => r.user_id === userStore.userId)
      const opponentResult = session.results.find(r => r.user_id !== userStore.userId)

      if (!myResult || !opponentResult) return null

      const myScore = myResult.score || 0
      const opponentScore = opponentResult.score || 0
      const delta = myScore - opponentScore

      let result = 'tie'
      if (delta > 0) result = 'won'
      else if (delta < 0) result = 'lost'

      return {
        id: session.id,
        result,
        opponentUsername: opponentResult.player_name || 'Unknown',
        letters: session.letters,
        myScore,
        opponentScore,
        delta,
        createdAt: myResult.played_at || session.created_at
      }
    }).filter(Boolean)

  } catch (error) {
    console.error('Error fetching match history:', error)
    userStore.showToast('Failed to load match history', 'error')
    allMatches.value = []
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

// Reset page when filters change
watch([resultFilter, opponentFilter], () => {
  page.value = 1
})

onMounted(() => {
  fetchMatches()
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
          <div class="page-eyebrow">{{ t('history.match.title') }}</div>
          <h1 class="page-title-display">{{ t('history.match.title') }}.</h1>
        </div>
        <div class="hist-stats">
          <div class="hist-stat">
            <span class="mono">{{ stats.wins }}</span>
            <span class="lbl">{{ t('history.match.won') }}</span>
          </div>
          <div class="hist-stat">
            <span class="mono">{{ stats.ties }}</span>
            <span class="lbl">{{ t('history.match.tied') }}</span>
          </div>
          <div class="hist-stat">
            <span class="mono">{{ stats.losses }}</span>
            <span class="lbl">{{ t('history.match.lost') }}</span>
          </div>
          <div class="hist-stat">
            <span class="mono">{{ stats.totalPoints.toLocaleString() }}</span>
            <span class="lbl">{{ t('history.match.totalPts') }}</span>
          </div>
        </div>
      </header>

      <!-- Filters -->
      <div class="hist-filters">
        <div class="hist-filter-grp">
          <span class="hist-filter-lbl muted">{{ t('history.match.result') }}</span>
          <div class="checkbox-row">
            <button
              class="chip-toggle"
              :data-active="resultFilter === 'all'"
              @click="resultFilter = 'all'"
            >
              {{ t('history.match.all') }}
            </button>
            <button
              class="chip-toggle"
              :data-active="resultFilter === 'won'"
              @click="resultFilter = 'won'"
            >
              {{ t('history.match.won_filter') }}
            </button>
            <button
              class="chip-toggle"
              :data-active="resultFilter === 'lost'"
              @click="resultFilter = 'lost'"
            >
              {{ t('history.match.lost_filter') }}
            </button>
            <button
              class="chip-toggle"
              :data-active="resultFilter === 'tie'"
              @click="resultFilter = 'tie'"
            >
              {{ t('history.match.tied_filter') }}
            </button>
          </div>
        </div>
        <div class="hist-filter-grp">
          <span class="hist-filter-lbl muted">{{ t('history.match.opponent') }}</span>
          <select v-model="opponentFilter" class="hist-select">
            <option v-for="opp in opponents" :key="opp" :value="opp">
              {{ opp === 'all' ? t('history.match.anyone') : opp }}
            </option>
          </select>
        </div>
      </div>

      <!-- Loading state -->
      <div v-if="loading" class="hist-empty">
        <div class="muted">{{ t('history.match.loading') }}</div>
      </div>

      <!-- Empty state -->
      <div v-else-if="filteredMatches.length === 0" class="hist-empty">
        <div class="muted">
          {{ allMatches.length === 0 ? t('history.match.noMatches') : t('history.match.noMatchesFiltered') }}
        </div>
        <button
          v-if="resultFilter !== 'all' || opponentFilter !== 'all'"
          class="btn btn--soft btn--sm"
          @click="resultFilter = 'all'; opponentFilter = 'all'"
        >
          {{ t('history.match.clearFilters') }}
        </button>
      </div>

      <!-- Table -->
      <div v-else class="hist-table">
        <!-- Header -->
        <div class="hist-row hist-row--head">
          <span class="hr-col-result">{{ t('history.match.result') }}</span>
          <span class="hr-col-with">{{ t('history.match.opponent') }}</span>
          <span class="hr-col-letters">{{ t('history.match.letters') }}</span>
          <span class="hr-col-score">{{ t('history.match.score') }}</span>
          <span class="hr-col-delta">{{ t('history.match.delta') }}</span>
          <span class="hr-col-date">{{ t('history.match.when') }}</span>
          <span class="hr-col-cta" />
        </div>

        <!-- Rows -->
        <div
          v-for="m in paginatedMatches"
          :key="m.id"
          :class="['hist-row', `rm-${m.result}`]"
        >
          <span :class="['hr-col-result', 'hr-result', `rm-${m.result}`]">
            {{ m.result === 'won' ? 'W' : m.result === 'lost' ? 'L' : 'T' }}
          </span>
          <span class="hr-col-with">
            <span class="hr-avatar">{{ m.opponentUsername[0].toUpperCase() }}</span>
            <span>{{ m.opponentUsername }}</span>
          </span>
          <span class="hr-col-letters mono">{{ m.letters }}</span>
          <span class="hr-col-score mono">
            {{ m.myScore.toLocaleString() }} <span class="muted">·</span> {{ m.opponentScore.toLocaleString() }}
          </span>
          <span :class="['hr-col-delta', 'mono', `rm-${m.result}`]">
            <template v-if="m.result === 'won'">+{{ Math.abs(m.delta).toLocaleString() }}</template>
            <template v-else-if="m.result === 'lost'">−{{ Math.abs(m.delta).toLocaleString() }}</template>
            <template v-else>0</template>
          </span>
          <span class="hr-col-date muted">{{ formatDate(m.createdAt) }}</span>
          <span class="hr-col-cta">
            <button class="btn btn--ghost btn--sm" @click="router.push(`/results/${m.id}`)">
              {{ t('history.match.view') }}
            </button>
          </span>
        </div>
      </div>

      <!-- Pagination -->
      <Pagination
        v-if="filteredMatches.length > 0"
        :page="page"
        :total-pages="totalPages"
        :showing="paginatedMatches.length"
        :total="filteredMatches.length"
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

.hist-select {
  padding: 6px 10px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-subtle);
  background: var(--bg-surface);
  font-family: var(--font-body);
  font-size: 13px;
  color: var(--fg-primary);
  cursor: pointer;
}

.hist-select:focus {
  outline: 0;
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-soft);
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
  grid-template-columns: 60px 160px 140px 160px 80px 1fr 80px;
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

.hr-result {
  display: inline-grid;
  place-items: center;
  width: 28px;
  height: 28px;
  border-radius: 8px;
  font-family: var(--font-mono);
  font-weight: 700;
  font-size: 12px;
  color: var(--milk);
}

.hr-result.rm-won {
  background: var(--success);
}

.hr-result.rm-lost {
  background: var(--danger);
}

.hr-result.rm-tie {
  background: var(--fg-faint);
}

.hr-col-with {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  font-weight: 600;
  color: var(--fg-primary);
}

.hr-avatar {
  width: 26px;
  height: 26px;
  border-radius: 999px;
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  display: grid;
  place-items: center;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 11px;
  flex-shrink: 0;
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

.hr-col-score {
  color: var(--fg-secondary);
}

.hr-col-delta {
  font-weight: 700;
  text-align: right;
}

.hr-col-delta.rm-won {
  color: var(--success);
}

.hr-col-delta.rm-lost {
  color: var(--danger);
}

.hr-col-delta.rm-tie {
  color: var(--fg-muted);
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
    grid-template-columns: 36px 1fr auto;
    grid-template-areas:
      "res with cta"
      "res letters score"
      ". delta date";
    row-gap: 4px;
    padding: var(--sp-3);
  }

  .hist-row--head {
    display: none;
  }

  .hr-col-result {
    grid-area: res;
  }

  .hr-col-with {
    grid-area: with;
  }

  .hr-col-letters {
    grid-area: letters;
    font-size: 12px;
  }

  .hr-col-score {
    grid-area: score;
    text-align: right;
    font-size: 11px;
  }

  .hr-col-delta {
    grid-area: delta;
    font-size: 12px;
    text-align: left;
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
