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
const resultFilter = ref('all') // all, won, lost, tie
const opponentFilter = ref('all')

// Pagination
const page = ref(1)
const pageSize = 15

// Data
const matches = ref([])
const totalMatches = ref(0)
const loading = ref(false)

// Stats
const stats = ref({
  wins: 0,
  ties: 0,
  losses: 0,
  totalPoints: 0
})

// Unique opponents for filter dropdown
const opponents = computed(() => {
  const unique = new Set()
  matches.value.forEach(m => {
    if (m.opponentUsername) {
      unique.add(m.opponentUsername)
    }
  })
  return Array.from(unique).sort()
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
      page: page.value.toString(),
      per_page: pageSize.toString()
    })

    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/api/v1/sessions/all?${params}`)

    if (!response.ok) {
      throw new Error('Failed to fetch match history')
    }

    const data = await response.json()

    // Filter only multiplayer sessions (with multiple results)
    const multiplayerSessions = (data.sessions || []).filter(s => s.results && s.results.length > 1)

    // Convert to match format
    matches.value = multiplayerSessions.map(session => {
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
        opponentUsername: opponentResult.player_name,
        letters: session.letters,
        myScore,
        opponentScore,
        delta,
        createdAt: session.created_at
      }
    }).filter(Boolean)

    // Apply filters
    let filtered = matches.value
    if (resultFilter.value !== 'all') {
      filtered = filtered.filter(m => m.result === resultFilter.value)
    }
    if (opponentFilter.value !== 'all') {
      filtered = filtered.filter(m => m.opponentUsername === opponentFilter.value)
    }

    matches.value = filtered
    totalMatches.value = filtered.length

    // Calculate stats
    calculateMockStats()
  } catch (error) {
    console.error('Error fetching match history:', error)
    userStore.showToast('Failed to load match history', 'error')

    // Fallback to mock data
    matches.value = generateMockMatches()
    totalMatches.value = matches.value.length
    calculateMockStats()
  } finally {
    loading.value = false
  }
}

// Mock data generator (fallback)
function generateMockMatches() {
  const results = ['won', 'lost', 'tie']
  const opponents = ['alice_m', 'bob_k', 'charlie_d', 'diana_s', 'evan_r']
  const letterSets = ['ABCDEFG', 'TESTING', 'EXAMPLE', 'WORDGAM', 'PLAYNOW']

  const mockData = []
  for (let i = 0; i < 25; i++) {
    const result = results[Math.floor(Math.random() * results.length)]
    const myScore = 800 + Math.floor(Math.random() * 1200)
    const theirScore = result === 'won'
      ? myScore - Math.floor(Math.random() * 300) - 50
      : result === 'lost'
      ? myScore + Math.floor(Math.random() * 300) + 50
      : myScore + Math.floor(Math.random() * 40) - 20

    mockData.push({
      id: `match-${i}`,
      result,
      opponentUsername: opponents[Math.floor(Math.random() * opponents.length)],
      letters: letterSets[Math.floor(Math.random() * letterSets.length)],
      myScore,
      opponentScore: theirScore,
      delta: myScore - theirScore,
      createdAt: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString()
    })
  }

  return mockData.sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
}

function calculateMockStats() {
  stats.value = {
    wins: matches.value.filter(m => m.result === 'won').length,
    ties: matches.value.filter(m => m.result === 'tie').length,
    losses: matches.value.filter(m => m.result === 'lost').length,
    totalPoints: matches.value.reduce((sum, m) => sum + (m.myScore || 0), 0)
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

// View match details
function viewMatch(matchId) {
  router.push(`/results/${matchId}`)
}

// Computed values
const totalPages = computed(() => Math.max(1, Math.ceil(totalMatches.value / pageSize)))
const showing = computed(() => Math.min(pageSize, totalMatches.value - (page.value - 1) * pageSize))

// Watch filters and reset to page 1
watch([resultFilter, opponentFilter], () => {
  page.value = 1
  fetchMatches()
})

// Watch page changes
watch(page, () => {
  fetchMatches()
})

onMounted(() => {
  fetchMatches()
})
</script>

<template>
  <div class="hist-wrap">
    <div class="hist-head">
      <div>
        <div class="page-eyebrow muted">History</div>
        <h1 class="display">Match History</h1>
        <p class="muted">Head-to-head games against friends.</p>
      </div>
    </div>

    <!-- Stats Overview -->
    <div class="hist-stats">
      <div class="hist-stat-card">
        <div class="hist-stat-num accent-text">{{ stats.wins }}</div>
        <div class="hist-stat-lbl muted">Wins</div>
      </div>
      <div class="hist-stat-card">
        <div class="hist-stat-num" style="color: var(--fg-muted)">{{ stats.ties }}</div>
        <div class="hist-stat-lbl muted">Ties</div>
      </div>
      <div class="hist-stat-card">
        <div class="hist-stat-num" style="color: var(--danger)">{{ stats.losses }}</div>
        <div class="hist-stat-lbl muted">Losses</div>
      </div>
      <div class="hist-stat-card">
        <div class="hist-stat-num mono">{{ stats.totalPoints.toLocaleString() }}</div>
        <div class="hist-stat-lbl muted">Total points</div>
      </div>
    </div>

    <!-- Filters -->
    <div class="hist-filters">
      <div class="hist-filter-group">
        <label class="hist-filter-label muted">Result</label>
        <div class="chip-toggle-group">
          <button
            :class="['chip-toggle', { 'is-active': resultFilter === 'all' }]"
            @click="resultFilter = 'all'"
          >
            All
          </button>
          <button
            :class="['chip-toggle', { 'is-active': resultFilter === 'won' }]"
            @click="resultFilter = 'won'"
          >
            Won
          </button>
          <button
            :class="['chip-toggle', { 'is-active': resultFilter === 'lost' }]"
            @click="resultFilter = 'lost'"
          >
            Lost
          </button>
          <button
            :class="['chip-toggle', { 'is-active': resultFilter === 'tie' }]"
            @click="resultFilter = 'tie'"
          >
            Tie
          </button>
        </div>
      </div>

      <div class="hist-filter-group">
        <label class="hist-filter-label muted">Opponent</label>
        <select v-model="opponentFilter" class="input">
          <option value="all">All opponents</option>
          <option v-for="opp in opponents" :key="opp" :value="opp">{{ opp }}</option>
        </select>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="hist-loading">
      <p class="muted">{{ t('common.loading') }}</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="matches.length === 0" class="hist-empty">
      <p class="muted">No matches found. Play some multiplayer games to see them here!</p>
    </div>

    <!-- Table -->
    <div v-else class="hist-table-wrap">
      <table class="hist-table">
        <thead>
          <tr>
            <th>Result</th>
            <th>Opponent</th>
            <th>Letters</th>
            <th class="hist-th-score">Your score</th>
            <th class="hist-th-score">Their score</th>
            <th class="hist-th-delta">Delta</th>
            <th>When</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="match in matches" :key="match.id" class="hist-row">
            <td>
              <span
                :class="[
                  'hist-result-badge',
                  `hist-result-badge--${match.result}`
                ]"
              >
                {{ match.result === 'won' ? '✓ Won' : match.result === 'lost' ? '✗ Lost' : '− Tie' }}
              </span>
            </td>
            <td>
              <div class="hist-opponent">
                <div class="hist-avatar">{{ match.opponentUsername?.[0]?.toUpperCase() || '?' }}</div>
                <span>{{ match.opponentUsername || 'Unknown' }}</span>
              </div>
            </td>
            <td>
              <span class="mono hist-letters">{{ match.letters }}</span>
            </td>
            <td class="hist-score">
              <span class="mono">{{ match.myScore?.toLocaleString() || '—' }}</span>
            </td>
            <td class="hist-score">
              <span class="mono">{{ match.opponentScore?.toLocaleString() || '—' }}</span>
            </td>
            <td class="hist-delta">
              <span
                v-if="match.delta !== undefined"
                :class="[
                  'mono',
                  'hist-delta-num',
                  match.delta > 0 ? 'hist-delta--pos' : match.delta < 0 ? 'hist-delta--neg' : ''
                ]"
              >
                {{ match.delta > 0 ? '+' : '' }}{{ match.delta }}
              </span>
              <span v-else class="mono">—</span>
            </td>
            <td>
              <span class="muted hist-when">{{ formatRelativeTime(match.createdAt) }}</span>
            </td>
            <td>
              <button class="btn btn-sm btn-ghost" @click="viewMatch(match.id)">
                View
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <Pagination
      v-if="!loading && matches.length > 0"
      :page="page"
      :total-pages="totalPages"
      :showing="showing"
      :total="totalMatches"
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
  min-width: 200px;
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

.hist-table th.hist-th-score,
.hist-table th.hist-th-delta {
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

.hist-delta {
  text-align: right;
}

/* Result Badge */
.hist-result-badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
  font-family: var(--font-mono);
}

.hist-result-badge--won {
  background: var(--success-soft, color-mix(in oklab, var(--success) 12%, transparent));
  color: var(--success);
  border: 1px solid color-mix(in oklab, var(--success) 30%, transparent);
}

.hist-result-badge--lost {
  background: color-mix(in oklab, var(--danger) 12%, transparent);
  color: var(--danger);
  border: 1px solid color-mix(in oklab, var(--danger) 30%, transparent);
}

.hist-result-badge--tie {
  background: var(--bg-card);
  color: var(--fg-muted);
  border: 1px solid var(--border-subtle);
}

/* Opponent */
.hist-opponent {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}

.hist-avatar {
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
  border: 2px solid var(--milk);
  box-shadow: 0 0 0 1px var(--border-default);
  flex-shrink: 0;
}

.hist-letters {
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--navy);
}

.hist-delta-num {
  font-size: 14px;
  font-weight: 700;
}

.hist-delta--pos {
  color: var(--success);
}

.hist-delta--neg {
  color: var(--danger);
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
    min-width: 700px;
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
