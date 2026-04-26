<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  who: {
    type: String,
    required: true,
    validator: (value) => ['you', 'them'].includes(value)
  },
  name: {
    type: String,
    required: true
  },
  score: {
    type: Number,
    required: true
  },
  found: {
    type: Number,
    required: true
  },
  total: {
    type: Number,
    required: true
  },
  winner: {
    type: Boolean,
    default: false
  },
  solo: {
    type: Boolean,
    default: false
  }
})

const pct = computed(() => Math.round((props.found * 100) / Math.max(1, props.total)))
</script>

<template>
  <div :class="['res-pb', `res-pb--${who}`, { 'is-winner': winner, 'is-solo': solo }]">
    <div class="res-pb-top">
      <div class="res-pb-avatar">{{ name[0]?.toUpperCase() }}</div>
      <div class="res-pb-name">{{ name }}</div>
      <div v-if="winner" class="res-pb-crown">👑</div>
    </div>
    <div class="res-pb-score mono">{{ score.toLocaleString() }}</div>
    <div class="res-pb-meta">
      <div>
        <span class="mono">{{ found }}</span>
        <span class="lbl">{{ t('results.playerBlock.words') }}</span>
      </div>
      <div class="sep" />
      <div>
        <span class="mono">{{ pct }}%</span>
        <span class="lbl">{{ t('results.playerBlock.found') }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.res-pb {
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 20px;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  position: relative;
  transition: all var(--dur-base);
}

.res-pb.is-winner {
  background: var(--navy);
  color: var(--milk);
  border-color: var(--navy);
}

.res-pb.is-winner .res-pb-name {
  color: var(--milk);
}

.res-pb.is-winner .res-pb-score {
  color: var(--milk);
}

.res-pb.is-winner .res-pb-meta .mono {
  color: var(--milk);
}

.res-pb.is-winner .res-pb-meta .lbl {
  color: color-mix(in oklab, var(--milk) 60%, transparent);
}

.res-pb.is-winner .sep {
  background: rgba(255, 255, 255, 0.18);
}

.res-pb-top {
  display: flex;
  align-items: center;
  gap: 12px;
}

.res-pb-avatar {
  width: 40px;
  height: 40px;
  border-radius: 999px;
  background: var(--grad-accent);
  color: var(--milk);
  font-family: var(--font-display);
  font-weight: 700;
  font-size: 16px;
  display: grid;
  place-items: center;
  border: 2px solid var(--milk);
  box-shadow: 0 0 0 1px var(--border-default);
}

.res-pb.is-winner .res-pb-avatar {
  border-color: var(--navy);
  box-shadow: 0 0 0 1px rgba(255, 255, 255, 0.18);
}

.res-pb-name {
  font-family: var(--font-display);
  font-weight: 600;
  font-size: 18px;
  color: var(--fg-primary);
  flex: 1;
}

.res-pb-crown {
  font-size: 22px;
}

.res-pb-score {
  font-family: var(--font-mono);
  font-size: 48px;
  font-weight: 700;
  letter-spacing: -1.5px;
  line-height: 1;
  color: var(--accent);
  margin: 4px 0;
}

.res-pb.is-winner .res-pb-score {
  color: var(--milk);
}

.res-pb-meta {
  display: flex;
  align-items: center;
  gap: 16px;
  padding-top: 12px;
  border-top: 1px solid var(--border-hairline);
}

.res-pb.is-winner .res-pb-meta {
  border-top-color: rgba(255, 255, 255, 0.12);
}

.res-pb-meta > div {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.res-pb-meta .mono {
  font-family: var(--font-mono);
  font-size: 18px;
  font-weight: 700;
  color: var(--fg-primary);
}

.res-pb-meta .lbl {
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 1.2px;
  color: var(--fg-muted);
  font-weight: 600;
}

.res-pb-meta .sep {
  width: 1px;
  align-self: stretch;
  background: var(--border-hairline);
}

@media (max-width: 720px) {
  .res-pb {
    padding: 18px;
  }

  .res-pb-score {
    font-size: 36px;
  }
}
</style>
