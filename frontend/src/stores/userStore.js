import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  // User state
  const userId = ref(null)
  const username = ref('player_one')
  const email = ref('player@example.com')
  const joinedDate = ref('April 2026')
  const isAuthenticated = ref(false)

  // Load user from localStorage on init
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
      } catch (e) {
        console.error('Failed to load user from localStorage', e)
      }
    }
  }

  // Load user on store creation
  loadUser()

  // Stats
  const gamesPlayed = ref(42)
  const bestScore = ref(3200)
  const longestWord = ref('КРОКОДИЛ')

  // Modals
  const showHelp = ref(false)
  const showSoloSettings = ref(false)

  // Solo game settings
  const soloTime = ref(60)
  const soloLetters = ref(7)
  const soloLang = ref('ru')

  // Auth
  const loginTab = ref('login')

  // Leaderboard
  const lbPeriod = ref('week')

  // Easter egg
  const easterEgg = ref(false)
  const konamiIdx = ref(0)

  // Actions
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

    // Save to localStorage
    localStorage.setItem('anagram_user', JSON.stringify(userData))
  }

  function signOut() {
    userId.value = null
    username.value = 'player_one'
    email.value = 'player@example.com'
    joinedDate.value = 'April 2026'
    isAuthenticated.value = false

    // Clear localStorage
    localStorage.removeItem('anagram_user')
  }

  return {
    // State
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

    // Actions
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
