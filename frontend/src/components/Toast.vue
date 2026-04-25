<script setup>
import { defineProps } from 'vue'

const props = defineProps({
  type: {
    type: String,
    default: 'default',
    validator: (value) => ['default', 'success', 'warn', 'error'].includes(value)
  },
  message: {
    type: String,
    required: true
  },
  number: {
    type: Number,
    default: null
  }
})

// Icons for different toast types
const getIcon = () => {
  switch (props.type) {
    case 'success':
      return `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
        <polyline points="20 6 9 17 4 12"></polyline>
      </svg>`
    case 'error':
      return `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="12" cy="12" r="10"></circle>
        <line x1="15" y1="9" x2="9" y2="15"></line>
        <line x1="9" y1="9" x2="15" y2="15"></line>
      </svg>`
    case 'warn':
      return `<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
        <path d="M10.29 3.86L1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"></path>
        <line x1="12" y1="9" x2="12" y2="13"></line>
        <line x1="12" y1="17" x2="12.01" y2="17"></line>
      </svg>`
    default:
      return ''
  }
}
</script>

<template>
  <div
    :class="[
      'toast',
      `toast--${type}`
    ]"
  >
    <span v-if="getIcon()" v-html="getIcon()"></span>
    <span>{{ message }}</span>
    <span v-if="number !== null" class="toast-num">{{ number }}</span>
  </div>
</template>
