<script setup>
import { ref } from 'vue'
import { useUserStore } from '../stores/userStore'

const userStore = useUserStore()

const leaderboard = ref([
  { name: 'wordmaster_42', score: 4200, words: 18 },
  { name: 'лексикон', score: 3800, words: 15 },
  { name: 'anagram_pro', score: 3100, words: 14 },
  { name: 'буквоед', score: 2900, words: 12 },
  { name: 'solver99', score: 2400, words: 11 },
  { name: 'слововед', score: 2100, words: 10 },
  { name: 'quicktype', score: 1800, words: 9 },
])

const periods = ['day', 'week', 'month', 'all']

function getPeriodLabel(period) {
  if (period === 'all') return 'All time'
  return period.charAt(0).toUpperCase() + period.slice(1)
}
</script>

<template>
  <div class="lb-page">
    <h2 class="page-title">Leaderboard</h2>

    <div class="lb-tabs">
      <button
        v-for="period in periods"
        :key="period"
        :class="['lb-tab', { active: userStore.lbPeriod === period }]"
        @click="userStore.setLbPeriod(period)"
      >
        {{ getPeriodLabel(period) }}
      </button>
    </div>

    <div
      v-for="(user, i) in leaderboard"
      :key="i"
      :class="['lb-row', { first: i === 0 }]"
    >
      <span :class="['lb-rank', `rank-${i}`]">
        {{ i === 0 ? '👑' : `#${i + 1}` }}
      </span>
      <span class="lb-name">{{ user.name }}</span>
      <span class="lb-words">{{ user.words }} words</span>
      <span class="lb-score">{{ user.score.toLocaleString() }}</span>
    </div>
  </div>
</template>

<style scoped>
.lb-page {
  max-width: 600px;
  margin: 40px auto;
  padding: 0 24px;
}

.page-title {
  font-family: 'Space Mono', monospace;
  font-size: 22px;
  font-weight: 700;
  margin: 0 0 24px;
  color: #e8e6e1;
}

.lb-tabs {
  display: flex;
  gap: 6px;
  margin-bottom: 24px;
}

.lb-tab {
  padding: 8px 16px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  border: none;
  background: rgba(255, 255, 255, 0.03);
  color: #666;
  font-family: 'Outfit', sans-serif;
  transition: all 0.2s;
}

.lb-tab.active {
  background: rgba(99, 230, 190, 0.12);
  color: var(--accent);
}

.lb-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 16px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.04);
  margin-bottom: 8px;
}

.lb-row.first {
  background: rgba(99, 230, 190, 0.06);
  border-color: rgba(99, 230, 190, 0.15);
}

.lb-rank {
  font-family: 'Space Mono', monospace;
  font-size: 16px;
  font-weight: 700;
  width: 32px;
  text-align: center;
  color: #555;
}

.lb-rank.rank-0 {
  color: var(--accent);
}

.lb-rank.rank-1 {
  color: #fbbf24;
}

.lb-rank.rank-2 {
  color: #fb923c;
}

.lb-name {
  flex: 1;
  font-size: 14px;
  font-weight: 500;
}

.lb-words {
  font-size: 12px;
  color: #555;
  margin-right: 8px;
}

.lb-score {
  font-family: 'Space Mono', monospace;
  font-size: 14px;
  font-weight: 700;
  color: var(--accent);
}
</style>
