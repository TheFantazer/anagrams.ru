<script setup>
import { useUserStore } from '../stores/userStore'

const userStore = useUserStore()
</script>

<template>
  <Teleport to="body">
    <div v-if="userStore.showHelp" class="overlay" @click="userStore.setShowHelp(false)">
      <div class="modal" @click.stop>
        <button class="modal-close" @click="userStore.setShowHelp(false)">&times;</button>
        <h3 class="modal-title">How to play</h3>
        <div class="modal-content">
          <p><span class="step">1.</span> You get a set of letters and a timer.</p>
          <p><span class="step">2.</span> Make as many words as you can (3+ letters).</p>
          <p><span class="step">3.</span> Longer words = more points.</p>
          <p><span class="step">4.</span> Share the link — your friend plays the same letters.</p>
          <p><span class="step">5.</span> Compare results. Best score wins.</p>
          <div class="scoring-box">
            <p>3 letters → 100 pts &nbsp;·&nbsp; 4 → 300 &nbsp;·&nbsp; 5 → 500 &nbsp;·&nbsp; 6+ → 1000+</p>
          </div>
        </div>
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

.modal-content {
  font-size: 14px;
  line-height: 1.8;
  color: #aaa;
}

.modal-content p {
  margin: 0 0 12px;
}

.step {
  color: var(--accent);
  font-weight: 600;
}

.scoring-box {
  background: rgba(99, 230, 190, 0.06);
  border-radius: 10px;
  padding: 16px;
  border: 1px solid rgba(99, 230, 190, 0.1);
  margin-top: 20px;
}

.scoring-box p {
  margin: 0;
  font-size: 12px;
  color: #888;
  font-family: 'Space Mono', monospace;
}
</style>
