<script setup>
import { ref, watch, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../stores/userStore'

const { t } = useI18n()
const userStore = useUserStore()

const leaderboard = ref([])
const loading = ref(false)

const periodOptions = computed(() => [
  { id: 'day', label: t('leaderboard.periods.today') },
  { id: 'week', label: t('leaderboard.periods.week') },
  { id: 'month', label: t('leaderboard.periods.month') },
  { id: 'all', label: t('leaderboard.periods.allTime') }
])

async function loadLeaderboard() {
  loading.value = true
  try {
    const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
    const response = await fetch(`${apiUrl}/api/v1/leaderboard?period=${userStore.lbPeriod}`)

    if (response.ok) {
      const data = await response.json()
      leaderboard.value = data || []
    } else {
      leaderboard.value = []
    }
  } catch (error) {
    console.error('Failed to load leaderboard:', error)
    leaderboard.value = []
  } finally {
    loading.value = false
  }
}

watch(() => userStore.lbPeriod, () => {
  loadLeaderboard()
})

onMounted(() => {
  loadLeaderboard()
})
</script>

<template>
  <div class="page">
    <div class="shell lb-wrap">
      <header class="page-head">
        <div>
          <div class="page-eyebrow">{{ $t('leaderboard.title') }}</div>
          <h1 class="page-title-display">{{ $t('leaderboard.subtitle') }}</h1>
        </div>
        <div class="lb-tabs">
          <button
            v-for="period in periodOptions"
            :key="period.id"
            class="chip-toggle"
            :data-active="userStore.lbPeriod === period.id"
            @click="userStore.setLbPeriod(period.id)"
          >
            {{ period.label }}
          </button>
        </div>
      </header>

      <!-- Loading/Empty state -->
      <div v-if="loading" style="text-align:center; padding:60px; color:var(--fg-muted)">
        {{ $t('leaderboard.loading') }}
      </div>
      <div v-else-if="leaderboard.length === 0" style="text-align:center; padding:60px; color:var(--fg-muted)">
        No results yet for this period
      </div>

      <!-- Podium (top 3) -->
      <div v-else-if="leaderboard.length > 0" class="lb-podium">
        <div
          v-for="(user, i) in leaderboard.slice(0, 3)"
          :key="i"
          :class="['podium', `podium-${i + 1}`]"
          :data-you="user.name === userStore.username"
        >
          <div class="podium-rank">{{ i === 0 ? '👑' : `0${i + 1}` }}</div>
          <div class="podium-name">{{ user.name }}</div>
          <div class="podium-score">{{ user.score.toLocaleString() }}</div>
          <div class="podium-words">{{ user.words }} {{ $t('leaderboard.columns.words').toLowerCase() }}</div>
        </div>
      </div>

      <!-- Table (rest of leaderboard) -->
      <div v-if="!loading && leaderboard.length > 0" class="lb-table">
        <div class="lb-row lb-row--head">
          <span style="width:40px">{{ $t('leaderboard.columns.rank') }}</span>
          <span style="flex:1">{{ $t('leaderboard.columns.player') }}</span>
          <span style="width:80px; text-align:right">{{ $t('leaderboard.columns.words') }}</span>
          <span style="width:100px; text-align:right">{{ $t('leaderboard.columns.score') }}</span>
        </div>
        <div
          v-for="(user, i) in leaderboard"
          :key="i"
          :class="['lb-row', { you: user.name === userStore.username }]"
        >
          <span :class="['lb-rank', `r-${i}`]">
            {{ i < 3 ? ['①', '②', '③'][i] : String(i + 1).padStart(2, '0') }}
          </span>
          <span class="lb-name">
            <span
              class="lb-avatar"
              :style="{ background: `hsl(${i * 37} 40% 40%)` }"
            >
              {{ user.name[0].toUpperCase() }}
            </span>
            {{ user.name }}
            <span v-if="user.name === userStore.username" class="lb-youtag">{{ $t('leaderboard.you') }}</span>
          </span>
          <span class="lb-words">{{ user.words }}</span>
          <span class="lb-score mono">{{ user.score.toLocaleString() }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* All styles are in pages.css and app.css */
</style>
