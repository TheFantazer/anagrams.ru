import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useGameStore = defineStore('game', () => {
  
  const RU_LETTERS = "АААААББВВВГГДДДЕЕЕЕЕЕЁЖЗИИИИЙКККЛЛЛМММНННННОООООППРРРСССТТТТУУУФХЦЧШЩЪЫЬЭЮЯ"
  const EN_LETTERS = "AAABCDDEEEEEFGHIIIJKLMNNOOOPQRRSSTTUUVWXYZ"

  const dictionaries = ref({
    ru: null,
    en: null
  })

  const gameLetters = ref([])
  const inputWord = ref('')
  const foundWords = ref([])
  const score = ref(0)
  const timeLeft = ref(60)
  const gameActive = ref(false)
  const shake = ref(false)
  const initialTime = ref(60)
  const validWords = ref([])
  const sessionId = ref(null)
  const usedLetterIndices = ref([])

  const settingsLang = ref('ru')
  const settingsLetters = ref(7)

  async function loadDictionary(lang) {
    if (dictionaries.value[lang]) {
      return dictionaries.value[lang]
    }

    try {
      const response = await fetch(`/dictionaries/${lang}.txt`)
      const text = await response.text()
      const words = text.split('\n').filter(w => w.trim().length > 0)

      const byLength = {}
      words.forEach(word => {
        const len = word.length
        if (!byLength[len]) {
          byLength[len] = []
        }
        byLength[len].push(word.toLowerCase())
      })

      dictionaries.value[lang] = byLength
      return byLength
    } catch (error) {
      console.error('Failed to load dictionary:', error)
      return null
    }
  }

  async function generateLettersFromDict(count, lang) {
    const dict = await loadDictionary(lang)

    if (!dict || !dict[count] || dict[count].length === 0) {
      return generateLettersFallback(count, lang)
    }

    const words = dict[count]
    const randomWord = words[Math.floor(Math.random() * words.length)]

    const letters = randomWord.toUpperCase().split('')
    for (let i = letters.length - 1; i > 0; i--) {
      const j = Math.floor(Math.random() * (i + 1))
      ;[letters[i], letters[j]] = [letters[j], letters[i]]
    }

    return letters
  }

  function generateLettersFallback(count, lang) {
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

  async function startGame(time, letters, lang) {
    try {
      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
      const response = await fetch(`${apiUrl}/api/v1/sessions`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          language: lang,
          letter_count: letters,
          time_limit: time
        })
      })

      if (!response.ok) {
        throw new Error('Failed to create session')
      }

      const session = await response.json()

      sessionId.value = session.id
      gameLetters.value = session.letters.toUpperCase().split('')
      validWords.value = session.valid_words.map(w => w.toUpperCase())
      inputWord.value = ''
      foundWords.value = []
      score.value = 0
      timeLeft.value = time
      initialTime.value = time
      gameActive.value = true
      usedLetterIndices.value = []
    } catch (error) {
      console.error('Failed to start game:', error)
      const gl = await generateLettersFromDict(letters, lang)
      gameLetters.value = gl
      validWords.value = []
      inputWord.value = ''
      foundWords.value = []
      score.value = 0
      timeLeft.value = time
      initialTime.value = time
      gameActive.value = true
      usedLetterIndices.value = []
    }
  }

  function addLetterByIndex(index) {
    if (!gameActive.value) return
    if (usedLetterIndices.value.includes(index)) return

    const letter = gameLetters.value[index]
    usedLetterIndices.value.push(index)
    inputWord.value += letter.toUpperCase()
  }

  function addLetter(letter) {
    if (!gameActive.value) return

    const upper = letter.toUpperCase()

    if (!availableLetters.value[upper] || availableLetters.value[upper] <= 0) {
      return
    }

    const availableIndex = gameLetters.value.findIndex((l, idx) =>
      l.toUpperCase() === upper && !usedLetterIndices.value.includes(idx)
    )

    if (availableIndex === -1) return

    usedLetterIndices.value.push(availableIndex)
    inputWord.value += upper
  }

  function removeLast() {
    if (inputWord.value.length > 0) {
      inputWord.value = inputWord.value.slice(0, -1)
      usedLetterIndices.value.pop()
    }
  }

  function calculatePoints(word) {
    const len = word.length
    if (len === 3) return 100
    if (len === 4) return 400
    if (len === 5) return 1200
    if (len === 6) return 2000
    if (len === 7) return 2800
    return 0
  }

  function submitWord() {
    if (inputWord.value.length < 3) {
      triggerShake()
      return
    }

    const upper = inputWord.value.toUpperCase()

    if (foundWords.value.includes(upper)) {
      triggerShake()
      return
    }

    const available = [...gameLetters.value.map(l => l.toUpperCase())]
    for (const ch of upper) {
      const idx = available.indexOf(ch)
      if (idx === -1) {
        triggerShake()
        return
      }
      available.splice(idx, 1)
    }

    if (validWords.value.length > 0 && !validWords.value.includes(upper)) {
      triggerShake()
      return
    }

    foundWords.value.push(upper)
    const pts = calculatePoints(upper)
    score.value += pts
    inputWord.value = ''
    usedLetterIndices.value = []
  }

  
  function triggerShake() {
    shake.value = true
    setTimeout(() => {
      shake.value = false
    }, 400)
  }


  function decreaseTime() {
    if (gameActive.value && timeLeft.value > 0) {
      timeLeft.value--
    }
    if (timeLeft.value === 0) {
      endGame()
    }
  }

  async function endGame() {
    if (!gameActive.value) return
    gameActive.value = false

    
    if (sessionId.value) {
      try {
        const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
        const durationMs = (initialTime.value - timeLeft.value) * 1000

        let userId = null
        const storedUser = localStorage.getItem('anagram_user')
        if (storedUser) {
          try {
            const userData = JSON.parse(storedUser)
            userId = userData.id
          } catch (e) {
            console.error('Failed to parse user data', e)
          }
        }

        await fetch(`${apiUrl}/api/v1/sessions/${sessionId.value}/results`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            user_id: userId,
            player_name: 'Player',
            fingerprint: generateFingerprint(),
            found_words: foundWords.value,
            duration_ms: durationMs
          })
        })
      } catch (error) {
        console.error('Failed to submit results:', error)
      }
    }
  }

  function generateFingerprint() {
    
    const nav = navigator
    const screen = window.screen
    const components = [
      nav.userAgent,
      nav.language,
      screen.colorDepth,
      screen.width + 'x' + screen.height,
      new Date().getTimezoneOffset()
    ]
    return btoa(components.join('|'))
  }

  function resetGame() {
    gameLetters.value = []
    inputWord.value = ''
    foundWords.value = []
    score.value = 0
    timeLeft.value = 60
    gameActive.value = false
    shake.value = false
    validWords.value = []
    sessionId.value = null
    usedLetterIndices.value = []
  }

  
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

  
  const timerPercentage = computed(() => {
    return timeLeft.value / initialTime.value
  })

  return {

    gameLetters,
    inputWord,
    foundWords,
    score,
    timeLeft,
    gameActive,
    shake,
    settingsLang,
    settingsLetters,
    validWords,
    sessionId,
    usedLetterIndices,

    availableLetters,
    timerPercentage,

    startGame,
    addLetter,
    addLetterByIndex,
    removeLast,
    submitWord,
    triggerShake,
    decreaseTime,
    endGame,
    resetGame,
  }
})
