<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../stores/userStore'
import Modal from '../components/Modal.vue'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()

const tab = ref('all') // all, requests, add
const friends = ref([])
const incomingRequests = ref([])
const outgoingRequests = ref([])
const searchQuery = ref('')
const addSearchQuery = ref('')
const searchResults = ref([])
const openProfile = ref(null)
const loading = ref(false)
const inviteCopied = ref(false)

const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

// Filtered friends based on search
const filteredFriends = computed(() => {
  if (!searchQuery.value.trim()) return friends.value
  const q = searchQuery.value.toLowerCase()
  return friends.value.filter(f =>
    f.username?.toLowerCase().includes(q)
  )
})

// Online friends count
const onlineFriendsCount = computed(() => {
  return filteredFriends.value.filter(f => f.online).length
})

// Mock suggested friends - TODO: Replace with real API
const suggestedFriends = ref([
  { name: 'alex_word', reason: 'Plays often', id: 999 },
  { name: 'maria_fast', reason: 'Top 100 player', id: 998 },
])

// Match history for friend profile
const friendMatches = ref([])
const loadingFriendMatches = ref(false)

onMounted(() => {
  if (userStore.isAuthenticated) {
    loadFriends()
    loadRequests()
  }
})

// Load friend matches when profile modal opens
watch(openProfile, (newProfile) => {
  if (newProfile && newProfile.id) {
    loadFriendMatches(newProfile.id)
  } else {
    friendMatches.value = []
  }
})

async function loadFriends() {
  if (!userStore.userId) return

  loading.value = true
  try {
    const response = await fetch(`${apiUrl}/api/v1/friends?user_id=${userStore.userId}`)
    if (response.ok) {
      const data = await response.json()
      // Add mock data for demonstration
      friends.value = data.map((f, i) => ({
        ...f,
        online: i % 3 === 0, // Mock online status
        score: 5000 + i * 200, // Mock score
        streak: 5 + i, // Mock streak
        bestWord: 'streami', // Mock best word
        joined: '2 weeks ago' // Mock join date
      }))
    }
  } catch (error) {
    console.error('Failed to load friends:', error)
  } finally {
    loading.value = false
  }
}

async function loadRequests() {
  if (!userStore.userId) return

  try {
    // Load incoming
    const incomingRes = await fetch(`${apiUrl}/api/v1/friends/requests/pending?user_id=${userStore.userId}`)
    if (incomingRes.ok) {
      const incoming = await incomingRes.json()
      incomingRequests.value = incoming.map((r, i) => ({
        ...r,
        name: r.from_username || 'User',
        mutual: i % 2, // Mock mutual friends
        sent: '2 days ago', // Mock sent time
        msg: i === 0 ? 'Let\'s play!' : null // Mock message
      }))
    }

    // Load outgoing
    const outgoingRes = await fetch(`${apiUrl}/api/v1/friends/requests/sent?user_id=${userStore.userId}`)
    if (outgoingRes.ok) {
      const outgoing = await outgoingRes.json()
      outgoingRequests.value = outgoing.map((r, i) => ({
        ...r,
        name: r.to_username || 'User',
        mutual: i % 3, // Mock mutual friends
        sent: '1 day ago' // Mock sent time
      }))
    }
  } catch (error) {
    console.error('Failed to load requests:', error)
  }
}

async function searchUsers() {
  if (!addSearchQuery.value.trim()) {
    searchResults.value = []
    return
  }

  try {
    const response = await fetch(`${apiUrl}/api/v1/users/search?q=${encodeURIComponent(addSearchQuery.value)}`)
    if (response.ok) {
      const results = await response.json()
      searchResults.value = results.filter(user => user.id !== userStore.userId)
    }
  } catch (error) {
    console.error('Failed to search users:', error)
  }
}

async function sendRequest(userId, name) {
  try {
    const response = await fetch(`${apiUrl}/api/v1/friends/requests?user_id=${userStore.userId}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ to_user_id: userId })
    })

    if (response.ok) {
      userStore.showToast(`Friend request sent to ${name}`, 'success')
      await loadRequests()
      searchResults.value = searchResults.value.filter(u => u.id !== userId)
    } else {
      userStore.showToast('Failed to send request', 'error')
    }
  } catch (error) {
    console.error('Failed to send request:', error)
    userStore.showToast('Failed to send request', 'error')
  }
}

async function acceptRequest(requestId, name) {
  try {
    const response = await fetch(`${apiUrl}/api/v1/friends/requests/${requestId}/accept?user_id=${userStore.userId}`, {
      method: 'POST'
    })

    if (response.ok) {
      userStore.showToast(`${name} added to friends`, 'success')
      await Promise.all([loadFriends(), loadRequests()])
    } else {
      userStore.showToast('Failed to accept request', 'error')
    }
  } catch (error) {
    console.error('Failed to accept request:', error)
    userStore.showToast('Failed to accept request', 'error')
  }
}

async function declineRequest(requestId, name) {
  try {
    const response = await fetch(`${apiUrl}/api/v1/friends/requests/${requestId}/reject?user_id=${userStore.userId}`, {
      method: 'POST'
    })

    if (response.ok) {
      userStore.showToast(`Request from ${name} declined`, 'info')
      await loadRequests()
    } else {
      userStore.showToast('Failed to decline request', 'error')
    }
  } catch (error) {
    console.error('Failed to decline request:', error)
    userStore.showToast('Failed to decline request', 'error')
  }
}

async function cancelRequest(requestId, name) {
  try {
    const response = await fetch(`${apiUrl}/api/v1/friends/requests/${requestId}?user_id=${userStore.userId}`, {
      method: 'DELETE'
    })

    if (response.ok) {
      userStore.showToast(`Request to ${name} cancelled`, 'info')
      await loadRequests()
    } else {
      userStore.showToast('Failed to cancel request', 'error')
    }
  } catch (error) {
    console.error('Failed to cancel request:', error)
    userStore.showToast('Failed to cancel request', 'error')
  }
}

function copyInviteLink() {
  const link = `anagrams.ru/u/${userStore.username}`
  navigator.clipboard?.writeText(link)
  inviteCopied.value = true
  userStore.showToast('Invite link copied', 'success')
  setTimeout(() => { inviteCopied.value = false }, 2200)
}

function challengeFriend(name) {
  userStore.showToast(`Challenge sent to ${name}!`, 'success')
  router.push('/play')
}

// Load matches against a specific friend
async function loadFriendMatches(friendId) {
  if (!userStore.userId || !friendId) return

  loadingFriendMatches.value = true
  friendMatches.value = []

  try {
    const response = await fetch(`${apiUrl}/api/v1/sessions/all?user_id=${userStore.userId}&page=1&per_page=100`)

    if (response.ok) {
      const data = await response.json()

      // Filter for completed matches with this specific friend
      const matchesWithFriend = (data.sessions || []).filter(session => {
        if (!session.results || session.results.length < 2) return false

        // Check if both user and friend have results in this session
        const hasUserResult = session.results.some(r => r.user_id === userStore.userId)
        const hasFriendResult = session.results.some(r => r.user_id === friendId)

        return hasUserResult && hasFriendResult
      })

      // Transform to match format
      friendMatches.value = matchesWithFriend.map(session => {
        const myResult = session.results.find(r => r.user_id === userStore.userId)
        const friendResult = session.results.find(r => r.user_id === friendId)

        if (!myResult || !friendResult) return null

        const delta = myResult.score - friendResult.score
        let result = 'tie'
        if (delta > 0) result = 'won'
        else if (delta < 0) result = 'lost'

        return {
          id: session.id,
          result,
          letters: session.letters,
          yourScore: myResult.score,
          theirScore: friendResult.score,
          date: formatRelativeTime(myResult.played_at || session.created_at)
        }
      }).filter(Boolean)
    }
  } catch (error) {
    console.error('Failed to load friend matches:', error)
  } finally {
    loadingFriendMatches.value = false
  }
}

// Helper to format relative time
function formatRelativeTime(dateStr) {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now - date
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffDays === 0) return 'today'
  if (diffDays === 1) return 'yesterday'
  if (diffDays < 7) return `${diffDays} days ago`
  if (diffDays < 30) return `${Math.floor(diffDays / 7)} weeks ago`
  return date.toLocaleDateString()
}

// Head-to-head stats for profile modal
function getH2HStats(friendName) {
  const wins = friendMatches.value.filter(m => m.result === 'won').length
  const losses = friendMatches.value.filter(m => m.result === 'lost').length
  const ties = friendMatches.value.filter(m => m.result === 'tie').length
  return { wins, losses, ties }
}
</script>

<template>
  <div class="page">
    <div class="shell fr-wrap">
      <header class="page-head">
        <div>
          <div class="page-eyebrow">{{ $t('friends.title') }}</div>
          <h1 class="page-title-display">
            {{ friends.length }} {{ friends.length === 1 ? $t('friends.friend') : $t('friends.friends') }}
            <template v-if="incomingRequests.length > 0">, <span style="color:var(--accent)">{{ incomingRequests.length }} waiting</span></template>.
          </h1>
        </div>
        <button class="btn btn--accent btn--sm" @click="tab = 'add'">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 5v14M5 12h14"/>
          </svg>
          {{$t('friends.addFriend')}}
        </button>
      </header>

      <!-- Tabs -->
      <div class="fr-tabs">
        <button
          :class="['fr-tab', { 'is-active': tab === 'all' }]"
          @click="tab = 'all'"
        >
          {{$t('friends.allFriends')}}
          <span v-if="friends.length > 0" class="fr-tab-count">{{ friends.length }}</span>
        </button>
        <button
          :class="['fr-tab', { 'is-active': tab === 'requests' }]"
          @click="tab = 'requests'"
        >
          {{$t('friends.requests')}}
          <span v-if="incomingRequests.length > 0" class="fr-tab-count is-accent">{{ incomingRequests.length }}</span>
        </button>
        <button
          :class="['fr-tab', { 'is-active': tab === 'add' }]"
          @click="tab = 'add'"
        >
          {{$t('friends.addFriend')}}
        </button>
      </div>

      <!-- ALL FRIENDS -->
      <section v-if="tab === 'all'">
        <div class="fr-search">
          <input
            v-model="searchQuery"
            class="input"
            placeholder="Search by username…"
          />
          <span class="fr-online-count muted">
            <span class="online-dot" /> {{ onlineFriendsCount }} {{$t('friends.online')}}
          </span>
        </div>
        <div class="fr-grid">
          <div v-for="friend in filteredFriends" :key="friend.id" class="fr-card">
            <button class="fr-card-main" @click="openProfile = friend">
              <div class="fr-card-avatar">
                {{ friend.username[0].toUpperCase() }}
                <span v-if="friend.online" class="online-dot online-dot--card" />
              </div>
              <div class="fr-card-meta">
                <div class="fr-card-name">{{ friend.username }}</div>
                <div class="fr-card-sub muted">
                  {{ friend.online ? 'Online now' : `Joined ${friend.joined}` }}
                </div>
              </div>
            </button>
            <div class="fr-card-stats">
              <div>
                <span class="mono">{{ friend.score?.toLocaleString() || '0' }}</span>
                <span class="lbl">pts</span>
              </div>
              <div class="fr-card-sep" />
              <div>
                <span class="mono">{{ friend.streak || 0 }}d</span>
                <span class="lbl">streak</span>
              </div>
            </div>
            <button class="btn btn--accent btn--sm btn--block" @click="challengeFriend(friend.username)">
              <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M13 2L4 14h7l-1 8 9-12h-7l1-8z"/>
              </svg>
              Challenge
            </button>
          </div>
          <div v-if="filteredFriends.length === 0" class="fr-empty">
            <div class="muted" style="text-align:center; padding:40px 20px; grid-column:1 / -1">
              {{ searchQuery ? `No friends match "${searchQuery}".` : $t('friends.noFriends') }}
            </div>
          </div>
        </div>
      </section>

      <!-- REQUESTS -->
      <section v-if="tab === 'requests'" class="fr-requests">
        <!-- Incoming -->
        <div class="fr-req-block">
          <div class="fr-req-head">
            <div class="page-eyebrow">Incoming</div>
            <h3 class="fr-req-title">
              {{ incomingRequests.length }} {{ incomingRequests.length === 1 ? 'person wants' : 'people want' }} to play
            </h3>
          </div>
          <div v-if="incomingRequests.length === 0" class="fr-empty-soft muted">
            No incoming requests.
          </div>
          <div v-else class="fr-req-list">
            <div v-for="request in incomingRequests" :key="request.id" class="fr-req fr-req--in">
              <div class="fr-req-who">
                <div class="fr-avatar fr-avatar--lg">{{ request.name[0].toUpperCase() }}</div>
                <div>
                  <div class="fr-req-name">{{ request.name }}</div>
                  <div class="fr-req-meta">
                    <template v-if="request.mutual > 0">
                      <span>
                        <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                          <circle cx="9" cy="7" r="4"/>
                          <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                          <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                        </svg>
                        {{ request.mutual }} mutual
                      </span>
                      <span class="dot-sep">·</span>
                    </template>
                    <span>sent {{ request.sent }}</span>
                  </div>
                  <div v-if="request.msg" class="fr-req-msg">&ldquo;{{ request.msg }}&rdquo;</div>
                </div>
              </div>
              <div class="fr-req-actions">
                <button class="btn btn--ghost btn--sm" @click="declineRequest(request.id, request.name)">Decline</button>
                <button class="btn btn--accent btn--sm" @click="acceptRequest(request.id, request.name)">
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M20 6L9 17l-5-5"/>
                  </svg>
                  Accept
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Outgoing -->
        <div class="fr-req-block">
          <div class="fr-req-head">
            <div class="page-eyebrow">Outgoing</div>
            <h3 class="fr-req-title">{{ outgoingRequests.length }} pending</h3>
          </div>
          <div v-if="outgoingRequests.length === 0" class="fr-empty-soft muted">
            No pending invites.
          </div>
          <div v-else class="fr-req-list">
            <div v-for="request in outgoingRequests" :key="request.id" class="fr-req fr-req--out">
              <div class="fr-req-who">
                <div class="fr-avatar fr-avatar--lg fr-avatar--muted">{{ request.name[0].toUpperCase() }}</div>
                <div>
                  <div class="fr-req-name">{{ request.name }}</div>
                  <div class="fr-req-meta">
                    <template v-if="request.mutual > 0">
                      <span>
                        <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                          <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                          <circle cx="9" cy="7" r="4"/>
                          <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                          <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
                        </svg>
                        {{ request.mutual }} mutual
                      </span>
                      <span class="dot-sep">·</span>
                    </template>
                    <span>sent {{ request.sent }}</span>
                  </div>
                </div>
              </div>
              <div class="fr-req-actions">
                <span class="fr-req-pending">
                  <span class="pulse-dot" /> Awaiting
                </span>
                <button class="btn btn--ghost btn--sm" @click="cancelRequest(request.id, request.name)">Cancel</button>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- ADD FRIENDS -->
      <section v-if="tab === 'add'" class="fr-add">
        <div class="fr-add-grid">
          <!-- Search -->
          <div class="card fr-add-card">
            <div class="page-eyebrow">01 · Search</div>
            <h3 class="fr-add-title">Find by username</h3>
            <p class="muted" style="font-size:13px; margin:0 0 14px">
              Exact match or partial — three letters minimum.
            </p>
            <div class="fr-search-input-wrap">
              <input
                v-model="addSearchQuery"
                @input="searchUsers"
                class="input"
                placeholder="@username"
              />
              <button class="btn btn--accent btn--sm" @click="searchUsers">Search</button>
            </div>

            <!-- Search Results -->
            <div v-if="searchResults.length > 0" style="margin-top:18px">
              <div class="fr-eyebrow-mini">Results</div>
              <div class="fr-suggested">
                <div v-for="user in searchResults" :key="user.id" class="fr-sugg">
                  <div class="fr-sugg-who">
                    <div class="fr-avatar">{{ user.username[0].toUpperCase() }}</div>
                    <div>
                      <div class="fr-sugg-name">{{ user.username }}</div>
                      <div class="fr-sugg-reason muted">{{ user.email }}</div>
                    </div>
                  </div>
                  <button class="btn btn--soft btn--sm" @click="sendRequest(user.id, user.username)">
                    <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M12 5v14M5 12h14"/>
                    </svg>
                    Add
                  </button>
                </div>
              </div>
            </div>

            <!-- Suggested -->
            <div v-if="suggestedFriends.length > 0" style="margin-top:18px">
              <div class="fr-eyebrow-mini">Suggested</div>
              <div class="fr-suggested">
                <div v-for="user in suggestedFriends" :key="user.id" class="fr-sugg">
                  <div class="fr-sugg-who">
                    <div class="fr-avatar">{{ user.name[0].toUpperCase() }}</div>
                    <div>
                      <div class="fr-sugg-name">{{ user.name }}</div>
                      <div class="fr-sugg-reason muted">{{ user.reason }}</div>
                    </div>
                  </div>
                  <button class="btn btn--soft btn--sm" @click="sendRequest(user.id, user.name)">
                    <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                      <path d="M12 5v14M5 12h14"/>
                    </svg>
                    Add
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Invite link -->
          <div class="card fr-add-card fr-add-card--invite">
            <div class="page-eyebrow">02 · Invite</div>
            <h3 class="fr-add-title">Share your link</h3>
            <p class="muted" style="font-size:13px; margin:0 0 14px">
              Anyone with this link can send you a friend request directly.
            </p>
            <div class="fr-link-box">
              <span class="mono">anagrams.ru/u/{{ userStore.username || 'user' }}</span>
              <button class="btn btn--primary btn--sm" @click="copyInviteLink">
                <svg v-if="!inviteCopied" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="9" y="9" width="13" height="13" rx="2"/>
                  <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                </svg>
                <svg v-else width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M20 6L9 17l-5-5"/>
                </svg>
                {{ inviteCopied ? 'Copied' : 'Copy' }}
              </button>
            </div>
            <div class="fr-share-row">
              <button class="btn btn--ghost btn--sm">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="18" cy="5" r="3"/>
                  <circle cx="6" cy="12" r="3"/>
                  <circle cx="18" cy="19" r="3"/>
                  <path d="M8.59 13.51l6.83 3.98M15.41 6.51l-6.82 3.98"/>
                </svg>
                Share
              </button>
              <span class="muted" style="font-size:12px">Link works for 30 days.</span>
            </div>
            <div class="fr-qr">
              <div class="fr-qr-grid">
                <span v-for="i in 49" :key="i" class="fr-qr-cell" :data-on="(i * 7 + 3) % 5 < 2 ? 'true' : 'false'" />
              </div>
              <div class="fr-qr-text muted">QR placeholder</div>
            </div>
          </div>
        </div>
      </section>
    </div>

    <!-- Profile Modal -->
    <Modal :open="!!openProfile" :on-close="() => openProfile = null" :title="openProfile?.username" wide>
      <div v-if="openProfile">
        <div class="fr-prof-head">
          <div class="fr-prof-avatar">
            {{ openProfile.username[0].toUpperCase() }}
            <span v-if="openProfile.online" class="online-dot online-dot--card" />
          </div>
          <div>
            <div class="fr-prof-name">{{ openProfile.username }}</div>
            <div class="muted" style="font-size:13px">
              {{ openProfile.online ? 'Online now' : `Joined ${openProfile.joined}` }} · {{ openProfile.streak }}d streak
            </div>
          </div>
        </div>

        <div class="fr-prof-stats">
          <div class="fr-prof-stat">
            <div class="fr-prof-stat-v mono">{{ openProfile.score?.toLocaleString() }}</div>
            <div class="fr-prof-stat-k">leaderboard score</div>
          </div>
          <div class="fr-prof-stat">
            <div class="fr-prof-stat-v mono">{{ openProfile.bestWord?.toLowerCase() }}</div>
            <div class="fr-prof-stat-k">best word</div>
          </div>
          <div class="fr-prof-stat">
            <div class="fr-prof-stat-v mono">
              {{ loadingFriendMatches ? '...' : friendMatches.length }}
            </div>
            <div class="fr-prof-stat-k">matches with you</div>
          </div>
        </div>

        <!-- Head-to-head Record -->
        <div v-if="loadingFriendMatches" class="fr-prof-loading">
          <span class="muted">Loading match history...</span>
        </div>
        <div v-else-if="friendMatches.length > 0" class="fr-prof-record">
          <div class="page-eyebrow" style="margin-bottom:8px">Head-to-head</div>
          <div class="fr-prof-bar">
            <div class="fr-prof-bar-w" :style="{ flex: getH2HStats().wins || 0.001 }">
              <span v-if="getH2HStats().wins > 0">{{ getH2HStats().wins }} W</span>
            </div>
            <div class="fr-prof-bar-t" :style="{ flex: getH2HStats().ties || 0.001 }">
              <span v-if="getH2HStats().ties > 0">{{ getH2HStats().ties }} T</span>
            </div>
            <div class="fr-prof-bar-l" :style="{ flex: getH2HStats().losses || 0.001 }">
              <span v-if="getH2HStats().losses > 0">{{ getH2HStats().losses }} L</span>
            </div>
          </div>
          <div class="fr-prof-recent">
            <div class="fr-eyebrow-mini">Recent</div>
            <div v-for="match in friendMatches.slice(0, 3)" :key="match.id" :class="['fr-prof-recent-row', `rm-${match.result}`]">
              <span :class="['fr-prof-recent-tag', `rm-${match.result}`]">
                {{ match.result === 'won' ? 'W' : match.result === 'lost' ? 'L' : 'T' }}
              </span>
              <span class="mono fr-prof-recent-letters">{{ match.letters }}</span>
              <span class="mono">{{ match.yourScore.toLocaleString() }} <span class="muted">·</span> {{ match.theirScore.toLocaleString() }}</span>
              <span class="muted" style="margin-left:auto; font-size:12px">{{ match.date }}</span>
            </div>
          </div>
        </div>
        <div v-else class="fr-prof-empty muted" style="padding:20px; text-align:center; font-size:13px">
          No matches played yet
        </div>

        <div class="fr-prof-actions">
          <button class="btn btn--ghost" @click="openProfile = null">Close</button>
          <button class="btn btn--accent" @click="challengeFriend(openProfile.username); openProfile = null">
            <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M13 2L4 14h7l-1 8 9-12h-7l1-8z"/>
            </svg>
            Challenge to a round
          </button>
        </div>
      </div>
    </Modal>
  </div>
</template>

<style scoped>
/* Import all styles from friends.css */
.fr-wrap {
  padding-top: var(--sp-8);
  padding-bottom: var(--sp-12);
}

/* Tabs */
.fr-tabs {
  display: flex;
  gap: 2px;
  border-bottom: 1px solid var(--border-hairline);
  margin-bottom: var(--sp-6);
}

.fr-tab {
  appearance: none;
  border: 0;
  background: transparent;
  padding: 12px 4px;
  margin-right: var(--sp-6);
  font-family: var(--font-body);
  font-size: 13px;
  font-weight: 600;
  color: var(--fg-muted);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  position: relative;
  transition: color var(--dur-base);
  letter-spacing: 0.2px;
}

.fr-tab:hover {
  color: var(--fg-secondary);
}

.fr-tab.is-active {
  color: var(--fg-primary);
}

.fr-tab.is-active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--accent);
}

.fr-tab-count {
  display: inline-grid;
  place-items: center;
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  border-radius: 999px;
  background: var(--bg-card);
  color: var(--fg-secondary);
  font-family: var(--font-mono);
  font-size: 10px;
  font-weight: 700;
}

.fr-tab-count.is-accent {
  background: var(--accent);
  color: var(--milk);
}

/* Search bar */
.fr-search {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
  margin-bottom: var(--sp-5);
}

.fr-search .input {
  flex: 1;
}

.fr-online-count {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-family: var(--font-mono);
  white-space: nowrap;
}

.online-dot {
  width: 7px;
  height: 7px;
  border-radius: 999px;
  background: var(--success);
}

.online-dot--card {
  position: absolute;
  bottom: -1px;
  right: -1px;
  border: 2px solid var(--bg-surface);
  width: 11px;
  height: 11px;
}

/* Friend grid */
.fr-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: var(--sp-3);
}

.fr-card {
  background: var(--bg-surface);
  border: 1px solid var(--border-hairline);
  border-radius: var(--radius-xl);
  padding: var(--sp-4);
  display: flex;
  flex-direction: column;
  gap: var(--sp-3);
  transition: all var(--dur-base) var(--ease-std);
}

.fr-card:hover {
  border-color: var(--border-default);
  transform: translateY(-2px);
  box-shadow: var(--shadow-sm);
}

.fr-card-main {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
  background: transparent;
  border: 0;
  padding: 0;
  text-align: left;
  cursor: pointer;
  font-family: var(--font-body);
}

.fr-card-avatar {
  width: 44px;
  height: 44px;
  position: relative;
  border-radius: 999px;
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  display: grid;
  place-items: center;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 16px;
  color: var(--fg-primary);
  flex-shrink: 0;
}

.fr-card-meta {
  flex: 1;
  min-width: 0;
}

.fr-card-name {
  font-weight: 600;
  font-size: 14px;
  color: var(--fg-primary);
  margin-bottom: 2px;
}

.fr-card-sub {
  font-size: 11px;
  font-family: var(--font-mono);
}

.fr-card-stats {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  padding: var(--sp-2) var(--sp-3);
  background: var(--bg-card);
  border-radius: var(--radius-md);
  font-size: 12px;
}

.fr-card-stats > div {
  display: flex;
  align-items: baseline;
  gap: 4px;
  flex: 1;
}

.fr-card-stats .mono {
  font-weight: 700;
  color: var(--fg-primary);
  font-size: 13px;
}

.fr-card-stats .lbl {
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.6px;
  color: var(--fg-muted);
  font-weight: 600;
}

.fr-card-sep {
  width: 1px;
  height: 20px;
  background: var(--border-subtle);
  flex: 0 0 auto !important;
}

/* Requests panels */
.fr-requests {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--sp-6);
}

@media (max-width: 820px) {
  .fr-requests {
    grid-template-columns: 1fr;
  }
}

.fr-req-block {
  background: var(--bg-surface);
  border: 1px solid var(--border-hairline);
  border-radius: var(--radius-xl);
  padding: var(--sp-5);
}

.fr-req-head {
  margin-bottom: var(--sp-4);
}

.fr-req-title {
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 18px;
  letter-spacing: -0.2px;
  margin: 4px 0 0;
  color: var(--fg-primary);
  text-transform: none;
}

.fr-req-list {
  display: flex;
  flex-direction: column;
  gap: var(--sp-3);
}

.fr-req {
  display: flex;
  align-items: flex-start;
  gap: var(--sp-3);
  padding: var(--sp-3);
  background: var(--bg-card);
  border-radius: var(--radius-md);
  border: 1px solid transparent;
}

.fr-req--in {
  border-color: var(--accent-soft);
}

.fr-req-who {
  display: flex;
  align-items: flex-start;
  gap: var(--sp-3);
  flex: 1;
  min-width: 0;
}

.fr-avatar {
  width: 36px;
  height: 36px;
  border-radius: 999px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  display: grid;
  place-items: center;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 13px;
  color: var(--fg-primary);
  flex-shrink: 0;
}

.fr-avatar--lg {
  width: 42px;
  height: 42px;
  font-size: 15px;
}

.fr-avatar--muted {
  color: var(--fg-muted);
}

.fr-req-name {
  font-weight: 600;
  font-size: 14px;
  color: var(--fg-primary);
}

.fr-req-meta {
  display: flex;
  align-items: center;
  gap: 4px;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--fg-muted);
  margin-top: 2px;
}

.fr-req-msg {
  margin-top: 6px;
  font-size: 12px;
  font-style: italic;
  color: var(--fg-secondary);
  line-height: 1.5;
}

.fr-req-actions {
  display: flex;
  align-items: center;
  gap: var(--sp-2);
  flex-shrink: 0;
}

.fr-req-pending {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  color: var(--fg-muted);
}

.pulse-dot {
  width: 6px;
  height: 6px;
  border-radius: 999px;
  background: var(--warning);
  animation: pulse-anim 1.6s var(--ease-std) infinite;
}

@keyframes pulse-anim {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.4;
    transform: scale(0.8);
  }
}

.fr-empty-soft {
  padding: var(--sp-5);
  text-align: center;
  font-size: 13px;
  background: var(--bg-card);
  border-radius: var(--radius-md);
}

/* Add friends */
.fr-add-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--sp-5);
}

@media (max-width: 820px) {
  .fr-add-grid {
    grid-template-columns: 1fr;
  }
}

.fr-add-card {
  padding: var(--sp-5);
}

.fr-add-title {
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 22px;
  letter-spacing: -0.3px;
  margin: 4px 0 6px;
  color: var(--fg-primary);
  text-transform: none;
}

.fr-search-input-wrap {
  display: flex;
  gap: var(--sp-2);
}

.fr-search-input-wrap .input {
  flex: 1;
}

.fr-eyebrow-mini {
  font-family: var(--font-body);
  font-size: 10px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1.2px;
  color: var(--fg-muted);
  margin-bottom: var(--sp-3);
}

.fr-suggested {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.fr-sugg {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px;
  background: var(--bg-card);
  border-radius: var(--radius-md);
}

.fr-sugg-who {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
}

.fr-sugg-name {
  font-weight: 600;
  font-size: 13px;
}

.fr-sugg-reason {
  font-size: 11px;
}

.fr-link-box {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px var(--sp-3);
  background: var(--bg-card);
  border: 1px dashed var(--border-default);
  border-radius: var(--radius-md);
  margin-bottom: var(--sp-3);
}

.fr-link-box .mono {
  font-size: 13px;
  color: var(--fg-secondary);
}

.fr-share-row {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
  margin-bottom: var(--sp-4);
}

/* QR placeholder */
.fr-qr {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: var(--sp-3);
  padding: var(--sp-3);
  background: var(--bg-card);
  border-radius: var(--radius-md);
}

.fr-qr-grid {
  display: grid;
  grid-template-columns: repeat(7, 18px);
  gap: 2px;
  padding: 8px;
  background: var(--bg-surface);
  border-radius: 6px;
}

.fr-qr-cell {
  width: 18px;
  height: 18px;
  background: transparent;
  border-radius: 2px;
}

.fr-qr-cell[data-on="true"] {
  background: var(--navy);
}

.fr-qr-text {
  margin-top: var(--sp-2);
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 600;
}

/* Profile modal */
.fr-prof-head {
  display: flex;
  align-items: center;
  gap: var(--sp-4);
  padding-bottom: var(--sp-4);
  border-bottom: 1px solid var(--border-hairline);
  margin-bottom: var(--sp-4);
}

.fr-prof-avatar {
  width: 64px;
  height: 64px;
  position: relative;
  border-radius: 999px;
  background: var(--bg-card);
  border: 1px solid var(--border-subtle);
  display: grid;
  place-items: center;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 24px;
  color: var(--fg-primary);
}

.fr-prof-name {
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 22px;
  letter-spacing: -0.3px;
  margin: 0;
}

.fr-prof-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--sp-2);
  margin-bottom: var(--sp-5);
}

.fr-prof-stat {
  background: var(--bg-card);
  border-radius: var(--radius-md);
  padding: var(--sp-3);
  text-align: center;
}

.fr-prof-stat-v {
  font-size: 18px;
  font-weight: 700;
  color: var(--fg-primary);
  margin-bottom: 4px;
}

.fr-prof-stat-k {
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 0.8px;
  color: var(--fg-muted);
  font-weight: 600;
}

.fr-prof-record {
  margin-bottom: var(--sp-5);
}

.fr-prof-bar {
  display: flex;
  height: 28px;
  border-radius: 8px;
  overflow: hidden;
  margin-bottom: var(--sp-4);
  background: var(--bg-card);
}

.fr-prof-bar-w {
  background: var(--success);
}

.fr-prof-bar-l {
  background: var(--danger);
}

.fr-prof-bar-t {
  background: var(--fg-faint);
}

.fr-prof-bar > div {
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--milk);
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.8px;
  min-width: 0;
  overflow: hidden;
  white-space: nowrap;
}

.fr-prof-recent {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.fr-prof-recent-row {
  display: flex;
  align-items: center;
  gap: var(--sp-3);
  padding: 8px 10px;
  background: var(--bg-card);
  border-radius: var(--radius-sm);
  font-size: 13px;
}

.fr-prof-recent-tag {
  width: 22px;
  height: 22px;
  border-radius: 6px;
  display: grid;
  place-items: center;
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  color: var(--milk);
  flex-shrink: 0;
}

.fr-prof-recent-tag.rm-won {
  background: var(--success);
}

.fr-prof-recent-tag.rm-lost {
  background: var(--danger);
}

.fr-prof-recent-tag.rm-tie {
  background: var(--fg-faint);
}

.fr-prof-recent-letters {
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 1px;
  color: var(--fg-secondary);
}

.fr-prof-actions {
  display: flex;
  justify-content: space-between;
  gap: var(--sp-3);
  padding-top: var(--sp-3);
  border-top: 1px solid var(--border-hairline);
}

.dot-sep {
  color: var(--fg-faint);
}
</style>
