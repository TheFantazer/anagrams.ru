<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '../stores/userStore'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()

const username = ref('')
const email = ref('')
const password = ref('')
const acceptedPrivacy = ref(false)
const error = ref('')
const loading = ref(false)

onMounted(() => {
  const script = document.createElement('script')
  script.src = 'https://telegram.org/js/telegram-widget.js?22'
  script.async = true
  script.setAttribute('data-telegram-login', import.meta.env.VITE_TELEGRAM_BOT_USERNAME || 'YOUR_BOT_USERNAME')
  script.setAttribute('data-size', 'large')
  script.setAttribute('data-auth-url', `${import.meta.env.VITE_API_URL || 'http://localhost:8080'}/api/v1/auth/telegram/callback`)
  script.setAttribute('data-request-access', 'write')

  const container = document.getElementById('telegram-login-container')
  if (container) {
    container.appendChild(script)
  }
})

async function handleSubmit() {
  error.value = ''
  loading.value = true

  try {
    if (userStore.loginTab === 'register') {
      await register()
    } else {
      await login()
    }
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}

async function register() {
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const response = await fetch(`${apiUrl}/api/v1/auth/register`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      username: username.value,
      email: email.value,
      password: password.value,
      accepted_privacy_policy: true
    })
  })

  const data = await response.json()

  if (!response.ok) {
    throw new Error(data.message || 'Registration failed')
  }

  userStore.setUser(data)
  router.push('/')
}

async function login() {
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  const response = await fetch(`${apiUrl}/api/v1/auth/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      username: username.value,
      password: password.value
    })
  })

  const data = await response.json()

  if (!response.ok) {
    throw new Error(data.message || 'Login failed')
  }

  userStore.setUser(data)

  // Проверяем, есть ли сохранённый путь для редиректа
  const redirectPath = sessionStorage.getItem('redirectAfterAuth')
  if (redirectPath) {
    sessionStorage.removeItem('redirectAfterAuth')
    router.push(redirectPath)
  } else {
    router.push('/')
  }
}

const eyebrowText = computed(() =>
  userStore.loginTab === 'login' ? t('auth.welcomeBack') : t('auth.joinTheBoard')
)

const titleText = computed(() =>
  userStore.loginTab === 'login' ? t('auth.signInTitle') : t('auth.signUpTitle')
)

const usernamePlaceholder = computed(() =>
  userStore.loginTab === 'login' ? t('auth.placeholders.usernameOrEmail') : t('auth.placeholders.username')
)

const usernameLabel = computed(() =>
  userStore.loginTab === 'login' ? t('auth.fields.usernameOrEmail') : t('auth.fields.username')
)

const submitButtonText = computed(() =>
  loading.value ? t('auth.actions.loading') : (userStore.loginTab === 'login' ? t('auth.actions.signIn') : t('auth.actions.createAccount'))
)

const isSubmitDisabled = computed(() => {
  if (loading.value) return true
  if (userStore.loginTab === 'register' && !acceptedPrivacy.value) return true
  return false
})

const toggleText = computed(() =>
  userStore.loginTab === 'login' ? t('auth.toggles.noAccount') : t('auth.toggles.hasAccount')
)

const orText = computed(() =>
  userStore.loginTab === 'login' ? t('auth.orUseUsername') : t('auth.orUseUsernameSignUp')
)
</script>

<template>
  <div class="page">
    <div class="shell auth-wrap">
      <div class="auth-card">
        <!-- Left side: Navy background with feature highlights -->
        <div class="auth-card-left">
          <div>
            <div class="auth-eyebrow">
              {{ eyebrowText }}
            </div>
            <h1 class="auth-title">
              {{ titleText }}
            </h1>
            <p class="muted" style="font-size:14px; max-width:320px">
              {{ $t('auth.description') }}
            </p>
          </div>
          <div class="auth-highlights">
            <div class="auth-hl">
              <span class="auth-hl-num">01</span> {{ $t('auth.features.puzzles') }}
            </div>
            <div class="auth-hl">
              <span class="auth-hl-num">02</span> {{ $t('auth.features.history') }}
            </div>
            <div class="auth-hl">
              <span class="auth-hl-num">03</span> {{ $t('auth.features.multiplayer') }}
            </div>
          </div>
        </div>

        <!-- Right side: Auth form -->
        <form class="auth-card-right" @submit.prevent="handleSubmit">
          <!-- Tabs with sliding indicator -->
          <div class="auth-tabs">
            <button
              type="button"
              class="auth-tab"
              :data-active="userStore.loginTab === 'login'"
              @click="userStore.setLoginTab('login')"
            >
              {{ $t('auth.tabs.signIn') }}
            </button>
            <button
              type="button"
              class="auth-tab"
              :data-active="userStore.loginTab === 'register'"
              @click="userStore.setLoginTab('register')"
            >
              {{ $t('auth.tabs.signUp') }}
            </button>
            <div
              class="auth-tabs-slide"
              :style="{ left: userStore.loginTab === 'login' ? '4px' : 'calc(50% + 2px)' }"
            />
          </div>

          <!-- Telegram Login Widget -->
          <div id="telegram-login-container" style="display: flex; justify-content: center; margin-bottom: 20px"></div>

          <div class="auth-or">
            <span>{{ orText }}</span>
          </div>

          <!-- Form fields -->
          <div class="field">
            <label class="field-label">{{ usernameLabel }}</label>
            <input
              v-model="username"
              class="input"
              type="text"
              :placeholder="usernamePlaceholder"
              autofocus
            />
          </div>

          <div v-if="userStore.loginTab === 'register'" class="field">
            <label class="field-label">{{ $t('auth.fields.email') }}</label>
            <input v-model="email" class="input" type="email" :placeholder="$t('auth.placeholders.email')" />
          </div>

          <div class="field">
            <label class="field-label">{{ $t('auth.fields.password') }}</label>
            <input v-model="password" class="input" type="password" :placeholder="$t('auth.placeholders.password')" />
          </div>

          <!-- Privacy Policy Checkbox (only for registration) -->
          <div v-if="userStore.loginTab === 'register'" style="margin: 16px 0">
            <label style="display: flex; align-items: start; gap: 8px; cursor: pointer">
              <input
                type="checkbox"
                v-model="acceptedPrivacy"
                style="margin-top: 2px; cursor: pointer"
              />
              <span style="font-size: 13px; line-height: 1.4">
                {{ $t('auth.acceptPrivacy') }}
                <router-link to="/privacy" style="color: var(--cocoa); text-decoration: underline">
                  {{ $t('auth.privacyPolicy') }}
                </router-link>
              </span>
            </label>
          </div>

          <!-- Error message -->
          <div v-if="error" class="auth-err">{{ error }}</div>

          <!-- Submit button -->
          <button type="submit" class="btn btn--accent btn--block btn--lg" :disabled="isSubmitDisabled">
            {{ submitButtonText }}
          </button>

          <!-- Toggle text -->
          <p
            class="muted"
            style="font-size:12px; text-align:center; margin-top:16px; cursor:pointer"
            @click="userStore.setLoginTab(userStore.loginTab === 'login' ? 'register' : 'login')"
          >
            {{ toggleText }}
          </p>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* All styles are in pages.css and app.css */
</style>
