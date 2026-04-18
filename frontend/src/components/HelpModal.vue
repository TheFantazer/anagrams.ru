<script setup>
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../stores/userStore'

const { t } = useI18n()
const userStore = useUserStore()
</script>

<template>
  <Teleport to="body">
    <div v-if="userStore.showHelp" class="modal-overlay" @click="userStore.setShowHelp(false)">
      <div class="modal-shell" @click.stop>
        <button class="modal-close-btn" @click="userStore.setShowHelp(false)">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M18 6L6 18M6 6l12 12"/>
          </svg>
        </button>

        <h2 class="modal-title">{{ $t('help.title') }}</h2>

        <ol class="help-list">
          <li>{{ $t('help.objective.description') }}</li>
          <li>{{ $t('help.rules.minLength') }}</li>
          <li>{{ $t('help.rules.useLetters') }}</li>
          <li>{{ $t('help.rules.noRepeats') }}</li>
        </ol>

        <div class="help-scoring">
          <div class="help-score-cell">
            <span class="mono">100</span>
            <span class="lbl">{{ $t('help.scoring.3letters') }}</span>
          </div>
          <div class="help-score-cell">
            <span class="mono">400</span>
            <span class="lbl">{{ $t('help.scoring.4letters') }}</span>
          </div>
          <div class="help-score-cell">
            <span class="mono">1.2k</span>
            <span class="lbl">{{ $t('help.scoring.5letters') }}</span>
          </div>
          <div class="help-score-cell">
            <span class="mono">2k</span>
            <span class="lbl">{{ $t('help.scoring.6letters') }}</span>
          </div>
          <div class="help-score-cell">
            <span class="mono">2.8k</span>
            <span class="lbl">{{ $t('help.scoring.7letters') }}</span>
          </div>
        </div>

        <div class="help-keys">
          <span><kbd class="kbd">↵</kbd> {{ $t('help.controls.submit').toLowerCase() }}</span>
          <span><kbd class="kbd">Esc</kbd> {{ $t('help.controls.clear').toLowerCase() }}</span>
          <span><kbd class="kbd">⌫</kbd> {{ $t('help.controls.delete').toLowerCase() }}</span>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.help-list {
  list-style: none;
  padding: 0;
  margin: 0 0 20px;
  counter-reset: h;
}
.help-list li {
  counter-increment: h;
  padding: 10px 0 10px 36px;
  position: relative;
  font-size: 14px;
  color: var(--fg-secondary);
  border-bottom: 1px solid var(--border-hairline);
}
.help-list li:last-child { border-bottom: 0; }
.help-list li::before {
  content: counter(h, decimal-leading-zero);
  position: absolute;
  left: 0;
  top: 10px;
  font-family: var(--font-mono);
  font-size: 11px;
  font-weight: 700;
  color: var(--accent);
}
.help-list b {
  color: var(--fg-primary);
  font-weight: 600;
}

.help-scoring {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 6px;
  margin-bottom: 20px;
  background: var(--bg-card);
  border-radius: 12px;
  padding: 10px;
}
.help-score-cell {
  text-align: center;
  padding: 8px 4px;
  background: var(--bg-surface);
  border-radius: 8px;
}
.help-score-cell .mono {
  font-family: var(--font-mono);
  font-weight: 700;
  font-size: 18px;
  color: var(--accent);
  display: block;
}
.help-score-cell .lbl {
  font-size: 10px;
  color: var(--fg-muted);
  text-transform: uppercase;
  letter-spacing: 1px;
  font-weight: 600;
}

.help-keys {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
  font-size: 12px;
  color: var(--fg-muted);
  align-items: center;
}
</style>
