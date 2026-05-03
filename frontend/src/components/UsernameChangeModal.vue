<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Modal from './Modal.vue'

const { t } = useI18n()

const props = defineProps({
  show: {
    type: Boolean,
    required: true
  },
  currentUsername: {
    type: String,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'confirm'])

const newUsername = ref('')
const error = ref('')

const isValid = computed(() => {
  if (newUsername.value.length < 3 || newUsername.value.length > 30) {
    return false
  }
  if (newUsername.value === props.currentUsername) {
    return false
  }
  return true
})

function handleConfirm() {
  if (!isValid.value) {
    error.value = t('settings.usernameChange.invalidLength')
    return
  }

  error.value = ''
  emit('confirm', newUsername.value)
}

function handleClose() {
  newUsername.value = ''
  error.value = ''
  emit('close')
}
</script>

<template>
  <Modal :show="show" @close="handleClose" maxWidth="460px">
    <template #title>
      {{ $t('settings.usernameChange.title') }}
    </template>

    <div style="margin-bottom: 24px">
      <p class="muted" style="font-size: 14px; margin-bottom: 16px">
        {{ $t('settings.usernameChange.warning') }}
      </p>

      <div class="field">
        <label class="field-label">{{ $t('settings.usernameChange.currentUsername') }}</label>
        <input
          type="text"
          class="input"
          :value="currentUsername"
          disabled
          style="opacity: 0.6"
        />
      </div>

      <div class="field">
        <label class="field-label">{{ $t('settings.usernameChange.newUsername') }}</label>
        <input
          v-model="newUsername"
          type="text"
          class="input"
          :placeholder="$t('settings.usernameChange.placeholder')"
          autofocus
          @keyup.enter="handleConfirm"
        />
      </div>

      <div v-if="error" style="color: var(--danger); font-size: 13px; margin-top: 8px">
        {{ error }}
      </div>
    </div>

    <div style="display: flex; gap: 12px">
      <button
        type="button"
        class="btn btn--ghost"
        @click="handleClose"
        :disabled="loading"
      >
        {{ $t('common.cancel') }}
      </button>
      <button
        type="button"
        class="btn btn--accent"
        style="flex: 1"
        @click="handleConfirm"
        :disabled="!isValid || loading"
      >
        {{ loading ? $t('common.loading') : $t('settings.usernameChange.confirm') }}
      </button>
    </div>
  </Modal>
</template>
