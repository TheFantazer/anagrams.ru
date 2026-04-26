<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  words: {
    type: Array,
    required: true,
    default: () => []
  },
  owner: {
    type: String,
    required: true,
    validator: (value) => ['you', 'them'].includes(value)
  }
})

// Group by length
function groupByLength(words) {
  const groups = {}
  words.forEach(w => {
    (groups[w.length] = groups[w.length] || []).push(w)
  })
  const lens = Object.keys(groups).map(Number).sort((a, b) => b - a)
  return lens.map(l => ({ length: l, words: groups[l] }))
}

const groups = computed(() => groupByLength(props.words))
</script>

<template>
  <div class="res-chip-groups">
    <div v-for="g in groups" :key="g.length" class="res-chip-group">
      <div class="res-chip-group-head">
        <span class="mono">{{ g.length }}</span>
        <span class="lbl">{{ g.length > 1 ? t('results.groupLabel.letters') : t('results.groupLabel.letter') }}</span>
        <span class="cnt">{{ g.words.length }}</span>
      </div>
      <div class="res-chip-row">
        <span v-for="(w, i) in g.words" :key="i" :class="['r-chip', `r-chip--${owner}`]">
          {{ w.toLowerCase() }}
        </span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.res-chip-groups {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.res-chip-group-head {
  display: flex;
  align-items: baseline;
  gap: 6px;
  margin-bottom: 6px;
}

.res-chip-group-head .mono {
  font-family: var(--font-mono);
  font-weight: 700;
  font-size: 11px;
  color: var(--navy);
  background: var(--bg-card);
  padding: 2px 8px;
  border-radius: 6px;
}

.res-chip-group-head .lbl {
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: var(--fg-muted);
  font-weight: 600;
}

.res-chip-group-head .cnt {
  margin-left: auto;
  font-family: var(--font-mono);
  font-size: 11px;
  color: var(--fg-faint);
}

.res-chip-row {
  display: flex;
  flex-wrap: wrap;
  gap: 5px;
}

.r-chip {
  padding: 4px 9px;
  border-radius: 7px;
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 1px;
}

.r-chip--you {
  background: var(--accent-soft);
  color: var(--accent);
  border: 1px solid var(--border-accent);
}

.r-chip--them {
  background: var(--bg-card);
  color: var(--navy);
  border: 1px solid var(--border-subtle);
}
</style>
