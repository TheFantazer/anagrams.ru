<script setup>
</script>

<template>
  <footer class="app-footer">
    <div class="shell">
      <router-link to="/privacy" class="footer-link">
        {{ $t('footer.privacyPolicy') }}
      </router-link>
      <span class="muted footer-copy">© {{ new Date().getFullYear() }} Anagrams.ru</span>
    </div>
  </footer>
</template>

<style scoped>
/* Footer styles are in app.css */
</style>
