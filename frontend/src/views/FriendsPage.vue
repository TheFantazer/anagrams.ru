<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../stores/userStore'

const { t } = useI18n()
const userStore = useUserStore()

const activeTab = ref('friends') // friends, pending, sent, search
const friends = ref([])
const pendingRequests = ref([])
const sentRequests = ref([])
const searchResults = ref([])
const searchQuery = ref('')
const loading = ref(false)
const searching = ref(false)
const userCache = ref({}) // Cache for user data by ID

const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

onMounted(() => {
  if (userStore.isAuthenticated) {
    loadFriends()
    loadPendingRequests()
    loadSentRequests()
  }
})

async function loadFriends() {
  if (!userStore.userId) return

  loading.value = true
  try {
    const response = await fetch(`${apiUrl}/api/v1/friends?user_id=${userStore.userId}`)
    if (response.ok) {
      friends.value = await response.json()
    }
  } catch (error) {
    console.error('Failed to load friends:', error)
  } finally {
    loading.value = false
  }
}

async function getUserInfo(userId) {
  if (userCache.value[userId]) {
    return userCache.value[userId]
  }

  try {
    const response = await fetch(`${apiUrl}/api/v1/users/${userId}`)
    if (response.ok) {
      const user = await response.json()

      userCache.value[userId] = user
      return user
    }
  } catch (error) {
    console.error('Failed to load user info:', error)
  }

  return null
}

async function loadPendingRequests() {
  if (!userStore.userId) return

  try {
    const response = await fetch(`${apiUrl}/api/v1/friends/requests/pending?user_id=${userStore.userId}`)
    if (response.ok) {
      const requests = await response.json()
      // Load user info for each request
      for (const request of requests) {
        request.fromUser = await getUserInfo(request.from_user_id)
      }
      pendingRequests.value = requests
    }
  } catch (error) {
    console.error('Failed to load pending requests:', error)
  }
}

async function loadSentRequests() {
  if (!userStore.userId) return

  try {
    const response = await fetch(`${apiUrl}/api/v1/friends/requests/sent?user_id=${userStore.userId}`)
    if (response.ok) {
      const requests = await response.json()
      // Load user info for each request
      for (const request of requests) {
        request.toUser = await getUserInfo(request.to_user_id)
      }
      sentRequests.value = requests
    }
  } catch (error) {
    console.error('Failed to load sent requests:', error)
  }
}

async function searchUsers() {
  if (!searchQuery.value.trim()) {
    searchResults.value = []
    return
  }

  searching.value = true
  try {
    const response = await fetch(`${apiUrl}/api/v1/users/search?q=${encodeURIComponent(searchQuery.value)}`)
    if (response.ok) {
      const results = await response.json()
      // Исключаем себя из результатов
      searchResults.value = results.filter(user => user.id !== userStore.userId)
    }
  } catch (error) {
    console.error('Failed to search users:', error)
  } finally {
    searching.value = false
  }
}

async function sendFriendRequest(userId) {
  try {
    const response = await fetch(`${apiUrl}/api/v1/friends/requests?user_id=${userStore.userId}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ to_user_id: userId })
    })

    if (response.ok) {
      alert(t('friends.requestSent'))
      await loadSentRequests()
      // Убираем из результатов поиска
      searchResults.value = searchResults.value.filter(u => u.id !== userId)
    } else {
      const error = await response.json()
      alert(error.message || t('friends.requestFailed'))
    }
  } catch (error) {
    console.error('Failed to send friend request:', error)
    alert(t('friends.requestFailed'))
  }
}

async function acceptRequest(requestId) {
  try {
    const response = await fetch(`${apiUrl}/api/v1/friends/requests/${requestId}/accept?user_id=${userStore.userId}`, {
      method: 'POST'
    })

    if (response.ok) {
      alert(t('friends.requestAccepted'))
      await Promise.all([loadFriends(), loadPendingRequests()])
    } else {
      alert(t('friends.requestFailed'))
    }
  } catch (error) {
    console.error('Failed to accept request:', error)
    alert(t('friends.requestFailed'))
  }
}

async function rejectRequest(requestId) {
  try {
    const response = await fetch(`${apiUrl}/api/v1/friends/requests/${requestId}/reject?user_id=${userStore.userId}`, {
      method: 'POST'
    })

    if (response.ok) {
      alert(t('friends.requestRejected'))
      await loadPendingRequests()
    } else {
      alert(t('friends.requestFailed'))
    }
  } catch (error) {
    console.error('Failed to reject request:', error)
    alert(t('friends.requestFailed'))
  }
}

async function removeFriend(friendId) {
  if (!confirm(t('friends.confirmRemove'))) return

  try {
    const response = await fetch(`${apiUrl}/api/v1/friends/${friendId}?user_id=${userStore.userId}`, {
      method: 'DELETE'
    })

    if (response.ok) {
      alert(t('friends.friendRemoved'))
      await loadFriends()
    } else {
      alert(t('friends.removeFailed'))
    }
  } catch (error) {
    console.error('Failed to remove friend:', error)
    alert(t('friends.removeFailed'))
  }
}
</script>

<template>
  <div class="page">
    <div class="shell multi-wrap">
      <header class="page-head" style="justify-content: center">
        <h1 class="page-title-display">{{ t('friends.title') }}</h1>
      </header>

      <!-- Tabs -->
      <div class="tabs">
        <button
          :class="['tab', { active: activeTab === 'friends' }]"
          @click="activeTab = 'friends'"
        >
          {{ t('friends.tabs.friends') }} ({{ friends.length }})
        </button>
        <button
          :class="['tab', { active: activeTab === 'pending' }]"
          @click="activeTab = 'pending'"
        >
          {{ t('friends.tabs.pending') }}
          <span v-if="pendingRequests.length > 0" class="badge">{{ pendingRequests.length }}</span>
        </button>
        <button
          :class="['tab', { active: activeTab === 'sent' }]"
          @click="activeTab = 'sent'"
        >
          {{ t('friends.tabs.sent') }} ({{ sentRequests.length }})
        </button>
        <button
          :class="['tab', { active: activeTab === 'search' }]"
          @click="activeTab = 'search'"
        >
          {{ t('friends.tabs.search') }}
        </button>
      </div>

      <!-- Friends List -->
      <div v-if="activeTab === 'friends'" class="tab-content">
        <div v-if="loading" class="loading">{{ t('common.loading') }}</div>
        <div v-else-if="friends.length === 0" class="empty-state">
          {{ t('friends.noFriends') }}
        </div>
        <div v-else class="friends-list">
          <div v-for="friend in friends" :key="friend.id" class="friend-card">
            <div class="friend-info">
              <div class="friend-name">{{ friend.username }}</div>
              <div class="friend-email">{{ friend.email }}</div>
            </div>
            <button @click="removeFriend(friend.id)" class="btn-remove">
              {{ t('friends.remove') }}
            </button>
          </div>
        </div>
      </div>

      <!-- Pending Requests -->
      <div v-if="activeTab === 'pending'" class="tab-content">
        <div v-if="pendingRequests.length === 0" class="empty-state">
          {{ t('friends.noPending') }}
        </div>
        <div v-else class="requests-list">
          <div v-for="request in pendingRequests" :key="request.id" class="request-card">
            <div class="request-info">
              <div class="request-label">{{ t('friends.requestFrom') }}</div>
              <div class="request-user">{{ request.fromUser?.username || request.from_user_id }}</div>
            </div>
            <div class="request-actions">
              <button @click="acceptRequest(request.id)" class="btn btn--sm btn--success">
                {{ t('friends.accept') }}
              </button>
              <button @click="rejectRequest(request.id)" class="btn btn--sm btn--danger">
                {{ t('friends.reject') }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Sent Requests -->
      <div v-if="activeTab === 'sent'" class="tab-content">
        <div v-if="sentRequests.length === 0" class="empty-state">
          {{ t('friends.noSent') }}
        </div>
        <div v-else class="requests-list">
          <div v-for="request in sentRequests" :key="request.id" class="request-card">
            <div class="request-info">
              <div class="request-label">{{ t('friends.requestTo') }}</div>
              <div class="request-user">{{ request.toUser?.username || request.to_user_id }}</div>
              {{ console.log(request.toUser)
              }}
              <div class="request-status" :class="request.status">
                {{ t(`friends.status.${request.status}`) }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Search Users -->
      <div v-if="activeTab === 'search'" class="tab-content">
        <div class="search-box">
          <input
            v-model="searchQuery"
            @input="searchUsers"
            type="text"
            :placeholder="t('friends.searchPlaceholder')"
            class="search-input"
          />
        </div>

        <div v-if="searching" class="loading">{{ t('common.searching') }}</div>
        <div v-else-if="searchQuery && searchResults.length === 0" class="empty-state">
          {{ t('friends.noResults') }}
        </div>
        <div v-else-if="searchResults.length > 0" class="search-results">
          <div v-for="user in searchResults" :key="user.id" class="user-card">
            <div class="user-info">
              <div class="user-name">{{ user.username }}</div>
              <div class="user-email">{{ user.email }}</div>
            </div>
            <button @click="sendFriendRequest(user.id)" class="btn-add">
              {{ t('friends.addFriend') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 32px;
  background: var(--bg-card);
  padding: 8px;
  border-radius: 12px;
  flex-wrap: wrap;
}

.tab {
  flex: 1;
  min-width: 120px;
  padding: 12px 16px;
  background: transparent;
  border: none;
  color: var(--fg-secondary);
  font-family: var(--font-body);
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  border-radius: 8px;
  transition: all var(--dur-base) var(--ease-std);
  position: relative;
}

.tab:hover {
  background: var(--bg-hover);
  color: var(--fg-primary);
}

.tab.active {
  background: var(--navy);
  color: var(--milk);
}

.badge {
  display: inline-block;
  background: var(--danger);
  color: var(--milk);
  padding: 2px 8px;
  border-radius: 12px;
  font-size: 11px;
  margin-left: 4px;
}

.tab-content {
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 16px;
  padding: 32px;
  min-height: 400px;
}

.loading, .empty-state {
  text-align: center;
  color: var(--fg-muted);
  padding: 48px 16px;
  font-size: 15px;
  font-family: var(--font-body);
}

.friends-list, .requests-list, .search-results {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.friend-card, .request-card, .user-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: var(--bg-card);
  border: 1px solid var(--border-hairline);
  border-radius: 12px;
  transition: all var(--dur-base) var(--ease-std);
}

.friend-card:hover, .request-card:hover, .user-card:hover {
  border-color: var(--border-default);
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

.friend-info, .request-info, .user-info {
  flex: 1;
}

.friend-name, .user-name {
  font-family: var(--font-body);
  font-size: 16px;
  font-weight: 600;
  color: var(--fg-primary);
  margin-bottom: 4px;
}

.friend-email, .user-email, .request-user {
  font-family: var(--font-body);
  font-size: 14px;
  color: var(--fg-secondary);
}

.request-label {
  font-family: var(--font-body);
  font-size: 11px;
  color: var(--fg-muted);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  margin-bottom: 4px;
}

.request-status {
  display: inline-block;
  margin-top: 8px;
  padding: 4px 12px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  font-family: var(--font-body);
}

.request-status.pending {
  background: var(--warning);
  color: var(--milk);
}

.request-status.accepted {
  background: var(--success);
  color: var(--milk);
}

.request-status.rejected {
  background: var(--danger-soft);
  color: var(--danger);
  border: 1px solid var(--danger-border);
}

.request-actions {
  display: flex;
  gap: 8px;
}

.btn--success {
  background: var(--success);
  color: var(--milk);
}

.btn--success:hover {
  background: var(--success);
  opacity: 0.9;
}

.btn--danger {
  background: var(--danger-soft);
  color: var(--danger);
  border: 1px solid var(--danger-border);
}

.btn--danger:hover {
  background: var(--danger);
  color: var(--milk);
}

.search-box {
  margin-bottom: 24px;
}

.search-input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid var(--border-default);
  border-radius: 10px;
  font-family: var(--font-body);
  font-size: 15px;
  color: var(--fg-primary);
  background: var(--bg-surface);
  transition: all var(--dur-base) var(--ease-std);
}

.search-input:focus {
  outline: none;
  border-color: var(--navy);
  box-shadow: 0 0 0 3px var(--accent-glow);
}

.search-input::placeholder {
  color: var(--fg-muted);
}

@media (max-width: 640px) {
  .friend-card, .request-card, .user-card {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .request-actions {
    width: 100%;
  }

  .request-actions button {
    flex: 1;
  }
}
</style>
