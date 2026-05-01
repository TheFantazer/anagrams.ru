<script setup>
import { useRouter } from 'vue-router'

const props = defineProps({
  text: {
    type: String,
    required: true
  },
  icon: {
    type: String,
    default: '/login-svgrepo-com.svg'
  }
})

const router = useRouter()
</script>

<template>
  <div class="ph-empty">
    <div class="icon">
      <img :src="icon" alt="">
    </div>

    <div class="ph-empty-text">
      <strong>{{ text }}</strong>
    </div>

    <button class="btn btn--accent btn--sm" @click="router.push('/auth')">
      Sign in
    </button>
  </div>
</template>

<style scoped>
.ph-empty {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  background: var(--bg-card);
  border: 1px dashed var(--border-default);
  border-radius: 16px;
}

.ph-empty-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  flex: 1;
}

.ph-empty-text strong {
  font-size: 14px;
  font-weight: 700;
  color: var(--fg-primary);
}

.icon {
  width: 48px;
  height: 48px;
  flex-shrink: 0;
  background: white;
  border-radius: 10px;
  display: grid;
  place-items: center;
}

.icon img {
  width: 70%;
  height: 70%;
  object-fit: contain;
}
</style>