<script setup>
import { ref, computed, onMounted } from 'vue'
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

// Form settings
const language = ref('ru')
const letterCount = ref(7)
const timeLimit = ref(60)
const hideLetters = ref(false)

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

async function createChallenge() {
  creating.value = true
  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

    // Build URL with optional user_id query param
    let url = `${apiUrl}/api/v1/sessions`
    if (userStore.isAuthenticated && userStore.user?.id) {
      url += `?user_id=${userStore.user.id}`
    }

    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        language: language.value,
        letter_count: letterCount.value,
        time_limit: timeLimit.value
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

async function loadActiveChallenges() {
  if (!userStore.isAuthenticated || !userStore.user?.id) {
    return
  }

  loadingChallenges.value = true
  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/api/v1/sessions/my?user_id=${userStore.user.id}`)

    if (!response.ok) {
      throw new Error('Failed to load challenges')
    }

    const sessions = await response.json()

    // Загружаем результаты для каждой сессии
    const sessionsWithResults = await Promise.all(
      sessions.map(async (session) => {
        try {
          const resultsResponse = await fetch(`${apiUrl}/api/v1/sessions/${session.id}/results?top=5`)
          if (resultsResponse.ok) {
            session.results = await resultsResponse.json()
          } else {
            session.results = []
          }
        } catch (error) {
          console.error(`Failed to load results for session ${session.id}:`, error)
          session.results = []
        }
        return session
      })
    )

    activeChallenges.value = sessionsWithResults
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
    // Автоматически создаём челлендж
    createChallenge()
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
                <path d="M6 4l14 8-14 8z"/>
              </svg>
              {{ creating ? $t('common.creating') : $t('multiplayer.createChallenge') }}
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
        <div class="multi-eye" style="margin-bottom:16px">
          <span class="multi-num">03</span>{{$t('multiplayer.card03.title')}}
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

        <div v-else class="challenges-grid">
          <div
            v-for="challenge in activeChallenges"
            :key="challenge.id"
            class="challenge-card"
            @click="router.push(`/play/${challenge.id}`)"
          >
            <div class="challenge-letters">
              <span v-for="(letter, i) in challenge.letters.split('')" :key="i" class="challenge-tile">
                {{ letter.toLowerCase() }}
              </span>
            </div>
            <div class="challenge-meta">
              <span class="challenge-badge">{{ challenge.language.toUpperCase() }}</span>
              <span class="challenge-badge">{{ challenge.letter_count }} letters</span>
              <span class="challenge-badge">{{ Math.floor(challenge.time_limit / 60) }}:{{ String(challenge.time_limit % 60).padStart(2, '0') }}</span>
            </div>

            <!-- Results -->
            <div v-if="challenge.results && challenge.results.length > 0" class="challenge-results">
              <div class="challenge-results-header">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="18" cy="5" r="3"/>
                  <circle cx="6" cy="12" r="3"/>
                  <circle cx="18" cy="19" r="3"/>
                  <path d="M8.59 13.51l6.83 3.98M15.41 6.51l-6.82 3.98"/>
                </svg>
                {{ challenge.results.length }} {{ challenge.results.length === 1 ? 'player' : 'players' }}
              </div>
              <div class="challenge-results-list">
                <div
                  v-for="(result, idx) in challenge.results.slice(0, 3)"
                  :key="result.id"
                  class="challenge-result-item"
                >
                  <span class="result-rank">#{{ idx + 1 }}</span>
                  <span class="result-name">{{ result.player_name || 'Anonymous' }}</span>
                  <span class="result-score">{{ result.score }}</span>
                </div>
              </div>
            </div>
            <div v-else class="challenge-no-results">
              No one played yet
            </div>

            <div class="challenge-date">
              {{ new Date(challenge.created_at).toLocaleDateString() }}
            </div>
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

.challenges-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.challenge-card {
  padding: 20px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.challenge-card:hover {
  border-color: var(--accent);
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.challenge-letters {
  display: flex;
  gap: 6px;
  margin-bottom: 12px;
  flex-wrap: wrap;
}

.challenge-tile {
  width: 32px;
  height: 38px;
  border-radius: 8px;
  background: var(--navy);
  color: var(--milk);
  display: grid;
  place-items: center;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 16px;
  box-shadow: 0 2px 0 var(--navy-2);
}

.challenge-meta {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 8px;
}

.challenge-badge {
  padding: 4px 10px;
  background: var(--bg-card);
  border-radius: 6px;
  font-size: 11px;
  font-weight: 600;
  color: var(--fg-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.challenge-date {
  font-size: 12px;
  color: var(--fg-muted);
  margin-top: 8px;
}

.challenge-results {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border-hairline);
}

.challenge-results-header {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  font-weight: 600;
  color: var(--fg-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin-bottom: 8px;
}

.challenge-results-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.challenge-result-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  padding: 6px 8px;
  background: var(--bg-card);
  border-radius: 6px;
}

.result-rank {
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  color: var(--accent);
  min-width: 24px;
}

.result-name {
  flex: 1;
  color: var(--fg-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.result-score {
  font-weight: 700;
  color: var(--fg-primary);
}

.challenge-no-results {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--border-hairline);
  font-size: 12px;
  color: var(--fg-muted);
  font-style: italic;
}

@media (max-width: 820px) {
  .multi-grid { grid-template-columns: 1fr; }
  .page-title-display { font-size: 30px; letter-spacing: -0.8px; }
  .challenges-grid { grid-template-columns: 1fr; }
}
</style>
