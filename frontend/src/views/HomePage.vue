<script setup>
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/userStore'
import { useGameStore } from '../stores/gameStore'

const router = useRouter()
const userStore = useUserStore()
const gameStore = useGameStore()

const modes = [
  {
    id: 'solo',
    icon: '◇',
    name: 'Solo game',
    desc: 'Custom time, letters & language',
    action: () => userStore.setShowSoloSettings(true)
  },
  {
    id: 'fast',
    icon: '⚡',
    name: 'Fast game',
    desc: 'Jump in instantly with defaults',
    action: () => {
      gameStore.startGame(60, 7, gameStore.settingsLang || 'ru')
      router.push('/game')
    }
  },
  {
    id: 'multi',
    icon: '◈',
    name: 'Multiplayer',
    desc: 'Challenge a friend via link',
    action: () => {},
    soon: true
  }
]
</script>

<template>
  <div class="home-center">
    <h1 class="title">ANAGRAM</h1>
    <p class="subtitle">Find all words. Beat your friends.</p>

    <div class="modes-grid">
      <div
        v-for="mode in modes"
        :key="mode.id"
        class="mode-card"
        @click="mode.action"
      >
        <span class="mode-icon">{{ mode.icon }}</span>
        <p class="mode-name">{{ mode.name }}</p>
        <p class="mode-desc">{{ mode.desc }}</p>
        <div v-if="mode.soon" class="soon-badge">SOON</div>
      </div>
    </div>

    <p class="version" title="Try the Konami code ↑↑↓↓←→←→BA">v0.1.0</p>
  </div>
</template>

<style scoped>
.home-center {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 140px);
  gap: 12px;
}

.title {
  font-family: 'Space Mono', monospace;
  font-size: 42px;
  font-weight: 700;
  letter-spacing: -1px;
  text-align: center;
  margin: 0;
  background: linear-gradient(135deg, #63e6be, #4ecdc4, #45b7aa);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  font-size: 15px;
  color: #666;
  text-align: center;
  margin: 4px 0 40px;
  font-weight: 300;
  letter-spacing: 0.5px;
}

.modes-grid {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  justify-content: center;
}

.mode-card {
  width: 200px;
  padding: 32px 24px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  cursor: pointer;
  transition: all 0.3s;
  text-align: center;
  position: relative;
  overflow: hidden;
}

.mode-card:hover {
  background: rgba(99, 230, 190, 0.08);
  border-color: rgba(99, 230, 190, 0.3);
  transform: translateY(-4px);
}

.mode-icon {
  font-size: 36px;
  margin-bottom: 16px;
  display: block;
}

.mode-name {
  font-family: 'Space Mono', monospace;
  font-size: 16px;
  font-weight: 700;
  margin: 0 0 8px;
}

.mode-desc {
  font-size: 12px;
  color: #777;
  line-height: 1.5;
  margin: 0;
}

.soon-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  font-size: 10px;
  padding: 3px 8px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.06);
  color: #555;
  font-weight: 600;
}

.version {
  font-size: 11px;
  color: #333;
  margin-top: 40px;
  font-family: 'Space Mono', monospace;
  letter-spacing: 2px;
  cursor: default;
}
</style>
