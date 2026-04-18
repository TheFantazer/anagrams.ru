<script setup>
import { computed } from 'vue'
import { useGameStore } from '../stores/gameStore'

const gameStore = useGameStore()

const radius = 42
const circumference = 2 * Math.PI * radius

const timerState = computed(() => {
  const pct = gameStore.timerPercentage
  if (pct > 0.3) return 'ok'
  if (pct > 0.1) return 'low'
  return 'crit'
})

const strokeColor = computed(() => {
  const state = timerState.value
  if (state === 'ok') return 'var(--accent)'
  if (state === 'low') return 'var(--warning)'
  return 'var(--danger)'
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
  <div
    :class="['timer-ring', {
      'low': timerState === 'low',
      'crit': timerState === 'crit'
    }]"
  >
    <svg width="96" height="96" viewBox="0 0 96 96" preserveAspectRatio="xMidYMid meet">
      <circle
        cx="48"
        cy="48"
        :r="radius"
        fill="none"
        stroke="var(--border-hairline)"
        stroke-width="4"
      />
      <circle
        cx="48"
        cy="48"
        :r="radius"
        fill="none"
        :stroke="strokeColor"
        stroke-width="4"
        stroke-linecap="round"
        :stroke-dasharray="circumference"
        :stroke-dashoffset="strokeDashoffset"
        style="transition: stroke-dashoffset 0.95s linear, stroke 0.3s"
      />
    </svg>
    <div class="timer-label">
      {{ formattedTime }}
    </div>
  </div>
</template>

<style scoped>
/* Styles are in game.css */
</style>
