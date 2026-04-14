<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useGameStore } from '../stores/gameStore'

const router = useRouter()
const gameStore = useGameStore()

const inputRef = ref(null)
const showAllWords = ref(false)
let timerInterval = null

watch(() => gameStore.gameActive, async (isActive) => {
  if (isActive) {
    await nextTick()
    inputRef.value?.focus()
    startTimer()
  } else {
    stopTimer()
  }
})

onMounted(async () => {
  if (gameStore.gameActive) {
    await nextTick()
    inputRef.value?.focus()
    startTimer()
  }
})

onUnmounted(() => {
  stopTimer()
})

function startTimer() {
  stopTimer() // stop previous timer if it was
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
    gameStore.submitWord()
  } else if (e.key === 'Backspace') {
    e.preventDefault()
    gameStore.removeLast()
  } else if (/^[a-zA-Zа-яА-ЯёЁ]$/.test(e.key)) {
    e.preventDefault()
    const upper = e.key.toUpperCase()
    gameStore.addLetter(upper)
  }
}


function isLetterUsed(letter, index) {
  return gameStore.usedLetterIndices.includes(index)
}

function handlePlayAgain() {
  gameStore.resetGame()
  showAllWords.value = false
  router.push('/')
}

function getMaskedWord(word) {
  return '?'.repeat(word.length)
}

function getDisplayWords() {
  if (!gameStore.validWords || gameStore.validWords.length === 0) {
    return []
  }

  const words = gameStore.validWords.map(word => {
    const found = gameStore.foundWords.includes(word)
    return {
      word: word,
      found: found,
      display: found || showAllWords.value ? word.toLowerCase() : getMaskedWord(word)
    }
  })

  return words.sort((a, b) => {
    if (a.word.length !== b.word.length) {
      return b.word.length - a.word.length
    }
    return a.word.localeCompare(b.word)
  })
}
</script>

<template>
  <div
    class="game-wrap"
    tabindex="0"
    @keydown="handleKeyDown"
  >
    <div v-if="gameStore.gameActive">
      <div class="score-display">{{ gameStore.score }} pts</div>

      <!-- Timer progress bar -->
      <div class="timer-bar-container">
        <div class="timer-bar" :style="{ width: (gameStore.timerPercentage * 100) + '%' }"></div>
      </div>

      <input
        ref="inputRef"
        :value="gameStore.inputWord"
        readonly
        :class="['word-input', { shake: gameStore.shake }]"
        placeholder="..."
      />

      <div class="letters-row">
        <div
          v-for="(letter, i) in gameStore.gameLetters"
          :key="i"
          :class="['letter-tile', { used: isLetterUsed(letter, i) }]"
          @click="!isLetterUsed(letter, i) && gameStore.addLetterByIndex(i)"
        >
          {{ letter.toUpperCase() }}
        </div>
      </div>

      <div class="button-row">
        <button class="btn-secondary" @click="gameStore.removeLast">
          &larr; Delete
        </button>
        <button class="btn-primary" @click="gameStore.submitWord">
          Submit &crarr;
        </button>
      </div>
    </div>

    <div v-else class="game-over">
      <p class="game-over-title">{{ gameStore.score > 0 ? "Time's up!" : "Game Over" }}</p>
      <p class="final-score">{{ gameStore.score }}</p>
      <p class="words-count">
        {{ gameStore.foundWords.length }} / {{ gameStore.validWords.length }} words found
      </p>

      <div class="words-area">
        <span
          v-for="(item, i) in getDisplayWords()"
          :key="i"
          :class="['word-chip', { found: item.found, hidden: !item.found && !showAllWords }]"
        >
          {{ item.display }}
        </span>
      </div>

      <div class="button-row">
        <button
          v-if="!showAllWords && gameStore.sessionId"
          class="btn-secondary"
          @click="showAllWords = true"
        >
          Show all words
        </button>
        <button class="btn-primary" @click="handlePlayAgain">
          Play again
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.game-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 100px);
  gap: 24px;
  padding: 20px 0;
  outline: none;
}

.score-display {
  font-family: 'Space Mono', monospace;
  font-size: 18px;
  color: var(--accent);
  font-weight: 700;
  text-align: center;
  margin-bottom: 16px;
}

.timer-bar-container {
  width: 100%;
  max-width: 400px;
  height: 6px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 24px;
}

.timer-bar {
  height: 100%;
  background: linear-gradient(90deg, var(--accent), var(--accent-hover));
  transition: width 0.3s linear;
  border-radius: 3px;
}

.word-input {
  font-family: 'Space Mono', monospace;
  font-size: 28px;
  font-weight: 700;
  text-align: center;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(99, 230, 190, 0.2);
  border-radius: 14px;
  padding: 16px 24px;
  color: #e8e6e1;
  min-width: 260px;
  outline: none;
  letter-spacing: 4px;
  caret-color: var(--accent);
}

.word-input.shake {
  animation: shake 0.4s ease;
}

.letters-row {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  justify-content: center;
  margin: 8px 0;
}

.letter-tile {
  width: 52px;
  height: 56px;
  border-radius: 12px;
  background: rgba(99, 230, 190, 0.08);
  border: 1.5px solid rgba(99, 230, 190, 0.25);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Space Mono', monospace;
  font-size: 22px;
  font-weight: 700;
  color: var(--accent);
  cursor: pointer;
  transition: all 0.15s;
  user-select: none;
}

.letter-tile.used {
  background: rgba(255, 255, 255, 0.02);
  border-color: rgba(255, 255, 255, 0.04);
  color: #333;
  cursor: default;
}

.letter-tile:not(.used):hover {
  transform: scale(1.1);
  background: rgba(99, 230, 190, 0.15);
}

.button-row {
  display: flex;
  gap: 10px;
  margin-top: 4px;
  justify-content: flex-end;
}

.btn-primary,
.btn-secondary {
  padding: 10px 28px;
  font-size: 13px;
  border-radius: 10px;
  border: none;
  cursor: pointer;
  font-family: 'Outfit', sans-serif;
  transition: all 0.2s;
  font-weight: 500;
}

.btn-primary {
  padding: 10px 28px;
  background: linear-gradient(135deg, var(--accent), var(--accent-hover));
  color: var(--bg-dark);
  font-weight: 600;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(99, 230, 190, 0.3);
}

.btn-secondary {
  padding: 10px 20px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #e8e6e1;
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.1);
}

.words-area {
  display: flex;
  flex-direction: column;
  gap: 8px;
  align-items: center;
  max-width: 500px;
  width: 100%;
}

.word-chip {
  padding: 12px 16px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.04);
  font-size: 14px;
  font-family: 'Space Mono', monospace;
  font-weight: 500;
  color: #999;
  width: 100%;
  text-align: left;
}

.word-chip.hidden {
  background: rgba(255, 255, 255, 0.02);
  border-color: rgba(255, 255, 255, 0.04);
  color: #555;
}

.word-chip.found {
  background: rgba(99, 230, 190, 0.06);
  border-color: rgba(99, 230, 190, 0.15);
  color: var(--accent);
}

.game-over {
  text-align: center;
  padding: 40px;
}

.game-over-title {
  font-family: 'Space Mono', monospace;
  font-size: 32px;
  font-weight: 700;
  color: var(--accent);
  margin: 0 0 8px;
}

.final-score {
  font-size: 48px;
  font-family: 'Space Mono', monospace;
  font-weight: 700;
  color: var(--accent);
  margin: 8px 0;
}

.words-count {
  color: #666;
  font-size: 14px;
  margin-bottom: 32px;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  10%, 30%, 50%, 70%, 90% { transform: translateX(-8px); }
  20%, 40%, 60%, 80% { transform: translateX(8px); }
}
</style>
