<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../stores/userStore'
import { useGameStore } from '../stores/gameStore'
import Modal from '../components/Modal.vue'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()
const gameStore = useGameStore()

// State
const showNewChallengeModal = ref(false)
const activeChallenges = ref([])
const loadingChallenges = ref(false)
const friends = ref([])
const loadingFriends = ref(false)
const selectedFriend = ref(null)
const searchFilter = ref('')
const modalStep = ref('pick') // 'pick' | 'sent'

// Recent completed matches
const recentMatches = ref([])
const loadingRecentMatches = ref(false)

// Computed
const yourTurnCount = computed(() => {
  return activeChallenges.value.filter(c => {
    const userResult = c.results?.find(r => r.user_id === userStore.userId)
    return !userResult && c.type === 'invited'
  }).length
})

const filteredFriends = computed(() => {
  if (!searchFilter.value) return friends.value
  return friends.value.filter(f =>
    f.username.toLowerCase().includes(searchFilter.value.toLowerCase())
  )
})

// Functions
async function loadActiveChallenges() {
  if (!userStore.isAuthenticated || !userStore.userId) return

  loadingChallenges.value = true
  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/api/v1/sessions/all?user_id=${userStore.userId}&page=1&per_page=50`)

    if (response.ok) {
      const data = await response.json()
      // Filter only active challenges (not finished)
      activeChallenges.value = (data.sessions || []).filter(c => {
        const hasAllPlayed = c.results && c.results.length >= 2
        return !hasAllPlayed
      })
    }
  } catch (error) {
    console.error('Failed to load challenges:', error)
  } finally {
    loadingChallenges.value = false
  }
}

async function loadRecentMatches() {
  if (!userStore.isAuthenticated || !userStore.userId) return

  loadingRecentMatches.value = true
  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/api/v1/sessions/all?user_id=${userStore.userId}&page=1&per_page=5`)

    if (response.ok) {
      const data = await response.json()
      // Filter completed multiplayer sessions
      const completed = (data.sessions || []).filter(s => {
        return s.results && s.results.length >= 2
      })

      // Transform to match format
      recentMatches.value = completed.slice(0, 5).map(session => {
        const myResult = session.results.find(r => r.user_id === userStore.userId)
        const opponentResult = session.results.find(r => r.user_id !== userStore.userId)

        if (!myResult || !opponentResult) return null

        const delta = myResult.score - opponentResult.score
        let result = 'tie'
        if (delta > 0) result = 'won'
        else if (delta < 0) result = 'lost'

        return {
          id: session.id,
          with: opponentResult.player_name || 'Opponent',
          letters: session.letters,
          yourScore: myResult.score,
          theirScore: opponentResult.score,
          result,
          date: formatDate(myResult.played_at || session.created_at)
        }
      }).filter(Boolean)
    }
  } catch (error) {
    console.error('Failed to load recent matches:', error)
  } finally {
    loadingRecentMatches.value = false
  }
}

async function loadFriends() {
  if (!userStore.userId) return

  loadingFriends.value = true
  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/api/v1/friends?user_id=${userStore.userId}`)
    if (response.ok) {
      friends.value = await response.json()
    }
  } catch (error) {
    console.error('Failed to load friends:', error)
  } finally {
    loadingFriends.value = false
  }
}

function startFastGame() {
  gameStore.startGame(userStore.soloTime, userStore.soloLetters, userStore.soloLang)
  router.push('/game')
}

function openCustomSetup() {
  userStore.setShowSoloSettings(true)
}

function openNewChallengeModal() {
  loadFriends()
  showNewChallengeModal.value = true
  modalStep.value = 'pick'
  selectedFriend.value = null
  searchFilter.value = ''
}

function closeNewChallengeModal() {
  showNewChallengeModal.value = false
  modalStep.value = 'pick'
  selectedFriend.value = null
  searchFilter.value = ''
}

async function sendChallenge() {
  if (!selectedFriend.value) return

  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

    // Create session
    const sessionResponse = await fetch(`${apiUrl}/api/v1/sessions?user_id=${userStore.userId}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        language: userStore.soloLang,
        letter_count: 7,
        time_limit: 60,
        hide_letters: true,
        invite_mode: 'friend',
        max_opponents: 1
      })
    })

    if (!sessionResponse.ok) throw new Error('Failed to create session')
    const session = await sessionResponse.json()

    // Send invite
    await fetch(`${apiUrl}/api/v1/sessions/${session.id}/invites?user_id=${userStore.userId}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ invited_user_id: selectedFriend.value.id })
    })

    modalStep.value = 'sent'
    await loadActiveChallenges()
  } catch (error) {
    console.error('Failed to send challenge:', error)
    userStore.showToast('Failed to send challenge', 'error')
  }
}

function playYourRound() {
  closeNewChallengeModal()
  // Navigate to the newly created challenge - would need session ID from sendChallenge
  userStore.showToast('Round starting…', 'default')
}

function hasUserPlayed(challenge) {
  return challenge.results?.some(r => r.user_id === userStore.userId) || false
}

function isYourTurn(challenge) {
  const userResult = challenge.results?.find(r => r.user_id === userStore.userId)
  return !userResult && challenge.type === 'invited'
}

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
  return date.toLocaleDateString()
}

onMounted(() => {
  loadActiveChallenges()
  loadRecentMatches()
})
</script>

<template>
  <div class="page">
    <div class="shell ph-wrap">
      <!-- Header -->
      <header class="page-head">
        <div>
          <div class="page-eyebrow">{{ $t('playHub.title') }}</div>
          <h1 class="page-title-display">
            <template v-if="yourTurnCount > 0">
              {{ $t('playHub.yourMove') }} <span style="color:var(--accent)">×{{ yourTurnCount }}</span>.
            </template>
            <template v-else>
              {{ $t('playHub.pickARound') }}
            </template>
          </h1>
        </div>
        <div class="row gap-2">
          <button class="btn btn--ghost btn--sm" @click="router.push('/match-history')">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"/>
              <path d="M12 6v6l4 2"/>
            </svg>
            {{ $t('playHub.matchHistory') }}
          </button>
          <button class="btn btn--ghost btn--sm" @click="router.push('/solo-history')">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="3" width="7" height="7"/>
              <rect x="14" y="3" width="7" height="7"/>
              <rect x="14" y="14" width="7" height="7"/>
              <rect x="3" y="14" width="7" height="7"/>
            </svg>
            {{ $t('playHub.soloHistory') }}
          </button>
        </div>
      </header>

      <!-- PRIMARY CTAs -->
      <section class="ph-ctas">
        <!-- Challenge friend card -->
        <div class="ph-card ph-card--multi">
          <div class="ph-card-corner">
            <span class="ph-num">01</span>
            <span class="ph-tag ph-tag--accent">{{$t('playHub.challengeFriend.tag')}}</span>
          </div>
          <h3 class="ph-card-title">{{ $t('playHub.challengeFriend.title') }}</h3>
          <p class="ph-card-sub">{{ $t('playHub.challengeFriend.subtitle') }}</p>
          <div class="ph-card-actions">
            <button class="btn btn--accent" @click="openNewChallengeModal">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M12 5v14M5 12h14"/>
              </svg>
              {{ $t('playHub.challengeFriend.newChallenge') }}
            </button>
            <button class="btn btn--soft" @click="router.push('/friends')">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="9" cy="7" r="4"/>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
              </svg>
              {{ $t('playHub.challengeFriend.friends') }}
            </button>
          </div>
          <div class="ph-card-meta">
            <span>
              <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="9" cy="7" r="4"/>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
              </svg>
              {{ $t('playHub.challengeFriend.friendsOnline', { count: 0 }) }}
            </span>
            <span class="dot-sep">·</span>
            <span>{{ $t('playHub.challengeFriend.waitingOnYou', { count: yourTurnCount }) }}</span>
          </div>
        </div>

        <!-- Daily card (MOCK) -->
        <div class="ph-card ph-card--daily">
          <div class="ph-card-corner">
            <span class="ph-num">02</span>
            <span class="ph-tag ph-tag--accent">{{$t('playHub.dailyPuzzle.tag')}}</span>
          </div>
          <h3 class="ph-card-title">{{ $t('playHub.dailyPuzzle.title') }}</h3>
          <p class="ph-card-sub">{{ $t('playHub.dailyPuzzle.subtitle') }}</p>
          <div class="ph-card-actions">
            <button class="btn btn--primary" disabled>
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M6 4l14 8-14 8z"/>
              </svg>
              {{ $t('playHub.dailyPuzzle.play') }}
            </button>
          </div>
          <div class="ph-daily-tiles">
            <span v-for="(L, i) in ['B','R','I','G','H','T','E']" :key="i" class="ph-daily-tile" :style="{ animationDelay: `${i * 0.04}s` }">{{ L }}</span>
          </div>
          <div class="ph-card-meta">
            <span>{{ $t('playHub.dailyPuzzle.resetsIn', { time: '9h 22m' }) }}</span>
          </div>
        </div>

        <!-- Solo / Warm up card -->
        <div class="ph-card ph-card--solo">
          <div class="ph-card-corner">
            <span class="ph-num">03</span>
            <span class="ph-tag">{{$t('playHub.warmUp.tag')}}</span>
          </div>
          <h3 class="ph-card-title">{{ $t('playHub.warmUp.title') }}</h3>
          <p class="ph-card-sub">{{ $t('playHub.warmUp.subtitle', { letters: userStore.soloLetters, time: userStore.soloTime, language: userStore.soloLang.toUpperCase() }) }}</p>
          <div class="ph-card-actions">
            <button class="btn btn--primary" @click="startFastGame">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M13 2L4 14h7l-1 8 9-12h-7l1-8z"/>
              </svg>
              {{ $t('playHub.warmUp.quickPlay') }}
            </button>
            <button class="btn btn--ghost btn--sm" @click="openCustomSetup">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="3"/>
                <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
              </svg>
              {{ $t('playHub.warmUp.custom') }}
            </button>
          </div>
          <div class="ph-card-meta">
            <span>
              <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"/>
                <path d="M12 6v6l4 2"/>
              </svg>
              {{ userStore.soloTime }}s
            </span>
            <span class="dot-sep">·</span>
            <span>{{ userStore.soloLetters }} {{$t('playHub.warmUp.letters')}}</span>
            <span class="dot-sep">·</span>
            <span>
              <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"/>
                <path d="M2 12h20M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"/>
              </svg>
              {{ userStore.soloLang.toUpperCase() }}
            </span>
          </div>
        </div>
      </section>

      <!-- ACTIVE CHALLENGES -->
      <section class="ph-section">
        <div class="ph-section-head">
          <div>
            <div class="page-eyebrow">{{ $t('playHub.activeChallenges.title') }}</div>
            <h2 class="ph-section-title">{{ $t('playHub.activeChallenges.inFlight', { count: activeChallenges.length }) }}</h2>
          </div>
          <div class="row gap-2 muted" style="font-size:12px">
            <span><span class="ph-dot ph-dot--turn" /> {{ $t('playHub.activeChallenges.yourTurn') }}</span>
            <span><span class="ph-dot ph-dot--wait" /> {{ $t('playHub.activeChallenges.waitingOnThem') }}</span>
          </div>
        </div>

        <div v-if="!userStore.isAuthenticated" class="ph-empty">
          <div class="ph-empty-illo">⌒</div>
          <div class="ph-empty-text">
            <strong>Sign in to see your challenges</strong>
            <span class="muted">Track your games and compete with friends.</span>
          </div>
          <button class="btn btn--accent btn--sm" @click="router.push('/auth')">
            Sign in
          </button>
        </div>

        <div v-else-if="activeChallenges.length === 0" class="ph-empty">
          <div class="ph-empty-illo">⌒</div>
          <div class="ph-empty-text">
            <strong>{{ $t('playHub.activeChallenges.noActive') }}</strong>
            <span class="muted">{{ $t('playHub.activeChallenges.sendChallenge') }}</span>
          </div>
          <button class="btn btn--accent btn--sm" @click="openNewChallengeModal">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M12 5v14M5 12h14"/>
            </svg>
            {{ $t('playHub.challengeFriend.newChallenge') }}
          </button>
        </div>

        <div v-else class="ph-challenges">
          <div
            v-for="challenge in activeChallenges"
            :key="challenge.id"
            :class="['ph-ch', { 'is-yours': isYourTurn(challenge), 'is-waiting': !isYourTurn(challenge) }]"
            @click="router.push(hasUserPlayed(challenge) ? `/results/${challenge.id}` : `/challenge/${challenge.id}`)"
          >
            <div class="ph-ch-side">
              <span :class="['ph-dot', isYourTurn(challenge) ? 'ph-dot--turn' : 'ph-dot--wait']" />
            </div>
            <div class="ph-ch-who">
              <div class="ph-ch-avatar">
                {{ challenge.type === 'created'
                  ? userStore.username?.charAt(0).toUpperCase()
                  : (challenge.creator_username?.charAt(0).toUpperCase() || '?') }}
              </div>
              <div>
                <div class="ph-ch-name">
                  {{ challenge.type === 'created'
                    ? 'You created'
                    : challenge.creator_username || 'Challenge' }}
                </div>
                <div class="ph-ch-meta">
                  <span>{{ formatDate(challenge.created_at) }}</span>
                  <span class="dot-sep">·</span>
                  <span>expires in 7d</span>
                </div>
              </div>
            </div>
            <div class="ph-ch-letters">
              <span v-for="(L, i) in challenge.letters.split('')" :key="i" class="ph-ch-tile">
                {{ hasUserPlayed(challenge) ? L.toLowerCase() : '?' }}
              </span>
            </div>
            <div class="ph-ch-status">
              <template v-if="isYourTurn(challenge)">
                <span class="ph-ch-label ph-ch-label--turn">{{ $t('playHub.activeChallenges.yourTurn') }}</span>
              </template>
              <template v-else-if="hasUserPlayed(challenge)">
                <span class="ph-ch-label ph-ch-label--wait">
                  {{ $t('playHub.activeChallenges.waitingOn', { name: challenge.creator_username || $t('playHub.activeChallenges.opponent') }) }}
                </span>
                <span class="ph-ch-score mono">
                  {{ challenge.results?.find(r => r.user_id === userStore.userId)?.score?.toLocaleString() || '—' }} pts
                </span>
              </template>
            </div>
            <div class="ph-ch-cta">
              <button v-if="isYourTurn(challenge)" class="btn btn--accent btn--sm" @click.stop="router.push(`/challenge/${challenge.id}`)">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M6 4l14 8-14 8z"/>
                </svg>
                Play
              </button>
              <button v-else class="btn btn--soft btn--sm" @click.stop="router.push(`/results/${challenge.id}`)">
                View
              </button>
            </div>
          </div>
        </div>
      </section>

      <!-- RECENT MATCHES -->
      <section v-if="recentMatches.length > 0" class="ph-section">
        <div class="ph-section-head">
          <div>
            <div class="page-eyebrow">{{ $t('playHub.recentMatches.title') }}</div>
            <h2 class="ph-section-title">{{ $t('playHub.recentMatches.lastFinished', { count: 5 }) }}</h2>
          </div>
          <button class="btn btn--ghost btn--sm" @click="router.push('/leaderboard')">
            {{ $t('playHub.recentMatches.seeAll') }}
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M5 12h14M12 5l7 7-7 7"/>
            </svg>
          </button>
        </div>
        <div class="ph-recent">
          <div
            v-for="m in recentMatches"
            :key="m.id"
            :class="['ph-recent-row', `rm-${m.result}`]"
          >
            <span class="ph-recent-result">{{ m.result === 'won' ? '↑' : m.result === 'lost' ? '↓' : '=' }}</span>
            <span class="ph-recent-with">{{ $t('playHub.recentMatches.vs') }} <strong>{{ m.with }}</strong></span>
            <span class="ph-recent-letters mono">{{ m.letters }}</span>
            <span class="ph-recent-score mono">
              {{ m.yourScore.toLocaleString() }} <span class="muted">·</span> {{ m.theirScore.toLocaleString() }}
            </span>
            <span class="ph-recent-delta mono">
              <template v-if="m.result === 'won'">+{{ Math.abs(m.yourScore - m.theirScore).toLocaleString() }}</template>
              <template v-else-if="m.result === 'lost'">−{{ Math.abs(m.yourScore - m.theirScore).toLocaleString() }}</template>
              <template v-else>EVEN</template>
            </span>
            <span class="ph-recent-date muted">{{ m.date }}</span>
          </div>
        </div>
      </section>
    </div>

    <!-- New Challenge Modal -->
    <Modal :show="showNewChallengeModal" @close="closeNewChallengeModal" :title="$t('playHub.newChallengeModal.title')" :subtitle="$t('playHub.newChallengeModal.subtitle')" max-width="500px">
      <div v-if="modalStep === 'sent'" class="nc-sent">
        <div class="nc-sent-icon">
          <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="20 6 9 17 4 12"/>
          </svg>
        </div>
        <h4 style="margin:12px 0 4px; font-family:var(--font-display); font-size:22px">
          {{ $t('playHub.newChallengeModal.sent', { name: selectedFriend?.username }) }}
        </h4>
        <p class="muted" style="font-size:13px; margin:0 0 18px; text-align:center">
          {{ $t('playHub.newChallengeModal.sentMessage') }}
        </p>
        <div class="row gap-2">
          <button class="btn btn--ghost" @click="closeNewChallengeModal">{{ $t('playHub.newChallengeModal.close') }}</button>
          <button class="btn btn--accent" @click="playYourRound">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M6 4l14 8-14 8z"/>
            </svg>
            {{ $t('playHub.newChallengeModal.playYourRound') }}
          </button>
        </div>
      </div>

      <div v-else class="nc-pick">
        <div class="nc-search">
          <input
            v-model="searchFilter"
            class="input"
            :placeholder="$t('playHub.newChallengeModal.searchPlaceholder')"
            autofocus
          />
        </div>
        <div v-if="loadingFriends" class="nc-loading">
          {{ $t('common.loading') }}
        </div>
        <div v-else-if="filteredFriends.length === 0" class="nc-empty muted">
          {{ $t('playHub.newChallengeModal.noFriends') }}
        </div>
        <div v-else class="nc-list">
          <button
            v-for="friend in filteredFriends"
            :key="friend.id"
            :class="['nc-item', { 'is-picked': selectedFriend?.id === friend.id }]"
            @click="selectedFriend = friend"
          >
            <div class="nc-avatar">{{ friend.username?.charAt(0).toUpperCase() }}</div>
            <div class="nc-meta">
              <div class="nc-name">{{ friend.username }}</div>
              <div class="nc-stat muted">Active player</div>
            </div>
            <svg v-if="selectedFriend?.id === friend.id" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
              <polyline points="20 6 9 17 4 12"/>
            </svg>
          </button>
        </div>
        <div class="nc-foot">
          <button class="btn btn--ghost" @click="closeNewChallengeModal">{{ $t('playHub.newChallengeModal.cancel') }}</button>
          <button class="btn btn--accent" :disabled="!selectedFriend" @click="sendChallenge">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M13 2L4 14h7l-1 8 9-12h-7l1-8z"/>
            </svg>
            {{ $t('playHub.newChallengeModal.send') }}
          </button>
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
/* Layout */
.ph-wrap {
  max-width: 1120px;
  margin: 0 auto;
  padding: 0 32px;
  padding-bottom: 80px;
}

@media (max-width: 720px) {
  .ph-wrap {
    padding: 0 20px;
    padding-bottom: 80px;
  }
}

/* Header */
.page-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 20px;
  margin-bottom: 32px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--border-hairline);
}

.page-eyebrow {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 2px;
  color: var(--fg-muted);
  font-weight: 600;
  margin-bottom: 8px;
}

.page-title-display {
  font-family: var(--font-display);
  font-size: 42px;
  font-weight: 700;
  letter-spacing: -1.4px;
  color: var(--fg-primary);
  margin: 0;
  line-height: 1;
}

/* Primary CTAs */
.ph-ctas {
  display: grid;
  grid-template-columns: 1.4fr 1fr 1fr;
  gap: 16px;
  margin-bottom: 48px;
}

.ph-card {
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 20px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  min-height: 240px;
  transition: all var(--dur-base) var(--ease-out);
  position: relative;
  overflow: hidden;
}

.ph-card:hover {
  transform: translateY(-2px);
  border-color: var(--border-default);
}

.ph-card--multi {
  background: linear-gradient(135deg, rgba(122,74,43,0.06) 0%, transparent 60%), var(--bg-surface);
  border-color: var(--accent-soft);
}

.ph-card--multi .ph-card-title {
  font-size: 32px;
}

.ph-card--daily {
  background: var(--navy);
  color: var(--milk);
  border-color: var(--navy);
}

.ph-card--daily .ph-card-title {
  color: var(--milk);
}

.ph-card--daily .ph-card-sub {
  color: rgba(251,246,236,0.7);
}

.ph-card--daily .ph-num {
  background: rgba(251,246,236,0.08);
  color: var(--milk);
  border: 1px solid rgba(251,246,236,0.14);
}

.ph-card--daily .ph-tag {
  background: rgba(251,246,236,0.12);
  color: var(--milk);
}

.ph-card--daily .ph-card-meta {
  color: rgba(251,246,236,0.5);
  border-top-color: rgba(251,246,236,0.14);
}

.ph-card--solo {
  background: var(--bg-surface);
}

.ph-card-corner {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.ph-num {
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 600;
  padding: 4px 10px;
  background: var(--bg-card);
  border: 1px solid var(--border-hairline);
  border-radius: 8px;
  color: var(--fg-primary);
}

.ph-tag {
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 1.2px;
  padding: 4px 10px;
  border-radius: 999px;
  background: var(--bg-card);
  color: var(--fg-secondary);
  font-weight: 600;
}

.ph-tag--accent {
  background: var(--accent);
  color: var(--milk);
}

.ph-card-title {
  font-family: var(--font-display);
  font-size: 26px;
  font-weight: 700;
  letter-spacing: -0.5px;
  margin: 0 0 10px;
  color: var(--fg-primary);
}

.ph-card-sub {
  font-size: 13px;
  line-height: 1.55;
  color: var(--fg-muted);
  margin: 0 0 auto;
  min-height: 40px;
}

.ph-card-actions {
  display: flex;
  gap: 8px;
  margin-top: 16px;
  flex-wrap: wrap;
}

.ph-card-meta {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border-hairline);
  font-size: 11px;
  color: var(--fg-muted);
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
}

.ph-card-meta svg {
  vertical-align: middle;
}

.dot-sep {
  opacity: 0.4;
}

/* Daily tiles */
.ph-daily-tiles {
  display: flex;
  gap: 4px;
  margin: 16px 0;
  flex-wrap: wrap;
}

.ph-daily-tile {
  width: 30px;
  height: 36px;
  background: rgba(251,246,236,0.08);
  border: 1px solid rgba(251,246,236,0.14);
  color: var(--milk);
  border-radius: 8px;
  display: grid;
  place-items: center;
  font-family: var(--font-mono);
  font-weight: 700;
  font-size: 14px;
  animation: pop-in 0.4s var(--ease-out) backwards;
}

@keyframes pop-in {
  from { transform: scale(0.6); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}

/* Section */
.ph-section {
  margin-top: 48px;
}

.ph-section-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 20px;
}

.ph-section-title {
  font-family: var(--font-display);
  font-size: 28px;
  font-weight: 700;
  letter-spacing: -0.8px;
  color: var(--fg-primary);
  margin: 0;
}

.ph-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 999px;
  margin-right: 6px;
  background: var(--fg-faint);
}

.ph-dot--turn {
  background: var(--milk);
  box-shadow: 0 0 0 3px rgba(251,246,236,0.25);
}

.ph-dot--wait {
  background: var(--fg-faint);
}

.ph-ch.is-waiting .ph-ch-side .ph-dot {
  background: var(--fg-faint);
  box-shadow: none;
}

/* Empty state */
.ph-empty {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  background: var(--bg-card);
  border: 1px dashed var(--border-default);
  border-radius: 16px;
}

.ph-empty-illo {
  width: 56px;
  height: 56px;
  background: var(--bg-surface);
  border-radius: 999px;
  display: grid;
  place-items: center;
  font-family: var(--font-display);
  font-size: 28px;
  color: var(--accent);
  flex-shrink: 0;
}

.ph-empty-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}

.ph-empty-text strong {
  font-size: 14px;
  font-weight: 700;
  color: var(--fg-primary);
}

.ph-empty-text .muted {
  font-size: 13px;
  color: var(--fg-secondary);
}

/* Challenges list */
.ph-challenges {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.ph-ch {
  display: grid;
  grid-template-columns: 6px 200px 1fr auto auto;
  gap: 16px;
  align-items: center;
  background: var(--bg-surface);
  border: 1px solid var(--border-hairline);
  border-radius: 14px;
  padding: 12px 16px 12px 0;
  cursor: pointer;
  transition: all var(--dur-base) var(--ease-out);
}

.ph-ch:hover {
  border-color: var(--border-default);
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

.ph-ch.is-yours {
  background: linear-gradient(90deg, rgba(122,74,43,0.04) 0%, var(--bg-surface) 30%);
}

.ph-ch-side {
  align-self: stretch;
  display: grid;
  place-items: center;
  border-radius: 14px 0 0 14px;
}

.ph-ch.is-yours .ph-ch-side {
  background: var(--accent);
}

.ph-ch-who {
  display: flex;
  align-items: center;
  gap: 10px;
}

.ph-ch-avatar {
  width: 38px;
  height: 38px;
  border-radius: 999px;
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  color: var(--fg-primary);
  display: grid;
  place-items: center;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 14px;
  flex-shrink: 0;
}

.ph-ch-name {
  font-weight: 600;
  font-size: 14px;
  color: var(--fg-primary);
}

.ph-ch-meta {
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--fg-muted);
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 2px;
}

.ph-ch-letters {
  display: flex;
  gap: 3px;
}

.ph-ch-tile {
  width: 26px;
  height: 30px;
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 6px;
  display: grid;
  place-items: center;
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 700;
  color: var(--fg-primary);
}

.ph-ch-status {
  display: flex;
  flex-direction: column;
  gap: 2px;
  align-items: flex-end;
  min-width: 140px;
}

.ph-ch-label {
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.8px;
}

.ph-ch-label--turn {
  color: var(--accent);
}

.ph-ch-label--wait {
  color: var(--fg-muted);
}

.ph-ch-score {
  font-size: 12px;
  color: var(--fg-secondary);
}

.ph-ch-cta {
  flex-shrink: 0;
}

/* Recent matches */
.ph-recent {
  display: flex;
  flex-direction: column;
  border: 1px solid var(--border-hairline);
  border-radius: 14px;
  background: var(--bg-surface);
  overflow: hidden;
}

.ph-recent-row {
  display: grid;
  grid-template-columns: 32px 1fr 100px 200px 100px auto;
  gap: 12px;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--border-hairline);
  font-size: 13px;
}

.ph-recent-row:last-child {
  border-bottom: 0;
}

.ph-recent-result {
  font-family: var(--font-mono);
  font-weight: 700;
  font-size: 18px;
  text-align: center;
  color: var(--fg-faint);
}

.rm-won .ph-recent-result {
  color: var(--success);
}

.rm-lost .ph-recent-result {
  color: var(--danger);
}

.rm-tie .ph-recent-result {
  color: var(--fg-muted);
}

.ph-recent-with {
  color: var(--fg-secondary);
  font-size: 13px;
}

.ph-recent-with strong {
  font-weight: 600;
  color: var(--fg-primary);
}

.ph-recent-letters {
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 1px;
}

.ph-recent-score {
  font-size: 13px;
  color: var(--fg-secondary);
}

.ph-recent-delta {
  font-size: 13px;
  font-weight: 700;
  text-align: right;
}

.rm-won .ph-recent-delta { color: var(--success); }
.rm-lost .ph-recent-delta { color: var(--danger); }
.rm-tie .ph-recent-delta { color: var(--fg-muted); }

.ph-recent-date {
  font-size: 11px;
  text-align: right;
}

/* New Challenge Modal */
.nc-sent {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24px 0;
}

.nc-sent-icon {
  width: 64px;
  height: 64px;
  border-radius: 999px;
  background: var(--accent);
  color: var(--milk);
  display: grid;
  place-items: center;
  margin-bottom: 12px;
  animation: an-fade-up 0.3s var(--ease-out);
}

@keyframes an-fade-up {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.nc-pick {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.nc-search .input {
  width: 100%;
}

.nc-loading,
.nc-empty {
  text-align: center;
  padding: 32px 16px;
  color: var(--fg-muted);
  font-size: 13px;
}

.nc-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-height: 320px;
  overflow-y: auto;
  padding-right: 4px;
}

.nc-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  background: transparent;
  border: 1px solid var(--border-hairline);
  border-radius: 12px;
  cursor: pointer;
  transition: all var(--dur-fast) var(--ease-std);
  text-align: left;
  appearance: none;
  font-family: var(--font-body);
}

.nc-item:hover {
  background: var(--bg-card);
  border-color: var(--border-subtle);
}

.nc-item.is-picked {
  background: var(--accent-soft);
  border-color: var(--accent);
}

.nc-avatar {
  width: 36px;
  height: 36px;
  border-radius: 999px;
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  color: var(--fg-primary);
  display: grid;
  place-items: center;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 13px;
  flex-shrink: 0;
}

.nc-meta {
  flex: 1;
  min-width: 0;
}

.nc-name {
  font-weight: 600;
  font-size: 14px;
  color: var(--fg-primary);
  display: flex;
  align-items: center;
  gap: 6px;
}

.nc-stat {
  font-size: 11px;
  margin-top: 2px;
}

.nc-foot {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border-hairline);
}

/* Responsive */
@media (max-width: 980px) {
  .ph-ctas {
    grid-template-columns: 1fr 1fr;
  }

  .ph-card--multi {
    grid-column: 1 / -1;
  }

  .ph-card--multi .ph-card-title {
    font-size: 28px;
  }
}

@media (max-width: 860px) {
  .page-title-display {
    font-size: 32px;
    letter-spacing: -1px;
  }
}

@media (max-width: 760px) {
  .ph-ch {
    grid-template-columns: 4px 1fr auto;
    grid-template-areas:
      "side who cta"
      "side letters letters"
      "side status status";
    padding: 12px 12px 12px 0;
    row-gap: 8px;
  }

  .ph-ch-side {
    grid-area: side;
  }

  .ph-ch-who {
    grid-area: who;
  }

  .ph-ch-letters {
    grid-area: letters;
  }

  .ph-ch-status {
    grid-area: status;
    align-items: flex-start;
  }

  .ph-ch-cta {
    grid-area: cta;
  }

  .ph-recent-row {
    grid-template-columns: 24px 1fr auto;
    grid-template-areas:
      "res with delta"
      "res letters score"
      ". date date";
    row-gap: 2px;
  }

  .ph-recent-result {
    grid-area: res;
  }

  .ph-recent-with {
    grid-area: with;
  }

  .ph-recent-letters {
    grid-area: letters;
    font-size: 11px;
  }

  .ph-recent-score {
    grid-area: score;
    text-align: right;
    font-size: 11px;
  }

  .ph-recent-delta {
    grid-area: delta;
    font-size: 13px;
  }

  .ph-recent-date {
    grid-area: date;
    text-align: left;
  }
}

@media (max-width: 640px) {
  .ph-ctas {
    grid-template-columns: 1fr;
  }

  .ph-card--multi {
    grid-column: auto;
  }
}
</style>
