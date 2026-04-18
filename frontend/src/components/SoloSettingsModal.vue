<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/userStore'
import { useGameStore } from '../stores/gameStore'

const router = useRouter()
const userStore = useUserStore()
const gameStore = useGameStore()

async function startSoloGame() {
  userStore.setShowSoloSettings(false)
  await gameStore.startGame(userStore.soloTime, userStore.soloLetters, userStore.soloLang)
  router.push('/game')
}
</script>

<template>
  <Teleport to="body">
    <div v-if="userStore.showSoloSettings" class="modal-overlay" @click="userStore.setShowSoloSettings(false)">
      <div class="modal-shell" @click.stop>
        <button class="modal-close-btn" @click="userStore.setShowSoloSettings(false)">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 6L6 18M6 6l12 12"/>
          </svg>
        </button>

        <h2 class="modal-title">Game settings</h2>

        <div class="field">
          <label class="field-label">Language</label>
          <div class="checkbox-row">
            <button
              v-for="lang in [{ id: 'en', label: 'English' }, { id: 'ru', label: 'Русский' }]"
              :key="lang.id"
              class="chip-toggle"
              :data-active="userStore.soloLang === lang.id"
              @click="userStore.soloLang = lang.id"
            >
              {{ lang.label }}
            </button>
          </div>
        </div>

        <div class="field">
          <label class="field-label">Letters</label>
          <div class="checkbox-row">
            <button
              v-for="n in [6, 7, 8, 9, 10]"
              :key="n"
              class="chip-toggle chip-toggle--mono"
              :data-active="userStore.soloLetters === n"
              @click="userStore.soloLetters = n"
            >
              {{ n }}
            </button>
          </div>
        </div>

        <div class="field" style="margin-bottom:24px">
          <label class="field-label">Time limit</label>
          <div class="checkbox-row">
            <button
              v-for="time in [{ val: 30, label: '30s' }, { val: 60, label: '1:00' }, { val: 90, label: '1:30' }, { val: 120, label: '2:00' }]"
              :key="time.val"
              class="chip-toggle chip-toggle--mono"
              :data-active="userStore.soloTime === time.val"
              @click="userStore.soloTime = time.val"
            >
              {{ time.label }}
            </button>
          </div>
        </div>

        <button class="btn btn--accent btn--block btn--lg" @click="startSoloGame">
          Start game
        </button>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
/* All styles are in app.css */
</style>
