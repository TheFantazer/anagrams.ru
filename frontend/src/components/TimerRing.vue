<script setup>
import { computed } from 'vue'
import { useGameStore } from '../stores/gameStore'

const gameStore = useGameStore()

const radius = 34
const circumference = 2 * Math.PI * radius

const strokeColor = computed(() => {
  const pct = gameStore.timerPercentage
  if (pct > 0.3) return '#63e6be'
  if (pct > 0.1) return '#fbbf24'
  return '#ef4444'
})

const strokeDashoffset = computed(() => {
  return circumference * (1 - gameStore.timerPercentage)
})

const formattedTime = computed(() => {
  const minutes = Math.floor(gameStore.timeLeft / 60)
  const seconds = gameStore.timeLeft % 60
  return `${minutes}:${String(seconds).padStart(2, '0')}`
})
</script>

<template>
  <div class="timer-ring">
    <svg width="80" height="80" viewBox="0 0 80 80">
      <circle
        cx="40"
        cy="40"
        :r="radius"
        fill="none"
        stroke="rgba(255,255,255,0.06)"
        stroke-width="4"
      />
      <circle
        cx="40"
        cy="40"
        :r="radius"
        fill="none"
        :stroke="strokeColor"
        stroke-width="4"
        stroke-linecap="round"
        :stroke-dasharray="circumference"
        :stroke-dashoffset="strokeDashoffset"
        transform="rotate(-90 40 40)"
        class="progress-ring"
      />
    </svg>
    <div class="timer-text" :style="{ color: strokeColor }">
      {{ formattedTime }}
    </div>
  </div>
</template>

<style scoped>
.timer-ring {
  position: relative;
  width: 80px;
  height: 80px;
}

.timer-text {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: 'Space Mono', monospace;
  font-size: 20px;
  font-weight: 700;
}

.progress-ring {
  transition: stroke-dashoffset 1s linear, stroke 0.5s;
}
</style>
