<script setup>
import { useUserStore } from '../stores/userStore'
import Toast from './Toast.vue'

const userStore = useUserStore()
</script>

<template>
  <Teleport to="body">
    <div v-if="userStore.toasts.length > 0" class="toast-stack">
      <Toast
        v-for="toast in userStore.toasts"
        :key="toast.id"
        :type="toast.type"
        :message="toast.message"
        :number="toast.number"
      />
    </div>
  </Teleport>
</template>
