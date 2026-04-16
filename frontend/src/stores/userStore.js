import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const userId = ref(null)
  const username = ref('')
  const email = ref('')
  const joinedDate = ref('')
  const isAuthenticated = ref(false)

  const gamesPlayed = ref(0)
  const bestScore = ref(0)
  const longestWord = ref('-')

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

  async function loadUserFromToken() {
    const token = localStorage.getItem('access_token')
    if (!token) return

    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      const userIdFromToken = payload.user_id

      if (!userIdFromToken) {
        console.error('No user_id in token')
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        return
      }

      const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
      const response = await fetch(`${apiUrl}/api/v1/auth/me?user_id=${userIdFromToken}`)

      if (response.ok) {
        const userData = await response.json()
        setUser(userData)
      } else {
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        localStorage.removeItem('anagram_user')
      }
    } catch (error) {
      console.error('Failed to load user from token', error)
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      localStorage.removeItem('anagram_user')
    }
  }

  function signOut() {
    userId.value = null
    username.value = ''
    email.value = ''
    joinedDate.value = ''
    isAuthenticated.value = false

    localStorage.removeItem('anagram_user')
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
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
    loadUserFromToken,
    setShowHelp,
    setShowSoloSettings,
    setLoginTab,
    setLbPeriod,
    triggerEasterEgg,
    checkKonami,
    signOut,
  }
})
