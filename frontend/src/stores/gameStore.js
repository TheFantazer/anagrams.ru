import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useGameStore = defineStore('game', () => {
  // Letter pools
  const RU_LETTERS = "АААААББВВВГГДДДЕЕЕЕЕЕЁЖЗИИИИЙКККЛЛЛМММНННННОООООППРРРСССТТТТУУУФХЦЧШЩЪЫЬЭЮЯ"
  const EN_LETTERS = "AAABCDDEEEEEFGHIIIJKLMNNOOOPQRRSSTTUUVWXYZ"

  // Game state
  const gameLetters = ref([])
  const inputWord = ref('')
  const foundWords = ref([])
  const score = ref(0)
  const timeLeft = ref(60)
  const gameActive = ref(false)
  const shake = ref(false)
  const initialTime = ref(60)

  // Settings
  const settingsLang = ref('ru')
  const settingsLetters = ref(7)

  // Generate random letters
  function generateLetters(count, lang) {
    const pool = lang === 'ru' ? RU_LETTERS : EN_LETTERS
    const letters = []
    const used = new Set()

    while (letters.length < count) {
      const idx = Math.floor(Math.random() * pool.length)
      if (!used.has(idx)) {
        used.add(idx)
        letters.push(pool[idx])
      }
      if (used.size >= pool.length) break
    }

    return letters
  }

  // Start new game
  function startGame(time, letters, lang) {
    const gl = generateLetters(letters, lang)
    gameLetters.value = gl
    inputWord.value = ''
    foundWords.value = []
    score.value = 0
    timeLeft.value = time
    initialTime.value = time
    gameActive.value = true
  }

  // Add letter to input
  function addLetter(letter) {
    if (gameActive.value) {
      inputWord.value += letter
    }
  }

  // Remove last letter
  function removeLast() {
    inputWord.value = inputWord.value.slice(0, -1)
  }

  // Calculate points for word
  function calculatePoints(word) {
    const len = word.length
    if (len === 3) return 100
    if (len === 4) return 300
    if (len === 5) return 500
    if (len >= 6) return 1000 + (len - 6) * 500
    return 0
  }

  // Submit word
  function submitWord() {
    if (inputWord.value.length < 3) {
      triggerShake()
      return
    }

    const upper = inputWord.value.toUpperCase()

    // Check if already found
    if (foundWords.value.includes(upper)) {
      triggerShake()
      return
    }

    // Check if can be made from available letters
    const available = [...gameLetters.value.map(l => l.toUpperCase())]
    for (const ch of upper) {
      const idx = available.indexOf(ch)
      if (idx === -1) {
        triggerShake()
        return
      }
      available.splice(idx, 1)
    }

    // Valid word!
    foundWords.value.push(upper)
    const pts = calculatePoints(upper)
    score.value += pts
    inputWord.value = ''
  }

  // Trigger shake animation
  function triggerShake() {
    shake.value = true
    setTimeout(() => {
      shake.value = false
    }, 400)
  }

  // Decrease timer
  function decreaseTime() {
    if (gameActive.value && timeLeft.value > 0) {
      timeLeft.value--
    }
    if (timeLeft.value === 0) {
      gameActive.value = false
    }
  }

  // End game
  function endGame() {
    gameActive.value = false
  }

  // Reset game
  function resetGame() {
    gameLetters.value = []
    inputWord.value = ''
    foundWords.value = []
    score.value = 0
    timeLeft.value = 60
    gameActive.value = false
    shake.value = false
  }

  // Computed: available letters (remaining after input)
  const availableLetters = computed(() => {
    const remaining = {}
    gameLetters.value.forEach(l => {
      const upper = l.toUpperCase()
      remaining[upper] = (remaining[upper] || 0) + 1
    })

    for (const ch of inputWord.value.toUpperCase()) {
      if (remaining[ch]) {
        remaining[ch]--
      }
    }

    return remaining
  })

  // Computed: timer percentage
  const timerPercentage = computed(() => {
    return timeLeft.value / initialTime.value
  })

  return {
    // State
    gameLetters,
    inputWord,
    foundWords,
    score,
    timeLeft,
    gameActive,
    shake,
    settingsLang,
    settingsLetters,

    // Computed
    availableLetters,
    timerPercentage,

    // Actions
    startGame,
    addLetter,
    removeLast,
    submitWord,
    triggerShake,
    decreaseTime,
    endGame,
    resetGame,
  }
})
