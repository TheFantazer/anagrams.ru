import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const userId = ref(null)
  const username = ref('player_one')
  const email = ref('player@example.com')
  const joinedDate = ref('April 2026')
  const isAuthenticated = ref(false)

  const gamesPlayed = ref(42)
  const bestScore = ref(3200)
  const longestWord = ref('КРОКОДИЛ')

  const showHelp = ref(false)
  const showSoloSettings = ref(false)

  const soloTime = ref(60)
  const soloLetters = ref(7)
  const soloLang = ref('ru')

  const loadUser = () => {
    const stored = localStorage.getItem('anagram_user')
    if (stored) {
      try {
        const userData = JSON.parse(stored)
        userId.value = userData.id
        username.value = userData.username
        email.value = userData.email
        joinedDate.value = new Date(userData.created_at).toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
        isAuthenticated.value = true

        if (userData.default_letter_count) {
          soloLetters.value = userData.default_letter_count
        }
        if (userData.default_language) {
          soloLang.value = userData.default_language
        }
        if (userData.default_time_limit) {
          soloTime.value = userData.default_time_limit
        }
      } catch (e) {
        console.error('Failed to load user from localStorage', e)
      }
    }
  }

  loadUser()

  const loginTab = ref('login')

  const lbPeriod = ref('week')

  const easterEgg = ref(false)
  const konamiIdx = ref(0)

  function setShowHelp(value) {
    showHelp.value = value
  }

  function setShowSoloSettings(value) {
    showSoloSettings.value = value
  }

  function setLoginTab(tab) {
    loginTab.value = tab
  }

  function setLbPeriod(period) {
    lbPeriod.value = period
  }

  function triggerEasterEgg() {
    easterEgg.value = true
    setTimeout(() => {
      easterEgg.value = false
    }, 4000)
  }

  function checkKonami(keyCode) {
    const KONAMI = [38, 38, 40, 40, 37, 39, 37, 39, 66, 65]

    if (keyCode === KONAMI[konamiIdx.value]) {
      const next = konamiIdx.value + 1
      if (next === KONAMI.length) {
        triggerEasterEgg()
        konamiIdx.value = 0
      } else {
        konamiIdx.value = next
      }
    } else {
      konamiIdx.value = 0
    }
  }

  function setUser(userData) {
    userId.value = userData.id
    username.value = userData.username
    email.value = userData.email || null
    joinedDate.value = new Date(userData.created_at).toLocaleDateString('en-US', { month: 'long', year: 'numeric' })
    isAuthenticated.value = true

    if (userData.default_letter_count) {
      soloLetters.value = userData.default_letter_count
    }
    if (userData.default_language) {
      soloLang.value = userData.default_language
    }
    if (userData.default_time_limit) {
      soloTime.value = userData.default_time_limit
    }

    localStorage.setItem('anagram_user', JSON.stringify(userData))
  }

  function signOut() {
    userId.value = null
    username.value = 'player_one'
    email.value = 'player@example.com'
    joinedDate.value = 'April 2026'
    isAuthenticated.value = false

    localStorage.removeItem('anagram_user')
  }

  return {
    userId,
    username,
    email,
    joinedDate,
    isAuthenticated,
    gamesPlayed,
    bestScore,
    longestWord,
    showHelp,
    showSoloSettings,
    soloTime,
    soloLetters,
    soloLang,
    loginTab,
    lbPeriod,
    easterEgg,

    setUser,
    setShowHelp,
    setShowSoloSettings,
    setLoginTab,
    setLbPeriod,
    triggerEasterEgg,
    checkKonami,
    signOut,
  }
})
