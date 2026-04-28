<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../stores/userStore'
import PlayerBlock from '../components/PlayerBlock.vue'
import WordChips from '../components/WordChips.vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loading = ref(true)
const error = ref(null)
const session = ref(null)
const results = ref([])
const tab = ref('you') // mobile tab toggle
const revealed = ref(false)

const sessionId = computed(() => route.params.sessionId)

const myResult = computed(() => {
  if (!results.value || !userStore.userId) return null
  return results.value.find(r => r.user_id === userStore.userId)
})

const opponentResult = computed(() => {
  if (!results.value || !userStore.userId) return null
  return results.value.find(r => r.user_id !== userStore.userId)
})

const bothPlayed = computed(() => !!myResult.value && !!opponentResult.value)
const onlyYou = computed(() => !!myResult.value && !opponentResult.value)
const onlyThey = computed(() => !myResult.value && !!opponentResult.value)

const winner = computed(() => {
  if (!bothPlayed.value) return null
  if (myResult.value.score > opponentResult.value.score) return 'you'
  if (myResult.value.score < opponentResult.value.score) return 'them'
  return 'tie'
})

const verdictLabel = computed(() => {
  if (!bothPlayed.value) return t('results.verdict.pending')
  if (winner.value === 'you') return t('results.verdict.won')
  if (winner.value === 'them') return t('results.verdict.lost')
  return t('results.verdict.tie')
})

const verdictClass = computed(() => {
  if (!bothPlayed.value) return 'pending'
  if (winner.value === 'you') return 'won'
  if (winner.value === 'them') return 'lost'
  return 'tie'
})

const opponentName = computed(() => {
  return opponentResult.value?.player_name || session.value?.creator_username || 'Opponent'
})

const allWords = computed(() => {
  return session.value?.valid_words?.map(w => w.toUpperCase()) || []
})

const myWords = computed(() => {
  return myResult.value?.found_words?.map(w => w.toUpperCase()) || []
})

const opponentWords = computed(() => {
  return opponentResult.value?.found_words?.map(w => w.toUpperCase()) || []
})

const mySet = computed(() => new Set(myWords.value))
const opponentSet = computed(() => new Set(opponentWords.value))

const sortedAll = computed(() => {
  return [...allWords.value].sort((a, b) => {
    if (a.length !== b.length) return b.length - a.length
    return a.localeCompare(b)
  })
})

// Group words by length
function groupByLength(words) {
  const groups = {}
  words.forEach(w => {
    (groups[w.length] = groups[w.length] || []).push(w)
  })
  const lens = Object.keys(groups).map(Number).sort((a, b) => b - a)
  return lens.map(l => ({ length: l, words: groups[l] }))
}

const myGroups = computed(() => groupByLength(myWords.value))
const opponentGroups = computed(() => groupByLength(opponentWords.value))

async function loadResults() {
  loading.value = true
  error.value = null

  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'

    // Load session data
    const sessionResponse = await fetch(`${apiUrl}/api/v1/sessions/${sessionId.value}`)
    if (!sessionResponse.ok) throw new Error('Failed to load session')
    session.value = await sessionResponse.json()

    // Load results
    const resultsResponse = await fetch(`${apiUrl}/api/v1/sessions/${sessionId.value}/results`)
    if (!resultsResponse.ok) throw new Error('Failed to load results')
    results.value = await resultsResponse.json()

  } catch (err) {
    console.error('Failed to load results:', err)
    error.value = 'Failed to load results'
  } finally {
    loading.value = false
  }
}

function goBack() {
  router.push('/play')
}

function playChallenge() {
  router.push(`/challenge/${sessionId.value}`)
}

function share() {
  const letters = session.value?.letters?.toLowerCase() || ''
  const text = `anagrams.ru/r/${letters}-${sessionId.value}`
  navigator.clipboard?.writeText(text).then(() => {
    userStore.showToast(t('results.resultCopied'), 'success')
  })
}

onMounted(() => {
  loadResults()
})
</script>

<template>
  <div class="page">
    <div class="shell res-wrap">
      <!-- Loading State -->
      <div v-if="loading" class="loading-state">
        <div class="spinner"></div>
        <p class="muted">{{ $t('common.loading') }}</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="error-state">
        <p class="muted">{{ error }}</p>
        <button class="btn btn--primary" @click="goBack">
          {{ $t('challenge.backToHome') }}
        </button>
      </div>

      <!-- Results -->
      <div v-else-if="session">
        <!-- Top bar -->
        <div class="res-topbar">
          <button class="btn btn--ghost btn--sm" @click="goBack">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M19 12H5M12 19l-7-7 7-7"/>
            </svg>
            {{ $t('results.backToMultiplayer') }}
          </button>
          <span :class="['res-verdict', `v-${verdictClass}`]">
            <span v-if="!bothPlayed" class="res-pulse" />
            {{ verdictLabel }}
          </span>
        </div>

        <!-- Header — set + meta -->
        <header class="res-head">
          <div class="res-head-left">
            <div class="page-eyebrow">{{ $t('results.challengeBy', { name: opponentName }) }}</div>
            <div class="res-set">
              <span v-for="(L, i) in session.letters.split('')" :key="i" class="res-tile">{{ L.toUpperCase() }}</span>
            </div>
            <div class="res-meta">
              <span>
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="12" r="10"/>
                  <path d="M12 6v6l4 2"/>
                </svg>
                {{ $t('results.finished') }}
              </span>
              <span class="dot-sep">·</span>
              <span>{{ $t('results.possibleWords', { count: allWords.length }) }}</span>
            </div>
          </div>
          <div class="res-head-right">
            <button class="btn btn--ghost btn--sm" @click="share">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="18" cy="5" r="3"/>
                <circle cx="6" cy="12" r="3"/>
                <circle cx="18" cy="19" r="3"/>
                <path d="M8.59 13.51l6.83 3.98M15.41 6.51l-6.82 3.98"/>
              </svg>
              {{ $t('results.share') }}
            </button>
            <button v-if="onlyThey || !myResult" class="btn btn--accent btn--sm" @click="playChallenge">
              <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M6 4l14 8-14 8z"/>
              </svg>
              {{ $t('results.playIt') }}
            </button>
          </div>
        </header>

        <!-- SCORE COMPARISON (both played) -->
        <section v-if="bothPlayed" class="res-versus">
          <PlayerBlock
            who="you"
            :name="$t('results.tabs.you')"
            :score="myResult.score"
            :found="myWords.length"
            :total="allWords.length"
            :winner="winner === 'you'"
          />
          <div class="res-vs">
            <div class="res-vs-bar">
              <div class="res-vs-bar-fill you" :style="{ flex: myResult.score }" />
              <div class="res-vs-bar-fill them" :style="{ flex: opponentResult.score }" />
            </div>
            <div class="res-vs-label">{{ $t('results.vs') }}</div>
            <div class="res-vs-delta mono">
              <template v-if="winner === 'you'">+{{ (myResult.score - opponentResult.score).toLocaleString() }}</template>
              <template v-else-if="winner === 'them'">−{{ (opponentResult.score - myResult.score).toLocaleString() }}</template>
              <template v-else>EVEN</template>
            </div>
          </div>
          <PlayerBlock
            who="them"
            :name="opponentName"
            :score="opponentResult.score"
            :found="opponentWords.length"
            :total="allWords.length"
            :winner="winner === 'them'"
          />
        </section>

        <!-- SOLO STATE (only you played) -->
        <section v-if="onlyYou" class="res-solo">
          <div class="res-solo-side res-solo-side--you">
            <PlayerBlock
              who="you"
              :name="$t('results.tabs.you')"
              :score="myResult.score"
              :found="myWords.length"
              :total="allWords.length"
              :winner="false"
              solo
            />
          </div>
          <div class="res-solo-side res-solo-side--wait">
            <div class="res-wait">
              <div class="res-wait-anim">
                <span /><span /><span />
              </div>
              <h3 class="res-wait-title">{{ $t('results.waiting.title', { name: opponentName }) }}</h3>
              <p class="muted" style="font-size:13px; margin:4px 0 16px; max-width:280px; text-align:center">
                {{ $t('results.waiting.description', { name: opponentName }) }}
              </p>
              <button class="btn btn--soft btn--sm" @click="share">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="9" y="9" width="13" height="13" rx="2"/>
                  <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                </svg>
                {{ $t('results.waiting.nudge') }}
              </button>
            </div>
          </div>
        </section>

        <!-- REDIRECT STATE (only they played) -->
        <section v-if="onlyThey" class="res-redirect">
          <div class="res-redirect-card">
            <h3 class="res-redirect-title">{{ $t('results.playItFirst.title') }}</h3>
            <p class="muted" style="font-size:13px; margin:4px 0 16px; max-width:320px; text-align:center">
              {{ $t('results.playItFirst.description', { name: opponentName }) }}
            </p>
            <button class="btn btn--accent btn--lg" @click="playChallenge">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M6 4l14 8-14 8z"/>
              </svg>
              {{ $t('results.playItFirst.button') }}
            </button>
          </div>
        </section>

        <!-- WORDS (hide when only opponent has played) -->
        <section v-if="!onlyThey && (myResult || opponentResult)" class="res-words">
          <div class="res-words-head">
            <div>
              <div class="page-eyebrow">{{ $t('results.words.eyebrow') }}</div>
              <h3 class="res-words-title">
                {{ bothPlayed ? $t('results.words.bothPlayed') : myResult ? $t('results.words.youPlayed') : $t('results.words.theyPlayed', { name: opponentName }) }}
              </h3>
            </div>
            <button v-if="bothPlayed" class="btn btn--soft btn--sm" @click="revealed = !revealed">
              {{ revealed ? $t('results.words.hideMissed') : $t('results.words.showMissed') }}
            </button>
          </div>

          <!-- Mobile tab toggle (only if both played) -->
          <div v-if="bothPlayed" class="res-mob-tabs">
            <button class="chip-toggle" :data-active="tab === 'you'" @click="tab = 'you'">
              {{ $t('results.tabs.you') }} · {{ myWords.length }}
            </button>
            <button class="chip-toggle" :data-active="tab === 'them'" @click="tab = 'them'">
              {{ opponentName }} · {{ opponentWords.length }}
            </button>
            <button class="chip-toggle" :data-active="tab === 'all'" @click="tab = 'all'">
              {{ $t('results.tabs.all') }} · {{ allWords.length }}
            </button>
          </div>

          <div :class="['res-cols', bothPlayed ? '' : 'res-cols--single']" :data-mob-tab="tab">
            <!-- YOUR column -->
            <div v-if="myResult" class="res-col res-col--you">
              <div class="res-col-head">
                <span class="res-col-name">{{ $t('results.tabs.you') }}</span>
                <span class="res-col-num mono">{{ myWords.length }}</span>
              </div>
              <WordChips :words="myWords" owner="you" />
            </div>

            <!-- THEIR column -->
            <div v-if="opponentResult" class="res-col res-col--them">
              <div class="res-col-head">
                <span class="res-col-name">{{ opponentName }}</span>
                <span class="res-col-num mono">{{ opponentWords.length }}</span>
              </div>
              <WordChips :words="opponentWords" owner="them" />
            </div>
          </div>

          <!-- ALL words (mobile-only via tab, desktop via toggle) -->
          <div v-if="revealed || tab === 'all'" :class="['res-all', tab === 'all' ? 'mob-only' : '']">
            <div class="res-col-head">
              <span class="res-col-name">{{ $t('results.words.allWords') }}</span>
              <span class="res-col-num mono">{{ sortedAll.length }}</span>
            </div>
            <div class="result-grid">
              <span
                v-for="(w, i) in sortedAll"
                :key="i"
                :class="[
                  'word-chip',
                  mySet.has(w) && opponentSet.has(w) ? 'both' :
                  mySet.has(w) ? 'found' :
                  opponentSet.has(w) ? 'them-only' : 'revealed'
                ]"
                :title="mySet.has(w) && opponentSet.has(w) ? 'both found' : mySet.has(w) ? 'you' : opponentSet.has(w) ? opponentName : 'missed'"
              >
                {{ w.toLowerCase() }}
              </span>
            </div>
            <div class="res-legend">
              <span><span class="lg-sw lg-both" /> {{ $t('results.legend.both') }}</span>
              <span><span class="lg-sw lg-you" /> {{ $t('results.legend.you') }}</span>
              <span><span class="lg-sw lg-them" /> {{ $t('results.legend.them', { name: opponentName }) }}</span>
              <span><span class="lg-sw lg-miss" /> {{ $t('results.legend.missed') }}</span>
            </div>
          </div>
        </section>

        <!-- Footer actions -->
        <div class="res-foot">
          <button class="btn btn--ghost" @click="goBack">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M19 12H5M12 19l-7-7 7-7"/>
            </svg>
            {{ $t('results.backToChallenges') }}
          </button>
          <div class="row gap-2">
            <button v-if="bothPlayed" class="btn btn--soft" @click="playChallenge">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M13 2L4 14h7l-1 8 9-12h-7l1-8z"/>
              </svg>
              {{ $t('results.rematch') }}
            </button>
            <button class="btn btn--primary" @click="share">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="18" cy="5" r="3"/>
                <circle cx="6" cy="12" r="3"/>
                <circle cx="18" cy="19" r="3"/>
                <path d="M8.59 13.51l6.83 3.98M15.41 6.51l-6.82 3.98"/>
              </svg>
              {{ $t('results.shareResult') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Imports from results.css */
.res-wrap {
  max-width: 1040px;
  margin: 0 auto;
  padding: 0 32px;
  padding-bottom: 80px;
}

@media (max-width: 720px) {
  .res-wrap {
    padding: 0 20px;
    padding-bottom: 80px;
  }
}

.loading-state,
.error-state {
  text-align: center;
  padding: 60px 20px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--border-hairline);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.res-topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.res-verdict {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 14px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1.5px;
}

.res-verdict.v-won {
  background: var(--success-soft);
  color: var(--success);
}

.res-verdict.v-lost {
  background: var(--danger-soft);
  color: var(--danger);
}

.res-verdict.v-tie {
  background: var(--bg-card);
  color: var(--fg-secondary);
}

.res-verdict.v-pending {
  background: var(--warning);
  color: var(--navy);
}

.res-pulse {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: currentColor;
  animation: blink 1.4s ease-in-out infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}

/* Header */
.res-head {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 20px;
  padding-bottom: 24px;
  margin-bottom: 24px;
  border-bottom: 1px solid var(--border-hairline);
}

.res-head-left .page-eyebrow {
  margin-bottom: 14px;
}

.res-set {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 14px;
}

.res-tile {
  width: 44px;
  height: 50px;
  border-radius: 11px;
  background: var(--navy);
  color: var(--milk);
  display: grid;
  place-items: center;
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 20px;
  box-shadow: 0 3px 0 var(--navy-2);
}

.res-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  font-size: 12px;
  color: var(--fg-muted);
}

.res-meta svg {
  vertical-align: -2px;
  margin-right: 4px;
}

.dot-sep {
  color: var(--fg-faint);
}

.res-head-right {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

/* Versus */
.res-versus {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 20px;
  align-items: stretch;
  margin-bottom: 32px;
}

.res-pb {
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 20px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  position: relative;
  transition: all var(--dur-base);
}

.res-pb.is-winner {
  background: var(--navy);
  color: var(--milk);
  border-color: var(--navy);
}

.res-pb.is-winner .res-pb-name {
  color: var(--milk);
}

.res-pb.is-winner .res-pb-score {
  color: var(--milk);
}

.res-pb.is-winner .res-pb-meta .mono {
  color: var(--milk);
}

.res-pb.is-winner .res-pb-meta .lbl {
  color: color-mix(in oklab, var(--milk) 60%, transparent);
}

.res-pb.is-winner .sep {
  background: rgba(255, 255, 255, 0.18);
}

.res-pb-top {
  display: flex;
  align-items: center;
  gap: 12px;
}

.res-pb-avatar {
  width: 40px;
  height: 40px;
  border-radius: 999px;
  background: var(--grad-accent);
  color: var(--milk);
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 16px;
  display: grid;
  place-items: center;
  border: 2px solid var(--milk);
  box-shadow: 0 0 0 1px var(--border-default);
}

.res-pb.is-winner .res-pb-avatar {
  border-color: var(--navy);
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.18);
}

.res-pb-name {
  font-family: var(--font-display);
  font-weight: 600;
  font-size: 18px;
  color: var(--fg-primary);
  flex: 1;
}

.res-pb-crown {
  font-size: 22px;
}

.res-pb-score {
  font-family: var(--font-mono);
  font-size: 48px;
  font-weight: 700;
  letter-spacing: -1.5px;
  line-height: 1;
  color: var(--accent);
  margin: 4px 0;
}

.res-pb.is-winner .res-pb-score {
  color: var(--milk);
}

.res-pb-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  padding-top: 12px;
  border-top: 1px solid var(--border-hairline);
}

.res-pb.is-winner .res-pb-meta {
  border-top-color: rgba(255, 255, 255, 0.12);
}

.res-pb-meta > div {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.res-pb-meta .mono {
  font-family: var(--font-mono);
  font-size: 18px;
  font-weight: 700;
  color: var(--fg-primary);
}

.res-pb-meta .lbl {
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 1.2px;
  color: var(--fg-muted);
  font-weight: 600;
}

.res-pb-meta .sep {
  width: 1px;
  align-self: stretch;
  background: var(--border-hairline);
}

.res-vs {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  min-width: 80px;
}

.res-vs-bar {
  width: 8px;
  flex: 1;
  display: flex;
  flex-direction: column;
  border-radius: 999px;
  overflow: hidden;
  background: var(--bg-card);
  min-height: 120px;
}

.res-vs-bar-fill.you {
  background: var(--accent);
}

.res-vs-bar-fill.them {
  background: var(--navy);
}

.res-vs-label {
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 14px;
  color: var(--fg-muted);
  text-transform: uppercase;
  letter-spacing: 2px;
}

.res-vs-delta {
  font-family: var(--font-mono);
  font-weight: 700;
  font-size: 14px;
  color: var(--accent);
  padding: 4px 10px;
  background: var(--accent-soft);
  border-radius: 999px;
}

/* Solo states */
.res-solo {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
  margin-bottom: 32px;
}

.res-solo-side--wait {
  background: var(--bg-surface);
  border: 1px dashed var(--border-default);
  border-radius: 20px;
  padding: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.res-wait {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

.res-wait-anim {
  display: flex;
  gap: 6px;
  margin-bottom: 6px;
}

.res-wait-anim span {
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: var(--accent);
  animation: wait-bounce 1.4s ease-in-out infinite;
}

.res-wait-anim span:nth-child(2) {
  animation-delay: 0.2s;
}

.res-wait-anim span:nth-child(3) {
  animation-delay: 0.4s;
}

@keyframes wait-bounce {
  0%, 100% {
    opacity: 0.3;
    transform: translateY(0);
  }
  50% {
    opacity: 1;
    transform: translateY(-6px);
  }
}

.res-wait-title {
  font-family: var(--font-display);
  font-size: 22px;
  font-weight: 700;
  letter-spacing: -0.5px;
  color: var(--fg-primary);
  margin: 0;
}

/* Redirect */
.res-redirect {
  display: grid;
  place-items: center;
  padding: 32px 16px;
}

.res-redirect-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 32px 28px;
  border: 1px solid var(--border-default);
  border-radius: 18px;
  background: var(--bg-surface);
  max-width: 420px;
  text-align: center;
}

.res-redirect-title {
  font-family: var(--font-display);
  font-size: 22px;
  margin: 0 0 4px;
  color: var(--fg-primary);
}

/* Words */
.res-words {
  margin-bottom: 32px;
}

.res-words-head {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  flex-wrap: wrap;
  gap: 12px;
  margin-bottom: 18px;
}

.res-words-head .page-eyebrow {
  margin-bottom: 6px;
}

.res-words-title {
  font-family: var(--font-display);
  font-size: 24px;
  font-weight: 700;
  letter-spacing: -0.5px;
  color: var(--fg-primary);
  margin: 0;
}

.res-mob-tabs {
  display: none;
  gap: 6px;
  margin-bottom: 14px;
  flex-wrap: wrap;
}

.res-cols {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

.res-cols--single {
  grid-template-columns: 1fr;
  max-width: 600px;
  margin: 0 auto;
}

.res-col {
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 16px;
  padding: 18px;
}

.res-col--you {
  border-top: 3px solid var(--accent);
}

.res-col--them {
  border-top: 3px solid var(--navy);
}

.res-col-head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: 14px;
  padding-bottom: 10px;
  border-bottom: 1px solid var(--border-hairline);
}

.res-col-name {
  font-family: var(--font-display);
  font-weight: 600;
  font-size: 14px;
  color: var(--fg-primary);
}

.res-col-num {
  font-family: var(--font-mono);
  font-weight: 700;
  font-size: 18px;
  color: var(--accent);
}

.res-col--them .res-col-num {
  color: var(--navy);
}

.res-chip-groups {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.res-chip-group-head {
  display: flex;
  align-items: baseline;
  gap: 6px;
  margin-bottom: 6px;
}

.res-chip-group-head .mono {
  font-family: var(--font-mono);
  font-weight: 700;
  font-size: 11px;
  color: var(--navy);
  background: var(--bg-card);
  padding: 2px 8px;
  border-radius: 6px;
}

.res-chip-group-head .lbl {
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: var(--fg-muted);
  font-weight: 600;
}

.res-chip-group-head .cnt {
  margin-left: auto;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--fg-faint);
}

.res-chip-row {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.r-chip {
  padding: 4px 9px;
  border-radius: 7px;
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.r-chip--you {
  background: var(--accent-soft);
  color: var(--accent);
  border: 1px solid var(--border-accent);
}

.r-chip--them {
  background: var(--bg-card);
  color: var(--navy);
  border: 1px solid var(--border-subtle);
}

.res-empty {
  padding: 24px 0;
  text-align: center;
  font-size: 13px;
}

/* All words section */
.res-all {
  margin-top: 24px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 16px;
  padding: 18px;
}

.res-all .result-grid {
  margin: 14px 0 12px;
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.word-chip {
  padding: 4px 9px;
  border-radius: 7px;
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
  background: var(--accent-soft);
  color: var(--accent);
  border: 1px solid var(--border-accent);
}

.word-chip.both {
  background: var(--navy);
  color: var(--milk);
  border-color: var(--navy);
}

.word-chip.them-only {
  background: var(--bg-card);
  color: var(--navy);
  border-color: var(--border-subtle);
}

.word-chip.revealed {
  background: var(--bg-surface);
  color: var(--fg-muted);
  border-color: var(--border-hairline);
}

.res-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 14px;
  font-size: 11px;
  color: var(--fg-muted);
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 600;
  padding-top: 12px;
  border-top: 1px solid var(--border-hairline);
}

.lg-sw {
  display: inline-block;
  width: 12px;
  height: 12px;
  border-radius: 4px;
  vertical-align: -2px;
  margin-right: 6px;
}

.lg-both {
  background: var(--navy);
}

.lg-you {
  background: var(--accent);
}

.lg-them {
  background: var(--bg-card);
  border: 1px solid var(--border-default);
}

.lg-miss {
  background: var(--bg-surface);
  border: 1px solid var(--border-default);
}

/* Footer */
.res-foot {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 24px;
  border-top: 1px solid var(--border-hairline);
  flex-wrap: wrap;
  gap: 12px;
}

/* MOBILE */
@media (max-width: 720px) {
  .res-versus {
    grid-template-columns: 1fr;
    gap: 12px;
  }

  .res-vs {
    flex-direction: row;
    min-width: 0;
    gap: 12px;
  }

  .res-vs-bar {
    width: 100%;
    flex-direction: row;
    min-height: 0;
    height: 8px;
    flex: 1;
  }

  .res-pb-score {
    font-size: 36px;
  }

  .res-pb {
    padding: 18px;
  }

  .res-tile {
    width: 38px;
    height: 44px;
    font-size: 17px;
  }

  .res-solo {
    grid-template-columns: 1fr;
  }

  .res-mob-tabs {
    display: flex;
  }

  .res-cols {
    grid-template-columns: 1fr;
  }

  .res-cols .res-col {
    display: none;
  }

  .res-cols[data-mob-tab="you"] .res-col--you {
    display: block;
  }

  .res-cols[data-mob-tab="them"] .res-col--them {
    display: block;
  }

  .res-cols[data-mob-tab="all"] .res-col {
    display: none;
  }

  .res-all {
    display: none;
  }

  .res-cols[data-mob-tab="all"] ~ .res-all {
    display: block;
  }

  .res-all.mob-only {
    display: block !important;
  }

  .res-foot {
    flex-direction: column;
    align-items: stretch;
  }

  .res-foot .row {
    justify-content: space-between;
  }
}
</style>
