<script setup>
import { defineProps, defineEmits, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: ''
  },
  subtitle: {
    type: String,
    default: ''
  },
  maxWidth: {
    type: String,
    default: '440px'
  },
  wide: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close'])

const handleEsc = (e) => {
  if (e.key === 'Escape' && props.show) {
    emit('close')
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleEsc)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEsc)
})

const handleOverlayClick = (e) => {
  if (e.target === e.currentTarget) {
    emit('close')
  }
}
</script>

<template>
  <Teleport to="body">
    <Transition name="fade">
      <div
        v-if="show"
        class="overlay"
        @click="handleOverlayClick"
      >
        <div
          class="modal"
          :style="{ maxWidth: wide ? '720px' : maxWidth }"
          @click.stop
        >
          <button
            class="modal-close"
            @click="emit('close')"
            type="button"
          >
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>

          <h2 v-if="title" class="modal-title">{{ title }}</h2>
          <p v-if="subtitle" class="modal-sub">{{ subtitle }}</p>

          <slot></slot>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
