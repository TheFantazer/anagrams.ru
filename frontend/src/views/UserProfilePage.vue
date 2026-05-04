<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../stores/userStore'
import '../assets/profile.css'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const username = route.params.username
const user = ref(null)
const loading = ref(true)
const notFound = ref(false)
const linkCopied = ref(false)
const requested = ref(false)
const accepted = ref(false)

// Friend request states
const friendshipStatus = ref(null) // 'friend' | 'incoming' | 'outgoing' | 'stranger' | null (self)

// Match history (mock for now)
const matches = ref([])
const friendsList = ref([])

const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// Check if viewing own profile
const isMe = computed(() => {
  return userStore.isAuthenticated && user.value && user.value.username === userStore.username
})

// Determine relationship
const rel = computed(() => {
  if (isMe.value) return null
  return friendshipStatus.value
})

// Head-to-head stats for friends
const wins = computed(() => matches.value.filter(m => m.result === 'won').length)
const losses = computed(() => matches.value.filter(m => m.result === 'lost').length)
const ties = computed(() => matches.value.filter(m => m.result === 'tie').length)

// Own profile stats
const ownStats = computed(() => {
  if (!isMe.value || !user.value) return null
  return [
    { v: user.value.score?.toLocaleString() || '0', k: t('profile.leaderboard') },
    { v: user.value.totalGames || 0, k: t('profile.gamesPlayed') },
    { v: `${user.value.winRate || 0}%`, k: t('profile.winRate') },
    { v: `${user.value.streak || 0}d`, k: t('profile.currentStreak') },
  ]
})

// Decorative letters from username or best word
const decoLetters = computed(() => {
  if (!user.value) return []
  const text = user.value.bestWord || user.value.username || ''
  return text.split('').slice(0, 8)
})

onMounted(async () => {
  await loadUser()
  if (user.value && !isMe.value) {
    await checkFriendshipStatus()
    if (rel.value === 'friend') {
      loadMatches()
    }
  }
})

async function loadUser() {
  loading.value = true
  try {
    const response = await fetch(`${apiUrl}/api/v1/users/username/${username}`)
    if (response.ok) {
      user.value = await response.json()
      notFound.value = false
    } else if (response.status === 404) {
      notFound.value = true
    }
  } catch (error) {
    console.error('Failed to load user:', error)
    notFound.value = true
  } finally {
    loading.value = false
  }
}

async function checkFriendshipStatus() {
  if (!userStore.userId) return

  try {
    // Check if friends
    const friendsRes = await fetch(`${apiUrl}/api/v1/friends?user_id=${userStore.userId}`)
    if (friendsRes.ok) {
      const friends = await friendsRes.json()
      const isFriend = friends.some(f => f.id === user.value.id)
      if (isFriend) {
        friendshipStatus.value = 'friend'
        friendsList.value = friends
        return
      }
    }

    // Check incoming requests
    const incomingRes = await fetch(`${apiUrl}/api/v1/friends/requests/pending?user_id=${userStore.userId}`)
    if (incomingRes.ok) {
      const incoming = await incomingRes.json()
      const hasIncoming = incoming.some(r => r.from_user_id === user.value.id)
      if (hasIncoming) {
        friendshipStatus.value = 'incoming'
        return
      }
    }

    // Check outgoing requests
    const outgoingRes = await fetch(`${apiUrl}/api/v1/friends/requests/sent?user_id=${userStore.userId}`)
    if (outgoingRes.ok) {
      const outgoing = await outgoingRes.json()
      const hasOutgoing = outgoing.some(r => r.to_user_id === user.value.id)
      if (hasOutgoing) {
        friendshipStatus.value = 'outgoing'
        return
      }
    }

    // Stranger
    friendshipStatus.value = 'stranger'
  } catch (error) {
    console.error('Failed to check friendship status:', error)
    friendshipStatus.value = 'stranger'
  }
}

async function loadMatches() {
  // Mock data - head-to-head match history endpoint not yet implemented
  // This will be replaced with real API call to /api/v1/matches/head-to-head/:userId
  matches.value = [
    { id: 1, letters: 'BEAUTY', yourScore: 8420, theirScore: 7100, result: 'won', date: 'Yesterday' },
    { id: 2, letters: 'SHELTER', yourScore: 5240, theirScore: 6820, result: 'lost', date: '2 days ago' },
    { id: 3, letters: 'CRESTED', yourScore: 7820, theirScore: 7820, result: 'tie', date: '3 days ago' },
  ]
}

async function handleAddFriend() {
  if (!userStore.userId) return

  try {
    const response = await fetch(`${apiUrl}/api/v1/friends/requests`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        to_user_id: user.value.id
      }),
      credentials: 'include'
    })

    if (response.ok) {
      requested.value = true
      friendshipStatus.value = 'outgoing'
      userStore.showToast(t('profile.friendRequestSent', { name: user.value.username }), 'success')
    } else {
      const error = await response.json()
      userStore.showToast(error.message || t('profile.friendRequestFailed'), 'error')
    }
  } catch (error) {
    console.error('Failed to send friend request:', error)
    userStore.showToast(t('profile.friendRequestFailed'), 'error')
  }
}

async function handleAccept() {
  if (!userStore.userId) return

  try {
    // Find the request ID
    const incomingRes = await fetch(`${apiUrl}/api/v1/friends/requests/pending?user_id=${userStore.userId}`)
    if (!incomingRes.ok) return

    const incoming = await incomingRes.json()
    const request = incoming.find(r => r.from_user_id === user.value.id)
    if (!request) return

    const response = await fetch(`${apiUrl}/api/v1/friends/requests/${request.id}/accept`, {
      method: 'POST',
      credentials: 'include'
    })

    if (response.ok) {
      accepted.value = true
      friendshipStatus.value = 'friend'
      userStore.showToast(t('profile.friendRequestAccepted', { name: user.value.username }), 'success')
    }
  } catch (error) {
    console.error('Failed to accept friend request:', error)
    userStore.showToast(t('profile.friendRequestFailed'), 'error')
  }
}

async function handleDecline() {
  if (!userStore.userId) return

  try {
    const incomingRes = await fetch(`${apiUrl}/api/v1/friends/requests/pending?user_id=${userStore.userId}`)
    if (!incomingRes.ok) return

    const incoming = await incomingRes.json()
    const request = incoming.find(r => r.from_user_id === user.value.id)
    if (!request) return

    const response = await fetch(`${apiUrl}/api/v1/friends/requests/${request.id}/reject`, {
      method: 'POST',
      credentials: 'include'
    })

    if (response.ok) {
      userStore.showToast(t('profile.requestDeclined', { name: user.value.username }), 'info')
      router.push('/friends')
    }
  } catch (error) {
    console.error('Failed to decline friend request:', error)
  }
}

async function handleChallenge() {
  if (!userStore.userId) return

  try {
    // Create a new game session
    const sessionResponse = await fetch(`${apiUrl}/api/v1/sessions?user_id=${userStore.userId}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        language: userStore.soloLang,
        letter_count: userStore.soloLetters,
        time_limit: userStore.soloTime
      })
    })

    if (!sessionResponse.ok) {
      throw new Error('Failed to create session')
    }

    const session = await sessionResponse.json()

    // Create invite for the user
    const inviteResponse = await fetch(`${apiUrl}/api/v1/sessions/${session.id}/invites?user_id=${userStore.userId}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({
        invited_user_id: user.value.id
      })
    })

    if (!inviteResponse.ok) {
      throw new Error('Failed to create invite')
    }

    userStore.showToast(t('profile.challengeSent', { name: user.value.username }), 'success')

    // Redirect to game page
    router.push(`/game?session_id=${session.id}`)
  } catch (error) {
    console.error('Failed to create challenge:', error)
    userStore.showToast(t('profile.challengeFailed'), 'error')
  }
}

function copyLink() {
  const link = `${window.location.origin}/u/${user.value.username}`
  navigator.clipboard?.writeText(link)
  linkCopied.value = true
  userStore.showToast(t('profile.linkCopied'), 'success')
  setTimeout(() => {
    linkCopied.value = false
  }, 2200)
}

function goBack() {
  if (rel.value === 'friend' || rel.value === 'incoming' || rel.value === 'outgoing') {
    router.push('/friends')
  } else {
    router.push('/')
  }
}

function formatJoinDate(dateStr) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString('en-US', { year: 'numeric', month: 'short' })
}
</script>

<template>
  <div class="page">
    <div v-if="loading" class="shell pf-wrap">
      <p class="muted">{{ t('common.loading') }}...</p>
    </div>

    <div v-else-if="notFound" class="shell pf-wrap" style="padding: 80px 0; text-align: center">
      <h2 class="page-title-display">{{ t('profile.userNotFound') }}</h2>
      <p class="muted">{{ t('profile.noPlayerNamed', { name: username }) }}</p>
      <button class="btn btn--accent" style="margin-top: 16px" @click="router.push('/')">
        {{ t('profile.backToHome') }}
      </button>
    </div>

    <div v-else-if="user" class="shell pf-wrap">
      <!-- Breadcrumb / back -->
      <div class="pf-back">
        <button class="back-link" @click="goBack">
          ← {{ t('profile.backTo') }}
        </button>
      </div>

      <!-- HEADER -->
      <header class="pf-head">
        <!-- Decorative background letters -->
        <div class="pf-head-deco" aria-hidden="true">
          <span
            v-for="(letter, i) in decoLetters"
            :key="i"
            :style="{
              left: `${(i * 13 + 8) % 90}%`,
              top: `${(i * 31 + 12) % 80}%`,
              animationDelay: `${i * 0.18}s`
            }"
          >
            {{ letter.toUpperCase() }}
          </span>
        </div>

        <div class="pf-head-main">
          <div class="pf-avatar-wrap">
            <div class="pf-avatar">{{ user.username[0].toUpperCase() }}</div>
            <span v-if="user.online" class="pf-avatar-online" />
          </div>

          <div class="pf-head-text">
            <div class="pf-eyebrow">
              <span v-if="isMe" class="pf-rel-tag pf-rel-tag--self">{{ t('profile.yourProfile') }}</span>
              <span v-else-if="rel === 'friend'" class="pf-rel-tag pf-rel-tag--friend">
                ✓ {{ t('profile.friend') }}
              </span>
              <span v-else-if="rel === 'incoming'" class="pf-rel-tag pf-rel-tag--in">
                {{ t('profile.wantsToBeFriends') }}
              </span>
              <span v-else-if="rel === 'outgoing'" class="pf-rel-tag pf-rel-tag--out">
                {{ t('profile.requestPending') }}
              </span>
              <span v-else-if="rel === 'stranger'" class="pf-rel-tag">
                {{ t('profile.notFriendYet') }}
              </span>

              <span v-if="user.online" class="pf-online-pill">
                <span class="online-dot" /> {{ t('profile.onlineNow') }}
              </span>
            </div>

            <h1 class="pf-name">{{ user.username }}</h1>

            <div class="pf-meta">
              <span>🏆 {{ (user.score || 0).toLocaleString() }} {{ t('profile.pts') }}</span>
              <span class="dot-sep">·</span>
              <span>⚡ {{ user.streak || 0 }}d {{ t('profile.streak') }}</span>
              <span class="dot-sep">·</span>
              <span>{{ t('profile.joined') }} {{ formatJoinDate(user.created_at) }}</span>
            </div>

            <p v-if="isMe && user.bio" class="pf-bio">{{ user.bio }}</p>
          </div>
        </div>

        <!-- CTAs row -->
        <div class="pf-actions">
          <!-- Own profile -->
          <template v-if="isMe">
            <button class="btn btn--accent" @click="router.push('/settings')">
              ⚙️ {{ t('profile.editProfile') }}
            </button>
            <button class="btn btn--soft" @click="copyLink">
              {{ linkCopied ? '✓ ' + t('profile.copied') : '🔗 ' + t('profile.shareProfile') }}
            </button>
          </template>

          <!-- Friend -->
          <template v-else-if="rel === 'friend'">
            <button class="btn btn--accent" @click="handleChallenge">
              ⚡ {{ t('profile.challenge') }}
            </button>
            <button class="btn btn--soft" @click="copyLink">
              {{ linkCopied ? '✓ ' + t('profile.copied') : '🔗 ' + t('profile.share') }}
            </button>
          </template>

          <!-- Incoming request -->
          <template v-else-if="rel === 'incoming' && !accepted">
            <button class="btn btn--accent" @click="handleAccept">
              ✓ {{ t('profile.acceptRequest') }}
            </button>
            <button class="btn btn--ghost" @click="handleDecline">
              {{ t('profile.decline') }}
            </button>
          </template>
          <template v-else-if="rel === 'incoming' && accepted">
            <button class="btn btn--accent" @click="handleChallenge">
              ⚡ {{ t('profile.challenge') }}
            </button>
          </template>

          <!-- Outgoing request -->
          <template v-else-if="rel === 'outgoing'">
            <button class="btn btn--soft" disabled>
              <span class="pulse-dot" /> {{ t('profile.awaitingReply') }}
            </button>
          </template>

          <!-- Stranger -->
          <template v-else-if="rel === 'stranger' && !requested">
            <template v-if="!userStore.isAuthenticated">
              <button class="btn btn--accent" @click="router.push('/auth')">
                👤 {{ t('profile.signInToAdd') }}
              </button>
              <span class="pf-action-hint muted">{{ t('profile.needAccountHint') }}</span>
            </template>
            <template v-else>
              <button class="btn btn--accent" @click="handleAddFriend">
                ➕ {{ t('profile.addFriend') }}
              </button>
              <button class="btn btn--soft" @click="handleChallenge">
                ⚡ {{ t('profile.challengeAnyway') }}
              </button>
            </template>
          </template>
          <template v-else-if="rel === 'stranger' && requested">
            <button class="btn btn--soft" disabled>
              <span class="pulse-dot" /> {{ t('profile.requestSent') }}
            </button>
          </template>
        </div>
      </header>

      <!-- STATS GRID -->
      <section v-if="isMe && ownStats" class="pf-stats-grid pf-stats-grid--own">
        <div v-for="stat in ownStats" :key="stat.k" class="pf-stat">
          <div class="pf-stat-v">{{ stat.v }}</div>
          <div class="pf-stat-k">{{ stat.k }}</div>
        </div>
      </section>

      <section v-else class="pf-stats-grid">
        <div class="pf-stat">
          <div class="pf-stat-v mono">{{ (user.score || 0).toLocaleString() }}</div>
          <div class="pf-stat-k">{{ t('profile.leaderboard') }}</div>
        </div>
        <div class="pf-stat">
          <div class="pf-stat-v mono">{{ (user.bestWord || 'N/A').toLowerCase() }}</div>
          <div class="pf-stat-k">{{ t('profile.bestWord') }}</div>
        </div>
        <div class="pf-stat">
          <div class="pf-stat-v mono">{{ user.streak || 0 }}d</div>
          <div class="pf-stat-k">{{ t('profile.streak') }}</div>
        </div>
        <div v-if="rel === 'friend'" class="pf-stat">
          <div class="pf-stat-v mono">{{ matches.length }}</div>
          <div class="pf-stat-k">{{ t('profile.matchesWithYou') }}</div>
        </div>
      </section>

      <!-- HEAD-TO-HEAD (only for friends with history) -->
      <section v-if="rel === 'friend' && matches.length > 0" class="pf-section">
        <div class="page-eyebrow">{{ t('profile.headToHead') }}</div>
        <h2 class="pf-section-title">
          {{ wins }} – {{ losses }}{{ ties > 0 ? ` – ${ties}` : '' }}
          <span class="muted" style="font-weight: 400; font-size: 14px">
            {{ t('profile.overRounds', { count: matches.length }) }}
          </span>
        </h2>

        <div class="fr-prof-bar" style="margin-top: 14px">
          <div class="fr-prof-bar-w" :style="{ flex: wins || 0.001 }">
            <span v-if="wins > 0">{{ wins }} W</span>
          </div>
          <div class="fr-prof-bar-t" :style="{ flex: ties || 0.001 }">
            <span v-if="ties > 0">{{ ties }} T</span>
          </div>
          <div class="fr-prof-bar-l" :style="{ flex: losses || 0.001 }">
            <span v-if="losses > 0">{{ losses }} L</span>
          </div>
        </div>

        <div class="fr-prof-recent" style="margin-top: 14px">
          <div class="fr-eyebrow-mini">{{ t('profile.recentRounds') }}</div>
          <div
            v-for="match in matches.slice(0, 5)"
            :key="match.id"
            :class="`fr-prof-recent-row rm-${match.result}`"
          >
            <span :class="`fr-prof-recent-tag rm-${match.result}`">
              {{ match.result === 'won' ? 'W' : match.result === 'lost' ? 'L' : 'T' }}
            </span>
            <span class="mono fr-prof-recent-letters">{{ match.letters }}</span>
            <span class="mono">
              {{ match.yourScore.toLocaleString() }}
              <span class="muted">·</span>
              {{ match.theirScore.toLocaleString() }}
            </span>
            <span class="muted" style="margin-left: auto; font-size: 12px">{{ match.date }}</span>
          </div>
        </div>
      </section>

      <!-- Invite footer for strangers -->
      <section v-if="rel === 'stranger' && !requested && userStore.isAuthenticated" class="pf-invite-footer">
        <div>
          <div class="page-eyebrow">{{ t('profile.whyYouAreHere') }}</div>
          <p class="muted" style="margin: 4px 0 0">
            {{ t('profile.inviteExplanation') }}
          </p>
        </div>
        <button class="btn btn--accent" @click="handleAddFriend">
          ➕ {{ t('profile.addFriend') }}
        </button>
      </section>
    </div>
  </div>
</template>

<style scoped>
.back-link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: none;
  border: none;
  color: var(--fg-secondary);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  padding: 10px 0;
  transition: all var(--dur-fast);
  text-decoration: none;
}
.back-link:hover {
  color: var(--accent);
  transform: translateX(-2px);
}
</style>
