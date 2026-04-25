<script setup>
import {ref, computed, onMounted, onActivated, watch} from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../stores/userStore'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const linkCopied = ref(false)
const creating = ref(false)
const createdSessionId = ref(null)
const sessionLetters = ref([])
const activeChallenges = ref([])
const loadingChallenges = ref(false)
const showFriendsModal = ref(false)
const friends = ref([])
const selectedFriends = ref([])
const loadingFriends = ref(false)

// Form settings
const language = ref('ru')
const letterCount = ref(7)
const timeLimit = ref(60)
const hideLetters = ref(false)

const perPage = 15
const currentPage = ref(1)
const totalChallenges = ref(0)

const totalPages = computed(() => {
  return Math.max(1, Math.ceil(totalChallenges.value / perPage))
})

const pages = computed(() => {
  const maxVisible = 10
  const half = Math.floor(maxVisible / 2)

  let start = Math.max(1, currentPage.value - half)
  let end = Math.min(totalPages.value, currentPage.value + half)

  if (end - start < maxVisible - 1) {
    if (start === 1) {
      end = Math.min(totalPages.value, start + maxVisible - 1)
    } else {
      start = Math.max(1, end - maxVisible + 1)
    }
  }

  return Array.from({ length: end - start + 1 }, (_, i) => start + i)
})

const availableLanguages = computed(() => [
  { id: 'ru', label: t('settings.gameDefaults.languages.ru') },
  { id: 'en', label: t('settings.gameDefaults.languages.en') }
])

const letterCounts = [
  { value: 5, label: '5' },
  { value: 6, label: '6' },
  { value: 7, label: '7' },
  { value: 8, label: '8' },
  { value: 9, label: '9' }
]

const timeLimits = [
  { value: 30, label: '30s' },
  { value: 60, label: '1:00' },
  { value: 90, label: '1:30' },
  { value: 120, label: '2:00' }
]

async function changePage(newPage) {
  if (newPage < 1 || newPage > totalPages.value) return
  await loadActiveChallenges(newPage)
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

async function sendInvites() {
  if (selectedFriends.value.length === 0) return

  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

    // Create a separate session for each selected friend (strict 1v1)
    const sessionPromises = selectedFriends.value.map(async (friendId) => {
      // Create session with invite_mode = "friend" and max_opponents = 1
      let url = `${apiUrl}/api/v1/sessions`
      if (userStore.userId) {
        url += `?user_id=${userStore.userId}`
      }

      const sessionResponse = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          language: language.value,
          letter_count: letterCount.value,
          time_limit: timeLimit.value,
          hide_letters: hideLetters.value,
          invite_mode: 'friend',
          max_opponents: 1
        })
      })

      if (!sessionResponse.ok) {
        throw new Error('Failed to create session')
      }

      const session = await sessionResponse.json()

      // Send invite for this session to the specific friend
      await fetch(`${apiUrl}/api/v1/sessions/${session.id}/invites?user_id=${userStore.userId}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ invited_user_id: friendId })
      })

      return session
    })

    await Promise.all(sessionPromises)

    showFriendsModal.value = false
    selectedFriends.value = []
    resetForm()
    alert(t('multiplayer.invitesSent'))

    // Reload challenges to show the newly created ones
    if (userStore.isAuthenticated) {
      await loadActiveChallenges()
    }
  } catch (error) {
    console.error('Failed to send invites:', error)
    alert(t('multiplayer.invitesFailed'))
  }
}

function openFriendsModal() {
  loadFriends()
  showFriendsModal.value = true
}

function toggleFriendSelection(friendId) {
  const index = selectedFriends.value.indexOf(friendId)
  if (index > -1) {
    selectedFriends.value.splice(index, 1)
  } else {
    selectedFriends.value.push(friendId)
  }
}

async function createChallenge() {
  creating.value = true
  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

    // Build URL with optional user_id query param
    let url = `${apiUrl}/api/v1/sessions`
    if (userStore.isAuthenticated && userStore.userId) {
      url += `?user_id=${userStore.userId}`
    }

    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        language: language.value,
        letter_count: letterCount.value,
        time_limit: timeLimit.value,
        hide_letters: hideLetters.value,
        invite_mode: 'link',
        max_opponents: 1
      })
    })

    if (!response.ok) {
      throw new Error('Failed to create session')
    }

    const session = await response.json()
    createdSessionId.value = session.id
    sessionLetters.value = session.letters.toUpperCase().split('')

    // Reload challenges after creating a new one
    if (userStore.isAuthenticated) {
      await loadActiveChallenges()
    }
  } catch (error) {
    console.error('Failed to create challenge:', error)
    alert('Failed to create challenge. Please try again.')
  } finally {
    creating.value = false
  }
}

async function createChallengeForFriend() {
  // Don't create a session yet - just open the modal
  // Sessions will be created in sendInvites (one per friend for strict 1v1)
  openFriendsModal()
}

async function loadActiveChallenges(page = 1) {
  if (!userStore.isAuthenticated || !userStore.userId) {
    return
  }

  loadingChallenges.value = true
  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

    const response = await fetch(`${apiUrl}/api/v1/sessions/all?user_id=${userStore.userId}&page=${page}&per_page=${perPage}`)

    if (!response.ok) {
      throw new Error('Failed to load sessions')
    }

    const data = await response.json()

    activeChallenges.value = data.sessions || []

    // Обновляем пагинацию на основе server response
    currentPage.value = data.page
    totalChallenges.value = data.total
  } catch (error) {
    console.error('Failed to load challenges:', error)
  } finally {
    loadingChallenges.value = false
  }
}

const shareLink = computed(() => {
  if (!createdSessionId.value) return ''
  return `${window.location.origin}/play/${createdSessionId.value}`
})

function copyLink() {
  if (!shareLink.value) return
  navigator.clipboard?.writeText(shareLink.value)
  linkCopied.value = true
  setTimeout(() => linkCopied.value = false, 2200)
}

function playChallenge() {
  if (!createdSessionId.value) return
  router.push(`/play/${createdSessionId.value}`)
}

function resetForm() {
  createdSessionId.value = null
  sessionLetters.value = []
  linkCopied.value = false
  hideLetters.value = false
}

function formatDate(dateStr) {
  const date = new Date(dateStr)
  const now = new Date()
  const diffMs = now - date
  const diffMins = Math.floor(diffMs / 60000)
  const diffHours = Math.floor(diffMs / 3600000)
  const diffDays = Math.floor(diffMs / 86400000)

  if (diffMins < 60) {
    return `${diffMins}m ${t('multiplayer.ago')}`
  } else if (diffHours < 24) {
    return `${diffHours}h ${t('multiplayer.ago')}`
  } else if (diffDays === 1) {
    return t('multiplayer.yesterday')
  } else if (diffDays < 7) {
    return `${diffDays}d ${t('multiplayer.ago')}`
  } else {
    return date.toLocaleDateString()
  }
}

function getStatusClass(challenge) {
  // Проверяем, сыграл ли текущий пользователь
  const userResult = challenge.results?.find(r => r.user_id === userStore.userId)

  if (challenge.type === 'created') {
    // Для созданных челленджей
    if (!challenge.results || challenge.results.length === 0) {
      return 'st-wait' // Никто еще не играл
    }
    // Если кто-то играл, показываем статус создателя
    if (userResult) {
      const topScore = Math.max(...challenge.results.map(r => r.score))
      return userResult.score === topScore ? 'st-won' : 'st-lost'
    }
    return 'st-wait' // Создатель еще не сыграл
  } else {
    // Для челленджей, в которых пригласили
    if (userResult) {
      // Пользователь уже сыграл
      const topScore = Math.max(...challenge.results.map(r => r.score))
      return userResult.score === topScore ? 'st-won' : 'st-lost'
    }
    return 'st-play' // Еще не сыграл
  }
}

function getStatusLabel(challenge) {
  // Проверяем, сыграл ли текущий пользователь
  const userResult = challenge.results?.find(r => r.user_id === userStore.userId)

  if (challenge.type === 'created') {
    // Для созданных челленджей
    if (!challenge.results || challenge.results.length === 0) {
      return t('multiplayer.waitingForPlayers')
    }
    if (userResult) {
      const topScore = Math.max(...challenge.results.map(r => r.score))
      return userResult.score === topScore
        ? t('multiplayer.youWon')
        : t('multiplayer.youLost')
    }
    // Если создатель не играл, но есть результаты (друг сыграл)
    return t('multiplayer.yourTurn')
  } else {
    // Для челленджей, в которых пригласили
    if (userResult) {
      const topScore = Math.max(...challenge.results.map(r => r.score))
      return userResult.score === topScore
        ? t('multiplayer.youWon')
        : t('multiplayer.youLost')
    }
    return t('multiplayer.yourTurn')
  }
}

function hasUserPlayed(challenge) {
  return challenge.results?.some(r => r.user_id === userStore.userId) || false
}

function getMyScore(challenge) {
  const myResult = challenge.results?.find(r => r.user_id === userStore.userId)
  return myResult?.score
}

function getOpponentScore(challenge) {
  const opponentResult = challenge.results?.find(r => r.user_id !== userStore.userId)
  return opponentResult?.score
}

onMounted(() => {
  // Загружаем активные челленджи
  loadActiveChallenges()

  // Проверяем query параметры для автоматического создания челленджа
  if (route.query.create === 'true') {
    if (route.query.language) {
      language.value = route.query.language
    }
    if (route.query.letterCount) {
      letterCount.value = parseInt(route.query.letterCount)
    }
    if (route.query.timeLimit) {
      timeLimit.value = parseInt(route.query.timeLimit)
    }
    if (route.query.hideLetters) {
      hideLetters.value = route.query.hideLetters === 'true'
    }
    // Автоматически создаём челлендж
    createChallenge()
  }
})

// Обновляем список при возвращении на страницу
onActivated(() => {
  loadActiveChallenges()
})

// Обновляем список когда переходим на страницу мультиплеера
watch(() => route.path, (newPath) => {
  if (newPath === '/multiplayer') {
    loadActiveChallenges()
  }
})
</script>

<template>
  <div class="page">
    <div class="shell multi-wrap">
      <header class="page-head">
        <div>
          <div class="page-eyebrow">{{ $t('multiplayer.title') }}</div>
          <h1 class="page-title-display">{{ $t('multiplayer.subtitle') }}</h1>
        </div>
      </header>


      <div class="multi-grid">
        <!-- Create Challenge Section -->
        <section class="card card--paper multi-new">
          <div class="multi-eye">
            <span class="multi-num">01</span>
            {{$t('multiplayer.card01.title')}}
          </div>
          <h3 style="font-family:var(--font-display);font-size:28px;font-weight:700;letter-spacing:-0.5px;margin:6px 0 10px;color:var(--fg-primary);text-transform:none">
            {{ $t('multiplayer.card01.header') }}
          </h3>
          <p class="muted" style="margin:0 0 20px;font-size:13px;max-width:360px">
            {{ $t('multiplayer.card01.subtitle') }}
          </p>

          <div v-if="!createdSessionId" class="multi-set">
            <span v-for="i in letterCount" :key="i" class="multi-tile">?</span>
          </div>
          <div v-else class="multi-set">
            <span v-for="(letter, i) in sessionLetters" :key="i" class="multi-tile">
              {{ hideLetters ? '?' : letter }}
            </span>
          </div>

          <div v-if="!createdSessionId" class="row gap-2" style="margin-top:16px;flex-wrap:wrap">
            <button
              v-for="count in letterCounts"
              :key="count.value"
              class="btn btn--soft btn--sm"
              :class="{ 'btn--active': letterCount === count.value }"
              @click="letterCount = count.value"
            >
              {{ count.label }} {{ $t('multiplayer.letters') }}
            </button>
            <div style="width:100%"></div>
            <button
              v-for="time in timeLimits"
              :key="time.value"
              class="btn btn--soft btn--sm"
              :class="{ 'btn--active': timeLimit === time.value }"
              @click="timeLimit = time.value"
            >
              {{ time.label }}
            </button>
            <button
              v-for="lang in availableLanguages"
              :key="lang.id"
              class="btn btn--soft btn--sm"
              :class="{ 'btn--active': language === lang.id }"
              @click="language = lang.id"
            >
              {{ lang.label }}
            </button>
            <div style="width:100%"></div>
            <button
              class="btn btn--soft btn--sm"
              :class="{ 'btn--active': hideLetters }"
              @click="hideLetters = !hideLetters"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <template v-if="hideLetters">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </template>
                <template v-else>
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                  <line x1="1" y1="1" x2="23" y2="23"/>
                </template>
              </svg>
              {{ hideLetters ? $t('multiplayer.showLetters') : $t('multiplayer.hideLetters') }}
            </button>
          </div>
          <div v-else class="row gap-2" style="margin-top:16px;flex-wrap:wrap">
            <button class="btn btn--soft btn--sm" disabled>{{ letterCount }} {{ $t('multiplayer.letters') }}</button>
            <button class="btn btn--soft btn--sm" disabled>{{ timeLimits.find(t => t.value === timeLimit)?.label }}</button>
            <button class="btn btn--soft btn--sm" disabled>{{ availableLanguages.find(l => l.id === language)?.label }}</button>
            <button
              class="btn btn--soft btn--sm"
              :class="{ 'btn--active': hideLetters }"
              @click="hideLetters = !hideLetters"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <template v-if="hideLetters">
                  <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24"/>
                  <line x1="1" y1="1" x2="23" y2="23"/>
                </template>
                <template v-else>
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </template>
              </svg>
              {{ hideLetters ? $t('multiplayer.showLetters') : $t('multiplayer.hideLetters') }}
            </button>
          </div>

          <div v-if="createdSessionId" class="multi-link">
            <span class="mono" style="color:var(--fg-secondary);word-break:break-all">{{ shareLink }}</span>
            <button class="btn btn--primary btn--sm" @click="copyLink">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <template v-if="!linkCopied">
                  <rect x="9" y="9" width="13" height="13" rx="2"/>
                  <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                </template>
                <path v-else d="M20 6L9 17l-5-5"/>
              </svg>
              {{ linkCopied ? $t('multiplayer.card01.copied') : $t('multiplayer.card01.copyLink') }}
            </button>
          </div>

          <div class="row gap-2" style="margin-top:16px;flex-wrap:wrap">
            <button
              v-if="!createdSessionId"
              class="btn btn--accent"
              @click="createChallenge"
              :disabled="creating"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"/>
                <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"/>
              </svg>
              {{ creating ? $t('common.creating') : $t('multiplayer.createWithLink') }}
            </button>
            <button
              v-if="!createdSessionId"
              class="btn btn--primary"
              @click="createChallengeForFriend"
              :disabled="creating"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/>
                <circle cx="9" cy="7" r="4"/>
                <path d="M23 21v-2a4 4 0 0 0-3-3.87"/>
                <path d="M16 3.13a4 4 0 0 1 0 7.75"/>
              </svg>
              {{ creating ? $t('common.creating') : $t('multiplayer.createForFriend') }}
            </button>
            <template v-else>
              <button class="btn btn--accent" @click="playChallenge">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M6 4l14 8-14 8z"/>
                </svg>
                {{ $t('multiplayer.card01.playBtn')}}
              </button>
              <button class="btn btn--ghost" @click="resetForm">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/>
                  <path d="M21 3v5h-5"/>
                  <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/>
                  <path d="M3 21v-5h5"/>
                </svg>
                {{ $t('common.newChallenge') }}
              </button>
            </template>
          </div>
        </section>

        <!-- How It Works Section -->
        <section class="multi-how">
          <div class="multi-eye"><span class="multi-num">02</span>{{ $t('multiplayer.card02.title') }}</div>
          <ol class="multi-steps">
            <li><b>{{ $t('multiplayer.card02.step1.title') }}.</b> {{ $t('multiplayer.card02.step1.description') }}</li>
            <li><b>{{ $t('multiplayer.card02.step2.title') }}.</b> {{ $t('multiplayer.card02.step2.description') }}</li>
            <li><b>{{ $t('multiplayer.card02.step3.title') }}.</b> {{ $t('multiplayer.card02.step3.description') }}</li>
            <li><b>{{ $t('multiplayer.card02.step4.title') }}</b> {{ $t('multiplayer.card02.step4.description') }}</li>
          </ol>
          <div class="multi-tip">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <circle cx="12" cy="12" r="10"/>
              <path d="M12 6v6l4 2"/>
            </svg>
            <span>{{$t('multiplayer.card02.challengesOpenFor')}}</span>
          </div>
        </section>
      </div>

      <!-- Active Challenges Section -->
      <section style="margin-top:32px">
        <div class="multi-eye" style="margin-bottom:16px;display:flex;justify-content:space-between;align-items:center">
          <div>
            <span class="multi-num">03</span>{{$t('multiplayer.card03.title')}}
          </div>
          <button
            v-if="userStore.isAuthenticated"
            class="btn btn--ghost btn--sm"
            @click="loadActiveChallenges()"
            :disabled="loadingChallenges"
          >
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/>
              <path d="M21 3v5h-5"/>
              <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/>
              <path d="M3 21v-5h5"/>
            </svg>
            {{ $t('multiplayer.card03.refresh') }}
          </button>
        </div>

        <div v-if="loadingChallenges" class="empty-state">
          <p class="muted">{{$t('common.loading')}}</p>
        </div>

        <div v-else-if="!userStore.isAuthenticated" class="empty-state">
          <p class="muted">Sign in to view your challenges</p>
        </div>

        <div v-else-if="activeChallenges.length === 0" class="empty-state">
          <p class="muted">{{$t('multiplayer.card03.description')}}</p>
        </div>

        <div v-else class="multi-challenges">
          <div
            v-for="challenge in activeChallenges"
            :key="challenge.id"
            class="ch-row"
            :class="{ 'ch-play-now': getStatusClass(challenge) === 'st-play' }"
            @click="router.push(hasUserPlayed(challenge) ? `/results/${challenge.id}` : `/play/${challenge.id}`)"
          >
            <!-- Creator/Type -->
            <div class="ch-who">
              <div class="ch-avatar">
                {{ challenge.type === 'created'
                  ? userStore.username?.charAt(0).toUpperCase()
                  : (challenge.creator_username?.charAt(0).toUpperCase() || '?') }}
              </div>
              <div>
                <div class="ch-name">
                  {{ challenge.type === 'created'
                    ? $t('multiplayer.createdByYou')
                    : (challenge.creator_username ? `${$t('multiplayer.from')} ${challenge.creator_username}` : $t('multiplayer.invited')) }}
                </div>
                <div class="ch-meta">{{ formatDate(challenge.created_at) }}</div>
              </div>
            </div>

            <!-- Letters -->
            <div class="ch-letters">
              <span v-for="(letter, i) in challenge.letters.split('')" :key="i" class="ch-tile">
                {{ (challenge.hide_letters && (!challenge.results || challenge.results.length < 2)) ? '?' : letter.toLowerCase() }}
              </span>
            </div>

            <!-- Scores -->
            <div class="ch-scores">
              <template v-if="challenge.hide_letters && (!challenge.results || challenge.results.length < 2)">
                <span class="muted">—</span>
              </template>
              <template v-else-if="challenge.results && challenge.results.length > 0">
                <span class="mono">{{ getMyScore(challenge)?.toLocaleString() || '—' }}</span>
                <span class="muted"> {{ $t('multiplayer.vs') }} </span>
                <span class="mono">{{ getOpponentScore(challenge)?.toLocaleString() || '—' }}</span>
              </template>
              <span v-else class="muted">—</span>
            </div>

            <!-- Status -->
            <span class="ch-status" :class="getStatusClass(challenge)">
              {{ getStatusLabel(challenge) }}
            </span>

            <!-- Action Button -->
            <button class="btn btn--sm btn--primary" @click.stop="router.push(hasUserPlayed(challenge) ? `/results/${challenge.id}` : `/play/${challenge.id}`)">
              {{ hasUserPlayed(challenge) ? $t('multiplayer.view') : $t('multiplayer.play') }}
            </button>
          </div>
          <div class="pagination">
            <button
                class="btn btn--ghost btn--sm"
                @click="changePage(currentPage - 1)"
                :disabled="currentPage === 1"
            >
              ←
            </button>

            <button
                v-for="page in pages"
                :key="page"
                class="btn btn--soft btn--sm"
                :class="{ 'btn--accent': currentPage === page }"
                @click="changePage(page)"
            >
              {{ page }}
            </button>

            <button
                class="btn btn--ghost btn--sm"
                @click="changePage(currentPage + 1)"
                :disabled="currentPage === totalPages"
            >
              →
            </button>
          </div>
        </div>
      </section>

      <!-- CTA to go back -->
      <div style="margin-top:32px;text-align:center">
        <button class="btn btn--primary btn--lg" @click="router.push('/')">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M19 12H5M12 19l-7-7 7-7"/>
          </svg>
          {{$t('multiplayer.backBtn')}}
        </button>
      </div>
    </div>

    <!-- Friends Modal -->
    <div v-if="showFriendsModal" class="modal-overlay" @click="showFriendsModal = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h2>{{ $t('multiplayer.inviteFriends') }}</h2>
          <button class="modal-close" @click="showFriendsModal = false">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"/>
              <line x1="6" y1="6" x2="18" y2="18"/>
            </svg>
          </button>
        </div>

        <div class="modal-body">
          <div v-if="loadingFriends" class="modal-loading">
            {{ $t('common.loading') }}
          </div>

          <div v-else-if="friends.length === 0" class="modal-empty">
            <p class="muted">{{ $t('multiplayer.noFriends') }}</p>
            <button class="btn btn--primary btn--sm" @click="router.push('/friends'); showFriendsModal = false">
              {{ $t('multiplayer.goToFriends') }}
            </button>
          </div>

          <div v-else class="friends-list">
            <div
              v-for="friend in friends"
              :key="friend.id"
              class="friend-item"
              :class="{ selected: selectedFriends.includes(friend.id) }"
              @click="toggleFriendSelection(friend.id)"
            >
              <div class="friend-info">
                <div class="friend-name">{{ friend.username }}</div>
                <div class="friend-email">{{ friend.email }}</div>
              </div>
              <div class="friend-checkbox">
                <svg v-if="selectedFriends.includes(friend.id)" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                  <polyline points="20 6 9 17 4 12"/>
                </svg>
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button class="btn btn--ghost" @click="showFriendsModal = false">
            {{ $t('common.cancel') }}
          </button>
          <button
            class="btn btn--accent"
            @click="sendInvites"
            :disabled="selectedFriends.length === 0"
          >
            {{ $t('multiplayer.sendInvites') }} ({{ selectedFriends.length }})
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.multi-wrap { max-width: 980px; margin: 0 auto; }

.coming-soon-banner {
  background: var(--warning);
  color: var(--navy);
  padding: 16px 20px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 32px;
  box-shadow: var(--shadow-md);
}

.coming-soon-banner strong {
  display: block;
  font-size: 14px;
  font-weight: 700;
  margin-bottom: 2px;
}

.coming-soon-banner .muted {
  font-size: 13px;
  color: var(--navy);
  opacity: 0.7;
}

.page-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 20px;
  margin-bottom: 28px;
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
  font-size: 40px;
  font-weight: 700;
  letter-spacing: -1.2px;
  color: var(--fg-primary);
  margin: 0;
  line-height: 1;
}

.multi-grid {
  display: grid;
  grid-template-columns: 1.3fr 1fr;
  gap: 16px;
}

.multi-eye {
  display: inline-flex; align-items: center; gap: 10px;
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: 2px;
  font-weight: 600;
  color: var(--fg-muted);
  margin-bottom: 12px;
}

.multi-num {
  font-family: var(--font-mono);
  padding: 3px 8px;
  background: var(--bg-surface);
  border: 1px solid var(--border-hairline);
  border-radius: 6px;
  color: var(--fg-primary);
}

.multi-new {
  padding: 32px;
}

.multi-set {
  display: flex; gap: 10px;
  margin-top: 8px;
  flex-wrap: wrap;
}

.multi-tile {
  width: 48px; height: 54px;
  border-radius: 12px;
  background: var(--navy);
  color: var(--milk);
  display: grid; place-items: center;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 22px;
  box-shadow: 0 3px 0 var(--navy-2);
}

.multi-link {
  margin-top: 20px;
  padding: 14px 16px;
  background: var(--bg-surface);
  border: 1px dashed var(--border-default);
  border-radius: 12px;
  display: flex; align-items: center; justify-content: space-between; gap: 12px;
  font-size: 13px;
  flex-wrap: wrap;
}

.multi-how {
  padding: 28px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 20px;
  height: fit-content;
}

.multi-steps {
  list-style: none;
  padding: 0;
  margin: 0;
  counter-reset: step;
}

.multi-steps li {
  counter-increment: step;
  padding: 12px 0 12px 40px;
  border-bottom: 1px solid var(--border-hairline);
  position: relative;
  font-size: 14px;
  color: var(--fg-secondary);
}

.multi-steps li:last-child { border-bottom: 0; }

.multi-steps li::before {
  content: counter(step, decimal-leading-zero);
  position: absolute; left: 0; top: 14px;
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 700;
  color: var(--accent);
}

.multi-steps li b { color: var(--fg-primary); font-weight: 600; }

.multi-tip {
  margin-top: 16px;
  padding: 12px 14px;
  background: var(--bg-card);
  border-radius: 10px;
  display: flex; gap: 10px; align-items: center;
  color: var(--fg-secondary);
  font-size: 12px;
}

.empty-state {
  padding: 40px 20px;
  text-align: center;
  background: var(--bg-surface);
  border: 1px dashed var(--border-default);
  border-radius: 14px;
}

.btn--active {
  background: var(--navy);
  color: var(--milk);
  border-color: var(--navy);
}

/* Challenges Table */
.multi-challenges {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.ch-row {
  display: grid;
  grid-template-columns: 160px 1fr 140px 140px auto;
  gap: 16px;
  align-items: center;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 14px;
  padding: 14px 18px;
  transition: all var(--dur-base);
  cursor: pointer;
}

.ch-row:hover {
  border-color: var(--border-default);
}

.ch-play-now {
  border-left: 3px solid var(--accent);
}

.ch-who {
  display: flex;
  align-items: center;
  gap: 10px;
}

.ch-avatar {
  width: 32px;
  height: 32px;
  border-radius: 999px;
  background: var(--grad-accent);
  color: var(--milk);
  display: grid;
  place-items: center;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 13px;
}

.ch-name {
  font-weight: 600;
  font-size: 14px;
  color: var(--fg-primary);
}

.ch-meta {
  font-size: 11px;
  color: var(--fg-muted);
}

.ch-letters {
  display: flex;
  gap: 4px;
}

.ch-tile {
  width: 22px;
  height: 26px;
  background: var(--bg-card);
  border: 1px solid var(--border-hairline);
  border-radius: 5px;
  display: grid;
  place-items: center;
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  color: var(--fg-primary);
}

.ch-scores {
  font-size: 13px;
  color: var(--fg-primary);
}

.ch-status {
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
  padding: 4px 10px;
  border-radius: 999px;
  background: var(--bg-card);
  color: var(--fg-muted);
  white-space: nowrap;
}

.st-wait {
  background: var(--bg-card);
  color: var(--fg-muted);
}

.st-play {
  background: var(--accent);
  color: var(--milk);
}

.st-won {
  background: var(--success-soft);
  color: var(--success);
}

.st-lost {
  background: var(--danger-soft);
  color: var(--danger);
}

/* Friends Modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(30, 42, 70, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
  backdrop-filter: blur(4px);
}

.modal-content {
  background: var(--bg-surface);
  border-radius: 20px;
  max-width: 500px;
  width: 100%;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  box-shadow: var(--shadow-lg);
  border: 1px solid var(--border-subtle);
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 28px;
  border-bottom: 1px solid var(--border-hairline);
}

.modal-header h2 {
  font-family: var(--font-display);
  font-size: 22px;
  font-weight: 700;
  color: var(--fg-primary);
  margin: 0;
}

.modal-close {
  background: transparent;
  border: none;
  color: var(--fg-muted);
  cursor: pointer;
  padding: 4px;
  border-radius: 8px;
  transition: all var(--dur-base) var(--ease-std);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-close:hover {
  background: var(--bg-hover);
  color: var(--fg-primary);
}

.modal-body {
  padding: 20px 28px;
  overflow-y: auto;
  flex: 1;
}

.modal-loading,
.modal-empty {
  text-align: center;
  padding: 40px 20px;
  color: var(--fg-muted);
}

.modal-empty p {
  margin-bottom: 16px;
}

.friends-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.friend-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  background: var(--bg-card);
  border: 2px solid transparent;
  border-radius: 12px;
  cursor: pointer;
  transition: all var(--dur-base) var(--ease-std);
}

.friend-item:hover {
  background: var(--bg-hover);
  border-color: var(--border-default);
}

.friend-item.selected {
  background: var(--accent-soft);
  border-color: var(--accent);
}

.friend-info {
  flex: 1;
}

.friend-name {
  font-family: var(--font-body);
  font-size: 15px;
  font-weight: 600;
  color: var(--fg-primary);
  margin-bottom: 2px;
}

.friend-email {
  font-family: var(--font-body);
  font-size: 13px;
  color: var(--fg-secondary);
}

.friend-checkbox {
  width: 24px;
  height: 24px;
  border: 2px solid var(--border-default);
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-surface);
  transition: all var(--dur-base) var(--ease-std);
}

.friend-item.selected .friend-checkbox {
  background: var(--accent);
  border-color: var(--accent);
  color: var(--milk);
}

.modal-footer {
  display: flex;
  gap: 12px;
  padding: 20px 28px;
  border-top: 1px solid var(--border-hairline);
  justify-content: flex-end;
}

@media (max-width: 820px) {
  .multi-grid { grid-template-columns: 1fr; }
  .page-title-display { font-size: 30px; letter-spacing: -0.8px; }

  .ch-row {
    grid-template-columns: 1fr 1fr;
    gap: 12px;
  }

  .ch-letters,
  .ch-scores,
  .ch-status {
    grid-column: 1 / -1;
  }

  .modal-content {
    max-height: 90vh;
  }
}

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 6px;
  margin-top: 16px;
  flex-wrap: wrap;
}

.pagination .btn {
  min-width: 32px;
  padding: 4px 8px;
  font-family: var(--font-mono);
}

.pagination .btn--accent {
  box-shadow: 0 2px 0 var(--border-accent, rgba(0,0,0,0.15));
}

.pagination .btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.pagination .btn--soft:hover {
  background: var(--bg-hover);
}
</style>
