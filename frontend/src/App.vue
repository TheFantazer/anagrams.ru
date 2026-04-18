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
  <div class="app">
    <!-- Background effects -->
    <div class="app-paper" />
    <div class="app-sun" />

    <!-- Main content wrapper -->
    <div class="app-main">
      <!-- Header -->
      <HeaderNav />

      <!-- Page content -->
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

<style scoped>
.easter-egg-banner {
  position: fixed;
  top: 84px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 200;
  background: var(--accent);
  color: var(--milk);
  padding: 12px 28px;
  border-radius: 999px;
  font-family: var(--font-mono);
  font-size: 13px;
  font-weight: 700;
  box-shadow: var(--shadow-banner);
  animation: toast-in 0.35s var(--ease-out);
}
</style>
