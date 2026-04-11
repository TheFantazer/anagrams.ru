import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  // User state
  const username = ref('player_one')
  const email = ref('player@example.com')
  const joinedDate = ref('April 2026')
  const isAuthenticated = ref(false)

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

  function signOut() {
    isAuthenticated.value = false
  }

  return {
    // State
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
    setShowHelp,
    setShowSoloSettings,
    setLoginTab,
    setLbPeriod,
    triggerEasterEgg,
    checkKonami,
    signOut,
  }
})
