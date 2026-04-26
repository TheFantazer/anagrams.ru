<script setup>
import { computed } from 'vue'

const props = defineProps({
  page: {
    type: Number,
    required: true
  },
  totalPages: {
    type: Number,
    required: true
  },
  showing: {
    type: Number,
    required: true
  },
  total: {
    type: Number,
    required: true
  },
  pageSize: {
    type: Number,
    default: 15
  }
})

const emit = defineEmits(['update:page'])

const from = computed(() => (props.page - 1) * props.pageSize + 1)
const to = computed(() => (props.page - 1) * props.pageSize + props.showing)

// Build page numbers with ellipsis
const pages = computed(() => {
  const result = []
  for (let i = 1; i <= props.totalPages; i++) {
    if (i === 1 || i === props.totalPages || Math.abs(i - props.page) <= 1) {
      result.push(i)
    } else if (result[result.length - 1] !== '…') {
      result.push('…')
    }
  }
  return result
})

function goToPage(page) {
  if (page !== '…' && page >= 1 && page <= props.totalPages) {
    emit('update:page', page)
  }
}
</script>

<template>
  <div v-if="total > 0" class="pagination">
    <div class="pagination-info muted">
      Showing <span class="mono">{{ from }}–{{ to }}</span> of <span class="mono">{{ total }}</span>
    </div>
    <div class="pagination-ctrls">
      <button
        class="page-btn"
        :disabled="page === 1"
        @click="goToPage(page - 1)"
      >
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
          <path d="M19 12H5M12 19l-7-7 7-7"/>
        </svg>
      </button>

      <template v-for="(p, i) in pages" :key="i">
        <span v-if="p === '…'" class="page-ellipsis">…</span>
        <button
          v-else
          :class="['page-btn', { 'is-active': p === page }]"
          @click="goToPage(p)"
        >
          {{ p }}
        </button>
      </template>

      <button
        class="page-btn"
        :disabled="page === totalPages"
        @click="goToPage(page + 1)"
      >
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
          <path d="M5 12h14M12 5l7 7-7 7"/>
        </svg>
      </button>
    </div>
  </div>
</template>

<style scoped>
.pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--sp-4);
  flex-wrap: wrap;
  padding-top: var(--sp-3);
}

.pagination-info {
  font-size: 12px;
}

.pagination-ctrls {
  display: flex;
  align-items: center;
  gap: 4px;
}

.page-btn {
  appearance: none;
  min-width: 32px;
  height: 32px;
  padding: 0 8px;
  border-radius: var(--radius-sm);
  border: 1px solid var(--border-hairline);
  background: var(--bg-surface);
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 700;
  color: var(--fg-secondary);
  cursor: pointer;
  display: inline-grid;
  place-items: center;
  transition: all var(--dur-fast) var(--ease-std);
}

.page-btn:hover:not(:disabled) {
  background: var(--bg-card);
  border-color: var(--border-default);
  color: var(--fg-primary);
}

.page-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.page-btn.is-active {
  background: var(--navy);
  border-color: var(--navy);
  color: var(--milk);
}

.page-ellipsis {
  padding: 0 6px;
  color: var(--fg-muted);
  font-size: 12px;
}
</style>
