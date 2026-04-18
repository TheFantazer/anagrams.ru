<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/userStore'
import { useGameStore } from '../stores/gameStore'

const router = useRouter()
const userStore = useUserStore()
const gameStore = useGameStore()

// Ambient floating letters
const ambientLetters = [
  { ch: 'S', left: '10%', top: '15%', rot: -16, delay: 0, size: 'lg' },
  { ch: 'T', left: '22%', top: '26%', rot: -2, delay: 0.25, size: 'md' },
  { ch: 'R', left: '34%', top: '37%', rot: 12, delay: 0.5, size: 'md' },
  { ch: 'E', left: '46%', top: '48%', rot: -16, delay: 0.75, size: 'lg' },
  { ch: 'A', left: '58%', top: '59%', rot: -2, delay: 1, size: 'md' },
  { ch: 'M', left: '70%', top: '25%', rot: 12, delay: 1.25, size: 'md' },
  { ch: 'I', left: '82%', top: '36%', rot: -16, delay: 1.5, size: 'lg' },
]

function startFastGame() {
  gameStore.startGame(userStore.soloTime, userStore.soloLetters, userStore.soloLang)
  router.push('/game')
}

function openCustomSetup() {
  userStore.setShowSoloSettings(true)
}
</script>

<template>
  <div class="page">
    <div class="shell">
      <!-- Hero Section -->
      <section class="home-hero">
        <div class="home-hero-letters" aria-hidden="true">
          <span
            v-for="(letter, i) in ambientLetters"
            :key="i"
            :class="['hero-letter', `hero-letter--${letter.size}`]"
            :style="{
              left: letter.left,
              top: letter.top,
              '--rot': `${letter.rot}deg`,
              animationDelay: `${letter.delay}s`
            }"
          >
            {{ letter.ch }}
          </span>
        </div>

        <div class="home-hero-text">
          <div class="home-eyebrow">
            <span class="dot-blink" /> A word puzzle for two
          </div>
          <h1 class="home-title">
            Find every word.<br/>
            <span class="home-title-accent">Beat your friends.</span>
          </h1>
          <p class="home-sub">
            Seven letters, one minute, every anagram hiding inside. Share the link — they play
            the same set. The longer word wins.
          </p>
          <div class="home-cta-row">
            <button class="btn btn--accent btn--lg" @click="startFastGame">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M13 2L4 14h7l-1 8 9-12h-7l1-8z"/>
              </svg>
              Fast game
            </button>
            <button class="btn btn--ghost btn--lg" @click="openCustomSetup">
              <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="3"/>
                <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"/>
              </svg>
              Custom setup
            </button>
          </div>
          <div class="home-kbd-hint">
            <span class="kbd">F</span> fast &nbsp;·&nbsp; <span class="kbd">C</span> custom &nbsp;·&nbsp; <span class="kbd">?</span> how to play
          </div>
        </div>
      </section>

      <!-- Modes Section -->
      <section class="home-modes">
        <div class="home-modes-head">
          <div>
            <h2 class="home-modes-title">Pick a mode</h2>
            <p class="muted" style="margin:0;font-size:13px">Solo for practice, multiplayer for bragging rights.</p>
          </div>
        </div>
        <div class="home-modes-grid">
          <!-- Solo -->
          <div class="mode-card" @click="openCustomSetup">
            <div class="mode-card-top">
              <span class="mode-num">01</span>
              <span class="mode-tag">Warm up</span>
            </div>
            <div class="mode-card-body">
              <h3 class="mode-name">Solo</h3>
              <p class="mode-desc">Custom letters, language, and timer. Perfect for warming up.</p>
            </div>
            <div class="mode-card-foot">
              <span>Start</span>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M5 12h14M12 5l7 7-7 7"/>
              </svg>
            </div>
          </div>

          <!-- Fast -->
          <div class="mode-card mode-card--accent" @click="startFastGame">
            <div class="mode-card-top">
              <span class="mode-num">02</span>
              <span class="mode-tag">Recommended</span>
            </div>
            <div class="mode-card-body">
              <h3 class="mode-name">Fast</h3>
              <p class="mode-desc">7 letters · 60 seconds · pure reflex. Jump straight in.</p>
            </div>
            <div class="mode-card-foot">
              <span>Start</span>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M5 12h14M12 5l7 7-7 7"/>
              </svg>
            </div>
          </div>

          <!-- Multiplayer -->
          <div class="mode-card" @click="router.push('/multiplayer')">
            <div class="mode-card-top">
              <span class="mode-num">03</span>
              <span class="mode-tag">New</span>
            </div>
            <div class="mode-card-body">
              <h3 class="mode-name">Multiplayer</h3>
              <p class="mode-desc">Challenge a friend with a link. You both play the same set — highest score wins.</p>
            </div>
            <div class="mode-card-foot">
              <span>Start</span>
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
                <path d="M5 12h14M12 5l7 7-7 7"/>
              </svg>
            </div>
          </div>
        </div>
      </section>

      <!-- Stats Strip -->
      <section class="home-strip">
        <div class="strip-item">
          <div class="strip-num">3+</div>
          <div class="strip-lbl">letter min.</div>
        </div>
        <div class="strip-sep" />
        <div class="strip-item">
          <div class="strip-num mono">100 → 2.8k</div>
          <div class="strip-lbl">points scale</div>
        </div>
        <div class="strip-sep" />
        <div class="strip-item">
          <div class="strip-num"><span class="flag-dot" />EN · RU</div>
          <div class="strip-lbl">two dictionaries</div>
        </div>
        <div class="strip-sep" />
        <div class="strip-item">
          <div class="strip-num">0</div>
          <div class="strip-lbl">apps to install</div>
        </div>
      </section>

      <!-- Sign In Invite (if not authenticated) -->
      <section v-if="!userStore.isAuthenticated" class="home-invite">
        <div>
          <h3 style="color:var(--fg-primary);font-family:var(--font-display);font-size:22px;letter-spacing:-0.5px;margin:0 0 6px;text-transform:none">
            Keep your streak.
          </h3>
          <p class="muted" style="margin:0;font-size:13px;max-width:440px">
            Sign in to save your best words, track stats, and appear on the leaderboard.
          </p>
        </div>
        <div class="row gap-2">
          <button class="btn btn--ghost" @click="router.push('/auth')">Sign in</button>
          <button class="btn btn--primary" @click="router.push('/auth')">Create account</button>
        </div>
      </section>
    </div>
  </div>
</template>

<style scoped>
/* Home Hero */
.home-hero {
  position: relative;
  padding: 80px 0 48px;
  display: grid;
  grid-template-columns: 1fr;
  align-items: center;
  min-height: 560px;
  overflow: hidden;
}

.home-hero-letters {
  position: absolute;
  inset: 0;
  pointer-events: none;
  z-index: 0;
}

.hero-letter {
  position: absolute;
  font-family: var(--font-display);
  font-weight: 800;
  color: var(--navy);
  opacity: 0.04;
  transform: rotate(var(--rot));
  animation: float-up 12s ease-in-out infinite alternate;
  user-select: none;
}

.hero-letter--md { font-size: 120px; }
.hero-letter--lg { font-size: 200px; }

@keyframes float-up {
  from { transform: translateY(0) rotate(var(--rot)); }
  to   { transform: translateY(-20px) rotate(calc(var(--rot) + 4deg)); }
}

.home-hero-text {
  position: relative;
  z-index: 1;
  max-width: 640px;
}

.home-eyebrow {
  display: inline-flex; align-items: center; gap: 8px;
  padding: 6px 14px 6px 10px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 999px;
  font-size: 12px;
  color: var(--fg-secondary);
  font-weight: 500;
  margin-bottom: 24px;
}

.dot-blink {
  width: 8px; height: 8px; border-radius: 999px;
  background: var(--accent);
  animation: blink 1.6s ease-in-out infinite;
}

@keyframes blink {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.4; transform: scale(0.8); }
}

.home-title {
  font-family: var(--font-display);
  font-size: 72px;
  font-weight: 700;
  line-height: 0.95;
  letter-spacing: -2.5px;
  color: var(--fg-primary);
  margin: 0 0 20px;
}

.home-title-accent {
  color: var(--accent);
  font-style: italic;
  font-weight: 500;
}

.home-sub {
  font-size: 17px;
  color: var(--fg-secondary);
  max-width: 540px;
  line-height: 1.55;
  margin: 0 0 32px;
}

.home-cta-row {
  display: flex; gap: 12px; margin-bottom: 20px; flex-wrap: wrap;
}

.home-kbd-hint {
  font-size: 12px;
  color: var(--fg-muted);
}

/* Modes */
.home-modes {
  margin-top: 32px;
  padding: 40px 0 24px;
  border-top: 1px solid var(--border-hairline);
}

.home-modes-head {
  margin-bottom: 24px;
}

.home-modes-title {
  font-family: var(--font-display);
  font-size: 32px;
  font-weight: 700;
  letter-spacing: -1px;
  color: var(--fg-primary);
  margin: 0 0 4px;
}

.home-modes-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.mode-card {
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 20px;
  padding: 24px;
  cursor: pointer;
  transition: all var(--dur-base) var(--ease-out);
  position: relative;
  display: flex;
  flex-direction: column;
  min-height: 220px;
}

.mode-card:hover {
  transform: translateY(-4px);
  border-color: var(--border-default);
  box-shadow: var(--shadow-md);
}

.mode-card--accent {
  background: var(--navy);
  color: var(--milk);
}

.mode-card--accent:hover { background: var(--navy-2); }

.mode-card-top {
  display: flex; align-items: center; justify-content: space-between;
  margin-bottom: 20px;
}

.mode-num {
  font-family: var(--font-mono);
  font-size: 12px;
  font-weight: 500;
  color: var(--fg-muted);
  letter-spacing: 1px;
}

.mode-card--accent .mode-num { color: rgba(251, 246, 236, 0.6); }

.mode-tag {
  font-size: 10px;
  text-transform: uppercase;
  letter-spacing: 1.2px;
  padding: 4px 10px;
  border-radius: 999px;
  background: var(--bg-card);
  color: var(--fg-secondary);
  font-weight: 600;
}

.mode-card--accent .mode-tag {
  background: var(--accent);
  color: var(--milk);
}

.mode-card-body {
  flex: 1;
}

.mode-name {
  font-family: var(--font-display);
  font-size: 32px;
  font-weight: 700;
  letter-spacing: -0.5px;
  margin: 0 0 10px;
  color: inherit;
}

.mode-desc {
  font-size: 13px;
  line-height: 1.55;
  color: var(--fg-muted);
  margin: 0;
}

.mode-card--accent .mode-desc { color: rgba(251, 246, 236, 0.72); }

.mode-card-foot {
  display: flex; align-items: center; justify-content: space-between;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid var(--border-hairline);
  font-size: 13px;
  font-weight: 600;
  color: var(--fg-secondary);
}

.mode-card--accent .mode-card-foot {
  color: var(--milk);
  border-top-color: rgba(251, 246, 236, 0.15);
}

.mode-card:hover .mode-card-foot svg { transform: translateX(4px); }
.mode-card-foot svg { transition: transform var(--dur-base) var(--ease-out); }

/* Strip */
.home-strip {
  margin-top: 48px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24px 32px;
  background: var(--bg-surface);
  border: 1px solid var(--border-subtle);
  border-radius: 20px;
  flex-wrap: wrap;
  gap: 20px;
}

.strip-item { min-width: 120px; }

.strip-num {
  font-family: var(--font-display);
  font-size: 22px;
  font-weight: 700;
  color: var(--fg-primary);
  letter-spacing: -0.5px;
  display: inline-flex; align-items: center; gap: 8px;
}

.strip-num.mono { font-family: var(--font-mono); }

.strip-lbl {
  font-size: 11px;
  text-transform: uppercase;
  letter-spacing: var(--ls-upper);
  color: var(--fg-muted);
  margin-top: 4px;
  font-weight: 600;
}

.strip-sep { width: 1px; align-self: stretch; background: var(--border-hairline); }

.flag-dot { width: 10px; height: 10px; border-radius: 999px; background: var(--accent); display: inline-block; }

/* Invite */
.home-invite {
  margin-top: 40px;
  padding: 28px 32px;
  border: 1px dashed var(--border-default);
  border-radius: 20px;
  background: linear-gradient(135deg, var(--milk), var(--bg-card));
  display: flex; justify-content: space-between; align-items: center;
  flex-wrap: wrap; gap: 16px;
}

/* Mobile Responsive */
@media (max-width: 860px) {
  .home-modes-grid { grid-template-columns: 1fr; }
  .home-title { font-size: 52px; letter-spacing: -1.5px; }
  .home-hero { min-height: auto; padding: 48px 0 32px; }
  .hero-letter--md { font-size: 80px; }
  .hero-letter--lg { font-size: 140px; }
}

@media (max-width: 720px) {
  .strip-sep { display: none; }
  .home-strip { gap: 16px; }
}
</style>
