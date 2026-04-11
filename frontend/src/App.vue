<script setup>
import { onMounted, onUnmounted } from 'vue'
import { RouterView } from 'vue-router'
import { useUserStore } from './stores/userStore'
import HeaderNav from './components/HeaderNav.vue'
import HelpModal from './components/HelpModal.vue'
import SoloSettingsModal from './components/SoloSettingsModal.vue'

const userStore = useUserStore()

function handleKeyDown(e) {
  userStore.checkKonami(e.keyCode)
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
</script>

<template>
  <div class="app-root">
    <!-- Background effects -->
    <div class="noise-overlay" />
    <div class="glow-effect" />

    <!-- Header -->
    <HeaderNav />

    <!-- Main content -->
    <div class="content">
      <RouterView />
    </div>

    <!-- Modals -->
    <HelpModal />
    <SoloSettingsModal />

    <!-- Easter egg banner -->
    <Teleport to="body">
      <div v-if="userStore.easterEgg" class="easter-egg-banner">
        🎮 +9999 pts — You found the secret!
      </div>
    </Teleport>
  </div>
</template>

<style>
* {
  box-sizing: border-box;
}

body {
  margin: 0;
  font-family: 'Outfit', sans-serif;
  background: var(--bg-dark);
  color: var(--text-primary);
  overflow-x: hidden;
}

::selection {
  background: rgba(99, 230, 190, 0.3);
}

input::placeholder {
  color: #333;
}

select option {
  background: #1a1a22;
  color: #e8e6e1;
}

::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.08);
  border-radius: 3px;
}
</style>

<style scoped>
.app-root {
  font-family: 'Outfit', sans-serif;
  background: var(--bg-dark);
  color: var(--text-primary);
  min-height: 100vh;
  overflow: hidden;
}

.noise-overlay {
  position: fixed;
  inset: 0;
  background: url('data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" width="200" height="200"><filter id="n"><feTurbulence baseFrequency="0.9" numOctaves="4" stitchTiles="stitch"/></filter><rect width="200" height="200" filter="url(%23n)" opacity="0.03"/></svg>');
  pointer-events: none;
  z-index: 0;
}

.glow-effect {
  position: fixed;
  top: -30%;
  left: 50%;
  transform: translateX(-50%);
  width: 800px;
  height: 500px;
  background: radial-gradient(ellipse, rgba(99, 230, 190, 0.06) 0%, transparent 70%);
  pointer-events: none;
  z-index: 0;
}

.content {
  position: relative;
  z-index: 5;
  max-width: 900px;
  margin: 0 auto;
  padding: 0 24px;
}

.easter-egg-banner {
  position: fixed;
  top: 80px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 200;
  background: linear-gradient(135deg, var(--accent), var(--accent-hover));
  color: var(--bg-dark);
  padding: 12px 28px;
  border-radius: 12px;
  font-family: 'Space Mono', monospace;
  font-size: 14px;
  font-weight: 700;
  box-shadow: 0 8px 32px rgba(99, 230, 190, 0.3);
  animation: slideDown 0.4s ease;
}
</style>
