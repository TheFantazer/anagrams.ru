<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/userStore'
import { useGameStore } from '../stores/gameStore'

const router = useRouter()
const userStore = useUserStore()
const gameStore = useGameStore()

function startSoloGame() {
  userStore.setShowSoloSettings(false)
  gameStore.startGame(userStore.soloTime, userStore.soloLetters, userStore.soloLang)
  router.push('/game')
}
</script>

<template>
  <Teleport to="body">
    <div v-if="userStore.showSoloSettings" class="overlay" @click="userStore.setShowSoloSettings(false)">
      <div class="modal" @click.stop>
        <button class="modal-close" @click="userStore.setShowSoloSettings(false)">&times;</button>
        <h3 class="modal-title">Game settings</h3>

        <label class="label">Language</label>
        <select v-model="userStore.soloLang" class="select">
          <option value="ru">Русский</option>
          <option value="en">English</option>
        </select>

        <label class="label">Time</label>
        <select v-model.number="userStore.soloTime" class="select">
          <option :value="30">30 seconds</option>
          <option :value="60">60 seconds</option>
          <option :value="120">120 seconds</option>
        </select>

        <label class="label">Letters</label>
        <select v-model.number="userStore.soloLetters" class="select">
          <option v-for="n in [6, 7, 8, 9, 10]" :key="n" :value="n">{{ n }} letters</option>
        </select>

        <button class="btn-primary" @click="startSoloGame">
          Start game
        </button>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.7);
  backdrop-filter: blur(8px);
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal {
  background: #141419;
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  padding: 32px;
  max-width: 420px;
  width: calc(100% - 48px);
  position: relative;
  max-height: 90vh;
  overflow: auto;
}

.modal-title {
  font-family: 'Space Mono', monospace;
  font-size: 20px;
  font-weight: 700;
  margin: 0 0 24px;
  color: var(--accent);
}

.modal-close {
  position: absolute;
  top: 16px;
  right: 16px;
  background: none;
  border: none;
  color: #666;
  font-size: 20px;
  cursor: pointer;
  padding: 4px;
  transition: color 0.2s;
}

.modal-close:hover {
  color: #e8e6e1;
}

.label {
  font-size: 12px;
  color: #888;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin: 0 0 8px;
  display: block;
  font-weight: 500;
}

.select {
  width: 100%;
  padding: 12px 16px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #e8e6e1;
  font-size: 14px;
  outline: none;
  font-family: 'Outfit', sans-serif;
  appearance: none;
  margin-bottom: 20px;
  cursor: pointer;
}

.btn-primary {
  width: 100%;
  padding: 14px;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--accent), var(--accent-hover));
  border: none;
  color: var(--bg-dark);
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  font-family: 'Outfit', sans-serif;
  transition: all 0.2s;
  letter-spacing: 0.5px;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(99, 230, 190, 0.3);
}
</style>
