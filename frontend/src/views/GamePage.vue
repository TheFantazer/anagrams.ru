<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useGameStore } from '../stores/gameStore'
import { useUserStore } from '../stores/userStore'
import TimerRing from '../components/TimerRing.vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const gameStore = useGameStore()
const userStore = useUserStore()

const wrapRef = ref(null)
const showAllWords = ref(false)
const shake = ref(false)
const winFx = ref(false)
const errHint = ref('')
const sessionId = ref(null)
const isMultiplayer = ref(false)
const sessionResults = ref([])
const loadingResults = ref(false)
const wordViewMode = ref('my') // 'my', 'opponent', 'common'

let timerInterval = null

watch(() => gameStore.gameActive, async (isActive) => {
  if (isActive) {
    await nextTick()
    wrapRef.value?.focus()
    startTimer()
  } else {
    stopTimer()
    // Загружаем результаты после завершения игры, если это мультиплеер
    if (isMultiplayer.value && sessionId.value) {
      await loadSessionResults()
    }
  }
})

onMounted(async () => {
  // Проверяем, есть ли session_id в query params
  sessionId.value = route.query.session_id
  isMultiplayer.value = !!sessionId.value

  if (sessionId.value) {
    // Загружаем сессию с сервера и запускаем игру
    await loadMultiplayerSession()
  } else if (gameStore.gameActive) {
    await nextTick()
    wrapRef.value?.focus()
    startTimer()
  }
})

onUnmounted(() => {
  stopTimer()
})

async function loadMultiplayerSession() {
  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/api/v1/sessions/${sessionId.value}`)

    if (!response.ok) {
      throw new Error('Failed to load session')
    }

    const session = await response.json()

    // Call /start endpoint before starting the game (if user is authenticated)
    if (userStore.userId) {
      const startResponse = await fetch(
        `${apiUrl}/api/v1/sessions/${sessionId.value}/start?user_id=${userStore.userId}`,
        { method: 'POST' }
      )

      if (!startResponse.ok) {
        if (startResponse.status === 403) {
          // Link already used by another player
          alert(t('multiplayer.linkAlreadyUsed'))
          router.push('/play')
          return
        }
        throw new Error('Failed to start session')
      }
    }

    // Запускаем игру с параметрами из сессии
    gameStore.startGame(
      session.letters.toUpperCase(),
      session.language,
      session.time_limit,
      sessionId.value  // передаем session_id
    )

    await nextTick()
    wrapRef.value?.focus()
  } catch (error) {
    console.error('Failed to load multiplayer session:', error)
    alert('Failed to load challenge. Please try again.')
    router.push('/')
  }
}

async function loadSessionResults() {
  loadingResults.value = true
  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/api/v1/sessions/${sessionId.value}/results`)

    if (response.ok) {
      sessionResults.value = await response.json()
    }
  } catch (error) {
    console.error('Failed to load session results:', error)
  } finally {
    loadingResults.value = false
  }
}

function startTimer() {
  stopTimer()
  timerInterval = setInterval(() => {
    gameStore.decreaseTime()
    if (gameStore.timeLeft === 0) {
      stopTimer()
    }
  }, 1000)
}

function stopTimer() {
  if (timerInterval) {
    clearInterval(timerInterval)
    timerInterval = null
  }
}

function handleKeyDown(e) {
  if (!gameStore.gameActive) return

  if (e.key === 'Enter') {
    e.preventDefault()
    submitWord()
  } else if (e.key === 'Backspace') {
    e.preventDefault()
    removeLast()
  } else if (e.key === 'Escape') {
    e.preventDefault()
    clearAll()
  } else if (/^[a-zA-Zа-яА-ЯёЁ]$/.test(e.key)) {
    e.preventDefault()
    addByKey(e.key)
  }
}

function addByKey(key) {
  const upper = key.toUpperCase()
  const index = gameStore.gameLetters.findIndex((letter, idx) =>
    letter === upper && !gameStore.usedLetterIndices.includes(idx)
  )
  if (index !== -1) {
    gameStore.addLetterByIndex(index)
  }
  errHint.value = ''
}

function removeLast() {
  gameStore.removeLast()
  errHint.value = ''
}

function clearAll() {
  while (gameStore.inputWord.length > 0) {
    gameStore.removeLast()
  }
  errHint.value = ''
}

function submitWord() {
  const word = gameStore.inputWord.toUpperCase()

  if (word.length < 3) {
    shake.value = true
    errHint.value = t('game.errors.tooShort')
    setTimeout(() => { shake.value = false }, 420)
    return
  }

  if (gameStore.foundWords.includes(word)) {
    shake.value = true
    errHint.value = t('game.errors.alreadyFound')
    setTimeout(() => { shake.value = false }, 420)
    clearAll()
    return
  }

  const result = gameStore.submitWord()
  if (result.valid) {
    winFx.value = true
    setTimeout(() => { winFx.value = false }, 620)
    errHint.value = ''
  } else {
    shake.value = true
    errHint.value = t('game.errors.notInDictionary')
    setTimeout(() => { shake.value = false }, 420)
  }
}

function isLetterUsed(index) {
  return gameStore.usedLetterIndices.includes(index)
}

function endGame() {
  gameStore.endGame()
}

function exitGame() {
  gameStore.resetGame()
  router.push('/')
}

async function handlePlayAgain() {
  showAllWords.value = false
  gameStore.resetGame()

  // Всегда создаём новую сессию с теми же параметрами
  // (для соло это обычное поведение, для мультиплеера - реванш)
  await gameStore.startGame(
    gameStore.lastGameTime,
    gameStore.lastGameLetters,
    gameStore.lastGameLang
  )
}

async function handleChallengeBack() {
  // Создаём ответную игру с теми же настройками и переходим на /multiplayer
  router.push({
    path: '/play',
    query: {
      create: 'true',
      language: gameStore.lastGameLang,
      letterCount: gameStore.lastGameLetters.toString(),
      timeLimit: gameStore.lastGameTime.toString(),
      hideLetters: gameStore.lastGameHideLetters.toString()
    }
  })
}

function copyResult() {
  const letters = gameStore.gameLetters?.join('').toLowerCase() || ''
  const text = `anagrams.ru — I scored ${gameStore.score.toLocaleString()} with ${gameStore.foundWords.length}/${gameStore.validWords?.length || 0} words. Can you beat it? → ${window.location.origin}/r/${letters}`

  navigator.clipboard?.writeText(text).then(() => {
    userStore.showToast(t('game.gameOver.resultCopied'), 'success')
  }).catch(() => {
    userStore.showToast('Failed to copy', 'error')
  })
}

const sortedWords = computed(() => {
  if (!gameStore.validWords) return []
  return [...gameStore.validWords].sort((a, b) => {
    if (a.length !== b.length) return b.length - a.length
    return a.localeCompare(b)
  })
})

const percentFound = computed(() => {
  if (!gameStore.validWords || gameStore.validWords.length === 0) return 0
  return Math.round((gameStore.foundWords.length / gameStore.validWords.length) * 100)
})

const longestWord = computed(() => {
  if (!gameStore.foundWords || gameStore.foundWords.length === 0) return ''
  return gameStore.foundWords.reduce((a, b) => a.length >= b.length ? a : b)
})

const gameOutcome = computed(() => {
  if (!isMultiplayer.value || sessionResults.value.length === 0) {
    return null
  }

  const myScore = gameStore.score
  const topScore = Math.max(...sessionResults.value.map(r => r.score))

  if (myScore > topScore) {
    return 'won'
  } else if (myScore < topScore) {
    return 'lost'
  } else {
    return 'tie'
  }
})

const opponentWords = computed(() => {
  if (!isMultiplayer.value || sessionResults.value.length === 0) {
    return []
  }

  // Находим результат соперника (первый результат который не мой)
  const myWords = gameStore.foundWords.map(w => w.toUpperCase())
  const opponentResult = sessionResults.value.find(r => {
    // Проверяем что это не мои слова
    const resultWords = (r.found_words || []).map(w => w.toUpperCase())
    return JSON.stringify(resultWords.sort()) !== JSON.stringify(myWords.sort())
  })

  return opponentResult?.found_words?.map(w => w.toUpperCase()) || []
})

const displayWords = computed(() => {
  const myWords = gameStore.foundWords.map(w => w.toUpperCase())
  const opWords = opponentWords.value
  const allValidWords = gameStore.validWords?.map(w => w.toUpperCase()) || []

  if (wordViewMode.value === 'my') {
    // Показываем только мои слова + пропущенные (если showAllWords)
    return allValidWords.map(word => ({
      word,
      found: myWords.includes(word),
      foundByOpponent: false,
      common: false
    }))
  } else if (wordViewMode.value === 'opponent') {
    // Показываем слова соперника
    return allValidWords.map(word => ({
      word,
      found: opWords.includes(word),
      foundByOpponent: opWords.includes(word),
      common: false
    }))
  } else {
    // common - показываем общие и уникальные
    return allValidWords.map(word => {
      const foundByMe = myWords.includes(word)
      const foundByOpp = opWords.includes(word)
      return {
        word,
        found: foundByMe || foundByOpp,
        foundByOpponent: foundByOpp && !foundByMe,
        common: foundByMe && foundByOpp
      }
    })
  }
})

const sortedDisplayWords = computed(() => {
  return [...displayWords.value].sort((a, b) => {
    if (a.word.length !== b.word.length) return b.word.length - a.word.length
    return a.word.localeCompare(b.word)
  })
})
</script>

<template>
  <div class="page">
    <!-- Active Game View -->
    <div
      v-if="gameStore.gameActive"
      ref="wrapRef"
      class="shell game-wrap"
      tabindex="0"
      @keydown="handleKeyDown"
      style="outline: none"
    >
      <!-- HUD: Score, Timer, Actions -->
      <div class="game-hud">
        <div class="hud-left">
          <div class="hud-stat">
            <div>
              <div class="hud-stat-label">{{ $t('game.score') }}</div>
              <div class="hud-stat-value accent">{{ gameStore.score.toLocaleString() }}</div>
            </div>
          </div>
          <div class="hud-stat">
            <div>
              <div class="hud-stat-label">{{ $t('game.found') }}</div>
              <div class="hud-stat-value">
                <template v-if="gameStore.hideLetters || gameStore.isDaily">
                  {{ gameStore.foundWords.length }}
                </template>
                <template v-else>
                  {{ gameStore.foundWords.length }}<span style="color:var(--fg-faint);font-size:13px">/{{ gameStore.validWords?.length || 0 }}</span>
                </template>
              </div>
            </div>
          </div>
        </div>

        <TimerRing />

        <div class="hud-right">
          <button class="btn btn--soft btn--sm" @click="endGame">{{ $t('game.endGame') }}</button>
          <button class="btn btn--ghost btn--sm" @click="exitGame">{{ $t('game.exit') }}</button>
        </div>
      </div>

      <!-- Input & Letters -->
      <div class="game-input-wrap">
        <div
          :class="['game-input', {
            'active': gameStore.inputWord,
            'empty': !gameStore.inputWord,
            'shake': shake,
            'win': winFx
          }]"
        >
          {{ gameStore.inputWord || $t('game.placeholder') }}
        </div>

        <div class="game-input-hint" :data-visible="!!errHint">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 6L6 18M6 6l12 12"/>
          </svg>
          {{ errHint }}
        </div>

        <!-- Letter Tiles -->
        <div class="letters-grid">
          <div
            v-for="(letter, i) in gameStore.gameLetters"
            :key="i"
            class="letter-tile"
            :data-used="isLetterUsed(i)"
            @click="!isLetterUsed(i) && gameStore.addLetterByIndex(i)"
          >
            {{ letter }}
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="letter-actions">
          <button class="btn btn--soft" @click="removeLast">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
              <path d="M21 4H8l-7 8 7 8h13a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2zM18 9l-6 6M12 9l6 6"/>
            </svg>
            {{ $t('game.delete') }} <span class="kbd">⌫</span>
          </button>
          <button class="btn btn--ghost" @click="clearAll">
            {{ $t('game.clear') }} <span class="kbd">Esc</span>
          </button>
          <button class="btn btn--accent" @click="submitWord">
            {{ $t('game.submit') }} <span class="kbd" style="background:rgba(255,255,255,0.18);color:var(--milk);border-color:transparent">↵</span>
          </button>
        </div>

        <!-- Progress Rail -->
        <div class="progress-rail">
          <div class="progress-rail-fill" :style="{ width: `${gameStore.timerPercentage * 100}%` }" />
        </div>
      </div>

      <!-- Found Words Rail -->
      <div class="found-rail">
        <div class="found-rail-head">
          <span class="title">{{ $t('game.foundWords.title') }}</span>
          <span class="count">
            <template v-if="gameStore.hideLetters">
              {{ gameStore.foundWords.length }}
            </template>
            <template v-else>
              {{ gameStore.foundWords.length }} / {{ gameStore.validWords?.length || 0 }}
            </template>
          </span>
        </div>
        <p v-if="gameStore.foundWords.length === 0" class="muted found-rail-empty">
          {{ $t('game.foundWords.empty') }}
        </p>
        <div v-else class="found-chips">
          <span v-for="(word, i) in gameStore.foundWords" :key="i" class="found-chip">
            {{ word.toLowerCase() }}
          </span>
        </div>
      </div>
    </div>

    <!-- Game Over View -->
    <div v-else class="shell over-wrap">
      <div class="over-eyebrow">{{ $t('game.gameOver.title') }}</div>
      <h1 class="over-title">{{ gameStore.score > 0 ? $t('game.gameOver.subtitle') : 'No words this round.' }}</h1>

      <!-- Victory/Defeat Banner -->
      <div v-if="gameOutcome" class="outcome-banner" :class="`outcome-${gameOutcome}`">
        <svg v-if="gameOutcome === 'won'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
        </svg>
        <svg v-else-if="gameOutcome === 'lost'" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M10 15l4-4m0 4l-4-4m13 1a9 9 0 11-18 0 9 9 0 0118 0z"/>
        </svg>
        <svg v-else width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M12 2v20M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6"/>
        </svg>
        <span>
          {{ gameOutcome === 'won' ? $t('game.gameOver.youWon') : (gameOutcome === 'lost' ? $t('game.gameOver.youLost') : $t('game.gameOver.tie')) }}
        </span>
      </div>

      <div class="over-score">{{ gameStore.score.toLocaleString() }}</div>

      <div class="over-meta">
        <div class="cell">
          <div class="cell-num">
            <template v-if="gameStore.hideLetters">
              {{ gameStore.foundWords.length }}
            </template>
            <template v-else>
              {{ gameStore.foundWords.length }}/{{ gameStore.validWords?.length || 0 }}
            </template>
          </div>
          <div class="cell-lbl">{{ $t('game.gameOver.wordsFound') }}</div>
        </div>
        <div v-if="!gameStore.hideLetters" class="cell">
          <div class="cell-num">{{ percentFound }}%</div>
          <div class="cell-lbl">{{ $t('game.gameOver.percentFound') }}</div>
        </div>
        <div class="cell">
          <div class="cell-num">{{ longestWord.length || 0 }}</div>
          <div class="cell-lbl">{{ $t('game.gameOver.longestWord') }}</div>
        </div>
        <div class="cell">
          <div class="cell-num mono" style="font-size:14px">{{ gameStore.hideLetters ? '?'.repeat(gameStore.gameLetters?.length || 0) : (gameStore.gameLetters?.join('').toLowerCase() || '') }}</div>
          <div class="cell-lbl">set</div>
        </div>
      </div>

      <div class="row gap-2" style="justify-content:center;flex-wrap:wrap;margin-bottom:12px">
        <button v-if="!gameStore.isDaily" class="btn btn--accent btn--lg" @click="handlePlayAgain">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M13 2L4 14h7l-1 8 9-12h-7l1-8z"/>
          </svg>
          {{ $t('game.gameOver.playAgain') }}
        </button>
        <button
          v-if="gameStore.lastGameWasMultiplayer"
          class="btn btn--primary btn--lg"
          @click="handleChallengeBack"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M3 12a9 9 0 0 1 9-9 9.75 9.75 0 0 1 6.74 2.74L21 8"/>
            <path d="M21 3v5h-5"/>
            <path d="M21 12a9 9 0 0 1-9 9 9.75 9.75 0 0 1-6.74-2.74L3 16"/>
            <path d="M3 21v-5h5"/>
          </svg>
          {{ $t('game.gameOver.challengeBack') }}
        </button>
        <button
          v-else-if="!gameStore.isDaily"
          class="btn btn--primary btn--lg"
          @click="router.push('/play')"
        >
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="18" cy="5" r="3"/>
            <circle cx="6" cy="12" r="3"/>
            <circle cx="18" cy="19" r="3"/>
            <path d="M8.59 13.51l6.83 3.98M15.41 6.51l-6.82 3.98"/>
          </svg>
          {{ $t('game.gameOver.challengeFriend') }}
        </button>
        <button class="btn btn--ghost btn--lg" @click="copyResult">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
          </svg>
          {{ $t('game.gameOver.copyResult') }}
        </button>
      </div>

      <div class="found-rail-head" style="margin-top:32px;border-top:1px solid var(--border-hairline);padding-top:24px">
        <span class="title">
          <template v-if="gameStore.hideLetters">
            All words
          </template>
          <template v-else>
            All words — {{ gameStore.validWords?.length || 0 }}
          </template>
        </span>
        <button v-if="!showAllWords" class="btn btn--sm btn--ghost" @click="showAllWords = true">{{ $t('game.gameOver.viewAllWords') }}</button>
      </div>

      <!-- Word View Mode Switcher (only for multiplayer) -->
      <div v-if="isMultiplayer && opponentWords.length > 0" class="word-view-switcher">
        <button
          :class="['btn', 'btn--sm', wordViewMode === 'my' ? 'btn--primary' : 'btn--ghost']"
          @click="wordViewMode = 'my'"
        >
          {{ $t('game.gameOver.myWords') }}
        </button>
        <button
          :class="['btn', 'btn--sm', wordViewMode === 'opponent' ? 'btn--primary' : 'btn--ghost']"
          @click="wordViewMode = 'opponent'"
        >
          {{ $t('game.gameOver.opponentWords') }}
        </button>
        <button
          :class="['btn', 'btn--sm', wordViewMode === 'common' ? 'btn--primary' : 'btn--ghost']"
          @click="wordViewMode = 'common'"
        >
          {{ $t('game.gameOver.allWords') }}
        </button>
      </div>

      <div class="result-grid">
        <span
          v-for="(wordObj, i) in (isMultiplayer && opponentWords.length > 0 ? sortedDisplayWords : sortedWords.map(w => ({ word: w, found: gameStore.foundWords.includes(w), foundByOpponent: false, common: false })))"
          :key="i"
          :class="['word-chip', {
            'found': wordObj.found || (typeof wordObj === 'string' && gameStore.foundWords.includes(wordObj)),
            'revealed': (typeof wordObj === 'object' ? !wordObj.found : !gameStore.foundWords.includes(wordObj)) && showAllWords,
            'opponent-only': wordObj.foundByOpponent && !wordObj.common,
            'common': wordObj.common
          }]"
        >
          {{ (wordObj.found || (typeof wordObj === 'string' && gameStore.foundWords.includes(wordObj)) || showAllWords)
              ? (typeof wordObj === 'object' ? wordObj.word : wordObj).toLowerCase()
              : '•'.repeat(typeof wordObj === 'object' ? wordObj.word.length : wordObj.length) }}
        </span>
      </div>

      <div style="margin-top:28px">
        <button class="btn btn--ghost" @click="router.push('/')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M19 12H5M12 19l-7-7 7-7"/>
          </svg>
          Back home
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* All styles are in game.css */
</style>
