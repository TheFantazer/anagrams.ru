<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/userStore'

const router = useRouter()
const userStore = useUserStore()

const username = ref('')
const email = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

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
      password: password.value
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
  router.push('/')
}

function handleGoogleLogin() {
  const apiUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080'
  window.location.href = `${apiUrl}/api/v1/auth/google`
}
</script>

<template>
  <div class="page">
    <div class="shell auth-wrap">
      <div class="auth-card">
        <!-- Left side: Navy background with feature highlights -->
        <div class="auth-card-left">
          <div>
            <div class="auth-eyebrow">
              {{ userStore.loginTab === 'login' ? 'Welcome back' : 'Join the board' }}
            </div>
            <h1 class="auth-title">
              {{ userStore.loginTab === 'login' ? 'Sign in to keep your streak.' : 'Create an account in 30 seconds.' }}
            </h1>
            <p class="muted" style="font-size:14px; max-width:320px">
              Save your best words, track daily stats, and show up on the leaderboard. No ads, no emails.
            </p>
          </div>
          <div class="auth-highlights">
            <div class="auth-hl">
              <span class="auth-hl-num">01</span> daily puzzles & streaks
            </div>
            <div class="auth-hl">
              <span class="auth-hl-num">02</span> history across devices
            </div>
            <div class="auth-hl">
              <span class="auth-hl-num">03</span> async multiplayer
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
              Sign in
            </button>
            <button
              type="button"
              class="auth-tab"
              :data-active="userStore.loginTab === 'register'"
              @click="userStore.setLoginTab('register')"
            >
              Sign up
            </button>
            <div
              class="auth-tabs-slide"
              :style="{ left: userStore.loginTab === 'login' ? '4px' : 'calc(50% + 2px)' }"
            />
          </div>

          <!-- Google OAuth button -->
          <button type="button" class="btn btn--ghost btn--block" style="padding:14px" @click="handleGoogleLogin">
            <svg width="18" height="18" viewBox="0 0 24 24">
              <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 01-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4"/>
              <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
              <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
              <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
            </svg>
            Continue with Google
          </button>

          <div class="auth-or">
            <span>or use your {{ userStore.loginTab === 'login' ? 'username or email' : 'username' }}</span>
          </div>

          <!-- Form fields -->
          <div class="field">
            <label class="field-label">{{ userStore.loginTab === 'login' ? 'Username or Email' : 'Username' }}</label>
            <input
              v-model="username"
              class="input"
              type="text"
              :placeholder="userStore.loginTab === 'login' ? 'Enter username or email' : 'vera.m'"
              autofocus
            />
          </div>

          <div v-if="userStore.loginTab === 'register'" class="field">
            <label class="field-label">Email</label>
            <input v-model="email" class="input" type="email" placeholder="you@example.com" />
          </div>

          <div class="field">
            <label class="field-label">Password</label>
            <input v-model="password" class="input" type="password" placeholder="••••••••" />
          </div>

          <!-- Error message -->
          <div v-if="error" class="auth-err">{{ error }}</div>

          <!-- Submit button -->
          <button type="submit" class="btn btn--accent btn--block btn--lg" :disabled="loading">
            {{ loading ? 'Loading...' : (userStore.loginTab === 'login' ? 'Sign in' : 'Create account') }}
          </button>

          <!-- Toggle text -->
          <p
            class="muted"
            style="font-size:12px; text-align:center; margin-top:16px; cursor:pointer"
            @click="userStore.setLoginTab(userStore.loginTab === 'login' ? 'register' : 'login')"
          >
            {{ userStore.loginTab === 'login' ? "Don't have an account? Sign up" : 'Already have one? Sign in' }}
          </p>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* All styles are in pages.css and app.css */
</style>
