<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useGameStore } from '../stores/gameStore'
import TimerRing from '../components/TimerRing.vue'

const router = useRouter()
const gameStore = useGameStore()

const wrapRef = ref(null)
const showAllWords = ref(false)
const shake = ref(false)
const winFx = ref(false)
const errHint = ref('')

let timerInterval = null

watch(() => gameStore.gameActive, async (isActive) => {
  if (isActive) {
    await nextTick()
    wrapRef.value?.focus()
    startTimer()
  } else {
    stopTimer()
  }
})

onMounted(async () => {
  if (gameStore.gameActive) {
    await nextTick()
    wrapRef.value?.focus()
    startTimer()
  }
})

onUnmounted(() => {
  stopTimer()
})

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
    errHint.value = 'too short — 3+ letters'
    setTimeout(() => { shake.value = false }, 420)
    return
  }

  if (gameStore.foundWords.includes(word)) {
    shake.value = true
    errHint.value = 'already found'
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
    errHint.value = 'not in dictionary'
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
  await gameStore.startGame(
    gameStore.lastGameTime,
    gameStore.lastGameLetters,
    gameStore.lastGameLang
  )
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
              <div class="hud-stat-label">Score</div>
              <div class="hud-stat-value accent">{{ gameStore.score.toLocaleString() }}</div>
            </div>
          </div>
          <div class="hud-stat">
            <div>
              <div class="hud-stat-label">Found</div>
              <div class="hud-stat-value">
                {{ gameStore.foundWords.length }}<span style="color:var(--fg-faint);font-size:13px">/{{ gameStore.validWords?.length || 0 }}</span>
              </div>
            </div>
          </div>
        </div>

        <TimerRing />

        <div class="hud-right">
          <button class="btn btn--soft btn--sm" @click="endGame">End game</button>
          <button class="btn btn--ghost btn--sm" @click="exitGame">Exit</button>
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
          {{ gameStore.inputWord || '•••' }}
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
            Delete <span class="kbd">⌫</span>
          </button>
          <button class="btn btn--ghost" @click="clearAll">
            Clear <span class="kbd">Esc</span>
          </button>
          <button class="btn btn--accent" @click="submitWord">
            Submit <span class="kbd" style="background:rgba(255,255,255,0.18);color:var(--milk);border-color:transparent">↵</span>
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
          <span class="title">Words you've found</span>
          <span class="count">{{ gameStore.foundWords.length }} / {{ gameStore.validWords?.length || 0 }}</span>
        </div>
        <p v-if="gameStore.foundWords.length === 0" class="muted found-rail-empty">
          Submit your first word to start the streak. Longer words score more.
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
      <div class="over-eyebrow">Time's up</div>
      <h1 class="over-title">{{ gameStore.score > 0 ? 'Nice game.' : 'No words this round.' }}</h1>
      <div class="over-score">{{ gameStore.score.toLocaleString() }}</div>

      <div class="over-meta">
        <div class="cell">
          <div class="cell-num">{{ gameStore.foundWords.length }}/{{ gameStore.validWords?.length || 0 }}</div>
          <div class="cell-lbl">words</div>
        </div>
        <div class="cell">
          <div class="cell-num">{{ percentFound }}%</div>
          <div class="cell-lbl">found</div>
        </div>
        <div class="cell">
          <div class="cell-num">{{ longestWord.length || 0 }}</div>
          <div class="cell-lbl">longest</div>
        </div>
        <div class="cell">
          <div class="cell-num mono" style="font-size:14px">{{ gameStore.gameLetters?.join('').toLowerCase() || '' }}</div>
          <div class="cell-lbl">set</div>
        </div>
      </div>

      <div class="row gap-2" style="justify-content:center;flex-wrap:wrap;margin-bottom:12px">
        <button class="btn btn--accent btn--lg" @click="handlePlayAgain">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <path d="M13 2L4 14h7l-1 8 9-12h-7l1-8z"/>
          </svg>
          Play again
        </button>
        <button class="btn btn--primary btn--lg" @click="router.push('/multiplayer')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="18" cy="5" r="3"/>
            <circle cx="6" cy="12" r="3"/>
            <circle cx="18" cy="19" r="3"/>
            <path d="M8.59 13.51l6.83 3.98M15.41 6.51l-6.82 3.98"/>
          </svg>
          Challenge a friend
        </button>
      </div>

      <div class="found-rail-head" style="margin-top:32px;border-top:1px solid var(--border-hairline);padding-top:24px">
        <span class="title">All words — {{ gameStore.validWords?.length || 0 }}</span>
        <button v-if="!showAllWords" class="btn btn--sm btn--ghost" @click="showAllWords = true">Show missed</button>
      </div>

      <div class="result-grid">
        <span
          v-for="(word, i) in sortedWords"
          :key="i"
          :class="['word-chip', {
            'found': gameStore.foundWords.includes(word),
            'revealed': !gameStore.foundWords.includes(word) && showAllWords
          }]"
        >
          {{ gameStore.foundWords.includes(word) || showAllWords ? word.toLowerCase() : '•'.repeat(word.length) }}
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
