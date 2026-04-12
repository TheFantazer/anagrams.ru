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

  // Сохраняем пользователя
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

  // Сохраняем пользователя
  userStore.setUser(data)
  router.push('/')
}
</script>

<template>
  <div class="auth-page">
    <div class="auth-header">
      <div class="logo-box">
        <span class="logo-text">A</span>
      </div>
      <h2 class="page-title">
        {{ userStore.loginTab === 'login' ? 'Welcome back' : 'Create account' }}
      </h2>
      <p class="subtitle">Find words. Beat friends.</p>
    </div>

    <div class="auth-tabs">
      <button
        :class="['auth-tab', { active: userStore.loginTab === 'login' }]"
        @click="userStore.setLoginTab('login')"
      >
        Sign in
      </button>
      <button
        :class="['auth-tab', { active: userStore.loginTab === 'register' }]"
        @click="userStore.setLoginTab('register')"
      >
        Sign up
      </button>
    </div>

    <button class="google-btn">
      <svg width="18" height="18" viewBox="0 0 24 24">
        <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 01-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z" fill="#4285F4"/>
        <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
        <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
        <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
      </svg>
      Continue with Google
    </button>

    <div class="divider">
      <div class="divider-line" />
      <span>or</span>
      <div class="divider-line" />
    </div>

    <div v-if="userStore.loginTab === 'register'">
      <label class="label">Username</label>
      <input v-model="username" class="input" placeholder="Choose a username" />
    </div>

    <template v-if="userStore.loginTab === 'login'">
      <label class="label">Username</label>
      <input v-model="username" class="input" type="text" placeholder="Enter your username" />
    </template>
    <template v-else>
      <label class="label">Email</label>
      <input v-model="email" class="input" type="email" placeholder="you@example.com" />
    </template>

    <label class="label">Password</label>
    <input v-model="password" class="input" type="password" placeholder="********" @keyup.enter="handleSubmit" />

    <div v-if="error" class="error-message">
      {{ error }}
    </div>

    <button class="btn-primary" @click="handleSubmit" :disabled="loading">
      <span v-if="loading">Loading...</span>
      <span v-else>{{ userStore.loginTab === 'login' ? 'Sign in' : 'Create account' }}</span>
    </button>

    <p class="toggle-text" @click="userStore.setLoginTab(userStore.loginTab === 'login' ? 'register' : 'login')">
      {{ userStore.loginTab === 'login'
        ? "Don't have an account? Sign up"
        : 'Already have an account? Sign in'
      }}
    </p>
  </div>
</template>

<style scoped>
.auth-page {
  max-width: 400px;
  margin: 60px auto;
  padding: 0 24px;
}

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo-box {
  width: 48px;
  height: 48px;
  margin: 0 auto 16px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-text {
  font-family: 'Space Mono', monospace;
  font-weight: 700;
  font-size: 26px;
  color: #fff;
}

.page-title {
  font-family: 'Space Mono', monospace;
  font-size: 22px;
  font-weight: 700;
  margin: 0 0 4px;
  color: #e8e6e1;
}

.subtitle {
  color: #555;
  font-size: 13px;
  margin: 0;
}

.auth-tabs {
  display: flex;
  gap: 0;
  margin-bottom: 32px;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.auth-tab {
  flex: 1;
  padding: 12px;
  text-align: center;
  cursor: pointer;
  font-size: 14px;
  font-weight: 600;
  background: rgba(255, 255, 255, 0.02);
  color: #666;
  transition: all 0.2s;
  border: none;
  font-family: 'Outfit', sans-serif;
}

.auth-tab.active {
  background: rgba(99, 230, 190, 0.1);
  color: var(--accent);
}

.google-btn {
  width: 100%;
  padding: 14px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #e8e6e1;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  font-family: 'Outfit', sans-serif;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  margin-bottom: 20px;
  transition: all 0.2s;
}

.google-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.divider {
  display: flex;
  align-items: center;
  gap: 16px;
  margin: 20px 0;
  color: #444;
  font-size: 12px;
}

.divider-line {
  flex: 1;
  height: 1px;
  background: rgba(255, 255, 255, 0.06);
}

.label {
  font-size: 12px;
  color: #888;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin: 0 0 8px;
  display: block;
  font-weight: 500;
}

.input {
  width: 100%;
  padding: 12px 16px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: #e8e6e1;
  font-size: 14px;
  outline: none;
  font-family: 'Outfit', sans-serif;
  margin-bottom: 20px;
  box-sizing: border-box;
}

.btn-primary {
  width: 100%;
  padding: 14px;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--accent), var(--accent-hover));
  border: none;
  color: var(--bg-dark);
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  font-family: 'Outfit', sans-serif;
  transition: all 0.2s;
  letter-spacing: 0.5px;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(99, 230, 190, 0.3);
}

.toggle-text {
  text-align: center;
  font-size: 12px;
  color: #444;
  margin-top: 16px;
  cursor: pointer;
  transition: color 0.2s;
}

.toggle-text:hover {
  color: #666;
}

.error-message {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #ef4444;
  padding: 12px 16px;
  border-radius: 10px;
  font-size: 13px;
  margin-bottom: 16px;
  text-align: center;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
}
</style>
